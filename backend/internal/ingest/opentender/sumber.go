package opentender

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const (
	formatTanggalRaw = "2006-01-02"
	penandaTahun     = "{tahun}"
)

type Konfigurasi struct {
	URLUnduhan string            `yaml:"url_unduhan"`
	URLRujukan string            `yaml:"url_rujukan"`
	Atribusi   string            `yaml:"atribusi"`
	Tahun      []int             `yaml:"tahun"`
	IDPembeli  []string          `yaml:"id_pembeli"`
	Klien      klien.Konfigurasi `yaml:"klien"`
}

type Pengunduh interface {
	Ambil(ctx context.Context, alamat string) ([]byte, error)
}

type Sumber struct {
	Konfigurasi Konfigurasi
	Klien       Pengunduh
	DirRaw      string
	Sekarang    func() time.Time
}

func (s Sumber) AmbilRiwayat(ctx context.Context, tahun int) (usecase.RiwayatSumber, error) {
	diambilPada := s.Sekarang()
	alamat := strings.ReplaceAll(s.Konfigurasi.URLUnduhan, penandaTahun, strconv.Itoa(tahun))
	isi, err := s.Klien.Ambil(ctx, alamat)
	if err != nil {
		return usecase.RiwayatSumber{}, err
	}
	if err := s.simpanMentah(tahun, isi, diambilPada); err != nil {
		return usecase.RiwayatSumber{}, err
	}
	riwayat, err := UraiOCDS(bytes.NewReader(isi), s.Konfigurasi.IDPembeli)
	if err != nil {
		return usecase.RiwayatSumber{}, fmt.Errorf("opentender %d: %w", tahun, err)
	}
	prov := provenance.Baru(provenance.Catatan{
		Sumber: domain.SumberOpentender, URL: s.Konfigurasi.URLRujukan, Label: domain.LabelResmi, Catatan: s.Konfigurasi.Atribusi,
	}, isi, diambilPada)
	return usecase.RiwayatSumber{Riwayat: riwayat, Provenance: prov}, nil
}

func (s Sumber) simpanMentah(tahun int, isi []byte, diambilPada time.Time) error {
	dir := filepath.Join(s.DirRaw, usecase.JobOpentender, diambilPada.Format(formatTanggalRaw))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("buat folder data mentah: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, strconv.Itoa(tahun)+".jsonl.gz"), isi, 0o644)
}
