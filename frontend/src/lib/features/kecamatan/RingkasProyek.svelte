<script lang="ts">
	import { LoaderCircle } from '@lucide/svelte';
	import { ambilPaketKecamatan, tangkap } from '$lib/api/klien';
	import type { Anggaran, PaketRingkas } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import PackageCard from '$lib/components/PackageCard.svelte';
	import { t } from '$lib/i18n';
	import { batasPaketRingkas } from './kecamatan';

	let { anggaran, kodeKecamatan }: { anggaran: Anggaran; kodeKecamatan: string } = $props();

	let daftar = $state<PaketRingkas[]>([]);
	let status = $state<'memuat' | 'siap' | 'gagal'>('memuat');

	$effect(() => {
		status = 'memuat';
		void tangkap(ambilPaketKecamatan(fetch, kodeKecamatan, { tahun: anggaran.tahun })).then(
			(hasil) => {
				daftar = hasil.ok ? hasil.data.data.slice(0, batasPaketRingkas) : [];
				status = hasil.ok ? 'siap' : 'gagal';
			}
		);
	});
</script>

<div class="grid grid-cols-2 gap-3 lg:gap-4">
	<KeyNumberCard label={t('money.count')} angka={anggaran.jumlah_paket} />
	<KeyNumberCard label={t('money.total', { year: anggaran.tahun })} angka={anggaran.total_hps} />
</div>
{#if status === 'memuat'}
	<p class="mt-4 flex items-center gap-2 text-muted" aria-busy="true">
		<LoaderCircle size={18} strokeWidth={1.75} aria-hidden="true" class="animate-spin" />
		{t('state.loading')}
	</p>
{:else if status === 'gagal'}
	<div class="mt-4"><ErrorState title={t('home.money.error')} /></div>
{:else if daftar.length === 0}
	<div class="mt-4"><EmptyState title={t('money.empty')} /></div>
{:else}
	<ul class="mt-4 grid gap-3 md:grid-cols-3 md:gap-4">
		{#each daftar as paket (paket.kode)}
			<li><PackageCard {paket} {kodeKecamatan} /></li>
		{/each}
	</ul>
{/if}
