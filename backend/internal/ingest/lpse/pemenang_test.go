package lpse

import (
	"testing"
	"time"
)

func TestUraiPemenang_EvaluasiAdaPemenang_NamaDanHasilNegosiasi(t *testing.T) {
	daftar, err := UraiPemenang(bacaFixture(t, "pemenang_ada.html"))

	if err != nil || len(daftar) != 1 {
		t.Fatalf("pemenang = %v, err = %v", daftar, err)
	}
	if daftar[0].Nama != "CV. TIRTA ARJUNO LESTARI" || daftar[0].Nilai != 196351675.60 {
		t.Errorf("pemenang = %+v", daftar[0])
	}
}

func TestUraiPemenang_PaketBatalTanpaPemenang_DaftarKosong(t *testing.T) {
	daftar, err := UraiPemenang(bacaFixture(t, "pemenang_kosong.html"))

	if err != nil || daftar == nil || len(daftar) != 0 {
		t.Errorf("pemenang = %v, err = %v", daftar, err)
	}
}

func TestUraiPemenang_BerkontrakAda_NilaiDariHargaKontrak(t *testing.T) {
	daftar, err := UraiPemenang(bacaFixture(t, "kontrak_ada.html"))

	if err != nil || len(daftar) != 1 || daftar[0].Nama != "PT HERDINA MITRA PRAKARSA" || daftar[0].Nilai != 196861000 {
		t.Errorf("pemenang = %v, err = %v", daftar, err)
	}
}

func TestUraiPemenang_KontrakBelumDibuat_DaftarKosong(t *testing.T) {
	daftar, err := UraiPemenang(bacaFixture(t, "kontrak_kosong.html"))

	if err != nil || len(daftar) != 0 {
		t.Errorf("pemenang = %v, err = %v", daftar, err)
	}
}

func TestUraiTanggalKontrak_JadwalNonTender_AkhirTahapTandaTangan(t *testing.T) {
	tanggal, ada, err := UraiTanggalKontrak(bacaFixture(t, "jadwal.html"))

	harap := time.Date(2026, time.January, 2, 23, 59, 0, 0, zonaWIB)
	if err != nil || !ada || !tanggal.Equal(harap) {
		t.Errorf("tanggal = %v, ada = %v, err = %v", tanggal, ada, err)
	}
}
