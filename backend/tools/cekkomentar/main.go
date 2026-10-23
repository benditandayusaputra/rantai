package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var prefiksDirektif = []string{"//go:", "//nolint:"}

func main() {
	temuan, err := cariKomentar(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, t := range temuan {
		fmt.Println(t)
	}
	if len(temuan) > 0 {
		fmt.Fprintf(os.Stderr, "%d komentar ditemukan, kode tidak boleh berisi komentar\n", len(temuan))
		os.Exit(1)
	}
}

func cariKomentar(akar string) ([]string, error) {
	temuan := []string{}
	berkasSet := token.NewFileSet()
	err := filepath.WalkDir(akar, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		komentarBerkas, err := komentarDalamBerkas(berkasSet, path)
		temuan = append(temuan, komentarBerkas...)
		return err
	})
	return temuan, err
}

func komentarDalamBerkas(berkasSet *token.FileSet, path string) ([]string, error) {
	berkas, err := parser.ParseFile(berkasSet, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("urai %s: %w", path, err)
	}
	temuan := []string{}
	for _, grup := range berkas.Comments {
		for _, komentar := range grup.List {
			if !isDirektif(komentar.Text) {
				temuan = append(temuan, berkasSet.Position(komentar.Pos()).String())
			}
		}
	}
	return temuan, nil
}

func isDirektif(teks string) bool {
	for _, prefiks := range prefiksDirektif {
		if strings.HasPrefix(teks, prefiks) {
			return true
		}
	}
	return false
}
