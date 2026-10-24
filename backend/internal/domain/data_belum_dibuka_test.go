package domain

import (
	"strings"
	"testing"
)

func konfigurasiDataBelumDibukaUji(t *testing.T) KonfigurasiDataBelumDibuka {
	return bacaYAMLUji[KonfigurasiDataBelumDibuka](t, "../../config/data_belum_dibuka.yaml")
}

func TestDataBelumDibuka_EmpatDataLengkapDenganPengganti(t *testing.T) {
	k := konfigurasiDataBelumDibukaUji(t)

	if len(k.Data) != 4 {
		t.Fatalf("jumlah data = %d, harap 4", len(k.Data))
	}
	for _, d := range k.Data {
		for _, teks := range []Teks{d.Nama, d.Alasan, d.Pengganti, d.Rincian} {
			if teks.ID == "" || teks.EN == "" || d.FiturPengganti == "" {
				t.Errorf("%s: isian kosong", d.Kode)
			}
		}
	}
}

func TestSusunPermohonan_SetiapDataDanBahasa_BerdasarUU14Tahun2008TanpaPlaceholder(t *testing.T) {
	k := konfigurasiDataBelumDibukaUji(t)
	terlarang := bacaKosakataTerlarang(t)
	for _, d := range k.Data {
		for _, bahasa := range []string{BahasaIndonesia, BahasaInggris} {
			dokumen, ada := k.SusunPermohonan(d.Kode, bahasa)
			teks := dokumen.Teks()

			switch {
			case !ada:
				t.Errorf("%s: permohonan tidak tersusun", d.Kode)
			case !strings.Contains(teks, "Undang-Undang Nomor 14 Tahun 2008"):
				t.Errorf("%s %s: tanpa dasar UU 14/2008", d.Kode, bahasa)
			case !strings.Contains(teks, d.Rincian.Pilih(bahasa)) || !strings.Contains(teks, k.PPID.Alamat):
				t.Errorf("%s %s: rincian atau alamat PPID hilang", d.Kode, bahasa)
			case strings.ContainsAny(teks, "{}"):
				t.Errorf("%s %s: placeholder tersisa", d.Kode, bahasa)
			case len(terlarang.Temukan(teks)) > 0:
				t.Errorf("%s %s: kosakata terlarang %v", d.Kode, bahasa, terlarang.Temukan(teks))
			}
		}
	}
}

func TestSusunPermohonan_KodeTidakDikenal_TidakAda(t *testing.T) {
	if _, ada := konfigurasiDataBelumDibukaUji(t).SusunPermohonan("gaji_anggota", BahasaIndonesia); ada {
		t.Error("kode tak dikenal menghasilkan permohonan")
	}
}
