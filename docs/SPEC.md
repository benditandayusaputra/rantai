# Spesifikasi Fitur: Rantai

Versi 1.0 · 23 Oktober 2026 · Pasangan dari `docs/PRD.md`. ID fitur di sini sama dengan ID di PRD.

## 0. Konvensi

### 0.1 Label kepastian

Setiap angka yang tampil membawa satu label:

| Label | Arti | Contoh |
| --- | --- | --- |
| `resmi` | Diambil apa adanya dari sumber pemerintah | Pagu dan HPS paket dari LPSE |
| `proksi` | Dihitung Rantai dari data resmi dengan rumus terbuka | Kursi SMP dari rombel kali batas murid |
| `estimasi` | Perkiraan karena data resmi tidak dipublikasikan | Akhir masa pemeliharaan |
| `simulasi` | Data buatan untuk demonstrasi | Laporan warga di paket demo |

### 0.2 Tipe Angka

Setiap nilai numerik di respons API dikirim sebagai objek `Angka`, bukan angka telanjang:

```json
{ "nilai": 7360, "satuan": "kursi", "label": "proksi", "provenance_id": "prv_01J..." }
```

### 0.3 Amplop respons API

```json
{ "code": 200, "message": "ok", "data": {}, "meta": { "page": 1, "per_page": 20, "total": 134, "total_pages": 7 } }
```

Error: `{ "code": 404, "message": "paket tidak ditemukan", "data": null }`. Pesan error berbahasa Indonesia yang bisa dipahami warga; frontend memetakan ke i18n berdasarkan `code` dan kunci error.

### 0.4 Prefiks dan versi

Semua endpoint di bawah `/api/v1`. Kode wilayah memakai kode administrasi Kemendagri (contoh format `31.72.01.1001`).

### 0.5 Path

`config/` di dokumen ini berarti `backend/config/`. `e2e/` berarti `frontend/e2e/`. `data/snapshot/` ada di akar repo.

## 1. Gerbang Data

Tugas bertanda **[GERBANG DATA]** di PLAN.md dimulai dengan verifikasi ini. Hasil verifikasi ditulis ke `docs/data-sources.md` (dibuat saat gerbang pertama) dengan URL, tanggal cek, field yang tersedia, dan isi robots.txt.

| ID | Yang diverifikasi | Bila tidak sesuai |
| --- | --- | --- |
| GD-01 | Halaman paket `spse.inaproc.id/jakarta/lelang/{kode}/pengumumanlelang` dan `nontender/{kode}/pengumumanpl`: field kode paket, kode RUP, satker, pagu, HPS, jenis kontrak, lokasi, jumlah peserta. Halaman daftar `lelang` dan `nontender` beserta paginasinya. robots.txt | Bila daftar dirender JavaScript, cari endpoint JSON yang dipakai halaman; bila tertutup, gunakan Data INAPROC Realisasi sebagai sumber daftar dan SPSE hanya untuk detail |
| GD-02 | Halaman pemenang tender dan pemenang berkontrak: nama penyedia, nilai kontrak, tanggal kontrak, jadwal pelaksanaan | Bila nilai kontrak tidak tersedia, tampilkan HPS saja dan sembunyikan perbandingan kontrak |
| GD-03 | SiRUP: detail paket per kode RUP (nama, satker, pagu, sumber dana, waktu pemilihan) | Bila tidak bisa diakses, Rantai Rupiah menampilkan mata rantai RUP sebagai "kode RUP tercatat, detail tidak tersedia" |
| GD-04 | Dapodik: halaman sekolah per kecamatan dengan rombel dan siswa per tingkat (kelas 6, 7, 9, 10), status negeri/swasta, koordinat | Cadangan 1: Data Induk Peserta Didik per wilayah (unduhan resmi). Cadangan 2: impor CSV manual. Label tetap `proksi` |
| GD-05 | Data madrasah (MI, MTs, MA) per kecamatan DKI | Bila tidak tersedia per sekolah, tampilkan madrasah sebagai agregat kecamatan dengan catatan keterbatasan |
| GD-06 | Lampiran juknis SPMB DKI 2026 berisi daya tampung per sekolah | Bila tidak ditemukan, LAY-05 menampilkan status "daya tampung resmi tidak ditemukan di sumber publik" per sekolah |
| GD-07 | KPU: daftar 10 dapil DPRD DKI 2024 dan kecamatan anggotanya; 106 calon terpilih per dapil | Cadangan: berita acara penetapan KPU DKI dan daftar resmi DPRD DKI |
| GD-08 | Situs DPRD DKI: komisi, bidang tiap komisi, dan keanggotaan komisi periode 2024 sampai 2029 | Bila keanggotaan komisi tidak tersedia, TJ-02 menampilkan anggota dapil tanpa filter komisi dan menyebut komisi pengawas saja |
| GD-09 | Ambang nilai pengadaan langsung untuk Pekerjaan Konstruksi menurut Perpres 16/2018 beserta perubahan terakhirnya (Perpres 12/2021, Perpres 46/2025) | Nilai di `config/sinyal.yaml` diperbarui; sinyal S-03 dimatikan bila ambang tidak bisa dipastikan |
| GD-10 | Batas kelurahan dan kecamatan DKI dalam GeoJSON atau Shapefile dengan kode Kemendagri | Cadangan: GADM tingkat 4 dengan tabel padanan kode |
| GD-11 | Opentender OCDS: unduhan per tahun, field buyer, supplier, nilai, tanggal | Bila unduhan gagal, UANG-07 hanya memakai data LPSE |

## 2. Model data

PostgreSQL 16 dengan PostGIS. Semua tabel punya `dibuat_pada` dan `diubah_pada` (`timestamptz`). Kolom geometri memakai SRID 4326.

### 2.1 Wilayah dan provenance

| Tabel | Kolom penting |
| --- | --- |
| `wilayah` | `kode` PK, `tingkat` (provinsi, kota, kecamatan, kelurahan), `nama`, `induk_kode`, `alias` text[], `geom` MultiPolygon, `titik_tengah` Point |
| `provenance` | `id` PK (ULID), `sumber` (lpse, sirup, dapodik, kpu, dprd, opentender, juknis_spmb, rantai, simulasi), `url`, `diambil_pada`, `label`, `hash_konten`, `catatan` |
| `ingest_run` | `id`, `sumber`, `mulai`, `selesai`, `status`, `jumlah_baru`, `jumlah_berubah`, `pesan_error` |

