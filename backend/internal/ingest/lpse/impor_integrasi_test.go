package lpse

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const barisNonTenderUji = `["%s","%s","Provinsi DKI Jakarta","%s","185,1 Jt","Pengadaan Langsung","Pekerjaan Konstruksi - TA 2026","5","Nilai Kontrak belum dibuat","1","0","0"]`

var fixturePengumuman = map[string]string{
	"11047442000": "pengumuman_konstruksi.html",
	"11039162000": "pengumuman_rt_rw.html",
}

type pengirimImpor struct {
	t            *testing.T
	tahap        string
	kodeTambahan []string
}

func (p *pengirimImpor) Ambil(_ context.Context, alamat string) ([]byte, error) {
	if strings.HasSuffix(alamat, "/pengumumanpl") {
		return bacaFixture(p.t, fixturePengumuman[filepath.Base(filepath.Dir(alamat))]), nil
	}
	return bacaFixture(p.t, "halaman_nontender.html"), nil
}

func (p *pengirimImpor) KirimForm(_ context.Context, alamat string, _ url.Values) ([]byte, error) {
	if strings.Contains(alamat, "/dt/lelang") {
		return bacaFixture(p.t, "daftar_tender.json"), nil
	}
	baris := []string{}
	for kode := range fixturePengumuman {
		baris = append(baris, fmt.Sprintf(barisNonTenderUji, kode, "Paket "+kode, p.tahap))
	}
	for _, kode := range p.kodeTambahan {
		baris = append(baris, fmt.Sprintf(barisNonTenderUji, kode, "Paket non-tender "+kode, p.tahap))
	}
	return []byte(`{"data":[` + strings.Join(baris, ",") + `]}`), nil
}

func siapkanImpor(t *testing.T) (pgx.Tx, *pengirimImpor, usecase.ImporDaftarPaket, usecase.ImporDetailPaket) {
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	pengirim := &pengirimImpor{t: t, tahap: "Paket Sudah Selesai"}
	sumber := Sumber{Klien: pengirim, URLDasar: "https://spse.uji/jakarta", Ukuran: 100, DirRaw: t.TempDir(), Sekarang: time.Now}
	repo := repository.PaketRepository{DB: tx}
	pencatat := repository.IngestRunRepository{DB: tx}
	daftar := usecase.ImporDaftarPaket{Pengambil: sumber, Penyimpan: repo, Pencatat: pencatat, Tahun: []int{2026}}
	detail := usecase.ImporDetailPaket{
		Pengambil: sumber, Penyimpan: repo, Pencatat: pencatat,
		JenisPengadaan: []string{"Pekerjaan Konstruksi"},
		AturanInstansi: domain.AturanInstansi{
			Jenis: []domain.AturanJenisInstansi{{Jenis: domain.JenisInstansiSudin, Awalan: []string{"suku dinas"}}},
			Kota:  map[string][]string{"31.71": {"jakarta pusat"}, "31.75": {"jaktim"}},
		},
	}
	return tx, pengirim, daftar, detail
}

func hitungBaris(t *testing.T, tx pgx.Tx, sql string) int {
	t.Helper()
	var jumlah int
	if err := tx.QueryRow(t.Context(), sql).Scan(&jumlah); err != nil {
		t.Fatal(err)
	}
	return jumlah
}

func TestImporDaftarPaket_DuaKali_TidakMenggandakanDanRunTercatat(t *testing.T) {
	tx, _, daftar, _ := siapkanImpor(t)

	for range 2 {
		if err := daftar.Jalankan(t.Context()); err != nil {
			t.Fatal(err)
		}
	}

	if n := hitungBaris(t, tx, "SELECT count(*) FROM paket"); n != 5 {
		t.Errorf("paket = %d, harap 5", n)
	}
	if n := hitungBaris(t, tx, "SELECT count(*) FROM paket_riwayat_tahap"); n != 5 {
		t.Errorf("riwayat tahap = %d, harap 5", n)
	}
	if n := hitungBaris(t, tx, "SELECT count(*) FROM ingest_run WHERE sumber = 'lpse_daftar' AND status = 'berhasil'"); n != 2 {
		t.Errorf("ingest_run = %d, harap 2", n)
	}
}

func TestImporDaftarPaket_TahapBerubah_RiwayatTahapBertambah(t *testing.T) {
	tx, pengirim, daftar, _ := siapkanImpor(t)
	if err := daftar.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	pengirim.tahap = "Paket Dibatalkan"
	if err := daftar.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	n := hitungBaris(t, tx, "SELECT count(*) FROM paket_riwayat_tahap WHERE tahap_lama = 'Paket Sudah Selesai' AND tahap_baru = 'Paket Dibatalkan'")
	if n != 2 {
		t.Errorf("perubahan tahap tercatat = %d, harap 2", n)
	}
}

func TestImporDetailPaket_SetelahDaftar_SatkerRUPDanHPSPresisi(t *testing.T) {
	tx, _, daftar, detail := siapkanImpor(t)
	for _, job := range []func(context.Context) error{daftar.Jalankan, detail.Jalankan, daftar.Jalankan, detail.Jalankan} {
		if err := job(t.Context()); err != nil {
			t.Fatal(err)
		}
	}

	var satker, kota, kodeRUP string
	var hps float64
	err := tx.QueryRow(t.Context(), `SELECT i.nama, i.kode_kota, p.kode_rup, p.hps FROM paket p JOIN instansi i ON i.id = p.instansi_id
WHERE p.kode = 'nontender-11047442000'`).Scan(&satker, &kota, &kodeRUP, &hps)
	if err != nil {
		t.Fatal(err)
	}
	if kota != "31.71" || kodeRUP != "63728635" || hps != 185101703.01 {
		t.Errorf("satker %q kota %q rup %q hps %v", satker, kota, kodeRUP, hps)
	}
	if n := hitungBaris(t, tx, "SELECT count(*) FROM provenance WHERE url LIKE '%/pengumumanpl'"); n != 2 {
		t.Errorf("detail diambil = %d kali, harap 2 (tidak diambil ulang bila tahap tetap)", n)
	}
}

func TestImporDaftarPaket_KodeSPSESamaTenderDanNonTender_KeduanyaTersimpan(t *testing.T) {
	tx, pengirim, daftar, _ := siapkanImpor(t)
	pengirim.kodeTambahan = []string{"52767127"}

	if err := daftar.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	n := hitungBaris(t, tx, "SELECT count(*) FROM paket WHERE kode_spse = '52767127' AND kode IN ('tender-52767127', 'nontender-52767127')")
	if n != 2 {
		t.Errorf("paket berkode SPSE sama = %d, harap 2", n)
	}
}
