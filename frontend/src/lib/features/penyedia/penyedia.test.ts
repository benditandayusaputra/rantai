import { describe, expect, it } from 'vitest';
import type { TahunPenyedia } from '$lib/api/tipe';
import { inisialPenyedia, tahunTerbaru, tinggiBatang } from './penyedia';

function tahun(t: number, nilai: number): TahunPenyedia {
	const angka = { nilai, satuan: 'rupiah', label: 'resmi' as const, provenance_id: 'prv' };
	return {
		tahun: t,
		sumber: 'lpse',
		jumlah_paket: { ...angka, nilai: 1, satuan: 'paket' },
		total_nilai: angka
	};
}

describe('penyedia', () => {
	it('memilih tahun terbaru dan menormalkan tinggi batang terhadap nilai terbesar', () => {
		const daftar = [tahun(2024, 200), tahun(2026, 50), tahun(2025, 0)];
		expect(tahunTerbaru(daftar)?.tahun).toBe(2026);
		expect(tahunTerbaru([])).toBeNull();
		expect([...tinggiBatang(daftar).values()]).toEqual([1, 0.25, 0]);
	});

	it('membuat inisial tanpa bentuk badan usaha', () => {
		expect(inisialPenyedia('CV. Karya Mandiri Jaya')).toBe('KM');
		expect(inisialPenyedia('PT Bangun')).toBe('B');
	});
});
