package ujidb

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	envURLBasisdataUji = "TEST_DATABASE_URL"
	kunciTesBasisdata  = 7301
)

func BukaTransaksi(t *testing.T) pgx.Tx {
	t.Helper()
	url := os.Getenv(envURLBasisdataUji)
	if url == "" {
		t.Skip(envURLBasisdataUji + " kosong, tes integrasi dilewati")
	}
	koneksi, err := pgx.Connect(t.Context(), url)
	if err != nil {
		t.Fatalf("sambung basis data uji: %v", err)
	}
	tx, err := koneksi.Begin(t.Context())
	if err != nil {
		t.Fatalf("mulai transaksi uji: %v", err)
	}
	if _, err := tx.Exec(t.Context(), "SELECT pg_advisory_xact_lock($1)", kunciTesBasisdata); err != nil {
		t.Fatalf("kunci basis data uji: %v", err)
	}
	t.Cleanup(func() {
		_ = tx.Rollback(context.Background())
		_ = koneksi.Close(context.Background())
	})
	return tx
}
