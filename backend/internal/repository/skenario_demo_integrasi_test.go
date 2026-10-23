package repository_test

import (
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const pathSkenarioDemoUji = "../../../data/snapshot/skenario_demo.yaml"

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
