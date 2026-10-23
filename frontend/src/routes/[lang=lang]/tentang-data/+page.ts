import { ambilSumberData, tangkap } from '$lib/api/klien';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ fetch }) => ({ sumber: tangkap(ambilSumberData(fetch)) });
