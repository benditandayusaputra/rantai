package llm

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
)

type PercakapanAlat struct {
	klien  Klien
	sistem string
	alat   []anthropic.ToolUnionParam
	pesan  []anthropic.MessageParam
}

func (k Klien) BukaPercakapan(sistem string, alat []tanya.Spek) tanya.Percakapan {
	return &PercakapanAlat{klien: k, sistem: sistem, alat: keParamAlat(alat), pesan: []anthropic.MessageParam{}}
}

func keParamAlat(daftar []tanya.Spek) []anthropic.ToolUnionParam {
	hasil := make([]anthropic.ToolUnionParam, len(daftar))
	for i, spek := range daftar {
		hasil[i] = anthropic.ToolUnionParam{OfTool: &anthropic.ToolParam{
			Name:        spek.Nama,
			Description: anthropic.String(spek.Deskripsi),
			InputSchema: anthropic.ToolInputSchemaParam{Properties: spek.Parameter, Required: spek.Wajib},
		}}
	}
	return hasil
}

func (p *PercakapanAlat) Tanya(ctx context.Context, isi string) (tanya.Giliran, error) {
	p.pesan = append(p.pesan, anthropic.NewUserMessage(anthropic.NewTextBlock(isi)))
	return p.kirim(ctx)
}

func (p *PercakapanAlat) Balas(ctx context.Context, hasil []tanya.HasilAlat) (tanya.Giliran, error) {
	blok := make([]anthropic.ContentBlockParamUnion, len(hasil))
	for i, h := range hasil {
		blok[i] = anthropic.NewToolResultBlock(h.IDPanggilan, h.Isi, h.IsError)
	}
	p.pesan = append(p.pesan, anthropic.NewUserMessage(blok...))
	return p.kirim(ctx)
}

func (p *PercakapanAlat) kirim(ctx context.Context) (tanya.Giliran, error) {
	respons, err := p.klien.api.Messages.New(ctx, anthropic.MessageNewParams{
		Model:        p.klien.cfg.Model,
		MaxTokens:    p.klien.cfg.MaksToken,
		System:       []anthropic.TextBlockParam{{Text: p.sistem}},
		OutputConfig: anthropic.OutputConfigParam{Effort: p.klien.cfg.Effort},
		Tools:        p.alat,
		Messages:     p.pesan,
	})
	if err != nil {
		return tanya.Giliran{}, fmt.Errorf("panggil model %s: %w", p.klien.cfg.Model, err)
	}
	if err := periksaAlasanBerhenti(respons); err != nil {
		return tanya.Giliran{}, err
	}
	p.pesan = append(p.pesan, respons.ToParam())
	return giliranDari(respons), nil
}

func giliranDari(respons *anthropic.Message) tanya.Giliran {
	giliran := tanya.Giliran{Teks: teksJawaban(respons), Panggilan: []tanya.PanggilanAlat{}}
	for _, blok := range respons.Content {
		if pakai, ok := blok.AsAny().(anthropic.ToolUseBlock); ok {
			giliran.Panggilan = append(giliran.Panggilan, tanya.PanggilanAlat{ID: pakai.ID, Nama: pakai.Name, Masukan: pakai.Input})
		}
	}
	return giliran
}
