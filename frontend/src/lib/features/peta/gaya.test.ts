import { describe, expect, it } from 'vitest';
import { buatGayaDasar, sumberDasarBawaan, warnaKelas, warnaSektor } from './gaya';

const token = (nama: string) => `var-${nama}`;

describe('gaya peta', () => {
	it('mewarnai peta dasar dari token desain', () => {
		const gaya = buatGayaDasar('https://ubin.uji', token);
		const warna = gaya.layers.map((l) => JSON.stringify(l.paint));
		expect(warna.join()).toContain('var-map-water');
		expect(warna.join()).toContain('var-map-road');
		expect(gaya.layers.some((l) => l.type === 'symbol')).toBe(false);
	});

	it('memakai layer skema Protomaps untuk ubin PMTiles bawaan', () => {
		const gaya = buatGayaDasar(sumberDasarBawaan, token);
		const layerSumber = gaya.layers.flatMap((l) =>
			'source-layer' in l ? [l['source-layer']] : []
		);
		expect(layerSumber).toEqual(['landuse', 'water', 'roads']);
	});

	it('memetakan kelas rasio ke skala kesenjangan', () => {
		expect(warnaKelas(token)).toEqual([
			'match',
			['get', 'kelas'],
			'sangat_kurang',
			'var-gap-1',
			'kurang',
			'var-gap-2',
			'hampir_cukup',
			'var-gap-3',
			'cukup',
			'var-gap-4',
			'tidak_ada_data',
			'var-sunken',
			'var-sunken'
		]);
	});

	it('memberi sektor tak dikenal warna lainnya', () => {
		expect(warnaSektor(token).at(-1)).toBe('var-sector-lain');
	});
});
