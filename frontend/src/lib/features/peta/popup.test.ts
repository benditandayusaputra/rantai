import { describe, expect, it } from 'vitest';
import { isiPopup } from './popup';

describe('isi popup penanda', () => {
	it('menautkan proyek ke halaman paketnya', () => {
		const isi = isiPopup(
			'proyek-titik',
			{ kode: 'P 1', nama: 'Rehab SDN', sektor: 'pendidikan' },
			'id'
		);
		expect(isi.judul).toBe('Rehab SDN');
		expect(isi.keterangan).toBe('Proyek · Pendidikan');
		expect(isi.tautan?.href).toBe('/id/paket/P%201');
	});

	it('menandai laporan simulasi dan menautkan ke laporannya', () => {
		const isi = isiPopup(
			'laporan-titik',
			{ id: 7, status: 'TERVERIFIKASI', kondisi: 'terhenti', is_simulation: true },
			'en'
		);
		expect(isi.keterangan).toContain('Simulasi');
		expect(isi.tautan?.href).toBe('/en/laporan/7');
	});

	it('menyebut bentuk dan status sekolah tanpa tautan', () => {
		const isi = isiPopup(
			'sekolah-titik',
			{ nama: 'SMAN 1', bentuk: 'SMA', status: 'swasta' },
			'id'
		);
		expect(isi).toEqual({ keterangan: 'SMA swasta', judul: 'SMAN 1' });
	});
});
