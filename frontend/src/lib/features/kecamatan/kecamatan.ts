import type { Angka, KelasRasio } from '$lib/api/tipe';
import { formatAngka, formatRupiah } from '$lib/format';
import { bahasaBawaan, t, type KunciPesan } from '$lib/i18n';

export const daftarBagian = ['tema', 'anggaran', 'wakilku'] as const;
export type Bagian = (typeof daftarBagian)[number];

export function bacaBagian(nilai: string | null): Bagian {
	return daftarBagian.find((b) => b === nilai) ?? 'tema';
}

export function tautanBagian(url: URL, bagian: Bagian): string {
	const tujuan = new URL(url);
	if (bagian === 'tema') tujuan.searchParams.delete('bagian');
	else tujuan.searchParams.set('bagian', bagian);
	return tujuan.pathname + tujuan.search;
}

export const kunciKelas: Record<KelasRasio, KunciPesan> = {
	sangat_kurang: 'gap.sangat-kurang',
	kurang: 'gap.kurang',
	hampir_cukup: 'gap.hampir-cukup',
	cukup: 'gap.cukup',
	tidak_ada_data: 'gap.tidak-ada-data'
};

const desimalRasio = 2;

export interface NilaiTampil {
	nilai: string;
	satuan: string;
}

export function tampilkanAngka(angka: Angka): NilaiTampil {
	if (angka.satuan === 'rupiah')
		return { nilai: formatRupiah(angka.nilai, bahasaBawaan), satuan: '' };
	if (angka.satuan === 'rasio') {
		return {
			nilai: formatAngka(angka.nilai, bahasaBawaan, desimalRasio),
			satuan: t('satuan.rasio')
		};
	}
	return { nilai: formatAngka(angka.nilai, bahasaBawaan), satuan: angka.satuan };
}

const polaGelarDepan = /^(?:\S+\.\s+)+/;

export function inisial(nama: string): string {
	return nama
		.replace(polaGelarDepan, '')
		.split(/\s+/)
		.filter((kata) => /^\p{Lu}/u.test(kata))
		.slice(0, 2)
		.map((kata) => kata[0])
		.join('');
}

export async function bagikanTautan(url: string): Promise<'dibagikan' | 'disalin'> {
	if (navigator.share) {
		await navigator.share({ url });
		return 'dibagikan';
	}
	await navigator.clipboard.writeText(url);
	return 'disalin';
}
