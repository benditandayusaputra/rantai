# DESIGN.md · Rantai

Panduan implementasi antarmuka Rantai, web transparansi publik untuk warga DKI Jakarta, untuk tim frontend yang memakai SvelteKit 2, Svelte 5 (runes), TypeScript, dan Tailwind CSS.

Dokumen ini adalah sumber kebenaran. Kanvas desain (halaman 00 Sistem desain sampai 18 Status khusus) adalah acuan visual. Bila ada perbedaan, ikuti dokumen ini. Dua perbedaan yang sudah diketahui:

- Nama token di kanvas lebih pendek (misalnya `--resmi`, `--s-edu`, `--g1`). Nama di dokumen ini (`--cert-resmi`, `--sector-edu`, `--gap-1`) yang dipakai di kode.
- Kanvas menulis uang sebagai "Rp 2,41 M" dengan spasi. Format yang berlaku ada di bagian 9: "Rp2,41 miliar" atau "Rp2,41 M" tanpa spasi.

Prinsip yang mengikat semua keputusan:

1. **Angka adalah bintang.** Angka besar, tabular, selalu dengan satuan dan lencana sumber dalam satu tatapan.
2. **Sumber selalu menempel.** Setiap angka punya label kepastian (resmi, proksi, estimasi, simulasi) yang bisa dibuka menjadi panel sumber.
3. **Rantai untuk alur dan status.** Mata rantai tersambung, putus, atau menunggu. Dipakai sebagai bahasa status, bukan hiasan.
4. **Netral.** Amber untuk sinyal yang perlu ditanyakan. Merah hanya untuk error sistem. Tidak ada warna partai sebagai warna utama.
5. **Ponsel 360 px adalah kasus utama.** Target sentuh minimal 44 × 44 px. Fokus papan ketik selalu terlihat. Setiap grafik dan peta punya alternatif tabel.

## Daftar isi

