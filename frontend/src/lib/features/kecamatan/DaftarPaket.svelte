<script lang="ts">
	import { Calendar, ChevronDown, Download, Layers, ListFilter } from '@lucide/svelte';
	import { untrack } from 'svelte';
	import { replaceState } from '$app/navigation';
	import { page } from '$app/state';
	import { alamatUnduhPaket, ambilPaketKecamatan, tangkap, type FilterPaket } from '$lib/api/klien';
	import type { PaketRingkas, TemaKecamatan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import PackageCard from '$lib/components/PackageCard.svelte';
	import { daftarTahun } from '$lib/features/paket/paket';
	import { bahasaAktif, t } from '$lib/i18n';

	interface Props {
		kodeKecamatan: string;
		tema: TemaKecamatan[];
		tahunAwal: number;
	}

	let { kodeKecamatan, tema, tahunAwal }: Props = $props();

	const jumlahPilihanTahun = 3;
	const kueri = page.url.searchParams;
	let filterTema = $state(kueri.get('tema') ?? '');
	let filterTahun = $state(Number(kueri.get('tahun')) || untrack(() => tahunAwal));
	let daftar = $state<PaketRingkas[]>([]);
	let halaman = $state(1);
	let jumlahHalaman = $state(0);
	let total = $state(0);
	let status = $state<'memuat' | 'siap' | 'gagal'>('memuat');
	let urutanPermintaan = 0;

	async function muat(filter: FilterPaket, tambahkan: boolean) {
		const nomor = ++urutanPermintaan;
		status = 'memuat';
		const hasil = await tangkap(ambilPaketKecamatan(fetch, kodeKecamatan, filter));
		if (nomor !== urutanPermintaan) return;
		if (!hasil.ok) {
			status = 'gagal';
			return;
		}
		daftar = tambahkan ? [...daftar, ...hasil.data.data] : hasil.data.data;
		halaman = hasil.data.meta?.page ?? 1;
		jumlahHalaman = hasil.data.meta?.total_pages ?? 1;
		total = hasil.data.meta?.total ?? daftar.length;
		status = 'siap';
	}

	const filterAktif = $derived({ tema: filterTema || undefined, tahun: filterTahun });

	$effect(() => {
		muat({ ...filterAktif, halaman: 1 }, false);
	});

	function simpanFilter() {
		const url = new URL(page.url);
		if (filterTema) url.searchParams.set('tema', filterTema);
		else url.searchParams.delete('tema');
		url.searchParams.set('tahun', String(filterTahun));
		replaceState(url, {});
	}
</script>

<section aria-labelledby="judul-daftar-paket" class="flex flex-col gap-4">
	<div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
		<h3 id="judul-daftar-paket" class="flex items-center gap-2 text-h3 font-bold">
			<ListFilter size={20} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
			{t('money.listTitle')}
		</h3>
		<div class="grid grid-cols-2 gap-3 md:flex md:w-auto">
			<label class="flex flex-col gap-1 text-small font-semibold">
				<span class="inline-flex items-center gap-1">
					<Calendar size={14} strokeWidth={1.75} aria-hidden="true" />{t('money.filterYear')}
				</span>
				<select bind:value={filterTahun} onchange={simpanFilter} class="field md:w-36">
					{#each daftarTahun(tahunAwal, jumlahPilihanTahun) as tahun (tahun)}
						<option value={tahun}>{tahun}</option>
					{/each}
				</select>
			</label>
			<label class="flex flex-col gap-1 text-small font-semibold">
				<span class="inline-flex items-center gap-1">
					<Layers size={14} strokeWidth={1.75} aria-hidden="true" />{t('money.filterTheme')}
				</span>
				<select bind:value={filterTema} onchange={simpanFilter} class="field md:w-48">
					<option value="">{t('money.allThemes')}</option>
					{#each tema as item (item.kode)}
						<option value={item.kode}>{item.nama[bahasaAktif]}</option>
					{/each}
				</select>
			</label>
		</div>
	</div>
	{#if status === 'gagal'}
		<ErrorState
			title={t('error.sectionFailed')}
			onretry={() => muat({ ...filterAktif, halaman: 1 }, false)}
		/>
	{:else if status === 'siap' && !daftar.length}
		<EmptyState title={t('money.listEmpty')} />
	{:else}
		<div class="flex flex-wrap items-center justify-between gap-2">
			<p class="text-small text-muted" aria-live="polite">
				{status === 'siap' ? t('money.count.total', { n: total }) : ''}
			</p>
			<a
				href={alamatUnduhPaket(kodeKecamatan, filterTahun)}
				download
				class="btn btn-outline btn-sm"
				aria-describedby="catatan-unduh-paket"
			>
				<Download size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('action.downloadCsv')}
			</a>
		</div>
		<p id="catatan-unduh-paket" class="text-caption text-muted">
			{t('money.csvNote', { year: filterTahun })}
		</p>
		<ul class="grid gap-3 md:grid-cols-2 md:gap-4" aria-busy={status === 'memuat'}>
			{#each daftar as paket (paket.kode)}
				<li><PackageCard {paket} {kodeKecamatan} /></li>
			{/each}
		</ul>
		{#if status === 'memuat'}
			<div class="skel h-32 rounded-lg" aria-hidden="true"></div>
		{:else if halaman < jumlahHalaman}
			<button
				type="button"
				class="btn btn-secondary self-center max-md:btn-block"
				onclick={() => muat({ ...filterAktif, halaman: halaman + 1 }, true)}
			>
				<ChevronDown size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('money.loadMore')}
			</button>
		{/if}
	{/if}
</section>
