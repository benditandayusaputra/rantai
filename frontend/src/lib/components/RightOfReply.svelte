<script lang="ts">
	import { FileText, Scale } from '@lucide/svelte';
	import type { HakJawab } from '$lib/api/tipe';
	import { formatTanggal } from '$lib/format';
	import { bahasaBawaan, t } from '$lib/i18n';
	import SimRibbon from './SimRibbon.svelte';

	let { hakJawab }: { hakJawab: HakJawab } = $props();

	const idJudul = $derived(`hak-jawab-${hakJawab.id}`);
	const isAnggota = $derived(hakJawab.target_jenis === 'umpan');
</script>

<article aria-labelledby={idJudul} class="overflow-hidden rounded-lg border border-line bg-surface">
	{#if hakJawab.is_simulation}<SimRibbon />{/if}
	<p class="flex items-center gap-2 bg-info-bg px-4 py-2 text-small font-semibold text-info">
		<Scale size={16} strokeWidth={1.75} aria-hidden="true" />
		{isAnggota ? t('reply.memberTitle') : t('reply.title')}
		<span class="font-normal">·</span>
		<time datetime={hakJawab.dibuat_pada} class="font-normal"
			>{formatTanggal(new Date(hakJawab.dibuat_pada), bahasaBawaan, 'long')}</time
		>
	</p>
	<div class="flex flex-col gap-3 p-4">
		<div>
			<h3 id={idJudul} class="font-bold">{hakJawab.pihak}</h3>
			<p class="text-small text-muted">
				{isAnggota ? t('reply.memberVerified') : t('reply.verified')}{hakJawab.nama_penjawab
					? ` · ${hakJawab.nama_penjawab}`
					: ''}
			</p>
		</div>
		<blockquote class="whitespace-pre-line font-serif text-body-lg">{hakJawab.isi}</blockquote>
		{#if hakJawab.tautan_dokumen}
			<a
				href={hakJawab.tautan_dokumen}
				target="_blank"
				rel="noopener noreferrer"
				class="inline-flex min-h-touch items-center gap-2 self-start font-semibold underline"
			>
				<FileText size={16} strokeWidth={1.75} aria-hidden="true" />{t('reply.attachment')}
				<span class="sr-only">{t('a11y.opensNewTab')}</span>
			</a>
		{/if}
		<p class="text-caption text-muted">{t('reply.unedited')}</p>
	</div>
</article>
