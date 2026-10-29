package domain

import "testing"

func indikatorKotaUji(kebutuhan, kapasitas int) IndikatorTersimpan {
	rasio := hitungRasio(kapasitas, kebutuhan)
	return IndikatorTersimpan{
		IndikatorKesenjangan: IndikatorKesenjangan{
			Kebutuhan: kebutuhan, Kapasitas: kapasitas, Rasio: rasio, Label: LabelResmi, Kelas: KelasCukup,
		},
		ProvenanceID: "sumber-uji",
	}
}

func TestRingkasKapasitasKota_MenggunakanRasioTotalBukanRataRata(t *testing.T) {
	kecamatan := []RingkasanWilayah{{Kode: "31.72.01"}, {Kode: "31.72.02"}}
	indikator := []IndikatorTersimpan{indikatorKotaUji(100, 100), indikatorKotaUji(900, 0)}
	indikator[1].Label = LabelEstimasi
	r := RingkasKapasitasKota(kecamatan, indikator)
	if r.Rasio == nil || *r.Rasio != 0.1 || *r.Kebutuhan != 1000 || *r.Kapasitas != 100 || *r.Kesenjangan != 900 {
		t.Fatalf("ringkasan salah: %+v", r)
	}
	if r.Label != LabelEstimasi || r.KecamatanKurang != 1 || r.KecamatanDenganData != 2 {
		t.Fatalf("label dan cakupan salah: %+v", r)
	}
}

func TestRingkasKapasitasKota_DataHilangTidakDianggapNol(t *testing.T) {
	kecamatan := []RingkasanWilayah{{Kode: "31.72.01"}, {Kode: "31.72.02"}}
	for _, indikator := range [][]IndikatorTersimpan{
		{}, {indikatorKotaUji(100, 100)}, {indikatorKotaUji(100, 100), indikatorKotaUji(0, 50)},
	} {
		r := RingkasKapasitasKota(kecamatan, indikator)
		if r.Rasio != nil || r.Kebutuhan != nil || r.Kapasitas != nil || r.Kesenjangan != nil {
			t.Fatalf("total parsial harus kosong: %+v", r)
		}
	}
}

func TestRingkasKapasitasKota_SurplusTotalTidakMenghapusKekuranganLokal(t *testing.T) {
	kecamatan := []RingkasanWilayah{{Kode: "31.72.01"}, {Kode: "31.72.02"}}
	r := RingkasKapasitasKota(kecamatan, []IndikatorTersimpan{indikatorKotaUji(100, 250), indikatorKotaUji(100, 50)})
	if r.KecamatanKurang != 1 || r.Kesenjangan == nil || *r.Kesenjangan != -100 || *r.Rasio != 1.5 {
		t.Fatalf("surplus dan kekurangan lokal: %+v", r)
	}
}
