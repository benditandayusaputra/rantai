package opentender

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	ukuranBarisMaks = 64 << 20
	panjangTahun    = 4
	peranPembeli    = "buyer"
)

type nilaiOCDS struct {
	Amount *float64 `json:"amount"`
}

type pihakOCDS struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	Roles                 []string `json:"roles"`
	AdditionalIdentifiers []struct {
		LegalName string `json:"legalName"`
	} `json:"additionalIdentifiers"`
}

type rilisOCDS struct {
	Buyer struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"buyer"`
	Tender struct {
		Title string `json:"title"`
	} `json:"tender"`
	Awards []struct {
		ID        string    `json:"id"`
		Date      string    `json:"date"`
		Title     string    `json:"title"`
		Value     nilaiOCDS `json:"value"`
		Suppliers []struct {
			Name string `json:"name"`
		} `json:"suppliers"`
	} `json:"awards"`
	Parties []pihakOCDS `json:"parties"`
}

func UraiOCDS(isiGzip io.Reader, idPembeli []string) ([]domain.RiwayatPenyedia, error) {
	baca, err := gzip.NewReader(isiGzip)
	if err != nil {
		return nil, fmt.Errorf("buka gzip ocds: %w", err)
	}
	pemindai := bufio.NewScanner(baca)
	pemindai.Buffer(make([]byte, 0, bufio.MaxScanTokenSize), ukuranBarisMaks)
	hasil := []domain.RiwayatPenyedia{}
	for pemindai.Scan() {
		var r rilisOCDS
		if err := json.Unmarshal(pemindai.Bytes(), &r); err != nil {
			return nil, fmt.Errorf("urai baris ocds: %w", err)
		}
		if slices.Contains(idPembeli, r.Buyer.ID) {
			hasil = append(hasil, riwayatDari(r)...)
		}
	}
	return hasil, pemindai.Err()
}

func riwayatDari(r rilisOCDS) []domain.RiwayatPenyedia {
	hasil := []domain.RiwayatPenyedia{}
	pembeli := satkerPembeli(r)
	for _, a := range r.Awards {
		tahun, err := strconv.Atoi(a.Date[:min(len(a.Date), panjangTahun)])
		if err != nil {
			continue
		}
		judul := a.Title
		if judul == "" {
			judul = r.Tender.Title
		}
		for i, s := range a.Suppliers {
			hasil = append(hasil, domain.RiwayatPenyedia{
				NamaPenyedia: s.Name, Tahun: tahun, Pembeli: pembeli, NamaPaket: judul,
				Nilai: a.Value.Amount, KodeSumber: fmt.Sprintf("%s#%d", a.ID, i),
			})
		}
	}
	return hasil
}

func satkerPembeli(r rilisOCDS) string {
	for _, p := range r.Parties {
		if p.ID == r.Buyer.ID && slices.Contains(p.Roles, peranPembeli) && len(p.AdditionalIdentifiers) > 0 && p.AdditionalIdentifiers[0].LegalName != "" {
			return p.AdditionalIdentifiers[0].LegalName
		}
	}
	return r.Buyer.Name
}
