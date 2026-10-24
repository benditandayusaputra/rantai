<script lang="ts">
	import { Info, ShieldCheck } from '@lucide/svelte';
	import { ambilDetailWilayah, ambilSekolahTerdekat, cariWilayah, tangkap } from '$lib/api/klien';
	import type { SekitarTitik, Titik, WilayahDitemukan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import RegionSearch from '$lib/components/RegionSearch.svelte';
	import UseLocationButton from '$lib/components/UseLocationButton.svelte';
	import {
		batasDaftarNegeri,
		pecahanTekanan,
		radiusBawaanM,
		radiusKm,
		terdekatDulu
	} from '$lib/features/cekSekolah/cekSekolah';
	import KartuSekolah from '$lib/features/cekSekolah/KartuSekolah.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { cariPosisi } from '$lib/features/lapor/lapor';
	import type { StatusLokasi } from '$lib/features/wilayah/lokasi';
	import { t } from '$lib/i18n';

	let kata = $state('');
	let statusLokasi = $state<StatusLokasi>('idle');
	let akurasiM = $state(0);
	let pusat = $state<{ label: string; titik: Titik } | null>(null);
	let sekitar = $state<SekitarTitik | null>(null);
	let status = $state<'awal' | 'memuat' | 'siap' | 'gagal'>('awal');

	const km = radiusKm(radiusBawaanM);
	const rasio = $derived(sekitar?.tekanan_radius.rasio ?? null);
	const alternatif = $derived(sekitar ? [...sekitar.spmb_bersama, ...sekitar.madrasah] : []);

	async function muat(label: string, titik: Titik) {
		pusat = { label, titik };
		status = 'memuat';
		const hasil = await tangkap(ambilSekolahTerdekat(fetch, titik, radiusBawaanM));
		sekitar = hasil.ok ? hasil.data : null;
		status = hasil.ok ? 'siap' : 'gagal';
	}

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

<div class="wrap flex flex-col gap-8 py-6 lg:py-10">
	<header class="flex max-w-prose flex-col gap-4">
		<h1 class="t-h1">{t('check.title')}</h1>
		<p class="t-lead">{t('check.lead')}</p>
		<div class="flex flex-col gap-3 md:flex-row md:items-start">
			<div class="flex-1">
				<RegionSearch
					bind:value={kata}
					cari={(q) => cariWilayah(fetch, q)}
					onselect={pilihWilayah}
				/>
			</div>
			<UseLocationButton state={statusLokasi} {akurasiM} onclick={pakaiLokasiku} />
		</div>
		<p class="flex gap-2 text-small text-muted">
			<ShieldCheck size={18} strokeWidth={1.75} class="mt-0.5 shrink-0" aria-hidden="true" />
			{t('check.privacy')}
		</p>
	</header>

	{#if status === 'awal'}
		<p class="text-muted">{t('check.start')}</p>
	{:else if status === 'memuat'}
		<p class="text-muted" aria-live="polite">{t('check.loading', { km })}</p>
	{:else if status === 'gagal' || !sekitar}
		<ErrorState title={t('check.failed')} />
	{:else}
		<div class="grid gap-8 lg:grid-cols-12" aria-live="polite">
			<section class="flex flex-col gap-4 lg:col-span-7" aria-labelledby="judul-tekanan">
				<div class="flex flex-col gap-3 rounded-lg border border-line bg-surface p-5">
					<p class="t-over">{pusat?.label}</p>
					<h2 id="judul-tekanan" class="t-h3">{t('check.pressure', { km })}</h2>
					{#if rasio}
						{@const pecahan = pecahanTekanan(rasio.nilai)}
						<p class="font-serif text-h2 lg:text-h2-lg">{t('services.sentence', pecahan)}</p>
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
				</div>
				{#if sekitar.keterbatasan.length}
					<ul class="list-disc space-y-1 pl-5 text-small text-muted">
						{#each sekitar.keterbatasan as kunci (kunci)}<li>{t(`pen.limit.${kunci}`)}</li>{/each}
					</ul>
				{/if}
			</section>
			<div class="flex flex-col gap-6 lg:col-span-5">
				<section class="flex flex-col gap-3" aria-labelledby="judul-negeri">
					<h2 id="judul-negeri" class="t-h3">{t('check.nearest')}</h2>
					{#if sekitar.smp_negeri.length}
						<ol class="flex flex-col gap-2">
							{#each terdekatDulu(sekitar.smp_negeri, batasDaftarNegeri) as sekolah, i (sekolah.npsn)}
								<KartuSekolah {sekolah} urutan={i + 1} />
							{/each}
						</ol>
					{:else}
						<EmptyState title={t('check.empty', { km })} />
					{/if}
				</section>
				{#if alternatif.length}
					<section class="flex flex-col gap-3" aria-labelledby="judul-alternatif">
						<h2 id="judul-alternatif" class="t-h3">{t('check.alternatives')}</h2>
						<ul class="flex flex-col gap-2">
							{#each terdekatDulu(alternatif, alternatif.length) as sekolah (sekolah.npsn)}
								<KartuSekolah {sekolah} />
							{/each}
						</ul>
					</section>
				{/if}
			</div>
		</div>
	{/if}
</div>
