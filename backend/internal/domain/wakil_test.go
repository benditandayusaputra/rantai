package domain

import (
	"errors"
	"slices"
	"testing"
)

func TestNamaTanpaGelar_GelarDepanDanBelakang_Dibuang(t *testing.T) {
	kasus := []struct {
		nama  string
		harap string
	}{
		{"H. Matnoor Tindoan, SH", "Matnoor Tindoan"},
		{"Hj. Ida Mahmudah", "Ida Mahmudah"},
		{"Mayjen TNI (Purn) H. Ferrial Sofyan, SE", "Ferrial Sofyan"},
		{"CN. DR. Lauw Siegvrieda, SH., MH", "Lauw Siegvrieda"},
		{"Drs. K.H. Muhammad Thamrin", "Muhammad Thamrin"},
		{"Dr. Dian Pratama, Sp.OG.", "Dian Pratama"},
		{"Bun Joi Phiau S.T", "Bun Joi Phiau"},
		{"Riano P. Ahmad", "Riano P. Ahmad"},
		{"Ima Mahdiah", "Ima Mahdiah"},
	}
	for _, k := range kasus {
		if hasil := NamaTanpaGelar(k.nama); hasil != k.harap {
			t.Errorf("NamaTanpaGelar(%q) = %q, harap %q", k.nama, hasil, k.harap)
		}
	}
}

func TestKunciNama_EjaanSpasiBerbeda_KunciSama(t *testing.T) {
	if KunciNama("Muhammad Al Fatih") != KunciNama("Muhammad Alfatih, S.H.") {
		t.Error("harap kunci sama untuk ejaan dengan dan tanpa spasi")
	}
}

func dapilUji() []Dapil {
	return []Dapil{
		{Nomor: 1, JumlahKursi: 1, Kecamatan: []KecamatanDapil{{"31.71.01", "Gambir"}}},
		{Nomor: 2, JumlahKursi: 1, Kecamatan: []KecamatanDapil{{"31.72.03", "Koja"}}},
	}
}

var kecamatanUji = map[string]string{"31.71.01": "Gambir", "31.72.03": "Koja"}

func TestValidasiCakupanDapil_SetiapKecamatanTepatSatu_TanpaError(t *testing.T) {
	if err := ValidasiCakupanDapil(dapilUji(), kecamatanUji); err != nil {
		t.Errorf("harap sah, dapat %v", err)
	}
}

func TestValidasiCakupanDapil_CakupanRusak_ErrorSesuaiAturan(t *testing.T) {
	kasus := []struct {
		nama  string
		ubah  func([]Dapil) []Dapil
		harap error
	}{
		{"kecamatan ganda", func(d []Dapil) []Dapil {
			d[1].Kecamatan = append(d[1].Kecamatan, KecamatanDapil{"31.71.01", "Gambir"})
			return d
		}, ErrKecamatanDapilGanda},
		{"kecamatan tanpa dapil", func(d []Dapil) []Dapil { return d[:1] }, ErrKecamatanTanpaDapil},
		{"kode tak dikenal", func(d []Dapil) []Dapil {
			d[0].Kecamatan = append(d[0].Kecamatan, KecamatanDapil{"31.99.01", "Entah"})
			return d
		}, ErrKecamatanDapilTidakDikenal},
		{"nama tidak sesuai", func(d []Dapil) []Dapil {
			d[0].Kecamatan[0].Nama = "Gambier"
			return d
		}, ErrKecamatanDapilTidakDikenal},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := ValidasiCakupanDapil(k.ubah(dapilUji()), kecamatanUji); !errors.Is(err, k.harap) {
				t.Errorf("harap %v, dapat %v", k.harap, err)
			}
		})
	}
}

func TestValidasiKursiDapil_AnggotaKurangDariKursi_Error(t *testing.T) {
	anggota := []Anggota{{ID: "a", NomorDapil: 1}}
	if err := ValidasiKursiDapil(dapilUji(), anggota); !errors.Is(err, ErrJumlahKursiTidakSesuai) {
		t.Errorf("harap ErrJumlahKursiTidakSesuai, dapat %v", err)
	}
}

