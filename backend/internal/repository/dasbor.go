package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/jackc/pgx/v5"
)

const (
	BatasInstansiDasbor = 10

	sqlSyaratPaketDasbor = `NOT p.is_batal AND NOT p.is_simulation`
	sqlTahunTerbitUmpan  = `extract(year FROM u.terbit_pada AT TIME ZONE 'Asia/Jakarta')::int`
	sqlSyaratUmpanDasbor = `NOT u.disembunyikan`

	sqlRujukanPaket = `SELECT pr.id, pr.label FROM paket p JOIN provenance pr ON pr.id = p.provenance_id
WHERE ` + sqlSyaratPaketDasbor + ` ORDER BY pr.diambil_pada DESC LIMIT 1`
	sqlRujukanUmpan = `SELECT pr.id, pr.label FROM umpan_kegiatan u JOIN provenance pr ON pr.id = u.provenance_id
WHERE ` + sqlSyaratUmpanDasbor + ` ORDER BY pr.diambil_pada DESC LIMIT 1`
	sqlRujukanReses = `SELECT pr.id, pr.label FROM aspirasi_reses a JOIN provenance pr ON pr.id = a.provenance_id
ORDER BY pr.diambil_pada DESC LIMIT 1`
	sqlRujukanIndikator = `SELECT coalesce((array_agg(pr.id ORDER BY pr.diambil_pada DESC))[1], ''), coalesce(max(i.label)::text, '')
FROM indikator_kecamatan i JOIN provenance pr ON pr.id = i.provenance_id WHERE i.tema_kode = $1`

	sqlPaketTahunan = `SELECT p.tahun_anggaran::text, '', count(*), coalesce(sum(p.hps), 0)
FROM paket p WHERE ` + sqlSyaratPaketDasbor + ` GROUP BY 1 ORDER BY 1`
	sqlPaketBulanan = `SELECT to_char(p.tanggal_dibuat, 'YYYY-MM'), '', count(*), coalesce(sum(p.hps), 0)
FROM paket p WHERE ` + sqlSyaratPaketDasbor + ` AND extract(year FROM p.tanggal_dibuat)::int = $1 GROUP BY 1 ORDER BY 1`
	sqlPaketPerJenis = `SELECT p.jenis_pengadaan, '', count(*), coalesce(sum(p.hps), 0)
FROM paket p WHERE ` + sqlSyaratPaketDasbor + ` AND p.tahun_anggaran = $1 GROUP BY 1 ORDER BY 4 DESC, 3 DESC`
	sqlPaketPerMetode = `SELECT p.metode, '', count(*), coalesce(sum(p.hps), 0)
FROM paket p WHERE ` + sqlSyaratPaketDasbor + ` AND p.tahun_anggaran = $1 GROUP BY 1 ORDER BY 4 DESC, 3 DESC`
	sqlPaketPerKota = `SELECT w.kode, w.nama, count(*), coalesce(sum(p.hps), 0)
FROM paket p JOIN wilayah w ON w.kode = left(p.kode_wilayah, 5)
WHERE ` + sqlSyaratPaketDasbor + ` AND p.tahun_anggaran = $1 GROUP BY w.kode, w.nama ORDER BY 4 DESC`
	sqlPaketPerKecamatan = `SELECT w.kode, w.nama, count(*), coalesce(sum(p.hps), 0)
FROM paket p JOIN wilayah w ON w.kode = left(p.kode_wilayah, 8) AND w.tingkat = 'kecamatan'
WHERE ` + sqlSyaratPaketDasbor + ` AND p.tahun_anggaran = $1 GROUP BY w.kode, w.nama ORDER BY 4 DESC`
	sqlPaketPerInstansi = `SELECT i.id::text, i.nama, count(*), coalesce(sum(p.hps), 0)
FROM paket p JOIN instansi i ON i.id = p.instansi_id
WHERE ` + sqlSyaratPaketDasbor + ` AND p.tahun_anggaran = $1 AND i.nama <> $3 GROUP BY i.id, i.nama ORDER BY 4 DESC LIMIT $2`

	sqlKegiatanTahunan = `SELECT ` + sqlTahunTerbitUmpan + `::text, '', count(*), 0::float8
FROM umpan_kegiatan u WHERE ` + sqlSyaratUmpanDasbor + ` GROUP BY 1 ORDER BY 1`
	sqlKegiatanBulanan = `SELECT to_char(u.terbit_pada AT TIME ZONE 'Asia/Jakarta', 'YYYY-MM'), '', count(*), 0::float8
FROM umpan_kegiatan u WHERE ` + sqlSyaratUmpanDasbor + ` AND ` + sqlTahunTerbitUmpan + ` = $1 GROUP BY 1 ORDER BY 1`
	sqlTopikKegiatan = `SELECT tag, '', count(*), 0::float8
FROM umpan_kegiatan u CROSS JOIN unnest(u.tag_isu) AS tag
WHERE ` + sqlSyaratUmpanDasbor + ` AND ` + sqlTahunTerbitUmpan + ` = $1 GROUP BY 1 ORDER BY 3 DESC, 1`
	sqlResesPerBidang = `SELECT a.bidang, '', count(*), 0::float8
FROM aspirasi_reses a WHERE extract(year FROM a.tanggal_paripurna)::int = $1 GROUP BY 1 ORDER BY 3 DESC, 1`

	sqlKursiPerPeriode = `SELECT i.periode, '', i.kelas, count(*), sum(i.kebutuhan), sum(i.kapasitas)
FROM indikator_kecamatan i WHERE i.tema_kode = $1 GROUP BY 1, 3 ORDER BY 1`
	sqlKursiPerKota = `WITH terbaru AS (SELECT max(periode) AS periode FROM indikator_kecamatan WHERE tema_kode = $1)
SELECT w.kode, w.nama, i.kelas, count(*), sum(i.kebutuhan), sum(i.kapasitas)
FROM indikator_kecamatan i JOIN terbaru t ON t.periode = i.periode JOIN wilayah w ON w.kode = left(i.kode_kecamatan, 5)
WHERE i.tema_kode = $1 GROUP BY w.kode, w.nama, i.kelas ORDER BY w.kode`
)

