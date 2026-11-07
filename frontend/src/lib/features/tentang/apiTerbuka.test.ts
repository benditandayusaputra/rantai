import { describe, expect, it } from 'vitest';
import { alamatLengkap, contohPath, titikAkhirTerbuka } from './apiTerbuka';

describe('API terbuka', () => {
	it('menyusun alamat absolut untuk setiap titik akhir', () => {
		for (const titik of titikAkhirTerbuka) {
			expect(alamatLengkap(titik, 'https://rantai.example')).toMatch(/^https?:\/\//);
		}
	});

	it('menampilkan path beserta kueri contoh', () => {
		const jaringan = titikAkhirTerbuka.find((t) => t.path === '/jaringan');
		expect(jaringan && contohPath(jaringan)).toMatch(/\/jaringan\?wilayah=31\.72$/);
	});
});
