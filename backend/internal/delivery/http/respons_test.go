package http

import (
	"errors"
	"fmt"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

func TestPetakanError_JenisError_KodeDanPesanSesuai(t *testing.T) {
	kasus := []struct {
		nama  string
		err   error
		kode  int
		pesan string
	}{
		{"tidak ditemukan", domain.TidakDitemukan("paket tidak ditemukan"), fiber.StatusNotFound, "paket tidak ditemukan"},
		{"input tidak valid", domain.InputTidakValid("kode wilayah tidak valid"), fiber.StatusBadRequest, "kode wilayah tidak valid"},
		{"domain terbungkus konteks", fmt.Errorf("ambil paket X: %w", domain.TidakDitemukan("paket tidak ditemukan")), fiber.StatusNotFound, "paket tidak ditemukan"},
		{"error fiber", fiber.ErrMethodNotAllowed, fiber.StatusMethodNotAllowed, "metode tidak diizinkan"},
		{"error fiber tanpa pesan khusus", fiber.ErrTeapot, fiber.StatusTeapot, "permintaan tidak valid"},
		{"error tak dikenal disembunyikan", errors.New("koneksi db putus"), fiber.StatusInternalServerError, "terjadi kesalahan di server"},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			kode, pesan := petakanError(k.err)

			if kode != k.kode || pesan != k.pesan {
				t.Errorf("dapat (%d, %q), harap (%d, %q)", kode, pesan, k.kode, k.pesan)
			}
		})
	}
}

func TestTanganiError_RuteTidakAda_AmplopError404(t *testing.T) {
	hasil := panggil(t, appUji(), "/api/v1/tidak-ada")

	if hasil.kode != fiber.StatusNotFound || hasil.respons.Code != fiber.StatusNotFound {
		t.Fatalf("kode = %d, amplop = %d", hasil.kode, hasil.respons.Code)
	}
	if hasil.respons.Data != nil || hasil.respons.Message != "halaman tidak ditemukan" {
		t.Errorf("amplop error tidak sesuai: %+v", hasil.respons)
	}
}
