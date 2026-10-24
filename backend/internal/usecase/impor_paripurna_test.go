package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type pengambilParipurnaTiruan struct{ data DataParipurna }

func (p pengambilParipurnaTiruan) Ambil(context.Context) (DataParipurna, error) { return p.data, nil }

type penyimpanParipurnaTiruan struct{ disimpan *DataParipurna }

func (p *penyimpanParipurnaTiruan) DaftarFraksi(context.Context) ([]string, error) {
	return []string{"Fraksi Partai Golkar"}, nil
}

func (p *penyimpanParipurnaTiruan) SimpanParipurna(_ context.Context, data DataParipurna) (HasilImpor, error) {
	p.disimpan = &data
	return HasilImpor{JumlahBaru: len(data.Sikap) + len(data.Aspirasi)}, nil
}

func sikapImporUji(fraksi, ringkasan string) domain.SikapFraksi {
	return domain.SikapFraksi{
		Fraksi: fraksi, Agenda: domain.AgendaAPBD, JudulAgenda: "Ranperda APBD 2027", Ringkasan: ringkasan,
		Tanggal: time.Date(2026, time.October, 19, 0, 0, 0, 0, time.UTC), URL: "https://dprd-dkijakartaprov.go.id/berita/",
	}
}

func jalankanImporParipurna(sikap ...domain.SikapFraksi) (*penyimpanParipurnaTiruan, *pencatatTiruan, error) {
	penyimpan, pencatat := &penyimpanParipurnaTiruan{}, &pencatatTiruan{}
	impor := ImporParipurna{
		Pengambil: pengambilParipurnaTiruan{data: DataParipurna{Sikap: sikap, Aspirasi: []domain.AspirasiReses{}}},
		Penyimpan: penyimpan, Pencatat: pencatat, Terlarang: domain.KosakataTerlarang{"korupsi"},
	}
	return penyimpan, pencatat, impor.Jalankan(context.Background())
}

func TestImporParipurna_EntriSah_DisimpanDanSelesai(t *testing.T) {
	penyimpan, pencatat, err := jalankanImporParipurna(sikapImporUji("Fraksi Partai Golkar", "Meminta subsidi berbasis kinerja."))

	if err != nil || penyimpan.disimpan == nil || pencatat.selesai == nil || pencatat.selesai.JumlahBaru != 1 {
		t.Errorf("err %v, disimpan %v, selesai %v", err, penyimpan.disimpan, pencatat.selesai)
	}
}

func TestImporParipurna_FraksiTidakDikenal_DitolakTanpaMenyimpan(t *testing.T) {
	penyimpan, pencatat, err := jalankanImporParipurna(sikapImporUji("Fraksi Contoh", "Meminta subsidi berbasis kinerja."))

	if !errors.Is(err, ErrFraksiTidakDikenal) || penyimpan.disimpan != nil || pencatat.gagal == nil {
		t.Errorf("err %v, disimpan %v", err, penyimpan.disimpan)
	}
}

func TestImporParipurna_RingkasanBerkosakataTerlarang_Ditolak(t *testing.T) {
	penyimpan, _, err := jalankanImporParipurna(sikapImporUji("Fraksi Partai Golkar", "Menyebut dugaan korupsi."))

	if !errors.Is(err, ErrKosakataTerlarangParipurna) || penyimpan.disimpan != nil {
		t.Errorf("err %v, disimpan %v", err, penyimpan.disimpan)
	}
}
