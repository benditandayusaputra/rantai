package wikimedia

import (
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const polaThumbnailUji = "https://upload.wikimedia.org/wikipedia/commons/thumb/%s/%s/%s/%dpx-%s"

func TestAlamatThumbnail_DirektoriDariMD5NamaBerkas(t *testing.T) {
	kasus := []struct{ berkas, harap string }{
		{"SMP Negeri 101 Jakarta 2013-08-16.jpg", "https://upload.wikimedia.org/wikipedia/commons/thumb/1/19/SMP_Negeri_101_Jakarta_2013-08-16.jpg/330px-SMP_Negeri_101_Jakarta_2013-08-16.jpg"},
		{"Sekolah Don Bosco I Kelapa Gading (2025) 01.jpg", "https://upload.wikimedia.org/wikipedia/commons/thumb/2/20/Sekolah_Don_Bosco_I_Kelapa_Gading_%282025%29_01.jpg/330px-Sekolah_Don_Bosco_I_Kelapa_Gading_%282025%29_01.jpg"},
	}
	for _, k := range kasus {
		t.Run(k.berkas, func(t *testing.T) {
			if hasil := AlamatThumbnail(polaThumbnailUji, NamaBerkas(k.berkas), 330); hasil != k.harap {
				t.Errorf("dapat %s", hasil)
			}
		})
	}
}

func TestAmbil_KonfigurasiRepo_SemuaEntriLengkapDenganKredit(t *testing.T) {
	var cfg Konfigurasi
	if err := config.BacaYAML("../../../config/foto_wikimedia.yaml", &cfg); err != nil {
		t.Fatal(err)
	}

	daftar, err := Sumber{Konfigurasi: cfg}.Ambil(t.Context())

	if err != nil || len(daftar) == 0 {
		t.Fatalf("daftar %d, err %v", len(daftar), err)
	}
	for _, f := range daftar {
		if f.Sumber != domain.SumberFotoWikimedia || f.Kredit == "" || f.Halaman == "" {
			t.Errorf("foto %s tidak lengkap: %+v", f.NPSN, f.Foto)
		}
	}
}

func TestAmbil_EntriTanpaBerkas_Ditolak(t *testing.T) {
	cfg := Konfigurasi{PolaThumbnail: polaThumbnailUji, Foto: []Entri{{NPSN: "20100379"}}}

	_, err := Sumber{Konfigurasi: cfg}.Ambil(t.Context())

	if !errors.Is(err, ErrEntriTidakLengkap) {
		t.Errorf("harap ErrEntriTidakLengkap, dapat %v", err)
	}
}
