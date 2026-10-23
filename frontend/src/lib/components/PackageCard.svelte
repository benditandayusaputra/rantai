<script lang="ts">
	import type { PaketRingkas } from '$lib/api/tipe';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { kunciSektor } from '$lib/features/paket/paket';
	import { bahasaBawaan, t } from '$lib/i18n';
	import CertaintyBadge from './CertaintyBadge.svelte';
	import SimRibbon from './SimRibbon.svelte';

	let { paket }: { paket: PaketRingkas } = $props();

	const nilai = $derived(paket.nilai_kontrak ?? paket.hps);
	const tampil = $derived(nilai ? tampilkanAngka(nilai) : null);
</script>

<a
	href="/{bahasaBawaan}/paket/{paket.kode}"
	class="flex flex-col overflow-hidden rounded-lg border bg-surface shadow-1 transition-colors duration-fast ease-standard hover:border-line-strong {paket.is_simulation
		? 'border-ctrl border-cert-sim'
		: 'border-line'}"
>
	{#if paket.is_simulation}<SimRibbon />{/if}
	<span class="flex flex-col gap-2 p-4 md:p-5">
		<span class="text-caption text-muted">
			{t(kunciSektor[paket.sektor])} · {paket.tahun_anggaran} · {paket.tahap}
		</span>
		<span class="font-bold">{paket.nama}</span>
		{#if tampil && nilai}
			<span class="flex flex-wrap items-center gap-2">
				<span class="n-md">{tampil.nilai}</span>
				<CertaintyBadge level={nilai.label} />
			</span>
		{/if}
		<span class="text-small text-muted">
			{paket.penyedia ? t('money.vendor', { name: paket.penyedia }) : t('pkg.noVendor')}
		</span>
		{#if paket.lokasi.nama_wilayah}
			<span class="text-small text-muted">{paket.lokasi.nama_wilayah}</span>
		{/if}
	</span>
</a>
