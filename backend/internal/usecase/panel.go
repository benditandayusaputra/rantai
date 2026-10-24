package usecase

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

const (
	formatTanggal       = "2006-01-02"
	catatanSerahTerima  = "Tanggal serah terima %s diisi petugas %s lewat Panel Instansi"
	catatanBeritaAcara  = "; berita acara nomor %s"
	rentangTahunAntrean = 1
)

type HakJawabTersimpan struct {
	domain.HakJawabBaru
	Penjawab domain.Pengguna
	Pemilik  domain.PemilikTarget
}

type SerahTerimaTersimpan struct {
	domain.SerahTerima
	Provenance domain.Provenance
}

type FilterPaketPanel struct {
	InstansiID     int64
	JenisPengadaan []string
	TahunMinimal   int
}

type PenyimpanPanel interface {
	PemilikTarget(ctx context.Context, jenis domain.JenisTargetHakJawab, id string) (domain.PemilikTarget, error)
	SimpanHakJawab(ctx context.Context, h HakJawabTersimpan) (domain.HakJawab, error)
	SimpanSerahTerima(ctx context.Context, s SerahTerimaTersimpan) error
	AntreanInstansi(ctx context.Context, f FilterPaketPanel) (domain.AntreanInstansi, error)
	AntreanAnggota(ctx context.Context, anggotaID string) ([]domain.UmpanPanel, error)
}

type PaketPemeliharaanPanel struct {
	domain.PaketJadwalPanel
	Estimasi         domain.EstimasiPemeliharaan
	IsSegeraBerakhir bool
}

type AntreanPanel struct {
	Pengguna     domain.Pengguna
	NamaInstansi string
	Sinyal       []domain.SinyalPanel
	Laporan      []domain.LaporanPanel
	Paket        []PaketPemeliharaanPanel
	Umpan        []domain.UmpanPanel
}

type Panel struct {
	Pengguna       PenyimpanPengguna
	Penyimpan      PenyimpanPanel
	Laporan        pengubahStatus
	Pemeliharaan   KonfigurasiPemeliharaan
	BuatProvenance PembuatProvenance
	Sekarang       func() time.Time
}

func (u Panel) Antrean(ctx context.Context, penggunaID int64) (AntreanPanel, error) {
	pengguna, err := u.penggunaPanel(ctx, penggunaID)
	if err != nil {
		return AntreanPanel{}, err
	}
	hasil := AntreanPanel{Pengguna: pengguna, Sinyal: []domain.SinyalPanel{}, Laporan: []domain.LaporanPanel{}, Paket: []PaketPemeliharaanPanel{}, Umpan: []domain.UmpanPanel{}}
	switch {
	case pengguna.Peran == domain.PeranInstansi && pengguna.InstansiID != nil:
		return u.antreanInstansi(ctx, hasil, *pengguna.InstansiID)
	case pengguna.Peran == domain.PeranAnggota && pengguna.AnggotaID != nil:
		hasil.Umpan, err = u.Penyimpan.AntreanAnggota(ctx, *pengguna.AnggotaID)
		return hasil, err
	}
	return hasil, nil
}

func (u Panel) antreanInstansi(ctx context.Context, hasil AntreanPanel, instansiID int64) (AntreanPanel, error) {
	antrean, err := u.Penyimpan.AntreanInstansi(ctx, FilterPaketPanel{
		InstansiID: instansiID, JenisPengadaan: u.Pemeliharaan.JenisPengadaan, TahunMinimal: u.Sekarang().Year() - rentangTahunAntrean,
	})
	if err != nil {
		return hasil, fmt.Errorf("baca antrean instansi %d: %w", instansiID, err)
	}
	hasil.NamaInstansi, hasil.Sinyal, hasil.Laporan = antrean.NamaInstansi, antrean.Sinyal, antrean.Laporan
	hasil.Paket = u.paketPerluPerhatian(antrean.Paket)
	return hasil, nil
}

func (u Panel) paketPerluPerhatian(daftar []domain.PaketJadwalPanel) []PaketPemeliharaanPanel {
	hasil := []PaketPemeliharaanPanel{}
	for _, p := range daftar {
		e := u.Pemeliharaan.Estimasi(p.Jadwal, u.Sekarang())
		if e.Status != domain.StatusPelaksanaan && e.Status != domain.StatusPemeliharaanBerjalan {
			continue
		}
		segera := e.Status == domain.StatusPemeliharaanBerjalan && e.SisaHari <= u.Pemeliharaan.AmbangPanelHari
		hasil = append(hasil, PaketPemeliharaanPanel{PaketJadwalPanel: p, Estimasi: e, IsSegeraBerakhir: segera})
	}
	slices.SortStableFunc(hasil, func(a, b PaketPemeliharaanPanel) int {
		return a.Estimasi.AkhirPemeliharaan.Compare(*b.Estimasi.AkhirPemeliharaan)
	})
	return hasil[:min(len(hasil), domain.BatasAntreanPanel)]
}

