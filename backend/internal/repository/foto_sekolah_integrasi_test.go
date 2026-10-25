package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	npsnFotoUji     = "90000301"
	urlFotoResmi    = "https://file.data.kemendikdasmen.go.id/sekolahkita/90/9000/90000301-1.jpg"
	urlFotoCadangan = "https://upload.wikimedia.org/wikipedia/commons/thumb/a/ab/Uji.jpg/330px-Uji.jpg"
)

func simpanSekolahDenganFoto(t *testing.T, tx pgx.Tx, foto *domain.Foto) {
	t.Helper()
	sekolah := usecase.SekolahTerpetakan{
		DataSekolah: domain.DataSekolah{
			Sekolah:     domain.Sekolah{NPSN: npsnFotoUji, Nama: "SMP UJI FOTO", Bentuk: domain.BentukSMP, Status: domain.StatusNegeri, Foto: foto},
			TahunAjaran: "2026/2027",
		},
		PenempatanSekolah: domain.PenempatanSekolah{KodeKecamatan: "31.72.06"},
	}
	prov := provenance.Baru(provenance.Catatan{Sumber: domain.SumberDapodik, Label: domain.LabelResmi}, nil, time.Now())
	data := usecase.DataSimpanSekolah{Sekolah: []usecase.SekolahTerpetakan{sekolah}, Provenance: prov}
	if _, err := (repository.SekolahRepository{DB: tx}).SimpanSekolah(t.Context(), data); err != nil {
		t.Fatal(err)
	}
}

func simpanFotoCadanganUji(t *testing.T, tx pgx.Tx, npsn string) {
	t.Helper()
	cadangan := []domain.FotoSekolah{{NPSN: npsn, Foto: domain.Foto{URL: urlFotoCadangan, Halaman: "https://commons.wikimedia.org/wiki/File:Uji.jpg", Kredit: "Uji · CC BY-SA 4.0"}}}
	if _, err := (repository.FotoSekolahRepository{DB: tx}).SimpanFotoCadangan(t.Context(), domain.SumberFotoWikimedia, cadangan); err != nil {
		t.Fatal(err)
	}
}

func fotoTersimpan(t *testing.T, tx pgx.Tx) (url, sumber string) {
	t.Helper()
	err := tx.QueryRow(t.Context(), `SELECT url, sumber FROM foto_sekolah WHERE npsn = $1`, npsnFotoUji).Scan(&url, &sumber)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		t.Fatal(err)
	}
	return url, sumber
}

func TestFotoSekolah_CadanganMengisiKekosongan_FotoResmiSelaluMenang(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	ujiwilayah.MuatDKI(t, tx)
	simpanSekolahDenganFoto(t, tx, nil)

	simpanFotoCadanganUji(t, tx, npsnFotoUji)
	if url, sumber := fotoTersimpan(t, tx); url != urlFotoCadangan || sumber != string(domain.SumberFotoWikimedia) {
		t.Fatalf("harap foto cadangan mengisi sekolah tanpa foto resmi, dapat %s (%s)", url, sumber)
	}

	resmi := domain.Foto{URL: urlFotoResmi, Sumber: domain.SumberFotoSekolahKita, Halaman: "https://sekolah.data.kemendikdasmen.go.id/profil-sekolah/x"}
	simpanSekolahDenganFoto(t, tx, &resmi)
	simpanFotoCadanganUji(t, tx, npsnFotoUji)

	if url, sumber := fotoTersimpan(t, tx); url != urlFotoResmi || sumber != string(domain.SumberFotoSekolahKita) {
		t.Errorf("harap foto resmi menimpa cadangan dan tidak tertimpa lagi, dapat %s (%s)", url, sumber)
	}
}

func TestFotoSekolah_NPSNTakDikenal_DilewatiTanpaGalat(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)

	simpanFotoCadanganUji(t, tx, "99999999")

	var jumlah int
	if err := tx.QueryRow(t.Context(), `SELECT count(*) FROM foto_sekolah WHERE npsn = '99999999'`).Scan(&jumlah); err != nil || jumlah != 0 {
		t.Errorf("jumlah %d, err %v", jumlah, err)
	}
}
