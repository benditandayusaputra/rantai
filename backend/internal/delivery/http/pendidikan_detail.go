package http

import (
	"cmp"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type detailPendidikanJSON struct {
	Kode                         string                    `json:"kode"`
	Nama                         string                    `json:"nama"`
	Periode                      string                    `json:"periode"`
	Kelas                        string                    `json:"kelas"`
	Label                        domain.Label              `json:"label"`
	ProvenanceID                 string                    `json:"provenance_id"`
	MetodeCalonLulusan           domain.MetodeCalonLulusan `json:"metode_calon_lulusan"`
	CalonLulusan                 *Angka                    `json:"calon_lulusan"`
	KursiNegeri                  *Angka                    `json:"kursi_negeri"`
	KursiSPMBBersama             *Angka                    `json:"kursi_spmb_bersama"`
	KursiMadrasah                *Angka                    `json:"kursi_madrasah"`
	Kesenjangan                  Angka                     `json:"kesenjangan"`
	KesenjanganDenganSPMBBersama *Angka                    `json:"kesenjangan_dengan_spmb_bersama"`
	Rasio                        *Angka                    `json:"rasio"`
	Keterbatasan                 []string                  `json:"keterbatasan"`
	Tren                         []titikTrenJSON           `json:"tren"`
}

type titikTrenJSON struct {
	Periode string `json:"periode"`
	Rasio   Angka  `json:"rasio"`
}

type sekolahJSON struct {
	NPSN        string     `json:"npsn"`
	Nama        string     `json:"nama"`
	Bentuk      string     `json:"bentuk"`
	Status      string     `json:"status"`
	Lokasi      *titikJSON `json:"lokasi"`
	Jalur       *string    `json:"jalur"`
	DayaTampung *Angka     `json:"daya_tampung"`
}

func (h handlerPendidikan) detailKecamatan(c fiber.Ctx) error {
	d, err := h.layanan.DetailPendidikan(c.Context(), c.Params("kode"))
	if err != nil {
		return err
	}
	return kirimData(c, keDetailPendidikanJSON(d))
}

func keDetailPendidikanJSON(d domain.DetailPendidikanKecamatan) detailPendidikanJSON {
	indikator := pembuatAngka{label: d.Label, provenanceID: d.ProvenanceID}
	r := d.Rincian
	return detailPendidikanJSON{
		Kode: d.KodeKecamatan, Nama: d.NamaKecamatan, Periode: d.Periode, Kelas: string(d.Kelas),
		Label: d.Label, ProvenanceID: d.ProvenanceID, MetodeCalonLulusan: r.MetodeCalonLulusan,
		CalonLulusan:                 rujukan(d.ProvenanceCalonLulusan, d.ProvenanceID).berlabel(r.CalonLulusan, satuanSiswa),
		KursiNegeri:                  rujukan(d.ProvenanceKursiNegeri, d.ProvenanceID).berlabel(r.KursiNegeri, satuanKursi),
		KursiSPMBBersama:             rujukan(d.ProvenanceKursiSPMBBersama, d.ProvenanceID).berlabel(r.KursiSPMBBersama, satuanKursi),
		KursiMadrasah:                indikator.berlabel(r.KursiMadrasah, satuanKursi),
		Kesenjangan:                  indikator.bulat(d.Kesenjangan, satuanSiswa),
		KesenjanganDenganSPMBBersama: indikator.berlabel(r.KesenjanganDenganSwasta, satuanSiswa),
		Rasio:                        indikator.opsional(d.Rasio, satuanRasio),
		Keterbatasan:                 d.Keterbatasan(),
		Tren:                         keTrenJSON(d.Tren),
	}
}

func keTrenJSON(tren []domain.TitikTren) []titikTrenJSON {
	hasil := make([]titikTrenJSON, len(tren))
	for i, titik := range tren {
		rasio := Angka{Nilai: titik.Rasio, Satuan: satuanRasio, Label: titik.Label, ProvenanceID: titik.ProvenanceID}
		hasil[i] = titikTrenJSON{Periode: titik.Periode, Rasio: rasio}
	}
	return hasil
}

func rujukan(provenanceID, cadangan string) pembuatAngka {
	return pembuatAngka{provenanceID: cmp.Or(provenanceID, cadangan)}
}

func (h handlerPendidikan) sekolahKecamatan(c fiber.Ctx) error {
	daftar, err := h.layanan.SekolahKecamatan(c.Context(), c.Params("kode"))
	if err != nil {
		return err
	}
	return kirimData(c, keDaftarSekolahJSON(daftar))
}

func keDaftarSekolahJSON(daftar []domain.SekolahDenganKursi) []sekolahJSON {
	hasil := make([]sekolahJSON, len(daftar))
	for i, s := range daftar {
		hasil[i] = keSekolahJSON(s)
	}
	return hasil
}

func keSekolahJSON(s domain.SekolahDenganKursi) sekolahJSON {
	hasil := sekolahJSON{NPSN: s.NPSN, Nama: s.Nama, Bentuk: string(s.Bentuk), Status: string(s.Status)}
	if s.Lokasi != nil {
		hasil.Lokasi = &titikJSON{Lat: s.Lokasi.Lat, Lng: s.Lokasi.Lng}
	}
	if s.Kursi != nil {
		jalur := string(s.Jalur)
		hasil.Jalur = &jalur
		hasil.DayaTampung = &Angka{Nilai: float64(*s.Kursi), Satuan: satuanKursi, Label: domain.LabelResmi, ProvenanceID: s.ProvenanceKursi}
	}
	return hasil
}
