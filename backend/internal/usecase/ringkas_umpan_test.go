package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
)

type penulisRingkasanTiruan struct {
	jawaban map[int64]string
	err     error
	urutan  int64
}

func (p *penulisRingkasanTiruan) Ringkas(context.Context, string, string) (string, error) {
	p.urutan++
	return p.jawaban[p.urutan], p.err
}

type penyimpanRingkasanTiruan struct {
	umpan     []domain.UmpanUntukRingkas
	tersimpan map[int64]string
}

func (p *penyimpanRingkasanTiruan) UmpanBelumDiringkas(context.Context, int) ([]domain.UmpanUntukRingkas, error) {
	return p.umpan, nil
}

func (p *penyimpanRingkasanTiruan) SimpanRingkasan(_ context.Context, id int64, ringkasan string) error {
	p.tersimpan[id] = ringkasan
	return nil
}

func umpanRingkasUji() *penyimpanRingkasanTiruan {
	return &penyimpanRingkasanTiruan{
		umpan: []domain.UmpanUntukRingkas{
			{ID: 1, Judul: "Kasus DBD tembus 1.314", Cuplikan: "Komisi D meminta pengawasan diperketat."},
			{ID: 2, Judul: "Rapat anggaran sekolah", Cuplikan: "Komisi E membahas anggaran."},
			{ID: 3, Judul: "Tinjau banjir", Cuplikan: "Anggota meninjau tanggul."},
		},
		tersimpan: map[int64]string{},
	}
}

func ringkasUji(penulis PenulisRingkasan, penyimpan *penyimpanRingkasanTiruan, pencatat *pencatatTiruan) RingkasUmpan {
	return RingkasUmpan{
		Umpan: penyimpan, Penulis: penulis, Pencatat: pencatat, Batas: 10,
		Aturan: domain.AturanRingkasan{Terlarang: domain.KosakataTerlarang{"korupsi"}, KataPenilaian: []string{"hebat"}},
	}
}

func TestRingkasUmpan_LLMTiruan_RingkasanTerlarangDanAngkaKaranganJatuhKeJudul(t *testing.T) {
	penulis := &penulisRingkasanTiruan{jawaban: map[int64]string{
		1: "Komisi D meminta pengawasan DBD diperketat karena kasus mencapai 1.314.",
		2: "Komisi E menduga korupsi anggaran sekolah.",
		3: "Anggota meninjau tanggul yang menampung 9000 warga.",
	}}
	penyimpan := umpanRingkasUji()
	pencatat := &pencatatTiruan{}

	if err := ringkasUji(penulis, penyimpan, pencatat).Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}

	if penyimpan.tersimpan[1] != penulis.jawaban[1] {
		t.Errorf("ringkasan sah tidak tersimpan: %q", penyimpan.tersimpan[1])
	}
	for _, id := range []int64{2, 3} {
		if ringkasan, ada := penyimpan.tersimpan[id]; !ada || ringkasan != RingkasanDitolak {
			t.Errorf("umpan %d harap ditandai ditolak, dapat %q (ada %v)", id, ringkasan, ada)
		}
	}
	if pencatat.selesai == nil || pencatat.selesai.JumlahBerubah != 3 {
		t.Errorf("run harap selesai dengan 3 umpan diproses, dapat %+v", pencatat.selesai)
	}
}

func TestRingkasUmpan_ModelTidakTersedia_RunGagalTanpaMenandaiUmpan(t *testing.T) {
	penulis := &penulisRingkasanTiruan{err: errors.New("koneksi ditolak")}
	penyimpan := umpanRingkasUji()
	pencatat := &pencatatTiruan{}

	err := ringkasUji(penulis, penyimpan, pencatat).Jalankan(t.Context())

	if err == nil || pencatat.gagal == nil || len(penyimpan.tersimpan) != 0 {
		t.Errorf("harap run gagal dan umpan tetap tanpa ringkasan, dapat err %v, tersimpan %v", err, penyimpan.tersimpan)
	}
}
