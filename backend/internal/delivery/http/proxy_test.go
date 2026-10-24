package http

import (
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/benditandayusaputra/rantai/backend/internal/config"
	"github.com/gofiber/fiber/v3"
)

const ipKlienAsli = "203.0.113.7"

func ipLewatJaringan(t *testing.T, proxy []string) string {
	t.Helper()
	app := NewApp(config.Server{ProxyTepercaya: proxy}, Layanan{})
	app.Get("/ip-uji", func(c fiber.Ctx) error { return c.SendString(c.IP()) })
	pendengar, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = app.Listener(pendengar, fiber.ListenConfig{DisableStartupMessage: true}) }()
	t.Cleanup(func() { _ = app.Shutdown() })
	permintaan, err := http.NewRequestWithContext(t.Context(), fiber.MethodGet, "http://"+pendengar.Addr().String()+"/ip-uji", nil)
	if err != nil {
		t.Fatal(err)
	}
	permintaan.Header.Set(fiber.HeaderXForwardedFor, ipKlienAsli)
	resp, err := http.DefaultClient.Do(permintaan)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	isi, _ := io.ReadAll(resp.Body)
	return string(isi)
}

func TestProxyTepercaya(t *testing.T) {
	kasus := []struct {
		nama  string
		proxy []string
		harap string
	}{
		{"dari proxy tepercaya memakai X-Forwarded-For", []string{"127.0.0.0/8"}, ipKlienAsli},
		{"dari luar daftar mengabaikan header", []string{"172.16.0.0/12"}, "127.0.0.1"},
		{"tanpa proxy mengabaikan header", nil, "127.0.0.1"},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			if ip := ipLewatJaringan(t, k.proxy); ip != k.harap {
				t.Errorf("ip %s, harap %s", ip, k.harap)
			}
		})
	}
}
