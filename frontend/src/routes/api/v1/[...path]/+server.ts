import { error } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import { saringRespons, susunPermintaanHulu } from '$lib/server/proxy';
import type { RequestHandler } from './$types';

export const fallback: RequestHandler = async ({ request, params, getClientAddress }) => {
	if (!env.API_HULU || !env.PROXY_RAHASIA) error(404);
	const hulu = susunPermintaanHulu(
		request,
		params.path,
		{ hulu: env.API_HULU, rahasia: env.PROXY_RAHASIA },
		getClientAddress()
	);
	return saringRespons(await fetch(hulu));
};
