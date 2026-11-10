package http

import (
	"fmt"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

func appUjiUmpan(t *testing.T) *fiber.App {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwakil.MuatDPRD(t, tx)
	ujiwakil.MuatTema(t, tx)
	ujiwakil.MuatUmpanDPRD(t, tx)
	return NewApp(config.Server{}, Layanan{Umpan: usecase.BacaUmpan{Pembaca: repository.UmpanRepository{DB: tx}}})
}

func TestUmpanAnggota_RilisResmi_TampilDenganCuplikanDanMeta(t *testing.T) {
	hasil := panggil(t, appUjiUmpan(t), "/api/v1/anggota/yuke-yurike-st-mm/umpan")

	daftar, _ := hasil.respons.Data.([]any)
	if hasil.kode != fiber.StatusOK || len(daftar) == 0 || hasil.respons.Meta == nil || hasil.respons.Meta.PerPage != usecase.UmpanPerHalaman {
		t.Fatalf("kode %d, data %v, meta %+v", hasil.kode, hasil.respons.Data, hasil.respons.Meta)
	}
	pertama, _ := daftar[0].(map[string]any)
	if cuplikan, _ := pertama["cuplikan"].(string); cuplikan == "" {
		t.Error("rilis resmi DPRD harus membawa cuplikannya")
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

const anggotaUmpanUji = "yuke-yurike-st-mm"

func simpanUmpanUji(t *testing.T, tx pgx.Tx, umpan ...domain.UmpanTerpaut) {
	t.Helper()
	data := usecase.DataSimpanUmpan{Umpan: umpan, Provenance: domain.Provenance{
		ID: "uji-umpan-isu", Sumber: domain.SumberMedia, URL: "https://contoh.test", DiambilPada: time.Now(), Label: domain.LabelResmi,
	}}
	if _, err := (repository.UmpanRepository{DB: tx}).SimpanUmpan(t.Context(), data); err != nil {
		t.Fatal(err)
	}
}

func umpanBertagUji(urut int, jenis domain.JenisSumberUmpan, anggotaID string, tag ...string) domain.UmpanTerpaut {
	return domain.UmpanTerpaut{
		ItemUmpan: domain.ItemUmpan{
			JenisSumber: jenis, MediaNama: "Media Uji", MediaDomain: "contoh.test", Judul: fmt.Sprintf("Umpan uji %d", urut),
			URL: fmt.Sprintf("https://contoh.test/umpan-%d", urut), TerbitPada: time.Date(2026, time.October, urut, 9, 0, 0, 0, time.UTC),
		},
		TagIsu: tag,
		Pautan: []domain.PautanUmpan{{AnggotaID: anggotaID, DasarCocok: domain.DasarNamaDanDPRD}},
	}
}

func TestIsuAnggota_RilisDanMedia_JumlahSumberTerpisah(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	ujiwakil.MuatDPRD(t, tx)
	simpanUmpanUji(t, tx,
		umpanBertagUji(1, domain.SumberUmpanResmiDPRD, anggotaUmpanUji, "sekolah"),
		umpanBertagUji(2, domain.SumberUmpanMedia, anggotaUmpanUji, "sekolah", "guru"),
		umpanBertagUji(3, domain.SumberUmpanMedia, anggotaUmpanUji, "sekolah"),
	)
	app := NewApp(config.Server{}, Layanan{Umpan: usecase.BacaUmpan{Pembaca: repository.UmpanRepository{DB: tx}}})

	hasil := panggil(t, app, "/api/v1/anggota/"+anggotaUmpanUji+"/isu")

	daftar, _ := hasil.respons.Data.([]any)
	if hasil.kode != fiber.StatusOK || len(daftar) == 0 {
		t.Fatalf("kode %d, data %v", hasil.kode, hasil.respons.Data)
	}
	sekolah, _ := daftar[0].(map[string]any)
	item, _ := sekolah["item"].([]any)
	if sekolah["tag"] != "sekolah" || sekolah["jumlah_sumber"] != float64(3) || sekolah["jumlah_rilis_resmi"] != float64(1) ||
		sekolah["jumlah_media"] != float64(2) || len(item) != 3 {
		t.Errorf("isu pertama = %v", sekolah)
	}
}

func TestIsuAnggota_AnggotaTidakAda_404(t *testing.T) {
	if hasil := panggil(t, appUjiUmpan(t), "/api/v1/anggota/tidak-ada/isu"); hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode = %d, harap 404", hasil.kode)
	}
}
