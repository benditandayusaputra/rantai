package lpse

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const (
	formatTanggalRaw = "2006-01-02"
	catatanDetail    = "Pengumuman, pemenang, dan jadwal paket non-tender SPSE LPSE Provinsi DKI Jakarta. Tanggal kontrak memakai akhir tahap Penandatanganan Kontrak di jadwal pemilihan."
	catatanDaftar    = "Daftar paket SPSE LPSE Provinsi DKI Jakarta. HPS paket tender berasal dari kolom daftar yang dibulatkan sumber ke satu desimal (contoh 5,8 M); detail tender tidak dibuka untuk publik."
)

type Konfigurasi struct {
	Klien                klien.Konfigurasi     `yaml:"klien"`
	URLDasar             string                `yaml:"url_dasar"`
	UkuranHalaman        int                   `yaml:"ukuran_halaman"`
	TahunAnggaran        []int                 `yaml:"tahun_anggaran"`
	JenisPengadaanDetail []string              `yaml:"jenis_pengadaan_detail"`
	Instansi             domain.AturanInstansi `yaml:"instansi"`
}

type Sumber struct {
	Klien    Pengirim
	URLDasar string
	Ukuran   int
	DirRaw   string
	Sekarang func() time.Time
}

func (s Sumber) AmbilDaftarPaket(ctx context.Context, tahun int) (usecase.DaftarPaketSumber, error) {
	diambilPada := s.Sekarang()
	perekam := &perekamRespons{Pengirim: s.Klien}
	daftar := SumberDaftar{Klien: perekam, URLDasar: s.URLDasar, UkuranHalaman: s.Ukuran}
	paket := []domain.PaketDaftar{}
	for _, jenis := range []Jenis{JenisTender, JenisNonTender} {
		baris, err := daftar.AmbilDaftar(ctx, jenis, tahun)
		if err != nil {
			return usecase.DaftarPaketSumber{}, err
		}
		paket = append(paket, keDaftarDomain(baris)...)
	}
	if err := s.simpanMentah(usecase.JobLPSEDaftar, diambilPada, perekam.isiBernama(strconv.Itoa(tahun))); err != nil {
		return usecase.DaftarPaketSumber{}, err
	}
	return usecase.DaftarPaketSumber{Paket: paket, Provenance: s.provenance(s.URLDasar+JenisNonTender.PathHalaman, catatanDaftar, perekam.gabungan(), diambilPada)}, nil
}

func (s Sumber) AmbilDetailPaket(ctx context.Context, kodeSPSE string) (usecase.DetailPaketSumber, error) {
	diambilPada := s.Sekarang()
	perekam := &perekamRespons{Pengirim: s.Klien}
	alamat := s.alamatNonTender("/nontender/", kodeSPSE, "/pengumumanpl")
	isi, err := perekam.Ambil(ctx, alamat)
	if err != nil {
		return usecase.DetailPaketSumber{}, err
	}
	pengumuman, err := UraiPengumuman(isi)
	if err != nil {
		return usecase.DetailPaketSumber{}, err
	}
	detail := keDetailDomain(pengumuman)
	if !detail.IsBatal {
		if err := s.lengkapiPemenang(ctx, perekam, &detail); err != nil {
			return usecase.DetailPaketSumber{}, err
		}
	}
	if err := s.simpanMentah(usecase.JobLPSEDetail, diambilPada, perekam.isiBernama(kodeSPSE)); err != nil {
		return usecase.DetailPaketSumber{}, err
	}
	return usecase.DetailPaketSumber{Detail: detail, Provenance: s.provenance(alamat, catatanDetail, perekam.gabungan(), diambilPada)}, nil
}

func (s Sumber) lengkapiPemenang(ctx context.Context, pengirim Pengirim, detail *domain.DetailPaket) error {
	pemenang, err := s.ambilPemenang(ctx, pengirim, detail.KodeSPSE)
	if err != nil || pemenang == nil {
		return err
	}
	detail.Pemenang = pemenang
	isi, err := pengirim.Ambil(ctx, s.alamatNonTender("/nontender/", detail.KodeSPSE, "/jadwal"))
	if errors.Is(err, klien.ErrTidakDitemukan) {
		return nil
	}
	if err != nil {
		return err
	}
	tanggal, ada, err := UraiTanggalKontrak(isi)
	detail.TanggalKontrak = nilaiBila(tanggal, ada)
	return err
}

