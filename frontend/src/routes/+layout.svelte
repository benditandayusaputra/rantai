<script lang="ts">
	import '@fontsource-variable/inter';
	import '@fontsource-variable/plus-jakarta-sans';
	import '@fontsource/jetbrains-mono/400.css';
	import '@fontsource/jetbrains-mono/500.css';
	import fontIsi from '@fontsource-variable/inter/files/inter-latin-wght-normal.woff2?url';
	import fontJudul from '@fontsource-variable/plus-jakarta-sans/files/plus-jakarta-sans-latin-wght-normal.woff2?url';
	import '../app.css';
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import { navigating, page } from '$app/state';
	import { pulihkanWilayahku } from '$lib/features/wilayah/wilayahku.svelte';
	import { aturBahasa, bahasaDari } from '$lib/i18n';

	let { children } = $props();

	const bahasa = $derived(bahasaDari(page.params.lang));
	aturBahasa(bahasaDari(page.params.lang));

	onMount(pulihkanWilayahku);

	$effect.pre(() => aturBahasa(bahasa));
	$effect(() => {
		document.documentElement.lang = bahasa;
	});
</script>

<svelte:head>
	<link rel="preload" as="font" type="font/woff2" href={fontIsi} crossorigin="anonymous" />
	<link rel="preload" as="font" type="font/woff2" href={fontJudul} crossorigin="anonymous" />
</svelte:head>

{#if navigating.to}
	<div
		class="fixed inset-x-0 top-0 z-tooltip h-0.75 origin-left animate-muat bg-primary motion-reduce:animate-none"
		aria-hidden="true"
		out:fade={{ duration: 200 }}
	></div>
{/if}

{#key bahasa}
	{@render children()}
{/key}
