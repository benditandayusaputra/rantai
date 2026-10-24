package http

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

func patchModerasi(t *testing.T, app *fiber.App, path, body, token string) hasilPanggilan {
	t.Helper()
	permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodPatch, path, bytes.NewBufferString(body))
	permintaan.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	permintaan.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	return kirim(t, app, permintaan)
}

func ambilModerasi(t *testing.T, app *fiber.App, token string) hasilPanggilan {
	t.Helper()
	permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodGet, "/api/v1/moderasi/antrean", nil)
	permintaan.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	return kirim(t, app, permintaan)
}

func jumlahLogModerasi(t *testing.T, tx pgx.Tx, target string) int {
	t.Helper()
	var jumlah int
	if err := tx.QueryRow(t.Context(), `SELECT count(*) FROM moderasi_log WHERE target_jenis = $1`, target).Scan(&jumlah); err != nil {
		t.Fatal(err)
	}
	return jumlah
}

func TestModerasi_Warga_Ditolak(t *testing.T) {
	app, _, _ := appUjiLaporan(t)

	if hasil := ambilModerasi(t, app, tokenLaporanUji); hasil.kode != fiber.StatusUnauthorized && hasil.kode != fiber.StatusForbidden {
		t.Errorf("kode = %d, harap ditolak", hasil.kode)
	}
}

func TestModerasi_TolakLaporan_DitolakTercatatDiRiwayatDanLog(t *testing.T) {
	app, tx, _ := appUjiLaporan(t)
	data, _ := kirimLaporan(t, app, titikKelurahan(t, tx, "Kelapa Gading Barat"), tokenLaporanUji).respons.Data.(map[string]any)
	antrean, _ := ambilModerasi(t, app, tokenModeratorUji).respons.Data.(map[string]any)
	path := fmt.Sprintf("/api/v1/moderasi/laporan/%v", data["id"])

	hasil := patchModerasi(t, app, path, `{"aksi": "tolak", "alasan": "duplikat", "pesan": "sama dengan laporan lain"}`, tokenModeratorUji)
	ulang := patchModerasi(t, app, path, `{"aksi": "verifikasi"}`, tokenModeratorUji)

	daftar, _ := antrean["laporan"].([]any)
	if !slices.ContainsFunc(daftar, func(l any) bool { m, _ := l.(map[string]any); return m["id"] == data["id"] }) {
		t.Errorf("antrean laporan tanpa laporan baru: %v", daftar)
	}
	if hasil.kode != fiber.StatusOK || ulang.kode != fiber.StatusBadRequest {
		t.Fatalf("tolak %d, verifikasi setelah ditolak %d", hasil.kode, ulang.kode)
	}
	detail, _ := panggil(t, app, fmt.Sprintf("/api/v1/laporan/%v", data["id"])).respons.Data.(map[string]any)
	riwayat, _ := detail["riwayat"].([]any)
	terakhir, _ := riwayat[len(riwayat)-1].(map[string]any)
	if detail["status"] != "DITOLAK" || terakhir["alasan"] != "duplikat: sama dengan laporan lain" || jumlahLogModerasi(t, tx, domain.TargetLaporan) != 1 {
		t.Errorf("detail %v", detail)
	}
}

func TestModerasi_SembunyikanFoto_HilangDariHalamanPublik(t *testing.T) {
	app, tx, _ := appUjiLaporan(t)
	data, _ := kirimLaporan(t, app, titikKelurahan(t, tx, "Kelapa Gading Barat"), tokenLaporanUji).respons.Data.(map[string]any)
	var fotoID int64
	if err := tx.QueryRow(t.Context(), `SELECT id FROM laporan_foto WHERE laporan_id = $1`, data["id"]).Scan(&fotoID); err != nil {
		t.Fatal(err)
	}

	hasil := patchModerasi(t, app, fmt.Sprintf("/api/v1/moderasi/laporan/%v", data["id"]),
		fmt.Sprintf(`{"aksi": "sembunyikan_foto", "foto_id": %d, "alasan": "wajah warga terlihat"}`, fotoID), tokenModeratorUji)

	detail, _ := panggil(t, app, fmt.Sprintf("/api/v1/laporan/%v", data["id"])).respons.Data.(map[string]any)
	if foto, _ := detail["foto"].([]any); hasil.kode != fiber.StatusOK || len(foto) != 0 || jumlahLogModerasi(t, tx, domain.TargetLaporanFoto) != 1 {
		t.Errorf("kode %d, foto publik %v", hasil.kode, detail["foto"])
	}
}

func TestModerasi_SembunyikanUmpan_TersembunyiDanTercatat(t *testing.T) {
	app, tx, _ := appUjiLaporan(t)
	ujiwakil.MuatDPRD(t, tx)
	simpanUmpanUji(t, tx, umpanBertagUji(1, domain.SumberUmpanMedia, anggotaUmpanUji, "sekolah"))
	var umpanID int64
	if err := tx.QueryRow(t.Context(), `SELECT id FROM umpan_kegiatan WHERE url = 'https://contoh.test/umpan-1'`).Scan(&umpanID); err != nil {
		t.Fatal(err)
	}

	tanpaAlasan := patchModerasi(t, app, fmt.Sprintf("/api/v1/moderasi/umpan/%d", umpanID), `{}`, tokenModeratorUji)
	hasil := patchModerasi(t, app, fmt.Sprintf("/api/v1/moderasi/umpan/%d", umpanID), `{"alasan": "salah atribusi"}`, tokenModeratorUji)

	var disembunyikan bool
	if err := tx.QueryRow(t.Context(), `SELECT disembunyikan FROM umpan_kegiatan WHERE id = $1`, umpanID).Scan(&disembunyikan); err != nil {
		t.Fatal(err)
	}
	if tanpaAlasan.kode != fiber.StatusBadRequest || hasil.kode != fiber.StatusOK || !disembunyikan || jumlahLogModerasi(t, tx, domain.TargetUmpan) != 1 {
		t.Errorf("tanpa alasan %d, hasil %d, disembunyikan %v", tanpaAlasan.kode, hasil.kode, disembunyikan)
	}
}
