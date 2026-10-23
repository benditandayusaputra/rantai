package domain

import (
	"hash/fnv"
	"maps"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	KodeProvinsiDKI  = "31"
	meterPerDerajat  = 111320.0
	kataKelurahan    = "kelurahan"
	kataKecamatan    = "kecamatan"
	pembagiHashSudut = math.MaxUint32
	resolusiJarak    = 1000
)

var (
	polaRT    = regexp.MustCompile(`\brt\.?\s*0*(\d{1,3})\b`)
	polaRW    = regexp.MustCompile(`\brw\.?\s*0*(\d{1,3})\b`)
	polaRTRW  = regexp.MustCompile(`\brt\.?\s*0*(\d{1,3})\s*/\s*0*(\d{1,3})\b`)
	polaKata  = regexp.MustCompile(`[\p{L}\p{N}]+`)
	polaSpasi = regexp.MustCompile(`\s+`)
)

type EntriGazetir struct {
	Kode          string
	Tingkat       Tingkat
	Nama          []string
	KodeKota      string
	KodeKecamatan string
	TitikTengah   Titik
}

type AturanLokasi struct {
	Singkatan           map[string]string `yaml:"singkatan"`
	PergeseranMaksMeter float64           `yaml:"pergeseran_maks_meter"`
}

type PenentuLokasi struct {
	Aturan   AturanLokasi
	indeks   map[string][]EntriGazetir
	maksKata int
}

type MasukanLokasi struct {
	KodePaket        string
	NamaPaket        string
	KodeKotaInstansi string
}

type LokasiPaket struct {
	Tingkat     Tingkat
	KodeWilayah string
	Titik       *Titik
	RT          *int
	RW          *int
}

type kecocokan struct {
	entri     EntriGazetir
	awal      int
	akhir     int
	berawalan bool
}

func PenentuLokasiBaru(aturan AturanLokasi, gazetir []EntriGazetir) PenentuLokasi {
	p := PenentuLokasi{Aturan: aturan, indeks: map[string][]EntriGazetir{}}
	for _, entri := range gazetir {
		entri.Nama = p.namaGazetir(entri.Nama)
		for _, nama := range entri.Nama {
			kunci := strings.ReplaceAll(nama, " ", "")
			p.indeks[kunci] = append(p.indeks[kunci], entri)
			p.maksKata = max(p.maksKata, len(strings.Fields(nama))+1)
		}
	}
	return p
}

func (p PenentuLokasi) namaGazetir(daftar []string) []string {
	hasil := []string{}
	for _, nama := range daftar {
		kata := polaKata.FindAllString(p.normalisasi(nama), -1)
		for len(kata) > 1 && (kata[0] == kataKelurahan || kata[0] == kataKecamatan) {
			kata = kata[1:]
		}
		if bersih := strings.Join(kata, " "); bersih != "" && !slices.Contains(hasil, bersih) {
			hasil = append(hasil, bersih)
		}
	}
	return hasil
}

func (p PenentuLokasi) Tentukan(m MasukanLokasi) LokasiPaket {
	teks := p.normalisasi(m.NamaPaket)
	lokasi := p.tentukanWilayah(m, teks)
	lokasi.RT, lokasi.RW = ekstrakRTRW(teks)
	return lokasi
}

func (p PenentuLokasi) normalisasi(nama string) string {
	teks := strings.ToLower(nama)
	for _, singkatan := range slices.Sorted(maps.Keys(p.Aturan.Singkatan)) {
		teks = strings.ReplaceAll(teks, singkatan, " "+p.Aturan.Singkatan[singkatan]+" ")
	}
	return strings.TrimSpace(polaSpasi.ReplaceAllString(teks, " "))
}

func ekstrakRTRW(teks string) (*int, *int) {
	if bagian := polaRTRW.FindStringSubmatch(teks); bagian != nil {
		return bilanganDari(bagian[1]), bilanganDari(bagian[2])
	}
	return bilanganPola(polaRT, teks), bilanganPola(polaRW, teks)
}

func bilanganPola(pola *regexp.Regexp, teks string) *int {
	if bagian := pola.FindStringSubmatch(teks); bagian != nil {
		return bilanganDari(bagian[1])
	}
	return nil
}

func bilanganDari(teks string) *int {
	bilangan, err := strconv.Atoi(teks)
	if err != nil {
		return nil
	}
	return &bilangan
}

