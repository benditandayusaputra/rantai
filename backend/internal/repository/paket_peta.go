package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlFilterPeta = `($5 = '' OR $5 = ANY (p.tema_kode)) AND ($6 = 0 OR p.tahun_anggaran = $6)`

	sqlTitikPaket = `SELECT p.kode, p.nama, coalesce(p.sektor, ''), p.is_simulation, ST_Y(p.lokasi), ST_X(p.lokasi), 0
FROM paket p
WHERE p.lokasi && ST_MakeEnvelope($1, $2, $3, $4, 4326) AND ` + sqlFilterPeta + `
ORDER BY p.kode`

	sqlAreaPaket = `SELECT w.kode, w.nama, ST_AsGeoJSON(ST_SimplifyPreserveTopology(w.geom, $7)), array_agg(p.kode ORDER BY p.kode)
FROM paket p JOIN wilayah w ON w.kode = p.kode_wilayah
WHERE p.lokasi_tingkat = 'kecamatan' AND w.geom && ST_MakeEnvelope($1, $2, $3, $4, 4326) AND ` + sqlFilterPeta + `
GROUP BY w.kode, w.nama, w.geom
ORDER BY w.kode`

	sqlPaketTerdekat = `SELECT p.kode, p.nama, coalesce(p.sektor, ''), p.is_simulation, ST_Y(p.lokasi), ST_X(p.lokasi),
    ST_Distance(p.lokasi::geography, t.titik::geography) AS jarak
FROM paket p, (SELECT ST_SetSRID(ST_MakePoint($2, $1), 4326) AS titik) AS t
WHERE p.lokasi IS NOT NULL AND ST_DWithin(p.lokasi::geography, t.titik::geography, $3)
ORDER BY jarak, p.kode
LIMIT $4`
)

func (r PaketRepository) PetaPaket(ctx context.Context, f usecase.FilterPetaPaket) ([]domain.TitikPaket, []domain.AreaPaket, error) {
	k := f.Kotak
	titik, err := r.bacaTitik(ctx, sqlTitikPaket, k.Barat, k.Selatan, k.Timur, k.Utara, f.Tema, f.Tahun)
	if err != nil {
		return nil, nil, fmt.Errorf("peta titik paket: %w", err)
	}
	baris, err := r.DB.Query(ctx, sqlAreaPaket, k.Barat, k.Selatan, k.Timur, k.Utara, f.Tema, f.Tahun, f.Toleransi)
	if err != nil {
		return nil, nil, fmt.Errorf("peta area paket: %w", err)
	}
	area, err := pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.AreaPaket, error) {
		var a domain.AreaPaket
		err := b.Scan(&a.KodeKecamatan, &a.NamaKecamatan, &a.GeometriGeoJSON, &a.KodePaket)
		return a, err
	})
	return titik, area, err
}

func (r PaketRepository) PaketTerdekat(ctx context.Context, titik domain.Titik, radiusMeter float64, batas int) ([]domain.TitikPaket, error) {
	return r.bacaTitik(ctx, sqlPaketTerdekat, titik.Lat, titik.Lng, radiusMeter, batas)
}

func (r PaketRepository) bacaTitik(ctx context.Context, sql string, argumen ...any) ([]domain.TitikPaket, error) {
	baris, err := r.DB.Query(ctx, sql, argumen...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.TitikPaket, error) {
		var t domain.TitikPaket
		err := b.Scan(&t.Kode, &t.Nama, &t.Sektor, &t.IsSimulation, &t.Titik.Lat, &t.Titik.Lng, &t.JarakMeter)
		return t, err
	})
}
