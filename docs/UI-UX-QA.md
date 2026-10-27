# QA UI/UX Rantai: mudah dipahami pengguna baru

Tanggal: 27 Oktober 2026. Fokus: warga umum yang baru pertama membuka Rantai, terutama melalui ponsel.

## Kesimpulan

Identitas visual sudah konsisten: tipografi, warna, pembagian bagian, dan tombol cukup rapi. Kekuatan produk adalah hubungan antara kondisi layanan, pengadaan, dan pihak yang berwenang, disertai sumber angka.

Hambatan terbesarnya adalah urutan pemahaman. Pengguna melihat banyak nama fitur sebelum mendapat jawaban sederhana: apa manfaat Rantai untuk saya, harus mulai dari mana, apa arti angka ini, dan apa langkah berikutnya? Prioritas perbaikan adalah kejelasan pintu masuk, bahasa, konteks wilayah, dan kesinambungan tindakan.

## Dasar dan batas pemeriksaan

- Meninjau source frontend, pencarian wilayah pada backend, PRD, dan design system.
- Menginspeksi enam screenshot dalam `docs/screenshots/`: beranda desktop/ponsel, kecamatan desktop/ponsel, paket desktop, dan peta desktop. Screenshot merupakan artefak repo; beberapa elemen source saat ini sudah lebih berkembang, termasuk tur, peta berlapis, dan pilihan jenjang. Screenshot bukan bukti tampilan produksi terkini.
- Menjalankan tes komponen pencarian dan fitur wilayah, navigasi, serta tur: **4 file, 26 tes lulus**.
- Browser interaktif tidak tersedia pada sesi ini; pembacaan pratinjau melalui alat web juga gagal. Audit ini tidak mengklaim pengujian klik end-to-end, rendering terkini, GPS perangkat, atau pengukuran kontras langsung.
- Prediksi kebingungan pengguna di bawah merupakan penilaian heuristik, belum hasil observasi peserta uji.
- Tidak mengubah implementasi aplikasi.

Prioritas: **P1** menghambat alur inti atau berisiko menimbulkan salah pengertian; **P2** memperlambat pemahaman atau penemuan fitur. Nomor temuan bukan ukuran frekuensi kejadian.

## Temuan

### 01 — Manfaat dan langkah pertama belum cukup tegas · P1

**Bukti:** beranda memulai dengan “Selamat datang di Rantai” dan “Mengenal kota, dimulai dari sekitar kita”, disusul hero “Kenali wilayahmu. Ikuti jejak uangnya.” Tombol hero menuju peta, sedangkan pencarian langsung menuju kecamatan. Pada screenshot ponsel, area pembuka dan hero menghabiskan hampir seluruh layar pertama.

**Dampak yang diperkirakan:** warga harus menyimpulkan sendiri perbedaan mencari wilayah dan menjelajahi peta, serta hasil apa yang akan diperoleh.

**Perbaikan:** jadikan pemilihan wilayah sebagai tindakan utama. Letakkan janji manfaat bersama pencarian, misalnya:

> Cek sekolah, proyek publik, dan wakil rakyat di wilayahmu.
>
> Pilih kecamatan atau kelurahan di DKI Jakarta untuk mulai.

Peta tetap tersedia sebagai pilihan pendukung. Pertahankan foto, tetapi kurangi tinggi hero pada ponsel bila menghalangi akses ke manfaat dan hasil.

**Kriteria penerimaan:** tanpa membaca tur, peserta bisa menjelaskan manfaat Rantai dan menunjukkan langkah pertama dalam 10 detik.

**Lokasi:** `frontend/src/routes/[lang=lang]/+page.svelte`, `frontend/src/lib/i18n/id.json` (`home.*`).

### 02 — Tur otomatis terlalu dini dan berorientasi menu · P2

**Bukti:** `TurPengenalan.svelte` memanggil `mulaiTur` setelah 800 ms pada kunjungan pertama. Definisi tur berisi sambutan dan tujuh tujuan: cari, peta, dasbor, tanya, lapor, pantau, tentang. Target yang tidak tersedia disaring sehingga jumlah aktual berbeda menurut layar/halaman. Tombol lewati dan pemutaran ulang sudah tersedia.

