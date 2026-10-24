package domain

import (
	"strconv"
	"strings"
	"time"
)

const (
	BahasaIndonesia = "id"
	BahasaInggris   = "en"
	desimalRasio    = 2
	awalanButir     = "- "
	pemisahIsian    = " ____________________"
)

var namaBulanIndonesia = [...]string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli",
	"Agustus", "September", "Oktober", "November", "Desember"}

type BagianAspirasi struct {
	Data       Teks `yaml:"data"`
	Wewenang   Teks `yaml:"wewenang"`
	Pengawas   Teks `yaml:"pengawas"`
	Pertanyaan Teks `yaml:"pertanyaan"`
	Jalur      Teks `yaml:"jalur"`
	Pengirim   Teks `yaml:"pengirim"`
}

type KalimatAspirasi struct {
	Indikator       Teks `yaml:"indikator"`
	IndikatorKosong Teks `yaml:"indikator_kosong"`
	RasioKosong     Teks `yaml:"rasio_kosong"`
	Anggaran        Teks `yaml:"anggaran"`
	AnggaranKosong  Teks `yaml:"anggaran_kosong"`
	Sumber          Teks `yaml:"sumber"`
	Instansi        Teks `yaml:"instansi"`
	DasarHukum      Teks `yaml:"dasar_hukum"`
	Komisi          Teks `yaml:"komisi"`
	AnggotaKomisi   Teks `yaml:"anggota_komisi"`
	AnggotaKosong   Teks `yaml:"anggota_kosong"`
	Jalur           Teks `yaml:"jalur"`
	Catatan         Teks `yaml:"catatan"`
}

type KosakataAspirasi struct {
	Bagian  BagianAspirasi  `yaml:"bagian"`
	Kalimat KalimatAspirasi `yaml:"kalimat"`
	Isian   []Teks          `yaml:"isian"`
	Label   map[Label]Teks  `yaml:"label"`
}

type DataBersumber[T any] struct {
	Nilai      T
	Provenance Provenance
}

type BahanAspirasi struct {
	Tema          TemaTersimpan
	Masalah       TanggungJawab
	Kecamatan     string
	Kota          string
	TahunAnggaran int
	Indikator     *DataBersumber[IndikatorKecamatan]
	Anggaran      *DataBersumber[RingkasanAnggaran]
	NamaDapil     string
	AnggotaKomisi []Wakil
	Jalur         []JalurAspirasi
}

func (t Teks) Pilih(bahasa string) string {
	if bahasa == BahasaInggris && t.EN != "" {
		return t.EN
	}
	return t.ID
}

func IsiPlaceholder(teks string, nilai map[string]string) string {
	pasangan := make([]string, 0, 2*len(nilai))
	for kunci, isi := range nilai {
		pasangan = append(pasangan, "{"+kunci+"}", isi)
	}
	return strings.NewReplacer(pasangan...).Replace(teks)
}

func FormatTanggal(t time.Time, bahasa string) string {
	if bahasa == BahasaInggris {
		return t.Format("2 January 2006")
	}
	return strconv.Itoa(t.Day()) + " " + namaBulanIndonesia[t.Month()] + " " + strconv.Itoa(t.Year())
}

func formatDesimal(nilai float64) string {
	return strings.Replace(strconv.FormatFloat(nilai, 'f', desimalRasio, 64), ".", ",", 1)
}

type penyusunAspirasi struct {
	bahan    BahanAspirasi
	kosakata KosakataAspirasi
	bahasa   string
}

func SusunDokumenAspirasi(bahan BahanAspirasi, kosakata KosakataAspirasi, bahasa string) Dokumen {
	p := penyusunAspirasi{bahan: bahan, kosakata: kosakata, bahasa: bahasa}
	b := p.kosakata.Bagian
	return Dokumen{
		Judul: p.isi(p.bahan.Tema.TemplatAspirasi.Judul, nil),
		Bagian: []BagianDokumen{
			{Judul: p.pilih(b.Data), Baris: p.barisData()},
			{Judul: p.pilih(b.Wewenang), Baris: p.barisWewenang()},
			{Judul: p.pilih(b.Pengawas), Baris: p.barisPengawas()},
			{Judul: p.pilih(b.Pertanyaan), Baris: p.barisPertanyaan()},
			{Judul: p.pilih(b.Jalur), Baris: p.barisJalur()},
			{Judul: p.pilih(b.Pengirim), Baris: p.barisPengirim()},
		},
	}
}

func (p penyusunAspirasi) pilih(t Teks) string {
	return t.Pilih(p.bahasa)
}

func (p penyusunAspirasi) isi(t Teks, tambahan map[string]string) string {
	nilai := map[string]string{
		"kecamatan":   p.bahan.Kecamatan,
		"kota":        p.bahan.Kota,
		"tahun_depan": strconv.Itoa(p.bahan.TahunAnggaran + 1),
		"tema":        p.pilih(p.bahan.Tema.Nama),
	}
	for kunci, isi := range tambahan {
		nilai[kunci] = isi
	}
	return IsiPlaceholder(p.pilih(t), nilai)
}

