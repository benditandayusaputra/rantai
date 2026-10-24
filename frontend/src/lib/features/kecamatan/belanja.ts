import type { BelanjaKecamatan, BelanjaTema } from '$lib/api/tipe';

export interface UkuranGrafik {
	lebar: number;
	tinggi: number;
	kiri: number;
	kanan: number;
	atas: number;
	bawah: number;
}

const lebarMinimum = 280;
const tinggiPonsel = 260;
const tinggiLebar = 360;
const ambangLayarLebar = 640;
const margin = { kiri: 84, kanan: 12, atas: 24, bawah: 52 };

export function ukuranUntuk(lebarKontainer: number): UkuranGrafik {
	const lebar = Math.max(lebarMinimum, Math.round(lebarKontainer));
	return { lebar, tinggi: lebar >= ambangLayarLebar ? tinggiLebar : tinggiPonsel, ...margin };
}

export interface TitikSebar {
	kode: string;
	nama: string;
	x: number;
	y: number;
	rasio: number;
	nilai: number;
	isAktif: boolean;
}

function batasAtas(nilai: number[]): number {
	const maksimum = Math.max(0, ...nilai);
	return maksimum > 0 ? maksimum * 1.05 : 1;
}

export function batasSumbu(data: BelanjaTema): { maksX: number; maksY: number } {
	const terukur = data.kecamatan.filter((k) => k.indikator?.rasio);
	return {
		maksX: batasAtas(terukur.map((k) => k.indikator?.rasio?.nilai ?? 0)),
		maksY: batasAtas(terukur.map((k) => k.belanja?.total_hps.nilai ?? 0))
	};
}

export function titikSebar(
	data: BelanjaTema,
	kodeAktif: string,
	ukuran: UkuranGrafik
): TitikSebar[] {
	const terukur = data.kecamatan.filter((k) => k.indikator?.rasio);
	const { maksX, maksY } = batasSumbu(data);
	const { lebar, tinggi, kiri, kanan, atas, bawah } = ukuran;
	const titik = terukur.map((k) => {
		const rasio = k.indikator?.rasio?.nilai ?? 0;
		const nilai = k.belanja?.total_hps.nilai ?? 0;
		return {
			kode: k.kecamatan.kode,
			nama: k.kecamatan.nama,
			x: kiri + (rasio / maksX) * (lebar - kiri - kanan),
			y: tinggi - bawah - (nilai / maksY) * (tinggi - atas - bawah),
			rasio,
			nilai,
			isAktif: k.kecamatan.kode === kodeAktif
		};
	});
	return [...titik.filter((t) => !t.isAktif), ...titik.filter((t) => t.isAktif)];
}

export function jumlahTanpaRasio(data: BelanjaTema): number {
	return data.kecamatan.filter((k) => !k.indikator?.rasio).length;
}

export function urutBaris(data: BelanjaTema, kodeAktif: string): BelanjaKecamatan[] {
	return [...data.kecamatan].sort((a, b) => {
		if (a.kecamatan.kode === kodeAktif) return -1;
		if (b.kecamatan.kode === kodeAktif) return 1;
		return a.kecamatan.nama.localeCompare(b.kecamatan.nama, 'id');
	});
}
