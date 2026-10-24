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
	kolomPengguna = `id, email, nama_tampilan, peran, instansi_id, anggota_id`

	sqlHitungOTPSejak = `SELECT count(*) FROM otp WHERE email = $1 AND dibuat_pada >= $2`

	sqlSimpanOTP = `INSERT INTO otp (email, kode_hash, kedaluwarsa, dibuat_pada) VALUES ($1, $2, $3, $4)`

	sqlCatatPercobaanOTP = `UPDATE otp SET percobaan = percobaan + 1
WHERE id = (SELECT id FROM otp WHERE email = $1 ORDER BY dibuat_pada DESC, id DESC LIMIT 1)
RETURNING id, kode_hash, kedaluwarsa, percobaan, dipakai`

	sqlTandaiOTPDipakai = `UPDATE otp SET dipakai = true WHERE id = $1 AND NOT dipakai`

	sqlPastikanPengguna = `INSERT INTO pengguna (email) VALUES ($1)
ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
RETURNING ` + kolomPengguna

	sqlAmbilPengguna = `SELECT ` + kolomPengguna + ` FROM pengguna WHERE id = $1`

	sqlAturPeran = `INSERT INTO pengguna (email, peran) VALUES ($1, $2)
ON CONFLICT (email) DO UPDATE SET peran = EXCLUDED.peran, instansi_id = NULL, anggota_id = NULL, diubah_pada = now()`

	sqlSimpanSesi = `INSERT INTO sesi_refresh (pengguna_id, token_hash, kedaluwarsa) VALUES ($1, $2, $3)`

	sqlCabutSesi = `WITH lama AS (
    SELECT id, pengguna_id, kedaluwarsa, dicabut FROM sesi_refresh WHERE token_hash = $1 FOR UPDATE)
UPDATE sesi_refresh s SET dicabut = true FROM lama WHERE s.id = lama.id
RETURNING lama.pengguna_id, lama.kedaluwarsa, lama.dicabut`

	sqlCabutSemuaSesi = `UPDATE sesi_refresh SET dicabut = true WHERE pengguna_id = $1 AND NOT dicabut`
)

type PenggunaRepository struct {
	DB Basisdata
}

func (r PenggunaRepository) HitungOTPSejak(ctx context.Context, email string, sejak time.Time) (int, error) {
	var jumlah int
	if err := r.DB.QueryRow(ctx, sqlHitungOTPSejak, email, sejak).Scan(&jumlah); err != nil {
		return 0, fmt.Errorf("hitung permintaan kode masuk: %w", err)
	}
	return jumlah, nil
}

func (r PenggunaRepository) SimpanOTP(ctx context.Context, otp usecase.OTPBaru) error {
	if _, err := r.DB.Exec(ctx, sqlSimpanOTP, otp.Email, otp.KodeHash, otp.Kedaluwarsa, otp.DibuatPada); err != nil {
		return fmt.Errorf("simpan kode masuk: %w", err)
	}
	return nil
}

func (r PenggunaRepository) CatatPercobaanOTP(ctx context.Context, email string) (domain.OTPTersimpan, error) {
	var o domain.OTPTersimpan
	err := r.DB.QueryRow(ctx, sqlCatatPercobaanOTP, email).Scan(&o.ID, &o.KodeHash, &o.Kedaluwarsa, &o.Percobaan, &o.IsDipakai)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.OTPTersimpan{}, domain.ErrOTPTidakBerlaku
	}
	if err != nil {
		return domain.OTPTersimpan{}, fmt.Errorf("catat percobaan kode masuk: %w", err)
	}
	return o, nil
}

func (r PenggunaRepository) TandaiOTPDipakai(ctx context.Context, id int64) error {
	tag, err := r.DB.Exec(ctx, sqlTandaiOTPDipakai, id)
	if err != nil {
		return fmt.Errorf("tandai kode masuk %d dipakai: %w", id, err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrOTPTidakBerlaku
	}
	return nil
}

func (r PenggunaRepository) PastikanPengguna(ctx context.Context, email string) (domain.Pengguna, error) {
	p, err := pindaiPengguna(r.DB.QueryRow(ctx, sqlPastikanPengguna, email))
	if err != nil {
		return domain.Pengguna{}, fmt.Errorf("pastikan pengguna: %w", err)
	}
	return p, nil
}

func (r PenggunaRepository) AmbilPengguna(ctx context.Context, id int64) (domain.Pengguna, error) {
	p, err := pindaiPengguna(r.DB.QueryRow(ctx, sqlAmbilPengguna, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Pengguna{}, domain.ErrPenggunaTidakDitemukan
	}
	if err != nil {
		return domain.Pengguna{}, fmt.Errorf("ambil pengguna %d: %w", id, err)
	}
	return p, nil
}

func (r PenggunaRepository) AturPeran(ctx context.Context, email string, peran domain.Peran) error {
	if _, err := r.DB.Exec(ctx, sqlAturPeran, email, peran); err != nil {
		return fmt.Errorf("atur peran %s: %w", peran, err)
	}
	return nil
}

func (r PenggunaRepository) SimpanSesi(ctx context.Context, sesi usecase.SesiBaru) error {
	if _, err := r.DB.Exec(ctx, sqlSimpanSesi, sesi.PenggunaID, sesi.TokenHash, sesi.Kedaluwarsa); err != nil {
		return fmt.Errorf("simpan sesi: %w", err)
	}
	return nil
}

func (r PenggunaRepository) CabutSesi(ctx context.Context, tokenHash []byte) (domain.SesiRefresh, error) {
	var s domain.SesiRefresh
	err := r.DB.QueryRow(ctx, sqlCabutSesi, tokenHash).Scan(&s.PenggunaID, &s.Kedaluwarsa, &s.IsDicabut)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.SesiRefresh{}, domain.ErrSesiTidakBerlaku
	}
	if err != nil {
		return domain.SesiRefresh{}, fmt.Errorf("cabut sesi: %w", err)
	}
	return s, nil
}

func (r PenggunaRepository) CabutSemuaSesi(ctx context.Context, penggunaID int64) error {
	if _, err := r.DB.Exec(ctx, sqlCabutSemuaSesi, penggunaID); err != nil {
		return fmt.Errorf("cabut semua sesi pengguna %d: %w", penggunaID, err)
	}
	return nil
}

func pindaiPengguna(baris pgx.Row) (domain.Pengguna, error) {
	var p domain.Pengguna
	err := baris.Scan(&p.ID, &p.Email, &p.NamaTampilan, &p.Peran, &p.InstansiID, &p.AnggotaID)
	return p, err
}
