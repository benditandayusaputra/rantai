import { describe, expect, it } from 'vitest';
import type { Notifikasi } from '$lib/api/tipe';
import { kelompokkanPerHari } from './pantau';

function kabar(id: number, terjadi_pada: string): Notifikasi {
	return {
		id,
		jenis: 'paket_baru',
		kode_kecamatan: '31.72.06',
		nama_kecamatan: 'Kelapa Gading',
		judul: 'Paket',
		tautan: '/paket/x',
		terjadi_pada,
		dibaca: false,
		is_simulation: false
	};
}

describe('kelompokkanPerHari', () => {
	it('memakai label Hari ini, Kemarin, lalu tanggal menurut waktu Jakarta', () => {
		const sekarang = new Date('2026-10-24T10:00:00+07:00');
		const kelompok = kelompokkanPerHari(
			[
				kabar(1, '2026-10-24T08:00:00+07:00'),
				kabar(2, '2026-10-24T00:30:00+07:00'),
				kabar(3, '2026-10-23T23:00:00+07:00'),
				kabar(4, '2026-10-20T09:00:00+07:00')
			],
			sekarang
		);
		expect(kelompok.map((k) => [k.hari, k.item.length])).toEqual([
			['Hari ini', 2],
			['Kemarin', 1],
			['20 Oktober 2026', 1]
		]);
	});
});
