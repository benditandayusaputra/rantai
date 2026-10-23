package domain

import (
	"errors"
	"testing"
)

func skenarioSah() SkenarioDemo {
	return SkenarioDemo{
		PaketSimulasi: PaketSimulasi{Kode: "simulasi-1", Penyedia: "CV Contoh Simulasi", Instansi: "Instansi Contoh Simulasi"},
		AkunDemo:      []AkunDemo{{Peran: "warga"}, {Peran: PeranAnggota}},
	}
}

func TestValidasiSkenarioDemo(t *testing.T) {
	kasus := []struct {
		nama  string
		ubah  func(s *SkenarioDemo)
		harap error
	}{
		{"skenario sah", func(*SkenarioDemo) {}, nil},
		{"kode berpola paket nyata", func(s *SkenarioDemo) { s.PaketSimulasi.Kode = "nontender-12345" }, ErrKodePaketBukanSimulasi},
		{"penyedia bernama nyata", func(s *SkenarioDemo) { s.PaketSimulasi.Penyedia = "CV Karya Mandiri" }, ErrNamaSimulasiTanpaPenanda},
		{"instansi bernama nyata", func(s *SkenarioDemo) { s.PaketSimulasi.Instansi = "Dinas Pendidikan" }, ErrNamaSimulasiTanpaPenanda},
		{"akun anggota tertaut anggota", func(s *SkenarioDemo) { s.AkunDemo[1].AnggotaID = "anggota-uji" }, ErrAkunDemoTertautAnggota},
		{"akun instansi boleh tanpa anggota", func(s *SkenarioDemo) { s.AkunDemo[0].Peran = "instansi" }, nil},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			s := skenarioSah()
			k.ubah(&s)
			if err := s.Validasi(); !errors.Is(err, k.harap) {
				t.Errorf("err = %v, harap %v", err, k.harap)
			}
		})
	}
}
