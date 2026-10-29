<script lang="ts">
	import { ArrowLeft, Building2, ChevronRight, House } from '@lucide/svelte';
	import { goto, invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import ChoiceChips from '$lib/components/ChoiceChips.svelte';
	import ErrorState from '$lib/components/ErrorState.svelte';
	import AnalisisKota from '$lib/features/kota/AnalisisKota.svelte';
	import DayaTampungKota from '$lib/features/kota/DayaTampungKota.svelte';
	import SebaranSekolah from '$lib/features/kota/SebaranSekolah.svelte';
	import type { CakupanSekolah } from '$lib/features/kota/pendidikan-kota';
	import { pilihanJenjang } from '$lib/features/pendidikan/jenjang';
	import { bahasaAktif, t } from '$lib/i18n';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	let cakupan = $state<CakupanSekolah>('negeri');
	const hasil = $derived(data.hasil);
	const kota = $derived(hasil.ok ? hasil.data.kabupaten_kota : []);
</script>

<svelte:head>
	<title
		>{hasil.ok ? `${hasil.data.wilayah.nama} · ${t('city.title')}` : t('city.title')} · {t(
			'brand.name'
		)}</title
	>
</svelte:head>

{#if !hasil.ok}
	<div class="wrap pt-8">
		<h1 class="t-h1 mb-5">{t('city.title')}</h1>
		<ErrorState
			title={t('city.failed')}
			description={hasil.galat.message}
			onretry={invalidateAll}
		/>
	</div>
{:else}
	<div class="border-b border-line bg-surface">
		<div class="wrap py-6 lg:py-8">
			<nav aria-label={t('nav.breadcrumb')}>
				<ol class="flex flex-wrap items-center gap-2 text-small text-muted">
					<li>
						<a
							class="inline-flex min-h-touch items-center gap-1 hover:underline"
							href="/{bahasaAktif}"><House size={14} aria-hidden="true" />{t('nav.home')}</a
						>
					</li>
					<li aria-hidden="true"><ChevronRight size={16} /></li>
					<li>
						<a
							class="inline-flex min-h-touch items-center hover:underline"
							href="/{bahasaAktif}/dasbor">{t('city.title')}</a
						>
					</li>
					<li aria-current="page">{hasil.data.wilayah.nama}</li>
				</ol>
			</nav>
			<p class="t-over mt-4 flex items-center gap-2 text-primary">
				<Building2 size={18} aria-hidden="true" />{t('city.scope')}
			</p>
			<h1 class="t-h1 mt-3">{hasil.data.wilayah.nama}</h1>
			<p class="t-lead mt-3 max-w-prose">{t('city.lead')}</p>
			<div class="mt-6 flex flex-wrap items-end gap-6">
				<div>
					<p id="judul-jenjang-kota" class="mb-3 text-small font-semibold">{t('level.label')}</p>
					<ChoiceChips
						labelledby="judul-jenjang-kota"
						items={pilihanJenjang(page.url, data.tema)}
					/>
				</div>
				{#if kota.length}
					<div class="min-w-0 flex-1">
						<label for="pilih-kota" class="mb-3 block text-small font-semibold"
							>{t('city.choose')}</label
						><select
							id="pilih-kota"
							class="w-full rounded-md border border-line-strong bg-surface p-3 text-body"
							onchange={(event) => {
								void goto(`/${bahasaAktif}/kota/${event.currentTarget.value}?tema=${data.tema}`);
							}}
							value={data.kode}
							>{#each kota as k (k.kode)}<option value={k.kode}>{k.nama}</option>{/each}</select
						>
					</div>
				{/if}
			</div>
		</div>
	</div>
	<div class="wrap flex min-w-0 flex-col gap-10 pt-8 lg:pt-10">
		<DayaTampungKota data={hasil.data} tema={data.tema} />
		<SebaranSekolah data={hasil.data} tema={data.tema} bind:cakupan />
		<AnalisisKota data={hasil.data} tema={data.tema} {cakupan} />
		<a href="/{bahasaAktif}#cari-wilayah" class="btn btn-outline self-start"
			><ArrowLeft size={18} aria-hidden="true" />{t('city.backSearch')}</a
		>
	</div>
{/if}
