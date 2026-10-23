package wilayah

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

func jalankanImporUji(t *testing.T, tx pgx.Tx) {
	t.Helper()
	var cfg Konfigurasi
	if err := config.BacaYAML("../../../config/wilayah.yaml", &cfg); err != nil {
		t.Fatal(err)
	}
	impor := usecase.ImporWilayah{
		Pengambil:   sumberUji(t),
		Penyimpan:   repository.WilayahRepository{DB: tx},
		Pencatat:    repository.IngestRunRepository{DB: tx},
		AturanAlias: cfg.AturanAlias(),
		Provinsi:    cfg.Provinsi,
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func hitungWilayahDKI(t *testing.T, tx pgx.Tx) map[string]int {
	t.Helper()
	baris, err := tx.Query(t.Context(), `SELECT tingkat::text, count(*) FROM wilayah WHERE kode LIKE '31%' GROUP BY tingkat`)
	if err != nil {
		t.Fatal(err)
	}
	jumlah := map[string]int{}
	for baris.Next() {
		var tingkat string
		var n int
		if err := baris.Scan(&tingkat, &n); err != nil {
			t.Fatal(err)
		}
		jumlah[tingkat] = n
	}
	return jumlah
}

func TestImporWilayah_FixtureDKI_SeluruhWilayahTersimpan(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	jalankanImporUji(t, tx)

	jumlah := hitungWilayahDKI(t, tx)
	harap := map[string]int{"provinsi": 1, "kota": 6, "kecamatan": 44, "kelurahan": 267}
	for tingkat, n := range harap {
		if jumlah[tingkat] != n {
			t.Errorf("jumlah %s = %d, harap %d", tingkat, jumlah[tingkat], n)
		}
	}
}

func TestImporWilayah_FixtureDKI_AliasSingkatanTersimpan(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	jalankanImporUji(t, tx)

	var punyaAlias bool
	err := tx.QueryRow(t.Context(),
		`SELECT 'Klp. Gading Barat' = ANY(alias) FROM wilayah WHERE kode = '31.72.06.1003'`).Scan(&punyaAlias)
	if err != nil || !punyaAlias {
		t.Errorf("alias Klp. Gading Barat tidak ada (err %v)", err)
	}
}

func TestImporWilayah_DijalankanDuaKali_TidakMenggandakanData(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	jalankanImporUji(t, tx)
	jalankanImporUji(t, tx)

	var baru, berubah int
	err := tx.QueryRow(t.Context(),
		`SELECT jumlah_baru, jumlah_berubah FROM ingest_run WHERE sumber = 'wilayah' ORDER BY id DESC LIMIT 1`).Scan(&baru, &berubah)
	if err != nil || baru != 0 || berubah != 0 {
		t.Errorf("run kedua baru = %d, berubah = %d, err = %v", baru, berubah, err)
	}
}
