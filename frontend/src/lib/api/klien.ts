import { env } from '$env/dynamic/public';
import type { Amplop, Provenance, RingkasanKecamatan, WilayahDitemukan } from './tipe';

const alamatBawaan = 'http://localhost:8080/api/v1';

export type Fetch = typeof fetch;
export type Kueri = Record<string, string | number | boolean | undefined>;

export class GalatAPI extends Error {
	constructor(
		readonly kode: number,
		message: string
	) {
		super(message);
		this.name = 'GalatAPI';
	}
}

function susunURL(path: string, kueri: Kueri): string {
	const url = new URL((env.PUBLIC_API_URL || alamatBawaan) + path);
	for (const [nama, nilai] of Object.entries(kueri)) {
		if (nilai !== undefined && nilai !== '') url.searchParams.set(nama, String(nilai));
	}
	return url.toString();
}

async function bacaAmplop<T>(respons: Response): Promise<Amplop<T> | null> {
	try {
		return (await respons.json()) as Amplop<T>;
	} catch {
		return null;
	}
}

export async function ambilAmplop<T>(
	fetchFn: Fetch,
	path: string,
	kueri: Kueri = {}
): Promise<Amplop<T>> {
	const respons = await fetchFn(susunURL(path, kueri));
	const amplop = await bacaAmplop<T>(respons);
	if (!respons.ok || !amplop) {
		throw new GalatAPI(respons.status, amplop?.message ?? respons.statusText);
	}
	return amplop;
}

export async function ambilData<T>(fetchFn: Fetch, path: string, kueri: Kueri = {}): Promise<T> {
	return (await ambilAmplop<T>(fetchFn, path, kueri)).data;
}

export function ambilProvenance(fetchFn: Fetch, id: string): Promise<Provenance> {
	return ambilData<Provenance>(fetchFn, `/provenance/${encodeURIComponent(id)}`);
}

export function cariWilayah(fetchFn: Fetch, q: string): Promise<WilayahDitemukan[]> {
	return ambilData<WilayahDitemukan[]>(fetchFn, '/wilayah/cari', { q });
}

export function ambilWilayahDiTitik(
	fetchFn: Fetch,
	lat: number,
	lng: number
): Promise<WilayahDitemukan> {
	return ambilData<WilayahDitemukan>(fetchFn, '/wilayah/titik', { lat, lng });
}

export function ambilRingkasanKecamatan(fetchFn: Fetch, kode: string): Promise<RingkasanKecamatan> {
	return ambilData<RingkasanKecamatan>(fetchFn, `/kecamatan/${encodeURIComponent(kode)}/ringkasan`);
}
