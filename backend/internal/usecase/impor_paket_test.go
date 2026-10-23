package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type sumberPaketTiruan struct {
	errDetail error
}

func (s sumberPaketTiruan) AmbilDaftarPaket(_ context.Context, tahun int) (DaftarPaketSumber, error) {
	return DaftarPaketSumber{Paket: []domain.PaketDaftar{{Kode: "1", TahunAnggaran: tahun}}}, nil
}

func (s sumberPaketTiruan) AmbilDetailPaket(_ context.Context, kode string) (DetailPaketSumber, error) {
	return DetailPaketSumber{Detail: domain.DetailPaket{Kode: kode, SatuanKerja: "SUKU DINAS PENDIDIKAN JAKARTA TIMUR"}}, s.errDetail
}

type penyimpanPaketTiruan struct {
	tahunTersimpan []int
	detail         []DataSimpanDetailPaket
}

func (p *penyimpanPaketTiruan) SimpanDaftarPaket(_ context.Context, data DaftarPaketSumber) (HasilImpor, error) {
	p.tahunTersimpan = append(p.tahunTersimpan, data.Paket[0].TahunAnggaran)
	return HasilImpor{JumlahBaru: 1}, nil
}

func (p *penyimpanPaketTiruan) DaftarPaketPerluDetail(context.Context, []string) ([]string, error) {
	return []string{"10", "11"}, nil
}

func (p *penyimpanPaketTiruan) SimpanDetailPaket(_ context.Context, data DataSimpanDetailPaket) error {
	p.detail = append(p.detail, data)
	return nil
}

func TestImporDaftarPaket_DuaTahun_HasilDijumlahkan(t *testing.T) {
	penyimpan := &penyimpanPaketTiruan{}
	pencatat := &pencatatTiruan{}
	u := ImporDaftarPaket{Pengambil: sumberPaketTiruan{}, Penyimpan: penyimpan, Pencatat: pencatat, Tahun: []int{2025, 2026}}

	if err := u.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	if len(penyimpan.tahunTersimpan) != 2 || pencatat.selesai == nil || pencatat.selesai.JumlahBaru != 2 {
		t.Errorf("tahun = %v, selesai = %+v", penyimpan.tahunTersimpan, pencatat.selesai)
	}
}

func TestImporDetailPaket_SatkerSudin_InstansiDiklasifikasi(t *testing.T) {
	penyimpan := &penyimpanPaketTiruan{}
	u := ImporDetailPaket{
		Pengambil: sumberPaketTiruan{},
		Penyimpan: penyimpan,
		Pencatat:  &pencatatTiruan{},
		AturanInstansi: domain.AturanInstansi{
			Jenis: []domain.AturanJenisInstansi{{Jenis: domain.JenisInstansiSudin, Awalan: []string{"suku dinas"}}},
			Kota:  map[string][]string{"31.75": {"jakarta timur"}},
		},
	}

	if err := u.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	instansi := penyimpan.detail[0].Instansi
	if len(penyimpan.detail) != 2 || instansi.Jenis != domain.JenisInstansiSudin || instansi.KodeKota != "31.75" {
		t.Errorf("detail = %d, instansi = %+v", len(penyimpan.detail), instansi)
	}
}

func TestImporDetailPaket_SumberMenolak_RunGagalDanBerhenti(t *testing.T) {
	ditolak := errors.New("403")
	penyimpan := &penyimpanPaketTiruan{}
	pencatat := &pencatatTiruan{}
	u := ImporDetailPaket{Pengambil: sumberPaketTiruan{errDetail: ditolak}, Penyimpan: penyimpan, Pencatat: pencatat}

	err := u.Jalankan(t.Context())

	if !errors.Is(err, ditolak) || !errors.Is(pencatat.gagal, ditolak) || len(penyimpan.detail) != 0 {
		t.Errorf("err = %v, gagal = %v, tersimpan = %d", err, pencatat.gagal, len(penyimpan.detail))
	}
}
