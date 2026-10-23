import { env } from '$env/dynamic/public';
import type {
	Amplop,
	DetailPendidikan,
	PetaPaket,
	PetaTema,
	Provenance,
	RingkasanKecamatan,
	Sekolah,
	SekitarTitik,
	Titik,
	WilayahDitemukan
} from './tipe';

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

export function ambilDetailPendidikan(fetchFn: Fetch, kode: string): Promise<DetailPendidikan> {
	return ambilData<DetailPendidikan>(
		fetchFn,
		`/tema/pendidikan/kecamatan/${encodeURIComponent(kode)}`
	);
}

export function ambilSekolahKecamatan(fetchFn: Fetch, kode: string): Promise<Sekolah[]> {
	return ambilData<Sekolah[]>(
		fetchFn,
		`/tema/pendidikan/kecamatan/${encodeURIComponent(kode)}/sekolah`
	);
}

export type Hasil<T> = { ok: true; data: T } | { ok: false; galat: GalatAPI };

function keGalatAPI(galat: unknown): GalatAPI {
	if (galat instanceof GalatAPI) return galat;
	return new GalatAPI(0, galat instanceof Error ? galat.message : String(galat));
}

export function tangkap<T>(permintaan: Promise<T>): Promise<Hasil<T>> {
	return permintaan.then(
		(data) => ({ ok: true, data }),
		(galat: unknown) => ({ ok: false, galat: keGalatAPI(galat) })
	);
}

export function ambilPetaTema(fetchFn: Fetch, kodeTema: string): Promise<PetaTema> {
	return ambilData<PetaTema>(fetchFn, `/peta/tema/${encodeURIComponent(kodeTema)}`);
}

export function ambilPetaPaket(fetchFn: Fetch, bbox: string): Promise<PetaPaket> {
	return ambilData<PetaPaket>(fetchFn, '/peta/paket', { bbox });
}

export function ambilSekolahTerdekat(
	fetchFn: Fetch,
	titik: Titik,
	radiusM: number
): Promise<SekitarTitik> {
	return ambilData<SekitarTitik>(fetchFn, '/sekolah/terdekat', {
		lat: titik.lat,
		lng: titik.lng,
		radius_m: radiusM
	});
}
