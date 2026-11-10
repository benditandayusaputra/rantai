<script lang="ts">
	import { Check, LoaderCircle, LocateFixed, LogOut, RefreshCw, ShieldCheck } from '@lucide/svelte';
	import ActionBar from '$lib/components/ActionBar.svelte';
	import { induk } from '$lib/features/wilayah/wilayah';
	import { t } from '$lib/i18n';
	import { draf } from './draf.svelte';
	import { cariPosisi, kenaliTempat, type HasilPosisi, type HasilTempat } from './lapor';

	let { onlanjut, onkeluar }: { onlanjut: () => void; onkeluar: () => void } = $props();

	let hasil = $state<HasilPosisi | null>(null);
	let tempat = $state<HasilTempat | null>(null);
	let isMencari = $state(false);

	const isPerluPerhatian = $derived(hasil?.status === 'kasar' || tempat?.status === 'luar');
	const isGagal = $derived(hasil?.status === 'gagal' || tempat?.status === 'gagal');

	async function cari() {
		isMencari = true;
		tempat = null;
		hasil = await cariPosisi(navigator.geolocation);
		if (hasil.status === 'akurat') tempat = await kenaliTempat(fetch, hasil.titik);
		isMencari = false;
	}

	function konfirmasi() {
		if (hasil?.status !== 'akurat') return;
		draf.titik = hasil.titik;
		draf.akurasiM = hasil.akurasiM;
		draf.paketKode = null;
		draf.paketTerdekat = null;
		onlanjut();
	}
</script>

<div class="flex flex-col gap-5">
	<div
		class="flex aspect-[2/1] items-center justify-center rounded-lg bg-map {isPerluPerhatian
			? 'text-attn'
			: 'text-primary'}"
		aria-hidden="true"
	>
		<span
			class="flex h-24 w-24 items-center justify-center rounded-pill border-link border-dashed {isPerluPerhatian
				? 'border-attn-fill bg-attn-bg'
				: 'border-primary bg-secondary'}"
		>
			<LocateFixed size={28} strokeWidth={1.75} />
		</span>
	</div>

	{#if hasil?.status === 'kasar'}
		<div
			role="alert"
			class="flex flex-col gap-2 rounded-lg border-ctrl border-attn-fill bg-attn-bg p-4"
		>
			<h1 class="t-h3 text-attn">{t('report.s1.poorTitle')}</h1>
			<p>{t('report.s1.poorBody', { m: Math.round(hasil.akurasiM) })}</p>
		</div>
	{:else if tempat?.status === 'luar'}
		<div
			role="alert"
			class="flex flex-col gap-2 rounded-lg border-ctrl border-attn-fill bg-attn-bg p-4"
		>
			<h1 class="t-h3 text-attn">{t('report.s1.outsideTitle')}</h1>
			<p>{t('report.s1.outsideBody')}</p>
		</div>
	{:else if hasil?.status === 'akurat' && tempat?.status === 'dikenali'}
		<div role="status" class="flex flex-col gap-2 rounded-lg border-ctrl border-primary p-4">
			<h1 class="t-h3">{t('report.s1.confirmTitle')}</h1>
			<p class="font-heading text-h3 font-semibold">{tempat.wilayah.nama}</p>
			<p class="text-small text-muted">
				{t('report.s1.confirmBody', {
					area: induk(tempat.wilayah),
					m: Math.round(hasil.akurasiM)
				})}
			</p>
		</div>
	{:else}
		<h1 class="text-balance font-heading text-h1">{t('report.s1.title')}</h1>
	{/if}

	<ul class="flex flex-col gap-3 rounded-lg bg-secondary p-4 text-secondary-on">
		{#each ['report.s1.radius', 'report.s1.rounded'] as const as kunci (kunci)}
			<li class="flex gap-3">
				<ShieldCheck size={20} strokeWidth={1.75} class="mt-0.5 shrink-0" aria-hidden="true" />
				<span>{t(kunci)}</span>
			</li>
		{/each}
	</ul>

	{#if hasil?.status === 'ditolak' || isGagal}
		<p role="alert" class="text-danger">
			{t(hasil?.status === 'ditolak' ? 'report.s1.denied' : 'report.s1.failed')}
		</p>
	{/if}
</div>

<div class="-mx-4 -mb-6 mt-auto pt-5 md:m-0">
	<ActionBar>
		{#if tempat?.status === 'dikenali' && !isMencari}
			<button type="button" class="btn btn-primary btn-lg btn-block" onclick={konfirmasi}>
				<Check size={20} strokeWidth={1.75} aria-hidden="true" />
				{t('report.s1.confirmYes')}
			</button>
			<button type="button" class="btn btn-text btn-block" onclick={cari}>
				<RefreshCw size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('report.s1.searchAgain')}
			</button>
		{:else}
			<button
				type="button"
				class="btn btn-primary btn-lg btn-block"
				disabled={isMencari}
				aria-busy={isMencari}
				onclick={cari}
			>
				{#if isMencari}
					<LoaderCircle size={20} strokeWidth={1.75} class="animate-spin" aria-hidden="true" />
					{t('report.s1.locating')}
				{:else if hasil}
					<RefreshCw size={20} strokeWidth={1.75} aria-hidden="true" />
					{t('report.s1.searchAgain')}
				{:else}
					<LocateFixed size={20} strokeWidth={1.75} aria-hidden="true" />
					{t('report.s1.allow')}
				{/if}
			</button>
			{#if tempat?.status === 'luar'}
				<button type="button" class="btn btn-text btn-block" onclick={onkeluar}>
					<LogOut size={18} strokeWidth={1.75} aria-hidden="true" />
					{t('report.exit')}
				</button>
			{/if}
		{/if}
	</ActionBar>
</div>
