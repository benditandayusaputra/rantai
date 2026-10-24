package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlKejadianPaketBaru = `INSERT INTO kejadian (jenis, kode_kecamatan, judul, tautan, kunci, terjadi_pada, is_simulation)
SELECT 'paket_baru', substring(p.kode_wilayah from '^[0-9]+\.[0-9]+\.[0-9]+'), p.nama, '/paket/' || p.kode,
    'paket_baru:' || p.kode, p.dibuat_pada, p.is_simulation
FROM paket p
WHERE p.lokasi_tingkat IN ('kelurahan', 'kecamatan') AND p.dibuat_pada >= $1 AND NOT p.is_batal
ON CONFLICT (kunci) DO NOTHING`

	sqlKejadianMelewatiJadwal = `INSERT INTO kejadian (jenis, kode_kecamatan, judul, tautan, kunci, terjadi_pada, is_simulation)
SELECT 'melewati_jadwal', substring(p.kode_wilayah from '^[0-9]+\.[0-9]+\.[0-9]+'), p.nama, '/paket/' || p.kode,
    'melewati_jadwal:' || s.id, s.dihitung_pada, p.is_simulation
FROM sinyal s JOIN paket p ON p.kode = s.paket_kode
WHERE s.jenis = 'S-05' AND s.aktif AND s.dihitung_pada >= $1 AND p.lokasi_tingkat IN ('kelurahan', 'kecamatan')
ON CONFLICT (kunci) DO NOTHING`

	sqlKejadianLaporanTerverifikasi = `INSERT INTO kejadian (jenis, kode_kecamatan, judul, tautan, kunci, terjadi_pada, is_simulation)
SELECT 'laporan_terverifikasi', substring(l.kode_wilayah from '^[0-9]+\.[0-9]+\.[0-9]+'), coalesce(p.nama, w.nama),
    '/laporan/' || l.id, 'laporan_terverifikasi:' || l.id, r.dibuat_pada, l.is_simulation
FROM laporan_riwayat r
JOIN laporan l ON l.id = r.laporan_id
JOIN wilayah w ON w.kode = l.kode_wilayah
LEFT JOIN paket p ON p.kode = l.paket_kode
WHERE r.status_baru = 'TERVERIFIKASI' AND r.dibuat_pada >= $1
ON CONFLICT (kunci) DO NOTHING`

	sqlKejadianIndikator = `INSERT INTO kejadian (jenis, kode_kecamatan, judul, tautan, kunci, terjadi_pada)
SELECT 'indikator_diperbarui', i.kode_kecamatan, t.nama_id, '/kecamatan/' || i.kode_kecamatan,
    'indikator:' || i.tema_kode || ':' || i.kode_kecamatan || ':' || i.periode, i.dibuat_pada
FROM indikator_kecamatan i JOIN tema t ON t.kode = i.tema_kode
WHERE i.dibuat_pada >= $1
ON CONFLICT (kunci) DO NOTHING`

	sqlSebarNotifikasi = `INSERT INTO notifikasi (pengguna_id, kejadian_id)
SELECT l.pengguna_id, k.id
FROM kejadian k JOIN langganan l ON l.kode_wilayah = k.kode_kecamatan AND l.kanal = 'feed'
WHERE k.dibuat_pada >= $1 AND k.terjadi_pada >= l.dibuat_pada - $2::interval
ON CONFLICT (pengguna_id, kejadian_id) DO NOTHING`

	sqlNotifikasiAwalLangganan = `INSERT INTO notifikasi (pengguna_id, kejadian_id)
SELECT $1, k.id FROM kejadian k WHERE k.kode_kecamatan = $2 AND k.terjadi_pada >= now() - $3::interval
ON CONFLICT (pengguna_id, kejadian_id) DO NOTHING`

	sqlWilayahKecamatanAda = `SELECT EXISTS (SELECT 1 FROM wilayah WHERE kode = $1 AND tingkat = 'kecamatan')`

	sqlSimpanLangganan = `WITH baru AS (
    INSERT INTO langganan (pengguna_id, kode_wilayah, kanal) VALUES ($1, $2, 'feed')
    ON CONFLICT (pengguna_id, kode_wilayah, kanal) DO UPDATE SET kanal = EXCLUDED.kanal
    RETURNING id, kode_wilayah, dibuat_pada)
SELECT baru.id, baru.kode_wilayah, w.nama, baru.dibuat_pada FROM baru JOIN wilayah w ON w.kode = baru.kode_wilayah`

	sqlDaftarLangganan = `SELECT l.id, l.kode_wilayah, w.nama, l.dibuat_pada
FROM langganan l JOIN wilayah w ON w.kode = l.kode_wilayah
WHERE l.pengguna_id = $1 AND l.kanal = 'feed'
ORDER BY w.nama`

	sqlHapusLangganan = `DELETE FROM langganan WHERE id = $1 AND pengguna_id = $2`

	sqlDaftarNotifikasi = `SELECT n.id, k.jenis, k.kode_kecamatan, w.nama, k.judul, k.tautan, k.terjadi_pada, n.dibaca_pada IS NOT NULL,
    k.is_simulation
FROM notifikasi n
JOIN kejadian k ON k.id = n.kejadian_id
JOIN wilayah w ON w.kode = k.kode_kecamatan
WHERE n.pengguna_id = $1
ORDER BY k.terjadi_pada DESC, n.id DESC
LIMIT $2`

	sqlJumlahBelumDibaca = `SELECT count(*) FROM notifikasi WHERE pengguna_id = $1 AND dibaca_pada IS NULL`

	sqlTandaiDibaca = `UPDATE notifikasi SET dibaca_pada = coalesce(dibaca_pada, now()) WHERE id = $1 AND pengguna_id = $2`
)

