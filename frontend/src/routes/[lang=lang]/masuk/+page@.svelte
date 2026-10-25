<script lang="ts">
	import {
		ArrowLeft,
		ArrowRight,
		Bell,
		LoaderCircle,
		Lock,
		Mail,
		PenLine,
		RefreshCw,
		Scale,
		ShieldCheck,
		UserCheck
	} from '@lucide/svelte';
	import { tick } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { kirimKodeMasuk, verifikasiKodeMasuk } from '$lib/api/auth';
	import { tangkap, type GalatAPI } from '$lib/api/klien';
	import Logo from '$lib/components/Logo.svelte';
	import IlustrasiMasuk from '$lib/components/IlustrasiMasuk.svelte';
	import OtpInput from '$lib/components/OtpInput.svelte';
	import {
		detikKirimUlang,
		detikTerkunci,
		formatHitungMundur,
		keadaanDariGalat,
		tujuanSetelahMasuk,
		type KeadaanKode
	} from '$lib/features/akun/masuk';
	import { simpanSesi } from '$lib/features/akun/sesi.svelte';
	import { bahasaBawaan, t } from '$lib/i18n';

	const idKode = 'kode-masuk';
	const idPesanKode = 'pesan-kode';
	const milidetik = 1000;

	let tahap = $state<'email' | 'kode'>('email');
	let email = $state('');
	let kode = $state('');
	let keadaan = $state<KeadaanKode>('siap');
	let pesanGalat = $state('');
	let isMengirim = $state(false);
	let isMemeriksa = $state(false);
	let kirimUlangPada = $state(0);
	let terkunciSampai = $state(0);
	let sekarang = $state(Date.now());

	const sisaKirimUlang = $derived((kirimUlangPada - sekarang) / milidetik);
	const sisaTerkunci = $derived((terkunciSampai - sekarang) / milidetik);
	const isTerkunci = $derived(keadaan === 'terkunci' && sisaTerkunci > 0);
	const galatKode = $derived(
		keadaan === 'salah' ? 'invalid' : keadaan === 'kedaluwarsa' ? 'expired' : null
	);

	$effect(() => {
		const detak = setInterval(() => (sekarang = Date.now()), milidetik);
		return () => clearInterval(detak);
	});

	function tanganiGalat(galat: GalatAPI) {
		keadaan = keadaanDariGalat(galat);
		pesanGalat = galat.message;
		if (keadaan === 'terkunci') terkunciSampai = Date.now() + detikTerkunci * milidetik;
		if (keadaan === 'salah') kode = '';
	}

	async function kirimKode(event?: SubmitEvent) {
		event?.preventDefault();
		isMengirim = true;
		const hasil = await tangkap(kirimKodeMasuk(fetch, email));
		isMengirim = false;
		if (!hasil.ok) return tanganiGalat(hasil.galat);
		tahap = 'kode';
		kode = '';
		keadaan = 'siap';
		kirimUlangPada = Date.now() + detikKirimUlang * milidetik;
		await tick();
		document.getElementById(idKode)?.focus();
	}

	async function verifikasi(kodeLengkap: string) {
		isMemeriksa = true;
		const hasil = await tangkap(verifikasiKodeMasuk(fetch, email, kodeLengkap));
		isMemeriksa = false;
		if (!hasil.ok) return tanganiGalat(hasil.galat);
		simpanSesi(hasil.data);
		await goto(tujuanSetelahMasuk(page.url.searchParams.get('lanjut'), bahasaBawaan));
	}

	function gantiEmail() {
		tahap = 'email';
		keadaan = 'siap';
	}

	const manfaat = [
		{ ikon: UserCheck, teks: t('auth.benefit.witness') },
		{ ikon: Bell, teks: t('auth.benefit.watch') },
		{ ikon: Scale, teks: t('auth.benefit.reply') }
	];
</script>

<svelte:head>
	<title>{t('auth.title')} · {t('brand.name')}</title>
</svelte:head>

<a
	href="#konten"
	class="sr-only z-tooltip rounded-md bg-primary px-4 py-3 font-bold text-primary-on focus:not-sr-only focus:fixed focus:left-4 focus:top-4"
	>{t('nav.skip')}</a
>

