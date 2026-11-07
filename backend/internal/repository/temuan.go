package repository

import (
	"context"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlIsPaketAda = `SELECT EXISTS (SELECT 1 FROM paket WHERE kode = $1)`

	sqlHitungTemuanSejak = `SELECT count(*) FROM temuan WHERE pengguna_id = $1 AND dibuat_pada >= $2`

	sqlSimpanTemuan = `INSERT INTO temuan (pengguna_id, paket_kode, isi, tautan_sumber) VALUES ($1, $2, $3, $4) RETURNING id`

	sqlPengirimTemuan = `SELECT pengguna_id FROM temuan WHERE id = $1 AND NOT disembunyikan`

	sqlSimpanTanggapan = `INSERT INTO temuan_tanggapan (temuan_id, pengguna_id, sikap, tautan_sumber) VALUES ($1, $2, $3, $4)
ON CONFLICT (temuan_id, pengguna_id) DO UPDATE SET sikap = EXCLUDED.sikap, tautan_sumber = EXCLUDED.tautan_sumber, diubah_pada = now()`

	sqlDaftarTemuan = `SELECT t.id, t.pengguna_id, t.paket_kode, t.isi, t.tautan_sumber, pg.nama_tampilan, t.dibuat_pada,
    count(tt.*) FILTER (WHERE tt.sikap = 'mendukung'), count(tt.*) FILTER (WHERE tt.sikap = 'membantah')
FROM temuan t
JOIN pengguna pg ON pg.id = t.pengguna_id
LEFT JOIN temuan_tanggapan tt ON tt.temuan_id = t.id
WHERE t.paket_kode = $1 AND NOT t.disembunyikan
GROUP BY t.id, pg.nama_tampilan
ORDER BY t.dibuat_pada DESC, t.id DESC
LIMIT $2`

	sqlSembunyikanTemuan = `UPDATE temuan SET disembunyikan = true, diubah_pada = now() WHERE id = $1`
)

type TemuanRepository struct {
	DB Basisdata
}

func (r TemuanRepository) IsPaketAda(ctx context.Context, kode string) (bool, error) {
	var ada bool
	return ada, r.DB.QueryRow(ctx, sqlIsPaketAda, kode).Scan(&ada)
}

func (r TemuanRepository) HitungTemuanSejak(ctx context.Context, penggunaID int64, sejak time.Time) (int, error) {
	var jumlah int
	return jumlah, r.DB.QueryRow(ctx, sqlHitungTemuanSejak, penggunaID, sejak).Scan(&jumlah)
}

func (r TemuanRepository) SimpanTemuan(ctx context.Context, m domain.TemuanMasuk) (int64, error) {
	var id int64
	return id, r.DB.QueryRow(ctx, sqlSimpanTemuan, m.PenggunaID, m.PaketKode, m.Isi, m.TautanSumber).Scan(&id)
}

func (r TemuanRepository) AmbilPengirimTemuan(ctx context.Context, id int64) (int64, bool, error) {
	var pengirim int64
	return hasilSatuBaris(pengirim, r.DB.QueryRow(ctx, sqlPengirimTemuan, id).Scan(&pengirim))
}

func (r TemuanRepository) SimpanTanggapan(ctx context.Context, m domain.TanggapanMasuk) error {
	_, err := r.DB.Exec(ctx, sqlSimpanTanggapan, m.TemuanID, m.PenggunaID, m.Sikap, m.TautanSumber)
	return err
}

func (r TemuanRepository) DaftarTemuan(ctx context.Context, paketKode string, batas int) ([]domain.Temuan, error) {
	return kumpulkan(ctx, r.DB, sqlDaftarTemuan, func(b pgx.CollectableRow) (domain.Temuan, error) {
		var t domain.Temuan
		return t, b.Scan(&t.ID, &t.PenggunaID, &t.PaketKode, &t.Isi, &t.TautanSumber, &t.NamaKontributor, &t.DibuatPada,
			&t.JumlahMendukung, &t.JumlahMembantah)
	}, paketKode, batas)
}

func (r TemuanRepository) SembunyikanTemuan(ctx context.Context, id int64) (bool, error) {
	tag, err := r.DB.Exec(ctx, sqlSembunyikanTemuan, id)
	return tag.RowsAffected() > 0, err
}

func (r TemuanRepository) CatatModerasi(ctx context.Context, c domain.CatatanModerasi) error {
	_, err := r.DB.Exec(ctx, sqlCatatModerasi, c.ModeratorID, c.TargetJenis, c.TargetID, c.Aksi, c.Alasan)
	return err
}
