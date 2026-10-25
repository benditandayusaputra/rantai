import { describe, expect, it } from 'vitest';
import { tujuanAktif } from './navigasi';

describe('tujuanAktif', () => {
	it.each([
		['/id', 'beranda'],
		['/id/', 'beranda'],
		['/id/peta', 'peta'],
		['/id/dasbor', 'dasbor'],
		['/id/tanya', 'tanya'],
		['/id/tentang-data', 'tentang'],
		['/id/lapor', 'lapor'],
		['/id/pantau', 'pantau'],
		['/id/lapor/terkirim', 'lapor'],
		['/id/laporan/51', null],
		['/id/kecamatan/31.72.01', null]
	])('%s menandai %s', (path, harapan) => {
		expect(tujuanAktif(path, 'id')).toBe(harapan);
	});
});
