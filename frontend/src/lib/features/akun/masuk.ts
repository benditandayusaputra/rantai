import type { GalatAPI } from '$lib/api/klien';

export type KeadaanKode = 'siap' | 'salah' | 'kedaluwarsa' | 'terkunci' | 'galat';

export const panjangKode = 6;
export const detikKirimUlang = 60;
export const detikTerkunci = 15 * 60;

const keadaanPerStatus: Record<number, KeadaanKode> = {
	401: 'salah',
	410: 'kedaluwarsa',
	429: 'terkunci'
};

export function keadaanDariGalat(galat: GalatAPI): KeadaanKode {
	return keadaanPerStatus[galat.kode] ?? 'galat';
}

export function formatHitungMundur(detik: number): string {
	const sisa = Math.max(0, Math.ceil(detik));
	return `${Math.floor(sisa / 60)}:${String(sisa % 60).padStart(2, '0')}`;
}

export function bersihkanKode(masukan: string): string {
	return masukan.replace(/\D/g, '').slice(0, panjangKode);
}

export function tujuanSetelahMasuk(lanjut: string | null, bahasa: string): string {
	const isPathLokal =
		lanjut?.startsWith('/') && !lanjut.startsWith('//') && !lanjut.startsWith('/\\');
	return isPathLokal && lanjut ? lanjut : `/${bahasa}`;
}
