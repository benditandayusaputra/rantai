package domain_test

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

func idPenyedia(id int64) *int64 { return &id }

func TestSusunJaringan_DuaPaketSatuInstansiSatuPenyedia_SimpulDigabung(t *testing.T) {
	relasi := []domain.RelasiPaket{
		{PaketKode: "A", PaketNama: "Rehab SD", Nilai: 300, ProvenanceID: "p1", InstansiID: 7, InstansiNama: "Sudin Pendidikan", PenyediaID: idPenyedia(9), PenyediaNama: "CV Maju"},
		{PaketKode: "B", PaketNama: "Rehab SMP", Nilai: 200, ProvenanceID: "p2", InstansiID: 7, InstansiNama: "Sudin Pendidikan", PenyediaID: idPenyedia(9), PenyediaNama: "CV Maju"},
	}
	j := domain.SusunJaringan(relasi)
	if len(j.Simpul) != 4 || len(j.Sisi) != 4 {
		t.Fatalf("simpul %d sisi %d, ingin 4 dan 4", len(j.Simpul), len(j.Sisi))
	}
	nilai := map[string]domain.SimpulJaringan{}
	for _, s := range j.Simpul {
		nilai[s.ID] = s
	}
	if s := nilai["instansi:7"]; s.Nilai != 500 || s.JumlahPaket != 2 || s.ProvenanceID != "p1" {
		t.Errorf("instansi %+v", s)
	}
	if s := nilai["penyedia:9"]; s.Nilai != 500 || s.JumlahPaket != 2 {
		t.Errorf("penyedia %+v", s)
	}
}

func TestSusunJaringan_PaketTanpaPenyedia_HanyaSisiInstansi(t *testing.T) {
	j := domain.SusunJaringan([]domain.RelasiPaket{{PaketKode: "A", InstansiID: 1, ProvenanceID: "p"}})
	if len(j.Simpul) != 2 || len(j.Sisi) != 1 || j.Sisi[0] != (domain.SisiJaringan{Dari: "instansi:1", Ke: "paket:A"}) {
		t.Fatalf("hasil %+v", j)
	}
}

func TestSusunJaringan_Kosong_KoleksiKosong(t *testing.T) {
	j := domain.SusunJaringan(nil)
	if j.Simpul == nil || j.Sisi == nil {
		t.Fatal("koleksi harus kosong, bukan nil")
	}
}

func TestFilterJaringan_Validasi(t *testing.T) {
	kasus := []struct {
		nama   string
		filter domain.FilterJaringan
		isSah  bool
	}{
		{"wilayah saja", domain.FilterJaringan{KodeWilayah: "31.72"}, true},
		{"penyedia saja", domain.FilterJaringan{PenyediaID: 3}, true},
		{"keduanya", domain.FilterJaringan{KodeWilayah: "31.72", PenyediaID: 3}, false},
		{"kosong", domain.FilterJaringan{}, false},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := k.filter.Validasi(); (err == nil) != k.isSah {
				t.Errorf("err %v, ingin sah %v", err, k.isSah)
			}
		})
	}
}
