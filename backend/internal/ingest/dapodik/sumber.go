package dapodik

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const (
	formatTanggalRaw = "2006-01-02"
	pathCari         = "/cari-sekolah"
	pathDetail       = "/full-detail/"
	pathRombel       = "/rombongan-belajar/"
	pathPesertaDidik = "/peserta-didik/"
)

type Konfigurasi struct {
	URLAPI             string                 `yaml:"url_api"`
	URLRujukan         string                 `yaml:"url_rujukan"`
	PolaHalamanProfil  string                 `yaml:"pola_halaman_profil"`
	Catatan            string                 `yaml:"catatan"`
	Pengunduh          klien.Konfigurasi      `yaml:"pengunduh"`
	UkuranHalaman      int                    `yaml:"ukuran_halaman"`
	Kota               []string               `yaml:"kota"`
	Bentuk             []domain.BentukSekolah `yaml:"bentuk"`
	BentukDenganRombel []domain.BentukSekolah `yaml:"bentuk_dengan_rombel"`
}

type Sumber struct {
	Konfigurasi Konfigurasi
	Pengunduh   Pengunduh
	DirRaw      string
	Sekarang    func() time.Time
}

type pengambilan struct {
	cfg   Konfigurasi
	unduh pengunduhBercache
	hash  hash.Hash
}

type kueriCari struct {
	Page             int    `json:"page"`
	Size             int    `json:"size"`
	Keyword          string `json:"keyword"`
	KabupatenKota    string `json:"kabupaten_kota"`
	BentukPendidikan string `json:"bentuk_pendidikan"`
	StatusSekolah    string `json:"status_sekolah"`
}

func (s Sumber) Ambil(ctx context.Context) (usecase.DataSumberSekolah, error) {
	diambilPada := s.Sekarang()
	dirCache := filepath.Join(s.DirRaw, usecase.JobDapodik, diambilPada.Format(formatTanggalRaw))
	if err := os.MkdirAll(dirCache, 0o755); err != nil {
		return usecase.DataSumberSekolah{}, fmt.Errorf("buat folder data mentah: %w", err)
	}
	p := pengambilan{
		cfg:   s.Konfigurasi,
		unduh: pengunduhBercache{pengunduh: s.Pengunduh, dirCache: dirCache},
		hash:  sha256.New(),
	}
	daftar, err := p.ambilSemua(ctx)
	if err != nil {
		return usecase.DataSumberSekolah{}, err
	}
	return usecase.DataSumberSekolah{Sekolah: daftar, Provenance: s.catatProvenance(p.hash.Sum(nil), diambilPada)}, nil
}

func (p pengambilan) ambilSemua(ctx context.Context) ([]domain.DataSekolah, error) {
	daftar := []domain.DataSekolah{}
	for _, kota := range p.cfg.Kota {
		for _, bentuk := range p.cfg.Bentuk {
			butir, err := p.daftarSekolah(ctx, kota, bentuk)
			if err != nil {
				return nil, err
			}
			for _, b := range butir {
				data, err := p.ambilSekolah(ctx, b)
				if err != nil {
					return nil, fmt.Errorf("ambil sekolah %s: %w", b.NPSN, err)
				}
				daftar = append(daftar, data)
			}
		}
	}
	return daftar, nil
}

func (p pengambilan) daftarSekolah(ctx context.Context, kota string, bentuk domain.BentukSekolah) ([]ButirCari, error) {
	semua := []ButirCari{}
	for halaman := 0; ; halaman++ {
		isi, err := p.cari(ctx, kueriCari{Page: halaman, Size: p.cfg.UkuranHalaman, KabupatenKota: kota, BentukPendidikan: string(bentuk)})
		if err != nil {
			return nil, err
		}
		hasil, err := UraiCari(isi)
		if err != nil {
			return nil, fmt.Errorf("daftar %s %s halaman %d: %w", bentuk, kota, halaman, err)
		}
		semua = append(semua, hasil.Sekolah...)
		if len(hasil.Sekolah) == 0 || len(semua) >= hasil.Total {
			return semua, nil
		}
	}
}

