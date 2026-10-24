<script lang="ts">
	import { LoaderCircle, LocateFixed, ShieldCheck } from '@lucide/svelte';
	import ActionBar from '$lib/components/ActionBar.svelte';
	import { t } from '$lib/i18n';
	import { draf } from './draf.svelte';
	import { kelasTombolUtama } from './gaya';
	import { cariPosisi, type HasilPosisi } from './lapor';

	let { onlanjut }: { onlanjut: () => void } = $props();

	let hasil = $state<HasilPosisi | null>(null);
	let isMencari = $state(false);

	async function izinkan() {
		isMencari = true;
		hasil = await cariPosisi(navigator.geolocation);
		isMencari = false;
		if (hasil.status !== 'akurat') return;
		draf.titik = hasil.titik;
		draf.akurasiM = hasil.akurasiM;
		draf.paketKode = null;
		draf.paketTerdekat = null;
		onlanjut();
	}
</script>

<div class="flex flex-col gap-5">
	<div
		class="flex aspect-[2/1] items-center justify-center rounded-lg bg-map {hasil?.status ===
		'kasar'
			? 'text-attn'
			: 'text-primary'}"
		aria-hidden="true"
	>
		<span
			class="flex h-24 w-24 items-center justify-center rounded-pill border-link border-dashed {hasil?.status ===
			'kasar'
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
	{:else}
		<h1 class="t-h1">{t('report.s1.title')}</h1>
	{/if}

	<ul class="flex flex-col gap-3">
		{#each ['report.s1.radius', 'report.s1.rounded'] as const as kunci (kunci)}
			<li class="flex gap-3">
				<ShieldCheck
					size={20}
					strokeWidth={1.75}
					class="mt-0.5 shrink-0 text-primary"
					aria-hidden="true"
				/>
				<span>{t(kunci)}</span>
			</li>
		{/each}
	</ul>

	{#if hasil?.status === 'ditolak' || hasil?.status === 'gagal'}
		<p role="alert" class="text-danger">
			{t(hasil.status === 'ditolak' ? 'report.s1.denied' : 'report.s1.failed')}
		</p>
	{/if}
</div>

<ActionBar>
	<button
		type="button"
		class={kelasTombolUtama}
		disabled={isMencari}
		aria-busy={isMencari}
		onclick={izinkan}
	>
		{#if isMencari}
			<LoaderCircle size={20} strokeWidth={1.75} class="animate-spin" aria-hidden="true" />
			{t('report.s1.locating')}
		{:else}
			{hasil ? t('action.retry') : t('report.s1.allow')}
		{/if}
	</button>
</ActionBar>
