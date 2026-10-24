import { describe, expect, it, vi } from 'vitest';
import { urlServer } from '$lib/api/klien';
import { bacaPilihan, muatAspirasi, tautanAspirasi } from './aspirasi';

vi.mock('$env/dynamic/public', () => ({ env: { PUBLIC_API_URL: 'http://api.uji/api/v1' } }));

const wilayah = { kode: '31.72.04', nama: 'Cilincing', tingkat: 'kecamatan', hierarki: [] };
const tema = {
	kode: 'pendidikan',
	nama: { id: 'Pendidikan', en: 'Education' },
	tanggung_jawab: []
};

function fetchTiruan() {
	return vi.fn(async (url: string | URL | Request, init?: RequestInit) => {
		const path = new URL(String(url)).pathname;
		const data =
			path === '/api/v1/aspirasi'
				? { teks: 'Isi dokumen', pdf_url: '/api/v1/aspirasi/pdf?masalah=kursi' }
				: path === '/api/v1/tema'
					? [tema]
					: path.startsWith('/api/v1/tema/')
						? tema
						: wilayah;
		expect(init?.method === 'POST').toBe(path === '/api/v1/aspirasi');
		return new Response(JSON.stringify({ code: 200, message: 'ok', data }), { status: 200 });
	});
}

describe('bacaPilihan', () => {
	it('memakai bahasa Indonesia untuk bahasa yang tidak dikenal', () => {
		expect(bacaPilihan(new URLSearchParams('w=31.72.04&masalah=kursi&bahasa=fr'))).toEqual({
			kode: '31.72.04',
			tema: '',
			masalah: 'kursi',
			bahasa: 'id'
		});
	});
});

describe('tautanAspirasi', () => {
	it('mengganti satu pilihan dan mempertahankan kueri lain', () => {
		const url = new URL('http://rantai.uji/id/aspirasi?w=31.72.04&masalah=kursi');
		expect(tautanAspirasi(url, { bahasa: 'en' })).toBe(
			'/id/aspirasi?w=31.72.04&masalah=kursi&bahasa=en'
		);
	});
});

describe('muatAspirasi', () => {
	it('tidak menyusun dokumen sebelum masalah dipilih dan memakai tema aktif pertama', async () => {
		const fetchFn = fetchTiruan();
		const hasil = await muatAspirasi(fetchFn, {
			kode: '31.72.04',
			tema: '',
			masalah: '',
			bahasa: 'id'
		});
		expect(hasil.tema.kode).toBe('pendidikan');
		expect(hasil.dokumen).toBeNull();
	});

	it('menyusun dokumen dengan kode kecamatan, tema, masalah, dan bahasa terpilih', async () => {
		const fetchFn = fetchTiruan();
		const hasil = await muatAspirasi(fetchFn, {
			kode: '31.72.04',
			tema: 'pendidikan',
			masalah: 'kursi',
			bahasa: 'en'
		});
		const badan = fetchFn.mock.calls.find(([, init]) => init?.method === 'POST')?.[1]?.body;
		expect(JSON.parse(String(badan))).toEqual({
			kode_kecamatan: '31.72.04',
			tema: 'pendidikan',
			masalah: 'kursi',
			bahasa: 'en'
		});
		expect(urlServer(hasil.dokumen?.pdf_url ?? '')).toBe(
			'http://api.uji/api/v1/aspirasi/pdf?masalah=kursi'
		);
	});
});
