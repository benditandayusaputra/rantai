package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlMasukanKesenjangan = `WITH periode AS (SELECT max(tahun_ajaran) AS t FROM calon_lulusan_kecamatan WHERE jenjang = $1),
tahun_kursi AS (SELECT max(tahun_ajaran) AS t FROM daya_tampung),
kursi AS (
    SELECT s.kode_kecamatan, d.jalur, s.bentuk, sum(d.kursi) AS kursi
    FROM daya_tampung d JOIN sekolah s ON s.npsn = d.npsn, tahun_kursi
    WHERE d.tahun_ajaran = tahun_kursi.t AND s.jenjang = $2
    GROUP BY s.kode_kecamatan, d.jalur, s.bentuk
),
kursi_kecamatan AS (
    SELECT kode_kecamatan,
        jsonb_object_agg(bentuk, kursi) FILTER (WHERE jalur = 'negeri') AS negeri,
        jsonb_object_agg(bentuk, kursi) FILTER (WHERE jalur = 'spmb_bersama') AS spmb_bersama
    FROM kursi
    GROUP BY kode_kecamatan
),
madrasah AS (
    SELECT s.kode_kecamatan, count(DISTINCT s.npsn) AS jumlah, count(DISTINCT r.npsn) AS dengan_rombel,
        coalesce(sum(r.jumlah_rombel) FILTER (WHERE r.tingkat = $4), 0) AS rombel_tingkat_pertama
    FROM sekolah s
    LEFT JOIN sekolah_rombel r ON r.npsn = s.npsn AND r.tahun_ajaran = (SELECT t FROM periode)
    WHERE s.bentuk = $3
    GROUP BY s.kode_kecamatan
)
SELECT w.kode, coalesce(p.t, ''), coalesce(c.jumlah, 0), coalesce(c.label, 'resmi')::text, coalesce(c.metode, ''),
    coalesce(k.negeri, '{}'::jsonb), coalesce(k.spmb_bersama, '{}'::jsonb),
    coalesce(m.jumlah, 0)::int, coalesce(m.dengan_rombel, 0)::int, coalesce(m.rombel_tingkat_pertama, 0)::int
FROM wilayah w
CROSS JOIN periode p
LEFT JOIN calon_lulusan_kecamatan c ON c.kode_kecamatan = w.kode AND c.tahun_ajaran = p.t AND c.jenjang = $1
LEFT JOIN kursi_kecamatan k ON k.kode_kecamatan = w.kode
LEFT JOIN madrasah m ON m.kode_kecamatan = w.kode
WHERE w.tingkat = 'kecamatan'
ORDER BY w.kode`

	sqlSimpanIndikator = `INSERT INTO indikator_kecamatan
    (tema_kode, kode_kecamatan, periode, kebutuhan, kapasitas, kesenjangan, rasio, kelas, label, rincian, provenance_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (tema_kode, kode_kecamatan, periode) DO UPDATE SET
    kebutuhan = EXCLUDED.kebutuhan, kapasitas = EXCLUDED.kapasitas, kesenjangan = EXCLUDED.kesenjangan,
    rasio = EXCLUDED.rasio, kelas = EXCLUDED.kelas, label = EXCLUDED.label, rincian = EXCLUDED.rincian,
    provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (indikator_kecamatan.kebutuhan, indikator_kecamatan.kapasitas, indikator_kecamatan.kelas,
        indikator_kecamatan.label, indikator_kecamatan.rincian)
    IS DISTINCT FROM (EXCLUDED.kebutuhan, EXCLUDED.kapasitas, EXCLUDED.kelas, EXCLUDED.label, EXCLUDED.rincian)
RETURNING (xmax = 0) AS baru`
)

type IndikatorRepository struct {
	DB Basisdata
}

func (r IndikatorRepository) MasukanKesenjangan(ctx context.Context, p domain.Peralihan) (string, []domain.MasukanKesenjangan, error) {
	baris, err := r.DB.Query(ctx, sqlMasukanKesenjangan, p.JenjangAsal, p.JenjangTujuan, p.BentukMadrasahTujuan, p.TingkatTujuanPertama())
	if err != nil {
		return "", nil, fmt.Errorf("kueri masukan kesenjangan %s: %w", p.Kode, err)
	}
	periode := ""
	daftar, err := pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.MasukanKesenjangan, error) {
		var m domain.MasukanKesenjangan
		err := b.Scan(&m.KodeKecamatan, &periode, &m.CalonLulusan, &m.LabelCalonLulusan, &m.MetodeCalonLulusan,
			&m.KursiNegeri, &m.KursiSPMBBersama, &m.JumlahMadrasah, &m.JumlahMadrasahDenganRombel, &m.RombelTingkatPertamaMadrasah)
		return m, err
	})
	return periode, daftar, err
}

func (r IndikatorRepository) SimpanIndikator(ctx context.Context, data usecase.DataIndikator) (usecase.HasilImpor, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("mulai transaksi indikator: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	for _, p := range data.Provenance {
		if err := simpanProvenance(ctx, tx, p); err != nil {
			return usecase.HasilImpor{}, err
		}
	}
	hasil := usecase.HasilImpor{}
	for _, indikator := range data.Daftar {
		baru, berubah, err := simpanSatuIndikator(ctx, tx, indikator)
		if err != nil {
			return usecase.HasilImpor{}, err
		}
		hasil.JumlahBaru += hitungJika(baru)
		hasil.JumlahBerubah += hitungJika(berubah)
	}
	return hasil, tx.Commit(ctx)
}

func simpanSatuIndikator(ctx context.Context, tx pgx.Tx, i usecase.IndikatorTercatat) (baru, berubah bool, err error) {
	rincian, err := json.Marshal(i.Rincian)
	if err != nil {
		return false, false, fmt.Errorf("enkode rincian %s: %w", i.KodeKecamatan, err)
	}
	err = tx.QueryRow(ctx, sqlSimpanIndikator, i.TemaKode, i.KodeKecamatan, i.Periode, i.Kebutuhan, i.Kapasitas,
		i.Kesenjangan, i.Rasio, i.Kelas, i.Label, rincian, i.ProvenanceID).Scan(&baru)
	_, ada, err := hasilSatuBaris(baru, err)
	if err != nil {
		return false, false, fmt.Errorf("simpan indikator %s: %w", i.KodeKecamatan, err)
	}
	return ada && baru, ada && !baru, nil
}
