import { describe, expect, it } from 'vitest';
import type { BelanjaKecamatan, BelanjaTema } from '$lib/api/tipe';
import { jumlahTanpaRasio, titikSebar, ukuranUntuk, urutBaris } from './belanja';

const ukuranGrafik = ukuranUntuk(640);

function angka(nilai: number) {
	return { nilai, satuan: 'rupiah', label: 'resmi' as const, provenance_id: 'prv' };
}

function kec(kode: string, nama: string, rasio: number | null, hps: number): BelanjaKecamatan {
	return {
		kecamatan: { kode, nama, tingkat: 'kecamatan' },
		kode_kota: '31.72',
		indikator:
			rasio === null
				? null
				: {
						periode: '2026/2027',
						kebutuhan: angka(1),
						kapasitas: angka(1),
						kesenjangan: angka(0),
						rasio: angka(rasio),
						kelas: 'kurang'
					},
		belanja: { tahun: 2026, jumlah_paket: angka(1), total_hps: angka(hps) }
	};
}

const data: BelanjaTema = {
	tema: 'pendidikan',
	tahun: 2026,
	luar_kecamatan: [],
	kecamatan: [
		kec('a', 'Cilincing', 0.5, 0),
		kec('b', 'Koja', 1, 100),
		kec('c', 'Tanjung Priok', null, 50)
	]
};

describe('belanja', () => {
	it('menggambar hanya kecamatan berasio dan menaruh kecamatan aktif paling akhir', () => {
		const titik = titikSebar(data, 'a', ukuranGrafik);
		expect(titik.map((t) => t.kode)).toEqual(['b', 'a']);
		expect(jumlahTanpaRasio(data)).toBe(1);
	});

	it('memakai tinggi lebih rendah di layar sempit tanpa menyusut di bawah 280 px', () => {
		expect(ukuranUntuk(200).lebar).toBe(280);
		expect(ukuranUntuk(320).tinggi).toBeLessThan(ukuranUntuk(900).tinggi);
	});

	it('memetakan nol belanja ke garis dasar dan nilai terbesar ke dalam area plot', () => {
		const [koja, cilincing] = titikSebar(data, 'a', ukuranGrafik);
		expect(cilincing?.y).toBe(ukuranGrafik.tinggi - ukuranGrafik.bawah);
		expect(koja?.y).toBeGreaterThan(ukuranGrafik.atas);
		expect(koja?.x).toBeLessThan(ukuranGrafik.lebar - ukuranGrafik.kanan);
	});

	it('menaruh kecamatan aktif di baris pertama tabel lalu urut nama', () => {
		expect(urutBaris(data, 'c').map((k) => k.kecamatan.nama)).toEqual([
			'Tanjung Priok',
			'Cilincing',
			'Koja'
		]);
	});
});
