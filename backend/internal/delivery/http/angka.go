package http

import (
	"encoding/json"
	"errors"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var ErrAngkaTanpaRujukan = errors.New("angka wajib membawa label kepastian dan provenance_id")

type Angka struct {
	Nilai        float64
	Satuan       string
	Label        domain.Label
	ProvenanceID string
}

type angkaJSON struct {
	Nilai        float64      `json:"nilai"`
	Satuan       string       `json:"satuan"`
	Label        domain.Label `json:"label"`
	ProvenanceID string       `json:"provenance_id"`
}

func (a Angka) MarshalJSON() ([]byte, error) {
	if !a.Label.IsValid() || a.ProvenanceID == "" {
		return nil, ErrAngkaTanpaRujukan
	}
	return json.Marshal(angkaJSON(a))
}
