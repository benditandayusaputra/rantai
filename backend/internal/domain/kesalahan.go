package domain

type KategoriKesalahan int

const (
	KategoriTidakDitemukan KategoriKesalahan = iota + 1
	KategoriInputTidakValid
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
