import { Droplet, GraduationCap, Package, Route } from '@lucide/svelte';
import type { KodeMata, MataRantai, Sektor, StatusPemeliharaan } from '$lib/api/tipe';
import { t, type KunciPesan } from '$lib/i18n';

export const kunciTahap: Record<KodeMata, KunciPesan> = {
	rencana: 'chain.stage.rencana',
	pemilihan: 'chain.stage.pemilihan',
	pemenang: 'chain.stage.pemenang',
	pelaksanaan: 'chain.stage.pelaksanaan',
	lokasi: 'chain.stage.lokasi',
	serah_terima: 'chain.stage.serah-terima',
	masa_pemeliharaan: 'chain.stage.pemeliharaan',
	bukti_warga: 'chain.stage.bukti'
};

export const kunciSektor: Record<Sektor, KunciPesan> = {
	pendidikan: 'sector.edu',
	air: 'sector.air',
	jalan: 'sector.jalan',
	'': 'sector.lain'
};

export const ikonSektor: Record<Sektor, typeof Package> = {
	pendidikan: GraduationCap,
	air: Droplet,
	jalan: Route,
	'': Package
};

export const kunciPemeliharaan: Record<StatusPemeliharaan, KunciPesan> = {
	pelaksanaan: 'maint.status.pelaksanaan',
	pemeliharaan: 'maint.status.pemeliharaan',
	pemeliharaan_berakhir: 'maint.status.berakhir',
	tidak_dapat_diestimasi: 'maint.status.tidak'
};

export function ringkasRantai(rantai: MataRantai[]): string {
	const tersambung = rantai.filter((m) => m.status === 'ada').length;
	const putus = rantai.filter((m) => m.status === 'putus').length;
	const parameter = { ok: tersambung, total: rantai.length, broken: putus };
	return putus ? t('chain.summaryBroken', parameter) : t('chain.summary', parameter);
}

export function daftarTahun(tahunIni: number, jumlah: number): number[] {
	return Array.from({ length: jumlah }, (_, i) => tahunIni - i);
}

const jumlahSegmenKodeKecamatan = 3;

export function kodeKecamatanDari(kodeWilayah: string): string {
	return kodeWilayah.split('.').slice(0, jumlahSegmenKodeKecamatan).join('.');
}
