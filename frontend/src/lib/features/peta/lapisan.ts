import type { SekitarTitik, SekolahTerdekat } from '$lib/api/tipe';
import { isPetaDasar, petaDasarBawaan, type PetaDasar } from './gaya';

export const daftarLapisan = ['tema', 'proyek', 'sekolah', 'laporan'] as const;
export type Lapisan = (typeof daftarLapisan)[number];
export type StatusLapisan = Record<Lapisan, boolean>;

const lapisanBawaan: StatusLapisan = { tema: true, proyek: false, sekolah: false, laporan: false };

export function bacaLapisan(nilai: string | null): StatusLapisan {
	if (nilai === null) return { ...lapisanBawaan };
	const aktif = new Set(nilai.split(','));
	return Object.fromEntries(daftarLapisan.map((l) => [l, aktif.has(l)])) as StatusLapisan;
}

export function tulisLapisan(status: StatusLapisan): string {
	return daftarLapisan.filter((l) => status[l]).join(',');
}

export function bacaPetaDasar(nilai: string | null): PetaDasar {
	return isPetaDasar(nilai) ? nilai : petaDasarBawaan;
}

export function urlDenganStatusPeta(url: URL, lapisan: StatusLapisan, dasar: PetaDasar): URL {
	const tujuan = new URL(url);
	tujuan.searchParams.set('lapisan', tulisLapisan(lapisan));
	if (dasar === petaDasarBawaan) tujuan.searchParams.delete('dasar');
	else tujuan.searchParams.set('dasar', dasar);
	return tujuan;
}

const parameterPetaPenuh = ['tema', 'lapisan', 'dasar'];

export function tautanPetaPenuh(url: URL, bahasa: string): string {
	const tujuan = new URL(`/${bahasa}/peta`, url);
	for (const nama of parameterPetaPenuh) {
		const nilai = url.searchParams.get(nama);
		if (nilai !== null) tujuan.searchParams.set(nama, nilai);
	}
	return tujuan.pathname + tujuan.search;
}

export const radiusSekolahMaksimalM = 10_000;
const meterPerDerajatLintang = 111_320;

export function radiusTampilanM(utara: number, selatan: number): number {
	const setengahTinggi = ((utara - selatan) / 2) * meterPerDerajatLintang;
	return Math.min(radiusSekolahMaksimalM, Math.max(1, Math.round(setengahTinggi)));
}

export function gabungSekolah(sekitar: SekitarTitik): SekolahTerdekat[] {
	return [...sekitar.negeri, ...sekitar.spmb_bersama, ...sekitar.madrasah];
}

export const pusatJakarta = { lat: -6.2, lng: 106.83 } as const;
export const batasJakarta: [[number, number], [number, number]] = [
	[106.68, -6.38],
	[106.98, -6.08]
];
export const zoomKota = 10;
export const zoomSekolah = 15;
