package webpush

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

const (
	umurPesanDetik  = 86400
	envKunciPublik  = "VAPID_KUNCI_PUBLIK"
	envKunciPrivat  = "VAPID_KUNCI_PRIVAT"
	envSubjek       = "VAPID_SUBJEK"
	batasKodeSukses = 300
)

type Pengirim struct {
	KunciPublik string
	KunciPrivat string
	Subjek      string
	HTTP        *http.Client
}

func DariEnv(klien *http.Client) (Pengirim, bool) {
	p := Pengirim{KunciPublik: os.Getenv(envKunciPublik), KunciPrivat: os.Getenv(envKunciPrivat), Subjek: os.Getenv(envSubjek), HTTP: klien}
	return p, p.KunciPublik != "" && p.KunciPrivat != "" && p.Subjek != ""
}

func (p Pengirim) Kirim(ctx context.Context, l domain.LanggananPush, pesan domain.PesanPush) error {
	isi, err := json.Marshal(pesan)
	if err != nil {
		return err
	}
	langganan := &webpush.Subscription{Endpoint: l.Endpoint, Keys: webpush.Keys{P256dh: l.Kunci.P256dh, Auth: l.Kunci.Auth}}
	respons, err := webpush.SendNotificationWithContext(ctx, isi, langganan, &webpush.Options{
		HTTPClient: p.HTTP, Subscriber: p.Subjek, VAPIDPublicKey: p.KunciPublik, VAPIDPrivateKey: p.KunciPrivat, TTL: umurPesanDetik,
	})
	if err != nil {
		return fmt.Errorf("kirim push: %w", err)
	}
	defer func() { _ = respons.Body.Close() }()
	switch {
	case respons.StatusCode == http.StatusNotFound || respons.StatusCode == http.StatusGone:
		return usecase.ErrLanggananPushKedaluwarsa
	case respons.StatusCode >= batasKodeSukses:
		return fmt.Errorf("layanan push menjawab %d", respons.StatusCode)
	default:
		return nil
	}
}
