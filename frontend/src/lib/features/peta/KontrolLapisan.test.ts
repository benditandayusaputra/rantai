import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import KontrolLapisan from './KontrolLapisan.svelte';

describe('KontrolLapisan', () => {
	it('menyalakan dan mematikan setiap lapisan dengan papan ketik', async () => {
		const pengguna = userEvent.setup();
		const onchange = vi.fn();
		const lapisan = { tema: true, proyek: false, sekolah: false, laporan: false };
		render(KontrolLapisan, { lapisan, onchange });
		const kotak = screen.getAllByRole('checkbox');
		expect(kotak).toHaveLength(4);
		for (const centang of kotak) {
			await pengguna.tab();
			expect(centang).toHaveFocus();
			const sebelum = (centang as HTMLInputElement).checked;
			await pengguna.keyboard(' ');
			expect((centang as HTMLInputElement).checked).toBe(!sebelum);
		}
		expect(onchange).toHaveBeenCalledTimes(4);
	});
});
