package repository

import (
	"context"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlFraksiAda = `SELECT EXISTS (SELECT 1 FROM anggota WHERE fraksi = $1)`

	sqlSikapFraksi = `SELECT fraksi, agenda, judul_agenda, tanggal, ringkasan, url, provenance_id
FROM sikap_fraksi
WHERE fraksi = $1
ORDER BY tanggal DESC, agenda, judul_agenda, url`

	sqlDapilAda = `SELECT EXISTS (SELECT 1 FROM dapil WHERE id = $1)`

	sqlKolomAspirasiReses = `SELECT a.id, a.masa_reses, a.tanggal_paripurna, a.bidang, a.uraian, coalesce(a.kode_wilayah, ''),
    coalesce(w.nama, ''), coalesce(w.tingkat::text, ''), a.url, a.provenance_id
FROM aspirasi_reses a
LEFT JOIN wilayah w ON w.kode = a.kode_wilayah`

	sqlAmbilAspirasiReses = sqlKolomAspirasiReses + ` WHERE a.id = $1`

	sqlAspirasiReses = sqlKolomAspirasiReses + `
WHERE ($1 = 0 OR a.kode_wilayah IS NULL OR EXISTS (
        SELECT 1 FROM dapil_wilayah dw
        WHERE dw.dapil_id = $1 AND (a.kode_wilayah = dw.kode_kecamatan OR a.kode_wilayah LIKE dw.kode_kecamatan || '.%')))
    AND ($2 = '' OR a.bidang = $2)
ORDER BY a.tanggal_paripurna DESC, a.kode_wilayah IS NULL, a.bidang, a.uraian`

	sqlPaketTerkait = `SELECT p.kode, p.nama, p.tahap, p.jenis_pengadaan, p.metode, p.tahun_anggaran, p.hps, p.nilai_kontrak,
    coalesce(p.sektor, ''), p.tema_kode, p.lokasi_tingkat::text, p.kode_wilayah, w.nama, p.is_simulation, p.provenance_id,
    coalesce(s.nama, ''), coalesce(p.kategori = ANY ($3), false)
FROM paket p
JOIN wilayah w ON w.kode = p.kode_wilayah
LEFT JOIN penyedia s ON s.id = p.penyedia_id
WHERE (p.kode_wilayah = $1 OR p.kode_wilayah LIKE $1 || '.%')
    AND p.tahun_anggaran = $2 AND NOT p.is_batal AND NOT p.is_simulation
    AND (p.kategori = ANY ($3) OR p.nama ~* ANY ($4))
ORDER BY p.hps DESC NULLS LAST, p.kode
LIMIT $5`
)

func (r ParipurnaRepository) IsFraksiAda(ctx context.Context, fraksi string) (bool, error) {
	var ada bool
	return ada, r.DB.QueryRow(ctx, sqlFraksiAda, fraksi).Scan(&ada)
}

func (r ParipurnaRepository) SikapFraksi(ctx context.Context, fraksi string) ([]domain.SikapFraksi, error) {
	baris, err := r.DB.Query(ctx, sqlSikapFraksi, fraksi)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, pgx.RowToStructByPos[domain.SikapFraksi])
}

func (r ParipurnaRepository) IsDapilAda(ctx context.Context, id int) (bool, error) {
	var ada bool
	return ada, r.DB.QueryRow(ctx, sqlDapilAda, id).Scan(&ada)
}

func (r ParipurnaRepository) AspirasiReses(ctx context.Context, filter domain.FilterAspirasiReses) ([]domain.AspirasiReses, error) {
	baris, err := r.DB.Query(ctx, sqlAspirasiReses, filter.DapilID, filter.Bidang)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, pgx.RowToStructByPos[domain.AspirasiReses])
}

func (r ParipurnaRepository) AmbilAspirasiReses(ctx context.Context, id int64) (domain.AspirasiReses, bool, error) {
	baris, err := r.DB.Query(ctx, sqlAmbilAspirasiReses, id)
	if err != nil {
		return domain.AspirasiReses{}, false, err
	}
	return hasilSatuBaris(pgx.CollectExactlyOneRow(baris, pgx.RowToStructByPos[domain.AspirasiReses]))
}

func (r ParipurnaRepository) PaketTerkait(ctx context.Context, k domain.KueriPaketTerkait) ([]domain.PaketTerkait, error) {
	baris, err := r.DB.Query(ctx, sqlPaketTerkait, k.KodeWilayah, k.TahunAnggaran, k.Kategori, k.PolaKataKunci, k.Batas)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.PaketTerkait, error) {
		var p domain.PaketTerkait
		err := b.Scan(&p.Kode, &p.Nama, &p.Tahap, &p.JenisPengadaan, &p.Metode, &p.TahunAnggaran, &p.HPS, &p.NilaiKontrak,
			&p.Sektor, &p.TemaKode, &p.Lokasi.Tingkat, &p.Lokasi.KodeWilayah, &p.Lokasi.NamaWilayah, &p.IsSimulation,
			&p.ProvenanceID, &p.NamaPenyedia, &p.IsCocokKategori)
		return p, err
	})
}
