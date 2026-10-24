import { describe, expect, it, vi } from 'vitest';
import type { Pengguna } from '$lib/api/tipe';
import { inisial, keluar, pulihkanSesi, sesiAktif } from './sesi.svelte';

vi.mock('$env/dynamic/public', () => ({ env: { PUBLIC_API_URL: 'http://api.uji/api/v1' } }));

const pengguna: Pengguna = {
	id: 1,
	email: 'warga@contoh.id',
	nama_tampilan: '',
	peran: 'warga',
	instansi_id: null,
	anggota_id: null
};

function amplop(status: number, data: unknown) {
	return new Response(JSON.stringify({ code: status, message: 'x', data }), { status });
}

describe('sesi', () => {
	it('memulihkan sesi lewat refresh dengan cookie sekali saja, lalu keluar menghapusnya', async () => {
		const fetchFn = vi.fn(async () => amplop(200, { token_akses: 'akses', pengguna }));

		await pulihkanSesi(fetchFn);
		await pulihkanSesi(fetchFn);

		expect(sesiAktif()?.token_akses).toBe('akses');
		expect(fetchFn).toHaveBeenCalledTimes(1);
		expect(fetchFn).toHaveBeenCalledWith(
			'http://api.uji/api/v1/auth/refresh',
			expect.objectContaining({ credentials: 'include', method: 'POST' })
		);

		await keluar(vi.fn(async () => amplop(200, null)));

		expect(sesiAktif()).toBeNull();
	});
});

describe('inisial', () => {
	it.each([
		['Warga Contoh Simulasi', 'WC'],
		['  budi ', 'B'],
		['', 'W']
	])('nama %s menjadi %s', (nama, harap) => {
		expect(inisial({ ...pengguna, nama_tampilan: nama })).toBe(harap);
	});
});
