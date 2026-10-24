package domain

import "time"

type SinyalPanel struct {
	ID             int64
	Jenis          JenisSinyal
	Alasan         string
	KodePaket      string
	NamaPaket      string
	JumlahHakJawab int
}

type LaporanPanel struct {
	ID           int64
	Status       StatusLaporan
	Kondisi      KondisiLaporan
	Catatan      string
	KodePaket    string
	NamaPaket    string
	NamaWilayah  string
	IsSimulation bool
	DibuatPada   time.Time
}

type PaketJadwalPanel struct {
	Kode             string
	Nama             string
	JenisPengadaan   string
	NomorBeritaAcara string
	IsSimulation     bool
	Jadwal           JadwalPaket
}

type UmpanPanel struct {
	ID             int64
	Judul          string
	MediaNama      string
	URL            string
	TerbitPada     time.Time
	JumlahHakJawab int
}

type AntreanInstansi struct {
	NamaInstansi string
	Sinyal       []SinyalPanel
	Laporan      []LaporanPanel
	Paket        []PaketJadwalPanel
}
