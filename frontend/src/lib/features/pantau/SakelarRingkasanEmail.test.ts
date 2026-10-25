import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import SakelarRingkasanEmail from './SakelarRingkasanEmail.svelte';

const api = vi.hoisted(() => ({
	ambilRingkasanEmail: vi.fn(async () => ({ aktif: false })),
	aktifkanRingkasanEmail: vi.fn(async () => ({ aktif: true })),
	matikanRingkasanEmail: vi.fn(async () => ({ aktif: false }))
}));

vi.mock('$lib/api/pantau', () => api);

describe('SakelarRingkasanEmail', () => {
	it('memuat status lalu menyalakan ringkasan saat sakelar diklik', async () => {
		render(SakelarRingkasanEmail, { token: 'token-uji' });
		const sakelar = await screen.findByRole('switch');
		expect(sakelar).not.toBeChecked();

		await userEvent.click(sakelar);

		expect(api.aktifkanRingkasanEmail).toHaveBeenCalledWith(expect.any(Function), 'token-uji');
		expect(await screen.findByRole('switch')).toBeChecked();
	});
});
