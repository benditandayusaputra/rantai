import { describe, expect, it } from 'vitest';
import { pecahJawaban, pecahRujukan, siapkanTautan } from './tanya';

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

describe('tautan di jawaban', () => {
	const tautan = siapkanTautan(
		[
			{ teks: 'Justin Adrian, S.H., M.H.', jenis: 'anggota', tujuan: 'justin-adrian-sh-mh' },
			{ teks: 'Dapil DKI Jakarta 5', jenis: 'dapil', tujuan: '31.75.07' },
			{ teks: 'Komisi E Pantau GOR', jenis: 'rilis', tujuan: 'https://dprd.contoh/komisi-e' },
			{ teks: 'Rilis Berbahaya', jenis: 'rilis', tujuan: 'javascript:alert(1)' },
			{ teks: 'Anggota Aneh', jenis: 'anggota', tujuan: '../../masuk' }
		],
		'id'
	);

	it('hanya menyiapkan tujuan yang sah dan menandai rilis sebagai tautan luar', () => {
		expect(tautan.map((item) => [item.href, item.isLuar])).toEqual([
			['/id/anggota/justin-adrian-sh-mh', false],
			['/id/kecamatan/31.75.07?bagian=wakilku', false],
			['https://dprd.contoh/komisi-e', true]
		]);
	});

	it('menautkan nama di antara teks dan rujukan tanpa mengubah teks lain', () => {
		const segmen = pecahJawaban(
			'Justin Adrian, S.H., M.H. [1] mewakili Dapil DKI Jakarta 5. Rilis Berbahaya.',
			[1],
			tautan
		);

		expect(segmen).toEqual([
			{
				jenis: 'tautan',
				isi: 'Justin Adrian, S.H., M.H.',
				href: '/id/anggota/justin-adrian-sh-mh',
				isLuar: false
			},
			{ jenis: 'teks', isi: ' ' },
			{ jenis: 'rujukan', nomor: 1 },
			{ jenis: 'teks', isi: ' mewakili ' },
			{
				jenis: 'tautan',
				isi: 'Dapil DKI Jakarta 5',
				href: '/id/kecamatan/31.75.07?bagian=wakilku',
				isLuar: false
			},
			{ jenis: 'teks', isi: '. Rilis Berbahaya.' }
		]);
	});
});
