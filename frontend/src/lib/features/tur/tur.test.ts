import { afterEach, describe, expect, it } from 'vitest';
import {
	isTurSelesai,
	jarakKartu,
	kunciTurSelesai,
	langkahTersedia,
	langkahTur,
	posisiKartu,
	tandaiTurSelesai,
	tepiLayar
} from './tur';

const layar = { lebar: 1280, tinggi: 800 };
const kartu = { lebar: 320, tinggi: 200 };

describe('posisiKartu', () => {
	it('menaruh kartu di bawah target bila muat', () => {
		const posisi = posisiKartu({ x: 400, y: 20, lebar: 80, tinggi: 40 }, kartu, layar);
		expect(posisi).toEqual({ x: 280, y: 60 + jarakKartu });
	});

	it('menaruh kartu di atas target yang dekat tepi bawah', () => {
		const posisi = posisiKartu({ x: 400, y: 740, lebar: 80, tinggi: 50 }, kartu, layar);
		expect(posisi.y).toBe(740 - jarakKartu - kartu.tinggi);
	});

	it('menjaga kartu tetap di dalam layar secara mendatar', () => {
		const kiri = posisiKartu({ x: 0, y: 20, lebar: 40, tinggi: 40 }, kartu, layar);
		const kanan = posisiKartu({ x: 1260, y: 20, lebar: 20, tinggi: 40 }, kartu, layar);
		expect(kiri.x).toBe(tepiLayar);
		expect(kanan.x).toBe(layar.lebar - tepiLayar - kartu.lebar);
	});
});

describe('langkahTersedia', () => {
	it('membuang langkah yang targetnya tidak ada di halaman', () => {
		const tersedia = langkahTersedia(langkahTur, (target) => target !== 'cari');
		expect(tersedia.map((l) => l.target)).not.toContain('cari');
		expect(tersedia[0]?.target).toBeNull();
	});
});

describe('penanda tur selesai', () => {
	afterEach(() => localStorage.clear());

	it('belum selesai pada kunjungan pertama dan selesai setelah ditandai', () => {
		expect(isTurSelesai()).toBe(false);
		tandaiTurSelesai();
		expect(isTurSelesai()).toBe(true);
		expect(localStorage.getItem(kunciTurSelesai)).toBe('1');
	});
});
