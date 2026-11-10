import { GalatAPI, susunURL, type Fetch } from './klien';
import type { Amplop } from './tipe';

const pemisahPeristiwa = '\n\n';

export interface PermintaanTanya {
	pertanyaan: string;
	kode_kecamatan_konteks: string;
	bahasa: string;
}

export interface RujukanTanya {
	nomor: number;
	provenance_id: string;
}

export interface TautanTanya {
	teks: string;
	jenis: string;
	tujuan: string;
}

export interface JawabanTanya {
	teks: string;
	rujukan: RujukanTanya[];
	tautan?: TautanTanya[];
}

export interface Peristiwa {
	nama: string;
	data: unknown;
}

export class GalatTanya extends Error {
	constructor(message: string) {
		super(message);
		this.name = 'GalatTanya';
	}
}

function uraiSatuPeristiwa(blok: string): Peristiwa | null {
	let nama = 'message';
	const data: string[] = [];
	for (const baris of blok.split('\n')) {
		if (baris.startsWith('event:')) nama = baris.slice(6).trim();
		else if (baris.startsWith('data:')) data.push(baris.slice(5).trimStart());
	}
	return data.length ? { nama, data: JSON.parse(data.join('\n')) } : null;
}

export function uraiPeristiwa(penyangga: string): { peristiwa: Peristiwa[]; sisa: string } {
	const blok = penyangga.split(pemisahPeristiwa);
	const sisa = blok.pop() ?? '';
	const peristiwa = blok.map(uraiSatuPeristiwa).filter((p): p is Peristiwa => p !== null);
	return { peristiwa, sisa };
}

async function tolakBukanAliran(respons: Response): Promise<never> {
	const amplop = (await respons.json().catch(() => null)) as Amplop<null> | null;
	throw new GalatAPI(respons.status, amplop?.message ?? respons.statusText);
}

function terapkan(p: Peristiwa, onAlat: (alat: string) => void): JawabanTanya | null {
	if (p.nama === 'alat') onAlat((p.data as { alat: string }).alat);
	if (p.nama === 'galat') throw new GalatTanya((p.data as { pesan: string }).pesan);
	return p.nama === 'jawaban' ? (p.data as JawabanTanya) : null;
}

export async function tanyaRantai(
	fetchFn: Fetch,
	permintaan: PermintaanTanya,
	onAlat: (alat: string) => void
): Promise<JawabanTanya> {
	const respons = await fetchFn(susunURL('/tanya', {}), {
		method: 'POST',
		headers: { 'Content-Type': 'application/json', Accept: 'text/event-stream' },
		body: JSON.stringify(permintaan)
	});
	if (!respons.ok || !respons.body) return tolakBukanAliran(respons);
	const pembaca = respons.body.pipeThrough(new TextDecoderStream()).getReader();
	let penyangga = '';
	for (;;) {
		const { value, done } = await pembaca.read();
		if (done) throw new GalatTanya('');
		const hasil = uraiPeristiwa(penyangga + value);
		penyangga = hasil.sisa;
		for (const p of hasil.peristiwa) {
			const jawaban = terapkan(p, onAlat);
			if (jawaban) return jawaban;
		}
	}
}
