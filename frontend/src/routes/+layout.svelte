<script lang="ts">
	import '@fontsource-variable/plus-jakarta-sans';
	import '@fontsource-variable/newsreader';
	import '@fontsource/jetbrains-mono/400.css';
	import '@fontsource/jetbrains-mono/500.css';
	import fontIsi from '@fontsource-variable/plus-jakarta-sans/files/plus-jakarta-sans-latin-wght-normal.woff2?url';
	import fontJudul from '@fontsource-variable/newsreader/files/newsreader-latin-wght-normal.woff2?url';
	import '../app.css';
	import { ikutiTemaSistem } from '$lib/features/tema/tema';

	let { children } = $props();

	$effect(ikutiTemaSistem);
</script>

<svelte:head>
	<link rel="preload" as="font" type="font/woff2" href={fontIsi} crossorigin="anonymous" />
	<link rel="preload" as="font" type="font/woff2" href={fontJudul} crossorigin="anonymous" />
</svelte:head>

{@render children()}
