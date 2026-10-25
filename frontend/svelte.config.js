import adapterNode from '@sveltejs/adapter-node';
import adapterVercel from '@sveltejs/adapter-vercel';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

const regionGulirBisaDifokus = (peringatan) =>
	peringatan.code === 'a11y_no_noninteractive_tabindex' &&
	peringatan.filename?.endsWith('components/TableScroll.svelte');

export default {
	preprocess: vitePreprocess(),
	compilerOptions: {
		warningFilter: (peringatan) => !regionGulirBisaDifokus(peringatan)
	},
	kit: {
		adapter: process.env.VERCEL ? adapterVercel({ runtime: 'nodejs22.x' }) : adapterNode(),
		env: { dir: '..' }
	}
};
