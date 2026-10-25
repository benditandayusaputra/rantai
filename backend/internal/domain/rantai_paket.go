package domain

import "time"

const (
	ProvenanceRumusPemeliharaan = "prv_rantai_pemeliharaan"
	ProvenanceLokasiPaket       = "prv_rantai_lokasi_paket"
)

type StatusMata string

type JenisSinyal string

const (
	StatusMataAda      StatusMata = "ada"
	StatusMataPutus    StatusMata = "putus"
	StatusMataMenunggu StatusMata = "menunggu"
	StatusMataEstimasi StatusMata = "estimasi"
)

const (
	MataRencana      = "rencana"
	MataPemilihan    = "pemilihan"
	MataPemenang     = "pemenang"
	MataPelaksanaan  = "pelaksanaan"
	MataLokasi       = "lokasi"
	MataSerahTerima  = "serah_terima"
	MataPemeliharaan = "masa_pemeliharaan"
	MataBuktiWarga   = "bukti_warga"
)

type RUPPaket struct {
	Kode               string
	Nama               string
	SumberDana         string
	Pagu               *float64
	IsDetailTersedia   bool
	PelaksanaanMulai   *time.Time
	PelaksanaanSelesai *time.Time
	ProvenanceID       string
}

type PenyediaPaket struct {
	ID   int64
	Nama string
}

type LokasiTersimpan struct {
	Tingkat     Tingkat
	KodeWilayah string
	NamaWilayah string
	Titik       *Titik
	RT          *int
	RW          *int
}

type SinyalTersimpan struct {
	ID         int64
	HakJawab   []HakJawab
	Jenis      JenisSinyal
	Alasan     string
	Pembanding map[string]any
}

type PaketLengkap struct {
	Kode               string
	KodeSPSE           string
	Nama               string
	Tahap              string
	JenisPengadaan     string
	Metode             string
	JenisKontrak       string
	IsTender           bool
	IsBatal            bool
	IsSimulation       bool
	TahunAnggaran      int
	Pagu               *float64
	HPS                *float64
	NilaiKontrak       *float64
	JumlahPeserta      *int
	TanggalDibuat      *time.Time
	TanggalKontrak     *time.Time
	TanggalSerahTerima *time.Time
	Instansi           string
	Kategori           string
	Sektor             string
	TemaKode           []string
	RUP                *RUPPaket
	Penyedia           *PenyediaPaket
	Lokasi             LokasiTersimpan
	Sinyal             []SinyalTersimpan
	HakJawab           []HakJawab
	SerahTerimaProvID  string
	NomorBeritaAcara   string
	JumlahLaporan      int
	ProvenanceID       string
}

type PaketRingkas struct {
	Kode           string
	Nama           string
	Tahap          string
	JenisPengadaan string
	Metode         string
	TahunAnggaran  int
	HPS            *float64
	NilaiKontrak   *float64
	Sektor         string
	TemaKode       []string
	Lokasi         LokasiTersimpan
	NamaPenyedia   string
	IsSimulation   bool
	ProvenanceID   string
}

type BarisUnduhPaket struct {
	PaketRingkas
	KodeSPSE    string
	IsTender    bool
	Pagu        *float64
	Instansi    string
	KodeRUP     string
	URLSumber   string
	Label       Label
	DiambilPada time.Time
}

type MataRantai struct {
	Kode   string
	Status StatusMata
	Alasan string
}

func SusunRantai(p PaketLengkap, pemeliharaan EstimasiPemeliharaan) []MataRantai {
	return []MataRantai{
		mataRencana(p),
		{Kode: MataPemilihan, Status: StatusMataAda},
		mataPemenang(p),
		mataPelaksanaan(p),
		mataLokasi(p),
		mataSerahTerima(p),
		mataPemeliharaan(pemeliharaan),
		mataBuktiWarga(p),
	}
}

