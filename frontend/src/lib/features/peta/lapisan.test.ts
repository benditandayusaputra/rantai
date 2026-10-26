import { describe, expect, it } from 'vitest';
import {
	bacaLapisan,
	bacaPetaDasar,
	radiusTampilanM,
	tautanPetaPenuh,
	tulisLapisan,
	urlDenganStatusPeta
} from './lapisan';

describe('status lapisan di URL', () => {
	it('menyalakan lapisan tema bila URL tidak menyebut lapisan', () => {
		expect(bacaLapisan(null)).toEqual({
			tema: true,
			proyek: false,
			sekolah: false,
			laporan: false
		});
	});

	it('membaca dan menulis daftar lapisan secara bolak-balik', () => {
		const status = bacaLapisan('proyek,sekolah,laporan');
		expect(status).toEqual({ tema: false, proyek: true, sekolah: true, laporan: true });
		expect(tulisLapisan(status)).toBe('proyek,sekolah,laporan');
	});
});

describe('radiusTampilanM', () => {
	it('mengambil setengah tinggi tampilan dan membatasinya 10 km', () => {
		expect(radiusTampilanM(-6.19, -6.21)).toBe(1113);
		expect(radiusTampilanM(-5.9, -6.5)).toBe(10_000);
	});
});

describe('status peta dasar di URL', () => {
	it('memakai peta dasar standar untuk nilai kosong atau tak dikenal', () => {
		expect(bacaPetaDasar(null)).toBe('standar');
		expect(bacaPetaDasar('luar-angkasa')).toBe('standar');
		expect(bacaPetaDasar('satelit')).toBe('satelit');
	});

	it('menulis lapisan dan peta dasar tanpa membuang parameter lain', () => {
		const url = urlDenganStatusPeta(
			new URL('https://rantai.uji/id?wilayah=31.72.04&dasar=jalan'),
			bacaLapisan('tema,sekolah'),
			'standar'
		);
		expect(url.search).toBe('?wilayah=31.72.04&lapisan=tema%2Csekolah');
	});

	it('meneruskan jenjang, lapisan, dan peta dasar ke peta penuh', () => {
		const url = new URL(
			'https://rantai.uji/id?wilayah=31.72.04&tema=pendidikan_sma_smk&dasar=satelit'
		);
		expect(tautanPetaPenuh(url, 'id')).toBe('/id/peta?tema=pendidikan_sma_smk&dasar=satelit');
	});
});
