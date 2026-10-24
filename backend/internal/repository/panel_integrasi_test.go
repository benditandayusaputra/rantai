package repository_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	kodePaketSimulasiUji    = "simulasi-ruang-kelas-1"
	pathPemeliharaanUji     = "../../config/pemeliharaan.yaml"
	emailPetugasSimulasiUji = "petugas@demo.rantai.test"
)

type panelUjiDB struct {
	tx      pgx.Tx
	panel   usecase.Panel
	petugas int64
}

func siapkanPanelDemo(t *testing.T) panelUjiDB {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	siapkanKursiUji(t, tx)
	jalankanHitungIndikator(t, tx)
	var terpilih string
	if err := muatSkenarioDemo(t, tx, &terpilih); err != nil {
		t.Fatal(err)
	}
	var pemeliharaan usecase.KonfigurasiPemeliharaan
	if err := config.BacaYAML(pathPemeliharaanUji, &pemeliharaan); err != nil {
		t.Fatal(err)
	}
	var petugas int64
	if err := tx.QueryRow(t.Context(), `SELECT id FROM pengguna WHERE email = $1`, emailPetugasSimulasiUji).Scan(&petugas); err != nil {
		t.Fatal(err)
	}
	return panelUjiDB{tx: tx, petugas: petugas, panel: usecase.Panel{
		Pengguna: repository.PenggunaRepository{DB: tx}, Penyimpan: repository.PanelRepository{DB: tx},
		Laporan: repository.LaporanRepository{DB: tx}, Pemeliharaan: pemeliharaan,
		BuatProvenance: buatProvenanceUji, Sekarang: time.Now,
	}}
}

func (p panelUjiDB) bacaPaket(t *testing.T) usecase.DetailPaket {
	t.Helper()
	baca := usecase.BacaPaket{Pembaca: repository.PaketRepository{DB: p.tx}, Pemeliharaan: p.panel.Pemeliharaan, Sekarang: time.Now}
	detail, err := baca.Ambil(t.Context(), kodePaketSimulasiUji)
	if err != nil {
		t.Fatal(err)
	}
	return detail
}

