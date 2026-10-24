package foto

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"golang.org/x/image/draw"
)

const (
	kualitasJPEG       = 85
	batasPikselDekode  = 50_000_000
	pesanGagalDekode   = "foto tidak dapat dibaca"
	pesanFotoKebesaran = "resolusi foto terlalu besar"
)

var (
	tandaJPEG = []byte{0xFF, 0xD8, 0xFF}
	tandaPNG  = []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1A, '\n'}
)

type Pemroses struct{}

func (Pemroses) Proses(isi []byte) (domain.FotoDiproses, error) {
	gambar, err := dekode(isi)
	if err != nil {
		return domain.FotoDiproses{}, err
	}
	utama := perkecil(gambar, domain.SisiTerpanjangFotoPiksel)
	isiUtama, err := enkodeJPEG(utama)
	if err != nil {
		return domain.FotoDiproses{}, err
	}
	isiThumbnail, err := enkodeJPEG(perkecil(gambar, domain.SisiTerpanjangThumbnailPiksel))
	if err != nil {
		return domain.FotoDiproses{}, err
	}
	return domain.FotoDiproses{Utama: isiUtama, Thumbnail: isiThumbnail, Lebar: utama.Bounds().Dx(), Tinggi: utama.Bounds().Dy()}, nil
}

func dekode(isi []byte) (image.Image, error) {
	pembaca, err := pembacaFormat(isi)
	if err != nil {
		return nil, err
	}
	cfg, err := pembaca.konfigurasi(bytes.NewReader(isi))
	if err != nil {
		return nil, domain.InputTidakValid(pesanGagalDekode)
	}
	if cfg.Width*cfg.Height > batasPikselDekode {
		return nil, domain.InputTidakValid(pesanFotoKebesaran)
	}
	gambar, err := pembaca.dekode(bytes.NewReader(isi))
	if err != nil {
		return nil, domain.InputTidakValid(pesanGagalDekode)
	}
	return gambar, nil
}

type pembacaGambar struct {
	konfigurasi func(r *bytes.Reader) (image.Config, error)
	dekode      func(r *bytes.Reader) (image.Image, error)
}

func pembacaFormat(isi []byte) (pembacaGambar, error) {
	switch {
	case bytes.HasPrefix(isi, tandaJPEG):
		return pembacaGambar{
			konfigurasi: func(r *bytes.Reader) (image.Config, error) { return jpeg.DecodeConfig(r) },
			dekode:      func(r *bytes.Reader) (image.Image, error) { return jpeg.Decode(r) },
		}, nil
	case bytes.HasPrefix(isi, tandaPNG):
		return pembacaGambar{
			konfigurasi: func(r *bytes.Reader) (image.Config, error) { return png.DecodeConfig(r) },
			dekode:      func(r *bytes.Reader) (image.Image, error) { return png.Decode(r) },
		}, nil
	default:
		return pembacaGambar{}, domain.ErrFormatFotoTidakSah
	}
}

func perkecil(gambar image.Image, sisiTerpanjang int) image.Image {
	lebar, tinggi := gambar.Bounds().Dx(), gambar.Bounds().Dy()
	if max(lebar, tinggi) <= sisiTerpanjang {
		hasil := kanvasPutih(lebar, tinggi)
		draw.Draw(hasil, hasil.Bounds(), gambar, gambar.Bounds().Min, draw.Over)
		return hasil
	}
	if lebar >= tinggi {
		lebar, tinggi = sisiTerpanjang, tinggi*sisiTerpanjang/lebar
	} else {
		lebar, tinggi = lebar*sisiTerpanjang/tinggi, sisiTerpanjang
	}
	hasil := kanvasPutih(max(lebar, 1), max(tinggi, 1))
	draw.CatmullRom.Scale(hasil, hasil.Bounds(), gambar, gambar.Bounds(), draw.Over, nil)
	return hasil
}

func kanvasPutih(lebar, tinggi int) *image.RGBA {
	hasil := image.NewRGBA(image.Rect(0, 0, lebar, tinggi))
	draw.Draw(hasil, hasil.Bounds(), image.White, image.Point{}, draw.Src)
	return hasil
}

func enkodeJPEG(gambar image.Image) ([]byte, error) {
	var keluaran bytes.Buffer
	if err := jpeg.Encode(&keluaran, gambar, &jpeg.Options{Quality: kualitasJPEG}); err != nil {
		return nil, fmt.Errorf("enkode foto: %w", err)
	}
	return keluaran.Bytes(), nil
}
