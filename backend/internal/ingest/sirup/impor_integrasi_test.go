package sirup

import (
	"context"
	"fmt"
	"path"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

type pengunduhFixture struct {
	t *testing.T
}

func (p pengunduhFixture) Ambil(_ context.Context, alamat string) ([]byte, error) {
	if path.Base(alamat) == "63728635" {
		return bacaFixture(p.t, "rup_63728635.html"), nil
	}
	return nil, fmt.Errorf("%s status 500: %w", alamat, klien.ErrStatusTidakOK)
}

const sqlPaketUji = `INSERT INTO provenance (id, sumber, diambil_pada, label) VALUES ('prv_uji_lpse', 'lpse', now(), 'resmi');
INSERT INTO instansi (nama, jenis) VALUES ('Provinsi DKI Jakarta', 'lainnya');
INSERT INTO rup (kode, nama, provenance_id) VALUES ('63728635', 'stub', 'prv_uji_lpse'), ('99999999', 'stub', 'prv_uji_lpse');
INSERT INTO paket (kode, kode_spse, nama, kode_rup, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap, provenance_id)
SELECT k, k, k, r, (SELECT id FROM instansi), false, 'Pekerjaan Konstruksi', 'Pengadaan Langsung', 2026, 'Paket Sudah Selesai', 'prv_uji_lpse'
FROM (VALUES ('nontender-1', '63728635'), ('nontender-2', '99999999'), ('nontender-3', NULL)) AS v (k, r)`

func jalankanImporRUP(t *testing.T) pgx.Tx {
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	if _, err := tx.Exec(t.Context(), sqlPaketUji); err != nil {
		t.Fatal(err)
	}
	impor := usecase.ImporRUP{
		Pengambil: Sumber{Klien: pengunduhFixture{t: t}, URLDetail: "https://sirup.uji/", DirRaw: t.TempDir(), Sekarang: time.Now},
		Penyimpan: repository.RUPRepository{DB: tx},
		Pencatat:  repository.IngestRunRepository{DB: tx},
		AturanInstansi: domain.AturanInstansi{
			Jenis: []domain.AturanJenisInstansi{{Jenis: domain.JenisInstansiSudin, Awalan: []string{"suku dinas"}}},
			Kota:  map[string][]string{"31.71": {"jakarta pusat"}},
		},
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	return tx
}

func TestImporRUP_KodeDitemukan_DetailDanJadwalTersimpan(t *testing.T) {
	tx := jalankanImporRUP(t)

	var pagu float64
	var kota string
	var selesai time.Time
	err := tx.QueryRow(t.Context(), `SELECT r.pagu, i.kode_kota, r.jadwal_pelaksanaan_selesai FROM rup r JOIN instansi i ON i.id = r.instansi_id
WHERE r.kode = '63728635'`).Scan(&pagu, &kota, &selesai)
	if err != nil {
		t.Fatal(err)
	}
	if pagu != 185101704 || kota != "31.71" || selesai.Format(time.DateOnly) != "2026-12-31" {
		t.Errorf("pagu %v kota %q selesai %v", pagu, kota, selesai)
	}
}

func TestImporRUP_DetailTidakTersediaDanTanpaRUP_PaketTetapTersimpan(t *testing.T) {
	tx := jalankanImporRUP(t)

	var jumlahPaket, tanpaRUP int
	var sumberStub string
	err := tx.QueryRow(t.Context(), `SELECT (SELECT count(*) FROM paket), (SELECT count(*) FROM paket WHERE kode_rup IS NULL),
    (SELECT p.sumber::text FROM rup r JOIN provenance p ON p.id = r.provenance_id WHERE r.kode = '99999999')`).Scan(&jumlahPaket, &tanpaRUP, &sumberStub)
	if err != nil {
		t.Fatal(err)
	}
	if jumlahPaket != 3 || tanpaRUP != 1 || sumberStub != "lpse" {
		t.Errorf("paket %d, tanpa rup %d, sumber stub %s", jumlahPaket, tanpaRUP, sumberStub)
	}
}
