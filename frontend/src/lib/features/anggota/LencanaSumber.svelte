<script lang="ts">
	import { Landmark, Newspaper } from '@lucide/svelte';
	import type { ItemUmpan } from '$lib/api/tipe';
	import { t } from '$lib/i18n';

	let { item }: { item: Pick<ItemUmpan, 'jenis_sumber' | 'media_nama'> } = $props();
</script>

{#if item.jenis_sumber === 'resmi_dprd'}
	<span
		class="inline-flex items-center gap-1 self-start rounded-xs bg-cert-resmi px-2 py-0.5 text-caption font-bold text-cert-resmi-on"
	>
		<Landmark size={14} strokeWidth={1.75} aria-hidden="true" />{t('member.source.resmi_dprd')}
	</span>
{:else}
	<span
		class="inline-flex items-center gap-1 self-start rounded-xs border-ctrl border-line-strong px-2 py-0.5 text-caption font-bold"
	>
		<Newspaper size={14} strokeWidth={1.75} aria-hidden="true" />{t('member.source.media', {
			nama: item.media_nama
		})}
	</span>
{/if}
