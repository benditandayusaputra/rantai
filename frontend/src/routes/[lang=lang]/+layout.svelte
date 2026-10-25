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
	class="sr-only z-tooltip rounded-md bg-primary px-4 py-3 font-bold text-primary-on focus:not-sr-only focus:fixed focus:left-4 focus:top-4"
	>{t('nav.skip')}</a
>

<div
	class="flex min-h-dvh flex-col max-lg:pb-[calc(theme(spacing.16)+env(safe-area-inset-bottom))]"
>
	<header class="glass sticky top-0 z-sticky border-b border-line">
		<div class="wrap flex h-16 items-center justify-between gap-4 lg:h-18">
			<Logo />
			<nav aria-label={t('nav.main')} class="hidden lg:block">
				<ul class="flex items-center gap-1">
					{#each tujuanUtama as tujuan (tujuan.kunci)}
						{@const Ikon = tujuan.ikon}
						<li>
							<a
								href="/{bahasaAktif}{tujuan.path}"
								aria-current={aktif === tujuan.kunci ? 'page' : undefined}
								class="inline-flex min-h-touch items-center gap-2 rounded-pill px-3.5 text-small font-semibold text-muted transition-colors duration-fast ease-standard hover:bg-sunken hover:text-fg aria-[current=page]:bg-secondary aria-[current=page]:text-secondary-on"
							>
								<Ikon size={18} strokeWidth={1.75} aria-hidden="true" />
								{t(tujuan.label)}
							</a>
						</li>
					{/each}
				</ul>
			</nav>
			<div class="flex items-center gap-2">
				<PengalihTema class="max-md:hidden" />
				<PengalihBahasa />
				{#if sesi}
					{#if isAkunPanel(sesi.pengguna.peran)}
						<a href="/{bahasaAktif}/panel" class="btn btn-ghost btn-sm hidden lg:inline-flex">
							<LayoutDashboard size={18} strokeWidth={1.75} aria-hidden="true" />
							{t('nav.panel')}
						</a>
					{/if}
					<span
						title={t('auth.account', { email: sesi.pengguna.email })}
						class="inline-flex h-10 w-10 items-center justify-center rounded-pill bg-primary text-small font-bold text-primary-on"
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
					<a href={tautanMasuk} class="btn btn-primary btn-sm">
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
					class="group flex h-full flex-col items-center justify-center gap-1 text-caption text-muted aria-[current=page]:font-bold aria-[current=page]:text-fg"
				>
					<span
						class="inline-flex h-7.5 w-14 items-center justify-center rounded-pill transition-colors duration-fast ease-standard group-aria-[current=page]:bg-secondary group-aria-[current=page]:text-secondary-on"
					>
						<Ikon size={20} strokeWidth={1.75} aria-hidden="true" />
					</span>
					<span class="truncate">{t(tujuan.labelPonsel ?? tujuan.label)}</span>
				</a>
			</li>
		{/each}
	</ul>
</nav>
