import {
	ambilPaketKecamatan,
	ambilPetaTema,
	ambilRingkasanKecamatan,
	tangkap
} from '$lib/api/klien';
import { bacaSorotan } from '$lib/features/beranda/beranda';
import { bacaTemaPendidikan } from '$lib/features/pendidikan/jenjang';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ url, fetch }) => {
	const sorotan = bacaSorotan(url.searchParams.get('wilayah'));
	const jenjang = bacaTemaPendidikan(url.searchParams.get('tema'));
	const [ringkasan, paket, peta] = await Promise.all([
		tangkap(ambilRingkasanKecamatan(fetch, sorotan.kode)),
		tangkap(ambilPaketKecamatan(fetch, sorotan.kode, {})),
		tangkap(ambilPetaTema(fetch, jenjang))
	]);
	return { sorotan, jenjang, ringkasan, paket, peta };
};
