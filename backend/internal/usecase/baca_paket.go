package usecase

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type FilterPaketKecamatan struct {
	KodeKecamatan string
	Tema          string
	Tahun         int
	Halaman       int
	PerHalaman    int
}

type FilterPetaPaket struct {
	Kotak     domain.Kotak
	Tema      string
	Tahun     int
	Toleransi float64
}

type PembacaPaket interface {
	AmbilPaket(ctx context.Context, kode string) (domain.PaketLengkap, bool, error)
	DaftarPaketKecamatan(ctx context.Context, f FilterPaketKecamatan) ([]domain.PaketRingkas, int, error)
	PetaPaket(ctx context.Context, f FilterPetaPaket) ([]domain.TitikPaket, []domain.AreaPaket, error)
	PaketTerdekat(ctx context.Context, titik domain.Titik, radiusMeter float64, batas int) ([]domain.TitikPaket, error)
}

type AturanBacaPaket struct {
	PerHalaman          int     `yaml:"per_halaman"`
	RadiusBawaanM       float64 `yaml:"radius_bawaan_m"`
	RadiusMaksM         float64 `yaml:"radius_maks_m"`
	BatasTerdekat       int     `yaml:"batas_terdekat"`
	ToleransiPeta       float64 `yaml:"toleransi_peta"`
	LuasKotakMaks       float64 `yaml:"luas_kotak_maks_derajat"`
	URLSPSE             string  `yaml:"url_spse"`
	URLSIRUP            string  `yaml:"url_sirup"`
	BatasDaftarPenyedia int     `yaml:"batas_daftar_penyedia"`
}

type BacaPaket struct {
	Pembaca      PembacaPaket
	Aturan       AturanBacaPaket
	Pemeliharaan KonfigurasiPemeliharaan
	Sekarang     func() time.Time
}

type DetailPaket struct {
	Paket        domain.PaketLengkap
	Rantai       []domain.MataRantai
	Pemeliharaan HasilPemeliharaan
	URLSumber    string
	URLRUP       string
}

type HalamanPaket struct {
	Paket      []domain.PaketRingkas
	Halaman    int
	PerHalaman int
	Total      int
}

type PetaPaket struct {
	Titik []domain.TitikPaket
	Area  []domain.AreaPaket
}

func (u BacaPaket) Ambil(ctx context.Context, kode string) (DetailPaket, error) {
	p, ada, err := u.Pembaca.AmbilPaket(ctx, kode)
	if err != nil {
		return DetailPaket{}, fmt.Errorf("ambil paket %s: %w", kode, err)
	}
	if !ada {
		return DetailPaket{}, domain.TidakDitemukan("paket tidak ditemukan")
	}
	pemeliharaan := u.estimasiPemeliharaan(p)
	return DetailPaket{
		Paket:        p,
		Rantai:       domain.SusunRantai(p, pemeliharaan.Estimasi),
		Pemeliharaan: pemeliharaan,
		URLSumber:    u.urlSumber(p),
		URLRUP:       u.urlRUP(p),
	}, nil
}

func (u BacaPaket) urlSumber(p domain.PaketLengkap) string {
	if p.IsSimulation {
		return ""
	}
	if p.IsTender {
		return u.Aturan.URLSPSE + "/lelang/" + p.KodeSPSE + "/pengumumanlelang"
	}
	return u.Aturan.URLSPSE + "/nontender/" + p.KodeSPSE + "/pengumumanpl"
}

func (u BacaPaket) urlRUP(p domain.PaketLengkap) string {
	if p.RUP == nil || p.IsSimulation {
		return ""
	}
	return u.Aturan.URLSIRUP + p.RUP.Kode
}

func (u BacaPaket) estimasiPemeliharaan(p domain.PaketLengkap) HasilPemeliharaan {
	if !slices.Contains(u.Pemeliharaan.JenisPengadaan, p.JenisPengadaan) {
		return HasilPemeliharaan{Estimasi: domain.EstimasiPemeliharaan{Status: domain.StatusTidakDapatDiestimasi, Label: domain.LabelEstimasi}}
	}
	jadwal := domain.JadwalPaket{Kategori: p.Kategori, TanggalSerahTerima: p.TanggalSerahTerima, TanggalKontrak: p.TanggalKontrak}
	if p.RUP != nil {
		jadwal.JadwalSelesai = p.RUP.PelaksanaanSelesai
	}
	return HasilPemeliharaan{IsBerlaku: true, Estimasi: u.Pemeliharaan.Estimasi(jadwal, u.Sekarang())}
}

func (u BacaPaket) DaftarKecamatan(ctx context.Context, f FilterPaketKecamatan) (HalamanPaket, error) {
	f.PerHalaman = u.Aturan.PerHalaman
	f.Halaman = max(f.Halaman, 1)
	daftar, total, err := u.Pembaca.DaftarPaketKecamatan(ctx, f)
	if err != nil {
		return HalamanPaket{}, err
	}
	return HalamanPaket{Paket: daftar, Halaman: f.Halaman, PerHalaman: f.PerHalaman, Total: total}, nil
}

func (u BacaPaket) Peta(ctx context.Context, f FilterPetaPaket) (PetaPaket, error) {
	if err := f.Kotak.Periksa(u.Aturan.LuasKotakMaks); err != nil {
		return PetaPaket{}, err
	}
	f.Toleransi = u.Aturan.ToleransiPeta
	titik, area, err := u.Pembaca.PetaPaket(ctx, f)
	if err != nil {
		return PetaPaket{}, err
	}
	return PetaPaket{Titik: titik, Area: area}, nil
}

func (u BacaPaket) Terdekat(ctx context.Context, titik domain.Titik, radiusMeter float64) ([]domain.TitikPaket, error) {
	if radiusMeter == 0 {
		radiusMeter = u.Aturan.RadiusBawaanM
	}
	if radiusMeter < 0 || radiusMeter > u.Aturan.RadiusMaksM {
		return nil, domain.InputTidakValid(fmt.Sprintf("radius_m wajib di antara 1 dan %.0f", u.Aturan.RadiusMaksM))
	}
	return u.Pembaca.PaketTerdekat(ctx, titik, radiusMeter, u.Aturan.BatasTerdekat)
}
