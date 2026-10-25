package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlSekolahSekitar = `WITH titik AS (SELECT ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography AS g),
tahun_kursi AS (SELECT max(tahun_ajaran) AS t FROM daya_tampung)
SELECT s.npsn, s.nama, s.bentuk, s.status, ST_Y(s.lokasi), ST_X(s.lokasi), ST_Distance(s.lokasi::geography, titik.g),
    coalesce(d.jalur, ''), d.kursi, coalesce(d.provenance_id, ''), coalesce(f.url, ''), coalesce(f.sumber, ''), coalesce(f.halaman, ''), coalesce(f.kredit, '')
FROM sekolah s
CROSS JOIN titik
CROSS JOIN tahun_kursi
LEFT JOIN daya_tampung d ON d.npsn = s.npsn AND d.tahun_ajaran = tahun_kursi.t
LEFT JOIN foto_sekolah f ON f.npsn = s.npsn
WHERE s.jenjang = 'menengah_pertama' AND s.lokasi IS NOT NULL AND ST_DWithin(s.lokasi::geography, titik.g, $3)
ORDER BY 7, s.npsn`

	sqlSiswaKelas6Sekitar = `SELECT coalesce(sum(r.jumlah_siswa), 0)::int, coalesce(max(r.provenance_id), '')
FROM sekolah s
JOIN sekolah_rombel r ON r.npsn = s.npsn AND r.tingkat = $4
    AND r.tahun_ajaran = (SELECT max(tahun_ajaran) FROM sekolah_rombel)
WHERE s.jenjang = 'dasar' AND s.lokasi IS NOT NULL
    AND ST_DWithin(s.lokasi::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3)`
)

type SekolahTerdekatRepository struct {
	DB Basisdata
}

func (r SekolahTerdekatRepository) SekitarTitik(ctx context.Context, titik domain.Titik, radiusMeter int) (domain.SekitarTitik, error) {
	sekolah, err := r.sekolahSekitar(ctx, titik, radiusMeter)
	if err != nil {
		return domain.SekitarTitik{}, err
	}
	hasil := domain.SekitarTitik{Sekolah: sekolah}
	err = r.DB.QueryRow(ctx, sqlSiswaKelas6Sekitar, titik.Lng, titik.Lat, radiusMeter, domain.TingkatSDTerakhir).
		Scan(&hasil.SiswaKelas6, &hasil.ProvenanceSiswa)
	if err != nil {
		return domain.SekitarTitik{}, fmt.Errorf("kueri siswa kelas 6 sekitar: %w", err)
	}
	return hasil, nil
}

func (r SekolahTerdekatRepository) sekolahSekitar(ctx context.Context, titik domain.Titik, radiusMeter int) ([]domain.SekolahTerdekat, error) {
	baris, err := r.DB.Query(ctx, sqlSekolahSekitar, titik.Lng, titik.Lat, radiusMeter)
	if err != nil {
		return nil, fmt.Errorf("kueri sekolah sekitar: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.SekolahTerdekat, error) {
		var s domain.SekolahTerdekat
		var lat, lng float64
		err := b.Scan(append([]any{&s.NPSN, &s.Nama, &s.Bentuk, &s.Status, &lat, &lng, &s.JarakMeter, &s.Jalur, &s.Kursi, &s.ProvenanceKursi},
			tujuanFoto(&s.Foto)...)...)
		s.Lokasi = &domain.Titik{Lat: lat, Lng: lng}
		return s, err
	})
}
