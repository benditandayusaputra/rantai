package klien

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	ErrDitolakServer  = errors.New("server menolak permintaan, pengambilan dihentikan")
	ErrDilarangRobots = errors.New("path dilarang robots.txt")
	ErrStatusTidakOK  = errors.New("status respons bukan 200")
	ErrTidakDitemukan = fmt.Errorf("halaman tidak ditemukan: %w", ErrStatusTidakOK)
)

type Konfigurasi struct {
	UserAgent      string        `yaml:"user_agent"`
	Jeda           time.Duration `yaml:"jeda_permintaan"`
	PercobaanUlang int           `yaml:"percobaan_ulang"`
	JedaUlang      time.Duration `yaml:"jeda_ulang"`
}

type Klien struct {
	http     *http.Client
	cfg      Konfigurasi
	kunci    sync.Mutex
	terakhir map[string]time.Time
	robots   map[string]aturanRobots
}

func Baru(klienHTTP *http.Client, cfg Konfigurasi) *Klien {
	return &Klien{
		http:     klienHTTP,
		cfg:      cfg,
		terakhir: map[string]time.Time{},
		robots:   map[string]aturanRobots{},
	}
}

func (k *Klien) Ambil(ctx context.Context, alamat string) ([]byte, error) {
	permintaan, err := http.NewRequestWithContext(ctx, http.MethodGet, alamat, nil)
	if err != nil {
		return nil, err
	}
	return k.kirimSantun(permintaan)
}

func (k *Klien) KirimForm(ctx context.Context, alamat string, form url.Values) ([]byte, error) {
	permintaan, err := http.NewRequestWithContext(ctx, http.MethodPost, alamat, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	permintaan.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return k.kirimSantun(permintaan)
}

func (k *Klien) KirimJSON(ctx context.Context, alamat string, badan []byte) ([]byte, error) {
	permintaan, err := http.NewRequestWithContext(ctx, http.MethodPost, alamat, bytes.NewReader(badan))
	if err != nil {
		return nil, err
	}
	permintaan.Header.Set("Content-Type", "application/json")
	return k.kirimSantun(permintaan)
}

func (k *Klien) kirimSantun(permintaan *http.Request) ([]byte, error) {
	aturan, err := k.aturanRobotsUntuk(permintaan.Context(), permintaan.URL)
	if err != nil {
		return nil, err
	}
	if !aturan.Mengizinkan(permintaan.URL.EscapedPath()) {
		return nil, fmt.Errorf("%s: %w", permintaan.URL, ErrDilarangRobots)
	}
	status, isi, err := k.kirimDenganUlang(permintaan)
	if err != nil {
		return nil, err
	}
	return isi, periksaStatus(permintaan.URL, status)
}

func (k *Klien) kirimDenganUlang(permintaan *http.Request) (int, []byte, error) {
	for percobaan := 1; ; percobaan++ {
		status, isi, err := k.kirim(permintaan)
		if !layakDiulang(permintaan.Context(), status, err) || percobaan > k.cfg.PercobaanUlang {
			return status, isi, err
		}
		if err := tunggu(permintaan.Context(), k.cfg.JedaUlang*time.Duration(percobaan)); err != nil {
			return 0, nil, err
		}
		if permintaan, err = salinPermintaan(permintaan); err != nil {
			return 0, nil, err
		}
	}
}

func layakDiulang(ctx context.Context, status int, err error) bool {
	if err != nil {
		return ctx.Err() == nil
	}
	return status >= http.StatusInternalServerError
}

func salinPermintaan(permintaan *http.Request) (*http.Request, error) {
	salinan := permintaan.Clone(permintaan.Context())
	if permintaan.GetBody == nil {
		return salinan, nil
	}
	badan, err := permintaan.GetBody()
	if err != nil {
		return nil, fmt.Errorf("salin badan permintaan %s: %w", permintaan.URL, err)
	}
	salinan.Body = badan
	return salinan, nil
}

func (k *Klien) kirim(permintaan *http.Request) (int, []byte, error) {
	if err := k.tungguGiliran(permintaan.Context(), permintaan.URL.Host); err != nil {
		return 0, nil, err
	}
	permintaan.Header.Set("User-Agent", k.cfg.UserAgent)
	respons, err := k.http.Do(permintaan)
	if err != nil {
		return 0, nil, fmt.Errorf("kirim %s: %w", permintaan.URL, err)
	}
	defer func() { _ = respons.Body.Close() }()
	isi, err := io.ReadAll(respons.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("baca %s: %w", permintaan.URL, err)
	}
	return respons.StatusCode, isi, nil
}

func periksaStatus(alamat *url.URL, status int) error {
	switch status {
	case http.StatusOK:
		return nil
	case http.StatusForbidden, http.StatusTooManyRequests:
		return fmt.Errorf("%s status %d: %w", alamat, status, ErrDitolakServer)
	case http.StatusNotFound:
		return fmt.Errorf("%s: %w", alamat, ErrTidakDitemukan)
	default:
		return fmt.Errorf("%s status %d: %w", alamat, status, ErrStatusTidakOK)
	}
}

func (k *Klien) tungguGiliran(ctx context.Context, host string) error {
	return tunggu(ctx, time.Until(k.pesanGiliran(host)))
}

func tunggu(ctx context.Context, lama time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(lama):
		return nil
	}
}

func (k *Klien) pesanGiliran(host string) time.Time {
	k.kunci.Lock()
	defer k.kunci.Unlock()
	giliran := time.Now()
	if berikutnya := k.terakhir[host].Add(k.cfg.Jeda); berikutnya.After(giliran) {
		giliran = berikutnya
	}
	k.terakhir[host] = giliran
	return giliran
}

func (k *Klien) aturanRobotsUntuk(ctx context.Context, alamat *url.URL) (aturanRobots, error) {
	k.kunci.Lock()
	aturan, ada := k.robots[alamat.Host]
	k.kunci.Unlock()
	if ada {
		return aturan, nil
	}
	aturan, err := k.ambilRobots(ctx, alamat)
	if err != nil {
		return aturanRobots{}, err
	}
	k.kunci.Lock()
	k.robots[alamat.Host] = aturan
	k.kunci.Unlock()
	return aturan, nil
}

func (k *Klien) ambilRobots(ctx context.Context, alamat *url.URL) (aturanRobots, error) {
	alamatRobots := url.URL{Scheme: alamat.Scheme, Host: alamat.Host, Path: "/robots.txt"}
	permintaan, err := http.NewRequestWithContext(ctx, http.MethodGet, alamatRobots.String(), nil)
	if err != nil {
		return aturanRobots{}, err
	}
	status, isi, err := k.kirim(permintaan)
	switch {
	case err != nil:
		return aturanRobots{}, fmt.Errorf("ambil robots.txt: %w", err)
	case status == http.StatusOK:
		return uraiRobots(string(isi), k.cfg.UserAgent), nil
	case status >= http.StatusInternalServerError:
		return aturanRobots{}, fmt.Errorf("robots.txt %s status %d: %w", alamat.Host, status, ErrDilarangRobots)
	default:
		return aturanRobots{}, nil
	}
}
