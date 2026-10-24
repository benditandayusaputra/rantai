package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var (
	hariIniPanel = time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	dinasPanel   = int64(3)
	petugasPanel = domain.Pengguna{ID: 1, Peran: domain.PeranInstansi, InstansiID: &dinasPanel}
	instansiLain = int64(9)
	petugasLain  = domain.Pengguna{ID: 2, Peran: domain.PeranInstansi, InstansiID: &instansiLain}
	wargaPanel   = domain.Pengguna{ID: 4, Peran: domain.PeranWarga}
	pemilikDinas = domain.PemilikTarget{InstansiID: &dinasPanel, NamaInstansi: "Dinas Uji"}
)

type panelPalsu struct {
	pengguna    map[int64]domain.Pengguna
	pemilik     domain.PemilikTarget
	hakJawab    []HakJawabTersimpan
	serahTerima []SerahTerimaTersimpan
	perubahan   []domain.PerubahanStatus
	antrean     domain.AntreanInstansi
}

func (p *panelPalsu) PastikanPengguna(context.Context, string) (domain.Pengguna, error) {
	return domain.Pengguna{}, nil
}

func (p *panelPalsu) AmbilPengguna(_ context.Context, id int64) (domain.Pengguna, error) {
	return p.pengguna[id], nil
}

func (p *panelPalsu) PemilikTarget(context.Context, domain.JenisTargetHakJawab, string) (domain.PemilikTarget, error) {
	return p.pemilik, nil
}

func (p *panelPalsu) SimpanHakJawab(_ context.Context, h HakJawabTersimpan) (domain.HakJawab, error) {
	p.hakJawab = append(p.hakJawab, h)
	return domain.HakJawab{ID: int64(len(p.hakJawab)), Isi: h.Isi}, nil
}

func (p *panelPalsu) SimpanSerahTerima(_ context.Context, s SerahTerimaTersimpan) error {
	p.serahTerima = append(p.serahTerima, s)
	return nil
}

func (p *panelPalsu) AntreanInstansi(context.Context, FilterPaketPanel) (domain.AntreanInstansi, error) {
	return p.antrean, nil
}

func (p *panelPalsu) AntreanAnggota(context.Context, string) ([]domain.UmpanPanel, error) {
	return []domain.UmpanPanel{{ID: 5}}, nil
}

func (p *panelPalsu) UbahStatus(_ context.Context, perubahan domain.PerubahanStatus) error {
	p.perubahan = append(p.perubahan, perubahan)
	return nil
}

func panelUji(pemilik domain.PemilikTarget) (Panel, *panelPalsu) {
	palsu := &panelPalsu{
		pengguna: map[int64]domain.Pengguna{petugasPanel.ID: petugasPanel, petugasLain.ID: petugasLain, wargaPanel.ID: wargaPanel},
		pemilik:  pemilik,
	}
	return Panel{
		Pengguna: palsu, Penyimpan: palsu, Laporan: palsu,
		Pemeliharaan: KonfigurasiPemeliharaan{
			AturanPemeliharaan: domain.AturanPemeliharaan{
				JenisMasaHari: map[string]int{"permanen": 180}, Bawaan: domain.MasaKategori{Jenis: "permanen", DurasiBawaanHari: 90},
			},
			AmbangPanelHari: 30,
		},
		BuatProvenance: func(c CatatanProvenance) domain.Provenance {
			return domain.Provenance{ID: "prv_uji", Sumber: c.Sumber, Label: c.Label, Catatan: c.Catatan}
		},
		Sekarang: func() time.Time { return hariIniPanel },
	}, palsu
}

func hakJawabUji(jenis domain.JenisTargetHakJawab, id string) domain.HakJawabBaru {
	return domain.HakJawabBaru{TargetJenis: jenis, TargetID: id, Isi: "  Pekerjaan sudah diperbaiki.  "}
}

func TestPanel_Antrean_WargaDitolak(t *testing.T) {
	u, _ := panelUji(pemilikDinas)

	_, err := u.Antrean(t.Context(), wargaPanel.ID)

	if !errors.Is(err, domain.ErrBukanAkunPanel) {
		t.Errorf("err %v", err)
	}
}

