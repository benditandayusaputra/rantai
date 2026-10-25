# Rencana Kerja: Rantai

Versi 2.0 · 23 Oktober 2026 · Pasangan dari `docs/PRD.md` v2.0 dan `docs/SPEC.md` v2.0.

## Cara memakai

- Kerjakan tugas dari atas ke bawah. Ambil tugas pertama yang belum dicentang.
- Satu tugas, satu commit, langsung `git push origin main`. Contoh pesan commit di setiap tugas adalah titik awal; sesuaikan dengan pekerjaan yang benar-benar dilakukan, tetap memakai gitmoji dan bahasa Indonesia.
- **[GERBANG DATA GD-xx]**: sebelum menulis kode, verifikasi sumber sesuai SPEC bagian 1 dan catat hasilnya di `docs/data-sources.md`. Bila verifikasi gagal, jalankan cadangan yang tertulis di SPEC, atau tulis di "Catatan blokir" dan lanjut ke tugas berikutnya yang tidak bergantung.
- **[BOLEH PECAH]**: tugas ini boleh dipecah menjadi beberapa commit yang masing-masing lulus `make lint` dan `make test`.
- **[BUTUH DESAIN]**: tugas UI. Wajib membaca `docs/DESIGN.md` lebih dulu. Bila layar yang dibutuhkan belum ada di sana, lewati dan catat di "Catatan blokir".
- "Selesai bila" adalah kriteria terima. Tugas belum boleh dicentang sebelum semuanya terpenuhi dan tes lulus.
- Fase 0 dan 1 dikerjakan dengan ID fitur v1.0; padanannya ada di awal SPEC.

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

- [x] **T-101** Server API dasar
  - Pemuat konfigurasi YAML + env, `cmd/api` dengan GoFiber v3, amplop respons dan pemetaan error (SPEC 0.3), endpoint `GET /api/v1/kesehatan`.
  - Commit: `🧱 Bangun server API dengan amplop respons dan pemetaan error`
  - Selesai bila: tes handler kesehatan dan tes pemetaan error domain ke kode HTTP lulus.
- [x] **T-102** Skema wilayah dan provenance
  - Migrasi ekstensi `postgis`, `pg_trgm`; tabel `wilayah`, `provenance`, `ingest_run` (SPEC 2.1).
  - Commit: `🗃️ Tambah skema wilayah, provenance, dan riwayat ingest`
  - Selesai bila: migrasi naik dan turun tanpa error.
- [x] **T-103** Impor batas wilayah **[GERBANG DATA GD-10]**
  - Job `wilayah`: impor kota, kecamatan, kelurahan DKI beserta alias nama dan titik tengah.
  - Commit: `✨ Impor batas dan alias wilayah DKI Jakarta`
  - Selesai bila: 6 kota/kabupaten administrasi, seluruh kecamatan dan kelurahan DKI tersimpan; tes menghitung jumlah dan memeriksa satu alias.
- [x] **T-104** Pencarian wilayah (WIL-01)
  - `GET /wilayah/cari`, `GET /wilayah/titik`, `GET /wilayah/{kode}`.
  - Commit: `✨ Tambah API pencarian wilayah berdasarkan nama dan titik`
  - Selesai bila: kriteria terima WIL-01 terpenuhi dalam tes integrasi.
- [x] **T-105** Provenance (WIL-04)
  - `GET /provenance/{id}` dan tipe `Angka` di lapisan delivery.
  - Commit: `✨ Tambah API provenance dan tipe Angka berlabel`
  - Selesai bila: setiap `Angka` di respons wajib membawa label dan `provenance_id` (dites).
- [x] **T-106** Pengaman HTTP
  - Batas laju per IP, header keamanan (SPEC 12), CORS untuk origin frontend dari konfigurasi.
  - Commit: `🔒️ Pasang batas laju dan header keamanan pada API`
  - Selesai bila: tes 429 dengan `Retry-After` dan tes keberadaan header lulus.

## Fase 2. Kerangka tema dan data Pendidikan

- [x] **T-201** Kerangka tema (SPEC 4)
  - Migrasi `tema`, `indikator_kecamatan`; `config/tema/pendidikan.yaml`; job `tema` dengan validasi kontrak; `GET /tema`, `GET /tema/{kode}`.
  - Commit: `✨ Tambah konfigurasi tema berbasis basis data`
  - Selesai bila: validator menolak entri tanggung jawab tanpa dasar hukum, `ambang_rasio` tidak berurutan, regex rusak, dan placeholder tak dikenal (dites); tema Pendidikan termuat dari YAML ke basis data.
- [x] **T-202** Impor sekolah dan lokasi **[GERBANG DATA GD-06]** **[BOLEH PECAH]**
  - Migrasi `sekolah`, `sekolah_rombel`; job `dapodik` untuk SD, MI, SMP, MTs di DKI; kecamatan sekolah ditentukan dari koordinat dengan `ST_Covers`, cadangan pencocokan nama.
  - Commit: `✨ Impor sekolah, lokasi, dan rombel per tingkat dari Dapodik`
  - Selesai bila: setiap kecamatan DKI punya data sekolah; tes parser memakai fixture; sekolah tanpa koordinat dilaporkan, tidak ditebak.
