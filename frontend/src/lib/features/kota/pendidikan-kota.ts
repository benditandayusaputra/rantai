import type { BentukSekolah, PendidikanKota } from '$lib/api/tipe';
import type { TemaPendidikan } from '$lib/features/pendidikan/jenjang';

export type CakupanSekolah = 'negeri' | 'swasta' | 'semua';
export const bentukSekolah = ['SD', 'SMP', 'SMA', 'SMK'] as const;

export function bentukTujuan(tema: TemaPendidikan): BentukSekolah[] {
	return tema === 'pendidikan' ? ['SMP'] : ['SMA', 'SMK'];
}

export function jumlahSekolah(
	data: PendidikanKota,
	kode: string,
	bentuk: readonly BentukSekolah[],
	cakupan: CakupanSekolah
): number {
	return data.sekolah
		.filter((s) => s.kode_kecamatan === kode && bentuk.includes(s.bentuk))
		.reduce(
			(total, s) =>
				total +
				(cakupan === 'swasta' ? 0 : s.negeri.nilai) +
				(cakupan === 'negeri' ? 0 : s.swasta.nilai),
			0
		);
}

export function susunSebaran(data: PendidikanKota, tema: TemaPendidikan, cakupan: CakupanSekolah) {
	return data.kecamatan
		.map((k) => ({
			kunci: k.kode,
			nama: k.nama,
			jumlah: jumlahSekolah(data, k.kode, bentukTujuan(tema), cakupan),
			nilai: 0
		}))
		.sort((a, b) => b.jumlah - a.jumlah || a.nama.localeCompare(b.nama));
}

export function rentangSebaran(data: ReturnType<typeof susunSebaran>) {
	if (data.length === 0) return null;
	const minimum = Math.min(...data.map((r) => r.jumlah));
	const maksimum = Math.max(...data.map((r) => r.jumlah));
	return {
		minimum,
		maksimum,
		terendah: data.filter((r) => r.jumlah === minimum).map((r) => r.nama),
		tertinggi: data.filter((r) => r.jumlah === maksimum).map((r) => r.nama)
	};
}

export function sumberSekolah(data: PendidikanKota) {
	return [...new Map(data.sekolah.map((s) => [s.negeri.provenance_id, s])).values()];
}