**Dampak yang diperkirakan:** modal dapat memotong orang yang baru membaca atau mulai mengetik. Mengingat banyak fitur sekaligus menambah beban sebelum pengguna melihat hasil.

**Perbaikan:** sediakan bantuan opsional “Cara memakai Rantai”. Panduan awal cukup mengantar satu tugas: pilih wilayah, baca ringkasan, pilih langkah berikutnya. Penjelasan fitur lain muncul saat fitur dibuka.

**Kriteria penerimaan:** pengguna baru dapat mencari wilayah tanpa menutup modal lebih dahulu; bantuan tetap mudah ditemukan dan dapat dilewati.

**Lokasi:** `frontend/src/lib/features/tur/TurPengenalan.svelte`, `frontend/src/lib/features/tur/tur.ts`.

### 03 — Janji pencarian tidak sesuai cakupan dan petunjuk tersembunyi · P1

**Bukti:** placeholder berbunyi “Cari kelurahan atau alamat”, tetapi kueri backend mencocokkan nama/alias wilayah. Ini bukan pencarian alamat jalan lengkap. Minimal tiga karakter ditulis dalam bantuan `sr-only`, sehingga pengguna visual tidak melihatnya. Penyaringan toleransi ejaan sudah ada pada backend.

**Dampak yang diperkirakan:** pengguna mengetik alamat lengkap dan gagal menemukan hasil, atau menganggap pencarian tidak bekerja ketika baru mengetik dua huruf.

**Perbaikan:** gunakan “Cari kelurahan atau kecamatan di DKI Jakarta” dengan contoh yang terlihat, seperti “Contoh: Kelapa Gading · minimal 3 huruf”. Pada hasil kosong, jelaskan cakupan wilayah dan berikan jalur coba nama kecamatan.

**Kriteria penerimaan:** teks tidak menjanjikan geocoding alamat lengkap; pengguna mengetahui cakupan dan syarat input sebelum gagal.

**Lokasi:** `frontend/src/lib/components/RegionSearch.svelte`, `frontend/src/lib/i18n/id.json` (`search.*`), `backend/internal/repository/wilayah_baca.go`.

### 04 — Menekan Enter belum memberi hasil bila opsi belum disorot · P1

**Bukti kode:** indeks aktif dimulai dari `-1`. Handler Enter hanya memanggil `pilih` jika indeks aktif sudah `>= 0`. Mengetik lalu langsung Enter tidak memilih saran. Tes yang tersedia memeriksa rangkaian ArrowDown → Enter.

**Dampak yang diperkirakan:** pengguna keyboard atau tombol pencarian pada keyboard ponsel dapat menganggap pencarian macet. Ini perlu direproduksi di browser/perangkat pada tahap berikutnya.

**Perbaikan:** tentukan perilaku Enter yang eksplisit. Pilih kecocokan tepat yang tunggal, atau arahkan pengguna ke daftar hasil dengan pesan yang jelas jika ambigu. Jangan memilih kelurahan lain secara diam-diam hanya karena berada di urutan pertama.

**Kriteria penerimaan:** ketik nama wilayah → Enter selalu menghasilkan navigasi yang jelas atau instruksi memilih; kasus nama ambigu tidak salah memilih wilayah.

**Lokasi:** `frontend/src/lib/components/RegionSearch.svelte:93`, `frontend/src/lib/components/RegionSearch.test.ts`.

### 05 — Contoh wilayah dan konteks pribadi perlu dibedakan · P1

**Bukti:** beranda menggunakan Duren Sawit sebagai sorotan bawaan dan menyediakan empat wilayah sorotan. Pada ponsel, label visual “Wilayah sorotan” disembunyikan. Tautan global Tanya menuju `/tanya` tanpa `w`, sementara loader Tanya hanya mengambil konteks wilayah jika parameter `w` tersedia. Tautan kontekstual dari kecamatan sudah mendukung wilayah.

**Dampak yang diperkirakan:** pengguna dapat mengira contoh merupakan wilayahnya, atau mengira Tanya masih mengetahui kecamatan yang baru dibuka padahal konteksnya kosong.

