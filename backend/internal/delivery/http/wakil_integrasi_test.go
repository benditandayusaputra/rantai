package http

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
)

var polaKunciKinerja = regexp.MustCompile(`(?i)"[a-z_]*(skor|peringkat|rating|kinerja|ranking|kehadiran)[a-z_]*"\s*:`)

func appUjiWakil(t *testing.T) *fiber.App {
	t.Helper()
	app, _ := appUjiWakilDanTransaksi(t)
	return app
}

func appUjiWakilDanTransaksi(t *testing.T) (*fiber.App, pgx.Tx) {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwakil.MuatDPRD(t, tx)
	muatTema(t, tx)
	tema := usecase.BacaTema{Pembaca: repository.TemaRepository{DB: tx}}
	wakil := usecase.BacaWakil{Pembaca: repository.WakilRepository{DB: tx}, Tema: tema, Umpan: repository.UmpanRepository{DB: tx}}
	return NewApp(config.Server{}, Layanan{Wakil: wakil}), tx
}

func daftarAnggota(t *testing.T, dapil any) []map[string]any {
	t.Helper()
	d, _ := dapil.(map[string]any)
	mentah, _ := d["anggota"].([]any)
	daftar := []map[string]any{}
	for _, item := range mentah {
		a, _ := item.(map[string]any)
		daftar = append(daftar, a)
	}
	return daftar
}

func TestWakilKecamatan_TemaPendidikan_KomisiELebihDahulu(t *testing.T) {
	hasil := panggil(t, appUjiWakil(t), "/api/v1/kecamatan/31.74.01/wakil?tema=pendidikan")

	data, _ := hasil.respons.Data.(map[string]any)
	anggota := daftarAnggota(t, data["dapil"])
	if hasil.kode != fiber.StatusOK || len(anggota) != 12 {
		t.Fatalf("kode %d, jumlah anggota %d", hasil.kode, len(anggota))
	}
	sudahLewatPengawas := false
	for _, a := range anggota {
		pengawas, _ := a["duduk_di_komisi_pengawas"].(bool)
		komisi, _ := a["komisi"].(map[string]any)
		if pengawas && (sudahLewatPengawas || komisi["kode"] != "E") {
			t.Errorf("anggota komisi pengawas %v tidak berada di depan atau bukan Komisi E", a["id"])
		}
		sudahLewatPengawas = sudahLewatPengawas || !pengawas
	}
	if pertama, _ := anggota[0]["duduk_di_komisi_pengawas"].(bool); !pertama {
		t.Error("harap anggota pertama duduk di Komisi E")
	}
}

func TestWakilKecamatan_KodeKelurahan_MemakaiKecamatanInduk(t *testing.T) {
	hasil := panggil(t, appUjiWakil(t), "/api/v1/kecamatan/31.72.06.1003/wakil")

	data, _ := hasil.respons.Data.(map[string]any)
	kecamatan, _ := data["kecamatan"].(map[string]any)
	dapil, _ := data["dapil"].(map[string]any)
	if kecamatan["kode"] != "31.72.06" || dapil["nomor"] != float64(2) {
		t.Errorf("kecamatan %v, dapil %v", kecamatan["kode"], dapil["nomor"])
	}
}

func TestWakilKecamatan_TemaTidakAda_404(t *testing.T) {
	if hasil := panggil(t, appUjiWakil(t), "/api/v1/kecamatan/31.74.01/wakil?tema=banjir"); hasil.kode != fiber.StatusNotFound {
		t.Errorf("kode = %d, harap 404", hasil.kode)
	}
}

func TestAmbilAnggota_PimpinanDPRD_TanpaKomisiDenganKeteranganDanJabatan(t *testing.T) {
	hasil := panggil(t, appUjiWakil(t), "/api/v1/anggota/ima-mahdiah")

	data, _ := hasil.respons.Data.(map[string]any)
	jabatan, _ := data["alat_kelengkapan"].([]any)
	if hasil.kode != fiber.StatusOK || data["komisi"] != nil || data["keterangan_komisi"] != "pimpinan_dprd" || len(jabatan) == 0 {
		t.Fatalf("profil = %v", data)
	}
	if pertama, _ := jabatan[0].(map[string]any); pertama["alat_kelengkapan"] != "Pimpinan DPRD" {
		t.Errorf("harap jabatan pimpinan di urutan pertama, dapat %v", pertama)
	}
}

func TestAmbilAnggota_AnggotaKomisi_BidangDanTautanTersedia(t *testing.T) {
	hasil := panggil(t, appUjiWakil(t), "/api/v1/anggota/h-matnoor-tindoan-sh")

	data, _ := hasil.respons.Data.(map[string]any)
	komisi, _ := data["komisi"].(map[string]any)
	bidang, _ := komisi["bidang"].([]any)
	if komisi["kode"] == nil || len(bidang) == 0 || data["url_profil"] == "" || data["lhkpn_url"] == "" {
		t.Errorf("profil = %v", data)
	}
}

