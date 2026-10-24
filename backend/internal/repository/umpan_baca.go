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

	sqlKolomUmpanTampil = `u.id, u.jenis_sumber, u.media_nama, u.media_domain, u.judul, coalesce(u.ringkasan, ''), u.url,
    u.terbit_pada, u.tag_isu, ua.dasar_cocok, u.provenance_id`

	sqlUmpanAnggota = `SELECT ` + sqlKolomUmpanTampil + `, count(*) OVER ()
FROM umpan_anggota ua
JOIN umpan_kegiatan u ON u.id = ua.umpan_id
WHERE ua.anggota_id = $1 AND NOT u.disembunyikan
    AND ($2 = '' OR u.jenis_sumber::text = $2)
    AND ($3 = '' OR $3 = ANY (u.tag_isu))
ORDER BY u.terbit_pada DESC, u.id DESC
LIMIT $4 OFFSET $5`

	sqlUmpanBertagAnggota = `SELECT ` + sqlKolomUmpanTampil + `
FROM umpan_anggota ua
JOIN umpan_kegiatan u ON u.id = ua.umpan_id
WHERE ua.anggota_id = $1 AND NOT u.disembunyikan AND cardinality(u.tag_isu) > 0
ORDER BY u.terbit_pada DESC, u.id DESC`

	sqlUmpanIsuTerbaru = `SELECT * FROM (
    SELECT ua.anggota_id, ` + sqlKolomUmpanTampil + `,
        row_number() OVER (PARTITION BY ua.anggota_id ORDER BY u.terbit_pada DESC, u.id DESC) AS urutan
    FROM umpan_anggota ua
    JOIN umpan_kegiatan u ON u.id = ua.umpan_id
    WHERE ua.anggota_id = ANY ($1) AND NOT u.disembunyikan AND u.tag_isu && $2
) terbaru
WHERE urutan <= $3
ORDER BY anggota_id, urutan`
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
	_, err = pgx.ForEachRow(baris, append(kolomUmpanTampil(&u), &halaman.Total), func() error {
		halaman.Item = append(halaman.Item, u)
		return nil
	})
	return halaman, err
}

func (r UmpanRepository) UmpanBertagAnggota(ctx context.Context, anggotaID string) ([]domain.UmpanTampil, error) {
	baris, err := r.DB.Query(ctx, sqlUmpanBertagAnggota, anggotaID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.UmpanTampil, error) {
		var u domain.UmpanTampil
		return u, b.Scan(kolomUmpanTampil(&u)...)
	})
}

func (r UmpanRepository) UmpanIsuTerbaru(ctx context.Context, anggotaID, tagIsu []string) (map[string][]domain.UmpanTampil, error) {
	baris, err := r.DB.Query(ctx, sqlUmpanIsuTerbaru, anggotaID, tagIsu, domain.BatasUmpanIsuPerWakil)
	if err != nil {
		return nil, err
	}
	hasil := map[string][]domain.UmpanTampil{}
	var pemilik string
	var urutan int
	var u domain.UmpanTampil
	_, err = pgx.ForEachRow(baris, append(append([]any{&pemilik}, kolomUmpanTampil(&u)...), &urutan), func() error {
		hasil[pemilik] = append(hasil[pemilik], u)
		return nil
	})
	return hasil, err
}

func kolomUmpanTampil(u *domain.UmpanTampil) []any {
	return []any{&u.ID, &u.JenisSumber, &u.MediaNama, &u.MediaDomain, &u.Judul, &u.Ringkasan,
		&u.URL, &u.TerbitPada, &u.TagIsu, &u.DasarCocok, &u.ProvenanceID}
}
