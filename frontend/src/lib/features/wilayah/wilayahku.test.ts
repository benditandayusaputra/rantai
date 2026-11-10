import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { Sesi, WilayahDitemukan } from '$lib/api/tipe';
import { simpanSesi } from '$lib/features/akun/sesi.svelte';
import {
	berandaPengguna,
	ingatWilayah,
	isBerandaWilayah,
	kunciWilayahku,
	pulihkanWilayahku,
	simpanWilayahku,
	wilayahku
} from './wilayahku.svelte';

vi.mock('$env/dynamic/public', () => ({ env: { PUBLIC_API_URL: 'http://api.uji/api/v1' } }));

const kelurahan = {
	kode: '31.72.06.1003',
	tingkat: 'kelurahan',
	hierarki: [{ kode: '31.72.06', nama: 'Kelapa Gading', tingkat: 'kecamatan' }]
} as WilayahDitemukan;

beforeEach(() => {
	localStorage.clear();
	pulihkanWilayahku();
});

describe('wilayahku', () => {
	it('mengingat kecamatan dari kelurahan terpilih dan memulihkannya dari penyimpanan', () => {
		ingatWilayah(kelurahan);

		expect(wilayahku()).toBe('31.72.06');
		expect(localStorage.getItem(kunciWilayahku)).toBe('31.72.06');
		pulihkanWilayahku();
		expect(wilayahku()).toBe('31.72.06');
	});

	it('tidak mengingat kota dan menolak kode yang bukan kecamatan Jakarta', () => {
		ingatWilayah({ kode: '31.72', tingkat: 'kota', hierarki: [] } as unknown as WilayahDitemukan);
		simpanWilayahku('32.01.01');
		localStorage.setItem(kunciWilayahku, '<script>');
		pulihkanWilayahku();

		expect(wilayahku()).toBeNull();
	});

	it('menjadikan kecamatan tersimpan sebagai beranda hanya setelah masuk', () => {
		simpanWilayahku('31.75.07');
		expect(berandaPengguna('id')).toBe('/id');
		expect(isBerandaWilayah()).toBe(false);

		simpanSesi({ token_akses: 'akses' } as Sesi);

		expect(berandaPengguna('en')).toBe('/en/kecamatan/31.75.07');
		expect(isBerandaWilayah()).toBe(true);
	});
});
