package repository

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlFraksiAda = `SELECT EXISTS (SELECT 1 FROM anggota WHERE fraksi = $1)`

	sqlSikapFraksi = `SELECT fraksi, agenda, judul_agenda, tanggal, ringkasan, url, provenance_id
FROM sikap_fraksi
WHERE fraksi = $1
ORDER BY tanggal DESC, agenda, judul_agenda, url`

	sqlDapilAda = `SELECT EXISTS (SELECT 1 FROM dapil WHERE id = $1)`

	sqlAspirasiReses = `SELECT a.masa_reses, a.tanggal_paripurna, a.bidang, a.uraian, coalesce(a.kode_wilayah, ''),
    coalesce(w.nama, ''), coalesce(w.tingkat::text, ''), a.url, a.provenance_id
FROM aspirasi_reses a
LEFT JOIN wilayah w ON w.kode = a.kode_wilayah
WHERE ($1 = 0 OR a.kode_wilayah IS NULL OR EXISTS (
        SELECT 1 FROM dapil_wilayah dw
        WHERE dw.dapil_id = $1 AND (a.kode_wilayah = dw.kode_kecamatan OR a.kode_wilayah LIKE dw.kode_kecamatan || '.%')))
    AND ($2 = '' OR a.bidang = $2)
ORDER BY a.tanggal_paripurna DESC, a.kode_wilayah IS NULL, a.bidang, a.uraian`
)

func (r ParipurnaRepository) IsFraksiAda(ctx context.Context, fraksi string) (bool, error) {
	var ada bool
	return ada, r.DB.QueryRow(ctx, sqlFraksiAda, fraksi).Scan(&ada)
}

func (r ParipurnaRepository) SikapFraksi(ctx context.Context, fraksi string) ([]domain.SikapFraksi, error) {
	baris, err := r.DB.Query(ctx, sqlSikapFraksi, fraksi)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, pgx.RowToStructByPos[domain.SikapFraksi])
}

func (r ParipurnaRepository) IsDapilAda(ctx context.Context, id int) (bool, error) {
	var ada bool
	return ada, r.DB.QueryRow(ctx, sqlDapilAda, id).Scan(&ada)
}

func (r ParipurnaRepository) AspirasiReses(ctx context.Context, filter domain.FilterAspirasiReses) ([]domain.AspirasiReses, error) {
	baris, err := r.DB.Query(ctx, sqlAspirasiReses, filter.DapilID, filter.Bidang)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, pgx.RowToStructByPos[domain.AspirasiReses])
}
