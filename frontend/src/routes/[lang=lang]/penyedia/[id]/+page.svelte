<script lang="ts">
	import ErrorState from '$lib/components/ErrorState.svelte';
	import HalamanPenyedia from '$lib/features/penyedia/HalamanPenyedia.svelte';
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

{#await data.penyedia}
	<div class="wrap flex flex-col gap-4 py-6 lg:py-10" aria-busy="true">
		<p class="sr-only" aria-live="polite">{t('state.loading')}</p>
		<div class="skel h-16 w-2/3 rounded-md"></div>
		<div class="grid grid-cols-2 gap-3 lg:max-w-prose">
			<div class="skel h-28 rounded-lg"></div>
			<div class="skel h-28 rounded-lg"></div>
		</div>
		<div class="skel h-64 rounded-lg"></div>
	</div>
{:then hasil}
	{#if hasil.ok}
		<HalamanPenyedia penyedia={hasil.data} />
	{:else}
		<div class="wrap py-10"><ErrorState title={t('vendor.notFound')} /></div>
	{/if}
{/await}
