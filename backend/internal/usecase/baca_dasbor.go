package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaDasbor interface {
	RekapKursi(ctx context.Context, kodeTema string) (domain.DasborLayanan, error)
	PaketTahunan(ctx context.Context) (domain.Seri, error)
	RincianUang(ctx context.Context, tahun int) (domain.DasborUang, error)
	KegiatanTahunan(ctx context.Context) (domain.Seri, error)
	RincianTanggungJawab(ctx context.Context, tahun int, instansiTanpaRincian string) (domain.DasborTanggungJawab, error)
}

type BacaDasbor struct {
	Tema                 BacaTema
	Pembaca              PembacaDasbor
	Sekarang             func() time.Time
	InstansiTanpaRincian string
}

func (u BacaDasbor) Layanan(ctx context.Context, kodeTema string) (domain.DasborLayanan, error) {
	if _, err := u.Tema.Ambil(ctx, kodeTema); err != nil {
		return domain.DasborLayanan{}, err
	}
	d, err := u.Pembaca.RekapKursi(ctx, kodeTema)
	if err != nil {
		return domain.DasborLayanan{}, fmt.Errorf("rekap kursi tema %s: %w", kodeTema, err)
	}
	return d, nil
}

func (u BacaDasbor) Uang(ctx context.Context, tahun int) (domain.DasborUang, error) {
	tahunan, err := u.Pembaca.PaketTahunan(ctx)
	if err != nil {
		return domain.DasborUang{}, fmt.Errorf("paket tahunan: %w", err)
	}
	tahun = u.pilihTahun(tahun, tahunan)
	d, err := u.Pembaca.RincianUang(ctx, tahun)
	if err != nil {
		return domain.DasborUang{}, fmt.Errorf("rincian uang %d: %w", tahun, err)
	}
	d.Tahun, d.Tahunan = tahun, tahunan
	d.Bulanan.Isi = domain.LengkapiBulan(tahun, d.Bulanan.Isi)
	return d, nil
}

func (u BacaDasbor) TanggungJawab(ctx context.Context, tahun int) (domain.DasborTanggungJawab, error) {
	tahunan, err := u.Pembaca.KegiatanTahunan(ctx)
	if err != nil {
		return domain.DasborTanggungJawab{}, fmt.Errorf("kegiatan tahunan: %w", err)
	}
	tahun = u.pilihTahun(tahun, tahunan)
	d, err := u.Pembaca.RincianTanggungJawab(ctx, tahun, u.InstansiTanpaRincian)
	if err != nil {
		return domain.DasborTanggungJawab{}, fmt.Errorf("rincian tanggung jawab %d: %w", tahun, err)
	}
	d.Tahun, d.KegiatanTahunan = tahun, tahunan
	d.KegiatanBulanan.Isi = domain.LengkapiBulan(tahun, d.KegiatanBulanan.Isi)
	return d, nil
}

func (u BacaDasbor) pilihTahun(diminta int, tahunan domain.Seri) int {
	if diminta != 0 {
		return diminta
	}
	if tahun := domain.TahunTerlengkap(tahunan.Isi); tahun != 0 {
		return tahun
	}
	return u.Sekarang().Year()
}
