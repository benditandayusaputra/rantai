# PRD: Rantai

Versi 2.0 · 23 Oktober 2026 · Menggantikan PRD v1.0

## 1. Ringkasan

Rantai menunjukkan kepada warga siapa wakilnya di DPRD, ke mana uang kotanya mengalir, dan apakah layanan dasar di wilayahnya cukup. Setiap angka bisa ditelusuri ke sumbernya.

- **Pintu masuk**: lokasi. Warga mengetik alamat, kelurahan, atau kecamatan.
- **Lingkup**: Provinsi DKI Jakarta, tingkat kecamatan.
- **Bentuk**: platform dengan inti tetap dan tema yang bisa ditambah. Tema pertama adalah Pendidikan.
- **Konteks**: karya cabang Web In Action, MIA 2026, tema "Designing Digital Futures, Powering Public Transparency". Rantai menyentuh keempat subtema: Open Government Data, Public Service Tracking, Civic Reporting and Accountability, dan Data-Driven Decision Support.
- **Tim**: Galih, Kevin, Bendi.

Perubahan utama dari v1.0: struktur tiga lensa menjadi inti plus tema, Wakilku diperluas dengan umpan kegiatan dua sumber, kursi sekolah memakai data resmi, dan pokir dikeluarkan karena tidak publik.

## 2. Masalah

Warga tidak punya satu tempat untuk menyambungkan layanan, uang, dan tanggung jawab di wilayahnya.

1. **Warga tidak kenal wakilnya.** Perhatian publik tertuju ke politik pusat. Warga jarang tahu siapa anggota DPRD dari dapilnya, di komisi mana ia duduk, dan apa yang pernah ia suarakan.
2. **Jejak uang publik terputus.** Rencana pengadaan, tender, dan anggaran tersedia terpisah. Data pelaksanaan seperti progres dan serah terima hampir tidak dipublikasikan.
3. **Layanan dasar tidak sebanding dengan kebutuhan.** Contoh: daya tampung SMP negeri DKI 2026/2027 adalah 73.289 kursi, tetapi warga tidak bisa melihat kesenjangannya per kecamatan.
4. **Informasi tentang wakil rakyat berat sebelah.** Berita resmi DPRD ditulis oleh Sekretariat DPRD sendiri dan tidak memuat sorotan kritis dari media.
5. **Sebagian data penting belum dibuka.** Kehadiran per anggota, pokir, dan laporan reses per anggota tidak tersedia untuk publik.

Celah yang diisi: belum ada alat publik independen yang menyambungkan ketiganya dalam satu wilayah, dengan setiap angka bisa ditelusuri ke sumbernya.

## 3. Tujuan

Rantai berhasil bila warga bisa memahami wilayahnya dan bertindak dalam satu kunjungan, dengan data yang bisa dipercaya.

| Tujuan | Ukuran keberhasilan |
| --- | --- |
| Warga memahami wilayahnya dalam satu kunjungan | Dari beranda ke halaman kecamatan dalam 3 klik atau kurang |
| Warga kenal wakilnya | Dari alamat ke daftar anggota DPRD dapilnya dalam 1 langkah |
| Setiap angka bisa dipercaya | 100% angka yang tampil punya lencana sumber, waktu ambil, dan label |
| Uang publik terlacak sampai lokasi | Minimal 70% paket Pekerjaan Konstruksi LPSE DKI 2025 sampai 2026 terpetakan ke kecamatan |
| Informasi wakil rakyat berimbang | Setiap anggota punya umpan dari rilis resmi dan media independen |
| Warga bisa bertindak | Dari masalah ke Paket Aspirasi dalam 2 menit |
| Platform bisa diperluas | Tema baru cukup satu konfigurasi dan satu importer data |
| Demo selalu lancar | Alur demo berjalan penuh tanpa internet |

## 4. Pengguna

Pengguna utama adalah warga umum; persona lain memperdalam fitur yang sama.

