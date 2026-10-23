package repository

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlUmpanBelumDiringkas = `SELECT id, judul, cuplikan FROM umpan_kegiatan
WHERE ringkasan IS NULL AND NOT disembunyikan
ORDER BY terbit_pada DESC
LIMIT $1`

	sqlSimpanRingkasan = `UPDATE umpan_kegiatan SET ringkasan = $2, diubah_pada = now() WHERE id = $1`

	sqlAnggotaAda = `SELECT EXISTS (SELECT 1 FROM anggota WHERE id = $1)`

	sqlUmpanAnggota = `SELECT u.id, u.jenis_sumber, u.media_nama, u.media_domain, u.judul, coalesce(u.ringkasan, ''), u.url,
    u.terbit_pada, u.tag_isu, ua.dasar_cocok, u.provenance_id, count(*) OVER ()
FROM umpan_anggota ua
JOIN umpan_kegiatan u ON u.id = ua.umpan_id
WHERE ua.anggota_id = $1 AND NOT u.disembunyikan
    AND ($2 = '' OR u.jenis_sumber::text = $2)
    AND ($3 = '' OR $3 = ANY (u.tag_isu))
ORDER BY u.terbit_pada DESC, u.id DESC
LIMIT $4 OFFSET $5`
)

func (r UmpanRepository) UmpanBelumDiringkas(ctx context.Context, batas int) ([]domain.UmpanUntukRingkas, error) {
	baris, err := r.DB.Query(ctx, sqlUmpanBelumDiringkas, batas)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, pgx.RowToStructByPos[domain.UmpanUntukRingkas])
}

func (r UmpanRepository) SimpanRingkasan(ctx context.Context, id int64, ringkasan string) error {
	_, err := r.DB.Exec(ctx, sqlSimpanRingkasan, id, ringkasan)
	return err
}

func (r UmpanRepository) IsAnggotaAda(ctx context.Context, id string) (bool, error) {
	var ada bool
	return ada, r.DB.QueryRow(ctx, sqlAnggotaAda, id).Scan(&ada)
}

func (r UmpanRepository) UmpanAnggota(ctx context.Context, filter domain.FilterUmpan) (domain.HalamanUmpan, error) {
	baris, err := r.DB.Query(ctx, sqlUmpanAnggota, filter.AnggotaID, string(filter.Jenis), filter.Tag,
		filter.PerHalaman, (filter.Halaman-1)*filter.PerHalaman)
	if err != nil {
		return domain.HalamanUmpan{}, err
	}
	halaman := domain.HalamanUmpan{Item: []domain.UmpanTampil{}}
	var u domain.UmpanTampil
	_, err = pgx.ForEachRow(baris, []any{&u.ID, &u.JenisSumber, &u.MediaNama, &u.MediaDomain, &u.Judul, &u.Ringkasan,
		&u.URL, &u.TerbitPada, &u.TagIsu, &u.DasarCocok, &u.ProvenanceID, &halaman.Total}, func() error {
		halaman.Item = append(halaman.Item, u)
		return nil
	})
	return halaman, err
}
