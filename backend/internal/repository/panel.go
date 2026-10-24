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
	sqlKecamatanPaket = `CASE WHEN p.lokasi_tingkat IN ('kelurahan', 'kecamatan')
    THEN substring(p.kode_wilayah from '^[0-9]+\.[0-9]+\.[0-9]+') END`

	sqlPemilikPaket = `SELECT p.instansi_id, i.nama, NULL::text[], ` + sqlKecamatanPaket + `, p.is_simulation, ''
FROM paket p JOIN instansi i ON i.id = p.instansi_id WHERE p.kode = $1`

	sqlPemilikSinyal = `SELECT p.instansi_id, i.nama, NULL::text[], ` + sqlKecamatanPaket + `, p.is_simulation, ''
FROM sinyal s JOIN paket p ON p.kode = s.paket_kode JOIN instansi i ON i.id = p.instansi_id
WHERE s.id = $1::bigint`

	sqlPemilikLaporan = `SELECT p.instansi_id, coalesce(i.nama, ''), NULL::text[],
    substring(l.kode_wilayah from '^[0-9]+\.[0-9]+\.[0-9]+'), l.is_simulation, l.status::text
FROM laporan l LEFT JOIN paket p ON p.kode = l.paket_kode LEFT JOIN instansi i ON i.id = p.instansi_id
WHERE l.id = $1::bigint`

	sqlPemilikUmpan = `SELECT NULL::bigint, '', array_agg(ua.anggota_id), NULL::text, false, ''
FROM umpan_kegiatan u JOIN umpan_anggota ua ON ua.umpan_id = u.id
WHERE u.id = $1::bigint AND NOT u.disembunyikan
GROUP BY u.id`

	sqlSimpanHakJawab = `INSERT INTO hak_jawab (target_jenis, target_id, instansi_id, anggota_id, isi, tautan_dokumen, penjawab_id,
    kode_kecamatan, is_simulation)
VALUES ($1, $2, $3, $4, $5, NULLIF($6, ''), $7, NULLIF($8, ''), $9)
RETURNING id`

	sqlHakJawabByID = `SELECT h.id, h.target_jenis, h.target_id, h.isi, coalesce(h.tautan_dokumen, ''), g.nama_tampilan,
    coalesce(i.nama, a.nama, ''), h.is_simulation, h.dibuat_pada
FROM hak_jawab h
JOIN pengguna g ON g.id = h.penjawab_id
LEFT JOIN instansi i ON i.id = h.instansi_id
LEFT JOIN anggota a ON a.id = h.anggota_id
WHERE h.id = $1`

	sqlSimpanSerahTerima = `UPDATE paket SET tanggal_serah_terima = $2, nomor_berita_acara = NULLIF($3, ''),
    serah_terima_provenance_id = $4, diubah_pada = now()
WHERE kode = $1`

	sqlNamaInstansi = `SELECT nama FROM instansi WHERE id = $1`

	sqlSinyalInstansi = `SELECT s.id, s.jenis, s.alasan, p.kode, p.nama,
    (SELECT count(*) FROM hak_jawab h WHERE h.target_jenis = 'sinyal' AND h.target_id = s.id::text)
FROM sinyal s JOIN paket p ON p.kode = s.paket_kode
WHERE p.instansi_id = $1 AND s.aktif
ORDER BY s.dihitung_pada DESC, s.id
LIMIT $2`

	sqlLaporanInstansi = `SELECT l.id, l.status, l.kondisi, l.catatan, p.kode, p.nama, w.nama, l.is_simulation, l.dibuat_pada
FROM laporan l JOIN paket p ON p.kode = l.paket_kode JOIN wilayah w ON w.kode = l.kode_wilayah
WHERE p.instansi_id = $1 AND l.status IN ('DITERUSKAN', 'DITANGGAPI')
ORDER BY l.status, l.diubah_pada
LIMIT $2`

	sqlPaketInstansi = `SELECT p.kode, p.nama, p.jenis_pengadaan, coalesce(p.nomor_berita_acara, ''), p.is_simulation,
    coalesce(p.kategori, ''), p.tanggal_serah_terima, r.jadwal_pelaksanaan_selesai, p.tanggal_kontrak
FROM paket p LEFT JOIN rup r ON r.kode = p.kode_rup
WHERE p.instansi_id = $1 AND NOT p.is_batal AND p.jenis_pengadaan = ANY ($2) AND p.tahun_anggaran >= $3`

	sqlUmpanPanelAnggota = `SELECT u.id, u.judul, u.media_nama, u.url, u.terbit_pada,
    (SELECT count(*) FROM hak_jawab h WHERE h.target_jenis = 'umpan' AND h.target_id = u.id::text)
FROM umpan_anggota ua JOIN umpan_kegiatan u ON u.id = ua.umpan_id
WHERE ua.anggota_id = $1 AND NOT u.disembunyikan
ORDER BY u.terbit_pada DESC, u.id DESC
LIMIT $2`
)

