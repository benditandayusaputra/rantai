package domain

type KategoriKesalahan int

const (
	KategoriTidakDitemukan KategoriKesalahan = iota + 1
	KategoriInputTidakValid
	KategoriTidakBerwenang
	KategoriTerlaluBanyak
)

type Kesalahan struct {
	Kategori KategoriKesalahan
	Pesan    string
}

func (k *Kesalahan) Error() string {
	return k.Pesan
}

func TidakDitemukan(pesan string) *Kesalahan {
	return &Kesalahan{Kategori: KategoriTidakDitemukan, Pesan: pesan}
}

func InputTidakValid(pesan string) *Kesalahan {
	return &Kesalahan{Kategori: KategoriInputTidakValid, Pesan: pesan}
}

func TidakBerwenang(pesan string) *Kesalahan {
	return &Kesalahan{Kategori: KategoriTidakBerwenang, Pesan: pesan}
}

func TerlaluBanyak(pesan string) *Kesalahan {
	return &Kesalahan{Kategori: KategoriTerlaluBanyak, Pesan: pesan}
}