| Persona | Kebutuhan | Fitur kunci |
| --- | --- | --- |
| **Warga umum** | Memahami wilayahnya tanpa istilah teknis, tahu harus bertanya ke siapa | Wilayah, Wakilku, Tanya Rantai, Paket Aspirasi |
| **Orang tua calon siswa** | Peluang realistis anak masuk SMP atau SMA negeri terdekat | Tema Pendidikan, Cek Sekolahku |
| **Warga pelapor** | Melaporkan proyek mangkrak atau rusak dan melihat tindak lanjutnya | Lapor Bukti, Pantau Wilayahku |
| **Jurnalis dan pegiat** | Menemukan pola, membandingkan, mengunduh data bersumber | Arus Anggaran, umpan Wakilku, unduh CSV |
| **Anggota DPRD dan petugas instansi** | Menjawab sorotan secara resmi | Hak jawab, Panel Instansi |
| **Moderator Rantai** | Menyaring laporan spam dan menjaga kualitas | Panel moderasi |

## 5. Arsitektur produk

Rantai terdiri dari inti tetap dan tema yang bisa dipasang; tema baru cukup satu konfigurasi dan satu importer data. Setiap tema memakai inti dan sumber data yang sama.

| Lapis | Isi |
| --- | --- |
| **Tema** | Pendidikan (kursi vs lulusan SD) aktif. Menyusul: Banjir (genangan vs drainase), Kesehatan (penduduk vs puskesmas), Jalan (laporan vs perawatan) |
| **Inti** | Wilayah (cari alamat, peta, halaman kecamatan) · Wakilku (profil, umpan dua sumber, sikap fraksi) · Arus Anggaran (paket LPSE per kecamatan, penyedia) · Aksi (Tanya Rantai, Paket Aspirasi) |
| **Sumber data** | Peta dan data indikator tema · KPU, situs DPRD, media · LPSE, SiRUP, Opentender |

Tema hanya memberi tahu inti apa yang harus dilihat. Peta, anggaran, umpan wakil rakyat, Tanya Rantai, dan Paket Aspirasi bekerja otomatis dari konfigurasi itu.

**Kontrak tema.** Setiap tema menjawab lima pertanyaan:

| Bagian | Pertanyaan | Dipakai oleh |
| --- | --- | --- |
| Indikator | Seberapa besar kesenjangan layanan di sini? | Wilayah, peta |
| Filter anggaran | Belanja mana yang terkait isu ini? | Arus Anggaran |
| Tanggung jawab | Instansi, komisi, dan dasar hukum mana yang berwenang? | Wakilku, Aksi |
| Tag isu | Berita wakil rakyat mana yang relevan? | Wakilku |
| Templat aspirasi | Apa yang bisa diminta warga? | Aksi |

Konfigurasi tema disimpan di basis data, bukan halaman yang di-hardcode.

## 6. Fitur inti

Empat modul inti berlaku untuk setiap tema dan setiap wilayah. Prioritas: **P0** wajib untuk demo utama, **P1** memperkuat penilaian, **P2** nilai tambah bila waktu cukup.

### 6.1 Wilayah

| ID | Fitur | Prioritas |
| --- | --- | --- |
| WIL-01 | Pencarian wilayah dari alamat, kelurahan, atau kecamatan | P0 |
| WIL-02 | Halaman kecamatan yang menyatukan Wakilku, Arus Anggaran, dan tema aktif | P0 |
| WIL-03 | Peta interaktif dengan lapisan per tema, proyek, dan laporan | P0 |
| WIL-04 | Lencana sumber di setiap angka | P0 |
| WIL-05 | Halaman Tentang Data: sumber, metode, lisensi, pembaruan terakhir | P0 |
| WIL-06 | Pantau Wilayahku: notifikasi per wilayah | P1 |

### 6.2 Wakilku

| ID | Fitur | Prioritas |
| --- | --- | --- |
| WAK-01 | Dari kecamatan ke dapil ke anggota DPRD DKI | P0 |
| WAK-02 | Profil anggota: dapil, fraksi, komisi, alat kelengkapan, tautan LHKPN | P0 |
| WAK-03 | Jejak Kegiatan: umpan dua sumber, rilis resmi DPRD dan media independen, masing-masing berlencana | P0 |
| WAK-04 | Isu yang Disuarakan: topik yang diangkat anggota, dihitung dari umpan dengan label jumlah sumber | P1 |
| WAK-05 | Sikap Fraksi: pandangan fraksi anggota atas raperda dan APBD | P1 |
| WAK-06 | Aspirasi Reses: aspirasi agregat per bidang dari dokumen paripurna | P1 |
| WAK-07 | Hak jawab anggota di setiap item umpan | P1 |