var pembuatKejadian = []string{sqlKejadianPaketBaru, sqlKejadianMelewatiJadwal, sqlKejadianLaporanTerverifikasi, sqlKejadianIndikator}

type PantauRepository struct {
	DB Basisdata
}

func (r PantauRepository) BuatKejadian(ctx context.Context, sejak time.Time) (int, error) {
	baru := 0
	err := dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		for _, kueri := range pembuatKejadian {
			tag, err := tx.Exec(ctx, kueri, sejak)
			if err != nil {
				return fmt.Errorf("buat kejadian: %w", err)
			}
			baru += int(tag.RowsAffected())
		}
		_, err := tx.Exec(ctx, sqlSebarNotifikasi, sejak, domain.JendelaKejadian)
		return err
	})
	return baru, err
}

func (r PantauRepository) IsWilayahAda(ctx context.Context, kode string) (bool, error) {
	var ada bool
	return ada, r.DB.QueryRow(ctx, sqlWilayahKecamatanAda, kode).Scan(&ada)
}

func (r PantauRepository) SimpanLangganan(ctx context.Context, penggunaID int64, kodeKecamatan string) (domain.Langganan, error) {
	var l domain.Langganan
	err := dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		err := tx.QueryRow(ctx, sqlSimpanLangganan, penggunaID, kodeKecamatan).Scan(&l.ID, &l.KodeWilayah, &l.NamaWilayah, &l.DibuatPada)
		if err != nil {
			return fmt.Errorf("simpan langganan: %w", err)
		}
		_, err = tx.Exec(ctx, sqlNotifikasiAwalLangganan, penggunaID, kodeKecamatan, domain.JendelaKejadian)
		return err
	})
	return l, err
}

func (r PantauRepository) DaftarLangganan(ctx context.Context, penggunaID int64) ([]domain.Langganan, error) {
	return kumpulkan(ctx, r.DB, sqlDaftarLangganan, pgx.RowToStructByPos[domain.Langganan], penggunaID)
}

func (r PantauRepository) HapusLangganan(ctx context.Context, penggunaID, id int64) (bool, error) {
	tag, err := r.DB.Exec(ctx, sqlHapusLangganan, id, penggunaID)
	return tag.RowsAffected() > 0, err
}

func (r PantauRepository) DaftarNotifikasi(ctx context.Context, penggunaID int64, batas int) (domain.HalamanNotifikasi, error) {
	item, err := kumpulkan(ctx, r.DB, sqlDaftarNotifikasi, pgx.RowToStructByPos[domain.Notifikasi], penggunaID, batas)
	if err != nil {
		return domain.HalamanNotifikasi{}, err
	}
	halaman := domain.HalamanNotifikasi{Item: item}
	return halaman, r.DB.QueryRow(ctx, sqlJumlahBelumDibaca, penggunaID).Scan(&halaman.JumlahBelumDibaca)
}

func (r PantauRepository) TandaiDibaca(ctx context.Context, penggunaID, id int64) (bool, error) {
	tag, err := r.DB.Exec(ctx, sqlTandaiDibaca, id, penggunaID)
	return tag.RowsAffected() > 0, err
}