### 2.2 Lensa Layanan

| Tabel | Kolom penting |
| --- | --- |
| `sekolah` | `npsn` PK, `nama`, `bentuk` (SD, MI, SMP, MTs, SMA, SMK, MA), `jenjang` (dasar, menengah_pertama, menengah_atas), `status` (negeri, swasta), `is_madrasah`, `kode_kelurahan`, `lokasi` Point, `provenance_id` |
| `sekolah_rombel` | `npsn`, `tahun_ajaran`, `tingkat` (1 sampai 12), `jumlah_rombel`, `jumlah_siswa`, `provenance_id`; PK (`npsn`, `tahun_ajaran`, `tingkat`) |
| `daya_tampung_resmi` | `npsn`, `tahun_ajaran`, `kuota`, `provenance_id` |
| `kesenjangan_kecamatan` | materialized view, lihat LAY-01 |

### 2.3 Lensa Uang

| Tabel | Kolom penting |
| --- | --- |
| `instansi` | `id`, `nama` (contoh: Suku Dinas Sumber Daya Air Kota Administrasi Jakarta Utara), `jenis` (dinas, sudin, badan, kelurahan, lainnya), `kode_kota`, `induk_id` |
| `rup` | `kode` PK, `nama`, `instansi_id`, `pagu`, `sumber_dana`, `tahun_anggaran`, `metode`, `provenance_id` |
| `penyedia` | `id`, `nama`, `nama_normal` (huruf kecil, tanpa CV/PT/tanda baca), `alias` text[] |
| `paket` | `kode` PK, `nama`, `kode_rup` FK nullable, `instansi_id`, `jenis_pengadaan`, `metode`, `tahun_anggaran`, `pagu`, `hps`, `nilai_kontrak` nullable, `jenis_kontrak`, `tanggal_dibuat`, `tahap`, `jumlah_peserta`, `penyedia_id` nullable, `jadwal_mulai`, `jadwal_selesai`, `tanggal_serah_terima` nullable, `kategori`, `sektor`, `lokasi_tingkat` (kelurahan, kecamatan, kota), `kode_wilayah`, `lokasi` Point, `rt`, `rw`, `is_simulation`, `provenance_id` |
| `riwayat_penyedia` | `penyedia_id`, `tahun`, `pembeli`, `nama_paket`, `nilai`, `sumber` (opentender, lpse), `provenance_id` |
| `sinyal` | `id`, `paket_kode`, `jenis` (S-01 sampai S-05), `alasan`, `pembanding` jsonb, `aktif`, `dihitung_pada` |
| `laporan` | `id`, `pengguna_id`, `paket_kode`, `kode_wilayah`, `lokasi_publik` Point (dibulatkan ke grid 50 m), `lokasi_asli` Point (privat), `akurasi_meter`, `kondisi` (belum_dimulai, sedang_dikerjakan, terhenti, selesai_baik, selesai_rusak), `catatan`, `status`, `is_simulation` |
| `laporan_foto` | `id`, `laporan_id`, `object_key`, `lebar`, `tinggi` |
| `laporan_konfirmasi` | `laporan_id`, `pengguna_id`, `kondisi`, `lokasi_asli`, `dibuat_pada`; unik (`laporan_id`, `pengguna_id`) |
| `laporan_riwayat` | `laporan_id`, `status_lama`, `status_baru`, `oleh_pengguna_id`, `alasan`, `dibuat_pada` |

### 2.4 Lensa Tanggung Jawab

| Tabel | Kolom penting |
| --- | --- |
| `dapil` | `id`, `nomor`, `nama`, `tahun_pemilu` |
| `dapil_wilayah` | `dapil_id`, `kode_kecamatan` |
| `komisi` | `id`, `kode` (A sampai E, sesuai GD-08), `nama`, `bidang` text[] |
| `anggota` | `id`, `nama`, `fraksi`, `dapil_id`, `komisi_id` nullable, `jabatan`, `lhkpn_url`, `provenance_id` |
| `kewenangan` | dimuat dari `config/kewenangan.yaml`: `jenis_masalah`, `sektor`, `instansi_pelaksana`, `pola_unit_wilayah`, `komisi_kode`, `dasar_hukum`, `dasar_hukum_url` |

### 2.5 Pengguna dan interaksi

| Tabel | Kolom penting |
| --- | --- |
| `pengguna` | `id`, `email` (unik, privat), `nama_tampilan`, `peran` (warga, instansi, moderator, admin), `instansi_id` nullable |
| `otp` | `email`, `kode_hash`, `kedaluwarsa`, `percobaan` |
| `sesi_refresh` | `id`, `pengguna_id`, `token_hash`, `kedaluwarsa`, `dicabut` |
| `hak_jawab` | `id`, `target_jenis` (sinyal, laporan, paket), `target_id`, `instansi_id`, `isi`, `penjawab_id`, `dibuat_pada` |
| `langganan` | `id`, `pengguna_id`, `kode_wilayah`, `kanal` (feed, webpush, email), `webpush_endpoint` jsonb |
| `notifikasi` | `id`, `pengguna_id`, `jenis`, `payload` jsonb, `dibaca` |
| `tanya_log` | `id`, `pertanyaan`, `jawaban`, `rujukan` jsonb, `bahasa`, `durasi_ms` (tanpa data pribadi) |

## 3. Pipeline data

Worker Go (`cmd/worker`) menjalankan job terjadwal. Setiap job: ambil, simpan mentah ke `data/raw/<sumber>/<tanggal>/`, normalisasi, cocokkan, tulis bersama `provenance`, catat `ingest_run`.

