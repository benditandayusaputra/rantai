package repository_test

import (
	"slices"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
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
	var anggaran usecase.AturanTurunanAnggaran
	if err := config.BacaYAML("../../config/kategori.yaml", &anggaran.Kategori); err != nil {
		t.Fatal(err)
	}
	if err := config.BacaYAML("../../config/sinyal.yaml", &anggaran.Sinyal); err != nil {
		t.Fatal(err)
	}
	return tx, usecase.TurunanPaket{
		Gazetir:        lokasi,
		Lokasi:         lokasi,
		Tema:           repository.TemaRepository{DB: tx},
		TemaPaket:      lokasi,
		KategoriSinyal: lokasi,
		Anggaran:       anggaran,
		Sekarang:       time.Now,
		Pencatat:       repository.IngestRunRepository{DB: tx},
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

func simpanTemaUji(t *testing.T, tx pgx.Tx) {
	t.Helper()
	tema := domain.Tema{Kode: "pendidikan", Nama: domain.Teks{ID: "Pendidikan", EN: "Education"}, Aktif: true, Urutan: 1}
	tema.FilterAnggaran = domain.FilterAnggaran{InstansiPola: []string{"dinas pendidikan"}, KataKunci: []string{"pintu air"}}
	prov := domain.Provenance{ID: "prv_uji_tema_paket", Sumber: domain.SumberRantai, Label: domain.LabelResmi, DiambilPada: time.Now()}
	if _, err := (repository.TemaRepository{DB: tx}).SimpanTema(t.Context(), usecase.DataSumberTema{Tema: []domain.Tema{tema}, Provenance: prov}); err != nil {
		t.Fatal(err)
	}
}

func TestTurunanPaket_FilterAnggaranTema_TemaKodeTerisi(t *testing.T) {
	tx, turunan := turunanUji(t)
	simpanTemaUji(t, tx)

	if err := turunan.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	var temaPintuAir, temaATK []string
	err := tx.QueryRow(t.Context(), `SELECT (SELECT tema_kode FROM paket WHERE kode = 'nontender-1'), (SELECT tema_kode FROM paket WHERE kode = 'nontender-2')`).Scan(&temaPintuAir, &temaATK)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(temaPintuAir, []string{"pendidikan"}) || len(temaATK) != 0 {
		t.Errorf("tema pintu air %v, tema ATK %v", temaPintuAir, temaATK)
	}
}

const sqlPaketSinyalUji = `INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap,
    hps, tanggal_dibuat, provenance_id)
SELECT 'nontender-pl' || n, 'pl' || n, 'Pembangunan Saluran RW ' || n, (SELECT id FROM instansi LIMIT 1), false,
    'Pekerjaan Konstruksi', 'Pengadaan Langsung', 2026, 'Paket Sudah Selesai', 390000000 + n, DATE '2026-03-01' + n, 'prv_uji_lokasi'
FROM generate_series(1, 3) AS n`

func TestTurunanPaket_TigaPaketMepetBatas_KategoriDanSinyalS03(t *testing.T) {
	tx, turunan := turunanUji(t)
	if _, err := tx.Exec(t.Context(), sqlPaketSinyalUji); err != nil {
		t.Fatal(err)
	}

	for range 2 {
		if err := turunan.Jalankan(t.Context()); err != nil {
			t.Fatal(err)
		}
	}

	var jumlah int
	var kategori string
	err := tx.QueryRow(t.Context(), `SELECT count(*), min(p.kategori) FROM sinyal s JOIN paket p ON p.kode = s.paket_kode
WHERE s.jenis = 'S-03' AND s.aktif`).Scan(&jumlah, &kategori)
	if err != nil || jumlah != 3 || kategori != "saluran" {
		t.Errorf("sinyal S-03 = %d, kategori %q, err %v", jumlah, kategori, err)
	}
	if _, err := tx.Exec(t.Context(), `UPDATE paket SET hps = 100000000 WHERE kode = 'nontender-pl1'`); err != nil {
		t.Fatal(err)
	}
	if err := turunan.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := tx.QueryRow(t.Context(), `SELECT count(*) FROM sinyal WHERE aktif`).Scan(&jumlah); err != nil || jumlah != 0 {
		t.Errorf("sinyal aktif setelah pola hilang = %d, err %v", jumlah, err)
	}
}