### 6.3 Arus Anggaran

| ID | Fitur | Prioritas |
| --- | --- | --- |
| ANG-01 | Impor paket LPSE DKI dan taut ke RUP SiRUP | P0 |
| ANG-02 | Penentuan lokasi paket dari nama paket | P0 |
| ANG-03 | Halaman paket: nilai, penyedia, jadwal, status | P0 |
| ANG-04 | Estimasi masa pemeliharaan dan hitung mundur | P0 |
| ANG-05 | Aspirasi ke Anggaran: cek apakah aspirasi reses muncul sebagai paket di wilayah terkait | P1 |
| ANG-06 | Perbandingan HPS paket sejenis dan sinyal "perlu ditanyakan" | P1 |
| ANG-07 | Profil penyedia dan riwayat paket | P1 |
| ANG-08 | Lapor Bukti: foto bergeotag tertaut ke paket, dengan verifikasi dan status | P1 |
| ANG-09 | Unduh data paket per wilayah dalam CSV | P2 |

### 6.4 Aksi

| ID | Fitur | Prioritas |
| --- | --- | --- |
| AKS-01 | Tanya Rantai: tanya jawab bahasa sehari-hari, jawaban hanya dari fungsi kueri tetap dan selalu bersumber | P0 |
| AKS-02 | Paket Aspirasi: ringkasan berbasis data ke kanal yang tepat, siap salin dan unduh PDF | P0 |
| AKS-03 | Jalur aspirasi terdekat: LMK, Dewan Kota, CRM atau JAKI | P1 |
| AKS-04 | Data yang Belum Dibuka: daftar data yang belum publik dengan templat permohonan informasi ke PPID | P1 |
| AKS-05 | Bagikan ringkasan wilayah sebagai gambar | P2 |

### 6.5 Fondasi

| ID | Fitur | Prioritas |
| --- | --- | --- |
| DSR-01 | Mode demo luring dari snapshot data | P0 |
| DSR-02 | Akun warga, anggota, instansi, moderator, admin; untuk demo cukup akun contoh | P0 |
| DSR-03 | Panel Instansi dan panel moderasi | P1 |
| DSR-04 | Bilingual Indonesia dan Inggris, tema terang dan gelap | P2 |

## 7. Tema Pendidikan

Tema pertama mengukur apakah kursi SMP negeri di sebuah kecamatan cukup untuk lulusan SD-nya, lalu menyambungkannya ke belanja sekolah dan ke wakil yang berwenang.

| Bagian | Isi |
| --- | --- |
| Indikator | Jumlah siswa kelas 6 dibanding daya tampung SMP per kecamatan, dipisah negeri, swasta SPMB Bersama, dan madrasah sejauh datanya tersedia |
| Filter anggaran | Paket dari Dinas Pendidikan; kata kunci seperti "ruang kelas", "rehab sekolah", "pembangunan SMP" |
| Tanggung jawab | Dinas Pendidikan Provinsi DKI, Komisi E DPRD DKI bidang Kesejahteraan Rakyat, dasar hukum kewenangan |
| Tag isu | sekolah, SPMB, zonasi, daya tampung, guru |
| Templat aspirasi | Surat tentang kekurangan kursi SMP di kecamatan, berisi angka dan sumbernya |

| ID | Fitur | Prioritas |
| --- | --- | --- |
| PEN-01 | Hitung kesenjangan kursi per kecamatan | P0 |
| PEN-02 | Peta kesenjangan kecamatan | P0 |
| PEN-03 | Detail kecamatan: daftar sekolah dan daya tampungnya | P0 |
| PEN-04 | Kesenjangan vs Belanja: kesenjangan kursi dibanding belanja pembangunan sekolah per kecamatan | P1 |
| PEN-05 | Anggota Komisi E yang pernah menyuarakan isu sekolah, tertaut ke umpan Wakilku | P1 |
| PEN-06 | Cek Sekolahku: sekolah negeri terdekat dan perkiraan persaingan | P1 |
| PEN-07 | Tren kesenjangan antar tahun | P2 |

