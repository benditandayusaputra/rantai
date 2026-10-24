package http

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const provenanceAlatUji = "prv_alat_uji"

var (
	kecamatanAlatUji = domain.RingkasanWilayah{Kode: "31.72.01", Nama: "Cilincing", Tingkat: domain.TingkatKecamatan}
	dapilAlatUji     = domain.DetailDapil{
		RingkasanDapil: domain.RingkasanDapil{ID: 1, Nomor: 1, JumlahKursi: 10, ProvenanceID: provenanceAlatUji, ProvenanceAnggota: provenanceAlatUji},
		Kecamatan:      []domain.RingkasanWilayah{kecamatanAlatUji},
		Anggota:        []domain.Wakil{{ID: "agt_1", Nama: "Anggota Uji", KodeKomisi: "E"}},
	}
	nilaiAlatUji = 1500000000.0
	kursiAlatUji = 288
)

type wilayahAlatTiruan struct{ LayananWilayah }

func (wilayahAlatTiruan) Cari(context.Context, string) ([]domain.WilayahDitemukan, error) {
	return []domain.WilayahDitemukan{{RingkasanWilayah: kecamatanAlatUji, ProvenanceID: provenanceAlatUji}}, nil
}

type ringkasanAlatTiruan struct{}

func (ringkasanAlatTiruan) Ringkas(context.Context, string) (domain.RingkasanKecamatan, error) {
	kecamatan := domain.DetailWilayah{WilayahDitemukan: domain.WilayahDitemukan{RingkasanWilayah: kecamatanAlatUji, ProvenanceID: provenanceAlatUji}}
	return domain.RingkasanKecamatan{Kecamatan: kecamatan, Wakil: domain.WakilKecamatan{Kecamatan: kecamatanAlatUji, DetailDapil: dapilAlatUji}}, nil
}

type indikatorAlatTiruan struct{ LayananIndikator }

func (indikatorAlatTiruan) DetailPendidikan(context.Context, string) (domain.DetailPendidikanKecamatan, error) {
	indikator := domain.IndikatorKesenjanganSMP{KodeKecamatan: kecamatanAlatUji.Kode, Kesenjangan: 120, Label: domain.LabelEstimasi}
	return domain.DetailPendidikanKecamatan{IndikatorTersimpan: domain.IndikatorTersimpan{IndikatorKesenjanganSMP: indikator, ProvenanceID: provenanceAlatUji}}, nil
}

func (indikatorAlatTiruan) SekolahKecamatan(context.Context, string) ([]domain.SekolahDenganKursi, error) {
	return []domain.SekolahDenganKursi{{NPSN: "20100001", Nama: "SMPN Uji", Kursi: &kursiAlatUji, ProvenanceKursi: provenanceAlatUji}}, nil
}

type paketAlatTiruan struct{ LayananPaket }

func (paketAlatTiruan) Ambil(context.Context, string) (usecase.DetailPaket, error) {
	paket := domain.PaketLengkap{Kode: "pkt_1", Nama: "Rehab SDN Uji", HPS: &nilaiAlatUji, ProvenanceID: provenanceAlatUji}
	return usecase.DetailPaket{Paket: paket, Rantai: []domain.MataRantai{}}, nil
}

func (paketAlatTiruan) DaftarKecamatan(context.Context, usecase.FilterPaketKecamatan) (usecase.HalamanPaket, error) {
	paket := domain.PaketRingkas{Kode: "pkt_1", Nama: "Rehab SDN Uji", HPS: &nilaiAlatUji, ProvenanceID: provenanceAlatUji}
	return usecase.HalamanPaket{Paket: []domain.PaketRingkas{paket}, Total: 1}, nil
}

type belanjaAlatTiruan struct{}

func (belanjaAlatTiruan) Hitung(context.Context, string, int) (domain.BelanjaTema, error) {
	anggaran := domain.RingkasanAnggaran{Tahun: 2026, JumlahPaket: 3, TotalHPS: nilaiAlatUji, ProvenanceID: provenanceAlatUji}
	return domain.BelanjaTema{
		Tema:      domain.TemaTersimpan{Tema: domain.Tema{Kode: domain.TemaPendidikan}},
		Kecamatan: []domain.BelanjaKecamatan{{Kecamatan: kecamatanAlatUji, Anggaran: &anggaran}},
	}, nil
}

