<script lang="ts">
	import {
		ArrowLeft,
		ArrowRight,
		Calendar,
		GraduationCap,
		Info,
		School,
		Wallet
	} from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import ChoiceChips from '$lib/components/ChoiceChips.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import GapClassChip from '$lib/components/GapClassChip.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import DaftarSekolah from '$lib/features/pendidikan/DaftarSekolah.svelte';
	import GrafikKursi from '$lib/features/pendidikan/GrafikKursi.svelte';
	import GrafikTren from '$lib/features/pendidikan/GrafikTren.svelte';
	import { daftarTemaPendidikan, kunciMetode, teksJenjang } from '$lib/features/pendidikan/jenjang';
	import RumusKesenjangan from '$lib/features/pendidikan/RumusKesenjangan.svelte';
	import { bahasaAktif, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const halamanKecamatan = $derived(`/${bahasaAktif}/kecamatan/${data.kode}`);
	const teks = $derived(teksJenjang[data.tema]);
	const pilihan = $derived(
		daftarTemaPendidikan.map((tema) => ({
			key: tema,
			href: `${halamanKecamatan}/${tema}`,
			label: t(teksJenjang[tema].jenjang),
			current: tema === data.tema
		}))
	);
</script>

<svelte:head>
	{#await data.detail}
		<title>{t('brand.name')}</title>
	{:then hasil}
		<title
			>{hasil.ok
				? `${t(teks.judul, { name: hasil.data.nama })} · ${t('brand.name')}`
				: t('brand.name')}</title
		>
	{/await}
</svelte:head>

<div class="border-b border-line bg-surface">
	<div class="wrap flex flex-col gap-4 pb-8 pt-6 lg:pb-10 lg:pt-10">
		<a href={halamanKecamatan} class="btn btn-ghost btn-sm -ml-3 self-start">
			<ArrowLeft size={18} strokeWidth={1.75} aria-hidden="true" />
			{t('pen.backToArea')}
		</a>
		<div class="flex flex-wrap items-center gap-3">
			<p id="judul-jenjang" class="text-small font-semibold text-muted">{t('level.label')}</p>
			<ChoiceChips labelledby="judul-jenjang" items={pilihan} />
		</div>
		{#await data.detail}
			<div class="flex items-center gap-4" aria-busy="true">
				<p class="sr-only">{t('state.loadingSummary')}</p>
				<div class="skel h-14 w-14 rounded-md"></div>
				<div class="skel h-10 w-72 max-w-full rounded-sm"></div>
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
				<header class="flex items-start gap-4">
					<span class="icon-tile hidden h-14 w-14 md:inline-flex">
						<GraduationCap size={28} strokeWidth={1.75} aria-hidden="true" />
					</span>
					<div class="flex min-w-0 flex-col gap-3">
						<p class="t-over">{t('bagian.tema')} · {t('sector.edu')} · {t(teks.jenjang)}</p>
						<h1 class="t-h1">{t(teks.judul, { name: detail.nama })}</h1>
						<div class="flex flex-wrap items-center gap-3">
							<GapClassChip kelas={detail.kelas} />
							<span class="inline-flex items-center gap-1 text-small text-muted">
								<Calendar size={16} strokeWidth={1.75} aria-hidden="true" />
								{t('tema.period', { period: detail.periode })}
							</span>
						</div>
					</div>
				</header>
			{/if}
		{/await}
	</div>
</div>

<div class="wrap flex flex-col gap-10 pt-8 lg:pt-10">
	{#await data.detail}
		<div class="skel h-56 w-full rounded-lg" aria-hidden="true"></div>
	{:then hasil}
		{#if hasil.ok}
			{@const detail = hasil.data}
			{@const metode = kunciMetode(data.tema, detail.metode_calon_lulusan)}
			<div class="grid gap-6 lg:grid-cols-12 lg:gap-5">
				<div class="flex flex-col gap-4 lg:col-span-8">
					<GrafikKursi {detail} tema={data.tema} />
					{#if metode}
						<p class="flex gap-2 text-small text-muted">
							<Info size={18} strokeWidth={1.75} aria-hidden="true" class="mt-px shrink-0" />
							{t(metode)}
						</p>
					{/if}
					{#if detail.tren.length}
						<GrafikTren tren={detail.tren} tema={data.tema} />
					{/if}
				</div>
				<aside class="grid grid-cols-2 gap-3 lg:col-span-4 lg:flex lg:flex-col lg:gap-4">
					<div class="col-span-2">
						<KeyNumberCard label={t('pen.gap')} angka={detail.kesenjangan} size="lg" />
					</div>
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
			<RumusKesenjangan {detail} />
			<section aria-labelledby="judul-keterbatasan" class="rounded-lg bg-info-bg p-4 md:p-6">
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
				href="{halamanKecamatan}?bagian=anggaran&tema={data.tema}"
				class="btn btn-outline self-start max-md:btn-block"
			>
				<Wallet size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('pen.toBudget')}
				<ArrowRight size={18} strokeWidth={1.75} aria-hidden="true" />
			</a>
		{/if}
	{/await}

	<section aria-labelledby="judul-sekolah" class="flex flex-col gap-4">
		<h2 id="judul-sekolah" class="t-h2 flex items-center gap-3">
			<span class="icon-tile"><School size={20} strokeWidth={1.75} aria-hidden="true" /></span>
			{t(teks.judulSekolah)}
		</h2>
		{#await data.sekolah}
			<div class="skel h-40 w-full rounded-lg" aria-busy="true"></div>
		{:then hasil}
			{#if hasil.ok}
				<DaftarSekolah sekolah={hasil.data} tema={data.tema} kodeKecamatan={data.kode} />
			{:else}
				<ErrorState title={t('error.sectionFailed')} onretry={invalidateAll} />
			{/if}
		{/await}
	</section>
</div>
