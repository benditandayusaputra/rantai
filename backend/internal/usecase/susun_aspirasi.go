package usecase

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PermintaanAspirasi struct {
	KodeKecamatan string
	Tema          string
	Masalah       string
	Bahasa        string
}

type PembacaTemaAspirasi interface {
	Ambil(ctx context.Context, kode string) (domain.TemaTersimpan, error)
}

type PembacaDetailWilayah interface {
	Ambil(ctx context.Context, kode string) (domain.DetailWilayah, error)
}

type PenghitungBelanja interface {
	Hitung(ctx context.Context, kodeTema string, tahun int) (domain.BelanjaTema, error)
}

type PembacaSumberAspirasi interface {
	Ambil(ctx context.Context, id string) (domain.Provenance, error)
}

type SusunAspirasi struct {
	Tema       PembacaTemaAspirasi
	Wilayah    PembacaDetailWilayah
	Belanja    PenghitungBelanja
	Wakil      PembacaWakilKecamatan
	Provenance PembacaSumberAspirasi
	Jalur      domain.KonfigurasiJalurAspirasi
	Kosakata   domain.KosakataAspirasi
}

func (u SusunAspirasi) Susun(ctx context.Context, p PermintaanAspirasi) (domain.Dokumen, error) {
	if err := periksaBahasa(p.Bahasa); err != nil {
		return domain.Dokumen{}, err
	}
	bahan, err := u.temaDanMasalah(ctx, p)
	if err != nil {
		return domain.Dokumen{}, err
	}
	kecamatan, err := ambilKecamatanDanKota(ctx, u.Wilayah, p.KodeKecamatan)
	if err != nil {
		return domain.Dokumen{}, err
	}
	bahan.Kecamatan, bahan.Kota = kecamatan.Nama, kecamatan.kota
	bahan.Jalur = u.Jalur.Untuk(kecamatan.placeholder())
	if err := u.lengkapiData(ctx, &bahan, kecamatan.Kode); err != nil {
		return domain.Dokumen{}, err
	}
	if err := u.lengkapiPengawas(ctx, &bahan, kecamatan.Kode); err != nil {
		return domain.Dokumen{}, err
	}
	return domain.SusunDokumenAspirasi(bahan, u.Kosakata, p.Bahasa), nil
}

func periksaBahasa(bahasa string) error {
	if bahasa != domain.BahasaIndonesia && bahasa != domain.BahasaInggris {
		return domain.InputTidakValid("bahasa harus id atau en")
	}
	return nil
}

func (u SusunAspirasi) temaDanMasalah(ctx context.Context, p PermintaanAspirasi) (domain.BahanAspirasi, error) {
	tema, err := u.Tema.Ambil(ctx, p.Tema)
	if err != nil {
		return domain.BahanAspirasi{}, err
	}
	i := slices.IndexFunc(tema.TanggungJawab, func(tj domain.TanggungJawab) bool { return tj.Masalah == p.Masalah })
	if i < 0 {
		return domain.BahanAspirasi{}, domain.InputTidakValid("masalah tidak dikenal untuk tema " + tema.Kode)
	}
	return domain.BahanAspirasi{Tema: tema, Masalah: tema.TanggungJawab[i]}, nil
}

func (u SusunAspirasi) lengkapiData(ctx context.Context, bahan *domain.BahanAspirasi, kodeKecamatan string) error {
	belanja, err := u.Belanja.Hitung(ctx, bahan.Tema.Kode, 0)
	if err != nil {
		return err
	}
	bahan.TahunAnggaran = belanja.Tahun
	i := slices.IndexFunc(belanja.Kecamatan, func(b domain.BelanjaKecamatan) bool { return b.Kecamatan.Kode == kodeKecamatan })
	if i < 0 {
		return nil
	}
	baris := belanja.Kecamatan[i]
	if bahan.Indikator, err = bersumber(ctx, u.Provenance, baris.Indikator, func(i domain.IndikatorKecamatan) string { return i.ProvenanceID }); err != nil {
		return err
	}
	bahan.Anggaran, err = bersumber(ctx, u.Provenance, baris.Anggaran, func(a domain.RingkasanAnggaran) string { return a.ProvenanceID })
	return err
}

func bersumber[T any](ctx context.Context, p PembacaSumberAspirasi, nilai *T, idProvenance func(T) string) (*domain.DataBersumber[T], error) {
	if nilai == nil {
		return nil, nil
	}
	prov, err := p.Ambil(ctx, idProvenance(*nilai))
	if err != nil {
		return nil, fmt.Errorf("sumber data aspirasi: %w", err)
	}
	return &domain.DataBersumber[T]{Nilai: *nilai, Provenance: prov}, nil
}

func (u SusunAspirasi) lengkapiPengawas(ctx context.Context, bahan *domain.BahanAspirasi, kodeKecamatan string) error {
	wakil, err := u.Wakil.WakilKecamatanKomisi(ctx, kodeKecamatan, []string{bahan.Masalah.Komisi})
	if err != nil {
		return err
	}
	bahan.NamaDapil = wakil.Nama
	bahan.AnggotaKomisi = []domain.Wakil{}
	for _, w := range wakil.Anggota {
		if w.IsDiKomisiPengawasTema {
			bahan.AnggotaKomisi = append(bahan.AnggotaKomisi, w)
		}
	}
	return nil
}

type kecamatanDenganKota struct {
	domain.DetailWilayah
	kota string
}

func (k kecamatanDenganKota) placeholder() map[string]string {
	return map[string]string{"kecamatan": k.Nama, "kota": k.kota}
}

func ambilKecamatanDanKota(ctx context.Context, wilayah PembacaDetailWilayah, kode string) (kecamatanDenganKota, error) {
	kodeKecamatan, err := kecamatanDariKode(kode)
	if err != nil {
		return kecamatanDenganKota{}, err
	}
	kecamatan, err := wilayah.Ambil(ctx, kodeKecamatan)
	if err != nil {
		return kecamatanDenganKota{}, err
	}
	hasil := kecamatanDenganKota{DetailWilayah: kecamatan}
	for _, leluhur := range kecamatan.Hierarki {
		if leluhur.Tingkat == domain.TingkatKota {
			hasil.kota = leluhur.Nama
		}
	}
	return hasil, nil
}

type JalurWilayah struct {
	Wilayah          domain.DetailWilayah
	DiverifikasiPada time.Time
	Jalur            []domain.JalurAspirasi
}

type BacaJalurAspirasi struct {
	Wilayah PembacaDetailWilayah
	Jalur   domain.KonfigurasiJalurAspirasi
}

func (u BacaJalurAspirasi) Untuk(ctx context.Context, kodeWilayah string) (JalurWilayah, error) {
	kecamatan, err := ambilKecamatanDanKota(ctx, u.Wilayah, kodeWilayah)
	if err != nil {
		return JalurWilayah{}, err
	}
	return JalurWilayah{
		Wilayah: kecamatan.DetailWilayah, DiverifikasiPada: u.Jalur.DiverifikasiPada, Jalur: u.Jalur.Untuk(kecamatan.placeholder()),
	}, nil
}
