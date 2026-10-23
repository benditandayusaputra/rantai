# Rencana Kerja: Rantai

Versi 1.0 · 23 Oktober 2026 · Pasangan dari `docs/PRD.md` dan `docs/SPEC.md`.

## Cara memakai

- Kerjakan tugas dari atas ke bawah. Ambil tugas pertama yang belum dicentang.
- Satu tugas, satu commit, langsung `git push origin main`. Contoh pesan commit di setiap tugas adalah titik awal; sesuaikan dengan pekerjaan yang benar-benar dilakukan, tetap memakai gitmoji dan bahasa Indonesia.
- **[GERBANG DATA GD-xx]**: sebelum menulis kode, verifikasi sumber sesuai SPEC bagian 1 dan catat hasilnya di `docs/data-sources.md`. Bila verifikasi gagal, jalankan cadangan yang tertulis di SPEC, atau tulis di "Catatan blokir" dan lanjut ke tugas berikutnya yang tidak bergantung.
- **[BOLEH PECAH]**: tugas ini boleh dipecah menjadi beberapa commit yang masing-masing lulus `make lint` dan `make test`.
- **[BUTUH DESAIN]**: tugas UI. Wajib membaca `docs/DESIGN.md` lebih dulu. Bila berkas itu belum ada, lewati dan catat di "Catatan blokir".
- "Selesai bila" adalah kriteria terima. Tugas belum boleh dicentang sebelum semuanya terpenuhi dan tes lulus.

## Fase 0. Fondasi repo

- [x] **T-001** Inisialisasi monorepo
  - Struktur folder sesuai CLAUDE.md bagian 3, `go mod init`, README kerangka (judul, deskripsi, tech stack, cara menjalankan, progres pengerjaan), `.env.example`.
  - Commit: `🎉 Inisialisasi struktur monorepo Rantai`
  - Selesai bila: `go build ./...` lulus di folder backend; README memuat lima bagian wajib FDS walau isinya masih ringkas.
- [x] **T-002** Lingkungan lokal
  - `docker-compose.yml` untuk PostgreSQL 16 + PostGIS 3 dan MinIO; `Makefile` dengan target `up`, `down`, `migrate`, `lint`, `test`.
  - Commit: `🔧 Tambah Docker Compose PostGIS dan MinIO serta Makefile`
  - Selesai bila: `make up` lalu `make migrate` berjalan dari repo bersih.
- [x] **T-003** Linter dan CI
  - Konfigurasi golangci-lint (termasuk aturan panjang fungsi dan larangan komentar selain direktif), GitHub Actions untuk lint dan tes backend.
  - Commit: `🔧 Atur golangci-lint dan alur CI lint serta tes`
  - Selesai bila: CI hijau di main.

## Fase 1. Kerangka backend

- [ ] **T-101** Server API dasar
  - Pemuat konfigurasi YAML + env, `cmd/api` dengan GoFiber v3, amplop respons dan pemetaan error (SPEC 0.3), endpoint `GET /api/v1/kesehatan`.
  - Commit: `🧱 Bangun server API dengan amplop respons dan pemetaan error`
  - Selesai bila: tes handler kesehatan dan tes pemetaan error domain ke kode HTTP lulus.
- [ ] **T-102** Skema wilayah dan provenance
  - Migrasi ekstensi `postgis`, `pg_trgm`; tabel `wilayah`, `provenance`, `ingest_run` (SPEC 2.1).
  - Commit: `🗃️ Tambah skema wilayah, provenance, dan riwayat ingest`
  - Selesai bila: migrasi naik dan turun tanpa error.
- [ ] **T-103** Impor batas wilayah **[GERBANG DATA GD-10]**
  - Job `wilayah`: impor kota, kecamatan, kelurahan DKI beserta alias nama dan titik tengah.
  - Commit: `✨ Impor batas dan alias wilayah DKI Jakarta`
  - Selesai bila: 6 kota/kabupaten administrasi, seluruh kecamatan dan kelurahan DKI tersimpan; tes menghitung jumlah dan memeriksa satu alias.
