import { describe, expect, it } from 'vitest';
import type { GeometriPeta } from '$lib/api/tipe';
import { formatBulan } from '$lib/format';
import {
	ambangKuantil,
	angkaDari,
	bacaLensa,
	bacaTahun,
	namaWilayahPendek,
	batasRapi,
	rataPerBulanTercatat,
	susunGaris,
	susunPeta,
	tanpaNilai,
	tautanLensa,
	tautanTahun,
	tingkatKuantil,
	ukuranGaris,
	ukuranPeta,
	type FiturWilayah
} from './dasbor';

function kotak(kode: string, barat: number, selatan: number, sisi = 0.1): FiturWilayah {
	const cincin = [
		[barat, selatan],
		[barat + sisi, selatan],
		[barat + sisi, selatan + sisi],
		[barat, selatan + sisi],
		[barat, selatan]
	];
	const geometry: GeometriPeta = { type: 'Polygon', coordinates: [cincin] };
	return { geometry, properties: { kode, nama: kode } };
}

describe('tautan dasbor', () => {
	const url = new URL('https://rantai.test/id/dasbor?lensa=uang&tahun=2025');

	it('membaca lensa dan tahun yang sah saja', () => {
		expect(bacaLensa('tanggung-jawab')).toBe('tanggung-jawab');
		expect(bacaLensa('lain')).toBe('layanan');
		expect(bacaTahun('2025')).toBe(2025);
		expect(bacaTahun('25')).toBeUndefined();
		expect(bacaTahun(null)).toBeUndefined();
	});

	it('mengganti lensa sambil melepas tahun', () => {
		expect(tautanLensa(url, 'layanan')).toBe('/id/dasbor');
		expect(tautanLensa(url, 'tanggung-jawab')).toBe('/id/dasbor?lensa=tanggung-jawab');
		expect(tautanTahun(url, 2026)).toBe('/id/dasbor?lensa=uang&tahun=2026');
	});
});

describe('angka dasbor', () => {
	it('tanpa provenance tidak menghasilkan angka', () => {
		expect(angkaDari({ label: '', provenance_id: '' }, 5, 'paket')).toBeNull();
		expect(angkaDari({ label: 'resmi', provenance_id: 'prv' }, 5, 'paket')).toEqual({
			nilai: 5,
			satuan: 'paket',
			label: 'resmi',
			provenance_id: 'prv'
		});
	});

	it('merata-rata hanya bulan yang tercatat', () => {
		const bulanan = [0, 0, 30, 10].map((jumlah, i) => ({ kunci: String(i), jumlah, nilai: 0 }));
		expect(rataPerBulanTercatat(bulanan)).toBe(20);
		expect(rataPerBulanTercatat([])).toBe(0);
	});

	it.each([
		[0, 1],
		[7, 8],
		[180, 200],
		[384, 400],
		[10_884, 12_000],
		[41_200_000_000, 50_000_000_000]
	])('membulatkan batas sumbu %d ke %d', (maksimum, batas) => {
		expect(batasRapi(maksimum)).toBe(batas);
	});

	it('memberi tingkat kuantil dan memisahkan nilai nol', () => {
		const ambang = ambangKuantil([0, 10, 20, 30, 40, 50], 5);
		expect(ambang).toEqual([20, 30, 40, 50]);
		expect(tingkatKuantil(0, ambang)).toBe(tanpaNilai);
		expect(tingkatKuantil(10, ambang)).toBe(0);
		expect(tingkatKuantil(50, ambang)).toBe(4);
		expect(ambangKuantil([0, 0], 5)).toEqual([]);
		expect(ambangKuantil([5, 5, 5, 9], 4)).toEqual([5, 9]);
	});

	it('menulis nama bulan pendek sesuai bahasa', () => {
		expect(formatBulan('2025-03', 'id')).toBe('Mar');
		expect(formatBulan('2025-08', 'id')).toBe('Agu');
		expect(formatBulan('2025-08', 'en')).toBe('Aug');
	});
});

describe('susunPeta', () => {
	it('memproyeksikan daratan ke kanan petak Kepulauan Seribu', () => {
		const susunan = susunPeta([
			kotak('31.71.01', 106.8, -6.2),
			kotak('31.72.01', 106.9, -6.2),
			kotak('31.01.01', 106.5, -5.6, 0.05)
		]);

		expect(susunan.sisipan).not.toBeNull();
		const daratan = susunan.area.find((a) => a.kode === '31.71.01');
		const seribu = susunan.area.find((a) => a.kode === '31.01.01');
		expect(daratan?.d.startsWith(`M${ukuranPeta.lebarSisipan + ukuranPeta.jarak}.0,`)).toBe(true);
		expect(seribu?.tengah[0]).toBeLessThan(ukuranPeta.lebarSisipan);
		expect(seribu?.d).toBe('M8,8h80v64h-80Z');
		expect(susunan.lebar).toBe(ukuranPeta.lebarSisipan + ukuranPeta.jarak + ukuranPeta.lebarUtama);
		const faktorLintang = Math.cos((6.15 * Math.PI) / 180);
		expect(susunan.tinggi).toBeCloseTo(ukuranPeta.lebarUtama / 2 / faktorLintang, 3);
	});

	it('tanpa Kepulauan Seribu tidak ada sisipan', () => {
		const susunan = susunPeta([kotak('31.71.01', 106.8, -6.2)]);

		expect(susunan.sisipan).toBeNull();
		expect(susunan.area[0]?.d).toMatch(/^M0\.0,/);
	});
});

describe('teks dasbor', () => {
	it('memendekkan nama kota administrasi', () => {
		expect(namaWilayahPendek('Kota Administrasi Jakarta Pusat')).toBe('Jakarta Pusat');
		expect(namaWilayahPendek('Kabupaten Administrasi Kepulauan Seribu')).toBe('Kepulauan Seribu');
		expect(namaWilayahPendek('DKI Jakarta')).toBe('DKI Jakarta');
	});
});

describe('susunGaris', () => {
	const { lebar, tinggi, kiri, kanan, bawah } = ukuranGaris;

	it('membentang dari tepi kiri ke tepi kanan area grafik', () => {
		const titik = susunGaris([10, 20, 30]);
		expect(titik.map((p) => p.x)).toEqual([kiri, lebar / 2, lebar - kanan]);
	});

	it('menaruh nilai nol di garis dasar dan nilai lebih besar lebih tinggi', () => {
		const [nol, kecil, besar] = susunGaris([0, 5, 10]).map((p) => p.y);
		expect(nol).toBe(tinggi - bawah);
		expect(besar).toBeLessThan(kecil ?? 0);
	});

	it('menaruh satu titik di tengah', () => {
		expect(susunGaris([7])[0]?.x).toBe(lebar / 2);
	});

	it('mengembalikan daftar kosong tanpa nilai', () => {
		expect(susunGaris([])).toEqual([]);
	});
});
