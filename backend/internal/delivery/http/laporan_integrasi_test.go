package http

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/foto"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

const tokenLaporanUji = "token-uji"

type aksesLaporanTiruan struct{ pengguna map[string]int64 }

func (a aksesLaporanTiruan) Periksa(token string) (domain.KlaimAkses, error) {
	id, ada := a.pengguna[token]
	if !ada {
		return domain.KlaimAkses{}, domain.ErrTokenAksesTidakBerlaku
	}
	return domain.KlaimAkses{PenggunaID: id, Peran: domain.PeranWarga}, nil
}

type urlFotoTiruan struct{}

func (urlFotoTiruan) URLBertanda(_ context.Context, kunci string) (string, error) {
	return "https://objek.test/" + kunci, nil
}

type penyimpanObjekTiruan struct{ isi map[string][]byte }

func (p penyimpanObjekTiruan) SimpanJPEG(_ context.Context, kunci string, isi []byte) error {
	p.isi[kunci] = isi
	return nil
}

func appUjiLaporan(t *testing.T) (*fiber.App, pgx.Tx, penyimpanObjekTiruan) {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	akses := aksesLaporanTiruan{pengguna: map[string]int64{}}
	for _, token := range []string{tokenLaporanUji, "token-saksi-1", "token-saksi-2"} {
		var id int64
		if err := tx.QueryRow(t.Context(), `INSERT INTO pengguna (email) VALUES ($1 || '@uji.test') RETURNING id`, token).Scan(&id); err != nil {
			t.Fatal(err)
		}
		akses.pengguna[token] = id
	}
	penyimpan := penyimpanObjekTiruan{isi: map[string][]byte{}}
	repo := repository.LaporanRepository{DB: tx}
	return NewApp(config.Server{}, Layanan{
		Laporan: usecase.TerimaLaporan{
			Pemroses: foto.Pemroses{}, Penyimpan: penyimpan, Repo: repo,
			Sekarang: time.Now, BuatKunci: func() string { return fmt.Sprint(len(penyimpan.isi)) },
		},
		BacaLaporan:       usecase.BacaLaporan{Pembaca: repo, URLFoto: urlFotoTiruan{}},
		KonfirmasiLaporan: usecase.KonfirmasiLaporan{Repo: repo, Sekarang: time.Now},
		Akses:             akses,
	}), tx, penyimpan
}

func titikKelurahan(t *testing.T, tx pgx.Tx, nama string) domain.Titik {
	t.Helper()
	var titik domain.Titik
	err := tx.QueryRow(t.Context(), `SELECT ST_Y(titik_tengah), ST_X(titik_tengah) FROM wilayah WHERE nama = $1 AND tingkat = 'kelurahan'`, nama).
		Scan(&titik.Lat, &titik.Lng)
	if err != nil {
		t.Fatal(err)
	}
	return titik
}

func kirimLaporan(t *testing.T, app *fiber.App, titik domain.Titik, token string) hasilPanggilan {
	t.Helper()
	var body bytes.Buffer
	penulis := multipart.NewWriter(&body)
	for kunci, nilai := range map[string]string{
		"lat": fmt.Sprint(titik.Lat), "lng": fmt.Sprint(titik.Lng), "akurasi": "15", "kondisi": "terhenti", "catatan": "Pekerjaan berhenti.",
	} {
		if err := penulis.WriteField(kunci, nilai); err != nil {
			t.Fatal(err)
		}
	}
	berkas, err := penulis.CreateFormFile("foto", "bukti.jpg")
	if err != nil {
		t.Fatal(err)
	}
	if err := jpeg.Encode(berkas, image.NewRGBA(image.Rect(0, 0, 40, 30)), nil); err != nil {
		t.Fatal(err)
	}
	if err := penulis.Close(); err != nil {
		t.Fatal(err)
	}
	permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodPost, "/api/v1/laporan", &body)
	permintaan.Header.Set(fiber.HeaderContentType, penulis.FormDataContentType())
	permintaan.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	return kirim(t, app, permintaan)
}

func TestTerimaLaporan_Multipart_TersimpanMenungguSaksiDenganLokasiPublikTerbulatkan(t *testing.T) {
	app, tx, penyimpan := appUjiLaporan(t)
	titik := titikKelurahan(t, tx, "Kelapa Gading Barat")

	hasil := kirimLaporan(t, app, titik, tokenLaporanUji)

	data, _ := hasil.respons.Data.(map[string]any)
	if hasil.kode != fiber.StatusCreated || data["status"] != "MENUNGGU_SAKSI" || len(penyimpan.isi) != 2 {
		t.Fatalf("kode %d, data %v, objek %d", hasil.kode, hasil.respons, len(penyimpan.isi))
	}
	var publik, asli domain.Titik
	var wilayah string
	var riwayat, jumlahFoto int
	err := tx.QueryRow(t.Context(), `SELECT ST_Y(lokasi_publik), ST_X(lokasi_publik), ST_Y(lokasi_asli), ST_X(lokasi_asli), kode_wilayah,
    (SELECT count(*) FROM laporan_riwayat r WHERE r.laporan_id = l.id), (SELECT count(*) FROM laporan_foto f WHERE f.laporan_id = l.id)
FROM laporan l WHERE id = $1`, data["id"]).Scan(&publik.Lat, &publik.Lng, &asli.Lat, &asli.Lng, &wilayah, &riwayat, &jumlahFoto)
	if err != nil {
		t.Fatal(err)
	}
	harapPublik := domain.BulatkanLokasiPublik(titik)
	if publik == asli || !hampirSama(publik, harapPublik) || riwayat != 2 || jumlahFoto != 1 || len(wilayah) != len("31.72.06.1003") {
		t.Errorf("publik %+v (harap %+v), asli %+v, wilayah %s, riwayat %d, foto %d", publik, harapPublik, asli, wilayah, riwayat, jumlahFoto)
	}
}

