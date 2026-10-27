import type { KunciPesan } from '$lib/i18n';

export interface LangkahTur {
	target: string | null;
	judul: KunciPesan;
	isi: KunciPesan;
}

export interface Kotak {
	x: number;
	y: number;
	lebar: number;
	tinggi: number;
}

export interface Ukuran {
	lebar: number;
	tinggi: number;
}

export const langkahTur: LangkahTur[] = [
	{ target: null, judul: 'tour.welcome.title', isi: 'tour.welcome.body' },
	{ target: 'cari', judul: 'tour.search.title', isi: 'tour.search.body' }
];

export const kunciTurSelesai = 'rantai-tour-done';
export const jarakKartu = 12;
export const tepiLayar = 16;

export function isTurSelesai(): boolean {
	try {
		return localStorage.getItem(kunciTurSelesai) !== null;
	} catch {
		return true;
	}
}

export function tandaiTurSelesai() {
	try {
		localStorage.setItem(kunciTurSelesai, '1');
	} catch {
		return;
	}
}

export function langkahTersedia(
	daftar: LangkahTur[],
	isTargetAda: (target: string) => boolean
): LangkahTur[] {
	return daftar.filter((langkah) => langkah.target === null || isTargetAda(langkah.target));
}

export function elemenTarget(target: string): HTMLElement | null {
	const kandidat = document.querySelectorAll<HTMLElement>(`[data-tur="${target}"]`);
	return (
		[...kandidat].find((el) => {
			const kotak = el.getBoundingClientRect();
			return kotak.width > 0 && kotak.height > 0;
		}) ?? null
	);
}

function batasi(nilai: number, bawah: number, atas: number): number {
	return Math.min(Math.max(nilai, bawah), atas);
}

export function posisiKartu(target: Kotak, kartu: Ukuran, layar: Ukuran): { x: number; y: number } {
	const yBawah = target.y + target.tinggi + jarakKartu;
	const isMuatDiBawah = yBawah + kartu.tinggi <= layar.tinggi - tepiLayar;
	const yAtas = Math.max(tepiLayar, target.y - jarakKartu - kartu.tinggi);
	const xTengah = target.x + target.lebar / 2 - kartu.lebar / 2;
	return {
		x: batasi(xTengah, tepiLayar, layar.lebar - tepiLayar - kartu.lebar),
		y: isMuatDiBawah ? yBawah : yAtas
	};
}
