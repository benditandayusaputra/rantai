import { describe, expect, it } from 'vitest';
import { bacaBagian, inisial, tampilkanAngka, tautanBagian } from './kecamatan';

describe('bagian halaman kecamatan', () => {
	it('membaca bagian dari URL dengan tema sebagai bawaan', () => {
		expect(bacaBagian('wakilku')).toBe('wakilku');
		expect(bacaBagian(null)).toBe('tema');
		expect(bacaBagian('lain')).toBe('tema');
	});

	it('mempertahankan kelurahan yang disorot saat berpindah bagian', () => {
		const url = new URL('http://x/id/kecamatan/31.72.06?kelurahan=31.72.06.1001&bagian=wakilku');
		expect(tautanBagian(url, 'anggaran')).toBe(
			'/id/kecamatan/31.72.06?kelurahan=31.72.06.1001&bagian=anggaran'
		);
		expect(tautanBagian(url, 'tema')).toBe('/id/kecamatan/31.72.06?kelurahan=31.72.06.1001');
	});
});

describe('tampilkanAngka', () => {
	const dasar = { label: 'resmi', provenance_id: 'p' } as const;

	it('menulis rupiah, rasio, dan hitungan sesuai satuannya', () => {
		expect(tampilkanAngka({ ...dasar, nilai: 1_200_000_000, satuan: 'rupiah' })).toEqual({
			nilai: 'Rp1,2 miliar',
			satuan: ''
		});
		expect(tampilkanAngka({ ...dasar, nilai: 0.476, satuan: 'rasio' }).nilai).toBe('0,48');
		expect(tampilkanAngka({ ...dasar, nilai: 1284, satuan: 'kursi' })).toEqual({
			nilai: '1.284',
			satuan: 'kursi'
		});
	});
});

describe('inisial', () => {
	it('mengambil dua huruf kapital pertama tanpa gelar kecil', () => {
		expect(inisial('Hj. Anggota Contoh Satu')).toBe('AC');
		expect(inisial('H. Ramly Hi Muhamad, S.Sos')).toBe('RH');
		expect(inisial('dr. Budi Santoso')).toBe('BS');
	});
});
