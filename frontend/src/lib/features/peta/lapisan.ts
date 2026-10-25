import type { SekitarTitik, SekolahTerdekat } from '$lib/api/tipe';

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
export const zoomKota = 10;
export const zoomSekolah = 15;