func TestPanel_HakJawabPaket_LangsungTampilPublikDanTidakBisaDihapus(t *testing.T) {
	p := siapkanPanelDemo(t)

	tulis, err := p.panel.TulisHakJawab(t.Context(), p.petugas, domain.HakJawabBaru{
		TargetJenis: domain.TargetHakJawabPaket, TargetID: kodePaketSimulasiUji, Isi: "Rehabilitasi selesai sesuai kontrak.",
	})
	if err != nil {
		t.Fatal(err)
	}
	detail := p.bacaPaket(t)

	if len(detail.Paket.HakJawab) != 1 || detail.Paket.HakJawab[0].Pihak != "Instansi Contoh Simulasi" || !detail.Paket.HakJawab[0].IsSimulation {
		t.Errorf("hak jawab publik %+v", detail.Paket.HakJawab)
	}
	for _, perintah := range []string{`DELETE FROM hak_jawab WHERE id = $1`, `UPDATE hak_jawab SET isi = 'diganti' WHERE id = $1`} {
		simpan, err := p.tx.Begin(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		if _, err := simpan.Exec(t.Context(), perintah, tulis.ID); err == nil {
			t.Errorf("%s berhasil, harap ditolak", perintah)
		}
		_ = simpan.Rollback(t.Context())
	}
}

func TestPanel_LaporanDiteruskan_MasukAntreanLaluDitanggapiDanSelesai(t *testing.T) {
	p := siapkanPanelDemo(t)

	antrean, err := p.panel.Antrean(t.Context(), p.petugas)
	if err != nil || len(antrean.Laporan) != 2 || antrean.Laporan[0].Status != domain.StatusDiteruskan {
		t.Fatalf("antrean %+v, err %v", antrean.Laporan, err)
	}
	id := antrean.Laporan[0].ID
	if _, err := p.panel.TulisHakJawab(t.Context(), p.petugas, domain.HakJawabBaru{
		TargetJenis: domain.TargetHakJawabLaporan, TargetID: strconv.FormatInt(id, 10), Isi: "Plafon sudah diganti.",
	}); err != nil {
		t.Fatal(err)
	}
	if err := p.panel.TandaiLaporanSelesai(t.Context(), p.petugas, id); err != nil {
		t.Fatal(err)
	}

	var status string
	var riwayat int
	err = p.tx.QueryRow(t.Context(), `SELECT l.status, (SELECT count(*) FROM laporan_riwayat r WHERE r.laporan_id = l.id
    AND r.status_baru IN ('DITANGGAPI', 'SELESAI') AND r.oleh_pengguna_id = $2) FROM laporan l WHERE l.id = $1`, id, p.petugas).Scan(&status, &riwayat)
	if err != nil || status != string(domain.StatusSelesai) || riwayat != 2 {
		t.Errorf("status %s, riwayat oleh petugas %d, err %v", status, riwayat, err)
	}
}

func TestPanel_SerahTerima_LabelPemeliharaanResmiDariInstansi(t *testing.T) {
	p := siapkanPanelDemo(t)
	sebelum := p.bacaPaket(t).Pemeliharaan

	err := p.panel.IsiSerahTerima(t.Context(), p.petugas, domain.SerahTerima{
		KodePaket: kodePaketSimulasiUji, Tanggal: time.Date(2026, 7, 3, 0, 0, 0, 0, time.UTC), NomorBeritaAcara: "BA-07/2026",
	})
	if err != nil {
		t.Fatal(err)
	}
	sesudah := p.bacaPaket(t)

	var catatan string
	if err := p.tx.QueryRow(t.Context(), `SELECT catatan FROM provenance WHERE id = $1`, sesudah.Pemeliharaan.ProvenanceID).Scan(&catatan); err != nil {
		t.Fatal(err)
	}
	if sebelum.Estimasi.Label != domain.LabelEstimasi || sesudah.Pemeliharaan.Estimasi.Label != domain.LabelResmi ||
		sesudah.Paket.NomorBeritaAcara != "BA-07/2026" || catatan == "" {
		t.Errorf("sebelum %+v, sesudah %+v, catatan %q", sebelum, sesudah.Pemeliharaan, catatan)
	}
}

func TestPanel_PetugasNyataTidakBisaMenanggapiPaketSimulasi(t *testing.T) {
	p := siapkanPanelDemo(t)
	var nyata int64
	err := p.tx.QueryRow(t.Context(), `WITH i AS (INSERT INTO instansi (nama, jenis) VALUES ('Dinas Uji Nyata', 'dinas') RETURNING id)
INSERT INTO pengguna (email, peran, instansi_id) SELECT 'petugas@contoh.id', 'instansi', id FROM i RETURNING id`).Scan(&nyata)
	if err != nil {
		t.Fatal(err)
	}

	_, err = p.panel.TulisHakJawab(t.Context(), nyata, domain.HakJawabBaru{
		TargetJenis: domain.TargetHakJawabPaket, TargetID: kodePaketSimulasiUji, Isi: "Tanggapan.",
	})

	if !errors.Is(err, domain.ErrBukanPenjawabTarget) {
		t.Errorf("err %v", err)
	}
}

func TestPantau_HakJawabBaru_MenjadiKejadianKecamatan(t *testing.T) {
	p := siapkanPanelDemo(t)
	sejak := time.Now().Add(-time.Minute)
	if _, err := p.panel.TulisHakJawab(t.Context(), p.petugas, domain.HakJawabBaru{
		TargetJenis: domain.TargetHakJawabPaket, TargetID: kodePaketSimulasiUji, Isi: "Tanggapan instansi.",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := (repository.PantauRepository{DB: p.tx}).BuatKejadian(t.Context(), sejak); err != nil {
		t.Fatal(err)
	}

	var tautan string
	var simulasi bool
	err := p.tx.QueryRow(t.Context(), `SELECT tautan, is_simulation FROM kejadian WHERE jenis = 'hak_jawab'`).Scan(&tautan, &simulasi)
	if err != nil || tautan != "/paket/"+kodePaketSimulasiUji || !simulasi {
		t.Errorf("tautan %s, simulasi %v, err %v", tautan, simulasi, err)
	}
}