PEN-04 adalah wawasan paling orisinal Rantai: apakah uang pembangunan sekolah mengalir ke kecamatan yang paling kekurangan kursi.

### 7.1 Tema berikutnya

| Tema | Indikator | Instansi dan komisi |
| --- | --- | --- |
| Banjir | Titik genangan dibanding proyek drainase | Dinas Sumber Daya Air, Komisi D |
| Kesehatan | Penduduk dibanding kapasitas puskesmas | Dinas Kesehatan, Komisi E |
| Jalan | Laporan kerusakan warga dibanding paket pemeliharaan jalan | Dinas Bina Marga, Komisi D |

Pemetaan komisi untuk tema berikutnya perlu dikonfirmasi saat tema tersebut dibangun.

## 8. Sumber data

Daya tampung sekolah, pengadaan, dan anggota DPRD sudah terverifikasi; jumlah siswa kelas 6 per kecamatan masih perlu dicek manual.

| Data | Sumber | Status | Catatan |
| --- | --- | --- | --- |
| Daya tampung SMP negeri per sekolah | Keputusan Kadis Pendidikan DKI No. 258 Tahun 2026 | Terverifikasi | Per sekolah dengan NPSN; total SMP negeri 73.289 kursi (ANTARA) |
| Kursi SPMB Bersama | Pengumuman Disdik DKI | Terverifikasi | 1.597 kursi SMP di sekolah swasta mitra |
| Lokasi sekolah ke kecamatan | Dapodik, dicocokkan lewat NPSN | Perlu dicek | Telusur Dapodik tersedia sampai tingkat sekolah |
| Siswa kelas 6 per kecamatan | Sekolah Kita atau Dapodik | Perlu dicek | Cadangan: Satu Data Jakarta, permintaan data, atau estimasi total SD dibagi 6 dengan label |
| Paket pengadaan | LPSE DKI lewat pyproc, SiRUP, ISB Satu Data LKPP | Layak | Portal sering tidak bisa diakses dari IP luar Indonesia; jalankan dari server Indonesia atau snapshot |
| Riwayat penyedia | Opentender 2014 sampai 2023 | Layak | Turunan dibuka dengan lisensi ODbL |
| Anggota DPRD dan dapil | KPU, situs DPRD DKI | Terverifikasi | 106 anggota, 10 dapil, 9 sampai 12 kursi per dapil |
| Komisi dan profil anggota | Situs DPRD DKI | Layak | Halaman per anggota dan per komisi |
| Rilis kegiatan DPRD | Berita situs DPRD DKI | Layak, perlu uji crawl | Ditulis Sekretariat DPRD |
| Liputan media | RSS berita per nama anggota | Layak | Hanya media terverifikasi Dewan Pers |
| Aspirasi reses | Dokumen dan berita paripurna | Agregat saja | Per bidang, bukan per anggota |
| Batas wilayah | Data kecamatan dan kelurahan DKI | Layak | Untuk peta dan pencarian |

## 9. Data yang belum dibuka

Empat data tidak publik, tetapi ini bukan penghalang: janji inti Rantai tidak bergantung padanya, dan setiap celah punya pengganti.

| Data | Kenapa tidak tersedia | Pengganti di Rantai |
| --- | --- | --- |
| Suara per anggota | Keputusan DPRD diambil per fraksi atau musyawarah | Sikap Fraksi, WAK-05 |
| Pokir per anggota | Diinput ke SIPD, sistem internal pemerintah; DKI dilaporkan tidak lagi mengakomodasi pokir sejak 2015 | Aspirasi ke Anggaran, ANG-05 |
| Kehadiran per anggota | Hanya jumlah hadir yang muncul di berita | Data yang Belum Dibuka, AKS-04 |
| Laporan reses per anggota | Hanya kompilasi agregat per bidang | Aspirasi Reses, WAK-06 |

