import type { TautanTanya } from '$lib/api/tanya';
import { t } from '$lib/i18n';

export interface TautanSiap {
	teks: string;
	href: string;
	isLuar: boolean;
}

export type Segmen =
	| { jenis: 'teks'; isi: string }
	| { jenis: 'rujukan'; nomor: number }
	| { jenis: 'tautan'; isi: string; href: string; isLuar: boolean };

const polaIdAnggota = /^[a-z0-9-]+$/;
const polaKodeKecamatan = /^31\.\d{2}\.\d{2}$/;
const polaURLAman = /^https:\/\//;
const polaKarakterRegex = /[.*+?^${}()|[\]\\]/g;

function alamatTautan(tautan: TautanTanya, bahasa: string): string | null {
	if (tautan.jenis === 'anggota' && polaIdAnggota.test(tautan.tujuan))
		return `/${bahasa}/anggota/${tautan.tujuan}`;
	if (tautan.jenis === 'dapil' && polaKodeKecamatan.test(tautan.tujuan))
		return `/${bahasa}/kecamatan/${tautan.tujuan}?bagian=wakilku`;
	return tautan.jenis === 'rilis' && polaURLAman.test(tautan.tujuan) ? tautan.tujuan : null;
}

export function siapkanTautan(daftar: TautanTanya[], bahasa: string): TautanSiap[] {
	return daftar
		.map((tautan) => ({ tautan, href: alamatTautan(tautan, bahasa) }))
		.filter((item): item is { tautan: TautanTanya; href: string } => item.href !== null)
		.map(({ tautan, href }) => ({ teks: tautan.teks, href, isLuar: tautan.jenis === 'rilis' }))
		.sort((a, b) => b.teks.length - a.teks.length);
}

function pecahTautan(teks: string, tautan: TautanSiap[]): Segmen[] {
	if (tautan.length === 0) return [{ jenis: 'teks', isi: teks }];
	const pola = new RegExp(
		tautan.map((item) => item.teks.replace(polaKarakterRegex, '\\$&')).join('|'),
		'g'
	);
	const segmen: Segmen[] = [];
	let posisi = 0;
	for (const cocok of teks.matchAll(pola)) {
		const siap = tautan.find((item) => item.teks === cocok[0]);
		if (!siap) continue;
		if (cocok.index > posisi) segmen.push({ jenis: 'teks', isi: teks.slice(posisi, cocok.index) });
		segmen.push({ jenis: 'tautan', isi: siap.teks, href: siap.href, isLuar: siap.isLuar });
		posisi = cocok.index + cocok[0].length;
	}
	if (posisi < teks.length) segmen.push({ jenis: 'teks', isi: teks.slice(posisi) });
	return segmen;
}

export function pecahJawaban(teks: string, nomorSah: number[], tautan: TautanSiap[]): Segmen[] {
	return pecahRujukan(teks, nomorSah).flatMap((segmen) =>
		segmen.jenis === 'teks' ? pecahTautan(segmen.isi, tautan) : [segmen]
	);
}

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
