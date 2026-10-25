package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const JobCalonLulusan = "calon_lulusan"

var ErrBelumAdaDataRombel = errors.New("belum ada data rombel sekolah, jalankan job dapodik lebih dulu")

type KonstantaPeralihan struct {
	BatasMuridPerRombel int `yaml:"batas_murid_per_rombel"`
	PembagiEstimasi     int `yaml:"pembagi_estimasi"`
}

type KonstantaPendidikan struct {
	URLSumberSiswa string                        `yaml:"url_sumber_siswa"`
	Peralihan      map[string]KonstantaPeralihan `yaml:"peralihan"`
}

func (k KonstantaPendidikan) Untuk(p domain.Peralihan) KonstantaPeralihan {
	return k.Peralihan[p.Kode]
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
	Jenjang     domain.Jenjang
	Daftar      []CalonLulusanTercatat
	Provenance  []domain.Provenance
}

type PembacaSiswa interface {
	SiswaPerKecamatan(ctx context.Context, peralihan domain.Peralihan) (string, []domain.SiswaKecamatan, error)
}

type PenyimpanCalonLulusan interface {
	SimpanCalonLulusan(ctx context.Context, data DataCalonLulusan) (HasilImpor, error)
}

type PelaporCalonLulusan interface {
	LaporkanKecamatanTanpaCalonLulusan(jenjang domain.Jenjang, kode []string)
}

type HitungCalonLulusan struct {
	Pembaca        PembacaSiswa
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
	total := HasilImpor{}
	for _, peralihan := range domain.DaftarPeralihan() {
		hasil, err := u.hitungPeralihan(ctx, peralihan)
		if err != nil {
			return HasilImpor{}, err
		}
		total.JumlahBaru += hasil.JumlahBaru
		total.JumlahBerubah += hasil.JumlahBerubah
	}
	return total, nil
}

func (u HitungCalonLulusan) hitungPeralihan(ctx context.Context, peralihan domain.Peralihan) (HasilImpor, error) {
	tahunAjaran, daftarSiswa, err := u.Pembaca.SiswaPerKecamatan(ctx, peralihan)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("baca siswa %s per kecamatan: %w", peralihan.BentukAsal, err)
	}
	if tahunAjaran == "" {
		return HasilImpor{}, ErrBelumAdaDataRombel
	}
	data, tanpaData := u.susunData(peralihan, tahunAjaran, daftarSiswa)
	hasil, err := u.Penyimpan.SimpanCalonLulusan(ctx, data)
	if err != nil {
		return HasilImpor{}, fmt.Errorf("simpan calon lulusan %s: %w", peralihan.BentukAsal, err)
	}
	u.Pelapor.LaporkanKecamatanTanpaCalonLulusan(peralihan.JenjangAsal, tanpaData)
	return hasil, nil
}

func (u HitungCalonLulusan) susunData(peralihan domain.Peralihan, tahunAjaran string, daftarSiswa []domain.SiswaKecamatan) (DataCalonLulusan, []string) {
	data := DataCalonLulusan{TahunAjaran: tahunAjaran, Jenjang: peralihan.JenjangAsal, Daftar: []CalonLulusanTercatat{}, Provenance: []domain.Provenance{}}
	estimasi := domain.AturanEstimasi{Metode: peralihan.MetodeEstimasi, Pembagi: u.Konstanta.Untuk(peralihan).PembagiEstimasi}
	provenancePerMetode := map[domain.MetodeCalonLulusan]domain.Provenance{}
	tanpaData := []string{}
	for _, siswa := range daftarSiswa {
		calon, ada := domain.HitungCalonLulusan(siswa, estimasi)
		if !ada {
			tanpaData = append(tanpaData, siswa.KodeKecamatan)
			continue
		}
		p, sudahAda := provenancePerMetode[calon.Metode]
		if !sudahAda {
			p = u.BuatProvenance(u.catatanMetode(peralihan, calon.Metode, tahunAjaran))
			provenancePerMetode[calon.Metode] = p
			data.Provenance = append(data.Provenance, p)
		}
		data.Daftar = append(data.Daftar, CalonLulusanTercatat{CalonLulusan: calon, ProvenanceID: p.ID})
	}
	return data, tanpaData
}

func (u HitungCalonLulusan) catatanMetode(peralihan domain.Peralihan, metode domain.MetodeCalonLulusan, tahunAjaran string) CatatanProvenance {
	catatan := fmt.Sprintf("Jumlah siswa tingkat %d seluruh %s di kecamatan, data Dapodik tahun ajaran %s. "+
		"%s tidak termasuk karena data rombel madrasah tidak tersedia di sumber.",
		peralihan.TingkatAsalTerakhir(), peralihan.BentukAsal, tahunAjaran, peralihan.BentukMadrasahAsal)
	label := domain.LabelResmi
	if metode == peralihan.MetodeEstimasi {
		catatan = peralihan.RumusEstimasi(u.Konstanta.Untuk(peralihan).PembagiEstimasi) + ". Data Dapodik tahun ajaran " + tahunAjaran + "."
		label = domain.LabelEstimasi
	}
	return CatatanProvenance{Sumber: domain.SumberRantai, URL: u.Konstanta.URLSumberSiswa, Label: label, Catatan: catatan}
}
