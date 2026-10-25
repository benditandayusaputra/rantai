package wikimedia

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var ErrEntriTidakLengkap = errors.New("entri foto wikimedia wajib punya npsn dan berkas")

type Entri struct {
	NPSN     string `yaml:"npsn"`
	Wikidata string `yaml:"wikidata"`
	Berkas   string `yaml:"berkas"`
	Pembuat  string `yaml:"pembuat"`
	Lisensi  string `yaml:"lisensi"`
}

type Konfigurasi struct {
	Sumber         string  `yaml:"sumber"`
	Dihimpun       string  `yaml:"dihimpun"`
	PolaHalaman    string  `yaml:"pola_halaman"`
	PolaThumbnail  string  `yaml:"pola_thumbnail"`
	LebarThumbnail int     `yaml:"lebar_thumbnail"`
	Foto           []Entri `yaml:"foto"`
}

type Sumber struct {
	Konfigurasi Konfigurasi
}

func (s Sumber) Ambil(_ context.Context) ([]domain.FotoSekolah, error) {
	daftar := make([]domain.FotoSekolah, 0, len(s.Konfigurasi.Foto))
	for _, e := range s.Konfigurasi.Foto {
		if e.NPSN == "" || e.Berkas == "" {
			return nil, fmt.Errorf("%+v: %w", e, ErrEntriTidakLengkap)
		}
		daftar = append(daftar, s.keFoto(e))
	}
	return daftar, nil
}

func (s Sumber) keFoto(e Entri) domain.FotoSekolah {
	nama := NamaBerkas(e.Berkas)
	return domain.FotoSekolah{NPSN: e.NPSN, Foto: domain.Foto{
		URL:     AlamatThumbnail(s.Konfigurasi.PolaThumbnail, nama, s.Konfigurasi.LebarThumbnail),
		Sumber:  domain.SumberFotoWikimedia,
		Halaman: fmt.Sprintf(s.Konfigurasi.PolaHalaman, url.PathEscape(nama)),
		Kredit:  domain.SusunKredit(e.Pembuat, e.Lisensi),
	}}
}

func NamaBerkas(berkas string) string {
	return strings.ReplaceAll(strings.TrimSpace(berkas), " ", "_")
}

func AlamatThumbnail(pola, nama string, lebar int) string {
	jejak := md5.Sum([]byte(nama))
	hex := hex.EncodeToString(jejak[:])
	tersandi := url.PathEscape(nama)
	return fmt.Sprintf(pola, hex[:1], hex[:2], tersandi, lebar, tersandi)
}
