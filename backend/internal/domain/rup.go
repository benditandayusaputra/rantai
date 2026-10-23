package domain

import (
	"errors"
	"time"
)

var ErrRUPTidakTersedia = errors.New("detail RUP tidak tersedia di sumber")

type RUP struct {
	Kode               string
	Nama               string
	SatuanKerja        string
	Pagu               *float64
	SumberDana         string
	TahunAnggaran      int
	JenisPengadaan     string
	Metode             string
	PelaksanaanMulai   *time.Time
	PelaksanaanSelesai *time.Time
}
