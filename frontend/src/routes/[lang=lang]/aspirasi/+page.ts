import { tangkap } from '$lib/api/klien';
import { bacaPilihan, muatAspirasi } from '$lib/features/aspirasi/aspirasi';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ url, fetch }) => {
	const pilihan = bacaPilihan(url.searchParams);
	return { pilihan, halaman: tangkap(muatAspirasi(fetch, pilihan)) };
};
