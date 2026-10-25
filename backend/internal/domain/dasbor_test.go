package domain

import "testing"

func TestLengkapiBulan_BulanKosongDiisiNol(t *testing.T) {
	kasus := []struct {
		nama   string
		isi    []Porsi
		maret  int
		jumlah int
	}{
		{"tanpa data", []Porsi{}, 0, 0},
		{"satu bulan", []Porsi{{Kunci: "2025-03", Jumlah: 7}}, 7, 7},
		{"bulan tahun lain diabaikan", []Porsi{{Kunci: "2024-03", Jumlah: 9}, {Kunci: "2025-12", Jumlah: 2}}, 0, 2},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			hasil := LengkapiBulan(2025, k.isi)

			total := 0
			for _, p := range hasil {
				total += p.Jumlah
			}
			if len(hasil) != JumlahBulanSetahun || hasil[0].Kunci != "2025-01" || hasil[11].Kunci != "2025-12" {
				t.Fatalf("kunci bulan %v", hasil)
			}
			if hasil[2].Jumlah != k.maret || total != k.jumlah {
				t.Errorf("maret %d, total %d", hasil[2].Jumlah, total)
			}
		})
	}
}

func TestTahunTerlengkap_PilihTahunDenganJumlahTerbanyak(t *testing.T) {
	kasus := []struct {
		nama    string
		tahunan []Porsi
		tahun   int
	}{
		{"tanpa data", []Porsi{}, 0},
		{"jumlah terbanyak menang", []Porsi{{Kunci: "2025", Jumlah: 900}, {Kunci: "2026", Jumlah: 30}}, 2025},
		{"seri dimenangkan tahun terbaru", []Porsi{{Kunci: "2026", Jumlah: 5}, {Kunci: "2024", Jumlah: 5}}, 2026},
		{"kunci bukan tahun diabaikan", []Porsi{{Kunci: "lain", Jumlah: 99}, {Kunci: "2023", Jumlah: 1}}, 2023},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if hasil := TahunTerlengkap(k.tahunan); hasil != k.tahun {
				t.Errorf("tahun %d, ingin %d", hasil, k.tahun)
			}
		})
	}
}
