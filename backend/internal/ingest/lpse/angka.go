package lpse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	polaRupiah       = regexp.MustCompile(`\d[\d.]*(,\d+)?`)
	polaNilaiRingkas = regexp.MustCompile(`^([\d.]+(?:,\d+)?)\s*(Rb|Jt|M|T)?$`)
	polaTahun        = regexp.MustCompile(`\b(20\d{2})\b`)
	polaBilangan     = regexp.MustCompile(`\d+`)
	polaTanggal      = regexp.MustCompile(`^(\d{1,2}) (\p{L}+) (\d{4})(?: (\d{1,2}):(\d{2}))?$`)
)

var pengaliSatuan = map[string]float64{"": 1, "Rb": 1e3, "Jt": 1e6, "M": 1e9, "T": 1e12}

var nomorBulan = map[string]time.Month{
	"januari": time.January, "februari": time.February, "maret": time.March, "april": time.April,
	"mei": time.May, "juni": time.June, "juli": time.July, "agustus": time.August,
	"september": time.September, "oktober": time.October, "november": time.November, "desember": time.December,
}

var zonaWIB = time.FixedZone("WIB", 7*60*60)

func uraiRupiah(teks string) (float64, bool) {
	cocok := polaRupiah.FindString(teks)
	if cocok == "" || !strings.Contains(teks, "Rp") {
		return 0, false
	}
	return uraiDesimalIndonesia(cocok)
}

func uraiNilaiRingkas(teks string) (float64, bool) {
	bagian := polaNilaiRingkas.FindStringSubmatch(strings.TrimSpace(teks))
	if bagian == nil {
		return 0, false
	}
	nilai, ok := uraiDesimalIndonesia(bagian[1])
	return nilai * pengaliSatuan[bagian[2]], ok
}

func uraiDesimalIndonesia(teks string) (float64, bool) {
	baku := strings.ReplaceAll(strings.ReplaceAll(teks, ".", ""), ",", ".")
	nilai, err := strconv.ParseFloat(baku, 64)
	return nilai, err == nil
}

func uraiDaftarTahun(teks string) []int {
	daftar := []int{}
	for _, cocok := range polaTahun.FindAllString(teks, -1) {
		tahun, _ := strconv.Atoi(cocok)
		daftar = append(daftar, tahun)
	}
	return daftar
}

func uraiBilanganPertama(teks string) (int, bool) {
	bilangan, err := strconv.Atoi(polaBilangan.FindString(teks))
	return bilangan, err == nil
}

func uraiTanggal(teks string) (time.Time, error) {
	bagian := polaTanggal.FindStringSubmatch(strings.TrimSpace(teks))
	if bagian == nil {
		return time.Time{}, fmt.Errorf("format tanggal tidak dikenal: %q", teks)
	}
	bulan, ada := nomorBulan[strings.ToLower(bagian[2])]
	if !ada {
		return time.Time{}, fmt.Errorf("nama bulan tidak dikenal: %q", bagian[2])
	}
	hari, _ := strconv.Atoi(bagian[1])
	tahun, _ := strconv.Atoi(bagian[3])
	jam, _ := strconv.Atoi(bagian[4])
	menit, _ := strconv.Atoi(bagian[5])
	return time.Date(tahun, bulan, hari, jam, menit, 0, 0, zonaWIB), nil
}
