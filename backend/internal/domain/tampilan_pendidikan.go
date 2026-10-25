package domain

type IndikatorTersimpan struct {
	IndikatorKesenjangan
	NamaKecamatan string
	ProvenanceID  string
}

type FiturPetaIndikator struct {
	IndikatorTersimpan
	GeometriGeoJSON []byte
}

type DetailPendidikanKecamatan struct {
	IndikatorTersimpan
	Peralihan              Peralihan
	Ambang                 AmbangRasio
	ProvenanceCalonLulusan string
	ProvenanceKursi        []ProvenanceKursi
	Tren                   []TitikTren
}

type ProvenanceKursi struct {
	Jalur        JalurDayaTampung
	Bentuk       BentukSekolah
	ProvenanceID string
}

func (d DetailPendidikanKecamatan) ProvenanceKursiBentuk(jalur JalurDayaTampung, bentuk BentukSekolah) string {
	for _, p := range d.ProvenanceKursi {
		if p.Jalur == jalur && p.Bentuk == bentuk {
			return p.ProvenanceID
		}
	}
	return ""
}

func (d DetailPendidikanKecamatan) ProvenanceKursiJalur(jalur JalurDayaTampung) string {
	if d.Peralihan.IsBanyakBentukTujuan() || len(d.Peralihan.BentukTujuan) == 0 {
		return ""
	}
	return d.ProvenanceKursiBentuk(jalur, d.Peralihan.BentukTujuan[0])
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
	Foto            Foto
}

const (
	KeterbatasanSPMBTidakDibatasiKecamatan = "spmb_tidak_dibatasi_kecamatan"
	KeterbatasanCalonLulusanTanpaMI        = "calon_lulusan_tanpa_mi"
	KeterbatasanCalonLulusanTanpaMTs       = "calon_lulusan_tanpa_mts"
	KeterbatasanKursiMadrasahTidakTersedia = "kursi_madrasah_tidak_tersedia"
	KeterbatasanKursiMATidakTersedia       = "kursi_ma_tidak_tersedia"
	KeterbatasanSMASMKLintasKecamatan      = "sma_smk_lintas_kecamatan"
	KeterbatasanTanpaDataLulusan           = "tidak_ada_data_lulusan"
	KeterbatasanJarakGarisLurus            = "jarak_garis_lurus"
)

type PetaTema struct {
	Tema  TemaTersimpan
	Fitur []FiturPetaIndikator
}

var (
	keterbatasanCalonLulusanTanpaMadrasah = map[BentukSekolah]string{
		BentukMI: KeterbatasanCalonLulusanTanpaMI, BentukMTs: KeterbatasanCalonLulusanTanpaMTs,
	}
	keterbatasanKursiMadrasah = map[BentukSekolah]string{
		BentukMTs: KeterbatasanKursiMadrasahTidakTersedia, BentukMA: KeterbatasanKursiMATidakTersedia,
	}
)

func (p Peralihan) KeterbatasanCalonLulusan() string {
	return keterbatasanCalonLulusanTanpaMadrasah[p.BentukMadrasahAsal]
}

func (d DetailPendidikanKecamatan) Keterbatasan() []string {
	daftar := []string{KeterbatasanSPMBTidakDibatasiKecamatan}
	if d.Peralihan.IsBanyakBentukTujuan() {
		daftar = append(daftar, KeterbatasanSMASMKLintasKecamatan)
	}
	if d.Rincian.MetodeCalonLulusan == MetodePerSekolah {
		daftar = append(daftar, d.Peralihan.KeterbatasanCalonLulusan())
	}
	if d.Rincian.KursiMadrasah.Nilai == nil {
		daftar = append(daftar, keterbatasanKursiMadrasah[d.Peralihan.BentukMadrasahTujuan])
	}
	if d.Kelas == KelasTanpaData {
		daftar = append(daftar, KeterbatasanTanpaDataLulusan)
	}
	return daftar
}
