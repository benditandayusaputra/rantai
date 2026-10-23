package http

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

func appUjiUmpan(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwakil.MuatDPRD(t, tx)
	ujiwakil.MuatTema(t, tx)
	ujiwakil.MuatUmpanDPRD(t, tx)
	return NewApp(config.Server{}, Layanan{Umpan: usecase.BacaUmpan{Pembaca: repository.UmpanRepository{DB: tx}}})
}

func TestUmpanAnggota_RilisResmi_TampilTanpaCuplikanDenganMeta(t *testing.T) {
	hasil := panggil(t, appUjiUmpan(t), "/api/v1/anggota/yuke-yurike-st-mm/umpan")

	daftar, _ := hasil.respons.Data.([]any)
	if hasil.kode != fiber.StatusOK || len(daftar) == 0 || hasil.respons.Meta == nil || hasil.respons.Meta.PerPage != usecase.UmpanPerHalaman {
		t.Fatalf("kode %d, data %v, meta %+v", hasil.kode, hasil.respons.Data, hasil.respons.Meta)
	}
	pertama, _ := daftar[0].(map[string]any)
	if _, ada := pertama["cuplikan"]; ada {
		t.Error("cuplikan tidak boleh dikirim API")
	}
	if pertama["jenis_sumber"] != "resmi_dprd" || pertama["dasar_cocok"] != "situs_resmi" || pertama["ringkasan"] != nil || pertama["provenance_id"] == "" {
		t.Errorf("item = %v", pertama)
	}
}

func TestUmpanAnggota_FilterJenisMedia_Kosong(t *testing.T) {
	hasil := panggil(t, appUjiUmpan(t), "/api/v1/anggota/yuke-yurike-st-mm/umpan?jenis=media")

	if daftar, _ := hasil.respons.Data.([]any); hasil.kode != fiber.StatusOK || len(daftar) != 0 || hasil.respons.Meta.Total != 0 {
		t.Errorf("kode %d, data %v", hasil.kode, hasil.respons.Data)
	}
}

func TestUmpanAnggota_MasukanTidakSah_KodeError(t *testing.T) {
	app := appUjiUmpan(t)
	kasus := map[string]int{
		"/api/v1/anggota/yuke-yurike-st-mm/umpan?jenis=blog":   fiber.StatusBadRequest,
		"/api/v1/anggota/yuke-yurike-st-mm/umpan?halaman=0":    fiber.StatusBadRequest,
		"/api/v1/anggota/yuke-yurike-st-mm/umpan?halaman=satu": fiber.StatusBadRequest,
		"/api/v1/anggota/tidak-ada/umpan":                      fiber.StatusNotFound,
	}
	for path, harap := range kasus {
		if hasil := panggil(t, app, path); hasil.kode != harap {
			t.Errorf("%s: kode %d, harap %d", path, hasil.kode, harap)
		}
	}
}
