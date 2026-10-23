package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pembacaJadwalTiruan struct {
	tersimpan JadwalPaketTersimpan
}

func (p pembacaJadwalTiruan) AmbilJadwalPaket(context.Context, string) (JadwalPaketTersimpan, bool, error) {
	return p.tersimpan, p.tersimpan.JenisPengadaan != "", nil
}

func hariKe(bulan time.Month, hari int) time.Time {
	return time.Date(2026, bulan, hari, 9, 0, 0, 0, time.UTC)
}

func statusUji(jadwal domain.JadwalPaket, jenis string, sekarang time.Time) StatusPemeliharaan {
	return StatusPemeliharaan{
		Pembaca: pembacaJadwalTiruan{tersimpan: JadwalPaketTersimpan{JenisPengadaan: jenis, Jadwal: jadwal}},
		Konfigurasi: KonfigurasiPemeliharaan{
			AturanPemeliharaan: domain.AturanPemeliharaan{
				JenisMasaHari: map[string]int{"permanen": 180, "semi_permanen": 90},
				Bawaan:        domain.MasaKategori{Jenis: "semi_permanen", DurasiBawaanHari: 60},
			},
			JenisPengadaan: []string{"Pekerjaan Konstruksi"},
		},
		Sekarang: func() time.Time { return sekarang },
	}
}

func TestStatusPemeliharaan_JamTiruan_EmpatStatus(t *testing.T) {
	selesai := hariKe(time.March, 1)
	kasus := []struct {
		jadwal   domain.JadwalPaket
		sekarang time.Time
		harap    domain.StatusPemeliharaan
	}{
		{domain.JadwalPaket{JadwalSelesai: &selesai}, hariKe(time.February, 1), domain.StatusPelaksanaan},
		{domain.JadwalPaket{JadwalSelesai: &selesai}, hariKe(time.April, 1), domain.StatusPemeliharaanBerjalan},
		{domain.JadwalPaket{JadwalSelesai: &selesai}, hariKe(time.July, 1), domain.StatusPemeliharaanBerakhir},
		{domain.JadwalPaket{}, hariKe(time.July, 1), domain.StatusTidakDapatDiestimasi},
	}
	for _, k := range kasus {
		hasil, err := statusUji(k.jadwal, "Pekerjaan Konstruksi", k.sekarang).Hitung(t.Context(), "nontender-1")
		if err != nil || !hasil.IsBerlaku || hasil.Estimasi.Status != k.harap || hasil.Estimasi.Label != domain.LabelEstimasi {
			t.Errorf("%v: %+v, err %v, harap %s", k.sekarang, hasil, err, k.harap)
		}
	}
}

func TestStatusPemeliharaan_InstansiIsiSerahTerima_LabelMenjadiResmi(t *testing.T) {
	selesai, serahTerima := hariKe(time.March, 1), hariKe(time.March, 20)
	jadwal := domain.JadwalPaket{JadwalSelesai: &selesai}
	sebelum, _ := statusUji(jadwal, "Pekerjaan Konstruksi", hariKe(time.April, 1)).Hitung(t.Context(), "nontender-1")

	jadwal.TanggalSerahTerima = &serahTerima
	sesudah, _ := statusUji(jadwal, "Pekerjaan Konstruksi", hariKe(time.April, 1)).Hitung(t.Context(), "nontender-1")

	if sebelum.Estimasi.Label != domain.LabelEstimasi || sesudah.Estimasi.Label != domain.LabelResmi || sesudah.Estimasi.SisaHari != 78 {
		t.Errorf("sebelum %+v, sesudah %+v", sebelum.Estimasi, sesudah.Estimasi)
	}
}

func TestStatusPemeliharaan_BukanKonstruksi_TidakBerlaku(t *testing.T) {
	hasil, err := statusUji(domain.JadwalPaket{}, "Jasa Lainnya", hariKe(time.April, 1)).Hitung(t.Context(), "nontender-1")

	if err != nil || hasil.IsBerlaku {
		t.Errorf("hasil %+v, err %v", hasil, err)
	}
}
