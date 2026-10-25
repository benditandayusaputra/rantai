package domain

import "testing"

func TestTahunAjaranDariSemester_SemesterGasal_TahunAjaranBerjalan(t *testing.T) {
	kasus := map[string]string{"20261": "2026/2027", "20262": "2026/2027", "20251": "2025/2026"}
	for semester, harap := range kasus {
		if hasil, err := TahunAjaranDariSemester(semester); err != nil || hasil != harap {
			t.Errorf("TahunAjaranDariSemester(%q) = %q, %v; harap %q", semester, hasil, err, harap)
		}
	}
}

func TestTahunAjaranDariSemester_FormatSalah_Error(t *testing.T) {
	for _, semester := range []string{"", "2026", "abcd1", "202611"} {
		if _, err := TahunAjaranDariSemester(semester); err == nil {
			t.Errorf("semester %q harus ditolak", semester)
		}
	}
}

func TestNormalkanNamaKecamatan_VariasiPenulisan_BentukSama(t *testing.T) {
	kasus := []struct{ sumber, wilayah string }{
		{"Kec. Pal Merah", "Pal Merah"},
		{"Kec. Palmerah", "Pal Merah"},
		{"Kec. Kramat Jati", "Kramatjati"},
		{"KEC. PULO GADUNG", "Pulogadung"},
		{"Kecamatan Kelapa Gading", "Kelapa Gading"},
	}
	for _, k := range kasus {
		if NormalkanNamaKecamatan(k.sumber) != NormalkanNamaKecamatan(k.wilayah) {
			t.Errorf("%q dan %q harus setara: %q vs %q", k.sumber, k.wilayah,
				NormalkanNamaKecamatan(k.sumber), NormalkanNamaKecamatan(k.wilayah))
		}
	}
}

func TestBentukSekolah_JenjangMadrasahDanRentangTingkat(t *testing.T) {
	kasus := []struct {
		bentuk            BentukSekolah
		jenjang           Jenjang
		madrasah          bool
		pertama, terakhir int
	}{
		{BentukSD, JenjangDasar, false, 1, 6},
		{BentukMI, JenjangDasar, true, 1, 6},
		{BentukSMP, JenjangMenengahPertama, false, 7, 9},
		{BentukMTs, JenjangMenengahPertama, true, 7, 9},
		{BentukSMA, JenjangMenengahAtas, false, 10, 12},
		{BentukSMK, JenjangMenengahAtas, false, 10, 12},
		{BentukMA, JenjangMenengahAtas, true, 10, 12},
	}
	for _, k := range kasus {
		pertama, terakhir := k.bentuk.RentangTingkat()
		if k.bentuk.Jenjang() != k.jenjang || k.bentuk.IsMadrasah() != k.madrasah || pertama != k.pertama || terakhir != k.terakhir {
			t.Errorf("%s: jenjang %s, madrasah %v, tingkat %d-%d", k.bentuk, k.bentuk.Jenjang(), k.bentuk.IsMadrasah(), pertama, terakhir)
		}
	}
}

func TestStatusDariSumber_HurufBesar_Dinormalkan(t *testing.T) {
	if status, err := StatusDariSumber("NEGERI"); err != nil || status != StatusNegeri {
		t.Errorf("dapat %q, %v", status, err)
	}
	if _, err := StatusDariSumber("LAINNYA"); err == nil {
		t.Error("status tak dikenal harus ditolak")
	}
}

func TestTempatkanSekolah_UrutanSumberLokasi(t *testing.T) {
	kasus := []struct {
		nama     string
		kandidat KandidatWilayah
		harap    PenempatanSekolah
		ada      bool
	}{
		{"titik menang atas nama", KandidatWilayah{"31.72.06", "31.72.06.1001", "31.72.01"},
			PenempatanSekolah{"31.72.06", "31.72.06.1001", MetodeLokasiTitik}, true},
		{"titik di luar batas memakai nama", KandidatWilayah{KecamatanDariNama: "31.72.01"},
			PenempatanSekolah{KodeKecamatan: "31.72.01", Metode: MetodeLokasiNama}, true},
		{"tanpa titik dan nama tidak cocok", KandidatWilayah{}, PenempatanSekolah{}, false},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			hasil, ada := TempatkanSekolah(k.kandidat)

			if hasil != k.harap || ada != k.ada {
				t.Errorf("dapat %+v, %v; harap %+v, %v", hasil, ada, k.harap, k.ada)
			}
		})
	}
}
