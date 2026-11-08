package djpk

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const (
	formatTanggalRaw = "2006-01-02"
	catatanAPBD      = "Ekspor resmi Portal Data SIKD DJPK Kemenkeu untuk APBD murni Pemerintah Provinsi DKI Jakarta. Realisasi kumulatif sampai bulan terakhir yang dilaporkan pemda ke SIKD."
)

type Konfigurasi struct {
	URLEkspor    string            `yaml:"url_ekspor"`
	URLRujukan   string            `yaml:"url_rujukan"`
	KodeProvinsi string            `yaml:"kode_provinsi"`
	KodePemda    string            `yaml:"kode_pemda"`
	Tahun        []int             `yaml:"tahun_anggaran"`
	Klien        klien.Konfigurasi `yaml:"klien"`
	Aturan       domain.AturanAPBD `yaml:"aturan"`
}

type Pengunduh interface {
	Ambil(ctx context.Context, alamat string) ([]byte, error)
}

type Sumber struct {
	Konfigurasi Konfigurasi
	Pengunduh   Pengunduh
	DirRaw      string
	Sekarang    func() time.Time
}

func (s Sumber) AmbilAPBD(ctx context.Context, tahun, bulanTerakhir int) (usecase.APBDSumber, error) {
	diambilPada := s.Sekarang()
	terbaru, isi, err := s.ambilBulan(ctx, tahun, bulanTerakhir)
	if err != nil {
		return usecase.APBDSumber{}, err
	}
	bulan, err := s.cariBulanRealisasi(ctx, tahun, bulanTerakhir, terbaru)
	if err != nil {
		return usecase.APBDSumber{}, err
	}
	if err := s.simpanMentah(tahun, bulan, isi, diambilPada); err != nil {
		return usecase.APBDSumber{}, err
	}
	catatan := provenance.Catatan{Sumber: domain.SumberDJPK, URL: s.alamat(s.Konfigurasi.URLRujukan, tahun, bulan), Label: domain.LabelResmi, Catatan: catatanAPBD}
	prov := provenance.Baru(catatan, isi, diambilPada)
	apbd := domain.APBDDaerah{Tahun: tahun, BulanRealisasi: bulan, Akun: terbaru, ProvenanceID: prov.ID}
	return usecase.APBDSumber{APBD: apbd, Provenance: prov}, nil
}

func (s Sumber) cariBulanRealisasi(ctx context.Context, tahun, bulan int, terbaru []domain.AkunAPBD) (int, error) {
	for ; bulan > 1; bulan-- {
		sebelumnya, _, err := s.ambilBulan(ctx, tahun, bulan-1)
		if err != nil {
			return 0, err
		}
		if !domain.IsRealisasiAPBDSama(terbaru, sebelumnya) {
			return bulan, nil
		}
	}
	return bulan, nil
}

func (s Sumber) ambilBulan(ctx context.Context, tahun, bulan int) ([]domain.AkunAPBD, []byte, error) {
	isi, err := s.Pengunduh.Ambil(ctx, s.alamat(s.Konfigurasi.URLEkspor, tahun, bulan))
	if err != nil {
		return nil, nil, fmt.Errorf("ambil apbd %d bulan %d: %w", tahun, bulan, err)
	}
	akun, err := UraiAPBD(isi)
	return akun, isi, err
}

func (s Sumber) alamat(dasar string, tahun, bulan int) string {
	kueri := url.Values{
		"type": {"apbd"}, "periode": {strconv.Itoa(bulan)}, "tahun": {strconv.Itoa(tahun)},
		"provinsi": {s.Konfigurasi.KodeProvinsi}, "pemda": {s.Konfigurasi.KodePemda},
	}
	return dasar + "?" + kueri.Encode()
}

func (s Sumber) simpanMentah(tahun, bulan int, isi []byte, diambilPada time.Time) error {
	dir := filepath.Join(s.DirRaw, usecase.JobAPBD, diambilPada.Format(formatTanggalRaw))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("buat folder data mentah: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, fmt.Sprintf("apbd_%d_%02d.xls", tahun, bulan)), isi, 0o644)
}