- [ ] **T-104** Pencarian wilayah (DASAR-01)
  - `GET /wilayah/cari`, `GET /wilayah/titik`, `GET /wilayah/{kode}`.
  - Commit: `✨ Tambah API pencarian wilayah berdasarkan nama dan titik`
  - Selesai bila: kriteria terima DASAR-01 terpenuhi dalam tes integrasi.
- [ ] **T-105** Provenance (DASAR-02)
  - `GET /provenance/{id}` dan tipe `Angka` di lapisan delivery.
  - Commit: `✨ Tambah API provenance dan tipe Angka berlabel`
  - Selesai bila: setiap `Angka` di respons wajib membawa label dan `provenance_id` (dites).
- [ ] **T-106** Pengaman HTTP
  - Batas laju per IP, header keamanan (SPEC 10), CORS untuk origin frontend dari konfigurasi.
  - Commit: `🔒️ Pasang batas laju dan header keamanan pada API`
  - Selesai bila: tes 429 dengan `Retry-After` dan tes keberadaan header lulus.

## Fase 2. Data Lensa Layanan

- [ ] **T-201** Impor sekolah dan rombel **[GERBANG DATA GD-04, GD-05]** **[BOLEH PECAH]**
  - Migrasi `sekolah`, `sekolah_rombel`, `daya_tampung_resmi`; job `dapodik` untuk SD, MI, SMP, MTs, SMA, SMK, MA di DKI.
  - Commit: `✨ Impor data sekolah dan rombel per tingkat dari Dapodik`
  - Selesai bila: setiap kecamatan DKI punya data sekolah; tes parser memakai fixture.
- [ ] **T-202** Rumus kesenjangan kursi (LAY-01)
  - `config/layanan.yaml`, usecase perhitungan, materialized view `kesenjangan_kecamatan`.
  - Commit: `✨ Hitung kursi proksi dan kesenjangan lulusan per kecamatan`
  - Selesai bila: tes tabel mencakup semua kasus tepi LAY-01; tidak ada angka 32 atau 36 tertulis langsung di kode.
- [ ] **T-203** API Lensa Layanan (LAY-02, LAY-03)
  - `GET /wilayah/{kode}/layanan`, `GET /kecamatan/{kode}/sekolah`, `GET /peta/kesenjangan`.
  - Commit: `✨ Tambah API kesenjangan kursi dan daftar sekolah kecamatan`
  - Selesai bila: GeoJSON choropleth di bawah 300 KB untuk seluruh DKI; tes integrasi lulus.
- [ ] **T-204** Sekolah terdekat (LAY-04)
  - `GET /sekolah/terdekat` dengan tekanan radius.
  - Commit: `✨ Tambah API sekolah terdekat dan tekanan kursi radius`
  - Selesai bila: tes dengan data uji tiga sekolah di jarak berbeda lulus.
- [ ] **T-205** Keterbukaan daya tampung (LAY-05) **[GERBANG DATA GD-06]**
  - Job `juknis_spmb` impor CSV, `GET /keterbukaan/daya-tampung`.
  - Commit: `✨ Tambah rapor keterbukaan daya tampung sekolah negeri`
  - Selesai bila: status `diumumkan` dan `tidak_ditemukan` dites.

## Fase 3. Data Lensa Uang

- [ ] **T-301** Klien SPSE yang santun **[GERBANG DATA GD-01]**
  - Klien HTTP dengan batas 1 permintaan per detik per domain, User-Agent dari konfigurasi, cek robots.txt, berhenti pada 403 atau 429. Parser halaman pengumuman tender dan non-tender.
  - Commit: `✨ Tambah klien SPSE berbatas laju dan parser pengumuman paket`
  - Selesai bila: parser lulus tes dengan minimal 5 fixture HTML (SPEC UANG-01); tidak ada permintaan jaringan di tes.
