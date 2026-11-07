package domain

import "strconv"

type JenisSimpul string

const (
	SimpulInstansi JenisSimpul = "instansi"
	SimpulPaket    JenisSimpul = "paket"
	SimpulPenyedia JenisSimpul = "penyedia"
)

const BatasPaketJaringan = 60

var ErrFilterJaringanTidakSah = InputTidakValid("isi salah satu dari wilayah atau penyedia")

type FilterJaringan struct {
	KodeWilayah string
	PenyediaID  int64
	Tahun       int
}

type RelasiPaket struct {
	PaketKode    string
	PaketNama    string
	Nilai        float64
	ProvenanceID string
	InstansiID   int64
	InstansiNama string
	PenyediaID   *int64
	PenyediaNama string
}

type SimpulJaringan struct {
	ID           string
	Jenis        JenisSimpul
	Kunci        string
	Label        string
	Nilai        float64
	JumlahPaket  int
	ProvenanceID string
}

type SisiJaringan struct {
	Dari string
	Ke   string
}

type JaringanRelasi struct {
	Simpul      []SimpulJaringan
	Sisi        []SisiJaringan
	IsTerpotong bool
}

func (f FilterJaringan) Validasi() error {
	adaWilayah := f.KodeWilayah != ""
	adaPenyedia := f.PenyediaID > 0
	if adaWilayah == adaPenyedia {
		return ErrFilterJaringanTidakSah
	}
	return nil
}

type penyusunJaringan struct {
	urutan []string
	simpul map[string]*SimpulJaringan
	sisi   []SisiJaringan
}

func SusunJaringan(relasi []RelasiPaket) JaringanRelasi {
	p := penyusunJaringan{urutan: []string{}, simpul: map[string]*SimpulJaringan{}, sisi: []SisiJaringan{}}
	for _, r := range relasi {
		p.tambahRelasi(r)
	}
	return p.hasil()
}

func (p *penyusunJaringan) tambahRelasi(r RelasiPaket) {
	paket := p.catat(SimpulJaringan{Jenis: SimpulPaket, Kunci: r.PaketKode, Label: r.PaketNama}, r)
	instansi := p.catat(SimpulJaringan{Jenis: SimpulInstansi, Kunci: strconv.FormatInt(r.InstansiID, 10), Label: r.InstansiNama}, r)
	p.sisi = append(p.sisi, SisiJaringan{Dari: instansi, Ke: paket})
	if r.PenyediaID == nil {
		return
	}
	penyedia := p.catat(SimpulJaringan{Jenis: SimpulPenyedia, Kunci: strconv.FormatInt(*r.PenyediaID, 10), Label: r.PenyediaNama}, r)
	p.sisi = append(p.sisi, SisiJaringan{Dari: paket, Ke: penyedia})
}

func (p *penyusunJaringan) catat(baru SimpulJaringan, r RelasiPaket) string {
	id := string(baru.Jenis) + ":" + baru.Kunci
	simpul, ada := p.simpul[id]
	if !ada {
		baru.ID = id
		baru.ProvenanceID = r.ProvenanceID
		simpul = &baru
		p.simpul[id] = simpul
		p.urutan = append(p.urutan, id)
	}
	simpul.Nilai += r.Nilai
	simpul.JumlahPaket++
	return id
}

func (p *penyusunJaringan) hasil() JaringanRelasi {
	simpul := make([]SimpulJaringan, len(p.urutan))
	for i, id := range p.urutan {
		simpul[i] = *p.simpul[id]
	}
	return JaringanRelasi{Simpul: simpul, Sisi: p.sisi}
}
