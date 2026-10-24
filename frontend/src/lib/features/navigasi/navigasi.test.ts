import { describe, expect, it } from 'vitest';
import { tujuanAktif } from './navigasi';

describe('tujuanAktif', () => {
	it.each([
		['/id', 'beranda'],
		['/id/', 'beranda'],
		['/id/peta', 'peta'],
		['/id/tanya', 'tanya'],
		['/id/tentang-data', 'tentang'],
		['/id/kecamatan/31.72.01', null]
	])('%s menandai %s', (path, harapan) => {
		expect(tujuanAktif(path, 'id')).toBe(harapan);
	});
});
