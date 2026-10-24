import { describe, expect, it, vi } from 'vitest';
import { bacaPilihanPermohonan, muatKeterbukaan, tautanPermohonan } from './keterbukaan';

vi.mock('$env/dynamic/public', () => ({ env: { PUBLIC_API_URL: 'http://api.uji/api/v1' } }));

const daftar = {
	ppid: { nama: 'PPID DPRD', alamat: 'Jl. Contoh', url: 'https://ppid.uji' },
	data: []
};

function fetchTiruan() {
	return vi.fn<typeof fetch>(async (url) => {
		const isPermohonan = String(url).endsWith('/permohonan');
		const data = isPermohonan ? { teks: 'Surat', pdf_url: '/api/v1/x.pdf' } : daftar;
		return new Response(JSON.stringify({ code: 200, message: 'ok', data }), { status: 200 });
	});
}

describe('tautanPermohonan', () => {
	it('mempertahankan data terpilih saat bahasa diganti dan menuju bagian surat', () => {
		const url = new URL('http://rantai.uji/id/data-belum-dibuka?minta=kehadiran_per_anggota');
		expect(tautanPermohonan(url, { bahasa: 'en' })).toBe(
			'/id/data-belum-dibuka?minta=kehadiran_per_anggota&bahasa=en#surat'
		);
	});
});

describe('bacaPilihanPermohonan', () => {
	it('memakai bahasa Indonesia untuk bahasa yang tidak dikenal', () => {
		expect(bacaPilihanPermohonan(new URLSearchParams('minta=a&bahasa=fr'))).toEqual({
			kode: 'a',
			bahasa: 'id'
		});
	});
});

describe('muatKeterbukaan', () => {
	it('tidak menyusun surat sebelum data dipilih', async () => {
		const fetchFn = fetchTiruan();
		const hasil = await muatKeterbukaan(fetchFn, { kode: '', bahasa: 'id' });
		expect(hasil.dokumen).toBeNull();
		expect(fetchFn).toHaveBeenCalledTimes(1);
	});

	it('menyusun surat untuk data terpilih dengan bahasa surat', async () => {
		const fetchFn = fetchTiruan();
		const hasil = await muatKeterbukaan(fetchFn, { kode: 'pokir_per_anggota', bahasa: 'en' });
		const [url, init] = fetchFn.mock.calls.find(([u]) => String(u).endsWith('/permohonan')) ?? [];
		expect(String(url)).toBe(
			'http://api.uji/api/v1/data-belum-dibuka/pokir_per_anggota/permohonan'
		);
		expect(JSON.parse(String(init?.body))).toEqual({ bahasa: 'en' });
		expect(hasil.dokumen?.teks).toBe('Surat');
	});
});
