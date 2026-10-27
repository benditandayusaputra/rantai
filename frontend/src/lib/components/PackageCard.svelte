<script lang="ts">
	import { ArrowRight, Hammer, MapPin } from '@lucide/svelte';
	import type { PaketRingkas } from '$lib/api/tipe';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { ikonSektor, kunciSektor } from '$lib/features/paket/paket';
	import { bahasaAktif, t } from '$lib/i18n';
	import CertaintyBadge from './CertaintyBadge.svelte';
	import SimRibbon from './SimRibbon.svelte';

	let { paket, kodeKecamatan }: { paket: PaketRingkas; kodeKecamatan?: string } = $props();

	const nilai = $derived(paket.nilai_kontrak ?? paket.hps);
	const tampil = $derived(nilai ? tampilkanAngka(nilai) : null);
	const IkonSektor = $derived(ikonSektor[paket.sektor]);
</script>

<a
	href={`/${bahasaAktif}/paket/${paket.kode}${kodeKecamatan ? `?w=${encodeURIComponent(kodeKecamatan)}` : ''}`}
	class="card-link group flex h-full flex-col overflow-hidden {paket.is_simulation
		? 'border-ctrl border-cert-sim'
		: ''}"
>
	{#if paket.is_simulation}<SimRibbon />{/if}
	<span class="flex flex-1 flex-col gap-3 p-4 md:p-5">
		<span class="flex items-start gap-3">
			<span class="icon-tile h-9 w-9">
				<IkonSektor size={18} strokeWidth={1.75} aria-hidden="true" />
			</span>
			<span class="flex min-w-0 flex-col gap-0.5">
				<span class="text-caption text-muted">
					{t(kunciSektor[paket.sektor])} · {paket.tahun_anggaran} · {paket.tahap}
				</span>
				<span class="break-words font-bold">{paket.nama}</span>
			</span>
			<ArrowRight
				size={18}
				strokeWidth={1.75}
				aria-hidden="true"
				class="ml-auto mt-0.5 shrink-0 text-muted transition-transform duration-fast ease-standard group-hover:translate-x-1 group-hover:text-fg"
			/>
		</span>
		{#if tampil && nilai}
			<span class="flex flex-wrap items-center gap-2">
				<span class="n-md">{tampil.nilai}</span>
				<CertaintyBadge level={nilai.label} />
			</span>
		{/if}
		<span class="mt-auto flex flex-col gap-1 border-t border-line pt-3 text-small text-muted">
			<span class="flex items-start gap-2">
				<Hammer size={16} strokeWidth={1.75} aria-hidden="true" class="mt-0.5 shrink-0" />
				{paket.penyedia ? t('money.vendor', { name: paket.penyedia }) : t('pkg.noVendor')}
			</span>
			{#if paket.lokasi.nama_wilayah}
				<span class="flex items-center gap-2">
					<MapPin size={16} strokeWidth={1.75} aria-hidden="true" class="shrink-0" />
					{paket.lokasi.nama_wilayah}
				</span>
			{/if}
		</span>
	</span>
</a>