1. [Token warna](#1-token-warna)
2. [Token tipografi, jarak, radius, bayangan, garis, animasi, z-index](#2-token-tipografi-jarak-radius-bayangan-garis-animasi-z-index)
3. [Konfigurasi Tailwind](#3-konfigurasi-tailwind)
4. [Font dan ikon](#4-font-dan-ikon)
5. [Komponen](#5-komponen)
6. [Tata letak setiap layar](#6-tata-letak-setiap-layar)
7. [Pola interaksi](#7-pola-interaksi)
8. [Teks antarmuka (i18n)](#8-teks-antarmuka-i18n)
9. [Aturan konten](#9-aturan-konten)

---

## 1. Token warna

### 1.1 Cara kerja tema

- Tema bawaan mengikuti pengaturan sistem perangkat. Pilihan manual pengguna disimpan di `localStorage` dengan kunci `rantai-theme` (`light`, `dark`, atau `system`).
- Atribut `data-theme` dipasang **hanya di elemen `<html>`**, sebelum halaman dilukis, agar tidak berkedip. Token alias (status laporan, ring fokus) bergantung pada aturan ini.
- Komponen tidak pernah memakai heks mentah. Selalu lewat token atau kelas Tailwind yang memetakan token.

`src/app.html`, di dalam `<head>` sebelum `%sveltekit.head%`:

```html
<script>
  (function () {
    var pref = 'system';
    try { pref = localStorage.getItem('rantai-theme') || 'system'; } catch (e) {}
    var dark = pref === 'dark' || (pref === 'system' && matchMedia('(prefers-color-scheme: dark)').matches);
    document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
  })();
</script>
```

Saat preferensi `system`, pasang listener `matchMedia('(prefers-color-scheme: dark)').addEventListener('change', ...)` di `+layout.svelte` untuk mengganti atribut secara langsung.

### 1.2 Token (CSS custom properties)

`src/lib/styles/tokens.css`:

```css
/* ============ TEMA TERANG (bawaan) ============ */
:root,
[data-theme='light'] {
  color-scheme: light;

  /* Permukaan */
  --bg: #EFF1EE;            /* latar halaman */
  --surface: #F9FAF8;       /* kartu, panel, tabel */
  --raised: #FFFFFF;        /* dialog, lembar bawah, menu, saran */
  --sunken: #E6E9E5;        /* kepala tabel, sumur, segmen */
  --glass: rgba(249, 250, 248, 0.82); /* hanya panel di atas peta, blur 12px */
  --scrim: rgba(14, 19, 22, 0.45);

  /* Teks dan garis */
  --text: #14191C;          /* 15,6:1 di --bg */
  --muted: #4E585E;         /* 6,4:1 di --bg */
  --line: #D5DAD6;          /* dekoratif, bukan untuk teks */
  --line-strong: #7D868B;   /* tepi input dan kontrol, 3,3:1 */
  --focus: #1D5FD1;         /* cincin fokus */

  /* Aksi */
  --primary: #0F5964;       /* 7,0:1 di --bg */
  --primary-hover: #0B4952;
  --primary-press: #083940;
  --on-primary: #FFFFFF;    /* 8,0:1 di --primary */
  --secondary: #E2ECEB;
  --secondary-hover: #D3E3E2;
  --on-secondary: #0F4E57;  /* 7,8:1 di --secondary, juga warna tautan */
  --dis-bg: #E3E6E2;
  --dis-text: #8A9296;      /* nonaktif, dikecualikan WCAG, selalu dengan alasan */

  /* Sinyal */
  --success: #1B7542;
  --on-success: #FFFFFF;
  --success-bg: #E2F2E8;
  --info: #2556B8;
  --info-bg: #E5ECFA;
  --attn: #855400;          /* teks sinyal "perlu ditanyakan" */
  --attn-fill: #E9A93A;     /* tepi chip, titik peta */
  --attn-bg: #FCF0D6;
  --on-attn-fill: #14191C;
  --danger: #B42318;        /* HANYA error sistem */
  --danger-bg: #FBE7E5;

  /* Label kepastian */
  --cert-resmi: #0F5964;
  --cert-resmi-on: #FFFFFF;
  --cert-proksi: #2556B8;
  --cert-proksi-bg: #E5ECFA;
  --cert-estimasi: #5A636A;
  --cert-sim: #8E1B5E;
  --cert-sim-bg: #FBE5F1;
  --cert-sim-stripe: #EBB9D5;

  /* Sektor proyek (turunan Okabe-Ito) */
  --sector-edu: #0072B2;    --sector-edu-on: #FFFFFF;
  --sector-air: #56B4E9;    --sector-air-on: #0E1316;
  --sector-san: #009E73;    --sector-san-on: #0E1316;
  --sector-jalan: #D55E00;  --sector-jalan-on: #FFFFFF;
  --sector-kes: #CC79A7;    --sector-kes-on: #0E1316;
  --sector-ruang: #E69F00;  --sector-ruang-on: #0E1316;
  --sector-lain: #9A9A9A;   --sector-lain-on: #0E1316;

  /* Skala choropleth kesenjangan kursi, 4 kelas */
  --gap-1: #2A3F7E;  --gap-1-on: #FFFFFF;  /* sangat kurang, + pola titik */
  --gap-2: #3D86A8;  --gap-2-on: #FFFFFF;  /* kurang */
  --gap-3: #7FBFC4;  --gap-3-on: #14191C;  /* hampir cukup */
  --gap-4: #C9E4DE;  --gap-4-on: #14191C;  /* cukup */
  --gap-dot: rgba(255, 255, 255, 0.45);

  /* Peta dasar */
  --map: #E6E9E5;
  --map-road: #FFFFFF;
  --map-water: #C9DCE3;
  --map-park: #D7E3D3;

  /* Kerangka memuat */
  --skel: #E2E5E1;
  --skel-hi: #F2F4F1;

  /* Bayangan (beda per tema) */
  --sh-1: 0 1px 2px rgba(20, 25, 28, 0.06), 0 1px 1px rgba(20, 25, 28, 0.04);
  --sh-2: 0 8px 20px -6px rgba(20, 25, 28, 0.14), 0 2px 4px rgba(20, 25, 28, 0.05);
  --sh-3: 0 24px 56px -16px rgba(20, 25, 28, 0.28), 0 4px 10px rgba(20, 25, 28, 0.06);
}

/* ============ TEMA GELAP ============ */
[data-theme='dark'] {
  color-scheme: dark;

  --bg: #0E1316;
  --surface: #151B1F;
  --raised: #1D252A;
  --sunken: #0A0E10;
  --glass: rgba(21, 27, 31, 0.8);
  --scrim: rgba(0, 0, 0, 0.6);

  --text: #E6EBE9;          /* 15,5:1 di --bg */
  --muted: #A1ACB1;         /* 8,1:1 di --bg */
  --line: #2A343A;
  --line-strong: #66737A;   /* 3,8:1 */
  --focus: #84A9FF;

  --primary: #6CC3C4;       /* 9,1:1 di --bg */
  --primary-hover: #86D1D1;
  --primary-press: #57AEAF;
  --on-primary: #062326;    /* 8,0:1 di --primary */
  --secondary: #1F3236;
  --secondary-hover: #27403F;
  --on-secondary: #A9DCDC;
  --dis-bg: #20282D;
  --dis-text: #5F6B71;

  --success: #5EC58C;
  --on-success: #062014;
  --success-bg: #13291D;
  --info: #8DB1FF;
  --info-bg: #162344;
  --attn: #F0B54A;
  --attn-fill: #E9A93A;
  --attn-bg: #2E2410;
  --on-attn-fill: #14191C;
  --danger: #FF8A7E;
  --danger-bg: #3A1714;

  --cert-resmi: #6CC3C4;
  --cert-resmi-on: #062326;
  --cert-proksi: #8DB1FF;
  --cert-proksi-bg: #162344;
  --cert-estimasi: #A1ACB1;
  --cert-sim: #F49ACB;
  --cert-sim-bg: #3A1229;
  --cert-sim-stripe: #5A1E40;

  --sector-edu: #3D9BE0;    --sector-edu-on: #0E1316;
  --sector-air: #7CC8EE;    --sector-air-on: #0E1316;
  --sector-san: #2FBF8F;    --sector-san-on: #0E1316;
  --sector-jalan: #F07C3A;  --sector-jalan-on: #0E1316;
  --sector-kes: #DE97BE;    --sector-kes-on: #0E1316;
  --sector-ruang: #F0B23A;  --sector-ruang-on: #0E1316;
  --sector-lain: #A0A4A8;   --sector-lain-on: #0E1316;

  /* Urutan terang-gelap dibalik: makin kontras dengan peta dasar, makin kurang */
  --gap-1: #BFD9F2;  --gap-1-on: #0E1316;
  --gap-2: #58A7C6;  --gap-2-on: #0E1316;
  --gap-3: #2C7381;  --gap-3-on: #E6EBE9;
  --gap-4: #23423F;  --gap-4-on: #E6EBE9;
  --gap-dot: rgba(14, 19, 22, 0.5);

  --map: #1A2125;
  --map-road: #283237;
  --map-water: #11262E;
  --map-park: #1C2A24;

  --skel: #1F272C;
  --skel-hi: #29333A;

  --sh-1: 0 1px 2px rgba(0, 0, 0, 0.4);
  --sh-2: 0 8px 20px -6px rgba(0, 0, 0, 0.55), 0 0 0 1px rgba(255, 255, 255, 0.04);
  --sh-3: 0 24px 56px -16px rgba(0, 0, 0, 0.7), 0 0 0 1px rgba(255, 255, 255, 0.05);
}

/* ============ ALIAS (berlaku di kedua tema) ============ */
:root {
  /* Cincin fokus: 2px jarak dengan warna permukaan, lalu 2px --focus */
  --ring: 0 0 0 2px var(--surface), 0 0 0 4px var(--focus);

  /* Status laporan warga */
  --st-sent-fg: var(--muted);         --st-sent-bg: transparent;           --st-sent-bd: var(--line-strong);  /* garis utuh */
  --st-wait-fg: var(--info);          --st-wait-bg: transparent;           --st-wait-bd: var(--info);         /* garis putus */
  --st-ver-fg: var(--success);        --st-ver-bg: var(--success-bg);      --st-ver-bd: transparent;
  --st-fwd-fg: var(--on-secondary);   --st-fwd-bg: var(--secondary);      --st-fwd-bd: transparent;
  --st-resp-fg: var(--info);          --st-resp-bg: var(--info-bg);        --st-resp-bd: transparent;
  --st-done-fg: var(--on-success);    --st-done-bg: var(--success);        --st-done-bd: transparent;
  --st-rej-fg: var(--surface);        --st-rej-bg: var(--text);            --st-rej-bd: transparent;          /* tinta netral, bukan merah */

  /* Mata rantai */
  --link-ok: var(--primary);
  --link-broken: var(--attn-fill);
  --link-wait: var(--line-strong);
  --link-est: var(--cert-estimasi);
}
```

### 1.3 Ringkasan token semantik

| Kelompok | Token | Terang | Gelap | Pemakaian |
|---|---|---|---|---|
| Latar | `--bg` | #EFF1EE | #0E1316 | Kanvas halaman |
| Permukaan | `--surface` | #F9FAF8 | #151B1F | Kartu, panel, tabel, bilah navigasi |
| Permukaan terangkat | `--raised` | #FFFFFF | #1D252A | Dialog, lembar bawah, menu, saran |
| Cekung | `--sunken` | #E6E9E5 | #0A0E10 | Kepala tabel, segmen, sumur |
| Teks utama | `--text` | #14191C | #E6EBE9 | Judul, isi, angka |
| Teks redup | `--muted` | #4E585E | #A1ACB1 | Keterangan, satuan, metadata |
| Garis | `--line` | #D5DAD6 | #2A343A | Pemisah, tepi kartu |
| Garis tegas | `--line-strong` | #7D868B | #66737A | Tepi input, kontrol, chip |
| Fokus | `--focus` | #1D5FD1 | #84A9FF | Cincin fokus |
| Aksi utama | `--primary` / `--on-primary` | #0F5964 / #FFFFFF | #6CC3C4 / #062326 | Tombol utama, tab aktif, mata rantai tersambung |
| Aksi sekunder | `--secondary` / `--on-secondary` | #E2ECEB / #0F4E57 | #1F3236 / #A9DCDC | Tombol sekunder, chip aktif, tautan |
| Berhasil | `--success` | #1B7542 | #5EC58C | Terverifikasi, tersimpan |
| Info | `--info` | #2556B8 | #8DB1FF | Catatan metodologi, hak jawab |
| Perhatian | `--attn` / `--attn-fill` | #855400 / #E9A93A | #F0B54A / #E9A93A | Sinyal "perlu ditanyakan", mata rantai putus |
| Bahaya | `--danger` | #B42318 | #FF8A7E | Error sistem saja |

### 1.4 Label kepastian

Label harus bisa dibedakan tanpa warna. Bentuk dan ikon wajib, teks wajib.

| Label | Token | Bentuk | Ikon (Lucide) | Arti |
|---|---|---|---|---|
| Resmi | `--cert-resmi`, `--cert-resmi-on` | Isi penuh, sudut 6px | `badge-check` | Diambil langsung dari sumber pemerintah tanpa diubah |
| Proksi | `--cert-proksi`, `--cert-proksi-bg` | Garis utuh 1,5px + isi tipis | `sigma` | Dihitung Rantai dengan rumus terbuka |
| Estimasi | `--cert-estimasi` | Garis putus 1,5px, tanpa isi | `approximately-equal` | Perkiraan karena data resmi tidak dipublikasikan |
| Simulasi | `--cert-sim`, `--cert-sim-bg`, `--cert-sim-stripe` | Arsir 45° + garis utuh, wajib pita | `flask-conical` | Data demonstrasi, tidak boleh dikutip |

Arsir simulasi: `repeating-linear-gradient(-45deg, var(--cert-sim-bg) 0 4px, var(--cert-sim-stripe) 4px 8px)`. Teks di atas arsir tetap 5:1 atau lebih (terang 5,0:1, gelap 6,1:1).

### 1.5 Sektor proyek

Warna sektor tidak pernah berdiri sendiri: selalu dengan ikon dan label, dan segmen grafik dipisah celah 2px.

| Sektor | Kunci | Token | Terang | Gelap | Ikon |
|---|---|---|---|---|---|
| Pendidikan | `edu` | `--sector-edu` | #0072B2 | #3D9BE0 | `graduation-cap` |
| Air | `air` | `--sector-air` | #56B4E9 | #7CC8EE | `droplet` |
| Sanitasi | `san` | `--sector-san` | #009E73 | #2FBF8F | `waves` |
| Jalan | `jalan` | `--sector-jalan` | #D55E00 | #F07C3A | `route` |
| Kesehatan | `kes` | `--sector-kes` | #CC79A7 | #DE97BE | `heart-pulse` |
| Ruang publik | `ruang` | `--sector-ruang` | #E69F00 | #F0B23A | `trees` |
| Lainnya | `lain` | `--sector-lain` | #9A9A9A | #A0A4A8 | `package` |

### 1.6 Status laporan warga

| Status | Kunci | Token | Bentuk | Ikon |
|---|---|---|---|---|
| Dikirim | `dikirim` | `--st-sent-*` | Garis utuh netral | `circle-arrow-up` |
| Menunggu saksi | `menunggu-saksi` | `--st-wait-*` | Garis putus info + hitungan "2/3" | `users` |
| Terverifikasi | `terverifikasi` | `--st-ver-*` | Isi tipis hijau | `badge-check` |
| Diteruskan | `diteruskan` | `--st-fwd-*` | Isi sekunder | `forward` |
| Ditanggapi | `ditanggapi` | `--st-resp-*` | Isi tipis info | `reply` |
| Selesai | `selesai` | `--st-done-*` | Isi penuh hijau | `circle-check` |
| Ditolak | `ditolak` | `--st-rej-*` | Isi tinta netral, selalu dengan alasan | `circle-x` |

### 1.7 Skala choropleth kesenjangan kursi

Rasio = calon lulusan jenjang bawah ÷ kursi kelas awal negeri di wilayah yang sama.

| Kelas | Kunci | Rentang rasio | Token | Terang | Gelap | Tambahan |
|---|---|---|---|---|---|---|
| Sangat kurang | `sangat-kurang` | lebih dari 2,0 | `--gap-1` | #2A3F7E | #BFD9F2 | Pola titik `--gap-dot` 6px |
| Kurang | `kurang` | 1,3 sampai 2,0 | `--gap-2` | #3D86A8 | #58A7C6 | |
| Hampir cukup | `hampir-cukup` | 1,0 sampai 1,3 | `--gap-3` | #7FBFC4 | #2C7381 | |
| Cukup | `cukup` | 1,0 atau kurang | `--gap-4` | #C9E4DE | #23423F | |

Aturan: makin kontras dengan peta dasar, makin kurang kursinya. Di tema gelap urutan terang dan gelap dibalik agar aturan ini tetap berlaku. Setiap area diberi garis tepi 2px warna `--map` agar kelas bersebelahan tetap terpisah.

### 1.8 Ilustrasi

Ilustrasi beranda (`IlustrasiKota`) dan halaman masuk (`IlustrasiMasuk`) adalah SVG inline. Bidang memakai token yang sudah ada (`--secondary`, `--secondary-hover`, `--raised`, `--line`, `--primary`, `--attn-fill`, `--gap-*`), sehingga ikut berganti tema. Hanya warna tokoh yang punya token sendiri:

| Token | Terang | Gelap | Kelas Tailwind |
|---|---|---|---|
| `--illu-skin` | #E9B48E | #D9A27C | `fill-illu-skin` |
| `--illu-skin-shade` | #C98B64 | #B07A56 | `fill-illu-shade` |
| `--illu-hair` | #2A211D | #120D0B | `fill-illu-hair`, `stroke-illu-hair` |
| `--illu-pants` | #2F3D45 | #46575F | `fill-illu-pants` |

Ilustrasi selalu `aria-hidden="true"`, tidak memuat teks, dan tidak menggambarkan angka yang bisa dikira data sungguhan.

---

## 2. Token tipografi, jarak, radius, bayangan, garis, animasi, z-index

`src/lib/styles/tokens.css` (lanjutan):

```css
:root {
  /* Keluarga huruf */
  --font-sans: 'Plus Jakarta Sans Variable', 'Plus Jakarta Sans', ui-sans-serif, system-ui, -apple-system, 'Segoe UI', sans-serif;
  --font-serif: 'Newsreader Variable', 'Newsreader', Georgia, 'Times New Roman', serif;
  --font-mono: 'JetBrains Mono', ui-monospace, Menlo, Consolas, monospace;

  /* Radius */
  --r-xs: 4px;   /* tooltip */
  --r-sm: 6px;   /* lencana */
  --r-md: 10px;  /* tombol, input */
  --r-lg: 14px;  /* kartu */
  --r-xl: 20px;  /* lembar bawah (sudut atas) */
  --r-pill: 999px;

  /* Ketebalan garis */
  --bw-hair: 1px;    /* pemisah, tepi kartu */
  --bw-ctrl: 1.5px;  /* input, chip, lencana */
  --bw-focus: 2px;   /* cincin fokus, input fokus */
  --bw-link: 3px;    /* mata rantai */

  /* Durasi */
  --dur-instant: 80ms;   /* tekan tombol */
  --dur-fast: 140ms;     /* hover, fokus, chip */
  --dur-base: 220ms;     /* tab, toggle, tooltip */
  --dur-enter: 320ms;    /* lembar bawah, dialog */
  --dur-link-step: 60ms; /* jeda antar mata rantai */

  /* Kurva */
  --ease-standard: cubic-bezier(0.2, 0, 0, 1);
  --ease-enter: cubic-bezier(0.05, 0.7, 0.1, 1);
  --ease-exit: cubic-bezier(0.3, 0, 0.8, 0.15);

  /* Lapisan */
  --z-base: 0;
  --z-raised: 10;
  --z-sticky: 100;   /* header, bilah bawah */
  --z-map-ui: 200;   /* panel di atas peta */
  --z-popover: 300;  /* saran pencarian, menu, panel sumber */
  --z-sheet: 400;    /* lembar bawah */
  --z-dialog: 500;   /* dialog dan scrim */
  --z-toast: 600;
  --z-tooltip: 700;
  --z-ribbon: 800;   /* pita mode demo */
}

@media (prefers-reduced-motion: reduce) {
  :root {
    --dur-fast: 120ms; --dur-base: 120ms; --dur-enter: 120ms; --dur-link-step: 0ms;
  }
}
```

### 2.1 Skala tipografi

Ukuran berubah di titik henti `lg` (1024px). Format: ukuran / tinggi baris dalam px.

| Token | Keluarga | Berat | Ponsel (< 1024) | Desktop (≥ 1024) | Spasi huruf | Pemakaian |
|---|---|---|---|---|---|---|
| `display` | Newsreader | 500 | 34 / 40 | 60 / 64 | -0,015em | Judul beranda, halaman status |
| `h1` | Newsreader | 500 | 28 / 34 | 44 / 50 | -0,01em | Judul halaman |
| `h2` | Plus Jakarta Sans | 700 | 22 / 28 | 30 / 36 | -0,01em | Judul bagian |
| `h3` | Plus Jakarta Sans | 700 | 18 / 24 | 22 / 28 | 0 | Judul kartu |
| `body-lg` | Plus Jakarta Sans | 400 | 17 / 26 | 18 / 28 | 0 | Paragraf pembuka |
| `body` | Plus Jakarta Sans | 400 | 16 / 24 | 16 / 24 | 0 | Isi |
| `small` | Plus Jakarta Sans | 400 sampai 600 | 14 / 20 | 14 / 20 | 0 | Metadata, isi kartu |
| `caption` | Plus Jakarta Sans | 500 | 12 / 16 | 12 / 16 | 0 | Keterangan, sumber |
| `overline` | Plus Jakarta Sans | 700, kapital | 12 / 16 | 12 / 16 | 0,08em | Label bagian |
| `num-xl` | Plus Jakarta Sans | 700, tnum | 44 / 48 | 72 / 72 | -0,03em | Angka utama halaman |
| `num-lg` | Plus Jakarta Sans | 700, tnum | 30 / 36 | 44 / 48 | -0,02em | Kartu angka kunci |
| `num-md` | Plus Jakarta Sans | 600, tnum | 20 / 26 | 24 / 30 | -0,01em | Angka di kartu ringkas |
| `mono` | JetBrains Mono | 400 | 13 / 20 | 13 / 20 | 0 | Rumus, kode RUP |

Aturan:

- Ukuran teks minimum 12px.
- Semua angka data memakai `font-variant-numeric: tabular-nums lining-nums` (kelas `.num`).
- Kalimat ringkas editorial ("Sekitar 2 lulusan untuk setiap 1 kursi") memakai Newsreader 500, 26 sampai 30px, dengan angka di dalamnya tetap Plus Jakarta Sans 700 tabular.
- Satuan di samping angka: 15px, 600, `--muted`.

### 2.2 Jarak (basis 4)

Skala bawaan Tailwind sudah berbasis 4px dan dipakai apa adanya.

| Token | px | Kelas Tailwind | Pemakaian umum |
|---|---|---|---|
| space-1 | 4 | `1` | Ikon ke teks rapat |
| space-2 | 8 | `2` | Antarchip, ikon ke teks |
| space-3 | 12 | `3` | Celah grid ponsel, isi lencana besar |
| space-4 | 16 | `4` | Margin ponsel, isi kartu ponsel |
| space-5 | 20 | `5` | Isi kartu tablet ke atas |
| space-6 | 24 | `6` | Antarblok dalam bagian, margin tablet |
| space-8 | 32 | `8` | Margin laptop |
| space-10 | 40 | `10` | Jarak bagian tablet |
| space-12 | 48 | `12` | |
| space-14 | 56 | `14` | Jarak bagian desktop |
| space-16 | 64 | `16` | |
| space-18 | 72 | `18` | Margin desktop 1440 |
| space-24 | 96 | `24` | |

### 2.3 Titik henti dan kisi

| Nama | Lebar minimum | Kolom | Margin | Celah | Lebar konten maks |
|---|---|---|---|---|---|
| (bawaan) | 0 (desain 360) | 4 | 16 | 12 | penuh |
| `md` | 768 | 8 | 24 | 16 | penuh |
| `lg` | 1024 | 12 | 32 | 20 | penuh |
| `xl` | 1360 (desain 1440) | 12 | 72 | 24 | 1296 |

Navigasi berganti dari bilah bawah ke bilah atas di `lg`. Tablet 768 tetap memakai pola ponsel (bilah bawah) dengan konten dua kolom.

### 2.4 Radius, bayangan, garis

| Token | Nilai | Pemakaian |
|---|---|---|
| `--r-xs` | 4px | Tooltip |
| `--r-sm` | 6px | Lencana kepastian, lencana provenance |
| `--r-md` | 10px | Tombol, input, opsi |
| `--r-lg` | 14px | Kartu, panel |
| `--r-xl` | 20px | Sudut atas lembar bawah |
| `--r-pill` | 999px | Chip, mata rantai, status |
| `--sh-1` | lihat token | Kartu yang bisa diklik |
| `--sh-2` | lihat token | Menu, saran, tooltip, panel peta |
| `--sh-3` | lihat token | Dialog, lembar bawah |
| `--ring` | 2px celah + 2px `--focus` | Semua fokus papan ketik (`:focus-visible`) |
| `--bw-link` | 3px | Mata rantai di stepper, mini rantai, linimasa |

Permukaan tembus pandang dengan blur (`--glass`, `backdrop-filter: blur(12px)`) hanya untuk panel yang melayang di atas peta. Di tempat lain permukaan selalu padat.

### 2.5 Animasi

| Situasi | Durasi | Kurva |
|---|---|---|
| Tekan tombol (translateY 1px) | `--dur-instant` | standard |
| Hover, fokus, chip | `--dur-fast` | standard |
| Tab, toggle, tooltip | `--dur-base` | standard |
| Lembar bawah, dialog, toast masuk | `--dur-enter` | enter |
| Elemen ditutup | `--dur-base` | exit |
| Mata rantai tersambung bertahap | 60ms × indeks, tiap mata `--dur-base` | enter |
| Shimmer kerangka | 1,4s linear, berulang | linear |

Saat `prefers-reduced-motion: reduce`: semua transisi menjadi pudar 120ms, rantai tampil sekaligus, shimmer dimatikan (kerangka statis).

### 2.6 Z-index

| Token | Nilai | Elemen |
|---|---|---|
| `--z-base` | 0 | Konten |
| `--z-raised` | 10 | Kartu terangkat, penanda terpilih |
| `--z-sticky` | 100 | Bilah atas, bilah bawah, header lengket |
| `--z-map-ui` | 200 | Kontrol lapisan, legenda, tombol peta |
| `--z-popover` | 300 | Saran pencarian, menu, panel sumber |
| `--z-sheet` | 400 | Lembar bawah |
| `--z-dialog` | 500 | Dialog dan scrim |
| `--z-toast` | 600 | Toast |
| `--z-tooltip` | 700 | Tooltip |
| `--z-ribbon` | 800 | Pita mode demo |

---

## 3. Konfigurasi Tailwind

Konfigurasi di bawah untuk Tailwind CSS v3.4. Warna dipetakan ke `var(--token)` sehingga tema berganti hanya dengan atribut `data-theme` di `<html>`, tanpa kelas `dark:` ganda. Varian `dark:` tetap tersedia untuk kasus langka (misalnya gambar khusus tema gelap).

Konsekuensi: pengubah opasitas (`bg-primary/50`) tidak berlaku untuk warna berbasis var heks. Bila butuh transparansi, pakai token khusus (`--glass`, `--scrim`) atau `color-mix()` di CSS.

`tailwind.config.ts`:

```ts
import type { Config } from 'tailwindcss';

const v = (name: string) => `var(--${name})`;

const sectors = ['edu', 'air', 'san', 'jalan', 'kes', 'ruang', 'lain'] as const;
const sectorColors = Object.fromEntries(
  sectors.map((s) => [s, { DEFAULT: v(`sector-${s}`), on: v(`sector-${s}-on`) }])
);

const statuses = ['sent', 'wait', 'ver', 'fwd', 'resp', 'done', 'rej'] as const;
const statusColors = Object.fromEntries(
  statuses.map((s) => [s, { fg: v(`st-${s}-fg`), bg: v(`st-${s}-bg`), bd: v(`st-${s}-bd`) }])
);

export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: ['selector', '[data-theme="dark"]'],
  theme: {
    screens: {
      md: '768px',
      lg: '1024px',
      xl: '1360px'
    },
    extend: {
      colors: {
        page: v('bg'),
        surface: v('surface'),
        raised: v('raised'),
        sunken: v('sunken'),
        glass: v('glass'),
        scrim: v('scrim'),
        fg: v('text'),
        muted: v('muted'),
        line: { DEFAULT: v('line'), strong: v('line-strong') },
        focus: v('focus'),
        primary: { DEFAULT: v('primary'), hover: v('primary-hover'), press: v('primary-press'), on: v('on-primary') },
        secondary: { DEFAULT: v('secondary'), hover: v('secondary-hover'), on: v('on-secondary') },
        dis: { bg: v('dis-bg'), fg: v('dis-text') },
        success: { DEFAULT: v('success'), on: v('on-success'), bg: v('success-bg') },
        info: { DEFAULT: v('info'), bg: v('info-bg') },
        attn: { DEFAULT: v('attn'), fill: v('attn-fill'), bg: v('attn-bg'), on: v('on-attn-fill') },
        danger: { DEFAULT: v('danger'), bg: v('danger-bg') },
        cert: {
          resmi: { DEFAULT: v('cert-resmi'), on: v('cert-resmi-on') },
          proksi: { DEFAULT: v('cert-proksi'), bg: v('cert-proksi-bg') },
          estimasi: v('cert-estimasi'),
          sim: { DEFAULT: v('cert-sim'), bg: v('cert-sim-bg'), stripe: v('cert-sim-stripe') }
        },
        sector: sectorColors,
        st: statusColors,
        gap: {
          1: { DEFAULT: v('gap-1'), on: v('gap-1-on') },
          2: { DEFAULT: v('gap-2'), on: v('gap-2-on') },
          3: { DEFAULT: v('gap-3'), on: v('gap-3-on') },
          4: { DEFAULT: v('gap-4'), on: v('gap-4-on') }
        },
        map: { DEFAULT: v('map'), road: v('map-road'), water: v('map-water'), park: v('map-park') },
        skel: { DEFAULT: v('skel'), hi: v('skel-hi') }
      },
      fontFamily: {
        sans: [v('font-sans')],
        serif: [v('font-serif')],
        mono: [v('font-mono')]
      },
      fontSize: {
        display: ['34px', { lineHeight: '40px', letterSpacing: '-0.015em', fontWeight: '500' }],
        'display-lg': ['60px', { lineHeight: '64px', letterSpacing: '-0.015em', fontWeight: '500' }],
        h1: ['28px', { lineHeight: '34px', letterSpacing: '-0.01em', fontWeight: '500' }],
        'h1-lg': ['44px', { lineHeight: '50px', letterSpacing: '-0.01em', fontWeight: '500' }],
        h2: ['22px', { lineHeight: '28px', letterSpacing: '-0.01em', fontWeight: '700' }],
        'h2-lg': ['30px', { lineHeight: '36px', letterSpacing: '-0.01em', fontWeight: '700' }],
        h3: ['18px', { lineHeight: '24px', fontWeight: '700' }],
        'h3-lg': ['22px', { lineHeight: '28px', fontWeight: '700' }],
        'body-lg': ['17px', { lineHeight: '26px' }],
        'body-lg-lg': ['18px', { lineHeight: '28px' }],
        body: ['16px', { lineHeight: '24px' }],
        small: ['14px', { lineHeight: '20px' }],
        caption: ['12px', { lineHeight: '16px', fontWeight: '500' }],
        overline: ['12px', { lineHeight: '16px', letterSpacing: '0.08em', fontWeight: '700' }],
        'num-xl': ['44px', { lineHeight: '48px', letterSpacing: '-0.03em', fontWeight: '700' }],
        'num-xl-lg': ['72px', { lineHeight: '72px', letterSpacing: '-0.03em', fontWeight: '700' }],
        'num-lg': ['30px', { lineHeight: '36px', letterSpacing: '-0.02em', fontWeight: '700' }],
        'num-lg-lg': ['44px', { lineHeight: '48px', letterSpacing: '-0.02em', fontWeight: '700' }],
        'num-md': ['20px', { lineHeight: '26px', letterSpacing: '-0.01em', fontWeight: '600' }],
        'num-md-lg': ['24px', { lineHeight: '30px', letterSpacing: '-0.01em', fontWeight: '600' }],
        mono: ['13px', { lineHeight: '20px' }]
      },
      spacing: {
        14: '56px',
        18: '72px'
      },
      maxWidth: {
        content: '1296px',
        prose: '720px'
      },
      borderRadius: {
        xs: v('r-xs'),
        sm: v('r-sm'),
        md: v('r-md'),
        lg: v('r-lg'),
        xl: v('r-xl'),
        pill: v('r-pill')
      },
      borderWidth: {
        ctrl: '1.5px',
        focus: '2px',
        link: '3px'
      },
      boxShadow: {
        1: v('sh-1'),
        2: v('sh-2'),
        3: v('sh-3'),
        ring: v('ring')
      },
      transitionDuration: {
        instant: '80ms',
        fast: '140ms',
        base: '220ms',
        enter: '320ms'
      },
      transitionTimingFunction: {
        standard: v('ease-standard'),
        enter: v('ease-enter'),
        exit: v('ease-exit')
      },
      zIndex: {
        raised: '10',
        sticky: '100',
        'map-ui': '200',
        popover: '300',
        sheet: '400',
        dialog: '500',
        toast: '600',
        tooltip: '700',
        ribbon: '800'
      },
      minHeight: { touch: '44px' },
      minWidth: { touch: '44px' },
      keyframes: {
        shimmer: { '0%': { backgroundPosition: '-320px 0' }, '100%': { backgroundPosition: '320px 0' } },
        blink: { '50%': { opacity: '0' } }
      },
      animation: {
        shimmer: 'shimmer 1.4s linear infinite',
        caret: 'blink 1s steps(1) infinite'
      }
    }
  },
  plugins: []
} satisfies Config;
```

Kelas komposit di `src/app.css` agar skala tipografi responsif tidak diulang di setiap komponen:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@import './lib/styles/tokens.css';

@layer base {
  html { @apply bg-page text-fg font-sans antialiased; }
  :focus-visible { outline: none; box-shadow: var(--ring); }
  ::selection { background: var(--secondary); color: var(--on-secondary); }
}

@layer components {
  .num { font-variant-numeric: tabular-nums lining-nums; font-feature-settings: 'tnum' 1, 'lnum' 1; }
  .t-display { @apply font-serif text-display text-balance lg:text-display-lg; }
  .t-h1 { @apply font-serif text-h1 text-balance lg:text-h1-lg; }
  .t-h2 { @apply text-h2 lg:text-h2-lg; }
  .t-h3 { @apply text-h3 lg:text-h3-lg; }
  .t-lead { @apply text-body-lg text-muted lg:text-body-lg-lg; text-wrap: pretty; }
  .t-over { @apply text-overline uppercase text-muted; }
  .n-xl { @apply num text-num-xl lg:text-num-xl-lg; }
  .n-lg { @apply num text-num-lg lg:text-num-lg-lg; }
  .n-md { @apply num text-num-md lg:text-num-md-lg; }
  .wrap { @apply mx-auto w-full max-w-content px-4 md:px-6 lg:px-8 xl:px-0; }
  .glass { background: var(--glass); backdrop-filter: blur(12px); -webkit-backdrop-filter: blur(12px); }
  .hatch-sim { background: repeating-linear-gradient(-45deg, var(--cert-sim-bg) 0 4px, var(--cert-sim-stripe) 4px 8px); }
  .skel { background: linear-gradient(90deg, var(--skel) 0, var(--skel-hi) 40%, var(--skel) 80%); background-size: 640px 100%; @apply animate-shimmer motion-reduce:animate-none; }
}
```

Kelas komponen bersama di `app.css`, dipakai alih-alih merangkai ulang kelas utilitas di setiap halaman:

| Kelas | Isi |
|---|---|
| `.btn` + `.btn-primary` / `.btn-secondary` / `.btn-outline` / `.btn-ghost` / `.btn-text` | Varian `Button` (5.2). Ukuran `.btn-sm` (36), bawaan (44), `.btn-lg` (52); `.btn-block` untuk lebar penuh; `.btn-icon` untuk tombol ikon 44 × 44. Ikon Lucide 18 sampai 20 diletakkan di dalam tombol sebelum label, atau sesudahnya untuk panah lanjut. |
| `.field` | Isian teks 48 px (`TextField`). |
| `.card`, `.card-link` | Kartu `--surface` bergaris `--line` dengan `--sh-1`; `.card-link` menambah naik 2 px dan `--sh-2` saat hover. |
| `.icon-tile` | Kotak ikon 44 × 44 `--secondary` dengan ikon `--primary`, untuk judul bagian dan kartu fitur. |
| `.chip` | Kapsul pilihan 44 px; `aria-pressed="true"` memberi isi `--secondary`. |
| `.link` | Tautan dalam teks bergaris bawah. |

Catatan untuk Tailwind v4: pindahkan nilai `theme.extend` ke blok `@theme` di CSS (misalnya `--color-page: var(--bg);`) dan pakai `@custom-variant dark (&:where([data-theme=dark], [data-theme=dark] *));`. Nama kelas yang dihasilkan sama.

---

## 4. Font dan ikon

### 4.1 Font

| Keluarga | Peran | Berat yang dipakai | Sumber | Lisensi |
|---|---|---|---|---|
| Plus Jakarta Sans | UI, isi, semua angka | 400, 500, 600, 700, 800 (800 hanya logo kata "rantai") | Tokotype, Google Fonts / Fontsource | SIL OFL 1.1 |
| Newsreader | Judul editorial (`display`, `h1`), kalimat ringkas, kutipan hak jawab | 400, 500, 400 italic; sumbu ukuran optik 6 sampai 72 | Production Type, Google Fonts / Fontsource | SIL OFL 1.1 |
| JetBrains Mono | Rumus, kode RUP, ID kejadian | 400, 500 | JetBrains, Google Fonts / Fontsource | SIL OFL 1.1 |

Plus Jakarta Sans dirancang untuk identitas kota Jakarta, sehingga rasa Jakarta hadir lewat huruf, bukan lewat ikon klise. Ketiganya mendukung penuh huruf Latin bahasa Indonesia dan Inggris.

**Cara memuat (disarankan: self-host lewat Fontsource)**, agar tidak ada permintaan ke pihak ketiga dari perangkat warga:

```bash
npm i @fontsource-variable/plus-jakarta-sans @fontsource-variable/newsreader @fontsource/jetbrains-mono
```

```ts
// src/routes/+layout.ts atau +layout.svelte (bagian <script>)
import '@fontsource-variable/plus-jakarta-sans';
import '@fontsource-variable/newsreader';       // periksa dokumentasi paket untuk berkas sumbu opsz dan italic
import '@fontsource/jetbrains-mono/400.css';
import '@fontsource/jetbrains-mono/500.css';
```

Pasang `<link rel="preload" as="font" type="font/woff2" crossorigin>` untuk berkas Plus Jakarta Sans subset latin agar teks utama tidak berganti font di tengah render. Gunakan `font-display: swap`.

**Alternatif (Google Fonts)**, sama dengan kanvas desain:

```html
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Newsreader:ital,opsz,wght@0,6..72,400;0,6..72,500;1,6..72,400&family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;500&display=swap">
```

Angka tabular: aktifkan lewat `.num`. Periksa saat build bahwa fitur `tnum` Plus Jakarta Sans ikut di berkas subset yang dimuat. Bila tidak tersedia, angka di tabel boleh jatuh ke JetBrains Mono sebagai cadangan, bukan ke angka proporsional.

### 4.2 Ikon

- Pustaka: **Lucide** (sumber terbuka, lisensi ISC), paket `@lucide/svelte` untuk Svelte 5.
- Garis 1,75px, ujung dan sambungan membulat. Ukuran 16 di lencana dan chip, 20 di UI, 24 di navigasi.
- Ikon dekoratif: `aria-hidden="true"`. Ikon tanpa teks wajib dibungkus `IconButton` yang punya `label` (aria-label) dan tooltip.
- Ikon tidak pernah menjadi satu-satunya pembawa makna. Status, label kepastian, dan sektor selalu disertai teks.

```bash
npm i @lucide/svelte
```

```svelte
<script lang="ts">
  import { BadgeCheck } from '@lucide/svelte';
</script>
<BadgeCheck size={16} strokeWidth={1.75} aria-hidden="true" />
```

Pemetaan ikon yang dipakai:

| Makna | Lucide | Makna | Lucide |
|---|---|---|---|
| Cari | `search` | Peta | `map` |
| Lokasi | `map-pin` | Pakai lokasiku | `locate-fixed` |
| Tersambung | `link-2` | Putus | `unlink-2` |
| Menunggu | `clock` | Resmi | `badge-check` |
| Proksi | `sigma` | Estimasi | `approximately-equal` |
| Simulasi | `flask-conical` | Sinyal | `circle-help` |
| Info | `info` | Error | `triangle-alert` |
| Tautan asli | `external-link` | Pendidikan | `graduation-cap` |
| Air | `droplet` | Sanitasi | `waves` |
| Jalan | `route` | Kesehatan | `heart-pulse` |
| Ruang publik | `trees` | Lainnya | `package` |
| DPRD | `landmark` | Instansi | `building-2` |
| Saksi | `users` | Jadi saksi | `user-check` |
| Foto bukti, Lapor | `camera` | Tanya | `message-circle` |
| Kirim | `send` | Teruskan | `forward` |
| Tanggapan | `reply` | Hak jawab | `scale` |
| Pemeliharaan | `hourglass` | Sekolah | `school` |
| Uang | `wallet` | Dokumen | `file-text` |
| Lapisan | `layers` | Tabel | `table` |
| Pantau | `bell` | Bagikan | `share-2` |
| Beranda | `house` | Akun | `user` |
| Memuat | `loader-circle` | Coba lagi | `refresh-cw` |
| Salin | `copy` | Unduh | `download` |
| Tema terang | `sun` | Tema gelap | `moon` |
| Bahasa | `globe` | Filter | `sliders-horizontal` |
| Tampilkan | `eye` | Sembunyikan | `eye-off` |
| Luring | `wifi-off` | Tanggal | `calendar` |
| Selesai | `circle-check` | Ditolak | `circle-x` |
| Dikirim | `circle-arrow-up` | Laporkan konten | `flag` |
| Pemenang | `trophy` | Pelaksanaan | `hammer` |
| Serah terima | `handshake` | Urutkan | `arrow-up-down` |
| Antrean | `inbox` | Grafik | `chart-column` |
| Lampiran | `paperclip` | Email | `mail` |
| Terkunci | `lock` | Kembali | `arrow-left` |
| Lanjut | `arrow-right`, `chevron-right` | Buka | `chevron-down` |
| Tutup | `x` | Tambah | `plus` |

Logo: dua mata rantai bersilang (kiri `--primary`, kanan `--text`) diikuti kata "rantai" Plus Jakarta Sans 800, -0,02em. Buat sebagai komponen `Logo` berbasis SVG inline dengan `stroke` dari `currentColor` dan var token, bukan berkas gambar.

---

## 5. Komponen

Konvensi:

- Semua komponen di `src/lib/components/`, satu berkas per komponen, nama PascalCase.
- Props memakai `$props()` Svelte 5 dengan tipe di bawah. `$bindable()` ditandai **bindable**. Konten anak memakai `Snippet`.
- "Status" berarti status visual yang wajib didesain dan diuji: bawaan, hover, fokus, aktif, nonaktif, memuat, error, bila relevan.
- Semua teks lewat i18n (bagian 8). Contoh label di bawah memakai kunci i18n.
- Fokus selalu `:focus-visible` dengan `--ring`. Jangan menghapus outline tanpa pengganti.

### 5.0 Tipe bersama

```ts
// src/lib/types/ui.ts
import type { Component, Snippet } from 'svelte';

export type Theme = 'light' | 'dark' | 'system';
export type Locale = 'id' | 'en';
export type Certainty = 'resmi' | 'proksi' | 'estimasi' | 'simulasi';
export type SectorKey = 'edu' | 'air' | 'san' | 'jalan' | 'kes' | 'ruang' | 'lain';
export type GapClass = 'sangat-kurang' | 'kurang' | 'hampir-cukup' | 'cukup';
export type Level = 'smp' | 'sma';
export type LinkStatus = 'ada' | 'putus' | 'menunggu' | 'estimasi';
export type ChainStage = 'rencana' | 'pemilihan' | 'pemenang' | 'pelaksanaan' | 'lokasi' | 'serah-terima' | 'pemeliharaan' | 'bukti';
export type ReportStatus = 'dikirim' | 'menunggu-saksi' | 'terverifikasi' | 'diteruskan' | 'ditanggapi' | 'selesai' | 'ditolak';
export type SignalKey = 'beda-paket-sejenis' | 'melewati-jadwal' | 'peserta-sedikit';
export type LoadState = 'idle' | 'loading' | 'ready' | 'empty' | 'error';
export type Icon = Component<{ size?: number; strokeWidth?: number }>;

export interface Source {
  name: string;          // "Dapodik"
  owner: string;         // "Kemendikdasmen"
  url?: string;
  retrievedAt: Date;
}

export interface Provenance {
  level: Certainty;
  source: Source;
  formula?: string;      // wajib bila level = 'proksi'
  inputs?: string[];
  coverage?: string;
  limitations?: string;
  methodHref?: string;
}
```

### 5.1 Kerangka aplikasi

#### `AppShell`
- **Props:** `nav: 'beranda' | 'peta' | 'tanya' | 'lapor' | 'akun' | null`, `mobileHeader?: Snippet`, `children: Snippet`, `demo?: boolean`.
- **Perilaku:** menyusun `DemoRibbon` (bila `demo`), `TopBar` (≥ 1024), `MobileHeader` (< 1024), `<main id="konten">`, `BottomNav` (< 1024), wadah toast.
- **Aksesibilitas:** tautan "Langsung ke konten" sebagai elemen pertama yang bisa difokus. Satu `<main>` per halaman. Urutan fokus: tautan lompat, pita demo, header, konten, bilah bawah.

#### `TopBar` (≥ 1024)
- **Props:** `active: 'peta' | 'tanya' | 'tentang' | 'pantau' | null`, `searchText?: string`, `user?: { initials: string } | null`.
- **Isi berurutan:** `Logo`, tombol pencarian wilayah (membuka `RegionSearch` di popover, pintasan `/`), tautan Peta, Tanya, Tentang Data, `LangSwitch`, `ThemeToggle`, tombol Masuk atau avatar.
- **Ukuran:** tinggi 72, margin 32 (1024) atau 72 (≥ 1360), `--surface` dengan garis bawah `--line`, lengket di atas (`z-sticky`).
- **Aksesibilitas:** `<header>` + `<nav aria-label={t('nav.main')}>`. Tautan aktif `aria-current="page"`.

#### `MobileHeader` (< 1024)
- **Props:** `title?: string`, `backHref?: string`, `actions?: Snippet`.
- **Perilaku:** tanpa `backHref` menampilkan `Logo`. Dengan `backHref` menampilkan tombol kembali 44 × 44 dan judul satu baris (terpotong elipsis). Tinggi 60 (360) atau 64 (768).
- **Aksesibilitas:** tombol kembali berlabel `t('action.back')`.

#### `BottomNav` (< 1024)
- **Props:** `active: 'beranda' | 'peta' | 'tanya' | 'lapor' | 'akun' | null`.
- **Isi:** lima tujuan sama lebar: Beranda, Peta, Tanya, Lapor, Akun. Ikon 20 di atas label 12px. Item aktif: kapsul `--secondary` 56 × 30 di belakang ikon, teks 700.
- **Ukuran:** tinggi 64 + `env(safe-area-inset-bottom)`. Disembunyikan di alur Lapor Bukti dan Masuk.
- **Aksesibilitas:** `<nav aria-label={t('nav.main')}>`, item aktif `aria-current="page"`.

#### `Breadcrumb`
- **Props:** `items: { label: string; href?: string }[]`.
- **Aksesibilitas:** `<nav aria-label={t('nav.breadcrumb')}>` dengan `<ol>`. Item terakhir `aria-current="page"`, ikon pemisah `aria-hidden`.

#### `Footer`
- **Props:** `sources: { name: string; owner: string; role: string; retrievedAt: Date }[]`.
- **Isi:** logo, pernyataan "Proyek warga independen, bukan situs resmi pemerintah", jejak sumber data (1 kolom ponsel, 2 kolom tablet, 3 kolom desktop), tautan Tentang Data, Metodologi, Lisensi, Privasi.

#### `ThemeToggle`, `LangSwitch`
- `ThemeToggle` props: `value: Theme` **bindable**. Satu `IconButton` yang berputar antara terang dan gelap; pilihan "ikuti sistem" ada di halaman Akun. Ikon menunjukkan tema tujuan (bulan saat terang).
- `LangSwitch` props: `value: Locale` **bindable**. Segmen dua tombol ID / EN, `role="group"`, tombol aktif `aria-pressed="true"`. Mengganti `lang` di `<html>`.

#### `Logo`
- **Props:** `size?: 'sm' | 'md'`, `withWordmark?: boolean = true`, `href?: string = '/'`.
- **Aksesibilitas:** sebagai tautan berlabel "Rantai, ke beranda". SVG `aria-hidden`.

### 5.2 Aksi dan input

#### `Button`
- **Props:** `variant: 'primary' | 'secondary' | 'outline' | 'text'`, `size?: 'sm' | 'md' | 'lg'` (36 / 44 / 52, bawaan `md`), `icon?: Icon`, `iconEnd?: Icon`, `href?: string`, `type?: 'button' | 'submit'`, `loading?: boolean`, `loadingLabel?: string`, `error?: boolean`, `disabled?: boolean`, `disabledReason?: string`, `fullWidth?: boolean`, `onclick?: (e: MouseEvent) => void`, `children: Snippet`.
- **Varian:** utama (`--primary`), sekunder (`--secondary`), garis (transparan + garis `--line-strong` 1,5px), teks (warna `--on-secondary`). Tidak ada tombol merah.
- **Status:** hover (utama `--primary-hover`, sekunder `--secondary-hover`, garis isi `--sunken` dan garis `--text`), fokus `--ring`, aktif (utama `--primary-press`, geser 1px), nonaktif (`--dis-bg` / `--dis-text`), memuat (ikon `loader-circle` berputar + `loadingLabel`, lebar dipertahankan), error (`--danger-bg`, garis `--danger`, label "Coba lagi").
- **Responsif:** di ponsel tombol aksi utama memakai `fullWidth`. Label boleh dua baris; teks tidak pernah dipotong.
- **Aksesibilitas:** `<button>` atau `<a>` bila ada `href`. Memuat: `aria-busy="true"`. `disabledReason` dirender sebagai teks bantu dan ditautkan lewat `aria-describedby`. Ukuran `sm` tetap punya area sentuh 44 lewat pseudo-elemen tak terlihat.

#### `IconButton`
- **Props:** `icon: Icon`, `label: string` (wajib), `variant?: 'ghost' | 'outline' | 'sunken' | 'glass' | 'primary'`, `tooltip?: boolean = true`, `href?: string`, `onclick?`.
- **Ukuran:** 44 × 44, ikon 20.
- **Aksesibilitas:** `aria-label={label}`, tooltip muncul saat hover 400ms dan saat fokus.

#### `TextLink`
- **Props:** `href: string`, `variant?: 'inline' | 'more'`, `external?: boolean`, `children: Snippet`.
- **Varian:** `inline` bergaris bawah 1px (hover 2px + isi `--secondary`), `more` tebal tanpa garis bawah dengan ikon panah, tinggi minimal 44.
- **Aksesibilitas:** `external` menambah ikon `external-link`, `target="_blank" rel="noopener"`, dan teks tersembunyi `t('a11y.opensNewTab')`.

#### `TextField`
- **Props:** `id: string`, `label: string`, `value: string` **bindable**, `type?: 'text' | 'email' | 'search' | 'tel'`, `placeholder?`, `help?: string`, `error?: string | null`, `validating?: boolean`, `disabled?: boolean`, `autocomplete?: string`, `inputmode?: string`.
- **Status:** bawaan (garis 1,5px `--line-strong`), hover (garis `--text`), fokus (garis 2px `--focus` + cincin lembut), nonaktif, memeriksa (ikon putar di kanan), error (garis 2px `--danger`, pesan dengan ikon).
- **Aksesibilitas:** `<label for>`, `aria-invalid` saat error, `aria-describedby` ke teks bantu dan pesan error. Tinggi 48, teks 16px (mencegah zoom otomatis di iOS).

#### `TextArea`
- **Props:** sama dengan `TextField` + `rows?: number`, `maxlength?: number`.
- **Perilaku:** penghitung "84 / 280" di kanan bawah, `aria-live="polite"` hanya saat sisa ≤ 20 karakter.

#### `OtpInput`
- **Props:** `value: string` **bindable**, `length?: number = 6`, `error?: 'invalid' | 'expired' | null`, `locked?: boolean`, `oncomplete?: (code: string) => void`.
- **Perilaku:** secara teknis satu `<input>` dengan `inputmode="numeric"`, `autocomplete="one-time-code"`, `maxlength={length}`; tampilannya enam kotak 56px. Tempel kode utuh harus berfungsi. `oncomplete` dipanggil saat digit ke-6 terisi.
- **Status:** fokus pada kotak aktif, error (semua kotak bergaris `--danger`), terkunci (semua kotak nonaktif).
- **Aksesibilitas:** satu label "Kode masuk". Pesan error `role="alert"`.

#### `RegionSearch`
- **Props:** `value: string` **bindable**, `suggestions: RegionSuggestion[]`, `state: LoadState`, `size?: 'md' | 'lg'`, `onselect: (r: RegionSuggestion) => void`, `onuselocation?: () => void`.
- `RegionSuggestion = { id: string; name: string; level: 'kelurahan' | 'kecamatan' | 'kota' | 'jalan'; parent: string; match: [number, number] }`.
- **Perilaku:** saran muncul setelah 2 karakter dan jeda ketik 200ms, dikelompokkan per tingkat (Kelurahan, Kecamatan, Jalan dan tempat), huruf yang cocok ditebalkan, baris terakhir "Pakai lokasiku saja". Kosong: "Tidak ada wilayah bernama ..." + saran ejaan. Error: baris error dengan tombol Coba lagi.
- **Ukuran:** tinggi 56 (60 di beranda desktop), saran maksimal 8 baris, tiap baris minimal 52.
- **Aksesibilitas:** pola combobox ARIA 1.2: input `role="combobox"`, `aria-expanded`, `aria-controls`, `aria-activedescendant`; daftar `role="listbox"`, opsi `role="option"` + `aria-selected`. Panah atas bawah memindah sorotan, Enter memilih, Esc menutup. Pintasan `/` memfokuskan kotak di desktop. Jumlah hasil diumumkan lewat `aria-live="polite"`.

#### `UseLocationButton`
- **Props:** `state: 'idle' | 'locating' | 'found' | 'coarse' | 'denied'`, `accuracyM?: number`, `placeName?: string`, `coarseThresholdM?: number = 100`, `onclick: () => void`.
- **Status:** bawaan (garis), mencari (putar, "Mencari lokasi...", batas 10 detik), ditemukan (`--success-bg`, nama kelurahan + akurasi), kasar (amber, "Akurasi 1,2 km", ajakan memilih manual), ditolak (nonaktif, teks bantu cara mengizinkan).
- **Aksesibilitas:** perubahan status diumumkan `aria-live="polite"`. Lokasi tidak pernah disimpan tanpa tindakan pengguna.

#### `SegmentedControl<T extends string>`
- **Props:** `options: { value: T; label: string; disabled?: boolean; disabledReason?: string }[]`, `value: T` **bindable**, `label: string`.
- **Tampilan:** wadah `--sunken` padding 4, segmen tinggi 40 (area wadah 48), segmen aktif `--raised` + `--sh-1`.
- **Aksesibilitas:** `role="radiogroup"` + `aria-label`, segmen `role="radio"` + `aria-checked`, roving tabindex, panah kiri kanan berpindah.

#### `LensTabs`
- **Props:** `active: 'layanan' | 'uang' | 'tanggung-jawab'`, `baseHref: string`, `counts?: Partial<Record<'layanan' | 'uang' | 'tanggung-jawab', number>>`.
- **Perilaku:** setiap tab adalah rute sendiri (dapat dibagikan dan dibuka langsung), bukan panel di halaman yang sama.
- **Responsif:** < 768 tiga kolom sama lebar, ikon di atas label 13px, tinggi 60. ≥ 768 ikon di samping label 15px, tinggi 56. ≥ 1024 lebar mengikuti isi.
- **Tampilan aktif:** garis bawah 3px `--primary`, ikon `--primary`, teks 700.
- **Aksesibilitas:** `<nav aria-label={t('lens.nav')}>` dengan tautan; tab aktif `aria-current="page"`. Karena berupa navigasi rute, jangan pakai `role="tablist"`.

#### `FilterChip` dan `ChipGroup`
- `FilterChip` props: `label: string`, `pressed: boolean` **bindable**, `count?: number`, `icon?: Icon`, `sector?: SectorKey`, `disabled?: boolean`.
- `ChipGroup` props: `label: string`, `children: Snippet`.
- **Tampilan:** tinggi 40, kapsul, garis 1,5px `--line-strong`. Terpilih: ikon centang + `--secondary` + garis `--on-secondary`. Hover: garis `--text` + isi `--sunken`. Nonaktif saat hitungan 0.
- **Responsif:** < 1024 satu baris yang bisa digulir horizontal dengan bayangan tepi; ≥ 1024 membungkus ke baris berikut.
- **Aksesibilitas:** `<button aria-pressed>`. `ChipGroup` adalah `role="group"` + `aria-label`. Hitungan ikut dibacakan ("Pendidikan, 12").

#### `Switch`
- **Props:** `checked: boolean` **bindable**, `label: string`, `description?: string`, `disabled?: boolean`, `disabledReason?: string`.
- **Tampilan:** trek 52 × 32, kenop 24 dengan ikon centang saat aktif, baris minimal 56.
- **Aksesibilitas:** `<button role="switch" aria-checked>` dengan label terlihat yang ditautkan.

#### `Checkbox`, `RadioCard`
- `Checkbox` props: `checked` **bindable**, `label`, `description?`, `disabled?`. Kotak 22, baris minimal 44.
- `RadioCard` props: `name: string`, `value: string`, `group: string` **bindable**, `title: string`, `description?: string`, `leading?: Snippet`. Kartu minimal 56, terpilih `--secondary` + garis 2px `--on-secondary`.
- **Aksesibilitas:** elemen asli `<input type="checkbox|radio">` disembunyikan secara visual, bukan `display:none`. Grup radio dibungkus `<fieldset>` + `<legend>`.

#### `Dialog`
- **Props:** `open: boolean` **bindable**, `title: string`, `description?: string`, `dismissible?: boolean = true`, `children?: Snippet`, `actions: Snippet`.
- **Tampilan:** lebar maksimal 480 (ponsel: lebar layar dikurangi 32), `--raised`, radius 16, `--sh-3`, scrim `--scrim`.
- **Aksesibilitas:** elemen `<dialog>` dengan `showModal()`, `aria-labelledby` ke judul. Fokus awal ke tombol paling aman (bukan tombol berisiko). Esc menutup bila `dismissible`. Fokus kembali ke pemicu saat ditutup.

#### `BottomSheet` (< 1024)
- **Props:** `open: boolean` **bindable**, `detent: 'peek' | 'half' | 'full'` **bindable**, `detents?: { peek: number; half: number; full: number } = { peek: 0.3, half: 0.55, full: 0.92 }`, `title: string`, `modal?: boolean = false`, `children: Snippet`.
- **Perilaku:** pegangan 40 × 5 di atas. Seret untuk berpindah titik tahan, lempar cepat untuk melompat. Non-modal di atas peta (peta tetap bisa digeser); modal untuk panel sumber dan filter.
- **Aksesibilitas:** `role="dialog"` + `aria-labelledby`. Pegangan adalah `<button>` berlabel `t('sheet.resize')`; panah atas bawah mengganti titik tahan. Modal: perangkap fokus + Esc menutup. ≥ 1024 konten yang sama dirender di `SidePanel`.

#### `SidePanel` (≥ 1024)
- **Props:** `open` **bindable**, `title`, `width?: number = 400`, `overlay?: boolean = false`, `children`.
- **Perilaku:** menempel kanan. Non-overlay menggeser konten (peta), overlay di atas konten dengan scrim tipis (panel Tanya Rantai di Wilayahku, lebar 440).

#### `Toast`
- **API:** `toast.show({ kind: 'success' | 'info' | 'error', message: string, action?: { label: string; onclick: () => void }, duration?: number })` dari store `src/lib/stores/toast.ts`.
- **Perilaku:** maksimal satu toast terlihat. Hilang setelah 5 detik kecuali `error`. Posisi: bawah tengah di atas bilah bawah (ponsel), kanan bawah (desktop).
- **Tampilan:** `--raised`, `--sh-2`, ikon dalam lingkaran berwarna di kiri. Tidak memakai garis aksen tepi kiri.
- **Aksesibilitas:** `role="status"` untuk sukses dan info, `role="alert"` untuk error. Tombol tutup 44 × 44.

#### `Tooltip`
- **Props:** `text: string`, `delay?: number = 400`, `children: Snippet` (pemicu).
- **Aksesibilitas:** `role="tooltip"`, pemicu `aria-describedby`. Muncul saat hover dan fokus, tekan lama di layar sentuh. Hanya label pendek; penjelasan panjang memakai `SourcePanel`.

### 5.3 Komponen data

#### `CertaintyBadge`
- **Props:** `level: Certainty`, `size?: 'sm' | 'md'` (22 / 28), `label?: string`.
- **Tampilan:** lihat tabel 1.4. Lebar mengikuti teks, tidak pernah dipotong.
- **Aksesibilitas:** teks label selalu terlihat, ikon `aria-hidden`. Bila bukan tombol, berupa `<span>`.

#### `ProvenanceBadge`
- **Props:** `provenance: Provenance`, `compact?: boolean`, `onopen?: () => void`.
- **Tampilan:** satu tombol gabungan tinggi 26: kiri `CertaintyBadge` tanpa sudut, kanan nama sumber + tanggal pendek + ikon `chevron-down` (atau `external-link` bila langsung ke sumber). `compact` (< 768 atau ruang sempit) hanya nama sumber.
- **Aksesibilitas:** `<button aria-haspopup="dialog" aria-expanded>`. Nama aksesibel lengkap: "Sumber: Dapodik, resmi, diambil 14 Oktober 2026". Area sentuh 44 lewat padding tak terlihat.

#### `SourcePanel`
- **Props:** `provenance: Provenance`, `value: string` (angka yang dijelaskan, sudah diformat), `open` **bindable**.
- **Isi berurutan:** lencana kepastian, angka, sumber dan pemilik, waktu diambil (tanggal + jam WIB), cakupan, rumus dalam blok mono (wajib untuk proksi), masukan, batasan, tautan "Buka data asli" atau "Baca metodologi".
- **Responsif:** ≥ 1024 popover tertambat ke lencana (lebar 360, `z-popover`); < 1024 `BottomSheet` modal setengah layar.
- **Aksesibilitas:** `role="dialog"`, fokus ke judul, Esc menutup, fokus kembali ke lencana.

#### `KeyNumberCard`
- **Props:** `label: string`, `value: string`, `unit?: string`, `provenance?: Provenance`, `comparison?: string`, `tone?: 'neutral' | 'attn'`, `simulated?: boolean`, `size?: 'md' | 'lg'`, `href?: string`.
- **Tampilan:** label `caption` redup, angka `n-lg` (atau `n-md`), satuan di sampingnya, `ProvenanceBadge` atau `CertaintyBadge`, garis pemisah, lalu kalimat pembanding netral ("Rata-rata DKI 1,4. Tahun lalu 2,0."). `tone='attn'` memberi latar `--attn-bg` untuk hitungan sinyal. `simulated` menambah `SimRibbon` di atas kartu.
- **Aturan:** tidak ada panah naik turun berwarna. Perbandingan selalu kalimat.
- **Aksesibilitas:** angka dan satuan dalam satu elemen teks agar dibacakan utuh ("2,1 lulusan per kursi").

#### `GapBar`
- **Props:** `title: string`, `graduates: number`, `seats: { negeri: number; swasta: number; madrasah: number }`, `ratioText: string`, `showLegend?: boolean = true`.
- **Tampilan:** baris "Calon lulusan" (batang `--text` 82%) dan baris "Kursi kelas awal" bertumpuk: negeri isi penuh `--primary`, swasta `--primary` 45%, madrasah garis vertikal + garis tepi, selisih kotak berarsir amber. Celah 2px antarsegmen. Nilai di kanan, 52px, tabular.
- **Aksesibilitas:** `role="img"` dengan `aria-label` ringkasan ("2.701 calon lulusan, 1.284 kursi negeri, 640 swasta, 301 madrasah, selisih 476"). Selalu ada tombol `TableToggle`. Empat segmen terbaca tanpa warna.

#### `ChainStepper` (Rantai Rupiah)
- **Props:** `links: ChainLinkData[]`, `orientation?: 'auto' | 'horizontal' | 'vertical' = 'auto'`, `animate?: boolean = true`, `summary?: string`.
- `ChainLinkData = { stage: ChainStage; status: LinkStatus; title: string; meta: string; detail?: string; href?: string }`.
- **Tampilan mata rantai:** kapsul bergaris 3px berisi ikon tahap, lencana status kecil di sudut kanan atas.
  - `ada`: garis utuh `--primary`, isi `--secondary`, lencana centang.
  - `putus`: garis putus `--attn-fill`, isi `--attn-bg`, lencana `unlink-2`.
  - `menunggu`: garis titik-titik `--line-strong`, tanpa isi, lencana jam.
  - `estimasi`: garis putus `--cert-estimasi`, lencana `approximately-equal`.
- **Sambungan antar mata:** putus bila salah satu sisi `putus` (dua ujung terpisah celah 10px); titik-titik bila salah satu sisi `menunggu`; selain itu kapsul utuh yang bertumpang dengan kedua mata.
- **Responsif:** `auto` = vertikal di < 1024 (mata 34 × 46 di kolom 44px, judul dan status di kanan, garis vertikal sebagai sambungan), horizontal di ≥ 1024 (mata 68 × 40, judul dan meta di bawah, rata tengah).
- **Aksesibilitas:** `<ol aria-label={summary}>`, contoh ringkasan "5 dari 8 mata rantai tersambung, 1 putus". Setiap `<li>` memuat teks status ("Putus"), bukan hanya ikon. Mata yang punya `href` adalah tautan.

#### `MiniChain`
- **Props:** `statuses: LinkStatus[]`, `label?: string`.
- **Tampilan:** mata 20 × 13 bergaris 3px, saling bertumpang 3px; mata putus diberi jarak 7px di kirinya.
- **Aksesibilitas:** `role="img"` + `aria-label` ("5 dari 8 tersambung, 1 putus").

#### `SignalChip` dan `SignalPanel`
- `SignalChip` props: `signal: SignalKey`, `size?: 'sm' | 'md'`, `expanded?: boolean` **bindable**, `controls?: string`.
- `SignalPanel` props: `signal: SignalKey`, `evidence: Snippet`, `actions?: Snippet`, `expanded` **bindable**.
- **Tampilan:** kapsul `--attn-bg`, garis 1,5px `--attn-fill`, teks `--attn` 700, ikon `circle-help`. Panel terbuka: judul netral, data pembanding, kalimat "Perbedaan bisa wajar...", tautan "Cara sinyal dihitung" dan pintu ke hak jawab.
- **Aturan:** judul hanya dari katalog netral (bagian 9). Di dalam kartu yang seluruhnya tautan, chip dirender sebagai `<span>` non-interaktif.
- **Aksesibilitas:** chip interaktif adalah `<button aria-expanded aria-controls>`.

#### `MaintenanceCounter`
- **Props:** `start?: Date`, `end?: Date`, `status: 'aktif' | 'segera' | 'berakhir' | 'belum'`, `totalDays?: number = 180`.
- **Tampilan:** ikon `hourglass`, angka hari tersisa `n-lg` + "hari lagi", batang kemajuan 8px, tanggal akhir. `segera` (≤ 30 hari) memakai amber dan ajakan "Lapor kerusakan sebelum masa ini habis". `belum` memakai batang titik-titik dan teks "Menunggu serah terima".
- **Aksesibilitas:** batang `role="progressbar"` dengan `aria-valuenow` (hari berjalan), `aria-valuemax`, dan `aria-valuetext` ("47 hari lagi, sampai 9 Desember 2026").

#### `StripPlot`
- **Props:** `values: number[]`, `current: number`, `domain?: [number, number]`, `format: (n: number) => string`, `title: string`.
- **Tampilan:** titik paket sejenis 11px `--muted` 70%, pita kuartil (Q1 sampai Q3) `--secondary`, garis median 2,5px `--on-secondary` berlabel, paket ini lingkaran 24px `--text` bercincin. Sumbu bawah 4 sampai 6 label tabular. Kalimat wajib di bawah grafik, Newsreader italic: "Perbandingan ini menunjukkan posisi nilai, bukan kewajaran."
- **Aksesibilitas:** `role="img"` + `aria-label` ("Paket ini Rp2,41 miliar. Median 14 paket sejenis Rp1,62 miliar."). `TableToggle` menampilkan daftar 14 paket.

#### `TableToggle`
- **Props:** `label?: string`, `pressed: boolean` **bindable**, `controls: string`.
- **Perilaku:** mengganti grafik atau peta dengan `DataTable` setara di tempat yang sama. Pilihan diingat per sesi.
- **Aksesibilitas:** `aria-pressed`, `aria-controls`.

#### `DataTable<T>`
- **Props:** `caption: string`, `columns: { key: keyof T; label: string; align?: 'start' | 'end'; sortable?: boolean; format?: (v: unknown, row: T) => string }[]`, `rows: T[]`, `sort?: { key: keyof T; dir: 'asc' | 'desc' }` **bindable**, `stickyFirstColumn?: boolean = true`, `highlight?: (row: T) => boolean`.
- **Tampilan:** kepala `--sunken` 12px kapital, baris 48, angka rata kanan dan tabular.
- **Responsif:** < 1024 tabel digulir horizontal dalam wadahnya, kolom pertama lengket.
- **Aksesibilitas:** `<table>` asli, `<caption>` (boleh tersembunyi visual), `<th scope>`, tombol urut di dalam `<th>` dengan `aria-sort` pada `<th>`.

#### `SectorIcon`
- **Props:** `sector: SectorKey`, `size?: 20 | 24 | 44`.
- **Tampilan:** kotak radius 6 berwarna `--sector-*` dengan ikon `--sector-*-on`. Selalu didampingi label di sekitarnya atau `aria-label`.

#### `PackageCard`
- **Props:** `pkg: { id: string; title: string; sector: SectorKey; year: number; rupCode: string; agency: string; value: number; valueLevel: Certainty; vendor: string; chain: LinkStatus[]; signals: SignalKey[]; maintenance?: { status: 'aktif' | 'segera' | 'berakhir' | 'belum'; daysLeft?: number }; simulated?: boolean }`.
- **Isi berurutan:** sektor + tahun dan kode RUP (mono), judul 17/23 700, instansi, nilai `n-md` + lencana, pelaksana, `MiniChain` + ringkasan, baris tag (sinyal dan status pemeliharaan).
- **Tampilan:** seluruh kartu adalah satu tautan ke Rantai Rupiah, `--sh-1`. `simulated`: wadah bergaris `--cert-sim` + `SimRibbon`.
- **Aksesibilitas:** satu `<a>` pembungkus; tidak ada elemen interaktif bersarang. Judul jadi nama tautan.

#### `SchoolCard`
- **Props:** `school: { id: string; name: string; status: 'negeri' | 'swasta' | 'madrasah'; kelurahan: string; seats: number; rombel: number; applicants?: number; distanceKm?: number }`, `variant?: 'card' | 'row'`.
- **Isi:** ikon sekolah, nama, status dan kelurahan, jarak (bila ada), tiga angka: kursi kelas awal, rombel, pendaftar tahun lalu. `row` dipakai di Cek Sekolahku dengan nomor urut dan keketatan.

#### `RepresentativeCard`
- **Props:** `rep: { id: string; initials: string; name: string; faction: string; commission: string; dapil: string; oversees: boolean }`, `variant?: 'compact' | 'full'`.
- **Isi:** avatar inisial bulat (bukan foto), nama, fraksi dan komisi, penanda kapsul "Duduk di komisi yang mengawasi urusan ini" bila `oversees`. `full` menambah daftar kecamatan dapil, tautan LHKPN dan situs DPRD.
- **Aturan:** tanpa skor, tanpa peringkat, tanpa warna partai.
- **Aksesibilitas:** avatar `aria-hidden`.

#### `AgencyCard`
- **Props:** `agency: { name: string; role: string; description: string; legalBasis: { text: string; href?: string }[]; website?: string }`.
- **Isi:** ikon instansi, peran singkat (caption 700), nama, uraian kewenangan, blok cekung "Dasar hukum" dengan tautan naskah, tombol situs instansi, lencana sumber.

#### `ReportStatusChip`
- **Props:** `status: ReportStatus`, `witnesses?: { count: number; needed: number }`, `size?: 'sm' | 'md'`.
- **Tampilan:** lihat tabel 1.6. Menunggu saksi menampilkan hitungan "2/3".

#### `ReportTimeline`
- **Props:** `events: { status: ReportStatus | 'cek-otomatis'; at?: Date; note?: string; done: boolean }[]`.
- **Perilaku:** memakai `ChainStepper` dengan `orientation="vertical"` (mata rantai menjadi simpul linimasa). Tahap belum terjadi berstatus `menunggu`.

#### `WitnessPanel`
- **Props:** `count: number`, `needed?: number = 3`, `canWitness: boolean`, `reason?: string`, `onwitness: () => void`.
- **Isi:** tiga kapsul kemajuan, tombol utama besar "Saya juga melihat ini", syarat (masuk, dalam radius 200 m, satu foto pendukung).
- **Aksesibilitas:** tombol nonaktif menjelaskan alasannya lewat `aria-describedby`.

#### `RightOfReply`
- **Props:** `reply: { agency: string; verified: boolean; publishedAt: Date; body: string; attachments: { name: string; href: string; sizeKb?: number }[]; daysAfterForward?: number }`.
- **Tampilan:** pita atas `--info-bg` dengan ikon `scale` dan tanggal, nama instansi + "Akun instansi terverifikasi", kutipan Newsreader 19/28, lampiran, catatan "Ditampilkan apa adanya tanpa disunting Rantai".
- **Aksesibilitas:** `<article aria-labelledby>`, kutipan dalam `<blockquote>`.

#### `ChatMessage`, `Citation`, `ChatComposer`
- `ChatMessage` props: `role: 'user' | 'assistant'`, `state?: 'streaming' | 'done' | 'no-data' | 'error'`, `children: Snippet`, `citations?: { n: number; level: Certainty; label: string; provenance: Provenance }[]`, `followUps?: { label: string; href?: string }[]`.
- `Citation` props: `n: number`, `level: Certainty`, `label: string`, `onopen: () => void`. Tombol kecil 22 × 20 bernomor, warna mengikuti label kepastian, area sentuh 44.
- `ChatComposer` props: `value` **bindable**, `context?: { label: string; onremove: () => void }`, `suggestions?: string[]`, `disabled?: boolean`, `onsubmit: (q: string) => void`.
- **Tampilan:** pengguna rata kanan `--secondary`, sudut kanan bawah 4px. Rantai rata kiri `--raised` bergaris `--line`, avatar logo kecil. `no-data`: latar `--sunken`, garis putus, ikon `unlink-2`, teks "Rantai belum punya data ini" + apa yang tersedia. Streaming: label "Menyusun jawaban dari 3 sumber..." + kursor berkedip.
- **Aksesibilitas:** daftar pesan `role="log"` + `aria-live="polite"`. Jawaban streaming diumumkan sekali saat selesai, bukan per token. `Citation` berlabel "Rujukan 1: Dapodik".

#### `PhotoGrid`
- **Props:** `photos: { id: string; src: string; takenAt: Date; accuracyM: number; blurred?: number; hidden?: boolean }[]`, `max?: number = 4`, `editable?: boolean`, `onadd?`, `onremove?`.
- **Tampilan:** kisi 2 kolom persegi; slot kosong bergaris putus dengan ikon kamera; keterangan jam dan akurasi di pojok.
- **Aksesibilitas:** tombol hapus per foto berlabel "Hapus foto 1". Foto bukti memakai `alt` dari catatan pelapor atau "Foto bukti 1, 12 Oktober 08.14".

### 5.4 Peta

#### `MapView`
- **Props:** `layers: { gap: boolean; projects: boolean; reports: boolean; schools: boolean }` **bindable**, `level: Level` **bindable**, `selected?: { kind: 'area' | 'project' | 'report' | 'school'; id: string } | null` **bindable**, `bounds?: [number, number, number, number]`, `onselect?`.
- **Teknologi:** MapLibre GL JS dengan ubin vektor. Gaya peta dasar membaca token `--map*` saat inisialisasi dan saat tema berganti (`setPaintProperty`), tanpa label POI komersial.
- **Aksesibilitas:** wadah `role="region"` + `aria-label`. Penanda adalah tombol HTML di atas kanvas agar bisa difokus; urutan Tab mengikuti daftar terurut. Selalu ada `TableToggle` yang menampilkan data setara.

#### `MapMarker`
- **Props:** `kind: 'project' | 'report' | 'school' | 'cluster' | 'me'`, `sector?: SectorKey`, `signal?: boolean`, `selected?: boolean`, `simulated?: boolean`, `reportState?: 'verified' | 'waiting'`, `count?: number`, `label: string`.
- **Bentuk membedakan jenis:** lingkaran 32 (proyek, warna dan ikon sektor, cincin `--surface` 2,5px), belah ketupat 24 (laporan: hijau terverifikasi, garis putus info menunggu), persegi 26 bergaris tinta (sekolah), lingkaran 44 putih dengan angka (klaster), titik biru berhalo (lokasiku). Terpilih: 44 + cincin `--text`. Sinyal: lencana "?" amber di sudut. Simulasi: isi arsir.
- **Aksesibilitas:** `<button aria-label={label}>`, area sentuh 44 untuk semua ukuran.

#### `AreaHatch`
- **Props:** `geometry: GeoJSON.Polygon`, `sector: SectorKey`, `label?: string`.
- **Tampilan:** arsir diagonal warna sektor 55% + garis tepi 2px. Dipakai saat lokasi proyek hanya diketahui di tingkat kecamatan. Jangan pernah membuat titik palsu.

#### `LayerControl`
- **Props:** `layers` **bindable** (sama dengan `MapView`), `counts: Record<string, number>`, `level: Level` **bindable**.
- **Responsif:** ≥ 1024 panel kaca kiri atas lebar 264. < 1024 `IconButton` "Lapisan" kaca 48 × 48 yang membuka `BottomSheet` modal.
- **Aksesibilitas:** `<fieldset>` + `<legend>`, setiap lapisan `Checkbox` asli.

#### `MapLegend`
- **Props:** `scale: 'gap'`, `level: Level`, `compact?: boolean`, `provenance: Provenance`.
- **Responsif:** ≥ 1024 panel kaca kiri bawah lebar 300 (empat kotak kelas + nama + rentang). < 1024 kapsul kaca ringkas di atas lembar bawah (empat kotak + "kurang ke cukup").

#### `MapDetailPanel`
- **Props:** `area?: { name: string; ratio: number; gapClass: GapClass; graduates: number; seats: number; counts: { projects: number; reports: number; schools: number } }`, `nearby?: PackageSummary[]`.
- **Responsif:** ≥ 1024 `SidePanel` 400 tetap di kanan (peta menyempit). < 1024 `BottomSheet` non-modal titik tahan `peek`.

### 5.5 Pita, status, alur

#### `DemoRibbon`
- **Props:** `onexit: () => void`.
- **Tampilan:** pita arsir magenta tinggi minimal 44, garis bawah `--cert-sim`, teks dalam kapsul `--cert-sim-bg` agar kontras, tombol "Keluar" padat `--cert-sim`. Lengket di paling atas, `z-ribbon`, tidak bisa ditutup selama mode aktif.
- **Aksesibilitas:** `role="status"`, berada sebelum header dalam urutan fokus.

#### `SimRibbon` dan `SimBox`
- `SimRibbon` props: `note?: string`. Pita arsir 32px "DATA SIMULASI" + catatan opsional.
- `SimBox` props: `note?: string`, `children: Snippet`. Wadah bergaris 1,5px `--cert-sim` dengan `SimRibbon` di atas.
- **Aturan:** setiap angka simulasi wajib berada di dalam `SimBox` atau berlencana `simulasi`. Ekspor PDF dan gambar bagikan menyertakan pita yang sama.

#### `Skeleton`
- **Props:** `shape?: 'text' | 'block' | 'circle' | 'pill'`, `width?: string`, `height?: string`.
- **Aturan:** kerangka meniru bentuk konten akhir. Muncul hanya bila memuat lebih dari 300ms. Wadahnya `aria-busy="true"` dan ada satu teks status yang terbaca ("Memuat ringkasan data terbaru...").

#### `EmptyState`
- **Props:** `illustration?: 'chain' | Icon`, `title: string`, `description: string`, `actions?: Snippet`.
- **Aturan:** jelaskan sebabnya dan beri satu jalan keluar. Ilustrasi bawaan: tiga mata rantai bergaris titik-titik.

#### `ErrorState`
- **Props:** `variant: 'inline' | 'page'`, `title: string`, `description?: string`, `incidentId?: string`, `onretry?: () => void`, `fallback?: Snippet`.
- **Tampilan:** ikon dalam lingkaran `--danger-bg`, tombol "Coba lagi". Bila data setara masih tersedia (misalnya tabel saat peta gagal), tawarkan lewat `fallback`.
- **Aksesibilitas:** `role="alert"`.

#### `OfflineBanner`
- **Props:** `savedAt: Date`.
- **Tampilan:** pita tinta penuh lebar (`--text` dengan teks `--surface`) di bawah header: "Kamu sedang luring. Menampilkan salinan tersimpan 2 jam lalu."
- **Aksesibilitas:** `role="status"`.

#### `StepHeader` (Lapor Bukti)
- **Props:** `step: number`, `total?: number = 6`, `onback: () => void`, `onclose: () => void`.
- **Tampilan:** tombol tutup (langkah 1) atau kembali, judul "Lapor Bukti", teks "Langkah 2 dari 6", dan baris enam kapsul kemajuan: selesai isi `--primary`, aktif garis `--primary` + isi `--secondary`, mendatang titik-titik.
- **Aksesibilitas:** `role="progressbar"` + `aria-valuenow`, `aria-valuemax`, `aria-valuetext` ("Langkah 2 dari 6, pilih proyek").

#### `WorkspaceShell` (Panel Instansi, moderasi)
- **Props:** `section: string`, `nav: { key: string; label: string; icon: Icon; count?: number; href: string }[]`, `org?: { name: string; verified: boolean }`, `user: { initials: string; name: string; role: string }`, `children`.
- **Tata letak:** kolom navigasi 240 `--surface` + konten. Tidak memakai `TopBar` atau `BottomNav` publik.

#### `QueueList`
- **Props:** `items: { id: string; title: string; meta: string; age: string; status: string; tag?: Snippet }[]`, `selectedId` **bindable**.
- **Aksesibilitas:** `role="listbox"` bila memilih memuat detail di sebelahnya; pintasan J dan K berpindah, V verifikasi, R tolak (hanya di panel moderasi, tercantum di header).

---

## 6. Tata letak setiap layar

Ketentuan umum untuk semua layar publik:

- **360:** 1 kolom, margin 16, celah 12. `MobileHeader` di atas, `BottomNav` di bawah.
- **768:** kisi 8 kolom, margin 24, celah 16. Kartu dan daftar menjadi 2 kolom. Navigasi tetap pola ponsel.
- **1024:** kisi 12 kolom, margin 32, celah 20. `TopBar` menggantikan `MobileHeader` dan `BottomNav`. Pola `split` = konten 8 kolom + aside 4 kolom.
- **1440:** sama dengan 1024, konten maksimal 1296 di tengah, margin 72, celah 24 sampai 48 antara konten dan aside.
- Aside desktop selalu ditaruh **setelah** konten utama di DOM, sehingga di ponsel muncul di bawah tanpa perlu diurutkan ulang. Pengecualian dicatat per layar.

Rute SvelteKit yang disarankan:

| Layar | Rute |
|---|---|
| Beranda | `/` |
| Wilayahku | `/w/[kode]` (Layanan), `/w/[kode]/uang`, `/w/[kode]/tanggung-jawab` |
| Peta | `/peta` |
| Rantai Rupiah | `/paket/[id]` |
| Detail kecamatan sekolah | `/w/[kode]/sekolah` |
| Cek Sekolahku | `/cek-sekolah` |
| Profil penyedia | `/penyedia/[id]` |
| Profil anggota DPRD | `/dprd/[id]` |
| Tanya Rantai | `/tanya` (dan panel di `/w/[kode]?tanya=1`) |
| Lapor Bukti | `/lapor` (langkah lewat `?langkah=1..6`, `/lapor/terkirim`) |
| Detail laporan | `/laporan/[id]` |
| Paket Aspirasi | `/aspirasi?w=[kode]&masalah=[kunci]` |
| Pantau | `/pantau` |
| Masuk | `/masuk` |
| Panel Instansi | `/instansi/...` |
| Panel moderasi | `/moderasi/...` |
| Tentang Data | `/tentang-data` |
| 404, error | `+error.svelte` |

### 6.1 Beranda

| Lebar | Kisi | Urutan blok | Disembunyikan atau dipindah |
|---|---|---|---|
| 360 | 1 kolom | Overline, judul display, paragraf pembuka, `RegionSearch` (lebar penuh), `UseLocationButton`, chip contoh wilayah, kartu "Contoh hasil", tiga kartu lensa bertumpuk, kartogram + legenda + tombol, contoh pertanyaan Tanya + cuplikan obrolan, kaki dengan jejak sumber | Pintasan `/` di kotak cari disembunyikan |
| 768 | 1 kolom hero, kartu lensa 2 kolom (kartu ketiga di baris baru) | Sama dengan 360 | Jejak sumber di kaki 2 kolom |
| 1024 | Hero `split` 8/4 rata tengah vertikal | Kiri: teks + cari + lokasiku + chip; kanan: kartu "Contoh hasil" (`--raised`, `--sh-2`). Lensa 3 kolom. Peta `split-l` 4/8 (teks kiri, kartogram kanan). Tanya `split` 8/4 (daftar pertanyaan kiri, cuplikan obrolan kanan). Kaki 3 kolom | Tidak ada |
| 1440 | Sama, margin 72 | Sama | Kartogram maksimal 640 lebar |

Kartogram: kisi 8 × 8, satu kotak per kecamatan (44 kecamatan) disusun mendekati letak geografis, label tiga huruf, kecamatan aktif bergaris 2,5px `--text`. Di 360 kotak sekitar 38px, label 10px; ≥ 768 label 12px.

### 6.2 Wilayahku (kepala bersama tiga tab)

| Lebar | Kepala |
|---|---|
| 360 | `MobileHeader` (kembali, "Kec. Cilincing", cari). Breadcrumb, judul h1, meta, tombol Pantau dan Bagikan (berdampingan, membungkus bila perlu). Tiga kartu ringkas lensa **disembunyikan** (informasinya muncul di tab). `LensTabs` tiga kolom ikon di atas label |
| 768 | Kartu ringkas lensa tampil 2 kolom (kartu ketiga baris baru), kartu tab aktif bergaris 2px `--primary`. Tab ikon di samping label |
| 1024 | Tombol tambahan "Tanya Rantai" (membuka panel geser). Kartu ringkas 3 kolom. Tab selebar isi |
| 1440 | Sama dengan 1024 |

#### 6.2a Tab Layanan

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom | Segmen SMP / SMA-SMK, kalimat ringkas (Newsreader 30/36) + lencana provenance, kartu `GapBar` (SMP lalu SMA/SMK) + tombol tabel + tautan detail sekolah, catatan info penerimaan lintas kecamatan, judul "SMP negeri di Cilincing (12)" + urutan, kartu sekolah bertumpuk, tombol "Lihat semua", lalu isi aside: angka kunci, kartu Cek Sekolahku, cara angka dihitung | Aside muncul setelah daftar sekolah |
| 768 | Kartu sekolah 2 kolom | Sama | |
| 1024 | `split` 8/4 | Kiri: segmen sampai daftar sekolah. Kanan: angka kunci (daftar 5 baris dengan lencana), Cek Sekolahku (`--secondary`), rumus | |
| 1440 | Sama | | |

#### 6.2b Tab Uang

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom, angka kunci 2 × 2 | 4 kartu angka kunci (nilai, jumlah, dalam pemeliharaan, sinyal), batang komposisi sektor, segmen tahun + urutan, chip sektor (gulir horizontal), kartu paket bertumpuk, satu kartu paket simulasi dalam `SimBox`, "Muat 31 paket lagi", lalu aside | Kartu ringkas lensa di kepala disembunyikan |
| 768 | Angka kunci 4 kolom, kartu paket 2 kolom | Sama | |
| 1024 | `split` 8/4 | Kanan: penyedia terbanyak (tanpa skor), kartu "Tentang sinyal" amber, tombol "Lihat paket di peta". Chip membungkus | |
| 1440 | Sama | | |

#### 6.2c Tab Tanggung Jawab

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom | Judul "Masalah apa yang ingin kamu sampaikan?", chip jenis masalah (radio, gulir), kartu "Rantai tanggung jawab" (`ChainStepper` vertikal 4 mata: instansi, komisi, wakil, jalur), kartu instansi berwenang (2), kartu Komisi E, wakil dapil (4 kartu ringkas), jalur aspirasi terdekat, kartu ajakan "Buat Paket Aspirasi" (tombol lebar penuh) | |
| 768 | Sama, kartu wakil tetap 1 kolom | | |
| 1024 | Rantai tanggung jawab horizontal; di bawahnya `split` 8/4 | Kiri: instansi + komisi. Kanan: wakil dapil + jalur aspirasi. Kartu ajakan selebar konten di bawah | |
| 1440 | Sama | | |

### 6.3 Peta

Frame peta selalu setinggi viewport (tidak menggulir halaman).

| Lebar | Susunan | Kontrol |
|---|---|---|
| 360 | Peta penuh di antara header dan `BottomNav` (tanpa `MobileHeader`; pencarian mengambang) | Atas: kotak cari kaca + `IconButton` Lapisan 48. Di bawahnya: segmen SMP/SMA kaca (kiri) dan tombol "Tabel" kaca (kanan). Bawah: legenda ringkas kaca di atas `BottomSheet` detail (peek). Tombol zoom disembunyikan (pakai cubit) |
| 768 | Sama dengan 360, lembar bawah lebih lebar (maks 560, rata tengah) | Sama |
| 1024 | `TopBar` + baris: peta (fleksibel) dan `SidePanel` detail 400 tetap di kanan | Kiri atas: `LayerControl` kaca 264. Kiri bawah: `MapLegend` kaca 300. Kanan atas: "Lihat sebagai tabel". Kanan bawah: zoom + ke lokasiku |
| 1440 | Sama | Sama |

Panel detail berisi kelurahan terpilih: angka rasio + kelas + lencana, `GapBar` ringkas, hitungan proyek, laporan, sekolah, dua proyek terdekat, tombol "Buka Wilayahku" dan "Tabel".

### 6.4 Rantai Rupiah (detail paket)

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom | Kepala: breadcrumb, sektor + kode RUP, judul h1, instansi, 4 angka (pagu, HPS, kontrak, pelaksana) dalam 2 × 2, tombol Pantau. Kartu Rantai Rupiah (vertikal 8 mata) + catatan mata rantai putus. Sinyal (panel terbuka dengan strip plot, lalu sinyal kedua tertutup). Hak jawab. Bukti warga (2 kolom). Aside: pemeliharaan, detail paket, sumber. **Bilah aksi** "Lapor Bukti" + "Aspirasi" 2 kolom di atas `BottomNav` | Kartu aksi di aside disembunyikan, diganti bilah aksi |
| 768 | Sama, angka 4 kolom | Tombol Bagikan tampil | |
| 1024 | Rantai horizontal selebar konten; di bawahnya `split` 8/4 | Kiri: sinyal, hak jawab, bukti warga. Kanan: kartu aksi (Lapor Bukti `lg` lebar penuh, Paket Aspirasi sekunder) di posisi teratas, lalu pemeliharaan, detail paket, sumber | Bilah aksi bawah disembunyikan |
| 1440 | Sama | | |

### 6.5 Detail kecamatan sekolah

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom | Kepala (breadcrumb, judul, segmen jenjang, tombol tahun ajaran). Kartu grafik batang per kelurahan (lulusan vs kursi negeri, rasio di kanan) + lencana. Angka kecamatan, batang kursi tersedia, catatan keketatan. Judul daftar + chip status. `DataTable` sekolah (gulir horizontal, kolom nama lengket). Unduh CSV + muat lagi | |
| 768 | Sama | Tabel masih menggulir bila perlu | |
| 1024 | `split` 8/4 untuk grafik + aside; tabel selebar konten | | Tabel tampil penuh tanpa gulir |
| 1440 | Sama | | |

### 6.6 Cek Sekolahku

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom | Kepala: judul, paragraf, kotak alamat + "Pakai lokasiku" (bertumpuk), segmen jenjang, catatan privasi. Kartu tekanan radius (kalimat Newsreader + kelas + lencana + catatan zonasi). Peta radius 3 km tinggi 280. Daftar SMP negeri terdekat (5 baris). Alternatif swasta dan madrasah | |
| 768 | Kotak alamat dan tombol lokasi berdampingan | Peta tinggi 360 | |
| 1024 | `split` 8/4 | Kiri: kartu tekanan + peta 360. Kanan: daftar negeri + alternatif | |
| 1440 | Sama | | |

### 6.7 Profil penyedia

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom, angka 2 × 2 | Kepala: avatar inisial persegi, nama, keterangan, catatan "tanpa skor". 4 angka kunci. Grafik batang per tahun (tinggi = nilai, angka di dalam = jumlah). Daftar paket. Aside: kategori, instansi pemberi, atribusi Opentender | |
| 768 | Angka 4 kolom | | |
| 1024 | `split` 8/4 | Kiri: grafik + daftar paket. Kanan: kategori, instansi pemberi, atribusi | |
| 1440 | Sama | | |

### 6.8 Profil anggota DPRD

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom | Kepala: avatar inisial 72, nama, fraksi dan masa jabatan, penanda komisi, tombol LHKPN dan Situs DPRD, catatan tanpa skor. Kartu dapil (chip kecamatan) dan kartu komisi. Urusan komisi di dapil ini. Data dasar + lencana. Aside: ajakan Paket Aspirasi, reses berikutnya (estimasi), wakil lain | |
| 768 | Kartu dapil dan komisi 2 kolom | | |
| 1024 | `split` 8/4 | | |
| 1440 | Sama | | |

### 6.9 Tanya Rantai

| Lebar | Susunan | Catatan |
|---|---|---|
| 360 | `MobileHeader` (kembali, judul, percakapan baru). Kolom pesan, chip saran tindak lanjut (gulir), `ChatComposer` menempel di atas `BottomNav`. Panel sumber dari rujukan membuka `BottomSheet` | Chip konteks wilayah di atas pesan |
| 768 | Sama, kolom pesan maksimal 860 rata tengah | |
| 1024 | Kolom pesan (maks 860) + `SidePanel` sumber 400 di kanan, terbuka saat rujukan diklik | |
| 1440 | Sama | |

**Panel geser di Wilayahku (≥ 1024):** `SidePanel` overlay 440 di kanan dengan scrim tipis 40%, berisi kepala (ikon, judul, chip konteks, buka layar penuh, tutup), pesan, chip saran, composer. Di < 1024 tombol "Tanya Rantai" di Wilayahku membuka rute `/tanya` dengan konteks wilayah.

### 6.10 Lapor Bukti (ponsel, 6 langkah)

Alur layar penuh: `BottomNav` disembunyikan. Struktur setiap langkah: `StepHeader` di atas, isi menggulir, bilah aksi lengket di bawah (tombol utama `lg` lebar penuh + tombol teks sekunder opsional).

| Langkah | Isi | Tombol |
|---|---|---|
| 1 Izin lokasi | Peta mini 200 dengan lingkaran akurasi, judul, dua poin privasi (radius 500 m, titik publik dibulatkan 50 m) | "Izinkan lokasi", "Pilih titik di peta" |
| 1 Akurasi buruk | Lingkaran akurasi amber 180 m, catatan amber | "Coba lagi", "Geser pin di peta" |
| 2 Pilih proyek | Akurasi saat ini, `RadioCard` proyek dalam 500 m (sektor, jarak, nilai, pelaksana), opsi "Proyek tidak ada di daftar" | "Lanjut" (nonaktif sampai memilih) |
| 3 Foto | `PhotoGrid` 2 × 2 (1 sampai 4), catatan buram otomatis | "Lanjut dengan 2 foto" |
| 4 Kondisi | 4 `RadioCard`: baik, rusak ringan, rusak berat, belum terlihat dikerjakan | "Lanjut" |
| 5 Catatan | `TextArea` 280 karakter, catatan larangan nama dan tuduhan | "Lanjut", "Lewati" |
| 6 Ringkasan | Daftar ringkasan dengan tombol "Ubah" per baris, lokasi publik, centang persetujuan | "Kirim laporan" |
| Sukses | Ilustrasi mata rantai tersambung, judul h1, penjelasan 3 saksi, linimasa vertikal (dikirim, menunggu saksi 0/3, terverifikasi) | "Ajak tetangga jadi saksi", "Lihat laporan" |

Di ≥ 768 alur yang sama tampil di kolom tengah maksimal 560 dengan latar `--bg`; di ≥ 1024 `TopBar` tetap tampil.

### 6.11 Detail laporan

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom | Kepala: breadcrumb, status + kondisi, judul h1, meta (waktu, kelurahan dibulatkan, anonim). Foto 2 kolom. Catatan pelapor (kutipan). Linimasa status. Laporkan konten. Aside: `WitnessPanel`, proyek terkait (mini rantai), peta lokasi dibulatkan | `WitnessPanel` dipindah ke **setelah foto** di < 1024 agar ajakan saksi terlihat tanpa menggulir jauh (urutan DOM: kiri, lalu aside; gunakan `order` CSS hanya untuk kasus ini) |
| 768 | Sama | | |
| 1024 | `split` 8/4 | Kanan: `WitnessPanel` teratas | |
| 1440 | Sama | | |

### 6.12 Paket Aspirasi

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom | Kepala: overline, judul, masalah. Kartu penerima (checkbox). Pratinjau dokumen (`--raised`, bilah atas "Pratinjau · 2 halaman A4"). Bilah aksi "Salin" + "Unduh PDF" | Pengaturan nama dan bahasa dokumen dipindah ke `IconButton` pengaturan di header (lembar bawah) |
| 768 | Sama | Kartu pengaturan tampil di bawah penerima | |
| 1024 | `split-l` 4/8 | Kiri: penerima + pengaturan. Kanan: pratinjau. Tombol Salin, Unduh PDF, Bagikan di kepala | Bilah aksi bawah disembunyikan |
| 1440 | Sama | | |

Isi dokumen berurutan: judul (Newsreader), penerima, 1 Masalah, 2 Data dan sumber (tabel + catatan kaki bernomor + lencana kepastian), 3 Instansi berwenang, 4 Komisi pengawas, 5 Wakil (terkandung di penerima), 6 Pertanyaan yang disarankan (bernomor, netral), kaki "Disusun dengan Rantai" + tautan bukti.

### 6.13 Pantau

| Lebar | Kisi | Urutan blok | Catatan |
|---|---|---|---|
| 360 | 1 kolom | Kepala: judul, ringkasan, chip yang dipantau (gulir) + "Tambah". Chip jenis kejadian. Kejadian dikelompokkan per hari (Hari ini, Kemarin, tanggal). Aside: pengaturan kabar (switch), catatan arti kejadian | `BottomNav` aktif di Akun |
| 768 | Sama | | |
| 1024 | `split` 8/4 | Tombol "Atur notifikasi" di kepala; tautan Pantau aktif di `TopBar` | |
| 1440 | Sama | | |

Kartu kejadian: ikon lingkaran 36 berwarna menurut jenis (paket baru sekunder, melewati jadwal dan pemeliharaan segera berakhir amber, laporan terverifikasi hijau, hak jawab info), jenis + jam, judul, meta.

### 6.14 Masuk

| Lebar | Susunan |
|---|---|
| 360 | Tanpa navigasi utama. Baris atas: logo, `LangSwitch`, tutup. Formulir di tengah vertikal, maks 460 |
| 768 | Sama, formulir di tengah |
| 1024 | Dua kolom sama lebar: kiri panel `--surface` (logo, ilustrasi dua mata rantai, kalimat display "Membaca data tidak perlu akun.", paragraf, catatan independen), kanan formulir |
| 1440 | Sama |

Status formulir: email, mengirim kode (tombol memuat), OTP (6 kotak, kirim ulang dengan hitung mundur), kode kedaluwarsa (kotak error + pesan + tombol "Kirim kode baru"), terlalu banyak percobaan (kotak nonaktif, catatan `--danger-bg` dengan ikon kunci, tombol nonaktif berhitung mundur 15 menit).

### 6.15 Panel Instansi (desktop)

| Lebar | Susunan |
|---|---|
| < 1024 | Tidak didukung: tampilkan pesan "Panel Instansi dirancang untuk layar 1024 px ke atas" + tautan ke ringkasan hanya-baca |
| 1024 | `WorkspaceShell`: navigasi 240 + konten. Kepala konten (judul, keterangan, cari). Tiga angka kunci. Di bawahnya kisi `auto-fit` minimal 420: antrean sinyal, lalu kolom formulir (hak jawab, isi tanggal serah terima) turun ke baris berikut |
| 1440 | Antrean (kiri) dan kolom formulir (kanan) berdampingan |

Formulir hak jawab: catatan sinyal amber, `TextArea` 1.200 karakter ("Tampil publik apa adanya"), lampiran, centang kewenangan, "Terbitkan hak jawab" + "Simpan draf". Formulir serah terima: tanggal, nomor berita acara, unggah PDF, pratinjau `MiniChain` yang akan tersambung, "Simpan dan sambungkan".

### 6.16 Panel moderasi (desktop)

| Lebar | Susunan |
|---|---|
| < 1024 | Tidak didukung, pesan seperti 6.15 |
| 1024 | `WorkspaceShell` + kepala dengan pintasan papan ketik. Dua kolom: antrean 260 sampai 320 (thumbnail, judul, meta, status, tanda sistem) dan detail (foto 2 kolom dengan tombol sembunyikan per foto, pemeriksaan otomatis, konteks, tolak dengan alasan, aksi) |
| 1440 | Sama, kartu pemeriksaan dan konteks berdampingan |

Aksi: "Verifikasi dan teruskan" (utama `lg`), "Tolak" (garis `lg`, wajib alasan dari daftar + pesan), "Lewati" (teks). Menyembunyikan foto memakai tombol garis, bukan merah.

### 6.17 Tentang Data

| Lebar | Kisi | Urutan blok |
|---|---|---|
| 360 | 1 kolom | Kepala: overline, judul display, paragraf, chip daftar isi (gulir). Tabel sumber (gulir horizontal). Empat label kepastian (bertumpuk). Rumus kursi. Definisi sinyal. Keterbatasan. Lisensi |
| 768 | Label kepastian 2 kolom, keterbatasan dan lisensi 2 kolom | |
| 1024 | Label 4 kolom; rumus dan sinyal `split` 8/4 | Tabel tampil penuh |
| 1440 | Sama | |

### 6.18 Status khusus

| Layar | 360 | 1024 dan 1440 |
|---|---|---|
| 404 | Ilustrasi mata rantai putus, "Galat 404", judul display "Mata rantai ini tidak ditemukan.", paragraf, kotak cari, "Ke Beranda" + "Buka peta". Navigasi utama tetap | Isi di kolom tengah maks 720 |
| Error server | Ikon dalam lingkaran `--danger-bg`, "Galat 500", judul, paragraf "bukan kesalahanmu", "Coba lagi" + "Status layanan", ID kejadian (mono) + Salin | Sama, kolom tengah |
| Offline | `OfflineBanner` di bawah header, ringkasan wilayah yang dipantau dari salinan (tag "Tersimpan ..."), laporan draf menunggu kirim (tombol kirim nonaktif), catatan fitur yang butuh koneksi | Sama |
| Mode demo aktif | `DemoRibbon` lengket di atas, catatan penjelasan mode demo, contoh isi dalam `SimBox`, "Coba Lapor Bukti" + "Kembali ke data asli" | Sama |

### 6.19 Status memuat dan kosong per layar

| Layar | Memuat | Kosong |
|---|---|---|
| Beranda | Hero tetap, kartu angka jadi kerangka | Pencarian tanpa hasil + "Maksudmu Semper Barat?" + lokasiku |
| Wilayahku Layanan | Kerangka kepala, kalimat, batang | Data kursi tahun ajaran ini belum terbit + lihat tahun lalu + "Kabari saat terbit" |
| Wilayahku Uang | Kerangka angka dan kartu paket | Filter tanpa paket + "Hapus filter" + paket tingkat kota |
| Wilayahku Tanggung Jawab | Kerangka chip, rantai, kartu | Masalah belum dipetakan + kategori terdekat + usulkan kategori |
| Peta | Ubin berkilau + "Memuat peta..." + tawaran tabel | Lapisan tanpa data di area ini + perkecil peta + Lapor |
| Rantai Rupiah | Kerangka 8 mata rantai | Paket baru diumumkan: 1 mata ada, sisanya menunggu + Pantau paket |
| Kecamatan sekolah | Kerangka batang dan baris tabel | Filter tanpa sekolah (contoh: tidak ada MTs negeri) + tautan kecamatan tetangga |
| Cek Sekolahku | Kerangka kalimat dan peta + "Menghitung sekolah dalam radius 3 km..." | Tidak ada sekolah negeri dalam radius + perluas radius |
| Profil penyedia | Kerangka angka dan grafik | Belum ada paket di DKI + cara pencocokan nama |
| Profil DPRD | Kerangka kartu | Penempatan komisi belum diumumkan (pergantian antarwaktu) |
| Tanya Rantai | Kerangka gelembung + "Memuat percakapan..." | Percakapan baru dengan 4 saran pertanyaan |
| Lapor Bukti | Langkah 2: kerangka opsi + "Mencari proyek dalam radius 500 m..." | Langkah 2: tidak ada proyek dalam 500 m + perluas ke 1 km + lanjut tanpa proyek |
| Detail laporan | Kerangka foto dan teks | Laporan disembunyikan moderator + alasan + status Ditolak |
| Paket Aspirasi | "Menyusun dokumen dari 3 sumber..." + kerangka dokumen | Belum ada penerima: pesan error di kartu penerima, pratinjau kosong bergaris putus, Unduh nonaktif |
| Pantau | Kerangka kartu kejadian | Belum memantau apa pun + cari wilayah + pantau wilayah lokasiku |
| Panel Instansi | Kerangka angka dan antrean | Tidak ada sinyal terbuka + tautan laporan diteruskan |
| Panel moderasi | Kerangka antrean dan foto | Antrean kosong + tautan yang menunggu saksi |
| Tentang Data | Kerangka tabel sumber + "Memeriksa status sumber..." | Status sumber gagal dimuat (inline error, daftar tetap berlaku) |

---

## 7. Pola interaksi

### 7.1 Navigasi

- Lima tujuan utama di ponsel (Beranda, Peta, Tanya, Lapor, Akun). Pantau dan pengaturan berada di Akun.
- Desktop: logo ke Beranda, pencarian wilayah global di `TopBar` (pintasan `/`), lalu Peta, Tanya, Tentang Data. Pantau muncul setelah masuk.
- Wilayahku: tab lensa adalah rute. Mengganti tab mempertahankan posisi gulir di kepala (gulir ke atas tab, bukan ke atas halaman).
- Tombol kembali di `MobileHeader` memakai `history.back()` bila riwayat berasal dari situs yang sama, selain itu ke rute induk.
- Status filter, segmen jenjang, tahun, dan lapisan peta disimpan di query string agar tautan bisa dibagikan apa adanya.
- Setelah navigasi, fokus dipindah ke `<h1>` halaman baru dan judul halaman diumumkan (SvelteKit `afterNavigate`).

### 7.2 Lembar bawah

- Tiga titik tahan: 30% (peek), 55% (half), 92% (full). Peta: non-modal, mulai di peek. Panel sumber, filter, lapisan: modal, mulai di half.
- Seret pegangan atau isi lembar saat isi berada di posisi gulir paling atas. Lempar dengan kecepatan lebih dari 0,5 px/ms melompat ke titik tahan berikutnya.
- Animasi `--dur-enter` dengan `--ease-enter`; tutup `--dur-base` dengan `--ease-exit`.
- Papan ketik: pegangan bisa difokus, panah atas bawah mengganti titik tahan, Esc menutup lembar modal.
- Saat lembar modal terbuka, konten di belakang `inert`.

### 7.3 Panel peta

- Memilih area atau titik: area mendapat garis tepi 3,5px `--text`, penanda membesar ke 44. Detail terbuka di lembar bawah (< 1024) atau `SidePanel` (≥ 1024). Peta menggeser agar pilihan tidak tertutup panel.
- Hover di desktop memunculkan tooltip ringkas (nama + angka) setelah 150ms; klik membuka panel.
- Mengganti lapisan atau jenjang tidak menutup panel; isi panel diperbarui di tempat.
- "Lihat sebagai tabel" mengganti kanvas peta dengan `DataTable` setara (kelurahan, calon lulusan, kursi, rasio, kelas). Tombol berubah menjadi "Lihat sebagai peta". Baris terpilih di peta disorot di tabel.
- Proyek dengan lokasi tingkat kecamatan tampil sebagai `AreaHatch`, tidak pernah sebagai titik.
- Klaster terbuka saat diklik (zoom ke batasnya). Di ponsel, ketuk dua kali zoom; tombol zoom hanya di desktop.

### 7.4 Animasi rantai

- Saat `ChainStepper` masuk ke viewport pertama kali (IntersectionObserver, ambang 40%): mata rantai muncul berurutan dari awal, jeda `--dur-link-step` (60ms) per indeks, masing-masing pudar + skala 0,92 ke 1 selama `--dur-base` dengan `--ease-enter`. Sambungan digambar setelah mata di kirinya muncul.
- Mata `putus`: berhenti sejenak 120ms, lalu kedua ujung sambungan bergeser 2px menjauh.
- Saat instansi menyimpan tanggal serah terima (Panel Instansi), pratinjau `MiniChain` menutup celah dalam satu transisi `--dur-base`.
- Hanya sekali per tampilan halaman. Dengan `prefers-reduced-motion`, semua tampil sekaligus tanpa geser.

### 7.5 Memuat

- Tampilkan kerangka hanya bila respons belum tiba setelah 300ms; bila tiba lebih cepat, langsung tampilkan konten.
- Kerangka meniru bentuk konten akhir. Wadah `aria-busy="true"`, satu teks status di dalamnya.
- Tombol aksi yang memicu permintaan berpindah ke status memuat tanpa berubah lebar, dan nonaktif sampai selesai.
- Streaming Tanya Rantai: label "Menyusun jawaban dari N sumber..." lalu teks muncul per potongan dengan kursor berkedip. Rujukan bernomor ditampilkan setelah jawaban selesai.
- Data yang sudah ada tetap ditampilkan saat memuat ulang (stale while revalidate), dengan cap waktu "Diperbarui 5 menit lalu".

### 7.6 Error

- **Bahaya (merah) hanya untuk kegagalan sistem.** Penolakan moderasi, sinyal, dan data yang belum terbit bukan error.
- Error di satu bagian tidak menjatuhkan halaman: bagian itu menampilkan `ErrorState` inline dengan "Coba lagi"; bagian lain tetap tampil.
- Bila ada data setara (tabel saat peta gagal, salinan tersimpan saat luring), tawarkan.
- Toast error bertahan sampai ditutup dan memuat aksi "Coba lagi" bila bisa diulang.
- Validasi formulir terjadi saat meninggalkan kolom dan saat kirim; pesan di bawah kolom, fokus berpindah ke kolom error pertama saat kirim gagal.
- OTP: kedaluwarsa setelah 10 menit; 5 kali salah mengunci 15 menit dengan hitung mundur.
- Error server menampilkan ID kejadian yang bisa disalin.
- Luring: aplikasi memakai service worker untuk menyimpan salinan halaman wilayah yang dipantau. Laporan yang dibuat saat luring disimpan sebagai draf dan dikirim otomatis saat kembali daring, dengan toast konfirmasi.

### 7.7 Mode demo

- Diaktifkan lewat `?demo=1` atau tombol di Tentang Data, disimpan per sesi.
- `DemoRibbon` lengket di atas semua layar, tidak bisa ditutup kecuali keluar dari mode demo.
- Semua angka berasal dari set data simulasi dan dibungkus `SimBox` atau berlencana simulasi. Formulir (Lapor Bukti, Paket Aspirasi, hak jawab) bisa dicoba tetapi tidak menyimpan atau mengirim apa pun.

### 7.8 Bagikan dan salin

- "Bagikan" memakai Web Share API bila tersedia, selain itu menyalin tautan dan menampilkan toast "Tautan wilayah disalin".
- Gambar bagikan (kartu wilayah) selalu memuat lencana sumber, tanggal data, dan pita simulasi bila ada.
- "Salin" Paket Aspirasi menyalin teks polos dengan catatan kaki sumber bernomor; tombol berubah menjadi "Teks disalin" dengan ikon centang selama 3 detik.

---

## 8. Teks antarmuka (i18n)

- Berkas: `src/lib/i18n/id.json` dan `src/lib/i18n/en.json`, kunci bertitik seperti di tabel. Bahasa bawaan `id`.
- Parameter memakai kurung kurawal: `{n}`, `{name}`. Bentuk jamak bahasa Inggris memakai `Intl.PluralRules` (kunci `_one` / `_other`); bahasa Indonesia tidak mengubah kata benda.
- Teks Inggris rata-rata 15% lebih panjang. Setiap komponen harus diuji dengan teks EN: tombol boleh dua baris, lencana dan tab tidak boleh terpotong.
- Huruf kapital: kalimat biasa (sentence case) untuk tombol, judul, dan label. Nama produk dan fitur tetap: Rantai, Rantai Rupiah, Tanya Rantai, Lapor Bukti, Paket Aspirasi, Cek Sekolahku, Wilayahku, Panel Instansi.

### 8.1 Navigasi dan umum

| Kunci | Indonesia | Inggris |
|---|---|---|
| `nav.main` | Navigasi utama | Main navigation |
| `nav.breadcrumb` | Jejak halaman | Breadcrumb |
| `nav.home` | Beranda | Home |
| `nav.map` | Peta | Map |
| `nav.ask` | Tanya | Ask |
| `nav.report` | Lapor | Report |
| `nav.account` | Akun | Account |
| `nav.about` | Tentang Data | About the data |
| `nav.watch` | Pantau | Watchlist |
| `nav.skip` | Langsung ke konten | Skip to content |
| `action.signIn` | Masuk | Sign in |
| `action.back` | Kembali | Back |
| `action.close` | Tutup | Close |
| `action.share` | Bagikan | Share |
| `action.watch` | Pantau | Watch |
| `action.watching` | Dipantau | Watching |
| `action.retry` | Coba lagi | Try again |
| `action.next` | Lanjut | Continue |
| `action.skip` | Lewati | Skip |
| `action.copy` | Salin | Copy |
| `action.copied` | Teks disalin | Text copied |
| `action.downloadPdf` | Unduh PDF | Download PDF |
| `action.downloadCsv` | Unduh CSV | Download CSV |
| `action.loadMore` | Muat {n} lagi | Load {n} more |
| `action.viewAll` | Lihat semua {n} | View all {n} |
| `action.viewAsTable` | Lihat sebagai tabel | View as table |
| `action.viewAsMap` | Lihat sebagai peta | View as map |
| `action.viewSource` | Lihat sumber | View source |
| `action.openOriginal` | Buka data asli | Open original data |
| `action.readMethod` | Baca metodologi | Read the methodology |
| `action.clearFilters` | Hapus filter | Clear filters |
| `action.edit` | Ubah | Edit |
| `theme.toggle` | Ganti tema | Switch theme |
| `theme.light` | Tema terang | Light theme |
| `theme.dark` | Tema gelap | Dark theme |
| `theme.system` | Ikuti sistem | Follow system |
| `lang.label` | Bahasa | Language |
| `footer.independent` | Proyek warga independen, bukan situs resmi pemerintah. | An independent civic project, not an official government site. |
| `footer.sources` | Jejak sumber data | Data sources |
| `a11y.opensNewTab` | (membuka situs lain) | (opens another site) |
| `sheet.resize` | Ubah tinggi lembar | Resize panel |

### 8.2 Beranda dan pencarian

| Kunci | Indonesia | Inggris |
|---|---|---|
| `home.overline` | Transparansi publik · DKI Jakarta | Public transparency · Jakarta |
| `home.title` | Satu wilayah, tiga lensa: layanan, uang, dan tanggung jawab. | One area, three lenses: services, money, and accountability. |
| `home.lead` | Ketik kelurahanmu. Rantai menunjukkan apakah kursi sekolah negeri cukup, proyek apa yang dibangun di sana, dan siapa yang berwenang menindaklanjutinya. | Type your neighbourhood. Rantai shows whether there are enough public school seats, which projects are being built there, and who is responsible for following up. |
| `home.try` | Coba: | Try: |
| `home.preview` | Contoh hasil | Sample result |
| `home.lensesTitle` | Tiga pertanyaan untuk setiap wilayah | Three questions for every area |
| `home.mapTitle` | Di mana lulusan SD paling sulit mendapat kursi SMP negeri? | Where do primary graduates struggle most to find a public junior high seat? |
| `home.askTitle` | Bertanya dengan bahasa sehari-hari | Ask in everyday language |
| `search.placeholder` | Cari kelurahan atau alamat | Search a neighbourhood or address |
| `search.group.kelurahan` | Kelurahan | Neighbourhoods |
| `search.group.kecamatan` | Kecamatan | Districts |
| `search.group.street` | Jalan dan tempat | Streets and places |
| `search.loading` | Mencari "{q}"... | Searching "{q}"... |
| `search.empty` | Tidak ada wilayah bernama "{q}" | No area named "{q}" |
| `search.emptyHint` | Coba ejaan lain, nama kecamatan, atau pakai lokasimu. | Try another spelling, a district name, or use your location. |
| `search.didYouMean` | Maksudmu {name}? | Did you mean {name}? |
| `search.error` | Saran gagal dimuat. Server tidak menjawab. | Suggestions failed to load. The server did not respond. |
| `search.results_other` | {n} hasil | {n} results |
| `location.use` | Pakai lokasiku | Use my location |
| `location.useOnly` | Pakai lokasiku saja | Just use my location |
| `location.locating` | Mencari lokasi... | Finding your location... |
| `location.found` | Akurasi {m} m | Accuracy {m} m |
| `location.coarse` | Akurasi {dist}. Pilih kelurahan dari daftar. | Accuracy {dist}. Pick your neighbourhood from the list. |
| `location.denied` | Izin lokasi ditolak | Location permission denied |
| `location.deniedHint` | Ketik nama kelurahan, atau buka izin lokasi di pengaturan peramban. | Type your neighbourhood, or allow location in your browser settings. |
| `location.privacy` | Lokasimu hanya dipakai untuk mencari, tidak disimpan. | Your location is only used to search and is not stored. |

### 8.3 Lensa dan Wilayahku

| Kunci | Indonesia | Inggris |
|---|---|---|
| `lens.nav` | Lensa wilayah | Area lenses |
| `lens.services` | Layanan | Services |
| `lens.money` | Uang | Money |
| `lens.accountability` | Tanggung Jawab | Accountability |
| `lens.services.q` | Apakah kursi SMP dan SMA negeri cukup untuk lulusan di kecamatanku? | Are there enough public junior and senior high seats for graduates in my district? |
| `lens.money.q` | Proyek apa yang dibangun di sini, berapa nilainya, siapa pelaksananya, dan apakah masih dalam pemeliharaan? | Which projects are built here, how much do they cost, who builds them, and are they still under maintenance? |
| `lens.accountability.q` | Instansi mana yang berwenang, komisi DPRD mana yang mengawasi, dan siapa wakil dapilku? | Which agency is responsible, which council committee oversees it, and who represents my electoral district? |
| `area.meta` | {n} kelurahan · Dapil {dapil} · data per {date} | {n} neighbourhoods · District {dapil} · data as of {date} |
| `level.smp` | SMP | Junior high |
| `level.sma` | SMA/SMK | Senior high |
| `services.sentence` | Sekitar {a} lulusan untuk setiap {b} kursi SMP negeri. | About {a} graduates for every {b} public junior high seat. |
| `services.graduates` | Calon lulusan | Expected graduates |
| `services.seats` | Kursi kelas awal | Entry-level seats |
| `services.public` | Negeri | Public |
| `services.private` | Swasta | Private |
| `services.madrasah` | Madrasah | Madrasah |
| `services.gap` | Selisih | Shortfall |
| `services.ratioUnit` | lulusan per kursi | graduates per seat |
| `services.zoningNote` | Penerimaan murid tidak dibatasi batas kecamatan. Angka ini menggambarkan tekanan di sekitar rumah, bukan kepastian anak tidak tertampung. | Admissions are not limited by district boundaries. This figure shows pressure near home, not certainty that a child will miss out. |
| `services.schoolsTitle` | SMP negeri di {area} | Public junior high schools in {area} |
| `services.checkMySchool` | Cek Sekolahku | Check My School |
| `money.total` | Nilai paket {year} | {year} package value |
| `money.count` | Jumlah paket | Number of packages |
| `money.maintenance` | Dalam pemeliharaan | Under maintenance |
| `money.signals` | Perlu ditanyakan | Worth asking about |
| `money.bySector` | Komposisi nilai per sektor | Value by sector |
| `money.topVendors` | Penyedia terbanyak di sini | Most frequent vendors here |
| `money.noScore` | Tanpa skor. | No scores. |
| `money.vendor` | Pelaksana: {name} | Contractor: {name} |
| `acc.issueTitle` | Masalah apa yang ingin kamu sampaikan? | What issue do you want to raise? |
| `acc.chainTitle` | Rantai tanggung jawab | Chain of responsibility |
| `acc.agency` | Instansi berwenang | Responsible agency |
| `acc.legalBasis` | Dasar hukum | Legal basis |
| `acc.committee` | Komisi DPRD yang mengawasi | Overseeing council committee |
| `acc.reps` | Wakil Dapil {dapil} | Representatives for District {dapil} |
| `acc.oversees` | Duduk di komisi yang mengawasi urusan ini | Sits on the committee that oversees this issue |
| `acc.channels` | Jalur aspirasi terdekat | Nearest channels for public input |
| `acc.createPacket` | Buat Paket Aspirasi | Create a Petition Pack |
| `issue.smpSeats` | Kursi SMP kurang | Not enough junior high seats |
| `issue.drain` | Saluran tersumbat | Blocked drain |
| `issue.road` | Jalan rusak | Damaged road |
| `issue.clinic` | Puskesmas penuh | Overcrowded health centre |
| `issue.park` | Taman tidak terawat | Neglected park |
| `issue.other` | Lainnya | Other |

### 8.4 Label, sektor, kelas, status

| Kunci | Indonesia | Inggris |
|---|---|---|
| `cert.resmi` | Resmi | Official |
| `cert.proksi` | Proksi | Proxy |
| `cert.estimasi` | Estimasi | Estimate |
| `cert.simulasi` | Simulasi | Simulated |
| `cert.resmi.desc` | Diambil langsung dari sumber pemerintah tanpa diubah. | Taken directly from a government source, unchanged. |
| `cert.proksi.desc` | Dihitung Rantai dengan rumus terbuka. | Calculated by Rantai with an open formula. |
| `cert.estimasi.desc` | Perkiraan karena data resmi tidak dipublikasikan. | An estimate because official data is not published. |
| `cert.simulasi.desc` | Data demonstrasi. Tidak boleh dikutip sebagai fakta. | Demonstration data. Do not cite as fact. |
| `source.title` | Sumber angka ini | Source of this figure |
| `source.retrieved` | Diambil | Retrieved |
| `source.coverage` | Cakupan | Coverage |
| `source.formula` | Rumus | Formula |
| `source.limits` | Batasan | Limitations |
| `source.label` | Sumber: {name}, {level}, diambil {date} | Source: {name}, {level}, retrieved {date} |
| `sector.edu` | Pendidikan | Education |
| `sector.air` | Air | Water |
| `sector.san` | Sanitasi | Sanitation |
| `sector.jalan` | Jalan | Roads |
| `sector.kes` | Kesehatan | Health |
| `sector.ruang` | Ruang publik | Public space |
| `sector.lain` | Lainnya | Other |
| `gap.sangat-kurang` | Sangat kurang | Severe shortage |
| `gap.kurang` | Kurang | Shortage |
| `gap.hampir-cukup` | Hampir cukup | Nearly enough |
| `gap.cukup` | Cukup | Enough |
| `report.dikirim` | Dikirim | Submitted |
| `report.menunggu-saksi` | Menunggu saksi {n}/{m} | Awaiting witnesses {n}/{m} |
| `report.terverifikasi` | Terverifikasi | Verified |
| `report.diteruskan` | Diteruskan | Forwarded |
| `report.ditanggapi` | Ditanggapi | Responded |
| `report.selesai` | Selesai | Resolved |
| `report.ditolak` | Ditolak | Rejected |

### 8.5 Rantai Rupiah, sinyal, pemeliharaan

| Kunci | Indonesia | Inggris |
|---|---|---|
| `chain.title` | Rantai Rupiah | Rupiah Chain |
| `chain.summary` | {ok} dari {total} mata rantai tersambung | {ok} of {total} links connected |
| `chain.summaryBroken` | {ok} dari {total} tersambung, {broken} putus | {ok} of {total} connected, {broken} broken |
| `chain.stage.rencana` | Rencana | Plan |
| `chain.stage.pemilihan` | Pemilihan | Procurement |
| `chain.stage.pemenang` | Pemenang | Winner |
| `chain.stage.pelaksanaan` | Pelaksanaan | Execution |
| `chain.stage.lokasi` | Lokasi | Location |
| `chain.stage.serah-terima` | Serah terima | Handover |
| `chain.stage.pemeliharaan` | Pemeliharaan | Maintenance |
| `chain.stage.bukti` | Bukti warga | Citizen evidence |
| `chain.status.ada` | Ada | Present |
| `chain.status.putus` | Putus | Broken |
| `chain.status.menunggu` | Menunggu | Pending |
| `chain.status.estimasi` | Estimasi | Estimated |
| `chain.brokenTitle` | Mata rantai putus: {stage} | Broken link: {stage} |
| `chain.askAgency` | Tanyakan ke instansi | Ask the agency |
| `pkg.budget` | Pagu | Budget ceiling |
| `pkg.hps` | HPS | Owner's estimate |
| `pkg.contract` | Nilai kontrak | Contract value |
| `pkg.vendor` | Pelaksana | Contractor |
| `signal.sectionTitle` | Sinyal yang perlu ditanyakan | Signals worth asking about |
| `signal.beda-paket-sejenis` | Nilai berbeda dari paket sejenis | Value differs from similar packages |
| `signal.melewati-jadwal` | Melewati jadwal | Past schedule |
| `signal.peserta-sedikit` | Peserta tender sedikit | Few bidders |
| `signal.disclaimer` | Perbedaan bisa wajar. Sinyal ini mengajak bertanya, bukan menyimpulkan. | Differences can be reasonable. This signal invites questions, not conclusions. |
| `signal.how` | Cara sinyal dihitung | How signals are calculated |
| `strip.caption` | Perbandingan ini menunjukkan posisi nilai, bukan kewajaran. | This comparison shows where the value sits, not whether it is reasonable. |
| `strip.median` | Median | Median |
| `strip.thisPackage` | Paket ini | This package |
| `strip.iqr` | Separuh paket di tengah | Middle half of packages |
| `maint.title` | Masa pemeliharaan | Maintenance period |
| `maint.daysLeft` | {n} hari lagi | {n} days left |
| `maint.soon` | Segera berakhir | Ending soon |
| `maint.soonHint` | Lapor kerusakan sebelum masa ini habis. | Report damage before this period ends. |
| `maint.ended` | Berakhir {date} | Ended {date} |
| `maint.pending` | Menunggu serah terima | Awaiting handover |
| `reply.title` | Hak jawab instansi | Agency's right of reply |
| `reply.verified` | Akun instansi terverifikasi | Verified agency account |
| `reply.unedited` | Ditampilkan apa adanya tanpa disunting Rantai. | Shown as submitted, not edited by Rantai. |
| `evidence.title` | Bukti warga | Citizen evidence |
| `evidence.blurNote` | Wajah dan pelat nomor diburamkan otomatis. Foto tampil publik setelah 3 saksi. | Faces and number plates are blurred automatically. Photos go public after 3 witnesses. |

### 8.6 Peta

| Kunci | Indonesia | Inggris |
|---|---|---|
| `map.layers` | Lapisan peta | Map layers |
| `map.layer.gap` | Kesenjangan kursi | Seat gap |
| `map.layer.projects` | Proyek | Projects |
| `map.layer.reports` | Laporan warga | Citizen reports |
| `map.layer.schools` | Sekolah | Schools |
| `map.legend` | Lulusan per kursi SMP negeri | Graduates per public junior high seat |
| `map.legendShort` | kurang ke cukup | short to enough |
| `map.loading` | Memuat peta... | Loading map... |
| `map.tableWhileLoading` | Data wilayah tetap bisa dibuka sebagai tabel selagi peta dimuat. | Area data can still be opened as a table while the map loads. |
| `map.emptyReports` | Belum ada laporan warga di area ini | No citizen reports in this area yet |
| `map.districtLevel` | Lokasi tingkat kecamatan | District-level location |
| `map.zoomIn` | Perbesar | Zoom in |
| `map.zoomOut` | Perkecil | Zoom out |
| `map.toMyLocation` | Ke lokasiku | Go to my location |
| `map.openArea` | Buka Wilayahku | Open area page |

### 8.7 Tanya Rantai

| Kunci | Indonesia | Inggris |
|---|---|---|
| `ask.title` | Tanya Rantai | Ask Rantai |
| `ask.emptyTitle` | Tanyakan apa saja tentang wilayahmu | Ask anything about your area |
| `ask.emptyLead` | Rantai menjawab hanya dari data yang dimilikinya, lengkap dengan rujukan bernomor. Bila datanya tidak ada, Rantai akan bilang. | Rantai only answers from data it holds, with numbered references. If the data is missing, Rantai will say so. |
| `ask.placeholder` | Tulis pertanyaanmu | Type your question |
| `ask.send` | Kirim | Send |
| `ask.context` | Konteks: {area} | Context: {area} |
| `ask.streaming` | Menyusun jawaban dari {n} sumber... | Drafting an answer from {n} sources... |
| `ask.noData` | Rantai belum punya data ini. | Rantai does not have this data yet. |
| `ask.caution` | Bisa keliru. Selalu cek rujukannya. | Can be wrong. Always check the references. |
| `ask.citation` | Rujukan {n}: {source} | Reference {n}: {source} |
| `ask.newChat` | Percakapan baru | New conversation |
| `ask.fullScreen` | Buka layar penuh | Open full screen |

### 8.8 Lapor Bukti dan detail laporan

| Kunci | Indonesia | Inggris |
|---|---|---|
| `report.flowTitle` | Lapor Bukti | Report Evidence |
| `report.step` | Langkah {n} dari {total} | Step {n} of {total} |
| `report.s1.title` | Izinkan lokasi agar laporan menempel ke proyek yang benar | Allow location so your report is linked to the right project |
| `report.s1.radius` | Lokasi dipakai untuk mencari proyek dalam radius 500 m darimu. | Your location is used to find projects within 500 m of you. |
| `report.s1.rounded` | Titik tepat tidak ditampilkan publik. Dibulatkan ke 50 m. | Your exact point is never shown publicly. It is rounded to 50 m. |
| `report.s1.allow` | Izinkan lokasi | Allow location |
| `report.s1.pickOnMap` | Pilih titik di peta | Pick a point on the map |
| `report.s1.poorTitle` | Lokasimu belum cukup akurat | Your location is not accurate enough |
| `report.s1.poorBody` | Akurasi lokasimu {m} m, terlalu kasar untuk mencocokkan proyek. Pindah ke tempat terbuka, atau geser pin di peta. | Your location accuracy is {m} m, too coarse to match a project. Move to an open area, or drag the pin on the map. |
| `report.s2.title` | Proyek mana yang kamu lihat? | Which project are you looking at? |
| `report.s2.notListed` | Proyek tidak ada di daftar | Project not in the list |
| `report.s2.searching` | Mencari proyek dalam radius 500 m... | Looking for projects within 500 m... |
| `report.s2.none` | Tidak ada proyek tercatat dalam 500 m | No recorded projects within 500 m |
| `report.s2.widen` | Perluas ke 1 km | Widen to 1 km |
| `report.s3.title` | Ambil 1 sampai 4 foto | Take 1 to 4 photos |
| `report.s3.take` | Ambil foto | Take photo |
| `report.s3.continueWith_other` | Lanjut dengan {n} foto | Continue with {n} photos |
| `report.s4.title` | Bagaimana kondisinya sekarang? | What condition is it in now? |
| `report.cond.good` | Baik dan berfungsi | Good and working |
| `report.cond.minor` | Rusak ringan | Minor damage |
| `report.cond.major` | Rusak berat atau tidak berfungsi | Major damage or not working |
| `report.cond.notBuilt` | Belum terlihat dikerjakan | No visible work yet |
| `report.s5.title` | Ceritakan yang terlihat | Describe what you see |
| `report.s5.hint` | Tulis apa yang terlihat, bukan dugaan. | Write what you see, not what you suspect. |
| `report.s5.rules` | Jangan cantumkan nama orang, nomor telepon, atau tuduhan. | Do not include people's names, phone numbers, or accusations. |
| `report.s6.title` | Periksa sebelum mengirim | Check before sending |
| `report.s6.consent` | Saya paham foto tampil publik setelah 3 saksi mengonfirmasi. | I understand photos go public after 3 witnesses confirm. |
| `report.s6.submit` | Kirim laporan | Send report |
| `report.done.title` | Laporan terkirim | Report sent |
| `report.done.body` | Agar terverifikasi, dibutuhkan 3 saksi: warga lain yang melihat kondisi yang sama dalam radius 200 m. | To be verified, it needs 3 witnesses: other residents who see the same condition within 200 m. |
| `report.done.invite` | Ajak tetangga jadi saksi | Invite neighbours to witness |
| `report.witness.cta` | Saya juga melihat ini | I see this too |
| `report.witness.rule` | Tombol aktif bila kamu sudah masuk dan berada dalam 200 m dari lokasi. | Available when you are signed in and within 200 m of the location. |
| `report.hidden` | Laporan ini disembunyikan | This report is hidden |
| `report.flag` | Laporkan konten tidak pantas | Report inappropriate content |

### 8.9 Paket Aspirasi, Pantau, Masuk

| Kunci | Indonesia | Inggris |
|---|---|---|
| `packet.title` | Paket Aspirasi | Petition Pack |
| `packet.recipients` | Ditujukan kepada | Addressed to |
| `packet.noRecipient` | Pilih minimal satu penerima agar paket bisa dibuat. | Choose at least one recipient to create the pack. |
| `packet.includeName` | Sertakan namaku | Include my name |
| `packet.anonymous` | Bawaan: anonim | Default: anonymous |
| `packet.docLang` | Bahasa dokumen | Document language |
| `packet.preview` | Pratinjau · {n} halaman A4 | Preview · {n} A4 pages |
| `packet.composing` | Menyusun dokumen dari {n} sumber... | Composing the document from {n} sources... |
| `packet.s.problem` | Masalah | Problem |
| `packet.s.data` | Data dan sumber | Data and sources |
| `packet.s.agency` | Instansi berwenang | Responsible agency |
| `packet.s.committee` | Komisi pengawas | Overseeing committee |
| `packet.s.questions` | Pertanyaan yang disarankan | Suggested questions |
| `packet.footer` | Disusun dengan Rantai. | Prepared with Rantai. |
| `watch.title` | Pantau | Watchlist |
| `watch.summary` | {areas} wilayah dan {pkgs} paket · {n} kejadian baru | {areas} areas and {pkgs} packages · {n} new events |
| `watch.add` | Tambah | Add |
| `watch.today` | Hari ini | Today |
| `watch.yesterday` | Kemarin | Yesterday |
| `watch.ev.newPackage` | Paket baru | New package |
| `watch.ev.pastSchedule` | Melewati jadwal | Past schedule |
| `watch.ev.maintEnding` | Pemeliharaan segera berakhir | Maintenance ending soon |
| `watch.ev.reportVerified` | Laporan terverifikasi | Report verified |
| `watch.ev.newReply` | Hak jawab baru | New right of reply |
| `watch.emptyTitle` | Belum ada yang kamu pantau | You are not watching anything yet |
| `watch.emailDigest` | Email ringkasan mingguan | Weekly email digest |
| `auth.title` | Masuk ke Rantai | Sign in to Rantai |
| `auth.lead` | Tanpa kata sandi. Kami kirim kode 6 digit ke emailmu. | No password. We will email you a 6 digit code. |
| `auth.email` | Email | Email |
| `auth.sendCode` | Kirim kode | Send code |
| `auth.sending` | Mengirim kode... | Sending code... |
| `auth.agency` | Masuk sebagai petugas instansi | Sign in as agency staff |
| `auth.codeTitle` | Masukkan kode | Enter the code |
| `auth.codeSent` | Kode 6 digit dikirim ke {email}. Berlaku 10 menit. | A 6 digit code was sent to {email}. Valid for 10 minutes. |
| `auth.codeLabel` | Kode masuk | Sign-in code |
| `auth.resendIn` | Kirim ulang kode dalam {time} | Resend code in {time} |
| `auth.changeEmail` | Ganti email | Change email |
| `auth.expired` | Kode sudah kedaluwarsa. Minta kode baru, lalu masukkan dalam 10 menit. | This code has expired. Request a new one and enter it within 10 minutes. |
| `auth.newCode` | Kirim kode baru | Send a new code |
| `auth.locked` | Terlalu banyak percobaan. Demi keamanan akunmu, coba lagi dalam 15 menit. | Too many attempts. For your account's safety, try again in 15 minutes. |
| `auth.lockedButton` | Coba lagi dalam {time} | Try again in {time} |
| `auth.noAccountNeeded` | Membaca data tidak perlu akun. | You do not need an account to read the data. |

### 8.10 Status dan error

| Kunci | Indonesia | Inggris |
|---|---|---|
| `state.loadingSummary` | Memuat ringkasan data terbaru... | Loading the latest data summary... |
| `state.updatedAgo` | Diperbarui {time} | Updated {time} |
| `state.notPublished` | Data {period} belum terbit | Data for {period} is not published yet |
| `state.notifyWhenPublished` | Kabari saat terbit | Notify me when published |
| `error.404.title` | Mata rantai ini tidak ditemukan. | This link in the chain is missing. |
| `error.404.body` | Halaman yang kamu cari mungkin sudah dipindah, atau tautannya terpotong saat dibagikan. | The page may have moved, or the link was cut off when shared. |
| `error.500.title` | Server Rantai sedang bermasalah | Rantai's server is having trouble |
| `error.500.body` | Ini bukan kesalahanmu. Data yang sudah tersimpan tetap aman. Coba lagi dalam beberapa menit. | This is not your fault. Saved data is safe. Please try again in a few minutes. |
| `error.incident` | ID kejadian | Incident ID |
| `error.sectionFailed` | Angka gagal dimuat | Figures failed to load |
| `offline.banner` | Kamu sedang luring. Menampilkan salinan tersimpan {time}. | You are offline. Showing a copy saved {time}. |
| `offline.draftPending` | {n} laporan menunggu dikirim | {n} report waiting to be sent |
| `demo.ribbon` | Mode demo aktif · semua angka di layar ini contoh | Demo mode · all figures on this screen are samples |
| `demo.exit` | Keluar | Exit |
| `sim.ribbon` | Data simulasi | Simulated data |
| `toast.reportSent` | Laporan terkirim. Butuh {n} saksi lagi agar terverifikasi. | Report sent. Needs {n} more witnesses to be verified. |
| `toast.linkCopied` | Tautan wilayah disalin. | Area link copied. |
| `toast.uploadFailed` | Foto gagal diunggah. Koneksi terputus di {pct}. | Photo failed to upload. Connection dropped at {pct}. |

---

## 9. Aturan konten

### 9.1 Kosakata netral

Rantai menyajikan data dan mengajak bertanya. Rantai tidak menuduh dan tidak menyimpulkan kesalahan siapa pun.

**Kata terlarang** di seluruh antarmuka, contoh data, dokumen ekspor, dan jawaban Tanya Rantai (dalam bahasa Indonesia maupun padanannya dalam bahasa Inggris seperti *corruption*, *fraud*, *mark-up*, *embezzlement*, *fictitious*):

| Jangan | Pakai |
|---|---|
| korupsi, penyelewengan | perlu ditanyakan |
| mark up, penggelembungan | berbeda dari paket sejenis |
| curang, kongkalikong | (tidak dipakai; tampilkan sinyal yang sesuai datanya, misalnya peserta tender sedikit) |
| fiktif | belum terlihat dikerjakan, data belum ada |
| mangkrak, terbengkalai | melewati jadwal |
| merugikan negara | (tidak dipakai; Rantai tidak menilai kerugian) |

Aturan tambahan:

- **Judul sinyal hanya dari katalog** (`signal.*`). Tidak ada judul sinyal yang ditulis bebas.
- **Setiap sinyal** didampingi data pembanding, kalimat "Perbedaan bisa wajar..." dan jalan ke hak jawab instansi.
- **Hak jawab** tampil setara dengan data, tidak disunting, tidak diletakkan sebagai catatan kaki.
- **Penolakan moderasi** ditulis sebagai fakta beralasan ("Lokasi foto berjarak 640 m dari proyek yang dipilih."), tanpa menilai pelapor.
- **Wakil rakyat:** tanpa skor, peringkat, atau warna partai. Avatar inisial, bukan foto.
- **Penyedia:** tanpa skor. Tampilkan riwayat, atribusi Opentender bila dipakai.
- **Nama di contoh dan data simulasi** selalu rekaan: "CV Contoh Simulasi", "PT Contoh Bangun Sejahtera", "Anggota Contoh Satu", "Fraksi Contoh A", "SMP Negeri Contoh 1". Nama wilayah memakai nama Jakarta yang nyata dengan angka contoh.
- **Laporan warga:** tulis apa yang terlihat, bukan dugaan. Catatan yang memuat nama orang, nomor telepon, atau tuduhan disunting moderator.
- **Gaya bahasa:** kalimat pendek, "kamu" untuk menyapa warga, kalimat aktif, tanpa jargon pengadaan tanpa penjelasan (HPS selalu punya tooltip "Harga Perkiraan Sendiri, perkiraan biaya yang dibuat instansi sebelum tender").
- **Tanda baca:** jangan memakai tanda pisah panjang (em dash) di teks antarmuka. Pakai koma, titik dua, kurung, atau kata "sampai" untuk rentang ("1,3 sampai 2,0").

### 9.2 Menulis angka

| Jenis | Indonesia (`id`) | Inggris (`en`) |
|---|---|---|
| Ribuan | 1.284 kursi | 1,284 seats |
| Desimal | 2,1 lulusan per kursi | 2.1 graduates per seat |
| Persen | 12,5% | 12.5% |
| Uang di bawah 1 juta | Rp845.000 | IDR 845,000 |
| Uang juta | Rp845 juta | IDR 845 million |
| Uang miliar | Rp1,2 miliar | IDR 1.2 billion |
| Uang triliun | Rp1,4 triliun | IDR 1.4 trillion |
| Uang ringkas (ruang sempit: sumbu grafik, kartu kecil, mini daftar) | Rp845 jt · Rp1,2 M · Rp1,4 T | IDR 845M · IDR 1.2B · IDR 1.4T |
| Uang lengkap (panel sumber, tabel, ekspor) | Rp2.408.700.000 | IDR 2,408,700,000 |
| Rentang | 1,3 sampai 2,0 | 1.3 to 2.0 |
| Perkiraan | ±3.150 calon lulusan | ±3,150 expected graduates |
| Jarak | 1,2 km · 500 m | 1.2 km · 500 m |
| Durasi | 47 hari lagi | 47 days left |

Aturan:

- "Rp" langsung menempel ke angka tanpa spasi. "IDR" diikuti spasi.
- Satu desimal untuk ringkasan (kartu, kalimat), dua desimal untuk detail paket (pagu, HPS, kontrak), nilai lengkap di panel sumber dan tabel.
- Satuan selalu ditulis. Angka tidak pernah tampil tanpa lencana sumber dalam satu tatapan.
- Semua angka data memakai kelas `.num` (tabular).
- Rasio ditulis sebagai kalimat manusiawi di kalimat ringkas ("Sekitar 2 lulusan untuk setiap 1 kursi"), angka presisi di kartu ("2,1 lulusan per kursi").

Pembantu format di `src/lib/format.ts`:

```ts
import type { Locale } from '$lib/types/ui';

const tag = (l: Locale) => (l === 'id' ? 'id-ID' : 'en-GB');

export const formatNumber = (n: number, l: Locale, digits = 0) =>
  new Intl.NumberFormat(tag(l), { minimumFractionDigits: digits, maximumFractionDigits: digits }).format(n);

const SCALES = [
  { v: 1e12, id: 'triliun', idShort: 'T', en: 'trillion', enShort: 'T' },
  { v: 1e9, id: 'miliar', idShort: 'M', en: 'billion', enShort: 'B' },
  { v: 1e6, id: 'juta', idShort: 'jt', en: 'million', enShort: 'M' }
] as const;

/** formatRupiah(1_200_000_000, 'id') -> "Rp1,2 miliar"; (…, 'en') -> "IDR 1.2 billion" */
export function formatRupiah(
  n: number,
  l: Locale,
  opts: { style?: 'words' | 'short' | 'full'; digits?: number } = {}
): string {
  const { style = 'words', digits = 1 } = opts;
  const prefix = l === 'id' ? 'Rp' : 'IDR ';
  if (style === 'full' || n < 1e6) return prefix + formatNumber(n, l);
  const s = SCALES.find((x) => n >= x.v)!;
  const num = formatNumber(n / s.v, l, Number.isInteger(n / s.v) ? 0 : digits);
  if (style === 'short') {
    const unit = l === 'id' ? s.idShort : s.enShort;
    return l === 'id' ? `${prefix}${num} ${unit}` : `${prefix}${num}${unit}`;
  }
  return `${prefix}${num} ${l === 'id' ? s.id : s.en}`;
}

export const formatPercent = (ratio: number, l: Locale, digits = 1) =>
  new Intl.NumberFormat(tag(l), { style: 'percent', maximumFractionDigits: digits }).format(ratio);
```

Catatan: `Intl.NumberFormat('id-ID', { style: 'currency' })` menghasilkan "Rp 1.200" dengan spasi tak terputus, jadi jangan dipakai untuk rupiah; pakai `formatRupiah`.

### 9.3 Tanggal dan waktu

| Jenis | Indonesia | Inggris | Pemakaian |
|---|---|---|---|
| Pendek | 12 Okt 2026 | 12 Oct 2026 | Kartu, lencana, tabel |
| Panjang | 12 Oktober 2026 | 12 October 2026 | Kalimat, dokumen, nama aksesibel |
| Dengan jam | 12 Okt 2026, 08.14 WIB | 12 Oct 2026, 08:14 WIB | Panel sumber, linimasa |
| Jam saja | 08.14 | 08:14 | Feed hari ini |
| Bulan | Mar 2026 | Mar 2026 | Rantai Rupiah |
| Rentang | 5 Mei sampai 2 Sep 2026 | 5 May to 2 Sep 2026 | Masa pelaksanaan |
| Relatif (di bawah 7 hari) | 2 jam lalu · kemarin · 3 hari lalu | 2 hours ago · yesterday · 3 days ago | Feed Pantau, cap "Diperbarui" |
| Tahun ajaran | 2026/2027 | 2026/27 | Data sekolah |

Aturan:

- Zona waktu selalu Asia/Jakarta dan ditulis "WIB" bila jam ditampilkan.
- Format 24 jam di kedua bahasa. Pemisah jam bahasa Indonesia titik (08.14), bahasa Inggris titik dua (08:14).
- Lebih dari 7 hari: tampilkan tanggal pendek, bukan waktu relatif.
- Tanggal di nama aksesibel memakai format panjang ("diambil 14 Oktober 2026").

```ts
export const formatDate = (d: Date, l: Locale, style: 'short' | 'long' = 'short') =>
  new Intl.DateTimeFormat(tag(l), {
    day: 'numeric',
    month: style === 'short' ? 'short' : 'long',
    year: 'numeric',
    timeZone: 'Asia/Jakarta'
  }).format(d).replace(/\.$/, '');

export const formatTime = (d: Date, l: Locale) => {
  const t = new Intl.DateTimeFormat('en-GB', { hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'Asia/Jakarta' }).format(d);
  return l === 'id' ? t.replace(':', '.') : t;
};

export const formatRelative = (d: Date, l: Locale, now = new Date()) => {
  const rtf = new Intl.RelativeTimeFormat(tag(l), { numeric: 'auto' });
  const mins = Math.round((d.getTime() - now.getTime()) / 60000);
  if (Math.abs(mins) < 60) return rtf.format(mins, 'minute');
  const hours = Math.round(mins / 60);
  if (Math.abs(hours) < 24) return rtf.format(hours, 'hour');
  const days = Math.round(hours / 24);
  return Math.abs(days) < 7 ? rtf.format(days, 'day') : formatDate(d, l);
};
```

Periksa keluaran `Intl` untuk singkatan bulan bahasa Indonesia ("Okt", "Agu", "Des") di peramban target; bila berbeda dari tabel, pakai daftar bulan sendiri.

### 9.4 Daftar periksa sebelum rilis layar

- [ ] Semua angka punya satuan dan lencana kepastian.
- [ ] Setiap grafik dan peta punya `TableToggle` yang menghasilkan data setara.
- [ ] Angka simulasi berada di `SimBox` atau berlencana simulasi.
- [ ] Tidak ada kata terlarang, tidak ada merah selain error sistem.
- [ ] Target sentuh 44 × 44, fokus terlihat, urutan fokus mengikuti urutan baca.
- [ ] Diuji di 360, 768, 1024, 1440, tema terang dan gelap, bahasa ID dan EN.
- [ ] Status memuat, kosong, dan error tersedia.
- [ ] Kontras teks minimal 4,5:1 (3:1 untuk elemen non-teks dan teks ≥ 24px).
