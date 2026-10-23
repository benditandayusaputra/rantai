import { describe, expect, it } from 'vitest';
import { t } from '.';

describe('t', () => {
	it('mengisi parameter bernama', () => {
		expect(t('search.loading', { q: 'gading' })).toBe('Mencari "gading"...');
	});

	it('membiarkan parameter yang tidak diberikan', () => {
		expect(t('chain.summary', { ok: 5 })).toBe('5 dari {total} mata rantai tersambung');
	});

	it('mengembalikan pesan tanpa parameter apa adanya', () => {
		expect(t('nav.main')).toBe('Navigasi utama');
	});
});
