package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var ErrBelumAdaCalonLulusan = errors.New("belum ada calon lulusan, jalankan job calon_lulusan lebih dulu")

type PembacaMasukanKesenjangan interface {
	MasukanKesenjangan(ctx context.Context, peralihan domain.Peralihan) (string, []domain.MasukanKesenjangan, error)
}

type KesenjanganKursi struct {
	Peralihan      domain.Peralihan
	Pembaca        PembacaMasukanKesenjangan
	Konstanta      KonstantaPendidikan
	BuatProvenance PembuatProvenance
}

func (p KesenjanganKursi) Hitung(ctx context.Context, tema domain.Tema) (DataIndikator, error) {
	periode, daftarMasukan, err := p.Pembaca.MasukanKesenjangan(ctx, p.Peralihan)
	if err != nil {
		return DataIndikator{}, fmt.Errorf("baca masukan kesenjangan %s: %w", p.Peralihan.Kode, err)
	}
	if periode == "" {
		return DataIndikator{}, ErrBelumAdaCalonLulusan
	}
	aturan := domain.AturanKesenjangan{
		BentukTujuan:        p.Peralihan.BentukTujuan,
		BatasMuridPerRombel: p.Konstanta.Untuk(p.Peralihan).BatasMuridPerRombel,
		Ambang:              tema.Indikator.AmbangRasio,
	}
	data := DataIndikator{Daftar: []IndikatorTercatat{}, Provenance: []domain.Provenance{}}
	provenancePerLabel := map[domain.Label]domain.Provenance{}
	for _, masukan := range daftarMasukan {
		indikator := domain.HitungKesenjanganKursi(masukan, aturan)
		indikator.TemaKode, indikator.Periode = tema.Kode, periode
		prv, sudahAda := provenancePerLabel[indikator.Label]
		if !sudahAda {
			prv = p.BuatProvenance(p.catatan(indikator.Label, periode, aturan.BatasMuridPerRombel))
			provenancePerLabel[indikator.Label] = prv
			data.Provenance = append(data.Provenance, prv)
		}
		data.Daftar = append(data.Daftar, IndikatorTercatat{IndikatorKesenjangan: indikator, ProvenanceID: prv.ID})
	}
	return data, nil
}

func (p KesenjanganKursi) catatan(label domain.Label, periode string, batasMuridPerRombel int) CatatanProvenance {
	asal, kursi := p.Peralihan.BentukAsal, p.rumusKursiNegeri()
	rumus := fmt.Sprintf("Periode %s. kesenjangan = calon lulusan %s − %s (Kepkadis 258/2026 sebagaimana diubah Kepkadis 315/2026); "+
		"rasio = %s ÷ calon lulusan %s; kursi madrasah (proksi) = rombel kelas %d %s × %d. "+
		"Label mengikuti label terlemah dari masukan.",
		periode, asal, kursi, kursi, asal, p.Peralihan.TingkatTujuanPertama(), p.Peralihan.BentukMadrasahTujuan, batasMuridPerRombel)
	return CatatanProvenance{Sumber: domain.SumberRantai, URL: p.Konstanta.URLSumberSiswa, Label: label, Catatan: rumus}
}

func (p KesenjanganKursi) rumusKursiNegeri() string {
	suku := make([]string, len(p.Peralihan.BentukTujuan))
	for i, bentuk := range p.Peralihan.BentukTujuan {
		suku[i] = fmt.Sprintf("kursi %s negeri", bentuk)
	}
	if len(suku) == 1 {
		return suku[0]
	}
	return "(" + strings.Join(suku, " + ") + ")"
}
