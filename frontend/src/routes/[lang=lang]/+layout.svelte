<script lang="ts">
	import { afterNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import Footer from '$lib/components/Footer.svelte';
	import Logo from '$lib/components/Logo.svelte';
	import { tujuanAktif, tujuanUtama } from '$lib/features/navigasi/navigasi';
	import { bahasaBawaan, t } from '$lib/i18n';

	let { children } = $props();

	const aktif = $derived(tujuanAktif(page.url.pathname, bahasaBawaan));

	afterNavigate(({ from }) => {
		if (!from) return;
		const judul = document.querySelector<HTMLElement>('main h1');
		judul?.setAttribute('tabindex', '-1');
		judul?.focus();
	});
</script>

<a
	href="#konten"
	class="sr-only z-tooltip rounded-md bg-primary px-4 py-3 font-bold text-primary-on focus:not-sr-only focus:fixed focus:left-4 focus:top-4"
	>{t('nav.skip')}</a
>

<header class="sticky top-0 z-sticky border-b border-line bg-surface">
	<div class="wrap flex h-16 items-center justify-between gap-4 lg:h-18">
		<Logo />
		<nav aria-label={t('nav.main')} class="hidden lg:block">
			<ul class="flex items-center gap-2">
				{#each tujuanUtama as tujuan (tujuan.kunci)}
					<li>
						<a
							href="/{bahasaBawaan}{tujuan.path}"
							aria-current={aktif === tujuan.kunci ? 'page' : undefined}
							class="inline-flex min-h-touch items-center rounded-md px-3 font-semibold hover:bg-sunken aria-[current=page]:text-primary aria-[current=page]:underline aria-[current=page]:decoration-2 aria-[current=page]:underline-offset-8"
							>{t(tujuan.label)}</a
						>
					</li>
				{/each}
			</ul>
		</nav>
	</div>
</header>

<main id="konten" tabindex="-1" class="pb-24 outline-none lg:pb-0">
	{@render children()}
</main>

<Footer />

<nav
	aria-label={t('nav.main')}
	class="fixed inset-x-0 bottom-0 z-sticky border-t border-line bg-surface pb-[env(safe-area-inset-bottom)] lg:hidden"
>
	<ul class="flex h-16">
		{#each tujuanUtama as tujuan (tujuan.kunci)}
			{@const Ikon = tujuan.ikon}
			<li class="flex-1">
				<a
					href="/{bahasaBawaan}{tujuan.path}"
					aria-current={aktif === tujuan.kunci ? 'page' : undefined}
					class="group flex h-full flex-col items-center justify-center gap-1 text-caption aria-[current=page]:font-bold"
				>
					<span
						class="inline-flex h-7.5 w-14 items-center justify-center rounded-pill group-aria-[current=page]:bg-secondary"
					>
						<Ikon size={20} strokeWidth={1.75} aria-hidden="true" />
					</span>
					{t(tujuan.label)}
				</a>
			</li>
		{/each}
	</ul>
</nav>
