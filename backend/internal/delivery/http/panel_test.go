package http

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/token"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
)

type panelPalsu struct {
	serahTerima domain.SerahTerima
	penjawab    int64
}

func (p *panelPalsu) Antrean(context.Context, int64) (usecase.AntreanPanel, error) {
	return usecase.AntreanPanel{}, nil
}

func (p *panelPalsu) TulisHakJawab(_ context.Context, id int64, b domain.HakJawabBaru) (domain.HakJawab, error) {
	p.penjawab = id
	return domain.HakJawab{ID: 1, TargetJenis: b.TargetJenis, TargetID: b.TargetID, Isi: b.Isi}, nil
}

func (p *panelPalsu) TandaiLaporanSelesai(context.Context, int64, int64) error { return nil }

func (p *panelPalsu) IsiSerahTerima(_ context.Context, _ int64, s domain.SerahTerima) error {
	p.serahTerima = s
	return nil
}

func appUjiPanel(t *testing.T) (*fiber.App, *panelPalsu, string) {
	t.Helper()
	akses, err := token.Baru(strings.Repeat("r", token.PanjangMinRahasia), time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	tokenSah, err := akses.Terbitkan(domain.KlaimAkses{PenggunaID: 11, Peran: domain.PeranInstansi})
	if err != nil {
		t.Fatal(err)
	}
	panel := &panelPalsu{}
	return NewApp(config.Server{}, Layanan{Panel: panel, Akses: akses}), panel, tokenSah
}

func TestPanel_TanpaToken_401(t *testing.T) {
	app, _, _ := appUjiPanel(t)

	hasil := kirimAuth(t, app, fiber.MethodGet, "/api/v1/panel/antrean", "", tanpaUbah)

	if hasil.kode != fiber.StatusUnauthorized {
		t.Errorf("kode %d", hasil.kode)
	}
}

func TestPanel_HakJawab_201DenganPenjawabDariToken(t *testing.T) {
	app, panel, tokenSah := appUjiPanel(t)

	hasil := kirimAuth(t, app, fiber.MethodPost, "/api/v1/hak-jawab", `{"target_jenis":"paket","target_id":"tender-1","isi":"Tanggapan."}`,
		func(h map[string]string) { h[fiber.HeaderAuthorization] = "Bearer " + tokenSah })

	if hasil.kode != fiber.StatusCreated || panel.penjawab != 11 {
		t.Errorf("kode %d, penjawab %d", hasil.kode, panel.penjawab)
	}
}

func TestPanel_SerahTerima_TanggalISO(t *testing.T) {
	kasus := []struct {
		nama, tanggal string
		harap         int
	}{
		{"tanggal sah", "2026-07-03", fiber.StatusOK},
		{"format lain", "03/07/2026", fiber.StatusBadRequest},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			app, panel, tokenSah := appUjiPanel(t)

			hasil := kirimAuth(t, app, fiber.MethodPatch, "/api/v1/panel/paket/tender-1/serah-terima", `{"tanggal":"`+k.tanggal+`"}`,
				func(h map[string]string) { h[fiber.HeaderAuthorization] = "Bearer " + tokenSah })

			if hasil.kode != k.harap || (k.harap == fiber.StatusOK && panel.serahTerima.KodePaket != "tender-1") {
				t.Errorf("kode %d, serah terima %+v", hasil.kode, panel.serahTerima)
			}
		})
	}
}
