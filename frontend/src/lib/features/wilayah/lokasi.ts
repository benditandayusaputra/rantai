import { ambilWilayahDiTitik } from '$lib/api/klien';
import type { WilayahDitemukan } from '$lib/api/tipe';

export const ambangAkurasiKasarM = 100;
const batasWaktuLokasiMs = 10_000;
const kodeIzinDitolak = 1;

export type StatusLokasi = 'idle' | 'locating' | 'found' | 'coarse' | 'denied' | 'failed';

export type HasilLokasi =
	| { status: 'found'; wilayah: WilayahDitemukan; akurasiM: number }
	| { status: 'coarse'; akurasiM: number }
	| { status: 'denied' }
	| { status: 'failed' };

export function mintaPosisi(geolokasi: Geolocation): Promise<GeolocationPosition> {
	return new Promise((selesai, gagal) =>
		geolokasi.getCurrentPosition(selesai, gagal, {
			enableHighAccuracy: true,
			timeout: batasWaktuLokasiMs,
			maximumAge: 0
		})
	);
}

export function isIzinDitolak(galat: unknown): boolean {
	return (
		typeof galat === 'object' && galat !== null && 'code' in galat && galat.code === kodeIzinDitolak
	);
}

async function petakanPosisi(posisi: GeolocationPosition): Promise<HasilLokasi> {
	const { latitude, longitude, accuracy } = posisi.coords;
	if (accuracy > ambangAkurasiKasarM) return { status: 'coarse', akurasiM: accuracy };
	const wilayah = await ambilWilayahDiTitik(fetch, latitude, longitude);
	return { status: 'found', wilayah, akurasiM: accuracy };
}

export async function cariLokasiku(geolokasi: Geolocation | undefined): Promise<HasilLokasi> {
	if (!geolokasi) return { status: 'failed' };
	try {
		return await petakanPosisi(await mintaPosisi(geolokasi));
	} catch (galat) {
		return isIzinDitolak(galat) ? { status: 'denied' } : { status: 'failed' };
	}
}
