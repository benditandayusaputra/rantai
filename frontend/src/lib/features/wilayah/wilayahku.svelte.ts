import type { WilayahDitemukan } from '$lib/api/tipe';
import { sesiAktif } from '$lib/features/akun/sesi.svelte';
import { kodeKecamatan } from './wilayah';

export const kunciWilayahku = 'rantai-wilayahku';
const polaKodeKecamatan = /^31\.\d{2}\.\d{2}$/;

const keadaan = $state<{ kode: string | null }>({ kode: null });

export function wilayahku(): string | null {
	return keadaan.kode;
}

export function pulihkanWilayahku(): void {
	try {
		const kode = localStorage.getItem(kunciWilayahku);
		keadaan.kode = kode && polaKodeKecamatan.test(kode) ? kode : null;
	} catch {
		keadaan.kode = null;
	}
}

export function simpanWilayahku(kode: string): void {
	if (!polaKodeKecamatan.test(kode)) return;
	keadaan.kode = kode;
	try {
		localStorage.setItem(kunciWilayahku, kode);
	} catch {
		return;
	}
}

export function ingatWilayah(wilayah: WilayahDitemukan): void {
	if (wilayah.tingkat !== 'kota') simpanWilayahku(kodeKecamatan(wilayah));
}

export function berandaPengguna(bahasa: string): string {
	return sesiAktif() && keadaan.kode ? `/${bahasa}/kecamatan/${keadaan.kode}` : `/${bahasa}`;
}

export function isBerandaWilayah(): boolean {
	return Boolean(sesiAktif() && keadaan.kode);
}
