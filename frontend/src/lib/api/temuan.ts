import { ambilData, kirimBerwenang, type Fetch } from './klien';
import type { SikapTemuan, Temuan } from './tipe';

export function ambilTemuan(fetchFn: Fetch, kodePaket: string): Promise<Temuan[]> {
	return ambilData<Temuan[]>(fetchFn, `/paket/${encodeURIComponent(kodePaket)}/temuan`);
}

export function kirimTemuan(
	fetchFn: Fetch,
	token: string,
	kodePaket: string,
	temuan: { isi: string; tautan_sumber: string }
): Promise<{ id: number }> {
	return kirimBerwenang(fetchFn, `/paket/${encodeURIComponent(kodePaket)}/temuan`, token, temuan);
}

export function tanggapiTemuan(
	fetchFn: Fetch,
	token: string,
	id: number,
	tanggapan: { sikap: SikapTemuan; tautan_sumber: string }
): Promise<null> {
	return kirimBerwenang<null>(fetchFn, `/temuan/${id}/tanggapan`, token, tanggapan);
}

export function sembunyikanTemuan(
	fetchFn: Fetch,
	token: string,
	id: number,
	alasan: string
): Promise<null> {
	return kirimBerwenang<null>(fetchFn, `/moderasi/temuan/${id}`, token, { alasan }, 'PATCH');
}
