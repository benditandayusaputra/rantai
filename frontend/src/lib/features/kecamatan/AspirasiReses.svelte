<script lang="ts">
	import { ExternalLink, Info, MapPin } from '@lucide/svelte';
	import { ambilAspirasiReses, tangkap } from '$lib/api/klien';
	import type { AspirasiReses } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import PaketTerkait from './PaketTerkait.svelte';
	import { t } from '$lib/i18n';

	let { idDapil }: { idDapil: number } = $props();

	let daftar = $state<AspirasiReses[] | null>(null);
	let isGagal = $state(false);
	let isSemua = $state(false);
	const batasAwal = 6;

	$effect(() => {
		void tangkap(ambilAspirasiReses(fetch, idDapil)).then((hasil) => {
			if (hasil.ok) daftar = hasil.data;
			else isGagal = true;
		});
	});

	const tampil = $derived(isSemua ? (daftar ?? []) : (daftar ?? []).slice(0, batasAwal));
	const perMasa = $derived(
		Object.entries(
			tampil.reduce<Record<string, AspirasiReses[]>>((kelompok, a) => {
				(kelompok[a.masa_reses] ??= []).push(a);
				return kelompok;
			}, {})
		)
	);
</script>

<section aria-labelledby="judul-reses" class="flex flex-col gap-4">
	<h3 id="judul-reses" class="t-h3">{t('reses.title')}</h3>
	<p class="flex gap-2 rounded-md bg-info-bg p-3 text-small">
		<Info size={18} strokeWidth={1.75} class="mt-0.5 shrink-0 text-info" aria-hidden="true" />
		{t('reses.note')}
	</p>
	{#if isGagal}
		<p class="text-small text-muted">{t('reses.failed')}</p>
	{:else if daftar === null}
		<p class="text-muted" aria-live="polite">{t('state.loading')}</p>
	{:else if daftar.length === 0}
		<EmptyState title={t('reses.empty')} />
	{:else}
		{#each perMasa as [masa, item] (masa)}
			<div class="flex flex-col gap-2">
				<p class="t-over">{masa}</p>
				<ul class="flex flex-col gap-2">
					{#each item as a (a.id)}
						<li class="flex flex-col gap-1 rounded-lg border border-line bg-surface p-4">
							<p class="flex flex-wrap items-center gap-2 text-small">
								<span class="rounded-xs bg-sunken px-2 py-0.5 text-caption font-bold uppercase"
									>{a.bidang}</span
								>
								<span class="inline-flex items-center gap-1 text-muted">
									<MapPin size={14} strokeWidth={1.75} aria-hidden="true" />{a.wilayah?.nama ??
										t('reses.allJakarta')}
								</span>
							</p>
							<p>{a.uraian}</p>
							{#if a.wilayah}<PaketTerkait idAspirasi={a.id} />{/if}
							<a
								href={a.url}
								target="_blank"
								rel="noopener"
								class="inline-flex items-center gap-1 self-start text-small font-semibold underline"
							>
								{t('stance.source')}<ExternalLink
									size={14}
									strokeWidth={1.75}
									aria-hidden="true"
								/><span class="sr-only">{t('a11y.opensNewTab')}</span>
							</a>
						</li>
					{/each}
				</ul>
			</div>
		{/each}
		{#if !isSemua && daftar.length > batasAwal}
			<button
				type="button"
				class="inline-flex min-h-touch items-center self-start rounded-md border-ctrl border-line-strong px-4 font-semibold hover:bg-sunken"
				onclick={() => (isSemua = true)}>{t('action.viewAll', { n: daftar.length })}</button
			>
		{/if}
	{/if}
</section>
