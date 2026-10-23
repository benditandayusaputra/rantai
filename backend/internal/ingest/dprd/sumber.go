package dprd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const formatTanggalRaw = "2006-01-02"

type Halaman struct {
	Nama    string `yaml:"nama"`
	Halaman string `yaml:"halaman"`
}

type HalamanKomisi struct {
	Kode    string `yaml:"kode"`
	Halaman `yaml:",inline"`
}

type KonfigurasiSumber struct {
	URLDasar        string          `yaml:"url_dasar"`
	URLRujukan      string          `yaml:"url_rujukan"`
	Catatan         string          `yaml:"catatan"`
	HalamanDapil    string          `yaml:"halaman_dapil"`
	HalamanProfil   string          `yaml:"halaman_profil"`
	Fraksi          []Halaman       `yaml:"fraksi"`
	Komisi          []HalamanKomisi `yaml:"komisi"`
	KataAwalUnit    []string        `yaml:"kata_awal_unit"`
	AlatKelengkapan []Halaman       `yaml:"alat_kelengkapan"`
}

type KonfigurasiDapilKPU struct {
	TahunPemilu int            `yaml:"tahun_pemilu"`
	URL         string         `yaml:"url"`
	Catatan     string         `yaml:"catatan"`
	Daftar      []domain.Dapil `yaml:"daftar"`
}

type KonfigurasiRilis struct {
	URLAPI       string        `yaml:"url_api"`
	URLRujukan   string        `yaml:"url_rujukan"`
	Catatan      string        `yaml:"catatan"`
	Kategori     []string      `yaml:"kategori"`
	PerHalaman   int           `yaml:"per_halaman"`
	BatasHalaman int           `yaml:"batas_halaman"`
	Mulai        time.Time     `yaml:"mulai"`
	Tumpang      time.Duration `yaml:"tumpang_tindih"`
	MediaNama    string        `yaml:"media_nama"`
	MediaDomain  string        `yaml:"media_domain"`
}

type Konfigurasi struct {
	Sumber    KonfigurasiSumber   `yaml:"sumber"`
	Pengunduh klien.Konfigurasi   `yaml:"pengunduh"`
	URLLHKPN  string              `yaml:"url_lhkpn"`
	Rilis     KonfigurasiRilis    `yaml:"rilis"`
	DapilKPU  KonfigurasiDapilKPU `yaml:"dapil_kpu"`
}

type Pengunduh interface {
	Ambil(ctx context.Context, url string) ([]byte, error)
}

type Sumber struct {
	Konfigurasi Konfigurasi
	Pengunduh   Pengunduh
	DirRaw      string
	Sekarang    func() time.Time
}

type pengunduh struct {
	Sumber
	diambilPada time.Time
	semuaIsi    []byte
}

func (s Sumber) Ambil(ctx context.Context) (usecase.DataSumberDPRD, error) {
	u := &pengunduh{Sumber: s, diambilPada: s.Sekarang()}
	data := usecase.DataSumberDPRD{}
	var err error
	if data.BarisDapil, err = u.ambilSemuaDapil(ctx); err != nil {
		return usecase.DataSumberDPRD{}, err
	}
	if data.AnggotaFraksi, err = u.ambilSemuaFraksi(ctx); err != nil {
		return usecase.DataSumberDPRD{}, err
	}
	if data.Komisi, data.Struktur, err = u.ambilSemuaKomisi(ctx); err != nil {
		return usecase.DataSumberDPRD{}, err
	}
	akd, err := u.ambilSemuaAlatKelengkapan(ctx)
	if err != nil {
		return usecase.DataSumberDPRD{}, err
	}
	data.Struktur = append(data.Struktur, akd...)
	data.ProvenanceDapil = s.provenanceDapil(u.diambilPada)
	data.ProvenanceAnggota = s.provenanceAnggota(u.semuaIsi, u.diambilPada)
	return data, nil
}

