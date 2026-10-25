<script lang="ts">
	import { ArrowRight, Sigma } from '@lucide/svelte';
	import type { DetailPendidikan } from '$lib/api/tipe';
	import { formatAngka } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import { susunRumus } from './pendidikan';

	let { detail }: { detail: DetailPendidikan } = $props();

	const id = $props.id();
	const desimalAmbang = 1;
	const baris = $derived(susunRumus(detail));
	const ambang = $derived(detail.rumus.ambang_rasio);
	const formatAmbang = (nilai: number) => formatAngka(nilai, bahasaAktif, desimalAmbang);
</script>

<section aria-labelledby="{id}-judul" class="card flex flex-col gap-4 p-4 md:p-6">
	<div class="flex items-center gap-3">
		<span class="icon-tile"><Sigma size={20} strokeWidth={1.75} aria-hidden="true" /></span>
		<h2 id="{id}-judul" class="t-h3">{t('pen.formula.title')}</h2>
	</div>
	<p class="max-w-prose text-small text-muted">{t('pen.formula.lead')}</p>
	<dl class="flex flex-col divide-y divide-line rounded-lg bg-sunken">
		{#each baris as item (item.kunci)}
			<div
				class="flex flex-col gap-1 p-3 md:flex-row md:items-baseline md:justify-between md:gap-6"
			>
				<dt class="text-small">{t(item.kunci, item.parameter)}</dt>
				<dd class="num font-mono text-mono font-bold md:text-right">
					{#if item.hitungan}{item.hitungan}{:else}<span class="font-normal text-muted"
							>{t('pen.formula.unavailable')}</span
						>{/if}
				</dd>
			</div>
		{/each}
	</dl>
	<p class="text-small">
		{t('pen.formula.class', {
			a: formatAmbang(ambang.sangat_kurang),
			b: formatAmbang(ambang.kurang),
			c: formatAmbang(ambang.hampir_cukup)
		})}
	</p>
	<a
		href="/{bahasaAktif}/tentang-data#rumus"
		class="link inline-flex items-center gap-1 self-start text-small"
	>
		{t('pen.formula.more')}
		<ArrowRight size={16} strokeWidth={1.75} aria-hidden="true" />
	</a>
</section>
