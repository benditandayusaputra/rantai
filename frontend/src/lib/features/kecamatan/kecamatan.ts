import type { Angka, IndikatorRingkas, KelasRasio } from '$lib/api/tipe';
import { formatAngka, formatRupiah } from '$lib/format';
import { bahasaAktif, isKunciPesan, t, type KunciPesan } from '$lib/i18n';

export const daftarBagian = ['tema', 'anggaran', 'wakilku'] as const;
export type Bagian = (typeof daftarBagian)[number];

export const batasWakilRingkas = 8;
export const batasPaketRingkas = 3;

export function bacaBagian(nilai: string | null): Bagian | null {
	return daftarBagian.find((b) => b === nilai) ?? null;
}

export function tautanBagian(url: URL, bagian: Bagian): string {
	const tujuan = new URL(url);
	tujuan.searchParams.set('bagian', bagian);
	return tujuan.pathname + tujuan.search;
}

export function statusLayanan(
	indikator: IndikatorRingkas | null
): 'perhatian' | 'memadai' | 'belum-dinilai' {
	if (
		!indikator ||
		indikator.kelas === 'tidak_ada_data' ||
		!Number.isFinite(indikator.kebutuhan.nilai) ||
		indikator.kebutuhan.nilai <= 0 ||
		!Number.isFinite(indikator.kapasitas.nilai) ||
		indikator.kapasitas.nilai < 0
	)
		return 'belum-dinilai';
	return indikator.kapasitas.nilai < indikator.kebutuhan.nilai ? 'perhatian' : 'memadai';
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
		return { nilai: formatRupiah(angka.nilai, bahasaAktif), satuan: '' };
	if (angka.satuan === 'rasio') {
		return {
			nilai: formatAngka(angka.nilai, bahasaAktif, desimalRasio),
			satuan: t('satuan.rasio')
		};
	}
	const kunciSatuan = `satuan.${angka.satuan}`;
	return {
		nilai: formatAngka(angka.nilai, bahasaAktif),
		satuan: isKunciPesan(kunciSatuan) ? t(kunciSatuan) : angka.satuan
	};
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
