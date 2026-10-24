package domain

import (
	"math"
	"slices"
	"time"
)

const (
	RadiusKonfirmasiMeter = 200.0
	JendelaKonfirmasi     = 14 * 24 * time.Hour
	radiusBumiMeter       = 6371000.0
)

var (
	ErrTransisiTidakSah        = InputTidakValid("perpindahan status laporan tidak sah")
	ErrKonfirmasiOlehPelapor   = InputTidakValid("pelapor tidak dapat menjadi saksi laporannya sendiri")
	ErrLaporanTidakMenunggu    = InputTidakValid("laporan ini tidak sedang menunggu saksi")
	ErrKonfirmasiSudahTercatat = InputTidakValid("kamu sudah mengonfirmasi laporan ini")
	transisiLaporanSah         = map[StatusLaporan][]StatusLaporan{
		StatusDikirim:       {StatusMenungguSaksi, StatusDitolak},
		StatusMenungguSaksi: {StatusTerverifikasi, StatusDitolak},
		StatusTerverifikasi: {StatusDiteruskan},
		StatusDiteruskan:    {StatusDitanggapi},
		StatusDitanggapi:    {StatusSelesai},
	}
)

type Konfirmasi struct {
	PenggunaID int64
	Kondisi    KondisiLaporan
	Titik      Titik
	DibuatPada time.Time
}

type LaporanUntukVerifikasi struct {
	ID         int64
	PelaporID  int64
	Status     StatusLaporan
	Kondisi    KondisiLaporan
	Titik      Titik
	DibuatPada time.Time
	Konfirmasi []Konfirmasi
}

func Transisi(dari, ke StatusLaporan) error {
	if !slices.Contains(transisiLaporanSah[dari], ke) {
		return ErrTransisiTidakSah
	}
	return nil
}

func (l LaporanUntukVerifikasi) JumlahSaksiSah() int {
	akun := map[int64]bool{l.PelaporID: true}
	for _, k := range l.Konfirmasi {
		if l.isKonfirmasiSah(k) {
			akun[k.PenggunaID] = true
		}
	}
	return len(akun)
}

func (l LaporanUntukVerifikasi) isKonfirmasiSah(k Konfirmasi) bool {
	return k.Kondisi == l.Kondisi &&
		JarakMeter(l.Titik, k.Titik) <= RadiusKonfirmasiMeter &&
		!k.DibuatPada.Before(l.DibuatPada) && k.DibuatPada.Sub(l.DibuatPada) <= JendelaKonfirmasi
}

func (l LaporanUntukVerifikasi) IsCukupSaksi() bool {
	return l.JumlahSaksiSah() >= JumlahKonfirmasiTerverifikasi
}

func (l LaporanUntukVerifikasi) PeriksaBisaDikonfirmasi(penggunaID int64) error {
	switch {
	case l.Status != StatusMenungguSaksi:
		return ErrLaporanTidakMenunggu
	case penggunaID == l.PelaporID:
		return ErrKonfirmasiOlehPelapor
	case slices.ContainsFunc(l.Konfirmasi, func(k Konfirmasi) bool { return k.PenggunaID == penggunaID }):
		return ErrKonfirmasiSudahTercatat
	default:
		return nil
	}
}

func JarakMeter(a, b Titik) float64 {
	radian := math.Pi / 180
	dLat, dLng := (b.Lat-a.Lat)*radian, (b.Lng-a.Lng)*radian
	h := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(a.Lat*radian)*math.Cos(b.Lat*radian)*math.Sin(dLng/2)*math.Sin(dLng/2)
	return 2 * radiusBumiMeter * math.Asin(math.Sqrt(h))
}

type KonfirmasiMasuk struct {
	LaporanID    int64
	PenggunaID   int64
	Titik        Titik
	AkurasiMeter float64
	Kondisi      KondisiLaporan
}

type PerubahanStatus struct {
	LaporanID      int64
	Dari           StatusLaporan
	Ke             StatusLaporan
	OlehPenggunaID *int64
	Alasan         string
}

type RiwayatLaporan struct {
	StatusLama StatusLaporan
	StatusBaru StatusLaporan
	Alasan     string
	DibuatPada time.Time
}

type LaporanPublik struct {
	ID           int64
	Status       StatusLaporan
	Kondisi      KondisiLaporan
	Catatan      string
	KodeWilayah  string
	NamaWilayah  string
	TitikPublik  Titik
	PaketKode    string
	PaketNama    string
	IsSimulation bool
	DibuatPada   time.Time
	Foto         []FotoTersimpan
	Riwayat      []RiwayatLaporan
	JumlahSaksi  int
}

func (k KonfirmasiMasuk) Validasi() error {
	switch {
	case k.AkurasiMeter < 0 || k.AkurasiMeter > BatasAkurasiLokasiMeter:
		return ErrAkurasiLokasiBuruk
	case !k.Kondisi.IsValid():
		return ErrKondisiLaporanTidakSah
	default:
		return nil
	}
}
