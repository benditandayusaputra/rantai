package domain

import (
	"slices"
	"testing"
)

func TestLabelAnggaran_LokasiDariNamaPaket_ProksiKecualiAdaSimulasi(t *testing.T) {
	if label := (RingkasanAnggaran{}).Label(); label != LabelProksi {
		t.Errorf("label = %s, harap proksi", label)
	}
	if label := (RingkasanAnggaran{AdaSimulasi: true}).Label(); label != LabelSimulasi {
		t.Errorf("label = %s, harap simulasi", label)
	}
}

func TestKomisiPengawasSemuaTema_DuaTema_GabunganUnik(t *testing.T) {
	daftar := []TemaTersimpan{
		{Tema: Tema{KonfigurasiTema: KonfigurasiTema{TanggungJawab: []TanggungJawab{{Komisi: "E"}}}}},
		{Tema: Tema{KonfigurasiTema: KonfigurasiTema{TanggungJawab: []TanggungJawab{{Komisi: "D"}, {Komisi: "E"}}}}},
	}

	if komisi := KomisiPengawasSemuaTema(daftar); !slices.Equal(komisi, []string{"D", "E"}) {
		t.Errorf("komisi = %v", komisi)
	}
}
