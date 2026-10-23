package provenance

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/oklog/ulid/v2"
)

const prefiksID = "prv_"

type Catatan struct {
	Sumber  domain.SumberData
	URL     string
	Label   domain.Label
	Catatan string
}

func Baru(catatan Catatan, isi []byte, diambilPada time.Time) domain.Provenance {
	return domain.Provenance{
		ID:          IDBaru(),
		Sumber:      catatan.Sumber,
		URL:         catatan.URL,
		DiambilPada: diambilPada,
		Label:       catatan.Label,
		HashKonten:  HashKonten(isi),
		Catatan:     catatan.Catatan,
	}
}

func IDBaru() string {
	return prefiksID + ulid.Make().String()
}

func HashKonten(isi []byte) string {
	jumlah := sha256.Sum256(isi)
	return hex.EncodeToString(jumlah[:])
}
