import { ambilBerwenang, kirimBerwenang, type Fetch } from './klien';
import type { AksiModerasi, AntreanModerasi } from './tipe';

export function ambilAntreanModerasi(fetchFn: Fetch, token: string): Promise<AntreanModerasi> {
	return ambilBerwenang<AntreanModerasi>(fetchFn, '/moderasi/antrean', token);
}

export function moderasiLaporan(
	fetchFn: Fetch,
	token: string,
	id: number,
	aksi: AksiModerasi
): Promise<null> {
	return kirimBerwenang<null>(fetchFn, `/moderasi/laporan/${id}`, token, aksi, 'PATCH');
}

export function sembunyikanUmpan(
	fetchFn: Fetch,
	token: string,
	id: number,
	alasan: string
): Promise<null> {
	return kirimBerwenang<null>(fetchFn, `/moderasi/umpan/${id}`, token, { alasan }, 'PATCH');
}
