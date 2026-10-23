# PRD: Rantai

Versi 1.0 · 23 Oktober 2026

## 1. Ringkasan

Rantai adalah web transparansi publik yang menjawab tiga pertanyaan warga tentang wilayahnya: apakah layanan dasarnya cukup, uang publiknya ke mana, dan siapa yang bertanggung jawab. Pintu masuknya selalu lokasi. Warga mengetik alamat atau kelurahan, lalu melihat wilayahnya lewat tiga lensa yang saling bertaut:

| Lensa | Pertanyaan warga | Sumber utama |
| --- | --- | --- |
| Layanan | Apakah kursi sekolah di wilayahku cukup untuk lulusan tahun ini? | Dapodik, data induk peserta didik, juknis SPMB |
| Uang | Proyek apa yang dibiayai di sini, berapa nilainya, sampai mana buktinya? | LPSE DKI, SiRUP, Opentender |
| Tanggung Jawab | Instansi mana yang berwenang dan siapa wakilku di DPRD? | KPU, DPRD DKI, regulasi kewenangan |

Lingkup awal: Provinsi DKI Jakarta. Sektor utama pendidikan, sektor kedua jalan dan saluran air.

Konteks: karya cabang Web In Action (WIA), MIA 2026, tema "Designing Digital Futures, Powering Public Transparency". Rantai menyentuh keempat subtema: Open Government Data, Public Service Tracking, Civic Reporting and Accountability, dan Data-Driven Decision Support.

## 2. Masalah

1. **Jejak uang publik putus setelah tender.** Rencana pengadaan (SiRUP), tender (SPSE), dan anggaran (SMARTAPBD) tersedia terpisah. Data fase pelaksanaan seperti progres dan serah terima hampir tidak dipublikasikan. Warga juga jarang tahu bahwa pekerjaan konstruksi permanen punya masa pemeliharaan minimal 6 bulan, saat kontraktor wajib memperbaiki kerusakan.
2. **Kursi sekolah tidak sebanding dengan lulusan.** Contoh 2026: di Kota Malang, 13.792 lulusan SD/MI hanya tersedia 7.360 kursi SMP negeri. Pemda wajib mengumumkan daya tampung untuk SPMB, tetapi per Mei 2026 baru 237 dari 514 kabupaten/kota yang mengunggah dokumen penetapannya.
3. **Warga tidak tahu siapa yang harus ditanya.** Perhatian publik tertuju ke politik pusat. Kanal aspirasi dibangun per lembaga, sehingga warga harus tahu sendiri lembaga mana yang berwenang atas masalahnya.

Celah yang diisi: belum ada alat publik independen yang menyambungkan layanan, uang, dan tanggung jawab dalam satu wilayah dengan setiap angka bisa ditelusuri ke sumbernya.

## 3. Tujuan

| Tujuan | Ukuran keberhasilan |
| --- | --- |
| Warga memahami kondisi wilayahnya dalam satu kunjungan | Dari beranda ke ringkasan tiga lensa sebuah kelurahan dalam 3 klik atau kurang |
| Setiap angka bisa dipercaya | 100% angka yang tampil punya provenance (sumber, waktu ambil, label) |
| Uang publik terlacak sampai lokasi | Minimal 70% paket Pekerjaan Konstruksi LPSE DKI 2025 sampai 2026 berhasil dipetakan ke kecamatan |
| Warga bisa bertindak | Dari masalah ke paket aspirasi berbasis data dalam 2 menit |
| Instansi bisa menjawab | Setiap sinyal dan laporan punya jalur hak jawab yang tampil berdampingan |
| Demo selalu lancar | Alur demo utama berjalan penuh tanpa internet (mode snapshot) |

## 4. Pengguna

