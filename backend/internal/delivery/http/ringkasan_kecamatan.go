package http

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

const (
	satuanPaket  = "paket"
	satuanRupiah = "rupiah"
	satuanRasio  = "rasio"
)

type LayananRingkasanKecamatan interface {
	Ringkas(ctx context.Context, kode string) (domain.RingkasanKecamatan, error)
}

type indikatorJSON struct {
	Periode     string `json:"periode"`
	Kebutuhan   Angka  `json:"kebutuhan"`
	Kapasitas   Angka  `json:"kapasitas"`
	Kesenjangan Angka  `json:"kesenjangan"`
	Rasio       *Angka `json:"rasio"`
	Kelas       string `json:"kelas"`
}

type temaKecamatanJSON struct {
	Kode      string         `json:"kode"`
	Nama      domain.Teks    `json:"nama"`
	Judul     domain.Teks    `json:"judul"`
	Indikator *indikatorJSON `json:"indikator"`
}

type anggaranJSON struct {
	Tahun       int   `json:"tahun"`
	JumlahPaket Angka `json:"jumlah_paket"`
	TotalHPS    Angka `json:"total_hps"`
}

type ringkasanKecamatanJSON struct {
	Kecamatan        wilayahDitemukanJSON  `json:"kecamatan"`
	KelurahanDisorot *ringkasanWilayahJSON `json:"kelurahan_disorot"`
	Tema             []temaKecamatanJSON   `json:"tema"`
	Anggaran         *anggaranJSON         `json:"anggaran"`
	Wakilku          wakilKecamatanJSON    `json:"wakilku"`
}

type handlerRingkasanKecamatan struct {
	layanan LayananRingkasanKecamatan
}

func (h handlerRingkasanKecamatan) ringkas(c fiber.Ctx) error {
	r, err := h.layanan.Ringkas(c.Context(), c.Params("kode"))
	if err != nil {
		return err
	}
	hasil := ringkasanKecamatanJSON{
		Kecamatan: keWilayahDitemukanJSON(r.Kecamatan.WilayahDitemukan),
		Tema:      make([]temaKecamatanJSON, len(r.Tema)),
		Anggaran:  keAnggaranJSON(r.Anggaran),
		Wakilku: wakilKecamatanJSON{
			Kecamatan: keRingkasanWilayahJSON(r.Wakil.Kecamatan), KomisiPengawas: r.Wakil.KomisiPengawas, Dapil: keDapilJSON(r.Wakil.DetailDapil),
		},
	}
	if r.KelurahanDisorot != nil {
		disorot := keRingkasanWilayahJSON(*r.KelurahanDisorot)
		hasil.KelurahanDisorot = &disorot
	}
	for i, t := range r.Tema {
		hasil.Tema[i] = keTemaKecamatanJSON(t)
	}
	return kirimData(c, hasil)
}

func keTemaKecamatanJSON(t domain.TemaKecamatan) temaKecamatanJSON {
	hasil := temaKecamatanJSON{Kode: t.Tema.Kode, Nama: t.Tema.Nama, Judul: t.Tema.Indikator.Judul}
	i := t.Indikator
	if i == nil {
		return hasil
	}
	angka := func(nilai float64, satuan string) Angka {
		return Angka{Nilai: nilai, Satuan: satuan, Label: i.Label, ProvenanceID: i.ProvenanceID}
	}
	hasil.Indikator = &indikatorJSON{
		Periode:     i.Periode,
		Kebutuhan:   angka(i.Kebutuhan, t.Tema.Indikator.SatuanKebutuhan),
		Kapasitas:   angka(i.Kapasitas, t.Tema.Indikator.SatuanKapasitas),
		Kesenjangan: angka(i.Kesenjangan, t.Tema.Indikator.SatuanKebutuhan),
		Kelas:       i.Kelas,
	}
	if i.Rasio != nil {
		rasio := angka(*i.Rasio, satuanRasio)
		hasil.Indikator.Rasio = &rasio
	}
	return hasil
}

func keAnggaranJSON(a *domain.RingkasanAnggaran) *anggaranJSON {
	if a == nil {
		return nil
	}
	return &anggaranJSON{
		Tahun:       a.Tahun,
		JumlahPaket: Angka{Nilai: float64(a.JumlahPaket), Satuan: satuanPaket, Label: a.Label(), ProvenanceID: a.ProvenanceID},
		TotalHPS:    Angka{Nilai: a.TotalHPS, Satuan: satuanRupiah, Label: a.Label(), ProvenanceID: a.ProvenanceID},
	}
}
