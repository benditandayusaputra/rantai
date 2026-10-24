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
	sqlDaftarFraksi = `SELECT DISTINCT fraksi FROM anggota ORDER BY fraksi`

	sqlSimpanSikapFraksi = `INSERT INTO sikap_fraksi (fraksi, agenda, judul_agenda, tanggal, ringkasan, url, provenance_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (fraksi, url) DO UPDATE SET
    agenda = EXCLUDED.agenda, judul_agenda = EXCLUDED.judul_agenda, tanggal = EXCLUDED.tanggal,
    ringkasan = EXCLUDED.ringkasan, provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (sikap_fraksi.agenda, sikap_fraksi.judul_agenda, sikap_fraksi.tanggal, sikap_fraksi.ringkasan)
    IS DISTINCT FROM (EXCLUDED.agenda, EXCLUDED.judul_agenda, EXCLUDED.tanggal, EXCLUDED.ringkasan)
RETURNING (xmax = 0) AS baru`

	sqlHapusSikapLama = `DELETE FROM sikap_fraksi
WHERE (fraksi, url) NOT IN (SELECT * FROM unnest($1::text[], $2::text[]))`

	sqlSimpanAspirasiReses = `INSERT INTO aspirasi_reses (masa_reses, tanggal_paripurna, bidang, uraian, kode_wilayah, url, provenance_id)
VALUES ($1, $2, $3, $4, nullif($5, ''), $6, $7)
ON CONFLICT (url, uraian) DO UPDATE SET
    masa_reses = EXCLUDED.masa_reses, tanggal_paripurna = EXCLUDED.tanggal_paripurna, bidang = EXCLUDED.bidang,
    kode_wilayah = EXCLUDED.kode_wilayah, provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (aspirasi_reses.masa_reses, aspirasi_reses.tanggal_paripurna, aspirasi_reses.bidang, aspirasi_reses.kode_wilayah)
    IS DISTINCT FROM (EXCLUDED.masa_reses, EXCLUDED.tanggal_paripurna, EXCLUDED.bidang, EXCLUDED.kode_wilayah)
RETURNING (xmax = 0) AS baru`

	sqlHapusAspirasiLama = `DELETE FROM aspirasi_reses
WHERE (url, uraian) NOT IN (SELECT * FROM unnest($1::text[], $2::text[]))`
)

type ParipurnaRepository struct {
	DB Basisdata
}

func (r ParipurnaRepository) DaftarFraksi(ctx context.Context) ([]string, error) {
	baris, err := r.DB.Query(ctx, sqlDaftarFraksi)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, pgx.RowTo[string])
}

func (r ParipurnaRepository) SimpanParipurna(ctx context.Context, data usecase.DataParipurna) (usecase.HasilImpor, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("mulai transaksi paripurna: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := simpanProvenance(ctx, tx, data.Provenance); err != nil {
		return usecase.HasilImpor{}, err
	}
	hasilSikap, err := simpanSikapFraksi(ctx, tx, data.Sikap, data.Provenance.ID)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("simpan sikap fraksi: %w", err)
	}
	hasilAspirasi, err := simpanAspirasiReses(ctx, tx, data.Aspirasi, data.Provenance.ID)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("simpan aspirasi reses: %w", err)
	}
	return usecase.HasilImpor{
		JumlahBaru:    hasilSikap.JumlahBaru + hasilAspirasi.JumlahBaru,
		JumlahBerubah: hasilSikap.JumlahBerubah + hasilAspirasi.JumlahBerubah,
	}, tx.Commit(ctx)
}

type kunciEntri struct {
	pertama []string
	kedua   []string
}

func (k *kunciEntri) tambah(pertama, kedua string) {
	k.pertama, k.kedua = append(k.pertama, pertama), append(k.kedua, kedua)
}

func simpanSikapFraksi(ctx context.Context, tx pgx.Tx, daftar []domain.SikapFraksi, provenanceID string) (usecase.HasilImpor, error) {
	hasil := usecase.HasilImpor{}
	kunci := kunciEntri{pertama: []string{}, kedua: []string{}}
	for _, s := range daftar {
		var baru bool
		err := tx.QueryRow(ctx, sqlSimpanSikapFraksi, s.Fraksi, s.Agenda, s.JudulAgenda, s.Tanggal, s.Ringkasan, s.URL, provenanceID).Scan(&baru)
		if err := catatUpsert(&hasil, baru, err); err != nil {
			return hasil, fmt.Errorf("%s %s: %w", s.Fraksi, s.URL, err)
		}
		kunci.tambah(s.Fraksi, s.URL)
	}
	terhapus, err := hapusEntriLama(ctx, tx, sqlHapusSikapLama, kunci)
	hasil.JumlahBerubah += terhapus
	return hasil, err
}

func simpanAspirasiReses(ctx context.Context, tx pgx.Tx, daftar []domain.AspirasiReses, provenanceID string) (usecase.HasilImpor, error) {
	hasil := usecase.HasilImpor{}
	kunci := kunciEntri{pertama: []string{}, kedua: []string{}}
	for _, a := range daftar {
		var baru bool
		err := tx.QueryRow(ctx, sqlSimpanAspirasiReses, a.MasaReses, a.TanggalParipurna, a.Bidang, a.Uraian, a.KodeWilayah, a.URL, provenanceID).Scan(&baru)
		if err := catatUpsert(&hasil, baru, err); err != nil {
			return hasil, fmt.Errorf("%s: %w", a.Uraian, err)
		}
		kunci.tambah(a.URL, a.Uraian)
	}
	terhapus, err := hapusEntriLama(ctx, tx, sqlHapusAspirasiLama, kunci)
	hasil.JumlahBerubah += terhapus
	return hasil, err
}

func catatUpsert(hasil *usecase.HasilImpor, baru bool, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	hasil.JumlahBaru += hitungJika(baru)
	hasil.JumlahBerubah += hitungJika(!baru)
	return nil
}

func hapusEntriLama(ctx context.Context, tx pgx.Tx, kueri string, kunci kunciEntri) (int, error) {
	tag, err := tx.Exec(ctx, kueri, kunci.pertama, kunci.kedua)
	if err != nil {
		return 0, fmt.Errorf("hapus entri lama: %w", err)
	}
	return int(tag.RowsAffected()), nil
}