- [ ] **T-302** Skema pengadaan
  - Migrasi `instansi`, `rup`, `penyedia`, `paket`, `riwayat_penyedia`, `sinyal` (SPEC 2.3).
  - Commit: `🗃️ Tambah skema instansi, RUP, penyedia, paket, dan sinyal`
  - Selesai bila: migrasi naik dan turun; indeks GiST di `paket.lokasi`.
- [ ] **T-303** Job impor paket (UANG-01) **[BOLEH PECAH]**
  - Job `lpse_daftar` dan `lpse_detail`, normalisasi nama penyedia, pencatatan perubahan tahap.
  - Commit: `✨ Jalankan impor harian paket LPSE DKI dengan pencatatan tahap`
  - Selesai bila: menjalankan job dua kali tidak menggandakan data; `ingest_run` tercatat.
- [ ] **T-304** Pemenang dan kontrak **[GERBANG DATA GD-02]**
  - Parser pemenang tender dan pemenang berkontrak.
  - Commit: `✨ Ambil pemenang dan nilai kontrak paket`
  - Selesai bila: fixture paket tanpa pemenang dan dengan pemenang dites.
- [ ] **T-305** Taut RUP **[GERBANG DATA GD-03]**
  - Job `sirup` untuk kode RUP yang ditemukan.
  - Commit: `✨ Tautkan paket ke rencana umum pengadaan`
  - Selesai bila: paket tanpa RUP tetap tersimpan dengan `kode_rup` kosong.
- [ ] **T-306** Geocoding nama paket (UANG-02)
  - Gazetir dari tabel wilayah, ekstraksi RT/RW, penyaringan kota instansi, pergeseran titik deterministik.
  - Commit: `✨ Petakan paket ke kelurahan dan kecamatan dari nama paket`
  - Selesai bila: tiga tes wajib UANG-02 lulus; cakupan dicetak oleh perintah `make laporan-geocode`.
- [ ] **T-307** Kategori dan pembanding (UANG-04)
  - `config/kategori.yaml`, penentuan kategori dan sektor, `GET /paket/{kode}/pembanding`.
  - Commit: `✨ Kategorikan paket dan bandingkan HPS paket sejenis`
  - Selesai bila: pembanding tidak tampil di bawah `min_pembanding`; tes kuartil lulus.
- [ ] **T-308** Estimasi masa pemeliharaan (UANG-06)
  - `config/pemeliharaan.yaml` dan usecase status.
  - Commit: `✨ Estimasi masa pemeliharaan paket konstruksi`
  - Selesai bila: keempat status dan perubahan label menjadi `resmi` dites dengan jam tiruan.
- [ ] **T-309** Sinyal pola (UANG-05) **[GERBANG DATA GD-09]**
  - `config/sinyal.yaml`, aturan S-01 sampai S-05, tes kosakata terlarang atas berkas i18n dan pesan backend.
  - Commit: `✨ Hitung sinyal pola pengadaan dengan ambang terkonfigurasi`
  - Selesai bila: setiap sinyal punya tes positif dan negatif; tes kosakata terlarang lulus.
- [ ] **T-310** API Lensa Uang (UANG-03)
  - `GET /paket/{kode}` berbentuk rantai, `GET /wilayah/{kode}/uang`, `GET /peta/paket`, `GET /paket/terdekat`.
  - Commit: `✨ Tambah API Rantai Rupiah dan daftar paket per wilayah`
  - Selesai bila: setiap mata rantai mengembalikan status ada, putus, menunggu, atau estimasi sesuai tabel UANG-03.
- [ ] **T-311** Riwayat penyedia (UANG-07) **[GERBANG DATA GD-11]**
  - Impor Opentender 2014 sampai 2023 untuk pembeli Pemprov DKI, `GET /penyedia/{id}`.
  - Commit: `✨ Impor riwayat Opentender dan tambah API profil penyedia`
  - Selesai bila: respons memuat atribusi ODbL; tidak ada skor atau peringkat.

## Fase 4. Data Lensa Tanggung Jawab