func TestPanel_AntreanInstansi_PaketBerjalanTerurutDanSegeraBerakhirDitandai(t *testing.T) {
	u, palsu := panelUji(pemilikDinas)
	selesai := func(hariLalu int) *time.Time { t := hariIniPanel.AddDate(0, 0, -hariLalu); return &t }
	palsu.antrean.Paket = []domain.PaketJadwalPanel{
		{Kode: "jauh", Jadwal: domain.JadwalPaket{JadwalSelesai: selesai(10)}},
		{Kode: "segera", Jadwal: domain.JadwalPaket{JadwalSelesai: selesai(170)}},
		{Kode: "berakhir", Jadwal: domain.JadwalPaket{JadwalSelesai: selesai(400)}},
		{Kode: "tanpa-jadwal"},
	}

	hasil, err := u.Antrean(t.Context(), petugasPanel.ID)

	if err != nil || len(hasil.Paket) != 2 || hasil.Paket[0].Kode != "segera" || !hasil.Paket[0].IsSegeraBerakhir || hasil.Paket[1].IsSegeraBerakhir {
		t.Errorf("paket %+v, err %v", hasil.Paket, err)
	}
}

func TestPanel_TulisHakJawab_PetugasPemilikTersimpanRapi(t *testing.T) {
	u, palsu := panelUji(pemilikDinas)

	_, err := u.TulisHakJawab(t.Context(), petugasPanel.ID, hakJawabUji(domain.TargetHakJawabSinyal, "12"))

	if err != nil || len(palsu.hakJawab) != 1 || palsu.hakJawab[0].Isi != "Pekerjaan sudah diperbaiki." || len(palsu.perubahan) != 0 {
		t.Errorf("tersimpan %+v, perubahan %+v, err %v", palsu.hakJawab, palsu.perubahan, err)
	}
}

func TestPanel_TulisHakJawab_InstansiLainDitolakTanpaMenyimpan(t *testing.T) {
	u, palsu := panelUji(pemilikDinas)

	_, err := u.TulisHakJawab(t.Context(), petugasLain.ID, hakJawabUji(domain.TargetHakJawabPaket, "tender-1"))

	if !errors.Is(err, domain.ErrBukanPenjawabTarget) || len(palsu.hakJawab) != 0 {
		t.Errorf("err %v, tersimpan %d", err, len(palsu.hakJawab))
	}
}

func TestPanel_TulisHakJawab_LaporanDiteruskanMenjadiDitanggapi(t *testing.T) {
	pemilik := pemilikDinas
	pemilik.StatusLaporan = domain.StatusDiteruskan
	u, palsu := panelUji(pemilik)

	_, err := u.TulisHakJawab(t.Context(), petugasPanel.ID, hakJawabUji(domain.TargetHakJawabLaporan, "7"))

	if err != nil || len(palsu.perubahan) != 1 || palsu.perubahan[0].LaporanID != 7 || palsu.perubahan[0].Ke != domain.StatusDitanggapi {
		t.Errorf("perubahan %+v, err %v", palsu.perubahan, err)
	}
}

func TestPanel_TandaiLaporanSelesai(t *testing.T) {
	kasus := []struct {
		nama  string
		dari  domain.StatusLaporan
		harap error
	}{
		{"sudah ditanggapi", domain.StatusDitanggapi, nil},
		{"belum ditanggapi", domain.StatusDiteruskan, domain.ErrTransisiTidakSah},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			pemilik := pemilikDinas
			pemilik.StatusLaporan = k.dari
			u, palsu := panelUji(pemilik)

			err := u.TandaiLaporanSelesai(t.Context(), petugasPanel.ID, 7)

			if !errors.Is(err, k.harap) || (err == nil && palsu.perubahan[0].Ke != domain.StatusSelesai) {
				t.Errorf("err %v, perubahan %+v", err, palsu.perubahan)
			}
		})
	}
}

func TestPanel_IsiSerahTerima_ProvenanceResmiDariInstansi(t *testing.T) {
	u, palsu := panelUji(pemilikDinas)
	s := domain.SerahTerima{KodePaket: "tender-1", Tanggal: hariIniPanel.AddDate(0, -1, 0), NomorBeritaAcara: "BA-12/2026"}

	err := u.IsiSerahTerima(t.Context(), petugasPanel.ID, s)

	if err != nil || len(palsu.serahTerima) != 1 {
		t.Fatalf("err %v", err)
	}
	prov := palsu.serahTerima[0].Provenance
	if prov.Label != domain.LabelResmi || prov.Sumber != domain.SumberInstansi || !strings.Contains(prov.Catatan, "Dinas Uji") || !strings.Contains(prov.Catatan, "BA-12/2026") {
		t.Errorf("provenance %+v", prov)
	}
}

func TestPanel_IsiSerahTerima_TanggalMasaDepanDitolak(t *testing.T) {
	u, palsu := panelUji(pemilikDinas)

	err := u.IsiSerahTerima(t.Context(), petugasPanel.ID, domain.SerahTerima{KodePaket: "tender-1", Tanggal: hariIniPanel.AddDate(0, 0, 2)})

	if !errors.Is(err, domain.ErrTanggalSerahTerimaTidakSah) || len(palsu.serahTerima) != 0 {
		t.Errorf("err %v", err)
	}
}
