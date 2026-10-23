import { ambilRingkasanKecamatan } from '$lib/api/klien';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params, url, fetch }) => ({
	kode: params.kode,
	ringkasan: ambilRingkasanKecamatan(fetch, url.searchParams.get('kelurahan') ?? params.kode)
});
