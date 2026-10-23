package domain

import "testing"

func TestNormalisasiNamaPenyedia(t *testing.T) {
	kasus := []struct {
		nama  string
		harap string
	}{
		{"CV. TIRTA ARJUNO LESTARI", "tirta arjuno lestari"},
		{"PT HERDINA MITRA PRAKARSA", "herdina mitra prakarsa"},
		{"PT. Wijaya Karya (Persero) Tbk.", "wijaya karya"},
		{"  ARTHA   UNGGUL-JAYA, CV ", "artha unggul jaya"},
		{"Cipta Karya Mandiri", "cipta karya mandiri"},
	}
	for _, k := range kasus {
		if hasil := NormalisasiNamaPenyedia(k.nama); hasil != k.harap {
			t.Errorf("%q = %q, harap %q", k.nama, hasil, k.harap)
		}
	}
}

func aturanInstansiUji() AturanInstansi {
	return AturanInstansi{
		Jenis: []AturanJenisInstansi{
			{Jenis: JenisInstansiSudin, Awalan: []string{"suku dinas", "sudin"}},
			{Jenis: JenisInstansiDinas, Awalan: []string{"dinas"}},
			{Jenis: JenisInstansiKelurahan, Awalan: []string{"kelurahan"}},
		},
		Kota: map[string][]string{
			"31.71": {"jakarta pusat", "jakpus"},
			"31.75": {"jakarta timur", "jaktim"},
		},
	}
}

func TestKlasifikasiInstansi(t *testing.T) {
	kasus := []struct {
		nama       string
		harapJenis JenisInstansi
		harapKota  string
	}{
		{"SUKU DINAS SUMBER DAYA AIR KOTA ADMINISTRASI JAKARTA PUSAT", JenisInstansiSudin, "31.71"},
		{"SUKU DINAS PEMUDA DAN OLAHRAGA KOTA - JAKTIM", JenisInstansiSudin, "31.75"},
		{"DINAS CIPTA KARYA, TATA RUANG DAN PERTANAHAN", JenisInstansiDinas, ""},
		{"KELURAHAN KAMPUNG MELAYU JAKARTA TIMUR", JenisInstansiKelurahan, "31.75"},
		{"BIRO KESEJAHTERAAN SOSIAL", JenisInstansiLainnya, ""},
		{"BADAN X JAKARTA PUSAT DAN JAKARTA TIMUR", JenisInstansiLainnya, ""},
	}
	for _, k := range kasus {
		hasil := aturanInstansiUji().Klasifikasi(k.nama)
		if hasil.Jenis != k.harapJenis || hasil.KodeKota != k.harapKota {
			t.Errorf("%q = %+v, harap %s %q", k.nama, hasil, k.harapJenis, k.harapKota)
		}
	}
}

func TestKodePaket_TenderDanNonTenderBerbeda(t *testing.T) {
	if KodePaket(true, "10033820000") != "tender-10033820000" || KodePaket(false, "10033820000") != "nontender-10033820000" {
		t.Error("kode paket tender dan non-tender tidak dibedakan")
	}
}
