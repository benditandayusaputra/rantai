package tanya

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	TipeTeks     = "string"
	TipeBilangan = "integer"
)

var (
	ErrAlatTidakDikenal = errors.New("alat tidak dikenal")
	ErrAlatGanda        = errors.New("nama alat terdaftar dua kali")
)

type Parameter struct {
	Tipe      string `json:"type"`
	Deskripsi string `json:"description"`
}

type Spek struct {
	Nama      string
	Deskripsi string
	Parameter map[string]Parameter
	Wajib     []string
}

type Alat struct {
	Spek
	jalankan func(ctx context.Context, masukan json.RawMessage) (any, error)
}

func BuatAlat[T any](spek Spek, jalankan func(ctx context.Context, masukan T) (any, error)) Alat {
	return Alat{Spek: spek, jalankan: func(ctx context.Context, mentah json.RawMessage) (any, error) {
		var masukan T
		if err := json.Unmarshal(mentah, &masukan); err != nil {
			return nil, domain.InputTidakValid("masukan alat " + spek.Nama + " tidak sesuai skema")
		}
		return jalankan(ctx, masukan)
	}}
}

type Registry struct {
	alat   map[string]Alat
	urutan []Spek
}

func BuatRegistry(daftar ...Alat) (Registry, error) {
	r := Registry{alat: make(map[string]Alat, len(daftar)), urutan: make([]Spek, 0, len(daftar))}
	for _, a := range daftar {
		if _, ada := r.alat[a.Nama]; ada {
			return Registry{}, fmt.Errorf("%w: %s", ErrAlatGanda, a.Nama)
		}
		r.alat[a.Nama] = a
		r.urutan = append(r.urutan, a.Spek)
	}
	return r, nil
}

func (r Registry) Daftar() []Spek {
	return r.urutan
}

func (r Registry) Jalankan(ctx context.Context, nama string, masukan json.RawMessage) (json.RawMessage, error) {
	a, ada := r.alat[nama]
	if !ada {
		return nil, fmt.Errorf("%w: %s", ErrAlatTidakDikenal, nama)
	}
	hasil, err := a.jalankan(ctx, masukan)
	if err != nil {
		return nil, fmt.Errorf("jalankan alat %s: %w", nama, err)
	}
	isi, err := json.Marshal(hasil)
	if err != nil {
		return nil, fmt.Errorf("serialisasi hasil alat %s: %w", nama, err)
	}
	return isi, nil
}
