package domain

import (
	"fmt"
	"slices"
	"sort"
	"time"
)

const (
	SinyalSatuPeserta        JenisSinyal = "S-01"
	SinyalNilaiDiAtasSejenis JenisSinyal = "S-02"
	SinyalMepetBatasPL       JenisSinyal = "S-03"
	SinyalPenyediaDominan    JenisSinyal = "S-04"
	SinyalMelewatiJadwal     JenisSinyal = "S-05"
	persen                               = 100
)

type AmbangPengadaanLangsung struct {
	BerlakuMulai time.Time `yaml:"berlaku_mulai"`
	Nilai        float64   `yaml:"nilai"`
	Dasar        string    `yaml:"dasar"`
}

type AturanSinyal struct {
	MinPembanding int `yaml:"min_pembanding"`
	S01           struct {
		Metode []string `yaml:"metode"`
	} `yaml:"s01"`
	S02 struct {
		PengaliIQR float64 `yaml:"pengali_iqr"`
	} `yaml:"s02"`
	S03 struct {
		Metode         string                    `yaml:"metode"`
		JenisPengadaan string                    `yaml:"jenis_pengadaan"`
		JumlahMin      int                       `yaml:"jumlah_min"`
		JendelaHari    int                       `yaml:"jendela_hari"`
		PorsiAmbang    float64                   `yaml:"porsi_ambang"`
		Ambang         []AmbangPengadaanLangsung `yaml:"ambang"`
	} `yaml:"s03"`
	S04 struct {
		PorsiNilaiMin  float64 `yaml:"porsi_nilai_min"`
		JumlahPaketMin int     `yaml:"jumlah_paket_min"`
	} `yaml:"s04"`
	S05 struct {
		TenggangHari int `yaml:"tenggang_hari"`
	} `yaml:"s05"`
}

type PaketSinyal struct {
	Kode                         string
	Metode                       string
	JenisPengadaan               string
	Kategori                     string
	TahunAnggaran                int
	InstansiID                   int64
	PenyediaID                   *int64
	JumlahPeserta                *int
	HPS                          *float64
	NilaiKontrak                 *float64
	TanggalDibuat                *time.Time
	JadwalSelesai                *time.Time
	LaporanTerverifikasiTertahan int
}

type Sinyal struct {
	KodePaket  string
	Jenis      JenisSinyal
	Alasan     string
	Pembanding map[string]any
}

func (a AturanSinyal) Hitung(semua []PaketSinyal, hariIni time.Time) []Sinyal {
	hasil := []Sinyal{}
	for _, p := range semua {
		hasil = append(hasil, a.satuPeserta(p)...)
		hasil = append(hasil, a.melewatiJadwal(p, hariIni)...)
	}
	hasil = append(hasil, a.nilaiDiAtasSejenis(semua)...)
	hasil = append(hasil, a.mepetBatasPL(semua)...)
	return append(hasil, a.penyediaDominan(semua)...)
}

func (a AturanSinyal) satuPeserta(p PaketSinyal) []Sinyal {
	if !slices.Contains(a.S01.Metode, p.Metode) || p.JumlahPeserta == nil || *p.JumlahPeserta != 1 {
		return nil
	}
	return []Sinyal{{
		KodePaket:  p.Kode,
		Jenis:      SinyalSatuPeserta,
		Alasan:     fmt.Sprintf("Paket dengan metode %s ini hanya diikuti 1 peserta.", p.Metode),
		Pembanding: map[string]any{"jumlah_peserta": 1, "metode": p.Metode},
	}}
}

func (a AturanSinyal) melewatiJadwal(p PaketSinyal, hariIni time.Time) []Sinyal {
	if p.JadwalSelesai == nil || p.LaporanTerverifikasiTertahan == 0 {
		return nil
	}
	batas := p.JadwalSelesai.AddDate(0, 0, a.S05.TenggangHari)
	if !hariIni.After(batas) {
		return nil
	}
	return []Sinyal{{
		KodePaket: p.Kode,
		Jenis:     SinyalMelewatiJadwal,
		Alasan: fmt.Sprintf("Jadwal selesai %s sudah lewat lebih dari %d hari dan ada %d laporan warga terverifikasi bahwa pekerjaan belum selesai.",
			p.JadwalSelesai.Format(time.DateOnly), a.S05.TenggangHari, p.LaporanTerverifikasiTertahan),
		Pembanding: map[string]any{"jadwal_selesai": p.JadwalSelesai.Format(time.DateOnly), "laporan_terverifikasi": p.LaporanTerverifikasiTertahan},
	}}
}

