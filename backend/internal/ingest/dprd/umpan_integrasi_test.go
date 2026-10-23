package dprd_test

import (
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb/ujiwakil"
	"github.com/jackc/pgx/v5"
)

func siapkanWakilDanTema(t *testing.T) pgx.Tx {
	t.Helper()
	tx := ujidb.BukaTransaksi(t)
	ujiwakil.MuatDPRD(t, tx)
	ujiwakil.MuatTema(t, tx)
	return tx
}

func TestImporUmpanDPRD_RilisMenyebutAnggota_DipautkanLewatSitusResmi(t *testing.T) {
	tx := siapkanWakilDanTema(t)

	ujiwakil.MuatUmpanDPRD(t, tx)

	var dasar string
	err := tx.QueryRow(t.Context(), `SELECT ua.dasar_cocok FROM umpan_anggota ua JOIN umpan_kegiatan u ON u.id = ua.umpan_id
WHERE u.judul = 'Kasus DBD Jakarta Selatan Tembus 1.314' AND ua.anggota_id = 'yuke-yurike-st-mm'`).Scan(&dasar)
	if err != nil || dasar != "situs_resmi" {
		t.Errorf("harap rilis DBD dipautkan ke Yuke Yurike, dapat %q, %v", dasar, err)
	}
}

func TestImporUmpanDPRD_IsiBerita_TidakDisimpan(t *testing.T) {
	tx := siapkanWakilDanTema(t)

	ujiwakil.MuatUmpanDPRD(t, tx)

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

	ujiwakil.MuatUmpanDPRD(t, tx)
	pertama := hitung(t, tx, `SELECT count(*) FROM umpan_anggota`)
	ujiwakil.MuatUmpanDPRD(t, tx)

	if kedua := hitung(t, tx, `SELECT count(*) FROM umpan_anggota`); pertama == 0 || kedua != pertama {
		t.Errorf("pautan pertama %d, kedua %d", pertama, kedua)
	}
}