| Persona | Kebutuhan | Fitur kunci |
| --- | --- | --- |
| **Warga umum** (orang tua, pekerja, mahasiswa) | Memahami wilayahnya tanpa istilah teknis, tahu harus bertanya ke siapa | Wilayahku, Tanya Rantai, Paket Aspirasi |
| **Orang tua calon siswa** | Peluang realistis anak masuk SMP atau SMA negeri terdekat | Cek Sekolahku, peta kesenjangan |
| **Warga pelapor** | Melaporkan proyek mangkrak atau rusak dan melihat tindak lanjutnya | Lapor Bukti, status laporan, Pantau Wilayahku |
| **Jurnalis dan pegiat** | Menemukan pola, membandingkan, mengunduh data bersumber | Perbandingan paket sejenis, sinyal, unduh CSV |
| **Petugas instansi** | Menjawab sinyal dan laporan secara resmi | Panel Instansi, hak jawab |
| **Moderator Rantai** | Menyaring laporan spam dan menjaga kualitas | Panel moderasi |

## 5. Ruang lingkup

### 5.1 Termasuk

- Seluruh kecamatan dan kelurahan DKI Jakarta.
- Paket pengadaan LPSE Provinsi DKI Jakarta tahun anggaran 2025 dan 2026, prioritas Pekerjaan Konstruksi.
- Riwayat penyedia dari Opentender (2014 sampai 2023).
- Sekolah SD, SMP, SMA, SMK negeri dan swasta dari Dapodik; madrasah sejauh datanya tersedia.
- 106 anggota DPRD DKI periode 2024 sampai 2029, 10 dapil, dan komisinya.
- Antarmuka bilingual (Indonesia dan Inggris), tema terang dan gelap.

### 5.2 Tidak termasuk

- Provinsi lain (arsitektur disiapkan agar bisa diperluas, tetapi data hanya DKI).
- Rincian komponen APBD per baris (data publiknya belum terverifikasi; lihat SPEC bagian Gerbang Data).
- Sengketa perdata antarwarga (data perkara memuat identitas para pihak).
- Skor, peringkat, atau rating politisi dan partai.
- Penyalinan data LHKPN; Rantai hanya menautkan ke e-Announcement KPK.
- Pengiriman aspirasi atas nama warga; Rantai menyusun, warga yang mengirim.

## 6. Kebutuhan fungsional

Prioritas: **P0** wajib ada untuk demo utama, **P1** memperkuat penilaian, **P2** nilai tambah bila waktu cukup. Rincian tiap fitur ada di `docs/SPEC.md` dengan ID yang sama.

### 6.1 Fondasi

| ID | Fitur | Prioritas |
| --- | --- | --- |
| DASAR-01 | Pencarian wilayah dari alamat, kelurahan, atau kecamatan | P0 |
| DASAR-02 | Lencana provenance di setiap angka | P0 |
| DASAR-03 | Halaman Wilayahku dengan tiga tab lensa | P0 |
| DASAR-04 | Peta interaktif dengan lapisan kursi, proyek, dan laporan | P0 |
| DASAR-05 | Akun warga (OTP email) dan peran instansi, moderator, admin | P0 |
| DASAR-06 | Bilingual ID/EN, tema terang/gelap, aksesibilitas WCAG 2.1 AA | P0 |
| DASAR-07 | Mode demo luring dari snapshot data | P0 |
| DASAR-08 | Halaman Tentang Data: sumber, metode, lisensi, pembaruan terakhir | P0 |

### 6.2 Lensa Layanan

| ID | Fitur | Prioritas |
| --- | --- | --- |
| LAY-01 | Perhitungan calon lulusan, kursi proksi, dan kesenjangan per kecamatan | P0 |
| LAY-02 | Peta kesenjangan kecamatan (choropleth) | P0 |
| LAY-03 | Detail kecamatan: negeri, swasta, madrasah, daftar sekolah | P0 |
| LAY-04 | Cek Sekolahku: sekolah negeri terdekat, kursi, perkiraan persaingan | P1 |
| LAY-05 | Rapor keterbukaan daya tampung resmi vs proksi | P1 |
| LAY-06 | Tren kesenjangan antar tahun | P2 |

