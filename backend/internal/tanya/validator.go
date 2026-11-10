package tanya

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var (
	ErrJawabanKosong            = errors.New("jawaban kosong")
	ErrJawabanTerlarang         = errors.New("jawaban memuat kosakata terlarang")
	ErrJawabanRujukanTakDikenal = errors.New("jawaban merujuk provenance yang tidak ada di hasil alat")
	ErrJawabanAngkaKarangan     = errors.New("jawaban memuat angka yang tidak ada di hasil alat")
	ErrJawabanMenyebutTag       = errors.New("jawaban menyebut nama tag masukan")
)

var polaNamaTagMasukan = regexp.MustCompile(`bahasa_antarmuka|kecamatan_aktif`)

type ValidatorJawaban struct {
	Terlarang domain.KosakataTerlarang
}

func (v ValidatorJawaban) Periksa(teks string, hasilAlat []json.RawMessage) error {
	potongan := make([]string, len(hasilAlat))
	for i, h := range hasilAlat {
		potongan[i] = string(h)
	}
	sumber := strings.Join(potongan, " ")
	tanpaRujukan := polaRujukanProvenance.ReplaceAllString(teks, "")
	switch {
	case strings.TrimSpace(tanpaRujukan) == "":
		return ErrJawabanKosong
	case len(v.Terlarang.Temukan(tanpaRujukan)) > 0:
		return ErrJawabanTerlarang
	case polaNamaTagMasukan.MatchString(tanpaRujukan):
		return ErrJawabanMenyebutTag
	case !isRujukanDikenal(teks, sumber):
		return ErrJawabanRujukanTakDikenal
	case !domain.IsAngkaBersumber(tanpaRujukan, sumber):
		return ErrJawabanAngkaKarangan
	default:
		return nil
	}
}

func isRujukanDikenal(teks, sumber string) bool {
	for _, cocok := range polaRujukanProvenance.FindAllStringSubmatch(teks, -1) {
		if !strings.Contains(sumber, `"`+cocok[1]+`"`) {
			return false
		}
	}
	return true
}
