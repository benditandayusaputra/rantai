package llm

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

var (
	ErrModelBelumDiatur = errors.New("model LLM belum diatur")
	ErrDitolakModel     = errors.New("permintaan ditolak model")
	ErrJawabanTerpotong = errors.New("jawaban model terpotong")
)

type Konfigurasi struct {
	Model     string
	MaksToken int64
	Effort    anthropic.OutputConfigEffort
}

type Klien struct {
	api anthropic.Client
	cfg Konfigurasi
}

func Baru(cfg Konfigurasi) (Klien, error) {
	if strings.TrimSpace(cfg.Model) == "" {
		return Klien{}, ErrModelBelumDiatur
	}
	return Klien{api: anthropic.NewClient(), cfg: cfg}, nil
}

func (k Klien) Tulis(ctx context.Context, sistem, isi string) (string, error) {
	respons, err := k.api.Messages.New(ctx, anthropic.MessageNewParams{
		Model:        k.cfg.Model,
		MaxTokens:    k.cfg.MaksToken,
		System:       []anthropic.TextBlockParam{{Text: sistem}},
		OutputConfig: anthropic.OutputConfigParam{Effort: k.cfg.Effort},
		Messages:     []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(isi))},
	})
	if err != nil {
		return "", fmt.Errorf("panggil model %s: %w", k.cfg.Model, err)
	}
	if err := periksaAlasanBerhenti(respons); err != nil {
		return "", err
	}
	return teksJawaban(respons), nil
}

func periksaAlasanBerhenti(respons *anthropic.Message) error {
	switch respons.StopReason {
	case anthropic.StopReasonRefusal:
		return ErrDitolakModel
	case anthropic.StopReasonMaxTokens:
		return ErrJawabanTerpotong
	default:
		return nil
	}
}

func teksJawaban(respons *anthropic.Message) string {
	var teks strings.Builder
	for _, blok := range respons.Content {
		if t, ok := blok.AsAny().(anthropic.TextBlock); ok {
			teks.WriteString(t.Text)
		}
	}
	return strings.TrimSpace(teks.String())
}
