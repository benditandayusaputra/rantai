import { Camera, House, Info, MapPin, MessageCircle } from '@lucide/svelte';
import type { KunciPesan } from '$lib/i18n';

export type Tujuan =
	'beranda' | 'wilayah' | 'peta' | 'dasbor' | 'tanya' | 'lapor' | 'pantau' | 'tentang';

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
}

export const tujuanUtama: TujuanNavigasi[] = [
	{ kunci: 'beranda', path: '', label: 'nav.home', ikon: House },
	{ kunci: 'wilayah', path: '/wilayah', label: 'nav.area', ikon: MapPin },
	{ kunci: 'tanya', path: '/tanya', label: 'nav.ask', ikon: MessageCircle },
	{ kunci: 'lapor', path: '/lapor', label: 'nav.report', ikon: Camera },
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
			{ path: '/peta', label: 'nav.map' },
			{ path: '/dasbor', label: 'nav.dashboard' },
			{ path: '/tanya', label: 'ask.title' },
			{ path: '/cek-sekolah', label: 'check.title' }
		]
	},
	{
		judul: 'nav.act',
		tautan: [
			{ path: '/lapor', label: 'report.flowTitle' },
			{ path: '/pantau', label: 'watch.title' }
		]
	},
	{
		judul: 'nav.openness',
		tautan: [
			{ path: '/tentang-data', label: 'nav.about' },
			{ path: '/data-belum-dibuka', label: 'unopened.title' }
		]
	}
];

export function tujuanAktif(pathname: string, bahasa: string): Tujuan | null {
	const sisa = pathname.slice(bahasa.length + 1);
	if (sisa === '' || sisa === '/') return 'beranda';
	if (/^\/(kecamatan|peta|dasbor|cek-sekolah|anggota|paket|penyedia)(?:\/|$)/.test(sisa))
		return 'wilayah';
	const isCocok = (path: string) => sisa === path || sisa.startsWith(`${path}/`);
	return tujuanUtama.find((tujuan) => tujuan.path && isCocok(tujuan.path))?.kunci ?? null;
}

export const tujuanPonsel = tujuanUtama.filter((tujuan) => !tujuan.isHanyaLayarLebar);

export function kodeWilayahDariURL(url: URL): string | null {
	const kandidat =
		url.pathname.match(/\/kecamatan\/(31\.\d{2}\.\d{2})(?:\/|$)/)?.[1] ?? url.searchParams.get('w');
	return kandidat && /^31\.\d{2}\.\d{2}(?:\.\d{4})?$/.test(kandidat)
		? kandidat.split('.').slice(0, 3).join('.')
		: null;
}

export function tautanDenganWilayah(path: string, bahasa: string, kode: string | null): string {
	if (path === '/wilayah') return kode ? `/${bahasa}/kecamatan/${kode}` : `/${bahasa}#cari-wilayah`;
	const tujuan = `/${bahasa}${path}`;
	return kode && (path === '/tanya' || path === '/cek-sekolah')
		? `${tujuan}?w=${encodeURIComponent(kode)}`
		: tujuan;
}
