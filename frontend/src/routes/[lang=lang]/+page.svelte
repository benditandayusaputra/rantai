<script lang="ts">
	import { GraduationCap, Landmark, Wallet } from '@lucide/svelte';
	import { goto } from '$app/navigation';
	import { cariWilayah } from '$lib/api/klien';
	import type { WilayahDitemukan } from '$lib/api/tipe';
	import RegionSearch from '$lib/components/RegionSearch.svelte';
	import UseLocationButton from '$lib/components/UseLocationButton.svelte';
	import { cariLokasiku, type StatusLokasi } from '$lib/features/wilayah/lokasi';
	import { contohWilayah, tujuanWilayah } from '$lib/features/wilayah/wilayah';
	import { t } from '$lib/i18n';

	let kata = $state('');
	let statusLokasi = $state<StatusLokasi>('idle');
	let akurasiM = $state(0);

	const lensa = [
		{ judul: t('lens.services'), tanya: t('lens.services.q'), ikon: GraduationCap },
		{ judul: t('lens.money'), tanya: t('lens.money.q'), ikon: Wallet },
		{ judul: t('lens.accountability'), tanya: t('lens.accountability.q'), ikon: Landmark }
	];

	function buka(wilayah: WilayahDitemukan) {
		goto(tujuanWilayah(wilayah));
	}

	async function pakaiLokasiku() {
		statusLokasi = 'locating';
		const hasil = await cariLokasiku(navigator.geolocation);
		statusLokasi = hasil.status;
		if ('akurasiM' in hasil) akurasiM = hasil.akurasiM;
		if (hasil.status === 'found') buka(hasil.wilayah);
	}
</script>

<svelte:head>
	<title>{t('brand.name')} · {t('home.title')}</title>
</svelte:head>

<section class="wrap pt-10 lg:pt-16">
	<div class="max-w-3xl">
		<p class="t-over">{t('home.overline')}</p>
		<h1 class="t-display mt-3">{t('home.title')}</h1>
		<p class="t-lead mt-4">{t('home.lead')}</p>
		<div class="mt-8 flex flex-col gap-4">
			<RegionSearch
				bind:value={kata}
				size="lg"
				cari={(q) => cariWilayah(fetch, q)}
				onselect={buka}
				onuselocation={pakaiLokasiku}
			/>
			<UseLocationButton state={statusLokasi} {akurasiM} onclick={pakaiLokasiku} />
			<p class="text-caption text-muted">{t('location.privacy')}</p>
			<div class="flex flex-wrap items-center gap-2">
				<span class="text-small text-muted">{t('home.try')}</span>
				{#each contohWilayah as nama (nama)}
					<button
						type="button"
						class="inline-flex min-h-touch items-center rounded-pill border-ctrl border-line-strong px-4 text-small font-semibold hover:border-fg hover:bg-sunken"
						onclick={() => (kata = nama)}>{nama}</button
					>
				{/each}
			</div>
		</div>
	</div>
</section>

<section class="wrap mt-14" aria-labelledby="judul-lensa">
	<h2 id="judul-lensa" class="t-h2">{t('home.lensesTitle')}</h2>
	<ul class="mt-6 grid gap-3 md:grid-cols-2 md:gap-4 lg:grid-cols-3 lg:gap-5">
		{#each lensa as item (item.judul)}
			{@const Ikon = item.ikon}
			<li class="rounded-lg border border-line bg-surface p-4 md:p-5">
				<Ikon size={24} strokeWidth={1.75} aria-hidden="true" class="text-primary" />
				<h3 class="t-h3 mt-3">{item.judul}</h3>
				<p class="mt-2 text-small text-muted">{item.tanya}</p>
			</li>
		{/each}
	</ul>
</section>
