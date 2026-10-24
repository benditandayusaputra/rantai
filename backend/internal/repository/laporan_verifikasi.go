package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	kodePelanggaranUnik = "23505"

	sqlLaporanUntukVerifikasi = `SELECT id, pengguna_id, status, kondisi, ST_Y(lokasi_asli), ST_X(lokasi_asli), dibuat_pada
FROM laporan WHERE id = $1`

	sqlKonfirmasiLaporan = `SELECT pengguna_id, kondisi, ST_Y(lokasi_asli), ST_X(lokasi_asli), dibuat_pada
FROM laporan_konfirmasi WHERE laporan_id = $1 ORDER BY dibuat_pada`

	sqlSimpanKonfirmasi = `INSERT INTO laporan_konfirmasi (laporan_id, pengguna_id, kondisi, lokasi_asli, dibuat_pada)
VALUES ($1, $2, $3, ST_SetSRID(ST_MakePoint($5, $4), 4326), $6)`

	sqlUbahStatusLaporan = `UPDATE laporan SET status = $3, diubah_pada = now() WHERE id = $1 AND status = $2`

	sqlCatatRiwayatLaporan = `INSERT INTO laporan_riwayat (laporan_id, status_lama, status_baru, oleh_pengguna_id, alasan)
VALUES ($1, $2, $3, $4, $5)`

	sqlLaporanPublik = `SELECT l.id, l.status, l.kondisi, l.catatan, l.kode_wilayah, w.nama, ST_Y(l.lokasi_publik), ST_X(l.lokasi_publik),
    coalesce(l.paket_kode, ''), coalesce(p.nama, ''), l.is_simulation, l.dibuat_pada
FROM laporan l
JOIN wilayah w ON w.kode = l.kode_wilayah
LEFT JOIN paket p ON p.kode = l.paket_kode
WHERE l.id = $1`

	sqlFotoLaporan = `SELECT object_key, thumbnail_key, lebar, tinggi FROM laporan_foto WHERE laporan_id = $1 ORDER BY urutan`

	sqlRiwayatLaporan = `SELECT coalesce(status_lama::text, ''), status_baru, alasan, dibuat_pada
FROM laporan_riwayat WHERE laporan_id = $1 ORDER BY dibuat_pada, id`
)

func (r LaporanRepository) AmbilUntukVerifikasi(ctx context.Context, id int64) (domain.LaporanUntukVerifikasi, bool, error) {
	var l domain.LaporanUntukVerifikasi
	err := r.DB.QueryRow(ctx, sqlLaporanUntukVerifikasi, id).Scan(&l.ID, &l.PelaporID, &l.Status, &l.Kondisi, &l.Titik.Lat, &l.Titik.Lng, &l.DibuatPada)
	l, ada, err := hasilSatuBaris(l, err)
	if !ada || err != nil {
		return l, ada, err
	}
	baris, err := r.DB.Query(ctx, sqlKonfirmasiLaporan, id)
	if err != nil {
		return l, false, err
	}
	l.Konfirmasi, err = pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.Konfirmasi, error) {
		var k domain.Konfirmasi
		return k, b.Scan(&k.PenggunaID, &k.Kondisi, &k.Titik.Lat, &k.Titik.Lng, &k.DibuatPada)
	})
	return l, true, err
}

func (r LaporanRepository) SimpanKonfirmasi(ctx context.Context, laporanID int64, k domain.Konfirmasi) error {
	_, err := r.DB.Exec(ctx, sqlSimpanKonfirmasi, laporanID, k.PenggunaID, k.Kondisi, k.Titik.Lat, k.Titik.Lng, k.DibuatPada)
	var kesalahanPG *pgconn.PgError
	if errors.As(err, &kesalahanPG) && kesalahanPG.Code == kodePelanggaranUnik {
		return domain.ErrKonfirmasiSudahTercatat
	}
	return err
}

func (r LaporanRepository) UbahStatus(ctx context.Context, p domain.PerubahanStatus) error {
	return dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, sqlUbahStatusLaporan, p.LaporanID, p.Dari, p.Ke)
		if err != nil {
			return fmt.Errorf("ubah status laporan %d: %w", p.LaporanID, err)
		}
		if tag.RowsAffected() == 0 {
			return domain.ErrTransisiTidakSah
		}
		_, err = tx.Exec(ctx, sqlCatatRiwayatLaporan, p.LaporanID, p.Dari, p.Ke, p.OlehPenggunaID, p.Alasan)
		return err
	})
}

func (r LaporanRepository) AmbilLaporanPublik(ctx context.Context, id int64) (domain.LaporanPublik, bool, error) {
	var l domain.LaporanPublik
	err := r.DB.QueryRow(ctx, sqlLaporanPublik, id).Scan(&l.ID, &l.Status, &l.Kondisi, &l.Catatan, &l.KodeWilayah, &l.NamaWilayah,
		&l.TitikPublik.Lat, &l.TitikPublik.Lng, &l.PaketKode, &l.PaketNama, &l.IsSimulation, &l.DibuatPada)
	l, ada, err := hasilSatuBaris(l, err)
	if !ada || err != nil {
		return l, ada, err
	}
	if l.Foto, err = kumpulkan(ctx, r.DB, sqlFotoLaporan, pgx.RowToStructByPos[domain.FotoTersimpan], id); err != nil {
		return l, false, fmt.Errorf("baca foto laporan %d: %w", id, err)
	}
	if l.Riwayat, err = kumpulkan(ctx, r.DB, sqlRiwayatLaporan, pgx.RowToStructByPos[domain.RiwayatLaporan], id); err != nil {
		return l, false, fmt.Errorf("baca riwayat laporan %d: %w", id, err)
	}
	return l, true, nil
}
