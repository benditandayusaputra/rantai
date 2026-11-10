import type { Dapil, IndikatorRingkas, TemaKecamatan, Wakil } from '$lib/api/tipe';
import type { TemaPendidikan } from '$lib/features/pendidikan/jenjang';

export interface WilayahSorotan {
	kode: string;
	nama: string;
	kota: string;
}

export const sorotanBawaan: WilayahSorotan = {
	kode: '31.75.07',
	nama: 'Duren Sawit',
	kota: 'Jakarta Timur'
};

export const daftarSorotan: WilayahSorotan[] = [
	sorotanBawaan,
	{ kode: '31.72.04', nama: 'Cilincing', kota: 'Jakarta Utara' },
	{ kode: '31.72.06', nama: 'Kelapa Gading', kota: 'Jakarta Utara' },
	{ kode: '31.71.07', nama: 'Tanah Abang', kota: 'Jakarta Pusat' },
	{ kode: '31.74.07', nama: 'Kebayoran Baru', kota: 'Jakarta Selatan' },
	{ kode: '31.73.01', nama: 'Cengkareng', kota: 'Jakarta Barat' }
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
