package main

import (
	"bufio"
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
)

const (
	pathSkenarioDemo  = "../data/snapshot/skenario_demo.yaml"
	pathCacheTanya    = "../data/snapshot/tanya_demo.json"
	alamatAPIBawaan   = "http://localhost:8080/api/v1"
	batasWaktuPerSoal = 3 * time.Minute
	izinBerkas        = 0o644
)

var (
	ErrKecamatanDemoKosong = errors.New("kecamatan_demo di skenario_demo.yaml masih kosong; jalankan make seed-demo lalu kunci kodenya")
	ErrTanpaJawaban        = errors.New("aliran berakhir tanpa peristiwa jawaban")
)

type peristiwaGalat struct {
	Pesan string `json:"pesan"`
}

type peristiwaAlat struct {
	Alat string `json:"alat"`
}

func main() {
	if err := jalankan(); err != nil {
		log.Fatal(err)
	}
}

func jalankan() error {
	var skenario domain.SkenarioDemo
	if err := config.BacaYAML(pathSkenarioDemo, &skenario); err != nil {
		return err
	}
	if skenario.KecamatanDemo == "" {
		return ErrKecamatanDemoKosong
	}
	alamat := cmp.Or(os.Getenv("RANTAI_API_URL"), alamatAPIBawaan)
	cache := make([]tanya.EntriCache, len(skenario.TanyaDemo))
	for i, pertanyaan := range skenario.TanyaDemo {
		entri, err := rekam(alamat, tanya.Permintaan{Pertanyaan: pertanyaan, KodeKecamatanKonteks: skenario.KecamatanDemo})
		if err != nil {
			return fmt.Errorf("rekam %q: %w", pertanyaan, err)
		}
		cache[i] = entri
		log.Printf("terekam: %s", pertanyaan)
	}
	isi, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pathCacheTanya, append(isi, '\n'), izinBerkas)
}

func rekam(alamat string, p tanya.Permintaan) (tanya.EntriCache, error) {
	ctx, batal := context.WithTimeout(context.Background(), batasWaktuPerSoal)
	defer batal()
	badan, err := json.Marshal(map[string]string{
		"pertanyaan": p.Pertanyaan, "kode_kecamatan_konteks": p.KodeKecamatanKonteks, "bahasa": tanya.BahasaIndonesia,
	})
	if err != nil {
		return tanya.EntriCache{}, err
	}
	permintaan, err := http.NewRequestWithContext(ctx, http.MethodPost, alamat+"/tanya", bytes.NewReader(badan))
	if err != nil {
		return tanya.EntriCache{}, err
	}
	permintaan.Header.Set("Content-Type", "application/json")
	respons, err := http.DefaultClient.Do(permintaan)
	if err != nil {
		return tanya.EntriCache{}, err
	}
	defer func() { _ = respons.Body.Close() }()
	if respons.StatusCode != http.StatusOK {
		return tanya.EntriCache{}, fmt.Errorf("status %s", respons.Status)
	}
	entri, err := bacaAliran(respons.Body)
	entri.Pertanyaan, entri.KodeKecamatanKonteks = p.Pertanyaan, p.KodeKecamatanKonteks
	return entri, err
}

func bacaAliran(r io.Reader) (tanya.EntriCache, error) {
	entri := tanya.EntriCache{Alat: []string{}}
	pemindai := bufio.NewScanner(r)
	var nama string
	for pemindai.Scan() {
		baris := pemindai.Text()
		switch {
		case strings.HasPrefix(baris, "event: "):
			nama = strings.TrimPrefix(baris, "event: ")
		case strings.HasPrefix(baris, "data: "):
			selesai, err := terapkan(&entri, nama, []byte(strings.TrimPrefix(baris, "data: ")))
			if err != nil || selesai {
				return entri, err
			}
		}
	}
	return entri, errors.Join(pemindai.Err(), ErrTanpaJawaban)
}

func terapkan(entri *tanya.EntriCache, nama string, data []byte) (bool, error) {
	switch nama {
	case "alat":
		var alat peristiwaAlat
		err := json.Unmarshal(data, &alat)
		entri.Alat = append(entri.Alat, alat.Alat)
		return false, err
	case "jawaban":
		return true, json.Unmarshal(data, entri)
	case "galat":
		var galat peristiwaGalat
		return true, errors.Join(json.Unmarshal(data, &galat), errors.New(galat.Pesan))
	default:
		return false, nil
	}
}
