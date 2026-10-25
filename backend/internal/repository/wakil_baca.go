package repository

import (
	"context"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	sqlKolomDapil = `d.id, d.nomor, d.nama, d.tahun_pemilu, d.jumlah_kursi, d.provenance_id,
    coalesce((SELECT a.provenance_id FROM anggota a WHERE a.dapil_id = d.id ORDER BY a.id LIMIT 1), '')`

	sqlDapilKecamatan = `SELECT w.kode, w.nama, w.tingkat, ` + sqlKolomDapil + `
FROM dapil_wilayah dw
JOIN dapil d ON d.id = dw.dapil_id
JOIN wilayah w ON w.kode = dw.kode_kecamatan
WHERE dw.kode_kecamatan = $1`

	sqlAmbilDapil = `SELECT ` + sqlKolomDapil + ` FROM dapil d WHERE d.id = $1`

	sqlKecamatanDapil = `SELECT w.kode, w.nama, w.tingkat
FROM dapil_wilayah dw JOIN wilayah w ON w.kode = dw.kode_kecamatan
WHERE dw.dapil_id = $1
ORDER BY w.nama`

	sqlAnggotaDapil = `SELECT a.id, a.nama, a.partai, a.fraksi, coalesce(a.komisi_kode, ''), coalesce(k.nama, ''), coalesce(a.foto_url, '')
FROM anggota a LEFT JOIN komisi k ON k.kode = a.komisi_kode
WHERE a.dapil_id = $1`

	sqlAmbilAnggota = `SELECT a.id, a.nama, a.partai, a.fraksi, coalesce(a.komisi_kode, ''), coalesce(k.nama, ''), coalesce(a.foto_url, ''),
    coalesce(k.bidang, '{}'), a.url_profil, a.lhkpn_url, a.provenance_id, ` + sqlKolomDapil + `
FROM anggota a
JOIN dapil d ON d.id = a.dapil_id
LEFT JOIN komisi k ON k.kode = a.komisi_kode
WHERE a.id = $1`

	sqlJabatanAnggota = `SELECT anggota_id, alat_kelengkapan, jabatan, ''
FROM anggota_akd
WHERE anggota_id = $1
ORDER BY alat_kelengkapan <> $2, alat_kelengkapan`
)

func (r WakilRepository) DapilKecamatan(ctx context.Context, kode string) (domain.RingkasanWilayah, domain.RingkasanDapil, bool, error) {
	var w domain.RingkasanWilayah
	var d domain.RingkasanDapil
	err := r.DB.QueryRow(ctx, sqlDapilKecamatan, kode).Scan(append([]any{&w.Kode, &w.Nama, &w.Tingkat}, kolomDapil(&d)...)...)
	_, ada, err := hasilSatuBaris(d, err)
	return w, d, ada, err
}

func (r WakilRepository) AmbilDapil(ctx context.Context, id int) (domain.RingkasanDapil, bool, error) {
	var d domain.RingkasanDapil
	return hasilSatuBaris(d, r.DB.QueryRow(ctx, sqlAmbilDapil, id).Scan(kolomDapil(&d)...))
}

func (r WakilRepository) IsiDapil(ctx context.Context, id int) ([]domain.RingkasanWilayah, []domain.Wakil, error) {
	baris, err := r.DB.Query(ctx, sqlKecamatanDapil, id)
	if err != nil {
		return nil, nil, err
	}
	kecamatan, err := pgx.CollectRows(baris, pgx.RowToStructByPos[domain.RingkasanWilayah])
	if err != nil {
		return nil, nil, fmt.Errorf("baca kecamatan dapil: %w", err)
	}
	baris, err = r.DB.Query(ctx, sqlAnggotaDapil, id)
	if err != nil {
		return nil, nil, err
	}
	anggota, err := pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.Wakil, error) {
		var w domain.Wakil
		return w, b.Scan(kolomWakil(&w)...)
	})
	return kecamatan, anggota, err
}

func (r WakilRepository) AmbilAnggota(ctx context.Context, id string) (domain.ProfilAnggota, bool, error) {
	var p domain.ProfilAnggota
	kolom := append(kolomWakil(&p.Wakil), &p.BidangKomisi, &p.URLProfil, &p.URLLHKPN, &p.ProvenanceID)
	err := r.DB.QueryRow(ctx, sqlAmbilAnggota, id).Scan(append(kolom, kolomDapil(&p.Dapil.RingkasanDapil)...)...)
	p, ada, err := hasilSatuBaris(p, err)
	if !ada || err != nil {
		return p, ada, err
	}
	baris, err := r.DB.Query(ctx, sqlJabatanAnggota, id, domain.AlatKelengkapanPimpinan)
	if err != nil {
		return p, false, err
	}
	p.AlatKelengkapan, err = pgx.CollectRows(baris, pgx.RowToStructByPos[domain.JabatanAKD])
	return p, true, err
}

func kolomDapil(d *domain.RingkasanDapil) []any {
	return []any{&d.ID, &d.Nomor, &d.Nama, &d.TahunPemilu, &d.JumlahKursi, &d.ProvenanceID, &d.ProvenanceAnggota}
}

func kolomWakil(w *domain.Wakil) []any {
	return []any{&w.ID, &w.Nama, &w.Partai, &w.Fraksi, &w.KodeKomisi, &w.NamaKomisi, &w.FotoURL}
}
