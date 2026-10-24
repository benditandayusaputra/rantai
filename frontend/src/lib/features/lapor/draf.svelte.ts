import { drafKosong, type DrafLaporan } from './lapor';

export const draf = $state<DrafLaporan>(drafKosong());

export function aturUlangDraf(): void {
	for (const foto of draf.foto) URL.revokeObjectURL(foto.pratinjau);
	Object.assign(draf, drafKosong());
}
