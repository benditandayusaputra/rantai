# Rantai

Web transparansi publik untuk DKI Jakarta yang menyambungkan layanan, uang, dan tanggung jawab dalam satu wilayah. Setiap angka bisa ditelusuri ke sumbernya.

## Deskripsi

Warga mengetik alamat atau kelurahan, lalu melihat wilayahnya lewat tiga lensa:

| Lensa | Pertanyaan warga |
| --- | --- |
| Layanan | Apakah kursi sekolah di wilayahku cukup untuk lulusan tahun ini? |
| Uang | Proyek apa yang dibiayai di sini, berapa nilainya, sampai mana buktinya? |
| Tanggung Jawab | Instansi mana yang berwenang, komisi DPRD mana yang mengawasi, dan siapa wakil dapilku? |

Karya cabang Web In Action, MIA 2026 (UPNVJ), tema "Designing Digital Futures, Powering Public Transparency".

## Tech stack

| Lapis | Teknologi |
| --- | --- |
| Backend API dan worker | Go, GoFiber v3, clean architecture |
| Database | PostgreSQL 16 + PostGIS 3 |
| Penyimpanan objek | MinIO (kompatibel S3) |
| Frontend | SvelteKit 2, Svelte 5, TypeScript, TailwindCSS |
| Peta | MapLibre GL JS |
| Tanya Rantai | Claude API dengan tool use |
| Pengujian | go test, Vitest, Playwright |
| Orkestrasi lokal | Docker Compose |

## Cara menjalankan

Prasyarat: Docker dengan Docker Compose, Go 1.25 ke atas, GNU Make.

```bash
cp .env.example .env
make up
make migrate
```

## Progres pengerjaan

- [x] Struktur monorepo
- [x] Lingkungan lokal Docker Compose
- [ ] Server API dan skema wilayah
- [ ] Lensa Layanan
- [ ] Lensa Uang
- [ ] Lensa Tanggung Jawab
- [ ] Frontend
- [ ] Laporan warga
- [ ] Tanya Rantai

## Pengambilan data

Pengambil data Rantai memakai User-Agent `RantaiBot/1.0 (+https://github.com/benditandayusaputra/rantai)`, mematuhi robots.txt, dan membatasi diri 1 permintaan per detik per domain. Pertanyaan atau keberatan dapat disampaikan lewat issue di repositori ini.
