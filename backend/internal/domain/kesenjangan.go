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

type MasukanKesenjanganSMP struct {
	KodeKecamatan          string
	CalonLulusan           int
	LabelCalonLulusan      Label
	MetodeCalonLulusan     MetodeCalonLulusan
	KursiNegeri            int
	KursiSPMBBersama       int
	JumlahMTs              int
	JumlahMTsDenganRombel  int
	RombelKelas7MTs        int
	BatasMuridPerRombelSMP int
}

type RincianKesenjanganSMP struct {
	MetodeCalonLulusan      MetodeCalonLulusan `json:"metode_calon_lulusan"`
	CalonLulusan            NilaiBerlabel      `json:"calon_lulusan"`
	KursiNegeri             NilaiBerlabel      `json:"kursi_negeri"`
	KursiSPMBBersama        NilaiBerlabel      `json:"kursi_spmb_bersama"`
	KursiMadrasah           NilaiBerlabel      `json:"kursi_madrasah"`
	KesenjanganDenganSwasta NilaiBerlabel      `json:"kesenjangan_dengan_spmb_bersama"`
	KesenjanganSemuaJalur   NilaiBerlabel      `json:"kesenjangan_semua_jalur"`
}

type IndikatorKesenjanganSMP struct {
	TemaKode      string
	KodeKecamatan string
	Periode       string
	Kebutuhan     int
	Kapasitas     int
	Kesenjangan   int
	Rasio         *float64
	Kelas         KelasRasio
	Label         Label
	Rincian       RincianKesenjanganSMP
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

func HitungKesenjanganKursiSMP(m MasukanKesenjanganSMP, ambang AmbangRasio) IndikatorKesenjanganSMP {
	label := LabelTerlemah(m.LabelCalonLulusan, LabelResmi)
	rasio := hitungRasio(m.KursiNegeri, m.CalonLulusan)
	return IndikatorKesenjanganSMP{
		KodeKecamatan: m.KodeKecamatan,
		Kebutuhan:     m.CalonLulusan,
		Kapasitas:     m.KursiNegeri,
		Kesenjangan:   m.CalonLulusan - m.KursiNegeri,
		Rasio:         rasio,
		Kelas:         ambang.Kelas(rasio),
		Label:         label,
		Rincian:       susunRincian(m, label),
	}
}

func hitungRasio(kapasitas, kebutuhan int) *float64 {
	if kebutuhan <= 0 {
		return nil
	}
	rasio := float64(kapasitas) / float64(kebutuhan)
	return &rasio
}

func susunRincian(m MasukanKesenjanganSMP, label Label) RincianKesenjanganSMP {
	kesenjanganSwasta := m.CalonLulusan - m.KursiNegeri - m.KursiSPMBBersama
	kursiMadrasah := kursiMadrasahProksi(m)
	return RincianKesenjanganSMP{
		MetodeCalonLulusan:      m.MetodeCalonLulusan,
		CalonLulusan:            NilaiBerlabel{Nilai: &m.CalonLulusan, Label: m.LabelCalonLulusan},
		KursiNegeri:             NilaiBerlabel{Nilai: &m.KursiNegeri, Label: LabelResmi},
		KursiSPMBBersama:        NilaiBerlabel{Nilai: &m.KursiSPMBBersama, Label: LabelResmi},
		KursiMadrasah:           NilaiBerlabel{Nilai: kursiMadrasah, Label: LabelProksi},
		KesenjanganDenganSwasta: NilaiBerlabel{Nilai: &kesenjanganSwasta, Label: label},
		KesenjanganSemuaJalur:   NilaiBerlabel{Nilai: kurangiJikaAda(kesenjanganSwasta, kursiMadrasah), Label: LabelTerlemah(label, LabelProksi)},
	}
}

func kurangiJikaAda(nilai int, pengurang *int) *int {
	if pengurang == nil {
		return nil
	}
	hasil := nilai - *pengurang
	return &hasil
}

func kursiMadrasahProksi(m MasukanKesenjanganSMP) *int {
	if m.JumlahMTs > 0 && m.JumlahMTsDenganRombel == 0 {
		return nil
	}
	kursi := m.RombelKelas7MTs * m.BatasMuridPerRombelSMP
	return &kursi
}
