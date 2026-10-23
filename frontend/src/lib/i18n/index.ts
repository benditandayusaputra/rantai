import id from './id.json';

export const bahasaTersedia = ['id'] as const;
export type Bahasa = (typeof bahasaTersedia)[number];
export const bahasaBawaan: Bahasa = 'id';

export type KunciPesan = keyof typeof id;
export type ParameterPesan = Record<string, string | number>;

const katalog: Record<Bahasa, Record<KunciPesan, string>> = { id };
const polaParameter = /\{(\w+)\}/g;

export function t(kunci: KunciPesan, parameter: ParameterPesan = {}): string {
	return katalog[bahasaBawaan][kunci].replace(polaParameter, (asli, nama: string) =>
		nama in parameter ? String(parameter[nama]) : asli
	);
}

export function isKunciPesan(kunci: string): kunci is KunciPesan {
	return kunci in katalog[bahasaBawaan];
}
