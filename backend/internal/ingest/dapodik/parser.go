package dapodik

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

var (
	ErrDetailKosong    = errors.New("detail sekolah kosong")
	ErrStatusSumber    = errors.New("sumber mengembalikan status gagal")
	statusSumberSukses = 200
)

type ButirCari struct {
	SekolahID string `json:"sekolah_id"`
	NPSN      string `json:"npsn"`
	Nama      string `json:"nama"`
	Bentuk    string `json:"bentuk_pendidikan"`
}

type HalamanCari struct {
	Total   int         `json:"total"`
	Sekolah []ButirCari `json:"data"`
}

type DetailSekolah struct {
	Sekolah    domain.Sekolah
	SemesterID string
}

type amplopSumber[T any] struct {
	StatusCode int `json:"status_code"`
	Data       T   `json:"data"`
}

type detailSumber struct {
	NPSN       string   `json:"npsn"`
	Nama       string   `json:"nama"`
	Bentuk     string   `json:"bentuk_pendidikan"`
	Status     string   `json:"status_sekolah"`
	Kecamatan  string   `json:"kecamatan"`
	Lintang    *float64 `json:"lintang"`
	Bujur      *float64 `json:"bujur"`
	SemesterID string   `json:"semester_id"`
}

func UraiCari(isi []byte) (HalamanCari, error) {
	var halaman struct {
		StatusCode int `json:"status_code"`
		HalamanCari
	}
	if err := json.Unmarshal(isi, &halaman); err != nil {
		return HalamanCari{}, fmt.Errorf("urai daftar sekolah: %w", err)
	}
	if halaman.StatusCode != statusSumberSukses {
		return HalamanCari{}, fmt.Errorf("daftar sekolah: %w (%d)", ErrStatusSumber, halaman.StatusCode)
	}
	if halaman.Sekolah == nil {
		halaman.Sekolah = []ButirCari{}
	}
	return halaman.HalamanCari, nil
}

func UraiDetail(isi []byte) (DetailSekolah, error) {
	amplop, err := uraiAmplop[struct {
		Sekolah []detailSumber `json:"sekolah"`
	}](isi)
	if err != nil {
		return DetailSekolah{}, err
	}
	if len(amplop.Sekolah) == 0 {
		return DetailSekolah{}, ErrDetailKosong
	}
	return keDetailSekolah(amplop.Sekolah[0])
}

func keDetailSekolah(d detailSumber) (DetailSekolah, error) {
	bentuk := domain.BentukSekolah(d.Bentuk)
	if !bentuk.IsValid() {
		return DetailSekolah{}, domain.InputTidakValid("bentuk sekolah tidak didukung: " + d.Bentuk)
	}
	status, err := domain.StatusDariSumber(d.Status)
	if err != nil {
		return DetailSekolah{}, err
	}
	return DetailSekolah{
		Sekolah: domain.Sekolah{
			NPSN: d.NPSN, Nama: d.Nama, Bentuk: bentuk, Status: status,
			NamaKecamatanSumber: d.Kecamatan, Lokasi: titikSah(d.Lintang, d.Bujur),
		},
		SemesterID: d.SemesterID,
	}, nil
}

func titikSah(lintang, bujur *float64) *domain.Titik {
	if lintang == nil || bujur == nil || (*lintang == 0 && *bujur == 0) {
		return nil
	}
	return &domain.Titik{Lat: *lintang, Lng: *bujur}
}

func UraiRombel(isi []byte, bentuk domain.BentukSekolah) (map[int]int, error) {
	return uraiPerTingkat(isi, bentuk, []string{"rombel_%d"})
}

func UraiPesertaDidik(isi []byte, bentuk domain.BentukSekolah) (map[int]int, error) {
	return uraiPerTingkat(isi, bentuk, []string{"pd_tkt_%d_l", "pd_tkt_%d_p"})
}

func uraiPerTingkat(isi []byte, bentuk domain.BentukSekolah, polaKunci []string) (map[int]int, error) {
	baris, err := uraiAmplop[[]map[string]json.RawMessage](isi)
	if err != nil {
		return nil, err
	}
	hasil := map[int]int{}
	if len(baris) == 0 {
		return hasil, nil
	}
	pertama, terakhir := bentuk.RentangTingkat()
	for tingkat := pertama; tingkat <= terakhir; tingkat++ {
		for _, pola := range polaKunci {
			var nilai int
			if mentah, ada := baris[0][fmt.Sprintf(pola, tingkat)]; ada && json.Unmarshal(mentah, &nilai) == nil {
				hasil[tingkat] += nilai
			}
		}
	}
	return hasil, nil
}

func uraiAmplop[T any](isi []byte) (T, error) {
	var amplop amplopSumber[T]
	if err := json.Unmarshal(isi, &amplop); err != nil {
		return amplop.Data, fmt.Errorf("urai respons sekolah: %w", err)
	}
	if amplop.StatusCode != statusSumberSukses {
		return amplop.Data, fmt.Errorf("%w (%d)", ErrStatusSumber, amplop.StatusCode)
	}
	return amplop.Data, nil
}
