package domain

import "testing"

var ambangUji = AmbangRasio{SangatKurang: 0.5, Kurang: 0.8, HampirCukup: 1.0}

const batasRombelUji = 30

func masukanUji() MasukanKesenjanganSMP {
	return MasukanKesenjanganSMP{
		KodeKecamatan: "31.72.06", CalonLulusan: 1000, LabelCalonLulusan: LabelResmi, MetodeCalonLulusan: MetodePerSekolah,
		KursiNegeri: 600, KursiSPMBBersama: 50, BatasMuridPerRombelSMP: batasRombelUji,
	}
}

func TestHitungKesenjanganKursiSMP_KasusTepi(t *testing.T) {
	kasus := []struct {
		nama        string
		ubah        func(*MasukanKesenjanganSMP)
		kesenjangan int
		rasio       *float64
		kelas       KelasRasio
		label       Label
	}{
		{"lulusan melebihi kursi", func(*MasukanKesenjanganSMP) {}, 400, ptr(0.6), KelasKurang, LabelResmi},
		{"tanpa SMP negeri", func(m *MasukanKesenjanganSMP) { m.KursiNegeri = 0 }, 1000, ptr(0), KelasSangatKurang, LabelResmi},
		{"calon lulusan nol", func(m *MasukanKesenjanganSMP) { m.CalonLulusan = 0 }, -600, nil, KelasTanpaData, LabelResmi},
		{"calon lulusan estimasi", func(m *MasukanKesenjanganSMP) { m.LabelCalonLulusan = LabelEstimasi }, 400, ptr(0.6), KelasKurang, LabelEstimasi},
		{"kursi melebihi lulusan", func(m *MasukanKesenjanganSMP) { m.KursiNegeri = 1200 }, -200, ptr(1.2), KelasCukup, LabelResmi},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			m := masukanUji()
			k.ubah(&m)

			hasil := HitungKesenjanganKursiSMP(m, ambangUji)

			if hasil.Kesenjangan != k.kesenjangan || !samaRasio(hasil.Rasio, k.rasio) || hasil.Kelas != k.kelas || hasil.Label != k.label {
				t.Errorf("kesenjangan %d, rasio %v, kelas %s, label %s", hasil.Kesenjangan, nilaiRasio(hasil.Rasio), hasil.Kelas, hasil.Label)
			}
			if hasil.Kebutuhan != m.CalonLulusan || hasil.Kapasitas != m.KursiNegeri {
				t.Errorf("kebutuhan %d, kapasitas %d", hasil.Kebutuhan, hasil.Kapasitas)
			}
		})
	}
}

func TestAmbangRasio_Kelas_BatasTepatMasukKelasAtas(t *testing.T) {
	kasus := map[float64]KelasRasio{0.49: KelasSangatKurang, 0.5: KelasKurang, 0.8: KelasHampirCukup, 0.99: KelasHampirCukup, 1.0: KelasCukup}
	for rasio, harap := range kasus {
		if kelas := ambangUji.Kelas(ptr(rasio)); kelas != harap {
			t.Errorf("rasio %v: kelas %s, harap %s", rasio, kelas, harap)
		}
	}
}

func TestLabelTerlemah_UrutanResmiProksiEstimasiSimulasi(t *testing.T) {
	kasus := []struct {
		masukan []Label
		harap   Label
	}{
		{[]Label{LabelResmi, LabelResmi}, LabelResmi},
		{[]Label{LabelResmi, LabelProksi}, LabelProksi},
		{[]Label{LabelProksi, LabelEstimasi, LabelResmi}, LabelEstimasi},
		{[]Label{LabelEstimasi, LabelSimulasi}, LabelSimulasi},
		{[]Label{}, LabelResmi},
	}
	for _, k := range kasus {
		if hasil := LabelTerlemah(k.masukan...); hasil != k.harap {
			t.Errorf("LabelTerlemah(%v) = %s, harap %s", k.masukan, hasil, k.harap)
		}
	}
}

func TestHitungKesenjanganKursiSMP_MadrasahTanpaRombel_KursiMadrasahTidakTersedia(t *testing.T) {
	m := masukanUji()
	m.JumlahMTs = 3

	rincian := HitungKesenjanganKursiSMP(m, ambangUji).Rincian

	if rincian.KursiMadrasah.Nilai != nil || rincian.KesenjanganSemuaJalur.Nilai != nil {
		t.Errorf("kursi madrasah %v, kesenjangan semua jalur %v", rincian.KursiMadrasah.Nilai, rincian.KesenjanganSemuaJalur.Nilai)
	}
	if *rincian.KesenjanganDenganSwasta.Nilai != 350 {
		t.Errorf("kesenjangan dengan SPMB Bersama = %d", *rincian.KesenjanganDenganSwasta.Nilai)
	}
}

func TestHitungKesenjanganKursiSMP_MadrasahDenganRombel_ProksiRombelKaliBatas(t *testing.T) {
	m := masukanUji()
	m.JumlahMTs, m.JumlahMTsDenganRombel, m.RombelKelas7MTs = 2, 2, 5

	rincian := HitungKesenjanganKursiSMP(m, ambangUji).Rincian

	if *rincian.KursiMadrasah.Nilai != 5*batasRombelUji || rincian.KursiMadrasah.Label != LabelProksi {
		t.Errorf("kursi madrasah = %+v", rincian.KursiMadrasah)
	}
	if *rincian.KesenjanganSemuaJalur.Nilai != 350-5*batasRombelUji || rincian.KesenjanganSemuaJalur.Label != LabelProksi {
		t.Errorf("kesenjangan semua jalur = %+v", rincian.KesenjanganSemuaJalur)
	}
}

func ptr(nilai float64) *float64 { return &nilai }

func samaRasio(a, b *float64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a-*b < 1e-9 && *b-*a < 1e-9
}

func nilaiRasio(r *float64) any {
	if r == nil {
		return nil
	}
	return *r
}
