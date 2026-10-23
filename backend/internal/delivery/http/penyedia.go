package http

import (
	"context"
	"strconv"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

type LayananPenyedia interface {
	Ambil(ctx context.Context, id int64) (usecase.ProfilPenyediaBeratribusi, error)
}

type tahunPenyediaJSON struct {
	Tahun       int    `json:"tahun"`
	Sumber      string `json:"sumber"`
	JumlahPaket Angka  `json:"jumlah_paket"`
	TotalNilai  Angka  `json:"total_nilai"`
}

type riwayatPenyediaJSON struct {
	Tahun        int    `json:"tahun"`
	Pembeli      string `json:"pembeli"`
	NamaPaket    string `json:"nama_paket"`
	Nilai        *Angka `json:"nilai"`
	Sumber       string `json:"sumber"`
	ProvenanceID string `json:"provenance_id"`
}

type profilPenyediaJSON struct {
	ID              int64                 `json:"id"`
	Nama            string                `json:"nama"`
	Alias           []string              `json:"alias"`
	PerTahun        []tahunPenyediaJSON   `json:"per_tahun"`
	KategoriUtama   []string              `json:"kategori_utama"`
	InstansiPemberi []string              `json:"instansi_pemberi"`
	Paket           []paketRingkasJSON    `json:"paket"`
	Riwayat         []riwayatPenyediaJSON `json:"riwayat_opentender"`
	Atribusi        string                `json:"atribusi"`
}

type handlerPenyedia struct {
	layanan LayananPenyedia
}

func (h handlerPenyedia) ambil(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return domain.InputTidakValid("id penyedia wajib berupa angka")
	}
	hasil, err := h.layanan.Ambil(c.Context(), id)
	if err != nil {
		return err
	}
	return kirimData(c, keProfilPenyediaJSON(hasil))
}

func keProfilPenyediaJSON(h usecase.ProfilPenyediaBeratribusi) profilPenyediaJSON {
	p := h.Profil
	j := profilPenyediaJSON{
		ID: p.ID, Nama: p.Nama, Alias: p.Alias, Atribusi: h.Atribusi,
		PerTahun: make([]tahunPenyediaJSON, len(p.PerTahun)), KategoriUtama: namaSaja(p.KategoriUtama),
		InstansiPemberi: namaSaja(p.InstansiPemberi), Paket: make([]paketRingkasJSON, len(p.Paket)),
		Riwayat: make([]riwayatPenyediaJSON, len(p.Riwayat)),
	}
	for i, t := range p.PerTahun {
		j.PerTahun[i] = tahunPenyediaJSON{
			Tahun: t.Tahun, Sumber: string(t.Sumber),
			JumlahPaket: Angka{Nilai: float64(t.JumlahPaket), Satuan: satuanPaket, Label: domain.LabelResmi, ProvenanceID: t.ProvenanceID},
			TotalNilai:  Angka{Nilai: t.TotalNilai, Satuan: satuanRupiah, Label: domain.LabelResmi, ProvenanceID: t.ProvenanceID},
		}
	}
	for i, paket := range p.Paket {
		j.Paket[i] = kePaketRingkasJSON(paket)
	}
	for i, r := range p.Riwayat {
		j.Riwayat[i] = riwayatPenyediaJSON{
			Tahun: r.Tahun, Pembeli: r.Pembeli, NamaPaket: r.NamaPaket, Sumber: string(r.Sumber), ProvenanceID: r.ProvenanceID,
			Nilai: angkaOpsional(r.Nilai, satuanRupiah, domain.LabelResmi, r.ProvenanceID),
		}
	}
	return j
}

func namaSaja(daftar []domain.JumlahPerNama) []string {
	nama := make([]string, len(daftar))
	for i, d := range daftar {
		nama[i] = d.Nama
	}
	return nama
}
