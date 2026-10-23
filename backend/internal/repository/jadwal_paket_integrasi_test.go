package repository_test

import (
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const sqlPaketJadwalUji = `INSERT INTO provenance (id, sumber, diambil_pada, label) VALUES ('prv_uji_jadwal', 'lpse', now(), 'resmi');
INSERT INTO instansi (nama, jenis) VALUES ('Instansi Uji Jadwal', 'lainnya');
INSERT INTO rup (kode, nama, jadwal_pelaksanaan_selesai, provenance_id) VALUES ('rup-uji', 'Rencana', '2026-03-31', 'prv_uji_jadwal');
INSERT INTO paket (kode, kode_spse, nama, kode_rup, instansi_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap, kategori, provenance_id)
VALUES ('nontender-jadwal', 'jadwal', 'Pembangunan Pintu Air', 'rup-uji', (SELECT id FROM instansi WHERE nama = 'Instansi Uji Jadwal'), false,
    'Pekerjaan Konstruksi', 'Pengadaan Langsung', 2026, 'Paket Sudah Selesai', 'pintu_air', 'prv_uji_jadwal')`

func TestStatusPemeliharaan_JadwalDariRUPDanKonfigurasiNyata_MasaPermanen(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	if _, err := tx.Exec(t.Context(), sqlPaketJadwalUji); err != nil {
		t.Fatal(err)
	}
	var cfg usecase.KonfigurasiPemeliharaan
	if err := config.BacaYAML("../../config/pemeliharaan.yaml", &cfg); err != nil {
		t.Fatal(err)
	}
	status := usecase.StatusPemeliharaan{
		Pembaca:     repository.PaketRepository{DB: tx},
		Konfigurasi: cfg,
		Sekarang:    func() time.Time { return time.Date(2026, time.April, 1, 0, 0, 0, 0, time.UTC) },
	}

	hasil, err := status.Hitung(t.Context(), "nontender-jadwal")

	if err != nil || hasil.Estimasi.Status != domain.StatusPemeliharaanBerjalan || hasil.Estimasi.AkhirPemeliharaan.Format(time.DateOnly) != "2026-09-27" {
		t.Errorf("hasil %+v, err %v", hasil, err)
	}
}
