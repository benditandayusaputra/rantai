package domain

import (
	"slices"
	"strings"
	"time"
)

type AksiModerasi string

const (
	AksiVerifikasi            AksiModerasi = "verifikasi"
	AksiTolak                 AksiModerasi = "tolak"
	AksiSembunyikanFoto       AksiModerasi = "sembunyikan_foto"
	AksiSembunyikanUmpan      AksiModerasi = "sembunyikan_umpan"
	TargetLaporan                          = "laporan"
	TargetLaporanFoto                      = "laporan_foto"
	TargetUmpan                            = "umpan"
	AlasanVerifikasiModerator              = "ditandai terverifikasi oleh moderator"
	BatasAntreanUmpanModerasi              = 50
)

var (
	AlasanTolakSah              = []string{"spam", "duplikat", "di_luar_lingkup"}
	ErrBukanModerator           = TidakBerwenang("hanya moderator yang dapat memakai panel ini")
	ErrAksiModerasiTidakSah     = InputTidakValid("aksi harus verifikasi, tolak, atau sembunyikan_foto")
	ErrAlasanTolakTidakSah      = InputTidakValid("alasan tolak harus spam, duplikat, atau di_luar_lingkup")
	ErrAlasanModerasiKosong     = InputTidakValid("alasan wajib diisi")
	ErrFotoModerasiTidakDipilih = InputTidakValid("pilih foto yang disembunyikan")
)

type PermintaanModerasiLaporan struct {
	LaporanID   int64
	ModeratorID int64
	Aksi        AksiModerasi
	Alasan      string
	Pesan       string
	FotoID      int64
}

type CatatanModerasi struct {
	ModeratorID int64
	TargetJenis string
	TargetID    int64
	Aksi        AksiModerasi
	Alasan      string
}

type FotoModerasi struct {
	ID            int64
	ObjectKey     string
	ThumbnailKey  string
	Disembunyikan bool
}

type LaporanModerasi struct {
	ID           int64
	Kondisi      KondisiLaporan
	Catatan      string
	NamaWilayah  string
	PaketNama    string
	IsSimulation bool
	DibuatPada   time.Time
	Foto         []FotoModerasi
}

type UmpanModerasi struct {
	ID          int64
	JenisSumber JenisSumberUmpan
	MediaNama   string
	Judul       string
	URL         string
	TerbitPada  time.Time
	Anggota     []string
}

type AntreanModerasi struct {
	Laporan []LaporanModerasi
	Umpan   []UmpanModerasi
}

func IsModerator(peran Peran) bool {
	return peran == PeranModerator || peran == PeranAdmin
}

func (p PermintaanModerasiLaporan) Validasi() error {
	switch p.Aksi {
	case AksiVerifikasi:
		return nil
	case AksiTolak:
		if !slices.Contains(AlasanTolakSah, p.Alasan) {
			return ErrAlasanTolakTidakSah
		}
		return nil
	case AksiSembunyikanFoto:
		return validasiSembunyikanFoto(p)
	default:
		return ErrAksiModerasiTidakSah
	}
}

func validasiSembunyikanFoto(p PermintaanModerasiLaporan) error {
	if p.FotoID <= 0 {
		return ErrFotoModerasiTidakDipilih
	}
	if strings.TrimSpace(p.Alasan) == "" {
		return ErrAlasanModerasiKosong
	}
	return nil
}

func (p PermintaanModerasiLaporan) StatusTujuan() StatusLaporan {
	if p.Aksi == AksiTolak {
		return StatusDitolak
	}
	return StatusTerverifikasi
}

func (p PermintaanModerasiLaporan) AlasanRiwayat() string {
	switch {
	case p.Aksi != AksiTolak:
		return AlasanVerifikasiModerator
	case strings.TrimSpace(p.Pesan) == "":
		return p.Alasan
	default:
		return p.Alasan + ": " + strings.TrimSpace(p.Pesan)
	}
}
