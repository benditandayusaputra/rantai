<script lang="ts">
	import {
		ArrowDown,
		ArrowRight,
		ExternalLink,
		Landmark,
		Map,
		MapPin,
		School,
		ShieldCheck,
		Wallet
	} from '@lucide/svelte';
	import { untrack } from 'svelte';
	import { goto } from '$app/navigation';
	import { cariWilayah } from '$lib/api/klien';
	import type { WilayahDitemukan } from '$lib/api/tipe';
	import RegionSearch from '$lib/components/RegionSearch.svelte';
	import UseLocationButton from '$lib/components/UseLocationButton.svelte';
	import { daftarSorotan } from '$lib/features/beranda/beranda';
	import { cariLokasiku, type StatusLokasi } from '$lib/features/wilayah/lokasi';
	import { sesiAktif } from '$lib/features/akun/sesi.svelte';
	import { berandaPengguna, ingatWilayah } from '$lib/features/wilayah/wilayahku.svelte';
	import { tujuanWilayah } from '$lib/features/wilayah/wilayah';
	import { bahasaAktif, t } from '$lib/i18n';

	let kata = $state('');
	let statusLokasi = $state<StatusLokasi>('idle');
	let akurasiM = $state(0);

	const sumberFoto = 'https://unsplash.com/photos/1555899434-94d1368aa7af';
	const penjelasan = [
		{ ikon: School, judul: 'home.explain.services', isi: 'home.explain.servicesBody' },
		{ ikon: Wallet, judul: 'home.explain.projects', isi: 'home.explain.projectsBody' },
		{ ikon: Landmark, judul: 'home.explain.people', isi: 'home.explain.peopleBody' }
	] as const;
	const pertanyaanUmum = [
		{ tanya: 'home.faq.q1', jawab: 'home.faq.a1' },
		{ tanya: 'home.faq.q2', jawab: 'home.faq.a2' },
		{ tanya: 'home.faq.q3', jawab: 'home.faq.a3' },
		{ tanya: 'home.faq.q4', jawab: 'home.faq.a4' },
		{ tanya: 'home.faq.q5', jawab: 'home.faq.a5' },
		{ tanya: 'home.faq.q6', jawab: 'home.faq.a6' }
	] as const;

	const jangkarCari = '#cari-wilayah';

	$effect(() => {
		if (!sesiAktif()) return;
		const beranda = untrack(() => berandaPengguna(bahasaAktif));
		if (beranda !== `/${bahasaAktif}` && location.hash !== jangkarCari)
			void goto(beranda, { replaceState: true });
	});

	function buka(wilayah: WilayahDitemukan) {
		ingatWilayah(wilayah);
		void goto(tujuanWilayah(wilayah));
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
	<title>{t('brand.name')} · {t('home.heroTitle')}</title>
</svelte:head>

<div class="wrap pt-6 lg:pt-8">
	<section
		aria-labelledby="judul-beranda"
		class="relative isolate overflow-hidden rounded-xl bg-foto text-foto-on"
	>
		<picture class="block">
			<source media="(max-width: 767px)" srcset="/foto/bundaran-hi-800.webp" />
			<img
				src="/foto/bundaran-hi.webp"
				alt={t('home.photoAlt')}
				width="1600"
				height="1067"
				fetchpriority="high"
				class="block h-auto w-full object-cover md:h-80 lg:h-96"
				data-testid="hero-image"
			/>
		</picture>
		<div
			class="bayang-foto pointer-events-none absolute inset-0 hidden md:block"
			aria-hidden="true"
		></div>
		<div
			class="relative flex flex-col items-start px-5 py-7 md:absolute md:inset-0 md:justify-center md:px-8 lg:px-12"
		>
			<p class="text-small font-semibold">{t('home.kicker')}</p>
			<h1 id="judul-beranda" class="t-display mt-3 max-w-xl">{t('home.title')}</h1>
			<p class="mt-4 max-w-lg text-body">{t('home.lead')}</p>
			<a href="#cari-wilayah" class="btn btn-lg mt-6 bg-foto-on text-foto hover:opacity-90"
				>{t('home.startSearch')}<ArrowDown size={18} aria-hidden="true" /></a
			>
			<a
				href={sumberFoto}
				target="_blank"
				rel="noopener"
				class="mt-4 inline-flex min-h-touch items-center gap-2 text-caption underline-offset-4 hover:underline md:absolute md:bottom-3 md:right-5"
				aria-label="{t('home.photoCredit')} {t('a11y.opensNewTab')}"
				>{t('home.photoCaption')}<ExternalLink size={14} aria-hidden="true" /></a
			>
		</div>
	</section>

	<section
		id="cari-wilayah"
		class="scroll-mt-28 border-b border-line py-8 lg:py-10"
		aria-labelledby="judul-cari"
	>
		<div class="grid gap-10 lg:grid-cols-2 lg:gap-12">
			<div>
				<p class="t-over text-primary">{t('home.startStep')}</p>
				<h2 id="judul-cari" class="t-h2 mt-2">{t('home.searchTitle')}</h2>
				<p class="mt-2 text-muted">{t('home.searchLead')}</p>
				<div class="mt-5">
					<RegionSearch
						bind:value={kata}
						cari={(q) => cariWilayah(fetch, q)}
						onselect={buka}
						onuselocation={pakaiLokasiku}
					/>
				</div>
				<div class="mt-4">
					<UseLocationButton state={statusLokasi} {akurasiM} onclick={pakaiLokasiku} />
				</div>
				<p class="mt-3 flex items-start gap-2 text-caption text-muted">
					<ShieldCheck size={16} class="mt-0.5 shrink-0" aria-hidden="true" />{t(
						'location.privacy'
					)}
				</p>
				<p class="mt-2 text-caption text-muted">{t('home.noAccount')}</p>
			</div>
			<div>
				<h2 class="t-h3">{t('home.areaPicker')}</h2>
				<p class="mt-2 text-small text-muted">{t('home.exampleLabel')}</p>
				<ul class="mt-4 grid gap-3 xs:grid-cols-2">
					{#each daftarSorotan as wilayah (wilayah.kode)}
						<li>
							<a
								href="/{bahasaAktif}/kecamatan/{wilayah.kode}"
								class="card-link group flex h-full items-center gap-3 p-4"
							>
								<span class="icon-tile h-9 w-9"><MapPin size={18} aria-hidden="true" /></span>
								<span class="flex min-w-0 flex-col">
									<span class="font-semibold">{wilayah.nama}</span>
									<span class="text-caption text-muted">{wilayah.kota}</span>
								</span>
								<ArrowRight
									size={18}
									aria-hidden="true"
									class="ml-auto shrink-0 text-muted group-hover:text-fg"
								/>
							</a>
						</li>
					{/each}
				</ul>
				<a href="/{bahasaAktif}/peta" class="btn btn-outline btn-sm mt-4"
					><Map size={16} aria-hidden="true" />{t('home.map.open')}</a
				>
			</div>
		</div>
	</section>

	<section class="py-8 lg:py-10" aria-labelledby="judul-fungsi">
		<h2 id="judul-fungsi" class="t-h2">{t('home.explain.title')}</h2>
		<ol class="mt-6 grid gap-6 md:grid-cols-3 md:gap-8">
			{#each penjelasan as item, i (item.judul)}
				{@const Ikon = item.ikon}
				<li class="border-t border-line pt-5">
					<p class="flex items-center gap-2 text-caption font-semibold text-muted">
						<Ikon size={18} aria-hidden="true" /><span aria-hidden="true">0{i + 1}</span>
					</p>
					<h3 class="t-h3 mt-3">{t(item.judul)}</h3>
					<p class="mt-2 text-small text-muted">{t(item.isi)}</p>
				</li>
			{/each}
		</ol>
		<p class="mt-6 text-small text-muted">
			{t('home.explain.scope')}
			<a
				href="/{bahasaAktif}/tentang-data"
				class="inline-flex min-h-touch items-center font-semibold underline underline-offset-4"
				>{t('home.trustLink')}</a
			>
			·
			<a
				href="/{bahasaAktif}/belajar"
				class="inline-flex min-h-touch items-center font-semibold underline underline-offset-4"
				>{t('home.learnLink')}</a
			>
		</p>
	</section>

	<section class="border-t border-line py-8 lg:py-10" aria-labelledby="judul-faq">
		<h2 id="judul-faq" class="t-h2">{t('home.faq.title')}</h2>
		<div class="mt-4 divide-y divide-line border-y border-line">
			{#each pertanyaanUmum as item (item.tanya)}
				<details>
					<summary class="min-h-touch cursor-pointer py-4 font-semibold">{t(item.tanya)}</summary>
					<p class="max-w-prose pb-5 text-muted">{t(item.jawab)}</p>
				</details>
			{/each}
		</div>
	</section>
</div>
