package domain

import (
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type JenisTargetHakJawab string

const (
	TargetHakJawabSinyal  JenisTargetHakJawab = "sinyal"
	TargetHakJawabLaporan JenisTargetHakJawab = "laporan"
	TargetHakJawabPaket   JenisTargetHakJawab = "paket"
	TargetHakJawabUmpan   JenisTargetHakJawab = "umpan"

	PanjangMaksHakJawab         = 2000
	PanjangMaksNomorBeritaAcara = 100
	PanjangMaksTautanDokumen    = 500
	AlasanDiteruskan            = "diteruskan otomatis ke antrean Panel Instansi pelaksana paket"
	AlasanDitanggapi            = "instansi memberi hak jawab"
	AlasanSelesaiInstansi       = "ditandai selesai oleh instansi"
	BatasAntreanPanel           = 50
)

var (
	ErrTargetHakJawabTidakSah     = InputTidakValid("sorotan yang ditanggapi harus sinyal, laporan, paket, atau umpan")
	ErrIsiHakJawabTidakSah        = InputTidakValid("hak jawab wajib diisi, maksimal 2.000 karakter")
	ErrTautanDokumenTidakSah      = InputTidakValid("tautan dokumen harus alamat http atau https")
	ErrBukanAkunPanel             = TidakBerwenang("panel hanya untuk akun instansi dan anggota DPRD")
	ErrBukanPenjawabTarget        = TidakBerwenang("akun ini tidak berwenang menanggapi sorotan ini")
	ErrBukanAkunInstansi          = TidakBerwenang("hanya petugas instansi yang dapat melakukan ini")
	ErrTanggalSerahTerimaTidakSah = InputTidakValid("tanggal serah terima wajib diisi dan tidak boleh di masa depan")
	ErrSorotanTidakDitemukan      = TidakDitemukan("sorotan yang ditanggapi tidak ditemukan")
	ErrNomorBeritaAcaraTidakSah   = InputTidakValid("nomor berita acara maksimal 100 karakter")
	targetHakJawabSah             = []JenisTargetHakJawab{TargetHakJawabSinyal, TargetHakJawabLaporan, TargetHakJawabPaket, TargetHakJawabUmpan}
	targetBerIDAngka              = []JenisTargetHakJawab{TargetHakJawabSinyal, TargetHakJawabLaporan, TargetHakJawabUmpan}
)

type HakJawabBaru struct {
	TargetJenis   JenisTargetHakJawab
	TargetID      string
	Isi           string
	TautanDokumen string
}

type PemilikTarget struct {
	StatusLaporan StatusLaporan
	InstansiID    *int64
	NamaInstansi  string
	AnggotaID     []string
	KodeKecamatan string
	IsSimulation  bool
}

type HakJawab struct {
	ID            int64
	TargetJenis   JenisTargetHakJawab
	TargetID      string
	Isi           string
	TautanDokumen string
	NamaPenjawab  string
	Pihak         string
	IsSimulation  bool
	DibuatPada    time.Time
}

type SerahTerima struct {
	KodePaket        string
	Tanggal          time.Time
	NomorBeritaAcara string
}

func (b HakJawabBaru) Rapikan() HakJawabBaru {
	b.Isi = strings.TrimSpace(b.Isi)
	b.TautanDokumen = strings.TrimSpace(b.TautanDokumen)
	b.TargetID = strings.TrimSpace(b.TargetID)
	return b
}

func (b HakJawabBaru) Validasi() error {
	switch {
	case !slices.Contains(targetHakJawabSah, b.TargetJenis) || b.TargetID == "" || !b.isIDTargetSah():
		return ErrTargetHakJawabTidakSah
	case b.Isi == "" || utf8.RuneCountInString(b.Isi) > PanjangMaksHakJawab:
		return ErrIsiHakJawabTidakSah
	case b.TautanDokumen != "" && !isTautanWeb(b.TautanDokumen):
		return ErrTautanDokumenTidakSah
	}
	return nil
}

func (b HakJawabBaru) isIDTargetSah() bool {
	if !slices.Contains(targetBerIDAngka, b.TargetJenis) {
		return true
	}
	_, err := strconv.ParseInt(b.TargetID, 10, 64)
	return err == nil
}

func isTautanWeb(tautan string) bool {
	u, err := url.Parse(tautan)
	return err == nil && len(tautan) <= PanjangMaksTautanDokumen && (u.Scheme == "https" || u.Scheme == "http") && u.Host != ""
}

func (p Pengguna) PeriksaAksesPanel() error {
	if p.Peran != PeranInstansi && p.Peran != PeranAnggota {
		return ErrBukanAkunPanel
	}
	return nil
}

func (p Pengguna) PeriksaBolehMenjawab(t PemilikTarget) error {
	if p.IsSimulation != t.IsSimulation {
		return ErrBukanPenjawabTarget
	}
	switch p.Peran {
	case PeranInstansi:
		if p.InstansiID != nil && t.InstansiID != nil && *p.InstansiID == *t.InstansiID {
			return nil
		}
	case PeranAnggota:
		if p.AnggotaID != nil && slices.Contains(t.AnggotaID, *p.AnggotaID) {
			return nil
		}
	}
	return ErrBukanPenjawabTarget
}

func (p Pengguna) PeriksaPetugasPemilik(t PemilikTarget) error {
	if p.Peran != PeranInstansi {
		return ErrBukanAkunInstansi
	}
	return p.PeriksaBolehMenjawab(t)
}

func (s SerahTerima) Validasi(hariIni time.Time) error {
	if s.Tanggal.IsZero() || s.Tanggal.After(hariIni) {
		return ErrTanggalSerahTerimaTidakSah
	}
	if utf8.RuneCountInString(s.NomorBeritaAcara) > PanjangMaksNomorBeritaAcara {
		return ErrNomorBeritaAcaraTidakSah
	}
	return nil
}
