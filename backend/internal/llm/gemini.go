package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
)

const (
	alamatGeminiBawaan      = "https://generativelanguage.googleapis.com/v1beta"
	envKunciGemini          = "GEMINI_API_KEY"
	peranPengguna           = "user"
	tipeObjek               = "OBJECT"
	alasanTerpotongGemini   = "MAX_TOKENS"
	batasBacaGalatGemini    = 1 << 16
	awalanIDPanggilanBuatan = "panggilan_"
	percobaanMaksGemini     = 3
	jedaUlangGemini         = 2 * time.Second
	jedaUlangMaksGemini     = 20 * time.Second
)

var (
	ErrTanpaKandidat    = errors.New("model tidak mengembalikan kandidat jawaban")
	alasanDitolakGemini = map[string]bool{
		"SAFETY": true, "RECITATION": true, "BLOCKLIST": true, "PROHIBITED_CONTENT": true, "SPII": true, "IMAGE_SAFETY": true,
	}
	statusBolehUlangGemini = map[int]bool{
		http.StatusTooManyRequests: true, http.StatusInternalServerError: true, http.StatusServiceUnavailable: true,
	}
)

type KlienGemini struct {
	http      *http.Client
	alamat    string
	kunci     string
	cfg       Konfigurasi
	jedaUlang time.Duration
}

type bagianGemini struct {
	Text             string             `json:"text,omitempty"`
	Thought          bool               `json:"thought,omitempty"`
	FunctionCall     *panggilanGemini   `json:"functionCall,omitempty"`
	FunctionResponse *responsAlatGemini `json:"functionResponse,omitempty"`
}

type panggilanGemini struct {
	ID   string          `json:"id,omitempty"`
	Name string          `json:"name"`
	Args json.RawMessage `json:"args,omitempty"`
}

type responsAlatGemini struct {
	ID       string         `json:"id,omitempty"`
	Name     string         `json:"name"`
	Response map[string]any `json:"response"`
}

type isiGemini struct {
	Role  string         `json:"role,omitempty"`
	Parts []bagianGemini `json:"parts"`
}

type skemaGemini struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description,omitempty"`
	Properties  map[string]skemaGemini `json:"properties,omitempty"`
	Required    []string               `json:"required,omitempty"`
}

type deklarasiAlatGemini struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  skemaGemini `json:"parameters"`
}

type alatGemini struct {
	FunctionDeclarations []deklarasiAlatGemini `json:"functionDeclarations"`
}

type konfigurasiGenerasiGemini struct {
	MaxOutputTokens int64 `json:"maxOutputTokens,omitempty"`
}

type permintaanGemini struct {
	SystemInstruction isiGemini                 `json:"systemInstruction"`
	Contents          []json.RawMessage         `json:"contents"`
	Tools             []alatGemini              `json:"tools,omitempty"`
	GenerationConfig  konfigurasiGenerasiGemini `json:"generationConfig"`
}

type kandidatGemini struct {
	Content      json.RawMessage `json:"content"`
	FinishReason string          `json:"finishReason"`
}

type responsGemini struct {
	Candidates     []kandidatGemini `json:"candidates"`
	PromptFeedback struct {
		BlockReason string `json:"blockReason"`
	} `json:"promptFeedback"`
}

type galatGemini struct {
	Error struct {
		Message string `json:"message"`
		Details []struct {
			RetryDelay string `json:"retryDelay"`
		} `json:"details"`
	} `json:"error"`
}

func isModelGemini(nama string) bool {
	return strings.HasPrefix(nama, "gemini") || strings.HasPrefix(nama, "gemma")
}

func baruGemini(cfg Konfigurasi) KlienGemini {
	return KlienGemini{
		http: http.DefaultClient, alamat: alamatGeminiBawaan, kunci: os.Getenv(envKunciGemini), cfg: cfg, jedaUlang: jedaUlangGemini,
	}
}

func (k KlienGemini) Tulis(ctx context.Context, sistem, isi string) (string, error) {
	pesan, err := isiPengguna(bagianGemini{Text: isi})
	if err != nil {
		return "", err
	}
	kandidat, err := k.kirim(ctx, permintaanGemini{SystemInstruction: instruksiSistem(sistem), Contents: []json.RawMessage{pesan}})
	if err != nil {
		return "", err
	}
	konten, err := uraiIsi(kandidat.Content)
	return teksGemini(konten), err
}

func (k KlienGemini) kirim(ctx context.Context, p permintaanGemini) (kandidatGemini, error) {
	p.GenerationConfig.MaxOutputTokens = k.cfg.MaksToken
	badan, err := json.Marshal(p)
	if err != nil {
		return kandidatGemini{}, fmt.Errorf("serialisasi permintaan model %s: %w", k.cfg.Model, err)
	}
	for percobaan := 1; ; percobaan++ {
		kandidat, jeda, err := k.kirimSekali(ctx, badan)
		if jeda == 0 || jeda > jedaUlangMaksGemini || percobaan >= percobaanMaksGemini || tunggu(ctx, jeda) != nil {
			return kandidat, err
		}
	}
}

