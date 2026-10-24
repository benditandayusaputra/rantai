<script lang="ts">
	import { Lock, X } from '@lucide/svelte';
	import { tick } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { kirimKodeMasuk, verifikasiKodeMasuk } from '$lib/api/auth';
	import { tangkap, type GalatAPI } from '$lib/api/klien';
	import Logo from '$lib/components/Logo.svelte';
	import LogoMark from '$lib/components/LogoMark.svelte';
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

	const kelasTombolUtama =
		'inline-flex min-h-touch w-full items-center justify-center gap-2 rounded-md bg-primary px-4 font-semibold text-primary-on hover:bg-primary-hover disabled:bg-dis-bg disabled:text-dis-fg';
	const kelasTombolTeks =
		'inline-flex min-h-touch items-center rounded-md px-2 font-semibold text-primary underline-offset-4 hover:underline disabled:text-dis-fg disabled:no-underline';
</script>

<svelte:head>
	<title>{t('auth.title')} · {t('brand.name')}</title>
</svelte:head>

<a
	href="#konten"
	class="sr-only z-tooltip rounded-md bg-primary px-4 py-3 font-bold text-primary-on focus:not-sr-only focus:fixed focus:left-4 focus:top-4"
	>{t('nav.skip')}</a
>

<div class="min-h-dvh lg:grid lg:grid-cols-2">
	<aside class="hidden flex-col justify-between gap-12 bg-surface p-12 lg:flex">
		<Logo />
		<div>
			<LogoMark class="h-24 w-24" />
			<p class="t-display mt-8">{t('auth.noAccountNeeded')}</p>
			<p class="mt-4 max-w-prose text-muted">{t('auth.panelLead')}</p>
		</div>
		<p class="text-small text-muted">{t('footer.independent')}</p>
	</aside>

	<div class="flex min-h-dvh flex-col">
		<header class="wrap flex h-16 items-center justify-between gap-4 lg:h-18">
			<span class="lg:invisible"><Logo /></span>
			<a
				href="/{bahasaBawaan}"
				aria-label={t('action.close')}
				class="inline-flex h-11 w-11 items-center justify-center rounded-md hover:bg-sunken"
			>
				<X size={24} strokeWidth={1.75} aria-hidden="true" />
			</a>
		</header>

		<main id="konten" tabindex="-1" class="wrap flex flex-1 items-center justify-center pb-16">
			<div class="w-full max-w-form">
				{#if tahap === 'email'}
					<h1 class="t-h1">{t('auth.title')}</h1>
					<p class="mt-2 text-muted">{t('auth.lead')}</p>
					<form class="mt-8 flex flex-col gap-4" onsubmit={kirimKode}>
						<div class="flex flex-col gap-2">
							<label for="email-masuk" class="font-semibold">{t('auth.email')}</label>
							<input
								id="email-masuk"
								type="email"
								required
								autocomplete="email"
								bind:value={email}
								aria-invalid={keadaan === 'galat' ? 'true' : undefined}
								aria-describedby={keadaan === 'siap' ? undefined : 'pesan-email'}
								class="h-12 rounded-md border-ctrl border-line-strong bg-raised px-4 text-body hover:border-fg aria-[invalid=true]:border-danger"
							/>
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
						<button type="submit" class={kelasTombolUtama} disabled={isMengirim || isTerkunci}>
							{#if isMengirim}
								{t('auth.sending')}
							{:else if isTerkunci}
								{t('auth.lockedButton', { time: formatHitungMundur(sisaTerkunci) })}
							{:else}
								{t('auth.sendCode')}
							{/if}
						</button>
					</form>
				{:else}
					<h1 class="t-h1">{t('auth.codeTitle')}</h1>
					<p class="mt-2 text-muted">{t('auth.codeSent', { email })}</p>
					<div class="mt-8 flex flex-col gap-3">
						<label for={idKode} class="font-semibold">{t('auth.codeLabel')}</label>
						<OtpInput
							id={idKode}
							bind:value={kode}
							error={galatKode}
							locked={isTerkunci || isMemeriksa}
							describedby={keadaan === 'siap' ? undefined : idPesanKode}
							oncomplete={verifikasi}
						/>
						<p class="sr-only" aria-live="polite">{isMemeriksa ? t('auth.verifying') : ''}</p>
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
					<div class="mt-6 flex flex-wrap items-center justify-between gap-2">
						{#if isTerkunci}
							<button type="button" class={kelasTombolTeks} disabled>
								{t('auth.lockedButton', { time: formatHitungMundur(sisaTerkunci) })}
							</button>
						{:else if sisaKirimUlang > 0 && keadaan !== 'kedaluwarsa'}
							<button type="button" class={kelasTombolTeks} disabled>
								{t('auth.resendIn', { time: formatHitungMundur(sisaKirimUlang) })}
							</button>
						{:else}
							<button
								type="button"
								class={kelasTombolTeks}
								disabled={isMengirim}
								onclick={() => kirimKode()}
							>
								{isMengirim ? t('auth.sending') : t('auth.newCode')}
							</button>
						{/if}
						<button type="button" class={kelasTombolTeks} onclick={gantiEmail}>
							{t('auth.changeEmail')}
						</button>
					</div>
				{/if}
			</div>
		</main>
	</div>
</div>
