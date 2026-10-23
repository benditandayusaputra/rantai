package domain

import (
	"cmp"
	"slices"
)

type RingkasanDapil struct {
	ID                int
	Nomor             int
	Nama              string
	TahunPemilu       int
	JumlahKursi       int
	ProvenanceID      string
	ProvenanceAnggota string
}

type Wakil struct {
	ID                     string
	Nama                   string
	Partai                 string
	Fraksi                 string
	KodeKomisi             string
	NamaKomisi             string
	IsDiKomisiPengawasTema bool
}

type DetailDapil struct {
	RingkasanDapil
	Kecamatan []RingkasanWilayah
	Anggota   []Wakil
}

type WakilKecamatan struct {
	Kecamatan      RingkasanWilayah
	KomisiPengawas []string
	DetailDapil
}

type ProfilAnggota struct {
	Wakil
	BidangKomisi    []string
	Dapil           DetailDapil
	AlatKelengkapan []JabatanAKD
	URLProfil       string
	URLLHKPN        string
	ProvenanceID    string
}

func UrutkanWakil(daftar []Wakil, komisiPengawas []string) []Wakil {
	hasil := make([]Wakil, len(daftar))
	for i, w := range daftar {
		w.IsDiKomisiPengawasTema = w.KodeKomisi != "" && slices.Contains(komisiPengawas, w.KodeKomisi)
		hasil[i] = w
	}
	slices.SortStableFunc(hasil, func(a, b Wakil) int {
		if a.IsDiKomisiPengawasTema != b.IsDiKomisiPengawasTema {
			return cmp.Compare(urutanPengawas(a), urutanPengawas(b))
		}
		return cmp.Compare(KunciNama(a.Nama), KunciNama(b.Nama))
	})
	return hasil
}

func urutanPengawas(w Wakil) int {
	if w.IsDiKomisiPengawasTema {
		return 0
	}
	return 1
}

func KomisiPengawas(tema Tema) []string {
	komisi := []string{}
	for _, tj := range tema.TanggungJawab {
		komisi = append(komisi, tj.Komisi)
	}
	slices.Sort(komisi)
	return slices.Compact(komisi)
}

const (
	KeteranganKomisiPimpinan      = "pimpinan_dprd"
	KeteranganKomisiTidakTersedia = "tidak_tersedia"
)

func (p ProfilAnggota) KeteranganKomisi() string {
	switch {
	case p.KodeKomisi != "":
		return ""
	case slices.ContainsFunc(p.AlatKelengkapan, func(j JabatanAKD) bool { return j.AlatKelengkapan == AlatKelengkapanPimpinan }):
		return KeteranganKomisiPimpinan
	default:
		return KeteranganKomisiTidakTersedia
	}
}
