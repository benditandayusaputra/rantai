# Rantai

Web transparansi publik untuk DKI Jakarta: siapa wakilku di DPRD, ke mana uang kotaku mengalir, dan apakah layanan dasar di wilayahku cukup. Setiap angka bisa ditelusuri ke sumbernya.

## Deskripsi

Warga mengetik alamat, kelurahan, atau kecamatan, lalu membuka halaman kecamatannya. Rantai terdiri dari inti tetap dan tema yang bisa ditambah:

| Modul inti | Isi |
| --- | --- |
| Wilayah | Pencarian, peta, halaman kecamatan, lencana sumber di setiap angka |
| Wakilku | Dapil dan anggota DPRD DKI, profil, jejak kegiatan dari rilis resmi dan media terverifikasi |
| Arus Anggaran | Paket pengadaan LPSE per kecamatan, penyedia, jadwal, masa pemeliharaan |
| Aksi | Tanya Rantai, Paket Aspirasi, Data yang Belum Dibuka |

Tema pertama adalah **Pendidikan**: apakah kursi SMP negeri di sebuah kecamatan cukup untuk lulusan SD-nya, ke mana belanja sekolah mengalir, dan siapa yang berwenang.

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

Prasyarat: Docker dengan Docker Compose, Go 1.26 ke atas, GNU Make.

```bash
cp .env.example .env
make up
make migrate
make impor-wilayah
make impor-tema
make impor-sekolah
make impor-daya-tampung
make hitung-calon-lulusan
make hitung-indikator
```

`make impor-sekolah` mengambil sekitar 4.000 sekolah dengan jeda 1 detik per permintaan, sehingga butuh beberapa jam. Setiap respons disimpan di `data/raw/dapodik/<tanggal>/`, jadi job yang terputus dapat dijalankan ulang pada hari yang sama tanpa mengulang permintaan.

### Mode demo

`make snapshot-demo` menyimpan seluruh data ke `data/snapshot/rantai.sql.gz` di mesin yang datanya sudah lengkap. Di mesin lain, `make up`, `make migrate`, lalu `make seed-demo` memuat snapshot itu dan skenario `data/snapshot/skenario_demo.yaml`. Dengan `RANTAI_MODE=demo` di `.env`, worker menolak job yang mengambil data dari internet.

### Basis data

Basis data utama (lokal dan produksi) adalah PostgreSQL di Neon dengan ekstensi PostGIS. Isi `DATABASE_URL` dengan URL pooler Neon dan `DATABASE_URL_MIGRASI` dengan URL endpoint langsung (tanpa `-pooler`), karena migrasi memakai advisory lock tingkat sesi. Tambahkan `&pool_max_conn_idle_time=30s&pool_max_conn_lifetime=5m&pool_health_check_period=15s` di akhir `DATABASE_URL` agar koneksi yang menganggur selama job mengambil data tidak diputus Neon di tengah jalan. Tanpa kedua variabel itu, semua perintah memakai PostgreSQL dari `make up`.

Tes integrasi selalu memakai basis data lokal terpisah `rantai_uji` agar tidak menyentuh data utama. Siapkan sekali dengan `make up && make siapkan-uji`, lalu `make test`.

## Deploy

API, worker, MinIO, dan Caddy berjalan di satu server dengan `docker-compose.prod.yml`; frontend di Vercel.

```bash
cp .env.example .env
docker compose -f docker-compose.prod.yml --profile alat run --rm migrate up
docker compose -f docker-compose.prod.yml up -d --build
docker compose -f docker-compose.prod.yml run --rm --entrypoint worker api <job>
```

Isi `.env` server dengan `API_DOMAIN` (misalnya `130-162-198-219.sslip.io`, Caddy mengambil sertifikat TLS otomatis), `PROXY_RAHASIA` acak, `AUTH_RAHASIA`, kredensial MinIO, `CORS_ORIGIN` domain Vercel, serta URL Neon. Port 80 dan 443 harus terbuka. Layanan `cadangan` menyimpan `pg_dump` harian ke `cadangan/` dan menyimpan tujuh hari terakhir.

Di Vercel, set root directory ke `frontend` dan isi `PUBLIC_API_URL=https://<domain-vercel>/api/v1`, `API_HULU=https://<API_DOMAIN>`, `PROXY_RAHASIA` yang sama dengan server, dan `PUBLIC_PETA_DASAR_URL`. Frontend meneruskan `/api/v1/*` ke API sehingga cookie sesi tetap satu situs, dan meneruskan IP pengunjung lewat header yang hanya dipercaya Caddy bila `PROXY_RAHASIA` cocok.

## Progres pengerjaan

- [x] Struktur monorepo
- [x] Lingkungan lokal Docker Compose
- [x] Server API dan skema wilayah
- [x] Kerangka tema dan Tema Pendidikan
- [ ] Arus Anggaran
- [ ] Wakilku
- [ ] Frontend
- [ ] Tanya Rantai dan Paket Aspirasi
- [ ] Akun dan hak jawab

## Sumber data

Batas wilayah: [cahyadsn/wilayah_boundaries](https://github.com/cahyadsn/wilayah_boundaries) (MIT), sesuai Kepmendagri No 300.2.2-2430 Tahun 2025. Rincian verifikasi tiap sumber ada di `docs/data-sources.md`.

## Pengambilan data

Pengambil data Rantai memakai User-Agent `RantaiBot/1.0 (+https://github.com/benditandayusaputra/rantai)`, mematuhi robots.txt, dan membatasi diri 1 permintaan per detik per domain. Pertanyaan atau keberatan dapat disampaikan lewat issue di repositori ini.
