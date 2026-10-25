package surel

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"mime"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

const (
	batasWaktuSMTP = 15 * time.Second
	ekstensiTLS    = "STARTTLS"
	subjekOTP      = "Kode masuk Rantai"
	templatPesan   = "From: Rantai <%s>\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\nMIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n%s"
	templatIsiOTP = "Kode masuk Rantai: %s\r\n\r\nKode berlaku %d menit. Abaikan email ini bila Anda tidak meminta kode.\r\n"
)

type SMTP struct {
	Host     string
	Port     string
	Pengguna string
	Sandi    string
	Dari     string
	UmurKode time.Duration
}

type PengirimLog struct{}

func (PengirimLog) KirimOTP(_ context.Context, _, kode string) error {
	log.Printf("kode masuk mode demo: %s", kode)
	return nil
}

func (PengirimLog) Kirim(_ context.Context, _, subjek, _ string) error {
	log.Printf("email mode demo: %s", subjek)
	return nil
}

func DariEnv(umurKode time.Duration) (SMTP, bool) {
	host := os.Getenv("SMTP_HOST")
	return SMTP{
		Host: host, Port: os.Getenv("SMTP_PORT"), Pengguna: os.Getenv("SMTP_PENGGUNA"),
		Sandi: os.Getenv("SMTP_SANDI"), Dari: os.Getenv("SMTP_DARI"), UmurKode: umurKode,
	}, host != ""
}

func (s SMTP) KirimOTP(ctx context.Context, email, kode string) error {
	return s.Kirim(ctx, email, subjekOTP, s.isiOTP(kode))
}

func (s SMTP) isiOTP(kode string) string {
	return fmt.Sprintf(templatIsiOTP, kode, int(s.UmurKode.Minutes()))
}

func (s SMTP) Kirim(ctx context.Context, email, subjek, isi string) error {
	koneksi, err := (&net.Dialer{Timeout: batasWaktuSMTP}).DialContext(ctx, "tcp", net.JoinHostPort(s.Host, s.Port))
	if err != nil {
		return fmt.Errorf("sambung SMTP: %w", err)
	}
	defer func() { _ = koneksi.Close() }()
	if err := koneksi.SetDeadline(time.Now().Add(batasWaktuSMTP)); err != nil {
		return fmt.Errorf("atur batas waktu SMTP: %w", err)
	}
	klien, err := smtp.NewClient(koneksi, s.Host)
	if err != nil {
		return fmt.Errorf("buka sesi SMTP: %w", err)
	}
	if err := s.amankan(klien); err != nil {
		return err
	}
	return s.kirim(klien, email, s.susunPesan(email, subjek, isi))
}

func (s SMTP) amankan(klien *smtp.Client) error {
	if ada, _ := klien.Extension(ekstensiTLS); ada {
		if err := klien.StartTLS(&tls.Config{ServerName: s.Host, MinVersion: tls.VersionTLS12}); err != nil {
			return fmt.Errorf("STARTTLS: %w", err)
		}
	}
	if s.Pengguna == "" {
		return nil
	}
	if err := klien.Auth(smtp.PlainAuth("", s.Pengguna, s.Sandi, s.Host)); err != nil {
		return fmt.Errorf("autentikasi SMTP: %w", err)
	}
	return nil
}

func (s SMTP) kirim(klien *smtp.Client, email string, pesan []byte) error {
	if err := klien.Mail(s.Dari); err != nil {
		return fmt.Errorf("SMTP MAIL: %w", err)
	}
	if err := klien.Rcpt(email); err != nil {
		return fmt.Errorf("SMTP RCPT: %w", err)
	}
	penulis, err := klien.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA: %w", err)
	}
	if _, err := penulis.Write(pesan); err != nil {
		return fmt.Errorf("tulis pesan SMTP: %w", err)
	}
	if err := penulis.Close(); err != nil {
		return fmt.Errorf("tutup pesan SMTP: %w", err)
	}
	return klien.Quit()
}

func (s SMTP) susunPesan(email, subjek, isi string) []byte {
	baris := strings.ReplaceAll(strings.ReplaceAll(isi, "\r\n", "\n"), "\n", "\r\n")
	return fmt.Appendf(nil, templatPesan, s.Dari, email, mime.QEncoding.Encode("utf-8", subjek), time.Now().Format(time.RFC1123Z), baris)
}
