import { describe, expect, it } from 'vitest';
import {
	buatGayaDasar,
	contohUbin,
	idLapisanRaster,
	isPetaDasar,
	sumberDasarBawaan,
	warnaKelas
} from './gaya';

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

	it('menyiapkan peta dasar jalan dan satelit dalam keadaan tersembunyi', () => {
		const gaya = buatGayaDasar(sumberDasarBawaan, token);
		const raster = gaya.layers.filter((l) => l.type === 'raster');
		expect(raster.map((l) => l.id)).toEqual([idLapisanRaster('jalan'), idLapisanRaster('satelit')]);
		expect(raster.every((l) => l.layout?.visibility === 'none')).toBe(true);
		expect(gaya.sources[idLapisanRaster('satelit')]).toMatchObject({ type: 'raster' });
	});

	it('mengenali nama peta dasar dan membuat contoh ubin Jakarta', () => {
		expect(isPetaDasar('satelit')).toBe(true);
		expect(isPetaDasar('gelap')).toBe(false);
		expect(contohUbin('jalan')).toBe('https://tile.openstreetmap.org/12/3263/2118.png');
		expect(contohUbin('satelit')).toMatch(/\/tile\/12\/2118\/3263$/);
	});
});
