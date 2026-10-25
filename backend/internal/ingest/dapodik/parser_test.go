package dapodik

import (
	"errors"
	"os"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

func bacaFixture(t *testing.T, nama string) []byte {
	t.Helper()
	isi, err := os.ReadFile("testdata/" + nama)
	if err != nil {
		t.Fatal(err)
	}
	return isi
}

func TestUraiCari_HalamanTerakhir_TotalDanButir(t *testing.T) {
	halaman, err := UraiCari(bacaFixture(t, "cari_sekolah.json"))

	if err != nil || halaman.Total != 14 || len(halaman.Sekolah) != 4 {
		t.Fatalf("total %d, butir %d, err %v", halaman.Total, len(halaman.Sekolah), err)
	}
	if b := halaman.Sekolah[3]; b.NPSN != "20104477" || b.SekolahID == "" || b.Bentuk != "SD" {
		t.Errorf("butir terakhir = %+v", b)
	}
}

func TestUraiCari_HalamanKosong_KoleksiKosong(t *testing.T) {
	halaman, err := UraiCari(bacaFixture(t, "cari_sekolah_kosong.json"))

	if err != nil || halaman.Sekolah == nil || len(halaman.Sekolah) != 0 {
		t.Errorf("harap koleksi kosong, dapat %v, %v", halaman.Sekolah, err)
	}
}

func TestUraiCari_StatusGagal_Error(t *testing.T) {
	_, err := UraiCari([]byte(`{"status_code":400,"message":"Validation error"}`))

	if !errors.Is(err, ErrStatusSumber) {
		t.Errorf("harap ErrStatusSumber, dapat %v", err)
	}
}

func TestUraiDetail_SDNegeri_IdentitasDanKoordinat(t *testing.T) {
	detail, err := UraiDetail(bacaFixture(t, "full_detail.json"))

	s := detail.Sekolah
	if err != nil || s.NPSN != "20104886" || s.Bentuk != domain.BentukSD || s.Status != domain.StatusNegeri {
		t.Fatalf("sekolah %+v, err %v", s, err)
	}
	if s.Lokasi == nil || s.Lokasi.Lat != -6.1744 || s.Lokasi.Lng != 106.9113 {
		t.Errorf("lokasi = %+v", s.Lokasi)
	}
	if s.NamaKecamatanSumber != "Kec. Kelapa Gading" || detail.SemesterID != "20261" {
		t.Errorf("kecamatan %q, semester %q", s.NamaKecamatanSumber, detail.SemesterID)
	}
	if detail.FotoURL != "https://file.data.kemendikdasmen.go.id/sekolahkita/20/2010/20104886-1.jpg" ||
		detail.SekolahID != "002BB794-2BF5-E011-9477-67A74022F833" {
		t.Errorf("foto sampul %q, sekolah_id %q", detail.FotoURL, detail.SekolahID)
	}
}

func TestUraiDetail_TanpaKoordinat_LokasiKosongTidakDitebak(t *testing.T) {
	detail, err := UraiDetail(bacaFixture(t, "full_detail_tanpa_koordinat.json"))

	if err != nil || detail.Sekolah.Lokasi != nil {
		t.Errorf("lokasi harus kosong, dapat %+v, %v", detail.Sekolah.Lokasi, err)
	}
	if detail.FotoURL != "" {
		t.Errorf("tanpa foto_sekolah, foto harus kosong, dapat %q", detail.FotoURL)
	}
}

func TestUraiDetail_MTsNegeri_Madrasah(t *testing.T) {
	detail, err := UraiDetail(bacaFixture(t, "full_detail_madrasah.json"))

	if err != nil || !detail.Sekolah.Bentuk.IsMadrasah() || detail.Sekolah.Status != domain.StatusNegeri {
		t.Errorf("sekolah %+v, err %v", detail.Sekolah, err)
	}
}

func TestUraiRombel_SD_RombelTingkatSatuSampaiEnam(t *testing.T) {
	rombel, err := UraiRombel(bacaFixture(t, "rombongan_belajar.json"), domain.BentukSD)

	if err != nil || len(rombel) != 6 || rombel[1] != 2 || rombel[6] != 2 {
		t.Errorf("rombel %v, err %v", rombel, err)
	}
}

func TestUraiPesertaDidik_SD_LakiDanPerempuanDijumlahkan(t *testing.T) {
	siswa, err := UraiPesertaDidik(bacaFixture(t, "peserta_didik.json"), domain.BentukSD)

	if err != nil || siswa[6] != 63 || siswa[1] != 63 || len(siswa) != 6 {
		t.Errorf("siswa %v, err %v", siswa, err)
	}
}

func TestUraiRombel_DataKosong_KoleksiKosong(t *testing.T) {
	rombel, err := UraiRombel(bacaFixture(t, "rombongan_belajar_kosong.json"), domain.BentukMTs)

	if err != nil || rombel == nil || len(rombel) != 0 {
		t.Errorf("rombel %v, err %v", rombel, err)
	}
}
