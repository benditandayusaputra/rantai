package paripurna

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const formatTanggal = time.DateOnly

var (
	ErrKolomTidakLengkap = errors.New("CSV paripurna tidak punya kolom wajib")
	kolomSikap           = []string{"fraksi", "agenda", "judul_agenda", "tanggal", "ringkasan", "url"}
	kolomAspirasi        = []string{"masa_reses", "tanggal_paripurna", "bidang", "uraian", "kode_wilayah", "url"}
)

type Konfigurasi struct {
	PathSikapFraksi   string `yaml:"path_sikap_fraksi"`
	PathAspirasiReses string `yaml:"path_aspirasi_reses"`
	URLSumber         string `yaml:"url_sumber"`
	Catatan           string `yaml:"catatan"`
}

type Sumber struct {
	Konfigurasi Konfigurasi
	Sekarang    func() time.Time
}

type baris map[string]string

func (s Sumber) Ambil(context.Context) (usecase.DataParipurna, error) {
	isiSikap, err := os.ReadFile(s.Konfigurasi.PathSikapFraksi)
	if err != nil {
		return usecase.DataParipurna{}, fmt.Errorf("baca %s: %w", s.Konfigurasi.PathSikapFraksi, err)
	}
	isiAspirasi, err := os.ReadFile(s.Konfigurasi.PathAspirasiReses)
	if err != nil {
		return usecase.DataParipurna{}, fmt.Errorf("baca %s: %w", s.Konfigurasi.PathAspirasiReses, err)
	}
	sikap, err := UraiSikapFraksi(isiSikap)
	if err != nil {
		return usecase.DataParipurna{}, fmt.Errorf("urai sikap fraksi: %w", err)
	}
	aspirasi, err := UraiAspirasiReses(isiAspirasi)
	if err != nil {
		return usecase.DataParipurna{}, fmt.Errorf("urai aspirasi reses: %w", err)
	}
	catatan := provenance.Catatan{Sumber: domain.SumberDPRD, URL: s.Konfigurasi.URLSumber, Label: domain.LabelResmi, Catatan: s.Konfigurasi.Catatan}
	return usecase.DataParipurna{
		Sikap: sikap, Aspirasi: aspirasi, Provenance: provenance.Baru(catatan, append(isiSikap, isiAspirasi...), s.Sekarang()),
	}, nil
}

func UraiSikapFraksi(isi []byte) ([]domain.SikapFraksi, error) {
	return uraiCSV(isi, kolomSikap, func(b baris) (domain.SikapFraksi, error) {
		tanggal, err := time.Parse(formatTanggal, b["tanggal"])
		return domain.SikapFraksi{
			Fraksi: b["fraksi"], Agenda: domain.AgendaSikap(b["agenda"]), JudulAgenda: b["judul_agenda"],
			Tanggal: tanggal, Ringkasan: b["ringkasan"], URL: b["url"],
		}, err
	})
}

func UraiAspirasiReses(isi []byte) ([]domain.AspirasiReses, error) {
	return uraiCSV(isi, kolomAspirasi, func(b baris) (domain.AspirasiReses, error) {
		tanggal, err := time.Parse(formatTanggal, b["tanggal_paripurna"])
		return domain.AspirasiReses{
			MasaReses: b["masa_reses"], TanggalParipurna: tanggal, Bidang: b["bidang"],
			Uraian: b["uraian"], KodeWilayah: b["kode_wilayah"], URL: b["url"],
		}, err
	})
}

func uraiCSV[T any](isi []byte, kolomWajib []string, ubah func(baris) (T, error)) ([]T, error) {
	semua, err := csv.NewReader(bytes.NewReader(isi)).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(semua) == 0 {
		return nil, ErrKolomTidakLengkap
	}
	posisi, err := posisiKolom(semua[0], kolomWajib)
	if err != nil {
		return nil, err
	}
	hasil := []T{}
	for nomor, rekaman := range semua[1:] {
		b := baris{}
		for kolom, i := range posisi {
			b[kolom] = strings.TrimSpace(rekaman[i])
		}
		entri, err := ubah(b)
		if err != nil {
			return nil, fmt.Errorf("baris %d: %w", nomor+2, err)
		}
		hasil = append(hasil, entri)
	}
	return hasil, nil
}

func posisiKolom(judul, kolomWajib []string) (map[string]int, error) {
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
