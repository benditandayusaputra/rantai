package domain

import (
	"testing"
	"time"
)

func tanggal(tahun int, bulan time.Month, hari int) *time.Time {
	t := time.Date(tahun, bulan, hari, 0, 0, 0, 0, time.UTC)
	return &t
}

func aturanPemeliharaanUji() AturanPemeliharaan {
	return AturanPemeliharaan{
		JenisMasaHari: map[string]int{"permanen": 180, "semi_permanen": 90},
		Kategori:      map[string]MasaKategori{"ruang_kelas": {Jenis: "permanen", DurasiBawaanHari: 120}},
		Bawaan:        MasaKategori{Jenis: "semi_permanen", DurasiBawaanHari: 60},
	}
}

func TestEstimasiPemeliharaan_EmpatStatus(t *testing.T) {
	jadwal := JadwalPaket{Kategori: "ruang_kelas", JadwalSelesai: tanggal(2026, time.March, 1)}
	kasus := []struct {
		nama     string
		jadwal   JadwalPaket
		hariIni  *time.Time
		harap    StatusPemeliharaan
		sisaHari int
	}{
		{"sebelum jadwal selesai", jadwal, tanggal(2026, time.February, 1), StatusPelaksanaan, 0},
		{"dalam masa pemeliharaan", jadwal, tanggal(2026, time.August, 27), StatusPemeliharaanBerjalan, 1},
		{"tepat di akhir pemeliharaan", jadwal, tanggal(2026, time.August, 28), StatusPemeliharaanBerakhir, 0},
		{"tanpa tanggal apa pun", JadwalPaket{Kategori: "ruang_kelas"}, tanggal(2026, time.August, 1), StatusTidakDapatDiestimasi, 0},
	}
	for _, k := range kasus {
		e := aturanPemeliharaanUji().Estimasi(k.jadwal, *k.hariIni)
		if e.Status != k.harap || e.SisaHari != k.sisaHari || e.Label != LabelEstimasi {
			t.Errorf("%s: %+v, harap %s sisa %d", k.nama, e, k.harap, k.sisaHari)
		}
	}
}

func TestEstimasiPemeliharaan_SerahTerimaDariInstansi_LabelResmiDanMenggantikanJadwal(t *testing.T) {
	jadwal := JadwalPaket{Kategori: "ruang_kelas", JadwalSelesai: tanggal(2026, time.March, 1), TanggalSerahTerima: tanggal(2026, time.April, 1)}

	e := aturanPemeliharaanUji().Estimasi(jadwal, *tanggal(2026, time.April, 2))

	if e.Label != LabelResmi || e.Status != StatusPemeliharaanBerjalan || !e.AkhirPemeliharaan.Equal(*tanggal(2026, time.September, 28)) {
		t.Errorf("estimasi = %+v, akhir %v", e, e.AkhirPemeliharaan)
	}
}

func TestEstimasiPemeliharaan_HanyaTanggalKontrak_DurasiBawaanKategoriLain(t *testing.T) {
	jadwal := JadwalPaket{Kategori: "pagar", TanggalKontrak: tanggal(2026, time.January, 1)}

	e := aturanPemeliharaanUji().Estimasi(jadwal, *tanggal(2026, time.January, 2))

	if !e.AkhirPelaksanaan.Equal(*tanggal(2026, time.March, 2)) || !e.AkhirPemeliharaan.Equal(*tanggal(2026, time.May, 31)) {
		t.Errorf("akhir pelaksanaan %v, pemeliharaan %v", e.AkhirPelaksanaan, e.AkhirPemeliharaan)
	}
}
