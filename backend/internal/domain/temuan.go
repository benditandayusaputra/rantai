package domain

import (
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

type SikapTemuan string

const (
	SikapMendukung SikapTemuan = "mendukung"
	SikapMembantah SikapTemuan = "membantah"
)

type StatusTemuan string

const (
	StatusTemuanMenunggu      StatusTemuan = "menunggu_konfirmasi"
	StatusTemuanTerkonfirmasi StatusTemuan = "terkonfirmasi"
	StatusTemuanDiperdebatkan StatusTemuan = "diperdebatkan"
)

const (
	PanjangMinIsiTemuan    = 20
	PanjangMaksIsiTemuan   = 1000
	PanjangMaksTautan      = 500
	AmbangKonfirmasiTemuan = 2
	BatasTemuanPerHari     = 5
	JendelaBatasTemuan     = 24 * time.Hour
	BatasTemuanPerPaket    = 50
	TargetTemuan           = "temuan"
	AksiSembunyikanTemuan  = AksiModerasi("sembunyikan_temuan")
)

var (
	ErrIsiTemuanTidakSah     = InputTidakValid("isi temuan 20 sampai 1000 karakter")
	ErrTautanTemuanTidakSah  = InputTidakValid("tautan sumber harus alamat web lengkap berawalan https:// atau http://")
	ErrKosakataTemuan        = InputTidakValid("tulis temuan sebagai fakta dan sumbernya, tanpa kata tuduhan")
	ErrSikapTemuanTidakSah   = InputTidakValid("sikap harus mendukung atau membantah")
	ErrTanggapiTemuanSendiri = InputTidakValid("temuan sendiri tidak dapat dikonfirmasi oleh pengirimnya")
	ErrBatasTemuanHarian     = TerlaluBanyak("batas 5 temuan per hari sudah tercapai; coba lagi besok")
	ErrTemuanTidakDitemukan  = TidakDitemukan("temuan tidak ditemukan")
)

type TemuanMasuk struct {
	PenggunaID   int64
	PaketKode    string
	Isi          string
	TautanSumber string
}

type TanggapanMasuk struct {
	TemuanID     int64
	PenggunaID   int64
	Sikap        SikapTemuan
	TautanSumber string
}

type Temuan struct {
	ID              int64
	PenggunaID      int64
	PaketKode       string
	Isi             string
	TautanSumber    string
	NamaKontributor string
	DibuatPada      time.Time
	JumlahMendukung int
	JumlahMembantah int
}

func (t TemuanMasuk) Validasi(terlarang KosakataTerlarang) error {
	panjang := utf8.RuneCountInString(strings.TrimSpace(t.Isi))
	switch {
	case panjang < PanjangMinIsiTemuan || panjang > PanjangMaksIsiTemuan:
		return ErrIsiTemuanTidakSah
	case len(terlarang.Temukan(t.Isi)) > 0:
		return ErrKosakataTemuan
	default:
		return ValidasiTautanSumber(t.TautanSumber)
	}
}

func (t TanggapanMasuk) Validasi() error {
	if t.Sikap != SikapMendukung && t.Sikap != SikapMembantah {
		return ErrSikapTemuanTidakSah
	}
	if t.TautanSumber == "" {
		return nil
	}
	return ValidasiTautanSumber(t.TautanSumber)
}

func ValidasiTautanSumber(tautan string) error {
	if len(tautan) > PanjangMaksTautan {
		return ErrTautanTemuanTidakSah
	}
	u, err := url.Parse(tautan)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		return ErrTautanTemuanTidakSah
	}
	return nil
}

func (t Temuan) Status() StatusTemuan {
	switch {
	case t.JumlahMembantah > 0:
		return StatusTemuanDiperdebatkan
	case t.JumlahMendukung >= AmbangKonfirmasiTemuan:
		return StatusTemuanTerkonfirmasi
	default:
		return StatusTemuanMenunggu
	}
}
