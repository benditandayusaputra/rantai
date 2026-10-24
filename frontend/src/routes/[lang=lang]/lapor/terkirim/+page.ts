import { tangkap } from '$lib/api/klien';
import { ambilLaporan } from '$lib/api/laporan';
import type { PageLoad } from './$types';

export const ssr = false;

export const load: PageLoad = ({ url, fetch }) => {
	const id = Number(url.searchParams.get('id'));
	return { id, laporan: tangkap(ambilLaporan(fetch, id)) };
};
