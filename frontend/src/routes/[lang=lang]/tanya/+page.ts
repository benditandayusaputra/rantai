import { ambilWilayah, tangkap } from '$lib/api/klien';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ url, fetch }) => {
	const kode = url.searchParams.get('w');
	return { wilayah: kode ? tangkap(ambilWilayah(fetch, kode)) : null };
};
