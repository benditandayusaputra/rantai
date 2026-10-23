package repository_test

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const sqlPaketLokasiUji = `INSERT INTO provenance (id, sumber, diambil_pada, label) VALUES ('prv_uji_lokasi', 'lpse', now(), 'resmi');
INSERT INTO instansi (nama, jenis, kode_kota) VALUES ('SUKU DINAS SUMBER DAYA AIR KOTA ADMINISTRASI JAKARTA UTARA', 'sudin', '31.72');
INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap, provenance_id)
SELECT k, k, n, (SELECT id FROM instansi), false, 'Pekerjaan Konstruksi', 'Pengadaan Langsung', 2026, 'Paket Sudah Selesai', 'prv_uji_lokasi'
FROM (VALUES
    ('nontender-1', 'Pembangunan Pintu Air Gg. Nusa Indah RT 04 RW 22 Gading Sengon (JPO Kali Betik), Kel. Kelapa Gading Barat, Kec. Kelapa Gading'),
    ('nontender-2', 'Pengadaan Alat Tulis Kantor')) AS v (k, n)`

func turunanUji(t *testing.T) (pgx.Tx, usecase.TurunanPaket) {
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	if _, err := tx.Exec(t.Context(), sqlPaketLokasiUji); err != nil {
		t.Fatal(err)
	}
	lokasi := repository.LokasiPaketRepository{DB: tx}
	return tx, usecase.TurunanPaket{
		Gazetir:  lokasi,
		Lokasi:   lokasi,
		Pencatat: repository.IngestRunRepository{DB: tx},
		AturanLokasi: domain.AturanLokasi{
			Singkatan:           map[string]string{"kel.": "kelurahan", "kec.": "kecamatan"},
			PergeseranMaksMeter: 150,
		},
	}
}

func TestTurunanPaket_NamaBerkelurahan_LokasiDanTitikTersimpan(t *testing.T) {
	tx, turunan := turunanUji(t)

	if err := turunan.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	var tingkat, nama string
	var rt, rw int
	var jarak float64
	err := tx.QueryRow(t.Context(), `SELECT p.lokasi_tingkat, w.nama, p.rt, p.rw, ST_Distance(p.lokasi::geography, w.titik_tengah::geography)
FROM paket p JOIN wilayah w ON w.kode = p.kode_wilayah WHERE p.kode = 'nontender-1'`).Scan(&tingkat, &nama, &rt, &rw, &jarak)
	if err != nil {
		t.Fatal(err)
	}
	if tingkat != "kelurahan" || nama != "Kelapa Gading Barat" || rt != 4 || rw != 22 || jarak > 150 {
		t.Errorf("tingkat %s, nama %s, rt %d, rw %d, jarak %.0f m", tingkat, nama, rt, rw, jarak)
	}
	var kodeKota string
	if err := tx.QueryRow(t.Context(), `SELECT kode_wilayah FROM paket WHERE kode = 'nontender-2' AND lokasi_tingkat = 'kota' AND lokasi IS NULL`).Scan(&kodeKota); err != nil || kodeKota != "31.72" {
		t.Errorf("paket tanpa wilayah: kota %q, err %v", kodeKota, err)
	}
}

func TestTurunanPaket_DijalankanDuaKali_TidakAdaPerubahanKedua(t *testing.T) {
	tx, turunan := turunanUji(t)
	for range 2 {
		if err := turunan.Jalankan(t.Context()); err != nil {
			t.Fatal(err)
		}
	}

	var berubah int
	if err := tx.QueryRow(t.Context(), `SELECT jumlah_berubah FROM ingest_run WHERE sumber = 'turunan' ORDER BY id DESC LIMIT 1`).Scan(&berubah); err != nil || berubah != 0 {
		t.Errorf("perubahan kedua = %d, err = %v", berubah, err)
	}
}
