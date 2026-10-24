package domain

import "time"

type JalurAspirasi struct {
	Kode       string `yaml:"kode"`
	Nama       Teks   `yaml:"nama"`
	Peran      Teks   `yaml:"peran"`
	Cara       Teks   `yaml:"cara"`
	DasarHukum string `yaml:"dasar_hukum"`
	URL        string `yaml:"url"`
}

type KonfigurasiJalurAspirasi struct {
	DiverifikasiPada time.Time       `yaml:"diverifikasi_pada"`
	Jalur            []JalurAspirasi `yaml:"jalur"`
}

func (k KonfigurasiJalurAspirasi) Untuk(nilai map[string]string) []JalurAspirasi {
	hasil := make([]JalurAspirasi, len(k.Jalur))
	for i, j := range k.Jalur {
		j.Nama = j.Nama.isi(nilai)
		j.Peran = j.Peran.isi(nilai)
		j.Cara = j.Cara.isi(nilai)
		hasil[i] = j
	}
	return hasil
}

func (t Teks) isi(nilai map[string]string) Teks {
	return Teks{ID: IsiPlaceholder(t.ID, nilai), EN: IsiPlaceholder(t.EN, nilai)}
}
