package djpk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var (
	ErrFormatBerubah = errors.New("format ekspor APBD DJPK berubah: kolom Akun, Anggaran, Realisasi tidak ditemukan")
	kolomWajib       = []string{"Akun", "Anggaran", "Realisasi"}
)

type bukuKerja struct {
	Baris []barisXML `xml:"Worksheet>Table>Row"`
}

type barisXML struct {
	Sel []string `xml:"Cell>Data"`
}

func UraiAPBD(isi []byte) ([]domain.AkunAPBD, error) {
	var buku bukuKerja
	if err := xml.Unmarshal(isi, &buku); err != nil {
		return nil, fmt.Errorf("urai ekspor apbd: %w", err)
	}
	if len(buku.Baris) == 0 || len(buku.Baris[0].Sel) < len(kolomWajib) || !slices.Equal(buku.Baris[0].Sel[:len(kolomWajib)], kolomWajib) {
		return nil, ErrFormatBerubah
	}
	akun := make([]domain.AkunAPBD, 0, len(buku.Baris)-1)
	for _, baris := range buku.Baris[1:] {
		a, err := uraiBaris(baris)
		if err != nil {
			return nil, err
		}
		akun = append(akun, a)
	}
	return domain.RapikanAkunAPBD(akun), nil
}

func uraiBaris(baris barisXML) (domain.AkunAPBD, error) {
	if len(baris.Sel) < len(kolomWajib) {
		return domain.AkunAPBD{}, ErrFormatBerubah
	}
	anggaran, errAnggaran := strconv.ParseFloat(baris.Sel[1], 64)
	realisasi, errRealisasi := strconv.ParseFloat(baris.Sel[2], 64)
	if err := errors.Join(errAnggaran, errRealisasi); err != nil {
		return domain.AkunAPBD{}, fmt.Errorf("angka akun %q: %w", baris.Sel[0], err)
	}
	return domain.AkunAPBD{Nama: strings.TrimSpace(baris.Sel[0]), Anggaran: anggaran, Realisasi: realisasi}, nil
}
