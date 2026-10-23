package lpse

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func bacaFixture(t *testing.T, nama string) []byte {
	t.Helper()
	isi, err := os.ReadFile(filepath.Join("testdata", nama))
	if err != nil {
		t.Fatal(err)
	}
	return isi
}

func uraiPengumumanFixture(t *testing.T, nama string) Pengumuman {
	t.Helper()
	p, err := UraiPengumuman(bacaFixture(t, nama))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestUraiDaftar_Tender_KolomTerpetakan(t *testing.T) {
	daftar, err := UraiDaftar(bacaFixture(t, "daftar_tender.json"), JenisTender)

	if err != nil || len(daftar) != 3 {
		t.Fatalf("daftar = %d baris, err = %v", len(daftar), err)
	}
	b := daftar[0]
	if b.Kode != "52767127" || b.Metode != "Tender" || b.JenisPengadaan != "Pekerjaan Konstruksi" || !b.IsTender {
		t.Errorf("identitas = %+v", b)
	}
	if b.HPSRingkas != 660.9e9 || b.NilaiKontrak != 620772529712.03 || !b.HasNilaiKontrak {
		t.Errorf("nilai: hps %v, kontrak %v", b.HPSRingkas, b.NilaiKontrak)
	}
	if !slices.Equal(b.TahunAnggaran, []int{2022, 2023, 2024, 2025, 2026}) {
		t.Errorf("tahun anggaran = %v", b.TahunAnggaran)
	}
}

func TestUraiDaftar_NonTenderKontrakBelumDibuat_TanpaNilaiKontrak(t *testing.T) {
	daftar, err := UraiDaftar(bacaFixture(t, "daftar_nontender.json"), JenisNonTender)

	if err != nil || len(daftar) != 5 {
		t.Fatalf("daftar = %d baris, err = %v", len(daftar), err)
	}
	if !daftar[0].HasNilaiKontrak || daftar[1].HasNilaiKontrak {
		t.Errorf("kontrak baris 1 = %v, baris 2 = %v", daftar[0].HasNilaiKontrak, daftar[1].HasNilaiKontrak)
	}
	if daftar[0].Metode != "Pengadaan Langsung" || daftar[0].HPSRingkas != 71.8e6 || daftar[0].IsTender {
		t.Errorf("baris 1 = %+v", daftar[0])
	}
}

func TestUraiDaftar_NamaBerlencana_LencanaDibuang(t *testing.T) {
	isi := []byte(`{"data":[["1","Pembangunan Pintu Air Kel. Semper Barat <span class=\"badge badge-warning\">Paket Batal</span>","Provinsi DKI Jakarta","Paket Sudah Selesai","129,7 Jt","Pengadaan Langsung","Pekerjaan Konstruksi - TA 2025","5","Nilai Kontrak belum dibuat","1","0","0"]]}`)

	daftar, err := UraiDaftar(isi, JenisNonTender)

	if err != nil || daftar[0].Nama != "Pembangunan Pintu Air Kel. Semper Barat" {
		t.Errorf("nama = %q, err = %v", daftar[0].Nama, err)
	}
}

func TestUraiToken_HalamanDaftar_TokenDitemukan(t *testing.T) {
	token, err := UraiToken(bacaFixture(t, "halaman_nontender.html"))

	if err != nil || len(token) != 40 {
		t.Errorf("token = %q, err = %v", token, err)
	}
}

func TestUraiPengumuman_NonTenderKonstruksi_SemuaFieldTerisi(t *testing.T) {
	p := uraiPengumumanFixture(t, "pengumuman_konstruksi.html")

	if p.Kode != "11047442000" || p.SatuanKerja != "SUKU DINAS SUMBER DAYA AIR KOTA ADMINISTRASI JAKARTA PUSAT" {
		t.Errorf("identitas = %q, %q", p.Kode, p.SatuanKerja)
	}
	if len(p.RUP) != 1 || p.RUP[0].Kode != "63728635" || p.RUP[0].SumberDana != "APBD" {
		t.Errorf("rup = %+v", p.RUP)
	}
	if p.Pagu != 185101704 || p.HPS != 185101703.01 || p.TahunAnggaran != 2026 || p.JumlahPeserta != 1 {
		t.Errorf("nilai: pagu %v, hps %v, ta %d, peserta %d", p.Pagu, p.HPS, p.TahunAnggaran, p.JumlahPeserta)
	}
	if p.JenisPengadaan != "Pekerjaan Konstruksi" || p.Metode != "Pengadaan Langsung" || p.JenisKontrak != "Harga Satuan" {
		t.Errorf("jenis = %q, %q, %q", p.JenisPengadaan, p.Metode, p.JenisKontrak)
	}
	if !p.TanggalDibuat.Equal(time.Date(2026, time.September, 9, 0, 0, 0, 0, zonaWIB)) || p.IsBatal {
		t.Errorf("tanggal = %v, batal = %v", p.TanggalDibuat, p.IsBatal)
	}
}

func TestUraiPengumuman_NamaDenganRTRW_NamaUtuh(t *testing.T) {
	p := uraiPengumumanFixture(t, "pengumuman_rt_rw.html")

	harap := "Belanja Jasa yang Diberikan kepada Masyarakat Untuk Perbaikan Fasilitas Olahraga Futsal RW 9 Kel. Ujung Menteng, Kec. Cakung"
	if p.Nama != harap || p.SatuanKerja != "SUKU DINAS PEMUDA DAN OLAHRAGA KOTA - JAKTIM" {
		t.Errorf("nama = %q, satker = %q", p.Nama, p.SatuanKerja)
	}
	if !slices.Equal(p.Lokasi, []string{"Jakarta Timur - Jakarta Timur (Kota)"}) {
		t.Errorf("lokasi = %v", p.Lokasi)
	}
}

func TestUraiPengumuman_PaketBatal_LencanaDanAlasanTerbaca(t *testing.T) {
	p := uraiPengumumanFixture(t, "pengumuman_batal.html")

	if !p.IsBatal || p.Nama != "Pembangunan Pintu Air Kp. Kandang RT 08 RW 04, Kel. Semper Barat, Kec. Cilincing" {
		t.Errorf("batal = %v, nama = %q", p.IsBatal, p.Nama)
	}
	if p.AlasanBatal == "" || p.TahunAnggaran != 2025 {
		t.Errorf("alasan = %q, ta = %d", p.AlasanBatal, p.TahunAnggaran)
	}
}

func TestUraiPengumuman_PaketBerkontrak_GabunganJenisKontrak(t *testing.T) {
	p := uraiPengumumanFixture(t, "pengumuman_berkontrak.html")

	if p.Kode != "10156350000" || p.JenisKontrak != "Gabungan Lumsum dan Harga Satuan" || p.RUP[0].Kode != "57315341" {
		t.Errorf("pengumuman = %+v", p)
	}
}

func TestUraiPengumuman_TanpaTabelRUP_RUPKosong(t *testing.T) {
	isi := []byte(`<table><tr><th class="bgwarning">Kode Paket</th><td>9</td></tr>
<tr><th class="bgwarning">Tanggal Pembuatan</th><td>1 Januari 2026</td></tr></table>`)

	p, err := UraiPengumuman(isi)

	if err != nil || p.RUP == nil || len(p.RUP) != 0 {
		t.Errorf("rup = %v, err = %v", p.RUP, err)
	}
}

func TestUraiPengumuman_HalamanLain_ErrBukanPengumuman(t *testing.T) {
	_, err := UraiPengumuman(bacaFixture(t, "halaman_nontender.html"))

	if !errors.Is(err, ErrBukanPengumuman) {
		t.Errorf("err = %v", err)
	}
}