<div class="min-h-dvh bg-page lg:grid lg:grid-cols-[1.05fr_1fr]">
	<aside
		class="relative hidden flex-col justify-between gap-10 overflow-hidden border-r border-line bg-surface p-10 lg:flex xl:p-14"
	>
		<Logo />
		<div class="flex flex-col gap-8">
			<IlustrasiMasuk class="mx-auto w-full max-w-md" />
			<div>
				<p class="t-h1">{t('auth.noAccountNeeded')}</p>
				<p class="mt-3 max-w-prose text-muted">{t('auth.panelLead')}</p>
			</div>
			<ul class="flex flex-col gap-3">
				{#each manfaat as item (item.teks)}
					{@const Ikon = item.ikon}
					<li class="flex items-center gap-3 font-semibold">
						<span class="icon-tile h-9 w-9"
							><Ikon size={18} strokeWidth={1.75} aria-hidden="true" /></span
						>
						{item.teks}
					</li>
				{/each}
			</ul>
		</div>
		<p class="text-small text-muted">{t('footer.independent')}</p>
	</aside>

	<div class="flex min-h-dvh flex-col">
		<header class="wrap flex h-16 items-center justify-between gap-4 lg:h-18">
			<span class="lg:invisible"><Logo /></span>
			<a href="/{bahasaBawaan}" class="btn btn-ghost btn-sm">
				<ArrowLeft size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('nav.home')}
			</a>
		</header>

		<main
			id="konten"
			tabindex="-1"
			class="wrap flex flex-1 flex-col items-center justify-center pb-16 outline-none"
		>
			<IlustrasiMasuk class="mb-2 w-56 max-w-full md:w-64 lg:hidden" />
			<div class="card w-full max-w-form p-6 shadow-2 md:p-8">
				{#if tahap === 'email'}
					<span class="icon-tile"><Mail size={20} strokeWidth={1.75} aria-hidden="true" /></span>
					<h1 class="t-h1 mt-4">{t('auth.title')}</h1>
					<p class="mt-2 text-muted">{t('auth.lead')}</p>
					<form class="mt-6 flex flex-col gap-4" onsubmit={kirimKode}>
						<div class="flex flex-col gap-2">
							<label for="email-masuk" class="font-semibold">{t('auth.email')}</label>
							<div class="relative">
								<Mail
									size={20}
									strokeWidth={1.75}
									aria-hidden="true"
									class="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-muted"
								/>
								<input
									id="email-masuk"
									type="email"
									required
									autocomplete="email"
									inputmode="email"
									placeholder={t('auth.emailPlaceholder')}
									bind:value={email}
									aria-invalid={keadaan === 'galat' ? 'true' : undefined}
									aria-describedby={keadaan === 'siap' ? undefined : 'pesan-email'}
									class="field pl-12"
								/>
							</div>
						</div>
						{#if keadaan === 'galat'}
							<p id="pesan-email" role="alert" class="text-small text-danger">{pesanGalat}</p>
						{:else if isTerkunci}
							<p
								id="pesan-email"
								role="alert"
								class="flex gap-2 rounded-md bg-danger-bg p-3 text-small text-danger"
							>
								<Lock size={20} strokeWidth={1.75} aria-hidden="true" />{t('auth.locked')}
							</p>
						{/if}
						<button
							type="submit"
							class="btn btn-primary btn-lg btn-block"
							disabled={isMengirim || isTerkunci}
							aria-busy={isMengirim}
						>
							{#if isMengirim}
								<LoaderCircle
									size={20}
									strokeWidth={1.75}
									aria-hidden="true"
									class="animate-spin"
								/>
								{t('auth.sending')}
							{:else if isTerkunci}
								<Lock size={20} strokeWidth={1.75} aria-hidden="true" />
								{t('auth.lockedButton', { time: formatHitungMundur(sisaTerkunci) })}
							{:else}
								{t('auth.sendCode')}
								<ArrowRight size={20} strokeWidth={1.75} aria-hidden="true" />
							{/if}
						</button>
					</form>
				{:else}
					<span class="icon-tile"
						><ShieldCheck size={20} strokeWidth={1.75} aria-hidden="true" /></span
					>
					<h1 class="t-h1 mt-4">{t('auth.codeTitle')}</h1>
					<p class="mt-2 text-muted">{t('auth.codeSent', { email })}</p>
					<div class="mt-6 flex flex-col gap-3">
						<label for={idKode} class="font-semibold">{t('auth.codeLabel')}</label>
						<OtpInput
							id={idKode}
							bind:value={kode}
							error={galatKode}
							locked={isTerkunci || isMemeriksa}
							describedby={keadaan === 'siap' ? undefined : idPesanKode}
							oncomplete={verifikasi}
						/>
						<p class="flex items-center gap-2 text-small text-muted" aria-live="polite">
							{#if isMemeriksa}
								<LoaderCircle
									size={16}
									strokeWidth={1.75}
									aria-hidden="true"
									class="animate-spin"
								/>
								{t('auth.verifying')}
							{/if}
						</p>
						{#if keadaan === 'salah'}
							<p id={idPesanKode} role="alert" class="text-small text-danger">
								{t('auth.wrongCode')}
							</p>
						{:else if keadaan === 'kedaluwarsa'}
							<p id={idPesanKode} role="alert" class="text-small text-danger">
								{t('auth.expired')}
							</p>
						{:else if isTerkunci}
							<p
								id={idPesanKode}
								role="alert"
								class="flex gap-2 rounded-md bg-danger-bg p-3 text-small text-danger"
							>
								<Lock size={20} strokeWidth={1.75} aria-hidden="true" />{t('auth.locked')}
							</p>
						{:else if keadaan === 'galat'}
							<p id={idPesanKode} role="alert" class="text-small text-danger">{pesanGalat}</p>
						{/if}
					</div>
					<div class="mt-4 flex flex-wrap items-center justify-between gap-2">
						{#if isTerkunci}
							<button type="button" class="btn btn-text" disabled>
								<Lock size={18} strokeWidth={1.75} aria-hidden="true" />
								{t('auth.lockedButton', { time: formatHitungMundur(sisaTerkunci) })}
							</button>
						{:else if sisaKirimUlang > 0 && keadaan !== 'kedaluwarsa'}
							<button type="button" class="btn btn-text" disabled>
								<RefreshCw size={18} strokeWidth={1.75} aria-hidden="true" />
								{t('auth.resendIn', { time: formatHitungMundur(sisaKirimUlang) })}
							</button>
						{:else}
							<button
								type="button"
								class="btn btn-text"
								disabled={isMengirim}
								onclick={() => kirimKode()}
							>
								<RefreshCw size={18} strokeWidth={1.75} aria-hidden="true" />
								{isMengirim ? t('auth.sending') : t('auth.newCode')}
							</button>
						{/if}
						<button type="button" class="btn btn-text" onclick={gantiEmail}>
							<PenLine size={18} strokeWidth={1.75} aria-hidden="true" />
							{t('auth.changeEmail')}
						</button>
					</div>
				{/if}
				<p class="mt-6 flex items-start gap-2 border-t border-line pt-4 text-caption text-muted">
					<ShieldCheck size={16} strokeWidth={1.75} aria-hidden="true" class="shrink-0" />
					{t('auth.privacy')}
				</p>
			</div>
		</main>
	</div>
</div>
