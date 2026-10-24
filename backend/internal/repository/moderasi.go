package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	batasAntreanLaporanModerasi = 100

	sqlAntreanLaporanModerasi = `SELECT l.id, l.kondisi, l.catatan, w.nama, coalesce(p.nama, ''), l.is_simulation, l.dibuat_pada
FROM laporan l
JOIN wilayah w ON w.kode = l.kode_wilayah
LEFT JOIN paket p ON p.kode = l.paket_kode
WHERE l.status = 'MENUNGGU_SAKSI'
ORDER BY l.dibuat_pada, l.id
LIMIT $1`

	sqlFotoAntreanModerasi = `SELECT laporan_id, id, object_key, thumbnail_key, disembunyikan
FROM laporan_foto WHERE laporan_id = ANY ($1) ORDER BY laporan_id, urutan`

	sqlUmpanAntreanModerasi = `SELECT u.id, u.jenis_sumber, u.media_nama, u.judul, u.url, u.terbit_pada, array_agg(a.nama ORDER BY a.nama)
FROM umpan_kegiatan u
JOIN umpan_anggota ua ON ua.umpan_id = u.id
JOIN anggota a ON a.id = ua.anggota_id
WHERE NOT u.disembunyikan AND u.jenis_sumber = 'media'
GROUP BY u.id
ORDER BY u.terbit_pada DESC, u.id DESC
LIMIT $1`

	sqlSembunyikanFoto = `UPDATE laporan_foto SET disembunyikan = true WHERE id = $1 AND laporan_id = $2`

	sqlSembunyikanUmpan = `UPDATE umpan_kegiatan SET disembunyikan = true, alasan_disembunyikan = $2, diubah_pada = now() WHERE id = $1`

	sqlCatatModerasi = `INSERT INTO moderasi_log (moderator_id, target_jenis, target_id, aksi, alasan) VALUES ($1, $2, $3, $4, $5)`
)

type ModerasiRepository struct {
	LaporanRepository
}

func (r ModerasiRepository) AntreanModerasi(ctx context.Context) (domain.AntreanModerasi, error) {
	laporan, err := kumpulkan(ctx, r.DB, sqlAntreanLaporanModerasi, func(b pgx.CollectableRow) (domain.LaporanModerasi, error) {
		var l domain.LaporanModerasi
		return l, b.Scan(&l.ID, &l.Kondisi, &l.Catatan, &l.NamaWilayah, &l.PaketNama, &l.IsSimulation, &l.DibuatPada)
	}, batasAntreanLaporanModerasi)
	if err != nil {
		return domain.AntreanModerasi{}, fmt.Errorf("baca antrean laporan: %w", err)
	}
	if err := r.lengkapiFotoAntrean(ctx, laporan); err != nil {
		return domain.AntreanModerasi{}, err
	}
	umpan, err := kumpulkan(ctx, r.DB, sqlUmpanAntreanModerasi, pgx.RowToStructByPos[domain.UmpanModerasi], domain.BatasAntreanUmpanModerasi)
	if err != nil {
		return domain.AntreanModerasi{}, fmt.Errorf("baca antrean umpan: %w", err)
	}
	return domain.AntreanModerasi{Laporan: laporan, Umpan: umpan}, nil
}

func (r ModerasiRepository) lengkapiFotoAntrean(ctx context.Context, laporan []domain.LaporanModerasi) error {
	indeks := map[int64]int{}
	id := make([]int64, len(laporan))
	for i, l := range laporan {
		indeks[l.ID], id[i] = i, l.ID
		laporan[i].Foto = []domain.FotoModerasi{}
	}
	baris, err := r.DB.Query(ctx, sqlFotoAntreanModerasi, id)
	if err != nil {
		return fmt.Errorf("baca foto antrean: %w", err)
	}
	var laporanID int64
	var f domain.FotoModerasi
	_, err = pgx.ForEachRow(baris, []any{&laporanID, &f.ID, &f.ObjectKey, &f.ThumbnailKey, &f.Disembunyikan}, func() error {
		laporan[indeks[laporanID]].Foto = append(laporan[indeks[laporanID]].Foto, f)
		return nil
	})
	return err
}

func (r ModerasiRepository) SembunyikanFoto(ctx context.Context, laporanID, fotoID int64) (bool, error) {
	tag, err := r.DB.Exec(ctx, sqlSembunyikanFoto, fotoID, laporanID)
	return tag.RowsAffected() > 0, err
}

func (r ModerasiRepository) SembunyikanUmpan(ctx context.Context, id int64, alasan string) (bool, error) {
	tag, err := r.DB.Exec(ctx, sqlSembunyikanUmpan, id, alasan)
	return tag.RowsAffected() > 0, err
}

func (r ModerasiRepository) CatatModerasi(ctx context.Context, c domain.CatatanModerasi) error {
	_, err := r.DB.Exec(ctx, sqlCatatModerasi, c.ModeratorID, c.TargetJenis, c.TargetID, c.Aksi, c.Alasan)
	return err
}