func (u Panel) TulisHakJawab(ctx context.Context, penggunaID int64, baru domain.HakJawabBaru) (domain.HakJawab, error) {
	baru = baru.Rapikan()
	if err := baru.Validasi(); err != nil {
		return domain.HakJawab{}, err
	}
	pengguna, pemilik, err := u.penjawabTarget(ctx, penggunaID, baru.TargetJenis, baru.TargetID)
	if err != nil {
		return domain.HakJawab{}, err
	}
	if err := pengguna.PeriksaBolehMenjawab(pemilik); err != nil {
		return domain.HakJawab{}, err
	}
	hakJawab, err := u.Penyimpan.SimpanHakJawab(ctx, HakJawabTersimpan{HakJawabBaru: baru, Penjawab: pengguna, Pemilik: pemilik})
	if err != nil {
		return domain.HakJawab{}, fmt.Errorf("simpan hak jawab: %w", err)
	}
	return hakJawab, u.tandaiDitanggapi(ctx, baru, pemilik, pengguna.ID)
}

func (u Panel) tandaiDitanggapi(ctx context.Context, baru domain.HakJawabBaru, pemilik domain.PemilikTarget, oleh int64) error {
	if baru.TargetJenis != domain.TargetHakJawabLaporan || domain.Transisi(pemilik.StatusLaporan, domain.StatusDitanggapi) != nil {
		return nil
	}
	id, _ := strconv.ParseInt(baru.TargetID, 10, 64)
	err := u.Laporan.UbahStatus(ctx, domain.PerubahanStatus{
		LaporanID: id, Dari: pemilik.StatusLaporan, Ke: domain.StatusDitanggapi, OlehPenggunaID: &oleh, Alasan: domain.AlasanDitanggapi,
	})
	if errors.Is(err, domain.ErrTransisiTidakSah) {
		return nil
	}
	return err
}

func (u Panel) TandaiLaporanSelesai(ctx context.Context, penggunaID, laporanID int64) error {
	pengguna, pemilik, err := u.penjawabTarget(ctx, penggunaID, domain.TargetHakJawabLaporan, strconv.FormatInt(laporanID, 10))
	if err != nil {
		return err
	}
	if err := pengguna.PeriksaPetugasPemilik(pemilik); err != nil {
		return err
	}
	if err := domain.Transisi(pemilik.StatusLaporan, domain.StatusSelesai); err != nil {
		return err
	}
	return u.Laporan.UbahStatus(ctx, domain.PerubahanStatus{
		LaporanID: laporanID, Dari: pemilik.StatusLaporan, Ke: domain.StatusSelesai, OlehPenggunaID: &pengguna.ID, Alasan: domain.AlasanSelesaiInstansi,
	})
}

func (u Panel) IsiSerahTerima(ctx context.Context, penggunaID int64, s domain.SerahTerima) error {
	if err := s.Validasi(u.Sekarang()); err != nil {
		return err
	}
	pengguna, pemilik, err := u.penjawabTarget(ctx, penggunaID, domain.TargetHakJawabPaket, s.KodePaket)
	if err != nil {
		return err
	}
	if err := pengguna.PeriksaPetugasPemilik(pemilik); err != nil {
		return err
	}
	return u.Penyimpan.SimpanSerahTerima(ctx, SerahTerimaTersimpan{SerahTerima: s, Provenance: u.provenanceSerahTerima(s, pemilik)})
}

func (u Panel) provenanceSerahTerima(s domain.SerahTerima, pemilik domain.PemilikTarget) domain.Provenance {
	catatan := fmt.Sprintf(catatanSerahTerima, s.Tanggal.Format(formatTanggal), pemilik.NamaInstansi)
	if s.NomorBeritaAcara != "" {
		catatan += fmt.Sprintf(catatanBeritaAcara, s.NomorBeritaAcara)
	}
	label, sumber := domain.LabelResmi, domain.SumberInstansi
	if pemilik.IsSimulation {
		label, sumber = domain.LabelSimulasi, domain.SumberSimulasi
	}
	return u.BuatProvenance(CatatanProvenance{Sumber: sumber, Label: label, Catatan: catatan})
}

func (u Panel) penjawabTarget(ctx context.Context, penggunaID int64, jenis domain.JenisTargetHakJawab, id string) (domain.Pengguna, domain.PemilikTarget, error) {
	pengguna, err := u.penggunaPanel(ctx, penggunaID)
	if err != nil {
		return domain.Pengguna{}, domain.PemilikTarget{}, err
	}
	pemilik, err := u.Penyimpan.PemilikTarget(ctx, jenis, id)
	return pengguna, pemilik, err
}

func (u Panel) penggunaPanel(ctx context.Context, penggunaID int64) (domain.Pengguna, error) {
	pengguna, err := u.Pengguna.AmbilPengguna(ctx, penggunaID)
	if err != nil {
		return domain.Pengguna{}, err
	}
	return pengguna, pengguna.PeriksaAksesPanel()
}
