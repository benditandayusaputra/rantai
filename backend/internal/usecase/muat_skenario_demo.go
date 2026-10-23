package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	JobSkenarioDemo      = "skenario_demo"
	catatanPaketSimulasi = "Paket simulasi untuk peragaan Lapor Bukti; penyedia dan instansi contoh, bukan data pengadaan sungguhan"
)

type PaketSimulasiTersimpan struct {
	domain.PaketSimulasi
	KodeKecamatan string
	Provenance    domain.Provenance
}

type PenyimpanSkenarioDemo interface {
	PilihKecamatanDemo(ctx context.Context, kodeKota string) (string, error)
	SimpanPaketSimulasi(ctx context.Context, paket PaketSimulasiTersimpan) error
}

type PelaporSkenarioDemo interface {
	LaporkanKecamatanDemoTerpilih(kode string)
}

type MuatSkenarioDemo struct {
	Skenario       domain.SkenarioDemo
	Penyimpan      PenyimpanSkenarioDemo
	Pelapor        PelaporSkenarioDemo
	BuatProvenance PembuatProvenance
}

func (u MuatSkenarioDemo) Jalankan(ctx context.Context) error {
	if err := u.Skenario.Validasi(); err != nil {
		return fmt.Errorf("skenario demo: %w", err)
	}
	kecamatan, err := u.kecamatanDemo(ctx)
	if err != nil {
		return err
	}
	return u.Penyimpan.SimpanPaketSimulasi(ctx, PaketSimulasiTersimpan{
		PaketSimulasi: u.Skenario.PaketSimulasi,
		KodeKecamatan: kecamatan,
		Provenance: u.BuatProvenance(CatatanProvenance{
			Sumber: domain.SumberSimulasi, Label: domain.LabelSimulasi, Catatan: catatanPaketSimulasi,
		}),
	})
}

func (u MuatSkenarioDemo) kecamatanDemo(ctx context.Context) (string, error) {
	if u.Skenario.KecamatanDemo != "" {
		return u.Skenario.KecamatanDemo, nil
	}
	kode, err := u.Penyimpan.PilihKecamatanDemo(ctx, u.Skenario.KotaDemo)
	if err != nil {
		return "", fmt.Errorf("pilih kecamatan demo di %s: %w", u.Skenario.KotaDemo, err)
	}
	u.Pelapor.LaporkanKecamatanDemoTerpilih(kode)
	return kode, nil
}
