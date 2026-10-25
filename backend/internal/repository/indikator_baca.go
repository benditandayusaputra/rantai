package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlKolomIndikator = `i.tema_kode, i.kode_kecamatan, w.nama, i.periode, i.kebutuhan::int, i.kapasitas::int, i.kesenjangan::int,
    i.rasio::float8, i.kelas, i.label, i.rincian, i.provenance_id`

	sqlFiturPetaTema = `WITH terbaru AS (SELECT max(periode) AS periode FROM indikator_kecamatan WHERE tema_kode = $1)
SELECT ` + sqlKolomIndikator + `, ST_AsGeoJSON(ST_SimplifyPreserveTopology(w.geom, $2), $3)
FROM indikator_kecamatan i JOIN wilayah w ON w.kode = i.kode_kecamatan, terbaru
WHERE i.tema_kode = $1 AND i.periode = terbaru.periode
ORDER BY i.kode_kecamatan`

	sqlDetailPendidikanKecamatan = `WITH terbaru AS (SELECT max(periode) AS periode FROM indikator_kecamatan WHERE tema_kode = $2)
SELECT ` + sqlKolomIndikator + `,
    coalesce((SELECT c.provenance_id FROM calon_lulusan_kecamatan c
        WHERE c.kode_kecamatan = i.kode_kecamatan AND c.tahun_ajaran = i.periode AND c.jenjang = $3), '')
FROM indikator_kecamatan i JOIN wilayah w ON w.kode = i.kode_kecamatan, terbaru
WHERE i.tema_kode = $2 AND i.kode_kecamatan = $1 AND i.periode = terbaru.periode`

	sqlProvenanceKursiJenjang = `SELECT DISTINCT ON (d.jalur, s.bentuk) d.jalur, s.bentuk, d.provenance_id
FROM daya_tampung d JOIN sekolah s ON s.npsn = d.npsn
WHERE s.jenjang = $1
ORDER BY d.jalur, s.bentuk, d.tahun_ajaran DESC, d.diubah_pada DESC`

	sqlTrenRasioKecamatan = `SELECT periode, rasio::float8, label, provenance_id FROM indikator_kecamatan
WHERE kode_kecamatan = $1 AND tema_kode = $2 AND rasio IS NOT NULL
ORDER BY periode`

	sqlSekolahJenjang = `WITH tahun_kursi AS (SELECT max(tahun_ajaran) AS t FROM daya_tampung)
SELECT s.npsn, s.nama, s.bentuk, s.status, ST_Y(s.lokasi), ST_X(s.lokasi),
    coalesce(d.jalur, ''), d.kursi, coalesce(d.provenance_id, ''), coalesce(f.url, ''), coalesce(f.sumber, ''), coalesce(f.halaman, ''), coalesce(f.kredit, '')
FROM sekolah s
CROSS JOIN tahun_kursi
LEFT JOIN daya_tampung d ON d.npsn = s.npsn AND d.tahun_ajaran = tahun_kursi.t
LEFT JOIN foto_sekolah f ON f.npsn = s.npsn
WHERE s.kode_kecamatan = $1 AND s.jenjang = $2
ORDER BY d.kursi IS NULL, s.status, s.bentuk, s.nama`
)

type IndikatorBacaRepository struct {
	DB   Basisdata
	Peta config.Peta
}

func (r IndikatorBacaRepository) FiturPetaTema(ctx context.Context, temaKode string) ([]domain.FiturPetaIndikator, error) {
	baris, err := r.DB.Query(ctx, sqlFiturPetaTema, temaKode, r.Peta.ToleransiPenyederhanaan, r.Peta.DesimalKoordinat)
	if err != nil {
		return nil, fmt.Errorf("kueri peta tema: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.FiturPetaIndikator, error) {
		var f domain.FiturPetaIndikator
		err := b.Scan(append(tujuanIndikator(&f.IndikatorTersimpan), &f.GeometriGeoJSON)...)
		return f, err
	})
}

func (r IndikatorBacaRepository) DetailPendidikanKecamatan(ctx context.Context, p domain.Peralihan, kode string) (domain.DetailPendidikanKecamatan, bool, error) {
	var d domain.DetailPendidikanKecamatan
	tujuan := append(tujuanIndikator(&d.IndikatorTersimpan), &d.ProvenanceCalonLulusan)
	err := r.DB.QueryRow(ctx, sqlDetailPendidikanKecamatan, kode, p.TemaKode, p.JenjangAsal).Scan(tujuan...)
	d, ada, err := hasilSatuBaris(d, err)
	if err != nil || !ada {
		return d, ada, err
	}
	d.ProvenanceKursi, err = r.provenanceKursi(ctx, p.JenjangTujuan)
	return d, true, err
}

func (r IndikatorBacaRepository) provenanceKursi(ctx context.Context, jenjang domain.Jenjang) ([]domain.ProvenanceKursi, error) {
	baris, err := r.DB.Query(ctx, sqlProvenanceKursiJenjang, jenjang)
	if err != nil {
		return nil, fmt.Errorf("kueri provenance kursi %s: %w", jenjang, err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.ProvenanceKursi, error) {
		var p domain.ProvenanceKursi
		err := b.Scan(&p.Jalur, &p.Bentuk, &p.ProvenanceID)
		return p, err
	})
}

func (r IndikatorBacaRepository) TrenRasio(ctx context.Context, temaKode, kode string) ([]domain.TitikTren, error) {
	baris, err := r.DB.Query(ctx, sqlTrenRasioKecamatan, kode, temaKode)
	if err != nil {
		return nil, fmt.Errorf("kueri tren rasio %s: %w", kode, err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.TitikTren, error) {
		var titik domain.TitikTren
		err := b.Scan(&titik.Periode, &titik.Rasio, &titik.Label, &titik.ProvenanceID)
		return titik, err
	})
}

func (r IndikatorBacaRepository) SekolahJenjang(ctx context.Context, jenjang domain.Jenjang, kode string) ([]domain.SekolahDenganKursi, error) {
	baris, err := r.DB.Query(ctx, sqlSekolahJenjang, kode, jenjang)
	if err != nil {
		return nil, fmt.Errorf("kueri sekolah %s: %w", jenjang, err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.SekolahDenganKursi, error) {
		var s domain.SekolahDenganKursi
		var lat, lng *float64
		err := b.Scan(append([]any{&s.NPSN, &s.Nama, &s.Bentuk, &s.Status, &lat, &lng, &s.Jalur, &s.Kursi, &s.ProvenanceKursi},
			tujuanFoto(&s.Foto)...)...)
		if lat != nil && lng != nil {
			s.Lokasi = &domain.Titik{Lat: *lat, Lng: *lng}
		}
		return s, err
	})
}

func tujuanFoto(f *domain.Foto) []any {
	return []any{&f.URL, &f.Sumber, &f.Halaman, &f.Kredit}
}

func tujuanIndikator(i *domain.IndikatorTersimpan) []any {
	return []any{&i.TemaKode, &i.KodeKecamatan, &i.NamaKecamatan, &i.Periode, &i.Kebutuhan, &i.Kapasitas, &i.Kesenjangan,
		&i.Rasio, &i.Kelas, &i.Label, &i.Rincian, &i.ProvenanceID}
}
