package foto

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var penandaEXIF = []byte("Exif\x00\x00MM\x00\x2aGPS-RAHASIA-106.9012")

func gambarUji(lebar, tinggi int) image.Image {
	gambar := image.NewRGBA(image.Rect(0, 0, lebar, tinggi))
	for x := range lebar {
		gambar.Set(x, tinggi/2, color.RGBA{R: 200, A: 255})
	}
	return gambar
}

func jpegDenganEXIF(t *testing.T, lebar, tinggi int) []byte {
	t.Helper()
	var polos bytes.Buffer
	if err := jpeg.Encode(&polos, gambarUji(lebar, tinggi), nil); err != nil {
		t.Fatal(err)
	}
	panjang := len(penandaEXIF) + 2
	segmenEXIF := append([]byte{0xFF, 0xE1, byte(panjang >> 8), byte(panjang)}, penandaEXIF...)
	isi := polos.Bytes()
	return append(append(append([]byte{}, isi[:2]...), segmenEXIF...), isi[2:]...)
}

func TestProses_JPEGBerEXIF_HasilTanpaEXIF(t *testing.T) {
	masukan := jpegDenganEXIF(t, 64, 48)
	if !bytes.Contains(masukan, []byte("GPS-RAHASIA")) {
		t.Fatal("fixture harus memuat EXIF")
	}

	hasil, err := Pemroses{}.Proses(masukan)

	if err != nil {
		t.Fatal(err)
	}
	for nama, isi := range map[string][]byte{"utama": hasil.Utama, "thumbnail": hasil.Thumbnail} {
		if bytes.Contains(isi, []byte("Exif")) || bytes.Contains(isi, []byte{0xFF, 0xE1}) || bytes.Contains(isi, []byte("GPS-RAHASIA")) {
			t.Errorf("foto %s masih memuat EXIF", nama)
		}
	}
}

func TestProses_FotoBesar_DiperkecilSisiTerpanjang(t *testing.T) {
	var png3000 bytes.Buffer
	if err := png.Encode(&png3000, gambarUji(3000, 1000)); err != nil {
		t.Fatal(err)
	}

	hasil, err := Pemroses{}.Proses(png3000.Bytes())

	if err != nil {
		t.Fatal(err)
	}
	thumbnail, err := jpeg.DecodeConfig(bytes.NewReader(hasil.Thumbnail))
	if err != nil || hasil.Lebar != domain.SisiTerpanjangFotoPiksel || hasil.Tinggi != 682 || thumbnail.Width != domain.SisiTerpanjangThumbnailPiksel {
		t.Errorf("utama %dx%d, thumbnail %+v, err %v", hasil.Lebar, hasil.Tinggi, thumbnail, err)
	}
}

func TestProses_FormatBukanJPEGAtauPNG_Ditolak(t *testing.T) {
	var isiGIF bytes.Buffer
	if err := gif.Encode(&isiGIF, gambarUji(8, 8), nil); err != nil {
		t.Fatal(err)
	}

	if _, err := (Pemroses{}).Proses(isiGIF.Bytes()); !errors.Is(err, domain.ErrFormatFotoTidakSah) {
		t.Errorf("err = %v", err)
	}
}

func TestProses_TandaJPEGIsiRusak_InputTidakValid(t *testing.T) {
	_, err := Pemroses{}.Proses(append(append([]byte{}, tandaJPEG...), 0x00, 0x01, 0x02))

	var kesalahan *domain.Kesalahan
	if !errors.As(err, &kesalahan) || kesalahan.Kategori != domain.KategoriInputTidakValid {
		t.Errorf("err = %v", err)
	}
}