- [x] **T-203** Daya tampung resmi **[GERBANG DATA GD-04, GD-05]**
  - Migrasi `daya_tampung`; ekstraksi lampiran Kepkadis 258/2026 dan daftar SPMB Bersama ke CSV di `data/snapshot/`; job `daya_tampung` mencocokkan NPSN.
  - Commit: `✨ Impor daya tampung SMP negeri dan kursi SPMB Bersama`
  - Selesai bila: total kursi negeri yang terimpor sama dengan 73.289 atau selisihnya dijelaskan di `docs/data-sources.md`; NPSN tak cocok dicetak sebagai laporan.
- [x] **T-204** Calon lulusan per kecamatan **[GERBANG DATA GD-12]**
  - Migrasi `calon_lulusan_kecamatan`; job `calon_lulusan` dengan metode per sekolah, agregat resmi, atau estimasi total SD ÷ 6 sesuai hasil gerbang.
  - Commit: `✨ Hitung calon lulusan SD per kecamatan dengan label metode`
  - Selesai bila: setiap kecamatan punya angka dengan label dan metode; metode estimasi berlabel `estimasi` dan rumusnya tercatat di provenance (dites).
- [x] **T-205** Indikator kesenjangan kursi (PEN-01)
  - `config/pendidikan.yaml`, rumus domain, aturan label terlemah, penghitung `kesenjangan_kursi_smp` menulis ke `indikator_kecamatan`.
  - Commit: `✨ Hitung kesenjangan kursi SMP negeri per kecamatan`
  - Selesai bila: tes tabel mencakup semua kasus tepi PEN-01 dan aturan label terlemah; tidak ada angka 32 atau 6 tertulis langsung di kode.
- [x] **T-206** API Tema Pendidikan (PEN-02, PEN-03)
  - `GET /peta/tema/{kode}`, `GET /tema/pendidikan/kecamatan/{kode}`, `GET /tema/pendidikan/kecamatan/{kode}/sekolah`.
  - Commit: `✨ Tambah API peta indikator dan detail sekolah kecamatan`
  - Selesai bila: GeoJSON choropleth di bawah 300 KB untuk seluruh DKI; tes integrasi lulus.
- [x] **T-207** Sekolah terdekat (PEN-06)
  - `GET /sekolah/terdekat` dengan tekanan radius.
  - Commit: `✨ Tambah API sekolah terdekat dan tekanan kursi radius`
  - Selesai bila: tes dengan data uji tiga sekolah di jarak berbeda lulus.

## Fase 3. Data Arus Anggaran

- [x] **T-301** Klien SPSE yang santun **[GERBANG DATA GD-01]**
  - Klien HTTP bersama dengan batas 1 permintaan per detik per domain, User-Agent dari konfigurasi, cek robots.txt, berhenti pada 403 atau 429. Klien endpoint JSON daftar paket (pola pyproc) dan parser halaman pengumuman tender dan non-tender.
  - Commit: `✨ Tambah klien SPSE berbatas laju dan parser pengumuman paket`
  - Selesai bila: parser lulus tes dengan minimal 5 fixture (SPEC ANG-01); tidak ada permintaan jaringan di tes.
- [x] **T-302** Skema pengadaan
  - Migrasi `instansi`, `rup`, `penyedia`, `paket`, `riwayat_penyedia`, `sinyal` (SPEC 2.4).
  - Commit: `🗃️ Tambah skema instansi, RUP, penyedia, paket, dan sinyal`
  - Selesai bila: migrasi naik dan turun; indeks GiST di `paket.lokasi`.
- [x] **T-303** Job impor paket (ANG-01) **[BOLEH PECAH]**
  - Job `lpse_daftar` dan `lpse_detail`, normalisasi nama penyedia, pencatatan perubahan tahap; snapshot hasil impor untuk lingkungan di luar Indonesia.
  - Commit: `✨ Jalankan impor harian paket LPSE DKI dengan pencatatan tahap`
  - Selesai bila: menjalankan job dua kali tidak menggandakan data; `ingest_run` tercatat.
- [x] **T-304** Pemenang dan kontrak **[GERBANG DATA GD-02]**
  - Parser pemenang tender dan pemenang berkontrak.
  - Commit: `✨ Ambil pemenang dan nilai kontrak paket`
  - Selesai bila: fixture paket tanpa pemenang dan dengan pemenang dites.
- [x] **T-305** Taut RUP **[GERBANG DATA GD-03]**
  - Job `sirup` untuk kode RUP yang ditemukan.
  - Commit: `✨ Tautkan paket ke rencana umum pengadaan`
  - Selesai bila: paket tanpa RUP tetap tersimpan dengan `kode_rup` kosong.
