import type { WilayahDitemukan } from '$lib/api/tipe';
import { bahasaAktif } from '$lib/i18n';

export const panjangKataMinimal = 3;
export const jedaKetikMs = 200;

const urutanTingkat = { kelurahan: 0, kecamatan: 1 } as const;
type TingkatHalaman = keyof typeof urutanTingkat;

export type WilayahHalaman = WilayahDitemukan & { tingkat: TingkatHalaman };

function isWilayahHalaman(w: WilayahDitemukan): w is WilayahHalaman {
	return w.tingkat in urutanTingkat;
}

export function saringHasilPencarian(hasil: WilayahDitemukan[]): WilayahHalaman[] {
	return hasil
		.filter(isWilayahHalaman)
		.sort((a, b) => urutanTingkat[a.tingkat] - urutanTingkat[b.tingkat]);
}

export function kodeKecamatan(w: WilayahDitemukan): string {
	return w.hierarki.find((h) => h.tingkat === 'kecamatan')?.kode ?? w.kode;
}

export function tujuanWilayah(w: WilayahDitemukan): string {
	const halaman = `/${bahasaAktif}/kecamatan/${kodeKecamatan(w)}`;
	return w.tingkat === 'kelurahan' ? `${halaman}?kelurahan=${w.kode}` : halaman;
}

export function induk(w: WilayahDitemukan): string {
	return w.hierarki
		.filter((h) => h.tingkat !== 'provinsi')
		.map((h) => h.nama)
		.reverse()
		.join(', ');
}

export function bagiSorotan(nama: string, kata: string): [string, string, string] {
	const awal = nama.toLowerCase().indexOf(kata.trim().toLowerCase());
	if (awal < 0 || !kata.trim()) return [nama, '', ''];
	const akhir = awal + kata.trim().length;
	return [nama.slice(0, awal), nama.slice(awal, akhir), nama.slice(akhir)];
}

export const contohWilayah = ['Cilincing', 'Kelapa Gading', 'Tanah Abang'];
