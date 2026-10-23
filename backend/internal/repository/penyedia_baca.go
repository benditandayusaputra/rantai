package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlAmbilPenyedia = `SELECT id, nama, alias FROM penyedia WHERE id = $1`

	sqlPenyediaPerTahun = `SELECT tahun, sumber, jumlah, total, provenance_id FROM (
    SELECT p.tahun_anggaran AS tahun, 'lpse' AS sumber, count(*) AS jumlah, coalesce(sum(coalesce(p.nilai_kontrak, p.hps)), 0) AS total,
        (array_agg(p.provenance_id ORDER BY pr.diambil_pada DESC))[1] AS provenance_id
    FROM paket p JOIN provenance pr ON pr.id = p.provenance_id
    WHERE p.penyedia_id = $1 AND NOT p.is_simulation
    GROUP BY 1
    UNION ALL
    SELECT r.tahun, 'opentender', count(*), coalesce(sum(r.nilai), 0), (array_agg(r.provenance_id ORDER BY pr.diambil_pada DESC))[1]
    FROM riwayat_penyedia r JOIN provenance pr ON pr.id = r.provenance_id
    WHERE r.penyedia_id = $1
    GROUP BY 1
) AS gabungan
ORDER BY tahun, sumber`

	sqlKategoriPenyedia = `SELECT kategori, count(*) FROM paket WHERE penyedia_id = $1 AND kategori IS NOT NULL
GROUP BY 1 ORDER BY 2 DESC, 1 LIMIT $2`

	sqlInstansiPenyedia = `SELECT nama, count(*) FROM (
    SELECT i.nama FROM paket p JOIN instansi i ON i.id = p.instansi_id WHERE p.penyedia_id = $1
    UNION ALL
    SELECT pembeli FROM riwayat_penyedia WHERE penyedia_id = $1
) AS pemberi
GROUP BY 1 ORDER BY 2 DESC, 1 LIMIT $2`

	sqlPaketPenyedia = `SELECT p.kode, p.nama, p.tahap, p.jenis_pengadaan, p.metode, p.tahun_anggaran, p.hps, p.nilai_kontrak,
    coalesce(p.sektor, ''), p.tema_kode, coalesce(p.lokasi_tingkat::text, ''), coalesce(p.kode_wilayah, ''), coalesce(w.nama, ''),
    p.is_simulation, p.provenance_id, '', 0
FROM paket p LEFT JOIN wilayah w ON w.kode = p.kode_wilayah
WHERE p.penyedia_id = $1 AND NOT p.is_simulation
ORDER BY p.tahun_anggaran DESC, p.kode
LIMIT $2`

	sqlRiwayatPenyedia = `SELECT tahun, pembeli, nama_paket, nilai, kode_sumber, sumber, provenance_id
FROM riwayat_penyedia WHERE penyedia_id = $1
ORDER BY tahun DESC, nilai DESC NULLS LAST, kode_sumber
LIMIT $2`
)

type PenyediaRepository struct {
	DB Basisdata
}

func (r PenyediaRepository) AmbilProfilPenyedia(ctx context.Context, id int64, batas int) (domain.ProfilPenyedia, bool, error) {
	var p domain.ProfilPenyedia
	err := r.DB.QueryRow(ctx, sqlAmbilPenyedia, id).Scan(&p.ID, &p.Nama, &p.Alias)
	p, ada, err := hasilSatuBaris(p, err)
	if err != nil || !ada {
		return p, ada, err
	}
	if err := r.lengkapiProfil(ctx, &p, batas); err != nil {
		return domain.ProfilPenyedia{}, false, fmt.Errorf("profil penyedia %d: %w", id, err)
	}
	return p, true, nil
}

func (r PenyediaRepository) lengkapiProfil(ctx context.Context, p *domain.ProfilPenyedia, batas int) error {
	var err error
	if p.PerTahun, err = kumpulkan(ctx, r.DB, sqlPenyediaPerTahun, func(b pgx.CollectableRow) (domain.RingkasanTahunPenyedia, error) {
		var t domain.RingkasanTahunPenyedia
		return t, b.Scan(&t.Tahun, &t.Sumber, &t.JumlahPaket, &t.TotalNilai, &t.ProvenanceID)
	}, p.ID); err != nil {
		return err
	}
	if p.KategoriUtama, err = kumpulkan(ctx, r.DB, sqlKategoriPenyedia, pindaiJumlahPerNama, p.ID, batas); err != nil {
		return err
	}
	if p.InstansiPemberi, err = kumpulkan(ctx, r.DB, sqlInstansiPenyedia, pindaiJumlahPerNama, p.ID, batas); err != nil {
		return err
	}
	if p.Paket, err = kumpulkan(ctx, r.DB, sqlPaketPenyedia, pindaiPaketRingkas, p.ID, batas); err != nil {
		return err
	}
	p.Riwayat, err = kumpulkan(ctx, r.DB, sqlRiwayatPenyedia, func(b pgx.CollectableRow) (domain.RiwayatTersimpan, error) {
		var t domain.RiwayatTersimpan
		return t, b.Scan(&t.Tahun, &t.Pembeli, &t.NamaPaket, &t.Nilai, &t.KodeSumber, &t.Sumber, &t.ProvenanceID)
	}, p.ID, batas)
	return err
}

func pindaiJumlahPerNama(b pgx.CollectableRow) (domain.JumlahPerNama, error) {
	var j domain.JumlahPerNama
	return j, b.Scan(&j.Nama, &j.Jumlah)
}

func pindaiPaketRingkas(b pgx.CollectableRow) (domain.PaketRingkas, error) {
	var p domain.PaketRingkas
	var total int
	return p, b.Scan(&p.Kode, &p.Nama, &p.Tahap, &p.JenisPengadaan, &p.Metode, &p.TahunAnggaran, &p.HPS, &p.NilaiKontrak,
		&p.Sektor, &p.TemaKode, &p.Lokasi.Tingkat, &p.Lokasi.KodeWilayah, &p.Lokasi.NamaWilayah, &p.IsSimulation,
		&p.ProvenanceID, &p.NamaPenyedia, &total)
}

func kumpulkan[T any](ctx context.Context, db Basisdata, sql string, pindai pgx.RowToFunc[T], argumen ...any) ([]T, error) {
	baris, err := db.Query(ctx, sql, argumen...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, pindai)
}