func TestAmbilDapil_Dapil8_KursiBerlabelResmi(t *testing.T) {
	app := appUjiWakil(t)
	kecamatan, _ := panggil(t, app, "/api/v1/kecamatan/31.74.01/wakil").respons.Data.(map[string]any)
	dapilKecamatan, _ := kecamatan["dapil"].(map[string]any)

	hasil := panggil(t, app, fmt.Sprintf("/api/v1/dapil/%v", dapilKecamatan["id"]))

	data, _ := hasil.respons.Data.(map[string]any)
	kursi, _ := data["jumlah_kursi"].(map[string]any)
	if hasil.kode != fiber.StatusOK || kursi["nilai"] != float64(12) || kursi["label"] != "resmi" || kursi["provenance_id"] == "" {
		t.Errorf("dapil = %v", data)
	}
}

func TestAPIWakil_SemuaEndpoint_TanpaAngkaKinerjaAtauPeringkat(t *testing.T) {
	app := appUjiWakil(t)
	for _, path := range []string{
		"/api/v1/kecamatan/31.74.01/wakil?tema=pendidikan",
		"/api/v1/anggota/ima-mahdiah",
		"/api/v1/anggota/h-matnoor-tindoan-sh",
	} {
		mentah, _ := json.Marshal(panggil(t, app, path).respons.Data)
		if polaKunciKinerja.Match(mentah) {
			t.Errorf("%s memuat angka kinerja: %s", path, polaKunciKinerja.Find(mentah))
		}
	}
}

func TestWakilKecamatan_DenganUmpan_TigaUmpanBertagTemaTerbaruUntukPengawas(t *testing.T) {
	app, tx := appUjiWakilDanTransaksi(t)
	awal, _ := panggil(t, app, "/api/v1/kecamatan/31.74.01/wakil?tema=pendidikan").respons.Data.(map[string]any)
	anggota := daftarAnggota(t, awal["dapil"])
	pengawas, _ := anggota[0]["id"].(string)
	bukanPengawas, _ := anggota[len(anggota)-1]["id"].(string)
	simpanUmpanUji(t, tx,
		umpanBertagUji(1, domain.SumberUmpanMedia, pengawas, "sekolah"),
		umpanBertagUji(2, domain.SumberUmpanResmiDPRD, pengawas, "zonasi"),
		umpanBertagUji(3, domain.SumberUmpanMedia, pengawas, "guru"),
		umpanBertagUji(4, domain.SumberUmpanResmiDPRD, pengawas, "spmb"),
		umpanBertagUji(5, domain.SumberUmpanMedia, pengawas, "banjir"),
		umpanBertagUji(6, domain.SumberUmpanMedia, bukanPengawas, "sekolah"),
	)

	hasil := panggil(t, app, "/api/v1/kecamatan/31.74.01/wakil?tema=pendidikan&dengan_umpan=true")

	data, _ := hasil.respons.Data.(map[string]any)
	umpanIsu, _ := data["umpan_isu_pengawas"].(map[string]any)
	item, _ := umpanIsu[pengawas].([]any)
	if hasil.kode != fiber.StatusOK || len(item) != domain.BatasUmpanIsuPerWakil {
		t.Fatalf("kode %d, umpan pengawas %v", hasil.kode, umpanIsu[pengawas])
	}
	if terbaru, _ := item[0].(map[string]any); terbaru["judul"] != "Umpan uji 4" || terbaru["jenis_sumber"] != "resmi_dprd" {
		t.Errorf("item terbaru = %v", terbaru)
	}
	if _, ada := umpanIsu[bukanPengawas]; ada {
		t.Error("anggota di luar komisi pengawas tidak boleh membawa umpan isu")
	}
}

func TestWakilKecamatan_TanpaDenganUmpan_TanpaKunciUmpan(t *testing.T) {
	data, _ := panggil(t, appUjiWakil(t), "/api/v1/kecamatan/31.74.01/wakil?tema=pendidikan").respons.Data.(map[string]any)

	if _, ada := data["umpan_isu_pengawas"]; ada {
		t.Error("umpan_isu_pengawas hanya dikirim bila dengan_umpan=true")
	}
}

func TestWakilKecamatan_DenganUmpanTidakSah_400(t *testing.T) {
	app := appUjiWakil(t)
	for _, path := range []string{
		"/api/v1/kecamatan/31.74.01/wakil?dengan_umpan=true",
		"/api/v1/kecamatan/31.74.01/wakil?tema=pendidikan&dengan_umpan=ya",
	} {
		if hasil := panggil(t, app, path); hasil.kode != fiber.StatusBadRequest {
			t.Errorf("%s: kode %d, harap 400", path, hasil.kode)
		}
	}
}