| Job | Jadwal | Isi |
| --- | --- | --- |
| `wilayah` | sekali, manual | Impor batas wilayah (GD-10), alias nama (contoh: "Kel. Kelapa Gading Barat", "Klp. Gading Barat") |
| `lpse_daftar` | harian 02.00 WIB | Daftar paket tender dan non-tender tahun berjalan dan sebelumnya |
| `lpse_detail` | harian, setelah daftar | Detail paket baru dan paket yang tahapnya berubah; pemenang (GD-02) |
| `sirup` | mingguan | Detail RUP untuk kode RUP yang ditemukan di paket |
| `opentender` | sekali | Riwayat 2014 sampai 2023 untuk pembeli Pemprov DKI |
| `dapodik` | per semester | Sekolah, rombel, dan siswa per tingkat (GD-04, GD-05) |
| `juknis_spmb` | tahunan, manual | Impor CSV daya tampung resmi hasil ekstraksi lampiran (GD-06) |
| `kpu_dprd` | sekali per periode | Dapil, anggota, komisi (GD-07, GD-08) |
| `turunan` | setelah job sumber | Kategori paket, geocoding, sinyal, refresh materialized view, kejadian notifikasi |

Aturan pengambilan: 1 permintaan per detik per domain, User-Agent `RantaiBot/1.0`, header `If-Modified-Since` bila didukung, berhenti dan catat error bila menerima 403 atau 429, patuhi robots.txt.

## 4. Fondasi

### DASAR-01 Pencarian wilayah

- Input teks: nama kelurahan, kecamatan, kota, atau alamat bebas. Pencocokan memakai `pg_trgm` atas `nama` dan `alias`, hasil diurutkan: kecocokan persis, lalu kemiripan, lalu tingkat (kelurahan dahulu).
- Tombol "Pakai lokasiku": `navigator.geolocation` lalu `GET /wilayah/titik`.
- Alamat bebas: ekstrak nama kelurahan atau kecamatan dari teks dengan gazetir yang sama dengan geocoding paket (UANG-02). Tidak memanggil layanan geocoding eksternal.

API: `GET /wilayah/cari?q=`, `GET /wilayah/titik?lat=&lng=`, `GET /wilayah/{kode}`.

Kriteria terima: "kelapa gading" mengembalikan Kecamatan Kelapa Gading dan tiga kelurahannya; titik di Jakarta Utara mengembalikan kelurahan yang benar; respons di bawah 150 ms.

### DASAR-02 Lencana provenance

Komponen kecil di samping setiap angka. Menampilkan label kepastian sebagai warna dan ikon. Klik membuka panel: sumber, tautan ke halaman asli, waktu diambil, dan rumus bila `proksi` atau `estimasi`.

API: `GET /provenance/{id}`.

### DASAR-03 Wilayahku

Halaman `/{bahasa}/wilayah/{kode}`. Header: nama wilayah, hierarki (kota, kecamatan), tombol Pantau, tombol Bagikan. Tiga tab: Layanan, Uang, Tanggung Jawab, masing-masing dengan URL sendiri (`?lensa=layanan`). Kartu ringkas di atas tab menampilkan satu angka kunci per lensa:

- Layanan: rasio kursi SMP negeri terhadap lulusan SD di kecamatan.
- Uang: jumlah dan total nilai paket tahun berjalan di wilayah.
- Tanggung Jawab: jumlah wakil dapil dan instansi utama.

Bila wilayah adalah kelurahan, Lensa Layanan memakai data kecamatannya dengan keterangan.

API: `GET /wilayah/{kode}/ringkasan`.

### DASAR-04 Peta

Halaman `/{bahasa}/peta`. MapLibre dengan peta dasar vektor bebas lisensi (tile OSM yang di-host sendiri atau penyedia gratis yang mengizinkan). Lapisan yang bisa dinyalakan:

| Lapisan | Bentuk | Sumber API |
| --- | --- | --- |
| Kesenjangan kursi | Choropleth kecamatan, pilihan jenjang SMP atau SMA/SMK | `GET /peta/kesenjangan?jenjang=` |
| Proyek | Titik untuk lokasi kelurahan, area arsir untuk lokasi kecamatan; warna per sektor | `GET /peta/paket?bbox=&kategori=&tahun=` |
| Laporan warga | Titik dengan status | `GET /peta/laporan?bbox=` |
| Sekolah | Titik per bentuk sekolah | `GET /peta/sekolah?bbox=&jenjang=` |

Setiap lapisan punya alternatif tabel yang bisa diakses papan ketik (syarat aksesibilitas). Respons GeoJSON disederhanakan (`ST_SimplifyPreserveTopology`) dan di-cache.

### DASAR-05 Akun dan peran

- Warga masuk dengan OTP 6 digit ke email. Kode berlaku 10 menit, maksimal 5 percobaan, maksimal 3 permintaan per email per jam.
- Token akses JWT 15 menit di memori; token refresh 7 hari di cookie `httpOnly`, `Secure`, `SameSite=Lax`, berotasi setiap pemakaian.
- Peran `instansi` dibuat admin dan terikat ke `instansi_id`. Peran `moderator` dan `admin` dibuat lewat perintah CLI di worker.
- Membaca seluruh data publik tidak memerlukan akun. Akun diperlukan untuk melapor, mengonfirmasi, berlangganan, dan memberi hak jawab.

API: `POST /auth/otp/kirim`, `POST /auth/otp/verifikasi`, `POST /auth/refresh`, `POST /auth/keluar`, `GET /auth/saya`.

### DASAR-06 Bahasa, tema, aksesibilitas

- Rute berawalan bahasa: `/id/...` dan `/en/...`. Nama wilayah, instansi, dan paket tidak diterjemahkan.
- Tema terang dan gelap; bawaan mengikuti sistem; pilihan disimpan.
- Semua elemen interaktif bisa dijangkau papan ketik dengan fokus terlihat; kontras minimal 4,5:1; peta dan grafik punya ringkasan teks dan tabel.

### DASAR-07 Mode demo luring

- `make seed-demo` memuat `data/snapshot/` (dump terkompresi) ke database lokal.
- Variabel `RANTAI_MODE=demo` mematikan job pengambilan data, menampilkan pita "Mode demo, data per <tanggal snapshot>", dan memakai cache jawaban Tanya Rantai untuk pertanyaan demo.
- Tile peta untuk wilayah DKI ikut di-bundle (PMTiles) agar peta tetap tampil tanpa internet.

