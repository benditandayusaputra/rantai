<script lang="ts">
	import { Info, Network } from '@lucide/svelte';
	import { ambilJaringan, tangkap, type FilterJaringan } from '$lib/api/klien';
	import type { Jaringan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import { t } from '$lib/i18n';
	import GrafJaringan from './GrafJaringan.svelte';

	interface Props {
		filter: FilterJaringan;
		idBagian: string;
	}

	let { filter, idBagian }: Props = $props();

	let jaringan = $state<Jaringan | null>(null);
	let status = $state<'memuat' | 'siap' | 'gagal'>('memuat');
	let urutanPermintaan = 0;

	async function muat(f: FilterJaringan) {
		const nomor = ++urutanPermintaan;
		status = 'memuat';
		const hasil = await tangkap(ambilJaringan(fetch, f));
		if (nomor !== urutanPermintaan) return;
		jaringan = hasil.ok ? hasil.data : null;
		status = hasil.ok ? 'siap' : 'gagal';
	}

	$effect(() => {
		muat(filter);
	});
</script>

<section aria-labelledby="judul-{idBagian}" class="flex flex-col gap-4">
	<div class="flex items-center gap-3">
		<span class="icon-tile"><Network size={20} strokeWidth={1.75} aria-hidden="true" /></span>
		<h2 id="judul-{idBagian}" class="t-h2">{t('network.title')}</h2>
	</div>
	<p class="max-w-prose text-muted">{t('network.lead')}</p>
	<p class="flex max-w-prose gap-2 rounded-md bg-info-bg p-3 text-small text-info">
		<Info size={18} strokeWidth={1.75} class="mt-0.5 shrink-0" aria-hidden="true" />
		{t('network.note')}
	</p>
	{#if status === 'memuat'}
		<div class="skel h-96 rounded-lg" aria-busy="true"></div>
	{:else if status === 'gagal' || !jaringan}
		<ErrorState title={t('network.failed')} onretry={() => muat(filter)} />
	{:else if !jaringan.simpul.length}
		<EmptyState title={t('network.empty')} />
	{:else}
		<GrafJaringan {jaringan} idIsi="isi-{idBagian}" />
		{#if jaringan.is_terpotong}
			<p class="text-small text-muted">
				{t('network.truncated', { jumlah: jaringan.batas_paket })}
			</p>
		{/if}
	{/if}
</section>
