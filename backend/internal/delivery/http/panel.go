package http

import (
	"context"
	"encoding/json"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

const formatTanggalMasukan = "2006-01-02"

type LayananPanel interface {
	Antrean(ctx context.Context, penggunaID int64) (usecase.AntreanPanel, error)
	TulisHakJawab(ctx context.Context, penggunaID int64, baru domain.HakJawabBaru) (domain.HakJawab, error)
	TandaiLaporanSelesai(ctx context.Context, penggunaID, laporanID int64) error
	IsiSerahTerima(ctx context.Context, penggunaID int64, s domain.SerahTerima) error
}

type hakJawabJSON struct {
	ID            int64     `json:"id"`
	TargetJenis   string    `json:"target_jenis"`
	TargetID      string    `json:"target_id"`
	Isi           string    `json:"isi"`
	TautanDokumen string    `json:"tautan_dokumen,omitempty"`
	NamaPenjawab  string    `json:"nama_penjawab"`
	Pihak         string    `json:"pihak"`
	IsSimulation  bool      `json:"is_simulation"`
	DibuatPada    time.Time `json:"dibuat_pada"`
}

type hakJawabMasukJSON struct {
	TargetJenis   string `json:"target_jenis"`
	TargetID      string `json:"target_id"`
	Isi           string `json:"isi"`
	TautanDokumen string `json:"tautan_dokumen"`
}

type serahTerimaMasukJSON struct {
	Tanggal          string `json:"tanggal"`
	NomorBeritaAcara string `json:"nomor_berita_acara"`
}

type sinyalPanelJSON struct {
	ID             int64  `json:"id"`
	Jenis          string `json:"jenis"`
	Alasan         string `json:"alasan"`
	KodePaket      string `json:"kode_paket"`
	NamaPaket      string `json:"nama_paket"`
	JumlahHakJawab int    `json:"jumlah_hak_jawab"`
}

type laporanPanelJSON struct {
	ID           int64     `json:"id"`
	Status       string    `json:"status"`
	Kondisi      string    `json:"kondisi"`
	Catatan      string    `json:"catatan"`
	KodePaket    string    `json:"kode_paket"`
	NamaPaket    string    `json:"nama_paket"`
	NamaWilayah  string    `json:"nama_wilayah"`
	IsSimulation bool      `json:"is_simulation"`
	DibuatPada   time.Time `json:"dibuat_pada"`
}

type paketPanelJSON struct {
	Kode             string           `json:"kode"`
	Nama             string           `json:"nama"`
	IsSimulation     bool             `json:"is_simulation"`
	SerahTerima      *string          `json:"tanggal_serah_terima"`
	NomorBeritaAcara string           `json:"nomor_berita_acara"`
	IsSegeraBerakhir bool             `json:"segera_berakhir"`
	Pemeliharaan     pemeliharaanJSON `json:"pemeliharaan"`
}

type umpanPanelJSON struct {
	ID             int64     `json:"id"`
	Judul          string    `json:"judul"`
	MediaNama      string    `json:"media_nama"`
	URL            string    `json:"url"`
	TerbitPada     time.Time `json:"terbit_pada"`
	JumlahHakJawab int       `json:"jumlah_hak_jawab"`
}

type antreanPanelJSON struct {
	Peran        string             `json:"peran"`
	NamaInstansi string             `json:"nama_instansi"`
	Sinyal       []sinyalPanelJSON  `json:"sinyal"`
	Laporan      []laporanPanelJSON `json:"laporan"`
	Paket        []paketPanelJSON   `json:"paket"`
	Umpan        []umpanPanelJSON   `json:"umpan"`
}

type handlerPanel struct {
	layanan LayananPanel
}

func daftarkanPanel(api fiber.Router, layanan Layanan) {
	if layanan.Panel == nil || layanan.Akses == nil {
		return
	}
	h := handlerPanel{layanan: layanan.Panel}
	masuk := wajibMasuk(layanan.Akses)
	api.Get("/panel/antrean", masuk, h.antrean)
	api.Patch("/panel/paket/:kode/serah-terima", masuk, h.serahTerima)
	api.Patch("/panel/laporan/:id/selesai", masuk, h.laporanSelesai)
	api.Post("/hak-jawab", masuk, h.hakJawab)
}

func (h handlerPanel) antrean(c fiber.Ctx) error {
	a, err := h.layanan.Antrean(c.Context(), klaimAkses(c).PenggunaID)
	if err != nil {
		return err
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return kirimData(c, keAntreanPanelJSON(a))
}

func (h handlerPanel) hakJawab(c fiber.Ctx) error {
	var masukan hakJawabMasukJSON
	if err := json.Unmarshal(c.Body(), &masukan); err != nil {
		return domain.InputTidakValid("badan permintaan wajib JSON berisi target_jenis, target_id, dan isi")
	}
	hasil, err := h.layanan.TulisHakJawab(c.Context(), klaimAkses(c).PenggunaID, domain.HakJawabBaru{
		TargetJenis: domain.JenisTargetHakJawab(masukan.TargetJenis), TargetID: masukan.TargetID,
		Isi: masukan.Isi, TautanDokumen: masukan.TautanDokumen,
	})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(Respons{Code: fiber.StatusCreated, Message: pesanOK, Data: keHakJawabJSON(hasil)})
}

func (h handlerPanel) serahTerima(c fiber.Ctx) error {
	var masukan serahTerimaMasukJSON
	if err := json.Unmarshal(c.Body(), &masukan); err != nil {
		return domain.InputTidakValid("badan permintaan wajib JSON berisi tanggal")
	}
	tanggal, err := time.Parse(formatTanggalMasukan, masukan.Tanggal)
	if err != nil {
		return domain.ErrTanggalSerahTerimaTidakSah
	}
	err = h.layanan.IsiSerahTerima(c.Context(), klaimAkses(c).PenggunaID, domain.SerahTerima{
		KodePaket: c.Params("kode"), Tanggal: tanggal, NomorBeritaAcara: masukan.NomorBeritaAcara,
	})
	if err != nil {
		return err
	}
	return kirimData(c, nil)
}

func (h handlerPanel) laporanSelesai(c fiber.Ctx) error {
	id, err := uraiIDLaporan(c)
	if err != nil {
		return err
	}
	if err := h.layanan.TandaiLaporanSelesai(c.Context(), klaimAkses(c).PenggunaID, id); err != nil {
		return err
	}
	return kirimData(c, nil)
}

func keHakJawabJSON(h domain.HakJawab) hakJawabJSON {
	return hakJawabJSON{
		ID: h.ID, TargetJenis: string(h.TargetJenis), TargetID: h.TargetID, Isi: h.Isi, TautanDokumen: h.TautanDokumen,
		NamaPenjawab: h.NamaPenjawab, Pihak: h.Pihak, IsSimulation: h.IsSimulation, DibuatPada: h.DibuatPada,
	}
}

func keDaftarHakJawabJSON(daftar []domain.HakJawab) []hakJawabJSON {
	hasil := make([]hakJawabJSON, len(daftar))
	for i, h := range daftar {
		hasil[i] = keHakJawabJSON(h)
	}
	return hasil
}

func keAntreanPanelJSON(a usecase.AntreanPanel) antreanPanelJSON {
	hasil := antreanPanelJSON{
		Peran: string(a.Pengguna.Peran), NamaInstansi: a.NamaInstansi,
		Sinyal: make([]sinyalPanelJSON, len(a.Sinyal)), Laporan: make([]laporanPanelJSON, len(a.Laporan)),
		Paket: make([]paketPanelJSON, len(a.Paket)), Umpan: make([]umpanPanelJSON, len(a.Umpan)),
	}
	for i, s := range a.Sinyal {
		hasil.Sinyal[i] = sinyalPanelJSON{ID: s.ID, Jenis: string(s.Jenis), Alasan: s.Alasan, KodePaket: s.KodePaket, NamaPaket: s.NamaPaket, JumlahHakJawab: s.JumlahHakJawab}
	}
	for i, l := range a.Laporan {
		hasil.Laporan[i] = laporanPanelJSON{
			ID: l.ID, Status: string(l.Status), Kondisi: string(l.Kondisi), Catatan: l.Catatan, KodePaket: l.KodePaket,
			NamaPaket: l.NamaPaket, NamaWilayah: l.NamaWilayah, IsSimulation: l.IsSimulation, DibuatPada: l.DibuatPada,
		}
	}
	for i, p := range a.Paket {
		hasil.Paket[i] = kePaketPanelJSON(p)
	}
	for i, u := range a.Umpan {
		hasil.Umpan[i] = umpanPanelJSON(u)
	}
	return hasil
}

func kePaketPanelJSON(p usecase.PaketPemeliharaanPanel) paketPanelJSON {
	provenance := domain.ProvenanceRumusPemeliharaan
	return paketPanelJSON{
		Kode: p.Kode, Nama: p.Nama, IsSimulation: p.IsSimulation, SerahTerima: tanggalJSON(p.Jadwal.TanggalSerahTerima),
		NomorBeritaAcara: p.NomorBeritaAcara, IsSegeraBerakhir: p.IsSegeraBerakhir,
		Pemeliharaan: kePemeliharaanJSON(usecase.HasilPemeliharaan{IsBerlaku: true, Estimasi: p.Estimasi, ProvenanceID: provenance}),
	}
}
