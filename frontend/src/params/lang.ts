import { bahasaTersedia, type Bahasa } from '$lib/i18n';

export function match(param: string): param is Bahasa {
	return (bahasaTersedia as readonly string[]).includes(param);
}
