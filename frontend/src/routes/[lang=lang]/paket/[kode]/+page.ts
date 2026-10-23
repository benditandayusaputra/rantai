import { ambilPaket, tangkap } from '$lib/api/klien';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params, fetch }) => ({
	paket: tangkap(ambilPaket(fetch, params.kode))
});
