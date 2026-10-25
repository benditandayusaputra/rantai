<script lang="ts">
	import { LayoutDashboard, LogIn, LogOut } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { afterNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import Footer from '$lib/components/Footer.svelte';
	import Logo from '$lib/components/Logo.svelte';
	import PengalihBahasa from '$lib/components/PengalihBahasa.svelte';
	import PengalihTema from '$lib/components/PengalihTema.svelte';
	import { inisial, keluar, pulihkanSesi, sesiAktif } from '$lib/features/akun/sesi.svelte';
	import { isAkunPanel } from '$lib/features/panel/panel';
	import { tujuanAktif, tujuanPonsel, tujuanUtama } from '$lib/features/navigasi/navigasi';
	import { bahasaAktif, t } from '$lib/i18n';

	let { children } = $props();

	const aktif = $derived(tujuanAktif(page.url.pathname, bahasaAktif));
	const sesi = $derived(sesiAktif());
	const tautanMasuk = $derived(
		`/${bahasaAktif}/masuk?lanjut=${encodeURIComponent(page.url.pathname + page.url.search)}`
	);

	onMount(() => {
		void pulihkanSesi(fetch);
	});

	afterNavigate(({ from }) => {
		if (!from) return;
		const judul = document.querySelector<HTMLElement>('main h1');
		judul?.setAttribute('tabindex', '-1');
		judul?.focus();
	});
</script>

<svelte:head>
	<meta name="description" content={t('home.lead')} />
</svelte:head>

<a
	href="#konten"
	class="sr-only z-tooltip rounded-md bg-surface px-4 py-3 font-semibold text-fg shadow-2 focus:not-sr-only focus:fixed focus:left-4 focus:top-4"
	>{t('nav.skip')}</a
>

<div
	class="flex min-h-dvh flex-col max-lg:pb-[calc(theme(spacing.16)+env(safe-area-inset-bottom))]"
>
	<header class="glass sticky top-0 z-sticky border-b border-line">
		<div class="wrap flex h-16 items-center gap-6 md:h-18 lg:h-kepala xl:gap-10">
			<Logo />
			<nav aria-label={t('nav.main')} class="hidden self-stretch lg:mx-auto lg:block">
				<ul class="flex h-full items-stretch gap-6 xl:gap-7.5">
					{#each tujuanUtama as tujuan (tujuan.kunci)}
						<li class="flex">
							<a
								href="/{bahasaAktif}{tujuan.path}"
								aria-current={aktif === tujuan.kunci ? 'page' : undefined}
								class="relative inline-flex items-center whitespace-nowrap text-small font-medium text-fg transition-colors duration-fast ease-standard after:absolute after:inset-x-0 after:bottom-0 after:h-0.75 after:bg-transparent hover:text-primary aria-[current=page]:font-semibold aria-[current=page]:after:bg-primary"
							>
								{t(tujuan.label)}
							</a>
						</li>
					{/each}
				</ul>
			</nav>
			<div class="ml-auto flex items-center gap-2 lg:ml-0">
				<PengalihTema class="max-md:hidden" />
				<PengalihBahasa />
				<span class="mx-1 hidden h-6 w-px bg-line md:block" aria-hidden="true"></span>
				{#if sesi}
					{#if isAkunPanel(sesi.pengguna.peran)}
						<a href="/{bahasaAktif}/panel" class="btn btn-ghost btn-sm hidden lg:inline-flex">
							<LayoutDashboard size={18} strokeWidth={1.75} aria-hidden="true" />
							{t('nav.panel')}
						</a>
					{/if}
					<span
						title={t('auth.account', { email: sesi.pengguna.email })}
						class="inline-flex h-9 w-9 items-center justify-center rounded-pill border border-line-strong bg-sunken text-small font-semibold text-fg"
					>
						<span aria-hidden="true">{inisial(sesi.pengguna)}</span>
						<span class="sr-only">{t('auth.account', { email: sesi.pengguna.email })}</span>
					</span>
					<button
						type="button"
						class="btn btn-ghost btn-icon"
						aria-label={t('action.signOut')}
						title={t('action.signOut')}
						onclick={() => keluar(fetch)}
					>
						<LogOut size={20} strokeWidth={1.75} aria-hidden="true" />
					</button>
				{:else}
					<a href={tautanMasuk} class="btn btn-outline btn-sm">
						<LogIn size={18} strokeWidth={1.75} aria-hidden="true" />
						{t('action.signIn')}
					</a>
				{/if}
			</div>
		</div>
	</header>

	<main id="konten" tabindex="-1" class="flex-1 outline-none">
		{@render children()}
	</main>

	<Footer />
</div>

<nav
	aria-label={t('nav.main')}
	class="glass fixed inset-x-0 bottom-0 z-sticky border-t border-line pb-[env(safe-area-inset-bottom)] lg:hidden"
>
	<ul class="flex h-16">
		{#each tujuanPonsel as tujuan (tujuan.kunci)}
			{@const Ikon = tujuan.ikon}
			<li class="min-w-0 flex-1">
				<a
					href="/{bahasaAktif}{tujuan.path}"
					aria-current={aktif === tujuan.kunci ? 'page' : undefined}
					class="relative flex h-full flex-col items-center justify-center gap-1 text-caption text-muted before:absolute before:inset-x-3 before:top-0 before:h-0.75 before:bg-transparent aria-[current=page]:font-semibold aria-[current=page]:text-fg aria-[current=page]:before:bg-primary"
				>
					<Ikon size={20} strokeWidth={1.75} aria-hidden="true" />
					<span class="max-w-full truncate px-1">{t(tujuan.labelPonsel ?? tujuan.label)}</span>
				</a>
			</li>
		{/each}
	</ul>
</nav>