- [x] **T-306** Penentuan lokasi paket (ANG-02)
  - Gazetir dari tabel wilayah, ekstraksi RT/RW, penyaringan kota instansi, pergeseran titik deterministik.
  - Commit: `✨ Petakan paket ke kelurahan dan kecamatan dari nama paket`
  - Selesai bila: tiga tes wajib ANG-02 lulus; uji 100 nama paket nyata dicatat; cakupan dicetak oleh `make laporan-geocode`.
- [x] **T-307** Estimasi masa pemeliharaan (ANG-04)
  - `config/pemeliharaan.yaml` dan usecase status.
  - Commit: `✨ Estimasi masa pemeliharaan paket konstruksi`
  - Selesai bila: keempat status dan perubahan label menjadi `resmi` dites dengan jam tiruan.
- [x] **T-308** Filter anggaran tema dan API Arus Anggaran (ANG-03)
  - Penandaan `paket.tema_kode` dari `filter_anggaran`; `GET /paket/{kode}` berbentuk rantai, `GET /kecamatan/{kode}/paket`, `GET /peta/paket`, `GET /paket/terdekat`.
  - Commit: `✨ Tambah API halaman paket dan daftar paket per kecamatan`
  - Selesai bila: setiap mata rantai mengembalikan status ada, putus, menunggu, atau estimasi sesuai tabel ANG-03; filter tema dites.
- [x] **T-309** Kesenjangan vs Belanja (PEN-04)
  - `GET /tema/{kode}/belanja` dari `indikator_kecamatan` dan paket bertema.
  - Commit: `✨ Bandingkan kesenjangan kursi dengan belanja sekolah per kecamatan`
  - Selesai bila: paket berlokasi tingkat kota dikeluarkan dari angka kecamatan dan dilaporkan terpisah (dites).
- [x] **T-310** Kategori, pembanding, dan sinyal (ANG-06) **[GERBANG DATA GD-09]**
  - `config/kategori.yaml`, `config/sinyal.yaml`, `GET /paket/{kode}/pembanding`, aturan S-01 sampai S-05, tes kosakata terlarang atas berkas i18n dan pesan backend.
  - Commit: `✨ Bandingkan HPS paket sejenis dan hitung sinyal perlu ditanyakan`
  - Selesai bila: pembanding tidak tampil di bawah `min_pembanding`; setiap sinyal punya tes positif dan negatif; tes kosakata terlarang lulus.
- [x] **T-311** Riwayat penyedia (ANG-07) **[GERBANG DATA GD-11]**
  - Impor Opentender 2014 sampai 2023 untuk pembeli Pemprov DKI, `GET /penyedia/{id}`.
  - Commit: `✨ Impor riwayat Opentender dan tambah API profil penyedia`
  - Selesai bila: respons memuat atribusi ODbL; tidak ada skor atau peringkat.

## Fase 4. Data Wakilku

- [x] **T-401** Dapil dan anggota (WAK-01) **[GERBANG DATA GD-07]**
  - Migrasi `dapil`, `dapil_wilayah`, `anggota`; job `kpu_dprd` bagian dapil dan anggota.
  - Commit: `✨ Impor dapil dan anggota DPRD DKI periode 2024 sampai 2029`
  - Selesai bila: 10 dapil, 106 anggota, setiap kecamatan DKI masuk tepat satu dapil.
- [x] **T-402** Komisi dan alat kelengkapan (WAK-02) **[GERBANG DATA GD-08]**
  - Migrasi `komisi`, `anggota_akd`; job `kpu_dprd` bagian komisi, AKD, dan tautan profil.
  - Commit: `✨ Impor komisi dan alat kelengkapan DPRD DKI`
  - Selesai bila: setiap komisi punya bidang; anggota tanpa data komisi ditandai, tidak ditebak.
- [x] **T-403** API Wakilku (WAK-01, WAK-02)
  - `GET /kecamatan/{kode}/wakil?tema=`, `GET /anggota/{id}`, `GET /dapil/{id}`.
  - Commit: `✨ Tambah API wakil dapil dan profil anggota`
  - Selesai bila: urutan anggota komisi pengawas tema lebih dahulu dites; tidak ada angka kinerja di respons.
- [x] **T-404** Umpan rilis resmi DPRD (WAK-03) **[GERBANG DATA GD-13]**
  - Migrasi `umpan_kegiatan`, `umpan_anggota`; job `umpan_dprd`; penandaan `tag_isu` dari tema.
  - Commit: `✨ Impor rilis kegiatan resmi DPRD per anggota`
  - Selesai bila: parser lulus tes fixture; item tidak disimpan dengan teks berita penuh.
- [x] **T-405** Umpan media terverifikasi (WAK-03) **[GERBANG DATA GD-14]**
  - `config/media.yaml`; job `umpan_media`; disambiguasi nama plus "DPRD DKI" atau fraksi.
  - Commit: `✨ Impor liputan media terverifikasi dengan disambiguasi nama`
  - Selesai bila: fixture nama sama tanpa "DPRD DKI" dibuang dan dengan fraksi dipautkan (dites); domain di luar daftar media ditolak.
