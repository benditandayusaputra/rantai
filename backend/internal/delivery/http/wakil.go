package http

import (
	"context"
	"strconv"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

const (
	satuanKursi   = "kursi"
	satuanAnggota = "orang"
)

type LayananWakil interface {
	WakilKecamatan(ctx context.Context, kode, kodeTema string) (domain.WakilKecamatan, error)
	WakilKecamatanDenganUmpan(ctx context.Context, kode, kodeTema string) (domain.WakilKecamatan, error)
	Dapil(ctx context.Context, id int) (domain.DetailDapil, error)
	Anggota(ctx context.Context, id string) (domain.ProfilAnggota, error)
}

type komisiJSON struct {
	Kode   string   `json:"kode"`
	Nama   string   `json:"nama"`
	Bidang []string `json:"bidang,omitempty"`
}

type wakilJSON struct {
	ID                    string      `json:"id"`
	Nama                  string      `json:"nama"`
	Partai                string      `json:"partai"`
	Fraksi                string      `json:"fraksi"`
	Komisi                *komisiJSON `json:"komisi"`
	DudukDiKomisiPengawas bool        `json:"duduk_di_komisi_pengawas"`
	FotoURL               string      `json:"foto_url,omitempty"`
}

type dapilJSON struct {
	ID            int                    `json:"id"`
	Nomor         int                    `json:"nomor"`
	Nama          string                 `json:"nama"`
	TahunPemilu   int                    `json:"tahun_pemilu"`
	JumlahKursi   Angka                  `json:"jumlah_kursi"`
	JumlahAnggota Angka                  `json:"jumlah_anggota"`
	Kecamatan     []ringkasanWilayahJSON `json:"kecamatan"`
	Anggota       []wakilJSON            `json:"anggota"`
	ProvenanceID  string                 `json:"provenance_id"`
}

type wakilKecamatanJSON struct {
	Kecamatan        ringkasanWilayahJSON   `json:"kecamatan"`
	KomisiPengawas   []string               `json:"komisi_pengawas"`
	Dapil            dapilJSON              `json:"dapil"`
	UmpanIsuPengawas map[string][]umpanJSON `json:"umpan_isu_pengawas,omitempty"`
}

type jabatanJSON struct {
	AlatKelengkapan string `json:"alat_kelengkapan"`
	Jabatan         string `json:"jabatan"`
}

type profilAnggotaJSON struct {
	wakilJSON
	KeteranganKomisi string        `json:"keterangan_komisi,omitempty"`
	AlatKelengkapan  []jabatanJSON `json:"alat_kelengkapan"`
	Dapil            dapilJSON     `json:"dapil"`
	URLProfil        string        `json:"url_profil"`
	URLLHKPN         string        `json:"lhkpn_url"`
	ProvenanceID     string        `json:"provenance_id"`
}

type handlerWakil struct {
	layanan LayananWakil
}

func (h handlerWakil) wakilKecamatan(c fiber.Ctx) error {
	hasil, err := h.ambilWakilKecamatan(c)
	if err != nil {
		return err
	}
	return kirimData(c, keWakilKecamatanJSON(hasil))
}

func (h handlerWakil) ambilWakilKecamatan(c fiber.Ctx) (domain.WakilKecamatan, error) {
	denganUmpan, err := uraiBoolOpsional(c.Query("dengan_umpan"))
	switch {
	case err != nil:
		return domain.WakilKecamatan{}, domain.InputTidakValid("dengan_umpan harus true atau false")
	case denganUmpan:
		return h.layanan.WakilKecamatanDenganUmpan(c.Context(), c.Params("kode"), c.Query("tema"))
	default:
		return h.layanan.WakilKecamatan(c.Context(), c.Params("kode"), c.Query("tema"))
	}
}

func uraiBoolOpsional(nilai string) (bool, error) {
	if nilai == "" {
		return false, nil
	}
	return strconv.ParseBool(nilai)
}

func keWakilKecamatanJSON(w domain.WakilKecamatan) wakilKecamatanJSON {
	var umpanIsu map[string][]umpanJSON
	if w.UmpanIsuPengawas != nil {
		umpanIsu = map[string][]umpanJSON{}
		for id, umpan := range w.UmpanIsuPengawas {
			umpanIsu[id] = keDaftarUmpanJSON(umpan)
		}
	}
	return wakilKecamatanJSON{
		Kecamatan:        keRingkasanWilayahJSON(w.Kecamatan),
		KomisiPengawas:   w.KomisiPengawas,
		Dapil:            keDapilJSON(w.DetailDapil),
		UmpanIsuPengawas: umpanIsu,
	}
}

func (h handlerWakil) dapil(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return domain.InputTidakValid("id dapil wajib berupa angka")
	}
	hasil, err := h.layanan.Dapil(c.Context(), id)
	if err != nil {
		return err
	}
	return kirimData(c, keDapilJSON(hasil))
}

func (h handlerWakil) anggota(c fiber.Ctx) error {
	p, err := h.layanan.Anggota(c.Context(), c.Params("id"))
	if err != nil {
		return err
	}
	wakil := keWakilJSON(p.Wakil)
	if wakil.Komisi != nil {
		wakil.Komisi.Bidang = p.BidangKomisi
	}
	jabatan := make([]jabatanJSON, len(p.AlatKelengkapan))
	for i, j := range p.AlatKelengkapan {
		jabatan[i] = jabatanJSON{AlatKelengkapan: j.AlatKelengkapan, Jabatan: j.Jabatan}
	}
	return kirimData(c, profilAnggotaJSON{
		wakilJSON: wakil, KeteranganKomisi: p.KeteranganKomisi(), AlatKelengkapan: jabatan,
		Dapil: keDapilJSON(p.Dapil), URLProfil: p.URLProfil, URLLHKPN: p.URLLHKPN, ProvenanceID: p.ProvenanceID,
	})
}

func keDapilJSON(d domain.DetailDapil) dapilJSON {
	kecamatan := make([]ringkasanWilayahJSON, len(d.Kecamatan))
	for i, k := range d.Kecamatan {
		kecamatan[i] = keRingkasanWilayahJSON(k)
	}
	anggota := make([]wakilJSON, len(d.Anggota))
	for i, w := range d.Anggota {
		anggota[i] = keWakilJSON(w)
	}
	return dapilJSON{
		ID: d.ID, Nomor: d.Nomor, Nama: d.Nama, TahunPemilu: d.TahunPemilu,
		JumlahKursi:   Angka{Nilai: float64(d.JumlahKursi), Satuan: satuanKursi, Label: domain.LabelResmi, ProvenanceID: d.ProvenanceID},
		JumlahAnggota: Angka{Nilai: float64(len(d.Anggota)), Satuan: satuanAnggota, Label: domain.LabelResmi, ProvenanceID: d.ProvenanceAnggota},
		Kecamatan:     kecamatan, Anggota: anggota, ProvenanceID: d.ProvenanceID,
	}
}

func keWakilJSON(w domain.Wakil) wakilJSON {
	hasil := wakilJSON{
		ID: w.ID, Nama: w.Nama, Partai: w.Partai, Fraksi: w.Fraksi,
		DudukDiKomisiPengawas: w.IsDiKomisiPengawasTema, FotoURL: w.FotoURL,
	}
	if w.KodeKomisi != "" {
		hasil.Komisi = &komisiJSON{Kode: w.KodeKomisi, Nama: w.NamaKomisi}
	}
	return hasil
}
