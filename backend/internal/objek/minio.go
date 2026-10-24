package objek

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	tipeJPEG            = "image/jpeg"
	endpointBawaan      = "localhost:9000"
	bucketBawaan        = "rantai-laporan"
	nilaiTLSAktif       = "true"
	envEndpoint         = "MINIO_ENDPOINT"
	envPengguna         = "MINIO_ROOT_USER"
	envSandi            = "MINIO_ROOT_PASSWORD"
	envBucket           = "MINIO_BUCKET"
	envTLS              = "MINIO_TLS"
	penggunaBawaanLokal = "rantai"
	sandiBawaanLokal    = "rantai-rahasia"
)

type Penyimpan struct {
	klien  *minio.Client
	bucket string
}

var ErrKredensialObjekKosong = errors.New(envPengguna + " dan " + envSandi + " wajib diisi di lingkungan")

func BaruDariEnv(ctx context.Context) (Penyimpan, error) {
	pengguna, sandi := os.Getenv(envPengguna), os.Getenv(envSandi)
	if pengguna == "" || sandi == "" {
		return Penyimpan{}, ErrKredensialObjekKosong
	}
	klien, err := minio.New(nilaiEnv(envEndpoint, endpointBawaan), &minio.Options{
		Creds:  credentials.NewStaticV4(pengguna, sandi, ""),
		Secure: os.Getenv(envTLS) == nilaiTLSAktif,
	})
	if err != nil {
		return Penyimpan{}, fmt.Errorf("buat klien penyimpanan objek: %w", err)
	}
	p := Penyimpan{klien: klien, bucket: nilaiEnv(envBucket, bucketBawaan)}
	return p, p.pastikanBucket(ctx)
}

func (p Penyimpan) pastikanBucket(ctx context.Context) error {
	ada, err := p.klien.BucketExists(ctx, p.bucket)
	if err != nil {
		return fmt.Errorf("periksa bucket %s: %w", p.bucket, err)
	}
	if ada {
		return nil
	}
	if err := p.klien.MakeBucket(ctx, p.bucket, minio.MakeBucketOptions{}); err != nil {
		return fmt.Errorf("buat bucket %s: %w", p.bucket, err)
	}
	return nil
}

func (p Penyimpan) SimpanJPEG(ctx context.Context, kunci string, isi []byte) error {
	_, err := p.klien.PutObject(ctx, p.bucket, kunci, bytes.NewReader(isi), int64(len(isi)), minio.PutObjectOptions{ContentType: tipeJPEG})
	if err != nil {
		return fmt.Errorf("simpan objek %s: %w", kunci, err)
	}
	return nil
}

func nilaiEnv(nama, bawaan string) string {
	if nilai := os.Getenv(nama); nilai != "" {
		return nilai
	}
	return bawaan
}
