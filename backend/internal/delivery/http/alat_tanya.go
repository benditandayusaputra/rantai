package http

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const halamanPertama = 1

var (
	parameterKodeKecamatan = tanya.Parameter{Tipe: tanya.TipeTeks, Deskripsi: "Kode kecamatan Kemendagri, contoh 31.72.01"}
	parameterTema          = tanya.Parameter{Tipe: tanya.TipeTeks, Deskripsi: "Kode tema: pendidikan (SD ke SMP) atau pendidikan_sma_smk (SMP ke SMA/SMK)"}
	parameterTahun         = tanya.Parameter{Tipe: tanya.TipeBilangan, Deskripsi: "Tahun anggaran; kosongkan untuk tahun berjalan"}
)

type masukanTeks struct {
	Teks string `json:"teks"`
}

type masukanKode struct {
	Kode string `json:"kode"`
}

type masukanKecamatanTema struct {
	KodeKecamatan string `json:"kode_kecamatan"`
	Tema          string `json:"tema"`
}

type masukanDaftarPaket struct {
	KodeKecamatan string `json:"kode_kecamatan"`
	Tema          string `json:"tema"`
	Tahun         int    `json:"tahun"`
}

type masukanTemaTahun struct {
	Tema  string `json:"tema"`
	Tahun int    `json:"tahun"`
}

type masukanUmpanAnggota struct {
	IDAnggota string `json:"id_anggota"`
	Tag       string `json:"tag"`
}

type daftarPaketTanyaJSON struct {
	Total int                `json:"total"`
	Paket []paketRingkasJSON `json:"paket"`
}

type tanggungJawabTemaJSON struct {
	Tema          string                 `json:"tema"`
	Nama          domain.Teks            `json:"nama"`
	TanggungJawab []domain.TanggungJawab `json:"tanggung_jawab"`
	ProvenanceID  string                 `json:"provenance_id"`
}

func RegistryTanya(l Layanan) (tanya.Registry, error) {
	return tanya.BuatRegistry(
		alatCariWilayah(l.Wilayah),
		alatRingkasanKecamatan(l.Ringkasan),
		alatIndikatorTema(l.Indikator),
		alatSekolahDiKecamatan(l.Indikator),
		alatDaftarPaket(l.Paket),
		alatDetailPaket(l.Paket),
		alatBelanjaTema(l.Belanja),
		alatTanggungJawabTema(l.Tema),
		alatWakilDapil(l.Wakil),
		alatUmpanAnggota(l.Umpan),
	)
}

func spekKode(nama, deskripsi, deskripsiKode string) tanya.Spek {
	return tanya.Spek{
		Nama: nama, Deskripsi: deskripsi, Wajib: []string{"kode"},
		Parameter: map[string]tanya.Parameter{"kode": {Tipe: tanya.TipeTeks, Deskripsi: deskripsiKode}},
	}
}

func alatCariWilayah(l LayananWilayah) tanya.Alat {
	spek := tanya.Spek{
		Nama: "cari_wilayah", Deskripsi: "Cari kota, kecamatan, atau kelurahan DKI Jakarta berdasarkan nama.", Wajib: []string{"teks"},
		Parameter: map[string]tanya.Parameter{"teks": {Tipe: tanya.TipeTeks, Deskripsi: "Nama wilayah yang dicari"}},
	}
	return tanya.BuatAlat(spek, func(ctx context.Context, m masukanTeks) (any, error) {
		hasil, err := l.Cari(ctx, m.Teks)
		return keDaftarWilayahJSON(hasil), err
	})
}

func alatRingkasanKecamatan(l LayananRingkasanKecamatan) tanya.Alat {
	spek := spekKode("ringkasan_kecamatan", "Angka kunci tema, arus anggaran tahun berjalan, dan wakil dapil sebuah kecamatan.", parameterKodeKecamatan.Deskripsi)
	return tanya.BuatAlat(spek, func(ctx context.Context, m masukanKode) (any, error) {
		hasil, err := l.Ringkas(ctx, m.Kode)
		return keRingkasanKecamatanJSON(hasil), err
	})
}

func alatIndikatorTema(l LayananIndikator) tanya.Alat {
	spek := tanya.Spek{
		Nama: "indikator_tema", Deskripsi: "Kebutuhan, kapasitas, kesenjangan, dan rincian indikator tema di satu kecamatan.",
		Parameter: map[string]tanya.Parameter{"kode_kecamatan": parameterKodeKecamatan, "tema": parameterTema},
		Wajib:     []string{"kode_kecamatan", "tema"},
	}
	return tanya.BuatAlat(spek, func(ctx context.Context, m masukanKecamatanTema) (any, error) {
		hasil, err := l.DetailPendidikan(ctx, m.Tema, m.KodeKecamatan)
		return keDetailPendidikanJSON(hasil), err
	})
}

