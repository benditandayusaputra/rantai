package domain

import "time"

type JenisKejadian string

const (
	KejadianPaketBaru            JenisKejadian = "paket_baru"
	KejadianMelewatiJadwal       JenisKejadian = "melewati_jadwal"
	KejadianPemeliharaanBerakhir JenisKejadian = "pemeliharaan_berakhir"
	KejadianLaporanTerverifikasi JenisKejadian = "laporan_terverifikasi"
	KejadianHakJawab             JenisKejadian = "hak_jawab"
	KejadianIndikatorDiperbarui  JenisKejadian = "indikator_diperbarui"
	KanalFeed                                  = "feed"
	JendelaKejadian                            = 7 * 24 * time.Hour
	NotifikasiPerHalaman                       = 30
)

var ErrLanggananTidakDitemukan = TidakDitemukan("langganan tidak ditemukan")

type Langganan struct {
	ID          int64
	KodeWilayah string
	NamaWilayah string
	DibuatPada  time.Time
}

type Notifikasi struct {
	ID            int64
	Jenis         JenisKejadian
	KodeKecamatan string
	NamaKecamatan string
	Judul         string
	Tautan        string
	TerjadiPada   time.Time
	IsDibaca      bool
	IsSimulation  bool
}

type HalamanNotifikasi struct {
	Item              []Notifikasi
	JumlahBelumDibaca int
}
