package domain

type RiwayatPenyedia struct {
	NamaPenyedia string
	Tahun        int
	Pembeli      string
	NamaPaket    string
	Nilai        *float64
	KodeSumber   string
}

type RiwayatTersimpan struct {
	RiwayatPenyedia
	Sumber       SumberData
	ProvenanceID string
}

type RingkasanTahunPenyedia struct {
	Tahun        int
	Sumber       SumberData
	JumlahPaket  int
	TotalNilai   float64
	ProvenanceID string
}

type JumlahPerNama struct {
	Nama   string
	Jumlah int
}

type ProfilPenyedia struct {
	ID              int64
	Nama            string
	Alias           []string
	PerTahun        []RingkasanTahunPenyedia
	KategoriUtama   []JumlahPerNama
	InstansiPemberi []JumlahPerNama
	Paket           []PaketRingkas
	Riwayat         []RiwayatTersimpan
}
