package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pembacaTemaTiruan struct {
	daftar []domain.TemaTersimpan
}

func (p pembacaTemaTiruan) DaftarTemaAktif(context.Context) ([]domain.TemaTersimpan, error) {
	return p.daftar, nil
}

type penyimpanIndikatorTiruan struct {
	data []DataIndikator
}

func (p *penyimpanIndikatorTiruan) SimpanIndikator(_ context.Context, data DataIndikator) (HasilImpor, error) {
	p.data = append(p.data, data)
	return HasilImpor{JumlahBaru: len(data.Daftar)}, nil
}

type pembacaMasukanTiruan struct {
	periode string
	daftar  []domain.MasukanKesenjangan
}

func (p pembacaMasukanTiruan) MasukanKesenjangan(context.Context, domain.Peralihan) (string, []domain.MasukanKesenjangan, error) {
	return p.periode, p.daftar, nil
}

func temaPendidikanUji() domain.TemaTersimpan {
	tema := domain.Tema{Kode: "pendidikan"}
	tema.Indikator = domain.Indikator{
		Penghitung:  domain.PenghitungKesenjanganKursiSMP,
		AmbangRasio: domain.AmbangRasio{SangatKurang: 0.5, Kurang: 0.8, HampirCukup: 1.0},
	}
	return domain.TemaTersimpan{Tema: tema}
}

func TestHitungIndikator_PenghitungBelumTersedia_Gagal(t *testing.T) {
	u := HitungIndikator{
		Tema:       pembacaTemaTiruan{daftar: []domain.TemaTersimpan{temaPendidikanUji()}},
		Penghitung: map[string]PenghitungIndikator{},
		Penyimpan:  &penyimpanIndikatorTiruan{},
		Pencatat:   &pencatatTiruan{},
	}

	if err := u.Jalankan(t.Context()); !errors.Is(err, ErrPenghitungBelumTersedia) {
		t.Errorf("harap ErrPenghitungBelumTersedia, dapat %v", err)
	}
}

func TestHitungIndikator_KesenjanganSMP_ProvenancePerLabelDanBatasDariKonfigurasi(t *testing.T) {
	penyimpan := &penyimpanIndikatorTiruan{}
	pembaca := pembacaMasukanTiruan{periode: "2026/2027", daftar: []domain.MasukanKesenjangan{
		{KodeKecamatan: "31.72.05", CalonLulusan: 100, LabelCalonLulusan: domain.LabelEstimasi, KursiNegeri: map[domain.BentukSekolah]int{domain.BentukSMP: 40}},
		{KodeKecamatan: "31.72.06", CalonLulusan: 100, LabelCalonLulusan: domain.LabelResmi, KursiNegeri: map[domain.BentukSekolah]int{domain.BentukSMP: 90},
			JumlahMadrasah: 1, JumlahMadrasahDenganRombel: 1, RombelTingkatPertamaMadrasah: 2},
	}}
	peralihan, _ := domain.PeralihanTema(domain.TemaPendidikan)
	u := HitungIndikator{
		Tema: pembacaTemaTiruan{daftar: []domain.TemaTersimpan{temaPendidikanUji()}},
		Penghitung: map[string]PenghitungIndikator{domain.PenghitungKesenjanganKursiSMP: KesenjanganKursi{
			Peralihan: peralihan, Pembaca: pembaca, Konstanta: konstantaUji,
			BuatProvenance: func(c CatatanProvenance) domain.Provenance {
				return domain.Provenance{ID: "prv_" + string(c.Label), Label: c.Label, Catatan: c.Catatan}
			},
		}},
		Penyimpan: penyimpan,
		Pencatat:  &pencatatTiruan{},
	}

	if err := u.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	data := penyimpan.data[0]
	if data.TemaKode != "pendidikan" || len(data.Provenance) != 2 || data.Daftar[0].ProvenanceID != "prv_estimasi" {
		t.Errorf("data = %+v", data)
	}
	if kursi := data.Daftar[1].Rincian.KursiMadrasah.Nilai; kursi == nil || *kursi != 60 {
		t.Errorf("kursi madrasah = %v", kursi)
	}
}

func TestKesenjanganKursi_SMASMK_RumusMenyebutKeduaBentukDanBatas36(t *testing.T) {
	peralihan, _ := domain.PeralihanTema(domain.TemaPendidikanSMASMK)
	pembaca := pembacaMasukanTiruan{periode: "2026/2027", daftar: []domain.MasukanKesenjangan{
		{KodeKecamatan: "31.72.05", CalonLulusan: 1000, LabelCalonLulusan: domain.LabelResmi,
			KursiNegeri: map[domain.BentukSekolah]int{domain.BentukSMA: 300, domain.BentukSMK: 200}},
	}}
	tema := domain.Tema{Kode: domain.TemaPendidikanSMASMK}
	tema.Indikator.AmbangRasio = domain.AmbangRasio{SangatKurang: 0.5, Kurang: 0.8, HampirCukup: 1.0}
	penghitung := KesenjanganKursi{
		Peralihan: peralihan, Pembaca: pembaca, Konstanta: konstantaUji,
		BuatProvenance: func(c CatatanProvenance) domain.Provenance { return domain.Provenance{ID: "prv_1", Catatan: c.Catatan} },
	}

	data, err := penghitung.Hitung(t.Context(), tema)

	if err != nil {
		t.Fatal(err)
	}
	indikator := data.Daftar[0]
	if indikator.Kapasitas != 500 || indikator.Kelas != domain.KelasKurang || indikator.Rincian.BatasMuridPerRombel != 36 {
		t.Errorf("indikator = %+v", indikator.IndikatorKesenjangan)
	}
	catatan := data.Provenance[0].Catatan
	if !strings.Contains(catatan, "(kursi SMA negeri + kursi SMK negeri)") || !strings.Contains(catatan, "kelas 10 MA × 36") {
		t.Errorf("catatan = %q", catatan)
	}
}
