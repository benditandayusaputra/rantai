package surel

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"time"
)

const (
	batasWaktuSMTP  = 15 * time.Second
	ekstensiTLS     = "STARTTLS"
	subjekOTP       = "Kode masuk Rantai"
	templatPesanOTP = "From: Rantai <%s>\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\nMIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		"Kode masuk Rantai: %s\r\n\r\nKode berlaku %d menit. Abaikan email ini bila Anda tidak meminta kode.\r\n"
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

func (s SMTP) KirimOTP(ctx context.Context, email, kode string) error {
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
	return s.kirim(klien, email, s.susunPesan(email, kode))
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

func (s SMTP) susunPesan(email, kode string) []byte {
	return fmt.Appendf(nil, templatPesanOTP, s.Dari, email, subjekOTP, time.Now().Format(time.RFC1123Z), kode, int(s.UmurKode.Minutes()))
}
