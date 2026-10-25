package http

import (
	"cmp"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type detailPendidikanJSON struct {
	Tema                         string                    `json:"tema"`
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
	KursiNegeriPerBentuk         []kursiBentukJSON         `json:"kursi_negeri_per_bentuk"`
	KursiSPMBBersamaPerBentuk    []kursiBentukJSON         `json:"kursi_spmb_bersama_per_bentuk"`
	Rumus                        rumusJSON                 `json:"rumus"`
	Keterbatasan                 []string                  `json:"keterbatasan"`
	Tren                         []titikTrenJSON           `json:"tren"`
}

type kursiBentukJSON struct {
	Bentuk string `json:"bentuk"`
	Kursi  Angka  `json:"kursi"`
}

type rumusJSON struct {
	BentukAsal           string             `json:"bentuk_asal"`
	TingkatAsal          int                `json:"tingkat_asal"`
	BentukTujuan         []string           `json:"bentuk_tujuan"`
	BentukMadrasahTujuan string             `json:"bentuk_madrasah_tujuan"`
	TingkatTujuan        int                `json:"tingkat_tujuan"`
	BatasMuridPerRombel  *int               `json:"batas_murid_per_rombel"`
	AmbangRasio          domain.AmbangRasio `json:"ambang_rasio"`
}

type titikTrenJSON struct {
	Periode string `json:"periode"`
	Rasio   Angka  `json:"rasio"`
}

type fotoJSON struct {
	URL     string `json:"url"`
	Sumber  string `json:"sumber"`
	Halaman string `json:"halaman"`
	Kredit  string `json:"kredit"`
}

type sekolahJSON struct {
	NPSN        string     `json:"npsn"`
	Nama        string     `json:"nama"`
	Bentuk      string     `json:"bentuk"`
	Status      string     `json:"status"`
	Lokasi      *titikJSON `json:"lokasi"`
	Jalur       *string    `json:"jalur"`
	DayaTampung *Angka     `json:"daya_tampung"`
	Foto        *fotoJSON  `json:"foto"`
}

func (h handlerPendidikan) detailKecamatan(c fiber.Ctx) error {
	d, err := h.layanan.DetailPendidikan(c.Context(), c.Params("kode"), c.Params("kecamatan"))
	if err != nil {
		return err
	}
	return kirimData(c, keDetailPendidikanJSON(d))
}

func keDetailPendidikanJSON(d domain.DetailPendidikanKecamatan) detailPendidikanJSON {
	indikator := pembuatAngka{label: d.Label, provenanceID: d.ProvenanceID}
	r := d.Rincian
	return detailPendidikanJSON{
		Tema: d.TemaKode, Kode: d.KodeKecamatan, Nama: d.NamaKecamatan, Periode: d.Periode, Kelas: string(d.Kelas),
		Label: d.Label, ProvenanceID: d.ProvenanceID, MetodeCalonLulusan: r.MetodeCalonLulusan,
		CalonLulusan:                 rujukan(d.ProvenanceCalonLulusan, d.ProvenanceID).berlabel(r.CalonLulusan, satuanSiswa),
		KursiNegeri:                  rujukan(d.ProvenanceKursiJalur(domain.JalurNegeri), d.ProvenanceID).berlabel(r.KursiNegeri, satuanKursi),
		KursiSPMBBersama:             rujukan(d.ProvenanceKursiJalur(domain.JalurSPMBBersama), d.ProvenanceID).berlabel(r.KursiSPMBBersama, satuanKursi),
		KursiMadrasah:                indikator.berlabel(r.KursiMadrasah, satuanKursi),
		Kesenjangan:                  indikator.bulat(d.Kesenjangan, satuanSiswa),
		KesenjanganDenganSPMBBersama: indikator.berlabel(r.KesenjanganDenganSwasta, satuanSiswa),
		Rasio:                        indikator.opsional(d.Rasio, satuanRasio),
		KursiNegeriPerBentuk:         keKursiBentukJSON(d, domain.JalurNegeri, r.KursiNegeriPerBentuk),
		KursiSPMBBersamaPerBentuk:    keKursiBentukJSON(d, domain.JalurSPMBBersama, r.KursiSPMBBersamaPerBentuk),
		Rumus:                        keRumusJSON(d),
		Keterbatasan:                 d.Keterbatasan(),
		Tren:                         keTrenJSON(d.Tren),
	}
}

func keKursiBentukJSON(d domain.DetailPendidikanKecamatan, jalur domain.JalurDayaTampung, daftar []domain.KursiBentuk) []kursiBentukJSON {
	hasil := make([]kursiBentukJSON, len(daftar))
	for i, k := range daftar {
		angka := rujukan(d.ProvenanceKursiBentuk(jalur, k.Bentuk), d.ProvenanceID).bulat(k.Kursi, satuanKursi)
		angka.Label = domain.LabelResmi
		hasil[i] = kursiBentukJSON{Bentuk: string(k.Bentuk), Kursi: angka}
	}
	return hasil
}

func keRumusJSON(d domain.DetailPendidikanKecamatan) rumusJSON {
	p := d.Peralihan
	tujuan := make([]string, len(p.BentukTujuan))
	for i, bentuk := range p.BentukTujuan {
		tujuan[i] = string(bentuk)
	}
	hasil := rumusJSON{
		BentukAsal: string(p.BentukAsal), TingkatAsal: p.TingkatAsalTerakhir(), BentukTujuan: tujuan,
		BentukMadrasahTujuan: string(p.BentukMadrasahTujuan), TingkatTujuan: p.TingkatTujuanPertama(), AmbangRasio: d.Ambang,
	}
	if batas := d.Rincian.BatasMuridPerRombel; batas > 0 {
		hasil.BatasMuridPerRombel = &batas
	}
	return hasil
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
	daftar, err := h.layanan.SekolahKecamatan(c.Context(), c.Params("kode"), c.Params("kecamatan"))
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
	hasil.Foto = keFotoJSON(s.Foto)
	return hasil
}

func keFotoJSON(f domain.Foto) *fotoJSON {
	if f.URL == "" {
		return nil
	}
	return &fotoJSON{URL: f.URL, Sumber: string(f.Sumber), Halaman: f.Halaman, Kredit: f.Kredit}
}