func (s Sumber) ambilPemenang(ctx context.Context, klien Pengirim, kodeSPSE string) (*domain.PemenangPaket, error) {
	for _, halaman := range []struct {
		path         string
		isBerkontrak bool
	}{{"/pemenangberkontrak", true}, {"/pemenang", false}} {
		isi, err := klien.Ambil(ctx, s.alamatNonTender("/evaluasinontender/", kodeSPSE, halaman.path))
		if err != nil {
			return nil, err
		}
		daftar, err := UraiPemenang(isi)
		if err != nil {
			return nil, err
		}
		if len(daftar) > 0 {
			return &domain.PemenangPaket{NamaPenyedia: daftar[0].Nama, Nilai: daftar[0].Nilai, IsBerkontrak: halaman.isBerkontrak}, nil
		}
	}
	return nil, nil
}

func (s Sumber) alamatNonTender(awalan, kodeSPSE, akhiran string) string {
	return s.URLDasar + awalan + url.PathEscape(kodeSPSE) + akhiran
}

func (s Sumber) provenance(alamat, catatan string, isi []byte, diambilPada time.Time) domain.Provenance {
	return provenance.Baru(provenance.Catatan{Sumber: domain.SumberLPSE, URL: alamat, Label: domain.LabelResmi, Catatan: catatan}, isi, diambilPada)
}

func (s Sumber) simpanMentah(job string, diambilPada time.Time, berkas map[string][]byte) error {
	dir := filepath.Join(s.DirRaw, job, diambilPada.Format(formatTanggalRaw))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("buat folder data mentah: %w", err)
	}
	for nama, isi := range berkas {
		if err := os.WriteFile(filepath.Join(dir, nama), isi, 0o644); err != nil {
			return fmt.Errorf("simpan data mentah %s: %w", nama, err)
		}
	}
	return nil
}

type perekamRespons struct {
	Pengirim
	respons [][]byte
}

func (p *perekamRespons) Ambil(ctx context.Context, alamat string) ([]byte, error) {
	return p.rekam(p.Pengirim.Ambil(ctx, alamat))
}

func (p *perekamRespons) KirimForm(ctx context.Context, alamat string, form url.Values) ([]byte, error) {
	return p.rekam(p.Pengirim.KirimForm(ctx, alamat, form))
}

func (p *perekamRespons) rekam(isi []byte, err error) ([]byte, error) {
	if err == nil {
		p.respons = append(p.respons, isi)
	}
	return isi, err
}

func (p *perekamRespons) isiBernama(awalan string) map[string][]byte {
	berkas := map[string][]byte{}
	for i, isi := range p.respons {
		berkas[fmt.Sprintf("%s_%03d.txt", awalan, i)] = isi
	}
	return berkas
}

func (p *perekamRespons) gabungan() []byte {
	return bytes.Join(p.respons, nil)
}

func keDaftarDomain(daftar []BarisDaftar) []domain.PaketDaftar {
	hasil := make([]domain.PaketDaftar, len(daftar))
	for i, b := range daftar {
		hasil[i] = domain.PaketDaftar{
			Kode: domain.KodePaket(b.IsTender, b.Kode), KodeSPSE: b.Kode, Nama: b.Nama, NamaInstansi: b.Instansi, Tahap: b.Tahap,
			JenisPengadaan: b.JenisPengadaan, Metode: b.Metode, IsTender: b.IsTender,
			TahunAnggaran: tahunTerakhir(b.TahunAnggaran),
			HPS:           nilaiBila(b.HPSRingkas, b.HPSRingkas > 0),
			NilaiKontrak:  nilaiBila(b.NilaiKontrak, b.HasNilaiKontrak),
		}
	}
	return hasil
}

func keDetailDomain(p Pengumuman) domain.DetailPaket {
	return domain.DetailPaket{
		Kode: domain.KodePaket(false, p.Kode), KodeSPSE: p.Kode, Nama: p.Nama, SatuanKerja: p.SatuanKerja, IsBatal: p.IsBatal,
		RUP: rupPertama(p.RUP), TanggalDibuat: p.TanggalDibuat, Tahap: p.Tahap,
		JenisPengadaan: p.JenisPengadaan, Metode: p.Metode, JenisKontrak: p.JenisKontrak,
		TahunAnggaran: p.TahunAnggaran,
		Pagu:          nilaiBila(p.Pagu, p.Pagu > 0),
		HPS:           nilaiBila(p.HPS, p.HPS > 0),
		JumlahPeserta: nilaiBila(p.JumlahPeserta, p.JumlahPeserta > 0),
	}
}

func rupPertama(daftar []domain.RUPTercatat) *domain.RUPTercatat {
	if len(daftar) == 0 {
		return nil
	}
	return &daftar[0]
}

func tahunTerakhir(daftar []int) int {
	if len(daftar) == 0 {
		return 0
	}
	return slices.Max(daftar)
}

func nilaiBila[T any](nilai T, ada bool) *T {
	if !ada {
		return nil
	}
	return &nilai
}
