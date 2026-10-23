package domain

type JalurDayaTampung string

const (
	JalurNegeri      JalurDayaTampung = "negeri"
	JalurSPMBBersama JalurDayaTampung = "spmb_bersama"
)

type BarisDayaTampung struct {
	NPSN  string
	Nama  string
	Kursi int
}

func (j JalurDayaTampung) IsValid() bool {
	return j == JalurNegeri || j == JalurSPMBBersama
}

func PisahkanTerpetakan(daftar []BarisDayaTampung, npsnTerdaftar map[string]bool) (terpetakan, tidakTerpetakan []BarisDayaTampung) {
	terpetakan, tidakTerpetakan = []BarisDayaTampung{}, []BarisDayaTampung{}
	for _, baris := range daftar {
		if baris.NPSN != "" && npsnTerdaftar[baris.NPSN] {
			terpetakan = append(terpetakan, baris)
			continue
		}
		tidakTerpetakan = append(tidakTerpetakan, baris)
	}
	return terpetakan, tidakTerpetakan
}

func TotalKursi(daftar []BarisDayaTampung) int {
	total := 0
	for _, baris := range daftar {
		total += baris.Kursi
	}
	return total
}
