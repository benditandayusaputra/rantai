package domain

import (
	"errors"
	"testing"
)

func anggotaStrukturUji() []Anggota {
	return []Anggota{
		{ID: "ima-mahdiah", Nama: "Ima Mahdiah"},
		{ID: "muhammad-hasan-abdillah", Nama: "Muhammad Hasan Abdilah"},
		{ID: "h-ismail-s-pd", Nama: "H. Ismail, S.Pd., M.H."},
		{ID: "ismail", Nama: "Ismail"},
		{ID: "hj-ida-mahmudah", Nama: "Hj. Ida Mahmudah"},
	}
}

func TestValidasiKomisi_TanpaBidangAtauKodeLuar_Error(t *testing.T) {
	kasus := []Komisi{{Kode: "E"}, {Kode: "F", Bidang: []string{"Dinas Sosial"}}}
	for _, k := range kasus {
		if err := k.Validasi(); !errors.Is(err, ErrKomisiTanpaBidang) {
			t.Errorf("komisi %+v: harap ErrKomisiTanpaBidang, dapat %v", k, err)
		}
	}
	if err := (Komisi{Kode: "E", Bidang: []string{"Dinas Pendidikan"}}).Validasi(); err != nil {
		t.Errorf("harap sah, dapat %v", err)
	}
}

func TestCocokkanStruktur_NamaBergelarDanEjaanBeda_DicocokkanLewatNamaLaluSlug(t *testing.T) {
	halaman := HalamanStruktur{AlatKelengkapan: "Komisi A", KodeKomisi: "A", Baris: []BarisStruktur{
		{Nama: "Ima Mahdiah", Jabatan: "Koordinator", Slug: "ima-mahdiah"},
		{Nama: "Muhammad Hasan Abdillah", Jabatan: "Anggota", Slug: "muhammad-hasan-abdillah"},
		{Nama: "Ida Mahmudah", Jabatan: "Anggota", Slug: "ima-mahdiah"},
	}}

	cocok, tidakCocok := CocokkanStruktur(anggotaStrukturUji(), halaman)

	harap := []string{"ima-mahdiah", "muhammad-hasan-abdillah", "hj-ida-mahmudah"}
	if len(tidakCocok) != 0 || len(cocok) != len(harap) {
		t.Fatalf("cocok %+v, tidak cocok %v", cocok, tidakCocok)
	}
	for i, id := range harap {
		if cocok[i].AnggotaID != id || cocok[i].KodeKomisi != "A" {
			t.Errorf("baris %d = %+v, harap anggota %s komisi A", i, cocok[i], id)
		}
	}
}

func TestCocokkanStruktur_NamaGandaTanpaSlugDanBukanAnggota_Dilaporkan(t *testing.T) {
	halaman := HalamanStruktur{AlatKelengkapan: "Badan Kehormatan", Baris: []BarisStruktur{
		{Nama: "Ismail", Jabatan: "Anggota"},
		{Nama: "Augustinus, SE., MM", Jabatan: "Sekretaris Bukan Anggota"},
	}}

	cocok, tidakCocok := CocokkanStruktur(anggotaStrukturUji(), halaman)

	if len(cocok) != 0 || len(tidakCocok) != 2 {
		t.Errorf("harap dua baris dilaporkan tanpa ditebak, dapat cocok %+v tidak cocok %v", cocok, tidakCocok)
	}
}

func TestTentukanKomisi_KoordinatorDanAKDLain_TidakMenjadiKomisi(t *testing.T) {
	jabatan := []JabatanAKD{
		{AnggotaID: "ima-mahdiah", Jabatan: JabatanKoordinator, KodeKomisi: "A"},
		{AnggotaID: "ima-mahdiah", AlatKelengkapan: "Badan Anggaran", Jabatan: "Wakil Ketua"},
		{AnggotaID: "ismail", Jabatan: "Anggota", KodeKomisi: "E"},
	}

	komisi, laporan := TentukanKomisi(jabatan)

	if len(laporan) != 0 || len(komisi) != 1 || komisi["ismail"] != "E" {
		t.Errorf("komisi = %v, laporan = %v", komisi, laporan)
	}
}

func TestTentukanKomisi_AnggotaDiDuaKomisi_DikosongkanDanDilaporkan(t *testing.T) {
	jabatan := []JabatanAKD{
		{AnggotaID: "ismail", Jabatan: "Anggota", KodeKomisi: "B"},
		{AnggotaID: "ismail", Jabatan: "Anggota", KodeKomisi: "E"},
	}

	komisi, laporan := TentukanKomisi(jabatan)

	if _, ada := komisi["ismail"]; ada || len(laporan) != 1 {
		t.Errorf("harap komisi dikosongkan dan dilaporkan, dapat %v, %v", komisi, laporan)
	}
}

func TestCocokkanStruktur_BarisTanpaSlug_TidakDicocokkanKeAnggotaTanpaSlug(t *testing.T) {
	anggota := []Anggota{{ID: "pandapotan-sinaga-se-mm", Nama: "Pandapotan Sinaga"}}
	halaman := HalamanStruktur{AlatKelengkapan: "Badan Kehormatan", Baris: []BarisStruktur{{Nama: "Augustinus, SE., MM", Jabatan: "Sekretaris DPRD"}}}

	cocok, tidakCocok := CocokkanStruktur(anggota, halaman)

	if len(cocok) != 0 || len(tidakCocok) != 1 {
		t.Errorf("harap dilaporkan, dapat cocok %+v", cocok)
	}
}

func TestTentukanKomisi_PimpinanDenganJabatanSalahKetik_TetapTanpaKomisi(t *testing.T) {
	jabatan := []JabatanAKD{
		{AnggotaID: "basri-baco-se-mm", AlatKelengkapan: AlatKelengkapanPimpinan, Jabatan: "Wakil Ketua"},
		{AnggotaID: "basri-baco-se-mm", Jabatan: "Koordintor", KodeKomisi: "B"},
	}

	komisi, _ := TentukanKomisi(jabatan)

	if _, ada := komisi["basri-baco-se-mm"]; ada {
		t.Errorf("pimpinan DPRD tidak boleh tercatat sebagai anggota komisi, dapat %v", komisi)
	}
}
