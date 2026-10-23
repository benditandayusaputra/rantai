package dprd

import (
	"context"
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

type Fraksi struct {
	Nama    string `yaml:"nama"`
	Halaman string `yaml:"halaman"`
}

type KonfigurasiSumber struct {
	URLDasar      string   `yaml:"url_dasar"`
	URLRujukan    string   `yaml:"url_rujukan"`
	Catatan       string   `yaml:"catatan"`
	HalamanDapil  string   `yaml:"halaman_dapil"`
	HalamanProfil string   `yaml:"halaman_profil"`
	Fraksi        []Fraksi `yaml:"fraksi"`
}

type KonfigurasiDapilKPU struct {
	TahunPemilu int            `yaml:"tahun_pemilu"`
	URL         string         `yaml:"url"`
	Catatan     string         `yaml:"catatan"`
	Daftar      []domain.Dapil `yaml:"daftar"`
}

type Konfigurasi struct {
	Sumber    KonfigurasiSumber   `yaml:"sumber"`
	Pengunduh klien.Konfigurasi   `yaml:"pengunduh"`
	URLLHKPN  string              `yaml:"url_lhkpn"`
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
	barisDapil, err := u.ambilSemuaDapil(ctx)
	if err != nil {
		return usecase.DataSumberDPRD{}, err
	}
	anggotaFraksi, err := u.ambilSemuaFraksi(ctx)
	if err != nil {
		return usecase.DataSumberDPRD{}, err
	}
	return usecase.DataSumberDPRD{
		BarisDapil:        barisDapil,
		AnggotaFraksi:     anggotaFraksi,
		ProvenanceDapil:   s.provenanceDapil(u.diambilPada),
		ProvenanceAnggota: s.provenanceAnggota(u.semuaIsi, u.diambilPada),
	}, nil
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
	return isi, u.simpanMentah(halaman, isi)
}

func (u *pengunduh) simpanMentah(halaman string, isi []byte) error {
	dir := filepath.Join(u.DirRaw, usecase.JobKPUDPRD, u.diambilPada.Format(formatTanggalRaw))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("buat folder data mentah: %w", err)
	}
	nama := strings.ReplaceAll(strings.Trim(halaman, "/"), "/", "_") + ".html"
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
