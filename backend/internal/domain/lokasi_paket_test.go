package domain

import "testing"

const (
	kodeJakut           = "31.72"
	kodeJaktim          = "31.75"
	kodeKecKelapaGading = "31.72.06"
	kodeKelKGBarat      = "31.72.06.1001"
	kodeKelKGTimur      = "31.72.06.1002"
	kodeKelPegangsaan   = "31.72.06.1003"
	kodeKecCakung       = "31.75.06"
	kodeKelUjungMenteng = "31.75.06.1005"
	kodeKecMenteng      = "31.71.07"
	kodeKelMenteng      = "31.71.07.1001"
)

func penentuUji() PenentuLokasi {
	return PenentuLokasiBaru(
		AturanLokasi{
			Singkatan:           map[string]string{"kel.": "kelurahan", "kec.": "kecamatan", "jl.": "jalan"},
			PergeseranMaksMeter: 150,
		},
		[]EntriGazetir{
			{Kode: kodeKecKelapaGading, Tingkat: TingkatKecamatan, Nama: []string{"Kelapa Gading", "Kec. Klp. Gading", "Kecamatan Kelapa Gading"}, KodeKota: kodeJakut},
			{Kode: kodeKelKGBarat, Tingkat: TingkatKelurahan, Nama: []string{"kelapa gading barat"}, KodeKota: kodeJakut, KodeKecamatan: kodeKecKelapaGading, TitikTengah: Titik{Lat: -6.16, Lng: 106.9}},
			{Kode: kodeKelKGTimur, Tingkat: TingkatKelurahan, Nama: []string{"kelapa gading timur"}, KodeKota: kodeJakut, KodeKecamatan: kodeKecKelapaGading},
			{Kode: kodeKelPegangsaan, Tingkat: TingkatKelurahan, Nama: []string{"pegangsaan dua"}, KodeKota: kodeJakut, KodeKecamatan: kodeKecKelapaGading},
			{Kode: kodeKecCakung, Tingkat: TingkatKecamatan, Nama: []string{"cakung"}, KodeKota: kodeJaktim},
			{Kode: "31.74.08.1001", Tingkat: TingkatKelurahan, Nama: []string{"Balekambang"}, KodeKota: "31.75", KodeKecamatan: kodeKecCakung},
			{Kode: kodeKelUjungMenteng, Tingkat: TingkatKelurahan, Nama: []string{"ujung menteng"}, KodeKota: kodeJaktim, KodeKecamatan: kodeKecCakung},
			{Kode: kodeKecMenteng, Tingkat: TingkatKecamatan, Nama: []string{"menteng"}, KodeKota: "31.71"},
			{Kode: kodeKelMenteng, Tingkat: TingkatKelurahan, Nama: []string{"menteng"}, KodeKota: "31.71", KodeKecamatan: kodeKecMenteng},
		},
	)
}

func tentukanUji(nama, kodeKota string) LokasiPaket {
	return penentuUji().Tentukan(MasukanLokasi{KodePaket: "10403722000", NamaPaket: nama, KodeKotaInstansi: kodeKota})
}

func TestTentukan_ContohWajib_KelurahanKelapaGadingBaratRT4RW22(t *testing.T) {
	nama := "Pembangunan Pintu Air Gg. Nusa Indah RT 04 RW 22 Gading Sengon (JPO Kali Betik), Kel. Kelapa Gading Barat, Kec. Kelapa Gading"

	lokasi := tentukanUji(nama, kodeJakut)

	if lokasi.Tingkat != TingkatKelurahan || lokasi.KodeWilayah != kodeKelKGBarat {
		t.Fatalf("lokasi = %+v", lokasi)
	}
	if lokasi.RT == nil || *lokasi.RT != 4 || lokasi.RW == nil || *lokasi.RW != 22 {
		t.Errorf("rt %v rw %v", lokasi.RT, lokasi.RW)
	}
}

func TestTentukan_DuaKelurahanBerbeda_TingkatKecamatanAtauKota(t *testing.T) {
	kasus := []struct {
		nama  string
		harap Tingkat
	}{
		{"Rehab Saluran Kel. Kelapa Gading Barat dan Kel. Pegangsaan Dua", TingkatKecamatan},
		{"Rehab Saluran Kel. Kelapa Gading Barat dan Kel. Ujung Menteng", TingkatKota},
	}
	for _, k := range kasus {
		if lokasi := tentukanUji(k.nama, ""); lokasi.Tingkat != k.harap && lokasi.Tingkat != TingkatProvinsi {
			t.Errorf("%q = %+v, harap %s", k.nama, lokasi, k.harap)
		}
	}
	if lokasi := tentukanUji(kasus[1].nama, kodeJakut); lokasi.Tingkat != TingkatKelurahan || lokasi.KodeWilayah != kodeKelKGBarat {
		t.Errorf("saring kota instansi gagal: %+v", lokasi)
	}
}