func (p pengambilan) cari(ctx context.Context, kueri kueriCari) ([]byte, error) {
	badan, err := json.Marshal(kueri)
	if err != nil {
		return nil, err
	}
	namaCache := fmt.Sprintf("cari_%s_%s_%d.json", slugKota(kueri.KabupatenKota), kueri.BentukPendidikan, kueri.Page)
	return p.ambil(ctx, permintaan{namaCache: namaCache, url: p.cfg.URLAPI + pathCari, badanJSON: badan})
}

func (p pengambilan) ambilSekolah(ctx context.Context, b ButirCari) (domain.DataSekolah, error) {
	isi, err := p.ambilPerSekolah(ctx, "detail", pathDetail, b.SekolahID)
	if err != nil {
		return domain.DataSekolah{}, err
	}
	detail, err := UraiDetail(isi)
	if err != nil {
		return domain.DataSekolah{}, err
	}
	tahunAjaran, err := domain.TahunAjaranDariSemester(detail.SemesterID)
	if err != nil {
		return domain.DataSekolah{}, err
	}
	data := domain.DataSekolah{Sekolah: detail.Sekolah, TahunAjaran: tahunAjaran, Rombel: []domain.RombelTingkat{}}
	data.Foto = p.fotoSekolah(detail)
	if !slices.Contains(p.cfg.BentukDenganRombel, detail.Sekolah.Bentuk) {
		return data, nil
	}
	data.Rombel, err = p.ambilRombel(ctx, b.SekolahID, detail.Sekolah.Bentuk)
	return data, err
}

func (p pengambilan) fotoSekolah(detail DetailSekolah) *domain.Foto {
	if detail.FotoURL == "" {
		return nil
	}
	return &domain.Foto{
		URL:     detail.FotoURL,
		Sumber:  domain.SumberFotoSekolahKita,
		Halaman: fmt.Sprintf(p.cfg.PolaHalamanProfil, detail.SekolahID),
	}
}

func (p pengambilan) ambilRombel(ctx context.Context, sekolahID string, bentuk domain.BentukSekolah) ([]domain.RombelTingkat, error) {
	isiRombel, err := p.ambilPerSekolah(ctx, "rombel", pathRombel, sekolahID)
	if err != nil {
		return nil, err
	}
	isiSiswa, err := p.ambilPerSekolah(ctx, "pd", pathPesertaDidik, sekolahID)
	if err != nil {
		return nil, err
	}
	rombel, err := UraiRombel(isiRombel, bentuk)
	if err != nil {
		return nil, err
	}
	siswa, err := UraiPesertaDidik(isiSiswa, bentuk)
	if err != nil {
		return nil, err
	}
	return gabungkanTingkat(rombel, siswa, bentuk), nil
}

func gabungkanTingkat(rombel, siswa map[int]int, bentuk domain.BentukSekolah) []domain.RombelTingkat {
	hasil := []domain.RombelTingkat{}
	if len(rombel) == 0 && len(siswa) == 0 {
		return hasil
	}
	pertama, terakhir := bentuk.RentangTingkat()
	for tingkat := pertama; tingkat <= terakhir; tingkat++ {
		hasil = append(hasil, domain.RombelTingkat{Tingkat: tingkat, JumlahRombel: rombel[tingkat], JumlahSiswa: siswa[tingkat]})
	}
	return hasil
}

func (p pengambilan) ambilPerSekolah(ctx context.Context, jenis, path, sekolahID string) ([]byte, error) {
	namaCache := fmt.Sprintf("%s_%s.json", jenis, sekolahID)
	return p.ambil(ctx, permintaan{namaCache: namaCache, url: p.cfg.URLAPI + path + sekolahID})
}

func (p pengambilan) ambil(ctx context.Context, req permintaan) ([]byte, error) {
	isi, err := p.unduh.ambil(ctx, req)
	if err != nil {
		return nil, err
	}
	_, _ = p.hash.Write(isi)
	return isi, nil
}

func (s Sumber) catatProvenance(ringkasan []byte, diambilPada time.Time) domain.Provenance {
	return provenance.Baru(provenance.Catatan{
		Sumber:  domain.SumberDapodik,
		URL:     s.Konfigurasi.URLRujukan,
		Label:   domain.LabelResmi,
		Catatan: s.Konfigurasi.Catatan,
	}, ringkasan, diambilPada)
}

func slugKota(kota string) string {
	return strings.NewReplacer(" ", "_", ".", "").Replace(strings.ToLower(kota))
}
