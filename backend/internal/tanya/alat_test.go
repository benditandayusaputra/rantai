package tanya

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type masukanGema struct {
	Teks string `json:"teks"`
}

func alatGema(nama string) Alat {
	return BuatAlat(Spek{Nama: nama}, func(_ context.Context, m masukanGema) (any, error) {
		return map[string]string{"teks": m.Teks, "provenance_id": "prv_1"}, nil
	})
}

func TestRegistryJalankan_AlatTerdaftar_HasilJSON(t *testing.T) {
	registry, err := BuatRegistry(alatGema("gema"))
	if err != nil {
		t.Fatal(err)
	}

	hasil, err := registry.Jalankan(t.Context(), "gema", json.RawMessage(`{"teks":"halo"}`))

	if harap := `{"provenance_id":"prv_1","teks":"halo"}`; err != nil || string(hasil) != harap {
		t.Errorf("dapat %s (err %v), harap %s", hasil, err, harap)
	}
}

func TestRegistryJalankan_AlatTidakDikenal_Error(t *testing.T) {
	registry, err := BuatRegistry(alatGema("gema"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = registry.Jalankan(t.Context(), "hapus_data", json.RawMessage(`{}`))

	if !errors.Is(err, ErrAlatTidakDikenal) {
		t.Errorf("err = %v, harap ErrAlatTidakDikenal", err)
	}
}

func TestRegistryJalankan_MasukanTidakSesuaiSkema_InputTidakValid(t *testing.T) {
	registry, err := BuatRegistry(alatGema("gema"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = registry.Jalankan(t.Context(), "gema", json.RawMessage(`{"teks":7}`))

	var kesalahan *domain.Kesalahan
	if !errors.As(err, &kesalahan) || kesalahan.Kategori != domain.KategoriInputTidakValid {
		t.Errorf("err = %v, harap input tidak valid", err)
	}
}

func TestBuatRegistry_NamaGanda_Ditolak(t *testing.T) {
	_, err := BuatRegistry(alatGema("gema"), alatGema("gema"))

	if !errors.Is(err, ErrAlatGanda) {
		t.Errorf("err = %v, harap ErrAlatGanda", err)
	}
}
