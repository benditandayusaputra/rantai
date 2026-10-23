import { describe, expect, it } from 'vitest';
import { bacaLapisan, radiusTampilanM, tulisLapisan } from './lapisan';

describe('status lapisan di URL', () => {
	it('menyalakan lapisan tema bila URL tidak menyebut lapisan', () => {
		expect(bacaLapisan(null)).toEqual({ tema: true, proyek: false, sekolah: false });
	});

	it('membaca dan menulis daftar lapisan secara bolak-balik', () => {
		const status = bacaLapisan('proyek,sekolah');
		expect(status).toEqual({ tema: false, proyek: true, sekolah: true });
		expect(tulisLapisan(status)).toBe('proyek,sekolah');
	});
});

describe('radiusTampilanM', () => {
	it('mengambil setengah tinggi tampilan dan membatasinya 10 km', () => {
		expect(radiusTampilanM(-6.19, -6.21)).toBe(1113);
		expect(radiusTampilanM(-5.9, -6.5)).toBe(10_000);
	});
});
