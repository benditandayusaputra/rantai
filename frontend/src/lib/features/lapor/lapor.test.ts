import { describe, expect, it, vi } from 'vitest';
import {
	bacaLangkah,
	drafKosong,
	isDrafTerisi,
	kenaliTempat,
	langkahTerjauh,
	nilaiPosisi,
	susunLinimasa,
	susunFormulir,
	ukuranSkala,
	type DrafLaporan
} from './lapor';

vi.mock('$env/dynamic/public', () => ({ env: { PUBLIC_API_URL: 'http://api.uji/api/v1' } }));

function jawab(status: number, data: unknown) {
	return vi.fn(
		async () => new Response(JSON.stringify({ code: status, message: 'x', data }), { status })
	);
}

function drafLengkap(): DrafLaporan {
	return {
		titik: { lat: -6.12, lng: 106.9 },
		akurasiM: 14.6,
		paketKode: 'nontender-1',
		paketTerdekat: [],
		foto: [{ id: 'a', berkas: new Blob(['x']), pratinjau: '', diambilPada: new Date() }],
		kondisi: 'terhenti',
		catatan: '  Pekerjaan berhenti.  ',
		setuju: true
	};
}

describe('lapor', () => {
	it('membaca langkah dari URL dan jatuh ke langkah 1 bila tidak sah', () => {
		expect(bacaLangkah('3')).toBe(3);
		expect([bacaLangkah(null), bacaLangkah('0'), bacaLangkah('7'), bacaLangkah('2.5')]).toEqual([
			1, 1, 1, 1
		]);
	});

	it('menolak akurasi lebih dari 100 m dengan status kasar', () => {
		expect(nilaiPosisi({ lat: 0, lng: 0 }, 100).status).toBe('akurat');
		expect(nilaiPosisi({ lat: 0, lng: 0 }, 180).status).toBe('kasar');
	});

	it('tidak membiarkan melompat melewati langkah yang belum lengkap', () => {
		const draf = drafKosong();
		expect(langkahTerjauh(draf)).toBe(1);
		draf.titik = { lat: -6.12, lng: 106.9 };
		draf.akurasiM = 180;
		expect(langkahTerjauh(draf)).toBe(1);
		draf.akurasiM = 20;
		draf.paketKode = '';
		expect(langkahTerjauh(draf)).toBe(3);
		expect(langkahTerjauh(drafLengkap())).toBe(6);
	});

	it('menyusun formulir multipart tanpa paket bila proyek tidak ada di daftar', () => {
		const formulir = susunFormulir({ ...drafLengkap(), paketKode: '' });
		expect(formulir.get('paket_kode')).toBeNull();
		expect(formulir.get('akurasi')).toBe('15');
		expect(formulir.get('catatan')).toBe('Pekerjaan berhenti.');
		expect(formulir.getAll('foto')).toHaveLength(1);
	});

	it('memperkecil foto hanya bila melebihi sisi terpanjang', () => {
		expect(ukuranSkala(4032, 3024, 2048)).toEqual({ lebar: 2048, tinggi: 1536 });
		expect(ukuranSkala(800, 600, 2048)).toEqual({ lebar: 800, tinggi: 600 });
	});
});

describe('susunLinimasa', () => {
	it('menambah tahap mendatang setelah status terakhir', () => {
		const linimasa = susunLinimasa([
			{
				status_lama: null,
				status_baru: 'DIKIRIM',
				alasan: '',
				dibuat_pada: '2026-10-01T02:00:00Z'
			},
			{
				status_lama: 'DIKIRIM',
				status_baru: 'MENUNGGU_SAKSI',
				alasan: '',
				dibuat_pada: '2026-10-01T02:00:01Z'
			}
		]);
		expect(linimasa.map((p) => [p.status, p.isTerjadi])).toEqual([
			['DIKIRIM', true],
			['MENUNGGU_SAKSI', true],
			['TERVERIFIKASI', false],
			['DITERUSKAN', false],
			['DITANGGAPI', false],
			['SELESAI', false]
		]);
	});

	it('berhenti di status ditolak tanpa tahap mendatang', () => {
		const linimasa = susunLinimasa([
			{
				status_lama: 'MENUNGGU_SAKSI',
				status_baru: 'DITOLAK',
				alasan: 'duplikat',
				dibuat_pada: '2026-10-02T00:00:00Z'
			}
		]);
		expect(linimasa).toHaveLength(1);
		expect(linimasa[0]?.catatan).toBe('duplikat');
	});
});

describe('konfirmasi lokasi laporan', () => {
	const titik = { lat: -6.158, lng: 106.905 };

	it.each([
		[200, { kode: '31.72.06.1001', nama: 'Kelapa Gading Timur' }, 'dikenali'],
		[404, null, 'luar'],
		[500, null, 'gagal']
	] as const)('jawaban %i dari pencarian titik menjadi %s', async (status, data, harapan) => {
		const hasil = await kenaliTempat(jawab(status, data), titik);

		expect(hasil.status).toBe(harapan);
	});

	it('menganggap draf terisi setelah ada foto, kondisi, atau catatan', () => {
		expect(isDrafTerisi(drafKosong())).toBe(false);
		expect(isDrafTerisi({ ...drafKosong(), catatan: 'retak' })).toBe(true);
		expect(isDrafTerisi(drafLengkap())).toBe(true);
	});
});
