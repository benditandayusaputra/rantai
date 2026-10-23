package domain

import "time"

type PengambilanSumber struct {
	Job            string
	StatusTerakhir string
	DijalankanPada time.Time
	BerhasilPada   *time.Time
}

type MediaTerverifikasi struct {
	Nama   string `yaml:"nama"`
	Domain string `yaml:"domain"`
}

type InfoSumberData struct {
	Pengambilan        []PengambilanSumber
	Media              []MediaTerverifikasi
	AtribusiOpentender string
}
