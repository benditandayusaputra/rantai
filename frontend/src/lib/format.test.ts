import { describe, expect, it } from 'vitest';
import { formatAngka, formatJam, formatJarak, formatRupiah, formatTanggal } from './format';

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

describe('formatJarak', () => {
	it('memakai meter di bawah 1 km dan km satu desimal di atasnya', () => {
		expect(formatJarak(480.4, 'id')).toBe('480 m');
		expect(formatJarak(1234, 'id')).toBe('1,2 km');
	});
});

describe('formatRupiah', () => {
	it.each([
		[845_000, 'Rp845.000'],
		[845_000_000, 'Rp845 juta'],
		[1_200_000_000, 'Rp1,2 miliar'],
		[1_400_000_000_000, 'Rp1,4 triliun']
	])('%d menjadi %s', (n, harapan) => {
		expect(formatRupiah(n, 'id')).toBe(harapan);
	});

	it('menulis nilai lengkap untuk panel sumber', () => {
		expect(formatRupiah(2_408_700_000, 'id', 'penuh')).toBe('Rp2.408.700.000');
	});
});
