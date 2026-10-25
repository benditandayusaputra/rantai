import { tangkap } from '$lib/api/klien';
import { bacaLensa, bacaTahun, muatDasbor } from '$lib/features/dasbor/dasbor';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url, fetch }) => {
	const lensa = bacaLensa(url.searchParams.get('lensa'));
	const tahun = bacaTahun(url.searchParams.get('tahun'));
	return { lensa, hasil: await tangkap(muatDasbor(fetch, lensa, tahun)) };
};