type temaAlatTiruan struct{ LayananTema }

func (temaAlatTiruan) Ambil(context.Context, string) (domain.TemaTersimpan, error) {
	tanggungJawab := []domain.TanggungJawab{{Masalah: "kursi_smp", InstansiPelaksana: "Dinas Pendidikan", Komisi: "E"}}
	tema := domain.Tema{Kode: domain.TemaPendidikan, KonfigurasiTema: domain.KonfigurasiTema{TanggungJawab: tanggungJawab}}
	return domain.TemaTersimpan{Tema: tema, ProvenanceID: provenanceAlatUji}, nil
}

type wakilAlatTiruan struct{ LayananWakil }

func (wakilAlatTiruan) WakilKecamatan(context.Context, string, string) (domain.WakilKecamatan, error) {
	return domain.WakilKecamatan{Kecamatan: kecamatanAlatUji, KomisiPengawas: []string{"E"}, DetailDapil: dapilAlatUji}, nil
}

type umpanAlatTiruan struct{}

func (umpanAlatTiruan) UmpanAnggota(context.Context, domain.FilterUmpan) (domain.HalamanUmpan, error) {
	item := domain.UmpanTampil{ID: 1, JenisSumber: domain.SumberUmpanResmiDPRD, Judul: "Rapat Komisi E", ProvenanceID: provenanceAlatUji}
	return domain.HalamanUmpan{Item: []domain.UmpanTampil{item}, Total: 1}, nil
}

func layananAlatUji() Layanan {
	return Layanan{
		Wilayah: wilayahAlatTiruan{}, Ringkasan: ringkasanAlatTiruan{}, Indikator: indikatorAlatTiruan{}, Paket: paketAlatTiruan{},
		Belanja: belanjaAlatTiruan{}, Tema: temaAlatTiruan{}, Wakil: wakilAlatTiruan{}, Umpan: umpanAlatTiruan{},
	}
}

func TestRegistryTanya_SetiapAlat_HasilMembawaProvenance(t *testing.T) {
	masukan := map[string]string{
		"cari_wilayah":         `{"teks":"cilincing"}`,
		"ringkasan_kecamatan":  `{"kode":"31.72.01"}`,
		"indikator_tema":       `{"kode_kecamatan":"31.72.01","tema":"pendidikan"}`,
		"sekolah_di_kecamatan": `{"kode":"31.72.01"}`,
		"daftar_paket":         `{"kode_kecamatan":"31.72.01","tema":"pendidikan","tahun":2026}`,
		"detail_paket":         `{"kode":"pkt_1"}`,
		"belanja_tema":         `{"tema":"pendidikan"}`,
		"tanggung_jawab_tema":  `{"tema":"pendidikan"}`,
		"wakil_dapil":          `{"kode_kecamatan":"31.72.01","tema":"pendidikan"}`,
		"umpan_anggota":        `{"id_anggota":"agt_1"}`,
	}
	registry, err := RegistryTanya(layananAlatUji())
	if err != nil {
		t.Fatal(err)
	}
	if jumlah := len(registry.Daftar()); jumlah != len(masukan) {
		t.Fatalf("jumlah alat = %d, harap %d", jumlah, len(masukan))
	}
	for _, spek := range registry.Daftar() {
		t.Run(spek.Nama, func(t *testing.T) {
			hasil, err := registry.Jalankan(t.Context(), spek.Nama, json.RawMessage(masukan[spek.Nama]))

			if err != nil {
				t.Fatalf("jalankan: %v", err)
			}
			if !strings.Contains(string(hasil), `"provenance_id":"`+provenanceAlatUji+`"`) {
				t.Errorf("hasil tanpa provenance_id: %s", hasil)
			}
		})
	}
}

func TestAlatIndikatorTema_TemaBelumPunyaRincian_TidakDitemukan(t *testing.T) {
	registry, err := RegistryTanya(layananAlatUji())
	if err != nil {
		t.Fatal(err)
	}

	_, err = registry.Jalankan(t.Context(), "indikator_tema", json.RawMessage(`{"kode_kecamatan":"31.72.01","tema":"kesehatan"}`))

	var kesalahan *domain.Kesalahan
	if !errors.As(err, &kesalahan) || kesalahan.Kategori != domain.KategoriTidakDitemukan {
		t.Errorf("err = %v, harap kesalahan tidak ditemukan", err)
	}
}
