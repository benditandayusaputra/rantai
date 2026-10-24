package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlJumlahLaporanSejak = `SELECT count(*) FROM laporan WHERE pengguna_id = $1 AND dibuat_pada >= $2`

	sqlJarakPaket = `SELECT ST_Distance(lokasi::geography, ST_SetSRID(ST_MakePoint($3, $2), 4326)::geography)
FROM paket WHERE kode = $1 AND NOT is_batal`

	sqlWilayahDiTitik = `SELECT kode FROM wilayah
WHERE tingkat IN ('kelurahan', 'kecamatan') AND ST_Covers(geom, ST_SetSRID(ST_MakePoint($2, $1), 4326))
ORDER BY tingkat = 'kelurahan' DESC
LIMIT 1`

	sqlSimpanLaporan = `INSERT INTO laporan (pengguna_id, paket_kode, kode_wilayah, lokasi_publik, lokasi_asli, akurasi_meter,
    kondisi, catatan, status, is_simulation)
VALUES ($1, nullif($2, ''), $3, ST_SetSRID(ST_MakePoint($5, $4), 4326), ST_SetSRID(ST_MakePoint($7, $6), 4326), $8,
    $9, $10, 'MENUNGGU_SAKSI', $11 OR coalesce((SELECT is_simulation FROM paket WHERE kode = $2), false))
RETURNING id`

	sqlSimpanFotoLaporan = `INSERT INTO laporan_foto (laporan_id, urutan, object_key, thumbnail_key, lebar, tinggi)
SELECT $1, x.urutan, x.object_key, x.thumbnail_key, x.lebar, x.tinggi
FROM unnest($2::smallint[], $3::text[], $4::text[], $5::int[], $6::int[]) AS x (urutan, object_key, thumbnail_key, lebar, tinggi)`

	sqlRiwayatLaporanBaru = `INSERT INTO laporan_riwayat (laporan_id, status_lama, status_baru, oleh_pengguna_id)
VALUES ($1, NULL, 'DIKIRIM', $2), ($1, 'DIKIRIM', 'MENUNGGU_SAKSI', NULL)`
)

type LaporanRepository struct {
	DB Basisdata
}

func (r LaporanRepository) JumlahLaporanSejak(ctx context.Context, penggunaID int64, sejak time.Time) (int, error) {
	var jumlah int
	return jumlah, r.DB.QueryRow(ctx, sqlJumlahLaporanSejak, penggunaID, sejak).Scan(&jumlah)
}

func (r LaporanRepository) JarakPaketMeter(ctx context.Context, kode string, titik domain.Titik) (*float64, bool, error) {
	var jarak *float64
	return hasilSatuBaris(jarak, r.DB.QueryRow(ctx, sqlJarakPaket, kode, titik.Lat, titik.Lng).Scan(&jarak))
}

func (r LaporanRepository) WilayahDiTitik(ctx context.Context, titik domain.Titik) (string, bool, error) {
	var kode string
	return hasilSatuBaris(kode, r.DB.QueryRow(ctx, sqlWilayahDiTitik, titik.Lat, titik.Lng).Scan(&kode))
}

func (r LaporanRepository) SimpanLaporan(ctx context.Context, l domain.LaporanBaru) (int64, error) {
	var id int64
	err := dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, sqlSimpanLaporan, l.PenggunaID, l.PaketKode, l.KodeWilayah, l.TitikPublik.Lat, l.TitikPublik.Lng,
			l.Titik.Lat, l.Titik.Lng, l.AkurasiMeter, l.Kondisi, l.Catatan, l.IsSimulation).Scan(&id)
		if err != nil {
			return fmt.Errorf("simpan laporan: %w", err)
		}
		if err := simpanFotoLaporan(ctx, tx, id, l.Foto); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, sqlRiwayatLaporanBaru, id, l.PenggunaID)
		return err
	})
	return id, err
}

func simpanFotoLaporan(ctx context.Context, tx pgx.Tx, laporanID int64, foto []domain.FotoTersimpan) error {
	urutan, objek, thumbnail, lebar, tinggi := []int16{}, []string{}, []string{}, []int{}, []int{}
	for i, f := range foto {
		urutan, objek, thumbnail = append(urutan, int16(i+1)), append(objek, f.ObjectKey), append(thumbnail, f.ThumbnailKey)
		lebar, tinggi = append(lebar, f.Lebar), append(tinggi, f.Tinggi)
	}
	if _, err := tx.Exec(ctx, sqlSimpanFotoLaporan, laporanID, urutan, objek, thumbnail, lebar, tinggi); err != nil {
		return fmt.Errorf("simpan foto laporan %d: %w", laporanID, err)
	}
	return nil
}
