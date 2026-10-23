package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlAgregatPaket = `count(*) AS jumlah, coalesce(sum(p.hps), 0) AS total, bool_or(p.is_simulation) AS simulasi,
    (array_agg(p.provenance_id ORDER BY pr.diambil_pada DESC))[1] AS provenance_id`

	sqlSyaratPaketBertema = `$1 = ANY (p.tema_kode) AND p.tahun_anggaran = $2 AND NOT p.is_batal`

	sqlBelanjaKecamatan = `WITH indikator AS (
    SELECT DISTINCT ON (kode_kecamatan) * FROM indikator_kecamatan WHERE tema_kode = $1 ORDER BY kode_kecamatan, periode DESC
), belanja AS (
    SELECT left(p.kode_wilayah, 8) AS kode_kecamatan, ` + sqlAgregatPaket + `
    FROM paket p JOIN provenance pr ON pr.id = p.provenance_id
    WHERE ` + sqlSyaratPaketBertema + ` AND p.lokasi_tingkat IN ('kecamatan', 'kelurahan')
    GROUP BY 1
), cadangan AS (
    SELECT id FROM provenance WHERE sumber = 'lpse' ORDER BY diambil_pada DESC LIMIT 1
)
SELECT w.kode, w.nama, w.induk_kode, i.periode, i.kebutuhan, i.kapasitas, i.kesenjangan, i.rasio, i.kelas, i.label, i.provenance_id,
    coalesce(b.jumlah, 0), coalesce(b.total, 0), coalesce(b.simulasi, false), coalesce(b.provenance_id, (SELECT id FROM cadangan), '')
FROM wilayah w
LEFT JOIN indikator i ON i.kode_kecamatan = w.kode
LEFT JOIN belanja b ON b.kode_kecamatan = w.kode
WHERE w.tingkat = 'kecamatan'
ORDER BY w.kode`

	sqlBelanjaLuarKecamatan = `SELECT w.kode, w.nama, w.tingkat, ` + sqlAgregatPaket + `
FROM paket p
JOIN provenance pr ON pr.id = p.provenance_id
JOIN wilayah w ON w.kode = p.kode_wilayah
WHERE ` + sqlSyaratPaketBertema + ` AND p.lokasi_tingkat IN ('kota', 'provinsi')
GROUP BY w.kode, w.nama, w.tingkat
ORDER BY w.kode`
)

type BelanjaTemaRepository struct {
	DB Basisdata
}

type indikatorOpsional struct {
	periode, kelas, provenanceID             *string
	kebutuhan, kapasitas, kesenjangan, rasio *float64
	label                                    *domain.Label
}

func (r BelanjaTemaRepository) BelanjaKecamatan(ctx context.Context, kodeTema string, tahun int) ([]domain.BelanjaKecamatan, error) {
	baris, err := r.DB.Query(ctx, sqlBelanjaKecamatan, kodeTema, tahun)
	if err != nil {
		return nil, fmt.Errorf("belanja kecamatan tema %s: %w", kodeTema, err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.BelanjaKecamatan, error) {
		var k domain.BelanjaKecamatan
		var i indikatorOpsional
		a := domain.RingkasanAnggaran{Tahun: tahun}
		err := b.Scan(&k.Kecamatan.Kode, &k.Kecamatan.Nama, &k.KodeKota, &i.periode, &i.kebutuhan, &i.kapasitas, &i.kesenjangan,
			&i.rasio, &i.kelas, &i.label, &i.provenanceID, &a.JumlahPaket, &a.TotalHPS, &a.AdaSimulasi, &a.ProvenanceID)
		k.Kecamatan.Tingkat = domain.TingkatKecamatan
		k.Indikator = i.keDomain()
		if a.ProvenanceID != "" {
			k.Anggaran = &a
		}
		return k, err
	})
}

func (i indikatorOpsional) keDomain() *domain.IndikatorKecamatan {
	if i.periode == nil {
		return nil
	}
	return &domain.IndikatorKecamatan{
		Periode: *i.periode, Kebutuhan: *i.kebutuhan, Kapasitas: *i.kapasitas, Kesenjangan: *i.kesenjangan,
		Rasio: i.rasio, Kelas: *i.kelas, Label: *i.label, ProvenanceID: *i.provenanceID,
	}
}

func (r BelanjaTemaRepository) BelanjaLuarKecamatan(ctx context.Context, kodeTema string, tahun int) ([]domain.BelanjaLuarKecamatan, error) {
	baris, err := r.DB.Query(ctx, sqlBelanjaLuarKecamatan, kodeTema, tahun)
	if err != nil {
		return nil, fmt.Errorf("belanja luar kecamatan tema %s: %w", kodeTema, err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.BelanjaLuarKecamatan, error) {
		l := domain.BelanjaLuarKecamatan{Anggaran: domain.RingkasanAnggaran{Tahun: tahun}}
		err := b.Scan(&l.Wilayah.Kode, &l.Wilayah.Nama, &l.Wilayah.Tingkat,
			&l.Anggaran.JumlahPaket, &l.Anggaran.TotalHPS, &l.Anggaran.AdaSimulasi, &l.Anggaran.ProvenanceID)
		return l, err
	})
}
