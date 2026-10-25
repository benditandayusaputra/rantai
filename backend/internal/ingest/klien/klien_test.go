package klien

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

const userAgentUji = "RantaiBot/1.0 (+uji)"

func serverUji(t *testing.T, robots string, status int) (*httptest.Server, *[]time.Time) {
	t.Helper()
	waktu := []time.Time{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			_, _ = w.Write([]byte(robots))
			return
		}
		waktu = append(waktu, time.Now())
		if r.UserAgent() != userAgentUji {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(r.Method + " " + r.FormValue("draw")))
	}))
	t.Cleanup(server.Close)
	return server, &waktu
}

func klienUji(server *httptest.Server, jeda time.Duration) *Klien {
	return Baru(server.Client(), Konfigurasi{UserAgent: userAgentUji, Jeda: jeda})
}

func TestAmbil_TigaPermintaanSatuDomain_BerjedaMinimal(t *testing.T) {
	const jeda = 50 * time.Millisecond
	server, _ := serverUji(t, "", http.StatusOK)
	k := klienUji(server, jeda)
	mulai := time.Now()

	for range 3 {
		if _, err := k.Ambil(t.Context(), server.URL+"/a"); err != nil {
			t.Fatal(err)
		}
	}

	if lama := time.Since(mulai); lama < 3*jeda {
		t.Errorf("lama = %v, harap minimal %v untuk robots.txt dan tiga permintaan", lama, 3*jeda)
	}
}

func TestAmbil_Status403Atau429_ErrDitolakServer(t *testing.T) {
	for _, status := range []int{http.StatusForbidden, http.StatusTooManyRequests} {
		server, _ := serverUji(t, "", status)

		_, err := klienUji(server, 0).Ambil(t.Context(), server.URL+"/a")

		if !errors.Is(err, ErrDitolakServer) {
			t.Errorf("status %d: err = %v", status, err)
		}
	}
}

func TestAmbil_Status404_ErrTidakDitemukanTetapStatusTidakOK(t *testing.T) {
	server, _ := serverUji(t, "", http.StatusNotFound)

	_, err := klienUji(server, 0).Ambil(t.Context(), server.URL+"/a")

	if !errors.Is(err, ErrTidakDitemukan) || !errors.Is(err, ErrStatusTidakOK) || errors.Is(err, ErrDitolakServer) {
		t.Errorf("err = %v", err)
	}
}

func TestAmbil_PathDilarangRobots_TidakMengirimPermintaan(t *testing.T) {
	server, waktu := serverUji(t, "User-agent: *\nDisallow: /rahasia\n", http.StatusOK)

	_, err := klienUji(server, 0).Ambil(t.Context(), server.URL+"/rahasia/1")

	if !errors.Is(err, ErrDilarangRobots) || len(*waktu) != 0 {
		t.Errorf("err = %v, permintaan = %d", err, len(*waktu))
	}
}

func TestKirimForm_Galat5xxSekali_DiulangDenganBadanSama(t *testing.T) {
	gagal := 1
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/robots.txt" && gagal > 0 {
			gagal--
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		_, _ = w.Write([]byte(r.FormValue("draw")))
	}))
	t.Cleanup(server.Close)
	k := Baru(server.Client(), Konfigurasi{UserAgent: userAgentUji, PercobaanUlang: 2})

	isi, err := k.KirimForm(t.Context(), server.URL+"/dt", url.Values{"draw": {"9"}})

	if err != nil || string(isi) != "9" {
		t.Errorf("isi = %q, err = %v", isi, err)
	}
}

func TestAmbil_Status403_TidakDiulang(t *testing.T) {
	server, waktu := serverUji(t, "", http.StatusForbidden)
	k := Baru(server.Client(), Konfigurasi{UserAgent: userAgentUji, PercobaanUlang: 3})

	_, err := k.Ambil(t.Context(), server.URL+"/a")

	if !errors.Is(err, ErrDitolakServer) || len(*waktu) != 1 {
		t.Errorf("err = %v, permintaan = %d", err, len(*waktu))
	}
}

func TestKirimForm_FormTerkirimDenganUserAgent(t *testing.T) {
	server, _ := serverUji(t, "", http.StatusOK)

	isi, err := klienUji(server, 0).KirimForm(t.Context(), server.URL+"/dt", url.Values{"draw": {"7"}})

	if err != nil || string(isi) != "POST 7" {
		t.Errorf("isi = %q, err = %v", isi, err)
	}
}

func TestKirimJSON_BadanDanJenisKontenTerkirim(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		badan, _ := io.ReadAll(r.Body)
		_, _ = w.Write([]byte(r.Method + " " + r.Header.Get("Content-Type") + " " + string(badan)))
	}))
	defer server.Close()

	isi, err := klienUji(server, 0).KirimJSON(t.Context(), server.URL+"/cari", []byte(`{"page":0}`))

	if err != nil || string(isi) != `POST application/json {"page":0}` {
		t.Errorf("isi = %q, err = %v", isi, err)
	}
}

func TestMengizinkan_AturanRobots(t *testing.T) {
	robots := "User-agent: *\nDisallow: /\n\nUser-agent: Googlebot\nAllow: /\n\nUser-agent: RantaiBot\nDisallow: /admin\nAllow: /admin/publik # komentar\n"
	kasus := []struct {
		path      string
		userAgent string
		harap     bool
	}{
		{"/paket", userAgentUji, true},
		{"/admin/x", userAgentUji, false},
		{"/admin/publik/x", userAgentUji, true},
		{"/paket", "BotLain/1.0", false},
	}
	for _, k := range kasus {
		if hasil := uraiRobots(robots, k.userAgent).Mengizinkan(k.path); hasil != k.harap {
			t.Errorf("%s %s = %v, harap %v", k.userAgent, k.path, hasil, k.harap)
		}
	}
}