### DASAR-08 Tentang Data

Halaman statis berisi: daftar sumber dan tanggal ambil terakhir (dari `ingest_run`), rumus Lensa Layanan, aturan geocoding, definisi tiap sinyal, label kepastian, keterbatasan yang diketahui, atribusi dan lisensi ODbL untuk data Opentender, serta tautan unduh turunan data.

## 5. Lensa Layanan

### LAY-01 Perhitungan kursi dan kesenjangan

Konstanta di `config/layanan.yaml`:

| Konstanta | Nilai | Dasar |
| --- | --- | --- |
| `batas_murid_per_rombel.smp` | 32 | Kepmendikdasmen 14/2026 |
| `batas_murid_per_rombel.sma` | 36 | Kepmendikdasmen 14/2026 |
| `batas_murid_per_rombel.smk` | 36 | Kepmendikdasmen 14/2026 |
| `ambang_rasio` | sangat_kurang < 0,5; kurang < 0,8; hampir_cukup < 1,0; cukup ≥ 1,0 | Kebijakan tampilan Rantai |

Rumus per kecamatan `k` dan tahun ajaran `t`, transisi SD ke SMP:

```
calon_lulusan(k, t)        = Σ jumlah_siswa tingkat 6 di SD dan MI dalam k
kursi_proksi(k, t, status) = Σ jumlah_rombel tingkat 7 × batas_murid_per_rombel.smp
                             untuk SMP dan MTs dalam k dengan status tersebut
kesenjangan(k, t, status)  = calon_lulusan(k, t) − kursi_proksi(k, t, status)
rasio(k, t, status)        = kursi_proksi(k, t, status) ÷ calon_lulusan(k, t)
```

Status dihitung terpisah: `negeri`, `swasta`, `madrasah`, dan `semua`. Transisi SMP ke SMA/SMK memakai tingkat 9 dan tingkat 10 dengan batas 36. Rombel tingkat 7 tahun berjalan dipakai sebagai proksi daya tampung tahun berikutnya; asumsi ini tertulis di DASAR-08.

Keterbatasan yang wajib ditampilkan: SPMB Jakarta tidak membatasi pendaftaran pada kecamatan, sehingga angka kecamatan menggambarkan tekanan lokal, bukan peluang pasti seorang anak. LAY-04 melengkapi dengan tampilan radius.

Kasus tepi yang wajib dites: kecamatan tanpa SMP negeri (kursi 0, rasio 0); calon lulusan 0 (rasio tidak didefinisikan, tampil "tidak ada data lulusan"); sekolah tanpa data rombel tingkat 7 (dikeluarkan dan dihitung sebagai "sekolah tanpa data").

Materialized view `kesenjangan_kecamatan` menyimpan hasil per kecamatan, tahun, jenjang, status. Di-refresh oleh job `turunan`.

### LAY-02 Peta kesenjangan

Choropleth kecamatan dengan skala berurutan empat kelas dari `ambang_rasio`. Pilihan: jenjang (SMP, SMA/SMK), status (negeri, semua). Tooltip: nama kecamatan, calon lulusan, kursi, selisih, label proksi. Legenda menjelaskan kelas dalam bahasa awam ("kursi negeri kurang dari separuh lulusan").

### LAY-03 Detail kecamatan

Halaman `/{bahasa}/kecamatan/{kode}/sekolah`. Isi:

- Grafik batang horizontal: calon lulusan dibanding kursi negeri, swasta, madrasah.
- Daftar sekolah jenjang tujuan: nama, bentuk, status, rombel tingkat 7, kursi proksi, siswa tingkat 7 aktual, tautan ke peta.
- Tautan ke Lensa Uang dengan filter sektor pendidikan di kecamatan yang sama.

API: `GET /wilayah/{kode}/layanan`, `GET /kecamatan/{kode}/sekolah?jenjang=`.

### LAY-04 Cek Sekolahku

Input alamat atau lokasi dan jenjang tujuan. Keluaran:

- Sekolah negeri jenjang tujuan dalam radius 3 km (konfigurasi), diurutkan dari terdekat, dengan jarak garis lurus dan kursi proksi.
- `tekanan_radius = calon_lulusan dalam radius ÷ kursi negeri dalam radius`, berlabel proksi, dijelaskan dengan kalimat ("sekitar 2 lulusan untuk setiap 1 kursi negeri di sekitarmu").
- Alternatif: sekolah swasta dan madrasah dalam radius yang sama.

API: `GET /sekolah/terdekat?lat=&lng=&jenjang=&radius_m=`.

### LAY-05 Rapor keterbukaan daya tampung

Per sekolah negeri: daya tampung resmi (bila ditemukan di juknis SPMB), kursi proksi Rantai, selisih, dan status `diumumkan` atau `tidak_ditemukan`. Ringkasan: persentase sekolah yang daya tampungnya ditemukan di sumber publik. Kalimat netral: "tidak ditemukan di sumber publik yang kami periksa", tidak menuduh tidak diumumkan.

API: `GET /keterbukaan/daya-tampung?jenjang=&kode_kota=`.

### LAY-06 Tren

Grafik garis rasio kursi per kecamatan antar tahun ajaran, tampil hanya bila data minimal dua tahun tersedia.

## 6. Lensa Uang

### UANG-01 Impor paket dan taut RUP

- Sumber: GD-01, GD-02, GD-03. Lingkup: semua jenis pengadaan TA 2025 dan 2026; prioritas tampil Pekerjaan Konstruksi.
- Paket ditautkan ke `rup` lewat kode RUP di halaman paket.
- `instansi` dibuat dari nama satuan kerja; kota instansi diambil dari nama ("Kota Administrasi Jakarta Utara") dan dipakai geocoding.
- `penyedia.nama_normal`: huruf kecil, buang "pt", "cv", "tbk", "persero", tanda baca, spasi ganda. Pencocokan penyedia lintas sumber memakai `nama_normal` persis; tidak ada pencocokan fuzzy otomatis.
- Perubahan tahap paket dicatat agar notifikasi "paket baru" dan "tahap berubah" bisa dibuat.

