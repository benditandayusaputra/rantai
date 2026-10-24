export type Pintasan = 'berikut' | 'sebelum' | 'verifikasi' | 'tolak';

const tombolPintasan: Record<string, Pintasan> = {
	j: 'berikut',
	k: 'sebelum',
	v: 'verifikasi',
	r: 'tolak'
};

export function bacaPintasan(event: KeyboardEvent): Pintasan | null {
	const target = event.target;
	const isMengetik =
		target instanceof HTMLElement && ['INPUT', 'TEXTAREA', 'SELECT'].includes(target.tagName);
	if (isMengetik || event.ctrlKey || event.metaKey || event.altKey) return null;
	return tombolPintasan[event.key.toLowerCase()] ?? null;
}

export function geser(indeks: number, arah: 1 | -1, jumlah: number): number {
	if (jumlah === 0) return 0;
	return Math.min(jumlah - 1, Math.max(0, indeks + arah));
}

export function sisihkan<T extends { id: number }>(daftar: T[], id: number): T[] {
	return daftar.filter((item) => item.id !== id);
}
