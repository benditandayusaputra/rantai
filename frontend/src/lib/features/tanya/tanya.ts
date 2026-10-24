import { t } from '$lib/i18n';

export type Segmen = { jenis: 'teks'; isi: string } | { jenis: 'rujukan'; nomor: number };

const polaNomorRujukan = /\[(\d+)\]/g;

export function pecahRujukan(teks: string, nomorSah: number[]): Segmen[] {
	const segmen: Segmen[] = [];
	let posisi = 0;
	for (const cocok of teks.matchAll(polaNomorRujukan)) {
		const nomor = Number(cocok[1]);
		if (!nomorSah.includes(nomor)) continue;
		if (cocok.index > posisi) segmen.push({ jenis: 'teks', isi: teks.slice(posisi, cocok.index) });
		segmen.push({ jenis: 'rujukan', nomor });
		posisi = cocok.index + cocok[0].length;
	}
	if (posisi < teks.length) segmen.push({ jenis: 'teks', isi: teks.slice(posisi) });
	return segmen;
}

export const panjangMaksPertanyaan = 500;

export function saranPertanyaan(): string[] {
	return [t('ask.suggest.1'), t('ask.suggest.2'), t('ask.suggest.3'), t('ask.suggest.4')];
}
