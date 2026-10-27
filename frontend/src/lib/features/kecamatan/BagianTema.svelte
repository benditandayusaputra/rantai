<script lang="ts">
	import { ArrowRight, Calendar, FileText, GraduationCap, LocateFixed } from '@lucide/svelte';
	import type { TemaKecamatan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import GapClassChip from '$lib/components/GapClassChip.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import SeatSummary from '$lib/components/SeatSummary.svelte';
	import { isTemaPendidikan, teksJenjang } from '$lib/features/pendidikan/jenjang';
	import { bahasaAktif, t } from '$lib/i18n';
	import SuaraPengawas from './SuaraPengawas.svelte';

	let { tema, kodeKecamatan }: { tema: TemaKecamatan[]; kodeKecamatan: string } = $props();

	const temaPendidikanPertama = $derived(tema.find((item) => isTemaPendidikan(item.kode))?.kode);

	const tautanCekSekolah = (kode: string) =>
		`/${bahasaAktif}/cek-sekolah?w=${kodeKecamatan}&tema=${kode}`;
	const tautanAspirasi = (kode: string) =>
		`/${bahasaAktif}/aspirasi?w=${kodeKecamatan}&tema=${kode}&masalah=${kode === 'pendidikan' ? 'kursi_smp_kurang' : 'kursi_sma_smk_kurang'}`;
</script>

<div class="flex flex-col gap-10">
	{#each tema as item (item.kode)}
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
				{#if item.indikator}
					<div class="flex flex-wrap items-center gap-3">
						<GapClassChip kelas={item.indikator.kelas} />
						<span class="inline-flex items-center gap-1 text-small text-muted">
							<Calendar size={16} strokeWidth={1.75} aria-hidden="true" />
							{t('tema.period', { period: item.indikator.periode })}
						</span>
					</div>
					<SeatSummary indikator={item.indikator} tema={item.kode} />
					<div class="grid gap-3 xs:grid-cols-2 lg:grid-cols-4 lg:gap-4">
						<KeyNumberCard label={t('tema.needs')} angka={item.indikator.kebutuhan} />
						<KeyNumberCard label={t('tema.capacity')} angka={item.indikator.kapasitas} />
						<KeyNumberCard label={t('tema.gap')} angka={item.indikator.kesenjangan} />
						{#if item.indikator.rasio}
							<KeyNumberCard label={t('tema.ratio')} angka={item.indikator.rasio} />
						{/if}
					</div>
				{:else}
					<EmptyState title={t('tema.empty')} />
				{/if}
				{#if isTemaPendidikan(item.kode)}
					{@const jenjang = t(teksJenjang[item.kode].jenjang)}
					<div class="flex flex-col gap-2 border-t border-line pt-4 md:flex-row md:flex-wrap">
						{#if item.indikator && item.indikator.kesenjangan.nilai > 0}
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
	{#if temaPendidikanPertama}
		<details class="border-t border-line pt-4">
			<summary class="min-h-touch cursor-pointer py-3 font-heading text-h3 font-semibold"
				>{t('services.oversightNews')}</summary
			>
			<div class="mt-4"><SuaraPengawas {kodeKecamatan} kodeTema={temaPendidikanPertama} /></div>
		</details>
	{/if}
</div>
