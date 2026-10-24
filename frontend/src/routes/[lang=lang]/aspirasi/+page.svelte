<script lang="ts">
	import { ChevronRight } from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import ChoiceChips from '$lib/components/ChoiceChips.svelte';
	import DocumentPreview from '$lib/components/DocumentPreview.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import { bahasaDokumen, tautanAspirasi } from '$lib/features/aspirasi/aspirasi';
	import { bahasaBawaan, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{t('packet.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="wrap pt-6 lg:pt-10">
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
						<a href="/{bahasaBawaan}/kecamatan/{wilayah.kode}?bagian=wakilku" class="underline"
							>{wilayah.nama}</a
						>
					</li>
					<li aria-hidden="true"><ChevronRight size={16} strokeWidth={1.75} /></li>
					<li aria-current="page">{t('packet.title')}</li>
				</ol>
			</nav>
			<p class="t-over mt-4">{tema.nama.id}</p>
			<h1 class="t-h1 mt-1">{t('packet.title')}</h1>
			<p class="mt-2 text-muted">{t('packet.forArea', { name: wilayah.nama })}</p>

			<div class="mt-8 grid gap-8 lg:grid-cols-12">
				<div class="flex flex-col gap-6 lg:col-span-4">
					<div class="flex flex-col gap-3">
						<h2 id="judul-masalah" class="t-h3">{t('acc.issueTitle')}</h2>
						<ChoiceChips
							labelledby="judul-masalah"
							items={tema.tanggung_jawab.map((tj) => ({
								key: tj.masalah,
								href: tautanAspirasi(page.url, { masalah: tj.masalah, tema: tema.kode }),
								label: tj.label.id,
								current: data.pilihan.masalah === tj.masalah
							}))}
						/>
					</div>
					<div class="flex flex-col gap-3">
						<h2 id="judul-bahasa" class="t-h3">{t('packet.docLang')}</h2>
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
					<p class="text-small text-muted">{t('packet.selfDeliver')}</p>
				</div>

				<section aria-labelledby="judul-pratinjau" class="lg:col-span-8">
					{#if dokumen}
						<DocumentPreview id="judul-pratinjau" title={t('packet.previewTitle')} {dokumen} />
					{:else}
						<h2 id="judul-pratinjau" class="sr-only">{t('packet.previewTitle')}</h2>
						<EmptyState title={t('packet.chooseIssue')} />
					{/if}
				</section>
			</div>
		{/if}
	{/await}
</div>
