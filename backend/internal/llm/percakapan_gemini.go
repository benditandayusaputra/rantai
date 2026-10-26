package llm

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/benditandayusaputra/rantai/backend/internal/tanya"
)

const (
	kunciHasilAlat = "hasil"
	kunciGalatAlat = "galat"
	argumenKosong  = "{}"
)

type PercakapanGemini struct {
	klien     KlienGemini
	sistem    string
	alat      []alatGemini
	riwayat   []json.RawMessage
	panggilan map[string]panggilanGemini
}

func (k KlienGemini) BukaPercakapan(sistem string, alat []tanya.Spek) tanya.Percakapan {
	return &PercakapanGemini{
		klien: k, sistem: sistem, alat: keAlatGemini(alat), riwayat: []json.RawMessage{}, panggilan: map[string]panggilanGemini{},
	}
}

func (p *PercakapanGemini) Tanya(ctx context.Context, isi string) (tanya.Giliran, error) {
	return p.kirimPengguna(ctx, []bagianGemini{{Text: isi}})
}

func (p *PercakapanGemini) Balas(ctx context.Context, hasil []tanya.HasilAlat) (tanya.Giliran, error) {
	bagian := make([]bagianGemini, len(hasil))
	for i, h := range hasil {
		asal := p.panggilan[h.IDPanggilan]
		bagian[i] = bagianGemini{FunctionResponse: &responsAlatGemini{ID: asal.ID, Name: asal.Name, Response: responsAlat(h)}}
	}
	return p.kirimPengguna(ctx, bagian)
}

func (p *PercakapanGemini) kirimPengguna(ctx context.Context, bagian []bagianGemini) (tanya.Giliran, error) {
	pesan, err := isiPengguna(bagian...)
	if err != nil {
		return tanya.Giliran{}, err
	}
	p.riwayat = append(p.riwayat, pesan)
	kandidat, err := p.klien.kirim(ctx, permintaanGemini{
		SystemInstruction: instruksiSistem(p.sistem), Contents: p.riwayat, Tools: p.alat,
	})
	if err != nil {
		return tanya.Giliran{}, err
	}
	isi, err := uraiIsi(kandidat.Content)
	if err != nil {
		return tanya.Giliran{}, err
	}
	p.riwayat = append(p.riwayat, kandidat.Content)
	return p.giliranDari(isi), nil
}

func (p *PercakapanGemini) giliranDari(isi isiGemini) tanya.Giliran {
	giliran := tanya.Giliran{Teks: teksGemini(isi), Panggilan: []tanya.PanggilanAlat{}}
	for _, bagian := range isi.Parts {
		if bagian.FunctionCall != nil {
			giliran.Panggilan = append(giliran.Panggilan, p.catatPanggilan(*bagian.FunctionCall))
		}
	}
	return giliran
}

func (p *PercakapanGemini) catatPanggilan(panggilan panggilanGemini) tanya.PanggilanAlat {
	id := panggilan.ID
	if id == "" {
		id = awalanIDPanggilanBuatan + strconv.Itoa(len(p.panggilan))
	}
	p.panggilan[id] = panggilan
	masukan := panggilan.Args
	if len(masukan) == 0 {
		masukan = json.RawMessage(argumenKosong)
	}
	return tanya.PanggilanAlat{ID: id, Nama: panggilan.Name, Masukan: masukan}
}

func responsAlat(h tanya.HasilAlat) map[string]any {
	switch {
	case h.IsError:
		return map[string]any{kunciGalatAlat: h.Isi}
	case json.Valid([]byte(h.Isi)):
		return map[string]any{kunciHasilAlat: json.RawMessage(h.Isi)}
	default:
		return map[string]any{kunciHasilAlat: h.Isi}
	}
}
