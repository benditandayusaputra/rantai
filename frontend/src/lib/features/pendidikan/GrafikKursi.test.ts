import { render, screen, within } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import type { Angka, DetailPendidikan } from '$lib/api/tipe';
import GrafikKursi from './GrafikKursi.svelte';

vi.mock('$lib/features/provenance/provenance', async (asli) => ({
	...(await asli<typeof import('$lib/features/provenance/provenance')>()),
	muatProvenance: vi.fn(() => new Promise(() => {}))
}));

const angka = (nilai: number, label: Angka['label'], satuan = 'kursi'): Angka => ({
	nilai,
	satuan,
	label,
	provenance_id: 'p'
});

const detail = {
	calon_lulusan: angka(2701, 'resmi', 'siswa'),
	kursi_negeri: angka(1284, 'resmi'),
	kursi_spmb_bersama: angka(640, 'resmi'),
	kursi_madrasah: angka(301, 'proksi')
} as DetailPendidikan;

describe('GrafikKursi', () => {
	it('menampilkan setiap jenis kursi dengan label kepastiannya', () => {
		render(GrafikKursi, { detail });
		const baris = screen.getAllByRole('listitem');
		expect(baris).toHaveLength(4);
		expect(within(baris[3] as HTMLElement).getByText('Kursi MTs')).toBeVisible();
		expect(within(baris[3] as HTMLElement).getAllByText('Proksi')[0]).toBeVisible();
		expect(within(baris[1] as HTMLElement).getAllByText('Resmi')[0]).toBeVisible();
	});

	it('mengganti grafik dengan tabel setara', async () => {
		render(GrafikKursi, { detail });
		await userEvent.setup().click(screen.getByRole('button', { name: 'Lihat sebagai tabel' }));
		const tabel = screen.getByRole('table');
		expect(within(tabel).getAllByRole('row')).toHaveLength(5);
		expect(within(tabel).getByRole('rowheader', { name: 'Calon lulusan SD' })).toBeVisible();
		expect(screen.getByRole('button', { name: 'Lihat sebagai grafik' })).toHaveAttribute(
			'aria-pressed',
			'true'
		);
	});
});
