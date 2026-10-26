<script lang="ts">
	import {
		ArrowRight,
		Bell,
		Camera,
		ExternalLink,
		FileLock,
		FileText,
		Landmark,
		MapPin,
		MessageCircle,
		School,
		Wallet
	} from '@lucide/svelte';
	import { goto, invalidateAll, replaceState } from '$app/navigation';
	import { navigating, page } from '$app/state';
	import { cariWilayah } from '$lib/api/klien';
	import type { WilayahDitemukan } from '$lib/api/tipe';
	import CertaintyBadge from '$lib/components/CertaintyBadge.svelte';
	import ChoiceChips from '$lib/components/ChoiceChips.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import FotoWakil from '$lib/components/FotoWakil.svelte';
	import GapClassChip from '$lib/components/GapClassChip.svelte';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import RegionSearch from '$lib/components/RegionSearch.svelte';
	import UseLocationButton from '$lib/components/UseLocationButton.svelte';
	import {
		batasPaketBeranda,
		daftarSorotan,
		indikatorPendidikan,
		persenTerhadap,
		tautanSorotan,
		wakilPengawas
	} from '$lib/features/beranda/beranda';
	import { kotakKelasPeta, namaWilayahPendek, urutanKelas } from '$lib/features/dasbor/dasbor';
	import { kunciKelas, tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import PetaBerlapis from '$lib/features/peta/PetaBerlapis.svelte';
	import {
		bacaLapisan,
		bacaPetaDasar,
		tautanPetaPenuh,
		urlDenganStatusPeta
	} from '$lib/features/peta/lapisan';
	import { pilihanJenjang, teksJenjang } from '$lib/features/pendidikan/jenjang';
	import { cariLokasiku, type StatusLokasi } from '$lib/features/wilayah/lokasi';
	import { tujuanWilayah } from '$lib/features/wilayah/wilayah';
	import { formatAngka } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	let kata = $state('');
	let statusLokasi = $state<StatusLokasi>('idle');
	let akurasiM = $state(0);
	let lapisan = $state(bacaLapisan(page.url.searchParams.get('lapisan')));
	let petaDasar = $state(bacaPetaDasar(page.url.searchParams.get('dasar')));

	const sumberFoto = 'https://unsplash.com/photos/1555899434-94d1368aa7af';
	const labelKepastian = ['resmi', 'proksi', 'estimasi', 'simulasi'] as const;
	const langkahLain = [
		{ path: '/cek-sekolah', judul: t('check.title'), isi: t('home.feature.check'), ikon: School },
		{ path: '/lapor', judul: t('report.flowTitle'), isi: t('home.feature.report'), ikon: Camera },
		{ path: '/pantau', judul: t('watch.title'), isi: t('home.feature.watch'), ikon: Bell }
	];

	const sorotan = $derived(data.sorotan);
	const halamanWilayah = $derived(`/${bahasaAktif}/kecamatan/${sorotan.kode}`);
	const isMemuatSorotan = $derived(navigating.to?.route.id === page.route.id);
	const teks = $derived(teksJenjang[data.jenjang]);
	const pilihanTema = $derived(pilihanJenjang(page.url, data.jenjang));

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

	function simpanStatusPeta() {
		replaceState(urlDenganStatusPeta(page.url, lapisan, petaDasar), page.state);
	}

	function bukaKecamatan(kode: string) {
		goto(`/${bahasaAktif}/kecamatan/${kode}`);
	}
</script>

<svelte:head>
	<title>{t('brand.name')} · {t('home.heroTitle')}</title>
</svelte:head>

<div class="wrap">
	<div
		class="flex flex-col gap-5 pb-6 pt-7 md:flex-row md:items-start md:justify-between lg:pb-7 lg:pt-8"
	>
		<div class="md:pt-2">
			<h1 class="t-h1">{t('home.title')}</h1>
			<p class="mt-1 text-muted">{t('home.subtitle')}</p>
		</div>
		<div class="flex w-full flex-col gap-2 md:max-w-md" data-tur="cari">
			<p class="text-caption text-muted" aria-hidden="true">{t('home.searchLabel')}</p>
			<RegionSearch
				bind:value={kata}
				cari={(q) => cariWilayah(fetch, q)}
				onselect={buka}
				onuselocation={pakaiLokasiku}
			/>
			<div class="flex flex-wrap items-start gap-x-4">
				<UseLocationButton state={statusLokasi} {akurasiM} onclick={pakaiLokasiku} />
				<p class="pt-3 text-caption text-muted">{t('location.privacy')}</p>
			</div>
		</div>
	</div>

	<section
		class="relative isolate flex min-h-foto overflow-hidden rounded-xl bg-foto text-foto-on"
		aria-labelledby="judul-pahlawan"
	>
		<picture>
			<source media="(max-width: 767px)" srcset="/foto/bundaran-hi-800.webp" />
			<img
				src="/foto/bundaran-hi.webp"
				alt={t('home.photoAlt')}
				width="1600"
				height="1067"
				fetchpriority="high"
				class="absolute inset-0 -z-10 h-full w-full object-cover object-[60%_58%] max-md:object-[67%_50%]"
			/>
		</picture>
		<div class="bayang-foto absolute inset-0 -z-10" aria-hidden="true"></div>
		<div class="flex max-w-xl flex-col items-start px-6 py-7 md:px-9 md:py-8">
			<p class="text-caption font-medium md:text-small">{t('home.kicker')}</p>
			<h2 id="judul-pahlawan" class="t-display mt-3">{t('home.heroTitle')}</h2>
			<p class="mt-3 max-w-md text-small md:text-body">{t('home.lead')}</p>
			<a href="/{bahasaAktif}/peta" class="btn mt-6 gap-5 bg-foto-on text-foto hover:opacity-90">
				{t('home.heroCta')}
				<ArrowRight size={18} strokeWidth={1.75} aria-hidden="true" />
			</a>
		</div>
		<a
			href={sumberFoto}
			target="_blank"
			rel="noopener"
			class="bayang-teks-foto absolute bottom-2 right-3 inline-flex min-h-touch items-center gap-2 text-caption hover:underline"
			aria-label="{t('home.photoCredit')} {t('a11y.opensNewTab')}"
		>
			{t('home.photoCaption')}
			<ExternalLink size={14} strokeWidth={1.75} aria-hidden="true" />
		</a>
	</section>

	<div
		class="mt-5 flex flex-col gap-2 border-b border-line md:flex-row md:items-center md:justify-between"
	>
		<nav aria-label={t('home.areaPicker')} class="-mb-px overflow-x-auto">
			<ul class="flex items-center gap-6">
				<li class="hidden shrink-0 text-small text-muted md:block" aria-hidden="true">
					{t('home.areaPicker')}
				</li>
				{#each daftarSorotan as wilayah (wilayah.kode)}
					<li class="shrink-0">
						<a
							href={tautanSorotan(page.url, wilayah.kode)}
							data-sveltekit-noscroll
							aria-current={wilayah.kode === sorotan.kode ? 'page' : undefined}
							class="inline-flex min-h-12 items-center border-b-2 border-transparent text-small text-muted transition-colors duration-fast ease-standard hover:text-fg aria-[current=page]:border-fg aria-[current=page]:font-semibold aria-[current=page]:text-fg"
						>
							{wilayah.nama}
						</a>
					</li>
				{/each}
			</ul>
		</nav>
		<a
			href="/{bahasaAktif}/tentang-data"
			class="inline-flex min-h-touch items-center gap-2 self-start text-caption text-muted hover:text-fg hover:underline md:self-auto"
		>
			<FileText size={14} strokeWidth={1.75} aria-hidden="true" />
			{t('home.trustLink')}
		</a>
	</div>

	<div
		class="transition-opacity duration-fast ease-standard {isMemuatSorotan ? 'opacity-60' : ''}"
		aria-busy={isMemuatSorotan}
	>
		{#if !data.ringkasan.ok}
			<div class="mt-8">
				<ErrorState
					title={t('home.loadError', { name: sorotan.nama })}
					description={data.ringkasan.galat.message}
					onretry={invalidateAll}
				/>
			</div>
		{:else}
			{@const ringkasan = data.ringkasan.data}
			{@const nama = ringkasan.kecamatan.nama}
			{@const kota = ringkasan.kecamatan.hierarki.find((h) => h.tingkat === 'kota')}
			{@const lokasi = kota ? `${nama} · ${namaWilayahPendek(kota.nama)}` : nama}
			{@const indikator = indikatorPendidikan(ringkasan.tema, data.jenjang)}
			{@const dapil = ringkasan.wakilku.dapil}
			{@const komisi = ringkasan.wakilku.komisi_pengawas.join(', ')}
			<section class="mt-8 lg:mt-10" aria-labelledby="judul-wilayah">
				<div class="flex flex-wrap items-end justify-between gap-x-6 gap-y-2">
					<div>
						<p class="t-over">{lokasi}</p>
						<h2 id="judul-wilayah" class="t-h2 mt-1">{t('home.area.title')}</h2>
					</div>
					<a
						href={halamanWilayah}
						class="inline-flex min-h-touch items-center gap-2 text-small font-semibold underline-offset-4 hover:underline"
					>
						{t('home.area.open', { name: nama })}
						<ArrowRight size={16} strokeWidth={1.75} aria-hidden="true" />
					</a>
				</div>
				<div class="mt-4 flex flex-wrap items-center gap-3">
					<p id="judul-jenjang-beranda" class="text-small font-semibold text-muted">
						{t('level.label')}
					</p>
					<ChoiceChips labelledby="judul-jenjang-beranda" items={pilihanTema} />
				</div>

				<div class="mt-5 grid gap-8 lg:grid-cols-[1.25fr_1fr] lg:gap-9">
					<div class="min-w-0">
						<div class="rounded-lg border border-line bg-surface p-4 md:p-5">
							<p class="text-small font-semibold">{t(teks.berandaPeta)}</p>
							{#if indikator}
								<p class="text-caption text-muted">
									{t('home.map.lead', { period: indikator.periode })}
								</p>
							{/if}
							<PetaBerlapis
								class="mt-4 h-96"
								tema={data.peta.ok ? data.peta.data : null}
								temaSekolah={data.jenjang}
								terpilih={sorotan.kode}
								onpilih={bukaKecamatan}
								bind:lapisan
								bind:petaDasar
								onchange={simpanStatusPeta}
								nama="peta-dasar-beranda"
								judulLegenda={teks.legenda}
							/>
							{#if !data.peta.ok}
								<p class="mt-3 text-caption text-muted">{t('home.map.error')}</p>
							{/if}
							<ul class="mt-4 flex flex-wrap gap-x-4 gap-y-2">
								{#each urutanKelas as kelas (kelas)}
									<li class="flex items-center gap-2 text-caption">
										<span class="h-3 w-5 rounded-xs {kotakKelasPeta[kelas]}" aria-hidden="true"
										></span>
										{t(kunciKelas[kelas])}
									</li>
								{/each}
								<li class="flex items-center gap-2 text-caption">
									<span class="h-3 w-5 rounded-xs border-2 border-fg" aria-hidden="true"></span>
									{t('home.map.selected')}
								</li>
							</ul>
						</div>
						<div class="mt-2 flex flex-wrap items-center justify-between gap-x-4">
							<p class="flex items-center gap-1.5 text-caption text-muted">
								<MapPin size={14} strokeWidth={1.75} aria-hidden="true" />
								{lokasi}
							</p>
							<a
								href={tautanPetaPenuh(page.url, bahasaAktif)}
								class="inline-flex min-h-touch items-center text-caption font-semibold underline-offset-4 hover:underline"
							>
								{t('home.map.open')}
							</a>
						</div>
					</div>

					<article
						class="flex flex-col border-line max-lg:border-t max-lg:pt-6 lg:border-l lg:pl-8"
						aria-labelledby="judul-pendidikan"
					>
						<div class="flex flex-wrap items-center justify-between gap-2">
							<p class="text-small font-semibold">{t('sector.edu')}</p>
							{#if indikator}
								<ProvenanceBadge
									angka={indikator.kesenjangan}
									nilai={tampilkanAngka(indikator.kesenjangan).nilai}
									compact
								/>
							{/if}
						</div>
						<h3 id="judul-pendidikan" class="t-h2 mt-3">{t(teks.berandaJudul)}</h3>
						{#if indikator}
							{@const acuan = Math.max(indikator.kebutuhan.nilai, indikator.kapasitas.nilai)}
							{@const selisih = indikator.kesenjangan.nilai}
							{@const batang = [
								{ label: t(teks.lulusan), angka: indikator.kebutuhan, warna: 'bg-chart' },
								{ label: t(teks.kursiNegeri), angka: indikator.kapasitas, warna: 'bg-chart-2' }
							]}
							<p class="mt-2 text-small text-muted">
								{t(teks.berandaLead, { name: nama, period: indikator.periode })}
							</p>
							<dl class="mt-5 flex flex-col gap-3">
								{#each batang as item (item.label)}
									<div class="flex flex-wrap items-baseline justify-between gap-x-3 text-small">
										<dt>{item.label}</dt>
										<dd class="num font-semibold">{tampilkanAngka(item.angka).nilai}</dd>
										<dd class="mt-1.5 h-2 w-full rounded-xs bg-sunken" aria-hidden="true">
											<div
												class="h-full rounded-xs {item.warna}"
												style="width: {persenTerhadap(item.angka.nilai, acuan)}%"
											></div>
										</dd>
									</div>
								{/each}
							</dl>
							<div class="mt-5 flex flex-wrap items-center gap-x-3 gap-y-2">
								<p class="n-xl">{formatAngka(Math.abs(selisih), bahasaAktif)}</p>
								<p class="text-small leading-tight">
									{selisih >= 0 ? t('home.edu.short') : t('home.edu.surplus')}
								</p>
								<GapClassChip kelas={indikator.kelas} />
							</div>
							<p class="mt-3 text-caption text-muted">{t('services.zoningNote')}</p>
						{:else}
							<p class="mt-3 text-small text-muted">{t('tema.empty')}</p>
						{/if}
						<a
							href="{halamanWilayah}/{data.jenjang}"
							class="mt-auto flex min-h-touch items-center justify-between gap-3 border-t border-line pt-3 text-small font-semibold underline-offset-4 hover:underline max-lg:mt-5"
						>
							{t('home.edu.more', { name: nama })}
							<ArrowRight size={18} strokeWidth={1.75} aria-hidden="true" />
						</a>
					</article>
				</div>
			</section>

			<div
				class="mt-10 grid gap-10 border-t border-line pt-8 lg:mt-12 lg:grid-cols-[1.2fr_1fr] lg:gap-10"
			>
				<section aria-labelledby="judul-anggaran">
					<p class="t-over">{t('bagian.anggaran')}</p>
					<h2 id="judul-anggaran" class="t-h2 mt-1">{t('home.money.title')}</h2>
					<p class="mt-1 text-small text-muted">{t('home.money.lead', { name: nama })}</p>
					<div
						class="mt-5 flex items-center justify-between gap-3 border-b border-line pb-2 text-caption text-muted"
					>
						<span>{t('home.money.list')}</span>
						{#if data.paket.ok && data.paket.data.meta}
							<span class="num">{t('home.money.total', { n: data.paket.data.meta.total })}</span>
						{/if}
					</div>
					{#if !data.paket.ok}
						<p class="py-5 text-small text-muted">{t('home.money.error')}</p>
					{:else if data.paket.data.data.length === 0}
						<p class="py-5 text-small text-muted">{t('money.empty')}</p>
					{:else}
						<ul>
							{#each data.paket.data.data.slice(0, batasPaketBeranda) as paket (paket.kode)}
								{@const nilai = paket.nilai_kontrak ?? paket.hps}
								<li class="border-b border-line">
									<a
										href="/{bahasaAktif}/paket/{paket.kode}"
										class="group flex items-center gap-3 py-4"
									>
										<span class="icon-tile h-10 w-10">
											<Wallet size={18} strokeWidth={1.75} aria-hidden="true" />
										</span>
										<span class="flex min-w-0 flex-1 flex-col gap-0.5">
											<span
												class="line-clamp-2 text-small font-semibold underline-offset-4 group-hover:underline"
												>{paket.nama}</span
											>
											<span class="flex flex-wrap items-center gap-2 text-caption text-muted">
												{paket.tahap} · {paket.tahun_anggaran}
												{#if paket.is_simulation}<CertaintyBadge level="simulasi" />{/if}
											</span>
										</span>
										{#if nilai}
											<span class="shrink-0 text-right">
												<span class="num block text-small font-semibold"
													>{tampilkanAngka(nilai).nilai}</span
												>
												<span class="block text-caption text-muted"
													>{paket.nilai_kontrak ? t('pkg.contract') : t('pkg.hps')}</span
												>
											</span>
										{/if}
										<ArrowRight
											size={16}
											strokeWidth={1.75}
											aria-hidden="true"
											class="shrink-0 text-muted max-md:hidden"
										/>
									</a>
								</li>
							{/each}
						</ul>
					{/if}
					<div class="mt-3 flex flex-wrap items-start justify-between gap-x-6 gap-y-2">
						<p class="max-w-md text-caption text-muted">{t('money.locationNote')}</p>
						<a
							href="{halamanWilayah}?bagian=anggaran"
							class="inline-flex min-h-touch items-center text-small font-semibold underline-offset-4 hover:underline"
						>
							{t('home.money.all')}
						</a>
					</div>
				</section>

				<section
					class="border-line max-lg:border-t max-lg:pt-8 lg:border-l lg:pl-8"
					aria-labelledby="judul-wakilku"
				>
					<p class="t-over">{t('bagian.wakilku')}</p>
					<h2 id="judul-wakilku" class="t-h2 mt-1">{t('home.reps.title')}</h2>
					<p class="mt-1 text-small text-muted">
						{t('home.reps.lead', { komisi, name: nama, dapil: dapil.nama })}
					</p>
					<a
						href="{halamanWilayah}?bagian=wakilku"
						class="group mt-5 flex items-center gap-3 rounded-group bg-sunken px-4 py-3"
					>
						<Landmark size={22} strokeWidth={1.75} aria-hidden="true" class="shrink-0" />
						<span class="flex min-w-0 flex-1 flex-col">
							<span class="text-caption text-muted">{t('home.reps.oversight')}</span>
							<span class="text-small font-semibold">{t('home.reps.commission', { komisi })}</span>
						</span>
						<ArrowRight
							size={18}
							strokeWidth={1.75}
							aria-hidden="true"
							class="shrink-0 transition-transform duration-fast ease-standard group-hover:translate-x-0.5"
						/>
					</a>
					{#if wakilPengawas(dapil).length}
						<p class="mt-5 text-caption text-muted">
							{t('home.reps.members', { dapil: dapil.nama })}
						</p>
						<ul class="mt-1">
							{#each wakilPengawas(dapil) as wakil (wakil.id)}
								<li class="border-b border-line">
									<a
										href="/{bahasaAktif}/anggota/{wakil.id}"
										class="group flex min-h-14 items-center gap-3 py-2"
									>
										<FotoWakil
											nama={wakil.nama}
											foto={wakil.foto_url}
											class="h-10 w-10 text-caption"
										/>
										<span class="flex min-w-0 flex-1 flex-col">
											<span
												class="truncate text-small font-semibold underline-offset-4 group-hover:underline"
												>{wakil.nama}</span
											>
											<span class="truncate text-caption text-muted">{wakil.fraksi}</span>
										</span>
									</a>
								</li>
							{/each}
						</ul>
					{/if}
					<div class="mt-3 flex flex-wrap items-start justify-between gap-x-6 gap-y-2">
						<p class="text-caption text-muted">
							{t('home.reps.noScore')}
							{t('photo.dprd')}
						</p>
						<a
							href="{halamanWilayah}?bagian=wakilku"
							class="inline-flex min-h-touch items-center text-small font-semibold underline-offset-4 hover:underline"
						>
							{t('home.reps.all')}
						</a>
					</div>
				</section>
			</div>
		{/if}
	</div>

	<section
		class="mt-10 flex flex-col gap-5 rounded-lg border border-line bg-secondary p-6 md:p-7 lg:mt-12 lg:flex-row lg:items-center lg:justify-between"
		aria-labelledby="judul-aksi"
	>
		<div>
			<p class="t-over">{t('home.act.label')}</p>
			<h2 id="judul-aksi" class="t-h2 mt-1">{t('home.act.title')}</h2>
			<p class="mt-1 text-small text-muted">{t('home.act.lead')}</p>
		</div>
		<div class="flex flex-col gap-2.5 md:flex-row">
			<a
				href="/{bahasaAktif}/tanya?w={encodeURIComponent(sorotan.kode)}"
				class="btn btn-outline btn-lg"
			>
				<MessageCircle size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('ask.title')}
			</a>
			<a
				href="/{bahasaAktif}/aspirasi?w={encodeURIComponent(sorotan.kode)}"
				class="btn btn-primary btn-lg"
			>
				<FileText size={18} strokeWidth={1.75} aria-hidden="true" />
				{t('acc.createPacket')}
			</a>
		</div>
	</section>

	<section class="mt-8" aria-labelledby="judul-langkah">
		<h2 id="judul-langkah" class="sr-only">{t('home.more')}</h2>
		<ul class="grid gap-x-8 md:grid-cols-3">
			{#each langkahLain as langkah (langkah.path)}
				{@const Ikon = langkah.ikon}
				<li class="border-b border-line">
					<a href="/{bahasaAktif}{langkah.path}" class="group flex items-start gap-3 py-4">
						<span class="icon-tile h-10 w-10">
							<Ikon size={18} strokeWidth={1.75} aria-hidden="true" />
						</span>
						<span class="flex min-w-0 flex-col gap-0.5">
							<span class="text-small font-semibold underline-offset-4 group-hover:underline"
								>{langkah.judul}</span
							>
							<span class="text-caption text-muted">{langkah.isi}</span>
						</span>
					</a>
				</li>
			{/each}
		</ul>
	</section>

	<section
		class="flex flex-col gap-3 py-7 md:flex-row md:items-center md:gap-4"
		aria-labelledby="judul-keterbukaan"
	>
		<FileLock size={22} strokeWidth={1.75} aria-hidden="true" class="shrink-0" />
		<div class="flex-1">
			<h2 id="judul-keterbukaan" class="text-small font-semibold">{t('home.open.title')}</h2>
			<p class="text-caption text-muted">{t('home.open.lead')}</p>
		</div>
		<a
			href="/{bahasaAktif}/data-belum-dibuka"
			class="inline-flex min-h-touch items-center gap-2 self-start text-small font-semibold underline-offset-4 hover:underline md:self-auto"
		>
			{t('unopened.title')}
			<ArrowRight size={16} strokeWidth={1.75} aria-hidden="true" />
		</a>
	</section>

	<section class="border-t border-line pt-8" aria-labelledby="judul-sumber">
		<div class="grid gap-6 lg:grid-cols-[1fr_2fr] lg:gap-10">
			<div>
				<h2 id="judul-sumber" class="t-h3">{t('home.trustTitle')}</h2>
				<p class="mt-2 text-small text-muted">{t('home.trustLead')}</p>
				<a
					href="/{bahasaAktif}/tentang-data"
					class="mt-3 inline-flex min-h-touch items-center gap-2 text-small font-semibold underline-offset-4 hover:underline"
				>
					{t('action.readMethod')}
					<ArrowRight size={16} strokeWidth={1.75} aria-hidden="true" />
				</a>
			</div>
			<ul class="grid gap-x-6 gap-y-4 md:grid-cols-2">
				{#each labelKepastian as level (level)}
					<li class="flex flex-col items-start gap-2">
						<CertaintyBadge {level} size="md" />
						<p class="text-small text-muted">{t(`cert.${level}.desc`)}</p>
					</li>
				{/each}
			</ul>
		</div>
	</section>
</div>
