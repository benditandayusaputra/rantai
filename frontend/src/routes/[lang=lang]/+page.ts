import {
	ambilPaketKecamatan,
	ambilPetaTema,
	ambilRingkasanKecamatan,
	tangkap
} from '$lib/api/klien';
import { bacaSorotan, temaBeranda } from '$lib/features/beranda/beranda';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url, fetch }) => {
	const sorotan = bacaSorotan(url.searchParams.get('wilayah'));
	const [ringkasan, paket, peta] = await Promise.all([
		tangkap(ambilRingkasanKecamatan(fetch, sorotan.kode)),
		tangkap(ambilPaketKecamatan(fetch, sorotan.kode, {})),
		tangkap(ambilPetaTema(fetch, temaBeranda))
	]);
	return { sorotan, ringkasan, paket, peta };
};
