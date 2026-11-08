package usecase

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pengambilAPBDTiruan struct {
	diminta []string
}

func (p *pengambilAPBDTiruan) AmbilAPBD(_ context.Context, tahun, bulanTerakhir int) (APBDSumber, error) {
	p.diminta = append(p.diminta, time.Date(tahun, time.Month(bulanTerakhir), 1, 0, 0, 0, 0, time.UTC).Format("2006-01"))
	return APBDSumber{APBD: domain.APBDDaerah{Tahun: tahun, BulanRealisasi: bulanTerakhir}}, nil
}

type penyimpanAPBDTiruan struct {
	tersimpan map[int]int
	disimpan  []int
}

func (p *penyimpanAPBDTiruan) BulanRealisasiTersimpan(_ context.Context, tahun int) (int, error) {
	return p.tersimpan[tahun], nil
}

func (p *penyimpanAPBDTiruan) SimpanAPBD(_ context.Context, sumber APBDSumber) error {
	p.disimpan = append(p.disimpan, sumber.APBD.Tahun)
	return nil
}

func TestImporAPBD_TahunLengkapDilewatiTahunBerjalanSampaiBulanIni(t *testing.T) {
	pengambil := &pengambilAPBDTiruan{}
	penyimpan := &penyimpanAPBDTiruan{tersimpan: map[int]int{2024: domain.BulanPenuhAPBD, 2026: 9}}
	pencatat := &pencatatTiruan{}
	impor := ImporAPBD{
		Pengambil: pengambil, Penyimpan: penyimpan, Pencatat: pencatat,
		Tahun:    []int{2024, 2025, 2026, 2027},
		Sekarang: func() time.Time { return time.Date(2026, time.November, 8, 2, 0, 0, 0, time.UTC) },
	}

	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(pengambil.diminta, []string{"2025-12", "2026-11"}) || !slices.Equal(penyimpan.disimpan, []int{2025, 2026}) {
		t.Errorf("diminta %v, disimpan %v", pengambil.diminta, penyimpan.disimpan)
	}
	if pencatat.selesai == nil || pencatat.selesai.JumlahBerubah != 2 {
		t.Errorf("hasil %+v", pencatat.selesai)
	}
}
