import {
	ambilDaftarTema,
	ambilTema,
	ambilWilayah,
	GalatAPI,
	susunAspirasi,
	type Fetch
} from '$lib/api/klien';
import type { Dokumen, TemaDetail, WilayahDitemukan } from '$lib/api/tipe';
import { t } from '$lib/i18n';

export const bahasaDokumen = ['id', 'en'] as const;
export type BahasaDokumen = (typeof bahasaDokumen)[number];

export interface PilihanAspirasi {
	kode: string;
	tema: string;
	masalah: string;
	bahasa: BahasaDokumen;
}

export interface HalamanAspirasi {
	wilayah: WilayahDitemukan;
	tema: TemaDetail;
	dokumen: Dokumen | null;
}

export function bacaPilihan(kueri: URLSearchParams): PilihanAspirasi {
	const bahasa = kueri.get('bahasa');
	return {
		kode: kueri.get('w') ?? '',
		tema: kueri.get('tema') ?? '',
		masalah: kueri.get('masalah') ?? '',
		bahasa: bahasa === 'en' ? 'en' : 'id'
	};
}

export function tautanAspirasi(url: URL, ubah: Partial<PilihanAspirasi>): string {
	const kueri = new URLSearchParams(url.search);
	const nama: Record<keyof PilihanAspirasi, string> = {
		kode: 'w',
		tema: 'tema',
		masalah: 'masalah',
		bahasa: 'bahasa'
	};
	for (const [kunci, nilai] of Object.entries(ubah) as [keyof PilihanAspirasi, string][]) {
		kueri.set(nama[kunci], nilai);
	}
	return `${url.pathname}?${kueri.toString()}`;
}

async function temaTerpilih(fetchFn: Fetch, kode: string): Promise<TemaDetail> {
	if (kode) return ambilTema(fetchFn, kode);
	const [pertama] = await ambilDaftarTema(fetchFn);
	if (!pertama) throw new GalatAPI(404, t('packet.noTheme'));
	return pertama;
}

export async function muatAspirasi(fetchFn: Fetch, p: PilihanAspirasi): Promise<HalamanAspirasi> {
	const [wilayah, tema] = await Promise.all([
		ambilWilayah(fetchFn, p.kode),
		temaTerpilih(fetchFn, p.tema)
	]);
	const dokumen = p.masalah
		? await susunAspirasi(fetchFn, {
				kode_kecamatan: wilayah.kode,
				tema: tema.kode,
				masalah: p.masalah,
				bahasa: p.bahasa
			})
		: null;
	return { wilayah, tema, dokumen };
}
