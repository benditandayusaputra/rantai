import type { MetodeCalonLulusan } from '$lib/api/tipe';
import { t, type KunciPesan } from '$lib/i18n';

export const daftarTemaPendidikan = ['pendidikan', 'pendidikan_sma_smk'] as const;
export type TemaPendidikan = (typeof daftarTemaPendidikan)[number];
export const temaPendidikanBawaan: TemaPendidikan = 'pendidikan';

interface TeksJenjang {
	jenjang: KunciPesan;
	judul: KunciPesan;
	judulGrafik: KunciPesan;
	lulusan: KunciPesan;
	kursiNegeri: KunciPesan;
	kursiSPMB: KunciPesan;
	kursiMadrasah: KunciPesan;
	metodePerSekolah: KunciPesan;
	judulSekolah: KunciPesan;
	kolomKursi: KunciPesan;
	sekolahKosong: KunciPesan;
	judulTren: KunciPesan;
	leadTren: KunciPesan;
	legenda: KunciPesan;
	tabelSekolah: KunciPesan;
	kalimatTekanan: KunciPesan;
	cekLead: KunciPesan;
	cekTerdekat: KunciPesan;
	cekKosong: KunciPesan;
	dasborTanya: KunciPesan;
	dasborLulusan: KunciPesan;
	dasborKursi: KunciPesan;
	belanjaLead: KunciPesan;
	belanjaSumbuX: KunciPesan;
}

export const teksJenjang: Record<TemaPendidikan, TeksJenjang> = {
	pendidikan: {
		jenjang: 'level.smp',
		judul: 'pen.title',
		judulGrafik: 'pen.chartTitle',
		lulusan: 'pen.graduates',
		kursiNegeri: 'pen.seats.negeri',
		kursiSPMB: 'pen.seats.spmb',
		kursiMadrasah: 'pen.seats.madrasah',
		metodePerSekolah: 'pen.method.per_sekolah',
		judulSekolah: 'pen.schoolsHeading',
		kolomKursi: 'pen.school.seats',
		sekolahKosong: 'pen.schoolsEmpty',
		judulTren: 'pen.trend.title',
		leadTren: 'pen.trend.lead',
		legenda: 'map.legend',
		tabelSekolah: 'map.table.sekolah',
		kalimatTekanan: 'services.sentence',
		cekLead: 'check.lead',
		cekTerdekat: 'check.nearest',
		cekKosong: 'check.empty',
		dasborTanya: 'dash.services.q',
		dasborLulusan: 'dash.svc.graduates',
		dasborKursi: 'dash.svc.seats',
		belanjaLead: 'spend.lead',
		belanjaSumbuX: 'spend.xAxis'
	},
	pendidikan_sma_smk: {
		jenjang: 'level.sma',
		judul: 'pen.sma.title',
		judulGrafik: 'pen.sma.chartTitle',
		lulusan: 'pen.sma.graduates',
		kursiNegeri: 'pen.sma.seats.negeri',
		kursiSPMB: 'pen.sma.seats.spmb',
		kursiMadrasah: 'pen.sma.seats.madrasah',
		metodePerSekolah: 'pen.sma.method.per_sekolah',
		judulSekolah: 'pen.sma.schoolsHeading',
		kolomKursi: 'pen.sma.school.seats',
		sekolahKosong: 'pen.sma.schoolsEmpty',
		judulTren: 'pen.sma.trend.title',
		leadTren: 'pen.sma.trend.lead',
		legenda: 'map.sma.legend',
		tabelSekolah: 'map.sma.table.sekolah',
		kalimatTekanan: 'services.sma.sentence',
		cekLead: 'check.sma.lead',
		cekTerdekat: 'check.sma.nearest',
		cekKosong: 'check.sma.empty',
		dasborTanya: 'dash.sma.services.q',
		dasborLulusan: 'dash.sma.svc.graduates',
		dasborKursi: 'dash.sma.svc.seats',
		belanjaLead: 'spend.sma.lead',
		belanjaSumbuX: 'spend.sma.xAxis'
	}
};

export function isTemaPendidikan(nilai: string | null | undefined): nilai is TemaPendidikan {
	return (daftarTemaPendidikan as readonly string[]).includes(nilai ?? '');
}

export function bacaTemaPendidikan(nilai: string | null): TemaPendidikan {
	return isTemaPendidikan(nilai) ? nilai : temaPendidikanBawaan;
}

export function tautanTemaPendidikan(url: URL, tema: TemaPendidikan): string {
	const tujuan = new URL(url);
	if (tema === temaPendidikanBawaan) tujuan.searchParams.delete('tema');
	else tujuan.searchParams.set('tema', tema);
	return tujuan.pathname + tujuan.search;
}

export function pilihanJenjang(url: URL, aktif: TemaPendidikan) {
	return daftarTemaPendidikan.map((tema) => ({
		key: tema,
		href: tautanTemaPendidikan(url, tema),
		label: t(teksJenjang[tema].jenjang),
		current: tema === aktif
	}));
}

export function kunciMetode(tema: TemaPendidikan, metode: MetodeCalonLulusan): KunciPesan | null {
	if (!metode) return null;
	return metode === 'per_sekolah' ? teksJenjang[tema].metodePerSekolah : `pen.method.${metode}`;
}
