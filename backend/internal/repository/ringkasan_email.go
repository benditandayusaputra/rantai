package repository

import (
	"context"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlRingkasanEmailAktif = `SELECT EXISTS (SELECT 1 FROM ringkasan_email WHERE pengguna_id = $1)`

	sqlAktifkanRingkasanEmail = `INSERT INTO ringkasan_email (pengguna_id) VALUES ($1) ON CONFLICT (pengguna_id) DO NOTHING`

	sqlMatikanRingkasanEmail = `DELETE FROM ringkasan_email WHERE pengguna_id = $1`

	sqlPenerimaRingkasan = `SELECT r.pengguna_id, p.email FROM ringkasan_email r
JOIN pengguna p ON p.id = r.pengguna_id
WHERE r.terakhir_dikirim IS NULL OR r.terakhir_dikirim <= $1
ORDER BY r.pengguna_id`

	sqlKabarRingkasan = `SELECT k.id, k.jenis, k.kode_kecamatan, w.nama, k.judul, k.tautan, k.terjadi_pada, false, k.is_simulation
FROM kejadian k
JOIN wilayah w ON w.kode = k.kode_kecamatan
WHERE k.terjadi_pada >= $2
    AND k.kode_kecamatan IN (SELECT kode_wilayah FROM langganan WHERE pengguna_id = $1 AND kanal = 'feed')
ORDER BY w.nama, k.terjadi_pada DESC, k.id DESC
LIMIT $3`

	sqlTandaiRingkasanTerkirim = `UPDATE ringkasan_email SET terakhir_dikirim = $2 WHERE pengguna_id = $1`
)

func (r PantauRepository) IsRingkasanEmailAktif(ctx context.Context, penggunaID int64) (bool, error) {
	var aktif bool
	err := r.DB.QueryRow(ctx, sqlRingkasanEmailAktif, penggunaID).Scan(&aktif)
	return aktif, err
}

func (r PantauRepository) AktifkanRingkasanEmail(ctx context.Context, penggunaID int64) error {
	_, err := r.DB.Exec(ctx, sqlAktifkanRingkasanEmail, penggunaID)
	return err
}

func (r PantauRepository) MatikanRingkasanEmail(ctx context.Context, penggunaID int64) error {
	_, err := r.DB.Exec(ctx, sqlMatikanRingkasanEmail, penggunaID)
	return err
}

func (r PantauRepository) PenerimaRingkasan(ctx context.Context, terkirimSebelum time.Time) ([]domain.PenerimaRingkasan, error) {
	return kumpulkan(ctx, r.DB, sqlPenerimaRingkasan, pgx.RowToStructByPos[domain.PenerimaRingkasan], terkirimSebelum)
}

func (r PantauRepository) KabarRingkasan(ctx context.Context, penggunaID int64, sejak time.Time, batas int) ([]domain.Notifikasi, error) {
	return kumpulkan(ctx, r.DB, sqlKabarRingkasan, pgx.RowToStructByPos[domain.Notifikasi], penggunaID, sejak, batas)
}

func (r PantauRepository) TandaiRingkasanTerkirim(ctx context.Context, penggunaID int64, pada time.Time) error {
	_, err := r.DB.Exec(ctx, sqlTandaiRingkasanTerkirim, penggunaID, pada)
	return err
}
