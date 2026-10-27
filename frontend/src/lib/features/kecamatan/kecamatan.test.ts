import type { IndikatorRingkas } from '$lib/api/tipe';
import { describe, expect, it } from 'vitest';
import { bacaBagian, inisial, statusLayanan, tampilkanAngka, tautanBagian } from './kecamatan';

describe('bagian halaman kecamatan', () => {
	it('membuka pilihan informasi bila bagian belum dipilih', () => {
		expect(bacaBagian('wakilku')).toBe('wakilku');
		expect(bacaBagian(null)).toBeNull();
		expect(bacaBagian('lain')).toBeNull();
	});

	it('mempertahankan kelurahan yang disorot saat berpindah bagian', () => {
		const url = new URL('http://x/id/kecamatan/31.72.06?kelurahan=31.72.06.1001&bagian=wakilku');
		expect(tautanBagian(url, 'anggaran')).toBe(
			'/id/kecamatan/31.72.06?kelurahan=31.72.06.1001&bagian=anggaran'
		);
		expect(tautanBagian(url, 'tema')).toBe(
			'/id/kecamatan/31.72.06?kelurahan=31.72.06.1001&bagian=tema'
		);
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

const indikator = (
	kebutuhan: number,
	kapasitas: number,
	kelas: IndikatorRingkas['kelas'] = 'cukup'
): IndikatorRingkas => {
	const angka = (nilai: number) => ({
		nilai,
		satuan: 'siswa',
		label: 'resmi' as const,
		provenance_id: 'p'
	});
	return {
		periode: '2026/2027',
		kebutuhan: angka(kebutuhan),
		kapasitas: angka(kapasitas),
		kesenjangan: angka(kebutuhan - kapasitas),
		rasio: null,
		kelas
	};
};

describe('statusLayanan', () => {
	it('tidak menganggap data kosong atau kebutuhan nol sebagai kondisi baik', () => {
		expect(statusLayanan(null)).toBe('belum-dinilai');
		expect(statusLayanan(indikator(0, 100))).toBe('belum-dinilai');
		expect(statusLayanan(indikator(100, 100, 'tidak_ada_data'))).toBe('belum-dinilai');
		expect(statusLayanan(indikator(100, Number.NaN))).toBe('belum-dinilai');
	});
	it('membedakan kekurangan, kapasitas seimbang, dan kapasitas lebih', () => {
		expect(statusLayanan(indikator(100, 49, 'sangat_kurang'))).toBe('perhatian');
		expect(statusLayanan(indikator(100, 99, 'hampir_cukup'))).toBe('perhatian');
		expect(statusLayanan(indikator(100, 100))).toBe('memadai');
		expect(statusLayanan(indikator(100, 130))).toBe('memadai');
	});
});
