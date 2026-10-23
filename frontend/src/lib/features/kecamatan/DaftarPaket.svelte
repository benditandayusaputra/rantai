<script lang="ts">
	import { untrack } from 'svelte';
	import { replaceState } from '$app/navigation';
	import { page } from '$app/state';
	import { ambilPaketKecamatan, tangkap, type FilterPaket } from '$lib/api/klien';
	import type { PaketRingkas, TemaKecamatan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import PackageCard from '$lib/components/PackageCard.svelte';
	import { daftarTahun } from '$lib/features/paket/paket';
	import { bahasaBawaan, t } from '$lib/i18n';

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

	const kelasPilihan =
		'min-h-touch rounded-md border-ctrl border-line-strong bg-raised px-3 text-body hover:border-fg';
</script>

<section aria-labelledby="judul-daftar-paket" class="flex flex-col gap-4">
	<h3 id="judul-daftar-paket" class="t-h3">{t('money.listTitle')}</h3>
	<div class="flex flex-wrap gap-4">
		<label class="flex flex-col gap-1 text-small font-semibold">
			{t('money.filterYear')}
			<select bind:value={filterTahun} onchange={simpanFilter} class={kelasPilihan}>
				{#each daftarTahun(tahunAwal, jumlahPilihanTahun) as tahun (tahun)}
					<option value={tahun}>{tahun}</option>
				{/each}
			</select>
		</label>
		<label class="flex flex-col gap-1 text-small font-semibold">
			{t('money.filterTheme')}
			<select bind:value={filterTema} onchange={simpanFilter} class={kelasPilihan}>
				<option value="">{t('money.allThemes')}</option>
				{#each tema as item (item.kode)}
					<option value={item.kode}>{item.nama[bahasaBawaan]}</option>
				{/each}
			</select>
		</label>
	</div>
	{#if status === 'gagal'}
		<ErrorState
			title={t('error.sectionFailed')}
			onretry={() => muat({ ...filterAktif, halaman: 1 }, false)}
		/>
	{:else if status === 'siap' && !daftar.length}
		<EmptyState title={t('money.listEmpty')} />
	{:else}
		<p class="text-small text-muted" aria-live="polite">
			{status === 'siap' ? t('money.count.total', { n: total }) : ''}
		</p>
		<ul class="grid gap-3 md:grid-cols-2 md:gap-4" aria-busy={status === 'memuat'}>
			{#each daftar as paket (paket.kode)}
				<li><PackageCard {paket} /></li>
			{/each}
		</ul>
		{#if status === 'memuat'}
			<div class="skel h-32 rounded-lg" aria-hidden="true"></div>
		{:else if halaman < jumlahHalaman}
			<button
				type="button"
				class="inline-flex min-h-touch items-center self-start rounded-md bg-secondary px-4 font-bold text-secondary-on hover:bg-secondary-hover"
				onclick={() => muat({ ...filterAktif, halaman: halaman + 1 }, true)}
				>{t('money.loadMore')}</button
			>
		{/if}
	{/if}
</section>
