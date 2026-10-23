package domain

type BelanjaKecamatan struct {
	Kecamatan RingkasanWilayah
	KodeKota  string
	Indikator *IndikatorKecamatan
	Anggaran  *RingkasanAnggaran
}

type BelanjaLuarKecamatan struct {
	Wilayah  RingkasanWilayah
	Anggaran RingkasanAnggaran
}

type BelanjaTema struct {
	Tema          TemaTersimpan
	Tahun         int
	Kecamatan     []BelanjaKecamatan
	LuarKecamatan []BelanjaLuarKecamatan
}
