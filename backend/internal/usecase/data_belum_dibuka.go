package usecase

import "github.com/benditandayusaputra/rantai/backend/internal/domain"

type BacaDataBelumDibuka struct {
	Konfigurasi domain.KonfigurasiDataBelumDibuka
}

func (u BacaDataBelumDibuka) Daftar() domain.KonfigurasiDataBelumDibuka {
	return u.Konfigurasi
}

func (u BacaDataBelumDibuka) Permohonan(kode, bahasa string) (domain.Dokumen, error) {
	if err := periksaBahasa(bahasa); err != nil {
		return domain.Dokumen{}, err
	}
	dokumen, ada := u.Konfigurasi.SusunPermohonan(kode, bahasa)
	if !ada {
		return domain.Dokumen{}, domain.TidakDitemukan("data belum dibuka tidak dikenal: " + kode)
	}
	return dokumen, nil
}
