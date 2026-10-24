package paripurna

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

func jalankanImporUji(t *testing.T, tx pgx.Tx, sumber Sumber) {
	t.Helper()
	impor := usecase.ImporParipurna{
		Pengambil: sumber,
		Penyimpan: repository.ParipurnaRepository{DB: tx},
		Pencatat:  repository.IngestRunRepository{DB: tx},
		Terlarang: domain.KosakataTerlarang{"korupsi"},
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func hasilImporTerakhir(t *testing.T, tx pgx.Tx) (baru, berubah int) {
	t.Helper()
	err := tx.QueryRow(t.Context(),
		`SELECT jumlah_baru, jumlah_berubah FROM ingest_run WHERE sumber = 'paripurna' ORDER BY id DESC LIMIT 1`).Scan(&baru, &berubah)
	if err != nil {
		t.Fatal(err)
	}
	return baru, berubah
}

func TestImporParipurna_DijalankanDuaKali_TidakMenggandakanData(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	ujiwakil.MuatDPRD(t, tx)

	jalankanImporUji(t, tx, sumberUji())
	jalankanImporUji(t, tx, sumberUji())

	if baru, berubah := hasilImporTerakhir(t, tx); baru != 0 || berubah != 0 {
		t.Errorf("baru %d, berubah %d", baru, berubah)
	}
}

func TestImporParipurna_EntriDihapusDariCSV_IkutTerhapus(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	ujiwakil.MuatDPRD(t, tx)
	jalankanImporUji(t, tx, sumberUji())
	sumber := sumberUji()
	sumber.Konfigurasi.PathAspirasiReses = "testdata/aspirasi_reses_satu.csv"

	jalankanImporUji(t, tx, sumber)

	var jumlah int
	if err := tx.QueryRow(t.Context(), `SELECT count(*) FROM aspirasi_reses`).Scan(&jumlah); err != nil || jumlah != 1 {
		t.Errorf("jumlah aspirasi %d, err %v", jumlah, err)
	}
}