type DasborRepository struct {
	DB Basisdata
}

type kueriSeri struct {
	sql  string
	tuju *domain.Seri
	args []any
}

func (r DasborRepository) PaketTahunan(ctx context.Context) (domain.Seri, error) {
	return r.seri(ctx, sqlRujukanPaket, sqlPaketTahunan)
}

func (r DasborRepository) KegiatanTahunan(ctx context.Context) (domain.Seri, error) {
	return r.seri(ctx, sqlRujukanUmpan, sqlKegiatanTahunan)
}

func (r DasborRepository) RincianUang(ctx context.Context, tahun int) (domain.DasborUang, error) {
	var d domain.DasborUang
	err := r.isiSeri(ctx, sqlRujukanPaket, []kueriSeri{
		{sqlPaketBulanan, &d.Bulanan, []any{tahun}},
		{sqlPaketPerJenis, &d.Jenis, []any{tahun}},
		{sqlPaketPerMetode, &d.Metode, []any{tahun}},
		{sqlPaketPerKota, &d.Kota, []any{tahun}},
		{sqlPaketPerKecamatan, &d.Kecamatan, []any{tahun}},
	})
	return d, err
}

func (r DasborRepository) RincianTanggungJawab(ctx context.Context, tahun int, instansiTanpaRincian string) (domain.DasborTanggungJawab, error) {
	var d domain.DasborTanggungJawab
	if err := r.isiSeri(ctx, sqlRujukanUmpan, []kueriSeri{
		{sqlKegiatanBulanan, &d.KegiatanBulanan, []any{tahun}},
		{sqlTopikKegiatan, &d.Topik, []any{tahun}},
	}); err != nil {
		return d, err
	}
	if err := r.isiSeri(ctx, sqlRujukanReses, []kueriSeri{{sqlResesPerBidang, &d.Reses, []any{tahun}}}); err != nil {
		return d, err
	}
	err := r.isiSeri(ctx, sqlRujukanPaket, []kueriSeri{
		{sqlPaketPerInstansi, &d.Instansi, []any{tahun, BatasInstansiDasbor, instansiTanpaRincian}},
	})
	return d, err
}