- [ ] **T-401** Dapil dan anggota **[GERBANG DATA GD-07]**
  - Migrasi `dapil`, `dapil_wilayah`, `anggota`; job `kpu_dprd`.
  - Commit: `✨ Impor dapil dan anggota DPRD DKI periode 2024 sampai 2029`
  - Selesai bila: 10 dapil, 106 anggota, setiap kecamatan DKI masuk tepat satu dapil.
- [ ] **T-402** Komisi **[GERBANG DATA GD-08]**
  - Migrasi `komisi`, keanggotaan komisi.
  - Commit: `✨ Impor komisi DPRD DKI dan keanggotaannya`
  - Selesai bila: setiap komisi punya bidang; anggota tanpa data komisi ditandai, tidak ditebak.
- [ ] **T-403** Peta kewenangan (TJ-01)
  - `config/kewenangan.yaml` dengan dasar hukum, pemuat dan validasi, `GET /kewenangan`, `GET /kewenangan/daftar`.
  - Commit: `✨ Tambah peta kewenangan masalah ke instansi dan komisi`
  - Selesai bila: pemuat menolak entri tanpa dasar hukum (dites); minimal 8 jenis masalah terisi.
- [ ] **T-404** API Tanggung Jawab (TJ-02, TJ-03, TJ-04)
  - `GET /wilayah/{kode}/tanggung-jawab`, `GET /anggota/{id}`, `GET /dapil/{id}`, konfigurasi jalur aspirasi.
  - Commit: `✨ Tambah API wakil dapil, profil anggota, dan jalur aspirasi`
  - Selesai bila: urutan anggota komisi pengawas lebih dahulu dites.
- [ ] **T-405** Ringkasan tiga lensa (DASAR-03)
  - `GET /wilayah/{kode}/ringkasan`.
  - Commit: `✨ Tambah ringkasan tiga lensa per wilayah`
  - Selesai bila: kelurahan memakai data Layanan kecamatannya dengan keterangan.

## Fase 5. Data demo

- [ ] **T-501** Mode demo (DASAR-07 sisi backend)
  - `data/snapshot/skenario_demo.yaml` (SPEC 11), `make seed-demo`, `RANTAI_MODE=demo` mematikan job pengambilan.
  - Commit: `🗃️ Tambah skenario dan snapshot data demo`
  - Selesai bila: semua data simulasi bertanda `is_simulation = true` dan tidak terhubung ke paket atau penyedia sungguhan (dites).

## Fase 6. Frontend inti

- [ ] **T-601** Kerangka SvelteKit **[BUTUH DESAIN]**
  - SvelteKit 2, Svelte 5 runes, TypeScript strict, Tailwind dari token DESIGN.md, tema terang dan gelap, ESLint, Prettier, svelte-check, Vitest. Tambahkan lint dan tes frontend ke CI.
  - Commit: `🎉 Inisialisasi frontend SvelteKit dengan token desain`
  - Selesai bila: `make lint` dan `make test` mencakup frontend; tidak ada warna atau ukuran ditulis langsung di komponen.
- [ ] **T-602** Dua bahasa (DASAR-06)
  - Rute `/id` dan `/en`, berkas pesan, pengalih bahasa.
  - Commit: `🌐 Tambah rute dua bahasa dan pengalih bahasa`
  - Selesai bila: tes memastikan setiap kunci ada di kedua bahasa.
- [ ] **T-603** Klien API dan lencana provenance (DASAR-02) **[BUTUH DESAIN]**
  - Klien bertipe, tipe `Angka`, komponen `LencanaProvenance` dengan panel sumber.
  - Commit: `✨ Tambah klien API bertipe dan lencana provenance`
  - Selesai bila: komponen dites untuk keempat label; panel bisa dibuka dengan papan ketik.
- [ ] **T-604** Beranda dan pencarian (DASAR-01) **[BUTUH DESAIN]**
  - Commit: `✨ Bangun beranda dengan pencarian wilayah dan lokasi`
  - Selesai bila: memilih hasil membuka Wilayahku; tombol lokasi menangani izin ditolak.
