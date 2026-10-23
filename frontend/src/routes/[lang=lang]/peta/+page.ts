import { ambilPetaTema, tangkap } from '$lib/api/klien';
import type { PageLoad } from './$types';

const temaPeta = 'pendidikan';

export const load: PageLoad = ({ fetch }) => ({ tema: tangkap(ambilPetaTema(fetch, temaPeta)) });
