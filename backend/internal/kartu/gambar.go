package kartu

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"
)

const (
	lebarKartu      = 1080
	tinggiKartu     = 1350
	tepi            = 72
	lebarIsi        = lebarKartu - 2*tepi
	tinggiKepala    = 168
	tinggiBlok      = 220
	jarakBlok       = 20
	radiusBlok      = 28
	tebalGaris      = 3
	pembulatanKurva = 0.5523
	dpiLayar        = 72
)

var (
	warnaLatar    = color.RGBA{0xf9, 0xfa, 0xf8, 0xff}
	warnaKartu    = color.RGBA{0xff, 0xff, 0xff, 0xff}
	warnaGaris    = color.RGBA{0xd5, 0xda, 0xd6, 0xff}
	warnaTeks     = color.RGBA{0x14, 0x19, 0x1c, 0xff}
	warnaRedup    = color.RGBA{0x4e, 0x58, 0x5e, 0xff}
	warnaPrimer   = color.RGBA{0x0f, 0x59, 0x64, 0xff}
	warnaSekunder = color.RGBA{0xe2, 0xec, 0xeb, 0xff}
	warnaAksen    = color.RGBA{0xe9, 0xa9, 0x3a, 0xff}
)

type pasanganWarna struct{ latar, teks color.RGBA }

var warnaLabel = map[domain.Label]pasanganWarna{
	domain.LabelResmi:    {warnaPrimer, warnaKartu},
	domain.LabelProksi:   {color.RGBA{0xe5, 0xec, 0xfa, 0xff}, color.RGBA{0x25, 0x56, 0xb8, 0xff}},
	domain.LabelEstimasi: {color.RGBA{0xe6, 0xe9, 0xe5, 0xff}, color.RGBA{0x5a, 0x63, 0x6a, 0xff}},
	domain.LabelSimulasi: {color.RGBA{0xfb, 0xe5, 0xf1, 0xff}, color.RGBA{0x8e, 0x1b, 0x5e, 0xff}},
}

type gaya struct {
	isTebal bool
	ukuran  float64
}

var (
	gayaMerek       = gaya{true, 52}
	gayaTagline     = gaya{false, 26}
	gayaJudul       = gaya{true, 64}
	gayaKota        = gaya{false, 34}
	gayaJudulBlok   = gaya{false, 28}
	gayaNilai       = gaya{true, 72}
	gayaNilaiKosong = gaya{true, 46}
	gayaSatuan      = gaya{false, 30}
	gayaKeterangan  = gaya{false, 24}
	gayaChip        = gaya{true, 22}
	gayaTautan      = gaya{true, 34}
	semuaGaya       = []gaya{gayaMerek, gayaTagline, gayaJudul, gayaKota, gayaJudulBlok, gayaNilai, gayaNilaiKosong, gayaSatuan, gayaKeterangan, gayaChip, gayaTautan}
)

type Pembuat struct {
	Konfigurasi Konfigurasi
	tebal       *opentype.Font
	biasa       *opentype.Font
}

type pelukis struct {
	kanvas *image.RGBA
	wajah  map[gaya]font.Face
}

func Baru(k Konfigurasi) (Pembuat, error) {
	tebal, err := opentype.Parse(gobold.TTF)
	if err != nil {
		return Pembuat{}, fmt.Errorf("baca huruf tebal: %w", err)
	}
	biasa, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return Pembuat{}, fmt.Errorf("baca huruf biasa: %w", err)
	}
	return Pembuat{Konfigurasi: k, tebal: tebal, biasa: biasa}, nil
}

func (p Pembuat) buatWajah() (map[gaya]font.Face, error) {
	wajah := make(map[gaya]font.Face, len(semuaGaya))
	for _, g := range semuaGaya {
		jenis := p.biasa
		if g.isTebal {
			jenis = p.tebal
		}
		w, err := opentype.NewFace(jenis, &opentype.FaceOptions{Size: g.ukuran, DPI: dpiLayar, Hinting: font.HintingFull})
		if err != nil {
			return nil, fmt.Errorf("buat wajah gaya %v: %w", g, err)
		}
		wajah[g] = w
	}
	return wajah, nil
}

func (p Pembuat) Gambar(r domain.RingkasanKecamatan, bahasa string) ([]byte, error) {
	wajah, err := p.buatWajah()
	if err != nil {
		return nil, err
	}
	isi := susunIsi(r, p.Konfigurasi, bahasa)
	kanvas := image.NewRGBA(image.Rect(0, 0, lebarKartu, tinggiKartu))
	draw.Draw(kanvas, kanvas.Bounds(), image.NewUniform(warnaLatar), image.Point{}, draw.Src)
	g := pelukis{kanvas: kanvas, wajah: wajah}
	g.kepala(isi.Tagline)
	y := g.judul(isi.Judul, isi.Kota)
	for _, b := range isi.Blok {
		y = g.blok(b, isi.Label, y) + jarakBlok
	}
	g.kaki(isi.Ajakan, isi.Tautan)
	var keluaran bytes.Buffer
	if err := png.Encode(&keluaran, kanvas); err != nil {
		return nil, fmt.Errorf("tulis png kartu %s: %w", r.Kecamatan.Kode, err)
	}
	return keluaran.Bytes(), nil
}

