package repository_test

import (
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/jackc/pgx/v5"
)

const sqlPaketRelasiUji = `INSERT INTO provenance (id, sumber, diambil_pada, label) VALUES ('prv_uji_relasi', 'lpse', now(), 'resmi');
INSERT INTO instansi (nama, jenis) VALUES ('Sudin Uji Relasi', 'sudin');
INSERT INTO penyedia (nama, nama_normal) VALUES ('CV Uji Relasi', 'cv uji relasi');
INSERT INTO paket (kode, kode_spse, nama, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap, provenance_id,
    kode_wilayah, penyedia_id, hps)
SELECT k, k, k, (SELECT id FROM instansi WHERE nama = 'Sudin Uji Relasi'), false, 'Pekerjaan Konstruksi', 'Pengadaan Langsung', 2026,
    'Paket Sudah Selesai', 'prv_uji_relasi', '31.72.06', p, h
FROM (VALUES
    ('relasi-1', (SELECT id FROM penyedia WHERE nama_normal = 'cv uji relasi'), 300),
    ('relasi-2', (SELECT id FROM penyedia WHERE nama_normal = 'cv uji relasi'), 200),
    ('relasi-tanpa-penyedia', NULL, 900)) AS v (k, p, h);
INSERT INTO pengguna (email, peran, nama_tampilan) VALUES ('pengirim@contoh.id', 'warga', 'Jurnalis Uji'), ('pemeriksa@contoh.id', 'warga', '')`

func siapkanRelasiUji(t *testing.T) pgx.Tx {
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	if _, err := tx.Exec(t.Context(), sqlPaketRelasiUji); err != nil {
		t.Fatal(err)
	}
	return tx
}

func idPenggunaUji(t *testing.T, tx pgx.Tx, email string) int64 {
	var id int64
	if err := tx.QueryRow(t.Context(), `SELECT id FROM pengguna WHERE email = $1`, email).Scan(&id); err != nil {
		t.Fatal(err)
	}
	return id
}

func TestBacaRelasiPaket_WilayahKota_HanyaPaketBerpenyediaUrutNilai(t *testing.T) {
	tx := siapkanRelasiUji(t)

	relasi, err := repository.JaringanRepository{DB: tx}.BacaRelasiPaket(t.Context(), domain.FilterJaringan{KodeWilayah: "31.72"}, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(relasi) != 2 || relasi[0].PaketKode != "relasi-1" || relasi[0].Nilai != 300 || relasi[0].PenyediaNama != "CV Uji Relasi" {
		t.Fatalf("relasi %+v", relasi)
	}
}

func TestTemuan_KirimTanggapiSembunyikan_HitunganDanVisibilitas(t *testing.T) {
	tx := siapkanRelasiUji(t)
	repo := repository.TemuanRepository{DB: tx}
	pengirim, pemeriksa := idPenggunaUji(t, tx, "pengirim@contoh.id"), idPenggunaUji(t, tx, "pemeriksa@contoh.id")

	id, err := repo.SimpanTemuan(t.Context(), domain.TemuanMasuk{
		PenggunaID: pengirim, PaketKode: "relasi-1", Isi: "Papan proyek menyebut nilai berbeda.", TautanSumber: "https://contoh.id",
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, sikap := range []domain.SikapTemuan{domain.SikapMembantah, domain.SikapMendukung} {
		if err := repo.SimpanTanggapan(t.Context(), domain.TanggapanMasuk{TemuanID: id, PenggunaID: pemeriksa, Sikap: sikap}); err != nil {
			t.Fatal(err)
		}
	}
	daftar, err := repo.DaftarTemuan(t.Context(), "relasi-1", 10)
	if err != nil || len(daftar) != 1 {
		t.Fatalf("daftar %+v err %v", daftar, err)
	}
	if d := daftar[0]; d.JumlahMendukung != 1 || d.JumlahMembantah != 0 || d.NamaKontributor != "Jurnalis Uji" {
		t.Errorf("temuan %+v", d)
	}
	jumlah, err := repo.HitungTemuanSejak(t.Context(), pengirim, time.Now().Add(-time.Hour))
	if err != nil || jumlah != 1 {
		t.Errorf("jumlah %d err %v", jumlah, err)
	}

	if ada, err := repo.SembunyikanTemuan(t.Context(), id); err != nil || !ada {
		t.Fatalf("sembunyikan ada %v err %v", ada, err)
	}
	if _, ada, _ := repo.AmbilPengirimTemuan(t.Context(), id); ada {
		t.Error("temuan tersembunyi masih bisa ditanggapi")
	}
	if daftar, _ := repo.DaftarTemuan(t.Context(), "relasi-1", 10); len(daftar) != 0 {
		t.Errorf("temuan tersembunyi masih tampil: %+v", daftar)
	}
}
