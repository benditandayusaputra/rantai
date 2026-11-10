<script lang="ts">
	import {
		ArrowLeft,
		ArrowRight,
		Building2,
		ChartLine,
		ChevronRight,
		GraduationCap,
		House,
		ImageDown,
		Landmark,
		MapPin,
		MapPinCheck,
		MessageCircle,
		Share2,
		Wallet
	} from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import { alamatKartuKecamatan } from '$lib/api/klien';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import BagianAnggaran from '$lib/features/kecamatan/BagianAnggaran.svelte';
	import BagianTema from '$lib/features/kecamatan/BagianTema.svelte';
	import BagianWakilku from '$lib/features/kecamatan/BagianWakilku.svelte';
	import RingkasLayanan from '$lib/features/kecamatan/RingkasLayanan.svelte';
	import RingkasProyek from '$lib/features/kecamatan/RingkasProyek.svelte';
	import RingkasWakil from '$lib/features/kecamatan/RingkasWakil.svelte';
	import TombolPantau from '$lib/features/pantau/TombolPantau.svelte';
	import {
		bacaBagian,
		bagikanTautan,
		daftarBagian,
		tautanBagian,
		type Bagian
	} from '$lib/features/kecamatan/kecamatan';
	import { simpanWilayahku, wilayahku } from '$lib/features/wilayah/wilayahku.svelte';
	import { bahasaAktif, t, type KunciPesan } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const hasil = $derived(data.ringkasan);
	const bagianAktif = $derived(bacaBagian(page.url.searchParams.get('bagian')));
	let pesanBagikan = $state('');
	const halamanWilayah = $derived(
		`/${bahasaAktif}/kecamatan/${data.kode}${page.url.searchParams.has('kelurahan') ? `?kelurahan=${encodeURIComponent(page.url.searchParams.get('kelurahan') ?? '')}` : ''}`
	);
	const infoBagian: Record<Bagian, { ikon: typeof Wallet; judul: KunciPesan; isi: KunciPesan }> = {
		tema: { ikon: GraduationCap, judul: 'bagian.tema', isi: 'hub.services.body' },
		anggaran: { ikon: Wallet, judul: 'bagian.anggaran', isi: 'hub.projects.body' },
		wakilku: { ikon: Landmark, judul: 'bagian.wakilku', isi: 'hub.people.body' }
	};

	async function bagikan() {
		pesanBagikan = '';
		try {
			if ((await bagikanTautan(page.url.href)) === 'disalin') pesanBagikan = t('share.copied');
		} catch {
			pesanBagikan = '';
		}
	}
</script>

<svelte:head>
	<title
		>{hasil.ok
			? `${t('area.title', { name: hasil.data.kecamatan.nama })} · ${t('brand.name')}`
			: t('brand.name')}</title
	>
</svelte:head>

