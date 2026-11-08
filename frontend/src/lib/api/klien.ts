import { env } from '$env/dynamic/public';
import type {
	Amplop,
	AspirasiReses,
	IsuDisuarakan,
	ItemUmpan,
	JenisSumberUmpan,
	ProfilAnggota,
	SikapFraksi,
	WakilDenganUmpan,
	BelanjaAPBD,
	BelanjaTema,
	DaftarDataBelumDibuka,
	DasborLayanan,
	DasborTanggungJawab,
	DasborUang,
	DetailPaket,
	DetailPendidikan,
	DetailWilayah,
	HasilPaketTerkait,
	Dokumen,
	InfoSumberData,
	Jaringan,
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

const asalRelatif = 'http://relatif.invalid';

function alamatDasar(): string {
	return env.PUBLIC_API_URL || alamatBawaan;
}

function keTeks(url: URL): string {
	return url.origin === asalRelatif ? url.pathname + url.search : url.toString();
}

export function susunURL(path: string, kueri: Kueri): string {
	const url = new URL(alamatDasar() + path, asalRelatif);
	for (const [nama, nilai] of Object.entries(kueri)) {
		if (nilai !== undefined && nilai !== '') url.searchParams.set(nama, String(nilai));
	}
	return keTeks(url);
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
	metode: 'POST' | 'PATCH' | 'PUT' = 'POST'
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

export function alamatUnduhPaket(kodeWilayah: string, tahun: number): string {
	return susunURL('/unduh/paket.csv', { kode_wilayah: kodeWilayah, tahun });
}

export function alamatKartuKecamatan(kode: string, bahasa: string): string {
	return susunURL(`/kecamatan/${encodeURIComponent(kode)}/kartu.png`, { bahasa });
}

export function urlServer(path: string): string {
	return keTeks(new URL(path, new URL(alamatDasar(), asalRelatif)));
}

export async function ambilBerwenang<T>(
	fetchFn: Fetch,
	path: string,
	token: string,
	metode: 'GET' | 'DELETE' = 'GET'
): Promise<T> {
	const respons = await fetchFn(susunURL(path, {}), {
		method: metode,
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

function pathTemaKecamatan(tema: string, kode: string): string {
	return `/tema/${encodeURIComponent(tema)}/kecamatan/${encodeURIComponent(kode)}`;
}

export function ambilDetailPendidikan(
	fetchFn: Fetch,
	tema: string,
	kode: string
): Promise<DetailPendidikan> {
	return ambilData<DetailPendidikan>(fetchFn, pathTemaKecamatan(tema, kode));
}

export function ambilSekolahKecamatan(
	fetchFn: Fetch,
	tema: string,
	kode: string
): Promise<Sekolah[]> {
	return ambilData<Sekolah[]>(fetchFn, `${pathTemaKecamatan(tema, kode)}/sekolah`);
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
	radiusM: number,
	tema: string
): Promise<SekitarTitik> {
	return ambilData<SekitarTitik>(fetchFn, '/sekolah/terdekat', {
		lat: titik.lat,
		lng: titik.lng,
		radius_m: radiusM,
		tema
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

export interface FilterJaringan {
	wilayah?: string;
	penyedia?: number;
	tahun?: number;
}

export function ambilJaringan(fetchFn: Fetch, filter: FilterJaringan): Promise<Jaringan> {
	return ambilData<Jaringan>(fetchFn, '/jaringan', { ...filter });
}

export function ambilDetailWilayah(fetchFn: Fetch, kode: string): Promise<DetailWilayah> {
	return ambilData<DetailWilayah>(fetchFn, `/wilayah/${encodeURIComponent(kode)}`);
}

export function ambilBelanjaTema(fetchFn: Fetch, kodeTema: string): Promise<BelanjaTema> {
	return ambilData<BelanjaTema>(fetchFn, `/tema/${encodeURIComponent(kodeTema)}/belanja`);
}

export function ambilAnggota(fetchFn: Fetch, id: string): Promise<ProfilAnggota> {
	return ambilData<ProfilAnggota>(fetchFn, `/anggota/${encodeURIComponent(id)}`);
}

export function ambilUmpanAnggota(
	fetchFn: Fetch,
	id: string,
	filter: { jenis?: JenisSumberUmpan; tag?: string; halaman?: number }
): Promise<Amplop<ItemUmpan[]>> {
	return ambilAmplop<ItemUmpan[]>(fetchFn, `/anggota/${encodeURIComponent(id)}/umpan`, filter);
}

export function ambilIsuAnggota(fetchFn: Fetch, id: string): Promise<IsuDisuarakan[]> {
	return ambilData<IsuDisuarakan[]>(fetchFn, `/anggota/${encodeURIComponent(id)}/isu`);
}

export function ambilWakilDenganUmpan(
	fetchFn: Fetch,
	kodeKecamatan: string,
	kodeTema: string
): Promise<WakilDenganUmpan> {
	return ambilData<WakilDenganUmpan>(
		fetchFn,
		`/kecamatan/${encodeURIComponent(kodeKecamatan)}/wakil`,
		{
			tema: kodeTema,
			dengan_umpan: true
		}
	);
}

export function ambilSikapFraksi(fetchFn: Fetch, fraksi: string): Promise<SikapFraksi[]> {
	return ambilData<SikapFraksi[]>(fetchFn, `/fraksi/${encodeURIComponent(fraksi)}/sikap`);
}

export function ambilAspirasiReses(fetchFn: Fetch, dapil: number): Promise<AspirasiReses[]> {
	return ambilData<AspirasiReses[]>(fetchFn, '/aspirasi-reses', { dapil });
}

export function ambilPaketTerkait(fetchFn: Fetch, idAspirasi: number): Promise<HasilPaketTerkait> {
	return ambilData<HasilPaketTerkait>(fetchFn, `/aspirasi-reses/${idAspirasi}/paket-terkait`);
}

export function ambilDasborLayanan(fetchFn: Fetch, kodeTema: string): Promise<DasborLayanan> {
	return ambilData<DasborLayanan>(fetchFn, `/tema/${encodeURIComponent(kodeTema)}/dasbor`);
}

export function ambilDasborUang(fetchFn: Fetch, tahun?: number): Promise<DasborUang> {
	return ambilData<DasborUang>(fetchFn, '/dasbor/uang', { tahun });
}

export function ambilDasborTanggungJawab(
	fetchFn: Fetch,
	tahun?: number
): Promise<DasborTanggungJawab> {
	return ambilData<DasborTanggungJawab>(fetchFn, '/dasbor/tanggung-jawab', { tahun });
}

export function ambilPendidikanKota(
	fetchFn: Fetch,
	tema: string,
	kode: string
): Promise<import('./tipe').PendidikanKota> {
	return ambilData(fetchFn, `/tema/${tema}/kota/${kode}`);
}

export function ambilBelanjaAPBD(fetchFn: Fetch, tahun: number): Promise<BelanjaAPBD> {
	return ambilData<BelanjaAPBD>(fetchFn, '/apbd', { tahun });
}
