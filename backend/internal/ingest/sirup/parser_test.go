package sirup

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

func bacaFixture(t *testing.T, nama string) []byte {
	t.Helper()
	isi, err := os.ReadFile(filepath.Join("testdata", nama))
	if err != nil {
		t.Fatal(err)
	}
	return isi
}

func TestUraiRUP_PaketKonstruksi_SemuaFieldTerisi(t *testing.T) {
	rup, err := UraiRUP(bacaFixture(t, "rup_63728635.html"))

	if err != nil {
		t.Fatal(err)
	}
	if rup.Kode != "63728635" || rup.SatuanKerja != "SUKU DINAS SUMBER DAYA AIR KOTA ADMINISTRASI JAKARTA PUSAT" || rup.TahunAnggaran != 2026 {
		t.Errorf("identitas = %+v", rup)
	}
	if rup.Pagu == nil || *rup.Pagu != 185101704 || rup.SumberDana != "APBD" || rup.JenisPengadaan != "Pekerjaan Konstruksi" || rup.Metode != "Pengadaan Langsung" {
		t.Errorf("anggaran = %+v", rup)
	}
	mulai := time.Date(2026, time.September, 1, 0, 0, 0, 0, zonaWIB)
	selesai := time.Date(2026, time.December, 31, 0, 0, 0, 0, zonaWIB)
	if rup.PelaksanaanMulai == nil || !rup.PelaksanaanMulai.Equal(mulai) || !rup.PelaksanaanSelesai.Equal(selesai) {
		t.Errorf("jadwal pelaksanaan = %v sampai %v", rup.PelaksanaanMulai, rup.PelaksanaanSelesai)
	}
}

func TestUraiRUP_JadwalPelaksanaanBukanJadwalPemanfaatan(t *testing.T) {
	rup, err := UraiRUP(bacaFixture(t, "rup_60262521.html"))

	if err != nil || rup.PelaksanaanSelesai.Month() != time.October || rup.PelaksanaanSelesai.Day() != 31 {
		t.Errorf("jadwal = %v, err = %v", rup.PelaksanaanSelesai, err)
	}
}

func TestUraiRUP_HalamanKesalahan_ErrRUPTidakTersedia(t *testing.T) {
	_, err := UraiRUP([]byte(`<html><body><table><tr><td>Ray ID</td><td>a4ef</td></tr></table></body></html>`))

	if !errors.Is(err, domain.ErrRUPTidakTersedia) {
		t.Errorf("err = %v", err)
	}
}
