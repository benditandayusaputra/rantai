package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	alasanTerverifikasiSimulasi = "dikonfirmasi tiga akun simulasi"

	sqlLokasiPaketSimulasi = `SELECT ST_Y(lokasi), ST_X(lokasi) FROM paket WHERE kode = $1 AND is_simulation AND lokasi IS NOT NULL`

	sqlHapusLaporanSimulasi = `DELETE FROM laporan WHERE is_simulation AND paket_kode = $1`

	sqlPenggunaSimulasi = `SELECT id FROM pengguna WHERE email = $1 AND is_simulation`

	sqlSimpanLaporanSimulasi = `INSERT INTO laporan (pengguna_id, paket_kode, kode_wilayah, lokasi_publik, lokasi_asli, akurasi_meter,
    kondisi, catatan, status, is_simulation)
SELECT $1, p.kode, p.kode_wilayah, ST_SetSRID(ST_MakePoint($3, $2), 4326), p.lokasi, 0, $4, $5, 'DITERUSKAN', true
FROM paket p WHERE p.kode = $6 AND p.is_simulation
RETURNING id`

	sqlRiwayatLaporanSimulasi = `INSERT INTO laporan_riwayat (laporan_id, status_lama, status_baru, oleh_pengguna_id, alasan)
VALUES ($1, NULL, 'DIKIRIM', $2, ''), ($1, 'DIKIRIM', 'MENUNGGU_SAKSI', NULL, ''), ($1, 'MENUNGGU_SAKSI', 'TERVERIFIKASI', NULL, $3), ($1, 'TERVERIFIKASI', 'DITERUSKAN', NULL, $4)`

	sqlKonfirmasiSimulasi = `INSERT INTO laporan_konfirmasi (laporan_id, pengguna_id, kondisi, lokasi_asli)
SELECT l.id, x.pengguna_id, l.kondisi, l.lokasi_asli FROM laporan l, unnest($2::bigint[]) AS x (pengguna_id) WHERE l.id = $1`
)

func (r SkenarioDemoRepository) LokasiPaketSimulasi(ctx context.Context, kode string) (domain.Titik, error) {
	var titik domain.Titik
	err := r.DB.QueryRow(ctx, sqlLokasiPaketSimulasi, kode).Scan(&titik.Lat, &titik.Lng)
	if errors.Is(err, pgx.ErrNoRows) {
		return titik, domain.TidakDitemukan("paket simulasi " + kode + " belum terpasang atau tanpa lokasi")
	}
	return titik, err
}

func (r SkenarioDemoRepository) SimpanLaporanSimulasi(ctx context.Context, data usecase.DataLaporanSimulasi) error {
	return dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, sqlHapusLaporanSimulasi, data.PaketKode); err != nil {
			return fmt.Errorf("hapus laporan simulasi lama: %w", err)
		}
		for _, l := range data.Laporan {
			if err := simpanSatuLaporanSimulasi(ctx, tx, data, l); err != nil {
				return err
			}
		}
		return nil
	})
}

func simpanSatuLaporanSimulasi(ctx context.Context, tx pgx.Tx, data usecase.DataLaporanSimulasi, l usecase.LaporanSimulasiTersimpan) error {
	pengguna, err := idPenggunaSimulasi(ctx, tx, append([]string{l.Pelapor}, l.Saksi...))
	if err != nil {
		return err
	}
	var id int64
	err = tx.QueryRow(ctx, sqlSimpanLaporanSimulasi, pengguna[0], data.TitikPublik.Lat, data.TitikPublik.Lng, l.Kondisi, l.Catatan, data.PaketKode).
		Scan(&id)
	if err != nil {
		return fmt.Errorf("simpan laporan simulasi: %w", err)
	}
	if err := simpanFotoLaporan(ctx, tx, id, l.Foto); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, sqlRiwayatLaporanSimulasi, id, pengguna[0], alasanTerverifikasiSimulasi, domain.AlasanDiteruskan); err != nil {
		return fmt.Errorf("simpan riwayat laporan simulasi: %w", err)
	}
	_, err = tx.Exec(ctx, sqlKonfirmasiSimulasi, id, pengguna)
	return err
}

func idPenggunaSimulasi(ctx context.Context, tx pgx.Tx, email []string) ([]int64, error) {
	id := make([]int64, len(email))
	for i, e := range email {
		err := tx.QueryRow(ctx, sqlPenggunaSimulasi, e).Scan(&id[i])
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", e, domain.ErrSimulasiTertautDataNyata)
		}
		if err != nil {
			return nil, fmt.Errorf("cari akun simulasi %s: %w", e, err)
		}
	}
	return id, nil
}
