import { describe, expect, it } from 'vitest';
import type { Angka, DetailPendidikan, Sekolah } from '$lib/api/tipe';
import {
	kunciJenisSekolah,
	persenLebar,
	skalaTren,
	susunBarisGrafik,
	susunRumus,
	ukuranGrafikTren
} from './pendidikan';

const angka = (nilai: number, label: Angka['label'] = 'resmi'): Angka => ({
	nilai,
	satuan: 'kursi',
	label,
	provenance_id: 'p'
});

const detail = {
	calon_lulusan: angka(2701),
	kursi_negeri: angka(1284),
	kursi_spmb_bersama: null,
	kursi_madrasah: angka(301, 'proksi')
} as DetailPendidikan;

describe('susunBarisGrafik', () => {
	it('melewati jenis kursi yang tidak tersedia dan mempertahankan labelnya', () => {
		const baris = susunBarisGrafik(detail, 'pendidikan');
		expect(baris.map((b) => b.kunci)).toEqual([
			'pen.graduates',
			'pen.seats.negeri',
			'pen.seats.madrasah'
		]);
		expect(baris[2]?.angka.label).toBe('proksi');
	});
});

const ambang = { sangat_kurang: 0.5, kurang: 0.8, hampir_cukup: 1 };

const detailSMASMK = {
	calon_lulusan: angka(4460),
	kursi_negeri: angka(576),
	kursi_spmb_bersama: angka(40),
	kursi_madrasah: null,
	kesenjangan: angka(3884),
	kesenjangan_dengan_spmb_bersama: angka(3844),
	rasio: { nilai: 0.129, satuan: 'rasio', label: 'resmi', provenance_id: 'p' },
	kursi_negeri_per_bentuk: [
		{ bentuk: 'SMA', kursi: angka(576) },
		{ bentuk: 'SMK', kursi: angka(0) }
	],
	kursi_spmb_bersama_per_bentuk: [
		{ bentuk: 'SMA', kursi: angka(0) },
		{ bentuk: 'SMK', kursi: angka(40) }
	],
	rumus: {
		bentuk_asal: 'SMP',
		tingkat_asal: 9,
		bentuk_tujuan: ['SMA', 'SMK'],
		bentuk_madrasah_tujuan: 'MA',
		tingkat_tujuan: 10,
		batas_murid_per_rombel: 36,
		ambang_rasio: ambang
	}
} as DetailPendidikan;

describe('susunBarisGrafik SMA/SMK', () => {
	it('memecah kursi negeri dan PMB Bersama per bentuk SMA dan SMK', () => {
		const baris = susunBarisGrafik(detailSMASMK, 'pendidikan_sma_smk');
		expect(baris.map((b) => b.kunci)).toEqual([
			'pen.sma.graduates',
			'pen.seats.negeri.SMA',
			'pen.seats.negeri.SMK',
			'pen.seats.spmb.SMA',
			'pen.seats.spmb.SMK'
		]);
		expect(baris.filter((b) => b.isKebutuhan)).toHaveLength(1);
	});
});

describe('susunRumus', () => {
	it('menampilkan rumus dengan angka kecamatan, termasuk penjumlahan kursi SMA dan SMK', () => {
		const baris = susunRumus(detailSMASMK);
		const perKunci = Object.fromEntries(baris.map((b) => [b.kunci, b]));
		expect(perKunci['pen.formula.graduates']?.parameter).toEqual({ tingkat: 9, bentuk: 'SMP' });
		expect(perKunci['pen.formula.seats']?.hitungan).toBe('576 + 0 = 576');
		expect(perKunci['pen.formula.gap']?.hitungan).toBe('4.460 − 576 = 3.884');
		expect(perKunci['pen.formula.ratio']?.hitungan).toBe('576 ÷ 4.460 = 0,13');
		expect(perKunci['pen.formula.gapSpmb']?.hitungan).toBe('3.884 − 40 = 3.844');
	});

	it('menandai kursi madrasah tidak tersedia tanpa menebak angka', () => {
		const madrasah = susunRumus(detailSMASMK).find((b) => b.kunci === 'pen.formula.madrasah');
		expect(madrasah?.parameter).toEqual({ madrasah: 'MA', tingkat: 10, batas: 36 });
		expect(madrasah?.hitungan).toBeNull();
	});
});

describe('persenLebar', () => {
	it('menskalakan terhadap nilai terbesar tanpa membagi nol', () => {
		expect(persenLebar(1284, 2701)).toBe(48);
		expect(persenLebar(0, 0)).toBe(0);
	});
});

describe('kunciJenisSekolah', () => {
	const dasar = { npsn: '1', nama: 'x', lokasi: null, daya_tampung: null } as const;
	it.each([
		[{ ...dasar, bentuk: 'SMP', status: 'negeri', jalur: 'negeri' }, 'pen.school.negeri'],
		[{ ...dasar, bentuk: 'SMP', status: 'swasta', jalur: 'spmb_bersama' }, 'pen.school.spmb'],
		[{ ...dasar, bentuk: 'SMP', status: 'swasta', jalur: null }, 'pen.school.swasta'],
		[{ ...dasar, bentuk: 'MTs', status: 'negeri', jalur: null }, 'pen.school.madrasah'],
		[{ ...dasar, bentuk: 'SMA', status: 'negeri', jalur: 'negeri' }, 'pen.school.sma.negeri'],
		[{ ...dasar, bentuk: 'SMK', status: 'swasta', jalur: 'spmb_bersama' }, 'pen.school.smk.spmb'],
		[{ ...dasar, bentuk: 'SMK', status: 'swasta', jalur: null }, 'pen.school.smk.swasta'],
		[{ ...dasar, bentuk: 'MA', status: 'negeri', jalur: null }, 'pen.school.ma']
	] as [Sekolah, string][])('%o menjadi %s', (sekolah, harapan) => {
		expect(kunciJenisSekolah(sekolah)).toBe(harapan);
	});
});

describe('skalaTren', () => {
	it('merentangkan titik dari tepi kiri ke tepi kanan', () => {
		const { titik } = skalaTren([1.2, 1.5, 1.1], (nilai) => nilai);
		expect(titik.at(0)?.x).toBe(ukuranGrafikTren.kiri);
		expect(titik.at(-1)?.x).toBe(ukuranGrafikTren.lebar - ukuranGrafikTren.kanan);
	});

	it('menaruh rasio lebih besar lebih tinggi di grafik', () => {
		const [rendah, tinggi] = skalaTren([1.2, 1.8], (nilai) => nilai).titik;
		expect(tinggi?.y).toBeLessThan(rendah?.y ?? 0);
	});

	it('selalu memuat garis satu lulusan per kursi walau semua rasio di bawahnya', () => {
		const { titik, ySeimbang } = skalaTren([0.4, 0.6], (nilai) => nilai);
		expect(ySeimbang).toBeGreaterThan(ukuranGrafikTren.atas);
		expect(titik.every((p) => p.y > ySeimbang)).toBe(true);
	});
});
