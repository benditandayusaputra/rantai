package token

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const rahasiaUji = "rahasia-uji-yang-panjangnya-cukup-32"

type jam struct {
	saat time.Time
}

func (j *jam) sekarang() time.Time { return j.saat }

func aksesUji() (Akses, *jam) {
	j := &jam{saat: time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)}
	return Akses{Rahasia: []byte(rahasiaUji), Umur: 15 * time.Minute, Sekarang: j.sekarang}, j
}

func terbitkan(t *testing.T, a Akses) string {
	t.Helper()
	token, err := a.Terbitkan(domain.KlaimAkses{PenggunaID: 42, Peran: domain.PeranInstansi})
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestAkses_TerbitLaluPeriksa_KlaimSama(t *testing.T) {
	a, _ := aksesUji()

	klaim, err := a.Periksa(terbitkan(t, a))

	if err != nil || klaim != (domain.KlaimAkses{PenggunaID: 42, Peran: domain.PeranInstansi}) {
		t.Errorf("klaim %+v, err %v", klaim, err)
	}
}

func TestAkses_TokenTidakBerlaku(t *testing.T) {
	a, j := aksesUji()
	token := terbitkan(t, a)
	bagian := strings.Split(token, ".")
	kepalaNone := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))
	kasus := []struct {
		nama  string
		token func() string
	}{
		{"lewat lima belas menit", func() string { j.saat = j.saat.Add(15 * time.Minute); return token }},
		{"isi diubah", func() string {
			return bagian[0] + "." + base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"1","peran":"admin","exp":9999999999}`)) + "." + bagian[2]
		}},
		{"algoritma none", func() string { return kepalaNone + "." + bagian[1] + "." }},
		{"rahasia lain", func() string {
			lain := Akses{Rahasia: []byte(strings.Repeat("x", PanjangMinRahasia)), Umur: time.Minute, Sekarang: j.sekarang}
			return terbitkan(t, lain)
		}},
		{"bukan jwt", func() string { return "acak" }},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			a, j = aksesUji()
			if _, err := a.Periksa(k.token()); !errors.Is(err, domain.ErrTokenAksesTidakBerlaku) {
				t.Errorf("err %v, harap %v", err, domain.ErrTokenAksesTidakBerlaku)
			}
		})
	}
}

func TestBaru_RahasiaPendek_Ditolak(t *testing.T) {
	if _, err := Baru("pendek", time.Minute); !errors.Is(err, ErrRahasiaTerlaluPendek) {
		t.Errorf("err %v", err)
	}
}
