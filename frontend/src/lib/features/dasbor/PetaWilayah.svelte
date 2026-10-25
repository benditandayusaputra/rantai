<script lang="ts">
	import type { Snippet } from 'svelte';
	import { bahasaAktif, t } from '$lib/i18n';
	import { susunPeta, type AreaPeta, type FiturWilayah, type IsiArea } from './dasbor';

	interface Props {
		fitur: FiturWilayah[];
		isiDari: (kode: string) => IsiArea;
		judul: string;
		legenda: Snippet;
	}

	let { fitur, isiDari, judul, legenda }: Props = $props();

	const susunan = $derived(susunPeta(fitur));
	let disorot = $state<AreaPeta | null>(null);
</script>

<figure class="flex flex-col gap-4">
	<div class="relative">
		<svg
			viewBox="0 0 {susunan.lebar} {susunan.tinggi}"
			class="block h-auto w-full"
			role="group"
			aria-label={judul}
		>
			{#if susunan.sisipan}
				{@const s = susunan.sisipan}
				<rect
					x={s.x + 0.5}
					y={s.y + 0.5}
					width={s.lebar - 1}
					height={s.tinggi - 1}
					rx="6"
					class="fill-none stroke-line"
					stroke-dasharray="4 4"
				/>
			{/if}
			{#each susunan.area as area (area.kode)}
				{@const isi = isiDari(area.kode)}
				<a
					href="/{bahasaAktif}/kecamatan/{area.kode}"
					aria-label="{area.nama}: {isi.keterangan}"
					onpointerenter={() => (disorot = area)}
					onpointerleave={() => (disorot = null)}
					onfocus={() => (disorot = area)}
					onblur={() => (disorot = null)}
					class="outline-none"
				>
					<path
						d={area.d}
						class="stroke-surface transition-opacity duration-fast ease-standard {isi.kelas}"
						fill-opacity={isi.opasitas}
						stroke-width="1.5"
						vector-effect="non-scaling-stroke"
					/>
				</a>
			{/each}
			{#if disorot}
				<path
					d={disorot.d}
					class="pointer-events-none fill-none stroke-fg"
					stroke-width="2.5"
					vector-effect="non-scaling-stroke"
				/>
			{/if}
		</svg>
		{#if susunan.sisipan}
			<span
				class="absolute -translate-x-1/2 pt-1 text-caption text-muted"
				style:left="{((susunan.sisipan.x + susunan.sisipan.lebar / 2) / susunan.lebar) * 100}%"
				style:top="{(susunan.sisipan.tinggi / susunan.tinggi) * 100}%">{t('dash.mapInset')}</span
			>
		{/if}
		{#if disorot}
			<div
				class="pointer-events-none absolute z-tooltip max-w-60 -translate-x-1/2 -translate-y-full rounded-md bg-raised px-3 py-2 text-small shadow-2"
				style:left="{(disorot.tengah[0] / susunan.lebar) * 100}%"
				style:top="calc({(disorot.tengah[1] / susunan.tinggi) * 100}% - 8px)"
				aria-hidden="true"
			>
				<p class="font-bold">{disorot.nama}</p>
				<p class="num">{isiDari(disorot.kode).keterangan}</p>
			</div>
		{/if}
	</div>
	<figcaption>{@render legenda()}</figcaption>
</figure>
