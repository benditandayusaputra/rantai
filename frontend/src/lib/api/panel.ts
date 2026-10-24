import { ambilBerwenang, kirimBerwenang, type Fetch } from './klien';
import type { AntreanPanel, HakJawab, HakJawabMasuk } from './tipe';

export function ambilAntreanPanel(fetchFn: Fetch, token: string): Promise<AntreanPanel> {
	return ambilBerwenang<AntreanPanel>(fetchFn, '/panel/antrean', token);
}

export function terbitkanHakJawab(
	fetchFn: Fetch,
	token: string,
	masukan: HakJawabMasuk
): Promise<HakJawab> {
	return kirimBerwenang<HakJawab>(fetchFn, '/hak-jawab', token, masukan);
}

export function isiSerahTerima(
	fetchFn: Fetch,
	token: string,
	kodePaket: string,
	isian: { tanggal: string; nomor_berita_acara: string }
): Promise<null> {
	return kirimBerwenang<null>(
		fetchFn,
		`/panel/paket/${encodeURIComponent(kodePaket)}/serah-terima`,
		token,
		isian,
		'PATCH'
	);
}

export function tandaiLaporanSelesai(fetchFn: Fetch, token: string, id: number): Promise<null> {
	return kirimBerwenang<null>(fetchFn, `/panel/laporan/${id}/selesai`, token, {}, 'PATCH');
}
