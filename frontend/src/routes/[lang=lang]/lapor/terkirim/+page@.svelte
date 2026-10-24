<script lang="ts">
	import { Link2, Share2 } from '@lucide/svelte';
	import ActionBar from '$lib/components/ActionBar.svelte';
	import ReportTimeline from '$lib/components/ReportTimeline.svelte';
	import SimRibbon from '$lib/components/SimRibbon.svelte';
	import { bagikanTautan } from '$lib/features/lapor/bagikan';
	import { kelasTombolTeks, kelasTombolUtama } from '$lib/features/lapor/gaya';
	import { jumlahSaksiDibutuhkan, susunLinimasa } from '$lib/features/lapor/lapor';
	import { bahasaBawaan, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let isDisalin = $state(false);
	const tautanLaporan = $derived(`/${bahasaBawaan}/laporan/${data.id}`);

	async function ajakSaksi() {
		const hasil = await bagikanTautan(
			t('report.done.shareText'),
			new URL(tautanLaporan, location.origin).toString()
		);
		isDisalin = hasil === 'disalin';
	}
</script>

<svelte:head>
	<title>{t('report.done.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="flex min-h-dvh flex-col bg-page">
	<main id="konten" class="mx-auto flex w-full max-w-form flex-1 flex-col gap-5 px-4 py-8">
		<span
			class="inline-flex h-16 w-16 items-center justify-center rounded-pill bg-secondary text-primary"
		>
			<Link2 size={32} strokeWidth={1.75} aria-hidden="true" />
		</span>
		<h1 class="t-h1">{t('report.done.title')}</h1>
		<p class="t-lead">{t('report.done.body')}</p>
		{#await data.laporan then hasil}
			{#if hasil.ok}
				{#if hasil.data.is_simulation}<SimRibbon />{/if}
				<section
					aria-label={t('report.timeline')}
					class="rounded-lg border-ctrl border-line bg-surface p-4"
				>
					<ReportTimeline
						events={susunLinimasa(hasil.data.riwayat).slice(0, jumlahSaksiDibutuhkan)}
						witnesses={hasil.data.saksi}
					/>
				</section>
			{/if}
		{/await}
		{#if isDisalin}<p role="status" class="text-success">{t('report.done.copied')}</p>{/if}
	</main>
	<ActionBar>
		<button type="button" class={kelasTombolUtama} onclick={ajakSaksi}>
			<Share2 size={20} strokeWidth={1.75} aria-hidden="true" />
			{t('report.done.invite')}
		</button>
		<a href={tautanLaporan} class={kelasTombolTeks}>{t('report.done.view')}</a>
	</ActionBar>
</div>