Kriteria terima: importer lulus tes dengan fixture HTML minimal 5 halaman paket (tender, non-tender, tanpa kode RUP, tanpa pemenang, dengan pemenang).

### UANG-02 Geocoding nama paket

Gazetir dibangun dari tabel `wilayah` (nama dan alias). Langkah:

1. Normalisasi nama paket: huruf kecil, ganti "kel." menjadi "kelurahan", "kec." menjadi "kecamatan", "jl." menjadi "jalan", rapikan spasi.
2. Ekstrak RT dan RW dengan pola `rt\s*0*(\d{1,3})` dan `rw\s*0*(\d{1,3})`.
3. Cari semua kemunculan nama kelurahan dan kecamatan dari gazetir, utamakan yang didahului kata "kelurahan" atau "kecamatan".
4. Saring kandidat dengan kota instansi: kandidat di luar kota instansi dibuang (sudin bekerja di kotanya sendiri; dinas provinsi tidak disaring).
5. Tentukan hasil:
   - tepat satu kelurahan cocok dan konsisten dengan kecamatan yang disebut: `lokasi_tingkat = kelurahan`, titik di `titik_tengah` kelurahan;
   - tidak ada kelurahan tetapi tepat satu kecamatan: `lokasi_tingkat = kecamatan`;
   - selain itu: `lokasi_tingkat = kota` dari instansi, atau `provinsi` untuk dinas provinsi.
6. Titik beberapa paket di kelurahan yang sama diberi pergeseran deterministik kecil (berbasis hash kode paket) agar tidak menumpuk, dan ditandai "lokasi kelurahan, bukan titik pasti".

Target: minimal 70% paket Pekerjaan Konstruksi mendapat tingkat kecamatan atau lebih rinci. Laporan cakupan geocoding tampil di Tentang Data.

Tes wajib: contoh nyata "Pembangunan Pintu Air Gg. Nusa Indah RT 04 RW 22 Gading Sengon (JPO Kali Betik), Kel. Kelapa Gading Barat, Kec. Kelapa Gading" menghasilkan kelurahan Kelapa Gading Barat, RT 4, RW 22; nama yang menyebut dua kelurahan berbeda menghasilkan tingkat kecamatan atau kota; nama tanpa wilayah menghasilkan tingkat kota.

### UANG-03 Rantai Rupiah

Halaman `/{bahasa}/paket/{kode}`. Inti visual: rantai horizontal (vertikal di layar sempit) dengan mata rantai:

| Mata rantai | Isi | Status bila data tidak ada |
| --- | --- | --- |
| Rencana (RUP) | Kode RUP, pagu, sumber dana, tautan SiRUP | Putus: "rencana tidak ditemukan" |
| Pemilihan | Metode, HPS, jumlah peserta, tautan SPSE | Selalu ada (sumber paket) |
| Pemenang | Penyedia, nilai kontrak | Menunggu: "belum ada pemenang" |
| Pelaksanaan | Jadwal mulai dan selesai | Estimasi atau tidak tersedia |
| Lokasi | Kelurahan atau kecamatan, RT/RW, peta mini | Tingkat kota |
| Serah terima | Tanggal serah terima (bila diisi instansi) | Estimasi |
| Masa pemeliharaan | Lihat UANG-06 | Estimasi |
| Bukti warga | Jumlah laporan per status, galeri | Kosong: ajakan melapor |

Bagian bawah: sinyal aktif beserta alasannya, hak jawab instansi, perbandingan dengan paket sejenis (UANG-04), riwayat penyedia ringkas, tombol Lapor Bukti, tombol Paket Aspirasi.

API: `GET /paket/{kode}`.

### UANG-04 Kategori dan perbandingan paket sejenis

Kategori dan sektor ditentukan dari nama paket dengan aturan regex di `config/kategori.yaml`, dievaluasi berurutan, kecocokan pertama menang. Contoh isi:

```yaml
- kategori: ruang_kelas
  sektor: pendidikan
  pola: ["(rehab|pembangunan|penambahan).*(ruang kelas|rkb|gedung sekolah)", "\\b(sdn|smpn|sman|smkn)\\b.*(rehab|pembangunan)"]
- kategori: pintu_air
  sektor: air
  pola: ["pintu air"]
- kategori: saluran
  sektor: air
  pola: ["saluran", "drainase", "box culvert", "u-ditch"]
- kategori: tangki_septik
  sektor: sanitasi
  pola: ["tangki septik", "septic tank", "ipal komunal"]
- kategori: jalan
  sektor: jalan
  pola: ["(perbaikan|peningkatan|pemeliharaan|pembangunan) jalan", "trotoar", "jalur pedestrian"]
```

Perbandingan: untuk satu paket, ambil paket lain dengan kategori, jenis pengadaan, dan tahun anggaran yang sama di seluruh DKI. Tampilkan strip plot nilai HPS dengan median, kuartil, dan posisi paket ini. Kalimat pendamping wajib: "Ukuran pekerjaan tiap paket bisa berbeda; perbandingan ini menunjukkan posisi nilai, bukan kewajaran." Perbandingan tidak ditampilkan bila jumlah pembanding kurang dari `min_pembanding` (bawaan 8).

API: `GET /paket/{kode}/pembanding`.

### UANG-05 Sinyal pola

Semua ambang di `config/sinyal.yaml`. Sinyal dihitung ulang oleh job `turunan`. Tampilan selalu: judul netral, alasan dengan angka, data pembanding, tautan sumber, dan ruang hak jawab.

| ID | Judul tampilan | Aturan | Ambang bawaan |
| --- | --- | --- | --- |
| S-01 | Hanya satu peserta | Metode tender atau seleksi, jumlah peserta = 1 | Tidak berlaku untuk pengadaan langsung dan penunjukan langsung |
| S-02 | Nilai jauh di atas paket sejenis | HPS > Q3 + 1,5 × IQR kelompok pembanding UANG-04 | Minimal 8 pembanding |
| S-03 | Beberapa paket sejenis bernilai mepet di bawah batas pengadaan langsung | ≥ 3 paket pengadaan langsung, instansi sama, kategori sama, dalam 30 hari, masing-masing ≥ 85% ambang pengadaan langsung | Ambang dari GD-09 |
| S-04 | Satu penyedia memenangkan sebagian besar paket sejenis | Penyedia yang sama memenangkan ≥ 40% nilai paket kategori tertentu di satu instansi dalam satu tahun anggaran, minimal 5 paket | |
| S-05 | Melewati jadwal | Hari ini > `jadwal_selesai` + 30 hari dan ada ≥ 1 laporan terverifikasi berkondisi belum_dimulai, sedang_dikerjakan, atau terhenti | |

