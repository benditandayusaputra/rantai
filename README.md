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
```

## Progres pengerjaan

- [x] Struktur monorepo
- [x] Lingkungan lokal Docker Compose
- [x] Server API dan skema wilayah
- [ ] Kerangka tema dan Tema Pendidikan
- [ ] Arus Anggaran
- [ ] Wakilku
- [ ] Frontend
- [ ] Tanya Rantai dan Paket Aspirasi
- [ ] Akun dan hak jawab

## Sumber data

Batas wilayah: [cahyadsn/wilayah_boundaries](https://github.com/cahyadsn/wilayah_boundaries) (MIT), sesuai Kepmendagri No 300.2.2-2430 Tahun 2025. Rincian verifikasi tiap sumber ada di `docs/data-sources.md`.

## Pengambilan data

Pengambil data Rantai memakai User-Agent `RantaiBot/1.0 (+https://github.com/benditandayusaputra/rantai)`, mematuhi robots.txt, dan membatasi diri 1 permintaan per detik per domain. Pertanyaan atau keberatan dapat disampaikan lewat issue di repositori ini.
