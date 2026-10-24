<script lang="ts">
	import { Info } from '@lucide/svelte';
	import type { TemaKecamatan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import GapClassChip from '$lib/components/GapClassChip.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import { bahasaBawaan, t } from '$lib/i18n';

	let { tema, kodeKecamatan }: { tema: TemaKecamatan[]; kodeKecamatan: string } = $props();
</script>

<div class="flex flex-col gap-10">
	{#each tema as item (item.kode)}
		<section aria-labelledby="tema-{item.kode}" class="flex flex-col gap-4">
			<p class="t-over">{item.nama[bahasaBawaan]}</p>
			<h3 id="tema-{item.kode}" class="t-h3">{item.judul[bahasaBawaan]}</h3>
			{#if item.indikator}
				<div class="flex flex-wrap items-center gap-3">
					<GapClassChip kelas={item.indikator.kelas} />
					<span class="text-small text-muted"
						>{t('tema.period', { period: item.indikator.periode })}</span
					>
				</div>
				<div class="grid grid-cols-2 gap-3 lg:grid-cols-4 lg:gap-5">
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
				<a
					href="/{bahasaBawaan}/kecamatan/{kodeKecamatan}/pendidikan"
					class="inline-flex min-h-touch items-center self-start font-bold text-secondary-on underline"
					>{t('tema.detail')}</a
				>
				<a
					href="/{bahasaBawaan}/cek-sekolah"
					class="inline-flex min-h-touch items-center self-start font-bold text-secondary-on underline"
					>{t('services.checkMySchool')}</a
				>
			{/if}
		</section>
	{/each}
</div>
