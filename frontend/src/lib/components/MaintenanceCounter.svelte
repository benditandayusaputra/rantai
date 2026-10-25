<script lang="ts">
	import { Hourglass } from '@lucide/svelte';
	import type { Pemeliharaan } from '$lib/api/tipe';
	import { kunciPemeliharaan } from '$lib/features/paket/paket';
	import { formatTanggal } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import CertaintyBadge from './CertaintyBadge.svelte';
	import ProvenanceBadge from './ProvenanceBadge.svelte';

	let { pemeliharaan }: { pemeliharaan: Pemeliharaan } = $props();

	const akhir = $derived(
		pemeliharaan.akhir_pemeliharaan ? new Date(pemeliharaan.akhir_pemeliharaan) : null
	);
</script>

<section aria-labelledby="judul-pemeliharaan" class="card flex flex-col gap-3 p-5">
	<h2 id="judul-pemeliharaan" class="flex items-center gap-3 font-bold">
		<span class="icon-tile h-9 w-9"
			><Hourglass size={18} strokeWidth={1.75} aria-hidden="true" /></span
		>
		{t('maint.title')}
	</h2>
	{#if pemeliharaan.berlaku}
		<div class="flex flex-wrap items-center gap-2">
			<span class="font-semibold">{t(kunciPemeliharaan[pemeliharaan.status])}</span>
			<CertaintyBadge level={pemeliharaan.label} />
		</div>
		{#if pemeliharaan.sisa_hari}
			<p class="n-lg">
				{t('maint.daysLeft', { n: pemeliharaan.sisa_hari.nilai })}
			</p>
			<ProvenanceBadge
				angka={pemeliharaan.sisa_hari}
				nilai={String(pemeliharaan.sisa_hari.nilai)}
				compact
			/>
		{/if}
		{#if akhir && pemeliharaan.status !== 'pemeliharaan_berakhir'}
			<p class="text-small text-muted">
				{t('maint.until', { date: formatTanggal(akhir, bahasaAktif, 'long') })}
			</p>
		{:else if akhir}
			<p class="text-small text-muted">
				{t('maint.ended', { date: formatTanggal(akhir, bahasaAktif, 'long') })}
			</p>
		{/if}
	{:else}
		<p class="text-small text-muted">{t('maint.notConstruction')}</p>
	{/if}
</section>
