package domain

type SekolahTerdekat struct {
	SekolahDenganKursi
	JarakMeter float64
}

type SekitarTitik struct {
	Titik           Titik
	RadiusMeter     int
	Sekolah         []SekolahTerdekat
	SiswaKelas6     int
	ProvenanceSiswa string
}

type SekolahTerdekatTerkelompok struct {
	SMPNegeri   []SekolahTerdekat
	SPMBBersama []SekolahTerdekat
	Madrasah    []SekolahTerdekat
}

func (s SekitarTitik) Kelompokkan() SekolahTerdekatTerkelompok {
	hasil := SekolahTerdekatTerkelompok{SMPNegeri: []SekolahTerdekat{}, SPMBBersama: []SekolahTerdekat{}, Madrasah: []SekolahTerdekat{}}
	for _, sekolah := range s.Sekolah {
		switch {
		case sekolah.Bentuk.IsMadrasah():
			hasil.Madrasah = append(hasil.Madrasah, sekolah)
		case sekolah.Status == StatusNegeri:
			hasil.SMPNegeri = append(hasil.SMPNegeri, sekolah)
		case sekolah.Jalur == JalurSPMBBersama:
			hasil.SPMBBersama = append(hasil.SPMBBersama, sekolah)
		}
	}
	return hasil
}

func KursiNegeriDalamRadius(smpNegeri []SekolahTerdekat) (kursi int, provenanceID string) {
	for _, s := range smpNegeri {
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

func TekananRadius(siswaKelas6, kursiNegeri int) *float64 {
	if kursiNegeri <= 0 {
		return nil
	}
	tekanan := float64(siswaKelas6) / float64(kursiNegeri)
	return &tekanan
}