### 6.3 Lensa Uang

| ID | Fitur | Prioritas |
| --- | --- | --- |
| UANG-01 | Impor paket LPSE DKI dan taut ke RUP | P0 |
| UANG-02 | Penentuan lokasi paket dari nama paket (geocoding) | P0 |
| UANG-03 | Halaman Rantai Rupiah per paket | P0 |
| UANG-04 | Kategori paket sejenis dan perbandingan HPS | P0 |
| UANG-05 | Sinyal pola ("perlu ditanyakan") | P1 |
| UANG-06 | Estimasi masa pemeliharaan dan hitung mundur | P0 |
| UANG-07 | Profil penyedia dan riwayat paket | P1 |
| UANG-08 | Lapor Bukti: foto bergeotag tertaut ke paket | P0 |
| UANG-09 | Verifikasi laporan dan alur status | P0 |
| UANG-10 | Unduh data paket per wilayah (CSV) | P2 |

### 6.4 Lensa Tanggung Jawab

| ID | Fitur | Prioritas |
| --- | --- | --- |
| TJ-01 | Peta kewenangan: masalah ke instansi, komisi DPRD, dasar hukum | P0 |
| TJ-02 | Wakilku: kecamatan ke dapil ke anggota DPRD | P0 |
| TJ-03 | Profil anggota: dapil, fraksi, komisi, tautan LHKPN | P1 |
| TJ-04 | Jalur aspirasi terdekat: LMK, Dewan Kota, CRM/JAKI | P1 |
| TJ-05 | Paket Aspirasi: ringkasan berbasis data siap salin dan unduh PDF | P0 |

### 6.5 Lintas lensa

| ID | Fitur | Prioritas |
| --- | --- | --- |
| LINTAS-01 | Tanya Rantai: tanya jawab bahasa sehari-hari dengan jawaban bersumber | P0 |
| LINTAS-02 | Hak jawab instansi dan Panel Instansi | P0 |
| LINTAS-03 | Pantau Wilayahku: langganan notifikasi per wilayah | P1 |
| LINTAS-04 | Panel moderasi laporan | P1 |
| LINTAS-05 | Bagikan ringkasan wilayah sebagai gambar | P2 |

## 7. Alur demo utama

Alur ini wajib berjalan sempurna, termasuk dalam mode luring, dan diuji end-to-end.

1. Warga membuka beranda, mengetik nama kecamatan di Jakarta Utara.
2. Wilayahku terbuka di tab Layanan: lulusan SD tahun ini jauh melebihi kursi SMP negeri.
3. Tab Uang: daftar paket sungguhan di kecamatan itu tampil dengan nilai, penyedia, jadwal, dan status masa pemeliharaan. Untuk memperagakan alur laporan, ada satu paket simulasi pembangunan ruang kelas dengan penyedia fiktif dan dua laporan warga simulasi, seluruhnya berlabel "Data simulasi". Laporan simulasi tidak pernah ditempelkan ke paket atau penyedia sungguhan.
4. Tab Tanggung Jawab: urusan SMP dipegang Dinas Pendidikan provinsi, diawasi komisi DPRD bidang terkait, dan anggota dari dapil itu yang duduk di komisi tersebut.
5. Warga bertanya di Tanya Rantai: "Kenapa SMP negeri di kecamatanku penuh?" Jawaban menyertakan angka yang bisa diklik ke sumbernya.
6. Warga membuat Paket Aspirasi dan mengunduhnya.
7. Pindah ke akun instansi: petugas memberi hak jawab, dan jawaban langsung tampil di halaman publik.

## 8. Kebutuhan non-fungsional

