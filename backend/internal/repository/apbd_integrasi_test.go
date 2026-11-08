package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

func sumberAPBDUji(id string, bulan int, belanja float64) usecase.APBDSumber {
	return usecase.APBDSumber{
		Provenance: domain.Provenance{ID: id, Sumber: domain.SumberDJPK, URL: "https://djpk.uji", DiambilPada: time.Now(), Label: domain.LabelResmi},
		APBD: domain.APBDDaerah{Tahun: 2026, BulanRealisasi: bulan, ProvenanceID: id, Akun: []domain.AkunAPBD{
			{Nama: "Belanja Daerah", Anggaran: 74285299931740, Realisasi: belanja},
			{Nama: "Belanja Modal", Anggaran: 10164767414640.5, Realisasi: 2121930697743},
		}},
	}
}

func TestAPBD_SimpanUlangMenggantiTahunYangSamaDanBacaTahunTerbaru(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	repo := repository.APBDRepository{DB: tx}

	for _, sumber := range []usecase.APBDSumber{sumberAPBDUji("prv_apbd_9", 9, 39604499169821), sumberAPBDUji("prv_apbd_10", 10, 43931198916303)} {
		if err := repo.SimpanAPBD(t.Context(), sumber); err != nil {
			t.Fatal(err)
		}
	}

	bulan, err := repo.BulanRealisasiTersimpan(t.Context(), 2026)
	if err != nil || bulan != 10 {
		t.Fatalf("bulan %d err %v", bulan, err)
	}
	apbd, err := repo.BacaAPBD(t.Context(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if apbd.Tahun != 2026 || len(apbd.Akun) != 2 || apbd.Akun[0].Realisasi != 43931198916303 || apbd.Akun[1].Anggaran != 10164767414640.5 || apbd.ProvenanceID != "prv_apbd_10" {
		t.Errorf("apbd %+v", apbd)
	}
	if _, err := repo.BacaAPBD(t.Context(), 2019); !errors.Is(err, domain.ErrAPBDTidakDitemukan) {
		t.Errorf("tahun tanpa data: %v", err)
	}
}
