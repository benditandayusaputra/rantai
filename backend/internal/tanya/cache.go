package tanya

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var ErrTidakTercache = errors.New("pertanyaan tidak ada di cache demo dan model tidak dipasang")

type EntriCache struct {
	Pertanyaan           string    `json:"pertanyaan"`
	KodeKecamatanKonteks string    `json:"kode_kecamatan_konteks"`
	Alat                 []string  `json:"alat"`
	Teks                 string    `json:"teks"`
	Rujukan              []Rujukan `json:"rujukan"`
}

type PenjawabTercache struct {
	Cache    []EntriCache
	Cadangan Penjawab
}

func (p PenjawabTercache) Periksa(m Permintaan) error {
	return p.Cadangan.Periksa(m)
}

func (p PenjawabTercache) Jawab(ctx context.Context, m Permintaan, kabar func(Kabar)) (Jawaban, error) {
	if entri, ada := p.cari(m); ada {
		for _, alat := range entri.Alat {
			kabar(Kabar{Alat: alat})
		}
		return Jawaban{Teks: entri.Teks, Rujukan: entri.Rujukan, HasilAlat: []json.RawMessage{}}, nil
	}
	if p.Cadangan.Model == nil {
		return Jawaban{}, ErrTidakTercache
	}
	return p.Cadangan.Jawab(ctx, m, kabar)
}

func (p PenjawabTercache) cari(m Permintaan) (EntriCache, bool) {
	pertanyaan := domain.NormalkanTeks(m.Pertanyaan)
	konteks := kecamatanKonteks(m.KodeKecamatanKonteks)
	for _, entri := range p.Cache {
		if domain.NormalkanTeks(entri.Pertanyaan) == pertanyaan && (konteks == "" || konteks == entri.KodeKecamatanKonteks) {
			return entri, true
		}
	}
	return EntriCache{}, false
}

func kecamatanKonteks(kode string) string {
	if tingkat, err := domain.TingkatDariKode(kode); err == nil && tingkat == domain.TingkatKelurahan {
		return domain.IndukKode(kode)
	}
	return kode
}
