# Spesifikasi Fitur: Rantai

Versi 2.0 · 23 Oktober 2026 · Pasangan dari `docs/PRD.md` v2.0. ID fitur di sini sama dengan ID di PRD.

Padanan ID dari v1.0, untuk membaca riwayat commit dan tugas lama:

| v1.0 | v2.0 | v1.0 | v2.0 |
| --- | --- | --- | --- |
| DASAR-01 | WIL-01 | UANG-01 | ANG-01 |
| DASAR-02 | WIL-04 | UANG-02 | ANG-02 |
| DASAR-03 Wilayahku | WIL-02 Halaman kecamatan | UANG-03 | ANG-03 |
| DASAR-04 | WIL-03 | UANG-04, UANG-05 | ANG-06 |
| DASAR-05 | DSR-02 | UANG-06 | ANG-04 |
| DASAR-06 | DSR-04 | UANG-07 | ANG-07 |
| DASAR-07 | DSR-01 | UANG-08, UANG-09 | ANG-08 |
| DASAR-08 | WIL-05 | UANG-10 | ANG-09 |
| LAY-01 | PEN-01 | TJ-01 | Kontrak tema, bagian Tanggung jawab |
| LAY-02 | PEN-02 | TJ-02 | WAK-01 |
| LAY-03 | PEN-03 | TJ-03 | WAK-02 |
| LAY-04 | PEN-06 | TJ-04 | AKS-03 |
| LAY-05 | Dihapus, kursi kini memakai data resmi | TJ-05 | AKS-02 |
| LAY-06 | PEN-07 | LINTAS-01 | AKS-01 |
| LINTAS-02 | DSR-03, WAK-07 | LINTAS-03 | WIL-06 |
| LINTAS-04 | DSR-03 | LINTAS-05 | AKS-05 |

## 0. Konvensi

### 0.1 Label kepastian

Setiap angka yang tampil membawa satu label:

| Label | Arti | Contoh |
| --- | --- | --- |
| `resmi` | Diambil dari sumber pemerintah, apa adanya atau dijumlahkan tanpa asumsi | Daya tampung SMP negeri dari Kepkadis 258/2026; pagu paket dari LPSE |
| `proksi` | Dihitung Rantai dari data resmi dengan rumus terbuka yang memuat asumsi | Kursi madrasah dari rombel kelas 7 kali batas murid |
| `estimasi` | Perkiraan karena data resmi tidak dipublikasikan | Akhir masa pemeliharaan; siswa kelas 6 dari total SD dibagi 6 |
| `simulasi` | Data buatan untuk demonstrasi | Laporan warga di paket demo |

Angka turunan memakai label terlemah dari masukannya, dengan urutan `resmi` > `proksi` > `estimasi` > `simulasi`. Contoh: kesenjangan dari kursi `resmi` dan calon lulusan `estimasi` berlabel `estimasi`.

### 0.2 Tipe Angka

Setiap nilai numerik di respons API dikirim sebagai objek `Angka`, bukan angka telanjang:

```json
{ "nilai": 7360, "satuan": "kursi", "label": "resmi", "provenance_id": "prv_01J..." }
```

`Angka` menolak diserialisasi tanpa label sah dan `provenance_id`, sehingga angka tanpa rujukan tidak pernah sampai ke respons.

### 0.3 Amplop respons API

```json
{ "code": 200, "message": "ok", "data": {}, "meta": { "page": 1, "per_page": 20, "total": 134, "total_pages": 7 } }
```

Error: `{ "code": 404, "message": "paket tidak ditemukan", "data": null }`. Pesan error berbahasa Indonesia yang bisa dipahami warga; frontend memetakan ke i18n berdasarkan `code` dan kunci error.

### 0.4 Prefiks dan versi

Semua endpoint di bawah `/api/v1`. Kode wilayah memakai kode administrasi Kemendagri (contoh format `31.72.01.1001`).

### 0.5 Path

`config/` di dokumen ini berarti `backend/config/`. `e2e/` berarti `frontend/e2e/`. `data/snapshot/` ada di akar repo.

### 0.6 Tingkat analisis

Lingkup analisis adalah kecamatan. Kelurahan dipakai untuk pencarian, penentuan lokasi paket, dan jalur aspirasi; halaman publik selalu bermuara di halaman kecamatan dengan kelurahan yang dicari disorot.

## 1. Gerbang Data

Tugas bertanda **[GERBANG DATA]** di PLAN.md dimulai dengan verifikasi ini. Hasil verifikasi ditulis ke `docs/data-sources.md` dengan URL, tanggal cek, field yang tersedia, dan isi robots.txt.

| ID | Yang diverifikasi | Bila tidak sesuai |
| --- | --- | --- |
| GD-01 | SPSE LPSE DKI (`spse.inaproc.id/jakarta`): endpoint JSON daftar tender dan non-tender yang dipakai halaman (pola yang sama dengan pustaka pyproc), halaman pengumuman paket (kode paket, kode RUP, satker, pagu, HPS, jenis kontrak, lokasi, jumlah peserta), robots.txt, dan apakah bisa diakses dari IP luar Indonesia | Bila daftar tertutup, pakai API Gateway Inaproc (pengganti ISB LKPP) sebagai sumber daftar dan SPSE hanya untuk detail. Bila IP luar diblokir, jalankan importer dari server Indonesia dan isi lingkungan lain dari snapshot. **Lolos sebagian**: detail tender tertutup untuk publik, lihat `docs/data-sources.md` |
| GD-02 | Halaman pemenang tender dan pemenang berkontrak: nama penyedia, nilai kontrak, tanggal kontrak, jadwal pelaksanaan | Bila nilai kontrak tidak tersedia, tampilkan HPS saja dan sembunyikan perbandingan kontrak |
| GD-03 | SiRUP atau ISB Satu Data LKPP: detail paket per kode RUP (nama, satker, pagu, sumber dana, waktu pemilihan) | Mata rantai RUP tampil sebagai "kode RUP tercatat, detail tidak tersedia" |
| GD-04 | Keputusan Kadis Pendidikan DKI No. 258 Tahun 2026: lampiran daya tampung SMP negeri per sekolah dengan NPSN; total 73.289 kursi | Bila lampiran tidak bisa diekstrak per sekolah, pakai total per kota dari dokumen yang sama dan tampilkan peringatan tingkat kota |
| GD-05 | Pengumuman Disdik DKI SPMB Bersama 2026: daftar SMP swasta mitra dan kursinya per sekolah; total 1.597 kursi | Bila tidak per sekolah, tampilkan total provinsi saja dan keluarkan dari angka kecamatan dengan catatan |
| GD-06 | Dapodik: data sekolah per NPSN (nama, bentuk, status, koordinat, kecamatan) untuk SD, MI, SMP, MTs di DKI; ketersediaan data madrasah | Cadangan: Data Referensi Kemendikdasmen per wilayah. Madrasah tanpa data per sekolah tampil sebagai catatan keterbatasan |
| GD-07 | KPU: daftar 10 dapil DPRD DKI 2024 dan kecamatan anggotanya; 106 calon terpilih per dapil (9 sampai 12 kursi per dapil) | Cadangan: berita acara penetapan KPU DKI dan daftar resmi DPRD DKI |
| GD-08 | Situs DPRD DKI: lima komisi dan bidangnya, keanggotaan komisi, alat kelengkapan dewan, fraksi, halaman profil per anggota periode 2024 sampai 2029 | Bila keanggotaan komisi tidak tersedia, WAK-01 menampilkan anggota dapil tanpa urutan komisi dan menyebut komisi pengawas saja |
| GD-09 | Ambang nilai pengadaan langsung Pekerjaan Konstruksi menurut Perpres 16/2018 beserta perubahan terakhirnya (Perpres 12/2021, Perpres 46/2025) | Sinyal S-03 dimatikan bila ambang tidak bisa dipastikan |
| GD-10 | Batas kelurahan dan kecamatan DKI dengan kode Kemendagri | Cadangan: GADM tingkat 4 dengan tabel padanan kode. **Lolos**, lihat `docs/data-sources.md` |
| GD-11 | Opentender OCDS: unduhan per tahun, field buyer, supplier, nilai, tanggal | ANG-07 hanya memakai data LPSE |
| GD-12 | Siswa kelas 6 per sekolah atau per kecamatan untuk tahun ajaran berjalan: Sekolah Kita atau Dapodik | Cadangan berurutan: Satu Data Jakarta, permintaan data resmi, estimasi `total siswa SD ÷ 6` berlabel `estimasi` |
| GD-13 | Berita situs DPRD DKI: daftar dan halaman berita, struktur judul, tanggal, nama anggota yang disebut, robots.txt | Cadangan: impor manual tautan rilis untuk anggota di skenario demo |
| GD-14 | Liputan media: daftar perusahaan pers terverifikasi Dewan Pers, RSS per media atau RSS pencarian per nama anggota, ketentuan pemakaian dan robots.txt | Bila RSS pencarian tidak diizinkan, pakai RSS kanal tiap media terverifikasi lalu saring nama secara lokal |
| GD-15 | Dokumen dan berita paripurna DPRD DKI: pandangan umum fraksi atas raperda dan APBD, kompilasi aspirasi reses per bidang | Bila tidak ada dokumen, WAK-05 dan WAK-06 memakai kurasi manual berbasis berita resmi dengan tautan sumber |

