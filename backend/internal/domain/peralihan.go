package domain

import (
	"fmt"
	"slices"
	"strings"
)

const (
	PeralihanSDKeSMP     = "sd_ke_smp"
	PeralihanSMPKeSMASMK = "smp_ke_sma_smk"
)

type Peralihan struct {
	Kode                 string
	TemaKode             string
	Penghitung           string
	JenjangAsal          Jenjang
	BentukAsal           BentukSekolah
	BentukMadrasahAsal   BentukSekolah
	JenjangTujuan        Jenjang
	BentukTujuan         []BentukSekolah
	BentukMadrasahTujuan BentukSekolah
	MetodeEstimasi       MetodeCalonLulusan
}

var daftarPeralihan = []Peralihan{
	{
		Kode: PeralihanSDKeSMP, TemaKode: TemaPendidikan, Penghitung: PenghitungKesenjanganKursiSMP,
		JenjangAsal: JenjangDasar, BentukAsal: BentukSD, BentukMadrasahAsal: BentukMI,
		JenjangTujuan: JenjangMenengahPertama, BentukTujuan: []BentukSekolah{BentukSMP}, BentukMadrasahTujuan: BentukMTs,
		MetodeEstimasi: MetodeEstimasiTotalSD,
	},
	{
		Kode: PeralihanSMPKeSMASMK, TemaKode: TemaPendidikanSMASMK, Penghitung: PenghitungKesenjanganKursiSMASMK,
		JenjangAsal: JenjangMenengahPertama, BentukAsal: BentukSMP, BentukMadrasahAsal: BentukMTs,
		JenjangTujuan: JenjangMenengahAtas, BentukTujuan: []BentukSekolah{BentukSMA, BentukSMK}, BentukMadrasahTujuan: BentukMA,
		MetodeEstimasi: MetodeEstimasiTotalSMP,
	},
}

func DaftarPeralihan() []Peralihan {
	return slices.Clone(daftarPeralihan)
}

func PeralihanTema(temaKode string) (Peralihan, error) {
	for _, p := range daftarPeralihan {
		if p.TemaKode == temaKode {
			return p, nil
		}
	}
	return Peralihan{}, TidakDitemukan("tema tanpa data peralihan jenjang: " + temaKode)
}

func (p Peralihan) TingkatAsalTerakhir() int {
	_, terakhir := p.BentukAsal.RentangTingkat()
	return terakhir
}

func (p Peralihan) TingkatTujuanPertama() int {
	pertama, _ := p.BentukMadrasahTujuan.RentangTingkat()
	return pertama
}

func (p Peralihan) IsBanyakBentukTujuan() bool {
	return len(p.BentukTujuan) > 1
}

func (p Peralihan) NamaTujuan() string {
	nama := make([]string, len(p.BentukTujuan))
	for i, bentuk := range p.BentukTujuan {
		nama[i] = string(bentuk)
	}
	return strings.Join(nama, "/")
}

func (p Peralihan) RumusEstimasi(pembagi int) string {
	return fmt.Sprintf("calon lulusan = total siswa %s di kecamatan ÷ %d, dibulatkan; dipakai bila jumlah siswa kelas %d per sekolah tidak tersedia",
		p.BentukAsal, pembagi, p.TingkatAsalTerakhir())
}
