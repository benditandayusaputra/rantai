import { tangkap } from '$lib/api/klien';
import { bacaPilihanPermohonan, muatKeterbukaan } from '$lib/features/keterbukaan/keterbukaan';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ url, fetch }) => {
	const pilihan = bacaPilihanPermohonan(url.searchParams);
	return { pilihan, halaman: tangkap(muatKeterbukaan(fetch, pilihan)) };
};
