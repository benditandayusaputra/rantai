import { Bell, Camera, House, Info, Map, MessageCircle } from '@lucide/svelte';
import type { KunciPesan } from '$lib/i18n';

export type Tujuan = 'beranda' | 'peta' | 'tanya' | 'lapor' | 'pantau' | 'tentang';

interface TujuanNavigasi {
	kunci: Tujuan;
	path: string;
	label: KunciPesan;
	ikon: typeof House;
	isHanyaLayarLebar?: boolean;
}

export const tujuanUtama: TujuanNavigasi[] = [
	{ kunci: 'beranda', path: '', label: 'nav.home', ikon: House },
	{ kunci: 'peta', path: '/peta', label: 'nav.map', ikon: Map },
	{ kunci: 'tanya', path: '/tanya', label: 'nav.ask', ikon: MessageCircle },
	{ kunci: 'lapor', path: '/lapor', label: 'nav.report', ikon: Camera },
	{ kunci: 'pantau', path: '/pantau', label: 'nav.watch', ikon: Bell, isHanyaLayarLebar: true },
	{ kunci: 'tentang', path: '/tentang-data', label: 'nav.about', ikon: Info }
];

export function tujuanAktif(pathname: string, bahasa: string): Tujuan | null {
	const sisa = pathname.slice(bahasa.length + 1);
	if (sisa === '' || sisa === '/') return 'beranda';
	const isCocok = (path: string) => sisa === path || sisa.startsWith(`${path}/`);
	return tujuanUtama.find((tujuan) => tujuan.path && isCocok(tujuan.path))?.kunci ?? null;
}

export const tujuanPonsel = tujuanUtama.filter((tujuan) => !tujuan.isHanyaLayarLebar);
