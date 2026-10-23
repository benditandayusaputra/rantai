package http

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananPembanding interface {
	Bandingkan(ctx context.Context, kode string) (domain.PembandingPaket, error)
}

type statistikJSON struct {
	Min    Angka `json:"min"`
	Q1     Angka `json:"q1"`
	Median Angka `json:"median"`
	Q3     Angka `json:"q3"`
	Maks   Angka `json:"maks"`
}

type nilaiPembandingJSON struct {
	Kode string `json:"kode"`
	HPS  Angka  `json:"hps"`
}

type pembandingJSON struct {
	Kode             string                `json:"kode"`
	Kategori         string                `json:"kategori"`
	JenisPengadaan   string                `json:"jenis_pengadaan"`
	TahunAnggaran    int                   `json:"tahun_anggaran"`
	IsTersedia       bool                  `json:"tersedia"`
	HPS              *Angka                `json:"hps"`
	JumlahPembanding *Angka                `json:"jumlah_pembanding"`
	Statistik        *statistikJSON        `json:"statistik"`
	Nilai            []nilaiPembandingJSON `json:"nilai"`
}

type handlerPembanding struct {
	layanan LayananPembanding
}

func (h handlerPembanding) bandingkan(c fiber.Ctx) error {
	b, err := h.layanan.Bandingkan(c.Context(), c.Params("kode"))
	if err != nil {
		return err
	}
	p := b.Paket
	hasil := pembandingJSON{
		Kode: p.Kode, Kategori: p.Kategori, JenisPengadaan: p.JenisPengadaan, TahunAnggaran: p.TahunAnggaran,
		IsTersedia: b.IsTersedia, HPS: angkaOpsional(p.HPS, satuanRupiah, domain.LabelResmi, p.ProvenanceID),
		Nilai: make([]nilaiPembandingJSON, len(b.Nilai)),
	}
	for i, n := range b.Nilai {
		hasil.Nilai[i] = nilaiPembandingJSON{Kode: n.Kode, HPS: Angka{Nilai: n.HPS, Satuan: satuanRupiah, Label: domain.LabelResmi, ProvenanceID: n.ProvenanceID}}
	}
	if b.IsTersedia {
		hasil.Statistik, hasil.JumlahPembanding = keStatistikJSON(b.Statistik, b.ProvenanceTerbaru)
	}
	return kirimData(c, hasil)
}

func keStatistikJSON(s domain.Statistik, provenanceID string) (*statistikJSON, *Angka) {
	rupiah := func(nilai float64) Angka {
		return Angka{Nilai: nilai, Satuan: satuanRupiah, Label: domain.LabelResmi, ProvenanceID: provenanceID}
	}
	jumlah := Angka{Nilai: float64(s.Jumlah), Satuan: satuanPaket, Label: domain.LabelResmi, ProvenanceID: provenanceID}
	return &statistikJSON{Min: rupiah(s.Min), Q1: rupiah(s.Q1), Median: rupiah(s.Median), Q3: rupiah(s.Q3), Maks: rupiah(s.Maks)}, &jumlah
}
