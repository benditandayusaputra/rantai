import en from './en.json';
import id from './id.json';

export const bahasaTersedia = ['id', 'en'] as const;
export type Bahasa = (typeof bahasaTersedia)[number];
export const bahasaBawaan: Bahasa = 'id';

export type KunciPesan = keyof typeof id;
export type ParameterPesan = Record<string, string | number>;

const katalog: Record<Bahasa, Record<KunciPesan, string>> = { id, en };
const polaParameter = /\{(\w+)\}/g;

export let bahasaAktif: Bahasa = bahasaBawaan;

export function isBahasa(nilai: string | undefined): nilai is Bahasa {
	return (bahasaTersedia as readonly (string | undefined)[]).includes(nilai);
}

export function bahasaDari(nilai: string | undefined): Bahasa {
	return isBahasa(nilai) ? nilai : bahasaBawaan;
}

export function aturBahasa(bahasa: Bahasa) {
	bahasaAktif = bahasa;
}

export function gantiBahasaPath(pathname: string, bahasa: Bahasa): string {
	const [, awal = '', ...sisa] = pathname.split('/');
	const bagian = isBahasa(awal) ? sisa : [awal, ...sisa].filter(Boolean);
	return ['', bahasa, ...bagian].join('/');
}

export function t(kunci: KunciPesan, parameter: ParameterPesan = {}): string {
	return katalog[bahasaAktif][kunci].replace(polaParameter, (asli, nama: string) =>
		nama in parameter ? String(parameter[nama]) : asli
	);
}

export function isKunciPesan(kunci: string): kunci is KunciPesan {
	return kunci in katalog[bahasaBawaan];
}