func (p penyusunAspirasi) barisData() []string {
	return []string{awalanButir + p.kalimatIndikator(), awalanButir + p.kalimatAnggaran()}
}

func (p penyusunAspirasi) kalimatIndikator() string {
	k := p.kosakata.Kalimat
	judul := map[string]string{"judul": p.pilih(p.bahan.Tema.Indikator.Judul)}
	if p.bahan.Indikator == nil {
		return p.isi(k.IndikatorKosong, judul)
	}
	i := p.bahan.Indikator.Nilai
	ind := p.bahan.Tema.Indikator
	judul["periode"] = i.Periode
	judul["kebutuhan"] = FormatRibuan(i.Kebutuhan)
	judul["satuan_kebutuhan"] = ind.SatuanKebutuhan
	judul["kapasitas"] = FormatRibuan(i.Kapasitas)
	judul["satuan_kapasitas"] = ind.SatuanKapasitas
	judul["kesenjangan"] = FormatRibuan(i.Kesenjangan)
	judul["rasio"] = p.rasio(i.Rasio)
	return p.isi(k.Indikator, judul) + " " + p.kalimatSumber(p.bahan.Indikator.Provenance, i.Label)
}

func (p penyusunAspirasi) rasio(rasio *float64) string {
	if rasio == nil {
		return p.pilih(p.kosakata.Kalimat.RasioKosong)
	}
	return formatDesimal(*rasio)
}

func (p penyusunAspirasi) kalimatAnggaran() string {
	k := p.kosakata.Kalimat
	nilai := map[string]string{"tahun": strconv.Itoa(p.bahan.TahunAnggaran)}
	if p.bahan.Anggaran == nil {
		return p.isi(k.AnggaranKosong, nilai)
	}
	a := p.bahan.Anggaran.Nilai
	nilai["jumlah"] = FormatRibuan(float64(a.JumlahPaket))
	nilai["total"] = FormatRupiah(a.TotalHPS)
	return p.isi(k.Anggaran, nilai) + " " + p.kalimatSumber(p.bahan.Anggaran.Provenance, a.Label())
}

func (p penyusunAspirasi) kalimatSumber(prov Provenance, label Label) string {
	sumber := string(prov.Sumber)
	if prov.URL != "" {
		sumber += " (" + prov.URL + ")"
	}
	return p.isi(p.kosakata.Kalimat.Sumber, map[string]string{
		"sumber":  sumber,
		"tanggal": FormatTanggal(prov.DiambilPada, p.bahasa),
		"label":   p.pilih(p.kosakata.Label[label]),
	})
}

func (p penyusunAspirasi) barisWewenang() []string {
	m := p.bahan.Masalah
	k := p.kosakata.Kalimat
	return []string{
		awalanButir + p.isi(k.Instansi, map[string]string{
			"instansi": m.InstansiPelaksana,
			"unit":     IsiPlaceholder(m.UnitWilayah, map[string]string{"kota": p.bahan.Kota}),
		}),
		awalanButir + p.isi(k.DasarHukum, map[string]string{"dasar": m.DasarHukum, "url": m.DasarHukumURL}),
	}
}

func (p penyusunAspirasi) barisPengawas() []string {
	k := p.kosakata.Kalimat
	nilai := map[string]string{"komisi": p.bahan.Masalah.Komisi, "dapil": p.bahan.NamaDapil}
	baris := []string{p.isi(k.Komisi, nilai)}
	if len(p.bahan.AnggotaKomisi) == 0 {
		return append(baris, p.isi(k.AnggotaKosong, nilai))
	}
	baris = append(baris, p.isi(k.AnggotaKomisi, nilai))
	for _, w := range p.bahan.AnggotaKomisi {
		baris = append(baris, awalanButir+w.Nama+" ("+w.Fraksi+")")
	}
	return baris
}

func (p penyusunAspirasi) barisPertanyaan() []string {
	baris := []string{}
	for i, pertanyaan := range p.bahan.Tema.TemplatAspirasi.Pertanyaan {
		baris = append(baris, strconv.Itoa(i+1)+". "+p.isi(pertanyaan, nil))
	}
	return baris
}

func (p penyusunAspirasi) barisJalur() []string {
	baris := []string{}
	for _, j := range p.bahan.Jalur {
		baris = append(baris, awalanButir+p.isi(p.kosakata.Kalimat.Jalur, map[string]string{
			"nama": p.pilih(j.Nama), "peran": p.pilih(j.Peran), "cara": p.pilih(j.Cara), "dasar": j.DasarHukum,
		}))
	}
	return baris
}

func (p penyusunAspirasi) barisPengirim() []string {
	baris := []string{p.pilih(p.kosakata.Kalimat.Catatan)}
	for _, isian := range p.kosakata.Isian {
		baris = append(baris, p.pilih(isian)+pemisahIsian)
	}
	return baris
}
