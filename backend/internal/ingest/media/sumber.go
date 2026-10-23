package media

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/klien"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const formatTanggalRaw = "2006-01-02"

var polaTagHTML = regexp.MustCompile(`<[^>]*>`)

type Media struct {
	Nama            string   `yaml:"nama"`
	Domain          string   `yaml:"domain"`
	StatusDewanPers string   `yaml:"status_dewan_pers"`
	RSS             []string `yaml:"rss"`
}

type Konfigurasi struct {
	Pengunduh     klien.Konfigurasi   `yaml:"pengunduh"`
	URLRujukan    string              `yaml:"url_rujukan"`
	Catatan       string              `yaml:"catatan"`
	Mulai         time.Time           `yaml:"mulai"`
	TumpangTindih time.Duration       `yaml:"tumpang_tindih"`
	KonteksDPRD   []string            `yaml:"konteks_dprd"`
	AliasFraksi   map[string][]string `yaml:"alias_fraksi"`
	Media         []Media             `yaml:"media"`
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

type umpanRSS struct {
	Item []itemRSS `xml:"channel>item"`
}

type itemRSS struct {
	Judul     string `xml:"title"`
	Tautan    string `xml:"link"`
	GUID      string `xml:"guid"`
	Terbit    string `xml:"pubDate"`
	Deskripsi string `xml:"description"`
}

func (k Konfigurasi) Aturan() domain.AturanMedia {
	domainSah := make([]string, len(k.Media))
	for i, m := range k.Media {
		domainSah[i] = m.Domain
	}
	return domain.AturanMedia{DomainSah: domainSah, KonteksDPRD: k.KonteksDPRD, AliasFraksi: k.AliasFraksi}
}

func (s Sumber) Ambil(ctx context.Context, sejak time.Time) (usecase.DataSumberUmpan, error) {
	diambilPada := s.Sekarang()
	semua, semuaIsi := []domain.ItemUmpan{}, []byte{}
	for _, m := range s.Konfigurasi.Media {
		for urutan, alamat := range m.RSS {
			isi, err := s.Pengunduh.Ambil(ctx, alamat)
			if err != nil {
				return usecase.DataSumberUmpan{}, err
			}
			if err := s.simpanMentah(diambilPada, fmt.Sprintf("%s_%d.xml", m.Domain, urutan), isi); err != nil {
				return usecase.DataSumberUmpan{}, err
			}
			item, err := UraiRSS(isi, m)
			if err != nil {
				return usecase.DataSumberUmpan{}, fmt.Errorf("%s: %w", alamat, err)
			}
			semua, semuaIsi = append(semua, terbitSejak(item, sejak)...), append(semuaIsi, isi...)
		}
	}
	return usecase.DataSumberUmpan{Item: semua, Provenance: s.provenance(semuaIsi, diambilPada)}, nil
}

func UraiRSS(isi []byte, m Media) ([]domain.ItemUmpan, error) {
	var umpan umpanRSS
	if err := xml.Unmarshal(isi, &umpan); err != nil {
		return nil, fmt.Errorf("urai RSS: %w", err)
	}
	item := make([]domain.ItemUmpan, 0, len(umpan.Item))
	for _, i := range umpan.Item {
		terbit, err := time.Parse(time.RFC1123Z, strings.TrimSpace(i.Terbit))
		if err != nil {
			return nil, fmt.Errorf("tanggal %q: %w", i.Terbit, err)
		}
		item = append(item, domain.ItemUmpan{
			JenisSumber: domain.SumberUmpanMedia, MediaNama: m.Nama, MediaDomain: m.Domain,
			Judul: teksBersih(i.Judul), Cuplikan: domain.PotongCuplikan(teksBersih(i.Deskripsi)),
			URL: tautan(i), TerbitPada: terbit,
		})
	}
	return item, nil
}

func tautan(i itemRSS) string {
	if strings.TrimSpace(i.Tautan) != "" {
		return strings.TrimSpace(i.Tautan)
	}
	return strings.TrimSpace(i.GUID)
}

func teksBersih(teks string) string {
	return strings.Join(strings.Fields(html.UnescapeString(polaTagHTML.ReplaceAllString(teks, " "))), " ")
}

func terbitSejak(daftar []domain.ItemUmpan, sejak time.Time) []domain.ItemUmpan {
	hasil := []domain.ItemUmpan{}
	for _, item := range daftar {
		if item.TerbitPada.After(sejak) {
			hasil = append(hasil, item)
		}
	}
	return hasil
}

func (s Sumber) simpanMentah(diambilPada time.Time, nama string, isi []byte) error {
	dir := filepath.Join(s.DirRaw, usecase.JobUmpanMedia, diambilPada.Format(formatTanggalRaw))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("buat folder data mentah: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, nama), isi, 0o644)
}

func (s Sumber) provenance(isi []byte, diambilPada time.Time) domain.Provenance {
	return provenance.Baru(provenance.Catatan{
		Sumber: domain.SumberMedia, URL: s.Konfigurasi.URLRujukan, Label: domain.LabelResmi, Catatan: s.Konfigurasi.Catatan,
	}, isi, diambilPada)
}
