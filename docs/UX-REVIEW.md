# Review revisi UI/UX

Branch: `ux/clearer-citizen-journey` · 27 Oktober 2026.

## Halaman yang bisa direview

- Beranda: http://localhost:5173/id
- Contoh kecamatan: http://localhost:5173/id/kecamatan/31.72.06
- Draf aspirasi dari masalah sekolah: http://localhost:5173/id/aspirasi?w=31.72.06&tema=pendidikan&masalah=kursi_smp_kurang
- Sekolah sekitar wilayah: http://localhost:5173/id/cek-sekolah?w=31.72.06&tema=pendidikan
- Tanya dengan wilayah terpilih: http://localhost:5173/id/tanya?w=31.72.06

## Alur terbaru

**Beranda → pilih wilayah → pilih informasi → baca rincian dan langkah berikutnya.**

1. Beranda hanya menjelaskan fungsi Rantai, menampilkan pencarian, dan memberi contoh wilayah. Hero memenuhi satu baris penuh dengan rasio asli gambar; tidak memakai crop atau tinggi tetap. Pada ponsel, teks ditempatkan di bawah foto agar gambar tetap utuh.
2. Beranda tidak lagi mengambil ringkasan kecamatan, peta, atau daftar pengadaan. Data wilayah muncul setelah pemilihan wilayah.
3. Halaman kecamatan menjadi pintu pilihan: **Layanan publik**, **Proyek pemerintah**, dan **Wakil rakyat & instansi**. Tidak ada kategori yang terbuka secara otomatis.
4. Layanan pendidikan menampilkan satu jenjang (SMP atau SMA/SMK), arti angka, dampak potensial, penanggung jawab, pengawas, anggaran terkait, serta pilihan peta tematik/tabel sekolah.
5. Kekurangan yang didukung data mendapat label **Perlu perhatian**. Kondisi memadai dijelaskan manfaatnya tanpa label isu. Data kosong atau tidak cukup mendapat penjelasan **Belum dapat dinilai**.
6. Halaman proyek mendahulukan daftar pekerjaan. Perbandingan grafik belanja dapat dibuka bila diperlukan. Detail proyek membedakan instansi, penyedia, jadwal rencana SiRUP, tahap pengadaan, dan informasi serah terima yang tersedia.
7. Progres fisik, tujuan rinci, dan dampak aktual yang belum tersedia dinyatakan belum tersedia; tidak disimpulkan dari anggaran atau tahap tender. Peta membedakan titik proyek dari lokasi pada tingkat wilayah.
8. Bagian wakil rakyat menjelaskan peran pelaksana dan pengawas sebelum profil anggota DPRD. Profil eksekutif dan direktori kesehatan belum ditambahkan karena datanya belum tersedia.
9. Navigasi utama menjadi Beranda, Wilayah, Tanya, Lapor, dan Tentang Data. Tentang tetap dapat dibuka tanpa memilih wilayah. Wilayah terpilih dibawa melalui URL ketika membuka data, proyek, profil wakil, sekolah, dan Tanya.
10. Aspirasi tetap bisa disusun dari masalah yang didukung data, dengan wilayah dan isu terpilih. Draf memperlihatkan bahwa dokumen belum dikirim dan menjelaskan kanal penyampaiannya.

## Alur review yang disarankan

Buka beranda, tekan “Cari wilayahmu”, cari “Kelapa Gading”, dan tekan Enter. Pilih **Layanan publik**. Baca arti angka, dampak, serta kewenangan; coba pengalih jenjang dan peta/tabel. Buka anggaran terkait, pilih proyek, lalu periksa perbedaan rencana dan progres yang diketahui. Kembali ke pilihan informasi untuk membuka **Wakil rakyat & instansi**.

Halaman langsung:

- Pilihan informasi: http://localhost:5173/id/kecamatan/31.72.06
- Layanan: http://localhost:5173/id/kecamatan/31.72.06?bagian=tema
- Proyek: http://localhost:5173/id/kecamatan/31.72.06?bagian=anggaran
- Wakil dan instansi: http://localhost:5173/id/kecamatan/31.72.06?bagian=wakilku

Coba lebar 390 px dan desktop. Rasio serta lebar hero diperiksa otomatis terhadap dimensi asli gambar. Uji pengguna tetap diperlukan untuk menilai kemudahan memahami isi.

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
- `npm test`: 48 file, 233 tes lulus.
- `npm run build`: lulus.
- Playwright alur UI dan aksesibilitas: 65 tes lulus, 1 dilewati karena pemeriksaan beberapa lebar layar hanya perlu dijalankan sekali di proyek desktop.
- Audit axe memeriksa pelanggaran serius/kritis WCAG 2.0/2.1 A/AA; hasil otomatis ini tidak menggantikan uji pengguna atau pemeriksaan aksesibilitas manual menyeluruh.
- Lebar 320, 390, 768, 1024, dan 1440 px diperiksa untuk overflow. Screenshot desktop/ponsel disimpan oleh Playwright di `frontend/test-results/`.
- Audit massal memakai konfigurasi API sementara di `/tmp` dengan kuota permintaan lebih tinggi agar tidak mengenai batas 120 permintaan/menit. Konfigurasi aplikasi di repo tetap memakai batas normal; server review dikembalikan ke konfigurasi normal setelah tes.
