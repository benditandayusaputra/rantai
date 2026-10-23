package repository

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlSimpanProvenance = `INSERT INTO provenance (id, sumber, url, diambil_pada, label, hash_konten, catatan)
VALUES ($1, $2, $3, $4, $5, $6, $7)`

	sqlSimpanProvinsi = `INSERT INTO wilayah (kode, tingkat, nama, induk_kode, alias, geom, titik_tengah, provenance_id)
SELECT $1, 'provinsi', $2, NULL, $3, g, ST_PointOnSurface(g), $5
FROM (
    SELECT ST_Multi(ST_CollectionExtract(ST_Union(ST_MakeValid(ST_SetSRID(ST_GeomFromGeoJSON(j), 4326))), 3)) AS g
    FROM unnest($4::text[]) AS j
) AS gabungan
ON CONFLICT (kode) DO UPDATE SET
    nama = EXCLUDED.nama, alias = EXCLUDED.alias, geom = EXCLUDED.geom,
    titik_tengah = EXCLUDED.titik_tengah, provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (wilayah.nama, wilayah.alias) IS DISTINCT FROM (EXCLUDED.nama, EXCLUDED.alias)
    OR NOT ST_OrderingEquals(wilayah.geom, EXCLUDED.geom)`

	sqlSimpanWilayah = `INSERT INTO wilayah (kode, tingkat, nama, induk_kode, alias, geom, titik_tengah, provenance_id)
SELECT $1, $2, $3, $4, $5, g, ST_PointOnSurface(g), $7
FROM (SELECT ST_Multi(ST_CollectionExtract(ST_MakeValid(ST_SetSRID(ST_GeomFromGeoJSON($6), 4326)), 3)) AS g) AS geometri
ON CONFLICT (kode) DO UPDATE SET
    tingkat = EXCLUDED.tingkat, nama = EXCLUDED.nama, induk_kode = EXCLUDED.induk_kode, alias = EXCLUDED.alias,
    geom = EXCLUDED.geom, titik_tengah = EXCLUDED.titik_tengah, provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (wilayah.tingkat, wilayah.nama, wilayah.induk_kode, wilayah.alias)
        IS DISTINCT FROM (EXCLUDED.tingkat, EXCLUDED.nama, EXCLUDED.induk_kode, EXCLUDED.alias)
    OR NOT ST_OrderingEquals(wilayah.geom, EXCLUDED.geom)
RETURNING (xmax = 0) AS baru`
)

type geoJSONMultiPolygon struct {
	Type        string              `json:"type"`
	Coordinates domain.MultiPolygon `json:"coordinates"`
}

type WilayahRepository struct {
	DB Basisdata
}

func (r WilayahRepository) SimpanWilayah(ctx context.Context, data usecase.DataSimpanWilayah) (usecase.HasilImpor, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("mulai transaksi wilayah: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	hasil, err := simpanSemua(ctx, tx, data)
	if err != nil {
		return usecase.HasilImpor{}, err
	}
	return hasil, tx.Commit(ctx)
}

func simpanSemua(ctx context.Context, tx pgx.Tx, data usecase.DataSimpanWilayah) (usecase.HasilImpor, error) {
	if err := simpanProvenance(ctx, tx, data.Provenance); err != nil {
		return usecase.HasilImpor{}, err
	}
	if err := simpanProvinsi(ctx, tx, data); err != nil {
		return usecase.HasilImpor{}, err
	}
	return simpanDaftarWilayah(ctx, tx, data.Wilayah, data.Provenance.ID)
}

func simpanProvenance(ctx context.Context, tx pgx.Tx, p domain.Provenance) error {
	_, err := tx.Exec(ctx, sqlSimpanProvenance, p.ID, p.Sumber, p.URL, p.DiambilPada, p.Label, p.HashKonten, p.Catatan)
	if err != nil {
		return fmt.Errorf("simpan provenance %s: %w", p.ID, err)
	}
	return nil
}

func simpanProvinsi(ctx context.Context, tx pgx.Tx, data usecase.DataSimpanWilayah) error {
	geometriKota, err := geometriTingkat(data.Wilayah, domain.TingkatKota)
	if err != nil {
		return err
	}
	p := data.Provinsi
	if _, err := tx.Exec(ctx, sqlSimpanProvinsi, p.Kode, p.Nama, p.Alias, geometriKota, data.Provenance.ID); err != nil {
		return fmt.Errorf("simpan provinsi %s: %w", p.Kode, err)
	}
	return nil
}

func simpanDaftarWilayah(ctx context.Context, tx pgx.Tx, daftar []domain.Wilayah, provenanceID string) (usecase.HasilImpor, error) {
	hasil := usecase.HasilImpor{}
	for _, w := range urutkanIndukDahulu(daftar) {
		baru, berubah, err := simpanSatuWilayah(ctx, tx, w, provenanceID)
		if err != nil {
			return usecase.HasilImpor{}, err
		}
		hasil.JumlahBaru += hitungJika(baru)
		hasil.JumlahBerubah += hitungJika(berubah)
	}
	return hasil, nil
}

func simpanSatuWilayah(ctx context.Context, tx pgx.Tx, w domain.Wilayah, provenanceID string) (baru, berubah bool, err error) {
	geometri, err := keGeoJSON(w.Geometri)
	if err != nil {
		return false, false, fmt.Errorf("enkode geometri %s: %w", w.Kode, err)
	}
	err = tx.QueryRow(ctx, sqlSimpanWilayah, w.Kode, w.Tingkat, w.Nama, w.IndukKode, w.Alias, geometri, provenanceID).Scan(&baru)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, false, nil
	}
	if err != nil {
		return false, false, fmt.Errorf("simpan wilayah %s: %w", w.Kode, err)
	}
	return baru, !baru, nil
}

func urutkanIndukDahulu(daftar []domain.Wilayah) []domain.Wilayah {
	return slices.SortedStableFunc(slices.Values(daftar), func(a, b domain.Wilayah) int {
		return cmp.Compare(len(a.Kode), len(b.Kode))
	})
}

func geometriTingkat(daftar []domain.Wilayah, tingkat domain.Tingkat) ([]string, error) {
	hasil := []string{}
	for _, w := range daftar {
		if w.Tingkat != tingkat {
			continue
		}
		geometri, err := keGeoJSON(w.Geometri)
		if err != nil {
			return nil, fmt.Errorf("enkode geometri %s: %w", w.Kode, err)
		}
		hasil = append(hasil, geometri)
	}
	return hasil, nil
}

func keGeoJSON(geometri domain.MultiPolygon) (string, error) {
	isi, err := json.Marshal(geoJSONMultiPolygon{Type: "MultiPolygon", Coordinates: geometri})
	return string(isi), err
}

func hitungJika(kondisi bool) int {
	if kondisi {
		return 1
	}
	return 0
}
