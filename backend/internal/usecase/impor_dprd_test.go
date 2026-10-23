package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pengambilDPRDTiruan struct {
	data DataSumberDPRD
}

func (p pengambilDPRDTiruan) Ambil(context.Context) (DataSumberDPRD, error) {
	return p.data, nil
}

type wakilTiruan struct {
	kecamatan map[string]string
	diterima  *DataSimpanWakil
}

func (w *wakilTiruan) DaftarKecamatan(context.Context) (map[string]string, error) {
	return w.kecamatan, nil
}

func (w *wakilTiruan) SimpanWakil(_ context.Context, data DataSimpanWakil) (HasilImpor, error) {
	w.diterima = &data
	return HasilImpor{JumlahBaru: len(data.Anggota)}, nil
}

func imporDPRDUji(sumber DataSumberDPRD, wakil *wakilTiruan, pencatat *pencatatTiruan) ImporDPRD {
	return ImporDPRD{
		Pengambil:        pengambilDPRDTiruan{data: sumber},
		PembacaKecamatan: wakil,
		Penyimpan:        wakil,
		Pencatat:         pencatat,
		Dapil:            []domain.Dapil{{Nomor: 1, JumlahKursi: 1, Kecamatan: []domain.KecamatanDapil{{Kode: "31.71.01", Nama: "Gambir"}}}},
		PolaURLProfil:    "https://dprd.example/profile/%s/",
	}
}

func sumberDPRDUji(slugFraksi string) DataSumberDPRD {
	return DataSumberDPRD{
		BarisDapil:    []domain.BarisDapil{{NomorDapil: 1, Nama: "Ima Mahdiah", Slug: "ima-mahdiah"}},
		AnggotaFraksi: []domain.AnggotaFraksi{{Fraksi: "F", BarisStruktur: domain.BarisStruktur{Nama: "Orang Lain", Slug: slugFraksi}}},
	}
}

func TestImporDPRD_AnggotaCocok_DisimpanDenganURLProfil(t *testing.T) {
	wakil := &wakilTiruan{kecamatan: map[string]string{"31.71.01": "Gambir"}}
	pencatat := &pencatatTiruan{}

	err := imporDPRDUji(sumberDPRDUji("ima-mahdiah"), wakil, pencatat).Jalankan(t.Context())

	if err != nil || wakil.diterima == nil || pencatat.selesai == nil {
		t.Fatalf("harap tersimpan dan run selesai, dapat %v", err)
	}
	if url := wakil.diterima.Anggota[0].URLProfil; url != "https://dprd.example/profile/ima-mahdiah/" {
		t.Errorf("url profil = %q", url)
	}
}

func TestImporDPRD_AnggotaTidakCocok_TidakDisimpanDanRunGagal(t *testing.T) {
	wakil := &wakilTiruan{kecamatan: map[string]string{"31.71.01": "Gambir"}}
	pencatat := &pencatatTiruan{}

	err := imporDPRDUji(sumberDPRDUji("orang-lain"), wakil, pencatat).Jalankan(t.Context())

	if !errors.Is(err, domain.ErrAnggotaTidakCocok) || wakil.diterima != nil || pencatat.gagal == nil {
		t.Errorf("harap gagal tanpa menyimpan, dapat %v", err)
	}
}

func TestImporDPRD_KecamatanTanpaDapil_TidakDisimpan(t *testing.T) {
	wakil := &wakilTiruan{kecamatan: map[string]string{"31.71.01": "Gambir", "31.72.03": "Koja"}}
	pencatat := &pencatatTiruan{}

	err := imporDPRDUji(sumberDPRDUji("ima-mahdiah"), wakil, pencatat).Jalankan(t.Context())

	if !errors.Is(err, domain.ErrKecamatanTanpaDapil) || wakil.diterima != nil {
		t.Errorf("harap ErrKecamatanTanpaDapil tanpa menyimpan, dapat %v", err)
	}
}
