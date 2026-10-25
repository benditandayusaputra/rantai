<script lang="ts">
	import { ChartLine } from '@lucide/svelte';
	import type { TitikTren } from '$lib/api/tipe';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import TableToggle from '$lib/components/TableToggle.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { t } from '$lib/i18n';
	import { skalaTren, ukuranGrafikTren } from './pendidikan';

	let { tren }: { tren: TitikTren[] } = $props();

	const id = $props.id();
	let isTabel = $state(false);

	const { lebar, tinggi, kiri, kanan } = ukuranGrafikTren;
	const skala = $derived(skalaTren(tren, (item) => item.rasio.nilai));
	const titik = $derived(
		skala.titik.map(({ item, x, y }) => ({ ...item, x, y, tampil: tampilkanAngka(item.rasio) }))
	);
	const garis = $derived(titik.map((p) => `${p.x},${p.y}`).join(' '));
</script>

<section aria-labelledby="{id}-judul" class="card p-4 md:p-6">
	<div class="flex flex-wrap items-start justify-between gap-3">
		<div class="flex items-center gap-3">
			<span class="icon-tile"><ChartLine size={20} strokeWidth={1.75} aria-hidden="true" /></span>
			<h2 id="{id}-judul" class="t-h3">{t('pen.trend.title')}</h2>
		</div>
		<TableToggle
			bind:pressed={isTabel}
			controls={id}
			labelKembali={t('action.viewAsChart')}
			ikonKembali={ChartLine}
		/>
	</div>
	<p class="mt-2 max-w-prose text-small text-muted">{t('pen.trend.lead')}</p>
	<div {id} class="mt-4">
		{#if isTabel}
			<table class="w-full text-small">
				<caption class="sr-only">{t('pen.trend.title')}</caption>
				<thead class="bg-sunken text-caption font-semibold">
					<tr>
						<th scope="col" class="p-3 text-left">{t('pen.trend.period')}</th>
						<th scope="col" class="p-3 text-right">{t('pen.trend.ratio')}</th>
					</tr>
				</thead>
				<tbody>
					{#each titik as item (item.periode)}
						<tr class="border-b border-line">
							<th scope="row" class="p-3 text-left font-normal">{item.periode}</th>
							<td class="num p-3 text-right">
								{item.tampil.nilai}
								<ProvenanceBadge angka={item.rasio} nilai={item.tampil.nilai} compact />
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		{:else}
			<svg viewBox="0 0 {lebar} {tinggi}" class="w-full text-caption" aria-hidden="true">
				<line
					x1={kiri}
					x2={lebar - kanan}
					y1={skala.ySeimbang}
					y2={skala.ySeimbang}
					stroke-width="1.5"
					stroke-dasharray="6 6"
					class="stroke-line-strong"
				/>
				<text x={lebar - kanan} y={skala.ySeimbang - 6} text-anchor="end" class="fill-muted"
					>{t('pen.trend.parity')}</text
				>
				<polyline
					points={garis}
					fill="none"
					stroke-width="3"
					stroke-linejoin="round"
					class="stroke-chart"
				/>
				{#each titik as item (item.periode)}
					<circle cx={item.x} cy={item.y} r="6" stroke-width="3" class="fill-raised stroke-chart" />
					<text x={item.x} y={item.y - 14} text-anchor="middle" class="num fill-fg font-bold"
						>{item.tampil.nilai}</text
					>
					<text x={item.x} y={tinggi - 14} text-anchor="middle" class="fill-muted"
						>{item.periode}</text
					>
				{/each}
			</svg>
			<ul class="sr-only">
				{#each titik as item (item.periode)}
					<li>{item.periode}: {item.tampil.nilai} {item.tampil.satuan}</li>
				{/each}
			</ul>
		{/if}
	</div>
</section>
