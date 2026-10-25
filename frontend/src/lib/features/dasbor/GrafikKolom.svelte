<script lang="ts">
	import type { Porsi } from '$lib/api/tipe';
	import { batasRapi, persenDari } from './dasbor';

	interface Props {
		isi: Porsi[];
		nilaiDari: (p: Porsi) => number;
		format: (n: number) => string;
		labelDari: (p: Porsi) => string;
		rincian?: (p: Porsi) => string;
		sorot?: string;
	}

	let { isi, nilaiDari, format, labelDari, rincian, sorot }: Props = $props();

	interface Kolom {
		porsi: Porsi;
		nilai: number;
		tinggi: number;
		posisi: string;
		isPuncak: boolean;
	}

	let disorot = $state<Kolom | null>(null);

	const garisBantu = [0.5, 1];

	function posisiLabel(i: number): string {
		if (i < isi.length / 3) return 'left-0';
		if (i >= (isi.length * 2) / 3) return 'right-0';
		return 'left-1/2 -translate-x-1/2';
	}

	const kolom = $derived.by((): Kolom[] => {
		const nilai = isi.map(nilaiDari);
		const tertinggi = Math.max(0, ...nilai);
		const maksimum = batasRapi(tertinggi);
		const puncak = nilai.indexOf(tertinggi);
		return isi.map((porsi, i) => {
			const n = nilaiDari(porsi);
			return {
				porsi,
				nilai: n,
				tinggi: persenDari(n, maksimum),
				posisi: posisiLabel(i),
				isPuncak: i === puncak && n > 0
			};
		});
	});
</script>

<figure class="pt-6">
	<div class="relative h-48 md:h-60">
		{#each garisBantu as g (g)}
			<div
				class="absolute inset-x-0 border-t border-dashed border-line"
				style:bottom="{g * 100}%"
				aria-hidden="true"
			></div>
		{/each}
		<ol class="absolute inset-0 flex items-end gap-1 border-b-ctrl border-line-strong md:gap-2">
			{#each kolom as k (k.porsi.kunci)}
				<li
					class="relative flex h-full min-w-0 flex-1 items-end justify-center"
					onpointerenter={() => (disorot = k)}
					onpointerleave={() => (disorot = null)}
				>
					<div
						class="w-full max-w-14 rounded-t-xs bg-chart transition-opacity duration-fast ease-standard {sorot &&
						sorot !== k.porsi.kunci
							? 'opacity-40'
							: ''} {disorot === k ? 'opacity-80' : ''}"
						class:min-h-0.5={k.nilai > 0}
						style:height="{k.tinggi}%"
					></div>
					{#if k.isPuncak && disorot === null}
						<span
							class="num absolute whitespace-nowrap pb-1 text-caption font-bold {k.posisi}"
							style:bottom="{k.tinggi}%"
							aria-hidden="true">{format(k.nilai)}</span
						>
					{/if}
					<span class="sr-only">{labelDari(k.porsi)}: {format(k.nilai)}</span>
				</li>
			{/each}
		</ol>
		{#if disorot}
			<div
				class="pointer-events-none absolute z-tooltip whitespace-nowrap rounded-md bg-raised px-3 py-2 text-small shadow-2 {disorot.posisi}"
				style:bottom="calc({disorot.tinggi}% + 8px)"
				aria-hidden="true"
			>
				<p class="font-bold">{labelDari(disorot.porsi)}</p>
				<p class="num">{format(disorot.nilai)}</p>
				{#if rincian}<p class="num text-muted">{rincian(disorot.porsi)}</p>{/if}
			</div>
		{/if}
	</div>
	<div class="mt-2 flex gap-1 md:gap-2" aria-hidden="true">
		{#each kolom as k (k.porsi.kunci)}
			<span
				class="min-w-0 flex-1 whitespace-nowrap text-center text-caption text-muted max-md:even:invisible {sorot ===
				k.porsi.kunci
					? 'font-bold text-fg'
					: ''}">{labelDari(k.porsi)}</span
			>
		{/each}
	</div>
</figure>
