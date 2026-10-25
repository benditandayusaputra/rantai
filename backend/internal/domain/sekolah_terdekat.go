package domain

type SekolahTerdekat struct {
	SekolahDenganKursi
	JarakMeter float64
}

type KueriSekitarTitik struct {
	Titik       Titik
	RadiusMeter int
	Peralihan   Peralihan
}

type SekitarTitik struct {
	KueriSekitarTitik
	Sekolah           []SekolahTerdekat
	SiswaTingkatAkhir int
	ProvenanceSiswa   string
}

type SekolahTerdekatTerkelompok struct {
	Negeri      []SekolahTerdekat
	SPMBBersama []SekolahTerdekat
	Madrasah    []SekolahTerdekat
}

func (s SekitarTitik) Kelompokkan() SekolahTerdekatTerkelompok {
	hasil := SekolahTerdekatTerkelompok{Negeri: []SekolahTerdekat{}, SPMBBersama: []SekolahTerdekat{}, Madrasah: []SekolahTerdekat{}}
	for _, sekolah := range s.Sekolah {
		switch {
		case sekolah.Bentuk.IsMadrasah():
			hasil.Madrasah = append(hasil.Madrasah, sekolah)
		case sekolah.Status == StatusNegeri:
			hasil.Negeri = append(hasil.Negeri, sekolah)
		case sekolah.Jalur == JalurSPMBBersama:
			hasil.SPMBBersama = append(hasil.SPMBBersama, sekolah)
		}
	}
	return hasil
}

func KursiNegeriDalamRadius(negeri []SekolahTerdekat) (kursi int, provenanceID string) {
	for _, s := range negeri {
		if s.Kursi == nil {
			continue
		}
		kursi += *s.Kursi
		if provenanceID == "" {
			provenanceID = s.ProvenanceKursi
		}
	}
	return kursi, provenanceID
}

func TekananRadius(calonLulusan, kursiNegeri int) *float64 {
	if kursiNegeri <= 0 {
		return nil
	}
	tekanan := float64(calonLulusan) / float64(kursiNegeri)
	return &tekanan
}
