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

type aksesLaporanTiruan struct{ penggunaID int64 }

func (a aksesLaporanTiruan) Periksa(token string) (domain.KlaimAkses, error) {
	if token != tokenLaporanUji {
		return domain.KlaimAkses{}, domain.ErrTokenAksesTidakBerlaku
	}
	return domain.KlaimAkses{PenggunaID: a.penggunaID, Peran: domain.PeranWarga}, nil
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
	var penggunaID int64
	if err := tx.QueryRow(t.Context(), `INSERT INTO pengguna (email) VALUES ('pelapor@uji.test') RETURNING id`).Scan(&penggunaID); err != nil {
		t.Fatal(err)
	}
	penyimpan := penyimpanObjekTiruan{isi: map[string][]byte{}}
	terima := usecase.TerimaLaporan{
		Pemroses: foto.Pemroses{}, Penyimpan: penyimpan, Repo: repository.LaporanRepository{DB: tx},
		Sekarang: time.Now, BuatKunci: func() string { return fmt.Sprint(len(penyimpan.isi)) },
	}
	return NewApp(config.Server{}, Layanan{Laporan: terima, Akses: aksesLaporanTiruan{penggunaID: penggunaID}}), tx, penyimpan
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
