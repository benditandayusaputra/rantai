package http

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
	"github.com/gofiber/fiber/v3"
)

const (
	tipeAliranPeristiwa = "text/event-stream"
	tanpaCache          = "no-cache"
	tanpaPenyanggaProxy = "no"
	peristiwaAlat       = "alat"
	peristiwaJawaban    = "jawaban"
	peristiwaGalat      = "galat"
	pesanGagalTanya     = "Tanya Rantai belum bisa menjawab sekarang. Coba lagi sebentar lagi atau buka halaman kecamatan."
)

type LayananTanya interface {
	Periksa(p tanya.Permintaan) error
	Jawab(ctx context.Context, p tanya.Permintaan, kabar func(tanya.Kabar)) (tanya.Jawaban, error)
}

type permintaanTanyaJSON struct {
	Pertanyaan           string `json:"pertanyaan"`
	KodeKecamatanKonteks string `json:"kode_kecamatan_konteks"`
	Bahasa               string `json:"bahasa"`
}

type kabarJSON struct {
	Alat string `json:"alat"`
}

type rujukanJSON struct {
	Nomor        int    `json:"nomor"`
	ProvenanceID string `json:"provenance_id"`
}

type jawabanTanyaJSON struct {
	Teks    string        `json:"teks"`
	Rujukan []rujukanJSON `json:"rujukan"`
}

type galatJSON struct {
	Pesan string `json:"pesan"`
}

type handlerTanya struct {
	layanan    LayananTanya
	batasWaktu time.Duration
}

func daftarkanTanya(api fiber.Router, layanan LayananTanya, cfg config.ServerTanya, rahasiaProxy string) {
	if layanan == nil {
		return
	}
	batas := pembatasLaju(cfg.BatasLaju, rahasiaProxy)
	api.Post("/tanya", batas, handlerTanya{layanan: layanan, batasWaktu: cfg.BatasWaktu}.tanya)
}

func (h handlerTanya) tanya(c fiber.Ctx) error {
	var masukan permintaanTanyaJSON
	if err := json.Unmarshal(c.Body(), &masukan); err != nil {
		return domain.InputTidakValid("badan permintaan wajib JSON berisi pertanyaan")
	}
	permintaan := tanya.Permintaan(masukan)
	if err := h.layanan.Periksa(permintaan); err != nil {
		return err
	}
	c.Set(fiber.HeaderContentType, tipeAliranPeristiwa)
	c.Set(fiber.HeaderCacheControl, tanpaCache)
	c.Set("X-Accel-Buffering", tanpaPenyanggaProxy)
	return c.SendStreamWriter(func(w *bufio.Writer) {
		if err := h.alirkan(w, permintaan); err != nil {
			log.Printf("alirkan jawaban tanya: %v", err)
		}
	})
}

func (h handlerTanya) alirkan(w *bufio.Writer, permintaan tanya.Permintaan) error {
	ctx, batal := context.WithTimeout(context.Background(), h.batasWaktu)
	defer batal()
	var errKirim error
	jawaban, err := h.layanan.Jawab(ctx, permintaan, func(k tanya.Kabar) {
		if errKirim == nil {
			errKirim = kirimPeristiwa(w, peristiwaAlat, kabarJSON(k))
		}
	})
	if err != nil {
		return errors.Join(fmt.Errorf("jawab pertanyaan: %w", err), kirimPeristiwa(w, peristiwaGalat, galatJSON{Pesan: pesanGagalTanya}))
	}
	if errKirim != nil {
		return errKirim
	}
	return kirimPeristiwa(w, peristiwaJawaban, keJawabanTanyaJSON(jawaban))
}

func keJawabanTanyaJSON(j tanya.Jawaban) jawabanTanyaJSON {
	rujukan := make([]rujukanJSON, len(j.Rujukan))
	for i, r := range j.Rujukan {
		rujukan[i] = rujukanJSON(r)
	}
	return jawabanTanyaJSON{Teks: j.Teks, Rujukan: rujukan}
}

func kirimPeristiwa(w *bufio.Writer, nama string, data any) error {
	isi, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("serialisasi peristiwa %s: %w", nama, err)
	}
	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", nama, isi); err != nil {
		return fmt.Errorf("tulis peristiwa %s: %w", nama, err)
	}
	return w.Flush()
}
