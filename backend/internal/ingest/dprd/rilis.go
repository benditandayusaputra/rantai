package dprd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"golang.org/x/net/html"
)

const (
	formatTanggalWordPress = "2006-01-02T15:04:05"
	kolomRilis             = "id,date_gmt,link,title,excerpt,content"
)

type teksTerender struct {
	Rendered string `json:"rendered"`
}

type postWordPress struct {
	TanggalGMT string       `json:"date_gmt"`
	Tautan     string       `json:"link"`
	Judul      teksTerender `json:"title"`
	Cuplikan   teksTerender `json:"excerpt"`
	Isi        teksTerender `json:"content"`
}

type SumberRilis struct {
	Konfigurasi KonfigurasiRilis
	Pengunduh   Pengunduh
	DirRaw      string
	Sekarang    func() time.Time
}

func (s SumberRilis) Ambil(ctx context.Context, sejak time.Time) (usecase.DataSumberUmpan, error) {
	diambilPada := s.Sekarang()
	dir := filepath.Join(s.DirRaw, usecase.JobUmpanDPRD, diambilPada.Format(formatTanggalRaw))
	semua, semuaIsi := []domain.ItemUmpan{}, []byte{}
	for halaman := 1; halaman <= s.Konfigurasi.BatasHalaman; halaman++ {
		isi, err := s.Pengunduh.Ambil(ctx, s.urlHalaman(sejak, halaman))
		if err != nil {
			return usecase.DataSumberUmpan{}, err
		}
		if err := simpanMentah(dir, fmt.Sprintf("rilis_%03d.json", halaman), isi); err != nil {
			return usecase.DataSumberUmpan{}, err
		}
		item, err := UraiRilis(isi, s.Konfigurasi)
		if err != nil {
			return usecase.DataSumberUmpan{}, fmt.Errorf("halaman %d: %w", halaman, err)
		}
		semua, semuaIsi = append(semua, item...), append(semuaIsi, isi...)
		if len(item) < s.Konfigurasi.PerHalaman {
			break
		}
	}
	return usecase.DataSumberUmpan{Item: semua, Provenance: s.provenance(semuaIsi, diambilPada)}, nil
}

func (s SumberRilis) urlHalaman(sejak time.Time, halaman int) string {
	kueri := url.Values{
		"categories": {strings.Join(s.Konfigurasi.Kategori, ",")},
		"after":      {sejak.UTC().Format(formatTanggalWordPress)},
		"per_page":   {strconv.Itoa(s.Konfigurasi.PerHalaman)},
		"page":       {strconv.Itoa(halaman)},
		"orderby":    {"date"},
		"order":      {"asc"},
		"_fields":    {kolomRilis},
	}
	return s.Konfigurasi.URLAPI + "?" + kueri.Encode()
}

func UraiRilis(isi []byte, cfg KonfigurasiRilis) ([]domain.ItemUmpan, error) {
	var daftar []postWordPress
	if err := json.Unmarshal(isi, &daftar); err != nil {
		return nil, fmt.Errorf("urai rilis DPRD: %w", err)
	}
	item := make([]domain.ItemUmpan, 0, len(daftar))
	for _, p := range daftar {
		terbit, err := time.Parse(formatTanggalWordPress, p.TanggalGMT)
		if err != nil {
			return nil, fmt.Errorf("tanggal rilis %s: %w", p.Tautan, err)
		}
		item = append(item, domain.ItemUmpan{
			JenisSumber: domain.SumberUmpanResmiDPRD, MediaNama: cfg.MediaNama, MediaDomain: cfg.MediaDomain,
			Judul: teksHTML(p.Judul.Rendered), Cuplikan: domain.PotongCuplikan(teksHTML(p.Cuplikan.Rendered)),
			Isi: teksHTML(p.Isi.Rendered), URL: p.Tautan, TerbitPada: terbit,
		})
	}
	return item, nil
}

func teksHTML(potongan string) string {
	akar, err := html.Parse(strings.NewReader(potongan))
	if err != nil {
		return ""
	}
	return teks(akar)
}

func (s SumberRilis) provenance(isi []byte, diambilPada time.Time) domain.Provenance {
	return provenance.Baru(provenance.Catatan{
		Sumber: domain.SumberDPRD, URL: s.Konfigurasi.URLRujukan, Label: domain.LabelResmi, Catatan: s.Konfigurasi.Catatan,
	}, isi, diambilPada)
}
