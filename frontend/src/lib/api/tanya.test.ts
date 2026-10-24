import { describe, expect, it, vi } from 'vitest';
import { GalatAPI } from './klien';
import { GalatTanya, tanyaRantai, uraiPeristiwa } from './tanya';

function aliran(...potongan: string[]): Response {
	const badan = new ReadableStream<Uint8Array>({
		start(kontrol) {
			for (const p of potongan) kontrol.enqueue(new TextEncoder().encode(p));
			kontrol.close();
		}
	});
	return new Response(badan, { status: 200, headers: { 'Content-Type': 'text/event-stream' } });
}

const permintaan = { pertanyaan: 'Kursi cukup?', kode_kecamatan_konteks: '', bahasa: 'id' };

describe('uraiPeristiwa', () => {
	it('memisahkan peristiwa lengkap dan menyimpan potongan yang belum selesai', () => {
		const hasil = uraiPeristiwa('event: alat\ndata: {"alat":"cari_wilayah"}\n\nevent: jaw');
		expect(hasil.peristiwa).toEqual([{ nama: 'alat', data: { alat: 'cari_wilayah' } }]);
		expect(hasil.sisa).toBe('event: jaw');
	});
});

describe('tanyaRantai', () => {
	it('melaporkan setiap alat lalu mengembalikan jawaban meski peristiwa terpotong antarpotongan', async () => {
		const fetchFn = vi.fn(async () =>
			aliran(
				'event: alat\ndata: {"alat":"indikator_tema"}\n\nevent: jawa',
				'ban\ndata: {"teks":"Kursi 288 [1].","rujukan":[{"nomor":1,"provenance_id":"prv_1"}]}\n\n'
			)
		);
		const alat: string[] = [];
		const jawaban = await tanyaRantai(fetchFn, permintaan, (a) => alat.push(a));
		expect(alat).toEqual(['indikator_tema']);
		expect(jawaban.rujukan).toEqual([{ nomor: 1, provenance_id: 'prv_1' }]);
	});

	it('melempar GalatTanya saat server mengirim peristiwa galat', async () => {
		const fetchFn = vi.fn(async () => aliran('event: galat\ndata: {"pesan":"belum bisa"}\n\n'));
		await expect(tanyaRantai(fetchFn, permintaan, () => {})).rejects.toEqual(
			new GalatTanya('belum bisa')
		);
	});

	it('melempar GalatAPI dengan pesan amplop saat permintaan ditolak', async () => {
		const fetchFn = vi.fn(
			async () =>
				new Response(JSON.stringify({ code: 400, message: 'pertanyaan wajib diisi', data: null }), {
					status: 400
				})
		);
		await expect(tanyaRantai(fetchFn, permintaan, () => {})).rejects.toEqual(
			new GalatAPI(400, 'pertanyaan wajib diisi')
		);
	});
});