- [ ] **T-605** Wilayahku (DASAR-03) **[BUTUH DESAIN]**
  - Header, kartu ringkas tiga lensa, tab dengan URL sendiri.
  - Commit: `✨ Bangun halaman Wilayahku dengan tiga tab lensa`
  - Selesai bila: tab bisa dibuka langsung dari URL; status kosong, memuat, dan error ada.
- [ ] **T-606** Tab Layanan dan detail kecamatan (LAY-03) **[BUTUH DESAIN]**
  - Commit: `✨ Tampilkan kesenjangan kursi dan daftar sekolah kecamatan`
  - Selesai bila: keterbatasan domisili tampil; grafik punya tabel alternatif.
- [ ] **T-607** Peta (DASAR-04, LAY-02) **[BUTUH DESAIN]** **[BOLEH PECAH]**
  - MapLibre, lapisan kesenjangan, proyek, sekolah; legenda; tabel alternatif.
  - Commit: `✨ Tambah peta kesenjangan kursi dan lapisan proyek`
  - Selesai bila: semua lapisan bisa dinyalakan dan dimatikan dengan papan ketik.
- [ ] **T-608** Tab Uang **[BUTUH DESAIN]**
  - Daftar paket wilayah dengan filter sektor dan tahun, chip sinyal, status pemeliharaan.
  - Commit: `✨ Tampilkan daftar paket dan status pemeliharaan per wilayah`
  - Selesai bila: paket simulasi berlabel "Data simulasi" di setiap tempat ia muncul.
- [ ] **T-609** Rantai Rupiah (UANG-03) **[BUTUH DESAIN]**
  - Halaman paket dengan rantai, sinyal, pembanding, hak jawab.
  - Commit: `✨ Bangun halaman Rantai Rupiah untuk detail paket`
  - Selesai bila: rantai horizontal di 1024 px ke atas dan vertikal di bawahnya; mata rantai putus terbaca jelas.
- [ ] **T-610** Tab Tanggung Jawab dan profil anggota (TJ-01 sampai TJ-04) **[BUTUH DESAIN]**
  - Commit: `✨ Tampilkan instansi berwenang, komisi, dan wakil dapil`
  - Selesai bila: tidak ada angka kinerja atau peringkat politisi di layar.
- [ ] **T-611** Tentang Data (DASAR-08) **[BUTUH DESAIN]**
  - Commit: `📝 Tambah halaman Tentang Data dengan rumus dan sumber`
  - Selesai bila: tanggal ambil terakhir tiap sumber diambil dari API, bukan ditulis tangan.

### Titik FDS (1 November 2026)

- [ ] **T-650** Siap FDS
  - Deploy pratinjau, README lengkap (judul, deskripsi, tech stack, cara menjalankan, progres pengerjaan dengan centang fitur), tangkapan layar beranda dan Wilayahku di `docs/screenshots/`.
  - Commit: `📝 Lengkapi README dan tangkapan layar untuk FDS`
  - Selesai bila: repo publik bisa dijalankan orang lain dari README dengan `make up`, `make migrate`, `make seed-demo`.

## Fase 7. Akun dan laporan warga

- [ ] **T-701** Autentikasi OTP (DASAR-05)
  - Migrasi `pengguna`, `otp`, `sesi_refresh`; endpoint `/auth/*`; pengirim email lewat SMTP dari env.
  - Commit: `🔒️ Tambah masuk dengan OTP email dan token berotasi`
  - Selesai bila: batas percobaan, kedaluwarsa, dan rotasi refresh dites.
- [ ] **T-702** Masuk di frontend **[BUTUH DESAIN]**
  - Commit: `✨ Tambah alur masuk OTP di frontend`
  - Selesai bila: token akses hanya di memori; muat ulang memakai refresh.
