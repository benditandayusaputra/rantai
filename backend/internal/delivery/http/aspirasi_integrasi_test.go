package http

import (
	"bytes"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const kodeKelapaGading = "31.72.06"

func appUjiAspirasi(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	muatTema(t, tx)
	ujiwakil.MuatDPRD(t, tx)
	if _, err := tx.Exec(t.Context(), sqlBelanjaUji); err != nil {
		t.Fatal(err)
	}
	var kosakata domain.KosakataAspirasi
	var jalur domain.KonfigurasiJalurAspirasi
	if err := config.BacaYAML("../../../config/aspirasi.yaml", &kosakata); err != nil {
		t.Fatal(err)
	}
	if err := config.BacaYAML("../../../config/jalur_aspirasi.yaml", &jalur); err != nil {
		t.Fatal(err)
	}
	tema := usecase.BacaTema{Pembaca: repository.TemaRepository{DB: tx}}
	wilayah := usecase.CariWilayah{Pembaca: repository.WilayahRepository{DB: tx}}
	tahunUji := func() time.Time { return time.Date(2026, time.October, 23, 0, 0, 0, 0, time.UTC) }
	return NewApp(config.Server{}, Layanan{
		Aspirasi: usecase.SusunAspirasi{
			Tema: tema, Wilayah: wilayah,
			Belanja:    usecase.BandingBelanja{Tema: tema, Pembaca: repository.BelanjaTemaRepository{DB: tx}, Sekarang: tahunUji},
			Wakil:      usecase.BacaWakil{Pembaca: repository.WakilRepository{DB: tx}, Tema: tema},
			Provenance: usecase.AmbilProvenance{Pembaca: repository.ProvenanceRepository{DB: tx}},
			Jalur:      jalur, Kosakata: kosakata,
		},
		Jalur: usecase.BacaJalurAspirasi{Wilayah: wilayah, Jalur: jalur},
	})
}

func kirimAspirasi(t *testing.T, app *fiber.App, badan string) hasilPanggilan {
	t.Helper()
	permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodPost, "/api/v1/aspirasi", strings.NewReader(badan))
	permintaan.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return kirim(t, app, permintaan)
}

func TestAspirasi_KelapaGading_TeksBerdataBersumberDanTautanPDF(t *testing.T) {
	hasil := kirimAspirasi(t, appUjiAspirasi(t), `{"kode_kecamatan":"`+kodeKelapaGading+`","tema":"pendidikan","masalah":"kursi_smp_kurang"}`)

	data, _ := hasil.respons.Data.(map[string]any)
	teks, _ := data["teks"].(string)
	pdfURL, _ := data["pdf_url"].(string)
	if hasil.kode != fiber.StatusOK {
		t.Fatalf("kode %d, pesan %s", hasil.kode, hasil.respons.Message)
	}
	for _, harap := range []string{"kebutuhan 1.200 siswa", "label estimasi", "2 paket dengan total HPS Rp150.000.000",
		"Sumber: lpse", "Komisi E", "Suku Dinas Pendidikan Kota Administrasi Jakarta Utara", "Dewan Kota/Kabupaten (Kota Administrasi Jakarta Utara)"} {
		if !strings.Contains(teks, harap) {
			t.Errorf("teks tidak memuat %q:\n%s", harap, teks)
		}
	}
	if !strings.HasPrefix(pdfURL, "/api/v1/aspirasi/pdf?") || !strings.Contains(pdfURL, "bahasa=id") {
		t.Errorf("pdf_url = %q", pdfURL)
	}
}

func TestAspirasi_UnduhPDF_BerkasPDFLampiran(t *testing.T) {
	app := appUjiAspirasi(t)
	permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodGet,
		"/api/v1/aspirasi/pdf?kode_kecamatan="+kodeKelapaGading+"&tema=pendidikan&masalah=kursi_smp_kurang&bahasa=en", nil)

	resp, err := app.Test(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	isi, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != fiber.StatusOK || resp.Header.Get(fiber.HeaderContentType) != tipePDF || !bytes.HasPrefix(isi, []byte("%PDF-")) ||
		!strings.Contains(resp.Header.Get(fiber.HeaderContentDisposition), "aspirasi-"+kodeKelapaGading+".pdf") {
		t.Errorf("kode %d, tipe %s, disposisi %s, awal %q", resp.StatusCode, resp.Header.Get(fiber.HeaderContentType),
			resp.Header.Get(fiber.HeaderContentDisposition), isi[:min(len(isi), 8)])
	}
}

func TestAspirasi_MasalahTidakDikenal_400(t *testing.T) {
	hasil := kirimAspirasi(t, appUjiAspirasi(t), `{"kode_kecamatan":"`+kodeKelapaGading+`","tema":"pendidikan","masalah":"banjir"}`)

	if hasil.kode != fiber.StatusBadRequest {
		t.Errorf("kode %d, pesan %s", hasil.kode, hasil.respons.Message)
	}
}

func TestJalurAspirasi_KodeKelurahan_EmpatJalurUntukKecamatannya(t *testing.T) {
	hasil := panggil(t, appUjiAspirasi(t), "/api/v1/jalur-aspirasi?kode_wilayah="+kodeKelapaGading+".1001")

	data, _ := hasil.respons.Data.(map[string]any)
	jalur, _ := data["jalur"].([]any)
	wilayah, _ := data["wilayah"].(map[string]any)
	if hasil.kode != fiber.StatusOK || len(jalur) != 4 || wilayah["kode"] != kodeKelapaGading || data["diverifikasi_pada"] != "2026-10-24" {
		t.Fatalf("kode %d, data %v", hasil.kode, data)
	}
	cara, _ := jalur[0].(map[string]any)["cara"].(map[string]any)
	if !strings.Contains(cara["id"].(string), "Kecamatan Kelapa Gading") {
		t.Errorf("cara LMK = %v", cara)
	}
}
