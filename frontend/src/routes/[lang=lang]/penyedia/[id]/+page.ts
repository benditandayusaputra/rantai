import { ambilPenyedia, tangkap } from '$lib/api/klien';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params, fetch }) => ({
	penyedia: tangkap(ambilPenyedia(fetch, Number(params.id)))
});
