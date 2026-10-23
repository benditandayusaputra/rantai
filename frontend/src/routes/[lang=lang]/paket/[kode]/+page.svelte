<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import HalamanPaket from '$lib/features/paket/HalamanPaket.svelte';
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	{#await data.paket then hasil}
		{#if hasil.ok}<title>{t('pkg.title', { name: hasil.data.nama })} · {t('brand.name')}</title
			>{/if}
	{/await}
</svelte:head>

{#await data.paket}
	<div class="wrap flex flex-col gap-4 pt-8" aria-busy="true">
		<p class="sr-only">{t('state.loadingSummary')}</p>
		<div class="skel h-10 w-3/4 rounded-sm"></div>
		<div class="grid grid-cols-2 gap-3 md:grid-cols-4">
			{#each [0, 1, 2, 3] as urutan (urutan)}<div class="skel h-28 rounded-lg"></div>{/each}
		</div>
		<div class="skel h-64 rounded-lg"></div>
	</div>
{:then hasil}
	{#if hasil.ok}
		<HalamanPaket paket={hasil.data} />
	{:else}
		<div class="wrap pt-10">
			<h1 class="sr-only">{t('pkg.notFound')}</h1>
			{#if hasil.galat.kode === 404}
				<EmptyState title={t('pkg.notFound')} />
			{:else}
				<ErrorState title={t('error.sectionFailed')} onretry={invalidateAll} />
			{/if}
		</div>
	{/if}
{/await}
