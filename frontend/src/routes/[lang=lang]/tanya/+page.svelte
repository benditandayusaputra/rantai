<script lang="ts">
	import { LoaderCircle, MessageCircle } from '@lucide/svelte';
	import PercakapanTanya from '$lib/features/tanya/PercakapanTanya.svelte';
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{t('ask.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="wrap pt-8 lg:pt-12">
	<div class="mx-auto max-w-chat">
		<header class="flex items-center gap-4 border-b border-line pb-6">
			<span class="icon-tile h-12 w-12"
				><MessageCircle size={24} strokeWidth={1.75} aria-hidden="true" /></span
			>
			<h1 class="t-h1">{t('ask.title')}</h1>
		</header>
		<div class="mt-8">
			{#await data.wilayah}
				<p aria-busy="true" class="flex items-center gap-2 text-muted">
					<LoaderCircle size={18} strokeWidth={1.75} aria-hidden="true" class="animate-spin" />
					{t('ask.thinking')}
				</p>
			{:then hasil}
				<PercakapanTanya
					konteks={hasil?.ok ? { kode: hasil.data.kode, nama: hasil.data.nama } : null}
				/>
			{/await}
		</div>
	</div>
</div>
