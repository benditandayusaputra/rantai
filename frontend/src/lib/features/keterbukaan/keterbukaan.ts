import { ambilDataBelumDibuka, susunPermohonan, type Fetch } from '$lib/api/klien';
import type { DaftarDataBelumDibuka, Dokumen } from '$lib/api/tipe';
import type { BahasaDokumen } from '$lib/features/aspirasi/aspirasi';

export interface PilihanPermohonan {
	kode: string;
	bahasa: BahasaDokumen;
}

export interface HalamanKeterbukaan {
	daftar: DaftarDataBelumDibuka;
	dokumen: Dokumen | null;
}

export function bacaPilihanPermohonan(kueri: URLSearchParams): PilihanPermohonan {
	return { kode: kueri.get('minta') ?? '', bahasa: kueri.get('bahasa') === 'en' ? 'en' : 'id' };
}

export function tautanPermohonan(url: URL, ubah: Partial<PilihanPermohonan>): string {
	const pilihan = { ...bacaPilihanPermohonan(url.searchParams), ...ubah };
	const kueri = new URLSearchParams({ minta: pilihan.kode, bahasa: pilihan.bahasa });
	return `${url.pathname}?${kueri.toString()}#surat`;
}

export async function muatKeterbukaan(
	fetchFn: Fetch,
	p: PilihanPermohonan
): Promise<HalamanKeterbukaan> {
	const [daftar, dokumen] = await Promise.all([
		ambilDataBelumDibuka(fetchFn),
		p.kode ? susunPermohonan(fetchFn, p.kode, p.bahasa) : Promise.resolve(null)
	]);
	return { daftar, dokumen };
}
