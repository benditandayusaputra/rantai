package kartu

import (
	"math"
	"strconv"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	bahasaBawaan  = "id"
	desimalRasio  = 1
	desimalUang   = 1
	satuMiliar    = 1e9
	satuJuta      = 1e6
	pemisahKomisi = ", "
)

type Teks struct {
	Tagline        string                  `yaml:"tagline"`
	Judul          string                  `yaml:"judul"`
	TemaKosong     string                  `yaml:"tema_kosong"`
	SatuanRasio    string                  `yaml:"satuan_rasio"`
	Periode        string                  `yaml:"periode"`
	Anggaran       string                  `yaml:"anggaran"`
	AnggaranKosong string                  `yaml:"anggaran_kosong"`
	JumlahPaket    string                  `yaml:"jumlah_paket"`
	Dapil          string                  `yaml:"dapil"`
	Kursi          string                  `yaml:"kursi"`
	Komisi         string                  `yaml:"komisi"`
	SumberTema     string                  `yaml:"sumber_tema"`
	SumberAnggaran string                  `yaml:"sumber_anggaran"`
	SumberWakil    string                  `yaml:"sumber_wakil"`
	Ajakan         string                  `yaml:"ajakan"`
	Miliar         string                  `yaml:"miliar"`
	Juta           string                  `yaml:"juta"`
	PemisahDesimal string                  `yaml:"pemisah_desimal"`
	Label          map[domain.Label]string `yaml:"label"`
}

type Konfigurasi struct {
	AlamatSitus string          `yaml:"-"`
	Teks        map[string]Teks `yaml:"teks"`
}

type blok struct {
	Judul      string
	Nilai      string
	Satuan     string
	Label      domain.Label
	Keterangan string
}

type isiKartu struct {
	Tagline string
	Judul   string
	Kota    string
	Blok    []blok
	Ajakan  string
	Tautan  string
	Label   map[domain.Label]string
}

func (k Konfigurasi) teks(bahasa string) (Teks, string) {
	if t, ada := k.Teks[bahasa]; ada {
		return t, bahasa
	}
	return k.Teks[bahasaBawaan], bahasaBawaan
}

func susunIsi(r domain.RingkasanKecamatan, k Konfigurasi, bahasa string) isiKartu {
	t, bahasa := k.teks(bahasa)
	alamat := strings.TrimPrefix(strings.TrimPrefix(k.AlamatSitus, "https://"), "http://")
	return isiKartu{
		Tagline: t.Tagline,
		Judul:   isi(t.Judul, "nama", r.Kecamatan.Nama),
		Kota:    namaKota(r.Kecamatan.Hierarki),
		Blok:    []blok{blokTema(r.Tema, t, bahasa), blokAnggaran(r.Anggaran, t), blokWakil(r.Wakil, t)},
		Ajakan:  t.Ajakan,
		Tautan:  alamat + "/" + bahasa + "/kecamatan/" + r.Kecamatan.Kode,
		Label:   t.Label,
	}
}

func namaKota(hierarki []domain.RingkasanWilayah) string {
	for _, w := range hierarki {
		if w.Tingkat == domain.TingkatKota {
			return w.Nama
		}
	}
	return ""
}

func blokTema(daftar []domain.TemaKecamatan, t Teks, bahasa string) blok {
	for _, tema := range daftar {
		if tema.Indikator == nil || tema.Indikator.Rasio == nil {
			continue
		}
		return blok{
			Judul:      teksBahasa(tema.Tema.Indikator.Judul, bahasa),
			Nilai:      desimal(*tema.Indikator.Rasio, desimalRasio, t.PemisahDesimal),
			Satuan:     t.SatuanRasio,
			Label:      tema.Indikator.Label,
			Keterangan: t.SumberTema + " · " + isi(t.Periode, "periode", tema.Indikator.Periode),
		}
	}
	judul := ""
	if len(daftar) > 0 {
		judul = teksBahasa(daftar[0].Tema.Indikator.Judul, bahasa)
	}
	return blok{Judul: judul, Nilai: t.TemaKosong}
}

func blokAnggaran(a *domain.RingkasanAnggaran, t Teks) blok {
	if a == nil {
		return blok{Judul: strings.TrimSpace(isi(t.Anggaran, "tahun", "")), Nilai: t.AnggaranKosong}
	}
	return blok{
		Judul:      isi(t.Anggaran, "tahun", strconv.Itoa(a.Tahun)),
		Nilai:      rupiahRingkas(a.TotalHPS, t),
		Label:      a.Label(),
		Keterangan: t.SumberAnggaran + " · " + isi(t.JumlahPaket, "n", strconv.Itoa(a.JumlahPaket)),
	}
}

func blokWakil(w domain.WakilKecamatan, t Teks) blok {
	keterangan := t.SumberWakil
	if len(w.KomisiPengawas) > 0 {
		keterangan += " · " + isi(t.Komisi, "komisi", strings.Join(w.KomisiPengawas, pemisahKomisi))
	}
	return blok{
		Judul:      isi(t.Dapil, "nomor", strconv.Itoa(w.Nomor)),
		Nilai:      isi(t.Kursi, "n", strconv.Itoa(w.JumlahKursi)),
		Label:      domain.LabelResmi,
		Keterangan: keterangan,
	}
}

func teksBahasa(teks domain.Teks, bahasa string) string {
	if bahasa == "en" && teks.EN != "" {
		return teks.EN
	}
	return teks.ID
}

func isi(templat, nama, nilai string) string {
	return strings.ReplaceAll(templat, "{"+nama+"}", nilai)
}

func rupiahRingkas(nilai float64, t Teks) string {
	switch {
	case nilai >= satuMiliar:
		return "Rp" + desimal(nilai/satuMiliar, desimalUang, t.PemisahDesimal) + " " + t.Miliar
	case nilai >= satuJuta:
		return "Rp" + desimal(nilai/satuJuta, desimalUang, t.PemisahDesimal) + " " + t.Juta
	default:
		return domain.FormatRupiah(nilai)
	}
}

func desimal(nilai float64, digit int, pemisah string) string {
	pangkat := math.Pow(10, float64(digit))
	teks := strconv.FormatFloat(math.Round(nilai*pangkat)/pangkat, 'f', -1, 64)
	return strings.Replace(teks, ".", pemisah, 1)
}
