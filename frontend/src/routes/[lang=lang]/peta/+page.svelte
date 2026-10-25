<script lang="ts">
	import { ArrowRight, ChevronDown, Layers, Map, MousePointerClick } from '@lucide/svelte';
	import { replaceState } from '$app/navigation';
	import { page } from '$app/state';
	import type { FiturLaporan, PetaPaket, SekolahTerdekat } from '$lib/api/tipe';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import GapClassChip from '$lib/components/GapClassChip.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import TableToggle from '$lib/components/TableToggle.svelte';
	import {
		bacaLapisan,
		pusatJakarta,
		tulisLapisan,
		zoomKota,
		zoomSekolah
	} from '$lib/features/peta/lapisan';
	import KontrolLapisan from '$lib/features/peta/KontrolLapisan.svelte';
	import LegendaKelas from '$lib/features/peta/LegendaKelas.svelte';
	import PetaInteraktif from '$lib/features/peta/PetaInteraktif.svelte';
	import TabelPeta from '$lib/features/peta/TabelPeta.svelte';
	import { bahasaBawaan, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const kueri = page.url.searchParams;
	const lat = Number(kueri.get('lat'));
	const lng = Number(kueri.get('lng'));
	const isAdaTitik = kueri.has('lat') && kueri.has('lng') && !isNaN(lat) && !isNaN(lng);
	const pusat = isAdaTitik ? { lat, lng } : pusatJakarta;
	const zoom = isAdaTitik ? zoomSekolah : zoomKota;
	const lapisanAwal = bacaLapisan(kueri.get('lapisan'));

	let lapisan = $state(isAdaTitik ? { ...lapisanAwal, sekolah: true } : lapisanAwal);
	let isTabel = $state(kueri.get('tampilan') === 'tabel');
	let isKontrolTerbuka = $state(false);
	let terpilih = $state<string | null>(null);
	let paket = $state<PetaPaket['features']>([]);
	let sekolah = $state<SekolahTerdekat[]>([]);
	let laporan = $state<FiturLaporan[]>([]);
	let status = $state<'memuat' | 'siap' | 'gagal'>('memuat');
	let komponenPeta = $state<PetaInteraktif>();
	let tema = $state<Awaited<typeof data.tema> | null>(null);

	$effect(() => {
		data.tema.then((hasil) => (tema = hasil));
	});

	const fiturTema = $derived(tema?.ok ? tema.data : null);
	const detailTerpilih = $derived(
		fiturTema?.features.find((f) => f.properties.kode === terpilih)?.properties ?? null
	);

	function simpanKeURL() {
		const url = new URL(page.url);
		url.searchParams.set('lapisan', tulisLapisan(lapisan));
		if (isTabel) url.searchParams.set('tampilan', 'tabel');
		else url.searchParams.delete('tampilan');
		replaceState(url, {});
	}

	function gantiTampilan(tabel: boolean) {
		isTabel = tabel;
		simpanKeURL();
	}

	$effect(() => {
		if (!isTabel) komponenPeta?.sesuaikanUkuran();
	});
</script>

<svelte:head>
	<title>{t('map.title')} · {t('brand.name')}</title>
</svelte:head>

<h1 class="sr-only">{t('map.title')}</h1>

<div class="relative">
	<div class="relative {isTabel ? 'hidden' : 'h-peta-ponsel lg:h-peta'}">
		<PetaInteraktif
			bind:this={komponenPeta}
			tema={fiturTema}
			{lapisan}
			{pusat}
			{zoom}
			{terpilih}
			onpilih={(kode) => (terpilih = kode)}
			bind:paket
			bind:sekolah
			bind:laporan
			bind:status
		/>
		{#if status === 'memuat'}
			<div class="skel absolute inset-0 flex items-center justify-center" aria-busy="true">
				<p class="rounded-md bg-raised px-4 py-3 text-small shadow-2">
					{t('map.loading')}
					{t('map.tableWhileLoading')}
				</p>
			</div>
		{:else if status === 'gagal'}
			<div class="absolute inset-0 flex items-center justify-center bg-page p-4">
				<ErrorState title={t('map.failed')} />
			</div>
		{/if}
	</div>

	<div
		class="z-map-ui flex flex-col gap-2 {isTabel
			? 'wrap mt-6'
			: 'absolute left-3 top-3 lg:glass lg:left-8 lg:top-8 lg:w-66 lg:gap-3 lg:rounded-lg lg:border lg:border-line lg:p-4 lg:shadow-2'}"
	>
		<div class="flex flex-col items-start gap-2 md:flex-row {isTabel ? '' : 'lg:order-last'}">
			<button
				type="button"
				class="btn btn-outline text-small shadow-2 lg:hidden"
				aria-expanded={isKontrolTerbuka}
				aria-controls="kontrol-lapisan"
				onclick={() => (isKontrolTerbuka = !isKontrolTerbuka)}
			>
				<Layers size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('map.layers')}
				<ChevronDown
					size={16}
					strokeWidth={1.75}
					aria-hidden="true"
					class="transition-transform duration-fast ease-standard {isKontrolTerbuka
						? 'rotate-180'
						: ''}"
				/>
			</button>
			<TableToggle
				bind:pressed={() => isTabel, gantiTampilan}
				controls="isi-peta"
				labelKembali={t('action.viewAsMap')}
				ikonKembali={Map}
			/>
		</div>
		<div
			id="kontrol-lapisan"
			class="flex-col gap-4 lg:flex {isKontrolTerbuka ? 'flex' : 'hidden'} {isTabel
				? 'card p-4 md:flex-row md:gap-10'
				: 'glass rounded-lg border border-line p-3 shadow-2 lg:border-0 lg:bg-transparent lg:p-0 lg:shadow-none lg:backdrop-blur-none'}"
		>
			<KontrolLapisan bind:lapisan onchange={simpanKeURL} />
			{#if lapisan.tema}<LegendaKelas />{/if}
		</div>
	</div>
</div>

<div id="isi-peta" class="wrap flex flex-col gap-6 py-6">
	{#if isTabel}
		<TabelPeta {lapisan} tema={fiturTema} {paket} {sekolah} {laporan} {terpilih} />
	{/if}
	<section aria-live="polite" class="flex flex-col gap-4">
		{#if detailTerpilih}
			<h2 class="t-h2">{detailTerpilih.nama}</h2>
			<div class="flex flex-wrap items-center gap-3">
				<GapClassChip kelas={detailTerpilih.kelas} />
				<span class="text-small text-muted">
					{t('tema.period', { period: detailTerpilih.periode })}
				</span>
			</div>
			<div class="grid grid-cols-2 gap-3 lg:grid-cols-4 lg:gap-5">
				<KeyNumberCard label={t('tema.needs')} angka={detailTerpilih.kebutuhan} />
				<KeyNumberCard label={t('tema.capacity')} angka={detailTerpilih.kapasitas} />
				<KeyNumberCard label={t('tema.gap')} angka={detailTerpilih.kesenjangan} />
				{#if detailTerpilih.rasio}
					<KeyNumberCard label={t('tema.ratio')} angka={detailTerpilih.rasio} />
				{/if}
			</div>
			<a href="/{bahasaBawaan}/kecamatan/{detailTerpilih.kode}" class="btn btn-primary self-start">
				{t('map.openArea')}
				<ArrowRight size={18} strokeWidth={1.75} aria-hidden="true" />
			</a>
		{:else}
			<p class="flex items-center gap-2 text-small text-muted">
				<MousePointerClick size={18} strokeWidth={1.75} aria-hidden="true" class="shrink-0" />
				{t('map.select')}
			</p>
		{/if}
		<p class="text-caption text-muted">{t('map.basemapNote')}</p>
	</section>
</div>
