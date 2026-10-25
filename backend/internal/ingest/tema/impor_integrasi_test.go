package tema

import (
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

func imporTemaUji(t *testing.T, tx pgx.Tx) {
	t.Helper()
	impor := usecase.ImporTema{
		Pengambil: Sumber{Dir: "../../../config/tema", URLRujukan: "https://contoh.test", Sekarang: time.Now},
		Penyimpan: repository.TemaRepository{DB: tx},
		Pencatat:  repository.IngestRunRepository{DB: tx},
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func TestImporTema_ConfigPendidikan_TersimpanDiBasisdata(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	imporTemaUji(t, tx)

	var namaID, komisi, sumber string
	err := tx.QueryRow(t.Context(), `SELECT t.nama_id, t.konfigurasi->'tanggung_jawab'->0->>'komisi', p.sumber::text
FROM tema t JOIN provenance p ON p.id = t.provenance_id WHERE t.kode = 'pendidikan' AND t.aktif`).Scan(&namaID, &komisi, &sumber)
	if err != nil || namaID != "Pendidikan SMP" || komisi != "E" || sumber != "rantai" {
		t.Errorf("nama %q, komisi %q, sumber %q, err %v", namaID, komisi, sumber, err)
	}
}

func TestImporTema_DijalankanDuaKali_TidakAdaPerubahan(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	imporTemaUji(t, tx)
	imporTemaUji(t, tx)

	var baru, berubah int
	err := tx.QueryRow(t.Context(),
		`SELECT jumlah_baru, jumlah_berubah FROM ingest_run WHERE sumber = 'tema' ORDER BY id DESC LIMIT 1`).Scan(&baru, &berubah)
	if err != nil || baru != 0 || berubah != 0 {
		t.Errorf("run kedua baru = %d, berubah = %d, err = %v", baru, berubah, err)
	}
}