## 2. Model data

PostgreSQL 16 dengan PostGIS. Semua tabel punya `dibuat_pada` dan `diubah_pada` (`timestamptz`). Kolom geometri memakai SRID 4326.

### 2.1 Wilayah dan provenance

| Tabel | Kolom penting |
| --- | --- |
| `wilayah` | `kode` PK, `tingkat` (provinsi, kota, kecamatan, kelurahan), `nama`, `induk_kode`, `alias` text[], `geom` MultiPolygon, `titik_tengah` Point, `provenance_id` |
| `provenance` | `id` PK (ULID berawalan `prv_`), `sumber` (kemendagri, lpse, sirup, dapodik, disdik, kpu, dprd, media, opentender, rantai, simulasi), `url`, `diambil_pada`, `label`, `hash_konten`, `catatan` (rumus untuk `proksi` dan `estimasi`) |
| `ingest_run` | `id`, `sumber` (nama job), `mulai`, `selesai`, `status`, `jumlah_baru`, `jumlah_berubah`, `pesan_error` |

Fungsi SQL `hierarki_wilayah(kode)` mengembalikan leluhur sebuah wilayah sebagai JSON berurutan dari provinsi.

### 2.2 Tema

| Tabel | Kolom penting |
| --- | --- |
| `tema` | `kode` PK, `nama_id`, `nama_en`, `aktif`, `urutan`, `konfigurasi` jsonb (indikator, filter anggaran, tanggung jawab, tag isu, templat aspirasi), `provenance_id` |
| `indikator_kecamatan` | `tema_kode`, `kode_kecamatan`, `periode`, `kebutuhan`, `kapasitas`, `kesenjangan`, `rasio` nullable, `kelas`, `label`, `rincian` jsonb, `provenance_id`; PK (`tema_kode`, `kode_kecamatan`, `periode`) |

### 2.3 Tema Pendidikan

| Tabel | Kolom penting |
| --- | --- |
| `sekolah` | `npsn` PK, `nama`, `bentuk` (SD, MI, SMP, MTs), `jenjang` (dasar, menengah_pertama), `status` (negeri, swasta), `is_madrasah`, `kode_kecamatan`, `kode_kelurahan` nullable, `lokasi` Point, `provenance_id` |
| `sekolah_rombel` | `npsn`, `tahun_ajaran`, `tingkat` (1 sampai 9), `jumlah_rombel`, `jumlah_siswa`, `provenance_id`; PK (`npsn`, `tahun_ajaran`, `tingkat`) |
| `daya_tampung` | `npsn`, `tahun_ajaran`, `jalur` (negeri, spmb_bersama), `kursi`, `provenance_id`; PK (`npsn`, `tahun_ajaran`, `jalur`) |
| `calon_lulusan_kecamatan` | `kode_kecamatan`, `tahun_ajaran`, `jumlah`, `label`, `metode` (per_sekolah, agregat_resmi, estimasi_total_sd), `provenance_id` |

### 2.4 Arus Anggaran

| Tabel | Kolom penting |
| --- | --- |
| `instansi` | `id`, `nama` (contoh: Suku Dinas Pendidikan Wilayah I Kota Administrasi Jakarta Utara), `jenis` (dinas, sudin, badan, kelurahan, lainnya), `kode_kota`, `induk_id` |
| `rup` | `kode` PK, `nama`, `instansi_id`, `pagu`, `sumber_dana`, `tahun_anggaran`, `metode`, `provenance_id` |
| `penyedia` | `id`, `nama`, `nama_normal` (huruf kecil, tanpa CV/PT/tanda baca), `alias` text[] |
| `paket` | `kode` PK, `nama`, `kode_rup` FK nullable, `instansi_id`, `jenis_pengadaan`, `metode`, `tahun_anggaran`, `pagu`, `hps`, `nilai_kontrak` nullable, `jenis_kontrak`, `tanggal_dibuat`, `tahap`, `jumlah_peserta`, `penyedia_id` nullable, `jadwal_mulai`, `jadwal_selesai`, `tanggal_serah_terima` nullable, `kategori`, `sektor`, `tema_kode` text[], `lokasi_tingkat` (kelurahan, kecamatan, kota, provinsi), `kode_wilayah`, `lokasi` Point, `rt`, `rw`, `is_simulation`, `provenance_id` |
| `riwayat_penyedia` | `penyedia_id`, `tahun`, `pembeli`, `nama_paket`, `nilai`, `sumber` (opentender, lpse), `provenance_id` |
| `sinyal` | `id`, `paket_kode`, `jenis` (S-01 sampai S-05), `alasan`, `pembanding` jsonb, `aktif`, `dihitung_pada` |
| `laporan` | `id`, `pengguna_id`, `paket_kode`, `kode_wilayah`, `lokasi_publik` Point (dibulatkan ke grid 50 m), `lokasi_asli` Point (privat), `akurasi_meter`, `kondisi` (belum_dimulai, sedang_dikerjakan, terhenti, selesai_baik, selesai_rusak), `catatan`, `status`, `is_simulation` |
| `laporan_foto` | `id`, `laporan_id`, `object_key`, `lebar`, `tinggi` |
| `laporan_konfirmasi` | `laporan_id`, `pengguna_id`, `kondisi`, `lokasi_asli`, `dibuat_pada`; unik (`laporan_id`, `pengguna_id`) |
| `laporan_riwayat` | `laporan_id`, `status_lama`, `status_baru`, `oleh_pengguna_id`, `alasan`, `dibuat_pada` |

### 2.5 Wakilku

| Tabel | Kolom penting |
| --- | --- |
| `dapil` | `id`, `nomor`, `nama`, `tahun_pemilu`, `jumlah_kursi`, `provenance_id`; unik (`tahun_pemilu`, `nomor`) |
| `dapil_wilayah` | `dapil_id`, `kode_kecamatan` (unik: satu kecamatan satu dapil) |
| `komisi` | `id`, `kode` (A sampai E), `nama`, `bidang` text[] |
| `anggota` | `id` (slug halaman profil DPRD), `nama`, `nama_lain` text[] (untuk disambiguasi), `partai`, `fraksi`, `dapil_id`, `komisi_id` nullable, `url_profil`, `lhkpn_url`, `provenance_id` |
| `anggota_akd` | `anggota_id`, `alat_kelengkapan` (contoh: Badan Anggaran, Bapemperda), `jabatan` |
| `umpan_kegiatan` | `id`, `jenis_sumber` (resmi_dprd, media), `media_nama`, `media_domain`, `judul`, `ringkasan` nullable, `url` unik, `terbit_pada`, `tag_isu` text[], `disembunyikan`, `provenance_id` |
| `umpan_anggota` | `umpan_id`, `anggota_id`, `dasar_cocok` (nama_dan_dprd, nama_dan_fraksi, situs_resmi) |
| `sikap_fraksi` | `id`, `fraksi`, `agenda` (raperda, apbd), `judul_agenda`, `tanggal`, `ringkasan`, `url`, `provenance_id` |
| `aspirasi_reses` | `id`, `masa_reses`, `dapil_id` nullable, `bidang`, `uraian`, `kode_wilayah` nullable, `provenance_id` |

### 2.6 Pengguna dan interaksi

