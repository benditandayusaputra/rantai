import { ambilRingkasanKecamatan, tangkap } from '$lib/api/klien';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, url, fetch }) => ({
	kode: params.kode,
	ringkasan: await tangkap(
		ambilRingkasanKecamatan(fetch, url.searchParams.get('kelurahan') ?? params.kode)
	)
});
