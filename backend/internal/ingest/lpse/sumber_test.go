package lpse

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type pengirimFixture struct {
	t *testing.T
}

func (p pengirimFixture) Ambil(_ context.Context, alamat string) ([]byte, error) {
	if filepath.Base(alamat) == "pengumumanpl" {
		return bacaFixture(p.t, "pengumuman_rt_rw.html"), nil
	}
	return bacaFixture(p.t, "halaman_nontender.html"), nil
}

func (p pengirimFixture) KirimForm(_ context.Context, alamat string, _ url.Values) ([]byte, error) {
	if filepath.Base(alamat) == "lelang?tahun=2026" {
		return bacaFixture(p.t, "daftar_tender.json"), nil
	}
	return bacaFixture(p.t, "daftar_nontender.json"), nil
}

func sumberUji(t *testing.T) Sumber {
	return Sumber{
		Klien:    pengirimFixture{t: t},
		URLDasar: "https://spse.uji/jakarta",
		Ukuran:   100,
		DirRaw:   t.TempDir(),
		Sekarang: func() time.Time { return time.Date(2026, 10, 23, 2, 0, 0, 0, zonaWIB) },
	}
}

func TestAmbilDaftarPaket_TenderDanNonTender_DigabungDenganProvenance(t *testing.T) {
	s := sumberUji(t)

	data, err := s.AmbilDaftarPaket(t.Context(), 2026)

	if err != nil || len(data.Paket) != 8 {
		t.Fatalf("paket = %d, err = %v", len(data.Paket), err)
	}
	tender := data.Paket[0]
	if !tender.IsTender || tender.TahunAnggaran != 2026 || tender.HPS == nil || tender.NilaiKontrak == nil {
		t.Errorf("tender = %+v", tender)
	}
	if data.Paket[3].NilaiKontrak == nil || data.Paket[4].NilaiKontrak != nil {
		t.Error("nilai kontrak non-tender tidak sesuai kolom daftar")
	}
	if data.Provenance.HashKonten == "" || data.Provenance.Catatan == "" {
		t.Errorf("provenance = %+v", data.Provenance)
	}
	mentah, _ := os.ReadDir(filepath.Join(s.DirRaw, "lpse_daftar", "2026-10-23"))
	if len(mentah) != 4 {
		t.Errorf("berkas mentah = %d, harap 4", len(mentah))
	}
}

func TestAmbilDetailPaket_Pengumuman_KeDomain(t *testing.T) {
	data, err := sumberUji(t).AmbilDetailPaket(t.Context(), "11039162000")

	if err != nil {
		t.Fatal(err)
	}
	d := data.Detail
	if d.RUP == nil || d.RUP.Kode != "63077471" || d.HPS == nil || *d.HPS != 197870202.32 || *d.JumlahPeserta != 1 {
		t.Errorf("detail = %+v", d)
	}
	if data.Provenance.URL != "https://spse.uji/jakarta/nontender/11039162000/pengumumanpl" {
		t.Errorf("url provenance = %s", data.Provenance.URL)
	}
}

type pengirimPencatat struct {
	pengirimFixture
	pengumuman string
	alamat     []string
}

func (p *pengirimPencatat) Ambil(ctx context.Context, alamat string) ([]byte, error) {
	p.alamat = append(p.alamat, alamat)
	if filepath.Base(alamat) == "pengumumanpl" {
		return bacaFixture(p.t, p.pengumuman), nil
	}
	return p.pengirimFixture.Ambil(ctx, alamat)
}

func TestAmbilDetailPaket_PaketBatal_HalamanPemenangTidakDiambil(t *testing.T) {
	pengirim := &pengirimPencatat{pengirimFixture: pengirimFixture{t: t}, pengumuman: "pengumuman_batal.html"}
	s := sumberUji(t)
	s.Klien = pengirim

	data, err := s.AmbilDetailPaket(t.Context(), "10403852000")

	if err != nil || data.Detail.Pemenang != nil || len(pengirim.alamat) != 1 {
		t.Errorf("pemenang = %v, permintaan = %v, err = %v", data.Detail.Pemenang, pengirim.alamat, err)
	}
}

func TestAmbilDetailPaket_TanpaPemenang_JadwalTidakDiambil(t *testing.T) {
	pengirim := &pengirimPencatat{pengirimFixture: pengirimFixture{t: t}, pengumuman: "pengumuman_konstruksi.html"}
	s := sumberUji(t)
	s.Klien = pengirim

	data, err := s.AmbilDetailPaket(t.Context(), "11047442000")

	if err != nil || data.Detail.Pemenang != nil || len(pengirim.alamat) != 3 {
		t.Errorf("pemenang = %v, permintaan = %v, err = %v", data.Detail.Pemenang, pengirim.alamat, err)
	}
}