| Tabel | Kolom penting |
| --- | --- |
| `pengguna` | `id`, `email` (unik, privat), `nama_tampilan`, `peran` (warga, anggota, instansi, moderator, admin), `instansi_id` nullable, `anggota_id` nullable |
| `otp` | `email`, `kode_hash`, `kedaluwarsa`, `percobaan` |
| `sesi_refresh` | `id`, `pengguna_id`, `token_hash`, `kedaluwarsa`, `dicabut` |
| `hak_jawab` | `id`, `target_jenis` (sinyal, laporan, paket, umpan), `target_id`, `instansi_id` nullable, `anggota_id` nullable, `isi`, `penjawab_id`, `dibuat_pada` |
| `langganan` | `id`, `pengguna_id`, `kode_wilayah`, `kanal` (feed, webpush, email), `webpush_endpoint` jsonb |
| `notifikasi` | `id`, `pengguna_id`, `jenis`, `payload` jsonb, `dibaca` |
| `tanya_log` | `id`, `pertanyaan`, `jawaban`, `rujukan` jsonb, `bahasa`, `durasi_ms` (tanpa data pribadi) |

## 3. Pipeline data

Worker Go (`cmd/worker <job>`) menjalankan job. Setiap job: ambil, simpan mentah ke `data/raw/<job>/<tanggal>/`, normalisasi, cocokkan, tulis bersama `provenance`, catat `ingest_run`. Menjalankan job dua kali tidak menggandakan data.

| Job | Jadwal | Isi |
| --- | --- | --- |
| `wilayah` | sekali, manual | Batas wilayah dan alias (GD-10) |
| `tema` | saat konfigurasi berubah | Muat `config/tema/*.yaml` ke tabel `tema` dengan validasi |
| `dapodik` | per semester | Sekolah, lokasi, rombel dan siswa per tingkat (GD-06) |
| `daya_tampung` | tahunan, manual | Impor CSV hasil ekstraksi Kepkadis 258/2026 dan SPMB Bersama (GD-04, GD-05) |
| `calon_lulusan` | per semester | Siswa kelas 6 per kecamatan (GD-12) |
| `lpse_daftar` | harian 02.00 WIB | Daftar paket tender dan non-tender tahun berjalan dan sebelumnya |
| `lpse_detail` | harian, setelah daftar | Detail paket baru dan paket yang tahapnya berubah; pemenang (GD-02) |
| `sirup` | mingguan | Detail RUP untuk kode RUP yang ditemukan di paket |
| `opentender` | sekali | Riwayat 2014 sampai 2023 untuk pembeli Pemprov DKI |
| `kpu_dprd` | sekali per periode | Dapil per kecamatan dari `config/dprd.yaml` (PKPU 6/2023), anggota, fraksi, komisi, alat kelengkapan dari situs DPRD DKI (GD-07, GD-08) |
| `umpan_dprd` | harian | Berita situs DPRD DKI yang menyebut anggota (GD-13) |
| `umpan_media` | harian | Liputan media terverifikasi per anggota (GD-14) |
| `ringkas_umpan` | setelah job umpan | Ringkasan netral satu kalimat untuk item umpan baru |
| `paripurna` | per masa sidang, manual | Impor CSV kurasi sikap fraksi dan aspirasi reses (GD-15) |
| `turunan` | setelah job sumber | Kategori dan tema paket, penentuan lokasi, indikator tema, sinyal, kejadian notifikasi |

Aturan pengambilan: 1 permintaan per detik per domain, User-Agent `RantaiBot/1.0 (+https://github.com/benditandayusaputra/rantai)`, header `If-Modified-Since` bila didukung, berhenti dan catat error bila menerima 403 atau 429, patuhi robots.txt. Portal pengadaan sering menolak IP luar Indonesia; job LPSE dan SiRUP dijalankan dari server di Indonesia, dan lingkungan lain memakai snapshot.

## 4. Kontrak tema

Tema adalah satu berkas `config/tema/<kode>.yaml` dan satu importer data indikator. Job `tema` memvalidasi berkas lalu menyimpannya di tabel `tema`; API dan frontend hanya membaca dari basis data, tidak ada halaman tema yang di-hardcode.

```yaml
kode: pendidikan
nama: { id: Pendidikan, en: Education }
aktif: true
urutan: 1
indikator:
  penghitung: kesenjangan_kursi_smp
  satuan_kebutuhan: siswa
  satuan_kapasitas: kursi
  judul: { id: "Kursi SMP negeri dibanding lulusan SD", en: "Public junior high seats vs primary graduates" }
  ambang_rasio: { sangat_kurang: 0.5, kurang: 0.8, hampir_cukup: 1.0 }
filter_anggaran:
  instansi_pola: ["dinas pendidikan", "suku dinas pendidikan"]
  kata_kunci: ["ruang kelas", "rehab sekolah", "pembangunan smp", "gedung sekolah", "\\b(sdn|smpn)\\b"]
tanggung_jawab:
  - masalah: kursi_smp_kurang
    label: { id: "Kursi SMP negeri kurang", en: "Not enough public junior high seats" }
    instansi_pelaksana: Dinas Pendidikan Provinsi DKI Jakarta
    unit_wilayah: "Suku Dinas Pendidikan {kota}"
    komisi: E
    dasar_hukum: "UU 23/2014 lampiran I.A pembagian urusan pendidikan; UU 2/2024 tentang Provinsi Daerah Khusus Jakarta"
    dasar_hukum_url: "https://peraturan.bpk.go.id/"
tag_isu: [sekolah, spmb, zonasi, daya tampung, guru]
templat_aspirasi:
  judul: { id: "Kekurangan kursi SMP negeri di Kecamatan {kecamatan}", en: "Shortage of public junior high seats in {kecamatan} District" }
  pertanyaan:
    - id: "Berapa rencana penambahan rombel SMP negeri di Kecamatan {kecamatan} untuk tahun ajaran {tahun_depan}?"
      en: "How many additional public junior high classes are planned in {kecamatan} District for the {tahun_depan} school year?"
```

| Bagian | Dipakai oleh | Aturan validasi |
| --- | --- | --- |
| `indikator` | Halaman kecamatan, peta, Tanya Rantai | `penghitung` harus terdaftar di kode; `ambang_rasio` naik berurutan |
| `filter_anggaran` | Arus Anggaran, PEN-04 | Minimal satu pola instansi atau kata kunci; regex harus bisa dikompilasi |
| `tanggung_jawab` | Wakilku, Paket Aspirasi | Setiap entri wajib punya `dasar_hukum`, `dasar_hukum_url`, dan `komisi` A sampai E |
| `tag_isu` | Umpan Wakilku, WAK-04, PEN-05 | Minimal satu tag, huruf kecil |
| `templat_aspirasi` | Paket Aspirasi | Hanya placeholder yang dikenal: `{kecamatan}`, `{kota}`, `{tahun_depan}` |

Penghitung indikator menulis ke `indikator_kecamatan` dalam bentuk umum (kebutuhan, kapasitas, kesenjangan, rasio, kelas), sehingga peta, ringkasan kecamatan, PEN-04, dan Tanya Rantai bekerja untuk tema apa pun tanpa kode tambahan.

Tema berikutnya (Banjir, Kesehatan, Jalan) masuk Rencana Pengembangan; pemetaan komisinya dikonfirmasi saat tema dibangun.

API: `GET /tema` (tema aktif berurutan), `GET /tema/{kode}`.

## 5. Wilayah

### WIL-01 Pencarian wilayah

- Input teks: nama kelurahan, kecamatan, kota, atau alamat bebas. Skor kemiripan memakai `pg_trgm` (`similarity` dan `word_similarity`) atas `nama` dan `alias`, ambang 0,45. Urutan: kecocokan persis, kemiripan, lalu tingkat (kelurahan dahulu). Bila sebuah kecamatan cocok persis, seluruh kelurahannya ikut tampil.
- Setiap hasil membawa `hierarki` (provinsi, kota, kecamatan) agar nama kembar bisa dibedakan.
- Tombol "Pakai lokasiku": `navigator.geolocation` lalu `GET /wilayah/titik`, yang mengembalikan kelurahan yang memuat titik.
- Tidak memanggil layanan geocoding eksternal.
- Memilih hasil kelurahan membuka halaman kecamatannya dengan kelurahan tersebut disorot.

API: `GET /wilayah/cari?q=`, `GET /wilayah/titik?lat=&lng=`, `GET /wilayah/{kode}`.