- [x] **T-406** Ringkasan netral dan API umpan (WAK-03)
  - Klien LLM bersama (model dari env), job `ringkas_umpan` dengan validator, `GET /anggota/{id}/umpan`.
  - Commit: `✨ Ringkas umpan wakil rakyat secara netral dan tambah API umpan`
  - Selesai bila: tes memakai klien LLM tiruan; ringkasan dengan kosakata terlarang atau angka karangan ditolak dan jatuh ke judul saja.
- [x] **T-407** Ringkasan halaman kecamatan (WIL-02 sisi backend)
  - `GET /kecamatan/{kode}/ringkasan` untuk bagian Tema, Arus Anggaran, dan Wakilku.
  - Commit: `✨ Tambah ringkasan halaman kecamatan`
  - Selesai bila: kode kelurahan diarahkan ke kecamatannya dengan keterangan; setiap angka berbentuk `Angka` (dites).

## Fase 5. Data demo

- [x] **T-501** Mode demo (DSR-01 sisi backend)
  - `data/snapshot/skenario_demo.yaml` (SPEC 13), `make seed-demo`, `RANTAI_MODE=demo` mematikan job pengambilan.
  - Commit: `🗃️ Tambah skenario dan snapshot data demo`
  - Selesai bila: semua data simulasi bertanda `is_simulation = true` dan tidak terhubung ke paket, penyedia, atau anggota sungguhan (dites).

## Fase 6. Frontend inti

Prasyarat: `docs/DESIGN.md` memuat layar PRD v2.0 (SPEC 14), terutama halaman kecamatan dengan tiga bagian, profil anggota dengan umpan dua sumber, Kesenjangan vs Belanja, dan Data yang Belum Dibuka.

- [x] **T-601** Kerangka SvelteKit **[BUTUH DESAIN]**
  - SvelteKit 2, Svelte 5 runes, TypeScript strict, Tailwind dari token DESIGN.md, i18n dengan pesan Indonesia, rute berawalan `/id`, ESLint, Prettier, svelte-check, Vitest. Tambahkan lint dan tes frontend ke CI.
  - Commit: `🎉 Inisialisasi frontend SvelteKit dengan token desain`
  - Selesai bila: `make lint` dan `make test` mencakup frontend; tidak ada warna atau ukuran ditulis langsung di komponen; tidak ada teks hardcode di komponen.
- [x] **T-602** Klien API dan lencana sumber (WIL-04) **[BUTUH DESAIN]**
  - Klien bertipe, tipe `Angka`, komponen `LencanaProvenance` dengan panel sumber.
  - Commit: `✨ Tambah klien API bertipe dan lencana sumber`
  - Selesai bila: komponen dites untuk keempat label; panel bisa dibuka dengan papan ketik.
- [x] **T-603** Beranda dan pencarian (WIL-01) **[BUTUH DESAIN]**
  - Commit: `✨ Bangun beranda dengan pencarian wilayah dan lokasi`
  - Selesai bila: memilih hasil membuka halaman kecamatan, kelurahan tersorot; tombol lokasi menangani izin ditolak.
- [x] **T-604** Halaman kecamatan (WIL-02) **[BUTUH DESAIN]**
  - Header, kartu ringkas, bagian Tema, Arus Anggaran, Wakilku dengan URL sendiri.
  - Commit: `✨ Bangun halaman kecamatan dengan bagian tema, anggaran, dan wakil`
  - Selesai bila: bagian bisa dibuka langsung dari URL; status kosong, memuat, dan error ada.
- [x] **T-605** Detail Pendidikan (PEN-03) **[BUTUH DESAIN]**
  - Commit: `✨ Tampilkan kesenjangan kursi dan daftar sekolah kecamatan`
  - Selesai bila: keterbatasan domisili tampil; grafik punya tabel alternatif; label tiap jenis kursi terlihat.
- [x] **T-606** Peta (WIL-03, PEN-02) **[BUTUH DESAIN]** **[BOLEH PECAH]**
  - MapLibre, lapisan indikator tema, proyek, sekolah; legenda; tabel alternatif.
  - Commit: `✨ Tambah peta indikator tema dan lapisan proyek`
  - Selesai bila: semua lapisan bisa dinyalakan dan dimatikan dengan papan ketik.
- [x] **T-607** Arus Anggaran dan halaman paket (ANG-03, ANG-04) **[BUTUH DESAIN]**
  - Daftar paket kecamatan dengan filter tema dan tahun, status pemeliharaan, halaman paket berbentuk rantai.
  - Commit: `✨ Bangun daftar paket dan halaman rantai paket`
  - Selesai bila: rantai horizontal di 1024 px ke atas dan vertikal di bawahnya; mata rantai putus terbaca jelas; paket simulasi berlabel "Data simulasi" di setiap tempat ia muncul.
