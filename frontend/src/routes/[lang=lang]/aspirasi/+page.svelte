<script lang="ts">
	import { Check, ChevronRight, Copy, Download } from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import { urlServer } from '$lib/api/klien';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import { bahasaDokumen, tautanAspirasi } from '$lib/features/aspirasi/aspirasi';
	import { bahasaBawaan, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const durasiTandaSalin = 3000;
	let isDisalin = $state(false);

	async function salin(teks: string) {
		await navigator.clipboard.writeText(teks);
		isDisalin = true;
		setTimeout(() => (isDisalin = false), durasiTandaSalin);
	}

	const kelasChip =
		'inline-flex min-h-touch items-center rounded-pill border-ctrl border-line-strong px-4 text-small hover:bg-sunken aria-[current=true]:border-primary aria-[current=true]:bg-secondary aria-[current=true]:font-bold aria-[current=true]:text-secondary-on';
	const kelasTombol =
		'inline-flex min-h-touch items-center justify-center gap-2 rounded-md px-4 font-semibold';
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
					<nav aria-labelledby="judul-masalah">
						<h2 id="judul-masalah" class="t-h3">{t('acc.issueTitle')}</h2>
						<ul class="mt-3 flex flex-wrap gap-2">
							{#each tema.tanggung_jawab as tj (tj.masalah)}
								<li>
									<a
										href={tautanAspirasi(page.url, { masalah: tj.masalah, tema: tema.kode })}
										aria-current={data.pilihan.masalah === tj.masalah}
										data-sveltekit-noscroll
										class={kelasChip}>{tj.label.id}</a
									>
								</li>
							{/each}
						</ul>
					</nav>
					<nav aria-labelledby="judul-bahasa">
						<h2 id="judul-bahasa" class="t-h3">{t('packet.docLang')}</h2>
						<ul class="mt-3 flex gap-2">
							{#each bahasaDokumen as bahasa (bahasa)}
								<li>
									<a
										href={tautanAspirasi(page.url, { bahasa })}
										aria-current={data.pilihan.bahasa === bahasa}
										data-sveltekit-noscroll
										class={kelasChip}>{t(`packet.lang.${bahasa}`)}</a
									>
								</li>
							{/each}
						</ul>
					</nav>
					<p class="text-small text-muted">{t('packet.selfDeliver')}</p>
				</div>

				<section aria-labelledby="judul-pratinjau" class="lg:col-span-8">
					{#if dokumen}
						<div class="flex flex-wrap items-center justify-between gap-3">
							<h2 id="judul-pratinjau" class="t-h3">{t('packet.previewTitle')}</h2>
							<div class="flex flex-wrap gap-2">
								<button
									type="button"
									class="{kelasTombol} border-ctrl border-line-strong hover:bg-sunken"
									onclick={() => salin(dokumen.teks)}
								>
									{#if isDisalin}
										<Check size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.copied')}
									{:else}
										<Copy size={20} strokeWidth={1.75} aria-hidden="true" />{t('action.copy')}
									{/if}
								</button>
								<a
									href={urlServer(dokumen.pdf_url)}
									download
									class="{kelasTombol} bg-primary text-primary-on"
								>
									<Download size={20} strokeWidth={1.75} aria-hidden="true" />{t(
										'action.downloadPdf'
									)}
								</a>
							</div>
						</div>
						<p class="sr-only" aria-live="polite">{isDisalin ? t('action.copied') : ''}</p>
						<pre
							class="mt-4 overflow-x-auto whitespace-pre-wrap rounded-lg border border-line bg-raised p-5 font-sans text-body shadow-1">{dokumen.teks}</pre>
					{:else}
						<h2 id="judul-pratinjau" class="sr-only">{t('packet.previewTitle')}</h2>
						<EmptyState title={t('packet.chooseIssue')} />
					{/if}
				</section>
			</div>
		{/if}
	{/await}
</div>
