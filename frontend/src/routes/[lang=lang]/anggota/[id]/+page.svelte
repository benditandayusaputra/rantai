<script lang="ts">
	import ErrorState from '$lib/components/ErrorState.svelte';
	import HalamanAnggota from '$lib/features/anggota/HalamanAnggota.svelte';
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

{#await data.anggota}
	<p class="wrap py-10 text-muted" aria-live="polite">{t('state.loading')}</p>
{:then hasil}
	{#if hasil.ok}
		<HalamanAnggota anggota={hasil.data} />
	{:else}
		<div class="wrap py-10"><ErrorState title={t('member.notFound')} /></div>
	{/if}
{/await}