Kriteria terima: "kelapa gading" mengembalikan Kecamatan Kelapa Gading dan tiga kelurahannya; titik di Jakarta Utara mengembalikan kelurahan yang benar; respons di bawah 150 ms. **Terpenuhi** (T-104).

### WIL-02 Halaman kecamatan

Halaman `/{bahasa}/kecamatan/{kode}`. Header: nama kecamatan, kota, kelurahan yang disorot (bila datang dari pencarian kelurahan), tombol Pantau, tombol Bagikan. Tiga bagian dengan URL sendiri (`?bagian=tema|anggaran|wakilku`, bawaan `tema`):

- **Tema aktif**: kartu indikator tiap tema aktif dari `indikator_kecamatan`, dengan kelas dan lencana sumber. Untuk Pendidikan: calon lulusan SD dibanding kursi SMP negeri.
- **Arus Anggaran**: jumlah dan total nilai paket tahun berjalan di kecamatan, dengan filter tema.
- **Wakilku**: dapil, jumlah anggota, dan anggota yang duduk di komisi pengawas tema aktif lebih dahulu.

API: `GET /kecamatan/{kode}/ringkasan` mengembalikan ringkasan ketiga bagian, setiap angka sebagai `Angka`.

### WIL-03 Peta

Halaman `/{bahasa}/peta`. MapLibre dengan peta dasar vektor bebas lisensi. Lapisan yang bisa dinyalakan:

| Lapisan | Bentuk | Sumber API |
| --- | --- | --- |
| Indikator tema | Choropleth kecamatan per tema aktif | `GET /peta/tema/{kode}` |
| Proyek | Titik untuk lokasi kelurahan, area arsir untuk lokasi kecamatan; warna per sektor | `GET /peta/paket?bbox=&tema=&tahun=` |
| Laporan warga | Titik dengan status | `GET /peta/laporan?bbox=` |
| Sekolah | Titik per bentuk sekolah (lapisan khusus tema Pendidikan) | `GET /peta/sekolah?bbox=` |

Setiap lapisan punya alternatif tabel yang bisa diakses papan ketik. Respons GeoJSON disederhanakan (`ST_SimplifyPreserveTopology`), di bawah 300 KB untuk seluruh DKI, dan di-cache.

### WIL-04 Lencana sumber

Komponen kecil di samping setiap angka. Menampilkan label kepastian sebagai warna dan ikon. Klik atau Enter membuka panel: sumber, tautan ke halaman asli, waktu diambil, dan rumus bila `proksi` atau `estimasi` (dari `provenance.catatan`).

API: `GET /provenance/{id}`. **Terpenuhi sisi backend** (T-105).

### WIL-05 Tentang Data

Halaman berisi: daftar sumber dan tanggal ambil terakhir (dari `ingest_run`, bukan ditulis tangan), rumus indikator tiap tema, aturan penentuan lokasi paket dan cakupannya, definisi tiap sinyal, label kepastian dan aturan label terlemah, daftar media terverifikasi dan aturan disambiguasi, pernyataan bias liputan, keterbatasan yang diketahui, atribusi dan lisensi ODbL untuk data Opentender, serta tautan unduh turunan data.

API: `GET /sumber-data`.

### WIL-06 Pantau Wilayahku

Warga berlangganan kecamatan. Kanal: feed di aplikasi (wajib), Web Push dengan VAPID, email ringkasan mingguan (cadangan P2). Kejadian:

| Kejadian | Pemicu |
| --- | --- |
| Paket baru | Paket baru di kecamatan |
| Melewati jadwal | Sinyal S-05 baru |
| Masa pemeliharaan segera berakhir | 30 hari sebelum akhir pemeliharaan |
| Laporan terverifikasi | Laporan di kecamatan berstatus TERVERIFIKASI |
| Hak jawab baru | Hak jawab untuk paket, laporan, atau umpan anggota dapil |
| Indikator diperbarui | Data indikator tema periode baru |

API: `POST /langganan`, `DELETE /langganan/{id}`, `GET /notifikasi`, `PATCH /notifikasi/{id}/dibaca`.

## 6. Tema Pendidikan

### PEN-01 Kesenjangan kursi SMP per kecamatan

Penghitung `kesenjangan_kursi_smp`. Rumus per kecamatan `k` dan tahun ajaran `t`:

```text
calon_lulusan(k, t)       = siswa kelas 6 SD dan MI dalam k (calon_lulusan_kecamatan)
kursi_negeri(k, t)        = Σ daya_tampung jalur negeri, SMP negeri dalam k             label resmi
kursi_spmb_bersama(k, t)  = Σ daya_tampung jalur spmb_bersama, SMP swasta mitra dalam k label resmi
kursi_madrasah(k, t)      = Σ jumlah_rombel tingkat 7 MTs × batas_murid_per_rombel_smp  label proksi
kesenjangan(k, t)         = calon_lulusan − kursi_negeri
rasio(k, t)               = kursi_negeri ÷ calon_lulusan
```

`indikator_kecamatan` menyimpan `kebutuhan = calon_lulusan`, `kapasitas = kursi_negeri`, serta rincian ketiga jenis kursi dan kesenjangan bila SPMB Bersama dan madrasah ikut dihitung. Label mengikuti aturan label terlemah (0.1).

Konstanta di `config/pendidikan.yaml`: `batas_murid_per_rombel_smp: 32` (Kepmendikdasmen 14/2026) dan pembagi cadangan `estimasi_total_sd: 6`. Kelas rasio dari `ambang_rasio` di konfigurasi tema. Tidak ada angka 32 atau 6 tertulis langsung di kode.

Keterbatasan yang wajib ditampilkan: SPMB Jakarta tidak membatasi pendaftaran pada kecamatan, sehingga angka kecamatan menggambarkan tekanan lokal, bukan peluang pasti seorang anak. PEN-06 melengkapi dengan tampilan radius.

Kasus tepi yang wajib dites: kecamatan tanpa SMP negeri (kursi 0, rasio 0); calon lulusan 0 (rasio tidak didefinisikan, tampil "tidak ada data lulusan"); sekolah di keputusan daya tampung yang NPSN-nya tidak ada di Dapodik (dikeluarkan dan dilaporkan sebagai "sekolah tidak terpetakan"); calon lulusan dari estimasi membuat label kesenjangan `estimasi`.

### PEN-02 Peta kesenjangan

Choropleth kecamatan dari `indikator_kecamatan` dengan skala berurutan empat kelas. Tooltip: nama kecamatan, calon lulusan, kursi negeri, selisih, label. Legenda dalam bahasa awam ("kursi negeri kurang dari separuh lulusan").

API: `GET /peta/tema/pendidikan`.

### PEN-03 Detail kecamatan

Bagian Tema di halaman kecamatan, dengan rincian di `/{bahasa}/kecamatan/{kode}/pendidikan`:

- Grafik batang horizontal: calon lulusan dibanding kursi negeri, SPMB Bersama, madrasah; dengan tabel alternatif.
- Daftar SMP dan MTs: nama, bentuk, status, daya tampung resmi atau kursi proksi dengan labelnya, tautan ke peta.
- Tautan ke Arus Anggaran dengan filter tema Pendidikan di kecamatan yang sama.

API: `GET /tema/pendidikan/kecamatan/{kode}`, `GET /tema/pendidikan/kecamatan/{kode}/sekolah`.

### PEN-04 Kesenjangan vs Belanja

Per kecamatan: rasio dari `indikator_kecamatan` dibanding total nilai paket yang lolos `filter_anggaran` tema pada tahun anggaran yang sama dan berlokasi di kecamatan itu (lokasi tingkat kelurahan atau kecamatan). Ditampilkan sebagai diagram sebar seluruh kecamatan DKI dengan kecamatan aktif disorot, plus tabel alternatif.

Kalimat pendamping wajib: "Belanja yang lokasinya hanya diketahui sampai tingkat kota tidak ikut dihitung; grafik ini menunjukkan sebaran, bukan penilaian." Jumlah dan nilai paket tingkat kota disebut terpisah.

API: `GET /tema/{kode}/belanja?tahun=`.

### PEN-05 Anggota Komisi E yang menyuarakan isu sekolah

Di bagian Tema: anggota dapil kecamatan yang duduk di komisi pengawas tema, beserta item umpan WAK-03 yang memuat `tag_isu` tema, maksimal tiga item terbaru per anggota dengan lencana jenis sumber. Tanpa perbandingan antaranggota.

