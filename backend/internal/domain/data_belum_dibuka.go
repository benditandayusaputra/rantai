package domain

type TemplatBagian struct {
	Judul Teks   `yaml:"judul"`
	Baris []Teks `yaml:"baris"`
}

type TemplatDokumen struct {
	Judul  Teks            `yaml:"judul"`
	Bagian []TemplatBagian `yaml:"bagian"`
}

type PPID struct {
	Nama   string `yaml:"nama"`
	Alamat string `yaml:"alamat"`
	URL    string `yaml:"url"`
}

type DataBelumDibuka struct {
	Kode           string `yaml:"kode"`
	Nama           Teks   `yaml:"nama"`
	Alasan         Teks   `yaml:"alasan"`
	Pengganti      Teks   `yaml:"pengganti"`
	FiturPengganti string `yaml:"fitur_pengganti"`
	Rincian        Teks   `yaml:"rincian"`
}

type KonfigurasiDataBelumDibuka struct {
	PPID       PPID              `yaml:"ppid"`
	Data       []DataBelumDibuka `yaml:"data"`
	Permohonan TemplatDokumen    `yaml:"permohonan"`
}

func (t TemplatDokumen) Isi(nilai map[string]string, bahasa string) Dokumen {
	dokumen := Dokumen{Judul: IsiPlaceholder(t.Judul.Pilih(bahasa), nilai), Bagian: make([]BagianDokumen, len(t.Bagian))}
	for i, bagian := range t.Bagian {
		baris := make([]string, len(bagian.Baris))
		for j, teks := range bagian.Baris {
			baris[j] = IsiPlaceholder(teks.Pilih(bahasa), nilai)
		}
		dokumen.Bagian[i] = BagianDokumen{Judul: bagian.Judul.Pilih(bahasa), Baris: baris}
	}
	return dokumen
}

func (k KonfigurasiDataBelumDibuka) SusunPermohonan(kode, bahasa string) (Dokumen, bool) {
	for _, d := range k.Data {
		if d.Kode == kode {
			return k.Permohonan.Isi(map[string]string{
				"nama": d.Nama.Pilih(bahasa), "rincian": d.Rincian.Pilih(bahasa),
				"ppid_nama": k.PPID.Nama, "ppid_alamat": k.PPID.Alamat, "ppid_url": k.PPID.URL,
			}, bahasa), true
		}
	}
	return Dokumen{}, false
}
