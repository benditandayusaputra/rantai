<script lang="ts">
	import { ChartNoAxesCombined, GraduationCap, Landmark, Wallet } from '@lucide/svelte';
	import { invalidateAll } from '$app/navigation';
	import { navigating, page } from '$app/state';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import LensaLayanan from '$lib/features/dasbor/LensaLayanan.svelte';
	import LensaTanggungJawab from '$lib/features/dasbor/LensaTanggungJawab.svelte';
	import LensaUang from '$lib/features/dasbor/LensaUang.svelte';
	import {
		daftarLensa,
		tahunTersedia,
		tautanLensa,
		tautanTahun,
		type Lensa
	} from '$lib/features/dasbor/dasbor';
	import { t, type KunciPesan } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const infoLensa: Record<Lensa, { ikon: typeof Wallet; nama: KunciPesan; tanya: KunciPesan }> = {
		layanan: { ikon: GraduationCap, nama: 'lens.services', tanya: 'dash.services.q' },
		uang: { ikon: Wallet, nama: 'lens.money', tanya: 'dash.money.q' },
		'tanggung-jawab': {
			ikon: Landmark,
			nama: 'lens.accountability',
			tanya: 'dash.accountability.q'
		}
	};

	const hasil = $derived(data.hasil);
	const pilihanTahun = $derived(hasil.ok ? tahunTersedia(hasil.data) : null);
	const IkonLensa = $derived(infoLensa[data.lensa].ikon);
</script>

<svelte:head>
	<title>{t('dash.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="border-b border-line bg-surface">
	<div class="wrap pt-6 lg:pt-10">
		<div class="flex items-start gap-4">
			<span class="icon-tile hidden h-14 w-14 md:inline-flex">
				<ChartNoAxesCombined size={28} strokeWidth={1.75} aria-hidden="true" />
			</span>
			<div class="min-w-0">
				<h1 class="t-h1">{t('dash.title')}</h1>
				<p class="t-lead mt-2 max-w-prose">{t('dash.lead')}</p>
			</div>
		</div>
		<nav aria-label={t('dash.lensNav')} class="mt-4">
			<ul class="grid grid-cols-3 md:flex md:gap-2">
				{#each daftarLensa as lensa (lensa)}
					{@const Ikon = infoLensa[lensa].ikon}
					<li>
						<a
							href={tautanLensa(page.url, lensa)}
							data-sveltekit-noscroll
							aria-current={data.lensa === lensa ? 'page' : undefined}
							class="-mb-px flex h-15 flex-col items-center justify-center gap-1 border-b-link border-transparent px-2 text-small font-semibold text-muted transition-colors duration-fast ease-standard hover:text-fg md:h-14 md:flex-row md:gap-2 md:px-4 md:text-body aria-[current=page]:border-primary aria-[current=page]:text-primary"
						>
							<Ikon size={20} strokeWidth={1.75} aria-hidden="true" />
							{t(infoLensa[lensa].nama)}
						</a>
					</li>
				{/each}
			</ul>
		</nav>
	</div>
</div>

<section
	class="wrap mt-8 transition-opacity duration-base ease-standard lg:mt-10"
	class:opacity-60={navigating.to}
	aria-labelledby="judul-lensa"
	aria-busy={Boolean(navigating.to)}
>
	<div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
		<div class="flex items-start gap-4">
			<span class="icon-tile hidden md:inline-flex">
				<IkonLensa size={20} strokeWidth={1.75} aria-hidden="true" />
			</span>
			<div>
				<h2 id="judul-lensa" class="t-h2">{t(infoLensa[data.lensa].nama)}</h2>
				<p class="t-lead mt-1 max-w-prose">{t(infoLensa[data.lensa].tanya)}</p>
			</div>
		</div>
		{#if pilihanTahun && pilihanTahun.daftar.length > 1}
			<nav aria-label={t('dash.yearNav')}>
				<ul class="flex flex-wrap gap-2">
					{#each pilihanTahun.daftar as tahun (tahun)}
						<li>
							<a
								href={tautanTahun(page.url, tahun)}
								data-sveltekit-noscroll
								aria-current={tahun === pilihanTahun.aktif ? 'true' : undefined}
								class="chip num aria-[current=true]:border-primary aria-[current=true]:bg-secondary aria-[current=true]:text-secondary-on"
								>{tahun}</a
							>
						</li>
					{/each}
				</ul>
			</nav>
		{/if}
	</div>

	<div class="mt-6 lg:mt-8">
		{#if !hasil.ok}
			<ErrorState
				title={t('dash.failed')}
				description={hasil.galat.message}
				onretry={invalidateAll}
			/>
		{:else if hasil.data.lensa === 'layanan'}
			<LensaLayanan dasbor={hasil.data.dasbor} peta={hasil.data.peta} />
		{:else if hasil.data.lensa === 'uang'}
			<LensaUang dasbor={hasil.data.dasbor} peta={hasil.data.peta} />
		{:else}
			<LensaTanggungJawab dasbor={hasil.data.dasbor} />
		{/if}
	</div>
</section>
