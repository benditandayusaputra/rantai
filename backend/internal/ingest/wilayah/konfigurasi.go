package wilayah

import (
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

type Konfigurasi struct {
	Sumber    KonfigurasiSumber                      `yaml:"sumber"`
	Provinsi  usecase.IdentitasProvinsi              `yaml:"provinsi"`
	Alias     map[domain.Tingkat]map[string][]string `yaml:"alias"`
	Singkatan map[string]string                      `yaml:"singkatan"`
}

func (k Konfigurasi) AturanAlias() domain.AturanAlias {
	return domain.AturanAlias{PrefiksPerTingkat: k.Alias, Singkatan: k.Singkatan}
}
