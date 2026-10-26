import { describe, expect, it } from 'vitest';
import {
	daftarPenanda,
	ekspresiPenandaLaporan,
	ekspresiPenandaProyek,
	namaPenandaProyek,
	penandaSekolah,
	svgPenanda
} from './penanda';

const token = (nama: string) => `var-${nama}`;

describe('penanda peta', () => {
	it('menyiapkan satu gambar untuk sekolah, tiap sektor proyek, dan tiap status laporan', () => {
		const nama = daftarPenanda(token).map((p) => p.nama);
		expect(nama).toEqual([
			penandaSekolah,
			'penanda-proyek-pendidikan',
			'penanda-proyek-air',
			'penanda-proyek-jalan',
			'penanda-proyek-lain',
			'penanda-laporan-tuntas',
			'penanda-laporan-proses',
			'penanda-laporan-simulasi'
		]);
	});

	it('mewarnai penanda proyek dengan token sektornya', () => {
		const air = daftarPenanda(token).find((p) => p.nama === namaPenandaProyek('air'));
		expect(air).toMatchObject({ warna: 'var-sector-air', warnaIkon: 'var-sector-air-on' });
	});

	it('memakai penanda lainnya untuk sektor kosong atau tak dikenal', () => {
		expect(namaPenandaProyek('')).toBe('penanda-proyek-lain');
		expect(ekspresiPenandaProyek.at(-1)).toBe('penanda-proyek-lain');
	});

	it('mendahulukan penanda simulasi untuk laporan simulasi', () => {
		expect(ekspresiPenandaLaporan.slice(0, 3)).toEqual([
			'case',
			['get', 'is_simulation'],
			'penanda-laporan-simulasi'
		]);
	});

	it('membungkus ikon dalam pin berwarna yang tajam di layar rapat', () => {
		const svg = svgPenanda('pin', 'var-sector-edu', 'var-surface', '<svg></svg>');
		expect(svg).toContain('width="60" height="76" viewBox="0 0 30 38"');
		expect(svg).toContain('fill="var-sector-edu" stroke="var-surface"');
		expect(svg).toContain('<g transform="translate(7 6.5)"><svg></svg></g>');
	});
});
