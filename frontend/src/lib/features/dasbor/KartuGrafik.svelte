<script lang="ts">
	import type { ChartColumn } from '@lucide/svelte';
	import type { Snippet } from 'svelte';
	import type { Angka } from '$lib/api/tipe';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import TableToggle from '$lib/components/TableToggle.svelte';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import { t } from '$lib/i18n';

	interface Props {
		judul: string;
		keterangan?: string;
		ikon: typeof ChartColumn;
		angka: Angka | null;
		kolom: string[];
		baris: string[][];
		class?: string;
		aksi?: Snippet;
		catatan?: Snippet;
		children: Snippet;
	}

	let {
		judul,
		keterangan,
		ikon: Ikon,
		angka,
		kolom,
		baris,
		class: kelas = '',
		aksi,
		catatan,
		children
	}: Props = $props();

	const id = $props.id();
	let isTabel = $state(false);
	const isKosong = $derived(baris.length === 0);
</script>

<section aria-labelledby="{id}-judul" class="card flex min-w-0 flex-col gap-4 p-4 md:p-6 {kelas}">
	<div class="flex flex-wrap items-start justify-between gap-3">
		<div class="flex min-w-0 items-start gap-3">
			<span class="icon-tile"><Ikon size={20} strokeWidth={1.75} aria-hidden="true" /></span>
			<div class="min-w-0">
				<h3 id="{id}-judul" class="t-h3">{judul}</h3>
			</div>
		</div>
		{#if !isKosong}<TableToggle bind:pressed={isTabel} controls={id} />{/if}
	</div>
	{#if aksi && !isKosong && !isTabel}{@render aksi()}{/if}
	<div {id} class="min-w-0">
		{#if isKosong}
			<p
				class="rounded-md border-ctrl border-dashed border-line-strong bg-sunken p-6 text-center text-small text-muted"
			>
				{t('dash.empty')}
			</p>
		{:else if isTabel}
			<div class="max-h-96 overflow-auto rounded-md border border-line">
				<table class="w-full text-small">
					<caption class="sr-only">{judul}</caption>
					<thead class="sticky top-0 bg-sunken text-caption font-semibold">
						<tr>
							{#each kolom as judulKolom, i (judulKolom)}
								<th scope="col" class="p-3 {i === 0 ? 'text-left' : 'text-right'}">{judulKolom}</th>
							{/each}
						</tr>
					</thead>
					<tbody>
						{#each baris as [kepala, ...sel], i (i)}
							<tr class="border-b border-line">
								<th scope="row" class="p-3 text-left font-semibold">{kepala}</th>
								{#each sel as isi, j (j)}<td class="num p-3 text-right">{isi}</td>{/each}
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{:else}
			{@render children()}
		{/if}
	</div>
	{#if keterangan}<p class="max-w-prose text-small text-muted">{keterangan}</p>{/if}
	{#if catatan}{@render catatan()}{/if}
	{#if angka && !isKosong}
		<div><ProvenanceBadge {angka} nilai={tampilkanAngka(angka).nilai} /></div>
	{/if}
</section>
