package repository

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	sqlIndikatorTerbaru = `SELECT periode, kebutuhan, kapasitas, kesenjangan, rasio, kelas, label, provenance_id
FROM indikator_kecamatan
WHERE tema_kode = $1 AND kode_kecamatan = $2
ORDER BY periode DESC
LIMIT 1`

	sqlRingkasanAnggaran = `SELECT count(*), coalesce(sum(p.hps), 0), coalesce(bool_or(p.is_simulation), false),
    coalesce((array_agg(p.provenance_id ORDER BY pr.diambil_pada DESC))[1], '')
FROM paket p
JOIN provenance pr ON pr.id = p.provenance_id
WHERE p.tahun_anggaran = $2 AND NOT p.is_batal
    AND p.lokasi_tingkat IN ('kecamatan', 'kelurahan')
    AND (p.kode_wilayah = $1 OR p.kode_wilayah LIKE $1 || '.%')`
)

type RingkasanKecamatanRepository struct {
	DB Basisdata
}

func (r RingkasanKecamatanRepository) IndikatorTerbaru(ctx context.Context, kodeTema, kodeKecamatan string) (domain.IndikatorKecamatan, bool, error) {
	var i domain.IndikatorKecamatan
	err := r.DB.QueryRow(ctx, sqlIndikatorTerbaru, kodeTema, kodeKecamatan).Scan(
		&i.Periode, &i.Kebutuhan, &i.Kapasitas, &i.Kesenjangan, &i.Rasio, &i.Kelas, &i.Label, &i.ProvenanceID,
	)
	return hasilSatuBaris(i, err)
}

func (r RingkasanKecamatanRepository) RingkasanAnggaran(ctx context.Context, kodeKecamatan string, tahun int) (domain.RingkasanAnggaran, bool, error) {
	a := domain.RingkasanAnggaran{Tahun: tahun}
	err := r.DB.QueryRow(ctx, sqlRingkasanAnggaran, kodeKecamatan, tahun).Scan(&a.JumlahPaket, &a.TotalHPS, &a.AdaSimulasi, &a.ProvenanceID)
	if err != nil {
		return domain.RingkasanAnggaran{}, false, err
	}
	return a, a.JumlahPaket > 0, nil
}