- [ ] **T-703** Terima laporan (UANG-08 sisi backend)
  - Migrasi laporan, pemrosesan foto (magic bytes, enkode ulang, buang EXIF, thumbnail), MinIO, pembulatan lokasi publik.
  - Commit: `✨ Terima laporan warga dengan foto tanpa metadata`
  - Selesai bila: tes memastikan foto hasil tidak memuat EXIF dan lokasi publik terbulatkan.
- [ ] **T-704** Verifikasi laporan (UANG-09)
  - Mesin status di domain, konfirmasi, riwayat.
  - Commit: `✨ Tambah mesin status verifikasi dan konfirmasi saksi`
  - Selesai bila: semua tes wajib UANG-09 lulus.
- [ ] **T-705** Alur Lapor Bukti seluler **[BUTUH DESAIN]** **[BOLEH PECAH]**
  - Commit: `✨ Bangun alur Lapor Bukti untuk layar ponsel`
  - Selesai bila: alur selesai di lebar 360 px; akurasi lokasi buruk ditolak dengan saran.
- [ ] **T-706** Linimasa dan lapisan laporan **[BUTUH DESAIN]**
  - Commit: `✨ Tampilkan linimasa status laporan dan lapisan peta laporan`
  - Selesai bila: laporan simulasi berlabel di peta dan linimasa.
- [ ] **T-707** Moderasi (LINTAS-04) **[BUTUH DESAIN]**
  - Commit: `✨ Tambah panel moderasi laporan`
  - Selesai bila: setiap aksi tercatat di riwayat.
- [ ] **T-708** Panel Instansi dan hak jawab (LINTAS-02) **[BUTUH DESAIN]** **[BOLEH PECAH]**
  - Commit: `✨ Tambah Panel Instansi dan hak jawab publik`
  - Selesai bila: hak jawab tidak bisa dihapus; isi tanggal serah terima mengubah label pemeliharaan menjadi `resmi`.

## Fase 8. Tanya Rantai

- [ ] **T-801** Fungsi alat
  - Sepuluh fungsi LINTAS-01 sebagai usecase yang sudah ada, dibungkus registry alat.
  - Commit: `✨ Daftarkan fungsi data untuk Tanya Rantai`
  - Selesai bila: setiap fungsi dites dan mengembalikan `provenance_id`.
- [ ] **T-802** Orkestrasi dan streaming
  - Klien LLM dengan tool use, model dari `TANYA_MODEL`, SSE ke klien, batas laju.
  - Commit: `✨ Jalankan Tanya Rantai dengan pemanggilan alat dan SSE`
  - Selesai bila: tes memakai klien LLM tiruan; kunci API hanya dari env.
- [ ] **T-803** Validator jawaban
  - Cek angka bersumber dan kosakata terlarang; pesan aman bila gagal.
  - Commit: `🔒️ Validasi angka dan kosakata pada jawaban Tanya Rantai`
  - Selesai bila: jawaban dengan angka karangan ditolak dalam tes.
- [ ] **T-804** Antarmuka Tanya **[BUTUH DESAIN]**
  - Halaman `/tanya` dan panel geser di Wilayahku, rujukan `[n]` yang bisa diklik.
  - Commit: `✨ Bangun antarmuka Tanya Rantai dengan rujukan sumber`
  - Selesai bila: rujukan membuka panel provenance.
- [ ] **T-805** Cache demo
  - Jawaban tercache untuk `tanya_demo.json` saat `RANTAI_MODE=demo`.
  - Commit: `✨ Sediakan jawaban tercache Tanya Rantai untuk mode demo`
  - Selesai bila: tiga pertanyaan demo terjawab tanpa internet.

## Fase 9. Fitur pelengkap

- [ ] **T-901** Paket Aspirasi (TJ-05) **[BUTUH DESAIN]**
  - `POST /aspirasi`, teks siap salin dan PDF dari server.
  - Commit: `✨ Susun Paket Aspirasi dalam teks dan PDF`
  - Selesai bila: setiap angka di PDF punya sumber dan tanggal.
