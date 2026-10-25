package domain

import (
	"slices"
	"testing"
)

var ambangUji = AmbangRasio{SangatKurang: 0.5, Kurang: 0.8, HampirCukup: 1.0}

const batasRombelUji = 30

var aturanSMPUji = AturanKesenjangan{BentukTujuan: []BentukSekolah{BentukSMP}, BatasMuridPerRombel: batasRombelUji, Ambang: ambangUji}

func masukanUji() MasukanKesenjangan {
	return MasukanKesenjangan{
		KodeKecamatan: "31.72.06", CalonLulusan: 1000, LabelCalonLulusan: LabelResmi, MetodeCalonLulusan: MetodePerSekolah,
		KursiNegeri: map[BentukSekolah]int{BentukSMP: 600}, KursiSPMBBersama: map[BentukSekolah]int{BentukSMP: 50},
	}
}

func TestHitungKesenjanganKursi_KasusTepi(t *testing.T) {
	kasus := []struct {
		nama        string
		ubah        func(*MasukanKesenjangan)
		kesenjangan int
		rasio       *float64
		kelas       KelasRasio
		label       Label
	}{
		{"lulusan melebihi kursi", func(*MasukanKesenjangan) {}, 400, ptr(0.6), KelasKurang, LabelResmi},
		{"tanpa SMP negeri", func(m *MasukanKesenjangan) { m.KursiNegeri = map[BentukSekolah]int{} }, 1000, ptr(0), KelasSangatKurang, LabelResmi},
		{"calon lulusan nol", func(m *MasukanKesenjangan) { m.CalonLulusan = 0 }, -600, nil, KelasTanpaData, LabelResmi},
		{"calon lulusan estimasi", func(m *MasukanKesenjangan) { m.LabelCalonLulusan = LabelEstimasi }, 400, ptr(0.6), KelasKurang, LabelEstimasi},
		{"kursi melebihi lulusan", func(m *MasukanKesenjangan) { m.KursiNegeri[BentukSMP] = 1200 }, -200, ptr(1.2), KelasCukup, LabelResmi},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			m := masukanUji()
			k.ubah(&m)

			hasil := HitungKesenjanganKursi(m, aturanSMPUji)

			if hasil.Kesenjangan != k.kesenjangan || !samaRasio(hasil.Rasio, k.rasio) || hasil.Kelas != k.kelas || hasil.Label != k.label {
				t.Errorf("kesenjangan %d, rasio %v, kelas %s, label %s", hasil.Kesenjangan, nilaiRasio(hasil.Rasio), hasil.Kelas, hasil.Label)
			}
			if hasil.Kebutuhan != m.CalonLulusan || hasil.Kapasitas != m.KursiNegeri[BentukSMP] {
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

func TestHitungKesenjanganKursi_MadrasahTanpaRombel_KursiMadrasahTidakTersedia(t *testing.T) {
	m := masukanUji()
	m.JumlahMadrasah = 3

	rincian := HitungKesenjanganKursi(m, aturanSMPUji).Rincian

	if rincian.KursiMadrasah.Nilai != nil || rincian.KesenjanganSemuaJalur.Nilai != nil {
		t.Errorf("kursi madrasah %v, kesenjangan semua jalur %v", rincian.KursiMadrasah.Nilai, rincian.KesenjanganSemuaJalur.Nilai)
	}
	if *rincian.KesenjanganDenganSwasta.Nilai != 350 {
		t.Errorf("kesenjangan dengan SPMB Bersama = %d", *rincian.KesenjanganDenganSwasta.Nilai)
	}
}

func TestHitungKesenjanganKursi_MadrasahDenganRombel_ProksiRombelKaliBatas(t *testing.T) {
	m := masukanUji()
	m.JumlahMadrasah, m.JumlahMadrasahDenganRombel, m.RombelTingkatPertamaMadrasah = 2, 2, 5

	rincian := HitungKesenjanganKursi(m, aturanSMPUji).Rincian

	if *rincian.KursiMadrasah.Nilai != 5*batasRombelUji || rincian.KursiMadrasah.Label != LabelProksi {
		t.Errorf("kursi madrasah = %+v", rincian.KursiMadrasah)
	}
	if *rincian.KesenjanganSemuaJalur.Nilai != 350-5*batasRombelUji || rincian.KesenjanganSemuaJalur.Label != LabelProksi {
		t.Errorf("kesenjangan semua jalur = %+v", rincian.KesenjanganSemuaJalur)
	}
}

func TestHitungKesenjanganKursi_SatuBentukTujuan_TanpaRincianPerBentuk(t *testing.T) {
	rincian := HitungKesenjanganKursi(masukanUji(), aturanSMPUji).Rincian

	if rincian.KursiNegeriPerBentuk != nil || rincian.KursiSPMBBersamaPerBentuk != nil || rincian.BatasMuridPerRombel != batasRombelUji {
		t.Errorf("rincian = %+v", rincian)
	}
}

func TestHitungKesenjanganKursi_SMADanSMK_KapasitasGabunganDanRincianPerBentuk(t *testing.T) {
	m := MasukanKesenjangan{
		KodeKecamatan: "31.72.06", CalonLulusan: 2000, LabelCalonLulusan: LabelResmi, MetodeCalonLulusan: MetodePerSekolah,
		KursiNegeri:      map[BentukSekolah]int{BentukSMA: 576, BentukSMK: 324},
		KursiSPMBBersama: map[BentukSekolah]int{BentukSMK: 40},
		JumlahMadrasah:   1,
	}
	aturan := AturanKesenjangan{BentukTujuan: []BentukSekolah{BentukSMA, BentukSMK}, BatasMuridPerRombel: 36, Ambang: ambangUji}

	hasil := HitungKesenjanganKursi(m, aturan)

	if hasil.Kapasitas != 900 || hasil.Kesenjangan != 1100 || !samaRasio(hasil.Rasio, ptr(0.45)) || hasil.Kelas != KelasSangatKurang {
		t.Errorf("kapasitas %d, kesenjangan %d, rasio %v, kelas %s", hasil.Kapasitas, hasil.Kesenjangan, nilaiRasio(hasil.Rasio), hasil.Kelas)
	}
	r := hasil.Rincian
	harapNegeri := []KursiBentuk{{BentukSMA, 576}, {BentukSMK, 324}}
	harapSPMB := []KursiBentuk{{BentukSMA, 0}, {BentukSMK, 40}}
	if !slices.Equal(r.KursiNegeriPerBentuk, harapNegeri) || !slices.Equal(r.KursiSPMBBersamaPerBentuk, harapSPMB) {
		t.Errorf("per bentuk negeri %v, SPMB %v", r.KursiNegeriPerBentuk, r.KursiSPMBBersamaPerBentuk)
	}
	if *r.KursiSPMBBersama.Nilai != 40 || *r.KesenjanganDenganSwasta.Nilai != 1060 || r.KursiMadrasah.Nilai != nil {
		t.Errorf("SPMB %d, kesenjangan dengan SPMB %d, madrasah %v", *r.KursiSPMBBersama.Nilai, *r.KesenjanganDenganSwasta.Nilai, r.KursiMadrasah.Nilai)
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
