package domain

type KelasRasio string

const (
	KelasSangatKurang KelasRasio = "sangat_kurang"
	KelasKurang       KelasRasio = "kurang"
	KelasHampirCukup  KelasRasio = "hampir_cukup"
	KelasCukup        KelasRasio = "cukup"
	KelasTanpaData    KelasRasio = "tidak_ada_data"
)

var urutanKekuatanLabel = map[Label]int{LabelResmi: 4, LabelProksi: 3, LabelEstimasi: 2, LabelSimulasi: 1}

type NilaiBerlabel struct {
	Nilai *int  `json:"nilai"`
	Label Label `json:"label"`
}

type MasukanKesenjangan struct {
	KodeKecamatan                string
	CalonLulusan                 int
	LabelCalonLulusan            Label
	MetodeCalonLulusan           MetodeCalonLulusan
	KursiNegeri                  map[BentukSekolah]int
	KursiSPMBBersama             map[BentukSekolah]int
	JumlahMadrasah               int
	JumlahMadrasahDenganRombel   int
	RombelTingkatPertamaMadrasah int
}

type AturanKesenjangan struct {
	BentukTujuan        []BentukSekolah
	BatasMuridPerRombel int
	Ambang              AmbangRasio
}

type KursiBentuk struct {
	Bentuk BentukSekolah `json:"bentuk"`
	Kursi  int           `json:"kursi"`
}

type RincianKesenjangan struct {
	MetodeCalonLulusan        MetodeCalonLulusan `json:"metode_calon_lulusan"`
	CalonLulusan              NilaiBerlabel      `json:"calon_lulusan"`
	KursiNegeri               NilaiBerlabel      `json:"kursi_negeri"`
	KursiSPMBBersama          NilaiBerlabel      `json:"kursi_spmb_bersama"`
	KursiMadrasah             NilaiBerlabel      `json:"kursi_madrasah"`
	KesenjanganDenganSwasta   NilaiBerlabel      `json:"kesenjangan_dengan_spmb_bersama"`
	KesenjanganSemuaJalur     NilaiBerlabel      `json:"kesenjangan_semua_jalur"`
	KursiNegeriPerBentuk      []KursiBentuk      `json:"kursi_negeri_per_bentuk,omitempty"`
	KursiSPMBBersamaPerBentuk []KursiBentuk      `json:"kursi_spmb_bersama_per_bentuk,omitempty"`
	BatasMuridPerRombel       int                `json:"batas_murid_per_rombel,omitempty"`
}

type IndikatorKesenjangan struct {
	TemaKode      string
	KodeKecamatan string
	Periode       string
	Kebutuhan     int
	Kapasitas     int
	Kesenjangan   int
	Rasio         *float64
	Kelas         KelasRasio
	Label         Label
	Rincian       RincianKesenjangan
}

func LabelTerlemah(daftar ...Label) Label {
	terlemah := LabelResmi
	for _, label := range daftar {
		if urutanKekuatanLabel[label] < urutanKekuatanLabel[terlemah] {
			terlemah = label
		}
	}
	return terlemah
}

func (a AmbangRasio) Kelas(rasio *float64) KelasRasio {
	switch {
	case rasio == nil:
		return KelasTanpaData
	case *rasio < a.SangatKurang:
		return KelasSangatKurang
	case *rasio < a.Kurang:
		return KelasKurang
	case *rasio < a.HampirCukup:
		return KelasHampirCukup
	default:
		return KelasCukup
	}
}

func HitungKesenjanganKursi(m MasukanKesenjangan, aturan AturanKesenjangan) IndikatorKesenjangan {
	label := LabelTerlemah(m.LabelCalonLulusan, LabelResmi)
	kursiNegeri := totalKursi(m.KursiNegeri)
	rasio := hitungRasio(kursiNegeri, m.CalonLulusan)
	return IndikatorKesenjangan{
		KodeKecamatan: m.KodeKecamatan,
		Kebutuhan:     m.CalonLulusan,
		Kapasitas:     kursiNegeri,
		Kesenjangan:   m.CalonLulusan - kursiNegeri,
		Rasio:         rasio,
		Kelas:         aturan.Ambang.Kelas(rasio),
		Label:         label,
		Rincian:       susunRincian(m, aturan, label),
	}
}

func hitungRasio(kapasitas, kebutuhan int) *float64 {
	if kebutuhan <= 0 {
		return nil
	}
	rasio := float64(kapasitas) / float64(kebutuhan)
	return &rasio
}

func susunRincian(m MasukanKesenjangan, aturan AturanKesenjangan, label Label) RincianKesenjangan {
	kursiNegeri, kursiSPMB := totalKursi(m.KursiNegeri), totalKursi(m.KursiSPMBBersama)
	kesenjanganSwasta := m.CalonLulusan - kursiNegeri - kursiSPMB
	kursiMadrasah := kursiMadrasahProksi(m, aturan.BatasMuridPerRombel)
	rincian := RincianKesenjangan{
		MetodeCalonLulusan:      m.MetodeCalonLulusan,
		CalonLulusan:            NilaiBerlabel{Nilai: &m.CalonLulusan, Label: m.LabelCalonLulusan},
		KursiNegeri:             NilaiBerlabel{Nilai: &kursiNegeri, Label: LabelResmi},
		KursiSPMBBersama:        NilaiBerlabel{Nilai: &kursiSPMB, Label: LabelResmi},
		KursiMadrasah:           NilaiBerlabel{Nilai: kursiMadrasah, Label: LabelProksi},
		KesenjanganDenganSwasta: NilaiBerlabel{Nilai: &kesenjanganSwasta, Label: label},
		KesenjanganSemuaJalur:   NilaiBerlabel{Nilai: kurangiJikaAda(kesenjanganSwasta, kursiMadrasah), Label: LabelTerlemah(label, LabelProksi)},
		BatasMuridPerRombel:     aturan.BatasMuridPerRombel,
	}
	if len(aturan.BentukTujuan) > 1 {
		rincian.KursiNegeriPerBentuk = kursiPerBentuk(m.KursiNegeri, aturan.BentukTujuan)
		rincian.KursiSPMBBersamaPerBentuk = kursiPerBentuk(m.KursiSPMBBersama, aturan.BentukTujuan)
	}
	return rincian
}

func totalKursi(perBentuk map[BentukSekolah]int) int {
	total := 0
	for _, kursi := range perBentuk {
		total += kursi
	}
	return total
}

func kursiPerBentuk(perBentuk map[BentukSekolah]int, urutan []BentukSekolah) []KursiBentuk {
	hasil := make([]KursiBentuk, len(urutan))
	for i, bentuk := range urutan {
		hasil[i] = KursiBentuk{Bentuk: bentuk, Kursi: perBentuk[bentuk]}
	}
	return hasil
}

func kurangiJikaAda(nilai int, pengurang *int) *int {
	if pengurang == nil {
		return nil
	}
	hasil := nilai - *pengurang
	return &hasil
}

func kursiMadrasahProksi(m MasukanKesenjangan, batasMuridPerRombel int) *int {
	if m.JumlahMadrasah > 0 && m.JumlahMadrasahDenganRombel == 0 {
		return nil
	}
	kursi := m.RombelTingkatPertamaMadrasah * batasMuridPerRombel
	return &kursi
}
