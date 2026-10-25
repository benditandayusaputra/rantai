package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestPeralihanTema_DuaTemaPendidikan_TingkatAsalDanTujuan(t *testing.T) {
	kasus := []struct {
		tema                     string
		asalTerakhir, tujuanAwal int
		madrasahTujuan           BentukSekolah
		banyakBentuk             bool
	}{
		{TemaPendidikan, TingkatSDTerakhir, TingkatSMPPertama, BentukMTs, false},
		{TemaPendidikanSMASMK, TingkatSMPTerakhir, TingkatSMAPertama, BentukMA, true},
	}
	for _, k := range kasus {
		p, err := PeralihanTema(k.tema)
		if err != nil {
			t.Fatalf("%s: %v", k.tema, err)
		}
		if p.TingkatAsalTerakhir() != k.asalTerakhir || p.TingkatTujuanPertama() != k.tujuanAwal ||
			p.BentukMadrasahTujuan != k.madrasahTujuan || p.IsBanyakBentukTujuan() != k.banyakBentuk {
			t.Errorf("%s: %+v", k.tema, p)
		}
	}
}

func TestPeralihanTema_TemaTanpaPeralihan_TidakDitemukan(t *testing.T) {
	_, err := PeralihanTema("banjir")

	var kesalahan *Kesalahan
	if !errors.As(err, &kesalahan) || kesalahan.Kategori != KategoriTidakDitemukan {
		t.Errorf("err = %v", err)
	}
}

func TestPeralihan_NamaTujuanDanRumusEstimasi(t *testing.T) {
	p, _ := PeralihanTema(TemaPendidikanSMASMK)

	if p.NamaTujuan() != "SMA/SMK" {
		t.Errorf("nama tujuan = %q", p.NamaTujuan())
	}
	if rumus := p.RumusEstimasi(3); !strings.Contains(rumus, "siswa SMP") || !strings.Contains(rumus, "÷ 3") || !strings.Contains(rumus, "kelas 9") {
		t.Errorf("rumus = %q", rumus)
	}
}
