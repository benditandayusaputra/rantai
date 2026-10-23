package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobCalonLulusan = "calon_lulusan"

var ErrBelumAdaDataRombel = errors.New("belum ada data rombel sekolah, jalankan job dapodik lebih dulu")

type KonstantaPendidikan struct {
	BatasMuridPerRombelSMP int    `yaml:"batas_murid_per_rombel_smp"`
	PembagiEstimasiTotalSD int    `yaml:"estimasi_total_sd"`
	URLSumberSiswa         string `yaml:"url_sumber_siswa"`
}

type CatatanProvenance struct {
	Sumber  domain.SumberData
	URL     string
	Label   domain.Label
	Catatan string
}

type PembuatProvenance func(catatan CatatanProvenance) domain.Provenance

type CalonLulusanTercatat struct {
	domain.CalonLulusan
	ProvenanceID string
}

type DataCalonLulusan struct {
	TahunAjaran string
	Daftar      []CalonLulusanTercatat
	Provenance  []domain.Provenance
}

type PembacaSiswaSD interface {
	SiswaSDPerKecamatan(ctx context.Context) (string, []domain.SiswaSDKecamatan, error)
}

type PenyimpanCalonLulusan interface {
	SimpanCalonLulusan(ctx context.Context, data DataCalonLulusan) (HasilImpor, error)
}

type PelaporCalonLulusan interface {
	LaporkanKecamatanTanpaCalonLulusan(kode []string)
}

type HitungCalonLulusan struct {
	Pembaca        PembacaSiswaSD
	Penyimpan      PenyimpanCalonLulusan
	Pencatat       PencatatIngest
	Pelapor        PelaporCalonLulusan
	Konstanta      KonstantaPendidikan
	BuatProvenance PembuatProvenance
}

func (u HitungCalonLulusan) Jalankan(ctx context.Context) error {
	id, err := u.Pencatat.Mulai(ctx, JobCalonLulusan)
	if err != nil {
		return fmt.Errorf("catat mulai job %s: %w", JobCalonLulusan, err)
	}
	hasil, err := u.hitung(ctx)
	if err != nil {
		return errors.Join(err, u.Pencatat.Gagal(ctx, id, err))
	}
	return u.Pencatat.Selesai(ctx, id, hasil)
}

func (u HitungCalonLulusan) hitung(ctx context.Context) (HasilImpor, error) {
	tahunAjaran, daftarSiswa, err := u.Pembaca.SiswaSDPerKecamatan(ctx)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("baca siswa SD per kecamatan: %w", err)
	}
	if tahunAjaran == "" {
		return HasilImpor{}, ErrBelumAdaDataRombel
	}
	data, tanpaData := u.susunData(tahunAjaran, daftarSiswa)
	hasil, err := u.Penyimpan.SimpanCalonLulusan(ctx, data)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("simpan calon lulusan: %w", err)
	}
	u.Pelapor.LaporkanKecamatanTanpaCalonLulusan(tanpaData)
	return hasil, nil
}

func (u HitungCalonLulusan) susunData(tahunAjaran string, daftarSiswa []domain.SiswaSDKecamatan) (DataCalonLulusan, []string) {
	data := DataCalonLulusan{TahunAjaran: tahunAjaran, Daftar: []CalonLulusanTercatat{}, Provenance: []domain.Provenance{}}
	provenancePerMetode := map[domain.MetodeCalonLulusan]domain.Provenance{}
	tanpaData := []string{}
	for _, siswa := range daftarSiswa {
		calon, ada := domain.HitungCalonLulusan(siswa, u.Konstanta.PembagiEstimasiTotalSD)
		if !ada {
			tanpaData = append(tanpaData, siswa.KodeKecamatan)
			continue
		}
		p, sudahAda := provenancePerMetode[calon.Metode]
		if !sudahAda {
			p = u.BuatProvenance(u.catatanMetode(calon, tahunAjaran))
			provenancePerMetode[calon.Metode] = p
			data.Provenance = append(data.Provenance, p)
		}
		data.Daftar = append(data.Daftar, CalonLulusanTercatat{CalonLulusan: calon, ProvenanceID: p.ID})
	}
	return data, tanpaData
}

func (u HitungCalonLulusan) catatanMetode(calon domain.CalonLulusan, tahunAjaran string) CatatanProvenance {
	catatan := fmt.Sprintf("Jumlah siswa tingkat 6 seluruh SD di kecamatan, data Dapodik tahun ajaran %s. "+
		"MI tidak termasuk karena data rombel madrasah tidak tersedia di sumber.", tahunAjaran)
	if calon.Metode == domain.MetodeEstimasiTotalSD {
		catatan = domain.RumusEstimasiTotalSD(u.Konstanta.PembagiEstimasiTotalSD) + ". Data Dapodik tahun ajaran " + tahunAjaran + "."
	}
	return CatatanProvenance{Sumber: domain.SumberRantai, URL: u.Konstanta.URLSumberSiswa, Label: calon.Label, Catatan: catatan}
}
