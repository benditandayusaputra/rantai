package domain

import (
	"testing"
	"time"
)

func statusMata(rantai []MataRantai, kode string) StatusMata {
	for _, m := range rantai {
		if m.Kode == kode {
			return m.Status
		}
	}
	return ""
}

func TestSusunRantai_StatusTiapMataSesuaiTabelANG03(t *testing.T) {
	selesai := time.Date(2026, time.December, 31, 0, 0, 0, 0, time.UTC)
	lengkap := PaketLengkap{
		RUP:                &RUPPaket{Kode: "1", IsDetailTersedia: true, PelaksanaanSelesai: &selesai},
		Penyedia:           &PenyediaPaket{ID: 1, Nama: "CV Uji"},
		Lokasi:             LokasiTersimpan{Tingkat: TingkatKelurahan},
		TanggalSerahTerima: &selesai,
		JumlahLaporan:      2,
	}
	kosong := PaketLengkap{Lokasi: LokasiTersimpan{Tingkat: TingkatKota}}
	kasus := []struct {
		nama         string
		paket        PaketLengkap
		estimasi     EstimasiPemeliharaan
		harapPerMata map[string]StatusMata
	}{
		{"semua data ada", lengkap, EstimasiPemeliharaan{Label: LabelResmi}, map[string]StatusMata{
			MataRencana: StatusMataAda, MataPemilihan: StatusMataAda, MataPemenang: StatusMataAda, MataPelaksanaan: StatusMataEstimasi,
			MataLokasi: StatusMataAda, MataSerahTerima: StatusMataAda, MataPemeliharaan: StatusMataAda, MataBuktiWarga: StatusMataAda,
		}},
		{"data minimal", kosong, EstimasiPemeliharaan{Label: LabelEstimasi}, map[string]StatusMata{
			MataRencana: StatusMataPutus, MataPemilihan: StatusMataAda, MataPemenang: StatusMataMenunggu, MataPelaksanaan: StatusMataPutus,
			MataLokasi: StatusMataPutus, MataSerahTerima: StatusMataEstimasi, MataPemeliharaan: StatusMataEstimasi, MataBuktiWarga: StatusMataMenunggu,
		}},
	}
	for _, k := range kasus {
		rantai := SusunRantai(k.paket, k.estimasi)
		if len(rantai) != len(k.harapPerMata) {
			t.Fatalf("%s: %d mata rantai", k.nama, len(rantai))
		}
		for mata, harap := range k.harapPerMata {
			if status := statusMata(rantai, mata); status != harap {
				t.Errorf("%s: %s = %s, harap %s", k.nama, mata, status, harap)
			}
		}
	}
}

func TestSusunRantai_RUPTanpaDetail_AdaDenganAlasan(t *testing.T) {
	rantai := SusunRantai(PaketLengkap{RUP: &RUPPaket{Kode: "1"}}, EstimasiPemeliharaan{})

	if rantai[0].Status != StatusMataAda || rantai[0].Alasan != "kode RUP tercatat, detail tidak tersedia" {
		t.Errorf("rencana = %+v", rantai[0])
	}
}

func TestSusunRantai_TenderTanpaRUP_AlasanSumberTertutup(t *testing.T) {
	rantai := SusunRantai(PaketLengkap{IsTender: true, NilaiKontrak: new(float64)}, EstimasiPemeliharaan{})

	if rantai[0].Status != StatusMataPutus || rantai[0].Alasan != "kode RUP paket tender tidak dibuka di SPSE" || rantai[2].Status != StatusMataAda {
		t.Errorf("rantai = %+v", rantai)
	}
}
