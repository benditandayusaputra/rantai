<script lang="ts">
	import { Info } from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import GapClassChip from '$lib/components/GapClassChip.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import DaftarSekolah from '$lib/features/pendidikan/DaftarSekolah.svelte';
	import GrafikKursi from '$lib/features/pendidikan/GrafikKursi.svelte';
	import { bahasaBawaan, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const halamanKecamatan = $derived(`/${bahasaBawaan}/kecamatan/${data.kode}`);
</script>

<svelte:head>
	{#await data.detail}
		<title>{t('brand.name')}</title>
	{:then hasil}
		<title
			>{hasil.ok
				? `${t('pen.title', { name: hasil.data.nama })} · ${t('brand.name')}`
				: t('brand.name')}</title
		>
	{/await}
</svelte:head>

<div class="wrap flex flex-col gap-8 pt-6 lg:pt-10">
	<a
		href={halamanKecamatan}
		class="inline-flex min-h-touch items-center self-start text-small font-semibold text-secondary-on underline"
		>{t('pen.backToArea')}</a
	>
	{#await data.detail}
		<div class="flex flex-col gap-4" aria-busy="true">
			<p class="sr-only">{t('state.loadingSummary')}</p>
			<div class="skel h-10 w-72 rounded-sm"></div>
			<div class="skel h-56 w-full rounded-lg"></div>
		</div>
	{:then hasil}
		{#if !hasil.ok}
			<h1 class="sr-only">{t('pen.notFound')}</h1>
			{#if hasil.galat.kode === 404}
				<EmptyState title={t('pen.notFound')} />
			{:else}
				<ErrorState title={t('error.sectionFailed')} onretry={invalidateAll} />
			{/if}
		{:else}
			{@const detail = hasil.data}
			<header class="flex flex-col gap-3">
				<p class="t-over">{t('bagian.tema')} · {t('sector.edu')}</p>
				<h1 class="t-h1">{t('pen.title', { name: detail.nama })}</h1>
				<div class="flex flex-wrap items-center gap-3">
					<GapClassChip kelas={detail.kelas} />
					<span class="text-small text-muted">{t('tema.period', { period: detail.periode })}</span>
				</div>
			</header>
			<div class="grid gap-8 lg:grid-cols-12 lg:gap-5">
				<div class="flex flex-col gap-6 lg:col-span-8">
					<GrafikKursi {detail} />
					{#if detail.metode_calon_lulusan}
						<p class="text-small text-muted">{t(`pen.method.${detail.metode_calon_lulusan}`)}</p>
					{/if}
				</div>
				<aside class="flex flex-col gap-4 lg:col-span-4">
					<KeyNumberCard label={t('pen.gap')} angka={detail.kesenjangan} size="lg" />
					{#if detail.kesenjangan_dengan_spmb_bersama}
						<KeyNumberCard
							label={t('pen.gapWithSpmb')}
							angka={detail.kesenjangan_dengan_spmb_bersama}
						/>
					{/if}
					{#if detail.rasio}
						<KeyNumberCard label={t('tema.ratio')} angka={detail.rasio} />
					{/if}
				</aside>
			</div>
			<section aria-labelledby="judul-keterbatasan" class="rounded-lg bg-info-bg p-4 md:p-5">
				<h2 id="judul-keterbatasan" class="flex items-center gap-2 font-bold">
					<Info size={20} strokeWidth={1.75} aria-hidden="true" class="text-info" />
					{t('pen.limits')}
				</h2>
				<ul class="mt-3 list-disc space-y-2 pl-5 text-small">
					{#each detail.keterbatasan as kunci (kunci)}
						<li>{t(`pen.limit.${kunci}`)}</li>
					{/each}
				</ul>
			</section>
			<a
				href="{halamanKecamatan}?bagian=anggaran&tema=pendidikan"
				class="inline-flex min-h-touch items-center self-start font-bold text-secondary-on underline"
				>{t('pen.toBudget')}</a
			>
		{/if}
	{/await}

	<section aria-labelledby="judul-sekolah" class="flex flex-col gap-4">
		<h2 id="judul-sekolah" class="t-h2">{t('pen.schoolsHeading')}</h2>
		{#await data.sekolah}
			<div class="skel h-40 w-full rounded-lg" aria-busy="true"></div>
		{:then hasil}
			{#if hasil.ok}
				<DaftarSekolah sekolah={hasil.data} />
			{:else}
				<ErrorState title={t('error.sectionFailed')} onretry={invalidateAll} />
			{/if}
		{/await}
	</section>
</div>
