package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlDaftarAnggotaUmpan = `SELECT id, nama, nama_lain, partai, fraksi FROM anggota ORDER BY id`

	sqlTerbitTerakhir = `SELECT max(terbit_pada) FROM umpan_kegiatan WHERE jenis_sumber = $1`

	sqlSimpanUmpan = `INSERT INTO umpan_kegiatan (jenis_sumber, media_nama, media_domain, judul, cuplikan, url, terbit_pada, tag_isu, provenance_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (url) DO UPDATE SET
    media_nama = EXCLUDED.media_nama, media_domain = EXCLUDED.media_domain, judul = EXCLUDED.judul,
    cuplikan = EXCLUDED.cuplikan, terbit_pada = EXCLUDED.terbit_pada, tag_isu = EXCLUDED.tag_isu,
    ringkasan = CASE WHEN (umpan_kegiatan.judul, umpan_kegiatan.cuplikan) IS DISTINCT FROM (EXCLUDED.judul, EXCLUDED.cuplikan)
        THEN NULL ELSE umpan_kegiatan.ringkasan END,
    provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (umpan_kegiatan.media_nama, umpan_kegiatan.judul, umpan_kegiatan.cuplikan, umpan_kegiatan.terbit_pada, umpan_kegiatan.tag_isu)
    IS DISTINCT FROM (EXCLUDED.media_nama, EXCLUDED.judul, EXCLUDED.cuplikan, EXCLUDED.terbit_pada, EXCLUDED.tag_isu)
RETURNING id, (xmax = 0) AS baru`

	sqlIDUmpan = `SELECT id FROM umpan_kegiatan WHERE url = $1`

	sqlHapusPautan = `DELETE FROM umpan_anggota WHERE umpan_id = $1`

	sqlSimpanPautan = `INSERT INTO umpan_anggota (umpan_id, anggota_id, dasar_cocok)
SELECT $1, x.anggota_id, x.dasar::dasar_cocok_umpan
FROM unnest($2::text[], $3::text[]) AS x (anggota_id, dasar)`
)

type UmpanRepository struct {
	DB Basisdata
}

func (r UmpanRepository) DaftarAnggota(ctx context.Context) ([]domain.Anggota, error) {
	baris, err := r.DB.Query(ctx, sqlDaftarAnggotaUmpan)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.Anggota, error) {
		var a domain.Anggota
		return a, b.Scan(&a.ID, &a.Nama, &a.NamaLain, &a.Partai, &a.Fraksi)
	})
}

func (r UmpanRepository) TerbitTerakhir(ctx context.Context, jenis domain.JenisSumberUmpan) (time.Time, bool, error) {
	var terakhir *time.Time
	if err := r.DB.QueryRow(ctx, sqlTerbitTerakhir, jenis).Scan(&terakhir); err != nil {
		return time.Time{}, false, fmt.Errorf("baca terbit terakhir %s: %w", jenis, err)
	}
	if terakhir == nil {
		return time.Time{}, false, nil
	}
	return *terakhir, true, nil
}

func (r UmpanRepository) SimpanUmpan(ctx context.Context, data usecase.DataSimpanUmpan) (usecase.HasilImpor, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("mulai transaksi umpan: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := simpanProvenance(ctx, tx, data.Provenance); err != nil {
		return usecase.HasilImpor{}, err
	}
	hasil := usecase.HasilImpor{}
	for _, u := range data.Umpan {
		baru, berubah, err := simpanSatuUmpan(ctx, tx, u, data.Provenance.ID)
		if err != nil {
			return usecase.HasilImpor{}, fmt.Errorf("simpan umpan %s: %w", u.URL, err)
		}
		hasil.JumlahBaru += hitungJika(baru)
		hasil.JumlahBerubah += hitungJika(berubah)
	}
	return hasil, tx.Commit(ctx)
}

func simpanSatuUmpan(ctx context.Context, tx pgx.Tx, u domain.UmpanTerpaut, provenanceID string) (baru, berubah bool, err error) {
	var id int64
	err = tx.QueryRow(ctx, sqlSimpanUmpan, u.JenisSumber, u.MediaNama, u.MediaDomain, u.Judul, u.Cuplikan,
		u.URL, u.TerbitPada, u.TagIsu, provenanceID).Scan(&id, &baru)
	tetap := errors.Is(err, pgx.ErrNoRows)
	if tetap {
		err = tx.QueryRow(ctx, sqlIDUmpan, u.URL).Scan(&id)
	}
	if err != nil {
		return false, false, err
	}
	return baru, !baru && !tetap, simpanPautan(ctx, tx, id, u.Pautan)
}

func simpanPautan(ctx context.Context, tx pgx.Tx, umpanID int64, pautan []domain.PautanUmpan) error {
	anggota, dasar := []string{}, []string{}
	for _, p := range pautan {
		anggota, dasar = append(anggota, p.AnggotaID), append(dasar, string(p.DasarCocok))
	}
	if _, err := tx.Exec(ctx, sqlHapusPautan, umpanID); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, sqlSimpanPautan, umpanID, anggota, dasar)
	return err
}
