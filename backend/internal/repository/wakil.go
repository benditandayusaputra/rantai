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
	sqlDaftarKecamatan = `SELECT kode, nama FROM wilayah WHERE tingkat = 'kecamatan'`

	sqlSimpanDapil = `INSERT INTO dapil (nomor, nama, tahun_pemilu, jumlah_kursi, provenance_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (tahun_pemilu, nomor) DO UPDATE SET
    nama = EXCLUDED.nama, jumlah_kursi = EXCLUDED.jumlah_kursi, provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (dapil.nama, dapil.jumlah_kursi) IS DISTINCT FROM (EXCLUDED.nama, EXCLUDED.jumlah_kursi)
RETURNING (xmax = 0) AS baru`

	sqlHapusDapilWilayah = `DELETE FROM dapil_wilayah WHERE dapil_id IN (SELECT id FROM dapil WHERE tahun_pemilu = $1)`

	sqlSimpanDapilWilayah = `INSERT INTO dapil_wilayah (dapil_id, kode_kecamatan)
SELECT d.id, x.kode
FROM unnest($2::int[], $3::text[]) AS x (nomor, kode)
JOIN dapil d ON d.tahun_pemilu = $1 AND d.nomor = x.nomor`

	sqlSimpanAnggota = `INSERT INTO anggota (id, nama, nama_lain, partai, fraksi, dapil_id, url_profil, lhkpn_url, provenance_id)
SELECT $1, $2, $3, $4, $5, d.id, $8, $9, $10
FROM dapil d WHERE d.tahun_pemilu = $6 AND d.nomor = $7
ON CONFLICT (id) DO UPDATE SET
    nama = EXCLUDED.nama, nama_lain = EXCLUDED.nama_lain, partai = EXCLUDED.partai, fraksi = EXCLUDED.fraksi,
    dapil_id = EXCLUDED.dapil_id, url_profil = EXCLUDED.url_profil, lhkpn_url = EXCLUDED.lhkpn_url,
    provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE (anggota.nama, anggota.nama_lain, anggota.partai, anggota.fraksi, anggota.dapil_id, anggota.url_profil, anggota.lhkpn_url)
    IS DISTINCT FROM (EXCLUDED.nama, EXCLUDED.nama_lain, EXCLUDED.partai, EXCLUDED.fraksi, EXCLUDED.dapil_id, EXCLUDED.url_profil, EXCLUDED.lhkpn_url)
RETURNING (xmax = 0) AS baru`

	sqlHapusAnggotaLama = `DELETE FROM anggota WHERE NOT (id = ANY ($1))`
)

type WakilRepository struct {
	DB Basisdata
}

func (r WakilRepository) DaftarKecamatan(ctx context.Context) (map[string]string, error) {
	baris, err := r.DB.Query(ctx, sqlDaftarKecamatan)
	if err != nil {
		return nil, err
	}
	kecamatan := map[string]string{}
	var kode, nama string
	_, err = pgx.ForEachRow(baris, []any{&kode, &nama}, func() error {
		kecamatan[kode] = nama
		return nil
	})
	return kecamatan, err
}

func (r WakilRepository) SimpanWakil(ctx context.Context, data usecase.DataSimpanWakil) (usecase.HasilImpor, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("mulai transaksi wakil: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	hasil, err := simpanSemuaWakil(ctx, tx, data)
	if err != nil {
		return usecase.HasilImpor{}, err
	}
	return hasil, tx.Commit(ctx)
}

func simpanSemuaWakil(ctx context.Context, tx pgx.Tx, data usecase.DataSimpanWakil) (usecase.HasilImpor, error) {
	if err := errors.Join(simpanProvenance(ctx, tx, data.ProvenanceDapil), simpanProvenance(ctx, tx, data.ProvenanceAnggota)); err != nil {
		return usecase.HasilImpor{}, err
	}
	hasilDapil, err := simpanDapil(ctx, tx, data)
	if err != nil {
		return usecase.HasilImpor{}, err
	}
	hasilAnggota, err := simpanAnggota(ctx, tx, data)
	if err != nil {
		return usecase.HasilImpor{}, err
	}
	return usecase.HasilImpor{
		JumlahBaru:    hasilDapil.JumlahBaru + hasilAnggota.JumlahBaru,
		JumlahBerubah: hasilDapil.JumlahBerubah + hasilAnggota.JumlahBerubah,
	}, nil
}

func simpanDapil(ctx context.Context, tx pgx.Tx, data usecase.DataSimpanWakil) (usecase.HasilImpor, error) {
	hasil := usecase.HasilImpor{}
	for _, d := range data.Dapil {
		baru, berubah, err := tulisBaris(tx.QueryRow(ctx, sqlSimpanDapil, d.Nomor, d.Nama, data.TahunPemilu, d.JumlahKursi, data.ProvenanceDapil.ID))
		if err != nil {
			return usecase.HasilImpor{}, fmt.Errorf("simpan dapil %d: %w", d.Nomor, err)
		}
		hasil.JumlahBaru += hitungJika(baru)
		hasil.JumlahBerubah += hitungJika(berubah)
	}
	return hasil, simpanDapilWilayah(ctx, tx, data)
}

func simpanDapilWilayah(ctx context.Context, tx pgx.Tx, data usecase.DataSimpanWakil) error {
	nomor, kode := []int{}, []string{}
	for _, d := range data.Dapil {
		for _, k := range d.Kecamatan {
			nomor, kode = append(nomor, d.Nomor), append(kode, k.Kode)
		}
	}
	if _, err := tx.Exec(ctx, sqlHapusDapilWilayah, data.TahunPemilu); err != nil {
		return fmt.Errorf("hapus kecamatan dapil lama: %w", err)
	}
	if _, err := tx.Exec(ctx, sqlSimpanDapilWilayah, data.TahunPemilu, nomor, kode); err != nil {
		return fmt.Errorf("simpan kecamatan dapil: %w", err)
	}
	return nil
}

func simpanAnggota(ctx context.Context, tx pgx.Tx, data usecase.DataSimpanWakil) (usecase.HasilImpor, error) {
	hasil := usecase.HasilImpor{}
	id := []string{}
	for _, a := range data.Anggota {
		baru, berubah, err := simpanSatuAnggota(ctx, tx, a, data)
		if err != nil {
			return usecase.HasilImpor{}, fmt.Errorf("simpan anggota %s: %w", a.ID, err)
		}
		hasil.JumlahBaru += hitungJika(baru)
		hasil.JumlahBerubah += hitungJika(berubah)
		id = append(id, a.ID)
	}
	if _, err := tx.Exec(ctx, sqlHapusAnggotaLama, id); err != nil {
		return usecase.HasilImpor{}, fmt.Errorf("hapus anggota lama: %w", err)
	}
	return hasil, nil
}

func simpanSatuAnggota(ctx context.Context, tx pgx.Tx, a domain.Anggota, data usecase.DataSimpanWakil) (baru, berubah bool, err error) {
	return tulisBaris(tx.QueryRow(ctx, sqlSimpanAnggota,
		a.ID, a.Nama, a.NamaLain, a.Partai, a.Fraksi, data.TahunPemilu, a.NomorDapil,
		a.URLProfil, data.URLLHKPN, data.ProvenanceAnggota.ID))
}

func tulisBaris(baris pgx.Row) (baru, berubah bool, err error) {
	err = baris.Scan(&baru)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, false, nil
	}
	return baru, err == nil && !baru, err
}
