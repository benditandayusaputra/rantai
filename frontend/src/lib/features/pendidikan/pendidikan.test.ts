import { describe, expect, it } from 'vitest';
import type { Angka, DetailPendidikan, Sekolah } from '$lib/api/tipe';
import { kunciJenisSekolah, persenLebar, susunBarisGrafik } from './pendidikan';

const angka = (nilai: number, label: Angka['label'] = 'resmi'): Angka => ({
	nilai,
	satuan: 'kursi',
	label,
	provenance_id: 'p'
});

const detail = {
	calon_lulusan: angka(2701),
	kursi_negeri: angka(1284),
	kursi_spmb_bersama: null,
	kursi_madrasah: angka(301, 'proksi')
} as DetailPendidikan;

describe('susunBarisGrafik', () => {
	it('melewati jenis kursi yang tidak tersedia dan mempertahankan labelnya', () => {
		const baris = susunBarisGrafik(detail);
		expect(baris.map((b) => b.kunci)).toEqual([
			'pen.graduates',
			'pen.seats.negeri',
			'pen.seats.madrasah'
		]);
		expect(baris[2]?.angka.label).toBe('proksi');
	});
});

describe('persenLebar', () => {
	it('menskalakan terhadap nilai terbesar tanpa membagi nol', () => {
		expect(persenLebar(1284, 2701)).toBe(48);
		expect(persenLebar(0, 0)).toBe(0);
	});
});

describe('kunciJenisSekolah', () => {
	const dasar = { npsn: '1', nama: 'x', lokasi: null, daya_tampung: null } as const;
	it.each([
		[{ ...dasar, bentuk: 'SMP', status: 'negeri', jalur: 'negeri' }, 'pen.school.negeri'],
		[{ ...dasar, bentuk: 'SMP', status: 'swasta', jalur: 'spmb_bersama' }, 'pen.school.spmb'],
		[{ ...dasar, bentuk: 'SMP', status: 'swasta', jalur: null }, 'pen.school.swasta'],
		[{ ...dasar, bentuk: 'MTs', status: 'negeri', jalur: null }, 'pen.school.madrasah']
	] as [Sekolah, string][])('%o menjadi %s', (sekolah, harapan) => {
		expect(kunciJenisSekolah(sekolah)).toBe(harapan);
	});
});
