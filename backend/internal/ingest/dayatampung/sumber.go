package dayatampung

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

var (
	ErrKolomTidakLengkap = errors.New("CSV daya tampung wajib punya kolom nama, npsn, dan kursi")
	kolomWajib           = []string{"nama", "npsn", "kursi"}
)

type KonfigurasiBerkas struct {
	Path    string                  `yaml:"path"`
	Jalur   domain.JalurDayaTampung `yaml:"jalur"`
	URL     string                  `yaml:"url"`
	Catatan string                  `yaml:"catatan"`
}

type Konfigurasi struct {
	TahunAjaran string              `yaml:"tahun_ajaran"`
	Berkas      []KonfigurasiBerkas `yaml:"berkas"`
}

type Sumber struct {
	Konfigurasi Konfigurasi
	Sekarang    func() time.Time
}

func (s Sumber) Ambil(context.Context) ([]usecase.BerkasDayaTampung, error) {
	hasil := []usecase.BerkasDayaTampung{}
	for _, cfg := range s.Konfigurasi.Berkas {
		if !cfg.Jalur.IsValid() {
			return nil, fmt.Errorf("jalur daya tampung tidak dikenal: %q", cfg.Jalur)
		}
		isi, err := os.ReadFile(cfg.Path)
		if err != nil {
			return nil, fmt.Errorf("baca %s: %w", cfg.Path, err)
		}
		baris, err := UraiCSV(isi)
		if err != nil {
			return nil, fmt.Errorf("urai %s: %w", cfg.Path, err)
		}
		hasil = append(hasil, usecase.BerkasDayaTampung{
			Jalur: cfg.Jalur, TahunAjaran: s.Konfigurasi.TahunAjaran, Baris: baris,
			Provenance: provenance.Baru(provenance.Catatan{
				Sumber: domain.SumberDisdik, URL: cfg.URL, Label: domain.LabelResmi, Catatan: cfg.Catatan,
			}, isi, s.Sekarang()),
		})
	}
	return hasil, nil
}

func UraiCSV(isi []byte) ([]domain.BarisDayaTampung, error) {
	semua, err := csv.NewReader(bytes.NewReader(isi)).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(semua) == 0 {
		return nil, ErrKolomTidakLengkap
	}
	posisi, err := posisiKolom(semua[0])
	if err != nil {
		return nil, err
	}
	hasil := []domain.BarisDayaTampung{}
	for nomor, rekaman := range semua[1:] {
		kursi, err := strconv.Atoi(strings.TrimSpace(rekaman[posisi["kursi"]]))
		if err != nil {
			return nil, fmt.Errorf("baris %d: kursi bukan angka: %w", nomor+2, err)
		}
		hasil = append(hasil, domain.BarisDayaTampung{
			NPSN:  strings.TrimSpace(rekaman[posisi["npsn"]]),
			Nama:  strings.TrimSpace(rekaman[posisi["nama"]]),
			Kursi: kursi,
		})
	}
	return hasil, nil
}

func posisiKolom(judul []string) (map[string]int, error) {
	posisi := map[string]int{}
	for i, kolom := range judul {
		posisi[strings.TrimSpace(kolom)] = i
	}
	for _, kolom := range kolomWajib {
		if _, ada := posisi[kolom]; !ada {
			return nil, fmt.Errorf("%w: kolom %q tidak ada", ErrKolomTidakLengkap, kolom)
		}
	}
	return posisi, nil
}
