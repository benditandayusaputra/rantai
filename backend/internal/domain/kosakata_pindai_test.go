package domain

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func bacaKosakataTerlarang(t *testing.T) KosakataTerlarang {
	t.Helper()
	isi, err := os.ReadFile("../../config/kosakata_terlarang.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var daftar KosakataTerlarang
	if err := yaml.Unmarshal(isi, &daftar); err != nil {
		t.Fatal(err)
	}
	return daftar
}

func TestKosakataTerlarang_TidakAdaDiPesanBackend(t *testing.T) {
	daftar := bacaKosakataTerlarang(t)
	err := filepath.WalkDir("../..", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		periksaLiteral(t, path, daftar)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestKosakataTerlarang_TidakAdaDiBerkasI18n(t *testing.T) {
	daftar := bacaKosakataTerlarang(t)
	berkas, _ := filepath.Glob("../../../frontend/src/lib/i18n/*")
	for _, path := range berkas {
		isi, err := os.ReadFile(path)
		if err != nil || strings.HasPrefix(filepath.Base(path), ".") {
			continue
		}
		if ditemukan := daftar.Temukan(string(isi)); len(ditemukan) > 0 {
			t.Errorf("%s memuat kosakata terlarang %v", path, ditemukan)
		}
	}
}

func periksaLiteral(t *testing.T, path string, daftar KosakataTerlarang) {
	t.Helper()
	berkas, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		t.Fatalf("urai %s: %v", path, err)
	}
	ast.Inspect(berkas, func(n ast.Node) bool {
		literal, ok := n.(*ast.BasicLit)
		if !ok || literal.Kind != token.STRING {
			return true
		}
		teks, _ := strconv.Unquote(literal.Value)
		if ditemukan := daftar.Temukan(teks); len(ditemukan) > 0 {
			t.Errorf("%s memuat kosakata terlarang %v: %q", path, ditemukan, teks)
		}
		return true
	})
}
