# Review revisi UI/UX

Branch: `ux/clearer-citizen-journey` · 27 Oktober 2026.

## Halaman yang bisa direview

- Beranda: http://localhost:5173/id
- Contoh kecamatan: http://localhost:5173/id/kecamatan/31.72.06
- Draf aspirasi dari masalah sekolah: http://localhost:5173/id/aspirasi?w=31.72.06&tema=pendidikan&masalah=kursi_smp_kurang
- Sekolah sekitar wilayah: http://localhost:5173/id/cek-sekolah?w=31.72.06&tema=pendidikan
- Tanya dengan wilayah terpilih: http://localhost:5173/id/tanya?w=31.72.06

## Perubahan utama

1. Beranda menjelaskan manfaat, menempatkan pencarian sebagai pintu masuk utama, dan menampilkan tiga langkah penggunaan. Foto hanya muncul pada layar besar agar pencarian cepat terlihat di ponsel.
2. Cakupan pencarian dijelaskan sebagai nama kelurahan/kecamatan DKI Jakarta. Petunjuk terlihat dan Enter dapat memilih kecocokan tepat; input ambigu meminta pilihan eksplisit.
3. Tur menjadi opsional dan lebih singkat. Contoh wilayah dibedakan dari wilayah pribadi.
4. Navigasi ponsel berisi lima tujuan. Label Tema menjadi Pendidikan dan Arus Anggaran menjadi Proyek publik. Pada ponsel sangat kecil, pilihan bahasa tersedia di footer agar header tetap muat.
5. Ringkasan pendidikan menjelaskan kursi per 100 calon lulusan serta batas interpretasinya. Setiap angka tetap memiliki sumber yang bisa dibuka. Kegiatan wakil rakyat dipindahkan ke bagian yang dapat dibuka setelah ringkasan sekolah.
6. Aspirasi tentang kekurangan kursi bisa dibuka langsung dari data dengan wilayah, jenjang, dan masalah terpilih. Cek Sekolahku dan Tanya membawa konteks wilayah melalui URL.
7. Cek Sekolahku menjelaskan ketika jarak dihitung dari titik tengah wilayah. Memilih wilayah baru memperbarui URL; memakai GPS membuang konteks wilayah lama tanpa memasukkan koordinat pribadi ke URL.
8. Rantai Rupiah menerangkan kelengkapan data. Status “Data belum tersedia” menggantikan “Putus” agar tidak dibaca sebagai kondisi pekerjaan.
9. Draf aspirasi menampilkan status belum dikirim, penerima, langkah penyampaian, dan informasi kanal dari API yang sudah ada. Tombol salin memiliki pesan pemulihan bila clipboard gagal.
10. Teks Indonesia dan Inggris diperbarui bersama. Warna dan identitas visual yang ada dipertahankan.

## Alur review yang disarankan

Buka beranda dengan browser baru, cari “Kelapa Gading”, lalu tekan Enter. Baca arti angka pendidikan, buka sumbernya, dan pilih “Susun aspirasi dari data ini”. Pastikan wilayah serta masalah sudah terisi dan jelas bahwa dokumen belum dikirim. Kembali ke kecamatan, buka Cek Sekolahku dan Tanya dari navigasi.

Coba pada lebar 390 px dan desktop. Untuk layar kecil 320 px, periksa header, kartu angka, dan navigasi bawah. Panduan penggunaan dapat dibuka melalui “Cara memakai Rantai”.

## Menjalankan kembali secara lokal

Database dan MinIO lokal harus sudah aktif dan berisi data proyek. Konfigurasi lokal tetap memakai berkas `.env` yang sudah ada; kredensial tidak disalin ke branch.

Terminal API, dari direktori `rantai/backend`:

```sh
set -a
source ../.env
set +a
export DATABASE_URL='postgres://rantai:rantai@localhost:5433/rantai?sslmode=disable'
export API_PORT=18090
export CORS_ORIGIN='http://localhost:5173,http://127.0.0.1:5173'
export RANTAI_MODE=demo
export SMTP_HOST=''
export PROXY_RAHASIA=rantai-local-review
export MINIO_ENDPOINT=localhost:9000
export MINIO_TLS=false
go run ./cmd/api
```

Terminal frontend, dari direktori `rantai/frontend`:

```sh
PUBLIC_API_URL=/api/v1 API_HULU=http://127.0.0.1:18090 PROXY_RAHASIA=rantai-local-review npm run dev -- --host 127.0.0.1 --port 5173 --strictPort
```

## Batas revisi

Alur Lapor tetap membutuhkan akun, GPS akurat, foto, dan konfirmasi saksi. Prasyaratnya sekarang dijelaskan lebih awal. Pemilihan titik laporan secara manual belum ditambahkan karena perlu dibedakan dari lokasi perangkat pada alur verifikasi.

Tanya menggunakan konfigurasi model yang sudah tersedia. Cache Tanya untuk demo luring belum tersedia dalam snapshot repo; jawaban model tetap memerlukan koneksi ke penyedianya. Audit awal lengkap tersedia di `UI-UX-QA.md`.

## Hasil verifikasi

- `npm run lint`: lulus, termasuk pemeriksaan tipe Svelte tanpa error atau warning.
- `npm test`: 48 file, 230 tes lulus.
- `npm run build`: lulus.
- Playwright alur UI dan aksesibilitas: 55 tes lulus, 1 dilewati karena pemeriksaan beberapa lebar layar hanya perlu dijalankan sekali di proyek desktop.
- Audit axe memeriksa pelanggaran serius/kritis WCAG 2.0/2.1 A/AA; hasil otomatis ini tidak menggantikan uji pengguna atau pemeriksaan aksesibilitas manual menyeluruh.
- Lebar 320, 390, 768, 1024, dan 1440 px diperiksa untuk overflow. Screenshot desktop/ponsel disimpan oleh Playwright di `frontend/test-results/`.
- Audit massal memakai konfigurasi API sementara di `/tmp` dengan kuota permintaan lebih tinggi agar tidak mengenai batas 120 permintaan/menit. Konfigurasi aplikasi di repo tetap memakai batas normal; server review dikembalikan ke konfigurasi normal setelah tes.
