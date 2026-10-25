<script lang="ts">
	import '@fontsource-variable/plus-jakarta-sans';
	import '@fontsource-variable/newsreader';
	import '@fontsource/jetbrains-mono/400.css';
	import '@fontsource/jetbrains-mono/500.css';
	import fontIsi from '@fontsource-variable/plus-jakarta-sans/files/plus-jakarta-sans-latin-wght-normal.woff2?url';
	import fontJudul from '@fontsource-variable/newsreader/files/newsreader-latin-wght-normal.woff2?url';
	import '../app.css';
	import { page } from '$app/state';
	import { ikutiTemaSistem } from '$lib/features/tema/tema';
	import { aturBahasa, bahasaDari } from '$lib/i18n';

	let { children } = $props();

	const bahasa = $derived(bahasaDari(page.params.lang));
	aturBahasa(bahasaDari(page.params.lang));

	$effect.pre(() => aturBahasa(bahasa));
	$effect(() => {
		document.documentElement.lang = bahasa;
	});
	$effect(ikutiTemaSistem);
</script>

<svelte:head>
	<link rel="preload" as="font" type="font/woff2" href={fontIsi} crossorigin="anonymous" />
	<link rel="preload" as="font" type="font/woff2" href={fontJudul} crossorigin="anonymous" />
</svelte:head>

{#key bahasa}
	{@render children()}
{/key}
