package llm

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
)

const (
	kunciUji = "kunci-rahasia-uji"
	modelUji = "gemini-uji"
)

type serverGeminiTiruan struct {
	respons    []string
	status     []int
	permintaan []map[string]any
	header     []http.Header
	path       []string
}

func (s *serverGeminiTiruan) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	badan, _ := io.ReadAll(r.Body)
	var isi map[string]any
	_ = json.Unmarshal(badan, &isi)
	s.permintaan = append(s.permintaan, isi)
	s.header = append(s.header, r.Header.Clone())
	s.path = append(s.path, r.URL.Path)
	urutan := min(len(s.permintaan), len(s.respons)) - 1
	if urutan < len(s.status) && s.status[urutan] != 0 {
		w.WriteHeader(s.status[urutan])
	}
	_, _ = io.WriteString(w, s.respons[urutan])
}

func klienUji(t *testing.T, server *serverGeminiTiruan) KlienGemini {
	t.Helper()
	srv := httptest.NewServer(server)
	t.Cleanup(srv.Close)
	return KlienGemini{
		http: srv.Client(), alamat: srv.URL, kunci: kunciUji, cfg: Konfigurasi{Model: modelUji, MaksToken: 100}, jedaUlang: time.Millisecond,
	}
}

func ambil(t *testing.T, nilai any, jalur ...any) any {
	t.Helper()
	for _, langkah := range jalur {
		switch k := langkah.(type) {
		case string:
			nilai = nilai.(map[string]any)[k]
		case int:
			nilai = nilai.([]any)[k]
		}
	}
	return nilai
}

func TestKlienGemini_Tulis_MengirimSistemDanKunciLaluMengembalikanTeks(t *testing.T) {
	server := &serverGeminiTiruan{respons: []string{
		`{"candidates":[{"content":{"role":"model","parts":[{"text":"ringkas","thought":true},{"text":" Ringkasan netral. "}]},"finishReason":"STOP"}]}`,
	}}
	teks, err := klienUji(t, server).Tulis(context.Background(), "aturan sistem", "isi berita")
	if err != nil {
		t.Fatal(err)
	}
	if teks != "Ringkasan netral." {
		t.Errorf("teks = %q", teks)
	}
	if got := server.header[0].Get("x-goog-api-key"); got != kunciUji {
		t.Errorf("header kunci = %q", got)
	}
	if server.path[0] != "/models/"+modelUji+":generateContent" {
		t.Errorf("path = %q", server.path[0])
	}
	if got := ambil(t, server.permintaan[0], "systemInstruction", "parts", 0, "text"); got != "aturan sistem" {
		t.Errorf("systemInstruction = %v", got)
	}
	if got := ambil(t, server.permintaan[0], "generationConfig", "maxOutputTokens"); got != float64(100) {
		t.Errorf("maxOutputTokens = %v", got)
	}
}

func TestPercakapanGemini_PanggilanAlat_SignatureDikirimBalikDanResponsBerupaObjek(t *testing.T) {
	kasus := []struct {
		nama      string
		idModel   string
		idDiharap string
	}{
		{nama: "model memberi id", idModel: `"id":"abc",`, idDiharap: "abc"},
		{nama: "model tanpa id", idModel: "", idDiharap: awalanIDPanggilanBuatan + "0"},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			server := &serverGeminiTiruan{respons: []string{
				`{"candidates":[{"content":{"role":"model","parts":[{"functionCall":{` + k.idModel + `"name":"cari_wilayah","args":{"teks":"Cilincing"}},"thoughtSignature":"sig-1"}]},"finishReason":"STOP"}]}`,
				`{"candidates":[{"content":{"role":"model","parts":[{"text":"Cilincing ada di Jakarta Utara."}]},"finishReason":"STOP"}]}`,
			}}
			spek := tanya.Spek{Nama: "cari_wilayah", Deskripsi: "Cari wilayah", Wajib: []string{"teks"},
				Parameter: map[string]tanya.Parameter{"teks": {Tipe: tanya.TipeTeks, Deskripsi: "Nama"}}}
			percakapan := klienUji(t, server).BukaPercakapan("sistem", []tanya.Spek{spek})

			giliran, err := percakapan.Tanya(context.Background(), "di mana Cilincing?")
			if err != nil {
				t.Fatal(err)
			}
			if len(giliran.Panggilan) != 1 || giliran.Panggilan[0].ID != k.idDiharap || giliran.Panggilan[0].Nama != "cari_wilayah" ||
				string(giliran.Panggilan[0].Masukan) != `{"teks":"Cilincing"}` {
				t.Fatalf("panggilan = %+v", giliran.Panggilan)
			}
			deklarasi := ambil(t, server.permintaan[0], "tools", 0, "functionDeclarations", 0, "parameters")
			if ambil(t, deklarasi, "type") != tipeObjek || ambil(t, deklarasi, "properties", "teks", "type") != "STRING" {
				t.Errorf("skema parameter = %v", deklarasi)
			}

			giliran, err = percakapan.Balas(context.Background(), []tanya.HasilAlat{{IDPanggilan: k.idDiharap, Isi: `[{"kode":"31.72.04"}]`}})
			if err != nil {
				t.Fatal(err)
			}
			if giliran.Teks != "Cilincing ada di Jakarta Utara." || len(giliran.Panggilan) != 0 {
				t.Errorf("giliran akhir = %+v", giliran)
			}
			isi := ambil(t, server.permintaan[1], "contents").([]any)
			if len(isi) != 3 || ambil(t, isi[1], "parts", 0, "thoughtSignature") != "sig-1" {
				t.Fatalf("riwayat tidak memuat signature: %v", isi)
			}
			respons := ambil(t, isi[2], "parts", 0, "functionResponse").(map[string]any)
			if respons["name"] != "cari_wilayah" || ambil(t, respons, "response", "hasil", 0, "kode") != "31.72.04" {
				t.Errorf("functionResponse = %v", respons)
			}
			if id, ada := respons["id"]; (k.idModel != "") != ada || (ada && id != k.idDiharap) {
				t.Errorf("id functionResponse = %v, ada %v", id, ada)
			}
		})
	}
}

