import { describe, expect, it } from 'vitest';
import type { MataRantai } from '$lib/api/tipe';
import { daftarTahun, kodeKecamatanDari, ringkasRantai } from './paket';

const mata = (status: MataRantai['status']): MataRantai => ({ kode: 'rencana', status });

describe('ringkasRantai', () => {
	it('menyebut mata rantai putus bila ada', () => {
		const rantai = [mata('ada'), mata('ada'), mata('putus'), mata('menunggu')];
		expect(ringkasRantai(rantai)).toBe('2 dari 4 tersambung, 1 putus');
	});

	it('memakai kalimat ringkas bila tidak ada yang putus', () => {
		expect(ringkasRantai([mata('ada'), mata('estimasi')])).toBe('1 dari 2 mata rantai tersambung');
	});
});

describe('daftarTahun', () => {
	it('menurun dari tahun berjalan', () => {
		expect(daftarTahun(2026, 3)).toEqual([2026, 2025, 2024]);
	});
});

describe('kodeKecamatanDari', () => {
	it('memotong kode kelurahan menjadi kode kecamatan', () => {
		expect(kodeKecamatanDari('31.72.06.1001')).toBe('31.72.06');
		expect(kodeKecamatanDari('31.72.06')).toBe('31.72.06');
	});
});
