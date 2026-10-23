package http

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

const (
	satuanSiswa        = "siswa"
	cacheKontrolPeta   = "public, max-age=3600"
	tipeFeatureCollect = "FeatureCollection"
	tipeFeature        = "Feature"
)

type LayananIndikator interface {
	PetaTema(ctx context.Context, temaKode string) (domain.PetaTema, error)
	DetailPendidikan(ctx context.Context, kodeKecamatan string) (domain.DetailPendidikanKecamatan, error)
	SekolahKecamatan(ctx context.Context, kodeKecamatan string) ([]domain.SekolahDenganKursi, error)
}

type handlerPendidikan struct {
	layanan LayananIndikator
}

func (h handlerPendidikan) petaTema(c fiber.Ctx) error {
	peta, err := h.layanan.PetaTema(c.Context(), c.Params("kode"))
	if err != nil {
		return err
	}
	koleksi := koleksiFiturJSON{Type: tipeFeatureCollect, Features: make([]fiturJSON, len(peta.Fitur))}
	for i, f := range peta.Fitur {
		koleksi.Features[i] = keFiturJSON(f, peta.Tema.Indikator)
	}
	c.Set(fiber.HeaderCacheControl, cacheKontrolPeta)
	return kirimData(c, koleksi)
}

func keFiturJSON(f domain.FiturPetaIndikator, indikator domain.Indikator) fiturJSON {
	angka := pembuatAngka{label: f.Label, provenanceID: f.ProvenanceID}
	return fiturJSON{
		Type:     tipeFeature,
		Geometry: f.GeometriGeoJSON,
		Properties: map[string]any{
			"kode": f.KodeKecamatan, "nama": f.NamaKecamatan, "periode": f.Periode, "kelas": string(f.Kelas),
			"label": f.Label, "provenance_id": f.ProvenanceID,
			"kebutuhan":   angka.bulat(f.Kebutuhan, indikator.SatuanKebutuhan),
			"kapasitas":   angka.bulat(f.Kapasitas, indikator.SatuanKapasitas),
			"kesenjangan": angka.bulat(f.Kesenjangan, indikator.SatuanKebutuhan),
			"rasio":       angka.opsional(f.Rasio, satuanRasio),
		},
	}
}

type pembuatAngka struct {
	label        domain.Label
	provenanceID string
}

func (p pembuatAngka) bulat(nilai int, satuan string) Angka {
	return Angka{Nilai: float64(nilai), Satuan: satuan, Label: p.label, ProvenanceID: p.provenanceID}
}

func (p pembuatAngka) opsional(nilai *float64, satuan string) *Angka {
	if nilai == nil {
		return nil
	}
	return &Angka{Nilai: *nilai, Satuan: satuan, Label: p.label, ProvenanceID: p.provenanceID}
}

func (p pembuatAngka) berlabel(nilai domain.NilaiBerlabel, satuan string) *Angka {
	if nilai.Nilai == nil {
		return nil
	}
	return &Angka{Nilai: float64(*nilai.Nilai), Satuan: satuan, Label: nilai.Label, ProvenanceID: p.provenanceID}
}