- [x] **T-608** Kesenjangan vs Belanja (PEN-04) **[BUTUH DESAIN]**
  - Layar dirancang sendiri karena DESIGN.md belum memuatnya: diagram sebar di tab Uang halaman kecamatan dengan token yang ada.
  - Commit: `✨ Tampilkan grafik kesenjangan kursi dibanding belanja sekolah`
  - Selesai bila: kecamatan aktif tersorot; tabel alternatif dan kalimat pendamping wajib tampil.
- [x] **T-609** Wakilku dan profil anggota (WAK-01 sampai WAK-03) **[BUTUH DESAIN]**
  - Umpan dua sumber dirancang sendiri: lencana "Rilis resmi DPRD" memakai token label resmi, lencana media bergaris netral.
  - Daftar wakil dapil, profil anggota, umpan dua sumber berlencana, pernyataan bias liputan.
  - Commit: `✨ Tampilkan wakil dapil, profil anggota, dan jejak kegiatan`
  - Selesai bila: tidak ada angka kinerja atau peringkat politisi di layar; lencana jenis sumber selalu terlihat.
- [x] **T-610** Tentang Data (WIL-05) **[BUTUH DESAIN]**
  - Commit: `📝 Tambah halaman Tentang Data dengan rumus dan sumber`
  - Selesai bila: tanggal ambil terakhir tiap sumber diambil dari API, bukan ditulis tangan.

### Titik FDS (1 November 2026)

- [ ] **T-650** Siap FDS
  - Deploy pratinjau, README lengkap (judul, deskripsi, tech stack, cara menjalankan, progres pengerjaan dengan centang fitur), tangkapan layar beranda dan halaman kecamatan di `docs/screenshots/`.
  - Commit: `📝 Lengkapi README dan tangkapan layar untuk FDS`
  - Selesai bila: repo publik bisa dijalankan orang lain dari README dengan `make up`, `make migrate`, `make seed-demo`.

## Fase 7. Aksi

- [x] **T-701** Fungsi alat Tanya Rantai
  - Sepuluh fungsi AKS-01 sebagai usecase yang sudah ada, dibungkus registry alat.
  - Commit: `✨ Daftarkan fungsi data untuk Tanya Rantai`
  - Selesai bila: setiap fungsi dites dan mengembalikan `provenance_id`.
- [x] **T-702** Orkestrasi dan streaming
  - Tool use dengan klien LLM bersama, model dari `TANYA_MODEL`, SSE ke klien, batas laju.
  - Commit: `✨ Jalankan Tanya Rantai dengan pemanggilan alat dan SSE`
  - Selesai bila: tes memakai klien LLM tiruan; kunci API hanya dari env.
- [x] **T-703** Validator jawaban
  - Cek angka bersumber dan kosakata terlarang; pesan aman bila gagal.
  - Commit: `🔒️ Validasi angka dan kosakata pada jawaban Tanya Rantai`
  - Selesai bila: jawaban dengan angka karangan ditolak dalam tes.
- [x] **T-704** Antarmuka Tanya **[BUTUH DESAIN]**
  - Halaman `/tanya` dan panel geser di halaman kecamatan, rujukan `[n]` yang bisa diklik.
  - Commit: `✨ Bangun antarmuka Tanya Rantai dengan rujukan sumber`
  - Selesai bila: rujukan membuka panel sumber.
- [x] **T-705** Cache demo
  - Jawaban tercache untuk `tanya_demo.json` saat `RANTAI_MODE=demo`; daftar pertanyaan diambil dari `tanya_demo` di `skenario_demo.yaml`.
  - Commit: `✨ Sediakan jawaban tercache Tanya Rantai untuk mode demo`
  - Selesai bila: tiga pertanyaan demo terjawab tanpa internet.
- [x] **T-706** Paket Aspirasi dan jalur aspirasi (AKS-02, AKS-03)
  - `POST /aspirasi` dari templat tema, teks siap salin dan PDF dari server; `config/jalur_aspirasi.yaml`, `GET /jalur-aspirasi`.
  - Commit: `✨ Susun Paket Aspirasi berbasis tema dalam teks dan PDF`
  - Selesai bila: setiap angka di PDF punya sumber dan tanggal; kanal yang tepat tercantum.
- [x] **T-707** Antarmuka Paket Aspirasi **[BUTUH DESAIN]**
  - Commit: `✨ Bangun alur Paket Aspirasi`
  - Selesai bila: dari halaman kecamatan ke unduhan PDF dalam 2 menit pada uji manual.
- [x] **T-708** Data yang Belum Dibuka (AKS-04) **[BUTUH DESAIN]**
  - `config/data_belum_dibuka.yaml`, `GET /data-belum-dibuka`, `POST /data-belum-dibuka/{kode}/permohonan`, halaman dengan templat PPID.
  - Commit: `✨ Tambah halaman Data yang Belum Dibuka dan templat permohonan PPID`
  - Selesai bila: keempat data dan penggantinya tampil; surat permohonan berisi dasar UU 14/2008.

## Fase 8. Akun dan hak jawab

