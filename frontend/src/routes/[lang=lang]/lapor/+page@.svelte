<script lang="ts">
	import {
		ArrowLeft,
		ArrowRight,
		Camera,
		LoaderCircle,
		MapPin,
		Users,
		Lock,
		LogIn,
		Send,
		SkipForward
	} from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { tangkap } from '$lib/api/klien';
	import { kirimLaporan } from '$lib/api/laporan';
	import ActionBar from '$lib/components/ActionBar.svelte';
	import IlustrasiLapor from '$lib/components/IlustrasiLapor.svelte';
	import LatarHias from '$lib/components/LatarHias.svelte';
	import Logo from '$lib/components/Logo.svelte';
	import StepHeader from '$lib/components/StepHeader.svelte';
	import { pulihkanSesi, sesiAktif } from '$lib/features/akun/sesi.svelte';
	import { aturUlangDraf, draf } from '$lib/features/lapor/draf.svelte';
	import {
		bacaLangkah,
		isLangkahLengkap,
		langkahTerjauh,
		susunFormulir,
		totalLangkah
	} from '$lib/features/lapor/lapor';
	import LangkahCatatan from '$lib/features/lapor/LangkahCatatan.svelte';
	import LangkahFoto from '$lib/features/lapor/LangkahFoto.svelte';
	import LangkahKondisi from '$lib/features/lapor/LangkahKondisi.svelte';
	import LangkahLokasi from '$lib/features/lapor/LangkahLokasi.svelte';
	import LangkahProyek from '$lib/features/lapor/LangkahProyek.svelte';
	import LangkahRingkasan from '$lib/features/lapor/LangkahRingkasan.svelte';
	import { bahasaAktif, t, type KunciPesan } from '$lib/i18n';

	let isSesiDiperiksa = $state(false);
	let isMengirim = $state(false);
	let pesanGalat = $state('');

	const sesi = $derived(sesiAktif());
	const langkah = $derived(
		Math.min(bacaLangkah(page.url.searchParams.get('langkah')), langkahTerjauh(draf))
	);
	const tautanMasuk = $derived(
		`/${bahasaAktif}/masuk?lanjut=${encodeURIComponent(`/${bahasaAktif}/lapor`)}`
	);

	const langkahLapor = [
		{ ikon: MapPin, teks: t('report.signIn.step.locate') },
		{ ikon: Camera, teks: t('report.signIn.step.photo') },
		{ ikon: Users, teks: t('report.signIn.step.witness') }
	];

	onMount(async () => {
		await pulihkanSesi(fetch);
		isSesiDiperiksa = true;
	});

	function keLangkah(tujuan: number) {
		pesanGalat = '';
		void goto(`?langkah=${tujuan}`);
	}

	function tutup() {
		aturUlangDraf();
		void goto(`/${bahasaAktif}`);
	}

	function lewatiCatatan() {
		draf.catatan = '';
		keLangkah(6);
	}

	async function kirim() {
		if (!sesi) return;
		isMengirim = true;
		const hasil = await tangkap(kirimLaporan(fetch, sesi.token_akses, susunFormulir(draf)));
		isMengirim = false;
		if (!hasil.ok) {
			pesanGalat = hasil.galat.message;
			return;
		}
		aturUlangDraf();
		await goto(`/${bahasaAktif}/lapor/terkirim?id=${hasil.data.id}`);
	}

	const labelUtama = $derived(
		langkah === 3 && draf.foto.length > 0
			? t('report.s3.continueWith', { n: draf.foto.length })
			: langkah === totalLangkah
				? t('report.s6.submit')
				: t('action.next')
	);
</script>

<svelte:head>
	<title>{t('report.flowTitle')} · {t('brand.name')}</title>
</svelte:head>

