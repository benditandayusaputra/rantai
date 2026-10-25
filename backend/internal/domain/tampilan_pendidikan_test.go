package domain

import "testing"

func TestKeterbatasan_DataPerSekolahTanpaMadrasah_TigaCatatan(t *testing.T) {
	m := masukanUji()
	m.JumlahMTs = 2
	detail := DetailPendidikanKecamatan{IndikatorTersimpan: IndikatorTersimpan{IndikatorKesenjanganSMP: HitungKesenjanganKursiSMP(m, ambangUji)}}

	daftar := detail.Keterbatasan()

	harap := []string{KeterbatasanSPMBTidakDibatasiKecamatan, KeterbatasanCalonLulusanTanpaMI, KeterbatasanKursiMadrasahTidakTersedia}
	if len(daftar) != len(harap) || daftar[0] != harap[0] || daftar[1] != harap[1] || daftar[2] != harap[2] {
		t.Errorf("keterbatasan = %v", daftar)
	}
}

func TestTrenLayakTampil_KurangDariDuaPeriode_Kosong(t *testing.T) {
	kasus := []struct {
		nama   string
		titik  []TitikTren
		jumlah int
	}{
		{"tanpa periode", nil, 0},
		{"satu periode", []TitikTren{{Periode: "2026/2027", Rasio: 1.4}}, 0},
		{"dua periode", []TitikTren{{Periode: "2025/2026", Rasio: 1.2}, {Periode: "2026/2027", Rasio: 1.4}}, 2},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			hasil := TrenLayakTampil(k.titik)
			if hasil == nil || len(hasil) != k.jumlah {
				t.Errorf("tren = %v, harap %d titik", hasil, k.jumlah)
			}
		})
	}
}
