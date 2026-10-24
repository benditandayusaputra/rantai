package repository_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/foto"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const pathSkenarioDemoUji = "../../../data/snapshot/skenario_demo.yaml"

type penyimpanFotoUji struct{ kunci map[string]bool }

func (p penyimpanFotoUji) SimpanJPEG(_ context.Context, kunci string, _ []byte) error {
	p.kunci[kunci] = true
	return nil
}

type pelaporKecamatanDemo struct {
	terpilih *string
}

func (p pelaporKecamatanDemo) LaporkanKecamatanDemoTerpilih(kode string) {
	*p.terpilih = kode
}

func muatSkenarioDemo(t *testing.T, tx pgx.Tx, terpilih *string) error {
	t.Helper()
	var skenario domain.SkenarioDemo
	if err := config.BacaYAML(pathSkenarioDemoUji, &skenario); err != nil {
		t.Fatal(err)
	}
	muat := usecase.MuatSkenarioDemo{
		Skenario:       skenario,
		Penyimpan:      repository.SkenarioDemoRepository{DB: tx},
		Pelapor:        pelaporKecamatanDemo{terpilih: terpilih},
		BuatProvenance: buatProvenanceUji,
		Pemroses:       foto.Pemroses{},
		PenyimpanFoto:  penyimpanFotoUji{kunci: map[string]bool{}},
		BacaFoto: func(nama string) ([]byte, error) {
			return os.ReadFile(filepath.Join(filepath.Dir(pathSkenarioDemoUji), nama))
		},
	}
	return muat.Jalankan(t.Context())
}

func TestMuatSkenarioDemo_KecamatanRasioTerendahDipilihDanPaketBertandaSimulasi(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanKursiUji(t, tx)
	jalankanHitungIndikator(t, tx)
	var terpilih string

	if err := muatSkenarioDemo(t, tx, &terpilih); err != nil {
		t.Fatal(err)
	}

	var simulasi int
	var tanpaTanda int
	err := tx.QueryRow(t.Context(), `SELECT count(*) FILTER (WHERE p.is_simulation AND p.kode_wilayah = $1),
    count(*) FILTER (WHERE NOT p.is_simulation)
FROM paket p JOIN provenance v ON v.id = p.provenance_id WHERE v.label = 'simulasi'`, terpilih).Scan(&simulasi, &tanpaTanda)
	if err != nil || terpilih != "31.72.06" || simulasi != 1 || tanpaTanda != 0 {
		t.Errorf("terpilih %s, paket simulasi %d, paket berlabel simulasi tanpa tanda %d, err %v", terpilih, simulasi, tanpaTanda, err)
	}
}

func TestMuatSkenarioDemo_DijalankanDuaKali_PihakSimulasiTidakTertautDataNyata(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanKursiUji(t, tx)
	jalankanHitungIndikator(t, tx)
	var terpilih string

	for range 2 {
		if err := muatSkenarioDemo(t, tx, &terpilih); err != nil {
			t.Fatal(err)
		}
	}

	var paketSimulasi, paketNyataTertaut int
	err := tx.QueryRow(t.Context(), `WITH pihak AS (
    SELECT penyedia_id, instansi_id FROM paket WHERE is_simulation)
SELECT (SELECT count(*) FROM paket WHERE is_simulation),
    (SELECT count(*) FROM paket p, pihak WHERE NOT p.is_simulation
        AND (p.penyedia_id = pihak.penyedia_id OR p.instansi_id = pihak.instansi_id))
    + (SELECT count(*) FROM riwayat_penyedia r, pihak WHERE r.penyedia_id = pihak.penyedia_id)`).Scan(&paketSimulasi, &paketNyataTertaut)
	if err != nil || paketSimulasi != 1 || paketNyataTertaut != 0 {
		t.Errorf("paket simulasi %d, data nyata tertaut %d, err %v", paketSimulasi, paketNyataTertaut, err)
	}
}

