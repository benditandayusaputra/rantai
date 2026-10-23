import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import type { Angka, Label } from '$lib/api/tipe';
import { t } from '$lib/i18n';
import CertaintyBadge from './CertaintyBadge.svelte';
import ProvenanceBadge from './ProvenanceBadge.svelte';

vi.mock('$lib/features/provenance/provenance', async (asli) => ({
	...(await asli<typeof import('$lib/features/provenance/provenance')>()),
	muatProvenance: vi.fn(async (id: string) => ({
		id,
		sumber: 'dapodik',
		url: 'https://dapo.kemendikdasmen.go.id',
		diambil_pada: '2026-10-14T01:14:00Z',
		label: 'resmi',
		hash_konten: '',
		catatan: 'jumlah rombel kelas 7 x 32'
	}))
}));

const semuaLabel: Label[] = ['resmi', 'proksi', 'estimasi', 'simulasi'];

describe('CertaintyBadge', () => {
	it.each(semuaLabel)('menampilkan teks label %s yang terlihat', (level) => {
		render(CertaintyBadge, { level });
		const lencana = screen.getByText(t(`cert.${level}`));
		expect(lencana).toBeVisible();
		expect(lencana).toHaveAttribute('data-label', level);
	});
});

describe('ProvenanceBadge', () => {
	const angka: Angka = { nilai: 1284, satuan: 'kursi', label: 'proksi', provenance_id: 'prv_1' };

	it('membuka panel sumber dengan papan ketik', async () => {
		const pengguna = userEvent.setup();
		render(ProvenanceBadge, { angka, nilai: '1.284' });
		const tombol = await screen.findByRole('button', { name: /Dapodik, Proksi/ });
		await pengguna.tab();
		expect(tombol).toHaveFocus();
		await pengguna.keyboard('{Enter}');
		expect(tombol).toHaveAttribute('aria-expanded', 'true');
		const panel = screen.getByRole('dialog', { hidden: true });
		expect(panel).toHaveAttribute('open');
		expect(panel).toHaveTextContent('jumlah rombel kelas 7 x 32');
		expect(panel).toHaveTextContent('14 Okt 2026, 08.14 WIB');
	});

	it('menutup panel lewat tombol tutup', async () => {
		const pengguna = userEvent.setup();
		render(ProvenanceBadge, { angka, nilai: '1.284' });
		await pengguna.click(await screen.findByRole('button', { name: /Dapodik/ }));
		await pengguna.click(screen.getByRole('button', { name: t('action.close'), hidden: true }));
		expect(screen.getByRole('dialog', { hidden: true })).not.toHaveAttribute('open');
	});
});
