<script lang="ts">
	import { Info, MapPin, School, Search, ShieldCheck, Sparkles } from '@lucide/svelte';
	import { untrack } from 'svelte';
	import { page } from '$app/state';
	import { ambilDetailWilayah, ambilSekolahTerdekat, cariWilayah, tangkap } from '$lib/api/klien';
	import type { SekitarTitik, Titik, WilayahDitemukan } from '$lib/api/tipe';
	import ChoiceChips from '$lib/components/ChoiceChips.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import RegionSearch from '$lib/components/RegionSearch.svelte';
	import UseLocationButton from '$lib/components/UseLocationButton.svelte';
	import {
		batasDaftarNegeri,
		pecahanTekanan,
		radiusCekSekolahM,
		radiusKm,
		terdekatDulu
	} from '$lib/features/cekSekolah/cekSekolah';
	import KartuSekolah from '$lib/features/cekSekolah/KartuSekolah.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { cariPosisi } from '$lib/features/lapor/lapor';
	import {
		bacaTemaPendidikan,
		pilihanJenjang,
		teksJenjang
	} from '$lib/features/pendidikan/jenjang';
	import type { StatusLokasi } from '$lib/features/wilayah/lokasi';
	import { t } from '$lib/i18n';

	let kata = $state('');
	let statusLokasi = $state<StatusLokasi>('idle');
	let akurasiM = $state(0);
	let pusat = $state<{ label: string; titik: Titik } | null>(null);
	let sekitar = $state<SekitarTitik | null>(null);
	let status = $state<'awal' | 'memuat' | 'siap' | 'gagal'>('awal');

	const jenjang = $derived(bacaTemaPendidikan(page.url.searchParams.get('tema')));
	const teks = $derived(teksJenjang[jenjang]);
	const radiusM = $derived(radiusCekSekolahM[jenjang]);
	const km = $derived(radiusKm(radiusM));
	const rasio = $derived(sekitar?.tekanan_radius.rasio ?? null);
	const alternatif = $derived(sekitar ? [...sekitar.spmb_bersama, ...sekitar.madrasah] : []);

	async function muat(label: string, titik: Titik) {
		pusat = { label, titik };
		status = 'memuat';
		const hasil = await tangkap(ambilSekolahTerdekat(fetch, titik, radiusM, jenjang));
		sekitar = hasil.ok ? hasil.data : null;
		status = hasil.ok ? 'siap' : 'gagal';
	}

	$effect(() => {
		void jenjang;
		untrack(() => pusat && muat(pusat.label, pusat.titik));
	});

	async function pilihWilayah(wilayah: WilayahDitemukan) {
		const detail = await tangkap(ambilDetailWilayah(fetch, wilayah.kode));
		if (detail.ok) await muat(t('check.near', { nama: wilayah.nama }), detail.data.titik_tengah);
		else status = 'gagal';
	}

	async function pakaiLokasiku() {
		statusLokasi = 'locating';
		const hasil = await cariPosisi(navigator.geolocation);
		if (hasil.status === 'ditolak' || hasil.status === 'gagal') {
			statusLokasi = hasil.status === 'ditolak' ? 'denied' : 'failed';
			return;
		}
		akurasiM = hasil.akurasiM;
		statusLokasi = hasil.status === 'akurat' ? 'found' : 'coarse';
		await muat(t('check.nearMe'), hasil.titik);
	}
</script>

<svelte:head>
	<title>{t('check.title')} · {t('brand.name')}</title>
</svelte:head>

<div class="border-b border-line bg-surface">
	<div
		class="wrap grid gap-6 pb-8 pt-8 lg:grid-cols-[1fr_1.1fr] lg:items-center lg:gap-12 lg:pb-12 lg:pt-12"
	>
		<header class="flex flex-col gap-3">
			<span class="icon-tile h-12 w-12">
				<School size={24} strokeWidth={1.75} aria-hidden="true" />
			</span>
			<h1 class="t-h1">{t('check.title')}</h1>
			<p class="t-lead max-w-prose">{t(teks.cekLead)}</p>
			<div class="flex flex-wrap items-center gap-3">
				<p id="judul-jenjang-cek" class="text-small font-semibold text-muted">{t('level.label')}</p>
				<ChoiceChips labelledby="judul-jenjang-cek" items={pilihanJenjang(page.url, jenjang)} />
			</div>
		</header>
		<div class="card flex flex-col gap-4 p-4 shadow-2 md:p-5">
			<RegionSearch bind:value={kata} cari={(q) => cariWilayah(fetch, q)} onselect={pilihWilayah} />
			<UseLocationButton state={statusLokasi} {akurasiM} onclick={pakaiLokasiku} />
			<p class="flex gap-2 border-t border-line pt-4 text-small text-muted">
				<ShieldCheck size={18} strokeWidth={1.75} class="mt-0.5 shrink-0" aria-hidden="true" />
				{t('check.privacy')}
			</p>
		</div>
	</div>
