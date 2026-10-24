<script lang="ts">
	import {
		ChevronRight,
		GraduationCap,
		Landmark,
		MessageCircle,
		Share2,
		Wallet
	} from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import BagianAnggaran from '$lib/features/kecamatan/BagianAnggaran.svelte';
	import BagianTema from '$lib/features/kecamatan/BagianTema.svelte';
	import BagianWakilku from '$lib/features/kecamatan/BagianWakilku.svelte';
	import PanelTanya from '$lib/features/tanya/PanelTanya.svelte';
	import {
		bacaBagian,
		bagikanTautan,
		daftarBagian,
		tautanBagian,
		type Bagian
	} from '$lib/features/kecamatan/kecamatan';
	import { bahasaBawaan, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const bagianAktif = $derived(bacaBagian(page.url.searchParams.get('bagian')));
	let pesanBagikan = $state('');
	let isTanyaTerbuka = $state(page.url.searchParams.get('tanya') === '1');

	const infoBagian: Record<Bagian, { ikon: typeof Wallet; tanya: string }> = {
		tema: { ikon: GraduationCap, tanya: t('lens.services.q') },
		anggaran: { ikon: Wallet, tanya: t('lens.money.q') },
		wakilku: { ikon: Landmark, tanya: t('lens.accountability.q') }
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
	{#await data.ringkasan then hasil}
		{#if hasil.ok}
			<title>{t('area.title', { name: hasil.data.kecamatan.nama })} · {t('brand.name')}</title>
		{/if}
	{/await}
</svelte:head>

{#await data.ringkasan}
	<div class="wrap flex flex-col gap-4 pt-8" aria-busy="true">
		<p class="sr-only">{t('state.loadingSummary')}</p>
		<div class="skel h-4 w-48 rounded-sm"></div>
		<div class="skel h-10 w-72 rounded-sm"></div>
		<div class="skel h-14 w-full rounded-md"></div>
		<div class="grid grid-cols-2 gap-3">
			<div class="skel h-28 rounded-lg"></div>
			<div class="skel h-28 rounded-lg"></div>
		</div>
	</div>
{:then hasil}
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
		<div class="wrap pt-6 lg:pt-10">
			<nav aria-label={t('nav.breadcrumb')}>
				<ol class="flex flex-wrap items-center gap-1 text-small text-muted">
					<li><a href="/{bahasaBawaan}" class="underline">{t('nav.home')}</a></li>
					<li aria-hidden="true"><ChevronRight size={16} strokeWidth={1.75} /></li>
					<li aria-current="page">{ringkasan.kecamatan.nama}</li>
				</ol>
			</nav>
			<h1 class="t-h1 mt-3">{t('area.title', { name: ringkasan.kecamatan.nama })}</h1>
			{#if kota}<p class="mt-2 text-muted">{kota.nama}</p>{/if}
			{#if ringkasan.kelurahan_disorot}
				<p
					class="mt-4 rounded-md border-ctrl border-primary bg-secondary p-4 text-small text-secondary-on"
				>
					{t('area.highlighted', { name: ringkasan.kelurahan_disorot.nama })}
				</p>
			{/if}
			<div class="mt-4 flex flex-wrap items-center gap-3">
				<button
					type="button"
					class="inline-flex min-h-touch items-center gap-2 rounded-md border-ctrl border-line-strong px-4 font-semibold hover:border-fg hover:bg-sunken"
					onclick={bagikan}
				>
					<Share2 size={20} strokeWidth={1.75} aria-hidden="true" />
					{t('action.share')}
				</button>
				<a
					href="/{bahasaBawaan}/tanya?w={encodeURIComponent(ringkasan.kecamatan.kode)}"
					class="inline-flex min-h-touch items-center gap-2 rounded-md border-ctrl border-line-strong px-4 font-semibold hover:border-fg hover:bg-sunken lg:hidden"
				>
					<MessageCircle size={20} strokeWidth={1.75} aria-hidden="true" />
					{t('ask.title')}
				</a>
				<button
					type="button"
					class="hidden min-h-touch items-center gap-2 rounded-md border-ctrl border-line-strong px-4 font-semibold hover:border-fg hover:bg-sunken lg:inline-flex"
					aria-haspopup="dialog"
					onclick={() => (isTanyaTerbuka = true)}
				>
					<MessageCircle size={20} strokeWidth={1.75} aria-hidden="true" />
					{t('ask.title')}
				</button>
				<p class="text-small text-success" aria-live="polite">{pesanBagikan}</p>
			</div>
			<PanelTanya
				bind:open={isTanyaTerbuka}
				konteks={{ kode: ringkasan.kecamatan.kode, nama: ringkasan.kecamatan.nama }}
			/>

			<nav aria-label={t('lens.nav')} class="mt-8 border-b border-line">
				<ul class="grid grid-cols-3 lg:flex lg:gap-6">
					{#each daftarBagian as bagian (bagian)}
						{@const Ikon = infoBagian[bagian].ikon}
						<li>
							<a
								href={tautanBagian(page.url, bagian)}
								data-sveltekit-noscroll
								aria-current={bagianAktif === bagian ? 'page' : undefined}
								class="-mb-px flex h-15 flex-col items-center justify-center gap-1 border-b-link border-transparent px-2 text-small md:h-14 md:flex-row md:gap-2 md:text-body aria-[current=page]:border-primary aria-[current=page]:font-bold aria-[current=page]:text-primary"
							>
								<Ikon size={20} strokeWidth={1.75} aria-hidden="true" />
								{t(`bagian.${bagian}`)}
							</a>
						</li>
					{/each}
				</ul>
			</nav>

			<section class="mt-8" aria-labelledby="judul-bagian">
				<h2 id="judul-bagian" class="t-h2">{t(`bagian.${bagianAktif}`)}</h2>
				<p class="t-lead mt-2 max-w-prose">{infoBagian[bagianAktif].tanya}</p>
				<div class="mt-6">
					{#if bagianAktif === 'tema'}
						<BagianTema tema={ringkasan.tema} kodeKecamatan={ringkasan.kecamatan.kode} />
					{:else if bagianAktif === 'anggaran'}
						<BagianAnggaran
							anggaran={ringkasan.anggaran}
							kodeKecamatan={ringkasan.kecamatan.kode}
							tema={ringkasan.tema}
						/>
					{:else}
						<BagianWakilku wakilku={ringkasan.wakilku} />
					{/if}
				</div>
			</section>
		</div>
	{/if}
{/await}
