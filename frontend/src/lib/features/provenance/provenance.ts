import { ambilProvenance } from '$lib/api/klien';
import type { Provenance } from '$lib/api/tipe';
import { isKunciPesan, t } from '$lib/i18n';

const cache = new Map<string, Promise<Provenance>>();

export function muatProvenance(id: string): Promise<Provenance> {
	const tersimpan = cache.get(id);
	if (tersimpan) return tersimpan;
	const permintaan = ambilProvenance(fetch, id);
	cache.set(id, permintaan);
	permintaan.catch(() => cache.delete(id));
	return permintaan;
}

export function namaSumber(kode: string): string {
	const kunci = `sumber.${kode}`;
	return isKunciPesan(kunci) ? t(kunci) : kode;
}

export function pemilikSumber(kode: string): string {
	const kunci = `sumber.${kode}.pemilik`;
	return isKunciPesan(kunci) ? t(kunci) : '';
}
