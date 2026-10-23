package llm

import (
	"context"
	"fmt"
)

type Penulis interface {
	Tulis(ctx context.Context, sistem, isi string) (string, error)
}

type PeringkasUmpan struct {
	Penulis      Penulis
	PromptSistem string
}

func (p PeringkasUmpan) Ringkas(ctx context.Context, judul, cuplikan string) (string, error) {
	return p.Penulis.Tulis(ctx, p.PromptSistem, fmt.Sprintf("<berita>\nJudul: %s\nCuplikan: %s\n</berita>", judul, cuplikan))
}