</div>

<div class="wrap flex flex-col gap-8 pt-8 lg:pt-10">
	{#if status === 'awal'}
		<div
			class="flex flex-col items-center gap-3 rounded-lg border-ctrl border-dashed border-line-strong px-6 py-12 text-center"
		>
			<span class="icon-tile h-12 w-12"
				><Search size={24} strokeWidth={1.75} aria-hidden="true" /></span
			>
			<p class="max-w-prose text-muted">{t('check.start')}</p>
		</div>
	{:else if status === 'memuat'}
		<div class="grid gap-6 lg:grid-cols-12" aria-busy="true">
			<p class="sr-only" aria-live="polite">{t('check.loading', { km })}</p>
			<div class="skel h-64 rounded-lg lg:col-span-7"></div>
			<div class="flex flex-col gap-3 lg:col-span-5">
				<div class="skel h-24 rounded-lg"></div>
				<div class="skel h-24 rounded-lg"></div>
				<div class="skel h-24 rounded-lg"></div>
			</div>
		</div>
	{:else if status === 'gagal' || !sekitar}
		<ErrorState title={t('check.failed')} />
	{:else}
		<div class="grid gap-8 lg:grid-cols-12" aria-live="polite">
			<section class="flex flex-col gap-4 lg:col-span-7" aria-labelledby="judul-tekanan">
				<div class="card flex flex-col gap-4 p-5 md:p-6 lg:sticky lg:top-24">
					<p class="t-over flex items-center gap-2">
						<MapPin size={14} strokeWidth={2} aria-hidden="true" />{pusat?.label}
					</p>
					<h2 id="judul-tekanan" class="t-h3">{t('check.pressure', { km })}</h2>
					{#if rasio}
						{@const pecahan = pecahanTekanan(rasio.nilai)}
						<p class="font-heading text-h2 lg:text-h2-lg">{t(teks.kalimatTekanan, pecahan)}</p>
						<div class="self-start">
							<ProvenanceBadge angka={rasio} nilai={tampilkanAngka(rasio).nilai} />
						</div>
					{:else}
						<p class="text-muted">{t('check.noRatio')}</p>
					{/if}
					<p class="flex gap-2 rounded-md bg-info-bg p-3 text-small">
						<Info
							size={18}
							strokeWidth={1.75}
							class="mt-0.5 shrink-0 text-info"
							aria-hidden="true"
						/>
						{t('services.zoningNote')}
					</p>
					{#if sekitar.keterbatasan.length}
						<ul class="list-disc space-y-1 pl-5 text-small text-muted">
							{#each sekitar.keterbatasan as kunci (kunci)}<li>{t(`pen.limit.${kunci}`)}</li>{/each}
						</ul>
					{/if}
				</div>
			</section>
			<div class="flex flex-col gap-8 lg:col-span-5">
				<section class="flex flex-col gap-3" aria-labelledby="judul-negeri">
					<h2 id="judul-negeri" class="flex items-center gap-2 text-h3 font-bold">
						<School size={20} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
						{t(teks.cekTerdekat)}
					</h2>
					{#if sekitar.negeri.length}
						<ol class="flex flex-col gap-2">
							{#each terdekatDulu(sekitar.negeri, batasDaftarNegeri) as sekolah, i (sekolah.npsn)}
								<KartuSekolah {sekolah} urutan={i + 1} tema={jenjang} />
							{/each}
						</ol>
					{:else}
						<EmptyState title={t(teks.cekKosong, { km })} />
					{/if}
				</section>
				{#if alternatif.length}
					<section class="flex flex-col gap-3" aria-labelledby="judul-alternatif">
						<h2 id="judul-alternatif" class="flex items-center gap-2 text-h3 font-bold">
							<Sparkles size={20} strokeWidth={1.75} aria-hidden="true" class="text-muted" />
							{t('check.alternatives')}
						</h2>
						<ul class="flex flex-col gap-2">
							{#each terdekatDulu(alternatif, alternatif.length) as sekolah (sekolah.npsn)}
								<KartuSekolah {sekolah} tema={jenjang} />
							{/each}
						</ul>
					</section>
				{/if}
			</div>
		</div>
	{/if}
</div>