func TestMuatSkenarioDemo_PenyediaSimulasiSudahMemenangkanPaketNyata_Ditolak(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanKursiUji(t, tx)
	jalankanHitungIndikator(t, tx)
	_, err := tx.Exec(t.Context(), `WITH penyedia_nyata AS (
    INSERT INTO penyedia (nama, nama_normal) VALUES ('CV Contoh Simulasi', $1) RETURNING id),
instansi_nyata AS (
    INSERT INTO instansi (nama, jenis) VALUES ('Instansi Uji Nyata', 'lainnya') RETURNING id)
INSERT INTO paket (kode, kode_spse, nama, instansi_id, penyedia_id, is_tender, jenis_pengadaan, metode, tahun_anggaran, tahap, provenance_id)
SELECT 'nontender-90000001', '90000001', 'Paket Uji Nyata', instansi_nyata.id, penyedia_nyata.id, false,
    'Pekerjaan Konstruksi', 'Pengadaan Langsung', 2026, 'Selesai', (SELECT provenance_id FROM wilayah WHERE kode = '31.72.06')
FROM penyedia_nyata, instansi_nyata`, domain.NormalisasiNamaPenyedia("CV Contoh Simulasi"))
	if err != nil {
		t.Fatal(err)
	}
	var terpilih string

	err = muatSkenarioDemo(t, tx, &terpilih)

	if !errors.Is(err, domain.ErrSimulasiTertautDataNyata) {
		t.Errorf("err = %v, harap %v", err, domain.ErrSimulasiTertautDataNyata)
	}
}

func TestMuatSkenarioDemo_IndikatorBelumDihitung_Gagal(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	if _, err := tx.Exec(t.Context(), `DELETE FROM indikator_kecamatan`); err != nil {
		t.Fatal(err)
	}
	var terpilih string

	err := muatSkenarioDemo(t, tx, &terpilih)

	if !errors.Is(err, domain.ErrKecamatanDemoBelumDihitung) {
		t.Errorf("err = %v, harap %v", err, domain.ErrKecamatanDemoBelumDihitung)
	}
}

func TestMuatSkenarioDemo_AkunDemoTerpasangSebagaiSimulasi(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanKursiUji(t, tx)
	jalankanHitungIndikator(t, tx)
	var terpilih string

	for range 2 {
		if err := muatSkenarioDemo(t, tx, &terpilih); err != nil {
			t.Fatal(err)
		}
	}

	var akun, simulasi, petugasDinas, tertautAnggota int
	err := tx.QueryRow(t.Context(), `SELECT count(*), count(*) FILTER (WHERE p.is_simulation),
    count(*) FILTER (WHERE p.peran = 'instansi' AND i.nama = 'Dinas Pendidikan Provinsi DKI Jakarta'),
    count(*) FILTER (WHERE p.anggota_id IS NOT NULL)
FROM pengguna p LEFT JOIN instansi i ON i.id = p.instansi_id WHERE p.email LIKE '%@demo.rantai.test'`).Scan(&akun, &simulasi, &petugasDinas, &tertautAnggota)
	if err != nil || akun != 6 || simulasi != 6 || petugasDinas != 1 || tertautAnggota != 0 {
		t.Errorf("akun %d, simulasi %d, petugas dinas %d, tertaut anggota %d, err %v", akun, simulasi, petugasDinas, tertautAnggota, err)
	}
}

func TestMuatSkenarioDemo_LaporanSimulasiTerverifikasiHanyaTertautDataSimulasi(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanKursiUji(t, tx)
	jalankanHitungIndikator(t, tx)
	var terpilih string

	for range 2 {
		if err := muatSkenarioDemo(t, tx, &terpilih); err != nil {
			t.Fatal(err)
		}
	}

	var laporan, terverifikasi, tertautNyata, konfirmasi, foto int
	err := tx.QueryRow(t.Context(), `SELECT count(*), count(*) FILTER (WHERE l.status = 'TERVERIFIKASI' AND l.is_simulation),
    count(*) FILTER (WHERE NOT p.is_simulation OR NOT g.is_simulation),
    (SELECT count(*) FROM laporan_konfirmasi k JOIN pengguna kg ON kg.id = k.pengguna_id WHERE kg.is_simulation),
    (SELECT count(*) FROM laporan_foto f JOIN laporan fl ON fl.id = f.laporan_id WHERE fl.is_simulation)
FROM laporan l JOIN paket p ON p.kode = l.paket_kode JOIN pengguna g ON g.id = l.pengguna_id`).
		Scan(&laporan, &terverifikasi, &tertautNyata, &konfirmasi, &foto)
	if err != nil || laporan != 2 || terverifikasi != 2 || tertautNyata != 0 || konfirmasi != 2*domain.JumlahKonfirmasiTerverifikasi || foto != 2 {
		t.Errorf("laporan %d, terverifikasi %d, tertaut nyata %d, konfirmasi %d, foto %d, err %v",
			laporan, terverifikasi, tertautNyata, konfirmasi, foto, err)
	}
}
