<script lang="ts">
	import TableScroll from '$lib/components/TableScroll.svelte';
	import {
		ArrowRight,
		BadgeCheck,
		CircleHelp,
		Database,
		FileLock,
		FileText,
		MapPin,
		Newspaper,
		Sigma,
		TriangleAlert
	} from '@lucide/svelte';
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
	const daftarIsi = [
		{ id: 'sumber', judul: 'about.sources', ikon: Database },
		{ id: 'label', judul: 'about.labels', ikon: BadgeCheck },
		{ id: 'rumus', judul: 'about.formula', ikon: Sigma },
		{ id: 'lokasi', judul: 'about.location', ikon: MapPin },
		{ id: 'sinyal', judul: 'about.signals', ikon: CircleHelp },
		{ id: 'media', judul: 'about.media', ikon: Newspaper },
		{ id: 'batasan', judul: 'about.limits', ikon: TriangleAlert },
		{ id: 'lisensi', judul: 'about.license', ikon: FileText }
	] satisfies { id: string; judul: KunciPesan; ikon: typeof Database }[];
	const ikonBagian = Object.fromEntries(daftarIsi.map((item) => [item.id, item.ikon]));
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

	const kelasBagian = 'flex min-w-0 scroll-mt-24 flex-col gap-4';
	const kelasKartu = 'card flex min-w-0 scroll-mt-24 flex-col gap-4 p-5 md:p-6';
</script>

<svelte:head>
	<title>{t('about.title')} · {t('brand.name')}</title>
</svelte:head>

