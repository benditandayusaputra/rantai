<script lang="ts">
	import type { Porsi } from '$lib/api/tipe';
	import { persenDari } from './dasbor';

	interface Props {
		isi: Porsi[];
		nilaiDari: (p: Porsi) => number;
		format: (n: number) => string;
		namaDari: (p: Porsi) => string;
		rincian?: (p: Porsi) => string;
	}

	let { isi, nilaiDari, format, namaDari, rincian }: Props = $props();

	const maksimum = $derived(Math.max(0, ...isi.map(nilaiDari)));
</script>

<ol class="flex flex-col gap-4">
	{#each isi as p (p.kunci)}
		<li class="flex flex-col gap-1">
			<div class="flex items-baseline justify-between gap-3 text-small">
				<span class="min-w-0 break-words font-semibold">{namaDari(p)}</span>
				<span class="num shrink-0 font-bold">{format(nilaiDari(p))}</span>
			</div>
			<div class="h-3 rounded-sm bg-sunken" aria-hidden="true">
				<div
					class="h-full rounded-sm bg-primary"
					style:width="{persenDari(nilaiDari(p), maksimum)}%"
				></div>
			</div>
			{#if rincian}<p class="num text-caption text-muted">{rincian(p)}</p>{/if}
		</li>
	{/each}
</ol>