func TestResponsAlat_GalatDibungkusObjek(t *testing.T) {
	got := responsAlat(tanya.HasilAlat{Isi: "kecamatan tidak ditemukan", IsError: true})
	if got[kunciGalatAlat] != "kecamatan tidak ditemukan" {
		t.Errorf("respons = %v", got)
	}
}

func TestKlienGemini_JawabanTidakNormal_MengembalikanGalat(t *testing.T) {
	kasus := []struct {
		nama    string
		status  int
		respons string
		galat   error
		memuat  string
	}{
		{nama: "terpotong", respons: `{"candidates":[{"content":{"parts":[{"text":"sebagian"}]},"finishReason":"MAX_TOKENS"}]}`, galat: ErrJawabanTerpotong},
		{nama: "diblokir keamanan", respons: `{"candidates":[{"finishReason":"SAFETY"}]}`, galat: ErrDitolakModel},
		{nama: "prompt diblokir", respons: `{"promptFeedback":{"blockReason":"PROHIBITED_CONTENT"}}`, galat: ErrDitolakModel},
		{nama: "tanpa kandidat", respons: `{}`, galat: ErrTanpaKandidat},
		{nama: "http gagal", status: http.StatusBadRequest, respons: `{"error":{"message":"API key not valid"}}`, memuat: "API key not valid"},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			server := &serverGeminiTiruan{respons: []string{k.respons}, status: []int{k.status}}
			_, err := klienUji(t, server).BukaPercakapan("sistem", nil).Tanya(context.Background(), "halo")
			if err == nil || (k.galat != nil && !errors.Is(err, k.galat)) || !strings.Contains(err.Error(), k.memuat) {
				t.Fatalf("err = %v, diharap %v / memuat %q", err, k.galat, k.memuat)
			}
			if strings.Contains(err.Error(), kunciUji) {
				t.Errorf("pesan galat membocorkan kunci: %v", err)
			}
		})
	}
}

func TestBaru_MemilihPenyediaDariNamaModel(t *testing.T) {
	kasus := []struct {
		model  string
		galat  error
		gemini bool
		claude bool
	}{
		{model: " ", galat: ErrModelBelumDiatur},
		{model: "gemini-2.5-flash", gemini: true},
		{model: "gemma-4-31b-it", gemini: true},
		{model: "claude-opus-5-5", claude: true},
	}
	for _, k := range kasus {
		t.Run(k.model, func(t *testing.T) {
			model, err := Baru(Konfigurasi{Model: k.model})
			if !errors.Is(err, k.galat) {
				t.Fatalf("err = %v, diharap %v", err, k.galat)
			}
			_, isGemini := model.(KlienGemini)
			_, isClaude := model.(Klien)
			if isGemini != k.gemini || isClaude != k.claude {
				t.Errorf("model %q: gemini %v, claude %v", k.model, isGemini, isClaude)
			}
		})
	}
}

func TestKlienGemini_GalatSementara_DiulangSesuaiSaranServer(t *testing.T) {
	const (
		sukses        = `{"candidates":[{"content":{"parts":[{"text":"oke"}]},"finishReason":"STOP"}]}`
		sibuk         = `{"error":{"message":"high demand"}}`
		kuotaSebentar = `{"error":{"message":"quota","details":[{"retryDelay":"0.001s"}]}}`
		kuotaLama     = `{"error":{"message":"quota","details":[{"retryDelay":"60s"}]}}`
	)
	kasus := []struct {
		nama       string
		respons    []string
		status     []int
		isSukses   bool
		permintaan int
	}{
		{nama: "503 lalu sukses", respons: []string{sibuk, sukses}, status: []int{503}, isSukses: true, permintaan: 2},
		{nama: "429 dengan saran jeda singkat", respons: []string{kuotaSebentar, sukses}, status: []int{429}, isSukses: true, permintaan: 2},
		{nama: "429 dengan saran jeda terlalu lama", respons: []string{kuotaLama, sukses}, status: []int{429}, permintaan: 1},
		{nama: "400 tidak diulang", respons: []string{sibuk, sukses}, status: []int{400}, permintaan: 1},
		{nama: "503 terus menerus", respons: []string{sibuk}, status: []int{503}, permintaan: percobaanMaksGemini},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			server := &serverGeminiTiruan{respons: k.respons, status: k.status}
			teks, err := klienUji(t, server).Tulis(context.Background(), "sistem", "halo")
			if (err == nil) != k.isSukses || (k.isSukses && teks != "oke") {
				t.Fatalf("teks %q, err %v", teks, err)
			}
			if len(server.permintaan) != k.permintaan {
				t.Errorf("jumlah permintaan = %d, diharap %d", len(server.permintaan), k.permintaan)
			}
		})
	}
}
