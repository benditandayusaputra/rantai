import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import type { TitikTren } from '$lib/api/tipe';
import GrafikTren from './GrafikTren.svelte';

vi.mock('$lib/features/provenance/provenance', async (asli) => ({
	...(await asli<typeof import('$lib/features/provenance/provenance')>()),
	muatProvenance: vi.fn(() => new Promise(() => {}))
}));

const titik = (periode: string, nilai: number): TitikTren => ({
	periode,
	rasio: { nilai, satuan: 'rasio', label: 'resmi', provenance_id: 'p' }
});

describe('GrafikTren', () => {
	it('menyediakan tabel alternatif setiap tahun ajaran dengan label kepastian', async () => {
		render(GrafikTren, { tren: [titik('2025/2026', 1.5), titik('2026/2027', 1.8)] });
		await userEvent.click(screen.getByRole('button', { name: 'Lihat sebagai tabel' }));

		const baris = screen.getAllByRole('row');
		expect(baris).toHaveLength(3);
		expect(screen.getByRole('rowheader', { name: '2025/2026' })).toBeVisible();
		expect(screen.getAllByRole('button', { name: /Resmi/ })).toHaveLength(2);
	});
});
