package domain

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const (
	SumberKPU  SumberData = "kpu"
	SumberDPRD SumberData = "dprd"
)

var (
	ErrKecamatanDapilTidakDikenal = errors.New("kecamatan dapil tidak dikenal atau namanya tidak sesuai data wilayah")
	ErrKecamatanDapilGanda        = errors.New("kecamatan masuk lebih dari satu dapil")
	ErrKecamatanTanpaDapil        = errors.New("kecamatan tidak masuk dapil mana pun")
	ErrAnggotaTidakCocok          = errors.New("anggota dapil tidak bisa dicocokkan ke fraksi")
	ErrJumlahKursiTidakSesuai     = errors.New("jumlah anggota dapil tidak sama dengan jumlah kursi")
	gelarDepan                    = []string{"h", "hj", "dr", "drs", "ir", "cn", "prof", "mayjen", "tni", "purn", "kh"}
	polaKurung                    = regexp.MustCompile(`\([^)]*\)`)
	polaGelarBelakang             = regexp.MustCompile(`^\p{L}{1,3}(\.\p{L}{1,4})+\.?$`)
	polaBukanHuruf                = regexp.MustCompile(`[^\p{L}]`)
)

type KecamatanDapil struct {
	Kode string `yaml:"kode"`
	Nama string `yaml:"nama"`
}

type Dapil struct {
	Nomor       int              `yaml:"nomor"`
	Nama        string           `yaml:"nama"`
	JumlahKursi int              `yaml:"jumlah_kursi"`
	Kecamatan   []KecamatanDapil `yaml:"kecamatan"`
}

type BarisDapil struct {
	NomorDapil int
	Nama       string
	Partai     string
	Slug       string
}

type BarisStruktur struct {
	Nama    string
	Jabatan string
	Slug    string
}

type AnggotaFraksi struct {
	Fraksi string
	BarisStruktur
}

type Anggota struct {
	ID         string
	Nama       string
	NamaLain   []string
	Partai     string
	Fraksi     string
	NomorDapil int
	URLProfil  string
}

func ValidasiCakupanDapil(daftar []Dapil, kecamatanDKI map[string]string) error {
	pelanggaran := []error{}
	tercakup := map[string]int{}
	for _, d := range daftar {
		for _, k := range d.Kecamatan {
			if kecamatanDKI[k.Kode] != k.Nama {
				pelanggaran = append(pelanggaran, fmt.Errorf("dapil %d %s %q: %w", d.Nomor, k.Kode, k.Nama, ErrKecamatanDapilTidakDikenal))
			}
			tercakup[k.Kode]++
		}
	}
	for kode, nama := range kecamatanDKI {
		pelanggaran = append(pelanggaran, periksaCakupan(kode, nama, tercakup[kode]))
	}
	return errors.Join(pelanggaran...)
}

func periksaCakupan(kode, nama string, jumlah int) error {
	switch {
	case jumlah == 0:
		return fmt.Errorf("%s %s: %w", kode, nama, ErrKecamatanTanpaDapil)
	case jumlah > 1:
		return fmt.Errorf("%s %s: %w", kode, nama, ErrKecamatanDapilGanda)
	default:
		return nil
	}
}

func ValidasiKursiDapil(daftar []Dapil, anggota []Anggota) error {
	jumlahPerDapil := map[int]int{}
	for _, a := range anggota {
		jumlahPerDapil[a.NomorDapil]++
	}
	pelanggaran := []error{}
	for _, d := range daftar {
		if jumlahPerDapil[d.Nomor] != d.JumlahKursi {
			pelanggaran = append(pelanggaran, fmt.Errorf("dapil %d punya %d anggota, kursi %d: %w",
				d.Nomor, jumlahPerDapil[d.Nomor], d.JumlahKursi, ErrJumlahKursiTidakSesuai))
		}
	}
	return errors.Join(pelanggaran...)
}

func NamaTanpaGelar(nama string) string {
	inti, _, _ := strings.Cut(nama, ",")
	kata := strings.Fields(polaKurung.ReplaceAllString(inti, " "))
	for len(kata) > 0 && slices.Contains(gelarDepan, kunciKata(kata[0])) {
		kata = kata[1:]
	}
	for len(kata) > 1 && polaGelarBelakang.MatchString(kata[len(kata)-1]) {
		kata = kata[:len(kata)-1]
	}
	return strings.Join(kata, " ")
}

func KunciNama(nama string) string {
	return kunciKata(NamaTanpaGelar(nama))
}

func kunciKata(kata string) string {
	return strings.ToLower(polaBukanHuruf.ReplaceAllString(kata, ""))
}