func (a AturanSinyal) nilaiDiAtasSejenis(semua []PaketSinyal) []Sinyal {
	hasil := []Sinyal{}
	for _, kelompok := range kelompokkanPaket(semua, kunciSejenis) {
		statistik := HitungStatistik(nilaiHPS(kelompok))
		if statistik.Jumlah < a.MinPembanding {
			continue
		}
		batas := statistik.BatasAtas(a.S02.PengaliIQR)
		for _, p := range kelompok {
			if p.HPS != nil && *p.HPS > batas {
				hasil = append(hasil, sinyalNilaiTinggi(p, statistik, batas))
			}
		}
	}
	return hasil
}

func sinyalNilaiTinggi(p PaketSinyal, s Statistik, batas float64) Sinyal {
	return Sinyal{
		KodePaket: p.Kode,
		Jenis:     SinyalNilaiDiAtasSejenis,
		Alasan: fmt.Sprintf("HPS %s berada di atas batas sebaran %s dari %d paket kategori %s sejenis pada tahun %d (median %s).",
			FormatRupiah(*p.HPS), FormatRupiah(batas), s.Jumlah, p.Kategori, p.TahunAnggaran, FormatRupiah(s.Median)),
		Pembanding: map[string]any{"statistik": s, "batas_atas": batas},
	}
}

func (a AturanSinyal) mepetBatasPL(semua []PaketSinyal) []Sinyal {
	hasil := []Sinyal{}
	for _, kelompok := range kelompokkanPaket(semua, kunciInstansiKategori) {
		mepet := a.paketMepet(kelompok)
		for _, deret := range deretDalamJendela(mepet, a.S03.JendelaHari) {
			if len(deret) >= a.S03.JumlahMin {
				hasil = append(hasil, a.sinyalMepet(deret)...)
			}
		}
	}
	return unikPerPaket(hasil)
}

func (a AturanSinyal) paketMepet(kelompok []PaketSinyal) []PaketSinyal {
	mepet := []PaketSinyal{}
	for _, p := range kelompok {
		if p.Metode != a.S03.Metode || p.JenisPengadaan != a.S03.JenisPengadaan || p.HPS == nil || p.TanggalDibuat == nil {
			continue
		}
		if ambang, ada := a.ambangPada(*p.TanggalDibuat); ada && *p.HPS >= a.S03.PorsiAmbang*ambang.Nilai && *p.HPS <= ambang.Nilai {
			mepet = append(mepet, p)
		}
	}
	sort.Slice(mepet, func(i, j int) bool { return mepet[i].TanggalDibuat.Before(*mepet[j].TanggalDibuat) })
	return mepet
}

func (a AturanSinyal) ambangPada(t time.Time) (AmbangPengadaanLangsung, bool) {
	terpilih, ada := AmbangPengadaanLangsung{}, false
	for _, ambang := range a.S03.Ambang {
		if !t.Before(ambang.BerlakuMulai) && (!ada || ambang.BerlakuMulai.After(terpilih.BerlakuMulai)) {
			terpilih, ada = ambang, true
		}
	}
	return terpilih, ada
}

func deretDalamJendela(urut []PaketSinyal, jendelaHari int) [][]PaketSinyal {
	deret := [][]PaketSinyal{}
	for i := range urut {
		batas := urut[i].TanggalDibuat.AddDate(0, 0, jendelaHari)
		j := i
		for j < len(urut) && !urut[j].TanggalDibuat.After(batas) {
			j++
		}
		deret = append(deret, urut[i:j])
	}
	return deret
}

