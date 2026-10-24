package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type AgendaSikap string

const (
	AgendaRaperda AgendaSikap = "raperda"
	AgendaAPBD    AgendaSikap = "apbd"
)

var (
	ErrAgendaSikapTidakSah    = errors.New("agenda wajib raperda atau apbd")
	ErrEntriParipurnaKosong   = errors.New("kolom wajib entri paripurna kosong")
	ErrTautanSumberTidakSah   = errors.New("tautan sumber wajib URL http atau https")
	ErrTanggalParipurnaKosong = errors.New("tanggal entri paripurna wajib diisi")
	ErrBidangTidakSah         = errors.New("bidang wajib berhuruf kecil")
)

type SikapFraksi struct {
	Fraksi       string
	Agenda       AgendaSikap
	JudulAgenda  string
	Tanggal      time.Time
	Ringkasan    string
	URL          string
	ProvenanceID string
}

type AspirasiReses struct {
	MasaReses        string
	TanggalParipurna time.Time
	Bidang           string
	Uraian           string
	KodeWilayah      string
	NamaWilayah      string
	TingkatWilayah   Tingkat
	URL              string
	ProvenanceID     string
}

type FilterAspirasiReses struct {
	DapilID int
	Bidang  string
}

func (s SikapFraksi) Validasi() error {
	switch {
	case s.Agenda != AgendaRaperda && s.Agenda != AgendaAPBD:
		return fmt.Errorf("%q: %w", s.Agenda, ErrAgendaSikapTidakSah)
	case isAdaKosong(s.Fraksi, s.JudulAgenda, s.Ringkasan):
		return ErrEntriParipurnaKosong
	default:
		return validasiSumberParipurna(s.Tanggal, s.URL)
	}
}

func (a AspirasiReses) Validasi() error {
	switch {
	case isAdaKosong(a.MasaReses, a.Bidang, a.Uraian):
		return ErrEntriParipurnaKosong
	case a.Bidang != strings.ToLower(a.Bidang):
		return fmt.Errorf("%q: %w", a.Bidang, ErrBidangTidakSah)
	default:
		return validasiSumberParipurna(a.TanggalParipurna, a.URL)
	}
}

func validasiSumberParipurna(tanggal time.Time, alamat string) error {
	if tanggal.IsZero() {
		return ErrTanggalParipurnaKosong
	}
	if !isURLSah(alamat) {
		return fmt.Errorf("%q: %w", alamat, ErrTautanSumberTidakSah)
	}
	return nil
}

func isAdaKosong(nilai ...string) bool {
	for _, n := range nilai {
		if strings.TrimSpace(n) == "" {
			return true
		}
	}
	return false
}
