import type { SisiJaringan } from '$lib/api/tipe';

export interface Posisi {
	x: number;
	y: number;
}

export interface Bidang {
	lebar: number;
	tinggi: number;
	tepi: number;
}

interface Titik extends Posisi {
	dx: number;
	dy: number;
}

const jumlahIterasi = 320;
const sudutEmas = Math.PI * (3 - Math.sqrt(5));
const tarikanPusat = 1.5;
const skalaJarakIdeal = 0.5;
const jarakMinimum = 0.01;

function titikAwal(jumlah: number, bidang: Bidang): Titik[] {
	const jariMaks = Math.min(bidang.lebar, bidang.tinggi) / 2 - bidang.tepi;
	return Array.from({ length: jumlah }, (_, i) => {
		const jari = jariMaks * Math.sqrt((i + 0.5) / jumlah);
		return {
			x: bidang.lebar / 2 + jari * Math.cos(i * sudutEmas),
			y: bidang.tinggi / 2 + jari * Math.sin(i * sudutEmas),
			dx: 0,
			dy: 0
		};
	});
}

function dorong(a: Titik, b: Titik, faktor: number) {
	const dx = (a.x - b.x) * faktor;
	const dy = (a.y - b.y) * faktor;
	a.dx += dx;
	a.dy += dy;
	b.dx -= dx;
	b.dy -= dy;
}

function tambahTolakan(titik: Titik[], k: number) {
	titik.forEach((a, i) => {
		for (const b of titik.slice(i + 1)) {
			const jarak = Math.max(Math.hypot(a.x - b.x, a.y - b.y), jarakMinimum);
			dorong(a, b, (k * k) / jarak / jarak);
		}
	});
}

function tambahTarikan(pasangan: [Titik, Titik][], k: number) {
	for (const [a, b] of pasangan) {
		dorong(a, b, -Math.hypot(a.x - b.x, a.y - b.y) / k);
	}
}

function batasi(nilai: number, ukuran: number, tepi: number): number {
	return Math.min(ukuran - tepi, Math.max(tepi, nilai));
}

function geser(titik: Titik[], suhu: number, bidang: Bidang) {
	for (const p of titik) {
		const dx = p.dx + (bidang.lebar / 2 - p.x) * tarikanPusat;
		const dy = p.dy + (bidang.tinggi / 2 - p.y) * tarikanPusat;
		const langkah = Math.min(Math.hypot(dx, dy), suhu) / Math.max(Math.hypot(dx, dy), jarakMinimum);
		p.x = batasi(p.x + dx * langkah, bidang.lebar, bidang.tepi);
		p.y = batasi(p.y + dy * langkah, bidang.tinggi, bidang.tepi);
		p.dx = 0;
		p.dy = 0;
	}
}

function pasanganTitik(menurutId: Map<string, Titik>, sisi: SisiJaringan[]): [Titik, Titik][] {
	return sisi.flatMap((s) => {
		const a = menurutId.get(s.dari);
		const b = menurutId.get(s.ke);
		return a && b ? [[a, b] as [Titik, Titik]] : [];
	});
}

export function susunTataLetak(
	id: string[],
	sisi: SisiJaringan[],
	bidang: Bidang
): Map<string, Posisi> {
	const titik = titikAwal(id.length, bidang);
	const menurutId = new Map(titik.map((p, i) => [id[i] ?? String(i), p]));
	const pasangan = pasanganTitik(menurutId, sisi);
	const k = skalaJarakIdeal * Math.sqrt((bidang.lebar * bidang.tinggi) / Math.max(id.length, 1));
	const suhuAwal = bidang.lebar / 10;
	for (let iterasi = 0; iterasi < jumlahIterasi; iterasi++) {
		tambahTolakan(titik, k);
		tambahTarikan(pasangan, k);
		geser(titik, suhuAwal * (1 - iterasi / jumlahIterasi), bidang);
	}
	return menurutId;
}
