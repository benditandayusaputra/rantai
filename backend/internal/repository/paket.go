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
	sqlSimpanInstansi = `INSERT INTO instansi (nama, jenis, kode_kota) VALUES ($1, $2, NULLIF($3, ''))
ON CONFLICT (nama) DO UPDATE SET jenis = EXCLUDED.jenis, kode_kota = EXCLUDED.kode_kota, diubah_pada = now()
RETURNING id`

	sqlSimpanPaketDaftar = `WITH lama AS (SELECT tahap FROM paket WHERE kode = $1)
INSERT INTO paket (kode, nama, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, hps, nilai_kontrak, tahap, provenance_id, kode_spse)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT (kode) DO UPDATE SET
    nama = CASE WHEN paket.detail_diambil_pada IS NULL THEN EXCLUDED.nama ELSE paket.nama END,
    hps = CASE WHEN paket.detail_diambil_pada IS NULL THEN EXCLUDED.hps ELSE paket.hps END,
    provenance_id = CASE WHEN paket.detail_diambil_pada IS NULL THEN EXCLUDED.provenance_id ELSE paket.provenance_id END,
    tahap = EXCLUDED.tahap,
    tahun_anggaran = GREATEST(paket.tahun_anggaran, EXCLUDED.tahun_anggaran),
    nilai_kontrak = COALESCE(EXCLUDED.nilai_kontrak, paket.nilai_kontrak),
    diubah_pada = now()
WHERE NOT paket.is_simulation AND (
    (paket.tahap, paket.nilai_kontrak, paket.tahun_anggaran)
        IS DISTINCT FROM (EXCLUDED.tahap, COALESCE(EXCLUDED.nilai_kontrak, paket.nilai_kontrak), GREATEST(paket.tahun_anggaran, EXCLUDED.tahun_anggaran))
    OR (paket.detail_diambil_pada IS NULL AND (paket.nama, paket.hps) IS DISTINCT FROM (EXCLUDED.nama, EXCLUDED.hps)))
RETURNING (xmax = 0) AS baru, (SELECT tahap FROM lama)`

	sqlCatatTahap = `INSERT INTO paket_riwayat_tahap (paket_kode, tahap_lama, tahap_baru, dicatat_pada) VALUES ($1, $2, $3, $4)`

	sqlPaketPerluDetail = `SELECT p.kode_spse FROM paket p
WHERE NOT p.is_tender AND NOT p.is_simulation AND p.jenis_pengadaan = ANY($1)
    AND (p.detail_diambil_pada IS NULL OR EXISTS (
        SELECT 1 FROM paket_riwayat_tahap r WHERE r.paket_kode = p.kode AND r.dicatat_pada > p.detail_diambil_pada))
ORDER BY p.tahun_anggaran DESC, p.kode`

	sqlSimpanRUPTercatat = `INSERT INTO rup (kode, nama, instansi_id, sumber_dana, tahun_anggaran, provenance_id)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (kode) DO NOTHING`

	sqlSimpanDetailPaket = `UPDATE paket SET
    nama = $2, instansi_id = $3, kode_rup = $4, is_batal = $5, tanggal_dibuat = $6, jenis_pengadaan = $7,
    metode = $8, jenis_kontrak = $9, tahun_anggaran = $10, pagu = $11, hps = $12, jumlah_peserta = $13,
    provenance_id = $14, detail_diambil_pada = $15, penyedia_id = $16,
    nilai_kontrak = coalesce($17, nilai_kontrak), tanggal_kontrak = $18, diubah_pada = now()
WHERE kode = $1 AND NOT is_simulation`

	sqlSimpanPenyedia = `INSERT INTO penyedia (nama, nama_normal, alias) VALUES ($1, $2, ARRAY[$1])
ON CONFLICT (nama_normal) DO UPDATE SET
    alias = CASE WHEN $1 = ANY (penyedia.alias) THEN penyedia.alias ELSE array_append(penyedia.alias, $1) END,
    diubah_pada = now()
RETURNING id`
)

type PaketRepository struct {
	DB Basisdata
}

func (r PaketRepository) SimpanDaftarPaket(ctx context.Context, data usecase.DaftarPaketSumber) (usecase.HasilImpor, error) {
	var hasil usecase.HasilImpor
	err := dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		var err error
		hasil, err = simpanDaftarPaket(ctx, tx, data)
		return err
	})
	return hasil, err
}

func simpanDaftarPaket(ctx context.Context, tx pgx.Tx, data usecase.DaftarPaketSumber) (usecase.HasilImpor, error) {
	hasil := usecase.HasilImpor{}
	if err := simpanProvenance(ctx, tx, data.Provenance); err != nil {
		return hasil, err
	}
	idInstansi := map[string]int64{}
	for _, p := range data.Paket {
		id, err := idInstansiDaftar(ctx, tx, idInstansi, p.NamaInstansi)
		if err != nil {
			return hasil, err
		}
		baru, berubah, err := simpanSatuPaketDaftar(ctx, tx, p, id, data.Provenance)
		if err != nil {
			return hasil, err
		}
		hasil.JumlahBaru += hitung(baru)
		hasil.JumlahBerubah += hitung(berubah)
	}
	return hasil, nil
}

