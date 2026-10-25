<script lang="ts">
	import { Info, Wallet } from '@lucide/svelte';
	import type { Anggaran, TemaKecamatan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import KeyNumberCard from '$lib/components/KeyNumberCard.svelte';
	import { t } from '$lib/i18n';
	import DaftarPaket from './DaftarPaket.svelte';
	import GrafikBelanja from './GrafikBelanja.svelte';

	interface Props {
		anggaran: Anggaran | null;
		kodeKecamatan: string;
		tema: TemaKecamatan[];
	}

	let { anggaran, kodeKecamatan, tema }: Props = $props();
</script>

{#if anggaran}
	<section aria-labelledby="judul-anggaran" class="card flex flex-col gap-5 p-4 md:p-6">
		<div class="flex items-center gap-3">
			<span class="icon-tile"><Wallet size={20} strokeWidth={1.75} aria-hidden="true" /></span>
			<h3 id="judul-anggaran" class="t-h3">{t('money.yearTitle', { year: anggaran.tahun })}</h3>
		</div>
		<div class="grid grid-cols-2 gap-3 lg:gap-4">
			<KeyNumberCard label={t('money.count')} angka={anggaran.jumlah_paket} size="lg" />
			<KeyNumberCard
				label={t('money.total', { year: anggaran.tahun })}
				angka={anggaran.total_hps}
				size="lg"
			/>
		</div>
		<p class="flex gap-2 text-small text-muted">
			<Info size={18} strokeWidth={1.75} aria-hidden="true" class="mt-px shrink-0" />
			{t('money.locationNote')}
		</p>
	</section>
{:else}
	<EmptyState title={t('money.empty')} />
{/if}

{#if tema.some((item) => item.kode === 'pendidikan')}
	<div class="mt-10">
		<GrafikBelanja {kodeKecamatan} kodeTema="pendidikan" />
	</div>
{/if}

<div class="mt-10">
	<DaftarPaket {kodeKecamatan} {tema} tahunAwal={anggaran?.tahun ?? new Date().getFullYear()} />
</div>