API: `GET /kecamatan/{kode}/wakil?tema=pendidikan&dengan_umpan=true`.

### PEN-06 Cek Sekolahku

Input alamat atau lokasi. Keluaran:

- SMP negeri dalam radius 3 km (konfigurasi), diurutkan dari terdekat, dengan jarak garis lurus dan daya tampung resmi.
- `tekanan_radius = siswa kelas 6 dalam radius ÷ kursi negeri dalam radius`, dijelaskan dengan kalimat ("sekitar 2 lulusan untuk setiap 1 kursi negeri di sekitarmu"). Label mengikuti aturan label terlemah.
- Alternatif: sekolah SPMB Bersama dan madrasah dalam radius yang sama.

API: `GET /sekolah/terdekat?lat=&lng=&radius_m=`.

### PEN-07 Tren

Grafik garis rasio per kecamatan antar tahun ajaran, tampil hanya bila data minimal dua tahun tersedia.

## 7. Wakilku

### WAK-01 Kecamatan ke dapil ke anggota

Dari kecamatan ke dapil (`dapil_wilayah`) ke anggota. Bila tema dipilih, anggota yang duduk di komisi pengawas tema itu tampil lebih dahulu dengan penanda "duduk di Komisi X yang mengawasi urusan ini". Urutan selanjutnya alfabetis. Tanpa angka kinerja, skor, atau peringkat.

API: `GET /kecamatan/{kode}/wakil?tema=`, `GET /dapil/{id}`.

Kriteria terima: setiap kecamatan DKI masuk tepat satu dapil; urutan anggota komisi pengawas lebih dahulu dites.

### WAK-02 Profil anggota

Nama, fraksi, dapil dan daftar kecamatannya, komisi dan bidangnya, alat kelengkapan dewan beserta jabatan, tautan ke e-Announcement LHKPN (halaman pencarian KPK), tautan ke halaman resmi DPRD. Tanpa foto bila lisensi foto tidak jelas; gunakan inisial. Data LHKPN tidak disalin.

API: `GET /anggota/{id}`.

### WAK-03 Jejak Kegiatan

Umpan dua sumber per anggota, tampil berdampingan:

| Jenis sumber | Asal | Lencana |
| --- | --- | --- |
| `resmi_dprd` | Berita situs DPRD DKI (GD-13), ditulis Sekretariat DPRD | "Rilis resmi DPRD" |
| `media` | Perusahaan pers terverifikasi Dewan Pers (GD-14), daftar tertutup di `config/media.yaml` | "Media: {nama media}" |

Aturan:

- **Disambiguasi**: artikel media dipautkan ke anggota hanya bila nama (atau `nama_lain`) muncul bersama "DPRD DKI" atau nama fraksinya dalam artikel yang sama (judul atau cuplikan). Item tanpa dasar cocok dibuang. Berita situs resmi dipautkan lewat nama anggota yang disebut.
- **Ringkasan, bukan salinan**: Rantai menyimpan dan menampilkan judul, ringkasan netral satu kalimat, nama media, tanggal, dan tautan. Teks berita penuh tidak disimpan.
- **Ringkasan netral**: job `ringkas_umpan` membuat satu kalimat (maksimal 160 karakter) dari judul dan cuplikan memakai klien LLM yang sama dengan Tanya Rantai, model dari env `RINGKAS_MODEL`. Validator menolak ringkasan yang memuat kosakata terlarang, kata penilaian, atau angka yang tidak ada di sumber. Bila ditolak atau model tidak tersedia, ringkasan kosong dan hanya judul yang tampil.
- **Tag isu**: setiap item diberi `tag_isu` dari pencocokan kata kunci `tag_isu` semua tema.
- **Bias liputan diakui**: bagian umpan selalu memuat pernyataan "Anggota yang lebih sering diliput akan terlihat lebih aktif. Umpan ini bukan ukuran kinerja."
- Moderator bisa menyembunyikan item yang salah atribusi; item tersembunyi tercatat dengan alasan.

API: `GET /anggota/{id}/umpan?jenis=&tag=&halaman=`.

Kriteria terima: tes disambiguasi dengan fixture artikel yang menyebut nama sama tanpa "DPRD DKI" (dibuang) dan dengan fraksi (dipautkan); tes validator ringkasan menolak kosakata terlarang.

### WAK-04 Isu yang Disuarakan

Untuk satu anggota: daftar `tag_isu` yang muncul di umpannya, masing-masing dengan label jumlah sumber ("disebut di 5 sumber: 2 rilis resmi, 3 media") dan tautan ke item. Tidak ada perbandingan atau urutan antaranggota. Kata kerja yang dipakai: "menyuarakan", bukan "memperjuangkan".

API: `GET /anggota/{id}/isu`.

### WAK-05 Sikap Fraksi

Pandangan umum fraksi anggota atas raperda dan APBD (GD-15), sebagai daftar agenda dengan ringkasan netral dan tautan dokumen. Menjadi pengganti suara per anggota yang tidak publik.

API: `GET /fraksi/{nama}/sikap`.

### WAK-06 Aspirasi Reses

Aspirasi agregat per bidang dari dokumen paripurna, difilter ke dapil kecamatan bila dokumen menyebut dapil. Tampil dengan catatan "kompilasi per bidang, bukan per anggota".

API: `GET /aspirasi-reses?dapil=&bidang=`.

### WAK-07 Hak jawab anggota

Akun berperan `anggota` bisa menulis hak jawab pada item umpan yang dipautkan kepadanya (`hak_jawab.target_jenis = umpan`). Hak jawab tampil di bawah item dengan nama anggota dan tanggal; tidak bisa dihapus, hanya ditambah tanggapan baru.

## 8. Arus Anggaran

### ANG-01 Impor paket dan taut RUP

- Sumber: GD-01, GD-02, GD-03. Lingkup: semua jenis pengadaan TA 2025 dan 2026; prioritas tampil Pekerjaan Konstruksi.
- Daftar paket diambil dari endpoint JSON yang dipakai halaman SPSE (pola yang sama dengan pustaka pyproc), diimplementasikan di Go; detail dari halaman pengumuman. Bila kelak memakai API Gateway Inaproc, sumber diganti tanpa mengubah usecase.
- Hasil GD-01: halaman detail dan evaluasi tender tertutup untuk publik. Paket tender disimpan dari kolom daftar saja (kode, nama, tahap, HPS ringkas, metode, jenis pengadaan, tahun anggaran, nilai kontrak); `kode_rup`, satker, pagu, jumlah peserta, dan pemenang tender kosong, dan instansinya K/L/PD. HPS ringkas dibulatkan sumber ke satu desimal juta atau miliar; provenance mencatat pembulatan ini. Paket non-tender memakai halaman pengumuman lengkap.
- Paket ditautkan ke `rup` lewat kode RUP di halaman paket.
- `instansi` dibuat dari nama satuan kerja; kota instansi diambil dari nama dan dipakai penentuan lokasi.
- `penyedia.nama_normal`: huruf kecil, buang "pt", "cv", "tbk", "persero", tanda baca, spasi ganda. Pencocokan penyedia lintas sumber memakai `nama_normal` persis; tidak ada pencocokan fuzzy otomatis.
- Perubahan tahap paket dicatat agar notifikasi "paket baru" dan "tahap berubah" bisa dibuat.
- `paket.tema_kode` diisi job `turunan` dari `filter_anggaran` setiap tema aktif.

Kriteria terima: importer lulus tes dengan fixture minimal 5 berkas (daftar tender, daftar non-tender, pengumuman non-tender, paket tanpa pemenang, paket dengan pemenang); kasus tanpa kode RUP dites dengan potongan HTML; tidak ada permintaan jaringan di tes.

### ANG-02 Penentuan lokasi paket

Gazetir dibangun dari tabel `wilayah` (nama dan alias). Langkah:

1. Normalisasi nama paket: huruf kecil, ganti "kel." menjadi "kelurahan", "kec." menjadi "kecamatan", "jl." menjadi "jalan", rapikan spasi.
2. Ekstrak RT dan RW dengan pola `rt\s*0*(\d{1,3})` dan `rw\s*0*(\d{1,3})`.
3. Cari semua kemunculan nama kelurahan dan kecamatan dari gazetir, utamakan yang didahului kata "kelurahan" atau "kecamatan".
4. Saring kandidat dengan kota instansi: kandidat di luar kota instansi dibuang (sudin bekerja di kotanya sendiri; dinas provinsi tidak disaring).
5. Tentukan hasil:
   - tepat satu kelurahan cocok dan konsisten dengan kecamatan yang disebut: `lokasi_tingkat = kelurahan`, titik di `titik_tengah` kelurahan;
   - tidak ada kelurahan tetapi tepat satu kecamatan: `lokasi_tingkat = kecamatan`;
   - selain itu: `lokasi_tingkat = kota` dari instansi dengan label "lokasi tingkat kota", atau `provinsi` untuk dinas provinsi.
