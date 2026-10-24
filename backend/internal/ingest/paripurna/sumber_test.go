package paripurna

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

func sumberUji() Sumber {
	return Sumber{
		Konfigurasi: Konfigurasi{
			PathSikapFraksi: "testdata/sikap_fraksi.csv", PathAspirasiReses: "testdata/aspirasi_reses.csv",
			URLSumber: "https://dprd-dkijakartaprov.go.id/", Catatan: "kurasi uji",
		},
		Sekarang: time.Now,
	}
}

func TestAmbil_DuaBerkas_EntriDanProvenanceResmiDPRD(t *testing.T) {
	data, err := sumberUji().Ambil(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(data.Sikap) != 2 || len(data.Aspirasi) != 2 || data.Provenance.Sumber != domain.SumberDPRD || data.Provenance.Label != domain.LabelResmi {
		t.Fatalf("data = %+v", data)
	}
	sikap := data.Sikap[0]
	if sikap.Agenda != domain.AgendaAPBD || sikap.Tanggal.Day() != 19 || sikap.Ringkasan != "Meminta subsidi transportasi berbasis kinerja, disertai peningkatan layanan." {
		t.Errorf("sikap = %+v", sikap)
	}
	if data.Aspirasi[0].KodeWilayah != "31.75.02.1006" || data.Aspirasi[1].KodeWilayah != "" {
		t.Errorf("aspirasi = %+v", data.Aspirasi)
	}
}

func TestUraiSikapFraksi_KolomKurang_Ditolak(t *testing.T) {
	if _, err := UraiSikapFraksi([]byte("fraksi,agenda\nFraksi Partai Golkar,apbd\n")); !errors.Is(err, ErrKolomTidakLengkap) {
		t.Errorf("err = %v", err)
	}
}

func TestUraiAspirasiReses_TanggalRusak_DitolakDenganNomorBaris(t *testing.T) {
	isi := "masa_reses,tanggal_paripurna,bidang,uraian,kode_wilayah,url\nReses,26/08/2026,kesehatan,Uraian,,https://contoh.test/\n"

	if _, err := UraiAspirasiReses([]byte(isi)); err == nil || err.Error()[:7] != "baris 2" {
		t.Errorf("err = %v", err)
	}
}
