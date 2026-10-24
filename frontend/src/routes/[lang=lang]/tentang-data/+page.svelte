<script lang="ts">
	import { CircleHelp } from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import type { JenisSinyal, Label } from '$lib/api/tipe';
	import CertaintyBadge from '$lib/components/CertaintyBadge.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import { namaJob, pemilikJob } from '$lib/features/tentang/tentang';
	import { formatJam, formatTanggal } from '$lib/format';
	import { bahasaBawaan, t, type KunciPesan } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const label: Label[] = ['resmi', 'proksi', 'estimasi', 'simulasi'];
	const sinyal: JenisSinyal[] = ['S-01', 'S-02', 'S-03', 'S-04', 'S-05'];
	const daftarIsi: { id: string; judul: KunciPesan }[] = [
		{ id: 'sumber', judul: 'about.sources' },
		{ id: 'label', judul: 'about.labels' },
		{ id: 'rumus', judul: 'about.formula' },
		{ id: 'lokasi', judul: 'about.location' },
		{ id: 'sinyal', judul: 'about.signals' },
		{ id: 'media', judul: 'about.media' },
		{ id: 'batasan', judul: 'about.limits' },
		{ id: 'lisensi', judul: 'about.license' }
	];
	const aturanLokasi: KunciPesan[] = [
		'about.location.1',
		'about.location.2',
		'about.location.3',
		'about.location.4'
	];
	const keterbatasan: KunciPesan[] = [
		'about.limit.1',
		'about.limit.2',
		'about.limit.3',
		'about.limit.4'
	];

	function waktu(iso: string): string {
		const d = new Date(iso);
		return t('source.retrievedAt', { date: formatTanggal(d, bahasaBawaan), time: formatJam(d) });
	}

	const kelasBagian = 'flex scroll-mt-24 flex-col gap-4';
</script>

