package wilayah

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var polaBaris = regexp.MustCompile(`\('(\d[\d.]*)','((?:[^']|'')*)',-?[\d.]+,-?[\d.]+,'(\[[^']*\])'\)`)

const (
	grupKode = 1 + iota
	grupNama
	grupPath
)

func UraiSQL(isi []byte) ([]domain.Wilayah, error) {
	daftar := []domain.Wilayah{}
	for _, baris := range polaBaris.FindAllSubmatch(isi, -1) {
		w, err := uraiBaris(string(baris[grupKode]), string(baris[grupNama]), baris[grupPath])
		if err != nil {
			return nil, err
		}
		daftar = append(daftar, w)
	}
	return daftar, nil
}

func uraiBaris(kode, nama string, path []byte) (domain.Wilayah, error) {
	tingkat, err := domain.TingkatDariKode(kode)
	if err != nil {
		return domain.Wilayah{}, err
	}
	geometri, err := uraiGeometri(path)
	if err != nil {
		return domain.Wilayah{}, fmt.Errorf("urai geometri %s: %w", kode, err)
	}
	return domain.Wilayah{
		Kode:      kode,
		Nama:      strings.ReplaceAll(nama, "''", "'"),
		Tingkat:   tingkat,
		IndukKode: domain.IndukKode(kode),
		Geometri:  geometri,
	}, nil
}

func uraiGeometri(path []byte) (domain.MultiPolygon, error) {
	var multi [][][][2]float64
	if err := json.Unmarshal(path, &multi); err == nil {
		return tukarLintangBujur(multi), nil
	}
	var poligon [][][2]float64
	if err := json.Unmarshal(path, &poligon); err != nil {
		return nil, err
	}
	return tukarLintangBujur([][][][2]float64{poligon}), nil
}

func tukarLintangBujur(multi [][][][2]float64) domain.MultiPolygon {
	hasil := make(domain.MultiPolygon, len(multi))
	for i, poligon := range multi {
		hasil[i] = make([][]domain.Koordinat, len(poligon))
		for j, cincin := range poligon {
			hasil[i][j] = make([]domain.Koordinat, len(cincin))
			for k, titik := range cincin {
				lintang, bujur := titik[0], titik[1]
				hasil[i][j][k] = domain.Koordinat{bujur, lintang}
			}
		}
	}
	return hasil
}
