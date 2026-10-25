package repository_test

import (
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/tema"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

func buatProvenanceUji(c usecase.CatatanProvenance) domain.Provenance {
	return provenance.Baru(provenance.Catatan(c), nil, time.Now())
}

func siapkanKursiUji(t *testing.T, tx pgx.Tx) {
	t.Helper()
	siapkanSDUji(t, tx)
	langkah := []string{
		`DELETE FROM daya_tampung`,
		`DELETE FROM calon_lulusan_kecamatan`,
		`INSERT INTO sekolah (npsn, nama, bentuk, jenjang, status, is_madrasah, kode_kecamatan, provenance_id)
VALUES ('90000101', 'SMP UJI NEGERI', 'SMP', 'menengah_pertama', 'negeri', false, '31.72.06', (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')),
       ('90000102', 'SMP UJI SWASTA', 'SMP', 'menengah_pertama', 'swasta', false, '31.72.06', (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06'))
ON CONFLICT (npsn) DO NOTHING`,
		`INSERT INTO daya_tampung (npsn, tahun_ajaran, jalur, kursi, provenance_id)
VALUES ('90000101', '2026/2027', 'negeri', 90, (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')),
       ('90000102', '2026/2027', 'spmb_bersama', 12, (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06'))`,
		`INSERT INTO sekolah_rombel (npsn, tahun_ajaran, tingkat, jumlah_rombel, jumlah_siswa, provenance_id)
SELECT n, '2026/2027', t, 3, 100, (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')
FROM unnest(ARRAY['90000101', '90000102']) AS n, generate_series(7, 9) AS t`,
	}
	for _, sql := range langkah {
		if _, err := tx.Exec(t.Context(), sql); err != nil {
			t.Fatal(err)
		}
	}
	jalankanHitungCalon(t, tx)
}

func jalankanHitungIndikator(t *testing.T, tx pgx.Tx) {
	t.Helper()
	imporTema := usecase.ImporTema{
		Pengambil: tema.Sumber{Dir: "../../config/tema", URLRujukan: "https://contoh.test", Sekarang: time.Now},
		Penyimpan: repository.TemaRepository{DB: tx}, Pencatat: repository.IngestRunRepository{DB: tx},
	}
	if err := imporTema.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	repo := repository.IndikatorRepository{DB: tx}
	hitung := usecase.HitungIndikator{
		Tema:       repository.TemaRepository{DB: tx},
		Penghitung: penghitungUji(repo),
		Penyimpan:  repo, Pencatat: repository.IngestRunRepository{DB: tx},
	}
	if err := hitung.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func penghitungUji(repo repository.IndikatorRepository) map[string]usecase.PenghitungIndikator {
	hasil := map[string]usecase.PenghitungIndikator{}
	for _, peralihan := range domain.DaftarPeralihan() {
		hasil[peralihan.Penghitung] = usecase.KesenjanganKursi{
			Peralihan: peralihan, Pembaca: repo, Konstanta: konstantaPendidikanUji, BuatProvenance: buatProvenanceUji,
		}
	}
	return hasil
}

func siapkanSMASMKUji(t *testing.T, tx pgx.Tx) {
	t.Helper()
	langkah := []string{
		`INSERT INTO sekolah (npsn, nama, bentuk, jenjang, status, is_madrasah, kode_kecamatan, provenance_id)
VALUES ('90000111', 'SMA UJI NEGERI', 'SMA', 'menengah_atas', 'negeri', false, '31.72.06', (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')),
       ('90000112', 'SMK UJI NEGERI', 'SMK', 'menengah_atas', 'negeri', false, '31.72.06', (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')),
       ('90000113', 'SMK UJI SWASTA', 'SMK', 'menengah_atas', 'swasta', false, '31.72.06', (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')),
       ('90000114', 'MA UJI', 'MA', 'menengah_atas', 'swasta', true, '31.72.06', (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06'))
ON CONFLICT (npsn) DO NOTHING`,
		`INSERT INTO daya_tampung (npsn, tahun_ajaran, jalur, kursi, provenance_id)
VALUES ('90000111', '2026/2027', 'negeri', 72, (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')),
       ('90000112', '2026/2027', 'negeri', 36, (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')),
       ('90000113', '2026/2027', 'spmb_bersama', 8, (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06'))`,
	}
	for _, sql := range langkah {
		if _, err := tx.Exec(t.Context(), sql); err != nil {
			t.Fatal(err)
		}
	}
}

func TestHitungIndikator_SMASMK_LulusanKelas9DibandingKursiSMADanSMKNegeri(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanKursiUji(t, tx)
	siapkanSMASMKUji(t, tx)

	jalankanHitungIndikator(t, tx)

	var kebutuhan, kapasitas, kesenjangan, kursiSMK, kursiSPMB int
	var kelas string
	var madrasah *int
	err := tx.QueryRow(t.Context(), `SELECT kebutuhan, kapasitas, kesenjangan, kelas,
    (rincian->'kursi_negeri_per_bentuk'->1->>'kursi')::int, (rincian->'kursi_spmb_bersama'->>'nilai')::int,
    (rincian->'kursi_madrasah'->>'nilai')::int
FROM indikator_kecamatan WHERE tema_kode = 'pendidikan_sma_smk' AND kode_kecamatan = '31.72.06' AND periode = '2026/2027'`).
		Scan(&kebutuhan, &kapasitas, &kesenjangan, &kelas, &kursiSMK, &kursiSPMB, &madrasah)
	if err != nil || kebutuhan != 200 || kapasitas != 108 || kesenjangan != 92 || kelas != "kurang" || kursiSMK != 36 || kursiSPMB != 8 || madrasah != nil {
		t.Errorf("kebutuhan %d, kapasitas %d, kesenjangan %d, kelas %s, smk %d, spmb %d, madrasah %v, err %v",
			kebutuhan, kapasitas, kesenjangan, kelas, kursiSMK, kursiSPMB, madrasah, err)
	}
}

func TestHitungIndikator_KelapaGading_KesenjanganDariCalonLulusanDanKursiNegeri(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanKursiUji(t, tx)

	jalankanHitungIndikator(t, tx)

	var kebutuhan, kapasitas, kesenjangan, kursiSPMB int
	var rasio float64
	var kelas, label string
	err := tx.QueryRow(t.Context(), `SELECT kebutuhan, kapasitas, kesenjangan, rasio, kelas, label::text,
    (rincian->'kursi_spmb_bersama'->>'nilai')::int
FROM indikator_kecamatan WHERE tema_kode = 'pendidikan' AND kode_kecamatan = '31.72.06' AND periode = '2026/2027'`).
		Scan(&kebutuhan, &kapasitas, &kesenjangan, &rasio, &kelas, &label, &kursiSPMB)
	if err != nil || kebutuhan != 120 || kapasitas != 90 || kesenjangan != 30 || rasio != 0.75 || kelas != "kurang" || label != "resmi" || kursiSPMB != 12 {
		t.Errorf("kebutuhan %d, kapasitas %d, kesenjangan %d, rasio %v, kelas %s, label %s, spmb %d, err %v",
			kebutuhan, kapasitas, kesenjangan, rasio, kelas, label, kursiSPMB, err)
	}
}

func TestHitungIndikator_SeluruhKecamatan_TersimpanTermasukTanpaData(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanKursiUji(t, tx)

	jalankanHitungIndikator(t, tx)

	var jumlah, tanpaData int
	err := tx.QueryRow(t.Context(), `SELECT count(*), count(*) FILTER (WHERE kelas = 'tidak_ada_data')
FROM indikator_kecamatan WHERE tema_kode = 'pendidikan' AND periode = '2026/2027'`).Scan(&jumlah, &tanpaData)
	if err != nil || jumlah != 44 || tanpaData == 0 {
		t.Errorf("jumlah %d, tanpa data %d, err %v", jumlah, tanpaData, err)
	}
}