Kosakata terlarang di seluruh antarmuka dan respons API: korupsi, curang, mark up, penggelembungan, kongkalikong, fiktif, penyelewengan. Tes otomatis memindai berkas i18n untuk kata-kata ini.

### UANG-06 Estimasi masa pemeliharaan

Konfigurasi per kategori di `config/pemeliharaan.yaml`: `permanen` (180 hari) atau `semi_permanen` (90 hari), sesuai ketentuan masa pemeliharaan minimal pekerjaan konstruksi.

```
akhir_pelaksanaan = tanggal_serah_terima bila diisi instansi
                    selain itu jadwal_selesai bila ada
                    selain itu tanggal_kontrak + durasi_bawaan_kategori
akhir_pemeliharaan = akhir_pelaksanaan + masa_pemeliharaan_kategori
```

Status: `pelaksanaan`, `pemeliharaan` (dengan hitung mundur hari), `pemeliharaan_berakhir`, `tidak_dapat_diestimasi`. Label `estimasi` kecuali tanggal serah terima diisi instansi lewat Panel Instansi, lalu label menjadi `resmi` dengan sumber instansi tersebut.

Kalimat warga: "Bila ada kerusakan sebelum <tanggal>, perbaikannya masih menjadi kewajiban penyedia (perkiraan)."

### UANG-07 Profil penyedia

Halaman `/{bahasa}/penyedia/{id}`: nama, jumlah dan total nilai paket per tahun (grafik batang), kategori utama, instansi pemberi paket terbanyak, daftar paket di Rantai, riwayat Opentender 2014 sampai 2023 dengan atribusi ODbL. Tanpa skor atau peringkat.

### UANG-08 Lapor Bukti

Alur seluler (layar sempit adalah target utama):

1. Warga masuk (OTP) bila belum.
2. Izin lokasi; akurasi harus ≤ 100 m, bila tidak ditolak dengan saran pindah ke ruang terbuka.
3. Daftar proyek dalam radius 500 m dari lokasi, diurutkan jarak dan kesamaan sektor; warga memilih satu, atau "proyek tidak ada di daftar" (laporan tanpa paket, tetap tersimpan per kelurahan).
4. Ambil atau unggah 1 sampai 4 foto (JPEG, PNG, HEIC; maksimal 8 MB per foto).
5. Pilih kondisi dan isi catatan opsional (maksimal 500 karakter).
6. Ringkasan dan kirim.

Pemrosesan server: verifikasi magic bytes, dekode lalu enkode ulang ke JPEG (otomatis membuang EXIF), batasi sisi terpanjang 2048 px, buat thumbnail 480 px, simpan ke MinIO dengan kunci acak. `lokasi_asli` disimpan privat; `lokasi_publik` dibulatkan ke grid 50 m. Batas: 5 laporan per akun per hari.

API: `POST /laporan` (multipart), `GET /laporan/{id}`, `GET /paket/terdekat?lat=&lng=&radius_m=`.

### UANG-09 Verifikasi laporan

Mesin status di `domain`:

```
DIKIRIM → MENUNGGU_SAKSI → TERVERIFIKASI → DITERUSKAN → DITANGGAPI → SELESAI
                └────────────→ DITOLAK (spam, duplikat, di luar lingkup)
```

- `DIKIRIM` ke `MENUNGGU_SAKSI` otomatis setelah lolos pemrosesan foto.
- `TERVERIFIKASI` bila salah satu terpenuhi:
  - minimal 3 akun berbeda (termasuk pelapor) mengonfirmasi kondisi yang sama, lokasi tiap konfirmasi dalam 200 m dari lokasi laporan, dalam 14 hari;
  - moderator atau petugas instansi terkait menandai terverifikasi.
- `DITERUSKAN` otomatis saat terverifikasi: laporan masuk antrean Panel Instansi milik instansi pelaksana paket, atau instansi dari peta kewenangan untuk laporan tanpa paket.
- `DITANGGAPI` saat instansi memberi hak jawab.
- `SELESAI` bila instansi menandai selesai dan tidak ada sanggahan terverifikasi dalam 14 hari.
- Setiap perpindahan dicatat di `laporan_riwayat` dan tampil sebagai linimasa publik.

Konfirmasi: `POST /laporan/{id}/konfirmasi` dengan lokasi dan kondisi. Pelapor tidak bisa mengonfirmasi laporannya sendiri dua kali.

Tes wajib: semua transisi sah, transisi tidak sah ditolak dengan `ErrTransisiTidakSah`, konfirmasi di luar 200 m tidak dihitung, konfirmasi lewat 14 hari tidak dihitung.

### UANG-10 Unduh data

`GET /unduh/paket.csv?kode_wilayah=&tahun=`: kolom paket utama dan kolom provenance (URL sumber, tanggal ambil). Keterangan lisensi dan atribusi dicantumkan di halaman unduh dan Tentang Data, bukan di dalam CSV.

## 7. Lensa Tanggung Jawab

### TJ-01 Peta kewenangan

`config/kewenangan.yaml` dikurasi tim, setiap entri wajib punya dasar hukum dan tautan. Contoh bentuk (isi instansi dan komisi diverifikasi di GD-08 dan saat kurasi):

