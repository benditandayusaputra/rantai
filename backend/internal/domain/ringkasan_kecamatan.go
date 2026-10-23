package domain

import "slices"

type IndikatorKecamatan struct {
	Periode      string
	Kebutuhan    float64
	Kapasitas    float64
	Kesenjangan  float64
	Rasio        *float64
	Kelas        string
	Label        Label
	ProvenanceID string
}

type TemaKecamatan struct {
	Tema      TemaTersimpan
	Indikator *IndikatorKecamatan
}

type RingkasanAnggaran struct {
	Tahun        int
	JumlahPaket  int
	TotalHPS     float64
	AdaSimulasi  bool
	ProvenanceID string
}

type RingkasanKecamatan struct {
	Kecamatan        DetailWilayah
	KelurahanDisorot *RingkasanWilayah
	Tema             []TemaKecamatan
	Anggaran         *RingkasanAnggaran
	Wakil            WakilKecamatan
}

func (r RingkasanAnggaran) Label() Label {
	if r.AdaSimulasi {
		return LabelSimulasi
	}
	return LabelProksi
}

func KomisiPengawasSemuaTema(daftar []TemaTersimpan) []string {
	komisi := []string{}
	for _, t := range daftar {
		komisi = append(komisi, KomisiPengawas(t.Tema)...)
	}
	slices.Sort(komisi)
	return slices.Compact(komisi)
}
