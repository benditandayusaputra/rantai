import type { JenisSimpul, SimpulJaringan, SisiJaringan } from '$lib/api/tipe';

const jariMinimum: Record<JenisSimpul, number> = { instansi: 9, paket: 4, penyedia: 8 };
const jariTambahanMaks: Record<JenisSimpul, number> = { instansi: 13, paket: 8, penyedia: 12 };
const batasLabel = 26;

export function cariTetangga(id: string, sisi: SisiJaringan[]): Set<string> {
	const tetangga = new Set<string>();
	for (const s of sisi) {
		if (s.dari === id) tetangga.add(s.ke);
		if (s.ke === id) tetangga.add(s.dari);
	}
	return tetangga;
}

export function hitungJari(simpul: SimpulJaringan, nilaiMaks: Record<JenisSimpul, number>): number {
	const maks = nilaiMaks[simpul.jenis];
	const porsi = maks > 0 ? Math.sqrt(simpul.nilai.nilai / maks) : 0;
	return jariMinimum[simpul.jenis] + porsi * jariTambahanMaks[simpul.jenis];
}

export function nilaiMaksPerJenis(simpul: SimpulJaringan[]): Record<JenisSimpul, number> {
	const maks: Record<JenisSimpul, number> = { instansi: 0, paket: 0, penyedia: 0 };
	for (const s of simpul) maks[s.jenis] = Math.max(maks[s.jenis], s.nilai.nilai);
	return maks;
}

export function tautanSimpul(simpul: SimpulJaringan, bahasa: string): string | null {
	if (simpul.jenis === 'paket') return `/${bahasa}/paket/${encodeURIComponent(simpul.kunci)}`;
	if (simpul.jenis === 'penyedia') return `/${bahasa}/penyedia/${simpul.kunci}`;
	return null;
}

export function ringkasLabel(label: string): string {
	return label.length > batasLabel ? `${label.slice(0, batasLabel - 1)}…` : label;
}

const lebarHurufLabel = 7;
const tinggiLabel = 14;
const naikPusatLabel = 5;

export interface Kotak {
	id: string;
	x: number;
	y: number;
	lebar: number;
	tinggi: number;
}

export function kotakLabel(id: string, x: number, yGaris: number, teks: string): Kotak {
	return {
		id,
		x,
		y: yGaris - naikPusatLabel,
		lebar: teks.length * lebarHurufLabel,
		tinggi: tinggiLabel
	};
}

export function saringLabelBertumpuk(
	penghalang: Kotak[],
	labelUrutPrioritas: Kotak[]
): Set<string> {
	const terpasang = [...penghalang];
	const tampil = new Set<string>();
	for (const label of labelUrutPrioritas) {
		if (terpasang.some((lain) => lain.id !== label.id && isBertumpuk(label, lain))) continue;
		terpasang.push(label);
		tampil.add(label.id);
	}
	return tampil;
}

function isBertumpuk(a: Kotak, b: Kotak): boolean {
	return (
		Math.abs(a.x - b.x) * 2 < a.lebar + b.lebar && Math.abs(a.y - b.y) * 2 < a.tinggi + b.tinggi
	);
}

export interface BarisRelasi {
	paket: SimpulJaringan;
	instansi: SimpulJaringan | null;
	penyedia: SimpulJaringan | null;
}

export function susunBarisRelasi(simpul: SimpulJaringan[], sisi: SisiJaringan[]): BarisRelasi[] {
	const menurutId = new Map(simpul.map((s) => [s.id, s]));
	const ujung = (id: string, jenis: JenisSimpul) =>
		[...cariTetangga(id, sisi)].map((t) => menurutId.get(t)).find((s) => s?.jenis === jenis) ??
		null;
	return simpul
		.filter((s) => s.jenis === 'paket')
		.map((paket) => ({
			paket,
			instansi: ujung(paket.id, 'instansi'),
			penyedia: ujung(paket.id, 'penyedia')
		}));
}
