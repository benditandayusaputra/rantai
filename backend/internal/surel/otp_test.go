package surel

import (
	"strings"
	"testing"
	"time"
)

func TestSusunPesan_KodeUmurDanHeaderTeks(t *testing.T) {
	s := SMTP{Dari: "noreply@rantai.id", UmurKode: 10 * time.Minute}

	pesan := string(s.susunPesan("warga@contoh.id", subjekOTP, s.isiOTP("048213")))

	for _, harap := range []string{"To: warga@contoh.id\r\n", "Subject: Kode masuk Rantai\r\n", "Kode masuk Rantai: 048213", "berlaku 10 menit", "charset=UTF-8"} {
		if !strings.Contains(pesan, harap) {
			t.Errorf("pesan tanpa %q:\n%s", harap, pesan)
		}
	}
}

func TestSusunPesan_IsiBanyakBaris_BarisCRLF(t *testing.T) {
	pesan := string(SMTP{Dari: "noreply@rantai.id"}.susunPesan("warga@contoh.id", "Ringkasan mingguan Rantai", "Halo,\n\nKelapa Gading\n"))

	if !strings.HasSuffix(pesan, "\r\n\r\nHalo,\r\n\r\nKelapa Gading\r\n") || !strings.Contains(pesan, "Subject: Ringkasan mingguan Rantai\r\n") {
		t.Errorf("pesan:\n%q", pesan)
	}
}
