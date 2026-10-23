package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	PrefiksKodeSimulasi = "simulasi-"
	PenandaNamaSimulasi = "Simulasi"
	PeranAnggota        = "anggota"
)

var (
	ErrKodePaketBukanSimulasi     = errors.New("kode paket simulasi wajib berawalan " + PrefiksKodeSimulasi)
	ErrNamaSimulasiTanpaPenanda   = errors.New("nama penyedia dan instansi simulasi wajib memuat kata " + PenandaNamaSimulasi)
	ErrAkunDemoTertautAnggota     = errors.New("akun anggota demo tidak boleh ditautkan ke anggota sungguhan")
	ErrSimulasiTertautDataNyata   = errors.New("data simulasi tertaut ke paket, penyedia, atau instansi sungguhan")
	ErrKecamatanDemoBelumDihitung = errors.New("belum ada indikator kesenjangan kursi lengkap untuk kota demo")
)

type SkenarioDemo struct {
	KotaDemo      string        `yaml:"kota_demo"`
	KecamatanDemo string        `yaml:"kecamatan_demo"`
	PaketSimulasi PaketSimulasi `yaml:"paket_simulasi"`
	AkunDemo      []AkunDemo    `yaml:"akun_demo"`
}

type PaketSimulasi struct {
	Kode           string    `yaml:"kode"`
	Nama           string    `yaml:"nama"`
	Instansi       string    `yaml:"instansi"`
	Penyedia       string    `yaml:"penyedia"`
	Kategori       string    `yaml:"kategori"`
	Sektor         string    `yaml:"sektor"`
	TemaKode       []string  `yaml:"tema_kode"`
	JenisPengadaan string    `yaml:"jenis_pengadaan"`
	Metode         string    `yaml:"metode"`
	TahunAnggaran  int       `yaml:"tahun_anggaran"`
	Pagu           float64   `yaml:"pagu"`
	HPS            float64   `yaml:"hps"`
	NilaiKontrak   float64   `yaml:"nilai_kontrak"`
	Tahap          string    `yaml:"tahap"`
	TanggalKontrak time.Time `yaml:"tanggal_kontrak"`
	JadwalMulai    time.Time `yaml:"jadwal_mulai"`
	JadwalSelesai  time.Time `yaml:"jadwal_selesai"`
}

type AkunDemo struct {
	Peran     string `yaml:"peran"`
	AnggotaID string `yaml:"anggota_id"`
}

func (s SkenarioDemo) Validasi() error {
	p := s.PaketSimulasi
	if !strings.HasPrefix(p.Kode, PrefiksKodeSimulasi) {
		return ErrKodePaketBukanSimulasi
	}
	if !strings.Contains(p.Penyedia, PenandaNamaSimulasi) || !strings.Contains(p.Instansi, PenandaNamaSimulasi) {
		return ErrNamaSimulasiTanpaPenanda
	}
	for _, akun := range s.AkunDemo {
		if akun.Peran == PeranAnggota && akun.AnggotaID != "" {
			return ErrAkunDemoTertautAnggota
		}
	}
	return nil
}
