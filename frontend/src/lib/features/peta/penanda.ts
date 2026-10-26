import { Camera, School } from '@lucide/svelte';
import type { ExpressionSpecification, Map as PetaMaplibre } from 'maplibre-gl';
import { mount, unmount } from 'svelte';
import type { Sektor, StatusLaporan } from '$lib/api/tipe';
import { ikonSektor } from '$lib/features/paket/paket';
import type { BacaToken } from './gaya';

type Ikon = typeof School;
type Bentuk = 'pin' | 'lencana';

export interface SpekPenanda {
	nama: string;
	ikon: Ikon;
	warna: string;
	warnaIkon: string;
	bentuk: Bentuk;
}

const rasioPiksel = 2;
const tebalGarisIkon = 2.25;
const tebalTepi = 2;
const ukuranBentuk: Record<Bentuk, { lebar: number; tinggi: number; ikon: number; atas: number }> =
	{
		pin: { lebar: 30, tinggi: 38, ikon: 16, atas: 6.5 },
		lencana: { lebar: 26, tinggi: 26, ikon: 14, atas: 6 }
	};
const jalurPin = 'M15 37C15 37 3 24 3 14.5a12 12 0 0 1 24 0C27 24 15 37 15 37Z';
const lingkaranLencana = { pusat: 13, jari: 11.5 };

export const tokenSektor: Record<Sektor, string> = {
	pendidikan: 'sector-edu',
	air: 'sector-air',
	jalan: 'sector-jalan',
	'': 'sector-lain'
};
const tokenIkonSektor: Record<Sektor, string> = {
	pendidikan: 'sector-edu-on',
	air: 'sector-air-on',
	jalan: 'sector-jalan-on',
	'': 'sector-lain-on'
};
export const daftarSektor = Object.keys(tokenSektor) as Sektor[];

export const statusLaporanTuntas: StatusLaporan[] = ['TERVERIFIKASI', 'SELESAI'];
export const penandaSekolah = 'penanda-sekolah';
const awalanPenandaProyek = 'penanda-proyek-';
const penandaLaporan = {
	tuntas: 'penanda-laporan-tuntas',
	proses: 'penanda-laporan-proses',
	simulasi: 'penanda-laporan-simulasi'
} as const;

export function namaPenandaProyek(sektor: Sektor): string {
	return `${awalanPenandaProyek}${sektor || 'lain'}`;
}

export const ekspresiPenandaProyek: ExpressionSpecification = [
	'match',
	['get', 'sektor'],
	['pendidikan', 'air', 'jalan'],
	['concat', awalanPenandaProyek, ['get', 'sektor']],
	namaPenandaProyek('')
];

export const ekspresiPenandaLaporan: ExpressionSpecification = [
	'case',
	['get', 'is_simulation'],
	penandaLaporan.simulasi,
	['match', ['get', 'status'], statusLaporanTuntas, penandaLaporan.tuntas, penandaLaporan.proses]
];

export const ukuranPenandaMenurutZoom: ExpressionSpecification = [
	'interpolate',
	['linear'],
	['zoom'],
	10,
	0.5,
	13,
	0.8,
	15,
	1
];

export function daftarPenanda(token: BacaToken): SpekPenanda[] {
	const proyek = daftarSektor.map((sektor) => ({
		nama: namaPenandaProyek(sektor),
		ikon: ikonSektor[sektor],
		warna: token(tokenSektor[sektor]),
		warnaIkon: token(tokenIkonSektor[sektor]),
		bentuk: 'pin' as const
	}));
	const laporan = (warna: string, warnaIkon: string) => ({
		ikon: Camera,
		warna: token(warna),
		warnaIkon: token(warnaIkon),
		bentuk: 'pin' as const
	});
	return [
		{
			nama: penandaSekolah,
			ikon: School,
			warna: token('text'),
			warnaIkon: token('surface'),
			bentuk: 'lencana'
		},
		...proyek,
		{ nama: penandaLaporan.tuntas, ...laporan('success', 'on-success') },
		{ nama: penandaLaporan.proses, ...laporan('info', 'surface') },
		{ nama: penandaLaporan.simulasi, ...laporan('cert-sim', 'surface') }
	];
}

function latarPenanda(bentuk: Bentuk, warna: string, tepi: string): string {
	const gaya = `fill="${warna}" stroke="${tepi}" stroke-width="${tebalTepi}"`;
	if (bentuk === 'pin') return `<path d="${jalurPin}" ${gaya}/>`;
	const { pusat, jari } = lingkaranLencana;
	return `<circle cx="${pusat}" cy="${pusat}" r="${jari}" ${gaya}/>`;
}

export function svgPenanda(bentuk: Bentuk, warna: string, tepi: string, svgIkon: string): string {
	const { lebar, tinggi, ikon, atas } = ukuranBentuk[bentuk];
	return (
		`<svg xmlns="http://www.w3.org/2000/svg" width="${lebar * rasioPiksel}" height="${tinggi * rasioPiksel}" viewBox="0 0 ${lebar} ${tinggi}">` +
		latarPenanda(bentuk, warna, tepi) +
		`<g transform="translate(${(lebar - ikon) / 2} ${atas})">${svgIkon}</g></svg>`
	);
}

function svgIkon(ikon: Ikon, warna: string, ukuran: number): string {
	const wadah = document.createElement('div');
	const komponen = mount(ikon, {
		target: wadah,
		props: { size: ukuran, color: warna, strokeWidth: tebalGarisIkon }
	});
	const svg = wadah.innerHTML;
	void unmount(komponen);
	return svg;
}

async function gambarPenanda(spek: SpekPenanda, tepi: string): Promise<HTMLImageElement> {
	const { lebar, tinggi, ikon } = ukuranBentuk[spek.bentuk];
	const svg = svgPenanda(spek.bentuk, spek.warna, tepi, svgIkon(spek.ikon, spek.warnaIkon, ikon));
	const gambar = new Image(lebar * rasioPiksel, tinggi * rasioPiksel);
	gambar.src = `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`;
	await gambar.decode();
	return gambar;
}

export async function pasangPenanda(peta: PetaMaplibre, token: BacaToken): Promise<void> {
	const siap = await Promise.all(
		daftarPenanda(token).map(async (spek) => ({
			nama: spek.nama,
			gambar: await gambarPenanda(spek, token('surface'))
		}))
	);
	for (const { nama, gambar } of siap) {
		if (peta.hasImage(nama)) peta.updateImage(nama, gambar);
		else peta.addImage(nama, gambar, { pixelRatio: rasioPiksel });
	}
}
