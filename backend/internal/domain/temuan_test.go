package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var terlarangUji = domain.KosakataTerlarang{"korupsi", "mark up"}

func TestTemuanMasuk_Validasi(t *testing.T) {
	isiSah := "Nilai kontrak di dokumen SPK berbeda dengan angka di LPSE."
	kasus := []struct {
		nama  string
		masuk domain.TemuanMasuk
		ingin error
	}{
		{"sah", domain.TemuanMasuk{Isi: isiSah, TautanSumber: "https://lpse.jakarta.go.id/x"}, nil},
		{"http juga sah", domain.TemuanMasuk{Isi: isiSah, TautanSumber: "http://contoh.id/dok.pdf"}, nil},
		{"isi terlalu pendek", domain.TemuanMasuk{Isi: "terlalu pendek", TautanSumber: "https://a.id"}, domain.ErrIsiTemuanTidakSah},
		{"isi terlalu panjang", domain.TemuanMasuk{Isi: strings.Repeat("a", 1001), TautanSumber: "https://a.id"}, domain.ErrIsiTemuanTidakSah},
		{"kata tuduhan", domain.TemuanMasuk{Isi: "Paket ini jelas korupsi karena nilainya tinggi", TautanSumber: "https://a.id"}, domain.ErrKosakataTemuan},
		{"tanpa skema", domain.TemuanMasuk{Isi: isiSah, TautanSumber: "lpse.jakarta.go.id"}, domain.ErrTautanTemuanTidakSah},
		{"skema javascript", domain.TemuanMasuk{Isi: isiSah, TautanSumber: "javascript:alert(1)"}, domain.ErrTautanTemuanTidakSah},
		{"kosong", domain.TemuanMasuk{Isi: isiSah}, domain.ErrTautanTemuanTidakSah},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := k.masuk.Validasi(terlarangUji); !errors.Is(err, k.ingin) {
				t.Errorf("err %v, ingin %v", err, k.ingin)
			}
		})
	}
}

func TestTanggapanMasuk_Validasi(t *testing.T) {
	kasus := []struct {
		nama  string
		masuk domain.TanggapanMasuk
		ingin error
	}{
		{"mendukung tanpa tautan", domain.TanggapanMasuk{Sikap: domain.SikapMendukung}, nil},
		{"membantah dengan tautan", domain.TanggapanMasuk{Sikap: domain.SikapMembantah, TautanSumber: "https://a.id"}, nil},
		{"sikap asing", domain.TanggapanMasuk{Sikap: "netral"}, domain.ErrSikapTemuanTidakSah},
		{"tautan rusak", domain.TanggapanMasuk{Sikap: domain.SikapMendukung, TautanSumber: "ftp://a.id"}, domain.ErrTautanTemuanTidakSah},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if err := k.masuk.Validasi(); !errors.Is(err, k.ingin) {
				t.Errorf("err %v, ingin %v", err, k.ingin)
			}
		})
	}
}

func TestTemuan_Status(t *testing.T) {
	kasus := []struct {
		nama                 string
		mendukung, membantah int
		ingin                domain.StatusTemuan
	}{
		{"belum ada tanggapan", 0, 0, domain.StatusTemuanMenunggu},
		{"satu dukungan", 1, 0, domain.StatusTemuanMenunggu},
		{"dua dukungan", 2, 0, domain.StatusTemuanTerkonfirmasi},
		{"ada bantahan", 3, 1, domain.StatusTemuanDiperdebatkan},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			got := domain.Temuan{JumlahMendukung: k.mendukung, JumlahMembantah: k.membantah}.Status()
			if got != k.ingin {
				t.Errorf("status %s, ingin %s", got, k.ingin)
			}
		})
	}
}
