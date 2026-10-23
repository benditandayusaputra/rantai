package dprd_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/dprd"
	"github.com/benditandayusaputra/rantai/backend/internal/ingest/tema"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
	"github.com/jackc/pgx/v5"
)

type pengunduhRilisFixture struct{}

func (pengunduhRilisFixture) Ambil(context.Context, string) ([]byte, error) {
	return os.ReadFile("testdata/rilis_halaman_1.json")
}

func jalankanUmpanDPRDUji(t *testing.T, tx pgx.Tx) {
	t.Helper()
	tema := usecase.BacaTema{Pembaca: repository.TemaRepository{DB: tx}}
	impor := usecase.ImporUmpan{
		Job:         usecase.JobUmpanDPRD,
		JenisSumber: domain.SumberUmpanResmiDPRD,
		MulaiBawaan: time.Date(2024, time.August, 26, 0, 0, 0, 0, time.UTC),
		Pengambil: dprd.SumberRilis{
			Konfigurasi: dprd.KonfigurasiRilis{PerHalaman: 100, BatasHalaman: 1, MediaNama: "Situs DPRD", MediaDomain: "dprd-dkijakartaprov.go.id"},
			Pengunduh:   pengunduhRilisFixture{},
			DirRaw:      t.TempDir(),
			Sekarang:    time.Now,
		},
		Anggota:   repository.UmpanRepository{DB: tx},
		Tema:      tema,
		Penyimpan: repository.UmpanRepository{DB: tx},
		Pencatat:  repository.IngestRunRepository{DB: tx},
		Pautkan:   domain.PencariNama.PautkanRilisResmi,
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
}

func siapkanWakilDanTema(t *testing.T) pgx.Tx {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwakil.MuatDPRD(t, tx)
	impor := usecase.ImporTema{
		Pengambil: tema.Sumber{Dir: "../../../config/tema", URLRujukan: "https://contoh.test", Sekarang: time.Now},
		Penyimpan: repository.TemaRepository{DB: tx},
		Pencatat:  repository.IngestRunRepository{DB: tx},
	}
	if err := impor.Jalankan(t.Context()); err != nil {
		t.Fatal(err)
	}
	return tx
}

func TestImporUmpanDPRD_RilisMenyebutAnggota_DipautkanLewatSitusResmi(t *testing.T) {
	tx := siapkanWakilDanTema(t)

	jalankanUmpanDPRDUji(t, tx)

	var dasar string
	err := tx.QueryRow(t.Context(), `SELECT ua.dasar_cocok FROM umpan_anggota ua JOIN umpan_kegiatan u ON u.id = ua.umpan_id
WHERE u.judul = 'Kasus DBD Jakarta Selatan Tembus 1.314' AND ua.anggota_id = 'yuke-yurike-st-mm'`).Scan(&dasar)
	if err != nil || dasar != "situs_resmi" {
		t.Errorf("harap rilis DBD dipautkan ke Yuke Yurike, dapat %q, %v", dasar, err)
	}
}

func TestImporUmpanDPRD_IsiBerita_TidakDisimpan(t *testing.T) {
	tx := siapkanWakilDanTema(t)

	jalankanUmpanDPRDUji(t, tx)

	isiPenuh := hitung(t, tx, `SELECT count(*) FROM umpan_kegiatan u
WHERE char_length(u.cuplikan) > 300 OR u::text LIKE '%Keterlibatan ketua RT, RW, lurah%'`)
	if isiPenuh != 0 {
		t.Errorf("%d item menyimpan teks berita penuh", isiPenuh)
	}
	if n := hitung(t, tx, `SELECT count(*) FROM umpan_kegiatan u WHERE NOT EXISTS (SELECT 1 FROM umpan_anggota ua WHERE ua.umpan_id = u.id)`); n != 0 {
		t.Errorf("%d item tanpa anggota ikut tersimpan", n)
	}
}

func TestImporUmpanDPRD_DijalankanDuaKali_TidakMenggandakan(t *testing.T) {
	tx := siapkanWakilDanTema(t)

	jalankanUmpanDPRDUji(t, tx)
	pertama := hitung(t, tx, `SELECT count(*) FROM umpan_anggota`)
	jalankanUmpanDPRDUji(t, tx)

	if kedua := hitung(t, tx, `SELECT count(*) FROM umpan_anggota`); pertama == 0 || kedua != pertama {
		t.Errorf("pautan pertama %d, kedua %d", pertama, kedua)
	}
}