```yaml
- masalah: kursi_smp_kurang
  sektor: pendidikan
  label_id: "Kursi SMP negeri kurang"
  label_en: "Not enough public junior high seats"
  instansi_pelaksana: "Dinas Pendidikan Provinsi DKI Jakarta"
  unit_wilayah: "Suku Dinas Pendidikan Wilayah {kota}"
  komisi: E
  dasar_hukum: "UU 23/2014 lampiran I.A pembagian urusan pendidikan; UU 2/2024 tentang Provinsi Daerah Khusus Jakarta"
  dasar_hukum_url: "https://peraturan.bpk.go.id/"
- masalah: jalan_rusak_kota
  sektor: jalan
  instansi_pelaksana: "Dinas Bina Marga Provinsi DKI Jakarta"
  unit_wilayah: "Suku Dinas Bina Marga {kota}"
  komisi: D
- masalah: saluran_tersumbat
  sektor: air
  instansi_pelaksana: "Dinas Sumber Daya Air Provinsi DKI Jakarta"
  unit_wilayah: "Suku Dinas Sumber Daya Air {kota}"
  komisi: D
```

Catatan: status jalan (nasional, provinsi, lingkungan) menentukan kewenangan. Bila status jalan tidak diketahui, tampilkan kemungkinan kewenangan beserta cara mengeceknya.

API: `GET /kewenangan?masalah=`, `GET /kewenangan/daftar`.

### TJ-02 Wakilku

Dari kecamatan wilayah ke dapil (`dapil_wilayah`) ke anggota. Bila masalah dipilih, anggota yang duduk di komisi pengawas masalah itu ditampilkan lebih dahulu dengan penanda "duduk di Komisi X yang mengawasi urusan ini". Urutan selanjutnya alfabetis. Tanpa angka kinerja.

API: `GET /wilayah/{kode}/tanggung-jawab?masalah=`.

### TJ-03 Profil anggota

Nama, fraksi, dapil dan daftar kecamatannya, komisi dan bidangnya, jabatan alat kelengkapan, tautan ke e-Announcement LHKPN (halaman pencarian KPK), tautan ke halaman resmi DPRD. Tanpa foto bila lisensi foto tidak jelas; gunakan inisial.

### TJ-04 Jalur aspirasi terdekat

Untuk kelurahan: Lembaga Musyawarah Kelurahan (penjelasan peran, anggotanya dipilih di tingkat RW, cara menghubungi lewat kantor kelurahan), Dewan Kota (satu wakil per kecamatan), kanal CRM dan JAKI untuk aduan layanan, serta penjelasan masa reses DPRD. Isi statis dari konfigurasi, diverifikasi saat kurasi.

### TJ-05 Paket Aspirasi

Warga memilih masalah dan wilayah. Rantai menyusun dokumen berisi:

1. Judul masalah dan wilayah.
2. Data pendukung dari Lensa Layanan dan Lensa Uang, setiap angka dengan sumber dan tanggal.
3. Instansi berwenang dan dasar hukumnya.
4. Komisi DPRD dan anggota dapil terkait.
5. Pertanyaan yang disarankan, netral dan spesifik (contoh: "Berapa rencana penambahan rombel SMP negeri di Kecamatan X untuk tahun ajaran 2027/2028?").
6. Ruang kosong untuk nama dan tanda tangan warga.

Keluaran: teks siap salin dan PDF (dibuat di server dengan pustaka PDF Go, bukan headless browser). Rantai tidak mengirim dokumen atas nama warga.

API: `POST /aspirasi` dengan `{kode_wilayah, masalah, bahasa}`, respons `{teks, pdf_url}`.

## 8. Lintas lensa

### LINTAS-01 Tanya Rantai

Antarmuka obrolan di `/{bahasa}/tanya` dan panel geser di Wilayahku (wilayah aktif ikut sebagai konteks).

Arsitektur: backend memanggil API LLM dengan tool use. Model dari env `TANYA_MODEL`, kunci dari env `TANYA_API_KEY`. Model hanya boleh menjawab dari hasil fungsi berikut:

| Fungsi | Masukan | Keluaran |
| --- | --- | --- |
| `cari_wilayah` | teks | daftar wilayah |
| `ringkasan_wilayah` | kode | angka kunci tiga lensa |
| `kesenjangan_kursi` | kode kecamatan, jenjang | calon lulusan, kursi per status, rasio |
| `sekolah_di_kecamatan` | kode, jenjang | daftar sekolah dan kursi |
| `daftar_paket` | kode wilayah, sektor, tahun | paket ringkas |
| `detail_paket` | kode | rantai paket |
| `statistik_kategori` | kategori, tahun | median, kuartil, jumlah |
| `kewenangan_masalah` | masalah | instansi, komisi, dasar hukum |
| `wakil_dapil` | kode kecamatan, masalah opsional | anggota dan komisi |
| `laporan_wilayah` | kode | jumlah laporan per status |

Setiap hasil fungsi menyertakan `provenance_id`. Aturan prompt sistem:

- Jawab dalam bahasa pertanyaan (Indonesia atau Inggris), kalimat pendek, istilah awam.
- Setiap angka wajib bersumber dari hasil fungsi dan ditandai rujukan `[n]`.
- Bila data tidak ada, katakan tidak tersedia dan sarankan halaman terkait.
- Tidak memberi penilaian, opini, atau rekomendasi memilih terhadap politisi dan partai; jawab dengan fakta jabatan dan jalur aspirasi.
- Sinyal dijelaskan sebagai hal yang layak ditanyakan, dengan kosakata netral.

Validator pasca-jawaban di server: setiap angka di teks jawaban harus muncul di hasil fungsi pada percakapan itu; bila gagal, jawaban diganti pesan aman dan dicatat. Kosakata terlarang UANG-05 juga diperiksa.

Streaming ke klien via SSE. Batas: 20 pertanyaan per IP per jam, panjang pertanyaan 500 karakter. Mode demo memakai cache jawaban untuk pertanyaan di `data/snapshot/tanya_demo.json`.

API: `POST /tanya` (SSE) dengan `{pertanyaan, kode_wilayah_konteks, bahasa}`.

### LINTAS-02 Hak jawab dan Panel Instansi

Panel di `/{bahasa}/panel` untuk peran instansi:

- Antrean: sinyal aktif pada paket milik instansinya, laporan berstatus DITERUSKAN, paket yang mendekati akhir masa pemeliharaan.
- Aksi: tulis hak jawab (maksimal 2.000 karakter, bisa melampirkan tautan dokumen), tandai laporan selesai, isi tanggal serah terima paket.
- Hak jawab tampil publik di samping sinyal atau laporan dengan nama instansi dan tanggal; tidak bisa dihapus, hanya ditambah tanggapan baru.

