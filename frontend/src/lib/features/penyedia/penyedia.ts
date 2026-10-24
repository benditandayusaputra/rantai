import type { TahunPenyedia } from '$lib/api/tipe';

export function tahunTerbaru(daftar: TahunPenyedia[]): TahunPenyedia | null {
	return daftar.reduce<TahunPenyedia | null>(
		(terbaru, item) => (!terbaru || item.tahun > terbaru.tahun ? item : terbaru),
		null
	);
}

export function tinggiBatang(daftar: TahunPenyedia[]): Map<number, number> {
	const maksimum = Math.max(0, ...daftar.map((item) => item.total_nilai.nilai));
	return new Map(
		daftar.map((item) => [item.tahun, maksimum > 0 ? item.total_nilai.nilai / maksimum : 0])
	);
}

export function inisialPenyedia(nama: string): string {
	const kata = nama
		.replace(/^(cv|pt|ud|koperasi)\.?\s+/i, '')
		.split(/\s+/)
		.filter(Boolean);
	return kata
		.slice(0, 2)
		.map((k) => k[0]?.toUpperCase() ?? '')
		.join('');
}