func (g pelukis) kepala(tagline string) {
	g.isiBulat(image.Rect(0, 0, lebarKartu, tinggiKepala), 0, warnaPrimer)
	g.cincin(image.Rect(tepi, 62, tepi+64, 106), warnaKartu)
	g.cincin(image.Rect(tepi+36, 62, tepi+100, 106), warnaAksen)
	g.tulis("rantai", gayaMerek, image.Pt(tepi+124, 104), warnaKartu)
	g.tulis(tagline, gayaTagline, image.Pt(lebarKartu-tepi-g.lebarTeks(tagline, gayaTagline), 100), warnaSekunder)
}

func (g pelukis) judul(judul, kota string) int {
	y := tinggiKepala + 96
	for _, baris := range g.bungkus(judul, gayaJudul, lebarIsi) {
		g.tulis(baris, gayaJudul, image.Pt(tepi, y), warnaTeks)
		y += 76
	}
	g.tulis(kota, gayaKota, image.Pt(tepi, y-22), warnaRedup)
	return y + 30
}

func (g pelukis) blok(b blok, label map[domain.Label]string, y int) int {
	kotak := image.Rect(tepi, y, lebarKartu-tepi, y+tinggiBlok)
	g.isiBulat(kotak, radiusBlok, warnaGaris)
	g.isiBulat(kotak.Inset(tebalGaris), radiusBlok-tebalGaris, warnaKartu)
	x := tepi + 40
	g.tulis(b.Judul, gayaJudulBlok, image.Pt(x, y+60), warnaRedup)
	if b.Label == "" {
		g.tulis(b.Nilai, gayaNilaiKosong, image.Pt(x, y+140), warnaRedup)
		return y + tinggiBlok
	}
	lebarNilai := g.tulis(b.Nilai, gayaNilai, image.Pt(x, y+144), warnaTeks)
	g.tulis(b.Satuan, gayaSatuan, image.Pt(x+lebarNilai+16, y+144), warnaRedup)
	lebarChip := g.chip(label[b.Label], warnaLabel[b.Label], image.Pt(x, y+168))
	g.tulis(b.Keterangan, gayaKeterangan, image.Pt(x+lebarChip+16, y+195), warnaRedup)
	return y + tinggiBlok
}

func (g pelukis) chip(teks string, warna pasanganWarna, kiriAtas image.Point) int {
	lebar := g.lebarTeks(teks, gayaChip) + 32
	g.isiBulat(image.Rectangle{Min: kiriAtas, Max: kiriAtas.Add(image.Pt(lebar, 38))}, 19, warna.latar)
	g.tulis(teks, gayaChip, kiriAtas.Add(image.Pt(16, 27)), warna.teks)
	return lebar
}

func (g pelukis) kaki(ajakan, tautan string) {
	g.isiBulat(image.Rect(tepi, 1190, lebarKartu-tepi, 1192), 0, warnaGaris)
	g.tulis(ajakan, gayaJudulBlok, image.Pt(tepi, 1246), warnaRedup)
	g.tulis(tautan, gayaTautan, image.Pt(tepi, 1296), warnaPrimer)
}

func (g pelukis) cincin(kotak image.Rectangle, warna color.RGBA) {
	radius := kotak.Dy() / 2
	g.isiBulat(kotak, radius, warna)
	g.isiBulat(kotak.Inset(9), radius-9, warnaPrimer)
}

func (g pelukis) lebarTeks(teks string, h gaya) int {
	return font.MeasureString(g.wajah[h], teks).Ceil()
}

func (g pelukis) tulis(teks string, h gaya, titik image.Point, warna color.RGBA) int {
	penulis := font.Drawer{Dst: g.kanvas, Src: image.NewUniform(warna), Face: g.wajah[h], Dot: fixed.P(titik.X, titik.Y)}
	penulis.DrawString(teks)
	return (penulis.Dot.X - fixed.I(titik.X)).Ceil()
}

func (g pelukis) bungkus(teks string, h gaya, lebarMaks int) []string {
	baris := []string{}
	sekarang := ""
	for _, kata := range strings.Fields(teks) {
		calon := strings.TrimSpace(sekarang + " " + kata)
		if sekarang != "" && g.lebarTeks(calon, h) > lebarMaks {
			baris = append(baris, sekarang)
			calon = kata
		}
		sekarang = calon
	}
	return append(baris, sekarang)
}

func (g pelukis) isiBulat(kotak image.Rectangle, radius int, warna color.RGBA) {
	lebar, tinggi := float32(kotak.Dx()), float32(kotak.Dy())
	r := float32(radius)
	k := r * pembulatanKurva
	ras := vector.NewRasterizer(kotak.Dx(), kotak.Dy())
	ras.MoveTo(r, 0)
	ras.LineTo(lebar-r, 0)
	ras.CubeTo(lebar-r+k, 0, lebar, r-k, lebar, r)
	ras.LineTo(lebar, tinggi-r)
	ras.CubeTo(lebar, tinggi-r+k, lebar-r+k, tinggi, lebar-r, tinggi)
	ras.LineTo(r, tinggi)
	ras.CubeTo(r-k, tinggi, 0, tinggi-r+k, 0, tinggi-r)
	ras.LineTo(0, r)
	ras.CubeTo(0, r-k, r-k, 0, r, 0)
	ras.ClosePath()
	ras.Draw(g.kanvas, kotak, image.NewUniform(warna), image.Point{})
}
