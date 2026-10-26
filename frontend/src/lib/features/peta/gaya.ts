import type {
	FilterSpecification,
	LayerSpecification,
	SourceSpecification,
	StyleSpecification
} from 'maplibre-gl';
import type { KelasRasio } from '$lib/api/tipe';

export const awalanPMTiles = 'pmtiles://';
export const sumberDasarBawaan = `${awalanPMTiles}/peta/dki.pmtiles`;

interface SkemaUbin {
	atribusi: string;
	taman: string;
	filterTaman?: FilterSpecification;
	air: string;
	jalan: string;
}

const skemaProtomaps: SkemaUbin = {
	atribusi: '© OpenStreetMap contributors, Protomaps',
	taman: 'landuse',
	filterTaman: ['in', ['get', 'kind'], ['literal', ['park', 'forest', 'grass', 'nature_reserve']]],
	air: 'water',
	jalan: 'roads'
};

const skemaOpenMapTiles: SkemaUbin = {
	atribusi: '© OpenStreetMap contributors, OpenFreeMap',
	taman: 'park',
	air: 'water',
	jalan: 'transportation'
};

export function skemaUbin(urlSumber: string): SkemaUbin {
	return urlSumber.startsWith(awalanPMTiles) ? skemaProtomaps : skemaOpenMapTiles;
}

export function atribusiDasar(urlSumber: string): string {
	return skemaUbin(urlSumber).atribusi;
}

export function urlSumberDasar(dariEnv: string | undefined): string {
	return dariEnv || sumberDasarBawaan;
}

export type BacaToken = (nama: string) => string;

export function bacaTokenDokumen(nama: string): string {
	return getComputedStyle(document.documentElement).getPropertyValue(`--${nama}`).trim();
}

export const daftarPetaDasar = ['standar', 'jalan', 'satelit'] as const;
export type PetaDasar = (typeof daftarPetaDasar)[number];
type PetaDasarRaster = Exclude<PetaDasar, 'standar'>;
export const petaDasarBawaan: PetaDasar = 'standar';

const ubinRaster: Record<PetaDasarRaster, { url: string; atribusi: string }> = {
	jalan: {
		url: 'https://tile.openstreetmap.org/{z}/{x}/{y}.png',
		atribusi: '© OpenStreetMap contributors'
	},
	satelit: {
		url: 'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}',
		atribusi: 'Esri, Maxar, Earthstar Geographics'
	}
};
const daftarPetaDasarRaster = Object.keys(ubinRaster) as PetaDasarRaster[];
const ukuranUbinRaster = 256;
const zoomUbinRasterMaks = 19;
const ubinContohJakarta = { z: '12', x: '3263', y: '2118' };

export function isPetaDasar(nilai: string | null): nilai is PetaDasar {
	return (daftarPetaDasar as readonly (string | null)[]).includes(nilai);
}

export function idLapisanRaster(dasar: PetaDasarRaster): string {
	return `dasar-${dasar}`;
}

export function isPetaDasarRaster(dasar: PetaDasar): dasar is PetaDasarRaster {
	return dasar !== 'standar';
}

export function contohUbin(dasar: PetaDasarRaster): string {
	return ubinRaster[dasar].url
		.replace('{z}', ubinContohJakarta.z)
		.replace('{x}', ubinContohJakarta.x)
		.replace('{y}', ubinContohJakarta.y);
}

function sumberRaster(): Record<string, SourceSpecification> {
	return Object.fromEntries(
		daftarPetaDasarRaster.map((dasar) => [
			idLapisanRaster(dasar),
			{
				type: 'raster',
				tiles: [ubinRaster[dasar].url],
				tileSize: ukuranUbinRaster,
				maxzoom: zoomUbinRasterMaks,
				attribution: ubinRaster[dasar].atribusi
			}
		])
	);
}

function lapisanRaster(): LayerSpecification[] {
	return daftarPetaDasarRaster.map((dasar) => ({
		id: idLapisanRaster(dasar),
		type: 'raster',
		source: idLapisanRaster(dasar),
		layout: { visibility: 'none' }
	}));
}

export function buatGayaDasar(urlSumber: string, token: BacaToken): StyleSpecification {
	const skema = skemaUbin(urlSumber);
	return {
		version: 8,
		sources: {
			dasar: { type: 'vector', url: urlSumber, attribution: skema.atribusi },
			...sumberRaster()
		},
		layers: [
			{ id: 'latar', type: 'background', paint: { 'background-color': token('map') } },
			{
				id: 'taman',
				type: 'fill',
				source: 'dasar',
				'source-layer': skema.taman,
				...(skema.filterTaman ? { filter: skema.filterTaman } : {}),
				paint: { 'fill-color': token('map-park') }
			},
			{
				id: 'air',
				type: 'fill',
				source: 'dasar',
				'source-layer': skema.air,
				paint: { 'fill-color': token('map-water') }
			},
			{
				id: 'jalan',
				type: 'line',
				source: 'dasar',
				'source-layer': skema.jalan,
				paint: { 'line-color': token('map-road'), 'line-width': 1 }
			},
			...lapisanRaster()
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
