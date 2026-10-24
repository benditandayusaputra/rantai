import { env } from '$env/dynamic/public';
import type {
	Amplop,
	DaftarDataBelumDibuka,
	DetailPaket,
	DetailPendidikan,
	DetailWilayah,
	Dokumen,
	InfoSumberData,
	PaketRingkas,
	PermintaanAspirasi,
	PetaPaket,
	PetaTema,
	ProfilPenyedia,
	Provenance,
	RingkasanKecamatan,
	Sekolah,
	SekitarTitik,
	TemaDetail,
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

export function susunURL(path: string, kueri: Kueri): string {
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

export async function kirimData<T>(
	fetchFn: Fetch,
	path: string,
	badan: unknown,
	opsi: RequestInit = {}
): Promise<T> {
	const respons = await fetchFn(susunURL(path, {}), {
		...opsi,
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(badan)
	});
	const amplop = await bacaAmplop<T>(respons);
	if (!respons.ok || !amplop) {
		throw new GalatAPI(respons.status, amplop?.message ?? respons.statusText);
	}
	return amplop.data;
}

export async function kirimBerwenang<T>(
	fetchFn: Fetch,
	path: string,
	token: string,
	badan: FormData | object,
	metode: 'POST' | 'PATCH' = 'POST'
): Promise<T> {
	const isFormulir = badan instanceof FormData;
	const respons = await fetchFn(susunURL(path, {}), {
		method: metode,
		headers: {
			Authorization: `Bearer ${token}`,
			...(isFormulir ? {} : { 'Content-Type': 'application/json' })
		},
		body: isFormulir ? badan : JSON.stringify(badan)
	});
	const amplop = await bacaAmplop<T>(respons);
	if (!respons.ok || !amplop) {
		throw new GalatAPI(respons.status, amplop?.message ?? respons.statusText);
	}
	return amplop.data;
}

export function urlServer(path: string): string {
	return new URL(path, env.PUBLIC_API_URL || alamatBawaan).toString();
}

export async function ambilBerwenang<T>(fetchFn: Fetch, path: string, token: string): Promise<T> {
	const respons = await fetchFn(susunURL(path, {}), {
		headers: { Authorization: `Bearer ${token}` }
	});
	const amplop = await bacaAmplop<T>(respons);
	if (!respons.ok || !amplop) {
		throw new GalatAPI(respons.status, amplop?.message ?? respons.statusText);
	}
	return amplop.data;
}

export async function ambilData<T>(fetchFn: Fetch, path: string, kueri: Kueri = {}): Promise<T> {
	return (await ambilAmplop<T>(fetchFn, path, kueri)).data;
}

export function ambilProvenance(fetchFn: Fetch, id: string): Promise<Provenance> {
	return ambilData<Provenance>(fetchFn, `/provenance/${encodeURIComponent(id)}`);
}

export function ambilWilayah(fetchFn: Fetch, kode: string): Promise<WilayahDitemukan> {
	return ambilData<WilayahDitemukan>(fetchFn, `/wilayah/${encodeURIComponent(kode)}`);
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

export interface FilterPaket {
	tema?: string;
	tahun?: number;
	halaman?: number;
}

export function ambilPaketKecamatan(
	fetchFn: Fetch,
	kode: string,
	filter: FilterPaket
): Promise<Amplop<PaketRingkas[]>> {
	return ambilAmplop<PaketRingkas[]>(fetchFn, `/kecamatan/${encodeURIComponent(kode)}/paket`, {
		tema: filter.tema,
		tahun: filter.tahun,
		halaman: filter.halaman
	});
}

export function ambilPaket(fetchFn: Fetch, kode: string): Promise<DetailPaket> {
	return ambilData<DetailPaket>(fetchFn, `/paket/${encodeURIComponent(kode)}`);
}

export function ambilSumberData(fetchFn: Fetch): Promise<InfoSumberData> {
	return ambilData<InfoSumberData>(fetchFn, '/sumber-data');
}

export function ambilDaftarTema(fetchFn: Fetch): Promise<TemaDetail[]> {
	return ambilData<TemaDetail[]>(fetchFn, '/tema');
}

export function ambilTema(fetchFn: Fetch, kode: string): Promise<TemaDetail> {
	return ambilData<TemaDetail>(fetchFn, `/tema/${encodeURIComponent(kode)}`);
}

export function susunAspirasi(fetchFn: Fetch, p: PermintaanAspirasi): Promise<Dokumen> {
	return kirimData<Dokumen>(fetchFn, '/aspirasi', p);
}

export function ambilDataBelumDibuka(fetchFn: Fetch): Promise<DaftarDataBelumDibuka> {
	return ambilData<DaftarDataBelumDibuka>(fetchFn, '/data-belum-dibuka');
}

export function susunPermohonan(fetchFn: Fetch, kode: string, bahasa: string): Promise<Dokumen> {
	return kirimData<Dokumen>(fetchFn, `/data-belum-dibuka/${encodeURIComponent(kode)}/permohonan`, {
		bahasa
	});
}

export function ambilPenyedia(fetchFn: Fetch, id: number): Promise<ProfilPenyedia> {
	return ambilData<ProfilPenyedia>(fetchFn, `/penyedia/${id}`);
}

export function ambilDetailWilayah(fetchFn: Fetch, kode: string): Promise<DetailWilayah> {
	return ambilData<DetailWilayah>(fetchFn, `/wilayah/${encodeURIComponent(kode)}`);
}
