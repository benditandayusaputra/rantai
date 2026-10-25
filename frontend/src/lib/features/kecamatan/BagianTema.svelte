<script lang="ts">
	import { ArrowRight, Calendar, GraduationCap, Info, LocateFixed } from '@lucide/svelte';
	import type { TemaKecamatan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import GapClassChip from '$lib/components/GapClassChip.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import { bahasaBawaan, t } from '$lib/i18n';
	import SuaraPengawas from './SuaraPengawas.svelte';

	let { tema, kodeKecamatan }: { tema: TemaKecamatan[]; kodeKecamatan: string } = $props();
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
						<p class="t-over">{item.nama[bahasaBawaan]}</p>
						<h3 id="tema-{item.kode}" class="t-h3">{item.judul[bahasaBawaan]}</h3>
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
					<div class="grid grid-cols-2 gap-3 lg:grid-cols-4 lg:gap-4">
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
				{#if item.kode === 'pendidikan'}
					<p class="flex gap-2 rounded-md bg-info-bg p-4 text-small text-fg">
						<Info size={20} strokeWidth={1.75} aria-hidden="true" class="shrink-0 text-info" />
						{t('services.zoningNote')}
					</p>
					<div class="flex flex-col gap-2 md:flex-row md:flex-wrap">
						<a
							href="/{bahasaBawaan}/kecamatan/{kodeKecamatan}/pendidikan"
							class="btn btn-primary max-md:btn-block"
						>
							{t('tema.detail')}
							<ArrowRight size={18} strokeWidth={1.75} aria-hidden="true" />
						</a>
						<a href="/{bahasaBawaan}/cek-sekolah" class="btn btn-outline max-md:btn-block">
							<LocateFixed size={18} strokeWidth={1.75} aria-hidden="true" />
							{t('services.checkMySchool')}
						</a>
					</div>
				{/if}
			</div>
			<SuaraPengawas {kodeKecamatan} kodeTema={item.kode} />
		</section>
	{/each}
</div>
