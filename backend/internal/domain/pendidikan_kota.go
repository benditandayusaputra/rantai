package domain

import "time"

type JumlahSekolahKecamatan struct {
	KodeKecamatan string
	Bentuk        BentukSekolah
	Negeri        int
	Swasta        int
	Label         Label
	ProvenanceID  string
	DiambilPada   time.Time
}

type PendidikanKota struct {
	KabupatenKota []RingkasanWilayah
	Wilayah       RingkasanWilayah
	Tema          string
	Periode       string
	Kecamatan     []RingkasanWilayah
	Indikator     []IndikatorTersimpan
	Sekolah       []JumlahSekolahKecamatan
	Ringkasan     RingkasanKapasitasKota
}

type RingkasanKapasitasKota struct {
	KecamatanDenganData int
	KecamatanKurang     int
	Kebutuhan           *int
	Kapasitas           *int
	Kesenjangan         *int
	Rasio               *float64
	Label               Label
	ProvenanceID        string
}

func RingkasKapasitasKota(kecamatan []RingkasanWilayah, indikator []IndikatorTersimpan) RingkasanKapasitasKota {
	r := RingkasanKapasitasKota{Label: LabelResmi}
	kebutuhan, kapasitas := 0, 0
	for _, i := range indikator {
		if i.Rasio == nil || i.Kebutuhan <= 0 || i.Kelas == KelasTanpaData {
			continue
		}
		r.KecamatanDenganData++
		kebutuhan, kapasitas = kebutuhan+i.Kebutuhan, kapasitas+i.Kapasitas
		r.Label = LabelTerlemah(r.Label, i.Label)
		if r.ProvenanceID == "" {
			r.ProvenanceID = i.ProvenanceID
		}
		if i.Kapasitas < i.Kebutuhan {
			r.KecamatanKurang++
		}
	}
	if len(kecamatan) > 0 && r.KecamatanDenganData == len(kecamatan) {
		selisih := kebutuhan - kapasitas
		r.Kebutuhan, r.Kapasitas, r.Kesenjangan = &kebutuhan, &kapasitas, &selisih
		r.Rasio = hitungRasio(kapasitas, kebutuhan)
	}
	return r
}
