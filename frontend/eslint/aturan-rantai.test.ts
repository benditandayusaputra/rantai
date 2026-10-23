import { Linter } from 'eslint';
import svelteParser from 'svelte-eslint-parser';
import { describe, expect, it } from 'vitest';
import aturanRantai from './aturan-rantai.js';

const linter = new Linter();

function periksa(kode: string, aturan: string): string[] {
	const pesan = linter.verify(
		kode,
		{
			files: ['**/*.svelte'],
			languageOptions: { parser: svelteParser },
			plugins: { rantai: aturanRantai },
			rules: { [`rantai/${aturan}`]: 'error' }
		},
		'Komponen.svelte'
	);
	return pesan.map((p) => p.message);
}

describe('aturan rantai', () => {
	it('menolak komentar skrip dan HTML', () => {
		const kode = '<script>\n// catatan\nlet a = 1;\n</script>\n<!-- catatan -->\n<p>{a}</p>';
		expect(periksa(kode, 'tanpa-komentar')).toHaveLength(2);
	});

	it('menolak warna heks dan ukuran arbitrer', () => {
		const kode = '<div class="text-[#ff0000] w-[13px]" style="color: rgb(1,2,3)"></div>';
		expect(periksa(kode, 'tanpa-nilai-mentah')).toHaveLength(3);
	});

	it('menerima kelas token', () => {
		expect(
			periksa('<a href="#konten" class="bg-primary p-4">{x}</a>', 'tanpa-nilai-mentah')
		).toEqual([]);
	});

	it('menolak teks mentah di templat dan atribut', () => {
		const kode = '<button aria-label="Tutup">Kirim</button>';
		expect(periksa(kode, 'tanpa-teks-mentah')).toHaveLength(2);
	});

	it('menerima teks dari i18n', () => {
		const kode = "<button aria-label={t('action.close')}>{t('action.next')}</button>";
		expect(periksa(kode, 'tanpa-teks-mentah')).toEqual([]);
	});
});
