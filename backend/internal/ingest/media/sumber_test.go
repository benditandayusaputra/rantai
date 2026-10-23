package media

import (
	"errors"
	"os"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/dprd"
)

func TestUraiRSS_Antara_JudulTautanDanTanggal(t *testing.T) {
	isi, err := os.ReadFile("testdata/rss_antara_metro.xml")
	if err != nil {
		t.Fatal(err)
	}

	item, err := UraiRSS(isi, Media{Nama: "ANTARA", Domain: "antaranews.com"})
	if err != nil {
		t.Fatal(err)
	}

	if len(item) != 20 {
		t.Fatalf("harap 20 item, dapat %d", len(item))
	}
	for _, i := range item {
		if i.JenisSumber != domain.SumberUmpanMedia || !strings.HasPrefix(i.URL, "https://www.antaranews.com/") ||
			i.Judul == "" || i.TerbitPada.IsZero() || strings.Contains(i.Cuplikan, "<img") || utf8.RuneCountInString(i.Cuplikan) > domain.BatasCuplikanUmpan {
			t.Errorf("item tidak lengkap atau belum bersih: %+v", i)
		}
	}
}

func TestUraiRSS_Detik_TautanDariElemenLinkAtauGUID(t *testing.T) {
	isi, err := os.ReadFile("testdata/rss_detik_berita.xml")
	if err != nil {
		t.Fatal(err)
	}

	item, err := UraiRSS(isi, Media{Nama: "detikcom", Domain: "detik.com"})
	if err != nil || len(item) != 12 {
		t.Fatalf("dapat %d item, %v", len(item), err)
	}
	if !strings.HasPrefix(item[0].URL, "https://news.detik.com/berita/") {
		t.Errorf("tautan = %q", item[0].URL)
	}
}

func TestTerbitSejak_ItemLama_Dibuang(t *testing.T) {
	batas := time.Date(2026, time.October, 10, 0, 0, 0, 0, time.UTC)
	daftar := []domain.ItemUmpan{{URL: "lama", TerbitPada: batas.Add(-time.Hour)}, {URL: "baru", TerbitPada: batas.Add(time.Hour)}}

	if hasil := terbitSejak(daftar, batas); len(hasil) != 1 || hasil[0].URL != "baru" {
		t.Errorf("hasil = %+v", hasil)
	}
}

func TestKonfigurasiMedia_AturanDariYAML_DomainDanAliasFraksiTerisi(t *testing.T) {
	var cfg Konfigurasi
	if err := config.BacaYAML("../../../config/media.yaml", &cfg); err != nil {
		t.Fatal(err)
	}

	aturan := cfg.Aturan()

	if !aturan.IsDomainSah("https://www.antaranews.com/berita/1") || aturan.IsDomainSah("https://contoh.blogspot.com/x") {
		t.Errorf("domain sah = %v", aturan.DomainSah)
	}
	if len(aturan.AliasFraksi) != 9 {
		t.Errorf("harap alias untuk 9 fraksi, dapat %d", len(aturan.AliasFraksi))
	}
}

func TestKonfigurasiMedia_AliasFraksi_SamaDenganFraksiDiKonfigurasiDPRD(t *testing.T) {
	var cfgMedia Konfigurasi
	var cfgDPRD dprd.Konfigurasi
	if err := errors.Join(config.BacaYAML("../../../config/media.yaml", &cfgMedia), config.BacaYAML("../../../config/dprd.yaml", &cfgDPRD)); err != nil {
		t.Fatal(err)
	}
	for _, f := range cfgDPRD.Sumber.Fraksi {
		if len(cfgMedia.AliasFraksi[f.Nama]) == 0 {
			t.Errorf("fraksi %q tidak punya alias di media.yaml", f.Nama)
		}
	}
}
