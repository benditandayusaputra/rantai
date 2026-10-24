package domain

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestLanggananPushValidasi(t *testing.T) {
	sah := LanggananPush{Endpoint: "https://fcm.googleapis.com/fcm/send/abc"}
	sah.Kunci.P256dh, sah.Kunci.Auth = "kunci", "rahasia"
	kasus := map[string]func(*LanggananPush){
		"endpoint http":  func(l *LanggananPush) { l.Endpoint = "http://contoh.test/push" },
		"tanpa auth":     func(l *LanggananPush) { l.Kunci.Auth = "" },
		"endpoint rusak": func(l *LanggananPush) { l.Endpoint = "::" },
	}
	if err := sah.Validasi(); err != nil {
		t.Fatalf("langganan sah ditolak: %v", err)
	}
	for nama, ubah := range kasus {
		l := sah
		ubah(&l)
		if err := l.Validasi(); !errors.Is(err, ErrLanggananPushTidakSah) {
			t.Errorf("%s: err = %v", nama, err)
		}
	}
}

func TestSusunPesanPush_JudulPanjang_Dipotong(t *testing.T) {
	pesan := SusunPesanPush(KejadianUntukPush{Jenis: KejadianPaketBaru, Judul: strings.Repeat("a", 200), Tautan: "/paket/x"})

	if utf8.RuneCountInString(pesan.Isi) != BatasPesanPushKarakter || pesan.Judul != "Paket baru" || pesan.Tautan != "/paket/x" {
		t.Errorf("pesan = %+v", pesan)
	}
}
