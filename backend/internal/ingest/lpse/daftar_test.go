package lpse

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

type pengirimTiruan struct {
	halamanDaftar []byte
	jumlahBaris   int
	formDiterima  []url.Values
}

func (p *pengirimTiruan) Ambil(_ context.Context, _ string) ([]byte, error) {
	return p.halamanDaftar, nil
}

func (p *pengirimTiruan) KirimForm(_ context.Context, _ string, form url.Values) ([]byte, error) {
	p.formDiterima = append(p.formDiterima, form)
	mulai, _ := strconv.Atoi(form.Get("start"))
	panjang, _ := strconv.Atoi(form.Get("length"))
	baris := []string{}
	for i := mulai; i < min(mulai+panjang, p.jumlahBaris); i++ {
		baris = append(baris, fmt.Sprintf(`["%d","Paket %d","Provinsi DKI Jakarta","Paket Sudah Selesai","1 Jt","Pengadaan Langsung","Jasa Lainnya - TA 2026","5","Nilai Kontrak belum dibuat"]`, i, i))
	}
	return []byte(`{"data":[` + strings.Join(baris, ",") + `]}`), nil
}

func TestAmbilDaftar_LimaBarisHalamanDua_TigaHalamanDiambil(t *testing.T) {
	pengirim := &pengirimTiruan{halamanDaftar: bacaFixture(t, "halaman_nontender.html"), jumlahBaris: 5}
	sumber := SumberDaftar{Klien: pengirim, URLDasar: "https://spse.uji/jakarta", UkuranHalaman: 2}

	daftar, err := sumber.AmbilDaftar(t.Context(), JenisNonTender, 2026)

	if err != nil || len(daftar) != 5 || len(pengirim.formDiterima) != 3 {
		t.Fatalf("daftar = %d, permintaan = %d, err = %v", len(daftar), len(pengirim.formDiterima), err)
	}
	if pengirim.formDiterima[0].Get("authenticityToken") == "" {
		t.Error("token tidak dikirim")
	}
}

func TestAmbilDaftar_ServerAbaikanPaginasi_BerhentiTanpaDuplikat(t *testing.T) {
	pengirim := &pengirimAbaikanPaginasi{pengirimTiruan{halamanDaftar: bacaFixture(t, "halaman_nontender.html"), jumlahBaris: 2}}
	sumber := SumberDaftar{Klien: pengirim, URLDasar: "https://spse.uji/jakarta", UkuranHalaman: 2}

	daftar, err := sumber.AmbilDaftar(t.Context(), JenisNonTender, 2026)

	if err != nil || len(daftar) != 2 {
		t.Errorf("daftar = %d, err = %v", len(daftar), err)
	}
}

type pengirimAbaikanPaginasi struct {
	pengirimTiruan
}

func (p *pengirimAbaikanPaginasi) KirimForm(ctx context.Context, alamat string, form url.Values) ([]byte, error) {
	form.Set("start", "0")
	return p.pengirimTiruan.KirimForm(ctx, alamat, form)
}
