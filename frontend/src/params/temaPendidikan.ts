import { isTemaPendidikan, type TemaPendidikan } from '$lib/features/pendidikan/jenjang';

export function match(param: string): param is TemaPendidikan {
	return isTemaPendidikan(param);
}
