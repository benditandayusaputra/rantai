package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type penyimpanPeranPalsu struct {
	email string
	peran domain.Peran
}

func (p *penyimpanPeranPalsu) AturPeran(_ context.Context, email string, peran domain.Peran) error {
	p.email, p.peran = email, peran
	return nil
}

func TestBeriPeran(t *testing.T) {
	kasus := []struct {
		nama  string
		peran domain.Peran
		harap error
	}{
		{"moderator", domain.PeranModerator, nil},
		{"admin", domain.PeranAdmin, nil},
		{"instansi lewat admin, bukan perintah", domain.PeranInstansi, domain.ErrPeranBukanUntukPerintah},
		{"peran asing", "pejabat", domain.ErrPeranBukanUntukPerintah},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			penyimpan := &penyimpanPeranPalsu{}

			err := BeriPeran{Penyimpan: penyimpan}.Jalankan(t.Context(), " Moderator@Contoh.ID", k.peran)

			if !errors.Is(err, k.harap) || (err == nil && penyimpan.email != "moderator@contoh.id") {
				t.Errorf("err %v, tersimpan %+v", err, penyimpan)
			}
		})
	}
}