func (r DasborRepository) seri(ctx context.Context, sqlRujukan, sqlIsi string) (domain.Seri, error) {
	var s domain.Seri
	err := r.isiSeri(ctx, sqlRujukan, []kueriSeri{{sql: sqlIsi, tuju: &s}})
	return s, err
}

func (r DasborRepository) isiSeri(ctx context.Context, sqlRujukan string, daftar []kueriSeri) error {
	label, provenanceID, err := r.rujukan(ctx, sqlRujukan)
	if err != nil {
		return err
	}
	for _, k := range daftar {
		isi, err := r.porsi(ctx, k.sql, k.args...)
		if err != nil {
			return err
		}
		*k.tuju = domain.Seri{Label: label, ProvenanceID: provenanceID, Isi: isi}
	}
	return nil
}

func (r DasborRepository) rujukan(ctx context.Context, sql string, args ...any) (domain.Label, string, error) {
	var id string
	var label domain.Label
	err := r.DB.QueryRow(ctx, sql, args...).Scan(&id, &label)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", nil
	}
	if err != nil {
		return "", "", fmt.Errorf("rujukan dasbor: %w", err)
	}
	return label, id, nil
}

func (r DasborRepository) porsi(ctx context.Context, sql string, args ...any) ([]domain.Porsi, error) {
	baris, err := r.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("porsi dasbor: %w", err)
	}
	return pgx.CollectRows(baris, func(b pgx.CollectableRow) (domain.Porsi, error) {
		var p domain.Porsi
		err := b.Scan(&p.Kunci, &p.Nama, &p.Jumlah, &p.Nilai)
		return p, err
	})
}

func (r DasborRepository) RekapKursi(ctx context.Context, kodeTema string) (domain.DasborLayanan, error) {
	d := domain.DasborLayanan{Tema: kodeTema}
	var err error
	if d.Label, d.ProvenanceID, err = r.rujukan(ctx, sqlRujukanIndikator, kodeTema); err != nil {
		return d, err
	}
	if d.Periode, err = r.rekap(ctx, sqlKursiPerPeriode, kodeTema); err != nil {
		return d, err
	}
	d.Kota, err = r.rekap(ctx, sqlKursiPerKota, kodeTema)
	return d, err
}

func (r DasborRepository) rekap(ctx context.Context, sql, kodeTema string) ([]domain.RekapKursi, error) {
	baris, err := r.DB.Query(ctx, sql, kodeTema)
	if err != nil {
		return nil, fmt.Errorf("rekap kursi: %w", err)
	}
	defer baris.Close()
	hasil := []domain.RekapKursi{}
	for baris.Next() {
		var kunci, nama string
		var kelas domain.KelasRasio
		var jumlah int
		var kebutuhan, kapasitas float64
		if err := baris.Scan(&kunci, &nama, &kelas, &jumlah, &kebutuhan, &kapasitas); err != nil {
			return nil, fmt.Errorf("baca rekap kursi: %w", err)
		}
		hasil = tambahKeRekap(hasil, domain.RekapKursi{Kunci: kunci, Nama: nama, Kebutuhan: kebutuhan, Kapasitas: kapasitas,
			PerKelas: map[domain.KelasRasio]int{kelas: jumlah}})
	}
	return hasil, baris.Err()
}

func tambahKeRekap(hasil []domain.RekapKursi, baru domain.RekapKursi) []domain.RekapKursi {
	if n := len(hasil); n > 0 && hasil[n-1].Kunci == baru.Kunci {
		akhir := &hasil[n-1]
		akhir.Kebutuhan += baru.Kebutuhan
		akhir.Kapasitas += baru.Kapasitas
		for kelas, jumlah := range baru.PerKelas {
			akhir.PerKelas[kelas] += jumlah
		}
		return hasil
	}
	return append(hasil, baru)
}
