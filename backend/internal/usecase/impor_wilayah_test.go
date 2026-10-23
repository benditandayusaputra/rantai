package usecase

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pengambilTiruan struct {
	data DataSumberWilayah
	err  error
}

func (p pengambilTiruan) Ambil(context.Context) (DataSumberWilayah, error) {
	return p.data, p.err
}

type penyimpanTiruan struct {
	diterima DataSimpanWilayah
}

func (p *penyimpanTiruan) SimpanWilayah(_ context.Context, data DataSimpanWilayah) (HasilImpor, error) {
	p.diterima = data
	return HasilImpor{JumlahBaru: len(data.Wilayah)}, nil
}

type pencatatTiruan struct {
	selesai *HasilImpor
	gagal   error
}

func (p *pencatatTiruan) Mulai(context.Context, string) (int64, error) { return 1, nil }

func (p *pencatatTiruan) Selesai(_ context.Context, _ int64, hasil HasilImpor) error {
	p.selesai = &hasil
	return nil
}

func (p *pencatatTiruan) Gagal(_ context.Context, _ int64, sebab error) error {
	p.gagal = sebab
	return nil
}

func imporUji(pengambil PengambilWilayah, penyimpan *penyimpanTiruan, pencatat *pencatatTiruan) ImporWilayah {
	return ImporWilayah{
		Pengambil: pengambil,
		Penyimpan: penyimpan,
		Pencatat:  pencatat,
		AturanAlias: domain.AturanAlias{PrefiksPerTingkat: map[domain.Tingkat]map[string][]string{
			domain.TingkatKelurahan: {"": {"Kel. "}},
			domain.TingkatProvinsi:  {"": {"Provinsi "}},
		}},
		Provinsi: IdentitasProvinsi{Kode: "31", Nama: "DKI Jakarta"},
	}
}

func TestImporWilayah_SumberBerhasil_AliasDilengkapiDanRunSelesai(t *testing.T) {
	sumber := pengambilTiruan{data: DataSumberWilayah{Wilayah: []domain.Wilayah{
		{Kode: "31.72.01.1005", Nama: "Pluit", Tingkat: domain.TingkatKelurahan},
	}}}
	penyimpan, pencatat := &penyimpanTiruan{}, &pencatatTiruan{}

	err := imporUji(sumber, penyimpan, pencatat).Jalankan(t.Context())

	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(penyimpan.diterima.Wilayah[0].Alias, []string{"Kel. Pluit"}) {
		t.Errorf("alias kelurahan = %q", penyimpan.diterima.Wilayah[0].Alias)
	}
	if penyimpan.diterima.Provinsi.Kode != "31" || !slices.Equal(penyimpan.diterima.Provinsi.Alias, []string{"Provinsi DKI Jakarta"}) {
		t.Errorf("provinsi = %+v", penyimpan.diterima.Provinsi)
	}
	if pencatat.selesai == nil || pencatat.selesai.JumlahBaru != 1 {
		t.Errorf("run tidak dicatat selesai: %+v", pencatat.selesai)
	}
}

func TestImporWilayah_SumberGagal_RunDicatatGagal(t *testing.T) {
	sebab := errors.New("sumber tidak terjangkau")
	penyimpan, pencatat := &penyimpanTiruan{}, &pencatatTiruan{}

	err := imporUji(pengambilTiruan{err: sebab}, penyimpan, pencatat).Jalankan(t.Context())

	if !errors.Is(err, sebab) || !errors.Is(pencatat.gagal, sebab) {
		t.Errorf("err = %v, gagal tercatat = %v", err, pencatat.gagal)
	}
	if pencatat.selesai != nil {
		t.Error("run gagal tidak boleh dicatat selesai")
	}
}
