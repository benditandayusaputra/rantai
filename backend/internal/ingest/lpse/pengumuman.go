package lpse

import (
	"errors"
	"fmt"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/halaman"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"golang.org/x/net/html"
)

const (
	lencanaBatal   = "Paket Batal"
	jumlahKolomRUP = 3
	tagLabel       = "th"
	kelasLabel     = "bgwarning"
)

var ErrBukanPengumuman = errors.New("halaman bukan pengumuman paket")

type Pengumuman struct {
	Kode           string
	Nama           string
	IsBatal        bool
	AlasanBatal    string
	RUP            []domain.RUPTercatat
	TanggalDibuat  time.Time
	Tahap          string
	Instansi       string
	SatuanKerja    string
	JenisPengadaan string
	Metode         string
	JenisKontrak   string
	TahunAnggaran  int
	Pagu           float64
	HPS            float64
	Lokasi         []string
	JumlahPeserta  int
}

func UraiPengumuman(isi []byte) (Pengumuman, error) {
	tabel, err := halaman.BacaTabelLabel(isi, tagLabel, kelasLabel)
	if err != nil {
		return Pengumuman{}, fmt.Errorf("urai html pengumuman: %w", err)
	}
	if tabel.Teks("Kode Paket") == "" {
		return Pengumuman{}, ErrBukanPengumuman
	}
	p := identitasPengumuman(tabel)
	if err := lengkapiNilai(&p, tabel); err != nil {
		return Pengumuman{}, fmt.Errorf("pengumuman %s: %w", p.Kode, err)
	}
	return p, nil
}

func identitasPengumuman(tabel halaman.TabelLabel) Pengumuman {
	nama := tabel["Nama Paket"]
	return Pengumuman{
		Kode:           tabel.Teks("Kode Paket"),
		Nama:           halaman.TeksTanpaLencana(nama),
		IsBatal:        memuatLencana(nama, lencanaBatal),
		AlasanBatal:    tabel.Teks("Alasan Pembatalan"),
		RUP:            uraiRUPTercatat(tabel["Rencana Umum Pengadaan"]),
		Tahap:          tabel.Teks("Tahap Paket Saat Ini"),
		Instansi:       tabel.Teks("K/L/PD/Instansi Lainnya"),
		SatuanKerja:    tabel.Teks("Satuan Kerja"),
		JenisPengadaan: tabel.Teks("Jenis Pengadaan"),
		Metode:         tabel.Teks("Metode Pengadaan"),
		JenisKontrak:   tabel.Teks("Jenis Kontrak"),
		Lokasi:         uraiDaftarButir(tabel["Lokasi Pekerjaan"]),
	}
}

func lengkapiNilai(p *Pengumuman, tabel halaman.TabelLabel) error {
	tanggal, err := uraiTanggal(tabel.Teks("Tanggal Pembuatan"))
	if err != nil {
		return err
	}
	p.TanggalDibuat = tanggal
	if tahun := uraiDaftarTahun(tabel.Teks("Tahun Anggaran")); len(tahun) > 0 {
		p.TahunAnggaran = tahun[len(tahun)-1]
	}
	p.Pagu, _ = uraiRupiah(tabel.Teks("Nilai Pagu Paket"))
	p.HPS, _ = uraiRupiah(tabel.Teks("Nilai HPS Paket"))
	p.JumlahPeserta, _ = uraiBilanganPertama(tabel.Teks("Peserta Non Tender") + tabel.Teks("Peserta Tender"))
	return nil
}

func uraiRUPTercatat(sel *html.Node) []domain.RUPTercatat {
	daftar := []domain.RUPTercatat{}
	if sel == nil {
		return daftar
	}
	for _, baris := range halaman.CariSemua(sel, func(n *html.Node) bool { return halaman.IsElemen(n, "tr") }) {
		kolom := halaman.SelBaris(baris)
		if len(kolom) == jumlahKolomRUP && kolom[0] != "" {
			daftar = append(daftar, domain.RUPTercatat{Kode: kolom[0], Nama: kolom[1], SumberDana: kolom[2]})
		}
	}
	return daftar
}

func uraiDaftarButir(sel *html.Node) []string {
	daftar := []string{}
	if sel == nil {
		return daftar
	}
	for _, butir := range halaman.CariSemua(sel, func(n *html.Node) bool { return halaman.IsElemen(n, "li") }) {
		daftar = append(daftar, halaman.Teks(butir))
	}
	return daftar
}

func memuatLencana(sel *html.Node, judul string) bool {
	if sel == nil {
		return false
	}
	for _, lencana := range halaman.CariSemua(sel, halaman.IsLencana) {
		if halaman.Teks(lencana) == judul {
			return true
		}
	}
	return false
}
