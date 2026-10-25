<script lang="ts">
	import {
		ArrowRight,
		Bell,
		Camera,
		FileLock,
		GraduationCap,
		Landmark,
		Map,
		MapPin,
		MessageCircle,
		School,
		Wallet
	} from '@lucide/svelte';
	import { goto } from '$app/navigation';
	import { cariWilayah } from '$lib/api/klien';
	import type { WilayahDitemukan } from '$lib/api/tipe';
	import CertaintyBadge from '$lib/components/CertaintyBadge.svelte';
	import IlustrasiKota from '$lib/components/IlustrasiKota.svelte';
	import RegionSearch from '$lib/components/RegionSearch.svelte';
	import UseLocationButton from '$lib/components/UseLocationButton.svelte';
	import { cariLokasiku, type StatusLokasi } from '$lib/features/wilayah/lokasi';
	import { contohWilayah, tujuanWilayah } from '$lib/features/wilayah/wilayah';
	import { bahasaAktif, t } from '$lib/i18n';

	let kata = $state('');
	let statusLokasi = $state<StatusLokasi>('idle');
	let akurasiM = $state(0);

	const lensa = [
		{ judul: t('lens.services'), tanya: t('lens.services.q'), ikon: GraduationCap },
		{ judul: t('lens.money'), tanya: t('lens.money.q'), ikon: Wallet },
		{ judul: t('lens.accountability'), tanya: t('lens.accountability.q'), ikon: Landmark }
	];
	const fitur = [
		{ path: '/peta', judul: t('nav.map'), isi: t('home.feature.map'), ikon: Map },
		{ path: '/tanya', judul: t('ask.title'), isi: t('home.feature.ask'), ikon: MessageCircle },
		{ path: '/cek-sekolah', judul: t('check.title'), isi: t('home.feature.check'), ikon: School },
		{ path: '/lapor', judul: t('report.flowTitle'), isi: t('home.feature.report'), ikon: Camera },
		{ path: '/pantau', judul: t('watch.title'), isi: t('home.feature.watch'), ikon: Bell },
		{
			path: '/data-belum-dibuka',
			judul: t('unopened.title'),
			isi: t('home.feature.unopened'),
			ikon: FileLock
		}
	];
	const labelKepastian = ['resmi', 'proksi', 'estimasi', 'simulasi'] as const;

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

<section class="relative overflow-hidden border-b border-line bg-surface">
	<div
		class="wrap grid items-center gap-10 pb-12 pt-10 md:pb-16 lg:grid-cols-[1.15fr_0.85fr] lg:gap-14 lg:pb-20 lg:pt-16"
	>
		<div>
			<p
				class="inline-flex items-center gap-2 rounded-pill bg-secondary px-3 py-1 text-overline uppercase text-secondary-on"
			>
				<MapPin size={14} strokeWidth={2} aria-hidden="true" />
				{t('home.overline')}
			</p>
			<h1 class="t-display mt-4">{t('home.title')}</h1>
			<p class="t-lead mt-4 max-w-prose">{t('home.lead')}</p>
			<div class="card mt-8 flex flex-col gap-4 p-4 shadow-2 md:p-5">
				<RegionSearch
					bind:value={kata}
					size="lg"
					cari={(q) => cariWilayah(fetch, q)}
					onselect={buka}
					onuselocation={pakaiLokasiku}
				/>
				<div class="flex flex-col gap-1">
					<UseLocationButton state={statusLokasi} {akurasiM} onclick={pakaiLokasiku} />
					<p class="text-caption text-muted">{t('location.privacy')}</p>
				</div>
				<div class="flex flex-wrap items-center gap-2 border-t border-line pt-4">
					<span class="text-small text-muted">{t('home.try')}</span>
					{#each contohWilayah as nama (nama)}
						<button type="button" class="chip" onclick={() => (kata = nama)}>
							<MapPin size={16} strokeWidth={1.75} aria-hidden="true" />
							{nama}
						</button>
					{/each}
				</div>
			</div>
		</div>
		<IlustrasiKota class="mx-auto hidden w-full max-w-lg md:block" />
	</div>
</section>

<section class="wrap mt-14 lg:mt-20" aria-labelledby="judul-lensa">
	<h2 id="judul-lensa" class="t-h2">{t('home.lensesTitle')}</h2>
	<ol class="mt-6 grid gap-4 md:grid-cols-3 lg:gap-5">
		{#each lensa as item, urutan (item.judul)}
			{@const Ikon = item.ikon}
			<li class="card relative flex flex-col gap-3 p-5 lg:p-6">
				<div class="flex items-center justify-between">
					<span class="icon-tile h-12 w-12"
						><Ikon size={24} strokeWidth={1.75} aria-hidden="true" /></span
					>
					<span class="num font-serif text-h1 text-line-strong" aria-hidden="true"
						>0{urutan + 1}</span
					>
				</div>
				<h3 class="t-h3">{item.judul}</h3>
				<p class="text-small text-muted">{item.tanya}</p>
			</li>
		{/each}
	</ol>
</section>

<section class="wrap mt-14 lg:mt-20" aria-labelledby="judul-fitur">
	<h2 id="judul-fitur" class="t-h2">{t('home.featuresTitle')}</h2>
	<ul class="mt-6 grid gap-4 md:grid-cols-2 lg:grid-cols-3 lg:gap-5">
		{#each fitur as item (item.path)}
			{@const Ikon = item.ikon}
			<li>
				<a
					href="/{bahasaAktif}{item.path}"
					class="card-link group flex h-full items-start gap-4 p-5"
				>
					<span class="icon-tile"><Ikon size={20} strokeWidth={1.75} aria-hidden="true" /></span>
					<span class="flex min-w-0 flex-1 flex-col gap-1">
						<span class="flex items-center justify-between gap-2 font-bold">
							{item.judul}
							<ArrowRight
								size={18}
								strokeWidth={1.75}
								aria-hidden="true"
								class="shrink-0 text-muted transition-transform duration-fast ease-standard group-hover:translate-x-1 group-hover:text-primary"
							/>
						</span>
						<span class="text-small text-muted">{item.isi}</span>
					</span>
				</a>
			</li>
		{/each}
	</ul>
</section>

<section class="wrap mt-14 lg:mt-20" aria-labelledby="judul-sumber">
	<div
		class="grid gap-6 rounded-xl bg-primary p-6 text-primary-on md:p-8 lg:grid-cols-[1fr_1.2fr] lg:items-center lg:p-10"
	>
		<div>
			<h2 id="judul-sumber" class="t-h2">{t('home.trustTitle')}</h2>
			<p class="mt-3 opacity-90">{t('home.trustLead')}</p>
			<a href="/{bahasaAktif}/tentang-data" class="btn mt-6 bg-raised text-fg hover:bg-surface">
				{t('action.readMethod')}
				<ArrowRight size={18} strokeWidth={1.75} aria-hidden="true" />
			</a>
		</div>
		<ul class="grid gap-3 md:grid-cols-2">
			{#each labelKepastian as level (level)}
				<li class="flex flex-col items-start gap-2 rounded-lg bg-raised p-4 text-fg">
					<CertaintyBadge {level} size="md" />
					<p class="text-small text-muted">{t(`cert.${level}.desc`)}</p>
				</li>
			{/each}
		</ul>
	</div>
</section>
