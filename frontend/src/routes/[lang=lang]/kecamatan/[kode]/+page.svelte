<script lang="ts">
	import {
		Building2,
		ChevronRight,
		GraduationCap,
		House,
		Landmark,
		MapPin,
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
	import TombolPantau from '$lib/features/pantau/TombolPantau.svelte';
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
	{#await data.ringkasan}
		<title>{t('brand.name')}</title>
	{:then hasil}
		<title
			>{hasil.ok
				? `${t('area.title', { name: hasil.data.kecamatan.nama })} · ${t('brand.name')}`
				: t('brand.name')}</title
		>
	{/await}
</svelte:head>

{#await data.ringkasan}
	<div class="border-b border-line bg-surface">
		<div class="wrap flex flex-col gap-4 pb-8 pt-6 lg:pt-10" aria-busy="true">
			<p class="sr-only">{t('state.loadingSummary')}</p>
			<div class="skel h-4 w-48 rounded-sm"></div>
			<div class="flex items-center gap-4">
				<div class="skel h-12 w-12 rounded-md"></div>
				<div class="skel h-10 w-72 max-w-full rounded-sm"></div>
			</div>
			<div class="skel h-11 w-80 max-w-full rounded-md"></div>
		</div>
	</div>
	<div class="wrap mt-8 grid grid-cols-2 gap-3 lg:grid-cols-4">
		<div class="skel h-28 rounded-lg"></div>
		<div class="skel h-28 rounded-lg"></div>
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
		{@const IkonBagian = infoBagian[bagianAktif].ikon}
		<div class="border-b border-line bg-surface">
			<div class="wrap pt-6 lg:pt-10">
				<nav aria-label={t('nav.breadcrumb')}>
					<ol class="flex flex-wrap items-center gap-1 text-small text-muted">
						<li>
							<a
								href="/{bahasaBawaan}"
								class="inline-flex items-center gap-1 underline-offset-4 hover:text-fg hover:underline"
							>
								<House size={14} strokeWidth={1.75} aria-hidden="true" />{t('nav.home')}
							</a>
						</li>
						<li aria-hidden="true"><ChevronRight size={16} strokeWidth={1.75} /></li>
						<li aria-current="page" class="font-semibold text-fg">{ringkasan.kecamatan.nama}</li>
					</ol>
				</nav>
				<div class="mt-4 flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
					<div class="flex items-start gap-4">
						<span class="icon-tile hidden h-14 w-14 md:inline-flex">
							<MapPin size={28} strokeWidth={1.75} aria-hidden="true" />
						</span>
						<div class="min-w-0">
							<h1 class="t-h1">{t('area.title', { name: ringkasan.kecamatan.nama })}</h1>
							{#if kota}
								<p class="mt-1 flex items-center gap-2 text-muted">
									<Building2 size={16} strokeWidth={1.75} aria-hidden="true" />{kota.nama}
								</p>
							{/if}
						</div>
					</div>
					<div class="flex flex-wrap items-center gap-2">
						<TombolPantau kodeKecamatan={ringkasan.kecamatan.kode} />
						<button
							type="button"
							class="btn btn-outline max-md:px-3 max-md:text-small"
							onclick={bagikan}
						>
							<Share2 size={18} strokeWidth={1.75} aria-hidden="true" />
							{t('action.share')}
						</button>
						<a
							href="/{bahasaBawaan}/tanya?w={encodeURIComponent(ringkasan.kecamatan.kode)}"
							class="btn btn-secondary max-md:px-3 max-md:text-small lg:hidden"
						>
							<MessageCircle size={18} strokeWidth={1.75} aria-hidden="true" />
							{t('ask.title')}
						</a>
						<button
							type="button"
							class="btn btn-secondary hidden lg:inline-flex"
							aria-haspopup="dialog"
							onclick={() => (isTanyaTerbuka = true)}
						>
							<MessageCircle size={18} strokeWidth={1.75} aria-hidden="true" />
							{t('ask.title')}
						</button>
						<p class="text-small text-success" aria-live="polite">{pesanBagikan}</p>
					</div>
				</div>
				{#if ringkasan.kelurahan_disorot}
					<p
						class="mt-4 flex items-start gap-2 rounded-md border-ctrl border-primary bg-secondary p-4 text-small text-secondary-on"
					>
						<MapPin size={18} strokeWidth={1.75} aria-hidden="true" class="mt-px shrink-0" />
						{t('area.highlighted', { name: ringkasan.kelurahan_disorot.nama })}
					</p>
				{/if}
				<PanelTanya
					bind:open={isTanyaTerbuka}
					konteks={{ kode: ringkasan.kecamatan.kode, nama: ringkasan.kecamatan.nama }}
				/>

				<nav aria-label={t('lens.nav')} class="mt-4">
					<ul class="grid grid-cols-3 md:flex md:gap-2">
						{#each daftarBagian as bagian (bagian)}
							{@const Ikon = infoBagian[bagian].ikon}
							<li>
								<a
									href={tautanBagian(page.url, bagian)}
									data-sveltekit-noscroll
									aria-current={bagianAktif === bagian ? 'page' : undefined}
									class="-mb-px flex h-15 flex-col items-center justify-center gap-1 border-b-link border-transparent px-2 text-small font-semibold text-muted transition-colors duration-fast ease-standard hover:text-fg md:h-14 md:flex-row md:gap-2 md:px-4 md:text-body aria-[current=page]:border-primary aria-[current=page]:text-primary"
								>
									<Ikon size={20} strokeWidth={1.75} aria-hidden="true" />
									{t(`bagian.${bagian}`)}
								</a>
							</li>
						{/each}
					</ul>
				</nav>
			</div>
		</div>

		<section class="wrap mt-8 lg:mt-10" aria-labelledby="judul-bagian">
			<div class="flex items-start gap-4">
				<span class="icon-tile hidden md:inline-flex">
					<IkonBagian size={20} strokeWidth={1.75} aria-hidden="true" />
				</span>
				<div>
					<h2 id="judul-bagian" class="t-h2">{t(`bagian.${bagianAktif}`)}</h2>
					<p class="t-lead mt-1 max-w-prose">{infoBagian[bagianAktif].tanya}</p>
				</div>
			</div>
			<div class="mt-6 lg:mt-8">
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
	{/if}
{/await}