6. Titik beberapa paket di kelurahan yang sama diberi pergeseran deterministik kecil (berbasis hash kode paket) agar tidak menumpuk, dan ditandai "lokasi kelurahan, bukan titik pasti".

Target: minimal 70% paket Pekerjaan Konstruksi mendapat tingkat kecamatan atau lebih rinci. Sebelum dipakai, aturan diuji pada 100 nama paket nyata. Laporan cakupan dicetak `make laporan-geocode` dan tampil di Tentang Data.

Tes wajib: "Pembangunan Pintu Air Gg. Nusa Indah RT 04 RW 22 Gading Sengon (JPO Kali Betik), Kel. Kelapa Gading Barat, Kec. Kelapa Gading" menghasilkan kelurahan Kelapa Gading Barat, RT 4, RW 22; nama yang menyebut dua kelurahan berbeda menghasilkan tingkat kecamatan atau kota; nama tanpa wilayah menghasilkan tingkat kota.

### ANG-03 Halaman paket

Halaman `/{bahasa}/paket/{kode}`. Inti visual adalah rantai horizontal (vertikal di bawah 1024 px) dengan mata rantai:

| Mata rantai | Isi | Status bila data tidak ada |
| --- | --- | --- |
| Rencana (RUP) | Kode RUP, pagu, sumber dana, tautan SiRUP | Putus: "rencana tidak ditemukan" |
| Pemilihan | Metode, HPS, jumlah peserta, tautan SPSE | Selalu ada (sumber paket) |
| Pemenang | Penyedia, nilai kontrak | Menunggu: "belum ada pemenang" |
| Pelaksanaan | Jadwal mulai dan selesai | Estimasi atau tidak tersedia |
| Lokasi | Kelurahan atau kecamatan, RT/RW, peta mini | Tingkat kota |
| Serah terima | Tanggal serah terima (bila diisi instansi) | Estimasi |
| Masa pemeliharaan | Lihat ANG-04 | Estimasi |
| Bukti warga | Jumlah laporan per status, galeri | Kosong: ajakan melapor |

Bagian bawah: tema yang cocok, sinyal aktif beserta alasannya, hak jawab instansi, perbandingan paket sejenis (ANG-06), riwayat penyedia ringkas, aspirasi reses yang kemungkinan terkait (ANG-05), tombol Lapor Bukti, tombol Paket Aspirasi.

API: `GET /paket/{kode}`, `GET /kecamatan/{kode}/paket?tema=&tahun=&halaman=`, `GET /paket/terdekat?lat=&lng=&radius_m=`.

Kriteria terima: setiap mata rantai mengembalikan status ada, putus, menunggu, atau estimasi sesuai tabel.

### ANG-04 Estimasi masa pemeliharaan

Konfigurasi per kategori di `config/pemeliharaan.yaml`: `permanen` (180 hari) atau `semi_permanen` (90 hari), sesuai ketentuan masa pemeliharaan minimal pekerjaan konstruksi.

```text
akhir_pelaksanaan = tanggal_serah_terima bila diisi instansi
                    selain itu jadwal_selesai bila ada
                    selain itu tanggal_kontrak + durasi_bawaan_kategori
akhir_pemeliharaan = akhir_pelaksanaan + masa_pemeliharaan_kategori
```

Status: `pelaksanaan`, `pemeliharaan` (dengan hitung mundur hari), `pemeliharaan_berakhir`, `tidak_dapat_diestimasi`. Label `estimasi` kecuali tanggal serah terima diisi instansi lewat Panel Instansi, lalu label menjadi `resmi` dengan sumber instansi tersebut.

Kalimat warga: "Bila ada kerusakan sebelum {tanggal}, perbaikannya masih menjadi kewajiban penyedia (perkiraan)."

### ANG-05 Aspirasi ke Anggaran

Untuk setiap `aspirasi_reses` yang menyebut wilayah: cari paket di wilayah itu pada tahun anggaran berikutnya yang kategorinya atau kata kunci bidangnya cocok. Hasil ditampilkan sebagai "paket yang kemungkinan terkait", bukan klaim bahwa paket berasal dari aspirasi. Bila tidak ada kecocokan, tampil "belum ditemukan paket yang terkait di data LPSE". Menjadi pengganti data pokir yang tidak publik.

API: `GET /aspirasi-reses/{id}/paket-terkait`.

### ANG-06 Perbandingan paket sejenis dan sinyal

**Kategori.** Kategori dan sektor ditentukan dari nama paket dengan aturan regex di `config/kategori.yaml`, dievaluasi berurutan, kecocokan pertama menang:

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
- kategori: jalan
  sektor: jalan
  pola: ["(perbaikan|peningkatan|pemeliharaan|pembangunan) jalan", "trotoar", "jalur pedestrian"]
```

**Perbandingan.** Untuk satu paket, ambil paket lain dengan kategori, jenis pengadaan, dan tahun anggaran yang sama di seluruh DKI. Tampilkan strip plot nilai HPS dengan median, kuartil, dan posisi paket ini. Kalimat pendamping wajib: "Ukuran pekerjaan tiap paket bisa berbeda; perbandingan ini menunjukkan posisi nilai, bukan kewajaran." Tidak ditampilkan bila pembanding kurang dari `min_pembanding` (bawaan 8).

**Sinyal "perlu ditanyakan".** Semua ambang di `config/sinyal.yaml`, dihitung ulang job `turunan`. Tampilan selalu: judul netral, alasan dengan angka, data pembanding, tautan sumber, dan ruang hak jawab.

| ID | Judul tampilan | Aturan | Ambang bawaan |
| --- | --- | --- | --- |
| S-01 | Hanya satu peserta | Metode tender atau seleksi, jumlah peserta = 1 | Tidak berlaku untuk pengadaan langsung dan penunjukan langsung |
| S-02 | Nilai jauh di atas paket sejenis | HPS > Q3 + 1,5 × IQR kelompok pembanding | Minimal 8 pembanding |
| S-03 | Beberapa paket sejenis bernilai mepet di bawah batas pengadaan langsung | ≥ 3 paket pengadaan langsung, instansi sama, kategori sama, dalam 30 hari, masing-masing ≥ 85% ambang | Ambang dari GD-09 |
| S-04 | Satu penyedia memenangkan sebagian besar paket sejenis | Penyedia sama memenangkan ≥ 40% nilai paket satu kategori di satu instansi dalam satu tahun anggaran, minimal 5 paket | |
| S-05 | Melewati jadwal | Hari ini > `jadwal_selesai` + 30 hari dan ada ≥ 1 laporan terverifikasi berkondisi belum_dimulai, sedang_dikerjakan, atau terhenti | |

**Kosakata terlarang** untuk teks buatan Rantai (berkas i18n, pesan backend, ringkasan umpan, jawaban Tanya Rantai): korupsi, curang, mark up, penggelembungan, kongkalikong, fiktif, penyelewengan, menyimpang, memperjuangkan. Judul berita media ditampilkan apa adanya sebagai kutipan berlencana dan tidak termasuk teks buatan Rantai. Tes otomatis memindai berkas i18n dan pesan backend.

API: `GET /paket/{kode}/pembanding`.

### ANG-07 Profil penyedia

Halaman `/{bahasa}/penyedia/{id}`: nama, jumlah dan total nilai paket per tahun (grafik batang dengan tabel), kategori utama, instansi pemberi paket terbanyak, daftar paket di Rantai, riwayat Opentender 2014 sampai 2023 dengan atribusi ODbL. Tanpa skor atau peringkat.

API: `GET /penyedia/{id}`.

### ANG-08 Lapor Bukti dan verifikasi

**Alur seluler** (layar sempit adalah target utama):

1. Warga masuk (OTP) bila belum.
2. Izin lokasi; akurasi harus ≤ 100 m, bila tidak ditolak dengan saran pindah ke ruang terbuka.
3. Daftar proyek dalam radius 500 m, diurutkan jarak; warga memilih satu, atau "proyek tidak ada di daftar" (laporan tanpa paket, tetap tersimpan per kelurahan).
4. Ambil atau unggah 1 sampai 4 foto (JPEG, PNG, HEIC; maksimal 8 MB per foto).
5. Pilih kondisi dan isi catatan opsional (maksimal 500 karakter).
6. Ringkasan dan kirim.

Pemrosesan server: verifikasi magic bytes, dekode lalu enkode ulang ke JPEG (membuang EXIF), batasi sisi terpanjang 2048 px, thumbnail 480 px, simpan ke MinIO dengan kunci acak. `lokasi_asli` privat; `lokasi_publik` dibulatkan ke grid 50 m. Batas: 5 laporan per akun per hari.

**Mesin status** di `domain`:

```text
DIKIRIM → MENUNGGU_SAKSI → TERVERIFIKASI → DITERUSKAN → DITANGGAPI → SELESAI
                └────────────→ DITOLAK (spam, duplikat, di luar lingkup)
