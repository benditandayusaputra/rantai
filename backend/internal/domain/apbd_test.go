package domain_test

import (
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var (
	akunBelanjaUji = domain.AkunAPBD{Nama: "Belanja Daerah", Anggaran: 100, Realisasi: 60}
	akunModalUji   = domain.AkunAPBD{Nama: "Belanja Modal", Anggaran: 30, Realisasi: 10}
	aturanAPBDUji  = domain.AturanAPBD{AkunBelanja: "Belanja Daerah", JenisBelanja: []string{"Belanja Pegawai", "Belanja Modal"}}
)

func TestRapikanAkunAPBD_BarisKembarBerurutan_Disatukan(t *testing.T) {
	rapi := domain.RapikanAkunAPBD([]domain.AkunAPBD{akunBelanjaUji, akunModalUji, akunModalUji, akunBelanjaUji})
	if len(rapi) != 3 || rapi[1] != akunModalUji || rapi[2] != akunBelanjaUji {
		t.Fatalf("akun %+v", rapi)
	}
}

func TestIsRealisasiAPBDSama(t *testing.T) {
	berubah := akunModalUji
	berubah.Realisasi = 11
	anggaranLain := akunModalUji
	anggaranLain.Anggaran = 99
	kasus := []struct {
		nama  string
		a, b  []domain.AkunAPBD
		ingin bool
	}{
		{"identik", []domain.AkunAPBD{akunModalUji}, []domain.AkunAPBD{akunModalUji}, true},
		{"hanya anggaran berbeda", []domain.AkunAPBD{akunModalUji}, []domain.AkunAPBD{anggaranLain}, true},
		{"realisasi bertambah", []domain.AkunAPBD{akunModalUji}, []domain.AkunAPBD{berubah}, false},
		{"jumlah baris berbeda", []domain.AkunAPBD{akunModalUji}, []domain.AkunAPBD{akunModalUji, akunBelanjaUji}, false},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if got := domain.IsRealisasiAPBDSama(k.a, k.b); got != k.ingin {
				t.Errorf("sama %v, ingin %v", got, k.ingin)
			}
		})
	}
}

func TestRingkasBelanja_JenisMengikutiUrutanAturanDanMelewatiYangTidakAda(t *testing.T) {
	apbd := domain.APBDDaerah{Tahun: 2026, BulanRealisasi: 10, Akun: []domain.AkunAPBD{akunModalUji, akunBelanjaUji}, ProvenanceID: "prv"}
	ringkas, err := apbd.RingkasBelanja(aturanAPBDUji)
	if err != nil {
		t.Fatal(err)
	}
	if ringkas.Belanja != akunBelanjaUji || len(ringkas.Jenis) != 1 || ringkas.Jenis[0] != akunModalUji || ringkas.BulanRealisasi != 10 {
		t.Errorf("ringkasan %+v", ringkas)
	}
}

func TestRingkasBelanja_TanpaAkunBelanja_TidakDitemukan(t *testing.T) {
	_, err := domain.APBDDaerah{Akun: []domain.AkunAPBD{akunModalUji}}.RingkasBelanja(aturanAPBDUji)
	if !errors.Is(err, domain.ErrAPBDTidakDitemukan) {
		t.Errorf("err %v", err)
	}
}
