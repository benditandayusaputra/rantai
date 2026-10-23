package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobDayaTampung = "daya_tampung"

type BerkasDayaTampung struct {
	Jalur       domain.JalurDayaTampung
	TahunAjaran string
	Baris       []domain.BarisDayaTampung
	Provenance  domain.Provenance
}

type LaporanDayaTampung struct {
	Jalur           domain.JalurDayaTampung
	KursiDokumen    int
	KursiTerimpor   int
	TidakTerpetakan []domain.BarisDayaTampung
}

type PengambilDayaTampung interface {
	Ambil(ctx context.Context) ([]BerkasDayaTampung, error)
}

type PenyimpanDayaTampung interface {
	NPSNTerdaftar(ctx context.Context) (map[string]bool, error)
	SimpanDayaTampung(ctx context.Context, berkas []BerkasDayaTampung) (HasilImpor, error)
}

type PelaporDayaTampung interface {
	LaporkanDayaTampung(laporan []LaporanDayaTampung)
}

type ImporDayaTampung struct {
	Pengambil PengambilDayaTampung
	Penyimpan PenyimpanDayaTampung
	Pencatat  PencatatIngest
	Pelapor   PelaporDayaTampung
}

func (u ImporDayaTampung) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, JobDayaTampung)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", JobDayaTampung, err)
	}
	hasil, err := u.impor(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u ImporDayaTampung) impor(ctx context.Context) (HasilImpor, error) {
	semuaBerkas, err := u.Pengambil.Ambil(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("ambil berkas daya tampung: %w", err)
	}
	terdaftar, err := u.Penyimpan.NPSNTerdaftar(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("baca NPSN sekolah: %w", err)
	}
	terpetakan, laporan := pisahkanBerkas(semuaBerkas, terdaftar)
	hasil, err := u.Penyimpan.SimpanDayaTampung(ctx, terpetakan)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("simpan daya tampung: %w", err)
	}
	u.Pelapor.LaporkanDayaTampung(laporan)
	return hasil, nil
}

func pisahkanBerkas(semuaBerkas []BerkasDayaTampung, terdaftar map[string]bool) ([]BerkasDayaTampung, []LaporanDayaTampung) {
	terpetakan := make([]BerkasDayaTampung, len(semuaBerkas))
	laporan := make([]LaporanDayaTampung, len(semuaBerkas))
	for i, berkas := range semuaBerkas {
		cocok, tidak := domain.PisahkanTerpetakan(berkas.Baris, terdaftar)
		terpetakan[i] = berkas
		terpetakan[i].Baris = cocok
		laporan[i] = LaporanDayaTampung{
			Jalur:           berkas.Jalur,
			KursiDokumen:    domain.TotalKursi(berkas.Baris),
			KursiTerimpor:   domain.TotalKursi(cocok),
			TidakTerpetakan: tidak,
		}
	}
	return terpetakan, laporan
}
