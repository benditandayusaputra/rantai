import { ambilBerwenang, kirimBerwenang, type Fetch } from './klien';
import type { HalamanNotifikasi, Langganan } from './tipe';

export function daftarLangganan(fetchFn: Fetch, token: string): Promise<Langganan[]> {
	return ambilBerwenang<Langganan[]>(fetchFn, '/langganan', token);
}

export function langgani(fetchFn: Fetch, token: string, kodeWilayah: string): Promise<Langganan> {
	return kirimBerwenang<Langganan>(fetchFn, '/langganan', token, { kode_wilayah: kodeWilayah });
}

export function berhentiLangganan(fetchFn: Fetch, token: string, id: number): Promise<null> {
	return ambilBerwenang<null>(fetchFn, `/langganan/${id}`, token, 'DELETE');
}

export function ambilNotifikasi(fetchFn: Fetch, token: string): Promise<HalamanNotifikasi> {
	return ambilBerwenang<HalamanNotifikasi>(fetchFn, '/notifikasi', token);
}

export function tandaiDibaca(fetchFn: Fetch, token: string, id: number): Promise<null> {
	return kirimBerwenang<null>(fetchFn, `/notifikasi/${id}/dibaca`, token, {}, 'PATCH');
}