{#snippet pratinjau(bagian: Bagian, tautan: string, labelTautan: KunciPesan)}
	{@const info = infoBagian[bagian]}
	{@const Ikon = info.ikon}
	{#if hasil.ok}
		<section aria-labelledby="judul-{bagian}">
			<div class="flex flex-wrap items-end justify-between gap-3">
				<div class="flex items-start gap-3">
					<span class="icon-tile"><Ikon size={20} aria-hidden="true" /></span>
					<div>
						<h2 id="judul-{bagian}" class="t-h2">{t(info.judul)}</h2>
						<p class="mt-1 max-w-prose text-small text-muted">{t(info.isi)}</p>
					</div>
				</div>
				<a href={tautan} class="btn btn-outline btn-sm"
					>{t(labelTautan)}<ArrowRight size={16} aria-hidden="true" /></a
				>
			</div>
			<div class="mt-5">
				{#if bagian === 'wakilku'}<RingkasWakil wakilku={hasil.data.wakilku} />
				{:else if bagian === 'tema'}<RingkasLayanan tema={hasil.data.tema} />
				{:else if hasil.data.anggaran}<RingkasProyek
						anggaran={hasil.data.anggaran}
						kodeKecamatan={hasil.data.kecamatan.kode}
					/>
				{:else}<EmptyState title={t('money.empty')} />{/if}
			</div>
		</section>
	{/if}
{/snippet}

{#if !hasil.ok}
	<div class="wrap pt-10">
		<h1 class="sr-only">{t('state.error.title')}</h1>
		<ErrorState
			title={t('state.error.title')}
			description={hasil.galat.message}
			onretry={invalidateAll}
		/>
	</div>
{:else}
	{@const ringkasan = hasil.data}
	{@const kota = ringkasan.kecamatan.hierarki.find((h) => h.tingkat === 'kota')}
	<div class="border-b border-line bg-surface">
		<div class="wrap py-6 lg:py-8">
			<nav aria-label={t('nav.breadcrumb')}>
				<ol class="flex flex-wrap items-center gap-2 text-small text-muted">
					<li>
						<a
							href="/{bahasaAktif}"
							class="inline-flex min-h-touch items-center gap-1 hover:underline"
							><House size={14} aria-hidden="true" />{t('nav.home')}</a
						>
					</li>
					<li aria-hidden="true"><ChevronRight size={16} /></li>
					<li>
						{#if bagianAktif}<a
								href={halamanWilayah}
								class="inline-flex min-h-touch items-center font-semibold text-fg hover:underline"
								>{ringkasan.kecamatan.nama}</a
							>{:else}<span aria-current="page" class="font-semibold text-fg"
								>{ringkasan.kecamatan.nama}</span
							>{/if}
					</li>
					{#if bagianAktif}<li aria-hidden="true"><ChevronRight size={16} /></li>
						<li aria-current="page">{t(infoBagian[bagianAktif].judul)}</li>{/if}
				</ol>
			</nav>
			<div class="mt-3 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
				<div>
					<p class="t-over flex items-center gap-2">
						<MapPin size={16} aria-hidden="true" />{t('hub.selectedArea')}
					</p>
					<h1 class="t-h1 mt-2">{t('area.title', { name: ringkasan.kecamatan.nama })}</h1>
					{#if kota}<p class="mt-2 flex items-center gap-2 text-small text-muted">
							<Building2 size={16} aria-hidden="true" /><a
								href="/{bahasaAktif}/kota/{kota.kode}"
								class="link">{kota.nama}</a
							>
						</p>{/if}
				</div>
				<div class="flex flex-wrap items-center gap-2 self-start">
					{#if wilayahku() === ringkasan.kecamatan.kode}
						<span class="chip"><MapPinCheck size={16} aria-hidden="true" />{t('area.mine')}</span>
					{:else}
						<button
							type="button"
							class="btn btn-secondary"
							onclick={() => simpanWilayahku(ringkasan.kecamatan.kode)}
							><MapPinCheck size={18} aria-hidden="true" />{t('area.setMine')}</button
						>
					{/if}
					<a href="/{bahasaAktif}#cari-wilayah" class="btn btn-outline">{t('area.change')}</a>
				</div>
			</div>
			{#if ringkasan.kelurahan_disorot}<p class="mt-4 rounded-md bg-info-bg p-3 text-small">
					{t('area.highlighted', { name: ringkasan.kelurahan_disorot.nama })}
				</p>{/if}
		</div>
	</div>

	{#if bagianAktif === null}
		<div class="wrap flex flex-col gap-10 pt-8 lg:gap-12 lg:pt-10">
			{@render pratinjau('wakilku', tautanBagian(page.url, 'wakilku'), 'hub.people.more')}
			{@render pratinjau('anggaran', tautanBagian(page.url, 'anggaran'), 'hub.projects.more')}
			{@render pratinjau(
				'tema',
				kota ? `/${bahasaAktif}/kota/${kota.kode}` : tautanBagian(page.url, 'tema'),
				'hub.services.more'
			)}
			<div
				class="flex flex-col gap-4 rounded-lg bg-sunken p-5 md:flex-row md:items-center md:justify-between md:p-6"
			>
				<div>
					<h2 class="t-h3">{t('hub.askTitle')}</h2>
					<p class="mt-1 text-small text-muted">{t('hub.askBody')}</p>
				</div>
				<a
					href="/{bahasaAktif}/tanya?w={ringkasan.kecamatan.kode}"
					class="btn btn-secondary shrink-0"
					><MessageCircle size={18} aria-hidden="true" />{t('ask.title')}</a
				>
			</div>
			<div
				class="flex flex-col gap-4 rounded-lg border border-line p-5 md:flex-row md:items-center md:justify-between md:p-6"
			>
				<div>
					<h2 class="t-h3">{t('hub.trend.title')}</h2>
					<p class="mt-1 text-small text-muted">{t('hub.trend.body')}</p>
				</div>
				<a href="/{bahasaAktif}/dasbor" class="btn btn-outline shrink-0"
					><ChartLine size={18} aria-hidden="true" />{t('hub.trend.action')}</a
				>
			</div>
			<div>
				<div class="flex flex-wrap items-center justify-between gap-3 border-t border-line pt-4">
					<TombolPantau kodeKecamatan={ringkasan.kecamatan.kode} />
					<div class="flex flex-wrap gap-2">
						<button type="button" class="btn btn-ghost" onclick={bagikan}
							><Share2 size={18} aria-hidden="true" />{t('action.share')}</button
						><a
							href={alamatKartuKecamatan(ringkasan.kecamatan.kode, bahasaAktif)}
							download
							class="btn btn-ghost"
							><ImageDown size={18} aria-hidden="true" />{t('action.shareImage')}</a
						>
					</div>
				</div>
				<p class="mt-2 text-small text-success" aria-live="polite">{pesanBagikan}</p>
			</div>
		</div>
	{:else}
		<div class="wrap pt-5">
			<a href={halamanWilayah} class="btn btn-ghost btn-sm -ml-3"
				><ArrowLeft size={18} aria-hidden="true" />{t('hub.back')}</a
			>
			<nav aria-label={t('lens.nav')} class="mt-4 border-b border-line">
				<ul class="grid grid-cols-3 gap-1 md:flex md:gap-3">
					{#each daftarBagian as bagian (bagian)}
						<li>
							<a
								href={bagian === 'tema' && kota
									? `/${bahasaAktif}/kota/${kota.kode}`
									: tautanBagian(page.url, bagian)}
								aria-current={bagianAktif === bagian ? 'page' : undefined}
								class="-mb-px flex min-h-15 items-center justify-center border-b-link border-transparent px-2 py-3 text-center text-small font-semibold text-muted hover:text-fg md:px-4 aria-[current=page]:border-primary aria-[current=page]:text-primary"
								>{t(infoBagian[bagian].judul)}</a
							>
						</li>
					{/each}
				</ul>
			</nav>
		</div>
		<section class="wrap pt-7 lg:pt-8" aria-labelledby="judul-bagian">
			<h2 id="judul-bagian" class="t-h2">{t(infoBagian[bagianAktif].judul)}</h2>
			<p class="mt-2 max-w-prose text-muted">{t(infoBagian[bagianAktif].isi)}</p>
			<div class="mt-6">
				{#if bagianAktif === 'tema'}<BagianTema
						tema={ringkasan.tema}
						kodeKecamatan={ringkasan.kecamatan.kode}
					/>
				{:else if bagianAktif === 'anggaran'}<BagianAnggaran
						anggaran={ringkasan.anggaran}
						kodeKecamatan={ringkasan.kecamatan.kode}
						tema={ringkasan.tema}
					/>
				{:else}<BagianWakilku wakilku={ringkasan.wakilku} tema={ringkasan.tema} />{/if}
			</div>
		</section>
	{/if}
{/if}
