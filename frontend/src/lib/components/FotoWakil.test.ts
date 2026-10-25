import { fireEvent, render } from '@testing-library/svelte';
import { describe, expect, it } from 'vitest';
import FotoWakil from './FotoWakil.svelte';

const nama = 'Fatimah Tania Nadira Alatas';

describe('FotoWakil', () => {
	it('menampilkan foto bila tersedia', () => {
		const { container } = render(FotoWakil, { nama, foto: 'https://contoh.test/foto.jpg' });

		expect(container.querySelector('img')).toHaveAttribute('src', 'https://contoh.test/foto.jpg');
	});

	it('kembali ke inisial saat foto gagal dimuat', async () => {
		const { container } = render(FotoWakil, { nama, foto: 'https://contoh.test/rusak.jpg' });

		await fireEvent.error(container.querySelector('img') as HTMLImageElement);

		expect(container.querySelector('img')).toBeNull();
		expect(container).toHaveTextContent('FT');
	});

	it('memakai inisial tanpa foto', () => {
		const { container } = render(FotoWakil, { nama });

		expect(container.querySelector('img')).toBeNull();
		expect(container).toHaveTextContent('FT');
	});
});
