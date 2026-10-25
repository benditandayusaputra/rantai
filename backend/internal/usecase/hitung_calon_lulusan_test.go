package usecase

import (
	"context"
	"strings"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pembacaSiswaTiruan struct {
	tahunAjaran string
	daftar      []domain.SiswaKecamatan
	diminta     []domain.BentukSekolah
}

func (p *pembacaSiswaTiruan) SiswaPerKecamatan(_ context.Context, peralihan domain.Peralihan) (string, []domain.SiswaKecamatan, error) {
	p.diminta = append(p.diminta, peralihan.BentukAsal)
	return p.tahunAjaran, p.daftar, nil
}

type penyimpanCalonTiruan struct {
	semua []DataCalonLulusan
	data  DataCalonLulusan
}

func (p *penyimpanCalonTiruan) SimpanCalonLulusan(_ context.Context, data DataCalonLulusan) (HasilImpor, error) {
	p.semua = append(p.semua, data)
	if len(p.semua) == 1 {
		p.data = data
	}
	return HasilImpor{}, nil
}

type pelaporCalonTiruan struct {
	tanpaData []string
}

func (p *pelaporCalonTiruan) LaporkanKecamatanTanpaCalonLulusan(jenjang domain.Jenjang, kode []string) {
	if jenjang == domain.JenjangDasar {
		p.tanpaData = kode
	}
}

var konstantaUji = KonstantaPendidikan{Peralihan: map[string]KonstantaPeralihan{
	domain.PeralihanSDKeSMP:     {BatasMuridPerRombel: 30, PembagiEstimasi: 5},
	domain.PeralihanSMPKeSMASMK: {BatasMuridPerRombel: 36, PembagiEstimasi: 3},
}}

func hitungCalonUji(pembaca PembacaSiswa, penyimpan *penyimpanCalonTiruan, pelapor *pelaporCalonTiruan) HitungCalonLulusan {
	nomor := 0
	return HitungCalonLulusan{
		Pembaca: pembaca, Penyimpan: penyimpan, Pencatat: &pencatatTiruan{}, Pelapor: pelapor,
		Konstanta: konstantaUji,
		BuatProvenance: func(c CatatanProvenance) domain.Provenance {
			nomor++
			return domain.Provenance{ID: "prv_" + string(rune('0'+nomor)), Sumber: c.Sumber, Label: c.Label, Catatan: c.Catatan}
		},
	}
}

func TestHitungCalonLulusan_MetodeEstimasi_LabelEstimasiDanRumusDiProvenance(t *testing.T) {
	penyimpan, pelapor := &penyimpanCalonTiruan{}, &pelaporCalonTiruan{}
	pembaca := &pembacaSiswaTiruan{tahunAjaran: "2026/2027", daftar: []domain.SiswaKecamatan{
		{KodeKecamatan: "31.72.05", SiswaTingkatAkhir: 0, TotalSiswa: 600},
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
	pembaca := &pembacaSiswaTiruan{tahunAjaran: "2026/2027", daftar: []domain.SiswaKecamatan{
		{KodeKecamatan: "31.72.05", TotalSiswa: 600},
		{KodeKecamatan: "31.72.06", SiswaTingkatAkhir: 1250, TotalSiswa: 7400},
		{KodeKecamatan: "31.72.07", SiswaTingkatAkhir: 900},
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
	err := hitungCalonUji(&pembacaSiswaTiruan{}, &penyimpanCalonTiruan{}, &pelaporCalonTiruan{}).Jalankan(t.Context())

	if err == nil || !strings.Contains(err.Error(), ErrBelumAdaDataRombel.Error()) {
		t.Errorf("harap ErrBelumAdaDataRombel, dapat %v", err)
	}
}

func TestHitungCalonLulusan_DuaPeralihan_LulusanSDDanSMPDisimpanPerJenjang(t *testing.T) {
	penyimpan, pelapor := &penyimpanCalonTiruan{}, &pelaporCalonTiruan{}
	pembaca := &pembacaSiswaTiruan{tahunAjaran: "2026/2027", daftar: []domain.SiswaKecamatan{
		{KodeKecamatan: "31.72.05", TotalSiswa: 900},
	}}

	if err := hitungCalonUji(pembaca, penyimpan, pelapor).Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	if len(penyimpan.semua) != 2 || pembaca.diminta[0] != domain.BentukSD || pembaca.diminta[1] != domain.BentukSMP {
		t.Fatalf("diminta %v, disimpan %d", pembaca.diminta, len(penyimpan.semua))
	}
	smp := penyimpan.semua[1]
	calon := smp.Daftar[0]
	if smp.Jenjang != domain.JenjangMenengahPertama || calon.Jumlah != 300 || calon.Metode != domain.MetodeEstimasiTotalSMP {
		t.Errorf("lulusan SMP = %+v, %+v", smp.Jenjang, calon)
	}
	if !strings.Contains(smp.Provenance[0].Catatan, "siswa SMP") || !strings.Contains(smp.Provenance[0].Catatan, "÷ 3") {
		t.Errorf("catatan = %q", smp.Provenance[0].Catatan)
	}
}
