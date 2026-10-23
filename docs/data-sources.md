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
- **Ruang kode**: kode paket tender dan non-tender tidak unik bersama. Uji nyata TA 2025 menemukan 12 kode yang dipakai paket tender dan non-tender berbeda (contoh 10033820000: konstruksi puskesmas Cipedak dan tenaga ahli pengembang perangkat lunak). Rantai memakai kode `tender-<kode>` dan `nontender-<kode>`, dengan kode asli di kolom `kode_spse`.
- **Uji nyata**: daftar TA 2025 berisi 606 paket tender dan 10.405 non-tender (sebagian besar jasa lainnya PJLP); satu halaman `length=1000` memerlukan beberapa puluh detik.
- **Fixture**: `backend/internal/ingest/lpse/testdata/`, hanya paket dengan pemenang badan usaha; paket jasa perorangan (PJLP) tidak dijadikan fixture karena memuat nama orang.

## GD-07 Dapil dan anggota DPRD DKI Jakarta 2024-2029

- **Status**: lolos dengan cadangan. Daftar dapil per kecamatan dari PKPU; daftar anggota terpilih dari situs resmi DPRD DKI (cadangan yang disebut SPEC), karena situs KPU tidak menyediakan daftar calon terpilih per dapil dalam bentuk yang bisa diambil.
- **Tanggal cek**: 23 Oktober 2026
- **Dapil**: PKPU 6/2023 tentang Daerah Pemilihan dan Alokasi Kursi Pemilu 2024 ([berkas batang tubuh](https://peraturan.go.id/files/peraturan-kpu-no-6-tahun-2023.pdf); lampiran per provinsi tidak ikut di berkas itu). Pembagian kecamatan dicocokkan dengan Wikipedia "DPRD Provinsi DKI Jakarta" dan pemberitaan media; hasilnya dikurasi ke `backend/config/dprd.yaml`.
  - 10 dapil, 106 kursi: 12, 9, 9, 10, 10, 10, 10, 12, 12, 12.
  - Dapil 2 memuat seluruh Kabupaten Kepulauan Seribu (2 kecamatan) ditambah Koja, Cilincing, Kelapa Gading.
  - Job `kpu_dprd` memvalidasi bahwa ke-44 kecamatan DKI masuk tepat satu dapil dan nama kecamatan sama dengan tabel `wilayah`. Basis data juga menjaga keunikan `dapil_wilayah.kode_kecamatan`.
- **Anggota**: `https://dprd-dkijakartaprov.go.id` (WordPress).
  - `/akd/dapil/dki-jakarta-{1..10}/`: nomor urut, nama bergelar, partai, tautan `/profile/{slug}/`. Total 106 baris.
  - `/fraksi/periode-2024-2029/{partai}/`: 9 fraksi, kartu anggota berisi nama (`h3`), jabatan fraksi (`h4`), tautan profil. Total 106. Anggota PPP tergabung di Fraksi NasDem; anggota Perindo di Fraksi Demokrat dan Perindo.
  - `/profile/{slug}/`: tempat dan tanggal lahir, agama, media sosial, pengalaman organisasi. Tidak diambil; Rantai hanya menautkan.
- **Kesalahan di situs DPRD yang ditangani**:
  - Teks "Meliputi Kecamatan" di halaman dapil keliru (Tebet muncul di dapil 7 dan 8; Kebayoran Lama, Pesanggrahan, Cilandak tidak ada; Kepulauan Seribu muncul di dapil 2 dan 3). Karena itu cakupan kecamatan tidak diambil dari situs DPRD.
  - Dua tautan profil di halaman dapil salah tempel (Ida Mahmudah memakai slug Ima Mahdiah; Andika Wisnuadji memakai slug Francine Widjojo). Ejaan nama juga berbeda antarhalaman (contoh "Abdilah" dan "Abdillah").
  - Pencocokan dapil ke fraksi: slug profil bila unik di kedua halaman, lalu nama tanpa gelar bila unik di kedua halaman. Sisanya tidak ditebak: job gagal dan mencetak nama yang tidak cocok.
- **Identitas anggota**: `anggota.id` adalah slug profil dari halaman fraksi (unik di sana), contoh `ima-mahdiah`.
- **LHKPN**: tautan ke halaman pengumuman e-LHKPN KPK (`https://elhkpn.kpk.go.id/portal/user/login#announ`); data LHKPN tidak disalin.
- **robots.txt** `dprd-dkijakartaprov.go.id`: hanya `Disallow: /wp-admin/` untuk semua agen. Job memberi jeda 1 detik antarhalaman (19 halaman) dan memakai User-Agent RantaiBot.

## GD-08 Komisi, alat kelengkapan, dan fraksi DPRD DKI Jakarta

- **Status**: lolos. Keanggotaan komisi tersedia, jadi WAK-01 bisa mengurutkan anggota komisi pengawas lebih dahulu.
- **Tanggal cek**: 23 Oktober 2026
- **Komisi** `/akd/komisi/periode-2024-2029/komisi-{a..e}-.../`: judul "KOMISI X | Bidang ...", paragraf "Meliputi: ..." berisi mitra kerja (dipakai sebagai `komisi.bidang`), lalu kartu anggota (`h3` nama, `h4` jabatan: Koordinator, Ketua, Wakil Ketua, Sekretaris, Anggota). Total kartu 106 = 101 anggota komisi + 5 koordinator.
  - Daftar mitra memuat nama dinas yang mengandung koma ("Dinas Pangan, Pertanian dan Kelautan"). Pemisahan memakai daftar kata awal unit kerja di `config/dprd.yaml`.
- **Alat kelengkapan** `/akd/pimpinan-dprd/`, `/akd/badan-musyawarah/`, `/akd/badan-anggaran/`, `/akd/badan-pembentukan-peraturan-daerah/`, `/akd/badan-kehormatan/`: kartu dengan struktur yang sama. Sekretaris DPRD tercantum sebagai "Sekretaris Bukan Anggota" tanpa tautan profil; dilaporkan job dan tidak disimpan.
- **Fraksi**: lihat GD-07.
- **Aturan komisi**: PP 12/2018 Pasal 46 mewajibkan setiap anggota DPRD, kecuali pimpinan DPRD, menjadi anggota satu komisi. Kelima pimpinan muncul di halaman komisi sebagai Koordinator, sehingga `anggota.komisi_kode` mereka kosong dan perannya disimpan di `anggota_akd`. Anggota yang tercatat di dua komisi juga dikosongkan dan dilaporkan; tidak ada yang ditebak.
- **Kesalahan di situs yang ditangani**:
  - Jabatan Basri Baco di Komisi B tertulis "Koordintor". Karena itu pimpinan dikenali dari halaman Pimpinan DPRD, bukan dari teks jabatan.
  - Kartu Pandapotan Sinaga di halaman fraksi tidak bertaut profil. ID diambil dari slug di halaman dapil karena slug itu unik di sana.
- **Pencocokan kartu ke anggota**: nama tanpa gelar (termasuk `nama_lain`) bila menunjuk tepat satu anggota, lalu slug profil. Slug kosong tidak pernah dipakai.
