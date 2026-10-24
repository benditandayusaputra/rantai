package http

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

const sqlPaketPantauUji = `INSERT INTO provenance (id, sumber, diambil_pada, label) VALUES ('prv_uji_pantau', 'lpse', now(), 'resmi');
INSERT INTO instansi (nama, jenis, kode_kota) VALUES ('SUKU DINAS UJI PANTAU', 'sudin', '31.72');
INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap,
    lokasi_tingkat, kode_wilayah, provenance_id)
SELECT 'nontender-pantau', 'pantau', 'Rehab Saluran Kelapa Gading Barat', (SELECT id FROM instansi WHERE nama = 'SUKU DINAS UJI PANTAU'),
    false, 'Pekerjaan Konstruksi', 'Pengadaan Langsung', 2026, 'Paket Sudah Selesai', 'kelurahan', w.kode, 'prv_uji_pantau'
FROM wilayah w WHERE w.nama = 'Kelapa Gading Barat' AND w.tingkat = 'kelurahan'`

func appUjiPantau(t *testing.T) (*fiber.App, pgx.Tx) {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	akses := aksesLaporanTiruan{pengguna: map[string]int64{}}
	for _, token := range []string{"token-pemantau", "token-lain"} {
		var id int64
		if err := tx.QueryRow(t.Context(), `INSERT INTO pengguna (email) VALUES ($1 || '@uji.test') RETURNING id`, token).Scan(&id); err != nil {
			t.Fatal(err)
		}
		akses.pengguna[token] = id
	}
	return NewApp(config.Server{}, Layanan{Pantau: usecase.Pantau{Repo: repository.PantauRepository{DB: tx}}, Akses: akses}), tx
}

func mintaPantau(t *testing.T, app *fiber.App, metode, path, body, token string) hasilPanggilan {
	t.Helper()
	permintaan := httptest.NewRequestWithContext(t.Context(), metode, path, bytes.NewBufferString(body))
	permintaan.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	permintaan.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	return kirim(t, app, permintaan)
}

func buatKejadianUji(t *testing.T, tx pgx.Tx) int {
	t.Helper()
	baru, err := repository.PantauRepository{DB: tx}.BuatKejadian(t.Context(), time.Now().Add(-domain.JendelaKejadian))
	if err != nil {
		t.Fatal(err)
	}
	return baru
}

func TestPantau_LanggananKelurahan_FeedPaketBaruDanTandaiDibaca(t *testing.T) {
	app, tx := appUjiPantau(t)
	var kelurahan string
	if err := tx.QueryRow(t.Context(), `SELECT kode FROM wilayah WHERE nama = 'Kelapa Gading Barat' AND tingkat = 'kelurahan'`).Scan(&kelurahan); err != nil {
		t.Fatal(err)
	}
	langganan := mintaPantau(t, app, fiber.MethodPost, "/api/v1/langganan", fmt.Sprintf(`{"kode_wilayah": %q}`, kelurahan), "token-pemantau")
	if _, err := tx.Exec(t.Context(), sqlPaketPantauUji); err != nil {
		t.Fatal(err)
	}

	pertama, kedua := buatKejadianUji(t, tx), buatKejadianUji(t, tx)
	feed, _ := mintaPantau(t, app, fiber.MethodGet, "/api/v1/notifikasi", "", "token-pemantau").respons.Data.(map[string]any)
	lain, _ := mintaPantau(t, app, fiber.MethodGet, "/api/v1/notifikasi", "", "token-lain").respons.Data.(map[string]any)

	data, _ := langganan.respons.Data.(map[string]any)
	item, _ := feed["item"].([]any)
	if langganan.kode != fiber.StatusCreated || data["kode_wilayah"] != "31.72.06" || pertama < 1 || kedua != 0 {
		t.Fatalf("langganan %v, kejadian baru %d lalu %d", langganan.respons, pertama, kedua)
	}
	if len(item) != 1 || feed["belum_dibaca"] != float64(1) {
		t.Fatalf("feed = %v", feed)
	}
	if lainItem, _ := lain["item"].([]any); len(lainItem) != 0 {
		t.Errorf("pengguna tanpa langganan menerima %v", lainItem)
	}
	notifikasi, _ := item[0].(map[string]any)
	mintaPantau(t, app, fiber.MethodPatch, fmt.Sprintf("/api/v1/notifikasi/%v/dibaca", notifikasi["id"]), "", "token-pemantau")
	setelah, _ := mintaPantau(t, app, fiber.MethodGet, "/api/v1/notifikasi", "", "token-pemantau").respons.Data.(map[string]any)
	if notifikasi["jenis"] != "paket_baru" || notifikasi["tautan"] != "/paket/nontender-pantau" || setelah["belum_dibaca"] != float64(0) {
		t.Errorf("notifikasi %v, setelah dibaca %v", notifikasi, setelah)
	}
}