<svelte:head>
	<title>{t('about.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="wrap flex flex-col gap-14 pt-10 lg:pt-16">
	<header class="flex max-w-3xl flex-col gap-4">
		<p class="t-over">{t('about.overline')}</p>
		<h1 class="t-display">{t('about.title')}</h1>
		<p class="t-lead">{t('about.lead')}</p>
		<nav aria-label={t('about.toc')}>
			<ul class="flex gap-2 overflow-x-auto pb-1 lg:flex-wrap">
				{#each daftarIsi as item (item.id)}
					<li class="shrink-0">
						<a
							href="#{item.id}"
							class="inline-flex min-h-touch items-center rounded-pill border-ctrl border-line-strong px-4 text-small font-semibold hover:border-fg hover:bg-sunken"
							>{t(item.judul)}</a
						>
					</li>
				{/each}
			</ul>
		</nav>
	</header>

	<section id="sumber" aria-labelledby="judul-sumber" class={kelasBagian}>
		<h2 id="judul-sumber" class="t-h2">{t('about.sources')}</h2>
		{#await data.sumber}
			<div class="skel h-48 rounded-lg" aria-busy="true"></div>
		{:then hasil}
			{#if !hasil.ok}
				<ErrorState title={t('about.sourcesFailed')} onretry={invalidateAll} />
			{:else if !hasil.data.pengambilan.length}
				<EmptyState title={t('about.sourcesEmpty')} />
			{:else}
				<div class="overflow-x-auto rounded-lg border border-line bg-surface">
					<table class="w-full min-w-max text-small">
						<caption class="sr-only">{t('about.sources')}</caption>
						<thead class="bg-sunken text-caption uppercase">
							<tr>
								<th scope="col" class="sticky left-0 bg-sunken p-3 text-left"
									>{t('about.col.source')}</th
								>
								<th scope="col" class="p-3 text-left">{t('about.col.owner')}</th>
								<th scope="col" class="p-3 text-left">{t('about.col.lastOk')}</th>
								<th scope="col" class="p-3 text-left">{t('about.col.status')}</th>
							</tr>
						</thead>
						<tbody>
							{#each hasil.data.pengambilan as item (item.job)}
								<tr class="border-b border-line">
									<th scope="row" class="sticky left-0 bg-surface p-3 text-left font-semibold"
										>{namaJob(item.job)}</th
									>
									<td class="p-3">{pemilikJob(item.job)}</td>
									<td class="num p-3">
										{item.berhasil_pada ? waktu(item.berhasil_pada) : t('about.never')}
									</td>
									<td class="p-3">{t(`about.status.${item.status_terakhir}`)}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		{/await}
	</section>

	<section id="label" aria-labelledby="judul-label" class={kelasBagian}>
		<h2 id="judul-label" class="t-h2">{t('about.labels')}</h2>
		<ul class="grid gap-3 md:grid-cols-2 lg:grid-cols-4 lg:gap-5">
			{#each label as level (level)}
				<li class="flex flex-col gap-3 rounded-lg border border-line bg-surface p-4 md:p-5">
					<CertaintyBadge {level} size="md" class="self-start" />
					<p class="text-small">{t(`cert.${level}.desc`)}</p>
				</li>
			{/each}
		</ul>
		<p class="max-w-prose text-small text-muted">{t('about.weakest')}</p>
	</section>

	<div class="grid gap-14 lg:grid-cols-12 lg:gap-5">
		<section id="rumus" aria-labelledby="judul-rumus" class="{kelasBagian} lg:col-span-8">
			<h2 id="judul-rumus" class="t-h2">{t('about.formula')}</h2>
			<pre class="overflow-x-auto rounded-lg bg-sunken p-4 font-mono text-mono">{t(
					'about.formula.code'
				)}</pre>
			<p class="text-small text-muted">{t('about.formulaNote')}</p>
		</section>
		<section id="lokasi" aria-labelledby="judul-lokasi" class="{kelasBagian} lg:col-span-4">
			<h2 id="judul-lokasi" class="t-h3">{t('about.location')}</h2>
			<ol class="list-decimal space-y-2 pl-5 text-small">
				{#each aturanLokasi as kunci (kunci)}<li>{t(kunci)}</li>{/each}
			</ol>
		</section>
	</div>

	<section id="sinyal" aria-labelledby="judul-sinyal" class={kelasBagian}>
		<h2 id="judul-sinyal" class="t-h2">{t('about.signals')}</h2>
		<p class="max-w-prose text-muted">{t('about.signals.lead')}</p>
		<dl class="grid gap-3 md:grid-cols-2">
			{#each sinyal as jenis (jenis)}
				<div class="rounded-lg border-ctrl border-attn-fill bg-attn-bg p-4">
					<dt class="flex items-center gap-2 font-bold text-attn">
						<CircleHelp size={20} strokeWidth={1.75} aria-hidden="true" />
						{t(`signal.${jenis}`)}
					</dt>
					<dd class="mt-2 text-small">{t(`about.signal.${jenis}`)}</dd>
				</div>
			{/each}
		</dl>
		<p class="font-serif italic text-muted">{t('signal.disclaimer')}</p>
	</section>

	<section id="media" aria-labelledby="judul-media" class={kelasBagian}>
		<h2 id="judul-media" class="t-h2">{t('about.media')}</h2>
		<p class="max-w-prose">{t('about.media.rule')}</p>
		<p class="max-w-prose rounded-md bg-info-bg p-4 text-small">{t('about.media.bias')}</p>
		{#await data.sumber then hasil}
			{#if hasil.ok}
				<ul class="flex flex-wrap gap-2">
					{#each hasil.data.media as media (media.domain)}
						<li class="rounded-pill bg-secondary px-3 py-1 text-small text-secondary-on">
							{media.nama} <span class="text-caption">({media.domain})</span>
						</li>
					{/each}
				</ul>
			{/if}
		{/await}
	</section>

	<div class="grid gap-14 md:grid-cols-2 md:gap-5">
		<section id="batasan" aria-labelledby="judul-batasan" class={kelasBagian}>
			<h2 id="judul-batasan" class="t-h2">{t('about.limits')}</h2>
			<ul class="list-disc space-y-2 pl-5 text-small">
				{#each keterbatasan as kunci (kunci)}<li>{t(kunci)}</li>{/each}
			</ul>
			<a href="/{bahasaBawaan}/data-belum-dibuka" class="font-bold text-secondary-on underline"
				>{t('unopened.link')}</a
			>
		</section>
		<section id="lisensi" aria-labelledby="judul-lisensi" class={kelasBagian}>
			<h2 id="judul-lisensi" class="t-h2">{t('about.license')}</h2>
			{#await data.sumber then hasil}
				{#if hasil.ok && hasil.data.atribusi_opentender}
					<p class="text-small">{hasil.data.atribusi_opentender}</p>
				{/if}
			{/await}
			<p class="text-small">{t('about.license.code')}</p>
			<p class="text-small">{t('map.basemapNote')}</p>
			<p class="text-small text-muted">{t('about.demo')}</p>
		</section>
	</div>
</div>
