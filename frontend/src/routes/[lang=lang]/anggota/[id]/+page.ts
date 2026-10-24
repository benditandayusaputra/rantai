import { ambilAnggota, tangkap } from '$lib/api/klien';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params, fetch }) => ({
	anggota: tangkap(ambilAnggota(fetch, params.id))
});
