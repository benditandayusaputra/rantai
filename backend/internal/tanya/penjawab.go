package tanya

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	BahasaIndonesia = "id"
	BahasaInggris   = "en"
)

var ErrBatasLangkah = errors.New("model terus memanggil alat melewati batas langkah")

type Permintaan struct {
	Pertanyaan           string
	KodeKecamatanKonteks string
	Bahasa               string
}

type PanggilanAlat struct {
	ID      string
	Nama    string
	Masukan json.RawMessage
}

type HasilAlat struct {
	IDPanggilan string
	Isi         string
	IsError     bool
}

type Giliran struct {
	Teks      string
	Panggilan []PanggilanAlat
}

type Percakapan interface {
	Tanya(ctx context.Context, isi string) (Giliran, error)
	Balas(ctx context.Context, hasil []HasilAlat) (Giliran, error)
}

type Model interface {
	BukaPercakapan(sistem string, alat []Spek) Percakapan
}

type Kabar struct {
	Alat string
}

type Rujukan struct {
	Nomor        int    `json:"nomor"`
	ProvenanceID string `json:"provenance_id"`
}

type Jawaban struct {
	Teks      string
	Rujukan   []Rujukan
	HasilAlat []json.RawMessage
}

type Aturan struct {
	PromptSistem          string      `yaml:"prompt_sistem"`
	BatasLangkah          int         `yaml:"batas_langkah"`
	PanjangMaksPertanyaan int         `yaml:"panjang_maks_pertanyaan"`
	PesanAman             domain.Teks `yaml:"pesan_aman"`
}

type Penjawab struct {
	Model     Model
	Alat      Registry
	Aturan    Aturan
	Validator ValidatorJawaban
	Laporkan  func(pesan string)
}

func (p Penjawab) Periksa(m Permintaan) error {
	panjang := utf8.RuneCountInString(strings.TrimSpace(m.Pertanyaan))
	switch {
	case panjang == 0 || panjang > p.Aturan.PanjangMaksPertanyaan:
		return domain.InputTidakValid(fmt.Sprintf("pertanyaan wajib diisi, paling banyak %d karakter", p.Aturan.PanjangMaksPertanyaan))
	case m.Bahasa != "" && m.Bahasa != BahasaIndonesia && m.Bahasa != BahasaInggris:
		return domain.InputTidakValid("bahasa harus id atau en")
	case m.KodeKecamatanKonteks != "":
		return periksaKodeKonteks(m.KodeKecamatanKonteks)
	default:
		return nil
	}
}

func periksaKodeKonteks(kode string) error {
	tingkat, err := domain.TingkatDariKode(kode)
	if err != nil {
		return err
	}
	if tingkat != domain.TingkatKecamatan && tingkat != domain.TingkatKelurahan {
		return domain.InputTidakValid("kode_kecamatan_konteks harus kode kecamatan atau kelurahan")
	}
	return nil
}

func (p Penjawab) Jawab(ctx context.Context, m Permintaan, kabar func(Kabar)) (Jawaban, error) {
	if err := p.Periksa(m); err != nil {
		return Jawaban{}, err
	}
	percakapan := p.Model.BukaPercakapan(p.Aturan.PromptSistem, p.Alat.Daftar())
	giliran, err := percakapan.Tanya(ctx, isiPertanyaan(m))
	if err != nil {
		return Jawaban{}, fmt.Errorf("kirim pertanyaan: %w", err)
	}
	semuaHasil := []json.RawMessage{}
	for langkah := 0; len(giliran.Panggilan) > 0; langkah++ {
		if langkah >= p.Aturan.BatasLangkah {
			return Jawaban{}, ErrBatasLangkah
		}
		hasil, isi, err := p.jalankanPanggilan(ctx, giliran.Panggilan, kabar)
		if err != nil {
			return Jawaban{}, err
		}
		semuaHasil = append(semuaHasil, isi...)
		if giliran, err = percakapan.Balas(ctx, hasil); err != nil {
			return Jawaban{}, fmt.Errorf("kirim hasil alat: %w", err)
		}
	}
	return p.saring(m, giliran.Teks, semuaHasil), nil
}

func (p Penjawab) saring(m Permintaan, teks string, hasilAlat []json.RawMessage) Jawaban {
	if err := p.Validator.Periksa(teks, hasilAlat); err != nil {
		p.Laporkan(fmt.Sprintf("jawaban tanya diganti pesan aman: %v; pertanyaan %q; jawaban %q", err, m.Pertanyaan, teks))
		return Jawaban{Teks: p.pesanAman(m.Bahasa), Rujukan: []Rujukan{}, HasilAlat: hasilAlat}
	}
	teksBernomor, rujukan := SusunRujukan(teks)
	return Jawaban{Teks: teksBernomor, Rujukan: rujukan, HasilAlat: hasilAlat}
}

func (p Penjawab) pesanAman(bahasa string) string {
	if bahasa == BahasaInggris {
		return p.Aturan.PesanAman.EN
	}
	return p.Aturan.PesanAman.ID
}

func isiPertanyaan(m Permintaan) string {
	bahasa := m.Bahasa
	if bahasa == "" {
		bahasa = BahasaIndonesia
	}
	return fmt.Sprintf("<bahasa_antarmuka>%s</bahasa_antarmuka>\n<kecamatan_aktif>%s</kecamatan_aktif>\n<pertanyaan>\n%s\n</pertanyaan>",
		bahasa, m.KodeKecamatanKonteks, strings.TrimSpace(m.Pertanyaan))
}

func (p Penjawab) jalankanPanggilan(ctx context.Context, daftar []PanggilanAlat, kabar func(Kabar)) ([]HasilAlat, []json.RawMessage, error) {
	hasil := make([]HasilAlat, len(daftar))
	berhasil := []json.RawMessage{}
	for i, panggilan := range daftar {
		kabar(Kabar{Alat: panggilan.Nama})
		isi, err := p.Alat.Jalankan(ctx, panggilan.Nama, panggilan.Masukan)
		pesan, isGagalAlat := pesanGagalAlat(err)
		switch {
		case isGagalAlat:
			hasil[i] = HasilAlat{IDPanggilan: panggilan.ID, Isi: pesan, IsError: true}
		case err != nil:
			return nil, nil, err
		default:
			hasil[i] = HasilAlat{IDPanggilan: panggilan.ID, Isi: string(isi)}
			berhasil = append(berhasil, isi)
		}
	}
	return hasil, berhasil, nil
}

func pesanGagalAlat(err error) (string, bool) {
	var kesalahan *domain.Kesalahan
	switch {
	case errors.As(err, &kesalahan):
		return kesalahan.Pesan, true
	case errors.Is(err, ErrAlatTidakDikenal):
		return err.Error(), true
	default:
		return "", false
	}
}
