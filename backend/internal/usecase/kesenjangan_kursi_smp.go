package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var ErrBelumAdaCalonLulusan = errors.New("belum ada calon lulusan, jalankan job calon_lulusan lebih dulu")

type PembacaMasukanKesenjangan interface {
	MasukanKesenjanganSMP(ctx context.Context) (string, []domain.MasukanKesenjanganSMP, error)
}

type KesenjanganKursiSMP struct {
	Pembaca        PembacaMasukanKesenjangan
	Konstanta      KonstantaPendidikan
	BuatProvenance PembuatProvenance
}

func (p KesenjanganKursiSMP) Hitung(ctx context.Context, tema domain.Tema) (DataIndikator, error) {
	periode, daftarMasukan, err := p.Pembaca.MasukanKesenjanganSMP(ctx)
	if err != nil {
		return DataIndikator{}, fmt.Errorf("baca masukan kesenjangan: %w", err)
	}
	if periode == "" {
		return DataIndikator{}, ErrBelumAdaCalonLulusan
	}
	data := DataIndikator{Daftar: []IndikatorTercatat{}, Provenance: []domain.Provenance{}}
	provenancePerLabel := map[domain.Label]domain.Provenance{}
	for _, masukan := range daftarMasukan {
		masukan.BatasMuridPerRombelSMP = p.Konstanta.BatasMuridPerRombelSMP
		indikator := domain.HitungKesenjanganKursiSMP(masukan, tema.Indikator.AmbangRasio)
		indikator.TemaKode, indikator.Periode = tema.Kode, periode
		prv, sudahAda := provenancePerLabel[indikator.Label]
		if !sudahAda {
			prv = p.BuatProvenance(p.catatan(indikator.Label, periode))
			provenancePerLabel[indikator.Label] = prv
			data.Provenance = append(data.Provenance, prv)
		}
		data.Daftar = append(data.Daftar, IndikatorTercatat{IndikatorKesenjanganSMP: indikator, ProvenanceID: prv.ID})
	}
	return data, nil
}

func (p KesenjanganKursiSMP) catatan(label domain.Label, periode string) CatatanProvenance {
	rumus := fmt.Sprintf("Periode %s. kesenjangan = calon lulusan SD − kursi SMP negeri (Kepkadis 258/2026); "+
		"rasio = kursi SMP negeri ÷ calon lulusan SD; kursi madrasah (proksi) = rombel kelas 7 MTs × %d. "+
		"Label mengikuti label terlemah dari masukan.", periode, p.Konstanta.BatasMuridPerRombelSMP)
	return CatatanProvenance{Sumber: domain.SumberRantai, URL: p.Konstanta.URLSumberSiswa, Label: label, Catatan: rumus}
}
