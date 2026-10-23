<script lang="ts">
	import type { Anggaran, TemaKecamatan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import { t } from '$lib/i18n';
	import DaftarPaket from './DaftarPaket.svelte';

	interface Props {
		anggaran: Anggaran | null;
		kodeKecamatan: string;
		tema: TemaKecamatan[];
	}

	let { anggaran, kodeKecamatan, tema }: Props = $props();
</script>

{#if anggaran}
	<section aria-labelledby="judul-anggaran" class="flex flex-col gap-4">
		<h3 id="judul-anggaran" class="t-h3">{t('money.yearTitle', { year: anggaran.tahun })}</h3>
		<div class="grid grid-cols-2 gap-3 lg:gap-5">
			<KeyNumberCard label={t('money.count')} angka={anggaran.jumlah_paket} size="lg" />
			<KeyNumberCard
				label={t('money.total', { year: anggaran.tahun })}
				angka={anggaran.total_hps}
				size="lg"
			/>
		</div>
		<p class="text-small text-muted">{t('money.locationNote')}</p>
	</section>
{:else}
	<EmptyState title={t('money.empty')} />
{/if}

<div class="mt-10">
	<DaftarPaket {kodeKecamatan} {tema} tahunAwal={anggaran?.tahun ?? new Date().getFullYear()} />
</div>
