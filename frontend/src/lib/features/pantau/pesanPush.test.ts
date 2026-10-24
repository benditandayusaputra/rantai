import { describe, expect, it } from 'vitest';
import { bacaPesanPush, kunciBase64KeBita, tujuanKlikPush } from './pesanPush';

describe('pesanPush', () => {
	it('membaca pesan sah dan menolak isi lain', () => {
		expect(bacaPesanPush('{"judul":"Paket baru","isi":"Rehab","tautan":"/paket/x"}')?.judul).toBe(
			'Paket baru'
		);
		expect(bacaPesanPush('{"judul":1}')).toBeNull();
		expect(bacaPesanPush('bukan json')).toBeNull();
	});

	it('hanya membuka path lokal berawalan bahasa', () => {
		expect(tujuanKlikPush('/paket/x')).toBe('/id/paket/x');
		expect(tujuanKlikPush('https://contoh.test')).toBe('/id/pantau');
		expect(tujuanKlikPush('//contoh.test')).toBe('/id/pantau');
	});

	it('mengubah kunci base64url VAPID menjadi bita', () => {
		expect(Array.from(kunciBase64KeBita('AQID_w'))).toEqual([1, 2, 3, 255]);
	});
});
