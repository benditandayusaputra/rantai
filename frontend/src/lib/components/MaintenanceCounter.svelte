<script lang="ts">
	import { Hourglass } from '@lucide/svelte';
	import type { Pemeliharaan } from '$lib/api/tipe';
	import { kunciPemeliharaan } from '$lib/features/paket/paket';
	import { formatTanggal } from '$lib/format';
	import { bahasaBawaan, t } from '$lib/i18n';
	import CertaintyBadge from './CertaintyBadge.svelte';
	import ProvenanceBadge from './ProvenanceBadge.svelte';

	let { pemeliharaan }: { pemeliharaan: Pemeliharaan } = $props();

	const akhir = $derived(
		pemeliharaan.akhir_pemeliharaan ? new Date(pemeliharaan.akhir_pemeliharaan) : null
	);
</script>

<section
	aria-labelledby="judul-pemeliharaan"
	class="flex flex-col gap-3 rounded-lg border border-line bg-surface p-4 md:p-5"
>
	<h2 id="judul-pemeliharaan" class="flex items-center gap-2 font-bold">
		<Hourglass size={20} strokeWidth={1.75} aria-hidden="true" />
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
				{t('maint.until', { date: formatTanggal(akhir, bahasaBawaan, 'long') })}
			</p>
		{:else if akhir}
			<p class="text-small text-muted">
				{t('maint.ended', { date: formatTanggal(akhir, bahasaBawaan, 'long') })}
			</p>
		{/if}
	{:else}
		<p class="text-small text-muted">{t('maint.notConstruction')}</p>
	{/if}
</section>
