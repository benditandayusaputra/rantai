import {
	Bell,
	Camera,
	ChartNoAxesCombined,
	FileLock,
	House,
	Info,
	Map,
	MessageCircle,
	School
} from '@lucide/svelte';
import type { KunciPesan } from '$lib/i18n';

export type Tujuan = 'beranda' | 'peta' | 'dasbor' | 'tanya' | 'lapor' | 'pantau' | 'tentang';

interface TujuanNavigasi {
	kunci: Tujuan;
	path: string;
	label: KunciPesan;
	labelPonsel?: KunciPesan;
	ikon: typeof House;
	isHanyaLayarLebar?: boolean;
}

interface TautanKaki {
	path: string;
	label: KunciPesan;
	ikon: typeof House;
}

export const tujuanUtama: TujuanNavigasi[] = [
	{ kunci: 'beranda', path: '', label: 'nav.home', ikon: House },
	{ kunci: 'peta', path: '/peta', label: 'nav.map', ikon: Map },
	{ kunci: 'dasbor', path: '/dasbor', label: 'nav.dashboard', ikon: ChartNoAxesCombined },
	{ kunci: 'tanya', path: '/tanya', label: 'nav.ask', ikon: MessageCircle },
	{ kunci: 'lapor', path: '/lapor', label: 'nav.report', ikon: Camera },
	{ kunci: 'pantau', path: '/pantau', label: 'nav.watch', ikon: Bell, isHanyaLayarLebar: true },
	{
		kunci: 'tentang',
		path: '/tentang-data',
		label: 'nav.about',
		labelPonsel: 'nav.aboutShort',
		ikon: Info
	}
];

export const kelompokKaki: { judul: KunciPesan; tautan: TautanKaki[] }[] = [
	{
		judul: 'nav.explore',
		tautan: [
			{ path: '/peta', label: 'nav.map', ikon: Map },
			{ path: '/dasbor', label: 'nav.dashboard', ikon: ChartNoAxesCombined },
			{ path: '/tanya', label: 'ask.title', ikon: MessageCircle },
			{ path: '/cek-sekolah', label: 'check.title', ikon: School }
		]
	},
	{
		judul: 'nav.act',
		tautan: [
			{ path: '/lapor', label: 'report.flowTitle', ikon: Camera },
			{ path: '/pantau', label: 'watch.title', ikon: Bell }
		]
	},
	{
		judul: 'nav.openness',
		tautan: [
			{ path: '/tentang-data', label: 'nav.about', ikon: Info },
			{ path: '/data-belum-dibuka', label: 'unopened.title', ikon: FileLock }
		]
	}
];

export function tujuanAktif(pathname: string, bahasa: string): Tujuan | null {
	const sisa = pathname.slice(bahasa.length + 1);
	if (sisa === '' || sisa === '/') return 'beranda';
	const isCocok = (path: string) => sisa === path || sisa.startsWith(`${path}/`);
	return tujuanUtama.find((tujuan) => tujuan.path && isCocok(tujuan.path))?.kunci ?? null;
}

export const tujuanPonsel = tujuanUtama.filter((tujuan) => !tujuan.isHanyaLayarLebar);