- [ ] **T-902** Cek Sekolahku (LAY-04) **[BUTUH DESAIN]**
  - Commit: `✨ Bangun fitur Cek Sekolahku berbasis radius`
- [ ] **T-903** Profil penyedia (UANG-07) **[BUTUH DESAIN]**
  - Commit: `✨ Bangun halaman profil penyedia`
- [ ] **T-904** Pantau Wilayahku, feed (LINTAS-03) **[BUTUH DESAIN]**
  - Migrasi `langganan`, `notifikasi`; pembuat kejadian di job `turunan`; feed di aplikasi.
  - Commit: `✨ Tambah langganan wilayah dan feed notifikasi`
- [ ] **T-905** Web Push (LINTAS-03)
  - VAPID dari env, service worker.
  - Commit: `✨ Kirim notifikasi Pantau lewat Web Push`
- [ ] **T-906** Rapor keterbukaan (LAY-05) **[BUTUH DESAIN]**
  - Commit: `✨ Tampilkan rapor keterbukaan daya tampung`

## Fase 10. Kualitas dan rilis

- [ ] **T-1001** Tes alur demo
  - `e2e/alur-demo.spec.ts` pada 390 px dan 1440 px dalam mode demo.
  - Commit: `✅ Tambah tes end-to-end alur demo utama`
- [ ] **T-1002** Aksesibilitas
  - Audit axe pada semua halaman utama, perbaiki temuan.
  - Commit: `♿️ Perbaiki temuan audit aksesibilitas halaman utama`
  - Selesai bila: tidak ada pelanggaran serius atau kritis.
- [ ] **T-1003** Performa
  - Commit: `⚡️ Percepat muat halaman Wilayahku dan peta`
  - Selesai bila: Lighthouse mobile 90 ke atas untuk Beranda dan Wilayahku; API p95 di bawah 300 ms pada data penuh.
- [ ] **T-1004** Deploy
  - Compose produksi dengan Caddy, cadangan basis data harian.
  - Commit: `🔧 Siapkan deploy produksi dengan Caddy`
- [ ] **T-1005** Snapshot luring final (DASAR-07)
  - Snapshot data, PMTiles wilayah DKI, uji jalan tanpa internet.
  - Commit: `🗃️ Perbarui snapshot data dan tile peta untuk demo luring`
  - Selesai bila: alur demo lulus dengan jaringan dimatikan.
- [ ] **T-1006** README final
  - Commit: `📝 Perbarui README dengan fitur, arsitektur, dan progres akhir`

## Cadangan P2 (setelah semua di atas)

- [ ] **T-1101** Tren rasio kursi (LAY-06) · `✨ Tampilkan tren rasio kursi antar tahun ajaran`
- [ ] **T-1102** Unduh CSV paket (UANG-10) · `✨ Tambah unduhan CSV paket per wilayah`
- [ ] **T-1103** Bagikan sebagai gambar (LINTAS-05) · `✨ Buat kartu ringkasan wilayah untuk dibagikan`
- [ ] **T-1104** Ringkasan mingguan via email (LINTAS-03) · `✨ Kirim ringkasan mingguan Pantau lewat email`

## Keluaran lomba di luar kode

Dikerjakan tim, bukan tugas commit. Dicatat di sini agar urutannya terlihat.

| Tenggat | Keluaran | Bergantung pada |
| --- | --- | --- |
| 1 November | PDF FDS maks 10 halaman (`WIA_FSDMIA2026_NamaTim`): konsep, sitemap, wireframe, fitur, tangkapan layar landing + 1 halaman | T-650 |
| 11 November | Proposal maks 30 halaman (`WIA_PenyisihanMIA2026_NamaTim`), surat orisinalitas bermaterai, video maks 5 menit di YouTube publik | Fase 7 dan 8 |
| 28 November | Deck maks 15 halaman, prototipe daring dan cadangan luring | T-1005 |

## Catatan blokir

Format: `- T-xxx · tanggal · alasan · tugas yang dikerjakan sebagai gantinya`

(belum ada)
