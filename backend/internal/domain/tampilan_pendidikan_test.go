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
