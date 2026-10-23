package lpse

import (
	"errors"
	"fmt"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"golang.org/x/net/html"
)

const (
	lencanaBatal   = "Paket Batal"
	jumlahKolomRUP = 3
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
	tabel, err := bacaTabelLabel(isi)
	if err != nil {
		return Pengumuman{}, fmt.Errorf("urai html pengumuman: %w", err)
	}
	if tabel.teks("Kode Paket") == "" {
		return Pengumuman{}, ErrBukanPengumuman
	}
	p := identitasPengumuman(tabel)
	if err := lengkapiNilai(&p, tabel); err != nil {
		return Pengumuman{}, fmt.Errorf("pengumuman %s: %w", p.Kode, err)
	}
	return p, nil
}

func identitasPengumuman(tabel tabelLabel) Pengumuman {
	nama := tabel["Nama Paket"]
	return Pengumuman{
		Kode:           tabel.teks("Kode Paket"),
		Nama:           teksTanpaLencana(nama),
		IsBatal:        memuatLencana(nama, lencanaBatal),
		AlasanBatal:    tabel.teks("Alasan Pembatalan"),
		RUP:            uraiRUPTercatat(tabel["Rencana Umum Pengadaan"]),
		Tahap:          tabel.teks("Tahap Paket Saat Ini"),
		Instansi:       tabel.teks("K/L/PD/Instansi Lainnya"),
		SatuanKerja:    tabel.teks("Satuan Kerja"),
		JenisPengadaan: tabel.teks("Jenis Pengadaan"),
		Metode:         tabel.teks("Metode Pengadaan"),
		JenisKontrak:   tabel.teks("Jenis Kontrak"),
		Lokasi:         uraiDaftarButir(tabel["Lokasi Pekerjaan"]),
	}
}

func lengkapiNilai(p *Pengumuman, tabel tabelLabel) error {
	tanggal, err := uraiTanggal(tabel.teks("Tanggal Pembuatan"))
	if err != nil {
		return err
	}
	p.TanggalDibuat = tanggal
	if tahun := uraiDaftarTahun(tabel.teks("Tahun Anggaran")); len(tahun) > 0 {
		p.TahunAnggaran = tahun[len(tahun)-1]
	}
	p.Pagu, _ = uraiRupiah(tabel.teks("Nilai Pagu Paket"))
	p.HPS, _ = uraiRupiah(tabel.teks("Nilai HPS Paket"))
	p.JumlahPeserta, _ = uraiBilanganPertama(tabel.teks("Peserta Non Tender") + tabel.teks("Peserta Tender"))
	return nil
}

func uraiRUPTercatat(sel *html.Node) []domain.RUPTercatat {
	daftar := []domain.RUPTercatat{}
	if sel == nil {
		return daftar
	}
	for _, baris := range cariSemua(sel, func(n *html.Node) bool { return isElemen(n, "tr") }) {
		kolom := selBaris(baris)
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
	for _, butir := range cariSemua(sel, func(n *html.Node) bool { return isElemen(n, "li") }) {
		daftar = append(daftar, teks(butir))
	}
	return daftar
}

func memuatLencana(sel *html.Node, judul string) bool {
	if sel == nil {
		return false
	}
	for _, lencana := range cariSemua(sel, isLencana) {
		if teks(lencana) == judul {
			return true
		}
	}
	return false
}
