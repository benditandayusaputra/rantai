import { ambilData, kirimBerwenang, type Fetch } from './klien';
import type {
	DetailLaporan,
	HasilKonfirmasi,
	KondisiLaporan,
	LaporanDiterima,
	PaketTerdekat,
	PetaLaporan,
	Titik
} from './tipe';

export function ambilPaketTerdekat(
	fetchFn: Fetch,
	titik: Titik,
	radiusM: number
): Promise<PaketTerdekat[]> {
	return ambilData<PaketTerdekat[]>(fetchFn, '/paket/terdekat', {
		lat: titik.lat,
		lng: titik.lng,
		radius_m: radiusM
	});
}

export function kirimLaporan(
	fetchFn: Fetch,
	token: string,
	formulir: FormData
): Promise<LaporanDiterima> {
	return kirimBerwenang<LaporanDiterima>(fetchFn, '/laporan', token, formulir);
}

export function ambilLaporan(fetchFn: Fetch, id: number): Promise<DetailLaporan> {
	return ambilData<DetailLaporan>(fetchFn, `/laporan/${id}`);
}

export interface KonfirmasiSaksi {
	titik: Titik;
	akurasiM: number;
	kondisi: KondisiLaporan;
}

export function konfirmasiLaporan(
	fetchFn: Fetch,
	token: string,
	id: number,
	k: KonfirmasiSaksi
): Promise<HasilKonfirmasi> {
	return kirimBerwenang<HasilKonfirmasi>(fetchFn, `/laporan/${id}/konfirmasi`, token, {
		lat: k.titik.lat,
		lng: k.titik.lng,
		akurasi: k.akurasiM,
		kondisi: k.kondisi
	});
}

export function ambilPetaLaporan(fetchFn: Fetch, bbox: string): Promise<PetaLaporan> {
	return ambilData<PetaLaporan>(fetchFn, '/peta/laporan', { bbox });
}
