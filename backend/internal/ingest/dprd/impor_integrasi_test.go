package dprd_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/jackc/pgx/v5"
)

func hitung(t *testing.T, tx pgx.Tx, sql string) int {
	t.Helper()
	var n int
	if err := tx.QueryRow(t.Context(), sql).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func TestImporDPRD_FixtureSitusResmi_SepuluhDapilDan106Anggota(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	ujiwakil.MuatDPRD(t, tx)

	if n := hitung(t, tx, `SELECT count(*) FROM dapil WHERE tahun_pemilu = 2024`); n != 10 {
		t.Errorf("jumlah dapil = %d, harap 10", n)
	}
	if n := hitung(t, tx, `SELECT count(*) FROM anggota`); n != 106 {
		t.Errorf("jumlah anggota = %d, harap 106", n)
	}
}

func TestImporDPRD_FixtureSitusResmi_SetiapKecamatanTepatSatuDapil(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	ujiwakil.MuatDPRD(t, tx)

	tanpaDapil := hitung(t, tx, `SELECT count(*) FROM wilayah w
WHERE w.tingkat = 'kecamatan' AND (SELECT count(*) FROM dapil_wilayah dw WHERE dw.kode_kecamatan = w.kode) <> 1`)
	if tanpaDapil != 0 {
		t.Errorf("%d kecamatan tidak masuk tepat satu dapil", tanpaDapil)
	}
}

func TestImporDPRD_DijalankanDuaKali_TidakMenggandakanData(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	ujiwakil.MuatDPRD(t, tx)
	ujiwakil.MuatDPRD(t, tx)

	if n := hitung(t, tx, `SELECT count(*) FROM anggota`); n != 106 {
		t.Errorf("jumlah anggota setelah dua kali impor = %d, harap 106", n)
	}
	if n := hitung(t, tx, `SELECT count(*) FROM dapil_wilayah`); n != 44 {
		t.Errorf("jumlah kecamatan dapil = %d, harap 44", n)
	}
}

func TestImporDPRD_AnggotaPPPdiFraksiNasDem_FraksiDariHalamanFraksi(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	ujiwakil.MuatDPRD(t, tx)

	var partai, fraksi, url string
	err := tx.QueryRow(t.Context(), `SELECT partai, fraksi, url_profil FROM anggota WHERE id = 'h-matnoor-tindoan-sh'`).Scan(&partai, &fraksi, &url)
	if err != nil {
		t.Fatal(err)
	}
	if partai != "Partai Persatuan Pembangunan" || fraksi != "Fraksi Partai NasDem" || url != "https://dprd-dkijakartaprov.go.id/profile/h-matnoor-tindoan-sh/" {
		t.Errorf("dapat partai %q fraksi %q url %q", partai, fraksi, url)
	}
}

func TestImporDPRD_FixtureKomisi_SetiapKomisiPunyaBidangDanAnggota(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	laporan := ujiwakil.MuatDPRD(t, tx)

	if n := hitung(t, tx, `SELECT count(*) FROM komisi WHERE cardinality(bidang) > 0`); n != 5 {
		t.Errorf("komisi berbidang = %d, harap 5", n)
	}
	if n := hitung(t, tx, `SELECT count(DISTINCT komisi_kode) FROM anggota`); n != 5 {
		t.Errorf("komisi yang punya anggota = %d, harap 5", n)
	}
	if n := hitung(t, tx, `SELECT count(*) FROM anggota WHERE komisi_kode IS NULL`); n != 5 {
		t.Errorf("anggota tanpa komisi = %d, harap 5 pimpinan DPRD", n)
	}
	if n := hitung(t, tx, `SELECT count(*) FROM anggota WHERE id = ''`); n != 0 {
		t.Errorf("%d anggota tanpa ID", n)
	}
	if !slices.ContainsFunc(laporan, func(p string) bool { return strings.Contains(p, "anggota tanpa data komisi (5)") }) {
		t.Errorf("harap anggota tanpa komisi dilaporkan, laporan %q", laporan)
	}
}

func TestImporDPRD_PimpinanDPRD_KoordinatorTanpaKomisiTetapiPunyaJabatan(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	ujiwakil.MuatDPRD(t, tx)

	var komisi *string
	if err := tx.QueryRow(t.Context(), `SELECT komisi_kode FROM anggota WHERE id = 'ima-mahdiah'`).Scan(&komisi); err != nil {
		t.Fatal(err)
	}
	jabatan := hitung(t, tx, `SELECT count(*) FROM anggota_akd WHERE anggota_id = 'ima-mahdiah'
AND (alat_kelengkapan, jabatan) IN (('Komisi A Bidang Pemerintahan', 'Koordinator'), ('Pimpinan DPRD', 'Wakil Ketua DPRD Provinsi DKI Jakarta'))`)
	if komisi != nil || jabatan != 2 {
		t.Errorf("komisi = %v, jabatan cocok = %d, harap tanpa komisi dan 2 jabatan", komisi, jabatan)
	}
}