func (k KlienGemini) kirimSekali(ctx context.Context, badan []byte) (kandidatGemini, time.Duration, error) {
	permintaan, err := k.buatPermintaan(ctx, badan)
	if err != nil {
		return kandidatGemini{}, 0, err
	}
	respons, err := k.http.Do(permintaan)
	if err != nil {
		return kandidatGemini{}, 0, fmt.Errorf("panggil model %s: %w", k.cfg.Model, err)
	}
	defer func() { _ = respons.Body.Close() }()
	if respons.StatusCode/100 != 2 {
		pesan, saranJeda := bacaGalat(respons.Body)
		jeda := k.jedaUlangUntuk(respons.StatusCode, saranJeda)
		return kandidatGemini{}, jeda, fmt.Errorf("panggil model %s: status %d: %s", k.cfg.Model, respons.StatusCode, pesan)
	}
	var hasil responsGemini
	if err := json.NewDecoder(respons.Body).Decode(&hasil); err != nil {
		return kandidatGemini{}, 0, fmt.Errorf("urai jawaban model %s: %w", k.cfg.Model, err)
	}
	kandidat, err := kandidatPertama(hasil)
	return kandidat, 0, err
}

func (k KlienGemini) jedaUlangUntuk(status int, saranServer time.Duration) time.Duration {
	switch {
	case !statusBolehUlangGemini[status]:
		return 0
	case saranServer > 0:
		return saranServer
	default:
		return k.jedaUlang
	}
}

func tunggu(ctx context.Context, jeda time.Duration) error {
	timer := time.NewTimer(jeda)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func (k KlienGemini) buatPermintaan(ctx context.Context, badan []byte) (*http.Request, error) {
	alamat := fmt.Sprintf("%s/models/%s:generateContent", k.alamat, k.cfg.Model)
	permintaan, err := http.NewRequestWithContext(ctx, http.MethodPost, alamat, bytes.NewReader(badan))
	if err != nil {
		return nil, fmt.Errorf("susun permintaan model %s: %w", k.cfg.Model, err)
	}
	permintaan.Header.Set("Content-Type", "application/json")
	permintaan.Header.Set("x-goog-api-key", k.kunci)
	return permintaan, nil
}

func bacaGalat(badan io.Reader) (string, time.Duration) {
	mentah, _ := io.ReadAll(io.LimitReader(badan, batasBacaGalatGemini))
	var galat galatGemini
	if json.Unmarshal(mentah, &galat) != nil || galat.Error.Message == "" {
		return strings.TrimSpace(string(mentah)), 0
	}
	return galat.Error.Message, saranJedaDari(galat)
}

func saranJedaDari(galat galatGemini) time.Duration {
	for _, rinci := range galat.Error.Details {
		if jeda, err := time.ParseDuration(rinci.RetryDelay); err == nil {
			return jeda
		}
	}
	return 0
}

func kandidatPertama(r responsGemini) (kandidatGemini, error) {
	switch {
	case r.PromptFeedback.BlockReason != "":
		return kandidatGemini{}, ErrDitolakModel
	case len(r.Candidates) == 0:
		return kandidatGemini{}, ErrTanpaKandidat
	case alasanDitolakGemini[r.Candidates[0].FinishReason]:
		return kandidatGemini{}, ErrDitolakModel
	case r.Candidates[0].FinishReason == alasanTerpotongGemini:
		return kandidatGemini{}, ErrJawabanTerpotong
	default:
		return r.Candidates[0], nil
	}
}

func instruksiSistem(sistem string) isiGemini {
	return isiGemini{Parts: []bagianGemini{{Text: sistem}}}
}

func isiPengguna(bagian ...bagianGemini) (json.RawMessage, error) {
	mentah, err := json.Marshal(isiGemini{Role: peranPengguna, Parts: bagian})
	if err != nil {
		return nil, fmt.Errorf("serialisasi pesan pengguna: %w", err)
	}
	return mentah, nil
}

func uraiIsi(mentah json.RawMessage) (isiGemini, error) {
	var isi isiGemini
	if len(mentah) == 0 {
		return isi, nil
	}
	if err := json.Unmarshal(mentah, &isi); err != nil {
		return isi, fmt.Errorf("urai isi jawaban model: %w", err)
	}
	return isi, nil
}

func teksGemini(isi isiGemini) string {
	var teks strings.Builder
	for _, bagian := range isi.Parts {
		if !bagian.Thought {
			teks.WriteString(bagian.Text)
		}
	}
	return strings.TrimSpace(teks.String())
}

func keAlatGemini(daftar []tanya.Spek) []alatGemini {
	if len(daftar) == 0 {
		return nil
	}
	deklarasi := make([]deklarasiAlatGemini, len(daftar))
	for i, spek := range daftar {
		deklarasi[i] = deklarasiAlatGemini{Name: spek.Nama, Description: spek.Deskripsi, Parameters: skemaParameter(spek)}
	}
	return []alatGemini{{FunctionDeclarations: deklarasi}}
}

func skemaParameter(spek tanya.Spek) skemaGemini {
	properti := make(map[string]skemaGemini, len(spek.Parameter))
	for nama, p := range spek.Parameter {
		properti[nama] = skemaGemini{Type: strings.ToUpper(p.Tipe), Description: p.Deskripsi}
	}
	return skemaGemini{Type: tipeObjek, Properties: properti, Required: spek.Wajib}
}
