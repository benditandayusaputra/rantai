package http

import (
	"bytes"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

func appUjiDataBelumDibuka(t *testing.T) *fiber.App {
	t.Helper()
	var k domain.KonfigurasiDataBelumDibuka
	if err := config.BacaYAML("../../../config/data_belum_dibuka.yaml", &k); err != nil {
		t.Fatal(err)
	}
	return NewApp(config.Server{}, Layanan{BelumBuka: usecase.BacaDataBelumDibuka{Konfigurasi: k}})
}

func mintaPermohonan(t *testing.T, kode, badan string) hasilPanggilan {
	t.Helper()
	permintaan := httptest.NewRequestWithContext(t.Context(), fiber.MethodPost, "/api/v1/data-belum-dibuka/"+kode+"/permohonan", strings.NewReader(badan))
	permintaan.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	return kirim(t, appUjiDataBelumDibuka(t), permintaan)
}

func TestDataBelumDibuka_Daftar_EmpatDataDenganPenggantiDanPPID(t *testing.T) {
	hasil := panggil(t, appUjiDataBelumDibuka(t), "/api/v1/data-belum-dibuka")

	data, _ := hasil.respons.Data.(map[string]any)
	daftar, _ := data["data"].([]any)
	ppid, _ := data["ppid"].(map[string]any)
	if hasil.kode != fiber.StatusOK || len(daftar) != 4 || !strings.Contains(ppid["alamat"].(string), "Kebon Sirih") {
		t.Fatalf("kode %d, data %v", hasil.kode, data)
	}
	for _, item := range daftar {
		pengganti, _ := item.(map[string]any)["pengganti"].(map[string]any)
		if pengganti["id"] == "" {
			t.Errorf("tanpa pengganti: %v", item)
		}
	}
}

func TestDataBelumDibuka_Permohonan_TeksBerdasarUUDanTautanPDF(t *testing.T) {
	hasil := mintaPermohonan(t, "kehadiran_per_anggota", `{"bahasa":"id"}`)

	data, _ := hasil.respons.Data.(map[string]any)
	teks, _ := data["teks"].(string)
	if hasil.kode != fiber.StatusOK || !strings.Contains(teks, "Undang-Undang Nomor 14 Tahun 2008") ||
		data["pdf_url"] != "/api/v1/data-belum-dibuka/kehadiran_per_anggota/permohonan/pdf?bahasa=id" {
		t.Errorf("kode %d, data %v", hasil.kode, data)
	}
}

func TestDataBelumDibuka_PermohonanTanpaBadan_BahasaIndonesia(t *testing.T) {
	hasil := mintaPermohonan(t, "pokir_per_anggota", "")

	data, _ := hasil.respons.Data.(map[string]any)
	if hasil.kode != fiber.StatusOK || !strings.HasPrefix(data["teks"].(string), "Permohonan Informasi Publik") {
		t.Errorf("kode %d, data %v", hasil.kode, data)
	}
}

func TestDataBelumDibuka_PermohonanTidakValid(t *testing.T) {
	kasus := []struct {
		nama, kode, badan string
		harap             int
	}{
		{"kode tak dikenal", "gaji_anggota", `{}`, fiber.StatusNotFound},
		{"bahasa asing", "suara_per_anggota", `{"bahasa":"fr"}`, fiber.StatusBadRequest},
		{"badan bukan JSON", "suara_per_anggota", `bahasa=id`, fiber.StatusBadRequest},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if hasil := mintaPermohonan(t, k.kode, k.badan); hasil.kode != k.harap {
				t.Errorf("kode %d, harap %d", hasil.kode, k.harap)
			}
		})
	}
}

func TestDataBelumDibuka_UnduhPDF_BerkasPDFLampiran(t *testing.T) {
	resp, err := appUjiDataBelumDibuka(t).Test(httptest.NewRequestWithContext(t.Context(), fiber.MethodGet,
		"/api/v1/data-belum-dibuka/laporan_reses_per_anggota/permohonan/pdf?bahasa=en", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	isi, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != fiber.StatusOK || !bytes.HasPrefix(isi, []byte("%PDF-")) ||
		!strings.Contains(resp.Header.Get(fiber.HeaderContentDisposition), "permohonan-laporan_reses_per_anggota.pdf") {
		t.Errorf("kode %d, disposisi %s", resp.StatusCode, resp.Header.Get(fiber.HeaderContentDisposition))
	}
}
