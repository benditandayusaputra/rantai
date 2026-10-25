package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

const (
	sqlAmbilPaket = `SELECT p.kode, p.kode_spse, p.nama, p.tahap, p.jenis_pengadaan, p.metode, coalesce(p.jenis_kontrak, ''),
    p.is_tender, p.is_batal, p.is_simulation, p.tahun_anggaran, p.pagu, p.hps, p.nilai_kontrak, p.jumlah_peserta,
    p.tanggal_dibuat, p.tanggal_kontrak, p.tanggal_serah_terima, i.nama, coalesce(p.kategori, ''), coalesce(p.sektor, ''),
    p.tema_kode, coalesce(p.lokasi_tingkat::text, ''), coalesce(p.kode_wilayah, ''), coalesce(w.nama, ''),
    ST_Y(p.lokasi), ST_X(p.lokasi), p.rt, p.rw, p.provenance_id,
    r.kode, r.nama, coalesce(r.sumber_dana, ''), r.pagu, rp.sumber = 'sirup', r.jadwal_pelaksanaan_mulai,
    r.jadwal_pelaksanaan_selesai, r.provenance_id, s.id, s.nama, coalesce(p.serah_terima_provenance_id, ''),
    coalesce(p.nomor_berita_acara, '')
FROM paket p
JOIN instansi i ON i.id = p.instansi_id
LEFT JOIN wilayah w ON w.kode = p.kode_wilayah
LEFT JOIN rup r ON r.kode = p.kode_rup
LEFT JOIN provenance rp ON rp.id = r.provenance_id
LEFT JOIN penyedia s ON s.id = p.penyedia_id
WHERE p.kode = $1`

	sqlSinyalPaket = `SELECT id, jenis, alasan, pembanding FROM sinyal WHERE paket_kode = $1 AND aktif ORDER BY jenis`

	sqlFilterKecamatan = `(p.kode_wilayah = $1 OR p.kode_wilayah LIKE $1 || '.%')
    AND ($2 = '' OR $2 = ANY (p.tema_kode)) AND ($3 = 0 OR p.tahun_anggaran = $3)`

	sqlDaftarPaketKecamatan = `SELECT p.kode, p.nama, p.tahap, p.jenis_pengadaan, p.metode, p.tahun_anggaran, p.hps, p.nilai_kontrak,
    coalesce(p.sektor, ''), p.tema_kode, p.lokasi_tingkat::text, p.kode_wilayah, w.nama, p.is_simulation, p.provenance_id,
    coalesce(s.nama, ''), count(*) OVER ()
FROM paket p
JOIN wilayah w ON w.kode = p.kode_wilayah
LEFT JOIN penyedia s ON s.id = p.penyedia_id
WHERE ` + sqlFilterKecamatan + `
ORDER BY p.tahun_anggaran DESC, p.hps DESC NULLS LAST, p.kode
LIMIT $4 OFFSET $5`
)

const sqlUnduhPaket = `SELECT p.kode, p.nama, p.tahap, p.jenis_pengadaan, p.metode, p.tahun_anggaran, p.hps, p.nilai_kontrak,
    coalesce(p.sektor, ''), p.tema_kode, p.lokasi_tingkat::text, p.kode_wilayah, w.nama, p.is_simulation, p.provenance_id,
    coalesce(s.nama, ''), p.kode_spse, p.is_tender, p.pagu, i.nama, coalesce(p.kode_rup, ''), v.label, v.diambil_pada
FROM paket p
JOIN wilayah w ON w.kode = p.kode_wilayah
JOIN instansi i ON i.id = p.instansi_id
JOIN provenance v ON v.id = p.provenance_id
LEFT JOIN penyedia s ON s.id = p.penyedia_id
WHERE (p.kode_wilayah = $1 OR p.kode_wilayah LIKE $1 || '.%') AND ($2 = 0 OR p.tahun_anggaran = $2)
ORDER BY p.tahun_anggaran DESC, p.hps DESC NULLS LAST, p.kode
LIMIT $3`

type rupBaca struct {
	kode, nama, sumberDana *string
	pagu                   *float64
	isDetail               *bool
	mulai, selesai         *time.Time
	provenanceID           *string
}

func (r PaketRepository) AmbilPaket(ctx context.Context, kode string) (domain.PaketLengkap, bool, error) {
	p, err := pindaiPaketLengkap(r.DB.QueryRow(ctx, sqlAmbilPaket, kode))
	p, ada, err := hasilSatuBaris(p, err)
	if err != nil || !ada {
		return p, ada, err
	}
	if p.Sinyal, err = r.sinyalPaket(ctx, kode); err != nil {
		return p, true, err
	}
	return p, true, r.lampirkanHakJawab(ctx, &p)
}

func (r PaketRepository) lampirkanHakJawab(ctx context.Context, p *domain.PaketLengkap) error {
	paket, err := hakJawabUntuk(ctx, r.DB, domain.TargetHakJawabPaket, []string{p.Kode})
	if err != nil {
		return err
	}
	p.HakJawab = hakJawabAtau(paket, p.Kode)
	idSinyal := make([]string, len(p.Sinyal))
	for i, s := range p.Sinyal {
		idSinyal[i] = strconv.FormatInt(s.ID, 10)
	}
	sinyal, err := hakJawabUntuk(ctx, r.DB, domain.TargetHakJawabSinyal, idSinyal)
	for i := range p.Sinyal {
		p.Sinyal[i].HakJawab = hakJawabAtau(sinyal, idSinyal[i])
	}
	return err
}

