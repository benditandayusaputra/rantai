package domain

import (
	"slices"
	"testing"
)

func wakilUji() []Wakil {
	return []Wakil{
		{ID: "zahra", Nama: "Hj. Zahrina Nurbaiti", KodeKomisi: "E"},
		{ID: "basri", Nama: "Basri Baco, SE., MM"},
		{ID: "agustina", Nama: "Agustina Hermanto", KodeKomisi: "E"},
		{ID: "dwi", Nama: "Dwi Rio Sambodo", KodeKomisi: "D"},
		{ID: "achmad", Nama: "H. Achmad Yani", KodeKomisi: "A"},
	}
}

func idWakil(daftar []Wakil) []string {
	id := []string{}
	for _, w := range daftar {
		id = append(id, w.ID)
	}
	return id
}

func TestUrutkanWakil_TemaDipilih_KomisiPengawasLebihDahuluLaluAlfabetis(t *testing.T) {
	hasil := UrutkanWakil(wakilUji(), []string{"E"})

	harap := []string{"agustina", "zahra", "achmad", "basri", "dwi"}
	if !slices.Equal(idWakil(hasil), harap) {
		t.Errorf("urutan = %v, harap %v", idWakil(hasil), harap)
	}
	if !hasil[0].IsDiKomisiPengawasTema || !hasil[1].IsDiKomisiPengawasTema || hasil[2].IsDiKomisiPengawasTema {
		t.Errorf("penanda komisi pengawas keliru: %+v", hasil)
	}
}

func TestUrutkanWakil_TanpaTema_AlfabetisTanpaGelarDanTanpaPenanda(t *testing.T) {
	hasil := UrutkanWakil(wakilUji(), []string{})

	harap := []string{"achmad", "agustina", "basri", "dwi", "zahra"}
	if !slices.Equal(idWakil(hasil), harap) {
		t.Errorf("urutan = %v, harap %v", idWakil(hasil), harap)
	}
	if slices.ContainsFunc(hasil, func(w Wakil) bool { return w.IsDiKomisiPengawasTema }) {
		t.Error("tanpa tema tidak boleh ada penanda komisi pengawas")
	}
}

func TestKomisiPengawas_DuaMasalahKomisiSama_SatuKomisi(t *testing.T) {
	tema := Tema{KonfigurasiTema: KonfigurasiTema{TanggungJawab: []TanggungJawab{{Komisi: "E"}, {Komisi: "D"}, {Komisi: "E"}}}}

	if komisi := KomisiPengawas(tema); !slices.Equal(komisi, []string{"D", "E"}) {
		t.Errorf("komisi pengawas = %v", komisi)
	}
}

func TestKeteranganKomisi_TanpaKomisi_MembedakanPimpinanDanDataKosong(t *testing.T) {
	kasus := []struct {
		profil ProfilAnggota
		harap  string
	}{
		{ProfilAnggota{Wakil: Wakil{KodeKomisi: "E"}}, ""},
		{ProfilAnggota{AlatKelengkapan: []JabatanAKD{{AlatKelengkapan: AlatKelengkapanPimpinan, Jabatan: "Ketua"}}}, KeteranganKomisiPimpinan},
		{ProfilAnggota{AlatKelengkapan: []JabatanAKD{{AlatKelengkapan: "Badan Anggaran"}}}, KeteranganKomisiTidakTersedia},
	}
	for _, k := range kasus {
		if hasil := k.profil.KeteranganKomisi(); hasil != k.harap {
			t.Errorf("keterangan = %q, harap %q", hasil, k.harap)
		}
	}
}