func hampirSama(a, b domain.Titik) bool {
	const toleransiDerajat = 1e-9
	return math.Abs(a.Lat-b.Lat) < toleransiDerajat && math.Abs(a.Lng-b.Lng) < toleransiDerajat
}

func TestTerimaLaporan_TanpaToken_401(t *testing.T) {
	app, tx, _ := appUjiLaporan(t)

	if hasil := kirimLaporan(t, app, titikKelurahan(t, tx, "Kelapa Gading Barat"), "salah"); hasil.kode != fiber.StatusUnauthorized {
		t.Errorf("kode = %d, harap 401", hasil.kode)
	}
}

func TestTerimaLaporan_TitikDiLuarDKI_400(t *testing.T) {
	app, _, penyimpan := appUjiLaporan(t)

	hasil := kirimLaporan(t, app, domain.Titik{Lat: -6.9175, Lng: 107.6191}, tokenLaporanUji)

	if hasil.kode != fiber.StatusBadRequest || len(penyimpan.isi) != 0 {
		t.Errorf("kode %d, objek %d", hasil.kode, len(penyimpan.isi))
	}
}

func konfirmasikan(t *testing.T, app *fiber.App, id any, titik domain.Titik, token string) hasilPanggilan {
	t.Helper()
	body := fmt.Sprintf(`{"lat": %v, "lng": %v, "akurasi": 12, "kondisi": "terhenti"}`, titik.Lat, titik.Lng)
	permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodPost, fmt.Sprintf("/api/v1/laporan/%v/konfirmasi", id), bytes.NewBufferString(body))
	permintaan.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	permintaan.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	return kirim(t, app, permintaan)
}

func TestKonfirmasiLaporan_DuaSaksi_TerverifikasiDanTampilDiLinimasa(t *testing.T) {
	app, tx, _ := appUjiLaporan(t)
	titik := titikKelurahan(t, tx, "Kelapa Gading Barat")
	data, _ := kirimLaporan(t, app, titik, tokenLaporanUji).respons.Data.(map[string]any)

	pertama, _ := konfirmasikan(t, app, data["id"], titik, "token-saksi-1").respons.Data.(map[string]any)
	kedua, _ := konfirmasikan(t, app, data["id"], titik, "token-saksi-2").respons.Data.(map[string]any)
	detail := panggil(t, app, fmt.Sprintf("/api/v1/laporan/%v", data["id"]))

	if pertama["status"] != "MENUNGGU_SAKSI" || kedua["status"] != "TERVERIFIKASI" {
		t.Fatalf("pertama %v, kedua %v", pertama, kedua)
	}
	isi, _ := detail.respons.Data.(map[string]any)
	riwayat, _ := isi["riwayat"].([]any)
	foto, _ := isi["foto"].([]any)
	saksi, _ := isi["saksi"].(map[string]any)
	if detail.kode != fiber.StatusOK || isi["status"] != "TERVERIFIKASI" || len(riwayat) != 3 || len(foto) != 1 || saksi["jumlah"] != float64(3) {
		t.Errorf("detail = %v", isi)
	}
	for _, kunciPrivat := range []string{"lokasi_asli", "pengguna_id", "email"} {
		if _, ada := isi[kunciPrivat]; ada {
			t.Errorf("detail publik membocorkan %s", kunciPrivat)
		}
	}
}

func TestKonfirmasiLaporan_OlehPelaporAtauGanda_Ditolak(t *testing.T) {
	app, tx, _ := appUjiLaporan(t)
	titik := titikKelurahan(t, tx, "Kelapa Gading Barat")
	data, _ := kirimLaporan(t, app, titik, tokenLaporanUji).respons.Data.(map[string]any)
	konfirmasikan(t, app, data["id"], titik, "token-saksi-1")

	pelapor := konfirmasikan(t, app, data["id"], titik, tokenLaporanUji)
	ganda := konfirmasikan(t, app, data["id"], titik, "token-saksi-1")

	if pelapor.kode != fiber.StatusBadRequest || ganda.kode != fiber.StatusBadRequest {
		t.Errorf("pelapor %d, ganda %d", pelapor.kode, ganda.kode)
	}
}

func TestAmbilLaporan_TidakAda_404(t *testing.T) {
	app, _, _ := appUjiLaporan(t)

	if hasil := panggil(t, app, "/api/v1/laporan/999999999"); hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode = %d, harap 404", hasil.kode)
	}
}
