package tanya

import (
	"errors"
	"reflect"
	"testing"
)

func cacheUji() []EntriCache {
	pertanyaan := []string{
		"Kenapa SMP negeri di kecamatanku penuh?",
		"Proyek sekolah apa yang dibiayai di kecamatanku tahun ini?",
		"Siapa wakil dapilku yang duduk di Komisi E?",
	}
	cache := make([]EntriCache, len(pertanyaan))
	for i, p := range pertanyaan {
		cache[i] = EntriCache{
			Pertanyaan: p, KodeKecamatanKonteks: "31.72.04", Alat: []string{"ringkasan_kecamatan"},
			Teks: "Jawaban demo [1].", Rujukan: []Rujukan{{Nomor: 1, ProvenanceID: "prv_demo"}},
		}
	}
	return cache
}

func TestPenjawabTercache_TigaPertanyaanDemoTanpaModel_TerjawabDariCache(t *testing.T) {
	p := PenjawabTercache{Cache: cacheUji(), Cadangan: Penjawab{Aturan: Aturan{PanjangMaksPertanyaan: 500}}}
	for _, entri := range cacheUji() {
		t.Run(entri.Pertanyaan, func(t *testing.T) {
			kabar := []string{}

			jawaban, err := p.Jawab(t.Context(), Permintaan{Pertanyaan: entri.Pertanyaan, KodeKecamatanKonteks: "31.72.04.1001"}, func(k Kabar) {
				kabar = append(kabar, k.Alat)
			})

			if err != nil || jawaban.Teks != entri.Teks || !reflect.DeepEqual(jawaban.Rujukan, entri.Rujukan) {
				t.Errorf("jawaban = %+v, err = %v", jawaban, err)
			}
			if !reflect.DeepEqual(kabar, entri.Alat) {
				t.Errorf("kabar = %v, harap %v", kabar, entri.Alat)
			}
		})
	}
}

func TestPenjawabTercache_HurufDanTandaBacaBerbeda_TetapCocok(t *testing.T) {
	p := PenjawabTercache{Cache: cacheUji()}

	_, err := p.Jawab(t.Context(), Permintaan{Pertanyaan: "  kenapa SMP negeri di kecamatanku PENUH "}, func(Kabar) {})

	if err != nil {
		t.Errorf("err = %v, harap cocok cache", err)
	}
}

func TestPenjawabTercache_KecamatanLainTanpaModel_TidakTercache(t *testing.T) {
	p := PenjawabTercache{Cache: cacheUji()}

	_, err := p.Jawab(t.Context(), Permintaan{Pertanyaan: "Kenapa SMP negeri di kecamatanku penuh?", KodeKecamatanKonteks: "31.71.01"}, func(Kabar) {})

	if !errors.Is(err, ErrTidakTercache) {
		t.Errorf("err = %v, harap ErrTidakTercache", err)
	}
}

func TestPenjawabTercache_TidakTercacheDenganModel_DijawabModel(t *testing.T) {
	cadangan, _ := penjawabUji(t, Giliran{Teks: "Data belum tersedia."})
	p := PenjawabTercache{Cache: cacheUji(), Cadangan: cadangan}

	jawaban, err := p.Jawab(t.Context(), Permintaan{Pertanyaan: "Berapa jumlah puskesmas?"}, func(Kabar) {})

	if err != nil || jawaban.Teks != "Data belum tersedia." {
		t.Errorf("jawaban = %+v, err = %v", jawaban, err)
	}
}
