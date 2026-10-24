<script lang="ts">
	import { ExternalLink } from '@lucide/svelte';
	import type { ItemUmpan } from '$lib/api/tipe';
	import RightOfReply from '$lib/components/RightOfReply.svelte';
	import { formatTanggal } from '$lib/format';
	import { bahasaBawaan, t } from '$lib/i18n';
	import LencanaSumber from './LencanaSumber.svelte';

	let { item }: { item: ItemUmpan } = $props();
</script>

<li class="flex flex-col gap-2 rounded-lg border border-line bg-surface p-4">
	<div class="flex flex-wrap items-center gap-2">
		<LencanaSumber {item} />
		<span class="text-small text-muted num"
			>{formatTanggal(new Date(item.terbit_pada), bahasaBawaan, 'long')}</span
		>
	</div>
	<a
		href={item.url}
		target="_blank"
		rel="noopener"
		class="inline-flex items-start gap-1 font-semibold underline"
	>
		{item.judul}<ExternalLink
			size={16}
			strokeWidth={1.75}
			class="mt-1 shrink-0"
			aria-hidden="true"
		/><span class="sr-only">{t('a11y.opensNewTab')}</span>
	</a>
	{#if item.ringkasan}<p class="text-small text-muted">{item.ringkasan}</p>{/if}
	{#each item.hak_jawab ?? [] as hakJawab (hakJawab.id)}
		<RightOfReply {hakJawab} />
	{/each}
</li>
