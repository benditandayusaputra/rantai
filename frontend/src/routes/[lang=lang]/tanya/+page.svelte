<script lang="ts">
	import PercakapanTanya from '$lib/features/tanya/PercakapanTanya.svelte';
	import { t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{t('ask.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="wrap pt-6 lg:pt-10">
	<div class="mx-auto max-w-chat">
		<h1 class="t-h1">{t('ask.title')}</h1>
		<div class="mt-6">
			{#await data.wilayah}
				<p aria-busy="true" class="text-muted">{t('ask.thinking')}</p>
			{:then hasil}
				<PercakapanTanya
					konteks={hasil?.ok ? { kode: hasil.data.kode, nama: hasil.data.nama } : null}
				/>
			{/await}
		</div>
	</div>
</div>
