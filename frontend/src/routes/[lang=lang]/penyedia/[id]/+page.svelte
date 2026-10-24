<script lang="ts">
	import ErrorState from '$lib/components/ErrorState.svelte';
	import HalamanPenyedia from '$lib/features/penyedia/HalamanPenyedia.svelte';
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

{#await data.penyedia}
	<p class="wrap py-10 text-muted" aria-live="polite">{t('state.loading')}</p>
{:then hasil}
	{#if hasil.ok}
		<HalamanPenyedia penyedia={hasil.data} />
	{:else}
		<div class="wrap py-10"><ErrorState title={t('vendor.notFound')} /></div>
	{/if}
{/await}
