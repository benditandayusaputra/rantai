package domain

import (
	"errors"
	"fmt"
	"slices"
)

const (
	JabatanKoordinator      = "Koordinator"
	AlatKelengkapanPimpinan = "Pimpinan DPRD"
)

var ErrKomisiTanpaBidang = errors.New("komisi wajib punya kode A sampai E dan minimal satu bidang")

type Komisi struct {
	Kode   string
	Nama   string
	Bidang []string
}

type JabatanAKD struct {
	AnggotaID       string
	AlatKelengkapan string
	Jabatan         string
	KodeKomisi      string
}

type HalamanStruktur struct {
	AlatKelengkapan string
	KodeKomisi      string
	Baris           []BarisStruktur
}

func (k Komisi) Validasi() error {
	if !slices.Contains(komisiSah, k.Kode) || len(k.Bidang) == 0 {
		return fmt.Errorf("komisi %q: %w", k.Kode, ErrKomisiTanpaBidang)
	}
	return nil
}

func CocokkanStruktur(anggota []Anggota, halaman HalamanStruktur) ([]JabatanAKD, []string) {
	perKunci := indeksKunciNama(anggota)
	perID := kelompokkan(anggota, func(a Anggota) string { return a.ID })
	cocok, tidakCocok := []JabatanAKD{}, []string{}
	for _, b := range halaman.Baris {
		id, ada := cariAnggota(b, perKunci, perID)
		if !ada {
			tidakCocok = append(tidakCocok, fmt.Sprintf("%s: %q (%s)", halaman.AlatKelengkapan, b.Nama, b.Jabatan))
			continue
		}
		cocok = append(cocok, JabatanAKD{AnggotaID: id, AlatKelengkapan: halaman.AlatKelengkapan, Jabatan: b.Jabatan, KodeKomisi: halaman.KodeKomisi})
	}
	return cocok, tidakCocok
}

func indeksKunciNama(anggota []Anggota) map[string][]string {
	perKunci := map[string][]string{}
	for _, a := range anggota {
		kunci := []string{}
		for _, nama := range append([]string{a.Nama}, a.NamaLain...) {
			kunci = append(kunci, KunciNama(nama))
		}
		slices.Sort(kunci)
		for _, k := range slices.Compact(kunci) {
			perKunci[k] = append(perKunci[k], a.ID)
		}
	}
	return perKunci
}

func cariAnggota(b BarisStruktur, perKunci map[string][]string, perID map[string][]Anggota) (string, bool) {
	if calon := perKunci[KunciNama(b.Nama)]; len(calon) == 1 {
		return calon[0], true
	}
	if b.Slug != "" && len(perID[b.Slug]) == 1 {
		return b.Slug, true
	}
	return "", false
}

func TentukanKomisi(jabatan []JabatanAKD) (map[string]string, []string) {
	pimpinan := idPimpinan(jabatan)
	komisi := map[string]string{}
	ganda := map[string]bool{}
	for _, j := range jabatan {
		if j.KodeKomisi == "" || j.Jabatan == JabatanKoordinator || pimpinan[j.AnggotaID] {
			continue
		}
		if kode, ada := komisi[j.AnggotaID]; ada && kode != j.KodeKomisi {
			ganda[j.AnggotaID] = true
		}
		komisi[j.AnggotaID] = j.KodeKomisi
	}
	laporan := []string{}
	for id := range ganda {
		delete(komisi, id)
		laporan = append(laporan, id+" tercatat di lebih dari satu komisi")
	}
	slices.Sort(laporan)
	return komisi, laporan
}

func idPimpinan(jabatan []JabatanAKD) map[string]bool {
	pimpinan := map[string]bool{}
	for _, j := range jabatan {
		if j.AlatKelengkapan == AlatKelengkapanPimpinan {
			pimpinan[j.AnggotaID] = true
		}
	}
	return pimpinan
}

func AnggotaTanpaKomisi(anggota []Anggota) []string {
	tanpa := []string{}
	for _, a := range anggota {
		if a.KodeKomisi == "" {
			tanpa = append(tanpa, a.ID)
		}
	}
	return tanpa
}