func alatSekolahDiKecamatan(l LayananIndikator) tanya.Alat {
	spek := tanya.Spek{
		Nama: "sekolah_di_kecamatan",
		Deskripsi: "Daftar sekolah jenjang tujuan di kecamatan beserta daya tampungnya: tema pendidikan untuk SMP dan MTs, " +
			"tema pendidikan_sma_smk untuk SMA, SMK, dan MA.",
		Parameter: map[string]tanya.Parameter{"kode_kecamatan": parameterKodeKecamatan, "tema": parameterTema},
		Wajib:     []string{"kode_kecamatan", "tema"},
	}
	return tanya.BuatAlat(spek, func(ctx context.Context, m masukanKecamatanTema) (any, error) {
		hasil, err := l.SekolahKecamatan(ctx, m.Tema, m.KodeKecamatan)
		return keDaftarSekolahJSON(hasil), err
	})
}

func alatDaftarPaket(l LayananPaket) tanya.Alat {
	spek := tanya.Spek{
		Nama: "daftar_paket", Deskripsi: "Paket pengadaan berlokasi di kecamatan, halaman pertama, bisa disaring tema dan tahun.",
		Parameter: map[string]tanya.Parameter{"kode_kecamatan": parameterKodeKecamatan, "tema": parameterTema, "tahun": parameterTahun},
		Wajib:     []string{"kode_kecamatan"},
	}
	return tanya.BuatAlat(spek, func(ctx context.Context, m masukanDaftarPaket) (any, error) {
		hasil, err := l.DaftarKecamatan(ctx, usecase.FilterPaketKecamatan{
			KodeKecamatan: m.KodeKecamatan, Tema: m.Tema, Tahun: m.Tahun, Halaman: halamanPertama,
		})
		return daftarPaketTanyaJSON{Total: hasil.Total, Paket: keDaftarPaketJSON(hasil.Paket)}, err
	})
}

func alatDetailPaket(l LayananPaket) tanya.Alat {
	spek := spekKode("detail_paket", "Rantai lengkap satu paket: rencana, pengadaan, kontrak, pemeliharaan, dan sinyal perlu ditanyakan.", "Kode paket dari daftar_paket")
	return tanya.BuatAlat(spek, func(ctx context.Context, m masukanKode) (any, error) {
		hasil, err := l.Ambil(ctx, m.Kode)
		return keDetailPaketJSON(hasil), err
	})
}

func alatBelanjaTema(l LayananBelanja) tanya.Alat {
	spek := tanya.Spek{
		Nama: "belanja_tema", Deskripsi: "Kesenjangan indikator tema dan belanja paket bertema per kecamatan di seluruh DKI.",
		Parameter: map[string]tanya.Parameter{"tema": parameterTema, "tahun": parameterTahun},
		Wajib:     []string{"tema"},
	}
	return tanya.BuatAlat(spek, func(ctx context.Context, m masukanTemaTahun) (any, error) {
		hasil, err := l.Hitung(ctx, m.Tema, m.Tahun)
		return keBelanjaTemaJSON(hasil), err
	})
}

func alatTanggungJawabTema(l LayananTema) tanya.Alat {
	spek := tanya.Spek{
		Nama: "tanggung_jawab_tema", Deskripsi: "Instansi pelaksana, komisi DPRD pengawas, dan dasar hukum untuk setiap masalah dalam tema.",
		Parameter: map[string]tanya.Parameter{"tema": parameterTema},
		Wajib:     []string{"tema"},
	}
	return tanya.BuatAlat(spek, func(ctx context.Context, m masukanTemaTahun) (any, error) {
		t, err := l.Ambil(ctx, m.Tema)
		return tanggungJawabTemaJSON{Tema: t.Kode, Nama: t.Nama, TanggungJawab: t.TanggungJawab, ProvenanceID: t.ProvenanceID}, err
	})
}

func alatWakilDapil(l LayananWakil) tanya.Alat {
	spek := tanya.Spek{
		Nama: "wakil_dapil", Deskripsi: "Dapil kecamatan dan anggota DPRD DKI-nya beserta komisi; anggota komisi pengawas tema didahulukan.",
		Parameter: map[string]tanya.Parameter{"kode_kecamatan": parameterKodeKecamatan, "tema": parameterTema},
		Wajib:     []string{"kode_kecamatan"},
	}
	return tanya.BuatAlat(spek, func(ctx context.Context, m masukanKecamatanTema) (any, error) {
		hasil, err := l.WakilKecamatan(ctx, m.KodeKecamatan, m.Tema)
		return keWakilKecamatanJSON(hasil), err
	})
}

func alatUmpanAnggota(l LayananUmpan) tanya.Alat {
	spek := tanya.Spek{
		Nama: "umpan_anggota", Deskripsi: "Kegiatan terbaru anggota DPRD dari rilis resmi DPRD dan media terverifikasi, dengan jenis sumber.",
		Parameter: map[string]tanya.Parameter{
			"id_anggota": {Tipe: tanya.TipeTeks, Deskripsi: "ID anggota dari wakil_dapil"},
			"tag":        {Tipe: tanya.TipeTeks, Deskripsi: "Tag isu opsional, contoh sekolah"},
		},
		Wajib: []string{"id_anggota"},
	}
	return tanya.BuatAlat(spek, func(ctx context.Context, m masukanUmpanAnggota) (any, error) {
		hasil, err := l.UmpanAnggota(ctx, domain.FilterUmpan{AnggotaID: m.IDAnggota, Tag: m.Tag, Halaman: halamanPertama})
		return keDaftarUmpanJSON(hasil.Item), err
	})
}