func TestGabungkanAnggota_SlugUnik_DicocokkanWalauNamaBerbeda(t *testing.T) {
	dapil := []BarisDapil{{NomorDapil: 1, Nama: "Muhammad Hasan Abdilah", Partai: "PKS", Slug: "muhammad-hasan-abdillah"}}
	fraksi := []AnggotaFraksi{{Fraksi: "Fraksi PKS", BarisStruktur: BarisStruktur{Nama: "Muhammad Hasan Abdillah", Slug: "muhammad-hasan-abdillah"}}}

	anggota, err := GabungkanAnggota(dapil, fraksi)

	if err != nil || len(anggota) != 1 || anggota[0].ID != "muhammad-hasan-abdillah" || anggota[0].Fraksi != "Fraksi PKS" {
		t.Fatalf("harap tercocok lewat slug, dapat %+v, %v", anggota, err)
	}
	if !slices.Contains(anggota[0].NamaLain, "Muhammad Hasan Abdillah") {
		t.Errorf("harap nama di halaman fraksi masuk nama_lain, dapat %v", anggota[0].NamaLain)
	}
}

func TestGabungkanAnggota_SlugGandaDiHalamanDapil_DicocokkanLewatNama(t *testing.T) {
	dapil := []BarisDapil{
		{NomorDapil: 3, Nama: "Hj. Ida Mahmudah", Slug: "ima-mahdiah"},
		{NomorDapil: 10, Nama: "Ima Mahdiah", Slug: "ima-mahdiah"},
	}
	fraksi := []AnggotaFraksi{
		{Fraksi: "F", BarisStruktur: BarisStruktur{Nama: "Ida Mahmudah", Slug: "hj-ida-mahmudah"}},
		{Fraksi: "F", BarisStruktur: BarisStruktur{Nama: "Ima Mahdiah", Slug: "ima-mahdiah"}},
	}

	anggota, err := GabungkanAnggota(dapil, fraksi)

	if err != nil || anggota[0].ID != "hj-ida-mahmudah" || anggota[1].ID != "ima-mahdiah" {
		t.Fatalf("harap tercocok lewat nama, dapat %+v, %v", anggota, err)
	}
}

func TestGabungkanAnggota_NamaGandaSlugBerbeda_DicocokkanLewatSlug(t *testing.T) {
	dapil := []BarisDapil{{NomorDapil: 1, Nama: "H. Ismail, S.Pd.", Slug: "h-ismail-s-pd"}, {NomorDapil: 5, Nama: "Ismail", Slug: "ismail"}}
	fraksi := []AnggotaFraksi{
		{Fraksi: "F", BarisStruktur: BarisStruktur{Nama: "Ismail", Slug: "ismail"}},
		{Fraksi: "F", BarisStruktur: BarisStruktur{Nama: "H. Ismail, S.Pd., MH", Slug: "h-ismail-s-pd"}},
	}

	anggota, err := GabungkanAnggota(dapil, fraksi)

	if err != nil || anggota[0].ID != "h-ismail-s-pd" || anggota[1].ID != "ismail" {
		t.Fatalf("harap tercocok lewat slug, dapat %+v, %v", anggota, err)
	}
}

func TestGabungkanAnggota_TanpaPasangan_DilaporkanTidakDitebak(t *testing.T) {
	dapil := []BarisDapil{{NomorDapil: 2, Nama: "Nama Asing", Slug: "nama-asing"}}
	fraksi := []AnggotaFraksi{{Fraksi: "F", BarisStruktur: BarisStruktur{Nama: "Nama Lain", Slug: "nama-lain"}}}

	anggota, err := GabungkanAnggota(dapil, fraksi)

	if !errors.Is(err, ErrAnggotaTidakCocok) || anggota != nil {
		t.Errorf("harap ErrAnggotaTidakCocok tanpa hasil, dapat %+v, %v", anggota, err)
	}
}
