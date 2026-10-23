import { redirect } from '@sveltejs/kit';
import { bahasaBawaan } from '$lib/i18n';

export function load() {
	redirect(307, `/${bahasaBawaan}`);
}
