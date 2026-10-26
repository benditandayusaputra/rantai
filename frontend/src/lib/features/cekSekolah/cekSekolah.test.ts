import { describe, expect, it } from 'vitest';
import type { SekolahTerdekat } from '$lib/api/tipe';
import { pecahanTekanan, radiusCekSekolahM, radiusKm, terdekatDulu } from './cekSekolah';

function sekolah(npsn: string, jarak_m: number): SekolahTerdekat {
	return {
		npsn,
		nama: npsn,
		bentuk: 'SMP',
		status: 'negeri',
		lokasi: null,
		jalur: 'negeri',
		daya_tampung: null,
		jarak_m
	};
}

describe('cekSekolah', () => {
	it('menyatakan rasio sebagai pecahan yang mudah dibaca', () => {
		expect(pecahanTekanan(2.04)).toEqual({ a: '2,0', b: '1' });
		expect(pecahanTekanan(0.5)).toEqual({ a: '1', b: '2,0' });
	});

	it('mengurutkan sekolah dari yang terdekat dan membatasi jumlahnya', () => {
		expect(
			terdekatDulu([sekolah('c', 900), sekolah('a', 100), sekolah('b', 400)], 2).map((s) => s.npsn)
		).toEqual(['a', 'b']);
	});

	it('menulis radius dalam kilometer', () => {
		expect(radiusKm(3000)).toBe('3');
	});

	it('memakai radius 5 km untuk SMP dan 10 km untuk SMA/SMK', () => {
		expect(radiusKm(radiusCekSekolahM.pendidikan)).toBe('5');
		expect(radiusKm(radiusCekSekolahM.pendidikan_sma_smk)).toBe('10');
	});
});
