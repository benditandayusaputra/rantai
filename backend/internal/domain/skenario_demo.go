package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	PrefiksKodeSimulasi = "simulasi-"
	PenandaNamaSimulasi = "Simulasi"
	AkhiranEmailUji     = ".test"
)

var (
	ErrKodePaketBukanSimulasi     = errors.New("kode paket simulasi wajib berawalan " + PrefiksKodeSimulasi)
	ErrNamaSimulasiTanpaPenanda   = errors.New("nama penyedia dan instansi simulasi wajib memuat kata " + PenandaNamaSimulasi)
	ErrAkunDemoTertautAnggota     = errors.New("akun anggota demo tidak boleh ditautkan ke anggota sungguhan")
	ErrAkunDemoBukanSimulasi      = errors.New("akun demo wajib beremail domain " + AkhiranEmailUji + " dan bernama tampilan memuat kata " + PenandaNamaSimulasi)
	ErrAkunInstansiTanpaInstansi  = errors.New("akun demo berperan instansi wajib menyebut instansi")
	ErrSimulasiTertautDataNyata   = errors.New("data simulasi tertaut ke paket, penyedia, atau instansi sungguhan")
	ErrKecamatanDemoBelumDihitung = errors.New("belum ada indikator kesenjangan kursi lengkap untuk kota demo")
	ErrLaporanSimulasiTidakSah    = errors.New("laporan simulasi wajib berkondisi sah, berfoto 1 sampai 4, dan cukup saksi untuk terverifikasi")
	ErrPelaporBukanAkunDemo       = errors.New("pelapor dan saksi laporan simulasi wajib akun demo yang berbeda")
)

const JumlahKonfirmasiTerverifikasi = 3

type SkenarioDemo struct {
	KotaDemo        string            `yaml:"kota_demo"`
	KecamatanDemo   string            `yaml:"kecamatan_demo"`
	PaketSimulasi   PaketSimulasi     `yaml:"paket_simulasi"`
	AkunDemo        []AkunDemo        `yaml:"akun_demo"`
	LaporanSimulasi []LaporanSimulasi `yaml:"laporan_simulasi"`
	TanyaDemo       []string          `yaml:"tanya_demo"`
}

type LaporanSimulasi struct {
	Kondisi KondisiLaporan `yaml:"kondisi"`
	Catatan string         `yaml:"catatan"`
	Pelapor string         `yaml:"pelapor"`
	Saksi   []string       `yaml:"saksi"`
	Foto    []string       `yaml:"foto"`
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
	Peran        Peran  `yaml:"peran"`
	Email        string `yaml:"email"`
	NamaTampilan string `yaml:"nama_tampilan"`
	Instansi     string `yaml:"instansi"`
	AnggotaID    string `yaml:"anggota_id"`
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
		if err := akun.Validasi(); err != nil {
			return err
		}
	}
	return s.validasiLaporanSimulasi()
}

func (s SkenarioDemo) validasiLaporanSimulasi() error {
	emailDemo := map[string]bool{}
	for _, akun := range s.AkunDemo {
		emailDemo[akun.Email] = true
	}
	for _, l := range s.LaporanSimulasi {
		if !l.Kondisi.IsValid() || len(l.Foto) == 0 || len(l.Foto) > BatasFotoLaporan || len(l.Saksi)+1 < JumlahKonfirmasiTerverifikasi {
			return ErrLaporanSimulasiTidakSah
		}
		if !isAkunDemoBerbeda(append([]string{l.Pelapor}, l.Saksi...), emailDemo) {
			return ErrPelaporBukanAkunDemo
		}
	}
	return nil
}

func isAkunDemoBerbeda(email []string, emailDemo map[string]bool) bool {
	dilihat := map[string]bool{}
	for _, e := range email {
		if !emailDemo[e] || dilihat[e] {
			return false
		}
		dilihat[e] = true
	}
	return true
}

func (a AkunDemo) Validasi() error {
	switch {
	case !a.Peran.IsDikenal():
		return ErrPeranTidakDikenal
	case a.AnggotaID != "":
		return ErrAkunDemoTertautAnggota
	case !strings.HasSuffix(a.Email, AkhiranEmailUji) || !strings.Contains(a.NamaTampilan, PenandaNamaSimulasi):
		return ErrAkunDemoBukanSimulasi
	case (a.Peran == PeranInstansi) != (a.Instansi != ""):
		return ErrAkunInstansiTanpaInstansi
	}
	return nil
}
