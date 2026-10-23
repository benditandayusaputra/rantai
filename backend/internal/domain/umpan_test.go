package domain

import (
	"slices"
	"testing"
	"unicode/utf8"
)

func anggotaUmpanUji() []Anggota {
	return []Anggota{
		{ID: "yuke", Nama: "Hj. Yuke Yurike, ST., MM", Fraksi: "Fraksi PDI Perjuangan"},
		{ID: "wibi", Nama: "Wibi Andrino, SH., MH", Fraksi: "Fraksi NasDem"},
		{ID: "ismail-1", Nama: "H. Ismail, S.Pd."},
		{ID: "ismail-2", Nama: "Ismail"},
		{ID: "fatih", Nama: "Muhammad Al Fatih, S.Ak", NamaLain: []string{"Muhammad Alfatih"}},
		{ID: "mujiyono", Nama: "Mujiyono, SE"},
		{ID: "johan", Nama: "H. Ali Muhammad Johan C., SE., MH"},
	}
}

func TestAnggotaDisebut_NamaTanpaGelarDiTeks_Ditemukan(t *testing.T) {
	pencari := BuatPencariNama(anggotaUmpanUji())

	disebut := pencari.AnggotaDisebut("ujar Ketua Komisi D DPRD DKI Jakarta Yuke Yurike, Rabu (14/10). Wibi Andrino juga hadir.")

	if !slices.Equal(disebut, []string{"wibi", "yuke"}) {
		t.Errorf("disebut = %v", disebut)
	}
}

func TestAnggotaDisebut_NamaSatuKataAtauGanda_TidakDipautkan(t *testing.T) {
	pencari := BuatPencariNama(anggotaUmpanUji())

	if disebut := pencari.AnggotaDisebut("Anggota DPRD Ismail menghadiri rapat."); len(disebut) != 0 {
		t.Errorf("nama satu kata yang dipakai dua anggota tidak boleh dipautkan, dapat %v", disebut)
	}
}

func TestAnggotaDisebut_PotonganNama_TidakDianggapSebutan(t *testing.T) {
	pencari := BuatPencariNama(anggotaUmpanUji())

	if disebut := pencari.AnggotaDisebut("Yuke Yurikesari berbicara"); len(disebut) != 0 {
		t.Errorf("harap tidak cocok di tengah kata, dapat %v", disebut)
	}
}

func TestPautkanRilisResmi_NamaDiIsi_DasarSitusResmi(t *testing.T) {
	item := ItemUmpan{Judul: "Penguatan Pengawasan Lalu Lintas", Isi: "kata Muhammad Alfatih di Jakarta"}

	pautan := BuatPencariNama(anggotaUmpanUji()).PautkanRilisResmi(item)

	if !slices.Equal(pautan, []PautanUmpan{{AnggotaID: "fatih", DasarCocok: DasarSitusResmi}}) {
		t.Errorf("pautan = %v", pautan)
	}
}

func TestTandaiTagIsu_KataKunciDuaKata_DitandaiTanpaPotonganKata(t *testing.T) {
	tag := TandaiTagIsu("Daya tampung SMP negeri dan SPMB 2026 dibahas; sekolahnya penuh", []string{"spmb", "daya tampung", "sekolah", "guru"})

	if !slices.Equal(tag, []string{"daya tampung", "spmb"}) {
		t.Errorf("tag = %v", tag)
	}
}

func TestPotongCuplikan_TeksPanjang_TidakMelebihiBatas(t *testing.T) {
	panjang := ""
	for range 100 {
		panjang += "kata  panjang "
	}
	if n := utf8.RuneCountInString(PotongCuplikan(panjang)); n != BatasCuplikanUmpan {
		t.Errorf("panjang cuplikan = %d, harap %d", n, BatasCuplikanUmpan)
	}
}

func TestAnggotaDisebut_NamaSatuKata_HanyaSetelahDPRDDKI(t *testing.T) {
	pencari := BuatPencariNama(anggotaUmpanUji())

	kasus := []struct {
		teks  string
		harap []string
	}{
		{"ujar Ketua Komisi A DPRD DKI Jakarta Mujiyono, Senin.", []string{"mujiyono"}},
		{"kata anggota DPRD DKI Mujiyono", []string{"mujiyono"}},
		{"Warga bernama Mujiyono mengeluhkan banjir", []string{}},
	}
	for _, k := range kasus {
		if disebut := pencari.AnggotaDisebut(k.teks); !slices.Equal(disebut, k.harap) {
			t.Errorf("%q: disebut %v, harap %v", k.teks, disebut, k.harap)
		}
	}
}

func TestAnggotaDisebut_InisialDiBelakangNama_TetapDitemukan(t *testing.T) {
	if disebut := BuatPencariNama(anggotaUmpanUji()).AnggotaDisebut("Ali Muhammad Johan meninjau pasar"); !slices.Equal(disebut, []string{"johan"}) {
		t.Errorf("disebut = %v", disebut)
	}
}
