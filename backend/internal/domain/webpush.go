package domain

import (
	"net/url"
	"strings"
)

const (
	KanalWebPush           = "webpush"
	BatasPesanPushKarakter = 120
)

var labelKejadian = map[JenisKejadian]string{
	KejadianPaketBaru:            "Paket baru",
	KejadianMelewatiJadwal:       "Melewati jadwal",
	KejadianPemeliharaanBerakhir: "Pemeliharaan segera berakhir",
	KejadianLaporanTerverifikasi: "Laporan terverifikasi",
	KejadianHakJawab:             "Hak jawab baru",
	KejadianIndikatorDiperbarui:  "Data indikator diperbarui",
}

var ErrLanggananPushTidakSah = InputTidakValid("langganan push wajib berisi endpoint https, p256dh, dan auth")

type LanggananPush struct {
	Endpoint string `json:"endpoint"`
	Kunci    struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

type PesanPush struct {
	Judul  string `json:"judul"`
	Isi    string `json:"isi"`
	Tautan string `json:"tautan"`
}

type TujuanPush struct {
	LanggananID int64
	Langganan   LanggananPush
}

type KejadianUntukPush struct {
	ID            int64
	Jenis         JenisKejadian
	KodeKecamatan string
	Judul         string
	Tautan        string
}

func (l LanggananPush) Validasi() error {
	alamat, err := url.Parse(l.Endpoint)
	if err != nil || alamat.Scheme != "https" || alamat.Host == "" || l.Kunci.P256dh == "" || l.Kunci.Auth == "" {
		return ErrLanggananPushTidakSah
	}
	return nil
}

func SusunPesanPush(k KejadianUntukPush) PesanPush {
	isi := []rune(strings.TrimSpace(k.Judul))
	if len(isi) > BatasPesanPushKarakter {
		isi = append(isi[:BatasPesanPushKarakter-1], '…')
	}
	return PesanPush{Judul: labelKejadian[k.Jenis], Isi: string(isi), Tautan: k.Tautan}
}