**Perbaikan:** tulis “Contoh wilayah: Duren Sawit” sebelum pengguna memilih lokasi, termasuk pada ponsel. Tampilkan wilayah aktif secara konsisten dan sediakan “Ubah wilayah”. Bawa konteks ke tindakan berikutnya; bila belum ada, minta pengguna memilih wilayah secara jelas.

**Kriteria penerimaan:** contoh tidak tampil sebagai pilihan pribadi; setelah membuka kecamatan lalu menekan Tanya dari navigasi global, kecamatan aktif tetap terlihat atau ada permintaan memilih yang eksplisit. Pengguna tetap dapat membuang/mengubah konteks.

**Lokasi:** `frontend/src/lib/features/beranda/beranda.ts:9`, `frontend/src/routes/[lang=lang]/+page.svelte`, `frontend/src/lib/features/navigasi/navigasi.ts`, `frontend/src/routes/[lang=lang]/tanya/+page.ts`.

### 06 — Nama menu masih mencerminkan struktur produk · P2

**Bukti:** tab kecamatan pertama adalah “Tema”, sedangkan isi aktif adalah pendidikan. Navigasi bawah ponsel berisi enam tujuan, termasuk “Dasbor” dan “Tentang”, sementara Pantau tidak masuk navigasi bawah. Ini menambah pilihan yang tidak langsung menjelaskan kebutuhan warga.

**Dampak yang diperkirakan:** orang harus membuka menu untuk memahami isinya. “Tema” dapat diasosiasikan dengan pengaturan tampilan, dan “Dasbor” belum menjelaskan informasi apa yang tersedia.

**Perbaikan:** beri nama tab sesuai isi, seperti “Pendidikan” atau “Layanan publik” bila kelak meliputi beberapa topik. Uji navigasi bawah empat atau lima tujuan yang paling sering dipakai; pindahkan informasi pendukung ke “Lainnya”. Pertimbangkan “Bandingkan wilayah” hanya bila cocok dengan isi dasbor.

**Kriteria penerimaan:** peserta menemukan informasi sekolah, proyek, dan wakil tanpa mencoba menu secara acak. Pantau tetap dapat ditemukan dekat konteks wilayah.

**Lokasi:** `frontend/src/lib/features/navigasi/navigasi.ts`, `frontend/src/lib/features/kecamatan/kecamatan.ts`, `frontend/src/routes/[lang=lang]/kecamatan/[kode]/+page.svelte`.

### 07 — Angka perlu diterjemahkan menjadi makna bagi warga · P1

**Bukti:** bagian pendidikan menyajikan calon lulusan, kursi, selisih, dan rasio sebagai kartu setara. Screenshot kecamatan memperlihatkan contoh 1.387 calon lulusan, 684 kursi, selisih 703, dan rasio 0,49. Catatan tentang penerimaan lintas batas kecamatan sudah ada. Label sumber berulang dan teksnya terpotong pada kartu ponsel di screenshot.

**Dampak yang diperkirakan:** pengguna harus menghitung sendiri artinya atau menganggap selisih sebagai jumlah anak yang pasti gagal masuk. Detail sumber dapat bersaing dengan kesimpulan utama.

**Perbaikan:** tampilkan satu ringkasan sebelum rincian, misalnya “Sekitar 49 kursi SMP negeri tersedia untuk setiap 100 calon lulusan SD di kecamatan ini.” Letakkan batas interpretasi di dekatnya. Jelaskan perbedaan angka resmi dan hasil perhitungan. Pertahankan akses sumber untuk setiap angka; gunakan tombol sumber ringkas yang utuh pada ponsel.

**Kriteria penerimaan:** peserta dapat menjelaskan rasio dan tidak menyamakan selisih dengan hasil seleksi individu; sumber setiap nilai tetap dapat dibuka.

**Lokasi:** `frontend/src/lib/features/kecamatan/BagianTema.svelte:48`, `frontend/src/lib/components/ProvenanceBadge.svelte:41`, `docs/screenshots/kecamatan-ponsel.png`.

### 08 — “Putus” dapat dibaca sebagai kondisi proyek · P1

**Bukti:** Rantai Rupiah memberi status “Ada”, “Putus”, “Estimasi”, dan “Menunggu” pada tahapan berbeda. Screenshot paket menampilkan rencana/pelaksanaan “Putus” dengan alasan dokumen atau jadwal belum tersedia. HPS juga ditampilkan sebagai singkatan.

