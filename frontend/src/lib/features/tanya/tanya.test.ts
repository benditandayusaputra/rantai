import { describe, expect, it } from 'vitest';
import { pecahRujukan } from './tanya';

describe('pecahRujukan', () => {
	it('memecah teks menjadi potongan teks dan nomor rujukan yang dikenal', () => {
		expect(pecahRujukan('Kursi 288 [1], lulusan 400 [2].', [1, 2])).toEqual([
			{ jenis: 'teks', isi: 'Kursi 288 ' },
			{ jenis: 'rujukan', nomor: 1 },
			{ jenis: 'teks', isi: ', lulusan 400 ' },
			{ jenis: 'rujukan', nomor: 2 },
			{ jenis: 'teks', isi: '.' }
		]);
	});

	it('membiarkan kurung bernomor tanpa rujukan sebagai teks biasa', () => {
		expect(pecahRujukan('Lihat [3] di atas.', [1])).toEqual([
			{ jenis: 'teks', isi: 'Lihat [3] di atas.' }
		]);
	});
});
