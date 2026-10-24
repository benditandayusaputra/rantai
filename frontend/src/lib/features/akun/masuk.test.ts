import { describe, expect, it } from 'vitest';
import { GalatAPI } from '$lib/api/klien';
import { bersihkanKode, formatHitungMundur, keadaanDariGalat, tujuanSetelahMasuk } from './masuk';

describe('keadaanDariGalat', () => {
	it.each([
		[401, 'salah'],
		[410, 'kedaluwarsa'],
		[429, 'terkunci'],
		[500, 'galat'],
		[0, 'galat']
	])('status %i menjadi %s', (status, harap) => {
		expect(keadaanDariGalat(new GalatAPI(status, 'x'))).toBe(harap);
	});
});

describe('formatHitungMundur', () => {
	it.each([
		[900, '15:00'],
		[61, '1:01'],
		[0.2, '0:01'],
		[-3, '0:00']
	])('%f detik menjadi %s', (detik, harap) => {
		expect(formatHitungMundur(detik)).toBe(harap);
	});
});

describe('bersihkanKode', () => {
	it('membuang selain angka dan memotong ke enam digit', () => {
		expect(bersihkanKode(' 12-34 5678')).toBe('123456');
	});
});

describe('tujuanSetelahMasuk', () => {
	it.each([
		['/id/peta', '/id/peta'],
		[null, '/id'],
		['https://jahat.test', '/id'],
		['//jahat.test', '/id'],
		['/\\jahat.test', '/id']
	])('lanjut %s menuju %s', (lanjut, harap) => {
		expect(tujuanSetelahMasuk(lanjut, 'id')).toBe(harap);
	});
});
