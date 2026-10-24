package surel

import (
	"strings"
	"testing"
	"time"
)

func TestSusunPesan_KodeUmurDanHeaderTeks(t *testing.T) {
	s := SMTP{Dari: "noreply@rantai.id", UmurKode: 10 * time.Minute}

	pesan := string(s.susunPesan("warga@contoh.id", "048213"))

	for _, harap := range []string{"To: warga@contoh.id\r\n", "Subject: Kode masuk Rantai\r\n", "Kode masuk Rantai: 048213", "berlaku 10 menit", "charset=UTF-8"} {
		if !strings.Contains(pesan, harap) {
			t.Errorf("pesan tanpa %q:\n%s", harap, pesan)
		}
	}
}
