import { House, Info, Map } from '@lucide/svelte';
import type { KunciPesan } from '$lib/i18n';

export type Tujuan = 'beranda' | 'peta' | 'tentang';

interface TujuanNavigasi {
	kunci: Tujuan;
	path: string;
	label: KunciPesan;
	ikon: typeof House;
}

export const tujuanUtama: TujuanNavigasi[] = [
	{ kunci: 'beranda', path: '', label: 'nav.home', ikon: House },
	{ kunci: 'peta', path: '/peta', label: 'nav.map', ikon: Map },
	{ kunci: 'tentang', path: '/tentang-data', label: 'nav.about', ikon: Info }
];

export function tujuanAktif(pathname: string, bahasa: string): Tujuan | null {
	const sisa = pathname.slice(bahasa.length + 1);
	if (sisa === '' || sisa === '/') return 'beranda';
	return tujuanUtama.find((tujuan) => tujuan.path && sisa.startsWith(tujuan.path))?.kunci ?? null;
}
