package sirup

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const (
	formatTanggalRaw = "2006-01-02"
	catatanRUP       = "Detail Rencana Umum Pengadaan dari SiRUP LKPP. Jadwal pelaksanaan kontrak berketelitian bulan; tanggal awal dan akhir bulan dipakai sebagai batas."
)

type Konfigurasi struct {
	URLDetail string            `yaml:"url_detail"`
	Klien     klien.Konfigurasi `yaml:"klien"`
}

type Pengunduh interface {
	Ambil(ctx context.Context, alamat string) ([]byte, error)
}

type Sumber struct {
	Klien     Pengunduh
	URLDetail string
	DirRaw    string
	Sekarang  func() time.Time
}

func (s Sumber) AmbilRUP(ctx context.Context, kode string) (usecase.RUPSumber, error) {
	diambilPada := s.Sekarang()
	alamat := s.URLDetail + url.PathEscape(kode)
	isi, err := s.Klien.Ambil(ctx, alamat)
	if errors.Is(err, klien.ErrStatusTidakOK) {
		return usecase.RUPSumber{}, fmt.Errorf("rup %s: %w", kode, domain.ErrRUPTidakTersedia)
	}
	if err != nil {
		return usecase.RUPSumber{}, err
	}
	if err := s.simpanMentah(kode, isi, diambilPada); err != nil {
		return usecase.RUPSumber{}, err
	}
	rup, err := UraiRUP(isi)
	if err != nil {
		return usecase.RUPSumber{}, fmt.Errorf("rup %s: %w", kode, err)
	}
	prov := provenance.Baru(provenance.Catatan{Sumber: domain.SumberSIRUP, URL: alamat, Label: domain.LabelResmi, Catatan: catatanRUP}, isi, diambilPada)
	return usecase.RUPSumber{RUP: rup, Provenance: prov}, nil
}

func (s Sumber) simpanMentah(kode string, isi []byte, diambilPada time.Time) error {
	dir := filepath.Join(s.DirRaw, usecase.JobSIRUP, diambilPada.Format(formatTanggalRaw))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("buat folder data mentah: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, kode+".html"), isi, 0o644)
}
