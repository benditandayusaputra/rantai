package opentender

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUraiOCDS_HanyaPembeliPemprovDKI_SatkerDanPemenangTerbaca(t *testing.T) {
	berkas, err := os.Open(filepath.Join("testdata", "sampel_2023.jsonl.gz"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = berkas.Close() }()

	riwayat, err := UraiOCDS(berkas, []string{"D69"})

	if err != nil || len(riwayat) != 4 {
		t.Fatalf("riwayat = %d, err = %v", len(riwayat), err)
	}
	r := riwayat[0]
	if r.NamaPenyedia != "PT. NINDYA KARYA (Persero)" || r.Tahun != 2023 || r.Pembeli != "DINAS SUMBER DAYA AIR" || r.Nilai == nil || *r.Nilai != 787168020000 {
		t.Errorf("riwayat pertama = %+v", r)
	}
	if r.KodeSumber != "52764127-1#0" {
		t.Errorf("kode sumber = %s", r.KodeSumber)
	}
}
