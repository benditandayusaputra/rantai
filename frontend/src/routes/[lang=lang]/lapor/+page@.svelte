<script lang="ts">
	import {
		ArrowLeft,
		ArrowRight,
		LoaderCircle,
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
	<main id="konten" class="flex min-h-dvh items-center justify-center bg-page px-4 py-10">
		<div class="card flex w-full max-w-form flex-col gap-4 p-6 shadow-2 md:p-8">
			<span class="icon-tile h-14 w-14"
				><Lock size={28} strokeWidth={1.75} aria-hidden="true" /></span
			>
			<h1 class="t-h1">{t('report.signIn.title')}</h1>
			<p class="text-muted">{t('report.signIn.body')}</p>
			<a href={tautanMasuk} class="btn btn-primary btn-lg btn-block">
				<LogIn size={20} strokeWidth={1.75} aria-hidden="true" />
				{t('action.signIn')}
			</a>
			<a href="/{bahasaAktif}" class="btn btn-text btn-block">
				<ArrowLeft size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('action.back')}
			</a>
		</div>
	</main>
{:else}
	<div class="flex min-h-dvh flex-col bg-page">
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