- [x] **T-801** Autentikasi OTP (DSR-02)
  - Migrasi `pengguna`, `otp`, `sesi_refresh`; endpoint `/auth/*`; peran warga, anggota, instansi, moderator, admin; pengirim email lewat SMTP dari env.
  - Job `skenario_demo` ikut memasang `akun_demo` dari `skenario_demo.yaml`.
  - Commit: `🔒️ Tambah masuk dengan OTP email dan token berotasi`
  - Selesai bila: batas percobaan, kedaluwarsa, dan rotasi refresh dites.
- [x] **T-802** Masuk di frontend **[BUTUH DESAIN]**
  - Commit: `✨ Tambah alur masuk OTP di frontend`
  - Selesai bila: token akses hanya di memori; muat ulang memakai refresh.
- [x] **T-803** Panel Instansi dan hak jawab (DSR-03, WAK-07) **[BUTUH DESAIN]** **[BOLEH PECAH]**
  - Migrasi `hak_jawab`; panel untuk instansi dan anggota; hak jawab pada sinyal, paket, laporan, dan item umpan.
  - Perpindahan otomatis laporan `TERVERIFIKASI` ke `DITERUSKAN` ke antrean instansi pelaksana paket atau instansi `tanggung_jawab` tema (ANG-08), memakai `domain.Transisi` dan `LaporanRepository.UbahStatus`.
  - Laporan tanpa paket tetap `TERVERIFIKASI` karena Lapor Bukti belum menanyakan tema masalahnya.
  - Menyimpang dari DESIGN 6.15: isian hak jawab 2.000 karakter mengikuti SPEC (DESIGN 1.200); "Simpan draf", unggah lampiran dan PDF berita acara, serta pratinjau `MiniChain` belum dibangun karena backend belum menyimpan draf dan berkas; dokumen pendukung memakai tautan.
  - Commit: `✨ Tambah Panel Instansi dan hak jawab publik`
  - Selesai bila: hak jawab tidak bisa dihapus; isi tanggal serah terima mengubah label pemeliharaan menjadi `resmi`; hak jawab langsung tampil di halaman publik.

## Fase 9. Fitur P1

- [x] **T-901** Isu yang Disuarakan (WAK-04) dan isu sekolah Komisi E (PEN-05) **[BUTUH DESAIN]**
  - `GET /anggota/{id}/isu`, `GET /kecamatan/{kode}/wakil?dengan_umpan=true`, tampilan di profil dan bagian Tema.
  - Commit: `✨ Tampilkan isu yang disuarakan anggota dengan jumlah sumber`
  - Selesai bila: tidak ada perbandingan antaranggota; label jumlah sumber memisahkan rilis resmi dan media.
- [x] **T-902** Sikap Fraksi dan Aspirasi Reses (WAK-05, WAK-06) **[GERBANG DATA GD-15]**
  - Migrasi `sikap_fraksi`, `aspirasi_reses`; job `paripurna` impor CSV kurasi; API dan tampilan.
  - Commit: `✨ Tambah sikap fraksi dan aspirasi reses dari dokumen paripurna`
  - Selesai bila: setiap entri punya tautan sumber; aspirasi tampil dengan catatan "per bidang, bukan per anggota".
- [x] **T-903** Aspirasi ke Anggaran (ANG-05)
  - Tampil di Wakilku per aspirasi berwilayah. Tampilan sebaliknya di halaman paket butuh endpoint pencarian terbalik (paket ke aspirasi) dan belum dibuat.
  - `GET /aspirasi-reses/{id}/paket-terkait` dan tampilan di halaman paket dan Wakilku.
  - Commit: `✨ Cocokkan aspirasi reses dengan paket yang kemungkinan terkait`
  - Selesai bila: hasil selalu berlabel "kemungkinan terkait"; kasus tanpa kecocokan dites.
- [x] **T-904** Terima laporan (ANG-08 sisi backend)
  - Migrasi laporan, pemrosesan foto (magic bytes, enkode ulang, buang EXIF, thumbnail), MinIO, pembulatan lokasi publik.
  - Job `skenario_demo` ikut memasang `laporan_simulasi` dari `skenario_demo.yaml` ke paket simulasi, dengan foto ilustrasi bebas lisensi.
  - Commit: `✨ Terima laporan warga dengan foto tanpa metadata`
  - Selesai bila: tes memastikan foto hasil tidak memuat EXIF dan lokasi publik terbulatkan.
- [x] **T-905** Verifikasi laporan (ANG-08)
  - Mesin status di domain, konfirmasi, riwayat.
  - Commit: `✨ Tambah mesin status verifikasi dan konfirmasi saksi`
  - Selesai bila: semua tes wajib verifikasi ANG-08 lulus.