API: `GET /panel/antrean`, `POST /hak-jawab`, `PATCH /panel/paket/{kode}/serah-terima`, `PATCH /panel/laporan/{id}/selesai`.

### LINTAS-03 Pantau Wilayahku

Warga berlangganan kelurahan atau kecamatan. Kanal: feed di aplikasi (wajib), Web Push dengan VAPID (P1), email ringkasan mingguan (P2). Kejadian:

| Kejadian | Pemicu |
| --- | --- |
| Paket baru | Paket baru di wilayah |
| Melewati jadwal | Sinyal S-05 baru |
| Masa pemeliharaan segera berakhir | 30 hari sebelum akhir pemeliharaan |
| Laporan terverifikasi | Laporan di wilayah berstatus TERVERIFIKASI |
| Hak jawab baru | Hak jawab untuk paket atau laporan di wilayah |
| Data kursi diperbarui | Refresh tahun ajaran baru |

API: `POST /langganan`, `DELETE /langganan/{id}`, `GET /notifikasi`, `PATCH /notifikasi/{id}/dibaca`.

### LINTAS-04 Panel moderasi

Antrean laporan MENUNGGU_SAKSI dan laporan yang diadukan; aksi tandai terverifikasi, tolak dengan alasan, sembunyikan foto. Semua aksi tercatat di `laporan_riwayat`.

API: `GET /moderasi/laporan`, `PATCH /moderasi/laporan/{id}`.

### LINTAS-05 Bagikan sebagai gambar

Kartu ringkasan wilayah 1080 × 1350 dirender server (SVG ke PNG) dengan angka kunci, label, sumber, dan tautan pendek.

## 9. Ringkasan endpoint

| Metode | Path | Akses |
| --- | --- | --- |
| GET | `/wilayah/cari`, `/wilayah/titik`, `/wilayah/{kode}`, `/wilayah/{kode}/ringkasan` | publik |
| GET | `/wilayah/{kode}/layanan`, `/wilayah/{kode}/uang`, `/wilayah/{kode}/tanggung-jawab` | publik |
| GET | `/peta/kesenjangan`, `/peta/paket`, `/peta/laporan`, `/peta/sekolah` | publik |
| GET | `/kecamatan/{kode}/sekolah`, `/sekolah/terdekat`, `/keterbukaan/daya-tampung` | publik |
| GET | `/paket/{kode}`, `/paket/{kode}/pembanding`, `/paket/terdekat`, `/penyedia/{id}` | publik |
| GET | `/kewenangan`, `/kewenangan/daftar`, `/dapil/{id}`, `/anggota/{id}` | publik |
| GET | `/provenance/{id}`, `/unduh/paket.csv`, `/sumber-data` | publik |
| POST | `/tanya`, `/aspirasi` | publik dengan batas laju |
| POST | `/auth/otp/kirim`, `/auth/otp/verifikasi`, `/auth/refresh`, `/auth/keluar` | publik dengan batas laju |
| GET | `/auth/saya`, `/notifikasi` | warga |
| POST | `/laporan`, `/laporan/{id}/konfirmasi`, `/langganan` | warga |
| GET, POST, PATCH | `/panel/*`, `/hak-jawab` | instansi |
| GET, PATCH | `/moderasi/*` | moderator |
| POST, GET | `/admin/ingest/{sumber}`, `/admin/ingest/runs` | admin |

## 10. Keamanan

- Header: CSP tanpa `unsafe-inline` untuk skrip, HSTS, `X-Content-Type-Options`, `Referrer-Policy: strict-origin-when-cross-origin`, `Permissions-Policy` yang hanya membuka geolocation dan camera untuk halaman Lapor.
- Batas laju per IP dan per akun di middleware; respons 429 dengan `Retry-After`.
- CSRF: token refresh di cookie `SameSite=Lax`, endpoint yang mengubah data memerlukan header `Authorization` dari token akses di memori.
- Unggahan: validasi magic bytes, batas ukuran, dekode dan enkode ulang, kunci objek acak, bucket privat dengan URL bertanda tangan berumur pendek untuk foto yang belum terverifikasi.
- Query terparameterisasi; tidak ada SQL dari input pengguna, termasuk dari Tanya Rantai.
- Log tanpa email dan koordinat asli.

## 11. Alur Demo dan data uji

Skenario demo dikunci di `data/snapshot/skenario_demo.yaml`:

- `kecamatan_demo`: dipilih setelah GD-04, yaitu kecamatan di Jakarta Utara dengan rasio kursi SMP negeri terendah yang datanya lengkap.
- `paket_simulasi`: satu paket kategori `ruang_kelas` di kecamatan demo, nama penyedia fiktif "CV Contoh Simulasi", `is_simulation = true`.
- `laporan_simulasi`: dua laporan dan konfirmasi yang cukup untuk status TERVERIFIKASI, seluruhnya `is_simulation = true`, foto ilustrasi bebas lisensi.
- `akun_demo`: satu warga, satu petugas instansi demo, satu moderator.
- `tanya_demo`: tiga pertanyaan dan jawaban tercache.

Tes Playwright `e2e/alur-demo.spec.ts` menjalankan Alur Demo PRD bagian 7 langkah 1 sampai 7 dalam mode demo, pada lebar 390 px dan 1440 px.

## 12. Bentuk `docs/DESIGN.md`

Dokumen desain yang dihasilkan dari prompt desain wajib memuat bagian berikut agar bisa dipakai langsung oleh frontend:

1. Token warna (terang dan gelap), termasuk warna semantik label kepastian, sektor, status laporan, dan skala choropleth.
2. Token tipografi, jarak, radius, bayangan, durasi animasi.
3. Konfigurasi `tailwind.config` dari token tersebut.
4. Daftar komponen dengan varian dan status.
5. Tata letak tiap layar di empat lebar (360, 768, 1024, 1440).
6. Teks antarmuka utama dalam bahasa Indonesia dan Inggris.
