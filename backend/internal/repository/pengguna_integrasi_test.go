package repository_test

import (
	"errors"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/repository"
	"github.com/benditandayusaputra/rantai/backend/internal/ujidb"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const emailPenggunaUji = "warga@contoh.id"

var saatUjiPengguna = time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)

func TestPenggunaRepository_PercobaanOTP_BertambahPadaKodeTerbaru(t *testing.T) {
	repo := repository.PenggunaRepository{DB: ujidb.BukaTransaksi(t)}
	for i, kode := range []string{"lama", "baru"} {
		otp := usecase.OTPBaru{Email: emailPenggunaUji, KodeHash: []byte(kode), Kedaluwarsa: saatUjiPengguna.Add(time.Hour), DibuatPada: saatUjiPengguna.Add(time.Duration(i) * time.Minute)}
		if err := repo.SimpanOTP(t.Context(), otp); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := repo.CatatPercobaanOTP(t.Context(), emailPenggunaUji); err != nil {
		t.Fatal(err)
	}
	kedua, err := repo.CatatPercobaanOTP(t.Context(), emailPenggunaUji)
	if err != nil {
		t.Fatal(err)
	}
	jumlah, errHitung := repo.HitungOTPSejak(t.Context(), emailPenggunaUji, saatUjiPengguna.Add(time.Second))

	if string(kedua.KodeHash) != "baru" || kedua.Percobaan != 2 || jumlah != 1 || errHitung != nil {
		t.Errorf("otp %+v, jumlah sejak %d, err %v", kedua, jumlah, errHitung)
	}
}

func TestPenggunaRepository_OTPDipakaiDuaKali_KeduaDitolak(t *testing.T) {
	repo := repository.PenggunaRepository{DB: ujidb.BukaTransaksi(t)}
	if err := repo.SimpanOTP(t.Context(), usecase.OTPBaru{Email: emailPenggunaUji, KodeHash: []byte("h"), Kedaluwarsa: saatUjiPengguna, DibuatPada: saatUjiPengguna}); err != nil {
		t.Fatal(err)
	}
	otp, err := repo.CatatPercobaanOTP(t.Context(), emailPenggunaUji)
	if err != nil {
		t.Fatal(err)
	}

	errPertama := repo.TandaiOTPDipakai(t.Context(), otp.ID)
	errKedua := repo.TandaiOTPDipakai(t.Context(), otp.ID)
	_, errEmailLain := repo.CatatPercobaanOTP(t.Context(), "lain@contoh.id")

	if errPertama != nil || !errors.Is(errKedua, domain.ErrOTPTidakBerlaku) || !errors.Is(errEmailLain, domain.ErrOTPTidakBerlaku) {
		t.Errorf("pertama %v, kedua %v, email lain %v", errPertama, errKedua, errEmailLain)
	}
}

func TestPenggunaRepository_PastikanPengguna_WargaBaruLaluSama(t *testing.T) {
	repo := repository.PenggunaRepository{DB: ujidb.BukaTransaksi(t)}

	pertama, err := repo.PastikanPengguna(t.Context(), emailPenggunaUji)
	if err != nil {
		t.Fatal(err)
	}
	kedua, err := repo.PastikanPengguna(t.Context(), emailPenggunaUji)
	if err != nil {
		t.Fatal(err)
	}
	diambil, err := repo.AmbilPengguna(t.Context(), pertama.ID)

	if pertama.Peran != domain.PeranWarga || kedua.ID != pertama.ID || diambil.Email != emailPenggunaUji || err != nil {
		t.Errorf("pertama %+v, kedua %+v, diambil %+v, err %v", pertama, kedua, diambil, err)
	}
}

func TestPenggunaRepository_CabutSesi_MengembalikanKeadaanSebelumDicabut(t *testing.T) {
	repo := repository.PenggunaRepository{DB: ujidb.BukaTransaksi(t)}
	pengguna, err := repo.PastikanPengguna(t.Context(), emailPenggunaUji)
	if err != nil {
		t.Fatal(err)
	}
	for _, token := range []string{"satu", "dua"} {
		if err := repo.SimpanSesi(t.Context(), usecase.SesiBaru{PenggunaID: pengguna.ID, TokenHash: []byte(token), Kedaluwarsa: saatUjiPengguna}); err != nil {
			t.Fatal(err)
		}
	}

	pertama, errPertama := repo.CabutSesi(t.Context(), []byte("satu"))
	ulang, errUlang := repo.CabutSesi(t.Context(), []byte("satu"))
	errSemua := repo.CabutSemuaSesi(t.Context(), pengguna.ID)
	dua, errDua := repo.CabutSesi(t.Context(), []byte("dua"))
	_, errAsing := repo.CabutSesi(t.Context(), []byte("asing"))

	if errors.Join(errPertama, errUlang, errSemua, errDua) != nil || pertama.IsDicabut || pertama.PenggunaID != pengguna.ID ||
		!ulang.IsDicabut || !dua.IsDicabut || !errors.Is(errAsing, domain.ErrSesiTidakBerlaku) {
		t.Errorf("pertama %+v, ulang %+v, dua %+v, asing %v", pertama, ulang, dua, errAsing)
	}
}

func TestPenggunaRepository_AturPeran_ModeratorDariWarga(t *testing.T) {
	repo := repository.PenggunaRepository{DB: ujidb.BukaTransaksi(t)}
	pengguna, err := repo.PastikanPengguna(t.Context(), emailPenggunaUji)
	if err != nil {
		t.Fatal(err)
	}

	if err := repo.AturPeran(t.Context(), emailPenggunaUji, domain.PeranModerator); err != nil {
		t.Fatal(err)
	}
	diambil, err := repo.AmbilPengguna(t.Context(), pengguna.ID)

	if err != nil || diambil.Peran != domain.PeranModerator {
		t.Errorf("peran %s, err %v", diambil.Peran, err)
	}
}
