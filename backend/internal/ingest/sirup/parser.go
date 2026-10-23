package sirup

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/halaman"
	"golang.org/x/net/html"
)

const (
	tagLabel           = "td"
	kelasLabel         = "label-left"
	kolomSumberDana    = 1
	jumlahKolomJadwal  = 2
	labelKode          = "Kode RUP"
	labelJadwalKontrak = "Jadwal Pelaksanaan Kontrak"
)

var (
	polaBulanTahun = regexp.MustCompile(`^(\p{L}+) (\d{4})$`)
	polaBulat      = regexp.MustCompile(`\d+`)
)

var nomorBulan = map[string]time.Month{
	"januari": time.January, "februari": time.February, "maret": time.March, "april": time.April,
	"mei": time.May, "juni": time.June, "juli": time.July, "agustus": time.August,
	"september": time.September, "oktober": time.October, "november": time.November, "desember": time.December,
}

var zonaWIB = time.FixedZone("WIB", 7*60*60)

func UraiRUP(isi []byte) (domain.RUP, error) {
	tabel, err := halaman.BacaTabelLabel(isi, tagLabel, kelasLabel)
	if err != nil {
		return domain.RUP{}, fmt.Errorf("urai html rup: %w", err)
	}
	if tabel.Teks(labelKode) == "" {
		return domain.RUP{}, domain.ErrRUPTidakTersedia
	}
	tahun, _ := strconv.Atoi(tabel.Teks("Tahun Anggaran"))
	mulai, selesai := uraiJadwal(tabel[labelJadwalKontrak])
	return domain.RUP{
		Kode:               tabel.Teks(labelKode),
		Nama:               tabel.Teks("Nama Paket"),
		SatuanKerja:        tabel.Teks("Satuan Kerja"),
		Pagu:               uraiBulat(tabel.Teks("Total Pagu")),
		SumberDana:         sumberDana(tabel["Sumber Dana"]),
		TahunAnggaran:      tahun,
		JenisPengadaan:     strings.TrimRight(tabel.Teks("Jenis Pengadaan"), ", "),
		Metode:             tabel.Teks("Metode Pemilihan"),
		PelaksanaanMulai:   mulai,
		PelaksanaanSelesai: selesai,
	}, nil
}

func sumberDana(sel *html.Node) string {
	daftar := []string{}
	for _, baris := range barisData(sel) {
		if len(baris) > kolomSumberDana && !strings.Contains(strings.Join(daftar, ","), baris[kolomSumberDana]) {
			daftar = append(daftar, baris[kolomSumberDana])
		}
	}
	return strings.Join(daftar, ", ")
}

func uraiJadwal(sel *html.Node) (*time.Time, *time.Time) {
	for _, baris := range barisData(sel) {
		if len(baris) == jumlahKolomJadwal {
			mulai, okMulai := awalBulan(baris[0])
			akhir, okAkhir := awalBulan(baris[1])
			if okMulai && okAkhir {
				selesai := akhir.AddDate(0, 1, -1)
				return &mulai, &selesai
			}
		}
	}
	return nil, nil
}

func barisData(sel *html.Node) [][]string {
	hasil := [][]string{}
	if sel == nil {
		return hasil
	}
	for _, tr := range halaman.CariSemua(sel, func(n *html.Node) bool { return halaman.IsElemen(n, "tr") }) {
		if kolom := halaman.SelBaris(tr); len(kolom) > 0 {
			hasil = append(hasil, kolom)
		}
	}
	return hasil
}

func awalBulan(teks string) (time.Time, bool) {
	bagian := polaBulanTahun.FindStringSubmatch(strings.TrimSpace(teks))
	if bagian == nil {
		return time.Time{}, false
	}
	bulan, ada := nomorBulan[strings.ToLower(bagian[1])]
	tahun, err := strconv.Atoi(bagian[2])
	if !ada || err != nil {
		return time.Time{}, false
	}
	return time.Date(tahun, bulan, 1, 0, 0, 0, 0, zonaWIB), true
}

func uraiBulat(teks string) *float64 {
	nilai, err := strconv.ParseFloat(polaBulat.FindString(teks), 64)
	if err != nil {
		return nil
	}
	return &nilai
}