{#if !isSesiDiperiksa}
	<p class="flex min-h-dvh items-center justify-center" aria-live="polite">
		<LoaderCircle size={28} strokeWidth={1.75} class="animate-spin text-muted" aria-hidden="true" />
		<span class="sr-only">{t('state.loading')}</span>
	</p>
{:else if !sesi}
	<div class="relative isolate flex min-h-dvh flex-col bg-page">
		<LatarHias />
		<header class="wrap flex h-16 items-center justify-between gap-4 lg:h-18">
			<Logo />
			<a href="/{bahasaAktif}" class="btn btn-ghost btn-sm">
				<ArrowLeft size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('nav.home')}
			</a>
		</header>
		<main
			id="konten"
			class="wrap flex flex-1 flex-col items-center justify-center gap-6 pb-12 lg:flex-row lg:gap-14"
		>
			<IlustrasiLapor class="w-60 max-w-full md:w-80 lg:w-full lg:max-w-lg lg:flex-1" />
			<div class="card w-full max-w-form p-6 shadow-3 md:p-8">
				<p class="t-over flex items-center gap-2 text-primary">
					<Camera size={16} strokeWidth={1.75} aria-hidden="true" />
					{t('report.flowTitle')}
				</p>
				<h1 class="t-h1 mt-2">{t('report.signIn.title')}</h1>
				<p class="mt-2 text-muted">{t('report.signIn.body')}</p>
				<ol class="mt-5 flex flex-col gap-3">
					{#each langkahLapor as item, i (item.teks)}
						{@const Ikon = item.ikon}
						<li class="flex items-center gap-3 text-small font-semibold">
							<span class="icon-tile relative h-10 w-10">
								<Ikon size={18} strokeWidth={1.75} aria-hidden="true" />
								<span
									class="num absolute -right-1.5 -top-1.5 inline-flex h-5 w-5 items-center justify-center rounded-pill bg-primary text-caption text-primary-on"
									aria-hidden="true">{i + 1}</span
								>
							</span>
							{item.teks}
						</li>
					{/each}
				</ol>
				<a href={tautanMasuk} class="btn btn-primary btn-lg btn-block mt-6">
					<LogIn size={20} strokeWidth={1.75} aria-hidden="true" />
					{t('action.signIn')}
				</a>
				<p class="mt-5 flex items-start gap-2 border-t border-line pt-4 text-caption text-muted">
					<Lock size={16} strokeWidth={1.75} aria-hidden="true" class="shrink-0" />
					{t('auth.privacy')}
				</p>
			</div>
		</main>
	</div>
{:else}
	<div class="relative isolate flex min-h-dvh flex-col bg-page">
		<LatarHias class="max-md:hidden" />
		<StepHeader
			step={langkah}
			total={totalLangkah}
			stepLabel={t(`report.stepName.${langkah}` as KunciPesan)}
			onback={() => keLangkah(langkah - 1)}
			onclose={tutup}
		/>
		<main id="konten" class="mx-auto flex w-full max-w-form flex-1 flex-col px-4 py-6 md:py-10">
			<div class="flex flex-1 flex-col md:card md:flex-none md:p-8">
				{#if langkah === 1}
					<LangkahLokasi onlanjut={() => keLangkah(2)} />
				{:else if langkah === 2}
					<LangkahProyek />
				{:else if langkah === 3}
					<LangkahFoto />
				{:else if langkah === 4}
					<LangkahKondisi />
				{:else if langkah === 5}
					<LangkahCatatan />
				{:else}
					<LangkahRingkasan onubah={keLangkah} />
				{/if}
				{#if pesanGalat}<p role="alert" class="mt-4 text-danger">{pesanGalat}</p>{/if}
			</div>
		</main>
		{#if langkah > 1}
			<ActionBar>
				<button
					type="button"
					class="btn btn-primary btn-lg btn-block"
					disabled={!isLangkahLengkap(langkah, draf) || isMengirim}
					aria-busy={isMengirim}
					onclick={() => (langkah === totalLangkah ? kirim() : keLangkah(langkah + 1))}
				>
					{#if isMengirim}
						<LoaderCircle size={20} strokeWidth={1.75} class="animate-spin" aria-hidden="true" />
						{t('report.s6.sending')}
					{:else if langkah === totalLangkah}
						<Send size={20} strokeWidth={1.75} aria-hidden="true" />
						{labelUtama}
					{:else}
						{labelUtama}
						<ArrowRight size={20} strokeWidth={1.75} aria-hidden="true" />
					{/if}
				</button>
				{#if langkah === 5}
					<button type="button" class="btn btn-text btn-block" onclick={lewatiCatatan}>
						<SkipForward size={18} strokeWidth={1.75} aria-hidden="true" />
						{t('action.skip')}
					</button>
				{/if}
			</ActionBar>
		{/if}
	</div>
{/if}
