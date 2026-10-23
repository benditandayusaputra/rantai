import { describe, expect, it, vi } from 'vitest';
import { ambilData, GalatAPI, tangkap } from './klien';

vi.mock('$env/dynamic/public', () => ({ env: { PUBLIC_API_URL: 'http://api.uji/api/v1' } }));

function fetchTiruan(status: number, badan: unknown) {
	return vi.fn(async () => new Response(JSON.stringify(badan), { status }));
}

describe('ambilData', () => {
	it('mengembalikan isi data dan menyusun kueri tanpa nilai kosong', async () => {
		const fetchFn = fetchTiruan(200, { code: 200, message: 'ok', data: { a: 1 } });
		await expect(
			ambilData(fetchFn, '/wilayah/cari', { q: 'gading', x: undefined })
		).resolves.toEqual({ a: 1 });
		expect(fetchFn).toHaveBeenCalledWith('http://api.uji/api/v1/wilayah/cari?q=gading');
	});

	it('melempar GalatAPI dengan kode dan pesan dari amplop', async () => {
		const fetchFn = fetchTiruan(404, { code: 404, message: 'paket tidak ditemukan', data: null });
		await expect(ambilData(fetchFn, '/paket/x')).rejects.toEqual(
			new GalatAPI(404, 'paket tidak ditemukan')
		);
	});
});

describe('tangkap', () => {
	it('mengubah penolakan menjadi hasil gagal tanpa melempar', async () => {
		await expect(tangkap(Promise.reject(new GalatAPI(404, 'x')))).resolves.toEqual({
			ok: false,
			galat: new GalatAPI(404, 'x')
		});
		await expect(tangkap(Promise.resolve(1))).resolves.toEqual({ ok: true, data: 1 });
	});
});
