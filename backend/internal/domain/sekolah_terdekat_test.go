package domain

import "testing"

func sekolahDekatUji(npsn string, bentuk BentukSekolah, status StatusSekolah, jalur JalurDayaTampung, kursi *int) SekolahTerdekat {
	return SekolahTerdekat{SekolahDenganKursi: SekolahDenganKursi{
		NPSN: npsn, Bentuk: bentuk, Status: status, Jalur: jalur, Kursi: kursi, ProvenanceKursi: "prv_" + npsn,
	}}
}

func TestKelompokkan_NegeriSPMBMadrasah_SwastaLainDiabaikan(t *testing.T) {
	kursi := 100
	sekitar := SekitarTitik{Sekolah: []SekolahTerdekat{
		sekolahDekatUji("1", BentukSMP, StatusNegeri, JalurNegeri, &kursi),
		sekolahDekatUji("2", BentukSMP, StatusSwasta, JalurSPMBBersama, &kursi),
		sekolahDekatUji("3", BentukMTs, StatusNegeri, "", nil),
		sekolahDekatUji("4", BentukSMP, StatusSwasta, "", nil),
	}}

	k := sekitar.Kelompokkan()

	if len(k.SMPNegeri) != 1 || len(k.SPMBBersama) != 1 || len(k.Madrasah) != 1 || k.Madrasah[0].NPSN != "3" {
		t.Errorf("kelompok = %+v", k)
	}
}

func TestKursiNegeriDalamRadius_SekolahTanpaKursiDilewati(t *testing.T) {
	a, b := 288, 252
	daftar := []SekolahTerdekat{
		sekolahDekatUji("1", BentukSMP, StatusNegeri, JalurNegeri, &a),
		sekolahDekatUji("2", BentukSMP, StatusNegeri, "", nil),
		sekolahDekatUji("3", BentukSMP, StatusNegeri, JalurNegeri, &b),
	}

	kursi, provenanceID := KursiNegeriDalamRadius(daftar)

	if kursi != 540 || provenanceID != "prv_1" {
		t.Errorf("kursi %d, provenance %q", kursi, provenanceID)
	}
}

func TestTekananRadius_TanpaKursiNegeri_TidakDidefinisikan(t *testing.T) {
	if TekananRadius(500, 0) != nil {
		t.Error("tekanan tanpa kursi harus kosong")
	}
	if tekanan := TekananRadius(500, 250); tekanan == nil || *tekanan != 2 {
		t.Errorf("tekanan = %v", tekanan)
	}
}
