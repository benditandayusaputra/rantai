package domain

import (
	"slices"
	"testing"
)

func TestKeterbatasan_DataPerSekolahTanpaMadrasah_TigaCatatan(t *testing.T) {
	m := masukanUji()
	m.JumlahMadrasah = 2
	peralihan, _ := PeralihanTema(TemaPendidikan)
	detail := DetailPendidikanKecamatan{
		IndikatorTersimpan: IndikatorTersimpan{IndikatorKesenjangan: HitungKesenjanganKursi(m, aturanSMPUji)},
		Peralihan:          peralihan,
	}

	daftar := detail.Keterbatasan()

	harap := []string{KeterbatasanSPMBTidakDibatasiKecamatan, KeterbatasanCalonLulusanTanpaMI, KeterbatasanKursiMadrasahTidakTersedia}
	if len(daftar) != len(harap) || daftar[0] != harap[0] || daftar[1] != harap[1] || daftar[2] != harap[2] {
		t.Errorf("keterbatasan = %v", daftar)
	}
}

func TestKeterbatasan_SMASMK_CatatanLintasKecamatanTanpaMTsDanMA(t *testing.T) {
	m := masukanUji()
	m.JumlahMadrasah = 1
	peralihan, _ := PeralihanTema(TemaPendidikanSMASMK)
	detail := DetailPendidikanKecamatan{
		IndikatorTersimpan: IndikatorTersimpan{IndikatorKesenjangan: HitungKesenjanganKursi(m, aturanSMPUji)},
		Peralihan:          peralihan,
	}

	daftar := detail.Keterbatasan()

	harap := []string{KeterbatasanSPMBTidakDibatasiKecamatan, KeterbatasanSMASMKLintasKecamatan, KeterbatasanCalonLulusanTanpaMTs, KeterbatasanKursiMATidakTersedia}
	if !slices.Equal(daftar, harap) {
		t.Errorf("keterbatasan = %v", daftar)
	}
}

func TestProvenanceKursiJalur_SatuBentukLangsungBanyakBentukKosong(t *testing.T) {
	smp, _ := PeralihanTema(TemaPendidikan)
	smaSMK, _ := PeralihanTema(TemaPendidikanSMASMK)
	provenance := []ProvenanceKursi{
		{Jalur: JalurNegeri, Bentuk: BentukSMP, ProvenanceID: "prv_smp"},
		{Jalur: JalurNegeri, Bentuk: BentukSMA, ProvenanceID: "prv_sma"},
		{Jalur: JalurNegeri, Bentuk: BentukSMK, ProvenanceID: "prv_smk"},
	}

	detailSMP := DetailPendidikanKecamatan{Peralihan: smp, ProvenanceKursi: provenance}
	detailSMA := DetailPendidikanKecamatan{Peralihan: smaSMK, ProvenanceKursi: provenance}

	if detailSMP.ProvenanceKursiJalur(JalurNegeri) != "prv_smp" || detailSMA.ProvenanceKursiJalur(JalurNegeri) != "" {
		t.Errorf("smp %q, sma/smk %q", detailSMP.ProvenanceKursiJalur(JalurNegeri), detailSMA.ProvenanceKursiJalur(JalurNegeri))
	}
	if detailSMA.ProvenanceKursiBentuk(JalurNegeri, BentukSMK) != "prv_smk" {
		t.Errorf("smk = %q", detailSMA.ProvenanceKursiBentuk(JalurNegeri, BentukSMK))
	}
}

func TestTrenLayakTampil_KurangDariDuaPeriode_Kosong(t *testing.T) {
	kasus := []struct {
		nama   string
		titik  []TitikTren
		jumlah int
	}{
		{"tanpa periode", nil, 0},
		{"satu periode", []TitikTren{{Periode: "2026/2027", Rasio: 1.4}}, 0},
		{"dua periode", []TitikTren{{Periode: "2025/2026", Rasio: 1.2}, {Periode: "2026/2027", Rasio: 1.4}}, 2},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			hasil := TrenLayakTampil(k.titik)
			if hasil == nil || len(hasil) != k.jumlah {
				t.Errorf("tren = %v, harap %d titik", hasil, k.jumlah)
			}
		})
	}
}
