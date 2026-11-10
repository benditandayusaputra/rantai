import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { t } from '$lib/i18n';
import PercakapanTanya from './PercakapanTanya.svelte';

const tanyaRantai = vi.hoisted(() => vi.fn());

vi.mock('$lib/api/tanya', async (asli) => ({
	...(await asli<typeof import('$lib/api/tanya')>()),
	tanyaRantai
}));

vi.mock('$lib/features/provenance/provenance', async (asli) => ({
	...(await asli<typeof import('$lib/features/provenance/provenance')>()),
	muatProvenance: vi.fn(async (id: string) => ({
		id,
		sumber: 'dapodik',
		url: '',
		diambil_pada: '2026-10-14T01:14:00Z',
		label: 'resmi',
		hash_konten: '',
		catatan: 'rombel kelas 7 x 32'
	}))
}));

describe('PercakapanTanya', () => {
	it('mengirim konteks kecamatan dan membuka panel sumber dari rujukan jawaban', async () => {
		tanyaRantai.mockImplementation(async (_f, _p, onAlat: (a: string) => void) => {
			onAlat('indikator_tema');
			return { teks: 'Ada 288 kursi [1].', rujukan: [{ nomor: 1, provenance_id: 'prv_1' }] };
		});
		const pengguna = userEvent.setup();
		render(PercakapanTanya, { konteks: { kode: '31.72.01', nama: 'Cilincing' } });

		await pengguna.type(
			screen.getByRole('textbox', { name: t('ask.placeholder') }),
			'Kursi cukup?'
		);
		await pengguna.click(screen.getByRole('button', { name: t('ask.send') }));
		const rujukan = await screen.findByRole('button', { name: 'Rujukan 1: Dapodik' });
		await pengguna.click(rujukan);

		expect(tanyaRantai.mock.calls[0]?.[1]).toMatchObject({ kode_kecamatan_konteks: '31.72.01' });
		const panel = screen.getByRole('dialog', { hidden: true });
		expect(panel).toHaveAttribute('open');
		expect(panel).toHaveTextContent('rombel kelas 7 x 32');
	});

	it('menampilkan pesan error yang bisa dipahami saat jawaban gagal', async () => {
		tanyaRantai.mockRejectedValue(new Error('jaringan putus'));
		const pengguna = userEvent.setup();
		render(PercakapanTanya, { konteks: { kode: '31.72.01', nama: 'Cilincing' } });

		await pengguna.click(screen.getByRole('button', { name: t('ask.suggest.1') }));

		expect(await screen.findByText(t('ask.error'))).toBeVisible();
	});

	it('menautkan nama dari jawaban ke halaman detailnya', async () => {
		tanyaRantai.mockResolvedValue({
			teks: 'Wakilmu Justin Adrian, S.H., M.H.',
			rujukan: [],
			tautan: [{ teks: 'Justin Adrian, S.H., M.H.', jenis: 'anggota', tujuan: 'justin-adrian' }]
		});
		const pengguna = userEvent.setup();
		render(PercakapanTanya, { konteks: { kode: '31.75.07', nama: 'Duren Sawit' } });

		await pengguna.click(screen.getByRole('button', { name: t('ask.suggest.3') }));

		expect(await screen.findByRole('link', { name: 'Justin Adrian, S.H., M.H.' })).toHaveAttribute(
			'href',
			'/id/anggota/justin-adrian'
		);
	});

	it('menahan pertanyaan sampai wilayah dipilih', () => {
		render(PercakapanTanya, {});

		expect(screen.getByRole('textbox', { name: t('ask.needArea') })).toBeDisabled();
		expect(screen.getByRole('button', { name: t('ask.suggest.1') })).toBeDisabled();
	});
});