func TestTentukan_TanpaWilayah_TingkatKotaInstansiAtauProvinsi(t *testing.T) {
	if lokasi := tentukanUji("Pengadaan Alat Tulis Kantor", kodeJaktim); lokasi.Tingkat != TingkatKota || lokasi.KodeWilayah != kodeJaktim {
		t.Errorf("sudin = %+v", lokasi)
	}
	if lokasi := tentukanUji("Pengadaan Alat Tulis Kantor", ""); lokasi.Tingkat != TingkatProvinsi || lokasi.KodeWilayah != KodeProvinsiDKI {
		t.Errorf("dinas provinsi = %+v", lokasi)
	}
}

func TestTentukan_KecamatanSaja_TingkatKecamatan(t *testing.T) {
	lokasi := tentukanUji("Perbaikan Fasilitas Olahraga RW 9 Kec. Cakung", kodeJaktim)

	if lokasi.Tingkat != TingkatKecamatan || lokasi.KodeWilayah != kodeKecCakung || lokasi.Titik != nil || *lokasi.RW != 9 {
		t.Errorf("lokasi = %+v", lokasi)
	}
}

func TestTentukan_KelurahanTidakKonsistenDenganKecamatan_TingkatKota(t *testing.T) {
	lokasi := tentukanUji("Perbaikan Saluran Kel. Ujung Menteng Kec. Kelapa Gading", "")

	if lokasi.Tingkat != TingkatProvinsi {
		t.Errorf("lokasi = %+v", lokasi)
	}
}

func TestTentukan_NamaSamaUntukKelurahanDanKecamatan_AwalanMenentukan(t *testing.T) {
	if lokasi := tentukanUji("Rehab Taman Kelurahan Menteng", ""); lokasi.Tingkat != TingkatKelurahan || lokasi.KodeWilayah != kodeKelMenteng {
		t.Errorf("kelurahan = %+v", lokasi)
	}
	if lokasi := tentukanUji("Rehab Taman Kecamatan Menteng", ""); lokasi.Tingkat != TingkatKecamatan || lokasi.KodeWilayah != kodeKecMenteng {
		t.Errorf("kecamatan = %+v", lokasi)
	}
}

func TestTentukan_AliasSingkatan_KecamatanCocok(t *testing.T) {
	lokasi := tentukanUji("Rehab Saluran di Klp. Gading", kodeJakut)

	if lokasi.Tingkat != TingkatKecamatan || lokasi.KodeWilayah != kodeKecKelapaGading {
		t.Errorf("lokasi = %+v", lokasi)
	}
}

func TestTentukan_RTRWGarisMiring_KeduanyaTerbaca(t *testing.T) {
	lokasi := tentukanUji("Perbaikan Lapangan Voli Jl. Tanjung Duren Selatan I RT.007/002 Kel. Ujung Menteng", kodeJaktim)

	if lokasi.RT == nil || *lokasi.RT != 7 || lokasi.RW == nil || *lokasi.RW != 2 {
		t.Errorf("rt %v rw %v", lokasi.RT, lokasi.RW)
	}
}

func TestTentukan_KodePaketSama_PergeseranDeterministikDanKecil(t *testing.T) {
	nama := "Pembangunan Pintu Air Kel. Kelapa Gading Barat"
	pertama := tentukanUji(nama, kodeJakut)
	kedua := tentukanUji(nama, kodeJakut)
	lain := penentuUji().Tentukan(MasukanLokasi{KodePaket: "kode-lain", NamaPaket: nama})

	if *pertama.Titik != *kedua.Titik || *pertama.Titik == *lain.Titik {
		t.Errorf("titik %v %v %v", *pertama.Titik, *kedua.Titik, *lain.Titik)
	}
	selisih := (pertama.Titik.Lat - (-6.16)) * meterPerDerajat
	if selisih > 150 || selisih < -150 {
		t.Errorf("pergeseran %v meter melebihi batas", selisih)
	}
}

func TestTentukan_EjaanBerspasiBerbeda_TetapCocok(t *testing.T) {
	lokasi := tentukanUji("Pembangunan Puskesmas di Kelurahan Bale Kambang", "")

	if lokasi.Tingkat != TingkatKelurahan || lokasi.KodeWilayah != "31.74.08.1001" {
		t.Errorf("lokasi = %+v", lokasi)
	}
}