**Dampak yang diperkirakan:** pengguna dapat mengira “Putus” berarti pekerjaan dihentikan atau dana bermasalah, padahal status tersebut menerangkan ketersediaan data. “Ada” juga belum menyatakan apa yang tersedia.

**Perbaikan:** beri konteks “Kelengkapan data proyek” sebelum diagram. Gunakan status seperti “Data tersedia”, “Data belum tersedia”, dan “Perkiraan dari data yang ada”. Tulis “Perkiraan harga pemerintah (HPS)” pada kemunculan awal. Status pelaksanaan fisik, bila diketahui, harus punya label tersendiri.

**Kriteria penerimaan:** peserta membedakan dokumen yang belum tersedia dari pekerjaan yang terhenti; mereka memahami beda HPS dan nilai kontrak.

**Lokasi:** `frontend/src/lib/components/ChainStepper.svelte`, `frontend/src/lib/i18n/id.json` (`chain.*`, `pkg.hps`), `docs/screenshots/paket-desktop.png`.

### 09 — Lapor buntu ketika lokasi perangkat tidak tersedia · P1

**Bukti kode:** langkah lokasi hanya menyediakan permintaan GPS dan coba lagi. Alur baru maju setelah status `akurat`. Teks terjemahan “Pilih titik di peta” tersedia, tetapi tidak dirender pada komponen langkah lokasi.

**Dampak:** pada cabang GPS ditolak/gagal/tidak cukup akurat, pengguna tidak memiliki jalur manual untuk meneruskan laporan. Desktop dan pengguna yang enggan memberi izin lokasi terdampak.

**Perbaikan:** sediakan pemilihan titik di peta atau pencarian lokasi proyek secara manual. Jelaskan prasyarat akun, foto, dan mekanisme saksi sejak pengantar. Lokasi pilihan manual perlu ditandai sesuai tingkat verifikasinya agar tidak otomatis disamakan dengan bukti keberadaan GPS.

**Kriteria penerimaan:** simulasi izin GPS ditolak tetap memungkinkan menyusun laporan dengan lokasi manual; ketentuan verifikasi lokasi dan saksi tetap jelas.

**Lokasi:** `frontend/src/lib/features/lapor/LangkahLokasi.svelte:13`, `frontend/src/routes/[lang=lang]/lapor/+page@.svelte`, `frontend/src/lib/i18n/id.json` (`report.s1.pickOnMap`).

### 10 — Perpindahan dari data ke tindakan belum cukup dekat · P1

**Bukti:** kartu pendidikan menawarkan detail sekolah dan Cek Sekolahku; tombol membuat aspirasi terdapat di bagian Wakilku. Link Cek Sekolahku membawa jenjang tetapi tidak wilayah, dan halaman tujuan mulai dengan pusat lokasi kosong. Halaman aspirasi menyediakan pratinjau, salin, unduh PDF, dan catatan agar warga menyampaikan sendiri dokumen melalui kanal tercantum; tahap penyerahan belum tampil sebagai blok tindakan tersendiri di luar dokumen.

**Dampak yang diperkirakan:** setelah memahami kekurangan kursi, pengguna perlu mencari jalur tindakan lagi, memilih ulang lokasi di Cek Sekolahku, atau tidak tahu apa yang harus dilakukan setelah mengunduh dokumen.

**Perbaikan:** dekat ringkasan tampilkan pilihan “Lihat sekolah sekitar” dan “Susun aspirasi tentang kursi sekolah”. Bawa wilayah, jenjang, dan masalah secara eksplisit. Untuk pencarian berbasis titik tengah wilayah, jelaskan bahwa lokasi tersebut perkiraan sampai pengguna memilih titik rumah. Setelah dokumen tersedia, tampilkan penerima, kanal yang tersedia, langkah penyampaian, serta status “Draf siap; belum dikirim”.

**Kriteria penerimaan:** dari ringkasan masalah ke draf aspirasi tidak perlu memilih ulang wilayah/isu; peserta tahu dokumen belum dikirim dan dapat menunjukkan langkah penyampaiannya.

