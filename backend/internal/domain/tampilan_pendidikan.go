package domain

type IndikatorTersimpan struct {
	IndikatorKesenjanganSMP
	NamaKecamatan string
	ProvenanceID  string
}

type FiturPetaIndikator struct {
	IndikatorTersimpan
	GeometriGeoJSON []byte
}

type DetailPendidikanKecamatan struct {
	IndikatorTersimpan
	ProvenanceCalonLulusan     string
	ProvenanceKursiNegeri      string
	ProvenanceKursiSPMBBersama string
	Tren                       []TitikTren
}

const MinimalPeriodeTren = 2

type TitikTren struct {
	Periode      string
	Rasio        float64
	Label        Label
	ProvenanceID string
}

func TrenLayakTampil(titik []TitikTren) []TitikTren {
	if len(titik) < MinimalPeriodeTren {
		return []TitikTren{}
	}
	return titik
}

type SekolahDenganKursi struct {
	NPSN            string
	Nama            string
	Bentuk          BentukSekolah
	Status          StatusSekolah
	Lokasi          *Titik
	Jalur           JalurDayaTampung
	Kursi           *int
	ProvenanceKursi string
}

const (
	KeterbatasanSPMBTidakDibatasiKecamatan = "spmb_tidak_dibatasi_kecamatan"
	KeterbatasanCalonLulusanTanpaMI        = "calon_lulusan_tanpa_mi"
	KeterbatasanKursiMadrasahTidakTersedia = "kursi_madrasah_tidak_tersedia"
	KeterbatasanTanpaDataLulusan           = "tidak_ada_data_lulusan"
	KeterbatasanJarakGarisLurus            = "jarak_garis_lurus"
)

type PetaTema struct {
	Tema  TemaTersimpan
	Fitur []FiturPetaIndikator
}

func (d DetailPendidikanKecamatan) Keterbatasan() []string {
	daftar := []string{KeterbatasanSPMBTidakDibatasiKecamatan}
	if d.Rincian.MetodeCalonLulusan == MetodePerSekolah {
		daftar = append(daftar, KeterbatasanCalonLulusanTanpaMI)
	}
	if d.Rincian.KursiMadrasah.Nilai == nil {
		daftar = append(daftar, KeterbatasanKursiMadrasahTidakTersedia)
	}
	if d.Kelas == KelasTanpaData {
		daftar = append(daftar, KeterbatasanTanpaDataLulusan)
	}
	return daftar
}
