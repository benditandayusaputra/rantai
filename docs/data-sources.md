# Verifikasi Sumber Data

Catatan hasil Gerbang Data (SPEC bagian 1). Setiap entri memuat URL, tanggal cek, field yang tersedia, dan robots.txt.

## GD-10 Batas wilayah DKI Jakarta

- **Status**: lolos. Cadangan GADM tidak diperlukan.
- **Tanggal cek**: 23 Oktober 2026
- **Sumber**: [cahyadsn/wilayah_boundaries](https://github.com/cahyadsn/wilayah_boundaries), kompilasi batas wilayah sesuai Kepmendagri No 300.2.2-2430 Tahun 2025.
- **Lisensi**: MIT. Pemberitahuan lisensi dipertahankan di berkas fixture.
- **URL berkas** (dasar `https://raw.githubusercontent.com/cahyadsn/wilayah_boundaries/main/db`):
  - `kab/wilayah_boundaries_kab_31.sql`: 6 kota/kabupaten administrasi
  - `kec/wilayah_boundaries_kec_31.sql`: 44 kecamatan
  - `kel/31/wilayah_boundaries_kel_31.{01,71,72,73,74,75}.sql`: 267 kelurahan
- **Format**: pernyataan SQL `INSERT INTO wilayah_boundaries(kode,nama,lat,lng,path)`.
  - `kode`: kode Kemendagri bertitik, contoh `31.72.06.1003`.
  - `nama`: nama tanpa awalan untuk kecamatan dan kelurahan; kota memakai nama resmi lengkap, contoh `Kota Administrasi Jakarta Utara`.
  - `lat`, `lng`: titik acuan; tidak dipakai, titik tengah dihitung ulang dengan `ST_PointOnSurface`.
  - `path`: larik JSON berurutan **[lintang, bujur]**, berupa Polygon (kedalaman 3) atau MultiPolygon (kedalaman 4). Importer menukar urutan ke [bujur, lintang] sesuai GeoJSON.
- **Provinsi**: tidak diambil dari sumber; geometrinya gabungan geometri keenam kota/kabupaten.
- **robots.txt**: `raw.githubusercontent.com/robots.txt` mengembalikan 404, tidak ada pembatasan. Job tetap memberi jeda 1 detik antarberkas dan memakai User-Agent RantaiBot.
- **Perbedaan dengan SPEC**: sumber provenance baru `kemendagri` ditambahkan ke SPEC 2.1.

## GD-01 SPSE LPSE Provinsi DKI Jakarta

- **Status**: lolos sebagian. Daftar tender dan non-tender serta seluruh halaman non-tender terbuka; halaman detail tender tertutup untuk publik.
- **Tanggal cek**: 23 Oktober 2026, dari IP Indonesia.
- **robots.txt** (`https://spse.inaproc.id/robots.txt`): `User-agent: *` dan `Allow: /`, tanpa `Disallow`.
- **Daftar paket** (pola pyproc): buka `https://spse.inaproc.id/jakarta/lelang` atau `/jakarta/nontender`, ambil `authenticityToken` dari skrip halaman (terikat cookie sesi), lalu `POST /jakarta/dt/lelang?tahun=<TA>` atau `/jakarta/dt/pl?tahun=<TA>` dengan form `draw`, `start`, `length`, `authenticityToken`. Respons JSON DataTables; `recordsTotal` selalu 2147483647, jadi paginasi berhenti saat halaman kurang dari `length`. `length=1000` diterima.
  - Kolom tender: 0 kode, 1 nama (bisa memuat `<span class="badge">`), 2 K/L/PD, 3 tahap, 4 HPS ringkas (`660,9 M`, `71,8 Jt`), 5 kualifikasi, 6 metode, 7 metode evaluasi, 8 `<jenis pengadaan> - TA <tahun,...>`, 9 versi SPSE, 10 nilai kontrak (`Rp. 620.772.529.712,03`).
  - Kolom non-tender: 0 kode, 1 nama, 2 K/L/PD, 3 tahap, 4 HPS ringkas, 5 metode, 6 `<jenis> - TA <tahun>`, 7 versi, 8 nilai kontrak atau `Nilai Kontrak belum dibuat`.
- **Pengumuman non-tender** `GET /jakarta/nontender/<kode>/pengumumanpl`: tabel `th.bgwarning` berisi Kode Paket, Nama Paket (lencana `Paket Batal` bila batal), Alasan Pembatalan, Rencana Umum Pengadaan (tabel kode RUP, nama, sumber dana), Tanggal Pembuatan, Tahap Paket Saat Ini, K/L/PD, Satuan Kerja, Jenis Pengadaan, Metode Pengadaan, Tahun Anggaran (`APBD 2026`), Nilai Pagu Paket, Nilai HPS Paket, Jenis Kontrak, Lokasi Pekerjaan (daftar `<li>`, biasanya hanya kota), Peserta Non Tender (`1 peserta`).
- **Detail tender** `GET /jakarta/lelang/<kode>/pengumumanlelang`, `/jadwal`, dan `/jakarta/evaluasi/<kode>/*`: 403 dengan `cf-mitigated: challenge` untuk klien otomatis, dan "Akses Ditolak! Anda tidak diizinkan membuka halaman ini" bahkan dari peramban biasa. Rantai tidak mencoba melewati pembatasan ini.
- **Akibat**: paket tender hanya memuat data kolom daftar: kode, nama, tahap, HPS ringkas (dibulatkan satu desimal oleh sumber), metode, jenis pengadaan, tahun anggaran, nilai kontrak. Satuan kerja, kode RUP, pagu, lokasi pekerjaan, jumlah peserta, dan pemenang tender tidak tersedia. Instansi tender dicatat sebagai K/L/PD (Provinsi DKI Jakarta), sehingga penentuan lokasi tidak menyaring kota instansi.
- **Cadangan yang diperiksa**: ISB LKPP dinonaktifkan per 31 Desember 2025 dan dipindah ke API Gateway Inaproc (`data.inaproc.id/api/v1`, termasuk `tender/pengumuman` dan `tender/tender-ekontrak`). Akses memerlukan akun peran Data Integrator, token yang terikat IP, dan surat permohonan ke LKPP; dicatat di Catatan blokir PLAN.md sebagai jalur untuk detail tender.
- **Fixture**: `backend/internal/ingest/lpse/testdata/`, hanya paket dengan pemenang badan usaha; paket jasa perorangan (PJLP) tidak dijadikan fixture karena memuat nama orang.