func idInstansiDaftar(ctx context.Context, tx pgx.Tx, tersimpan map[string]int64, nama string) (int64, error) {
	if id, ada := tersimpan[nama]; ada {
		return id, nil
	}
	id, err := simpanInstansi(ctx, tx, domain.Instansi{Nama: nama, Jenis: domain.JenisInstansiLainnya})
	tersimpan[nama] = id
	return id, err
}

func simpanInstansi(ctx context.Context, tx pgx.Tx, i domain.Instansi) (int64, error) {
	var id int64
	if err := tx.QueryRow(ctx, sqlSimpanInstansi, i.Nama, i.Jenis, i.KodeKota).Scan(&id); err != nil {
		return 0, fmt.Errorf("simpan instansi %s: %w", i.Nama, err)
	}
	return id, nil
}

func simpanSatuPaketDaftar(ctx context.Context, tx pgx.Tx, p domain.PaketDaftar, idInstansi int64, prov domain.Provenance) (bool, bool, error) {
	var baru bool
	var tahapLama *string
	err := tx.QueryRow(ctx, sqlSimpanPaketDaftar, p.Kode, p.Nama, idInstansi, p.IsTender, p.JenisPengadaan,
		p.Metode, p.TahunAnggaran, p.HPS, p.NilaiKontrak, p.Tahap, prov.ID, p.KodeSPSE).Scan(&baru, &tahapLama)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, false, nil
	}
	if err != nil {
		return false, false, fmt.Errorf("simpan paket %s: %w", p.Kode, err)
	}
	if baru || *tahapLama != p.Tahap {
		if _, err := tx.Exec(ctx, sqlCatatTahap, p.Kode, tahapLama, p.Tahap, prov.DiambilPada); err != nil {
			return false, false, fmt.Errorf("catat tahap paket %s: %w", p.Kode, err)
		}
	}
	return baru, !baru, nil
}

func (r PaketRepository) DaftarPaketPerluDetail(ctx context.Context, jenisPengadaan []string) ([]string, error) {
	baris, err := r.DB.Query(ctx, sqlPaketPerluDetail, jenisPengadaan)
	if err != nil {
		return nil, fmt.Errorf("daftar paket perlu detail: %w", err)
	}
	return pgx.CollectRows(baris, pgx.RowTo[string])
}

func (r PaketRepository) SimpanDetailPaket(ctx context.Context, data usecase.DataSimpanDetailPaket) error {
	return dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		if err := simpanProvenance(ctx, tx, data.Provenance); err != nil {
			return err
		}
		idInstansi, err := simpanInstansi(ctx, tx, data.Instansi)
		if err != nil {
			return err
		}
		kodeRUP, err := simpanRUPTercatat(ctx, tx, data, idInstansi)
		if err != nil {
			return err
		}
		idPenyedia, err := simpanPenyedia(ctx, tx, data.Penyedia)
		if err != nil {
			return err
		}
		return simpanDetail(ctx, tx, data, relasiDetail{instansi: idInstansi, penyedia: idPenyedia, kodeRUP: kodeRUP})
	})
}

func simpanRUPTercatat(ctx context.Context, tx pgx.Tx, data usecase.DataSimpanDetailPaket, idInstansi int64) (*string, error) {
	rup := data.Detail.RUP
	if rup == nil {
		return nil, nil
	}
	_, err := tx.Exec(ctx, sqlSimpanRUPTercatat, rup.Kode, rup.Nama, idInstansi, rup.SumberDana, data.Detail.TahunAnggaran, data.Provenance.ID)
	if err != nil {
		return nil, fmt.Errorf("simpan rup %s: %w", rup.Kode, err)
	}
	return &rup.Kode, nil
}

type relasiDetail struct {
	instansi int64
	penyedia *int64
	kodeRUP  *string
}

func simpanPenyedia(ctx context.Context, tx pgx.Tx, p *domain.Penyedia) (*int64, error) {
	if p == nil {
		return nil, nil
	}
	var id int64
	if err := tx.QueryRow(ctx, sqlSimpanPenyedia, p.Nama, p.NamaNormal).Scan(&id); err != nil {
		return nil, fmt.Errorf("simpan penyedia %s: %w", p.Nama, err)
	}
	return &id, nil
}

func nilaiKontrak(p *domain.PemenangPaket) *float64 {
	if p == nil || !p.IsBerkontrak {
		return nil
	}
	return &p.Nilai
}

func simpanDetail(ctx context.Context, tx pgx.Tx, data usecase.DataSimpanDetailPaket, relasi relasiDetail) error {
	d := data.Detail
	_, err := tx.Exec(ctx, sqlSimpanDetailPaket, d.Kode, d.Nama, relasi.instansi, relasi.kodeRUP, d.IsBatal, tanggalAtauKosong(d.TanggalDibuat),
		d.JenisPengadaan, d.Metode, d.JenisKontrak, d.TahunAnggaran, d.Pagu, d.HPS, d.JumlahPeserta,
		data.Provenance.ID, data.Provenance.DiambilPada, relasi.penyedia, nilaiKontrak(d.Pemenang), d.TanggalKontrak)
	if err != nil {
		return fmt.Errorf("simpan detail paket %s: %w", d.Kode, err)
	}
	return nil
}

func dalamTransaksi(ctx context.Context, db Basisdata, kerja func(pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("mulai transaksi: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := kerja(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func tanggalAtauKosong(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func hitung(benar bool) int {
	if benar {
		return 1
	}
	return 0
}
