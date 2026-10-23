package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlSimpanKategoriPaket = `UPDATE paket SET kategori = nullif($2, ''), sektor = nullif($3, ''), diubah_pada = now()
WHERE kode = $1 AND (kategori, sektor) IS DISTINCT FROM (nullif($2, ''), nullif($3, ''))`

	sqlPaketSinyal = `SELECT p.kode, p.metode, p.jenis_pengadaan, coalesce(p.kategori, ''), p.tahun_anggaran, p.instansi_id, p.penyedia_id,
    p.jumlah_peserta, p.hps, p.nilai_kontrak, p.tanggal_dibuat, coalesce(p.jadwal_selesai, r.jadwal_pelaksanaan_selesai)
FROM paket p LEFT JOIN rup r ON r.kode = p.kode_rup
WHERE NOT p.is_batal AND NOT p.is_simulation`

	sqlSimpanSinyal = `INSERT INTO sinyal (paket_kode, jenis, alasan, pembanding, aktif, dihitung_pada)
VALUES ($1, $2, $3, $4, true, $5)
ON CONFLICT (paket_kode, jenis) DO UPDATE SET
    alasan = EXCLUDED.alasan, pembanding = EXCLUDED.pembanding, aktif = true, dihitung_pada = EXCLUDED.dihitung_pada, diubah_pada = now()
WHERE (sinyal.alasan, sinyal.pembanding, sinyal.aktif) IS DISTINCT FROM (EXCLUDED.alasan, EXCLUDED.pembanding, true)`

	sqlNonaktifkanSinyal = `UPDATE sinyal SET aktif = false, diubah_pada = now()
WHERE aktif AND NOT EXISTS (
    SELECT 1 FROM unnest($1::text[], $2::jenis_sinyal[]) AS t (kode, jenis) WHERE t.kode = sinyal.paket_kode AND t.jenis = sinyal.jenis)`
)

func (r LokasiPaketRepository) SimpanKategoriPaket(ctx context.Context, kategoriPerPaket map[string]domain.Kategori) (int, error) {
	berubah := 0
	err := dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		for kode, k := range kategoriPerPaket {
			tag, err := tx.Exec(ctx, sqlSimpanKategoriPaket, kode, k.Kategori, k.Sektor)
			if err != nil {
				return fmt.Errorf("simpan kategori paket %s: %w", kode, err)
			}
			berubah += int(tag.RowsAffected())
		}
		return nil
	})
	return berubah, err
}

func (r LokasiPaketRepository) DaftarPaketSinyal(ctx context.Context) ([]domain.PaketSinyal, error) {
	baris, err := r.DB.Query(ctx, sqlPaketSinyal)
	if err != nil {
		return nil, fmt.Errorf("daftar paket untuk sinyal: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.PaketSinyal, error) {
		var p domain.PaketSinyal
		err := b.Scan(&p.Kode, &p.Metode, &p.JenisPengadaan, &p.Kategori, &p.TahunAnggaran, &p.InstansiID, &p.PenyediaID,
			&p.JumlahPeserta, &p.HPS, &p.NilaiKontrak, &p.TanggalDibuat, &p.JadwalSelesai)
		return p, err
	})
}

func (r LokasiPaketRepository) SimpanSinyal(ctx context.Context, daftar []domain.Sinyal, dihitungPada time.Time) (int, error) {
	berubah := 0
	err := dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		kode, jenis := make([]string, len(daftar)), make([]string, len(daftar))
		for i, s := range daftar {
			n, err := simpanSatuSinyal(ctx, tx, s, dihitungPada)
			if err != nil {
				return err
			}
			berubah += n
			kode[i], jenis[i] = s.KodePaket, string(s.Jenis)
		}
		tag, err := tx.Exec(ctx, sqlNonaktifkanSinyal, kode, jenis)
		if err != nil {
			return fmt.Errorf("nonaktifkan sinyal lama: %w", err)
		}
		berubah += int(tag.RowsAffected())
		return nil
	})
	return berubah, err
}

func simpanSatuSinyal(ctx context.Context, tx pgx.Tx, s domain.Sinyal, dihitungPada time.Time) (int, error) {
	pembanding, err := json.Marshal(s.Pembanding)
	if err != nil {
		return 0, fmt.Errorf("serialisasi pembanding sinyal %s: %w", s.KodePaket, err)
	}
	tag, err := tx.Exec(ctx, sqlSimpanSinyal, s.KodePaket, s.Jenis, s.Alasan, pembanding, dihitungPada)
	if err != nil {
		return 0, fmt.Errorf("simpan sinyal %s %s: %w", s.KodePaket, s.Jenis, err)
	}
	return int(tag.RowsAffected()), nil
}
