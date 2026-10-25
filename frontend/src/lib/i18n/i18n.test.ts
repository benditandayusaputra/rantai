import { afterEach, describe, expect, it } from 'vitest';
import en from './en.json';
import id from './id.json';
import { aturBahasa, bahasaBawaan, gantiBahasaPath, t } from '.';

afterEach(() => aturBahasa(bahasaBawaan));

describe('t', () => {
	it('mengisi parameter bernama', () => {
		expect(t('search.loading', { q: 'gading' })).toBe('Mencari "gading"...');
	});

	it('membiarkan parameter yang tidak diberikan', () => {
		expect(t('chain.summary', { ok: 5 })).toBe('5 dari {total} mata rantai tersambung');
	});

	it('mengembalikan pesan tanpa parameter apa adanya', () => {
		expect(t('nav.main')).toBe('Navigasi utama');
	});
});

describe('katalog Inggris', () => {
	const polaParameter = /\{(\w+)\}/g;
	const parameterDari = (teks: string) => [...teks.matchAll(polaParameter)].map((m) => m[1]).sort();

	it('memuat kunci yang sama dengan parameter yang sama', () => {
		const katalogEn: Record<string, string> = en;
		for (const [kunci, teks] of Object.entries(id)) {
			expect(parameterDari(katalogEn[kunci] ?? ''), kunci).toEqual(parameterDari(teks));
		}
		expect(Object.keys(en)).toHaveLength(Object.keys(id).length);
	});

	it('dipakai t setelah bahasa aktif diganti', () => {
		aturBahasa('en');
		expect(t('nav.main')).toBe(en['nav.main']);
	});
});

describe('gantiBahasaPath', () => {
	it.each([
		['/id', 'en', '/en'],
		['/id/kecamatan/31.72.06', 'en', '/en/kecamatan/31.72.06'],
		['/en/peta', 'id', '/id/peta'],
		['/', 'en', '/en']
	] as const)('%s ke %s menjadi %s', (path, bahasa, harapan) => {
		expect(gantiBahasaPath(path, bahasa)).toBe(harapan);
	});
});
