package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlPilihKecamatanDemo = `SELECT i.kode_kecamatan FROM indikator_kecamatan i JOIN wilayah w ON w.kode = i.kode_kecamatan
WHERE i.tema_kode = $1 AND w.induk_kode = $2 AND i.rasio IS NOT NULL AND i.label = 'resmi'
ORDER BY i.periode DESC, i.rasio, i.kode_kecamatan
LIMIT 1`

	sqlSimpanPenyediaSimulasi = `INSERT INTO penyedia (nama, nama_normal, alias) VALUES ($1, $2, ARRAY[$1])
ON CONFLICT (nama_normal) DO UPDATE SET diubah_pada = now()
RETURNING id`

	sqlSimulasiTertautDataNyata = `SELECT EXISTS (SELECT 1 FROM paket WHERE NOT is_simulation AND (penyedia_id = $1 OR instansi_id = $2))
    OR EXISTS (SELECT 1 FROM riwayat_penyedia WHERE penyedia_id = $1)
    OR EXISTS (SELECT 1 FROM rup WHERE instansi_id = $2)`

	sqlSimpanPaketSimulasi = `INSERT INTO paket (kode, kode_spse, nama, instansi_id, penyedia_id, is_tender, jenis_pengadaan, metode,
    tahun_anggaran, pagu, hps, nilai_kontrak, tanggal_kontrak, tahap, jadwal_mulai, jadwal_selesai, kategori, sektor,
    tema_kode, lokasi_tingkat, kode_wilayah, lokasi, is_simulation, provenance_id)
SELECT $1, $1, $2, $3, $4, false, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
    'kecamatan', w.kode, w.titik_tengah, true, $18
FROM wilayah w WHERE w.kode = $19 AND w.tingkat = 'kecamatan'
ON CONFLICT (kode) DO UPDATE SET
    nama = EXCLUDED.nama, instansi_id = EXCLUDED.instansi_id, penyedia_id = EXCLUDED.penyedia_id,
    jenis_pengadaan = EXCLUDED.jenis_pengadaan, metode = EXCLUDED.metode, tahun_anggaran = EXCLUDED.tahun_anggaran,
    pagu = EXCLUDED.pagu, hps = EXCLUDED.hps, nilai_kontrak = EXCLUDED.nilai_kontrak,
    tanggal_kontrak = EXCLUDED.tanggal_kontrak, tahap = EXCLUDED.tahap, jadwal_mulai = EXCLUDED.jadwal_mulai,
    jadwal_selesai = EXCLUDED.jadwal_selesai, kategori = EXCLUDED.kategori, sektor = EXCLUDED.sektor,
    tema_kode = EXCLUDED.tema_kode, kode_wilayah = EXCLUDED.kode_wilayah, lokasi = EXCLUDED.lokasi,
    provenance_id = EXCLUDED.provenance_id, diubah_pada = now()
WHERE paket.is_simulation
RETURNING kode`
)

type SkenarioDemoRepository struct {
	DB Basisdata
}

func (r SkenarioDemoRepository) PilihKecamatanDemo(ctx context.Context, kodeKota string) (string, error) {
	var kode string
	err := r.DB.QueryRow(ctx, sqlPilihKecamatanDemo, domain.TemaPendidikan, kodeKota).Scan(&kode)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrKecamatanDemoBelumDihitung
	}
	if err != nil {
		return "", fmt.Errorf("pilih kecamatan demo: %w", err)
	}
	return kode, nil
}

func (r SkenarioDemoRepository) SimpanPaketSimulasi(ctx context.Context, paket usecase.PaketSimulasiTersimpan) error {
	return dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		if err := simpanProvenance(ctx, tx, paket.Provenance); err != nil {
			return err
		}
		idInstansi, idPenyedia, err := simpanPihakSimulasi(ctx, tx, paket.PaketSimulasi)
		if err != nil {
			return err
		}
		return simpanPaketSimulasi(ctx, tx, paket, pihakPaket{instansi: idInstansi, penyedia: idPenyedia})
	})
}

type pihakPaket struct {
	instansi int64
	penyedia int64
}

func simpanPihakSimulasi(ctx context.Context, tx pgx.Tx, p domain.PaketSimulasi) (int64, int64, error) {
	idInstansi, err := simpanInstansi(ctx, tx, domain.Instansi{Nama: p.Instansi, Jenis: domain.JenisInstansiLainnya})
	if err != nil {
		return 0, 0, err
	}
	var idPenyedia int64
	if err := tx.QueryRow(ctx, sqlSimpanPenyediaSimulasi, p.Penyedia, domain.NormalisasiNamaPenyedia(p.Penyedia)).Scan(&idPenyedia); err != nil {
		return 0, 0, fmt.Errorf("simpan penyedia simulasi %s: %w", p.Penyedia, err)
	}
	var tertaut bool
	if err := tx.QueryRow(ctx, sqlSimulasiTertautDataNyata, idPenyedia, idInstansi).Scan(&tertaut); err != nil {
		return 0, 0, fmt.Errorf("periksa keterpisahan data simulasi: %w", err)
	}
	if tertaut {
		return 0, 0, fmt.Errorf("%s, %s: %w", p.Penyedia, p.Instansi, domain.ErrSimulasiTertautDataNyata)
	}
	return idInstansi, idPenyedia, nil
}

func simpanPaketSimulasi(ctx context.Context, tx pgx.Tx, paket usecase.PaketSimulasiTersimpan, pihak pihakPaket) error {
	p := paket.PaketSimulasi
	var kode string
	err := tx.QueryRow(ctx, sqlSimpanPaketSimulasi, p.Kode, p.Nama, pihak.instansi, pihak.penyedia, p.JenisPengadaan, p.Metode,
		p.TahunAnggaran, p.Pagu, p.HPS, p.NilaiKontrak, p.TanggalKontrak, p.Tahap, p.JadwalMulai, p.JadwalSelesai,
		p.Kategori, p.Sektor, p.TemaKode, paket.Provenance.ID, paket.KodeKecamatan).Scan(&kode)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.TidakDitemukan("kecamatan demo " + paket.KodeKecamatan + " tidak ada di tabel wilayah")
	}
	if err != nil {
		return fmt.Errorf("simpan paket simulasi %s: %w", p.Kode, err)
	}
	return nil
}
