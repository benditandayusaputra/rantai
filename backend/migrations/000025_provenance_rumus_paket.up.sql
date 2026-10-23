INSERT INTO provenance (id, sumber, url, diambil_pada, label, catatan) VALUES
    ('prv_rantai_pemeliharaan', 'rantai', 'https://github.com/benditandayusaputra/rantai/blob/main/backend/config/pemeliharaan.yaml',
     '2026-10-23T00:00:00+07:00', 'estimasi',
     'akhir_pelaksanaan = tanggal serah terima dari instansi, atau akhir bulan jadwal pelaksanaan kontrak SiRUP, atau tanggal kontrak + durasi bawaan kategori; akhir_pemeliharaan = akhir_pelaksanaan + masa pemeliharaan kategori (permanen 180 hari, semi permanen 90 hari)'),
    ('prv_rantai_lokasi_paket', 'rantai', 'https://github.com/benditandayusaputra/rantai/blob/main/backend/config/geocode.yaml',
     '2026-10-23T00:00:00+07:00', 'estimasi',
     'Titik paket adalah titik tengah kelurahan yang disebut di nama paket, digeser deterministik maksimal 150 m; jarak dihitung garis lurus')
ON CONFLICT (id) DO NOTHING;
