<script lang="ts">
	import { ChevronRight, FileText, Globe, Info, ListChecks } from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import ChoiceChips from '$lib/components/ChoiceChips.svelte';
	import DocumentPreview from '$lib/components/DocumentPreview.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import { bahasaDokumen, tautanAspirasi } from '$lib/features/aspirasi/aspirasi';
	import LangkahPenyampaian from '$lib/features/aspirasi/LangkahPenyampaian.svelte';
	import { bahasaAktif, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{t('packet.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="wrap pt-8 lg:pt-12">
	{#await data.halaman}
		<div class="flex flex-col gap-4" aria-busy="true">
			<p class="text-muted">{t('packet.composingPlain')}</p>
			<div class="skel h-10 w-72 rounded-sm"></div>
			<div class="skel h-96 w-full rounded-lg"></div>
		</div>
	{:then hasil}
		{#if !hasil.ok}
			<h1 class="sr-only">{t('packet.failed')}</h1>
			<ErrorState
				title={t('packet.failed')}
				description={hasil.galat.message}
				onretry={invalidateAll}
			/>
		{:else}
			{@const { wilayah, tema, dokumen } = hasil.data}
			<nav aria-label={t('nav.breadcrumb')}>
				<ol class="flex flex-wrap items-center gap-1 text-small text-muted">
					<li>
						<a
							href="/{bahasaAktif}/kecamatan/{wilayah.kode}?bagian=wakilku"
							class="underline underline-offset-4 hover:text-fg">{wilayah.nama}</a
						>
					</li>
					<li aria-hidden="true"><ChevronRight size={16} strokeWidth={1.75} /></li>
					<li aria-current="page">{t('packet.title')}</li>
				</ol>
			</nav>
			<header class="mt-4 flex items-start gap-4">
				<span class="icon-tile h-12 w-12"
					><FileText size={24} strokeWidth={1.75} aria-hidden="true" /></span
				>
				<div>
					<p class="t-over">{tema.nama[bahasaAktif]}</p>
					<h1 class="t-h1 mt-1">{t('packet.title')}</h1>
					<p class="mt-2 text-muted">{t('packet.forArea', { name: wilayah.nama })}</p>
					<p class="mt-2 max-w-prose text-small text-muted">{t('packet.intro')}</p>
				</div>
			</header>

			<div class="mt-8 grid gap-8 lg:grid-cols-12">
				<div class="card flex flex-col gap-6 self-start p-5 lg:sticky lg:top-24 lg:col-span-4">
					<div class="flex flex-col gap-3">
						<h2 id="judul-masalah" class="t-h3 flex items-center gap-2">
							<ListChecks size={20} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
							{t('acc.issueTitle')}
						</h2>
						<ChoiceChips
							labelledby="judul-masalah"
							items={tema.tanggung_jawab.map((tj) => ({
								key: tj.masalah,
								href: tautanAspirasi(page.url, { masalah: tj.masalah, tema: tema.kode }),
								label: tj.label[bahasaAktif],
								current: data.pilihan.masalah === tj.masalah
							}))}
						/>
					</div>
					<div class="flex flex-col gap-3">
						<h2 id="judul-bahasa" class="t-h3 flex items-center gap-2">
							<Globe size={20} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
							{t('packet.docLang')}
						</h2>
						<ChoiceChips
							labelledby="judul-bahasa"
							items={bahasaDokumen.map((bahasa) => ({
								key: bahasa,
								href: tautanAspirasi(page.url, { bahasa }),
								label: t(`packet.lang.${bahasa}`),
								current: data.pilihan.bahasa === bahasa
							}))}
						/>
					</div>
					<p class="flex items-start gap-2 border-t border-line pt-4 text-small text-muted">
						<Info size={16} strokeWidth={1.75} aria-hidden="true" class="mt-0.5 shrink-0" />
						{t('packet.selfDeliver')}
					</p>
				</div>

				<section aria-labelledby="judul-pratinjau" class="lg:col-span-8">
					{#if dokumen}
						<DocumentPreview id="judul-pratinjau" title={t('packet.previewTitle')} {dokumen} />
						<LangkahPenyampaian
							kodeWilayah={wilayah.kode}
							penerima={tema.tanggung_jawab.find((tj) => tj.masalah === data.pilihan.masalah)
								?.instansi_pelaksana ?? ''}
						/>
					{:else}
						<h2 id="judul-pratinjau" class="sr-only">{t('packet.previewTitle')}</h2>
						<EmptyState title={t('packet.chooseIssue')} />
					{/if}
				</section>
			</div>
		{/if}
	{/await}
</div>
