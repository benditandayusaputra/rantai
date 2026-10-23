package domain

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
)

const PenghitungKesenjanganKursiSMP = "kesenjangan_kursi_smp"

var (
	ErrTemaTanpaKode                = errors.New("kode tema wajib diisi")
	ErrTeksTidakLengkap             = errors.New("teks wajib tersedia dalam bahasa Indonesia dan Inggris")
	ErrPenghitungTidakTerdaftar     = errors.New("penghitung indikator tidak terdaftar")
	ErrAmbangRasioTidakBerurutan    = errors.New("ambang_rasio harus positif dan naik berurutan")
	ErrFilterAnggaranKosong         = errors.New("filter_anggaran wajib punya minimal satu pola instansi atau kata kunci")
	ErrRegexFilterRusak             = errors.New("pola filter_anggaran bukan regex yang sah")
	ErrTanggungJawabTanpaDasarHukum = errors.New("tanggung jawab wajib punya dasar_hukum dan dasar_hukum_url yang sah")
	ErrKomisiTidakSah               = errors.New("komisi wajib salah satu dari A sampai E")
	ErrTagIsuTidakSah               = errors.New("tag_isu wajib minimal satu dan berhuruf kecil")
	ErrPlaceholderTidakDikenal      = errors.New("placeholder tidak dikenal")
	komisiSah                       = []string{"A", "B", "C", "D", "E"}
	placeholderDikenal              = []string{"kecamatan", "kota", "tahun_depan"}
	polaPlaceholder                 = regexp.MustCompile(`\{([^{}]*)\}`)
	skemaURLDasarHukum              = []string{"http", "https"}
)

type Teks struct {
	ID string `yaml:"id" json:"id"`
	EN string `yaml:"en" json:"en"`
}

type AmbangRasio struct {
	SangatKurang float64 `yaml:"sangat_kurang" json:"sangat_kurang"`
	Kurang       float64 `yaml:"kurang" json:"kurang"`
	HampirCukup  float64 `yaml:"hampir_cukup" json:"hampir_cukup"`
}

type Indikator struct {
	Penghitung      string      `yaml:"penghitung" json:"penghitung"`
	SatuanKebutuhan string      `yaml:"satuan_kebutuhan" json:"satuan_kebutuhan"`
	SatuanKapasitas string      `yaml:"satuan_kapasitas" json:"satuan_kapasitas"`
	Judul           Teks        `yaml:"judul" json:"judul"`
	AmbangRasio     AmbangRasio `yaml:"ambang_rasio" json:"ambang_rasio"`
}

type FilterAnggaran struct {
	InstansiPola []string `yaml:"instansi_pola" json:"instansi_pola"`
	KataKunci    []string `yaml:"kata_kunci" json:"kata_kunci"`
}

type TanggungJawab struct {
	Masalah           string `yaml:"masalah" json:"masalah"`
	Label             Teks   `yaml:"label" json:"label"`
	InstansiPelaksana string `yaml:"instansi_pelaksana" json:"instansi_pelaksana"`
	UnitWilayah       string `yaml:"unit_wilayah" json:"unit_wilayah"`
	Komisi            string `yaml:"komisi" json:"komisi"`
	DasarHukum        string `yaml:"dasar_hukum" json:"dasar_hukum"`
	DasarHukumURL     string `yaml:"dasar_hukum_url" json:"dasar_hukum_url"`
}

type TemplatAspirasi struct {
	Judul      Teks   `yaml:"judul" json:"judul"`
	Pertanyaan []Teks `yaml:"pertanyaan" json:"pertanyaan"`
}

type KonfigurasiTema struct {
	Indikator       Indikator       `yaml:"indikator" json:"indikator"`
	FilterAnggaran  FilterAnggaran  `yaml:"filter_anggaran" json:"filter_anggaran"`
	TanggungJawab   []TanggungJawab `yaml:"tanggung_jawab" json:"tanggung_jawab"`
	TagIsu          []string        `yaml:"tag_isu" json:"tag_isu"`
	TemplatAspirasi TemplatAspirasi `yaml:"templat_aspirasi" json:"templat_aspirasi"`
}

type Tema struct {
	Kode            string `yaml:"kode" json:"kode"`
	Nama            Teks   `yaml:"nama" json:"nama"`
	Aktif           bool   `yaml:"aktif" json:"aktif"`
	Urutan          int    `yaml:"urutan" json:"urutan"`
	KonfigurasiTema `yaml:",inline"`
}

type TemaTersimpan struct {
	Tema
	ProvenanceID string
}

