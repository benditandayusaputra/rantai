import { describe, expect, it } from 'vitest';
import { formatAngka, formatJam, formatTanggal } from './format';

describe('format', () => {
	it('memakai pemisah ribuan titik dan desimal koma', () => {
		expect(formatAngka(1284, 'id')).toBe('1.284');
		expect(formatAngka(2.14, 'id', 1)).toBe('2,1');
	});

	it('menulis tanggal pendek dan panjang dalam zona Jakarta', () => {
		const d = new Date('2026-10-13T20:00:00Z');
		expect(formatTanggal(d, 'id')).toBe('14 Okt 2026');
		expect(formatTanggal(d, 'id', 'long')).toBe('14 Oktober 2026');
	});

	it('menulis jam 24 jam dengan pemisah titik', () => {
		expect(formatJam(new Date('2026-10-14T01:14:00Z'))).toBe('08.14');
	});
});