Halaman Data yang Belum Dibuka menampilkan daftar ini beserta templat permohonan informasi ke PPID DPRD DKI berdasarkan UU Keterbukaan Informasi Publik. Celah data menjadi alat advokasi.

Di luar produk, tim mengajukan permohonan informasi kehadiran dan laporan reses ke PPID sekarang. Jawabannya kemungkinan tidak datang sebelum 11 November, tetapi permohonannya bisa disebut saat pitching.

## 10. Prinsip kepercayaan dan netralitas

Rantai hanya berguna bila dipercaya oleh warga dan oleh pihak yang disorot.

- **Setiap angka bersumber.** Lencana berisi sumber, waktu ambil, dan label. Angka tanpa rujukan tidak tampil.
- **Dua sumber untuk wakil rakyat.** Rilis resmi dan media independen tampil berdampingan dengan lencana jenis sumber.
- **Daftar media tertutup.** Hanya perusahaan pers terverifikasi Dewan Pers.
- **Tidak menilai.** Tidak ada skor, peringkat, atau rating politisi dan partai. Rantai menampilkan judul, ringkasan netral satu kalimat, media, tanggal, dan tautan.
- **Kosakata netral.** Pakai "menyuarakan" dan "perlu ditanyakan", bukan "memperjuangkan" atau "menyimpang".
- **Hak jawab di semua sorotan.** Anggota dan instansi bisa menjawab, dan jawabannya tampil berdampingan.
- **Bias liputan diakui.** Halaman Wakilku menyatakan bahwa anggota yang lebih sering diliput akan terlihat lebih aktif.
- **Ringkasan, bukan salinan.** Tidak ada teks berita penuh.
- **Data simulasi** selalu berlabel dan tidak pernah ditempelkan ke paket atau penyedia sungguhan.
- **Tanya Rantai tidak mengarang.** Jawaban hanya dari fungsi kueri tetap; angka tanpa rujukan ditolak validator.
- **Disambiguasi nama.** Artikel media dicocokkan pada nama plus "DPRD DKI" atau fraksi.

## 11. Alur demo utama

Alur ini wajib berjalan sempurna, termasuk dalam mode luring, dan diuji end-to-end. Ceritanya: kursi SMP di kecamatanku kurang, uangnya ke mana, siapa yang bertanggung jawab, dan apa yang bisa kulakukan.

1. Warga membuka beranda dan mengetik nama kecamatan di Jakarta Utara.
2. **Tema Pendidikan**: lulusan SD tahun ini jauh melebihi kursi SMP negeri, dengan angka resmi per sekolah.
3. **Arus Anggaran**: paket pembangunan sekolah di kecamatan itu tampil dengan nilai, penyedia, jadwal, dan masa pemeliharaan. Grafik Kesenjangan vs Belanja menunjukkan posisi kecamatan ini dibanding kecamatan lain.
4. **Wakilku**: urusan SMP dipegang Dinas Pendidikan dan diawasi Komisi E. Anggota dapil itu yang duduk di Komisi E tampil dengan umpan dua sumber dan isu yang pernah disuarakan.
5. **Tanya Rantai**: "Kenapa SMP negeri di kecamatanku penuh?" Jawaban menyertakan angka yang bisa diklik ke sumbernya.
6. **Paket Aspirasi**: warga membuat surat berbasis data ke kanal yang tepat dan mengunduhnya.
7. **Hak jawab**: pindah ke akun instansi contoh; petugas memberi jawaban dan jawabannya langsung tampil di halaman publik.
8. **Penutup**: halaman Data yang Belum Dibuka sebagai jembatan ke Rencana Pengembangan.

Untuk memperagakan Lapor Bukti, ada satu paket simulasi dengan penyedia fiktif dan dua laporan simulasi, seluruhnya berlabel "Data simulasi".

## 12. Kebutuhan non-fungsional

