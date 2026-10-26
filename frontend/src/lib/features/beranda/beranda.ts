import type { Dapil, IndikatorRingkas, TemaKecamatan, Wakil } from '$lib/api/tipe';
import type { TemaPendidikan } from '$lib/features/pendidikan/jenjang';

export interface WilayahSorotan {
	kode: string;
	nama: string;
}

export const sorotanBawaan: WilayahSorotan = { kode: '31.75.07', nama: 'Duren Sawit' };

export const daftarSorotan: WilayahSorotan[] = [
	sorotanBawaan,
	{ kode: '31.72.04', nama: 'Cilincing' },
	{ kode: '31.72.06', nama: 'Kelapa Gading' },
	{ kode: '31.71.07', nama: 'Tanah Abang' }
];

export const batasPaketBeranda = 3;
export const batasWakilBeranda = 3;
const persenPenuh = 100;

export function bacaSorotan(kode: string | null): WilayahSorotan {
	return daftarSorotan.find((w) => w.kode === kode) ?? sorotanBawaan;
}

export function tautanSorotan(url: URL, kode: string): string {
	const tujuan = new URL(url);
	if (kode === sorotanBawaan.kode) tujuan.searchParams.delete('wilayah');
	else tujuan.searchParams.set('wilayah', kode);
	return tujuan.pathname + tujuan.search;
}

export function indikatorPendidikan(
	tema: TemaKecamatan[],
	jenjang: TemaPendidikan
): IndikatorRingkas | null {
	return tema.find((item) => item.kode === jenjang)?.indikator ?? null;
}

export function persenTerhadap(nilai: number, acuan: number): number {
	if (acuan <= 0) return 0;
	return Math.min(persenPenuh, Math.max(0, (nilai / acuan) * persenPenuh));
}

export function wakilPengawas(dapil: Dapil): Wakil[] {
	return dapil.anggota.filter((w) => w.duduk_di_komisi_pengawas).slice(0, batasWakilBeranda);
}
