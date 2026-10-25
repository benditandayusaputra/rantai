package http

import (
	"context"
	"strconv"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

type LayananSekolahTerdekat interface {
	Cari(ctx context.Context, permintaan usecase.PermintaanSekolahTerdekat) (domain.SekitarTitik, error)
}

type sekolahTerdekatJSON struct {
	sekolahJSON
	JarakMeter float64 `json:"jarak_m"`
}

type tekananRadiusJSON struct {
	CalonLulusan *Angka `json:"calon_lulusan"`
	KursiNegeri  *Angka `json:"kursi_negeri"`
	Rasio        *Angka `json:"rasio"`
}

type sekitarTitikJSON struct {
	Tema          string                `json:"tema"`
	Titik         titikJSON             `json:"titik"`
	RadiusMeter   int                   `json:"radius_m"`
	Negeri        []sekolahTerdekatJSON `json:"negeri"`
	SPMBBersama   []sekolahTerdekatJSON `json:"spmb_bersama"`
	Madrasah      []sekolahTerdekatJSON `json:"madrasah"`
	TekananRadius tekananRadiusJSON     `json:"tekanan_radius"`
	Keterbatasan  []string              `json:"keterbatasan"`
}

type handlerSekolahTerdekat struct {
	layanan LayananSekolahTerdekat
}

func (h handlerSekolahTerdekat) cari(c fiber.Ctx) error {
	titik, err := uraiTitik(c.Query("lat"), c.Query("lng"))
	if err != nil {
		return err
	}
	radius, err := uraiRadius(c.Query("radius_m"))
	if err != nil {
		return err
	}
	sekitar, err := h.layanan.Cari(c.Context(), usecase.PermintaanSekolahTerdekat{Titik: titik, RadiusMeter: radius, TemaKode: c.Query("tema")})
	if err != nil {
		return err
	}
	return kirimData(c, keSekitarTitikJSON(sekitar))
}

func uraiRadius(nilai string) (int, error) {
	if nilai == "" {
		return 0, nil
	}
	radius, err := strconv.Atoi(nilai)
	if err != nil {
		return 0, domain.InputTidakValid("radius_m wajib berupa bilangan bulat")
	}
	return radius, nil
}

func keSekitarTitikJSON(s domain.SekitarTitik) sekitarTitikJSON {
	kelompok := s.Kelompokkan()
	kursi, provenanceKursi := domain.KursiNegeriDalamRadius(kelompok.Negeri)
	return sekitarTitikJSON{
		Tema:          s.Peralihan.TemaKode,
		Titik:         titikJSON(s.Titik),
		RadiusMeter:   s.RadiusMeter,
		Negeri:        keDaftarSekolahTerdekat(kelompok.Negeri),
		SPMBBersama:   keDaftarSekolahTerdekat(kelompok.SPMBBersama),
		Madrasah:      keDaftarSekolahTerdekat(kelompok.Madrasah),
		TekananRadius: keTekananRadiusJSON(s, kursi, provenanceKursi),
		Keterbatasan:  []string{s.Peralihan.KeterbatasanCalonLulusan(), domain.KeterbatasanJarakGarisLurus},
	}
}

func keTekananRadiusJSON(s domain.SekitarTitik, kursi int, provenanceKursi string) tekananRadiusJSON {
	hasil := tekananRadiusJSON{}
	if s.ProvenanceSiswa != "" {
		siswa := pembuatAngka{label: domain.LabelResmi, provenanceID: s.ProvenanceSiswa}
		jumlahSiswa := siswa.bulat(s.SiswaTingkatAkhir, satuanSiswa)
		hasil.CalonLulusan = &jumlahSiswa
		hasil.Rasio = siswa.opsional(domain.TekananRadius(s.SiswaTingkatAkhir, kursi), satuanRasio)
	}
	if provenanceKursi != "" {
		jumlahKursi := pembuatAngka{label: domain.LabelResmi, provenanceID: provenanceKursi}.bulat(kursi, satuanKursi)
		hasil.KursiNegeri = &jumlahKursi
	}
	return hasil
}

func keDaftarSekolahTerdekat(daftar []domain.SekolahTerdekat) []sekolahTerdekatJSON {
	hasil := make([]sekolahTerdekatJSON, len(daftar))
	for i, s := range daftar {
		hasil[i] = sekolahTerdekatJSON{sekolahJSON: keSekolahJSON(s.SekolahDenganKursi), JarakMeter: s.JarakMeter}
	}
	return hasil
}
