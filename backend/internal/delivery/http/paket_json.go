package http

import (
	"encoding/json"

	"github.com/benditandayusaputra/rantai/backend/internal/domain"
	"github.com/benditandayusaputra/rantai/backend/internal/usecase"
)

type lokasiPaketJSON struct {
	Tingkat     string     `json:"tingkat"`
	KodeWilayah string     `json:"kode_wilayah"`
	NamaWilayah string     `json:"nama_wilayah"`
	Titik       *titikJSON `json:"titik,omitempty"`
	RT          *int       `json:"rt,omitempty"`
	RW          *int       `json:"rw,omitempty"`
}

type rupJSON struct {
	Kode               string  `json:"kode"`
	Nama               string  `json:"nama"`
	SumberDana         string  `json:"sumber_dana"`
	Pagu               *Angka  `json:"pagu"`
	IsDetailTersedia   bool    `json:"detail_tersedia"`
	PelaksanaanMulai   *string `json:"pelaksanaan_mulai"`
	PelaksanaanSelesai *string `json:"pelaksanaan_selesai"`
	URL                string  `json:"url"`
	ProvenanceID       string  `json:"provenance_id"`
}

type penyediaJSON struct {
	ID   int64  `json:"id"`
	Nama string `json:"nama"`
}

type pemeliharaanJSON struct {
	IsBerlaku         bool    `json:"berlaku"`
	Status            string  `json:"status"`
	AkhirPelaksanaan  *string `json:"akhir_pelaksanaan"`
	AkhirPemeliharaan *string `json:"akhir_pemeliharaan"`
	SisaHari          *Angka  `json:"sisa_hari"`
	Label             string  `json:"label"`
	ProvenanceID      string  `json:"provenance_id"`
}

type mataRantaiJSON struct {
	Kode   string `json:"kode"`
	Status string `json:"status"`
	Alasan string `json:"alasan,omitempty"`
}

type sinyalJSON struct {
	Jenis      string         `json:"jenis"`
	Alasan     string         `json:"alasan"`
	Pembanding map[string]any `json:"pembanding"`
}

type detailPaketJSON struct {
	Kode           string           `json:"kode"`
	KodeSPSE       string           `json:"kode_spse"`
	Nama           string           `json:"nama"`
	Tahap          string           `json:"tahap"`
	JenisPengadaan string           `json:"jenis_pengadaan"`
	Metode         string           `json:"metode"`
	JenisKontrak   string           `json:"jenis_kontrak"`
	IsTender       bool             `json:"is_tender"`
	IsBatal        bool             `json:"is_batal"`
	IsSimulation   bool             `json:"is_simulation"`
	TahunAnggaran  int              `json:"tahun_anggaran"`
	Instansi       string           `json:"instansi"`
	Kategori       string           `json:"kategori"`
	Sektor         string           `json:"sektor"`
	Tema           []string         `json:"tema"`
	Pagu           *Angka           `json:"pagu"`
	HPS            *Angka           `json:"hps"`
	NilaiKontrak   *Angka           `json:"nilai_kontrak"`
	JumlahPeserta  *Angka           `json:"jumlah_peserta"`
	TanggalDibuat  *string          `json:"tanggal_dibuat"`
	TanggalKontrak *string          `json:"tanggal_kontrak"`
	SerahTerima    *string          `json:"tanggal_serah_terima"`
	URLSumber      string           `json:"url_sumber"`
	RUP            *rupJSON         `json:"rup"`
	Penyedia       *penyediaJSON    `json:"penyedia"`
	Lokasi         lokasiPaketJSON  `json:"lokasi"`
	Pemeliharaan   pemeliharaanJSON `json:"pemeliharaan"`
	Rantai         []mataRantaiJSON `json:"rantai"`
	Sinyal         []sinyalJSON     `json:"sinyal"`
	ProvenanceID   string           `json:"provenance_id"`
}

type paketRingkasJSON struct {
	Kode           string          `json:"kode"`
	Nama           string          `json:"nama"`
	Tahap          string          `json:"tahap"`
	JenisPengadaan string          `json:"jenis_pengadaan"`
	Metode         string          `json:"metode"`
	TahunAnggaran  int             `json:"tahun_anggaran"`
	HPS            *Angka          `json:"hps"`
	NilaiKontrak   *Angka          `json:"nilai_kontrak"`
	Sektor         string          `json:"sektor"`
	Tema           []string        `json:"tema"`
	Lokasi         lokasiPaketJSON `json:"lokasi"`
	Penyedia       string          `json:"penyedia"`
	IsSimulation   bool            `json:"is_simulation"`
	ProvenanceID   string          `json:"provenance_id"`
}

