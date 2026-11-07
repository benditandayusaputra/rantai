import { describe, expect, it } from 'vitest';
import { isIsiSah, isTautanSah } from './temuan';

describe('temuan', () => {
	it.each([
		['https://lpse.jakarta.go.id/x', true],
		['http://contoh.id', true],
		['javascript:alert(1)', false],
		['lpse.jakarta.go.id', false],
		['', false]
	])('tautan %s sah: %s', (tautan, harap) => {
		expect(isTautanSah(tautan)).toBe(harap);
	});

	it('menolak isi terlalu pendek setelah dipangkas', () => {
		expect(isIsiSah('   pendek   ')).toBe(false);
		expect(isIsiSah('Papan proyek menyebut nilai berbeda.')).toBe(true);
	});
});
