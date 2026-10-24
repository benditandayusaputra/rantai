import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import OtpInput from './OtpInput.svelte';

describe('OtpInput', () => {
	it('menerima tempelan kode utuh, membuang selain angka, lalu memanggil oncomplete', async () => {
		const pengguna = userEvent.setup();
		const oncomplete = vi.fn();
		render(OtpInput, { id: 'kode', value: '', oncomplete });
		const masukan = screen.getByRole('textbox');

		await pengguna.click(masukan);
		await pengguna.paste('12 34-56');

		expect(masukan).toHaveValue('123456');
		expect(oncomplete).toHaveBeenCalledWith('123456');
	});

	it('menandai error dan menonaktifkan isian saat terkunci', () => {
		render(OtpInput, { id: 'kode', value: '', error: 'invalid', locked: true });
		const masukan = screen.getByRole('textbox');

		expect(masukan).toBeDisabled();
		expect(masukan).toHaveAttribute('aria-invalid', 'true');
	});
});
