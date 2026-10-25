import { afterEach, describe, expect, it } from 'vitest';
import { kunciTema, pilihTema, temaLawan, temaTerpasang } from './tema';

afterEach(() => {
	localStorage.clear();
	document.documentElement.removeAttribute('data-theme');
});

describe('pilihTema', () => {
	it('memasang tema di dokumen dan mengingatnya', () => {
		pilihTema('dark');
		expect(temaTerpasang()).toBe('dark');
		expect(localStorage.getItem(kunciTema)).toBe('dark');
	});

	it('berganti ke tema lawan dari yang terpasang', () => {
		pilihTema('light');
		pilihTema(temaLawan(temaTerpasang()));
		expect(temaTerpasang()).toBe('dark');
	});
});
