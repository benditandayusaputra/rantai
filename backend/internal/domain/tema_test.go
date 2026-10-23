package domain

import (
	"errors"
	"testing"
)

var penghitungUji = []string{PenghitungKesenjanganKursiSMP}

func temaSahUji() Tema {
	return Tema{
		Kode:   "pendidikan",
		Nama:   Teks{ID: "Pendidikan", EN: "Education"},
		Aktif:  true,
		Urutan: 1,
		KonfigurasiTema: KonfigurasiTema{
			Indikator: Indikator{
				Penghitung:  PenghitungKesenjanganKursiSMP,
				Judul:       Teks{ID: "Kursi", EN: "Seats"},
				AmbangRasio: AmbangRasio{SangatKurang: 0.5, Kurang: 0.8, HampirCukup: 1.0},
			},
			FilterAnggaran: FilterAnggaran{KataKunci: []string{`\b(sdn|smpn)\b`}},
			TanggungJawab: []TanggungJawab{{
				Masalah:       "kursi_smp_kurang",
				Label:         Teks{ID: "Kursi kurang", EN: "Not enough seats"},
				UnitWilayah:   "Suku Dinas Pendidikan {kota}",
				Komisi:        "E",
				DasarHukum:    "UU 23/2014",
				DasarHukumURL: "https://peraturan.bpk.go.id/",
			}},
			TagIsu: []string{"sekolah", "daya tampung"},
			TemplatAspirasi: TemplatAspirasi{
				Judul:      Teks{ID: "Kursi di {kecamatan}", EN: "Seats in {kecamatan}"},
				Pertanyaan: []Teks{{ID: "Rencana {tahun_depan}?", EN: "Plan for {tahun_depan}?"}},
			},
		},
	}
}

func TestValidasiTema_KonfigurasiSah_TanpaError(t *testing.T) {
	if err := temaSahUji().Validasi(penghitungUji); err != nil {
		t.Errorf("harap sah, dapat %v", err)
	}
}

func TestValidasiTema_KonfigurasiRusak_ErrorSesuaiAturan(t *testing.T) {
	kasus := []struct {
		nama  string
		ubah  func(*Tema)
		harap error
	}{
		{"tanpa kode", func(t *Tema) { t.Kode = "" }, ErrTemaTanpaKode},
		{"nama tanpa bahasa Inggris", func(t *Tema) { t.Nama.EN = "" }, ErrTeksTidakLengkap},
		{"penghitung tidak terdaftar", func(t *Tema) { t.Indikator.Penghitung = "banjir" }, ErrPenghitungTidakTerdaftar},
		{"ambang menurun", func(t *Tema) { t.Indikator.AmbangRasio.Kurang = 0.4 }, ErrAmbangRasioTidakBerurutan},
		{"ambang sama", func(t *Tema) { t.Indikator.AmbangRasio.HampirCukup = 0.8 }, ErrAmbangRasioTidakBerurutan},
		{"ambang nol", func(t *Tema) { t.Indikator.AmbangRasio = AmbangRasio{} }, ErrAmbangRasioTidakBerurutan},
		{"filter kosong", func(t *Tema) { t.FilterAnggaran = FilterAnggaran{} }, ErrFilterAnggaranKosong},
		{"regex kata kunci rusak", func(t *Tema) { t.FilterAnggaran.KataKunci = []string{"(sdn"} }, ErrRegexFilterRusak},
		{"regex instansi rusak", func(t *Tema) { t.FilterAnggaran.InstansiPola = []string{"dinas["} }, ErrRegexFilterRusak},
		{"tanpa dasar hukum", func(t *Tema) { t.TanggungJawab[0].DasarHukum = "" }, ErrTanggungJawabTanpaDasarHukum},
		{"tanpa url dasar hukum", func(t *Tema) { t.TanggungJawab[0].DasarHukumURL = "" }, ErrTanggungJawabTanpaDasarHukum},
		{"url dasar hukum bukan http", func(t *Tema) { t.TanggungJawab[0].DasarHukumURL = "peraturan.bpk.go.id" }, ErrTanggungJawabTanpaDasarHukum},
		{"komisi F", func(t *Tema) { t.TanggungJawab[0].Komisi = "F" }, ErrKomisiTidakSah},
		{"tanpa tag isu", func(t *Tema) { t.TagIsu = nil }, ErrTagIsuTidakSah},
		{"tag isu huruf besar", func(t *Tema) { t.TagIsu = []string{"SPMB"} }, ErrTagIsuTidakSah},
		{"placeholder judul tak dikenal", func(t *Tema) { t.TemplatAspirasi.Judul.ID = "Kursi di {kelurahan}" }, ErrPlaceholderTidakDikenal},
		{"placeholder pertanyaan tak dikenal", func(t *Tema) { t.TemplatAspirasi.Pertanyaan[0].EN = "{tahun}" }, ErrPlaceholderTidakDikenal},
		{"placeholder unit wilayah tak dikenal", func(t *Tema) { t.TanggungJawab[0].UnitWilayah = "Sudin {wilayah}" }, ErrPlaceholderTidakDikenal},
		{"pertanyaan tanpa bahasa Indonesia", func(t *Tema) { t.TemplatAspirasi.Pertanyaan[0].ID = "" }, ErrTeksTidakLengkap},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			tema := temaSahUji()
			k.ubah(&tema)

			if err := tema.Validasi(penghitungUji); !errors.Is(err, k.harap) {
				t.Errorf("harap %v, dapat %v", k.harap, err)
			}
		})
	}
}

func TestValidasiTema_BeberapaPelanggaran_SemuaDilaporkan(t *testing.T) {
	tema := temaSahUji()
	tema.TanggungJawab[0].DasarHukum = ""
	tema.FilterAnggaran.KataKunci = []string{"(sdn"}

	err := tema.Validasi(penghitungUji)

	if !errors.Is(err, ErrTanggungJawabTanpaDasarHukum) || !errors.Is(err, ErrRegexFilterRusak) {
		t.Errorf("harap kedua pelanggaran dilaporkan, dapat %v", err)
	}
}
