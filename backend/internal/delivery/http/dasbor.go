package http

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananDasbor interface {
	Layanan(ctx context.Context, kodeTema string) (domain.DasborLayanan, error)
	Uang(ctx context.Context, tahun int) (domain.DasborUang, error)
	TanggungJawab(ctx context.Context, tahun int) (domain.DasborTanggungJawab, error)
}

type porsiJSON struct {
	Kunci  string  `json:"kunci"`
	Nama   string  `json:"nama,omitempty"`
	Jumlah int     `json:"jumlah"`
	Nilai  float64 `json:"nilai"`
}

type seriJSON struct {
	Label        domain.Label `json:"label"`
	ProvenanceID string       `json:"provenance_id"`
	Isi          []porsiJSON  `json:"isi"`
}

type rekapKursiJSON struct {
	Kunci     string                    `json:"kunci"`
	Nama      string                    `json:"nama,omitempty"`
	Kebutuhan float64                   `json:"kebutuhan"`
	Kapasitas float64                   `json:"kapasitas"`
	PerKelas  map[domain.KelasRasio]int `json:"per_kelas"`
}

type dasborLayananJSON struct {
	Tema         string           `json:"tema"`
	Label        domain.Label     `json:"label"`
	ProvenanceID string           `json:"provenance_id"`
	Periode      []rekapKursiJSON `json:"periode"`
	Kota         []rekapKursiJSON `json:"kota"`
}

type dasborUangJSON struct {
	Tahun     int      `json:"tahun"`
	Tahunan   seriJSON `json:"tahunan"`
	Bulanan   seriJSON `json:"bulanan"`
	Jenis     seriJSON `json:"jenis"`
	Metode    seriJSON `json:"metode"`
	Kota      seriJSON `json:"kota"`
	Kecamatan seriJSON `json:"kecamatan"`
}

type dasborTanggungJawabJSON struct {
	Tahun           int      `json:"tahun"`
	KegiatanTahunan seriJSON `json:"kegiatan_tahunan"`
	KegiatanBulanan seriJSON `json:"kegiatan_bulanan"`
	Topik           seriJSON `json:"topik"`
	Reses           seriJSON `json:"reses"`
	Instansi        seriJSON `json:"instansi"`
}

type handlerDasbor struct {
	layanan LayananDasbor
}

func daftarkanDasbor(api fiber.Router, h handlerDasbor) {
	api.Get("/tema/:kode/dasbor", h.layananPublik)
	api.Get("/dasbor/uang", h.uang)
	api.Get("/dasbor/tanggung-jawab", h.tanggungJawab)
}

func (h handlerDasbor) layananPublik(c fiber.Ctx) error {
	d, err := h.layanan.Layanan(c.Context(), c.Params("kode"))
	if err != nil {
		return err
	}
	return kirimData(c, dasborLayananJSON{
		Tema: d.Tema, Label: d.Label, ProvenanceID: d.ProvenanceID,
		Periode: keDaftarRekapJSON(d.Periode), Kota: keDaftarRekapJSON(d.Kota),
	})
}

func (h handlerDasbor) uang(c fiber.Ctx) error {
	tahun, err := tahunDasbor(c)
	if err != nil {
		return err
	}
	d, err := h.layanan.Uang(c.Context(), tahun)
	if err != nil {
		return err
	}
	return kirimData(c, dasborUangJSON{
		Tahun: d.Tahun, Tahunan: keSeriJSON(d.Tahunan), Bulanan: keSeriJSON(d.Bulanan), Jenis: keSeriJSON(d.Jenis),
		Metode: keSeriJSON(d.Metode), Kota: keSeriJSON(d.Kota), Kecamatan: keSeriJSON(d.Kecamatan),
	})
}

func (h handlerDasbor) tanggungJawab(c fiber.Ctx) error {
	tahun, err := tahunDasbor(c)
	if err != nil {
		return err
	}
	d, err := h.layanan.TanggungJawab(c.Context(), tahun)
	if err != nil {
		return err
	}
	return kirimData(c, dasborTanggungJawabJSON{
		Tahun: d.Tahun, KegiatanTahunan: keSeriJSON(d.KegiatanTahunan), KegiatanBulanan: keSeriJSON(d.KegiatanBulanan),
		Topik: keSeriJSON(d.Topik), Reses: keSeriJSON(d.Reses), Instansi: keSeriJSON(d.Instansi),
	})
}

func tahunDasbor(c fiber.Ctx) (int, error) {
	tahun, err := bilanganOpsional(c.Query("tahun"))
	if err != nil {
		return 0, domain.InputTidakValid("tahun wajib berupa angka")
	}
	return tahun, nil
}

func keSeriJSON(s domain.Seri) seriJSON {
	isi := make([]porsiJSON, len(s.Isi))
	for i, p := range s.Isi {
		isi[i] = porsiJSON(p)
	}
	return seriJSON{Label: s.Label, ProvenanceID: s.ProvenanceID, Isi: isi}
}

func keDaftarRekapJSON(daftar []domain.RekapKursi) []rekapKursiJSON {
	hasil := make([]rekapKursiJSON, len(daftar))
	for i, r := range daftar {
		hasil[i] = rekapKursiJSON(r)
	}
	return hasil
}
