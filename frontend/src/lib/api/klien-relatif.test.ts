import { describe, expect, it, vi } from 'vitest';
import { susunURL, urlServer } from './klien';

vi.mock('$env/dynamic/public', () => ({ env: { PUBLIC_API_URL: '/api/v1' } }));

describe('alamat API relatif', () => {
	it('menghasilkan path dan kueri tanpa origin agar mengikuti domain halaman', () => {
		expect(susunURL('/wilayah/cari', { q: 'gading', x: undefined })).toBe(
			'/api/v1/wilayah/cari?q=gading'
		);
		expect(urlServer('/api/v1/aspirasi/pdf?bahasa=id')).toBe('/api/v1/aspirasi/pdf?bahasa=id');
	});
});
