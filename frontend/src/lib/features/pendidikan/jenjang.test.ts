import { describe, expect, it } from 'vitest';
import { bacaTemaPendidikan, kunciMetode, pilihanJenjang, tautanTemaPendidikan } from './jenjang';

describe('bacaTemaPendidikan', () => {
	it('menerima tema pendidikan yang dikenal dan jatuh ke SMP untuk selainnya', () => {
		expect(bacaTemaPendidikan('pendidikan_sma_smk')).toBe('pendidikan_sma_smk');
		expect(bacaTemaPendidikan('banjir')).toBe('pendidikan');
		expect(bacaTemaPendidikan(null)).toBe('pendidikan');
	});
});

describe('tautanTemaPendidikan', () => {
	const url = new URL('https://rantai.test/id/peta?lapisan=tema,sekolah&tema=pendidikan_sma_smk');

	it('mempertahankan parameter lain dan menghapus tema bawaan dari URL', () => {
		expect(tautanTemaPendidikan(url, 'pendidikan')).toBe('/id/peta?lapisan=tema%2Csekolah');
		expect(tautanTemaPendidikan(url, 'pendidikan_sma_smk')).toBe(
			'/id/peta?lapisan=tema%2Csekolah&tema=pendidikan_sma_smk'
		);
	});

	it('menandai jenjang aktif di pilihan', () => {
		const pilihan = pilihanJenjang(url, 'pendidikan_sma_smk');
		expect(pilihan.map((p) => [p.label, p.current])).toEqual([
			['SMP', false],
			['SMA/SMK', true]
		]);
	});
});

describe('kunciMetode', () => {
	it('memakai kalimat per sekolah sesuai jenjang asal', () => {
		expect(kunciMetode('pendidikan', 'per_sekolah')).toBe('pen.method.per_sekolah');
		expect(kunciMetode('pendidikan_sma_smk', 'per_sekolah')).toBe('pen.sma.method.per_sekolah');
		expect(kunciMetode('pendidikan_sma_smk', 'estimasi_total_smp')).toBe(
			'pen.method.estimasi_total_smp'
		);
		expect(kunciMetode('pendidikan', '')).toBeNull();
	});
});
