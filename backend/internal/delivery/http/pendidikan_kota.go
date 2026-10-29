package http

import (
	"context"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type LayananPendidikanKota interface {
	Baca(ctx context.Context, tema, kode string) (domain.PendidikanKota, error)
}

type handlerPendidikanKota struct {
	layanan LayananPendidikanKota
}

type pendidikanKotaJSON struct {
	KabupatenKota []ringkasanWilayahJSON    `json:"kabupaten_kota"`
	Wilayah       ringkasanWilayahJSON      `json:"wilayah"`
	Tema          string                    `json:"tema"`
	Periode       string                    `json:"periode"`
	Kecamatan     []kecamatanPendidikanJSON `json:"kecamatan"`
	Sekolah       []jumlahSekolahJSON       `json:"sekolah"`
	Ringkasan     kapasitasKotaJSON         `json:"ringkasan"`
}

type kecamatanPendidikanJSON struct {
	Kode      string         `json:"kode"`
	Nama      string         `json:"nama"`
	Indikator *indikatorJSON `json:"indikator"`
}

type jumlahSekolahJSON struct {
	KodeKecamatan string    `json:"kode_kecamatan"`
	Bentuk        string    `json:"bentuk"`
	Negeri        Angka     `json:"negeri"`
	Swasta        Angka     `json:"swasta"`
	DiambilPada   time.Time `json:"diambil_pada"`
}

type kapasitasKotaJSON struct {
	TotalKecamatan      int    `json:"total_kecamatan"`
	KecamatanDenganData int    `json:"kecamatan_dengan_data"`
	KecamatanKurang     int    `json:"kecamatan_kurang"`
	Kebutuhan           *Angka `json:"kebutuhan"`
	Kapasitas           *Angka `json:"kapasitas"`
	Kesenjangan         *Angka `json:"kesenjangan"`
	Rasio               *Angka `json:"rasio"`
}

func (h handlerPendidikanKota) baca(c fiber.Ctx) error {
	d, err := h.layanan.Baca(c.Context(), c.Params("kode"), c.Params("kota"))
	if err != nil {
		return err
	}
	return kirimData(c, kePendidikanKotaJSON(d))
}

func kePendidikanKotaJSON(d domain.PendidikanKota) pendidikanKotaJSON {
	return pendidikanKotaJSON{
		Wilayah: keRingkasanWilayahJSON(d.Wilayah), Tema: d.Tema, Periode: d.Periode,
		KabupatenKota: keDaftarRingkasanKota(d.KabupatenKota),
		Kecamatan:     keKecamatanPendidikanJSON(d), Sekolah: keJumlahSekolahJSON(d.Sekolah),
		Ringkasan: keKapasitasKotaJSON(d.Ringkasan, len(d.Kecamatan)),
	}
}

func keKecamatanPendidikanJSON(d domain.PendidikanKota) []kecamatanPendidikanJSON {
	menurutKode := make(map[string]domain.IndikatorTersimpan, len(d.Indikator))
	for _, i := range d.Indikator {
		menurutKode[i.KodeKecamatan] = i
	}
	hasil := make([]kecamatanPendidikanJSON, len(d.Kecamatan))
	for n, w := range d.Kecamatan {
		hasil[n] = kecamatanPendidikanJSON{Kode: w.Kode, Nama: w.Nama}
		if i, ada := menurutKode[w.Kode]; ada {
			hasil[n].Indikator = keIndikatorKotaJSON(i)
		}
	}
	return hasil
}

func keIndikatorKotaJSON(i domain.IndikatorTersimpan) *indikatorJSON {
	a := pembuatAngka{label: i.Label, provenanceID: i.ProvenanceID}
	return &indikatorJSON{
		Periode: i.Periode, Kebutuhan: a.bulat(i.Kebutuhan, satuanSiswa),
		Kapasitas: a.bulat(i.Kapasitas, satuanKursi), Kesenjangan: a.bulat(i.Kesenjangan, satuanSiswa),
		Rasio: a.opsional(i.Rasio, satuanRasio), Kelas: string(i.Kelas),
	}
}

func keJumlahSekolahJSON(daftar []domain.JumlahSekolahKecamatan) []jumlahSekolahJSON {
	hasil := make([]jumlahSekolahJSON, len(daftar))
	for i, s := range daftar {
		a := pembuatAngka{label: s.Label, provenanceID: s.ProvenanceID}
		hasil[i] = jumlahSekolahJSON{
			KodeKecamatan: s.KodeKecamatan, Bentuk: string(s.Bentuk),
			Negeri: a.bulat(s.Negeri, "sekolah"), Swasta: a.bulat(s.Swasta, "sekolah"), DiambilPada: s.DiambilPada,
		}
	}
	return hasil
}

func keKapasitasKotaJSON(r domain.RingkasanKapasitasKota, total int) kapasitasKotaJSON {
	a := pembuatAngka{label: r.Label, provenanceID: r.ProvenanceID}
	opsional := func(nilai *int, satuan string) *Angka {
		if nilai == nil {
			return nil
		}
		angka := a.bulat(*nilai, satuan)
		return &angka
	}
	return kapasitasKotaJSON{
		TotalKecamatan: total, KecamatanDenganData: r.KecamatanDenganData, KecamatanKurang: r.KecamatanKurang,
		Kebutuhan: opsional(r.Kebutuhan, satuanSiswa), Kapasitas: opsional(r.Kapasitas, satuanKursi),
		Kesenjangan: opsional(r.Kesenjangan, satuanSiswa), Rasio: a.opsional(r.Rasio, satuanRasio),
	}
}

func keDaftarRingkasanKota(daftar []domain.RingkasanWilayah) []ringkasanWilayahJSON {
	hasil := make([]ringkasanWilayahJSON, len(daftar))
	for i, w := range daftar {
		hasil[i] = keRingkasanWilayahJSON(w)
	}
	return hasil
}
