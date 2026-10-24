import { describe, expect, it } from 'vitest';
import type { AntreanPanel } from '$lib/api/tipe';
import {
	bagianUntuk,
	daftarSorotan,
	isAkunPanel,
	isSiapTerbit,
	jumlahSegeraBerakhir,
	tanggalLokal
} from './panel';

const antrean = {
	peran: 'instansi',
	nama_instansi: 'Dinas Uji',
	sinyal: [
		{
			id: 4,
			jenis: 'S-01',
			alasan: 'Satu peserta',
			kode_paket: 'tender-1',
			nama_paket: 'Rehab SD',
			jumlah_hak_jawab: 0
		}
	],
	laporan: [],
	paket: [
		{ kode: 'a', segera_berakhir: true },
		{ kode: 'b', segera_berakhir: false }
	],
	umpan: []
} as unknown as AntreanPanel;

describe('panel', () => {
	it('hanya instansi dan anggota yang memakai panel', () => {
		expect(isAkunPanel('instansi')).toBe(true);
		expect(isAkunPanel('anggota')).toBe(true);
		expect(isAkunPanel('warga')).toBe(false);
		expect(isAkunPanel(undefined)).toBe(false);
	});

	it('anggota hanya melihat umpan, instansi melihat tiga antrean', () => {
		expect(bagianUntuk('anggota')).toEqual(['umpan']);
		expect(bagianUntuk('instansi')).toEqual(['sinyal', 'laporan', 'paket']);
	});

	it('sorotan sinyal membawa id sinyal dan alasan sebagai catatan', () => {
		expect(daftarSorotan(antrean, 'sinyal', 'id')).toEqual([
			{
				jenis: 'sinyal',
				id: '4',
				judul: 'Rehab SD',
				catatan: 'Satu peserta',
				keterangan: 'Hanya satu peserta · 0 hak jawab',
				tautanPublik: '/id/paket/tender-1',
				isPenting: true,
				isSimulation: false
			}
		]);
		expect(jumlahSegeraBerakhir(antrean)).toBe(1);
	});

	it('hak jawab siap terbit bila berisi, dalam batas, dan kewenangan dicentang', () => {
		expect(isSiapTerbit('Tanggapan', true)).toBe(true);
		expect(isSiapTerbit('   ', true)).toBe(false);
		expect(isSiapTerbit('Tanggapan', false)).toBe(false);
		expect(isSiapTerbit('x'.repeat(2001), true)).toBe(false);
	});

	it('tanggal lokal untuk batas isian tanggal', () => {
		expect(tanggalLokal(new Date(2026, 6, 3))).toBe('2026-07-03');
	});
});
