package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const ambangKemiripan = 0.45

const (
	sqlCariWilayah = `WITH skor AS (
    SELECT w.kode, w.nama, w.tingkat, w.induk_kode,
        $1 = lower(w.nama) OR $1 = ANY (SELECT lower(a) FROM unnest(w.alias) AS a) AS persis,
        (SELECT max(greatest(similarity(lower(n), $1), word_similarity(lower(n), $1), word_similarity($1, lower(n))))
         FROM unnest(array_append(w.alias, w.nama)) AS n) AS nilai
    FROM wilayah w
),
cocok AS (SELECT * FROM skor WHERE persis OR nilai >= $2),
anak AS (
    SELECT s.* FROM skor s JOIN cocok c ON s.induk_kode = c.kode
    WHERE c.persis AND c.tingkat = 'kecamatan'
),
gabungan AS (SELECT * FROM cocok UNION SELECT * FROM anak)
SELECT kode, nama, tingkat, hierarki_wilayah(kode)
FROM gabungan
ORDER BY persis DESC, nilai DESC, tingkat DESC, nama
LIMIT $3`

	sqlKelurahanDiTitik = `SELECT kode, nama, tingkat, hierarki_wilayah(kode)
FROM wilayah
WHERE tingkat = 'kelurahan' AND ST_Covers(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326))
ORDER BY kode
LIMIT 1`

	sqlAmbilWilayah = `SELECT kode, nama, tingkat, hierarki_wilayah(kode), alias,
    ST_Y(titik_tengah), ST_X(titik_tengah), provenance_id
FROM wilayah
WHERE kode = $1`
)

func (r WilayahRepository) CariWilayah(ctx context.Context, kueri string, batas int) ([]domain.WilayahDitemukan, error) {
	baris, err := r.DB.Query(ctx, sqlCariWilayah, kueri, ambangKemiripan, batas)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.WilayahDitemukan, error) {
		return pindaiWilayahDitemukan(b)
	})
}

func (r WilayahRepository) KelurahanDiTitik(ctx context.Context, titik domain.Titik) (domain.WilayahDitemukan, bool, error) {
	wilayah, err := pindaiWilayahDitemukan(r.DB.QueryRow(ctx, sqlKelurahanDiTitik, titik.Lng, titik.Lat))
	return hasilSatuBaris(wilayah, err)
}

func (r WilayahRepository) AmbilWilayah(ctx context.Context, kode string) (domain.DetailWilayah, bool, error) {
	var w domain.DetailWilayah
	err := r.DB.QueryRow(ctx, sqlAmbilWilayah, kode).Scan(
		&w.Kode, &w.Nama, &w.Tingkat, &w.Hierarki, &w.Alias,
		&w.TitikTengah.Lat, &w.TitikTengah.Lng, &w.ProvenanceID,
	)
	return hasilSatuBaris(w, err)
}

func pindaiWilayahDitemukan(baris pgx.Row) (domain.WilayahDitemukan, error) {
	var w domain.WilayahDitemukan
	err := baris.Scan(&w.Kode, &w.Nama, &w.Tingkat, &w.Hierarki)
	return w, err
}

func hasilSatuBaris[T any](nilai T, err error) (T, bool, error) {
	var kosong T
	if errors.Is(err, pgx.ErrNoRows) {
		return kosong, false, nil
	}
	if err != nil {
		return kosong, false, fmt.Errorf("baca baris: %w", err)
	}
	return nilai, true, nil
}
