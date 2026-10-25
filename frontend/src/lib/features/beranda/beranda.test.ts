import { describe, expect, it } from 'vitest';
import type { Dapil, TemaKecamatan, Wakil } from '$lib/api/tipe';
import {
	bacaSorotan,
	batasWakilBeranda,
	indikatorPendidikan,
	persenTerhadap,
	sorotanBawaan,
	tautanSorotan,
	wakilPengawas
} from './beranda';

describe('wilayah sorotan beranda', () => {
	it('membaca kode dari URL dan kembali ke wilayah pertama bila tidak dikenal', () => {
		expect(bacaSorotan('31.72.04').nama).toBe('Cilincing');
		expect(bacaSorotan(null)).toBe(sorotanBawaan);
		expect(bacaSorotan('99.99.99')).toBe(sorotanBawaan);
	});

	it('menghapus parameter untuk wilayah bawaan dan memasangnya untuk yang lain', () => {
		const url = new URL('http://x/id?wilayah=31.72.04');
		expect(tautanSorotan(url, sorotanBawaan.kode)).toBe('/id');
		expect(tautanSorotan(url, '31.72.06')).toBe('/id?wilayah=31.72.06');
	});
});

describe('indikatorPendidikan', () => {
	const indikator = { kelas: 'kurang' } as TemaKecamatan['indikator'];

	it('mengambil indikator tema pendidikan saja', () => {
		const tema = [
			{ kode: 'kesehatan', indikator: null },
			{ kode: 'pendidikan', indikator }
		] as TemaKecamatan[];
		expect(indikatorPendidikan(tema)).toBe(indikator);
	});

	it('mengembalikan null bila tema pendidikan belum dihitung', () => {
		expect(indikatorPendidikan([])).toBeNull();
	});
});

describe('persenTerhadap', () => {
	it('menghitung porsi batang dan membatasinya pada 0 sampai 100', () => {
		expect(persenTerhadap(4140, 6340)).toBeCloseTo(65.3, 1);
		expect(persenTerhadap(8000, 6340)).toBe(100);
		expect(persenTerhadap(5, 0)).toBe(0);
	});
});

describe('wakilPengawas', () => {
	const wakil = (id: string, isPengawas: boolean) =>
		({ id, duduk_di_komisi_pengawas: isPengawas }) as Wakil;

	it('memilih anggota komisi pengawas dengan batas tampilan beranda', () => {
		const anggota = [
			wakil('a', true),
			wakil('b', false),
			wakil('c', true),
			wakil('d', true),
			wakil('e', true)
		];
		const hasil = wakilPengawas({ anggota } as Dapil);
		expect(hasil.map((w) => w.id)).toEqual(['a', 'c', 'd']);
		expect(hasil).toHaveLength(batasWakilBeranda);
	});
});
