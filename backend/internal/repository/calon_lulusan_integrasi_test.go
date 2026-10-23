package repository_test

import (
	"time"

	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

type pelaporDiam struct{}

func (pelaporDiam) LaporkanKecamatanTanpaCalonLulusan([]string) {}

func siapkanSDUji(t *testing.T, tx pgx.Tx) {
	t.Helper()
	ujiwilayah.MuatDKI(t, tx)
	langkah := []string{
		`DELETE FROM sekolah_rombel`,
		`INSERT INTO sekolah (npsn, nama, bentuk, jenjang, status, is_madrasah, kode_kecamatan, provenance_id)
SELECT n, 'SD UJI ' || n, 'SD', 'dasar', 'negeri', false, '31.72.06', (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')
FROM unnest(ARRAY['90000001', '90000002']) AS n ON CONFLICT (npsn) DO NOTHING`,
		`INSERT INTO sekolah_rombel (npsn, tahun_ajaran, tingkat, jumlah_rombel, jumlah_siswa, provenance_id)
SELECT n, '2026/2027', t, 2, 60, (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')
FROM unnest(ARRAY['90000001', '90000002']) AS n, generate_series(1, 6) AS t`,
	}
	for _, sql := range langkah {
		if _, err := tx.Exec(t.Context(), sql); err != nil {
			t.Fatal(err)
		}
	}
}

func jalankanHitungCalon(t *testing.T, tx pgx.Tx) {
	t.Helper()
	repo := repository.CalonLulusanRepository{DB: tx}
	hitung := usecase.HitungCalonLulusan{
		Pembaca: repo, Penyimpan: repo, Pencatat: repository.IngestRunRepository{DB: tx}, Pelapor: pelaporDiam{},
		Konstanta: usecase.KonstantaPendidikan{PembagiEstimasiTotalSD: 5},
		BuatProvenance: func(c usecase.CatatanProvenance) domain.Provenance {
			return provenance.Baru(provenance.Catatan(c), nil, time.Now())
		},
	}
	if err := hitung.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestHitungCalonLulusan_DuaSDKelapaGading_JumlahKelasEnamResmi(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanSDUji(t, tx)

	jalankanHitungCalon(t, tx)

	var jumlah int
	var label, metode string
	err := tx.QueryRow(t.Context(), `SELECT jumlah, label::text, metode FROM calon_lulusan_kecamatan
WHERE kode_kecamatan = '31.72.06' AND tahun_ajaran = '2026/2027'`).Scan(&jumlah, &label, &metode)
	if err != nil || jumlah != 120 || label != "resmi" || metode != "per_sekolah" {
		t.Errorf("jumlah %d, label %q, metode %q, err %v", jumlah, label, metode, err)
	}
}
