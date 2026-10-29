import { error } from '@sveltejs/kit';
import { ambilPendidikanKota, tangkap } from '$lib/api/klien';
import { bacaTemaPendidikan } from '$lib/features/pendidikan/jenjang';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, url, fetch }) => {
	if (!/^31\.\d{2}$/.test(params.kode)) error(404);
	const tema = bacaTemaPendidikan(url.searchParams.get('tema'));
	const hasil = await tangkap(ambilPendidikanKota(fetch, tema, params.kode));
	if (!hasil.ok && hasil.galat.kode === 404) error(404);
	return { hasil, tema, kode: params.kode };
};