type paketTerdekatJSON struct {
	Kode         string    `json:"kode"`
	Nama         string    `json:"nama"`
	Sektor       string    `json:"sektor"`
	IsSimulation bool      `json:"is_simulation"`
	Titik        titikJSON `json:"titik"`
	Jarak        Angka     `json:"jarak"`
}

type fiturJSON struct {
	Type       string          `json:"type"`
	Geometry   json.RawMessage `json:"geometry"`
	Properties map[string]any  `json:"properties"`
}

type koleksiFiturJSON struct {
	Type     string      `json:"type"`
	Features []fiturJSON `json:"features"`
}

func keDetailPaketJSON(d usecase.DetailPaket) detailPaketJSON {
	p := d.Paket
	label := labelPaket(p.IsSimulation)
	j := detailPaketJSON{
		Kode: p.Kode, KodeSPSE: p.KodeSPSE, Nama: p.Nama, Tahap: p.Tahap, JenisPengadaan: p.JenisPengadaan,
		Metode: p.Metode, JenisKontrak: p.JenisKontrak, IsTender: p.IsTender, IsBatal: p.IsBatal, IsSimulation: p.IsSimulation,
		TahunAnggaran: p.TahunAnggaran, Instansi: p.Instansi, Kategori: p.Kategori, Sektor: p.Sektor, Tema: p.TemaKode,
		Pagu:          angkaOpsional(p.Pagu, satuanRupiah, label, p.ProvenanceID),
		HPS:           angkaOpsional(p.HPS, satuanRupiah, label, p.ProvenanceID),
		NilaiKontrak:  angkaOpsional(p.NilaiKontrak, satuanRupiah, label, p.ProvenanceID),
		JumlahPeserta: angkaBulatOpsional(p.JumlahPeserta, satuanPeserta, label, p.ProvenanceID),
		TanggalDibuat: tanggalJSON(p.TanggalDibuat), TanggalKontrak: tanggalJSON(p.TanggalKontrak),
		SerahTerima: tanggalJSON(p.TanggalSerahTerima), URLSumber: d.URLSumber,
		RUP: keRUPJSON(p.RUP, d.URLRUP, label), Penyedia: kePenyediaJSON(p.Penyedia),
		Lokasi: keLokasiPaketJSON(p.Lokasi), Pemeliharaan: kePemeliharaanJSON(d.Pemeliharaan),
		Rantai: keRantaiJSON(d.Rantai), Sinyal: keSinyalJSON(p.Sinyal), ProvenanceID: p.ProvenanceID,
	}
	return j
}

func angkaBulatOpsional(nilai *int, satuan string, label domain.Label, provenanceID string) *Angka {
	if nilai == nil {
		return nil
	}
	pecahan := float64(*nilai)
	return angkaOpsional(&pecahan, satuan, label, provenanceID)
}

func keRUPJSON(r *domain.RUPPaket, url string, label domain.Label) *rupJSON {
	if r == nil {
		return nil
	}
	return &rupJSON{
		Kode: r.Kode, Nama: r.Nama, SumberDana: r.SumberDana, Pagu: angkaOpsional(r.Pagu, satuanRupiah, label, r.ProvenanceID),
		IsDetailTersedia: r.IsDetailTersedia, PelaksanaanMulai: tanggalJSON(r.PelaksanaanMulai),
		PelaksanaanSelesai: tanggalJSON(r.PelaksanaanSelesai), URL: url, ProvenanceID: r.ProvenanceID,
	}
}

func kePenyediaJSON(p *domain.PenyediaPaket) *penyediaJSON {
	if p == nil {
		return nil
	}
	return &penyediaJSON{ID: p.ID, Nama: p.Nama}
}

func keLokasiPaketJSON(l domain.LokasiTersimpan) lokasiPaketJSON {
	j := lokasiPaketJSON{Tingkat: string(l.Tingkat), KodeWilayah: l.KodeWilayah, NamaWilayah: l.NamaWilayah, RT: l.RT, RW: l.RW}
	if l.Titik != nil {
		j.Titik = &titikJSON{Lat: l.Titik.Lat, Lng: l.Titik.Lng}
	}
	return j
}

