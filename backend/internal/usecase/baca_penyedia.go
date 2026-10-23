package usecase

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type PembacaPenyedia interface {
	AmbilProfilPenyedia(ctx context.Context, id int64, batas int) (domain.ProfilPenyedia, bool, error)
}

type BacaPenyedia struct {
	Pembaca  PembacaPenyedia
	Batas    int
	Atribusi string
}

type ProfilPenyediaBeratribusi struct {
	Profil   domain.ProfilPenyedia
	Atribusi string
}

func (u BacaPenyedia) Ambil(ctx context.Context, id int64) (ProfilPenyediaBeratribusi, error) {
	p, ada, err := u.Pembaca.AmbilProfilPenyedia(ctx, id, u.Batas)
	if err != nil {
		return ProfilPenyediaBeratribusi{}, fmt.Errorf("ambil penyedia %d: %w", id, err)
	}
	if !ada {
		return ProfilPenyediaBeratribusi{}, domain.TidakDitemukan("penyedia tidak ditemukan")
	}
	return ProfilPenyediaBeratribusi{Profil: p, Atribusi: u.Atribusi}, nil
}
