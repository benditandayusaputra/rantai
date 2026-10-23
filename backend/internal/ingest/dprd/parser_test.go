package dprd

import (
	"os"
	"slices"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

func bacaFixture(t *testing.T, nama string) []byte {
	t.Helper()
	isi, err := os.ReadFile("testdata/" + nama)
	if err != nil {
		t.Fatal(err)
	}
	return isi
}

func TestUraiHalamanDapil_Dapil1_DuaBelasAnggotaDenganPartai(t *testing.T) {
	daftar, err := UraiHalamanDapil(bacaFixture(t, "dapil_dki-jakarta-1.html"), 1)
	if err != nil {
		t.Fatal(err)
	}

	if len(daftar) != 12 {
		t.Fatalf("harap 12 anggota, dapat %d", len(daftar))
	}
	harap := domain.BarisDapil{NomorDapil: 1, Nama: "Muhammad Hasan Abdilah", Partai: "Partai Keadilan Sejahtera", Slug: "muhammad-hasan-abdillah"}
	if daftar[0] != harap {
		t.Errorf("baris pertama = %+v, harap %+v", daftar[0], harap)
	}
}

func TestUraiHalamanStruktur_FraksiNasDem_MemuatAnggotaPPP(t *testing.T) {
	daftar, err := UraiHalamanStruktur(bacaFixture(t, "fraksi_partai-nasdem.html"))
	if err != nil {
		t.Fatal(err)
	}

	if len(daftar) != 11 {
		t.Fatalf("harap 11 anggota fraksi, dapat %d", len(daftar))
	}
	harap := domain.BarisStruktur{Nama: "Matnoor Tindoan", Jabatan: "Anggota", Slug: "h-matnoor-tindoan-sh"}
	if !slices.Contains(daftar, harap) {
		t.Errorf("harap %+v ada di daftar fraksi", harap)
	}
}

func TestUraiBidangKomisi_KomisiB_NamaDinasBerkomaTetapUtuh(t *testing.T) {
	kataAwal := []string{"Dinas", "Badan", "Biro", "Perusahaan"}

	bidang, err := UraiBidangKomisi(bacaFixture(t, "komisi_komisi-b-bidang-perekonomian.html"), kataAwal)
	if err != nil {
		t.Fatal(err)
	}

	harap := []string{
		"Dinas Pangan, Pertanian dan Kelautan",
		"Dinas Perindustrian, Perdagangan, Koperasi, Usaha Kecil dan Menengah",
		"Dinas Pariwisata dan Ekonomi Kreatif",
		"Dinas Penanaman Modal dan Pelayanan Terpadu Satu Pintu",
		"Dinas Perhubungan",
		"Dinas Tenaga Kerja, Transmigrasi dan Energi",
		"Badan Pembinaan Badan Usaha Milik Daerah (BUMD)",
		"Perusahaan Daerah",
	}
	if !slices.Equal(bidang, harap) {
		t.Errorf("bidang = %q", bidang)
	}
}
