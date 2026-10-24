import type { StyleSpecification } from 'maplibre-gl';
import type { KelasRasio } from '$lib/api/tipe';

export const sumberDasarBawaan = 'https://tiles.openfreemap.org/planet';
export const atribusiDasar = '© OpenStreetMap contributors, OpenFreeMap';

export type BacaToken = (nama: string) => string;

export function bacaTokenDokumen(nama: string): string {
	return getComputedStyle(document.documentElement).getPropertyValue(`--${nama}`).trim();
}

export function buatGayaDasar(urlSumber: string, token: BacaToken): StyleSpecification {
	return {
		version: 8,
		sources: {
			dasar: { type: 'vector', url: urlSumber, attribution: atribusiDasar }
		},
		layers: [
			{ id: 'latar', type: 'background', paint: { 'background-color': token('map') } },
			{
				id: 'taman',
				type: 'fill',
				source: 'dasar',
				'source-layer': 'park',
				paint: { 'fill-color': token('map-park') }
			},
			{
				id: 'air',
				type: 'fill',
				source: 'dasar',
				'source-layer': 'water',
				paint: { 'fill-color': token('map-water') }
			},
			{
				id: 'jalan',
				type: 'line',
				source: 'dasar',
				'source-layer': 'transportation',
				paint: { 'line-color': token('map-road'), 'line-width': 1 }
			}
		]
	};
}

export const tokenKelas: Record<KelasRasio, string> = {
	sangat_kurang: 'gap-1',
	kurang: 'gap-2',
	hampir_cukup: 'gap-3',
	cukup: 'gap-4',
	tidak_ada_data: 'sunken'
};

export function warnaKelas(token: BacaToken): unknown[] {
	const pasangan = Object.entries(tokenKelas).flatMap(([kelas, nama]) => [kelas, token(nama)]);
	return ['match', ['get', 'kelas'], ...pasangan, token('sunken')];
}

const tokenSektor: Record<string, string> = {
	pendidikan: 'sector-edu',
	air: 'sector-air',
	jalan: 'sector-jalan'
};

export function warnaSektor(token: BacaToken): unknown[] {
	const pasangan = Object.entries(tokenSektor).flatMap(([sektor, nama]) => [sektor, token(nama)]);
	return ['match', ['get', 'sektor'], ...pasangan, token('sector-lain')];
}

export function warnaLaporan(token: BacaToken): unknown[] {
	return [
		'case',
		['get', 'is_simulation'],
		token('cert-sim'),
		['match', ['get', 'status'], ['TERVERIFIKASI', 'SELESAI'], token('success'), token('info')]
	];
}