| Aspek | Kebutuhan |
| --- | --- |
| Performa | LCP di bawah 2,5 detik pada 4G; respons API p95 di bawah 300 ms untuk halaman wilayah; Lighthouse 90 ke atas untuk Performance, Accessibility, Best Practices, SEO |
| Responsif | Tata letak benar di lebar 360, 768, 1024, dan 1440 piksel |
| Aksesibilitas | WCAG 2.1 AA: kontras, navigasi papan ketik, label pembaca layar, peta punya alternatif tabel |
| Keandalan | Situs tetap berjalan bila sumber pemerintah sedang tidak bisa diakses; data dari snapshot terakhir dengan tanggal ambil tampil |
| Keamanan | JWT dual-token, rate limit per IP dan per akun, CSP ketat, validasi dan pemindaian unggahan, query terparameterisasi, rahasia di env |
| Privasi | EXIF foto dihapus untuk publik; koordinat asli hanya untuk verifikasi; email warga tidak pernah tampil publik |
| Netralitas | Tanpa skor politisi; sinyal selalu dengan alasan dan ruang hak jawab; kosakata tuduhan dilarang |
| Keterbukaan | Turunan data Opentender dibuka dengan lisensi ODbL; metode perhitungan dijelaskan di Tentang Data |
| Kualitas kode | Clean Code tanpa komentar, lint bersih, tes untuk semua aturan bisnis (lihat CLAUDE.md) |

## 9. Risiko dan mitigasi

| Risiko | Dampak | Mitigasi |
| --- | --- | --- |
| Struktur halaman SPSE berubah atau membatasi akses | Data paket berhenti diperbarui | Importer berbasis fixture, cache, snapshot; verifikasi di gerbang data; laju 1 permintaan per detik |
| Data rombel per tingkat tidak bisa diambil dari Dapodik | Lensa Layanan tidak akurat | Gerbang data di awal; cadangan: data induk peserta didik per wilayah dan impor CSV manual |
| Geocoding nama paket salah | Proyek muncul di kecamatan yang salah | Hanya cocok bila nama wilayah ditemukan pasti; selain itu ditaruh di tingkat kota dengan label "lokasi tingkat kota" |
| Sinyal terbaca sebagai tuduhan | Risiko hukum dan kepercayaan | Kosakata netral, alasan dan pembanding wajib, hak jawab, ambang konservatif |
| Laporan warga palsu atau spam | Kepercayaan turun | Akun OTP, batas laju, verifikasi saksi independen, moderasi |
| Tanya Rantai mengarang angka | Kepercayaan turun | Jawaban hanya dari fungsi kueri tetap; angka tanpa rujukan ditolak validator |
| Internet putus saat final | Demo gagal | Mode demo luring dari snapshot, termasuk jawaban Tanya Rantai yang di-cache untuk pertanyaan demo |

## 10. Keluaran lomba yang didukung produk

| Tahap | Kebutuhan dari panduan | Yang disediakan repo |
| --- | --- | --- |
| First Draft (1 November 2026) | Repo publik dengan README berisi judul dan deskripsi, tech stack, cara menjalankan, progres; PDF maksimal 10 halaman dengan konsep, sitemap, wireframe, fitur, screenshot landing page dan 1 halaman lain | README terpelihara di setiap fase; halaman beranda dan Wilayahku diprioritaskan |
| Final Submission (11 November 2026) | Proposal maksimal 30 halaman, surat orisinalitas, video pitching maksimal 5 menit | Alur demo stabil untuk direkam, halaman Tentang Data sebagai bahan proposal |
| Final (28 November 2026) | Pitch deck maksimal 15 halaman, presentasi 10 menit dan tanya jawab 5 menit, prototipe daring dengan cadangan luring | Mode demo luring, skenario demo terkunci di seed |

Bobot penilaian penyisihan sebagai pengingat prioritas: Fungsionalitas 25%, Orisinalitas 20%, Desain 20%, Konten dan Informasi 20%, Kesesuaian Tema 15%.
