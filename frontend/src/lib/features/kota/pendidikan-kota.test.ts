import { describe, expect, it } from 'vitest';
import type { PendidikanKota } from '$lib/api/tipe';
import { jumlahSekolah, rentangSebaran, susunSebaran } from './pendidikan-kota';

const angka = (nilai: number) => ({
	nilai,
	satuan: 'sekolah',
	label: 'resmi' as const,
	provenance_id: 'dapodik'
});
const data = {
	kecamatan: [
		{ kode: 'a', nama: 'A' },
		{ kode: 'b', nama: 'B' },
		{ kode: 'c', nama: 'C' }
	],
	sekolah: [
		{ kode_kecamatan: 'a', bentuk: 'SMA', negeri: angka(2), swasta: angka(3) },
		{ kode_kecamatan: 'a', bentuk: 'SMK', negeri: angka(1), swasta: angka(5) },
		{ kode_kecamatan: 'b', bentuk: 'SMA', negeri: angka(1), swasta: angka(1) },
		{ kode_kecamatan: 'b', bentuk: 'SMP', negeri: angka(9), swasta: angka(9) }
	]
} as PendidikanKota;

describe('sebaran kabupaten/kota', () => {
	it('memisahkan status, menggabungkan SMA/SMK, dan tidak mencampur jenjang', () => {
		expect(jumlahSekolah(data, 'a', ['SMA', 'SMK'], 'negeri')).toBe(3);
		expect(jumlahSekolah(data, 'a', ['SMA', 'SMK'], 'swasta')).toBe(8);
		expect(jumlahSekolah(data, 'a', ['SMA', 'SMK'], 'semua')).toBe(11);
		expect(susunSebaran(data, 'pendidikan_sma_smk', 'negeri').map((r) => r.jumlah)).toEqual([
			3, 1, 0
		]);
	});

	it('tetap menampilkan kecamatan tanpa catatan sekolah dan semua nama saat seri', () => {
		const sebaran = susunSebaran(data, 'pendidikan', 'negeri');
		expect(sebaran.map((r) => r.kunci)).toEqual(['b', 'a', 'c']);
		expect(rentangSebaran(sebaran)).toEqual({
			minimum: 0,
			maksimum: 9,
			terendah: ['A', 'C'],
			tertinggi: ['B']
		});
	});

	it('tidak membuat kesimpulan untuk daftar kecamatan kosong', () => {
		expect(rentangSebaran([])).toBeNull();
	});
});