func (u *pengunduh) ambilSemuaKomisi(ctx context.Context) ([]domain.Komisi, []domain.HalamanStruktur, error) {
	komisi, halaman := []domain.Komisi{}, []domain.HalamanStruktur{}
	for _, k := range u.Konfigurasi.Sumber.Komisi {
		isi, err := u.ambil(ctx, k.Halaman.Halaman)
		if err != nil {
			return nil, nil, err
		}
		bidang, errBidang := UraiBidangKomisi(isi, u.Konfigurasi.Sumber.KataAwalUnit)
		baris, errBaris := UraiHalamanStruktur(isi)
		if err := errors.Join(errBidang, errBaris); err != nil {
			return nil, nil, fmt.Errorf("%s: %w", k.Nama, err)
		}
		komisi = append(komisi, domain.Komisi{Kode: k.Kode, Nama: k.Nama, Bidang: bidang})
		halaman = append(halaman, domain.HalamanStruktur{AlatKelengkapan: k.Nama, KodeKomisi: k.Kode, Baris: baris})
	}
	return komisi, halaman, nil
}

func (u *pengunduh) ambilSemuaAlatKelengkapan(ctx context.Context) ([]domain.HalamanStruktur, error) {
	halaman := []domain.HalamanStruktur{}
	for _, a := range u.Konfigurasi.Sumber.AlatKelengkapan {
		isi, err := u.ambil(ctx, a.Halaman)
		if err != nil {
			return nil, err
		}
		baris, err := UraiHalamanStruktur(isi)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", a.Nama, err)
		}
		halaman = append(halaman, domain.HalamanStruktur{AlatKelengkapan: a.Nama, Baris: baris})
	}
	return halaman, nil
}

func (u *pengunduh) ambilSemuaDapil(ctx context.Context) ([]domain.BarisDapil, error) {
	semua := []domain.BarisDapil{}
	for _, d := range u.Konfigurasi.DapilKPU.Daftar {
		isi, err := u.ambil(ctx, fmt.Sprintf(u.Konfigurasi.Sumber.HalamanDapil, d.Nomor))
		if err != nil {
			return nil, err
		}
		baris, err := UraiHalamanDapil(isi, d.Nomor)
		if err != nil {
			return nil, err
		}
		semua = append(semua, baris...)
	}
	return semua, nil
}

func (u *pengunduh) ambilSemuaFraksi(ctx context.Context) ([]domain.AnggotaFraksi, error) {
	semua := []domain.AnggotaFraksi{}
	for _, f := range u.Konfigurasi.Sumber.Fraksi {
		isi, err := u.ambil(ctx, f.Halaman)
		if err != nil {
			return nil, err
		}
		baris, err := UraiHalamanStruktur(isi)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f.Nama, err)
		}
		for _, b := range baris {
			semua = append(semua, domain.AnggotaFraksi{Fraksi: f.Nama, BarisStruktur: b})
		}
	}
	return semua, nil
}

func (u *pengunduh) ambil(ctx context.Context, halaman string) ([]byte, error) {
	isi, err := u.Pengunduh.Ambil(ctx, u.Konfigurasi.Sumber.URLDasar+halaman)
	if err != nil {
		return nil, err
	}
	u.semuaIsi = append(u.semuaIsi, isi...)
	nama := strings.ReplaceAll(strings.Trim(halaman, "/"), "/", "_") + ".html"
	return isi, simpanMentah(filepath.Join(u.DirRaw, usecase.JobKPUDPRD, u.diambilPada.Format(formatTanggalRaw)), nama, isi)
}

func simpanMentah(dir, nama string, isi []byte) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("buat folder data mentah: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, nama), isi, 0o644)
}

func (s Sumber) provenanceDapil(diambilPada time.Time) domain.Provenance {
	kpu := s.Konfigurasi.DapilKPU
	return provenance.Baru(provenance.Catatan{
		Sumber: domain.SumberKPU, URL: kpu.URL, Label: domain.LabelResmi, Catatan: kpu.Catatan,
	}, []byte(fmt.Sprint(kpu.Daftar)), diambilPada)
}

func (s Sumber) provenanceAnggota(isi []byte, diambilPada time.Time) domain.Provenance {
	sumber := s.Konfigurasi.Sumber
	return provenance.Baru(provenance.Catatan{
		Sumber: domain.SumberDPRD, URL: sumber.URLRujukan, Label: domain.LabelResmi, Catatan: sumber.Catatan,
	}, isi, diambilPada)
}