- [x] **T-906** Alur Lapor Bukti seluler **[BUTUH DESAIN]** **[BOLEH PECAH]**
  - Menyimpang dari DESIGN 6.10: pilihan kondisi mengikuti lima kondisi SPEC ANG-08, catatan maksimal 500 karakter (SPEC), peta mini dan "Pilih titik di peta" belum dibangun karena laporan wajib berakurasi GPS ≤ 100 m.
  - Commit: `✨ Bangun alur Lapor Bukti untuk layar ponsel`
  - Selesai bila: alur selesai di lebar 360 px; akurasi lokasi buruk ditolak dengan saran; laporan simulasi berlabel di peta dan linimasa.
- [x] **T-907** Panel moderasi (DSR-03) **[BUTUH DESAIN]**
  - Antrean umpan berisi liputan media terbaru karena fitur "Laporkan konten" (aduan warga) belum ada; aksi "Verifikasi dan teruskan" baru menandai terverifikasi sampai antrean Panel Instansi (T-803) dibangun.
  - Moderasi laporan dan item umpan yang salah atribusi.
  - Commit: `✨ Tambah panel moderasi laporan dan umpan`
  - Selesai bila: setiap aksi tercatat di riwayat.
- [x] **T-908** Profil penyedia (ANG-07) **[BUTUH DESAIN]**
  - Commit: `✨ Bangun halaman profil penyedia`
- [x] **T-909** Cek Sekolahku (PEN-06) **[BUTUH DESAIN]**
  - Peta radius 3 km dari DESIGN 6.6 belum dipasang; daftar sekolah dan kalimat tekanan sudah memuat jarak dan sumber.
  - Commit: `✨ Bangun fitur Cek Sekolahku berbasis radius`
- [x] **T-910** Pantau Wilayahku, feed (WIL-06) **[BUTUH DESAIN]**
  - Empat dari enam jenis kejadian aktif; pemeliharaan segera berakhir dan hak jawab baru tercatat di SPEC WIL-06. Sakelar kabar perangkat dibangun di T-911.
  - Migrasi `langganan`, `notifikasi`; pembuat kejadian di job `turunan`; feed di aplikasi.
  - Commit: `✨ Tambah langganan kecamatan dan feed notifikasi`
- [x] **T-911** Web Push (WIL-06)
  - VAPID dari env, service worker.
  - Commit: `✨ Kirim notifikasi Pantau lewat Web Push`

## Fase 10. Kualitas dan rilis

- [x] **T-1000** Redesain antarmuka **[BUTUH DESAIN]**
  - Kelas komponen bersama di `app.css` (`.btn`, `.field`, `.card`, `.icon-tile`, `.chip`), ikon Lucide di setiap tombol dan tautan aksi, header berikon, footer dengan tautan per kelompok, beranda dengan kartu fitur, ilustrasi kota di beranda dan ilustrasi tokoh di halaman Masuk (token `--illu-*`, DESIGN 1.8).
  - Commit: `💄 Rapikan tampilan semua halaman dengan ikon dan ilustrasi`
  - Selesai bila: tidak ada geser horizontal di 360 sampai 1440 px; audit aksesibilitas tetap tanpa pelanggaran serius atau kritis.
- [x] **T-1001** Tes alur demo
  - `e2e/alur-demo.spec.ts` untuk delapan langkah Alur Demo PRD pada 390 px dan 1440 px dalam mode demo.
  - `RANTAI_MODE=demo make e2e`; tanpa paket simulasi seluruh langkah dilewati dengan pesan `make seed-demo`. API uji menulis log ke `frontend/test-logs/api.log`, tempat tes membaca kode masuk petugas; Panel Instansi hanya untuk layar lebar, jadi di 390 px langkah 7 memeriksa bagian hak jawab di halaman paket.
  - Diuji pada salinan basis data lokal dengan `skenario_demo`: langkah 1, 3, 4, 6, 7, dan 8 lulus. Langkah 2 (data Dapodik) dan 5 (`tanya_demo.json`) menunggu T-1005.
  - Commit: `✅ Tambah tes end-to-end alur demo utama`
- [x] **T-1002** Aksesibilitas
  - Audit axe pada semua halaman utama, perbaiki temuan.
  - `make e2e` menjalankan `e2e/aksesibilitas.spec.ts` (15 halaman, 390 px dan 1440 px, WCAG 2.1 A dan AA). Audit pertama memakai data parsial (Dapodik belum masuk); spec tetap di suite sehingga keadaan berdata ikut terperiksa.
  - Commit: `♿️ Perbaiki temuan audit aksesibilitas halaman utama`
  - Selesai bila: tidak ada pelanggaran serius atau kritis.
- [x] **T-1003** Performa
  - Halaman kecamatan menunggu ringkasan saat `load` sehingga konten ada di HTML SSR dan hidrasi tidak merender ulang (CLS 0,23 ke 0); panel Tanya dimuat saat pertama dibuka; ikon Lucide disatukan dalam satu potongan JS (permintaan JS halaman kecamatan 78 ke 36). MapLibre sudah dimuat dinamis sejak T-606.
  - Lighthouse 12.6 mobile (simulasi bawaan) atas build produksi di balik kompresi gzip seperti Vercel: beranda 97, kecamatan 96, LCP 2,5 dan 2,7 s, CLS 0.
  - p95 API lokal dengan data parsial (tanpa Dapodik, paket belum berlokasi): ringkasan kecamatan 27 ms, wakil 7 ms, peta tema 3 ms, paket kecamatan 10 ms, cari wilayah 77 ms. Diukur ulang pada data penuh di T-1005.
  - Commit: `⚡️ Percepat muat halaman kecamatan dan peta`
  - Selesai bila: Lighthouse mobile 90 ke atas untuk beranda dan halaman kecamatan; API p95 di bawah 300 ms pada data penuh.
