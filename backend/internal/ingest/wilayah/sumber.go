package wilayah

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const formatTanggalRaw = "2006-01-02"

type KonfigurasiSumber struct {
	URLDasar       string        `yaml:"url_dasar"`
	Berkas         []string      `yaml:"berkas"`
	URLRujukan     string        `yaml:"url_rujukan"`
	Catatan        string        `yaml:"catatan"`
	UserAgent      string        `yaml:"user_agent"`
	JedaPermintaan time.Duration `yaml:"jeda_permintaan"`
}

type Sumber struct {
	Konfigurasi KonfigurasiSumber
	Klien       *http.Client
	DirRaw      string
	Sekarang    func() time.Time
}

func (s Sumber) Ambil(ctx context.Context) (usecase.DataSumberWilayah, error) {
	diambilPada := s.Sekarang()
	daftar := []domain.Wilayah{}
	semuaIsi := []byte{}
	for urutan, berkas := range s.Konfigurasi.Berkas {
		isi, err := s.ambilBerkas(ctx, urutan, berkas, diambilPada)
		if err != nil {
			return usecase.DataSumberWilayah{}, err
		}
		wilayahBerkas, err := UraiSQL(isi)
		if err != nil {
			return usecase.DataSumberWilayah{}, fmt.Errorf("urai %s: %w", berkas, err)
		}
		daftar = append(daftar, wilayahBerkas...)
		semuaIsi = append(semuaIsi, isi...)
	}
	return usecase.DataSumberWilayah{Wilayah: daftar, Provenance: s.catatProvenance(semuaIsi, diambilPada)}, nil
}

func (s Sumber) ambilBerkas(ctx context.Context, urutan int, berkas string, diambilPada time.Time) ([]byte, error) {
	if urutan > 0 {
		if err := tunggu(ctx, s.Konfigurasi.JedaPermintaan); err != nil {
			return nil, err
		}
	}
	isi, err := s.unduh(ctx, s.Konfigurasi.URLDasar+"/"+berkas)
	if err != nil {
		return nil, err
	}
	return isi, s.simpanMentah(path.Base(berkas), isi, diambilPada)
}

func (s Sumber) unduh(ctx context.Context, url string) ([]byte, error) {
	permintaan, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	permintaan.Header.Set("User-Agent", s.Konfigurasi.UserAgent)
	respons, err := s.Klien.Do(permintaan)
	if err != nil {
		return nil, fmt.Errorf("unduh %s: %w", url, err)
	}
	defer func() { _ = respons.Body.Close() }()
	if respons.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unduh %s: status %d", url, respons.StatusCode)
	}
	return io.ReadAll(respons.Body)
}

func (s Sumber) simpanMentah(nama string, isi []byte, diambilPada time.Time) error {
	dir := filepath.Join(s.DirRaw, usecase.JobWilayah, diambilPada.Format(formatTanggalRaw))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("buat folder data mentah: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, nama), isi, 0o644)
}

func (s Sumber) catatProvenance(isi []byte, diambilPada time.Time) domain.Provenance {
	return provenance.Baru(provenance.Catatan{
		Sumber:  domain.SumberKemendagri,
		URL:     s.Konfigurasi.URLRujukan,
		Label:   domain.LabelResmi,
		Catatan: s.Konfigurasi.Catatan,
	}, isi, diambilPada)
}

func tunggu(ctx context.Context, jeda time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(jeda):
		return nil
	}
}
