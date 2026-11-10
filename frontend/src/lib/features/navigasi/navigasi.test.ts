import { describe, expect, it } from 'vitest';
import { kodeWilayahDariURL, tautanDenganWilayah, tujuanAktif, tujuanPonsel } from './navigasi';

describe('tujuanAktif', () => {
	it.each([
		['/id', 'beranda'],
		['/id/', 'beranda'],
		['/id/peta', 'wilayah'],
		['/id/dasbor', 'wilayah'],
		['/id/tanya', 'tanya'],
		['/id/tentang-data', 'tentang'],
		['/id/lapor', 'lapor'],
		['/id/belajar', 'belajar'],
		['/id/pantau', null],
		['/id/lapor/terkirim', 'lapor'],
		['/id/laporan/51', null],
		['/id/kecamatan/31.72.01', 'wilayah'],
		['/id/kota/31.72', 'wilayah']
	])('%s menandai %s', (path, harapan) => {
		expect(tujuanAktif(path, 'id')).toBe(harapan);
	});
});

describe('konteks wilayah saat berpindah halaman', () => {
	it.each([
		['/id/kecamatan/31.72.06', '31.72.06'],
		['/id/kota/31.72', '31.72'],
		['/id/kecamatan/31.72.06/pendidikan', '31.72.06'],
		['/en/tanya?w=31.72.06', '31.72.06'],
		['/id/cek-sekolah?w=31.72.06.1001', '31.72.06'],
		['/id/tanya?w=di-luar-jakarta', null],
		['/id', null]
	])('%s membawa konteks %s', (path, kode) => {
		expect(kodeWilayahDariURL(new URL(path, 'http://localhost'))).toBe(kode);
	});

	it('membawa kode kecamatan ke Tanya dan tidak menambahkannya ke halaman umum', () => {
		expect(tautanDenganWilayah('/wilayah', 'id', '31.72.06')).toBe('/id/kecamatan/31.72.06');
		expect(tautanDenganWilayah('/wilayah', 'id', null)).toBe('/id/peta');
		expect(tautanDenganWilayah('/tanya', 'id', '31.72.06')).toBe('/id/tanya?w=31.72.06');
		expect(tautanDenganWilayah('/tanya', 'en', null)).toBe('/en/tanya');
		expect(tautanDenganWilayah('/tentang-data', 'id', '31.72.06')).toBe('/id/tentang-data');
	});

	it('menampilkan lima tujuan utama di navigasi ponsel', () => {
		expect(tujuanPonsel.map((tujuan) => tujuan.kunci)).toEqual([
			'beranda',
			'wilayah',
			'tanya',
			'lapor',
			'tentang'
		]);
	});
});

it('mempertahankan konteks kota tanpa mengirimkannya sebagai kecamatan ke Tanya', () => {
	expect(tautanDenganWilayah('/wilayah', 'id', '31.72')).toBe('/id/kota/31.72');
	expect(tautanDenganWilayah('/tanya', 'id', '31.72')).toBe('/id/tanya');
});