var sqlPemilikPerJenis = map[domain.JenisTargetHakJawab]string{
	domain.TargetHakJawabPaket: sqlPemilikPaket, domain.TargetHakJawabSinyal: sqlPemilikSinyal,
	domain.TargetHakJawabLaporan: sqlPemilikLaporan, domain.TargetHakJawabUmpan: sqlPemilikUmpan,
}

type PanelRepository struct {
	DB Basisdata
}

func (r PanelRepository) PemilikTarget(ctx context.Context, jenis domain.JenisTargetHakJawab, id string) (domain.PemilikTarget, error) {
	var t domain.PemilikTarget
	var kecamatan *string
	err := r.DB.QueryRow(ctx, sqlPemilikPerJenis[jenis], id).Scan(&t.InstansiID, &t.NamaInstansi, &t.AnggotaID, &kecamatan, &t.IsSimulation, &t.StatusLaporan)
	if errors.Is(err, pgx.ErrNoRows) {
		return t, domain.ErrSorotanTidakDitemukan
	}
	if err != nil {
		return t, fmt.Errorf("cari pemilik %s %s: %w", jenis, id, err)
	}
	if kecamatan != nil {
		t.KodeKecamatan = *kecamatan
	}
	return t, nil
}

func (r PanelRepository) SimpanHakJawab(ctx context.Context, h usecase.HakJawabTersimpan) (domain.HakJawab, error) {
	var instansiID *int64
	var anggotaID *string
	if h.Penjawab.Peran == domain.PeranInstansi {
		instansiID = h.Penjawab.InstansiID
	} else {
		anggotaID = h.Penjawab.AnggotaID
	}
	var id int64
	err := r.DB.QueryRow(ctx, sqlSimpanHakJawab, h.TargetJenis, h.TargetID, instansiID, anggotaID, h.Isi, h.TautanDokumen,
		h.Penjawab.ID, h.Pemilik.KodeKecamatan, h.Pemilik.IsSimulation).Scan(&id)
	if err != nil {
		return domain.HakJawab{}, err
	}
	baris, err := r.DB.Query(ctx, sqlHakJawabByID, id)
	if err != nil {
		return domain.HakJawab{}, err
	}
	return pgx.CollectExactlyOneRow(baris, pgx.RowToStructByPos[domain.HakJawab])
}

func (r PanelRepository) SimpanSerahTerima(ctx context.Context, s usecase.SerahTerimaTersimpan) error {
	return dalamTransaksi(ctx, r.DB, func(tx pgx.Tx) error {
		if err := simpanProvenance(ctx, tx, s.Provenance); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, sqlSimpanSerahTerima, s.KodePaket, s.Tanggal, s.NomorBeritaAcara, s.Provenance.ID)
		if err != nil {
			return fmt.Errorf("simpan serah terima %s: %w", s.KodePaket, err)
		}
		return nil
	})
}

func (r PanelRepository) AntreanInstansi(ctx context.Context, f usecase.FilterPaketPanel) (domain.AntreanInstansi, error) {
	var a domain.AntreanInstansi
	if err := r.DB.QueryRow(ctx, sqlNamaInstansi, f.InstansiID).Scan(&a.NamaInstansi); err != nil {
		return a, fmt.Errorf("baca instansi %d: %w", f.InstansiID, err)
	}
	var err error
	if a.Sinyal, err = kumpulkan(ctx, r.DB, sqlSinyalInstansi, pgx.RowToStructByPos[domain.SinyalPanel], f.InstansiID, domain.BatasAntreanPanel); err != nil {
		return a, fmt.Errorf("baca sinyal instansi: %w", err)
	}
	if a.Laporan, err = kumpulkan(ctx, r.DB, sqlLaporanInstansi, pgx.RowToStructByPos[domain.LaporanPanel], f.InstansiID, domain.BatasAntreanPanel); err != nil {
		return a, fmt.Errorf("baca laporan instansi: %w", err)
	}
	a.Paket, err = kumpulkan(ctx, r.DB, sqlPaketInstansi, pindaiPaketPanel, f.InstansiID, f.JenisPengadaan, f.TahunMinimal)
	return a, err
}

func pindaiPaketPanel(b pgx.CollectableRow) (domain.PaketJadwalPanel, error) {
	var p domain.PaketJadwalPanel
	return p, b.Scan(&p.Kode, &p.Nama, &p.JenisPengadaan, &p.NomorBeritaAcara, &p.IsSimulation,
		&p.Jadwal.Kategori, &p.Jadwal.TanggalSerahTerima, &p.Jadwal.JadwalSelesai, &p.Jadwal.TanggalKontrak)
}

func (r PanelRepository) AntreanAnggota(ctx context.Context, anggotaID string) ([]domain.UmpanPanel, error) {
	return kumpulkan(ctx, r.DB, sqlUmpanPanelAnggota, pgx.RowToStructByPos[domain.UmpanPanel], anggotaID, domain.BatasAntreanPanel)
}
