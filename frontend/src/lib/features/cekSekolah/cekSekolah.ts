import type { SekolahTerdekat } from '$lib/api/tipe';
import { formatAngka } from '$lib/format';
import { bahasaBawaan } from '$lib/i18n';

export const radiusBawaanM = 3000;
export const batasDaftarNegeri = 5;
const meterPerKilometer = 1000;
const desimalRasio = 1;

export function radiusKm(meter: number): string {
	return formatAngka(meter / meterPerKilometer, bahasaBawaan);
}

export function pecahanTekanan(rasio: number): { a: string; b: string } {
	if (rasio >= 1) return { a: formatAngka(rasio, bahasaBawaan, desimalRasio), b: '1' };
	return { a: '1', b: formatAngka(1 / rasio, bahasaBawaan, desimalRasio) };
}

export function terdekatDulu(daftar: SekolahTerdekat[], batas: number): SekolahTerdekat[] {
	return [...daftar].sort((x, y) => x.jarak_m - y.jarak_m).slice(0, batas);
}
