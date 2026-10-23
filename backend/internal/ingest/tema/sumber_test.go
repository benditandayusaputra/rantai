package tema

import (
	"os"
	"strings"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const pathTemaPendidikan = "../../../config/tema/pendidikan.yaml"

func TestUraiYAML_TemaPendidikan_SahMenurutKontrak(t *testing.T) {
	isi, err := os.ReadFile(pathTemaPendidikan)
	if err != nil {
		t.Fatal(err)
	}

	tema, err := UraiYAML(isi)

	if err != nil {
		t.Fatal(err)
	}
	if err := tema.Validasi(usecase.PenghitungTerdaftar); err != nil {
		t.Errorf("tema pendidikan tidak sah: %v", err)
	}
	if tema.TanggungJawab[0].Komisi != "E" || tema.Indikator.AmbangRasio.HampirCukup != 1.0 {
		t.Errorf("isi tema tidak terbaca utuh: %+v", tema)
	}
}

func TestUraiYAML_FieldSalahKetik_Ditolak(t *testing.T) {
	_, err := UraiYAML([]byte("kode: pendidikan\nindikatr:\n  penghitung: kesenjangan_kursi_smp\n"))

	if err == nil || !strings.Contains(err.Error(), "indikatr") {
		t.Errorf("harap field tak dikenal ditolak, dapat %v", err)
	}
}
