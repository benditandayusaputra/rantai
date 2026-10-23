package lpse

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const pemisahJenisTahun = " - TA "

var (
	polaToken              = regexp.MustCompile(`authenticityToken = '([0-9a-f]+)'`)
	ErrTokenTidakDitemukan = errors.New("authenticityToken tidak ditemukan di halaman daftar")
)

type Jenis struct {
	Nama         string
	PathHalaman  string
	PathData     string
	IsTender     bool
	KolomMetode  int
	KolomJenis   int
	KolomKontrak int
}

var (
	JenisTender = Jenis{
		Nama: "tender", PathHalaman: "/lelang", PathData: "/dt/lelang", IsTender: true,
		KolomMetode: 6, KolomJenis: 8, KolomKontrak: 10,
	}
	JenisNonTender = Jenis{
		Nama: "nontender", PathHalaman: "/nontender", PathData: "/dt/pl",
		KolomMetode: 5, KolomJenis: 6, KolomKontrak: 8,
	}
)

const (
	kolomKode = iota
	kolomNama
	kolomInstansi
	kolomTahap
	kolomHPSRingkas
)

type BarisDaftar struct {
	Kode            string
	Nama            string
	Instansi        string
	Tahap           string
	Metode          string
	JenisPengadaan  string
	IsTender        bool
	HPSRingkas      float64
	TahunAnggaran   []int
	NilaiKontrak    float64
	HasNilaiKontrak bool
}

type Pengirim interface {
	Ambil(ctx context.Context, alamat string) ([]byte, error)
	KirimForm(ctx context.Context, alamat string, form url.Values) ([]byte, error)
}

type SumberDaftar struct {
	Klien         Pengirim
	URLDasar      string
	UkuranHalaman int
}

type responsDaftar struct {
	Data [][]*string `json:"data"`
}

func (s SumberDaftar) AmbilDaftar(ctx context.Context, jenis Jenis, tahun int) ([]BarisDaftar, error) {
	token, err := s.ambilToken(ctx, jenis)
	if err != nil {
		return nil, err
	}
	semua := []BarisDaftar{}
	sudahAda := map[string]bool{}
	for mulai := 0; ; mulai += s.UkuranHalaman {
		halaman, err := s.ambilHalaman(ctx, jenis, tahun, token, mulai)
		if err != nil {
			return nil, err
		}
		baru := saringBaru(halaman, sudahAda)
		semua = append(semua, baru...)
		if len(halaman) < s.UkuranHalaman || len(baru) == 0 {
			return semua, nil
		}
	}
}

func (s SumberDaftar) ambilToken(ctx context.Context, jenis Jenis) (string, error) {
	isi, err := s.Klien.Ambil(ctx, s.URLDasar+jenis.PathHalaman)
	if err != nil {
		return "", fmt.Errorf("ambil halaman daftar %s: %w", jenis.Nama, err)
	}
	return UraiToken(isi)
}

func (s SumberDaftar) ambilHalaman(ctx context.Context, jenis Jenis, tahun int, token string, mulai int) ([]BarisDaftar, error) {
	alamat := s.URLDasar + jenis.PathData + "?tahun=" + strconv.Itoa(tahun)
	form := url.Values{
		"draw":              {"1"},
		"start":             {strconv.Itoa(mulai)},
		"length":            {strconv.Itoa(s.UkuranHalaman)},
		"authenticityToken": {token},
	}
	isi, err := s.Klien.KirimForm(ctx, alamat, form)
	if err != nil {
		return nil, fmt.Errorf("ambil daftar %s %d mulai %d: %w", jenis.Nama, tahun, mulai, err)
	}
	return UraiDaftar(isi, jenis)
}

func saringBaru(halaman []BarisDaftar, sudahAda map[string]bool) []BarisDaftar {
	baru := []BarisDaftar{}
	for _, baris := range halaman {
		if !sudahAda[baris.Kode] {
			sudahAda[baris.Kode] = true
			baru = append(baru, baris)
		}
	}
	return baru
}

func UraiToken(isi []byte) (string, error) {
	cocok := polaToken.FindSubmatch(isi)
	if cocok == nil {
		return "", ErrTokenTidakDitemukan
	}
	return string(cocok[1]), nil
}

func UraiDaftar(isi []byte, jenis Jenis) ([]BarisDaftar, error) {
	var respons responsDaftar
	if err := json.Unmarshal(isi, &respons); err != nil {
		return nil, fmt.Errorf("urai json daftar %s: %w", jenis.Nama, err)
	}
	daftar := make([]BarisDaftar, 0, len(respons.Data))
	for _, baris := range respons.Data {
		if len(baris) <= jenis.KolomKontrak {
			return nil, fmt.Errorf("baris daftar %s hanya %d kolom", jenis.Nama, len(baris))
		}
		daftar = append(daftar, uraiBaris(kolomTeks(baris), jenis))
	}
	return daftar, nil
}

func uraiBaris(kolom []string, jenis Jenis) BarisDaftar {
	jenisPengadaan, tahun, _ := strings.Cut(kolom[jenis.KolomJenis], pemisahJenisTahun)
	hps, _ := uraiNilaiRingkas(kolom[kolomHPSRingkas])
	kontrak, hasKontrak := uraiRupiah(kolom[jenis.KolomKontrak])
	return BarisDaftar{
		Kode:            kolom[kolomKode],
		Nama:            namaTanpaLencana(kolom[kolomNama]),
		Instansi:        kolom[kolomInstansi],
		Tahap:           kolom[kolomTahap],
		Metode:          kolom[jenis.KolomMetode],
		JenisPengadaan:  jenisPengadaan,
		IsTender:        jenis.IsTender,
		HPSRingkas:      hps,
		TahunAnggaran:   uraiDaftarTahun(tahun),
		NilaiKontrak:    kontrak,
		HasNilaiKontrak: hasKontrak,
	}
}

func kolomTeks(baris []*string) []string {
	kolom := make([]string, len(baris))
	for i, sel := range baris {
		if sel != nil {
			kolom[i] = strings.TrimSpace(*sel)
		}
	}
	return kolom
}

func namaTanpaLencana(nama string) string {
	dokumen, err := html.Parse(strings.NewReader(nama))
	if err != nil {
		return rapikan(nama)
	}
	return teksTanpaLencana(dokumen)
}
