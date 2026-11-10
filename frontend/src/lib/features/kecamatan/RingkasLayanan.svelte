<script lang="ts">
	import { GraduationCap } from '@lucide/svelte';
	import type { TemaKecamatan } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import SeatSummary from '$lib/components/SeatSummary.svelte';
	import { bahasaAktif, t } from '$lib/i18n';
	import { statusLayanan } from './kecamatan';

	let { tema }: { tema: TemaKecamatan[] } = $props();
</script>

{#if tema.length === 0}
	<EmptyState title={t('services.unassessed')} description={t('services.unassessedBody')} />
{:else}
	<ul class="grid gap-3 md:grid-cols-2 md:gap-4">
		{#each tema as item (item.kode)}
			<li class="card flex h-full flex-col gap-4 p-4 md:p-5">
				<div class="flex items-start gap-3">
					<span class="icon-tile"
						><GraduationCap size={20} strokeWidth={1.75} aria-hidden="true" /></span
					>
					<div class="flex min-w-0 flex-col gap-1">
						<p class="t-over">{item.nama[bahasaAktif]}</p>
						<h3 class="t-h3">{item.judul[bahasaAktif]}</h3>
					</div>
				</div>
				{#if item.indikator && statusLayanan(item.indikator) !== 'belum-dinilai'}
					<SeatSummary indikator={item.indikator} tema={item.kode} />
				{:else}
					<p class="text-small text-muted">{t('services.unassessedBody')}</p>
				{/if}
			</li>
		{/each}
	</ul>
{/if}
