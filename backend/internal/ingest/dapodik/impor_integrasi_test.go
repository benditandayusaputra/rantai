package dapodik

import (
	"context"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

type pengambilFixture struct {
	t *testing.T
}

func (p pengambilFixture) Ambil(context.Context) (usecase.DataSumberSekolah, error) {
	daftar := []domain.DataSekolah{}
	for _, nama := range []string{"full_detail.json", "full_detail_tanpa_koordinat.json", "full_detail_madrasah.json"} {
		detail, err := UraiDetail(bacaFixture(p.t, nama))
		if err != nil {
			return usecase.DataSumberSekolah{}, err
		}
		rombel := []domain.RombelTingkat{{Tingkat: 6, JumlahRombel: 2, JumlahSiswa: 63}}
		daftar = append(daftar, domain.DataSekolah{Sekolah: detail.Sekolah, TahunAjaran: "2026/2027", Rombel: rombel})
	}
	prv := domain.Provenance{ID: provenance.IDBaru(), Sumber: domain.SumberDapodik, Label: domain.LabelResmi}
	return usecase.DataSumberSekolah{Sekolah: daftar, Provenance: prv}, nil
}

type pelaporTiruan struct {
	laporan usecase.LaporanImporSekolah
}

func (p *pelaporTiruan) LaporkanSekolah(l usecase.LaporanImporSekolah) { p.laporan = l }

func jalankanImporUji(t *testing.T, tx pgx.Tx) usecase.LaporanImporSekolah {
	t.Helper()
	repo := repository.SekolahRepository{DB: tx}
	pelapor := &pelaporTiruan{}
	impor := usecase.ImporSekolah{
		Pengambil: pengambilFixture{t: t}, Pemeta: repo, Penyimpan: repo,
		Pencatat: repository.IngestRunRepository{DB: tx}, Pelapor: pelapor,
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	return pelapor.laporan
}

func TestImporSekolah_Koordinat_KecamatanDanKelurahanDariTitik(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)

	jalankanImporUji(t, tx)

	var kecamatan, kelurahan string
	err := tx.QueryRow(t.Context(), `SELECT kode_kecamatan, coalesce(kode_kelurahan, '') FROM sekolah WHERE npsn = '20104886'`).Scan(&kecamatan, &kelurahan)
	if err != nil || kecamatan != "31.72.06" || kelurahan == "" {
		t.Errorf("kecamatan %q, kelurahan %q, err %v", kecamatan, kelurahan, err)
	}
}

func TestImporSekolah_TanpaKoordinat_DilaporkanDanLokasiKosong(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)

	laporan := jalankanImporUji(t, tx)

	var kecamatan string
	var lokasiKosong bool
	err := tx.QueryRow(t.Context(), `SELECT kode_kecamatan, lokasi IS NULL FROM sekolah WHERE npsn = '99999901'`).Scan(&kecamatan, &lokasiKosong)
	if err != nil || kecamatan != "31.72.06" || !lokasiKosong {
		t.Errorf("kecamatan %q, lokasi kosong %v, err %v", kecamatan, lokasiKosong, err)
	}
	if len(laporan.TanpaKoordinat) != 1 || len(laporan.TidakTerpetakan) != 0 {
		t.Errorf("laporan = %+v", laporan)
	}
}

func TestImporSekolah_Madrasah_TersimpanSebagaiMadrasahMenengahPertama(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)

	jalankanImporUji(t, tx)

	var jenjang string
	var madrasah bool
	err := tx.QueryRow(t.Context(), `SELECT jenjang, is_madrasah FROM sekolah WHERE npsn = '20178210'`).Scan(&jenjang, &madrasah)
	if err != nil || jenjang != "menengah_pertama" || !madrasah {
		t.Errorf("jenjang %q, madrasah %v, err %v", jenjang, madrasah, err)
	}
}

func TestImporSekolah_DijalankanDuaKali_TidakMenggandakanData(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)

	jalankanImporUji(t, tx)
	jalankanImporUji(t, tx)

	var baru, berubah, jumlahRombel int
	err := tx.QueryRow(t.Context(), `SELECT jumlah_baru, jumlah_berubah, (SELECT count(*) FROM sekolah_rombel WHERE npsn = '20104886')
FROM ingest_run WHERE sumber = 'dapodik' ORDER BY id DESC LIMIT 1`).Scan(&baru, &berubah, &jumlahRombel)
	if err != nil || baru != 0 || berubah != 0 || jumlahRombel != 1 {
		t.Errorf("baru %d, berubah %d, rombel %d, err %v", baru, berubah, jumlahRombel, err)
	}
}
