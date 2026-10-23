package dapodik

import (
	"context"
	"os"
	"path/filepath"
)

type Pengunduh interface {
	Ambil(ctx context.Context, alamat string) ([]byte, error)
	KirimJSON(ctx context.Context, alamat string, badan []byte) ([]byte, error)
}

type permintaan struct {
	namaCache string
	url       string
	badanJSON []byte
}

type pengunduhBercache struct {
	pengunduh Pengunduh
	dirCache  string
}

func (p pengunduhBercache) ambil(ctx context.Context, req permintaan) ([]byte, error) {
	pathCache := filepath.Join(p.dirCache, req.namaCache)
	if isi, err := os.ReadFile(pathCache); err == nil {
		return isi, nil
	}
	isi, err := p.unduh(ctx, req)
	if err != nil {
		return nil, err
	}
	return isi, os.WriteFile(pathCache, isi, 0o644)
}

func (p pengunduhBercache) unduh(ctx context.Context, req permintaan) ([]byte, error) {
	if req.badanJSON != nil {
		return p.pengunduh.KirimJSON(ctx, req.url, req.badanJSON)
	}
	return p.pengunduh.Ambil(ctx, req.url)
}
