<script lang="ts">
	import {
		ArrowRight,
		Calendar,
		FileText,
		GraduationCap,
		LocateFixed,
		Wallet
	} from '@lucide/svelte';
	import type { TemaKecamatan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ChoiceChips from '$lib/components/ChoiceChips.svelte';
	import { page } from '$app/state';
	import { statusLayanan } from './kecamatan';
	import PenanggungJawabLayanan from './PenanggungJawabLayanan.svelte';
	import DataLayanan from './DataLayanan.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import SeatSummary from '$lib/components/SeatSummary.svelte';
	import {
		bacaTemaPendidikan,
		pilihanJenjang,
		isTemaPendidikan,
		teksJenjang
	} from '$lib/features/pendidikan/jenjang';
	import { bahasaAktif, t } from '$lib/i18n';
	import SuaraPengawas from './SuaraPengawas.svelte';

	let { tema, kodeKecamatan }: { tema: TemaKecamatan[]; kodeKecamatan: string } = $props();

	const jenjangAktif = $derived(bacaTemaPendidikan(page.url.searchParams.get('tema')));
	const temaTampil = $derived(tema.filter((item) => item.kode === jenjangAktif));
	const pilihan = $derived(pilihanJenjang(page.url, jenjangAktif));

	const tautanCekSekolah = (kode: string) =>
		`/${bahasaAktif}/cek-sekolah?w=${kodeKecamatan}&tema=${kode}`;
	const tautanAspirasi = (kode: string) =>
		`/${bahasaAktif}/aspirasi?w=${kodeKecamatan}&tema=${kode}&masalah=${kode === 'pendidikan' ? 'kursi_smp_kurang' : 'kursi_sma_smk_kurang'}`;
</script>

<div class="flex flex-col gap-6">
	<div class="rounded-lg border border-line bg-sunken p-4 md:p-5">
		<h3 class="t-h3">{t('city.districtContext.title')}</h3>
		<p class="mt-2 text-small text-muted">{t('city.districtContext.body')}</p>
		<a
			href="/{bahasaAktif}/kota/{kodeKecamatan
				.split('.')
				.slice(0, 2)
				.join('.')}?tema={jenjangAktif}"
			class="btn btn-outline mt-4"
			>{t('city.districtContext.action')}<ArrowRight size={18} aria-hidden="true" /></a
		>
	</div>
	<div>
		<p id="judul-pilih-layanan" class="mb-3 text-small font-semibold">{t('services.choose')}</p>
		<ChoiceChips labelledby="judul-pilih-layanan" items={pilihan} />
		<p class="mt-3 text-caption text-muted">{t('services.availableScope')}</p>
	</div>
	{#each temaTampil as item (item.kode)}
		{@const kondisi = statusLayanan(item.indikator)}
		<section aria-labelledby="tema-{item.kode}" class="flex flex-col gap-6">
			<div class="card flex flex-col gap-5 p-4 md:p-6">
				<div class="flex items-start gap-3">
					<span class="icon-tile">
						<GraduationCap size={20} strokeWidth={1.75} aria-hidden="true" />
					</span>
					<div class="flex min-w-0 flex-col gap-1">
						<p class="t-over">{item.nama[bahasaAktif]}</p>
						<h3 id="tema-{item.kode}" class="t-h3">{item.judul[bahasaAktif]}</h3>
					</div>
				</div>
				{#if item.indikator && kondisi !== 'belum-dinilai'}
					<p class="inline-flex items-center gap-1 text-small text-muted">
						<Calendar size={16} strokeWidth={1.75} aria-hidden="true" />
						{t('tema.period', { period: item.indikator.periode })}
					</p>
					<div class="grid gap-3 xs:grid-cols-2 lg:grid-cols-4 lg:gap-4">
						<KeyNumberCard label={t('tema.needs')} angka={item.indikator.kebutuhan} />
						<KeyNumberCard label={t('tema.capacity')} angka={item.indikator.kapasitas} />
						<KeyNumberCard label={t('tema.gap')} angka={item.indikator.kesenjangan} />
						{#if item.indikator.rasio}
							<KeyNumberCard label={t('tema.ratio')} angka={item.indikator.rasio} />
						{/if}
					</div>
				{:else}
					<EmptyState title={t('services.unassessed')} description={t('services.unassessedBody')} />
				{/if}
			</div>
			{#if isTemaPendidikan(item.kode)}
				<DataLayanan {kodeKecamatan} kodeTema={item.kode} />
			{/if}
			<div class="card flex flex-col gap-5 p-4 md:p-6">
				{#if item.indikator && kondisi !== 'belum-dinilai'}
					<SeatSummary indikator={item.indikator} tema={item.kode} />
					<div class="rounded-md bg-sunken p-4">
						<h4 class="text-small font-semibold">{t('services.impactTitle')}</h4>
						<p class="mt-2 text-small text-muted">
							{t(kondisi === 'perhatian' ? 'services.impactAttention' : 'services.impactGood')}
						</p>
					</div>
					{#if kondisi === 'perhatian'}
						<p
							class="self-start rounded-xs bg-attn-bg px-3 py-1 text-caption font-semibold text-attn"
						>
							{t('services.issue')}
						</p>
					{/if}
				{/if}
				{#if isTemaPendidikan(item.kode)}
					{@const jenjang = t(teksJenjang[item.kode].jenjang)}
					<div class="flex flex-col gap-2 border-t border-line pt-4 md:flex-row md:flex-wrap">
						{#if kondisi === 'perhatian'}
							<a href={tautanAspirasi(item.kode)} class="btn btn-primary max-md:btn-block">
								<FileText size={18} aria-hidden="true" />{t('services.actOnIssue')}
							</a>
						{/if}
						<a
							href="/{bahasaAktif}/kecamatan/{kodeKecamatan}/{item.kode}"
							class="btn btn-outline max-md:btn-block"
						>
							{t('tema.detailJenjang', { jenjang })}
							<ArrowRight size={18} strokeWidth={1.75} aria-hidden="true" />
						</a>
						<a href={tautanCekSekolah(item.kode)} class="btn btn-outline max-md:btn-block">
							<LocateFixed size={18} strokeWidth={1.75} aria-hidden="true" />
							{t('services.checkMySchool')}<span class="sr-only">: {jenjang}</span>
						</a>
					</div>
				{/if}
			</div>
		</section>
	{/each}
	{#if temaTampil.length}
		<PenanggungJawabLayanan kodeTema={jenjangAktif} />
		<div class="rounded-lg border border-line p-4 md:p-5">
			<h3 class="t-h3">{t('services.budgetTitle')}</h3>
			<p class="mt-2 max-w-prose text-small text-muted">{t('services.budgetBody')}</p>
			<a
				href="/{bahasaAktif}/kecamatan/{kodeKecamatan}?bagian=anggaran&tema={jenjangAktif}"
				class="btn btn-outline mt-4"
				><Wallet size={18} aria-hidden="true" />{t('services.budgetAction')}</a
			>
		</div>
		<details class="border-t border-line pt-4">
			<summary class="min-h-touch cursor-pointer py-3 font-heading text-h3 font-semibold"
				>{t('services.oversightNews')}</summary
			>
			<div class="mt-4"><SuaraPengawas {kodeKecamatan} kodeTema={jenjangAktif} /></div>
		</details>
	{/if}
</div>
