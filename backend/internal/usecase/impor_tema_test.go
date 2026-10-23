package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pengambilTemaTiruan struct {
	data DataSumberTema
}

func (p pengambilTemaTiruan) Ambil(context.Context) (DataSumberTema, error) {
	return p.data, nil
}

type penyimpanTemaTiruan struct {
	dipanggil bool
}

func (p *penyimpanTemaTiruan) SimpanTema(context.Context, DataSumberTema) (HasilImpor, error) {
	p.dipanggil = true
	return HasilImpor{}, nil
}

func temaRusakUji(kode string) domain.Tema {
	return domain.Tema{Kode: kode, Nama: domain.Teks{ID: "Uji", EN: "Test"}}
}

func TestImporTema_TemaTidakSah_TidakDisimpanDanRunGagal(t *testing.T) {
	penyimpan, pencatat := &penyimpanTemaTiruan{}, &pencatatTiruan{}
	impor := ImporTema{
		Pengambil: pengambilTemaTiruan{data: DataSumberTema{Tema: []domain.Tema{temaRusakUji("uji")}}},
		Penyimpan: penyimpan,
		Pencatat:  pencatat,
	}

	err := impor.Jalankan(t.Context())

	if !errors.Is(err, domain.ErrPenghitungTidakTerdaftar) || penyimpan.dipanggil || pencatat.gagal == nil {
		t.Errorf("err %v, disimpan %v, gagal tercatat %v", err, penyimpan.dipanggil, pencatat.gagal)
	}
}

func TestImporTema_KodeGanda_Ditolak(t *testing.T) {
	impor := ImporTema{
		Pengambil: pengambilTemaTiruan{data: DataSumberTema{Tema: []domain.Tema{temaRusakUji("a"), temaRusakUji("a")}}},
		Penyimpan: &penyimpanTemaTiruan{},
		Pencatat:  &pencatatTiruan{},
	}

	if err := impor.Jalankan(t.Context()); !errors.Is(err, ErrKodeTemaGanda) {
		t.Errorf("harap ErrKodeTemaGanda, dapat %v", err)
	}
}