**Lokasi:** `frontend/src/lib/features/kecamatan/BagianTema.svelte:21`, `frontend/src/lib/features/kecamatan/BagianWakilku.svelte:33`, `frontend/src/routes/[lang=lang]/cek-sekolah/+page.svelte:34`, `frontend/src/routes/[lang=lang]/aspirasi/+page.svelte`, `frontend/src/lib/components/DocumentPreview.svelte`.

## Urutan pengerjaan yang disarankan

1. **Perbaikan langsung:** sesuaikan janji pencarian, tampilkan contoh/bantuan input, perjelas nama tab dan status kelengkapan data, serta jadikan tur opsional.
2. **Alur inti:** pertegas tindakan pertama di beranda, bawa wilayah/jenjang/isu ke tujuan berikutnya, dan tambahkan ringkasan angka dalam bahasa sehari-hari.
3. **Pemulihan dan aksi:** selesaikan pilihan lokasi manual untuk Lapor dan tahap penyampaian aspirasi setelah dokumen siap.
4. **Penyempurnaan mobile:** validasi navigasi bawah, label sumber, ruang layar pertama, serta peta pada perangkat nyata setelah alur inti dibenahi.

## Rancangan alur yang lebih mudah dipahami

> Pilih wilayah → baca ringkasan kondisi → lihat sumber/detail bila perlu → pilih pertanyaan atau masalah → susun aspirasi/laporan → pahami langkah berikutnya.

Di halaman wilayah, tiga pertanyaan utama dapat menjadi pengantar:

- **Sekolah:** cukupkah kursinya?
- **Proyek publik:** apa yang dibiayai di sini?
- **Pihak berwenang:** siapa yang bisa saya tanyai?

Setiap bagian sebaiknya dimulai dari satu jawaban ringkas, lalu rincian dan tindakan yang relevan. Tur, grafik, peta, dan sumber memperdalam pemahaman setelah jawaban utama terbaca.

## Validasi dengan pengguna

Lakukan sesi awal dengan lima orang yang belum pernah melihat proyek, termasuk warga nonteknis dan orang tua calon siswa. Jumlah ini untuk menemukan masalah awal, bukan mengukur keberhasilan seluruh populasi. Beri tugas tanpa menyebut menu yang harus dipakai:

| Tugas | Target awal yang diusulkan |
| --- | --- |
| Jelaskan fungsi Rantai setelah melihat layar pertama | Minimal 4 dari 5 dapat menyebut manfaat konkret dalam 10 detik |
| Cari wilayah sendiri | Minimal 4 dari 5 selesai tanpa arahan dalam 30 detik |
| Jelaskan arti kekurangan kursi | Minimal 4 dari 5 memahami indikator agregat dan tidak menganggapnya kepastian gagal seleksi |
| Temukan proyek di wilayah dan sumber nilainya | Minimal 4 dari 5 selesai tanpa bantuan dalam 60 detik |
| Cari siapa yang berwenang dan siapkan aspirasi | Minimal 4 dari 5 mencapai draf dalam 2 menit dan tahu draf belum dikirim |
| Lanjutkan penyusunan laporan ketika GPS ditolak | Jalur pemulihan dapat ditemukan dan digunakan tanpa mengubah izin browser |

Catat waktu, klik pertama, salah tafsir, langkah mundur, dan permintaan bantuan. Jalankan dengan serta tanpa tur untuk menilai apakah antarmuka utamanya sudah menjelaskan fungsi produk.

## Catatan QA teknis berikutnya

- Konfigurasi Playwright saat ini menandai tur selesai melalui localStorage. Karena itu, jalur E2E yang ada tidak menguji pengalaman kunjungan pertama saat tur otomatis aktif.
- Tes pencarian yang lulus memeriksa ArrowDown → Enter dan hasil kosong; belum memvalidasi ekspektasi mengetik lalu langsung Enter pada perangkat nyata.
- Audit axe yang ada menyaring dampak serius/kritis. Kelulusan tes otomatis tidak membuktikan seluruh aksesibilitas atau kemudahan memahami produk.
- Sesi browser berikutnya perlu memeriksa minimal lebar 320/390/768/1440 px, zoom 200%, navigasi keyboard, jaringan lambat/gagal, tidak ada hasil/data, serta izin lokasi ditolak. Ini daftar yang belum dijalankan pada audit ini.
