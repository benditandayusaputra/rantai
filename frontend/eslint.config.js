import js from '@eslint/js';
import svelte from 'eslint-plugin-svelte';
import globals from 'globals';
import ts from 'typescript-eslint';
import aturanRantai from './eslint/aturan-rantai.js';

export default ts.config(
	{ ignores: ['.svelte-kit/', 'build/', 'node_modules/'] },
	js.configs.recommended,
	...ts.configs.strict,
	...svelte.configs.recommended,
	{
		languageOptions: { globals: { ...globals.browser, ...globals.node } }
	},
	{
		files: ['**/*.svelte', '**/*.svelte.ts'],
		languageOptions: {
			parserOptions: { parser: ts.parser, extraFileExtensions: ['.svelte'] }
		},
		rules: { 'no-useless-assignment': 'off' }
	},
	{
		plugins: { rantai: aturanRantai },
		rules: {
			'rantai/tanpa-komentar': 'error',
			'@typescript-eslint/no-explicit-any': 'error',
			'svelte/no-navigation-without-resolve': [
				'error',
				{ ignoreLinks: true, ignoreGoto: true, ignoreReplaceState: true }
			]
		}
	},
	{
		files: ['src/**/*.svelte'],
		rules: {
			'rantai/tanpa-nilai-mentah': 'error',
			'rantai/tanpa-teks-mentah': 'error'
		}
	}
);