| Aspek | Kebutuhan |
| --- | --- |
| Performa | LCP di bawah 2,5 detik pada 4G; API p95 di bawah 300 ms untuk halaman wilayah; Lighthouse 90 ke atas |
| Responsif | Tata letak benar di lebar 360, 768, 1024, dan 1440 piksel |
| Aksesibilitas | WCAG 2.1 AA; peta punya alternatif tabel |
| Keandalan | Situs tetap jalan saat sumber pemerintah tidak bisa diakses; data dari snapshot terakhir dengan tanggal ambil |
| Perluasan | Tema disimpan sebagai konfigurasi di basis data, bukan halaman yang di-hardcode |
| Keamanan | JWT dual-token, rate limit per IP dan akun, CSP ketat, validasi dan pemindaian unggahan, query terparameterisasi |
| Privasi | EXIF foto dihapus untuk publik; email warga tidak pernah tampil |
| Hukum | Ringkasan dan tautan saja untuk berita; tidak ada karakterisasi oleh Rantai; hak jawab di semua sorotan |
| Keterbukaan | Metode perhitungan dijelaskan di Tentang Data; turunan Opentender berlisensi ODbL |
| Kualitas kode | Lint bersih dan tes untuk semua aturan bisnis |

## 13. Ruang lingkup

### 13.1 Termasuk

- Seluruh kecamatan dan kelurahan DKI Jakarta.
- 106 anggota DPRD DKI periode 2024 sampai 2029, 10 dapil, lima komisi, dan fraksinya.
- Paket pengadaan LPSE DKI tahun anggaran 2025 dan 2026, prioritas Pekerjaan Konstruksi.
- Riwayat penyedia dari Opentender 2014 sampai 2023.
- Tema Pendidikan: SD, SMP negeri dan swasta dari Dapodik; madrasah sejauh datanya tersedia.

### 13.2 Tidak termasuk

- Provinsi lain. Arsitektur siap diperluas, tetapi data hanya DKI.
- Tema selain Pendidikan pada versi lomba. Banjir, Kesehatan, dan Jalan masuk Rencana Pengembangan.
- Rincian APBD per baris, karena keterbukaannya belum terverifikasi.
- Pokir, suara, dan kehadiran per anggota, karena tidak publik.
- Sengketa perdata antarwarga, karena memuat identitas para pihak.
- Skor atau peringkat politisi dan partai.
- Penyalinan data LHKPN; Rantai hanya menautkan ke e-Announcement KPK.
- Pengiriman aspirasi atas nama warga; Rantai menyusun, warga yang mengirim.

## 14. Risiko dan mitigasi

Risiko terbesar adalah data siswa kelas 6 per kecamatan; seluruh tema Pendidikan bergantung padanya.

| Risiko | Dampak | Mitigasi |
| --- | --- | --- |
| Data kelas 6 per kecamatan tidak tersedia publik | Indikator Pendidikan tidak akurat | Cek Sekolah Kita minggu ini; cadangan Satu Data Jakarta, permintaan data, atau estimasi berlabel |
| Portal SPSE berubah atau memblokir akses luar negeri | Data paket berhenti diperbarui | Jalankan importer dari server Indonesia; cache dan snapshot; laju 1 permintaan per detik |
| Geocoding nama paket salah | Proyek muncul di kecamatan yang salah | Cocok hanya bila nama wilayah ditemukan pasti; selain itu label "lokasi tingkat kota"; uji 100 nama paket dulu |
| Berita resmi DPRD bias ke pencitraan | Wakilku menjadi corong humas | Umpan dua sumber dengan media terverifikasi Dewan Pers |
| Nama anggota tertukar dengan orang lain | Kegiatan salah atribusi | Cocokkan nama plus "DPRD DKI" atau fraksi dalam artikel yang sama |
| Sorotan terbaca sebagai tuduhan | Risiko hukum dan kepercayaan | Tidak ada karakterisasi, kosakata netral, hak jawab |
| Situs DPRD tidak bisa di-crawl | Umpan resmi kosong | Uji minggu ini; cadangan impor manual untuk anggota demo |
| Laporan warga palsu | Kepercayaan turun | Akun OTP, batas laju, verifikasi, moderasi |
| Tanya Rantai mengarang angka | Kepercayaan turun | Fungsi kueri tetap dan validator rujukan |
| Internet putus saat final | Demo gagal | Mode luring dari snapshot, termasuk jawaban Tanya Rantai yang di-cache |
