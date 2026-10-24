package domain

import (
	"errors"
	"testing"
)

func skenarioSah() SkenarioDemo {
	return SkenarioDemo{
		PaketSimulasi: PaketSimulasi{Kode: "simulasi-1", Penyedia: "CV Contoh Simulasi", Instansi: "Instansi Contoh Simulasi"},
		AkunDemo: []AkunDemo{
			{Peran: PeranWarga, Email: "warga@demo.rantai.test", NamaTampilan: "Warga Contoh Simulasi"},
			{Peran: PeranAnggota, Email: "anggota@demo.rantai.test", NamaTampilan: "Anggota Contoh Simulasi"},
			{Peran: PeranWarga, Email: "saksi@demo.rantai.test", NamaTampilan: "Saksi Contoh Simulasi"},
		},
		LaporanSimulasi: []LaporanSimulasi{{
			Kondisi: KondisiSelesaiBaik, Pelapor: "warga@demo.rantai.test",
			Saksi: []string{"anggota@demo.rantai.test", "saksi@demo.rantai.test"}, Foto: []string{"foto.jpg"},
		}},
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
		{"akun instansi dengan instansi", func(s *SkenarioDemo) {
			s.AkunDemo[0].Peran, s.AkunDemo[0].Instansi = PeranInstansi, "Instansi Uji Simulasi"
		}, nil},
		{"akun instansi tertaut instansi nyata", func(s *SkenarioDemo) { s.AkunDemo[0].Peran, s.AkunDemo[0].Instansi = PeranInstansi, "Dinas Pendidikan" }, ErrAkunDemoBukanSimulasi},
		{"akun instansi tanpa instansi", func(s *SkenarioDemo) { s.AkunDemo[0].Peran = PeranInstansi }, ErrAkunInstansiTanpaInstansi},
		{"akun warga dengan instansi", func(s *SkenarioDemo) { s.AkunDemo[0].Instansi = "Instansi Uji Simulasi" }, ErrAkunInstansiTanpaInstansi},
		{"peran tak dikenal", func(s *SkenarioDemo) { s.AkunDemo[0].Peran = "pejabat" }, ErrPeranTidakDikenal},
		{"email bukan domain uji", func(s *SkenarioDemo) { s.AkunDemo[0].Email = "warga@gmail.com" }, ErrAkunDemoBukanSimulasi},
		{"nama tampilan tanpa penanda", func(s *SkenarioDemo) { s.AkunDemo[0].NamaTampilan = "Budi" }, ErrAkunDemoBukanSimulasi},
		{"laporan simulasi kurang saksi", func(s *SkenarioDemo) { s.LaporanSimulasi[0].Saksi = s.LaporanSimulasi[0].Saksi[:1] }, ErrLaporanSimulasiTidakSah},
		{"laporan simulasi tanpa foto", func(s *SkenarioDemo) { s.LaporanSimulasi[0].Foto = nil }, ErrLaporanSimulasiTidakSah},
		{"pelapor bukan akun demo", func(s *SkenarioDemo) { s.LaporanSimulasi[0].Pelapor = "warga@gmail.com" }, ErrPelaporBukanAkunDemo},
		{"saksi sama dengan pelapor", func(s *SkenarioDemo) { s.LaporanSimulasi[0].Saksi[0] = "warga@demo.rantai.test" }, ErrPelaporBukanAkunDemo},
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
