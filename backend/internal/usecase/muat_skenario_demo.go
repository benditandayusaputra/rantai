package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	JobSkenarioDemo      = "skenario_demo"
	prefiksKunciSimulasi = "simulasi/"
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
	SimpanAkunDemo(ctx context.Context, akun []domain.AkunDemo) error
	LokasiPaketSimulasi(ctx context.Context, kode string) (domain.Titik, error)
	SimpanLaporanSimulasi(ctx context.Context, data DataLaporanSimulasi) error
}

type LaporanSimulasiTersimpan struct {
	domain.LaporanSimulasi
	Foto []domain.FotoTersimpan
}

type DataLaporanSimulasi struct {
	PaketKode   string
	Titik       domain.Titik
	TitikPublik domain.Titik
	Laporan     []LaporanSimulasiTersimpan
}

type PelaporSkenarioDemo interface {
	LaporkanKecamatanDemoTerpilih(kode string)
}

type MuatSkenarioDemo struct {
	Skenario       domain.SkenarioDemo
	Penyimpan      PenyimpanSkenarioDemo
	Pelapor        PelaporSkenarioDemo
	BuatProvenance PembuatProvenance
	Pemroses       PemrosesFoto
	PenyimpanFoto  PenyimpanFoto
	BacaFoto       func(nama string) ([]byte, error)
}

func (u MuatSkenarioDemo) Jalankan(ctx context.Context) error {
	if err := u.Skenario.Validasi(); err != nil {
		return fmt.Errorf("skenario demo: %w", err)
	}
	kecamatan, err := u.kecamatanDemo(ctx)
	if err != nil {
		return err
	}
	err = u.Penyimpan.SimpanPaketSimulasi(ctx, PaketSimulasiTersimpan{
		PaketSimulasi: u.Skenario.PaketSimulasi,
		KodeKecamatan: kecamatan,
		Provenance: u.BuatProvenance(CatatanProvenance{
			Sumber: domain.SumberSimulasi, Label: domain.LabelSimulasi, Catatan: catatanPaketSimulasi,
		}),
	})
	if err != nil {
		return err
	}
	if err := u.Penyimpan.SimpanAkunDemo(ctx, u.Skenario.AkunDemo); err != nil {
		return err
	}
	return u.muatLaporanSimulasi(ctx)
}

func (u MuatSkenarioDemo) muatLaporanSimulasi(ctx context.Context) error {
	kode := u.Skenario.PaketSimulasi.Kode
	titik, err := u.Penyimpan.LokasiPaketSimulasi(ctx, kode)
	if err != nil {
		return fmt.Errorf("lokasi paket simulasi %s: %w", kode, err)
	}
	laporan := make([]LaporanSimulasiTersimpan, len(u.Skenario.LaporanSimulasi))
	for i, l := range u.Skenario.LaporanSimulasi {
		foto, err := u.unggahFotoSimulasi(ctx, prefiksKunciSimulasi+kode+"/"+strconv.Itoa(i+1), l.Foto)
		if err != nil {
			return fmt.Errorf("foto laporan simulasi %d: %w", i+1, err)
		}
		laporan[i] = LaporanSimulasiTersimpan{LaporanSimulasi: l, Foto: foto}
	}
	return u.Penyimpan.SimpanLaporanSimulasi(ctx, DataLaporanSimulasi{
		PaketKode: kode, Titik: titik, TitikPublik: domain.BulatkanLokasiPublik(titik), Laporan: laporan,
	})
}

func (u MuatSkenarioDemo) unggahFotoSimulasi(ctx context.Context, dasarKunci string, berkas []string) ([]domain.FotoTersimpan, error) {
	tersimpan := make([]domain.FotoTersimpan, len(berkas))
	for i, nama := range berkas {
		isi, err := u.BacaFoto(nama)
		if err != nil {
			return nil, fmt.Errorf("baca %s: %w", nama, err)
		}
		diproses, err := u.Pemroses.Proses(isi)
		if err != nil {
			return nil, fmt.Errorf("proses %s: %w", nama, err)
		}
		tersimpan[i] = kunciFoto(dasarKunci+"-"+strconv.Itoa(i+1), diproses)
		if err := unggahFoto(ctx, u.PenyimpanFoto, tersimpan[i], diproses); err != nil {
			return nil, err
		}
	}
	return tersimpan, nil
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
