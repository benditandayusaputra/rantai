import { tangkap } from '$lib/api/klien';
import { bacaLensa, bacaTahun, muatDasbor } from '$lib/features/dasbor/dasbor';
import { bacaTemaPendidikan } from '$lib/features/pendidikan/jenjang';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url, fetch }) => {
	const lensa = bacaLensa(url.searchParams.get('lensa'));
	const tahun = bacaTahun(url.searchParams.get('tahun'));
	const jenjang = bacaTemaPendidikan(url.searchParams.get('tema'));
	return { lensa, jenjang, hasil: await tangkap(muatDasbor(fetch, lensa, tahun, jenjang)) };
};