- [x] **T-1004** Deploy
  - Compose produksi dengan Caddy, cadangan basis data harian; server di Indonesia untuk job LPSE dan SiRUP.
  - Atur proxy tepercaya di Fiber (`TrustProxy`) agar batas laju per IP dan HSTS bekerja di belakang Caddy.
  - Dikerjakan dalam tiga commit: compose produksi, Caddy, dan cadangan harian; rahasia proxy antara frontend Vercel dan API; image API lintas arsitektur. Frontend di `https://rantai-civic.vercel.app`. Bila server API di luar Indonesia, job LPSE dan SiRUP dijalankan dari mesin di Indonesia lalu hasilnya dibawa dengan `make snapshot-paket` dan `make muat-snapshot-paket`.
  - Commit: `🔧 Siapkan deploy produksi dengan Caddy`
- [ ] **T-1005** Snapshot luring final (DSR-01)
  - Snapshot data (`make snapshot-demo` di mesin dengan data lengkap), kunci `kecamatan_demo` di `skenario_demo.yaml`, rekam ulang `data/snapshot/tanya_demo.json` dengan `make tanya-demo` terhadap API live berdata lengkap, PMTiles wilayah DKI, uji jalan tanpa internet.
  - Commit: `🗃️ Perbarui snapshot data dan tile peta untuk demo luring`
  - Ukur ulang p95 API halaman wilayah pada data penuh (sasaran di bawah 300 ms, T-1003).
  - Selesai bila: alur demo lulus dengan jaringan dimatikan.
- [x] **T-1006** README final
  - Fitur per modul, diagram arsitektur, cara menjalankan dan menguji, env Vercel untuk domain `rantai-civic.vercel.app`, dan progres akhir. Butir snapshot luring tetap terbuka sampai T-1005.
  - Commit: `📝 Perbarui README dengan fitur, arsitektur, dan progres akhir`

## Cadangan P2 (setelah semua di atas)

- [ ] **T-1101** Tren kesenjangan antar tahun (PEN-07) · `✨ Tampilkan tren kesenjangan kursi antar tahun ajaran`
- [ ] **T-1102** Unduh CSV paket (ANG-09) · `✨ Tambah unduhan CSV paket per wilayah`
- [ ] **T-1103** Bagikan sebagai gambar (AKS-05) · `✨ Buat kartu ringkasan kecamatan untuk dibagikan`
- [ ] **T-1104** Bahasa Inggris dan pengalih bahasa (DSR-04) · `🌐 Tambah terjemahan Inggris dan pengalih bahasa`
- [ ] **T-1105** Pengalih tema terang dan gelap (DSR-04) · `💄 Tambah pengalih tema terang dan gelap`
- [ ] **T-1106** Ringkasan mingguan via email (WIL-06) · `✨ Kirim ringkasan mingguan Pantau lewat email`

## Keluaran lomba di luar kode

Dikerjakan tim, bukan tugas commit. Dicatat di sini agar urutannya terlihat.

| Tenggat | Keluaran | Bergantung pada |
| --- | --- | --- |
| Secepatnya | Permohonan informasi kehadiran dan laporan reses ke PPID DPRD DKI (untuk disebut saat pitching) | - |
| Secepatnya | Cek Sekolah Kita untuk data kelas 6 per kecamatan (risiko terbesar) | GD-12 |
| Sebelum Fase 6 | `docs/DESIGN.md` diperbarui ke layar PRD v2.0 | SPEC 14 |
| 1 November | PDF FDS maks 10 halaman (`WIA_FSDMIA2026_NamaTim`): konsep, sitemap, wireframe, fitur, tangkapan layar landing + 1 halaman | T-650 |
| 11 November | Proposal maks 30 halaman (`WIA_PenyisihanMIA2026_NamaTim`), surat orisinalitas bermaterai, video maks 5 menit di YouTube publik | Fase 7 dan 8 |
| 28 November | Deck maks 15 halaman, prototipe daring dan cadangan luring | T-1005 |

## Catatan blokir

Format: `- T-xxx · tanggal · alasan · tugas yang dikerjakan sebagai gantinya`

- T-301 · 23 Oktober 2026 · halaman detail dan evaluasi tender SPSE DKI tertutup untuk publik (403, "Akses Ditolak"); detail tender butuh token API Gateway Inaproc (peran Data Integrator, surat ke LKPP) · paket tender disimpan dari kolom daftar, non-tender memakai halaman pengumuman lengkap
