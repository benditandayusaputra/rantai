package domain

import "slices"

const SumberDJPK SumberData = "djpk"

const BulanPenuhAPBD = 12

var ErrAPBDTidakDitemukan = TidakDitemukan("data APBD untuk tahun ini belum tersedia")

type AkunAPBD struct {
	Nama      string
	Anggaran  float64
	Realisasi float64
}

type APBDDaerah struct {
	Tahun          int
	BulanRealisasi int
	Akun           []AkunAPBD
	ProvenanceID   string
}

type AturanAPBD struct {
	AkunBelanja  string   `yaml:"akun_belanja"`
	JenisBelanja []string `yaml:"jenis_belanja"`
}

type RingkasanBelanjaAPBD struct {
	Tahun          int
	BulanRealisasi int
	Belanja        AkunAPBD
	Jenis          []AkunAPBD
	ProvenanceID   string
}

func RapikanAkunAPBD(akun []AkunAPBD) []AkunAPBD {
	return slices.Compact(slices.Clone(akun))
}

func IsRealisasiAPBDSama(a, b []AkunAPBD) bool {
	return slices.EqualFunc(a, b, func(x, y AkunAPBD) bool {
		return x.Nama == y.Nama && x.Realisasi == y.Realisasi
	})
}

func (a APBDDaerah) RingkasBelanja(aturan AturanAPBD) (RingkasanBelanjaAPBD, error) {
	belanja, ada := a.cariAkun(aturan.AkunBelanja)
	if !ada {
		return RingkasanBelanjaAPBD{}, ErrAPBDTidakDitemukan
	}
	jenis := make([]AkunAPBD, 0, len(aturan.JenisBelanja))
	for _, nama := range aturan.JenisBelanja {
		if akun, ada := a.cariAkun(nama); ada {
			jenis = append(jenis, akun)
		}
	}
	return RingkasanBelanjaAPBD{
		Tahun: a.Tahun, BulanRealisasi: a.BulanRealisasi, Belanja: belanja, Jenis: jenis, ProvenanceID: a.ProvenanceID,
	}, nil
}

func (a APBDDaerah) cariAkun(nama string) (AkunAPBD, bool) {
	indeks := slices.IndexFunc(a.Akun, func(akun AkunAPBD) bool { return akun.Nama == nama })
	if indeks < 0 {
		return AkunAPBD{}, false
	}
	return a.Akun[indeks], true
}
