package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const PanjangMinRahasia = 32

var (
	ErrRahasiaTerlaluPendek = fmt.Errorf("rahasia token minimal %d byte", PanjangMinRahasia)
	kepalaHS256             = base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
)

type klaimJWT struct {
	Subjek      string       `json:"sub"`
	Peran       domain.Peran `json:"peran"`
	Diterbitkan int64        `json:"iat"`
	Kedaluwarsa int64        `json:"exp"`
}

type Akses struct {
	Rahasia  []byte
	Umur     time.Duration
	Sekarang func() time.Time
}

func Baru(rahasia string, umur time.Duration) (Akses, error) {
	if len(rahasia) < PanjangMinRahasia {
		return Akses{}, ErrRahasiaTerlaluPendek
	}
	return Akses{Rahasia: []byte(rahasia), Umur: umur, Sekarang: time.Now}, nil
}

func (a Akses) Terbitkan(k domain.KlaimAkses) (string, error) {
	sekarang := a.Sekarang()
	isi, err := json.Marshal(klaimJWT{
		Subjek: strconv.FormatInt(k.PenggunaID, 10), Peran: k.Peran,
		Diterbitkan: sekarang.Unix(), Kedaluwarsa: sekarang.Add(a.Umur).Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("serialisasi klaim akses: %w", err)
	}
	badan := kepalaHS256 + "." + base64.RawURLEncoding.EncodeToString(isi)
	return badan + "." + a.tandatangani(badan), nil
}

func (a Akses) Periksa(token string) (domain.KlaimAkses, error) {
	bagian := strings.Split(token, ".")
	if len(bagian) != 3 || bagian[0] != kepalaHS256 {
		return domain.KlaimAkses{}, domain.ErrTokenAksesTidakBerlaku
	}
	badan := bagian[0] + "." + bagian[1]
	if !hmac.Equal([]byte(bagian[2]), []byte(a.tandatangani(badan))) {
		return domain.KlaimAkses{}, domain.ErrTokenAksesTidakBerlaku
	}
	klaim, err := uraiKlaim(bagian[1])
	if err != nil || a.Sekarang().Unix() >= klaim.Kedaluwarsa {
		return domain.KlaimAkses{}, domain.ErrTokenAksesTidakBerlaku
	}
	id, err := strconv.ParseInt(klaim.Subjek, 10, 64)
	if err != nil {
		return domain.KlaimAkses{}, domain.ErrTokenAksesTidakBerlaku
	}
	return domain.KlaimAkses{PenggunaID: id, Peran: klaim.Peran}, nil
}

func (a Akses) tandatangani(badan string) string {
	mac := hmac.New(sha256.New, a.Rahasia)
	mac.Write([]byte(badan))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func uraiKlaim(isi string) (klaimJWT, error) {
	mentah, err := base64.RawURLEncoding.DecodeString(isi)
	if err != nil {
		return klaimJWT{}, err
	}
	var klaim klaimJWT
	if err := json.Unmarshal(mentah, &klaim); err != nil {
		return klaimJWT{}, err
	}
	if klaim.Kedaluwarsa == 0 {
		return klaimJWT{}, errors.New("klaim tanpa kedaluwarsa")
	}
	return klaim, nil
}
