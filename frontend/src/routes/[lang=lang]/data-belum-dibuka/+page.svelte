<script lang="ts">
	import { ExternalLink, FileText, Lock, Replace } from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import ChoiceChips from '$lib/components/ChoiceChips.svelte';
	import DocumentPreview from '$lib/components/DocumentPreview.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import { bahasaDokumen } from '$lib/features/aspirasi/aspirasi';
	import { tautanPermohonan } from '$lib/features/keterbukaan/keterbukaan';
	import { bahasaBawaan, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>{t('unopened.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="wrap pt-6 lg:pt-10">
	<p class="t-over">{t('unopened.over')}</p>
	<h1 class="t-h1 mt-1">{t('unopened.title')}</h1>
	<p class="t-lead mt-3 max-w-prose">{t('unopened.lead')}</p>

	{#await data.halaman}
		<div class="mt-8 grid gap-4 lg:grid-cols-2" aria-busy="true">
			<p class="sr-only">{t('unopened.loading')}</p>
			{#each [0, 1, 2, 3] as i (i)}<div class="skel h-48 rounded-lg"></div>{/each}
		</div>
	{:then hasil}
		{#if !hasil.ok}
			<div class="mt-8">
				<ErrorState
					title={t('unopened.failed')}
					description={hasil.galat.message}
					onretry={invalidateAll}
				/>
			</div>
		{:else}
			{@const { daftar, dokumen } = hasil.data}
			{@const terpilih = daftar.data.find((d) => d.kode === data.pilihan.kode)}
			<ul class="mt-8 grid gap-4 lg:grid-cols-2">
				{#each daftar.data as item (item.kode)}
					<li
						class="flex flex-col gap-4 rounded-lg border p-5 {item.kode === data.pilihan.kode
							? 'border-primary bg-raised'
							: 'border-line bg-surface'}"
					>
						<h2 class="t-h3">{item.nama[bahasaBawaan]}</h2>
						<div class="flex flex-col gap-1">
							<p class="t-over flex items-center gap-2">
								<Lock size={16} strokeWidth={1.75} aria-hidden="true" />{t('unopened.why')}
							</p>
							<p class="text-small">{item.alasan[bahasaBawaan]}</p>
						</div>
						<div class="flex flex-col gap-1">
							<p class="t-over flex items-center gap-2">
								<Replace size={16} strokeWidth={1.75} aria-hidden="true" />{t(
									'unopened.substitute'
								)}
							</p>
							<p class="text-small font-semibold">{item.pengganti[bahasaBawaan]}</p>
						</div>
						<a
							href={tautanPermohonan(page.url, { kode: item.kode })}
							data-sveltekit-noscroll
							aria-current={item.kode === data.pilihan.kode}
							class="mt-auto inline-flex min-h-touch items-center justify-center gap-2 self-start rounded-md border-ctrl border-line-strong px-4 font-semibold hover:bg-sunken aria-[current=true]:border-primary aria-[current=true]:bg-secondary aria-[current=true]:text-secondary-on max-md:w-full"
						>
							<FileText size={20} strokeWidth={1.75} aria-hidden="true" />{t('unopened.request')}
						</a>
					</li>
				{/each}
			</ul>

			{#if terpilih && dokumen}
				<section
					id="surat"
					aria-labelledby="judul-surat"
					class="mt-10 grid scroll-mt-24 gap-8 lg:grid-cols-12"
				>
					<div class="flex flex-col gap-6 lg:col-span-4">
						<div class="flex flex-col gap-2 rounded-lg border border-line bg-surface p-5">
							<p class="t-over">{t('unopened.to')}</p>
							<p class="font-semibold">{daftar.ppid.nama}</p>
							<p class="text-small text-muted">{daftar.ppid.alamat}</p>
							<a
								href={daftar.ppid.url}
								target="_blank"
								rel="noopener"
								class="inline-flex min-h-touch items-center gap-2 font-bold text-secondary-on underline"
							>
								{t('unopened.ppidSite')}
								<ExternalLink size={16} strokeWidth={1.75} aria-hidden="true" />
								<span class="sr-only">{t('a11y.opensNewTab')}</span>
							</a>
						</div>
						<div class="flex flex-col gap-3">
							<h3 id="judul-bahasa-surat" class="t-h3">{t('packet.docLang')}</h3>
							<ChoiceChips
								labelledby="judul-bahasa-surat"
								items={bahasaDokumen.map((bahasa) => ({
									key: bahasa,
									href: tautanPermohonan(page.url, { bahasa }),
									label: t(`packet.lang.${bahasa}`),
									current: data.pilihan.bahasa === bahasa
								}))}
							/>
						</div>
						<p class="text-small text-muted">{t('unopened.selfDeliver')}</p>
					</div>
					<div class="lg:col-span-8">
						<DocumentPreview
							id="judul-surat"
							title={t('unopened.letterTitle', { name: terpilih.nama[bahasaBawaan] })}
							{dokumen}
						/>
					</div>
				</section>
			{/if}
		{/if}
	{/await}
</div>
