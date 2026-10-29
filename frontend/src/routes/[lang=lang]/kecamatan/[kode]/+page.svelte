<script lang="ts">
	import {
		ArrowLeft,
		ArrowRight,
		Building2,
		ChevronRight,
		GraduationCap,
		House,
		ImageDown,
		Landmark,
		MapPin,
		MessageCircle,
		Share2,
		Wallet
	} from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import { alamatKartuKecamatan } from '$lib/api/klien';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import BagianAnggaran from '$lib/features/kecamatan/BagianAnggaran.svelte';
	import BagianTema from '$lib/features/kecamatan/BagianTema.svelte';
	import BagianWakilku from '$lib/features/kecamatan/BagianWakilku.svelte';
	import TombolPantau from '$lib/features/pantau/TombolPantau.svelte';
	import {
		bacaBagian,
		bagikanTautan,
		daftarBagian,
		tautanBagian,
		type Bagian
	} from '$lib/features/kecamatan/kecamatan';
	import { bahasaAktif, t, type KunciPesan } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const hasil = $derived(data.ringkasan);
	const bagianAktif = $derived(bacaBagian(page.url.searchParams.get('bagian')));
	let pesanBagikan = $state('');
	const halamanWilayah = $derived(
		`/${bahasaAktif}/kecamatan/${data.kode}${page.url.searchParams.has('kelurahan') ? `?kelurahan=${encodeURIComponent(page.url.searchParams.get('kelurahan') ?? '')}` : ''}`
	);
	const infoBagian: Record<
		Bagian,
		{ ikon: typeof Wallet; judul: KunciPesan; isi: KunciPesan; cakupan: KunciPesan }
	> = {
		tema: {
			ikon: GraduationCap,
			judul: 'bagian.tema',
			isi: 'hub.services.body',
			cakupan: 'hub.services.scope'
		},
		anggaran: {
			ikon: Wallet,
			judul: 'bagian.anggaran',
			isi: 'hub.projects.body',
			cakupan: 'hub.projects.scope'
		},
		wakilku: {
			ikon: Landmark,
			judul: 'bagian.wakilku',
			isi: 'hub.people.body',
			cakupan: 'hub.people.scope'
		}
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
				<a href="/{bahasaAktif}#cari-wilayah" class="btn btn-outline self-start"
					>{t('area.change')}</a
				>
			</div>
			{#if ringkasan.kelurahan_disorot}<p class="mt-4 rounded-md bg-info-bg p-3 text-small">
					{t('area.highlighted', { name: ringkasan.kelurahan_disorot.nama })}
				</p>{/if}
		</div>
	</div>

	{#if bagianAktif === null}
		<section class="wrap pt-8 lg:pt-10" aria-labelledby="judul-pilih-informasi">
			<h2 id="judul-pilih-informasi" class="t-h2">{t('hub.title')}</h2>
			<p class="mt-2 max-w-prose text-muted">{t('hub.lead', { name: ringkasan.kecamatan.nama })}</p>
			<ul class="mt-6 grid gap-4 lg:grid-cols-3 lg:gap-6">
				{#each daftarBagian as bagian (bagian)}
					{@const info = infoBagian[bagian]}
					{@const Ikon = info.ikon}
					<li>
						<a
							href={bagian === 'tema' && kota
								? `/${bahasaAktif}/kota/${kota.kode}`
								: tautanBagian(page.url, bagian)}
							class="card-link group flex h-full flex-col p-5 md:p-6"
						>
							<span class="icon-tile h-12 w-12"><Ikon size={24} aria-hidden="true" /></span>
							<h3 class="t-h3 mt-5">{t(info.judul)}</h3>
							<p class="mt-2 text-small text-muted">{t(info.isi)}</p>
							<p class="mt-4 text-caption text-muted">{t(info.cakupan)}</p>
							<span
								class="mt-auto flex min-h-touch items-center justify-between gap-3 border-t border-line pt-4 text-small font-semibold group-hover:text-primary"
								><span>{t('hub.open')}</span><ArrowRight size={18} aria-hidden="true" /></span
							>
						</a>
					</li>
				{/each}
			</ul>
			<div
				class="mt-8 flex flex-col gap-4 rounded-lg bg-sunken p-5 md:flex-row md:items-center md:justify-between md:p-6"
			>
				<div>
					<h3 class="t-h3">{t('hub.askTitle')}</h3>
					<p class="mt-1 text-small text-muted">{t('hub.askBody')}</p>
				</div>
				<a
					href="/{bahasaAktif}/tanya?w={ringkasan.kecamatan.kode}"
					class="btn btn-secondary shrink-0"
					><MessageCircle size={18} aria-hidden="true" />{t('ask.title')}</a
				>
			</div>
			<div class="mt-5 flex flex-wrap items-center justify-between gap-3 border-t border-line pt-4">
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
		</section>
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
