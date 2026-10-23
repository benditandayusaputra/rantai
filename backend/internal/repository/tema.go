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
	sqlSimpanTema = `INSERT INTO tema (kode, nama_id, nama_en, aktif, urutan, konfigurasi, provenance_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (kode) DO UPDATE SET
    nama_id = EXCLUDED.nama_id, nama_en = EXCLUDED.nama_en, aktif = EXCLUDED.aktif, urutan = EXCLUDED.urutan,
    konfigurasi = EXCLUDED.konfigurasi, provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (tema.nama_id, tema.nama_en, tema.aktif, tema.urutan, tema.konfigurasi)
    IS DISTINCT FROM (EXCLUDED.nama_id, EXCLUDED.nama_en, EXCLUDED.aktif, EXCLUDED.urutan, EXCLUDED.konfigurasi)
RETURNING (xmax = 0) AS baru`

	sqlKolomTema = `SELECT kode, nama_id, nama_en, aktif, urutan, konfigurasi, provenance_id FROM tema`

	sqlDaftarTemaAktif = sqlKolomTema + ` WHERE aktif ORDER BY urutan, kode`

	sqlAmbilTemaAktif = sqlKolomTema + ` WHERE aktif AND kode = $1`
)

type TemaRepository struct {
	DB Basisdata
}

func (r TemaRepository) SimpanTema(ctx context.Context, data usecase.DataSumberTema) (usecase.HasilImpor, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("mulai transaksi tema: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := simpanProvenance(ctx, tx, data.Provenance); err != nil {
		return usecase.HasilImpor{}, err
	}
	hasil := usecase.HasilImpor{}
	for _, tema := range data.Tema {
		baru, berubah, err := simpanSatuTema(ctx, tx, tema, data.Provenance.ID)
		if err != nil {
			return usecase.HasilImpor{}, err
		}
		hasil.JumlahBaru += hitungJika(baru)
		hasil.JumlahBerubah += hitungJika(berubah)
	}
	return hasil, tx.Commit(ctx)
}

func simpanSatuTema(ctx context.Context, tx pgx.Tx, t domain.Tema, provenanceID string) (baru, berubah bool, err error) {
	konfigurasi, err := json.Marshal(t.KonfigurasiTema)
	if err != nil {
		return false, false, fmt.Errorf("enkode konfigurasi tema %s: %w", t.Kode, err)
	}
	err = tx.QueryRow(ctx, sqlSimpanTema, t.Kode, t.Nama.ID, t.Nama.EN, t.Aktif, t.Urutan, konfigurasi, provenanceID).Scan(&baru)
	_, ada, err := hasilSatuBaris(baru, err)
	if err != nil {
		return false, false, fmt.Errorf("simpan tema %s: %w", t.Kode, err)
	}
	return ada && baru, ada && !baru, nil
}

func (r TemaRepository) DaftarTemaAktif(ctx context.Context) ([]domain.TemaTersimpan, error) {
	baris, err := r.DB.Query(ctx, sqlDaftarTemaAktif)
	if err != nil {
		return nil, fmt.Errorf("kueri tema aktif: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.TemaTersimpan, error) {
		return pindaiTema(b)
	})
}

func (r TemaRepository) AmbilTemaAktif(ctx context.Context, kode string) (domain.TemaTersimpan, bool, error) {
	tema, err := pindaiTema(r.DB.QueryRow(ctx, sqlAmbilTemaAktif, kode))
	return hasilSatuBaris(tema, err)
}

func pindaiTema(baris pgx.Row) (domain.TemaTersimpan, error) {
	var t domain.TemaTersimpan
	err := baris.Scan(&t.Kode, &t.Nama.ID, &t.Nama.EN, &t.Aktif, &t.Urutan, &t.KonfigurasiTema, &t.ProvenanceID)
	return t, err
}
