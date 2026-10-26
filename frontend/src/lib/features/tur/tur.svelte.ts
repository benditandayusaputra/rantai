import { tandaiTurSelesai } from './tur';

export const tur = $state({ isTerbuka: false });

export function mulaiTur() {
	tur.isTerbuka = true;
}

export function tutupTur() {
	tur.isTerbuka = false;
	tandaiTurSelesai();
}
