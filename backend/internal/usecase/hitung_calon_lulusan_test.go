package usecase

import (
	"context"
	"strings"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pembacaSiswaTiruan struct {
	tahunAjaran string
	daftar      []domain.SiswaSDKecamatan
}

func (p pembacaSiswaTiruan) SiswaSDPerKecamatan(context.Context) (string, []domain.SiswaSDKecamatan, error) {
	return p.tahunAjaran, p.daftar, nil
}

type penyimpanCalonTiruan struct {
	data DataCalonLulusan
}

func (p *penyimpanCalonTiruan) SimpanCalonLulusan(_ context.Context, data DataCalonLulusan) (HasilImpor, error) {
	p.data = data
	return HasilImpor{}, nil
}

type pelaporCalonTiruan struct {
	tanpaData []string
}

func (p *pelaporCalonTiruan) LaporkanKecamatanTanpaCalonLulusan(kode []string) { p.tanpaData = kode }

func hitungCalonUji(pembaca PembacaSiswaSD, penyimpan *penyimpanCalonTiruan, pelapor *pelaporCalonTiruan) HitungCalonLulusan {
	nomor := 0
	return HitungCalonLulusan{
		Pembaca: pembaca, Penyimpan: penyimpan, Pencatat: &pencatatTiruan{}, Pelapor: pelapor,
		Konstanta: KonstantaPendidikan{PembagiEstimasiTotalSD: 5},
		BuatProvenance: func(c CatatanProvenance) domain.Provenance {
			nomor++
			return domain.Provenance{ID: "prv_" + string(rune('0'+nomor)), Sumber: c.Sumber, Label: c.Label, Catatan: c.Catatan}
		},
	}
}

func TestHitungCalonLulusan_MetodeEstimasi_LabelEstimasiDanRumusDiProvenance(t *testing.T) {
	penyimpan, pelapor := &penyimpanCalonTiruan{}, &pelaporCalonTiruan{}
	pembaca := pembacaSiswaTiruan{tahunAjaran: "2026/2027", daftar: []domain.SiswaSDKecamatan{
		{KodeKecamatan: "31.72.05", SiswaKelas6: 0, TotalSiswa: 600},
	}}

	if err := hitungCalonUji(pembaca, penyimpan, pelapor).Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	calon, p := penyimpan.data.Daftar[0], penyimpan.data.Provenance[0]
	if calon.Label != domain.LabelEstimasi || calon.Jumlah != 120 || calon.ProvenanceID != p.ID {
		t.Errorf("calon = %+v", calon)
	}
	if p.Label != domain.LabelEstimasi || !strings.Contains(p.Catatan, "÷ 5") {
		t.Errorf("provenance = %+v", p)
	}
}

func TestHitungCalonLulusan_DuaMetode_ProvenanceTerpisahDanKecamatanKosongDilaporkan(t *testing.T) {
	penyimpan, pelapor := &penyimpanCalonTiruan{}, &pelaporCalonTiruan{}
	pembaca := pembacaSiswaTiruan{tahunAjaran: "2026/2027", daftar: []domain.SiswaSDKecamatan{
		{KodeKecamatan: "31.72.05", TotalSiswa: 600},
		{KodeKecamatan: "31.72.06", SiswaKelas6: 1250, TotalSiswa: 7400},
		{KodeKecamatan: "31.72.07", SiswaKelas6: 900},
		{KodeKecamatan: "31.01.01"},
	}}

	if err := hitungCalonUji(pembaca, penyimpan, pelapor).Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	if len(penyimpan.data.Provenance) != 2 || len(penyimpan.data.Daftar) != 3 {
		t.Errorf("provenance %d, daftar %d", len(penyimpan.data.Provenance), len(penyimpan.data.Daftar))
	}
	if len(pelapor.tanpaData) != 1 || pelapor.tanpaData[0] != "31.01.01" {
		t.Errorf("tanpa data = %v", pelapor.tanpaData)
	}
}

func TestHitungCalonLulusan_BelumAdaRombel_Gagal(t *testing.T) {
	err := hitungCalonUji(pembacaSiswaTiruan{}, &penyimpanCalonTiruan{}, &pelaporCalonTiruan{}).Jalankan(t.Context())

	if err == nil || !strings.Contains(err.Error(), ErrBelumAdaDataRombel.Error()) {
		t.Errorf("harap ErrBelumAdaDataRombel, dapat %v", err)
	}
}
