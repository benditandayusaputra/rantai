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