```

- `DIKIRIM` ke `MENUNGGU_SAKSI` otomatis setelah lolos pemrosesan foto.
- `TERVERIFIKASI` bila minimal 3 akun berbeda (termasuk pelapor) mengonfirmasi kondisi yang sama dalam 200 m dan 14 hari, atau moderator atau petugas instansi terkait menandai terverifikasi.
- `DITERUSKAN` otomatis saat terverifikasi ke antrean Panel Instansi milik instansi pelaksana paket, atau instansi dari `tanggung_jawab` tema yang cocok untuk laporan tanpa paket.
- `DITANGGAPI` saat instansi memberi hak jawab. `SELESAI` bila instansi menandai selesai dan tidak ada sanggahan terverifikasi dalam 14 hari.
- Setiap perpindahan dicatat di `laporan_riwayat` dan tampil sebagai linimasa publik.

Tes wajib: foto hasil tanpa EXIF; lokasi publik terbulatkan; semua transisi sah; transisi tidak sah ditolak dengan `ErrTransisiTidakSah`; konfirmasi di luar 200 m atau lewat 14 hari tidak dihitung.

API: `POST /laporan` (multipart), `GET /laporan/{id}`, `POST /laporan/{id}/konfirmasi`.

### ANG-09 Unduh data

`GET /unduh/paket.csv?kode_wilayah=&tahun=`: kolom paket utama dan kolom provenance (URL sumber, tanggal ambil). Keterangan lisensi dan atribusi dicantumkan di halaman unduh dan Tentang Data.

## 9. Aksi

### AKS-01 Tanya Rantai

Antarmuka obrolan di `/{bahasa}/tanya` dan panel geser di halaman kecamatan (kecamatan aktif ikut sebagai konteks).

Backend memanggil API LLM dengan tool use. Model dari env `TANYA_MODEL`, kunci dari env `TANYA_API_KEY`. Model hanya boleh menjawab dari hasil fungsi berikut:

| Fungsi | Masukan | Keluaran |
| --- | --- | --- |
| `cari_wilayah` | teks | daftar wilayah |
| `ringkasan_kecamatan` | kode | angka kunci tema, anggaran, wakilku |
| `indikator_tema` | kode kecamatan, tema | kebutuhan, kapasitas, kesenjangan, rincian |
| `sekolah_di_kecamatan` | kode | daftar SMP dan daya tampung |
| `daftar_paket` | kode kecamatan, tema, tahun | paket ringkas |
| `detail_paket` | kode | rantai paket |
| `belanja_tema` | tema, tahun | kesenjangan dan belanja per kecamatan |
| `tanggung_jawab_tema` | tema | instansi, komisi, dasar hukum |
| `wakil_dapil` | kode kecamatan, tema opsional | anggota dan komisi |
| `umpan_anggota` | id anggota, tag opsional | item umpan ringkas dengan jenis sumber |

Setiap hasil fungsi menyertakan `provenance_id`. Aturan prompt sistem:

- Jawab dalam bahasa pertanyaan, kalimat pendek, istilah awam.
- Setiap angka wajib bersumber dari hasil fungsi dan ditandai rujukan `[n]`.
- Bila data tidak ada, katakan tidak tersedia dan sarankan halaman terkait.
- Tidak memberi penilaian, opini, atau rekomendasi memilih terhadap politisi dan partai; jawab dengan fakta jabatan, umpan berlencana, dan jalur aspirasi.
- Sinyal dijelaskan sebagai hal yang perlu ditanyakan, dengan kosakata netral.

Validator pasca-jawaban di server: setiap angka di teks jawaban harus muncul di hasil fungsi pada percakapan itu; bila gagal, jawaban diganti pesan aman dan dicatat. Kosakata terlarang ANG-06 juga diperiksa.

Streaming ke klien via SSE. Batas: 20 pertanyaan per IP per jam, panjang pertanyaan 500 karakter. Mode demo memakai cache jawaban untuk pertanyaan di `data/snapshot/tanya_demo.json`.

API: `POST /tanya` (SSE) dengan `{pertanyaan, kode_kecamatan_konteks, bahasa}`.

### AKS-02 Paket Aspirasi

Warga memilih kecamatan, tema, dan masalah dari `tanggung_jawab` tema. Rantai menyusun dokumen dari `templat_aspirasi`:

1. Judul masalah dan kecamatan.
2. Data pendukung dari indikator tema dan Arus Anggaran, setiap angka dengan sumber dan tanggal.
3. Instansi berwenang dan dasar hukumnya.
4. Komisi DPRD pengawas dan anggota dapil yang duduk di komisi itu.
5. Pertanyaan yang disarankan, netral dan spesifik.
6. Kanal yang tepat (AKS-03) dan ruang kosong untuk nama dan tanda tangan warga.

Keluaran: teks siap salin dan PDF (dibuat di server dengan pustaka PDF Go). Rantai tidak mengirim dokumen atas nama warga.

API: `POST /aspirasi` dengan `{kode_kecamatan, tema, masalah, bahasa}`, respons `{teks, pdf_url}`.

Kriteria terima: setiap angka di PDF punya sumber dan tanggal.

### AKS-03 Jalur aspirasi terdekat

Untuk kelurahan dan kecamatan: Lembaga Musyawarah Kelurahan (peran, cara menghubungi lewat kantor kelurahan), Dewan Kota (satu wakil per kecamatan), kanal CRM dan JAKI untuk aduan layanan, serta penjelasan masa reses DPRD. Isi dari `config/jalur_aspirasi.yaml`, diverifikasi saat kurasi.

API: `GET /jalur-aspirasi?kode_wilayah=`.

### AKS-04 Data yang Belum Dibuka

Halaman `/{bahasa}/data-belum-dibuka` dari `config/data_belum_dibuka.yaml`:

| Data | Kenapa tidak tersedia | Pengganti di Rantai |
| --- | --- | --- |
| Suara per anggota | Keputusan DPRD diambil per fraksi atau musyawarah | Sikap Fraksi, WAK-05 |
| Pokir per anggota | Diinput ke SIPD, sistem internal pemerintah | Aspirasi ke Anggaran, ANG-05 |
| Kehadiran per anggota | Hanya jumlah hadir yang muncul di berita | Permohonan informasi dari halaman ini |
| Laporan reses per anggota | Hanya kompilasi agregat per bidang | Aspirasi Reses, WAK-06 |

Setiap baris punya tombol "Minta data ini" yang menyusun surat permohonan informasi ke PPID DPRD DKI berdasarkan UU 14/2008 tentang Keterbukaan Informasi Publik, sebagai teks siap salin dan PDF. Rantai tidak mengirim permohonan atas nama warga.

API: `GET /data-belum-dibuka`, `POST /data-belum-dibuka/{kode}/permohonan` dengan `{bahasa}`.

### AKS-05 Bagikan sebagai gambar

Kartu ringkasan kecamatan 1080 × 1350 dirender server (SVG ke PNG) dengan angka kunci, label, sumber, dan tautan pendek.

## 10. Fondasi

### DSR-01 Mode demo luring

- `make seed-demo` memuat `data/snapshot/` (dump terkompresi) ke database lokal.
- `RANTAI_MODE=demo` mematikan job pengambilan data, menampilkan pita "Mode demo, data per {tanggal snapshot}", dan memakai cache jawaban Tanya Rantai serta ringkasan umpan yang sudah tersimpan.
- Tile peta untuk wilayah DKI ikut di-bundle (PMTiles) agar peta tetap tampil tanpa internet.

### DSR-02 Akun dan peran

- Warga masuk dengan OTP 6 digit ke email. Kode berlaku 10 menit, maksimal 5 percobaan, maksimal 3 permintaan per email per jam.
- Token akses JWT 15 menit di memori; token refresh 7 hari di cookie `httpOnly`, `Secure`, `SameSite=Lax`, berotasi setiap pemakaian.
- Peran `instansi` (terikat `instansi_id`) dan `anggota` (terikat `anggota_id`) dibuat admin. Peran `moderator` dan `admin` dibuat lewat perintah CLI di worker.
- Untuk demo cukup akun contoh dari `skenario_demo.yaml`.
- Membaca seluruh data publik tidak memerlukan akun.

API: `POST /auth/otp/kirim`, `POST /auth/otp/verifikasi`, `POST /auth/refresh`, `POST /auth/keluar`, `GET /auth/saya`.

### DSR-03 Panel Instansi dan panel moderasi

Panel di `/{bahasa}/panel` untuk peran instansi dan anggota:

- Antrean instansi: sinyal aktif pada paket milik instansinya, laporan berstatus DITERUSKAN, paket yang mendekati akhir masa pemeliharaan.
- Antrean anggota: item umpan yang dipautkan kepadanya.
- Aksi: tulis hak jawab (maksimal 2.000 karakter, boleh melampirkan tautan dokumen), tandai laporan selesai, isi tanggal serah terima paket.
- Hak jawab tampil publik berdampingan dengan sorotannya, dengan nama penjawab dan tanggal; tidak bisa dihapus.

Panel moderasi: antrean laporan MENUNGGU_SAKSI serta laporan atau item umpan yang diadukan; aksi tandai terverifikasi, tolak dengan alasan, sembunyikan foto atau item umpan. Semua aksi tercatat.

API: `GET /panel/antrean`, `POST /hak-jawab`, `PATCH /panel/paket/{kode}/serah-terima`, `PATCH /panel/laporan/{id}/selesai`, `GET /moderasi/antrean`, `PATCH /moderasi/laporan/{id}`, `PATCH /moderasi/umpan/{id}`.

### DSR-04 Bahasa dan tema tampilan

- Rute berawalan bahasa `/id/...` dan `/en/...`; nama wilayah, instansi, paket, dan judul berita tidak diterjemahkan. Seluruh teks antarmuka lewat `src/lib/i18n` sejak awal; terjemahan Inggris dan pengalih bahasa menyusul sebagai P2.
- Tema terang dan gelap mengikuti token `docs/DESIGN.md`; bawaan mengikuti sistem; pengalih menyusul sebagai P2.
- Aksesibilitas WCAG 2.1 AA tidak termasuk P2: fokus terlihat, kontras 4,5:1, peta dan grafik punya ringkasan teks dan tabel.

## 11. Ringkasan endpoint

| Metode | Path | Akses |
| --- | --- | --- |
| GET | `/kesehatan`, `/provenance/{id}`, `/sumber-data` | publik |
| GET | `/wilayah/cari`, `/wilayah/titik`, `/wilayah/{kode}`, `/kecamatan/{kode}/ringkasan` | publik |
| GET | `/tema`, `/tema/{kode}`, `/tema/{kode}/belanja`, `/peta/tema/{kode}` | publik |
| GET | `/tema/pendidikan/kecamatan/{kode}`, `/tema/pendidikan/kecamatan/{kode}/sekolah`, `/sekolah/terdekat`, `/peta/sekolah` | publik |
| GET | `/kecamatan/{kode}/wakil`, `/dapil/{id}`, `/anggota/{id}`, `/anggota/{id}/umpan`, `/anggota/{id}/isu`, `/fraksi/{nama}/sikap`, `/aspirasi-reses` | publik |
| GET | `/kecamatan/{kode}/paket`, `/paket/{kode}`, `/paket/{kode}/pembanding`, `/paket/terdekat`, `/penyedia/{id}`, `/aspirasi-reses/{id}/paket-terkait`, `/peta/paket`, `/peta/laporan`, `/unduh/paket.csv` | publik |
| GET | `/jalur-aspirasi`, `/data-belum-dibuka` | publik |
| POST | `/tanya`, `/aspirasi`, `/data-belum-dibuka/{kode}/permohonan` | publik dengan batas laju |
| POST | `/auth/otp/kirim`, `/auth/otp/verifikasi`, `/auth/refresh`, `/auth/keluar` | publik dengan batas laju |
| GET | `/auth/saya`, `/notifikasi` | warga |
| POST | `/laporan`, `/laporan/{id}/konfirmasi`, `/langganan` | warga |
| GET, POST, PATCH | `/panel/*`, `/hak-jawab` | instansi, anggota |
| GET, PATCH | `/moderasi/*` | moderator |
| POST, GET | `/admin/ingest/{job}`, `/admin/ingest/runs` | admin |

## 12. Keamanan

- Header: CSP tanpa `unsafe-inline` untuk skrip, HSTS di balik TLS, `X-Content-Type-Options`, `Referrer-Policy: strict-origin-when-cross-origin`, `Permissions-Policy` yang hanya membuka geolocation dan camera untuk halaman Lapor. API memakai CSP `default-src 'none'`.
- Batas laju per IP dan per akun di middleware; respons 429 dengan `Retry-After`. Di belakang reverse proxy, alamat IP hanya diambil dari proxy tepercaya.
- CORS hanya untuk origin frontend dari konfigurasi atau env `CORS_ORIGIN`.
- CSRF: token refresh di cookie `SameSite=Lax`, endpoint yang mengubah data memerlukan header `Authorization` dari token akses di memori.
- Unggahan: validasi magic bytes, batas ukuran, dekode dan enkode ulang, kunci objek acak, bucket privat dengan URL bertanda tangan berumur pendek untuk foto yang belum terverifikasi.
- Query terparameterisasi; tidak ada SQL dari input pengguna, termasuk dari Tanya Rantai.
- Log tanpa email dan koordinat asli.

## 13. Alur demo dan data uji

Skenario demo dikunci di `data/snapshot/skenario_demo.yaml`:

- `kecamatan_demo`: kecamatan di Jakarta Utara dengan rasio kursi SMP negeri terendah yang datanya lengkap, dipilih setelah PEN-01 terhitung.
- `paket_simulasi`: satu paket kategori `ruang_kelas` di kecamatan demo, penyedia fiktif "CV Contoh Simulasi", `is_simulation = true`.
- `laporan_simulasi`: dua laporan dan konfirmasi yang cukup untuk status TERVERIFIKASI, seluruhnya `is_simulation = true`, foto ilustrasi bebas lisensi.
- `akun_demo`: satu warga, satu petugas instansi Dinas Pendidikan, satu akun anggota contoh yang tidak ditautkan ke anggota sungguhan, satu moderator.
- `tanya_demo`: tiga pertanyaan dan jawaban tercache, termasuk "Kenapa SMP negeri di kecamatanku penuh?".

Data simulasi tidak pernah ditautkan ke paket, penyedia, anggota, atau warga sungguhan (dites).

Tes Playwright `e2e/alur-demo.spec.ts` menjalankan Alur Demo PRD bagian 11 langkah 1 sampai 8 dalam mode demo, pada lebar 390 px dan 1440 px.

## 14. Bentuk `docs/DESIGN.md`

Dokumen desain wajib memuat bagian berikut agar bisa dipakai langsung oleh frontend:

1. Token warna (terang dan gelap), termasuk warna semantik label kepastian, sektor, status laporan, jenis sumber umpan (rilis resmi, media), dan skala choropleth.
2. Token tipografi, jarak, radius, bayangan, durasi animasi.
3. Konfigurasi `tailwind.config` dari token tersebut.
4. Daftar komponen dengan varian dan status.
5. Tata letak tiap layar di empat lebar (360, 768, 1024, 1440): beranda, halaman kecamatan dengan tiga bagian, detail Pendidikan, peta, halaman paket, profil anggota dengan umpan dua sumber, Kesenjangan vs Belanja, Tanya Rantai, Paket Aspirasi, Data yang Belum Dibuka, Tentang Data, panel.
6. Teks antarmuka utama dalam bahasa Indonesia dan Inggris.
