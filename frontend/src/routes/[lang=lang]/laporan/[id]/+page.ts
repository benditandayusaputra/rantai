import { tangkap } from '$lib/api/klien';
import { ambilLaporan } from '$lib/api/laporan';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params, fetch }) => {
	const id = Number(params.id);
	return { id, laporan: tangkap(ambilLaporan(fetch, id)) };
};
