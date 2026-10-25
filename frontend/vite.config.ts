import { sveltekit } from '@sveltejs/kit/vite';
import { svelteTesting } from '@testing-library/svelte/vite';
import { defineConfig } from 'vitest/config';

const potonganIkon = (id: string) => (id.includes('/@lucide/svelte/') ? 'ikon' : undefined);

export default defineConfig({
	plugins: [sveltekit(), svelteTesting()],
	build: {
		rollupOptions: { output: { manualChunks: potonganIkon } }
	},
	test: {
		environment: 'jsdom',
		include: ['src/**/*.test.ts', 'eslint/**/*.test.ts'],
		setupFiles: ['./src/vitest-setup.ts']
	}
});
