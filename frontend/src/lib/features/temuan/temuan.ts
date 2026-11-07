export const panjangMinIsi = 20;
export const panjangMaksIsi = 1000;
export const panjangMaksTautan = 500;

export function isTautanSah(tautan: string): boolean {
	if (tautan.length > panjangMaksTautan) return false;
	try {
		const url = new URL(tautan);
		return (url.protocol === 'https:' || url.protocol === 'http:') && url.hostname !== '';
	} catch {
		return false;
	}
}

export function isIsiSah(isi: string): boolean {
	const panjang = isi.trim().length;
	return panjang >= panjangMinIsi && panjang <= panjangMaksIsi;
}
