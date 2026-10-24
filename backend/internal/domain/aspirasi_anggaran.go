package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type StatusPaketTerkait string

const (
	StatusKemungkinanTerkait StatusPaketTerkait = "kemungkinan_terkait"
	StatusBelumDitemukan     StatusPaketTerkait = "belum_ditemukan"
	StatusTanpaWilayah       StatusPaketTerkait = "tanpa_wilayah"
)

var (
	ErrKataKunciAnggaranTidakSah = errors.New("kata kunci bidang wajib berhuruf kecil dan tidak kosong")
	ErrBatasPaketTerkaitTidakSah = errors.New("batas_paket wajib lebih dari nol")
)

type AturanBidangAnggaran struct {
	Kategori  []string `yaml:"kategori"`
	KataKunci []string `yaml:"kata_kunci"`
}

type KonfigurasiAspirasiAnggaran struct {
	Bidang     map[string]AturanBidangAnggaran `yaml:"bidang"`
	BatasPaket int                             `yaml:"batas_paket"`
}

type KueriPaketTerkait struct {
	KodeWilayah   string
	TahunAnggaran int
	Kategori      []string
	PolaKataKunci []string
	Batas         int
}

type PaketTerkait struct {
	PaketRingkas
	IsCocokKategori bool
}

type HasilPaketTerkait struct {
	Aspirasi      AspirasiReses
	TahunAnggaran int
	Paket         []PaketTerkait
}

func (k KonfigurasiAspirasiAnggaran) Validasi() error {
	if k.BatasPaket <= 0 {
		return ErrBatasPaketTerkaitTidakSah
	}
	for bidang, aturan := range k.Bidang {
		for _, kata := range aturan.KataKunci {
			if strings.TrimSpace(kata) == "" || kata != strings.ToLower(kata) {
				return fmt.Errorf("bidang %s kata %q: %w", bidang, kata, ErrKataKunciAnggaranTidakSah)
			}
		}
	}
	return nil
}

func (k KonfigurasiAspirasiAnggaran) KueriUntuk(a AspirasiReses) KueriPaketTerkait {
	aturan := k.Bidang[a.Bidang]
	pola := make([]string, len(aturan.KataKunci))
	for i, kata := range aturan.KataKunci {
		pola[i] = `\m` + regexp.QuoteMeta(kata) + `\M`
	}
	kategori := aturan.Kategori
	if kategori == nil {
		kategori = []string{}
	}
	return KueriPaketTerkait{
		KodeWilayah: a.KodeWilayah, TahunAnggaran: a.TanggalParipurna.Year() + 1,
		Kategori: kategori, PolaKataKunci: pola, Batas: k.BatasPaket,
	}
}

func (h HasilPaketTerkait) Status() StatusPaketTerkait {
	switch {
	case h.Aspirasi.KodeWilayah == "":
		return StatusTanpaWilayah
	case len(h.Paket) == 0:
		return StatusBelumDitemukan
	default:
		return StatusKemungkinanTerkait
	}
}
