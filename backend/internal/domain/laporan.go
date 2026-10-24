package domain

import (
	"math"
	"time"
	"unicode/utf8"
)

type KondisiLaporan string

const (
	KondisiBelumDimulai     KondisiLaporan = "belum_dimulai"
	KondisiSedangDikerjakan KondisiLaporan = "sedang_dikerjakan"
	KondisiTerhenti         KondisiLaporan = "terhenti"
	KondisiSelesaiBaik      KondisiLaporan = "selesai_baik"
	KondisiSelesaiRusak     KondisiLaporan = "selesai_rusak"
)

type StatusLaporan string

const (
	StatusDikirim       StatusLaporan = "DIKIRIM"
	StatusMenungguSaksi StatusLaporan = "MENUNGGU_SAKSI"
	StatusTerverifikasi StatusLaporan = "TERVERIFIKASI"
	StatusDiteruskan    StatusLaporan = "DITERUSKAN"
	StatusDitanggapi    StatusLaporan = "DITANGGAPI"
	StatusSelesai       StatusLaporan = "SELESAI"
	StatusDitolak       StatusLaporan = "DITOLAK"
)

const (
	BatasFotoLaporan              = 4
	BatasUkuranFotoByte           = 8 << 20
	BatasCatatanLaporan           = 500
	BatasAkurasiLokasiMeter       = 100.0
	BatasLaporanPerHari           = 5
	JendelaBatasLaporan           = 24 * time.Hour
	RadiusProyekLaporanMeter      = 500.0
	GridLokasiPublikMeter         = 50.0
	meterPerDerajatLintang        = 111320.0
	SisiTerpanjangFotoPiksel      = 2048
	SisiTerpanjangThumbnailPiksel = 480
)

var (
	ErrAkurasiLokasiBuruk     = InputTidakValid("akurasi lokasi lebih dari 100 m; pindah ke ruang terbuka lalu coba lagi")
	ErrKondisiLaporanTidakSah = InputTidakValid("kondisi harus salah satu dari belum_dimulai, sedang_dikerjakan, terhenti, selesai_baik, selesai_rusak")
	ErrCatatanTerlaluPanjang  = InputTidakValid("catatan paling banyak 500 karakter")
	ErrJumlahFotoTidakSah     = InputTidakValid("lampirkan 1 sampai 4 foto")
	ErrFotoTerlaluBesar       = InputTidakValid("ukuran setiap foto paling besar 8 MB")
	ErrFormatFotoTidakSah     = InputTidakValid("foto harus berformat JPEG atau PNG")
	ErrBatasLaporanHarian     = TerlaluBanyak("batas 5 laporan per hari sudah tercapai; coba lagi besok")
	ErrLokasiDiLuarDKI        = InputTidakValid("lokasi berada di luar DKI Jakarta")
	ErrPaketLaporanJauh       = InputTidakValid("proyek yang dipilih berjarak lebih dari 500 m dari lokasimu")
)

type LaporanMasuk struct {
	PenggunaID   int64
	PaketKode    string
	Titik        Titik
	AkurasiMeter float64
	Kondisi      KondisiLaporan
	Catatan      string
	Foto         [][]byte
}

type FotoDiproses struct {
	Utama     []byte
	Thumbnail []byte
	Lebar     int
	Tinggi    int
}

type FotoTersimpan struct {
	ObjectKey    string
	ThumbnailKey string
	Lebar        int
	Tinggi       int
}

type LaporanBaru struct {
	LaporanMasuk
	KodeWilayah  string
	TitikPublik  Titik
	Foto         []FotoTersimpan
	IsSimulation bool
}

func (k KondisiLaporan) IsValid() bool {
	switch k {
	case KondisiBelumDimulai, KondisiSedangDikerjakan, KondisiTerhenti, KondisiSelesaiBaik, KondisiSelesaiRusak:
		return true
	default:
		return false
	}
}

func (l LaporanMasuk) Validasi() error {
	switch {
	case l.AkurasiMeter < 0 || l.AkurasiMeter > BatasAkurasiLokasiMeter:
		return ErrAkurasiLokasiBuruk
	case !l.Kondisi.IsValid():
		return ErrKondisiLaporanTidakSah
	case utf8.RuneCountInString(l.Catatan) > BatasCatatanLaporan:
		return ErrCatatanTerlaluPanjang
	case len(l.Foto) == 0 || len(l.Foto) > BatasFotoLaporan:
		return ErrJumlahFotoTidakSah
	default:
		return validasiUkuranFoto(l.Foto)
	}
}

func validasiUkuranFoto(foto [][]byte) error {
	for _, isi := range foto {
		if len(isi) > BatasUkuranFotoByte {
			return ErrFotoTerlaluBesar
		}
	}
	return nil
}

func BulatkanLokasiPublik(t Titik) Titik {
	langkahLintang := GridLokasiPublikMeter / meterPerDerajatLintang
	lintang := math.Round(t.Lat/langkahLintang) * langkahLintang
	langkahBujur := GridLokasiPublikMeter / (meterPerDerajatLintang * math.Cos(lintang*math.Pi/180))
	return Titik{Lat: lintang, Lng: math.Round(t.Lng/langkahBujur) * langkahBujur}
}