func TestPantau_HapusLanggananOrangLain_404(t *testing.T) {
	app, _ := appUjiPantau(t)
	langganan, _ := mintaPantau(t, app, fiber.MethodPost, "/api/v1/langganan", `{"kode_wilayah": "31.72.06"}`, "token-pemantau").respons.Data.(map[string]any)

	orangLain := mintaPantau(t, app, fiber.MethodDelete, fmt.Sprintf("/api/v1/langganan/%v", langganan["id"]), "", "token-lain")
	pemilik := mintaPantau(t, app, fiber.MethodDelete, fmt.Sprintf("/api/v1/langganan/%v", langganan["id"]), "", "token-pemantau")

	if orangLain.kode != fiber.StatusNotFound || pemilik.kode != fiber.StatusOK {
		t.Errorf("orang lain %d, pemilik %d", orangLain.kode, pemilik.kode)
	}
}

func TestPantau_KodeBukanWilayahDKI_404(t *testing.T) {
	app, _ := appUjiPantau(t)

	if hasil := mintaPantau(t, app, fiber.MethodPost, "/api/v1/langganan", `{"kode_wilayah": "32.73.01"}`, "token-pemantau"); hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode = %d, harap 404", hasil.kode)
	}
}

type pengirimPushUji struct{ jumlah *int }

func (p pengirimPushUji) Kirim(context.Context, domain.LanggananPush, domain.PesanPush) error {
	*p.jumlah++
	return nil
}

func TestWebPush_AktifkanLaluKejadianBaru_TerkirimSekali(t *testing.T) {
	_, tx := appUjiPantau(t)
	repo := repository.PantauRepository{DB: tx}
	app := NewApp(config.Server{}, Layanan{
		Pantau: usecase.Pantau{Repo: repo}, WebPush: usecase.LayananPush{Repo: repo, KunciPublik: "kunci-uji"},
		Akses: aksesLaporanTiruan{pengguna: map[string]int64{"token-pemantau": idPengguna(t, tx, "token-pemantau")}},
	})
	mintaPantau(t, app, fiber.MethodPost, "/api/v1/langganan", `{"kode_wilayah": "31.72.06"}`, "token-pemantau")
	aktif := mintaPantau(t, app, fiber.MethodPost, "/api/v1/webpush/langganan",
		`{"endpoint": "https://push.contoh.test/abc", "keys": {"p256dh": "kunci", "auth": "rahasia"}}`, "token-pemantau")
	if _, err := tx.Exec(t.Context(), sqlPaketPantauUji); err != nil {
		t.Fatal(err)
	}
	buatKejadianUji(t, tx)
	jumlah := 0
	dorong := usecase.DorongPush{Repo: repo, Pengirim: pengirimPushUji{jumlah: &jumlah}}

	pertama, errPertama := dorong.Dorong(t.Context())
	kedua, errKedua := dorong.Dorong(t.Context())

	data, _ := aktif.respons.Data.(map[string]any)
	if aktif.kode != fiber.StatusOK || data["jumlah_kecamatan"] != float64(1) {
		t.Fatalf("aktifkan push: %d %v", aktif.kode, aktif.respons)
	}
	if errPertama != nil || errKedua != nil || pertama < 1 || kedua != 0 || jumlah != pertama {
		t.Errorf("pertama %d (%v), kedua %d (%v), terkirim %d", pertama, errPertama, kedua, errKedua, jumlah)
	}
}

func idPengguna(t *testing.T, tx pgx.Tx, token string) int64 {
	t.Helper()
	var id int64
	if err := tx.QueryRow(t.Context(), `SELECT id FROM pengguna WHERE email = $1 || '@uji.test'`, token).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}
