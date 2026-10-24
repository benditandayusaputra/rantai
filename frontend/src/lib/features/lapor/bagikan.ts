export type HasilBagikan = 'dibagikan' | 'disalin' | 'batal';

export async function bagikanTautan(teks: string, alamat: string): Promise<HasilBagikan> {
	if (navigator.share) {
		try {
			await navigator.share({ text: teks, url: alamat });
			return 'dibagikan';
		} catch {
			return 'batal';
		}
	}
	await navigator.clipboard.writeText(`${teks} ${alamat}`);
	return 'disalin';
}
