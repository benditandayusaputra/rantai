package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlBacaGazetir = `SELECT kode, tingkat, array_prepend(nama, alias),
    CASE WHEN tingkat = 'kelurahan' THEN left(induk_kode, 5) ELSE induk_kode END,
    CASE WHEN tingkat = 'kelurahan' THEN induk_kode ELSE '' END,
    ST_Y(titik_tengah), ST_X(titik_tengah)
FROM wilayah
WHERE tingkat IN ('kecamatan', 'kelurahan')
ORDER BY kode`

	sqlPaketUntukLokasi = `SELECT p.kode, p.nama, coalesce(i.kode_kota, '')
FROM paket p JOIN instansi i ON i.id = p.instansi_id
WHERE NOT p.is_simulation
ORDER BY p.kode`

	sqlSimpanLokasiPaket = `UPDATE paket SET lokasi_tingkat = $2, kode_wilayah = $3,
    lokasi = CASE WHEN $4::float8 IS NULL THEN NULL ELSE ST_SetSRID(ST_MakePoint($5, $4), 4326) END,
    rt = $6, rw = $7, diubah_pada = now()
WHERE kode = $1 AND (lokasi_tingkat, kode_wilayah, rt, rw) IS DISTINCT FROM ($2::tingkat_wilayah, $3, $6, $7)`

	sqlCakupanLokasi = `SELECT coalesce(lokasi_tingkat::text, 'belum'), count(*)
FROM paket
WHERE jenis_pengadaan = $1 AND NOT is_simulation
GROUP BY 1`
)

type LokasiPaketRepository struct {
	DB Basisdata
}

func (r LokasiPaketRepository) BacaGazetir(ctx context.Context) ([]domain.EntriGazetir, error) {
	baris, err := r.DB.Query(ctx, sqlBacaGazetir)
	if err != nil {
		return nil, fmt.Errorf("baca gazetir: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.EntriGazetir, error) {
		var e domain.EntriGazetir
		err := b.Scan(&e.Kode, &e.Tingkat, &e.Nama, &e.KodeKota, &e.KodeKecamatan, &e.TitikTengah.Lat, &e.TitikTengah.Lng)
		return e, err
	})
}

func (r LokasiPaketRepository) DaftarPaketUntukLokasi(ctx context.Context) ([]domain.MasukanLokasi, error) {
	baris, err := r.DB.Query(ctx, sqlPaketUntukLokasi)
	if err != nil {
		return nil, fmt.Errorf("daftar paket untuk lokasi: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.MasukanLokasi, error) {
		var m domain.MasukanLokasi
		err := b.Scan(&m.KodePaket, &m.NamaPaket, &m.KodeKotaInstansi)
		return m, err
	})
}

func (r LokasiPaketRepository) SimpanLokasiPaket(ctx context.Context, daftar []usecase.LokasiTerhitung) (int, error) {
	berubah := 0
	err := dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		for _, l := range daftar {
			n, err := simpanLokasi(ctx, tx, l)
			if err != nil {
				return err
			}
			berubah += n
		}
		return nil
	})
	return berubah, err
}

func simpanLokasi(ctx context.Context, tx pgx.Tx, l usecase.LokasiTerhitung) (int, error) {
	var lat, lng *float64
	if l.Lokasi.Titik != nil {
		lat, lng = &l.Lokasi.Titik.Lat, &l.Lokasi.Titik.Lng
	}
	tag, err := tx.Exec(ctx, sqlSimpanLokasiPaket, l.KodePaket, l.Lokasi.Tingkat, l.Lokasi.KodeWilayah, lat, lng, l.Lokasi.RT, l.Lokasi.RW)
	if err != nil {
		return 0, fmt.Errorf("simpan lokasi paket %s: %w", l.KodePaket, err)
	}
	return int(tag.RowsAffected()), nil
}

func (r LokasiPaketRepository) CakupanLokasi(ctx context.Context, jenisPengadaan string) (map[string]int, error) {
	baris, err := r.DB.Query(ctx, sqlCakupanLokasi, jenisPengadaan)
	if err != nil {
		return nil, fmt.Errorf("cakupan lokasi: %w", err)
	}
	defer baris.Close()
	cakupan := map[string]int{}
	for baris.Next() {
		var tingkat string
		var jumlah int
		if err := baris.Scan(&tingkat, &jumlah); err != nil {
			return nil, err
		}
		cakupan[tingkat] = jumlah
	}
	return cakupan, baris.Err()
}
