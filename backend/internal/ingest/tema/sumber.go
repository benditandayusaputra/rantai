package tema

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/provenance"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"gopkg.in/yaml.v3"
)

const (
	polaBerkasTema = "*.yaml"
	catatanSumber  = "Konfigurasi tema disusun Rantai dari peraturan yang dirujuk di setiap entri tanggung jawab"
)

type Sumber struct {
	Dir        string
	URLRujukan string
	Sekarang   func() time.Time
}

func (s Sumber) Ambil(context.Context) (usecase.DataSumberTema, error) {
	berkas, err := filepath.Glob(filepath.Join(s.Dir, polaBerkasTema))
	if err != nil {
		return usecase.DataSumberTema{}, fmt.Errorf("cari berkas tema: %w", err)
	}
	daftar := []domain.Tema{}
	semuaIsi := []byte{}
	for _, path := range berkas {
		isi, err := os.ReadFile(path)
		if err != nil {
			return usecase.DataSumberTema{}, fmt.Errorf("baca %s: %w", path, err)
		}
		tema, err := UraiYAML(isi)
		if err != nil {
			return usecase.DataSumberTema{}, fmt.Errorf("urai %s: %w", path, err)
		}
		daftar = append(daftar, tema)
		semuaIsi = append(semuaIsi, isi...)
	}
	return usecase.DataSumberTema{Tema: daftar, Provenance: s.catatProvenance(semuaIsi)}, nil
}

func UraiYAML(isi []byte) (domain.Tema, error) {
	pengurai := yaml.NewDecoder(bytes.NewReader(isi))
	pengurai.KnownFields(true)
	var tema domain.Tema
	if err := pengurai.Decode(&tema); err != nil {
		return domain.Tema{}, err
	}
	return tema, nil
}

func (s Sumber) catatProvenance(isi []byte) domain.Provenance {
	return provenance.Baru(provenance.Catatan{
		Sumber:  domain.SumberRantai,
		URL:     s.URLRujukan,
		Label:   domain.LabelProksi,
		Catatan: catatanSumber,
	}, isi, s.Sekarang())
}
