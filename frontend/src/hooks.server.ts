import type { Handle } from '@sveltejs/kit';
import { bahasaDari } from '$lib/i18n';

const penandaBahasa = '%rantai.lang%';

export const handle: Handle = ({ event, resolve }) => {
	const bahasa = bahasaDari(event.params.lang);
	return resolve(event, { transformPageChunk: ({ html }) => html.replace(penandaBahasa, bahasa) });
};
