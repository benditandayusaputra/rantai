import { bahasaBawaan } from '$lib/i18n';

export interface PesanPush {
	judul: string;
	isi: string;
	tautan: string;
}

function isPesanPush(nilai: unknown): nilai is PesanPush {
	if (typeof nilai !== 'object' || nilai === null) return false;
	const p = nilai as Record<string, unknown>;
	return typeof p.judul === 'string' && typeof p.isi === 'string' && typeof p.tautan === 'string';
}

export function bacaPesanPush(teks: string): PesanPush | null {
	try {
		const nilai: unknown = JSON.parse(teks);
		return isPesanPush(nilai) ? nilai : null;
	} catch {
		return null;
	}
}

export function tujuanKlikPush(tautan: string): string {
	const isPathLokal = tautan.startsWith('/') && !tautan.startsWith('//');
	return isPathLokal ? `/${bahasaBawaan}${tautan}` : `/${bahasaBawaan}/pantau`;
}

export function kunciBase64KeBita(kunci: string): Uint8Array<ArrayBuffer> {
	const isian = '='.repeat((4 - (kunci.length % 4)) % 4);
	const biner = atob((kunci + isian).replace(/-/g, '+').replace(/_/g, '/'));
	const hasil = new Uint8Array(new ArrayBuffer(biner.length));
	for (let i = 0; i < biner.length; i += 1) hasil[i] = biner.charCodeAt(i);
	return hasil;
}
