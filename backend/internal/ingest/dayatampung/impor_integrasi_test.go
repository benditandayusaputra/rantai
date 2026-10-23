package dayatampung

import (
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwilayah"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

type pelaporTiruan struct {
	laporan []usecase.LaporanDayaTampung
}

func (p *pelaporTiruan) LaporkanDayaTampung(l []usecase.LaporanDayaTampung) { p.laporan = l }

func siapkanSekolahUji(t *testing.T, tx pgx.Tx) {
	t.Helper()
	ujiwilayah.MuatDKI(t, tx)
	_, err := tx.Exec(t.Context(), `INSERT INTO sekolah (npsn, nama, bentuk, jenjang, status, is_madrasah, kode_kecamatan, provenance_id)
SELECT n, 'SMP UJI ' || n, 'SMP', 'menengah_pertama', 'negeri', false, '31.01.01', (SELECT provenance_id FROM wilayah WHERE kode = '31.01.01')
FROM unnest(ARRAY['20100251', '20107183']) AS n
ON CONFLICT (npsn) DO NOTHING`)
	if err != nil {
		t.Fatal(err)
	}
}

func jalankanImporUji(t *testing.T, tx pgx.Tx) []usecase.LaporanDayaTampung {
	t.Helper()
	pelapor := &pelaporTiruan{}
	impor := usecase.ImporDayaTampung{
		Pengambil: Sumber{
			Konfigurasi: Konfigurasi{TahunAjaran: "2026/2027", Berkas: []KonfigurasiBerkas{
				{Path: "testdata/daya_tampung_negeri.csv", Jalur: domain.JalurNegeri, URL: "https://contoh.test/sk258.pdf"},
			}},
			Sekarang: time.Now,
		},
		Penyimpan: repository.DayaTampungRepository{DB: tx},
		Pencatat:  repository.IngestRunRepository{DB: tx},
		Pelapor:   pelapor,
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	return pelapor.laporan
}

func TestImporDayaTampung_NPSNTakTerdaftar_DikeluarkanDanDilaporkan(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanSekolahUji(t, tx)

	laporan := jalankanImporUji(t, tx)

	var total int
	err := tx.QueryRow(t.Context(), `SELECT coalesce(sum(kursi), 0) FROM daya_tampung
WHERE jalur = 'negeri' AND tahun_ajaran = '2026/2027' AND npsn IN ('20100251', '20107183')`).Scan(&total)
	if err != nil || total != 288 {
		t.Errorf("total kursi = %d, err %v", total, err)
	}
	l := laporan[0]
	if l.KursiDokumen != 396 || l.KursiTerimpor != 288 || len(l.TidakTerpetakan) != 1 || l.TidakTerpetakan[0].Kursi != 108 {
		t.Errorf("laporan = %+v", l)
	}
}

func TestImporDayaTampung_DijalankanDuaKali_TidakMenggandakanData(t *testing.T) {
	tx := ujidb.BukaTransaksi(t)
	siapkanSekolahUji(t, tx)

	jalankanImporUji(t, tx)
	jalankanImporUji(t, tx)

	var baru, berubah int
	err := tx.QueryRow(t.Context(),
		`SELECT jumlah_baru, jumlah_berubah FROM ingest_run WHERE sumber = 'daya_tampung' ORDER BY id DESC LIMIT 1`).Scan(&baru, &berubah)
	if err != nil || baru != 0 || berubah != 0 {
		t.Errorf("baru %d, berubah %d, err %v", baru, berubah, err)
	}
}
