package tanya

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type percakapanTiruan struct {
	giliran  []Giliran
	diterima [][]HasilAlat
	isiTanya string
	posisi   int
}

func (p *percakapanTiruan) berikutnya() (Giliran, error) {
	if p.posisi >= len(p.giliran) {
		return Giliran{}, errors.New("giliran tiruan habis")
	}
	p.posisi++
	return p.giliran[p.posisi-1], nil
}

func (p *percakapanTiruan) Tanya(_ context.Context, isi string) (Giliran, error) {
	p.isiTanya = isi
	return p.berikutnya()
}

func (p *percakapanTiruan) Balas(_ context.Context, hasil []HasilAlat) (Giliran, error) {
	p.diterima = append(p.diterima, hasil)
	return p.berikutnya()
}

type modelTiruan struct {
	percakapan *percakapanTiruan
}

func (m modelTiruan) BukaPercakapan(string, []Spek) Percakapan {
	return m.percakapan
}

func alatGagal() Alat {
	return BuatAlat(Spek{Nama: "gagal"}, func(context.Context, masukanGema) (any, error) {
		return nil, domain.TidakDitemukan("kecamatan tidak ditemukan")
	})
}

func penjawabUji(t *testing.T, giliran ...Giliran) (Penjawab, *percakapanTiruan) {
	t.Helper()
	registry, err := BuatRegistry(alatGema("gema"), alatGagal())
	if err != nil {
		t.Fatal(err)
	}
	percakapan := &percakapanTiruan{giliran: giliran}
	return Penjawab{Model: modelTiruan{percakapan: percakapan}, Alat: registry, Aturan: Aturan{BatasLangkah: 3, PanjangMaksPertanyaan: 500}}, percakapan
}

func panggil(nama string) Giliran {
	return Giliran{Panggilan: []PanggilanAlat{{ID: "toolu_" + nama, Nama: nama, Masukan: json.RawMessage(`{"teks":"halo"}`)}}}
}

func TestJawab_AlatLaluTeks_JawabanBerujukanDanKabarAlat(t *testing.T) {
	p, percakapan := penjawabUji(t, panggil("gema"), Giliran{Teks: "Ada 7360 kursi [prv_1]."})
	kabar := []string{}

	jawaban, err := p.Jawab(t.Context(), Permintaan{Pertanyaan: "Berapa kursi SMP?", KodeKecamatanKonteks: "31.72.01"}, func(k Kabar) {
		kabar = append(kabar, k.Alat)
	})

	if err != nil {
		t.Fatal(err)
	}
	if jawaban.Teks != "Ada 7360 kursi [1]." || !reflect.DeepEqual(jawaban.Rujukan, []Rujukan{{Nomor: 1, ProvenanceID: "prv_1"}}) {
		t.Errorf("jawaban = %+v", jawaban)
	}
	if len(jawaban.HasilAlat) != 1 || !reflect.DeepEqual(kabar, []string{"gema"}) {
		t.Errorf("hasil alat = %d, kabar = %v", len(jawaban.HasilAlat), kabar)
	}
	if !strings.Contains(percakapan.isiTanya, "<kecamatan_aktif>31.72.01</kecamatan_aktif>") {
		t.Errorf("konteks kecamatan tidak dikirim: %s", percakapan.isiTanya)
	}
}

func TestJawab_AlatMengembalikanKesalahanDomain_DikirimSebagaiHasilGagal(t *testing.T) {
	p, percakapan := penjawabUji(t, panggil("gagal"), Giliran{Teks: "Data belum tersedia."})

	_, err := p.Jawab(t.Context(), Permintaan{Pertanyaan: "Kursi di Atlantis?"}, func(Kabar) {})

	harap := []HasilAlat{{IDPanggilan: "toolu_gagal", Isi: "kecamatan tidak ditemukan", IsError: true}}
	if err != nil || len(percakapan.diterima) != 1 || !reflect.DeepEqual(percakapan.diterima[0], harap) {
		t.Errorf("err = %v, diterima = %+v", err, percakapan.diterima)
	}
}

func TestJawab_ModelTerusMemanggilAlat_BerhentiDiBatasLangkah(t *testing.T) {
	p, _ := penjawabUji(t, panggil("gema"), panggil("gema"), panggil("gema"), panggil("gema"))

	_, err := p.Jawab(t.Context(), Permintaan{Pertanyaan: "Ulangi terus"}, func(Kabar) {})

	if !errors.Is(err, ErrBatasLangkah) {
		t.Errorf("err = %v, harap ErrBatasLangkah", err)
	}
}

func TestPeriksa_PermintaanTidakSah_InputTidakValid(t *testing.T) {
	kasus := []struct {
		nama       string
		permintaan Permintaan
	}{
		{"kosong", Permintaan{Pertanyaan: "   "}},
		{"lebih dari 500 karakter", Permintaan{Pertanyaan: strings.Repeat("a", 501)}},
		{"bahasa tak dikenal", Permintaan{Pertanyaan: "Halo", Bahasa: "fr"}},
		{"kode kota", Permintaan{Pertanyaan: "Halo", KodeKecamatanKonteks: "31.72"}},
	}
	p, _ := penjawabUji(t)
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			var kesalahan *domain.Kesalahan
			if err := p.Periksa(k.permintaan); !errors.As(err, &kesalahan) || kesalahan.Kategori != domain.KategoriInputTidakValid {
				t.Errorf("err = %v, harap input tidak valid", err)
			}
		})
	}
}

func TestPeriksa_TepatLimaRatusKarakter_Diterima(t *testing.T) {
	p, _ := penjawabUji(t)

	if err := p.Periksa(Permintaan{Pertanyaan: strings.Repeat("é", 500), Bahasa: BahasaInggris}); err != nil {
		t.Errorf("err = %v, harap diterima", err)
	}
}

func TestSusunRujukan_ProvenanceBerulang_NomorSamaSesuaiUrutanMuncul(t *testing.T) {
	teks, rujukan := SusunRujukan("Kursi 288 [prv_b], lulusan 400 [prv_rantai_x], kursi lain [prv_b].")

	harap := []Rujukan{{Nomor: 1, ProvenanceID: "prv_b"}, {Nomor: 2, ProvenanceID: "prv_rantai_x"}}
	if teks != "Kursi 288 [1], lulusan 400 [2], kursi lain [1]." || !reflect.DeepEqual(rujukan, harap) {
		t.Errorf("teks = %q, rujukan = %+v", teks, rujukan)
	}
}
