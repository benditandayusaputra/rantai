package djpk_test

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/djpk"
)

type pengunduhFixture struct {
	berkasPerBulan map[string]string
	diminta        []string
}

func (p *pengunduhFixture) Ambil(_ context.Context, alamat string) ([]byte, error) {
	u, err := url.Parse(alamat)
	if err != nil {
		return nil, err
	}
	bulan := u.Query().Get("periode")
	p.diminta = append(p.diminta, bulan)
	return os.ReadFile(filepath.Join("testdata", p.berkasPerBulan[bulan]))
}

func bacaFixture(t *testing.T, nama string) []byte {
	t.Helper()
	isi, err := os.ReadFile(filepath.Join("testdata", nama))
	if err != nil {
		t.Fatal(err)
	}
	return isi
}

func TestUraiAPBD_BarisKembarDisatukanDanAngkaUtuh(t *testing.T) {
	akun, err := djpk.UraiAPBD(bacaFixture(t, "apbd_2025_12.xls"))
	if err != nil {
		t.Fatal(err)
	}
	if len(akun) != 29 {
		t.Fatalf("jumlah akun %d, ingin 29", len(akun))
	}
	ingin := domain.AkunAPBD{Nama: "Belanja Daerah", Anggaran: 82663062617784, Realisasi: 76098159818559}
	if !slices.Contains(akun, ingin) || akun[0].Nama != "Pendapatan Daerah" {
		t.Errorf("akun %+v", akun[:3])
	}
}

func TestUraiAPBD_KolomBerubah_Ditolak(t *testing.T) {
	rusak := []byte(`<Workbook><Worksheet><Table><Row><Cell><Data>Uraian</Data></Cell></Row></Table></Worksheet></Workbook>`)
	if _, err := djpk.UraiAPBD(rusak); !errors.Is(err, djpk.ErrFormatBerubah) {
		t.Errorf("err %v", err)
	}
}

func TestAmbilAPBD_BulanDepanMengulangDataTerakhir_BulanRealisasiDariPerubahanTerakhir(t *testing.T) {
	pengunduh := &pengunduhFixture{berkasPerBulan: map[string]string{
		"11": "apbd_2026_10.xls", "10": "apbd_2026_10.xls", "9": "apbd_2026_09.xls",
	}}
	dirRaw := t.TempDir()
	sumber := djpk.Sumber{
		Konfigurasi: djpk.Konfigurasi{URLEkspor: "https://djpk.uji/csv", URLRujukan: "https://djpk.uji/apbd", KodeProvinsi: "09", KodePemda: "00"},
		Pengunduh:   pengunduh, DirRaw: dirRaw,
		Sekarang: func() time.Time { return time.Date(2026, time.November, 8, 2, 0, 0, 0, time.UTC) },
	}

	hasil, err := sumber.AmbilAPBD(t.Context(), 2026, 11)
	if err != nil {
		t.Fatal(err)
	}

	if hasil.APBD.BulanRealisasi != 10 || !slices.Equal(pengunduh.diminta, []string{"11", "10", "9"}) {
		t.Errorf("bulan %d, diminta %v", hasil.APBD.BulanRealisasi, pengunduh.diminta)
	}
	if hasil.Provenance.Sumber != domain.SumberDJPK || hasil.APBD.ProvenanceID != hasil.Provenance.ID {
		t.Errorf("provenance %+v", hasil.Provenance)
	}
	if _, err := os.Stat(filepath.Join(dirRaw, "apbd_djpk", "2026-11-08", "apbd_2026_10.xls")); err != nil {
		t.Errorf("data mentah: %v", err)
	}
}
