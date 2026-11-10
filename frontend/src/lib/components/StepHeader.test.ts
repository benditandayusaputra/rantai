import { render, screen } from '@testing-library/svelte';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import StepHeader from './StepHeader.svelte';

describe('StepHeader', () => {
	it('mengumumkan kemajuan langkah dan memakai tombol kembali setelah langkah 1', async () => {
		const onback = vi.fn();
		render(StepHeader, { step: 2, stepLabel: 'pilih proyek', onback, onclose: vi.fn() });

		const kemajuan = screen.getByRole('progressbar');
		expect(kemajuan).toHaveAttribute('aria-valuenow', '2');
		expect(kemajuan).toHaveAttribute('aria-valuetext', 'Langkah 2 dari 6, pilih proyek');
		await userEvent.setup().click(screen.getByRole('button', { name: 'Kembali' }));
		expect(onback).toHaveBeenCalled();
	});

	it('tetap menyediakan tombol tutup setelah langkah 1', async () => {
		const onclose = vi.fn();
		render(StepHeader, { step: 6, stepLabel: 'ringkasan', onback: vi.fn(), onclose });

		await userEvent.setup().click(screen.getByRole('button', { name: 'Tutup' }));
		expect(onclose).toHaveBeenCalled();
	});

	it('memakai tombol tutup di langkah 1', () => {
		render(StepHeader, { step: 1, stepLabel: 'izin lokasi', onback: vi.fn(), onclose: vi.fn() });

		expect(screen.getByRole('button', { name: 'Tutup' })).toBeInTheDocument();
	});
});
