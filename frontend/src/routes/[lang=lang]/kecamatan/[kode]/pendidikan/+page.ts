import { ambilDetailPendidikan, ambilSekolahKecamatan, tangkap } from '$lib/api/klien';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params, fetch }) => ({
	kode: params.kode,
	detail: tangkap(ambilDetailPendidikan(fetch, params.kode)),
	sekolah: tangkap(ambilSekolahKecamatan(fetch, params.kode))
});