func (p PenentuLokasi) tentukanWilayah(m MasukanLokasi, teks string) LokasiPaket {
	cocok := saringKota(p.cariKecocokan(polaKata.FindAllString(teks, -1)), m.KodeKotaInstansi)
	kelurahan, kecamatan := pisahTingkat(cocok)
	if len(kelurahan) == 1 && konsisten(kelurahan[0], kecamatan) {
		titik := p.geser(kelurahan[0].TitikTengah, m.KodePaket)
		return LokasiPaket{Tingkat: TingkatKelurahan, KodeWilayah: kelurahan[0].Kode, Titik: &titik}
	}
	if kodeKecamatan := kecamatanTunggal(kelurahan, kecamatan); kodeKecamatan != "" {
		return LokasiPaket{Tingkat: TingkatKecamatan, KodeWilayah: kodeKecamatan}
	}
	if m.KodeKotaInstansi != "" {
		return LokasiPaket{Tingkat: TingkatKota, KodeWilayah: m.KodeKotaInstansi}
	}
	return LokasiPaket{Tingkat: TingkatProvinsi, KodeWilayah: KodeProvinsiDKI}
}

func (p PenentuLokasi) cariKecocokan(kata []string) []kecocokan {
	semua := []kecocokan{}
	for awal := range kata {
		for akhir := awal + 1; akhir <= min(len(kata), awal+p.maksKata); akhir++ {
			for _, entri := range p.indeks[strings.Join(kata[awal:akhir], "")] {
				semua = append(semua, cocokkan(kata, awal, akhir, entri)...)
			}
		}
	}
	return utamakanBerawalan(buangYangTercakup(semua))
}

func cocokkan(kata []string, awal, akhir int, entri EntriGazetir) []kecocokan {
	awalan := ""
	if awal > 0 {
		awalan = kata[awal-1]
	}
	if bertentangan(awalan, entri.Tingkat) {
		return nil
	}
	return []kecocokan{{entri: entri, awal: awal, akhir: akhir, berawalan: awalan == string(entri.Tingkat)}}
}

func bertentangan(awalan string, tingkat Tingkat) bool {
	return (awalan == kataKelurahan && tingkat != TingkatKelurahan) || (awalan == kataKecamatan && tingkat != TingkatKecamatan)
}

func buangYangTercakup(daftar []kecocokan) []kecocokan {
	hasil := []kecocokan{}
	for _, c := range daftar {
		tercakup := slices.ContainsFunc(daftar, func(lain kecocokan) bool {
			return lain.awal <= c.awal && lain.akhir >= c.akhir && lain.akhir-lain.awal > c.akhir-c.awal
		})
		if !tercakup {
			hasil = append(hasil, c)
		}
	}
	return hasil
}

func utamakanBerawalan(daftar []kecocokan) []kecocokan {
	berawalan := slices.DeleteFunc(slices.Clone(daftar), func(c kecocokan) bool { return !c.berawalan })
	if len(berawalan) > 0 {
		return berawalan
	}
	return daftar
}

func saringKota(daftar []kecocokan, kodeKota string) []kecocokan {
	if kodeKota == "" {
		return daftar
	}
	return slices.DeleteFunc(daftar, func(c kecocokan) bool { return c.entri.KodeKota != kodeKota })
}

func pisahTingkat(daftar []kecocokan) ([]EntriGazetir, []EntriGazetir) {
	kelurahan, kecamatan := []EntriGazetir{}, []EntriGazetir{}
	for _, c := range daftar {
		if c.entri.Tingkat == TingkatKelurahan {
			kelurahan = tambahUnik(kelurahan, c.entri)
		} else {
			kecamatan = tambahUnik(kecamatan, c.entri)
		}
	}
	return kelurahan, kecamatan
}

func tambahUnik(daftar []EntriGazetir, entri EntriGazetir) []EntriGazetir {
	if slices.ContainsFunc(daftar, func(e EntriGazetir) bool { return e.Kode == entri.Kode }) {
		return daftar
	}
	return append(daftar, entri)
}

func konsisten(kelurahan EntriGazetir, kecamatan []EntriGazetir) bool {
	return !slices.ContainsFunc(kecamatan, func(k EntriGazetir) bool { return k.Kode != kelurahan.KodeKecamatan })
}

func kecamatanTunggal(kelurahan, kecamatan []EntriGazetir) string {
	kode := []string{}
	for _, k := range kecamatan {
		kode = append(kode, k.Kode)
	}
	for _, k := range kelurahan {
		kode = append(kode, k.KodeKecamatan)
	}
	slices.Sort(kode)
	if kode = slices.Compact(kode); len(kode) == 1 {
		return kode[0]
	}
	return ""
}

func (p PenentuLokasi) geser(titik Titik, kodePaket string) Titik {
	h := fnv.New32a()
	_, _ = h.Write([]byte(kodePaket))
	nilai := h.Sum32()
	sudut := 2 * math.Pi * float64(nilai) / pembagiHashSudut
	jarak := p.Aturan.PergeseranMaksMeter * float64(nilai%resolusiJarak) / resolusiJarak / meterPerDerajat
	return Titik{
		Lat: titik.Lat + jarak*math.Sin(sudut),
		Lng: titik.Lng + jarak*math.Cos(sudut)/math.Cos(titik.Lat*math.Pi/180),
	}
}
