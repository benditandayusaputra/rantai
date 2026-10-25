import { ambilDetailPendidikan, ambilSekolahKecamatan, tangkap } from '$lib/api/klien';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params, fetch }) => ({
	kode: params.kode,
	tema: params.tema,
	detail: tangkap(ambilDetailPendidikan(fetch, params.tema, params.kode)),
	sekolah: tangkap(ambilSekolahKecamatan(fetch, params.tema, params.kode))
});