{#snippet judulBagian(id: string, kunci: KunciPesan, kelasJudul = 't-h2')}
	{@const Ikon = ikonBagian[id]}
	<div class="flex items-center gap-3">
		<span class="icon-tile"><Ikon size={20} strokeWidth={1.75} aria-hidden="true" /></span>
		<h2 id="judul-{id}" class={kelasJudul}>{t(kunci)}</h2>
	</div>
{/snippet}

<section class="border-b border-line bg-surface">
	<header class="wrap flex flex-col gap-4 pb-10 pt-10 lg:pb-14 lg:pt-16">
		<p class="t-over">{t('about.overline')}</p>
		<h1 class="t-display max-w-3xl">{t('about.title')}</h1>
		<p class="t-lead max-w-prose">{t('about.lead')}</p>
	</header>
</section>

<div class="wrap mt-8 grid gap-8 lg:mt-12 lg:grid-cols-[theme(spacing.66)_1fr] lg:gap-12">
	<nav aria-label={t('about.toc')} class="min-w-0 lg:sticky lg:top-24 lg:self-start">
		<p class="t-over mb-3 hidden lg:block">{t('about.toc')}</p>
		<ul class="-mx-4 flex gap-2 overflow-x-auto px-4 pb-1 lg:mx-0 lg:flex-col lg:gap-1 lg:px-0">
			{#each daftarIsi as item (item.id)}
				{@const Ikon = item.ikon}
				<li class="shrink-0">
					<a
						href="#{item.id}"
						class="chip lg:w-full lg:justify-start lg:rounded-md lg:border-transparent lg:bg-transparent lg:px-3"
					>
						<Ikon size={16} strokeWidth={1.75} aria-hidden="true" />{t(item.judul)}
					</a>
				</li>
			{/each}
		</ul>
	</nav>

	<div class="flex min-w-0 flex-col gap-12">
		<section id="sumber" aria-labelledby="judul-sumber" class={kelasBagian}>
			{@render judulBagian('sumber', 'about.sources')}
			{#await data.sumber}
				<div class="skel h-48 rounded-lg" aria-busy="true"></div>
			{:then hasil}
				{#if !hasil.ok}
					<ErrorState title={t('about.sourcesFailed')} onretry={invalidateAll} />
				{:else if !hasil.data.pengambilan.length}
					<EmptyState title={t('about.sourcesEmpty')} />
				{:else}
					<div class="card overflow-hidden">
						<TableScroll label={t('about.sources')}>
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
						</TableScroll>
					</div>
				{/if}
			{/await}
		</section>

		<section id="label" aria-labelledby="judul-label" class={kelasBagian}>
			{@render judulBagian('label', 'about.labels')}
			<ul class="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
				{#each label as level (level)}
					<li class="card flex flex-col gap-3 p-4 md:p-5">
						<CertaintyBadge {level} size="md" class="self-start" />
						<p class="text-small">{t(`cert.${level}.desc`)}</p>
					</li>
				{/each}
			</ul>
			<p class="max-w-prose text-small text-muted">{t('about.weakest')}</p>
		</section>

		<div class="grid gap-5">
			<section id="rumus" aria-labelledby="judul-rumus" class={kelasKartu}>
				{@render judulBagian('rumus', 'about.formula', 't-h3')}
				<pre
					class="whitespace-pre-wrap break-words rounded-lg bg-sunken p-4 font-mono text-mono">{t(
						'about.formula.code'
					)}</pre>
				<p class="text-small text-muted">{t('about.formulaNote')}</p>
			</section>
			<section id="lokasi" aria-labelledby="judul-lokasi" class={kelasKartu}>
				{@render judulBagian('lokasi', 'about.location', 't-h3')}
				<ol class="list-decimal space-y-2 pl-5 text-small">
					{#each aturanLokasi as kunci (kunci)}<li>{t(kunci)}</li>{/each}
				</ol>
			</section>
		</div>

		<section id="sinyal" aria-labelledby="judul-sinyal" class={kelasBagian}>
			{@render judulBagian('sinyal', 'about.signals')}
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
			{@render judulBagian('media', 'about.media')}
			<p class="max-w-prose">{t('about.media.rule')}</p>
			<p class="max-w-prose rounded-md bg-info-bg p-4 text-small">{t('about.media.bias')}</p>
			{#await data.sumber then hasil}
				{#if hasil.ok}
					<ul class="flex flex-wrap gap-2">
						{#each hasil.data.media as media (media.domain)}
							<li
								class="inline-flex items-center gap-2 rounded-pill bg-secondary px-3 py-1 text-small text-secondary-on"
							>
								<Newspaper size={14} strokeWidth={1.75} aria-hidden="true" />
								{media.nama} <span class="text-caption">({media.domain})</span>
							</li>
						{/each}
					</ul>
				{/if}
			{/await}
		</section>

		<div class="grid gap-5 xl:grid-cols-2">
			<section id="batasan" aria-labelledby="judul-batasan" class={kelasKartu}>
				{@render judulBagian('batasan', 'about.limits', 't-h3')}
				<ul class="list-disc space-y-2 pl-5 text-small">
					{#each keterbatasan as kunci (kunci)}<li>{t(kunci)}</li>{/each}
				</ul>
				<a href="/{bahasaBawaan}/data-belum-dibuka" class="btn btn-outline mt-auto self-start">
					<FileLock size={18} strokeWidth={1.75} aria-hidden="true" />
					{t('unopened.link')}
					<ArrowRight size={18} strokeWidth={1.75} aria-hidden="true" />
				</a>
			</section>
			<section id="lisensi" aria-labelledby="judul-lisensi" class={kelasKartu}>
				{@render judulBagian('lisensi', 'about.license', 't-h3')}
				{#await data.sumber then hasil}
					{#if hasil.ok && hasil.data.atribusi_opentender}
						<p class="text-small">{hasil.data.atribusi_opentender}</p>
					{/if}
				{/await}
				<p class="text-small">{t('about.license.code')}</p>
				<p class="text-small">{t('about.license.csv')}</p>
				<p class="text-small">{t('map.basemapNote')}</p>
				<p class="text-small text-muted">{t('about.demo')}</p>
			</section>
		</div>
	</div>
</div>