func (a AturanSinyal) sinyalMepet(deret []PaketSinyal) []Sinyal {
	ambang, _ := a.ambangPada(*deret[0].TanggalDibuat)
	kode := make([]string, len(deret))
	for i, p := range deret {
		kode[i] = p.Kode
	}
	hasil := make([]Sinyal, len(deret))
	for i, p := range deret {
		hasil[i] = Sinyal{
			KodePaket: p.Kode,
			Jenis:     SinyalMepetBatasPL,
			Alasan: fmt.Sprintf("%d paket pengadaan langsung kategori %s dari instansi yang sama dibuat dalam %d hari, masing-masing bernilai minimal %d%% dari batas pengadaan langsung %s (%s).",
				len(deret), p.Kategori, a.S03.JendelaHari, int(a.S03.PorsiAmbang*persen), FormatRupiah(ambang.Nilai), ambang.Dasar),
			Pembanding: map[string]any{"paket_terkait": kode, "batas_pengadaan_langsung": ambang.Nilai, "dasar": ambang.Dasar},
		}
	}
	return hasil
}

func (a AturanSinyal) penyediaDominan(semua []PaketSinyal) []Sinyal {
	hasil := []Sinyal{}
	for _, kelompok := range kelompokkanPaket(semua, kunciInstansiKategoriTahun) {
		if len(kelompok) < a.S04.JumlahPaketMin {
			continue
		}
		total, perPenyedia := nilaiPerPenyedia(kelompok)
		for _, p := range kelompok {
			if p.PenyediaID != nil && total > 0 && perPenyedia[*p.PenyediaID]/total >= a.S04.PorsiNilaiMin {
				hasil = append(hasil, sinyalDominan(p, perPenyedia[*p.PenyediaID]/total, len(kelompok)))
			}
		}
	}
	return hasil
}

func sinyalDominan(p PaketSinyal, porsi float64, jumlahPaket int) Sinyal {
	return Sinyal{
		KodePaket: p.Kode,
		Jenis:     SinyalPenyediaDominan,
		Alasan: fmt.Sprintf("Penyedia pemenang paket ini memenangkan %d%% nilai dari %d paket kategori %s di instansi yang sama pada tahun %d.",
			int(porsi*persen), jumlahPaket, p.Kategori, p.TahunAnggaran),
		Pembanding: map[string]any{"porsi_nilai": porsi, "jumlah_paket": jumlahPaket},
	}
}

func nilaiPerPenyedia(kelompok []PaketSinyal) (float64, map[int64]float64) {
	total, perPenyedia := 0.0, map[int64]float64{}
	for _, p := range kelompok {
		nilai := nilaiPaket(p)
		total += nilai
		if p.PenyediaID != nil {
			perPenyedia[*p.PenyediaID] += nilai
		}
	}
	return total, perPenyedia
}

func nilaiPaket(p PaketSinyal) float64 {
	switch {
	case p.NilaiKontrak != nil:
		return *p.NilaiKontrak
	case p.HPS != nil:
		return *p.HPS
	default:
		return 0
	}
}

func nilaiHPS(kelompok []PaketSinyal) []float64 {
	nilai := []float64{}
	for _, p := range kelompok {
		if p.HPS != nil {
			nilai = append(nilai, *p.HPS)
		}
	}
	return nilai
}

func kunciSejenis(p PaketSinyal) string {
	return fmt.Sprintf("%s|%s|%d", p.Kategori, p.JenisPengadaan, p.TahunAnggaran)
}

func kunciInstansiKategori(p PaketSinyal) string {
	return fmt.Sprintf("%d|%s", p.InstansiID, p.Kategori)
}

func kunciInstansiKategoriTahun(p PaketSinyal) string {
	return fmt.Sprintf("%d|%s|%d", p.InstansiID, p.Kategori, p.TahunAnggaran)
}

func kelompokkanPaket(semua []PaketSinyal, kunci func(PaketSinyal) string) map[string][]PaketSinyal {
	kelompok := map[string][]PaketSinyal{}
	for _, p := range semua {
		if p.Kategori != "" {
			kelompok[kunci(p)] = append(kelompok[kunci(p)], p)
		}
	}
	return kelompok
}

func unikPerPaket(daftar []Sinyal) []Sinyal {
	terlihat := map[string]bool{}
	hasil := []Sinyal{}
	for _, s := range daftar {
		if !terlihat[s.KodePaket] {
			terlihat[s.KodePaket] = true
			hasil = append(hasil, s)
		}
	}
	return hasil
}