func kePemeliharaanJSON(h usecase.HasilPemeliharaan) pemeliharaanJSON {
	e := h.Estimasi
	j := pemeliharaanJSON{
		IsBerlaku: h.IsBerlaku, Status: string(e.Status), AkhirPelaksanaan: tanggalJSON(e.AkhirPelaksanaan),
		AkhirPemeliharaan: tanggalJSON(e.AkhirPemeliharaan), Label: string(e.Label), ProvenanceID: domain.ProvenanceRumusPemeliharaan,
	}
	if e.Status == domain.StatusPemeliharaanBerjalan {
		j.SisaHari = &Angka{Nilai: float64(e.SisaHari), Satuan: satuanHari, Label: e.Label, ProvenanceID: domain.ProvenanceRumusPemeliharaan}
	}
	return j
}

func keRantaiJSON(rantai []domain.MataRantai) []mataRantaiJSON {
	hasil := make([]mataRantaiJSON, len(rantai))
	for i, m := range rantai {
		hasil[i] = mataRantaiJSON{Kode: m.Kode, Status: string(m.Status), Alasan: m.Alasan}
	}
	return hasil
}

func keSinyalJSON(daftar []domain.SinyalTersimpan) []sinyalJSON {
	hasil := make([]sinyalJSON, len(daftar))
	for i, s := range daftar {
		hasil[i] = sinyalJSON{Jenis: string(s.Jenis), Alasan: s.Alasan, Pembanding: s.Pembanding}
	}
	return hasil
}

func kePaketRingkasJSON(p domain.PaketRingkas) paketRingkasJSON {
	label := labelPaket(p.IsSimulation)
	return paketRingkasJSON{
		Kode: p.Kode, Nama: p.Nama, Tahap: p.Tahap, JenisPengadaan: p.JenisPengadaan, Metode: p.Metode,
		TahunAnggaran: p.TahunAnggaran, HPS: angkaOpsional(p.HPS, satuanRupiah, label, p.ProvenanceID),
		NilaiKontrak: angkaOpsional(p.NilaiKontrak, satuanRupiah, label, p.ProvenanceID), Sektor: p.Sektor,
		Tema: p.TemaKode, Lokasi: keLokasiPaketJSON(p.Lokasi), Penyedia: p.NamaPenyedia,
		IsSimulation: p.IsSimulation, ProvenanceID: p.ProvenanceID,
	}
}

func kePaketTerdekatJSON(t domain.TitikPaket) paketTerdekatJSON {
	return paketTerdekatJSON{
		Kode: t.Kode, Nama: t.Nama, Sektor: t.Sektor, IsSimulation: t.IsSimulation,
		Titik: titikJSON{Lat: t.Titik.Lat, Lng: t.Titik.Lng},
		Jarak: Angka{Nilai: t.JarakMeter, Satuan: satuanMeter, Label: domain.LabelEstimasi, ProvenanceID: domain.ProvenanceLokasiPaket},
	}
}

func kePetaJSON(p usecase.PetaPaket) koleksiFiturJSON {
	fitur := make([]fiturJSON, 0, len(p.Titik)+len(p.Area))
	for _, t := range p.Titik {
		geometri, _ := json.Marshal(map[string]any{"type": "Point", "coordinates": []float64{t.Titik.Lng, t.Titik.Lat}})
		fitur = append(fitur, fiturJSON{Type: "Feature", Geometry: geometri, Properties: map[string]any{
			"jenis_lokasi": "titik", "kode": t.Kode, "nama": t.Nama, "sektor": t.Sektor, "is_simulation": t.IsSimulation,
		}})
	}
	for _, a := range p.Area {
		fitur = append(fitur, fiturJSON{Type: "Feature", Geometry: json.RawMessage(a.GeometriGeoJSON), Properties: map[string]any{
			"jenis_lokasi": "area", "kode_kecamatan": a.KodeKecamatan, "nama_kecamatan": a.NamaKecamatan, "paket": a.KodePaket,
		}})
	}
	return koleksiFiturJSON{Type: "FeatureCollection", Features: fitur}
}
