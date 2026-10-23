import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import type { WilayahDitemukan } from '$lib/api/tipe';
import RegionSearch from './RegionSearch.svelte';

const provinsi = { kode: '31', nama: 'DKI Jakarta', tingkat: 'provinsi' } as const;
const kecamatan = { kode: '31.72.06', nama: 'Kelapa Gading', tingkat: 'kecamatan' } as const;
const hasil: WilayahDitemukan[] = [
	{ ...kecamatan, hierarki: [provinsi] },
	{
		kode: '31.72.06.1001',
		nama: 'Kelapa Gading Barat',
		tingkat: 'kelurahan',
		hierarki: [provinsi, kecamatan]
	}
];

describe('RegionSearch', () => {
	it('memilih saran dengan panah dan Enter, kelurahan lebih dahulu', async () => {
		const pengguna = userEvent.setup();
		const onselect = vi.fn();
		render(RegionSearch, { cari: async () => hasil, onselect });
		const kotak = screen.getByRole('combobox');
		await pengguna.type(kotak, 'gading');
		expect(await screen.findByRole('option', { name: /Kelapa Gading Barat/ })).toBeVisible();
		await pengguna.keyboard('{ArrowDown}{Enter}');
		expect(onselect).toHaveBeenCalledWith(hasil[1]);
	});

	it('memberi tahu saat tidak ada wilayah yang cocok', async () => {
		const pengguna = userEvent.setup();
		render(RegionSearch, { cari: async () => [], onselect: vi.fn() });
		await pengguna.type(screen.getByRole('combobox'), 'xyz');
		const pesan = await screen.findAllByText('Tidak ada wilayah bernama "xyz"');
		expect(pesan[0]).toBeVisible();
	});
});
