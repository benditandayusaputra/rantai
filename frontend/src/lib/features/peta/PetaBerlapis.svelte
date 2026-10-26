<script lang="ts">
	import { Layers, Maximize2, Minimize2 } from '@lucide/svelte';
	import type { PetaTema } from '$lib/api/tipe';
	import { t, type KunciPesan } from '$lib/i18n';
	import type { PetaDasar } from './gaya';
	import KontrolLapisan from './KontrolLapisan.svelte';
	import LegendaKelas from './LegendaKelas.svelte';
	import LegendaPenanda from './LegendaPenanda.svelte';
	import PetaInteraktif from './PetaInteraktif.svelte';
	import PilihPetaDasar from './PilihPetaDasar.svelte';
	import { batasJakarta, pusatJakarta, zoomKota, type StatusLapisan } from './lapisan';

	interface Props {
		tema: PetaTema | null;
		temaSekolah: string;
		terpilih?: string | null;
		onpilih: (kode: string) => void;
		lapisan: StatusLapisan;
		petaDasar: PetaDasar;
		onchange?: () => void;
		nama: string;
		judulLegenda: KunciPesan;
		class?: string;
	}

	let {
		tema,
		temaSekolah,
		terpilih = null,
		onpilih,
		lapisan = $bindable(),
		petaDasar = $bindable(),
		onchange = () => {},
		nama,
		judulLegenda,
		class: kelas = 'h-96'
	}: Props = $props();

	let status = $state<'memuat' | 'siap' | 'gagal'>('memuat');
	let isLayarPenuh = $state(false);
	let komponenPeta = $state<PetaInteraktif>();

	$effect(() => {
		if (!isLayarPenuh) return;
		const luapanAwal = document.body.style.overflow;
		document.body.style.overflow = 'hidden';
		return () => {
			document.body.style.overflow = luapanAwal;
		};
	});

	$effect(() => {
		void isLayarPenuh;
		komponenPeta?.sesuaikanUkuran();
	});

	function tutupDenganEscape(event: KeyboardEvent) {
		if (event.key === 'Escape' && isLayarPenuh) isLayarPenuh = false;
	}
</script>

<svelte:window onkeydown={tutupDenganEscape} />

<div
	class={isLayarPenuh
		? 'fixed inset-0 isolate z-dialog bg-page'
		: `relative isolate overflow-hidden rounded-md border border-line ${kelas}`}
	role={isLayarPenuh ? 'dialog' : undefined}
	aria-modal={isLayarPenuh ? 'true' : undefined}
	aria-label={isLayarPenuh ? t('map.region') : undefined}
>
	<PetaInteraktif
		bind:this={komponenPeta}
		{tema}
		{temaSekolah}
		{lapisan}
		{petaDasar}
		pusat={pusatJakarta}
		zoom={zoomKota}
		batas={batasJakarta}
		isKooperatif={!isLayarPenuh}
		{terpilih}
		{onpilih}
		bind:status
	/>
	{#if status === 'memuat'}
		<div class="skel absolute inset-0" aria-busy="true">
			<p class="sr-only">{t('map.loading')}</p>
		</div>
	{:else if status === 'gagal'}
		<p
			class="absolute inset-0 flex items-center justify-center bg-sunken p-4 text-small text-muted"
		>
			{t('home.map.error')}
		</p>
	{/if}
	<details class="absolute left-2 top-2 z-map-ui">
		<summary
			class="btn btn-outline btn-sm glass cursor-pointer list-none shadow-2 [&::-webkit-details-marker]:hidden"
		>
			<Layers size={16} strokeWidth={1.75} aria-hidden="true" />
			{t('map.layers')}
		</summary>
		<div
			class="glass mt-2 flex w-64 flex-col gap-4 overflow-y-auto rounded-lg border border-line p-3 shadow-2 {isLayarPenuh
				? 'max-h-[calc(100dvh-theme(spacing.24))]'
				: 'max-h-80'}"
		>
			<PilihPetaDasar bind:petaDasar {onchange} {nama} />
			<KontrolLapisan bind:lapisan {onchange} />
			{#if lapisan.tema}<LegendaKelas judul={judulLegenda} />{/if}
			<LegendaPenanda {lapisan} />
		</div>
	</details>
	<button
		type="button"
		class="btn btn-outline btn-sm glass absolute right-2 top-2 z-map-ui shadow-2"
		onclick={() => (isLayarPenuh = !isLayarPenuh)}
	>
		{#if isLayarPenuh}
			<Minimize2 size={16} strokeWidth={1.75} aria-hidden="true" />
			{t('map.fullscreen.exit')}
		{:else}
			<Maximize2 size={16} strokeWidth={1.75} aria-hidden="true" />
			<span class="max-md:sr-only">{t('map.fullscreen')}</span>
		{/if}
	</button>
</div>