func pindaiPaketLengkap(baris pgx.Row) (domain.PaketLengkap, error) {
	var p domain.PaketLengkap
	var lat, lng *float64
	var rup rupBaca
	var idPenyedia *int64
	var namaPenyedia *string
	err := baris.Scan(&p.Kode, &p.KodeSPSE, &p.Nama, &p.Tahap, &p.JenisPengadaan, &p.Metode, &p.JenisKontrak,
		&p.IsTender, &p.IsBatal, &p.IsSimulation, &p.TahunAnggaran, &p.Pagu, &p.HPS, &p.NilaiKontrak, &p.JumlahPeserta,
		&p.TanggalDibuat, &p.TanggalKontrak, &p.TanggalSerahTerima, &p.Instansi, &p.Kategori, &p.Sektor,
		&p.TemaKode, &p.Lokasi.Tingkat, &p.Lokasi.KodeWilayah, &p.Lokasi.NamaWilayah,
		&lat, &lng, &p.Lokasi.RT, &p.Lokasi.RW, &p.ProvenanceID,
		&rup.kode, &rup.nama, &rup.sumberDana, &rup.pagu, &rup.isDetail, &rup.mulai,
		&rup.selesai, &rup.provenanceID, &idPenyedia, &namaPenyedia, &p.SerahTerimaProvID, &p.NomorBeritaAcara)
	if lat != nil && lng != nil {
		p.Lokasi.Titik = &domain.Titik{Lat: *lat, Lng: *lng}
	}
	if idPenyedia != nil {
		p.Penyedia = &domain.PenyediaPaket{ID: *idPenyedia, Nama: *namaPenyedia}
	}
	p.RUP = rup.keDomain()
	return p, err
}

func (r rupBaca) keDomain() *domain.RUPPaket {
	if r.kode == nil {
		return nil
	}
	return &domain.RUPPaket{
		Kode: *r.kode, Nama: *r.nama, SumberDana: *r.sumberDana, Pagu: r.pagu,
		IsDetailTersedia: r.isDetail != nil && *r.isDetail,
		PelaksanaanMulai: r.mulai, PelaksanaanSelesai: r.selesai, ProvenanceID: *r.provenanceID,
	}
}

func (r PaketRepository) sinyalPaket(ctx context.Context, kode string) ([]domain.SinyalTersimpan, error) {
	baris, err := r.DB.Query(ctx, sqlSinyalPaket, kode)
	if err != nil {
		return nil, fmt.Errorf("baca sinyal paket %s: %w", kode, err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.SinyalTersimpan, error) {
		var s domain.SinyalTersimpan
		var pembanding []byte
		if err := b.Scan(&s.ID, &s.Jenis, &s.Alasan, &pembanding); err != nil {
			return s, err
		}
		return s, json.Unmarshal(pembanding, &s.Pembanding)
	})
}

func (r PaketRepository) BarisUnduhPaket(ctx context.Context, f usecase.FilterUnduhPaket) ([]domain.BarisUnduhPaket, error) {
	baris, err := r.DB.Query(ctx, sqlUnduhPaket, f.KodeWilayah, f.Tahun, f.Batas)
	if err != nil {
		return nil, fmt.Errorf("baris unduh paket %s: %w", f.KodeWilayah, err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.BarisUnduhPaket, error) {
		var p domain.BarisUnduhPaket
		err := b.Scan(&p.Kode, &p.Nama, &p.Tahap, &p.JenisPengadaan, &p.Metode, &p.TahunAnggaran, &p.HPS, &p.NilaiKontrak,
			&p.Sektor, &p.TemaKode, &p.Lokasi.Tingkat, &p.Lokasi.KodeWilayah, &p.Lokasi.NamaWilayah, &p.IsSimulation,
			&p.ProvenanceID, &p.NamaPenyedia, &p.KodeSPSE, &p.IsTender, &p.Pagu, &p.Instansi, &p.KodeRUP, &p.Label, &p.DiambilPada)
		return p, err
	})
}

func (r PaketRepository) DaftarPaketKecamatan(ctx context.Context, f usecase.FilterPaketKecamatan) ([]domain.PaketRingkas, int, error) {
	baris, err := r.DB.Query(ctx, sqlDaftarPaketKecamatan, f.KodeKecamatan, f.Tema, f.Tahun, f.PerHalaman, (f.Halaman-1)*f.PerHalaman)
	if err != nil {
		return nil, 0, fmt.Errorf("daftar paket kecamatan %s: %w", f.KodeKecamatan, err)
	}
	total := 0
	daftar, err := pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.PaketRingkas, error) {
		var p domain.PaketRingkas
		err := b.Scan(&p.Kode, &p.Nama, &p.Tahap, &p.JenisPengadaan, &p.Metode, &p.TahunAnggaran, &p.HPS, &p.NilaiKontrak,
			&p.Sektor, &p.TemaKode, &p.Lokasi.Tingkat, &p.Lokasi.KodeWilayah, &p.Lokasi.NamaWilayah, &p.IsSimulation,
			&p.ProvenanceID, &p.NamaPenyedia, &total)
		return p, err
	})
	return daftar, total, err
}