func mataRencana(p PaketLengkap) MataRantai {
	switch {
	case p.RUP == nil && p.IsTender:
		return MataRantai{Kode: MataRencana, Status: StatusMataPutus, Alasan: "kode RUP paket tender tidak dibuka di SPSE"}
	case p.RUP == nil:
		return MataRantai{Kode: MataRencana, Status: StatusMataPutus, Alasan: "rencana tidak ditemukan"}
	case !p.RUP.IsDetailTersedia:
		return MataRantai{Kode: MataRencana, Status: StatusMataAda, Alasan: "kode RUP tercatat, detail tidak tersedia"}
	default:
		return MataRantai{Kode: MataRencana, Status: StatusMataAda}
	}
}

func mataPemenang(p PaketLengkap) MataRantai {
	if p.Penyedia == nil && p.NilaiKontrak == nil {
		return MataRantai{Kode: MataPemenang, Status: StatusMataMenunggu, Alasan: "belum ada pemenang"}
	}
	return MataRantai{Kode: MataPemenang, Status: StatusMataAda}
}

func mataPelaksanaan(p PaketLengkap) MataRantai {
	if p.RUP != nil && p.RUP.PelaksanaanSelesai != nil {
		return MataRantai{Kode: MataPelaksanaan, Status: StatusMataEstimasi, Alasan: "jadwal dari rencana pengadaan, berketelitian bulan"}
	}
	return MataRantai{Kode: MataPelaksanaan, Status: StatusMataPutus, Alasan: "jadwal pelaksanaan tidak tersedia"}
}

func mataLokasi(p PaketLengkap) MataRantai {
	switch p.Lokasi.Tingkat {
	case TingkatKelurahan, TingkatKecamatan:
		return MataRantai{Kode: MataLokasi, Status: StatusMataAda}
	case TingkatKota:
		return MataRantai{Kode: MataLokasi, Status: StatusMataPutus, Alasan: "lokasi tingkat kota"}
	default:
		return MataRantai{Kode: MataLokasi, Status: StatusMataPutus, Alasan: "lokasi tidak disebut di nama paket"}
	}
}

func mataSerahTerima(p PaketLengkap) MataRantai {
	if p.TanggalSerahTerima != nil {
		return MataRantai{Kode: MataSerahTerima, Status: StatusMataAda}
	}
	return MataRantai{Kode: MataSerahTerima, Status: StatusMataEstimasi, Alasan: "tanggal serah terima belum diisi instansi"}
}

func mataPemeliharaan(e EstimasiPemeliharaan) MataRantai {
	if e.Label == LabelResmi {
		return MataRantai{Kode: MataPemeliharaan, Status: StatusMataAda}
	}
	return MataRantai{Kode: MataPemeliharaan, Status: StatusMataEstimasi, Alasan: "perkiraan dari jadwal dan ketentuan masa pemeliharaan"}
}

func mataBuktiWarga(p PaketLengkap) MataRantai {
	if p.JumlahLaporan > 0 {
		return MataRantai{Kode: MataBuktiWarga, Status: StatusMataAda}
	}
	return MataRantai{Kode: MataBuktiWarga, Status: StatusMataMenunggu, Alasan: "belum ada laporan warga"}
}

type TitikPaket struct {
	Kode         string
	Nama         string
	Sektor       string
	IsSimulation bool
	Titik        Titik
	JarakMeter   float64
}

type AreaPaket struct {
	KodeKecamatan   string
	NamaKecamatan   string
	GeometriGeoJSON string
	KodePaket       []string
}

type Kotak struct {
	Barat, Selatan, Timur, Utara float64
}

var ErrKotakTidakWajar = InputTidakValid("bbox wajib berupa barat,selatan,timur,utara dengan luas wajar")

func (k Kotak) Periksa(luasMaks float64) error {
	if k.Barat >= k.Timur || k.Selatan >= k.Utara || (k.Timur-k.Barat)*(k.Utara-k.Selatan) > luasMaks {
		return ErrKotakTidakWajar
	}
	return nil
}
