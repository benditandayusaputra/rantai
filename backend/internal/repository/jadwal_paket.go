package repository

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const sqlAmbilJadwalPaket = `SELECT p.jenis_pengadaan, coalesce(p.kategori, ''), p.tanggal_serah_terima,
    coalesce(p.jadwal_selesai, r.jadwal_pelaksanaan_selesai), p.tanggal_kontrak
FROM paket p LEFT JOIN rup r ON r.kode = p.kode_rup
WHERE p.kode = $1`

func (r PaketRepository) AmbilJadwalPaket(ctx context.Context, kode string) (usecase.JadwalPaketTersimpan, bool, error) {
	var t usecase.JadwalPaketTersimpan
	err := r.DB.QueryRow(ctx, sqlAmbilJadwalPaket, kode).Scan(
		&t.JenisPengadaan, &t.Jadwal.Kategori, &t.Jadwal.TanggalSerahTerima, &t.Jadwal.JadwalSelesai, &t.Jadwal.TanggalKontrak,
	)
	return hasilSatuBaris(t, err)
}
