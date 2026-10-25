import { ambilPetaTema, tangkap } from '$lib/api/klien';
import { bacaTemaPendidikan } from '$lib/features/pendidikan/jenjang';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ fetch, url }) => {
	const jenjang = bacaTemaPendidikan(url.searchParams.get('tema'));
	return { jenjang, tema: tangkap(ambilPetaTema(fetch, jenjang)) };
};
