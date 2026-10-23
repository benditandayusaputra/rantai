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
		Tema: repository.TemaRepository{DB: tx},
		Penghitung: map[string]usecase.PenghitungIndikator{domain.PenghitungKesenjanganKursiSMP: usecase.KesenjanganKursiSMP{
			Pembaca: repo, Konstanta: usecase.KonstantaPendidikan{BatasMuridPerRombelSMP: 30}, BuatProvenance: buatProvenanceUji,
		}},
		Penyimpan: repo, Pencatat: repository.IngestRunRepository{DB: tx},
	}
	if err := hitung.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
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
