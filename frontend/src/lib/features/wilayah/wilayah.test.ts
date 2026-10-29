import { describe, expect, it, vi } from 'vitest';
import type { WilayahDitemukan } from '$lib/api/tipe';
import { cariLokasiku } from './lokasi';
import { bagiSorotan, induk, saringHasilPencarian, tujuanWilayah } from './wilayah';

vi.mock('$lib/api/klien', () => ({
	ambilWilayahDiTitik: vi.fn(async () => kelurahan)
}));

const provinsi = { kode: '31', nama: 'DKI Jakarta', tingkat: 'provinsi' } as const;
const kota = { kode: '31.72', nama: 'Jakarta Utara', tingkat: 'kota' } as const;
const kecamatan: WilayahDitemukan = {
	kode: '31.72.06',
	nama: 'Kelapa Gading',
	tingkat: 'kecamatan',
	hierarki: [provinsi, kota]
};
const kelurahan: WilayahDitemukan = {
	kode: '31.72.06.1001',
	nama: 'Kelapa Gading Barat',
	tingkat: 'kelurahan',
	hierarki: [provinsi, kota, { kode: '31.72.06', nama: 'Kelapa Gading', tingkat: 'kecamatan' }]
};

function geolokasiTiruan(hasil: { akurasi?: number; kodeGalat?: number }): Geolocation {
	return {
		getCurrentPosition: (berhasil, gagal) =>
			hasil.kodeGalat
				? gagal?.({ code: hasil.kodeGalat } as GeolocationPositionError)
				: berhasil({
						coords: { latitude: -6.16, longitude: 106.9, accuracy: hasil.akurasi }
					} as GeolocationPosition)
	} as Geolocation;
}

describe('pencarian wilayah', () => {
	it('menyertakan kota dan mendahulukan kelurahan', () => {
		const hasil = saringHasilPencarian([{ ...kota, hierarki: [provinsi] }, kecamatan, kelurahan]);
		expect(hasil.map((w) => w.kode)).toEqual(['31.72.06.1001', '31.72.06', '31.72']);
	});

	it('mengarahkan kelurahan ke halaman kecamatan dengan kelurahan disorot', () => {
		expect(tujuanWilayah(kelurahan)).toBe('/id/kecamatan/31.72.06?kelurahan=31.72.06.1001');
		expect(tujuanWilayah(kecamatan)).toBe('/id/kecamatan/31.72.06');
		expect(tujuanWilayah({ ...kota, hierarki: [provinsi] })).toBe('/id/kota/31.72');
	});

	it('menulis induk dari yang terdekat tanpa provinsi', () => {
		expect(induk(kelurahan)).toBe('Kelapa Gading, Jakarta Utara');
	});

	it('memisahkan bagian nama yang cocok', () => {
		expect(bagiSorotan('Kelapa Gading Barat', 'gading')).toEqual(['Kelapa ', 'Gading', ' Barat']);
		expect(bagiSorotan('Cilincing', 'gading')).toEqual(['Cilincing', '', '']);
	});
});

describe('cariLokasiku', () => {
	it('melaporkan izin ditolak', async () => {
		expect(await cariLokasiku(geolokasiTiruan({ kodeGalat: 1 }))).toEqual({ status: 'denied' });
	});

	it('melaporkan gagal saat waktu habis atau geolokasi tidak tersedia', async () => {
		expect(await cariLokasiku(geolokasiTiruan({ kodeGalat: 3 }))).toEqual({ status: 'failed' });
		expect(await cariLokasiku(undefined)).toEqual({ status: 'failed' });
	});

	it('menolak akurasi kasar tanpa menebak wilayah', async () => {
		expect(await cariLokasiku(geolokasiTiruan({ akurasi: 1200 }))).toEqual({
			status: 'coarse',
			akurasiM: 1200
		});
	});

	it('mengembalikan kelurahan di titik bila akurat', async () => {
		expect(await cariLokasiku(geolokasiTiruan({ akurasi: 20 }))).toEqual({
			status: 'found',
			wilayah: kelurahan,
			akurasiM: 20
		});
	});
});
