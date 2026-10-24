import { render, screen } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import ReportStatusChip from './ReportStatusChip.svelte';

describe('ReportStatusChip', () => {
	it('menampilkan hitungan saksi hanya bila datanya ada', () => {
		const { unmount } = render(ReportStatusChip, {
			status: 'MENUNGGU_SAKSI',
			witnesses: { jumlah: 2, dibutuhkan: 3 }
		});
		expect(screen.getByText('Menunggu saksi 2/3')).toBeInTheDocument();
		unmount();

		render(ReportStatusChip, { status: 'MENUNGGU_SAKSI' });
		expect(screen.getByText('Menunggu saksi')).toBeInTheDocument();
	});
});
