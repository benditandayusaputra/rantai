<script lang="ts">
	import { ChartColumn } from '@lucide/svelte';
	import type { DetailPendidikan } from '$lib/api/tipe';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import TableToggle from '$lib/components/TableToggle.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { t } from '$lib/i18n';
	import { persenLebar, susunBarisGrafik } from './pendidikan';

	let { detail }: { detail: DetailPendidikan } = $props();

	const id = $props.id();
	let isTabel = $state(false);

	const baris = $derived(
		susunBarisGrafik(detail).map((b) => ({ ...b, tampil: tampilkanAngka(b.angka) }))
	);
	const maksimum = $derived(Math.max(0, ...baris.map((b) => b.angka.nilai)));
</script>

<section aria-labelledby="{id}-judul" class="card p-4 md:p-6">
	<div class="flex flex-wrap items-start justify-between gap-3">
		<div class="flex items-center gap-3">
			<span class="icon-tile"><ChartColumn size={20} strokeWidth={1.75} aria-hidden="true" /></span>
			<h2 id="{id}-judul" class="t-h3">{t('pen.chartTitle')}</h2>
		</div>
		<TableToggle bind:pressed={isTabel} controls={id} />
	</div>
	<div {id} class="mt-6">
		{#if isTabel}
			<table class="w-full text-small">
				<caption class="sr-only">{t('pen.chartTitle')}</caption>
				<thead class="bg-sunken text-caption uppercase">
					<tr>
						<th scope="col" class="p-3 text-left">{t('pen.school.type')}</th>
						<th scope="col" class="p-3 text-right">{t('pen.value')}</th>
					</tr>
				</thead>
				<tbody>
					{#each baris as item (item.kunci)}
						<tr class="border-b border-line">
							<th scope="row" class="p-3 text-left font-normal">{t(item.kunci)}</th>
							<td class="num p-3 text-right">
								{item.tampil.nilai}
								{item.tampil.satuan}
								<ProvenanceBadge angka={item.angka} nilai={item.tampil.nilai} compact />
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{:else}
			<ul class="flex flex-col gap-4">
				{#each baris as item (item.kunci)}
					<li class="flex flex-col gap-1">
						<div class="flex flex-wrap items-center justify-between gap-2">
							<span class="text-small font-semibold">{t(item.kunci)}</span>
							<ProvenanceBadge angka={item.angka} nilai={item.tampil.nilai} compact />
						</div>
						<div class="flex items-center gap-3">
							<div class="h-6 flex-1 rounded-sm bg-sunken" aria-hidden="true">
								<div
									class="h-full rounded-sm {item.kunci === 'pen.graduates'
										? 'bg-fg'
										: item.angka.label === 'resmi'
											? 'bg-primary'
											: 'border-ctrl border-dashed border-primary bg-secondary'}"
									style="width: {persenLebar(item.angka.nilai, maksimum)}%"
								></div>
							</div>
							<span class="num w-24 text-right font-bold"
								>{item.tampil.nilai} <span class="sr-only">{item.tampil.satuan}</span></span
							>
						</div>
					</li>
				{/each}
			</ul>
		{/if}
	</div>
</section>