func (t Tema) Validasi(penghitungTerdaftar []string) error {
	pelanggaran := []error{t.validasiIdentitas()}
	pelanggaran = append(pelanggaran, t.Indikator.validasi(penghitungTerdaftar)...)
	pelanggaran = append(pelanggaran, t.FilterAnggaran.validasi()...)
	for _, tj := range t.TanggungJawab {
		pelanggaran = append(pelanggaran, tj.validasi()...)
	}
	pelanggaran = append(pelanggaran, validasiTagIsu(t.TagIsu), t.TemplatAspirasi.validasi())
	if err := errors.Join(pelanggaran...); err != nil {
		return fmt.Errorf("tema %q tidak valid: %w", t.Kode, err)
	}
	return nil
}

func (t Tema) validasiIdentitas() error {
	if t.Kode == "" {
		return ErrTemaTanpaKode
	}
	return t.Nama.validasi("nama")
}

func (t Teks) validasi(bagian string) error {
	if strings.TrimSpace(t.ID) == "" || strings.TrimSpace(t.EN) == "" {
		return fmt.Errorf("%s: %w", bagian, ErrTeksTidakLengkap)
	}
	return nil
}

func (i Indikator) validasi(penghitungTerdaftar []string) []error {
	pelanggaran := []error{i.Judul.validasi("indikator.judul")}
	if !slices.Contains(penghitungTerdaftar, i.Penghitung) {
		pelanggaran = append(pelanggaran, fmt.Errorf("%q: %w", i.Penghitung, ErrPenghitungTidakTerdaftar))
	}
	if !i.AmbangRasio.isBerurutan() {
		pelanggaran = append(pelanggaran, ErrAmbangRasioTidakBerurutan)
	}
	return pelanggaran
}

func (a AmbangRasio) isBerurutan() bool {
	return a.SangatKurang > 0 && a.SangatKurang < a.Kurang && a.Kurang < a.HampirCukup
}

func (f FilterAnggaran) validasi() []error {
	if len(f.InstansiPola) == 0 && len(f.KataKunci) == 0 {
		return []error{ErrFilterAnggaranKosong}
	}
	pelanggaran := []error{}
	for _, pola := range slices.Concat(f.InstansiPola, f.KataKunci) {
		if _, err := regexp.Compile(pola); err != nil {
			pelanggaran = append(pelanggaran, fmt.Errorf("%q: %w", pola, ErrRegexFilterRusak))
		}
	}
	return pelanggaran
}

func (tj TanggungJawab) validasi() []error {
	pelanggaran := []error{tj.Label.validasi("tanggung_jawab." + tj.Masalah + ".label")}
	if strings.TrimSpace(tj.DasarHukum) == "" || !isURLSah(tj.DasarHukumURL) {
		pelanggaran = append(pelanggaran, fmt.Errorf("%q: %w", tj.Masalah, ErrTanggungJawabTanpaDasarHukum))
	}
	if !slices.Contains(komisiSah, tj.Komisi) {
		pelanggaran = append(pelanggaran, fmt.Errorf("%q komisi %q: %w", tj.Masalah, tj.Komisi, ErrKomisiTidakSah))
	}
	return append(pelanggaran, validasiPlaceholder(tj.UnitWilayah))
}

func isURLSah(alamat string) bool {
	u, err := url.ParseRequestURI(alamat)
	return err == nil && u.Host != "" && slices.Contains(skemaURLDasarHukum, u.Scheme)
}

func validasiTagIsu(tagIsu []string) error {
	if len(tagIsu) == 0 {
		return ErrTagIsuTidakSah
	}
	for _, tag := range tagIsu {
		if tag == "" || tag != strings.ToLower(tag) {
			return fmt.Errorf("%q: %w", tag, ErrTagIsuTidakSah)
		}
	}
	return nil
}

func (t TemplatAspirasi) validasi() error {
	pelanggaran := []error{}
	for _, teks := range append([]Teks{t.Judul}, t.Pertanyaan...) {
		pelanggaran = append(pelanggaran, teks.validasi("templat_aspirasi"),
			validasiPlaceholder(teks.ID), validasiPlaceholder(teks.EN))
	}
	return errors.Join(pelanggaran...)
}

func validasiPlaceholder(teks string) error {
	for _, cocok := range polaPlaceholder.FindAllStringSubmatch(teks, -1) {
		if !slices.Contains(placeholderDikenal, cocok[1]) {
			return fmt.Errorf("{%s}: %w", cocok[1], ErrPlaceholderTidakDikenal)
		}
	}
	return nil
}
