<script lang="ts">
	import { FlaskConical } from '@lucide/svelte';
	import { tangkap } from '$lib/api/klien';
	import { ambilPaketTerdekat } from '$lib/api/laporan';
	import type { PaketTerdekat } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import RadioCard from '$lib/components/RadioCard.svelte';
	import { formatJarak } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import { draf } from './draf.svelte';
	import { paketTidakTerdaftar, radiusProyekM } from './lapor';

	let pesanGalat = $state('');
	const daftar = $derived<PaketTerdekat[] | null>(draf.paketTerdekat);

	$effect(() => {
		const titik = draf.titik;
		if (!titik || draf.paketTerdekat) return;
		void tangkap(ambilPaketTerdekat(fetch, titik, radiusProyekM)).then((hasil) => {
			if (hasil.ok) draf.paketTerdekat = hasil.data;
			else pesanGalat = hasil.galat.message;
		});
	});
</script>

<fieldset class="flex flex-col gap-3">
	<legend class="mb-2 text-balance font-heading text-h1">{t('report.s2.title')}</legend>
	<p class="text-small text-muted">{t('report.s1.accuracy', { m: Math.round(draf.akurasiM) })}</p>
	{#if pesanGalat}
		<p role="alert" class="text-danger">{pesanGalat}</p>
	{:else if daftar === null}
		<p aria-live="polite" class="text-muted">{t('report.s2.searching')}</p>
	{:else if daftar.length === 0}
		<EmptyState title={t('report.s2.none')} />
	{/if}
	{#each daftar ?? [] as paket (paket.kode)}
		<RadioCard
			name="paket"
			value={paket.kode}
			bind:group={draf.paketKode}
			title={paket.nama}
			description={t('report.s2.distance', { jarak: formatJarak(paket.jarak.nilai, bahasaAktif) })}
		>
			{#snippet leading()}
				{#if paket.is_simulation}
					<span
						class="inline-flex shrink-0 items-center gap-1 rounded-xs bg-cert-sim-bg px-1 text-caption font-bold text-cert-sim"
					>
						<FlaskConical size={14} strokeWidth={1.75} aria-hidden="true" />{t('sim.ribbon')}
					</span>
				{/if}
			{/snippet}
		</RadioCard>
	{/each}
	<RadioCard
		name="paket"
		value={paketTidakTerdaftar}
		bind:group={draf.paketKode}
		title={t('report.s2.notListed')}
		description={t('report.s2.notListedHint')}
	/>
</fieldset>
