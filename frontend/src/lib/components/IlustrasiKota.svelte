<script lang="ts">
	let { class: kelas = 'w-full' }: { class?: string } = $props();

	const gedung = [
		{ x: 112, y: 246, l: 40, t: 144 },
		{ x: 156, y: 292, l: 30, t: 98 },
		{ x: 190, y: 318, l: 26, t: 72 },
		{ x: 324, y: 226, l: 44, t: 164 },
		{ x: 372, y: 280, l: 34, t: 110 },
		{ x: 410, y: 310, l: 28, t: 80 }
	];
	const jendela = [0, 1, 2, 3, 4].flatMap((baris) =>
		[0, 1].map((kolom) => ({ x: 334 + kolom * 16, y: 242 + baris * 22 }))
	);
	const batang = [
		{ x: 386, tinggi: 22, warna: 'fill-gap-3' },
		{ x: 402, tinggi: 38, warna: 'fill-gap-2' },
		{ x: 418, tinggi: 16, warna: 'fill-gap-4' },
		{ x: 434, tinggi: 30, warna: 'fill-gap-1' }
	];
</script>

<svg viewBox="0 0 520 440" class={kelas} fill="none" aria-hidden="true">
	<circle cx="270" cy="220" r="190" class="fill-secondary" />
	<ellipse cx="270" cy="392" rx="210" ry="12" class="fill-secondary-hover" />

	<g class="fill-secondary-hover">
		{#each gedung as g (g.x)}
			<rect x={g.x} y={g.y} width={g.l} height={g.t} rx="3" />
		{/each}
		<polygon points="216,392 324,392 306,366 234,366" />
		<polygon points="262,366 278,366 275,148 265,148" />
		<polygon points="255,148 285,148 280,137 260,137" />
	</g>
	<g class="fill-secondary">
		{#each jendela as j (`${j.x}-${j.y}`)}
			<rect x={j.x} y={j.y} width="8" height="10" rx="1.5" />
		{/each}
	</g>
	<path d="M270 137c-11-10-6-23 0-32 6 9 11 22 0 32z" class="fill-attn-fill" />

	<path
		d="M176 134c50-46 120-56 162-36"
		stroke-width="2.5"
		stroke-dasharray="4 7"
		stroke-linecap="round"
		class="stroke-line-strong"
	/>
	<path d="M150 204c-22-26-32-42-32-58a32 32 0 0 1 64 0c0 16-10 32-32 58z" class="fill-primary" />
	<circle cx="150" cy="146" r="12" class="fill-raised" />

	<rect x="336" y="62" width="156" height="100" rx="16" class="fill-line" />
	<rect x="330" y="56" width="156" height="100" rx="16" class="fill-raised stroke-line" />
	<rect x="344" y="70" width="30" height="30" rx="8" class="fill-secondary" />
	<path d="M359 77l12 6-12 6-12-6z" class="fill-primary" />
	<path d="M352 86v6c4 3 10 3 14 0v-6" stroke-width="2.5" class="stroke-primary" />
	<rect x="344" y="114" width="34" height="7" rx="3.5" class="fill-line" />
	<rect x="344" y="128" width="24" height="7" rx="3.5" class="fill-line" />
	{#each batang as b (b.x)}
		<rect x={b.x} y={142 - b.tinggi} width="11" height={b.tinggi} rx="3" class={b.warna} />
	{/each}

	<rect x="40" y="256" width="146" height="82" rx="16" class="fill-line" />
	<rect x="34" y="250" width="146" height="82" rx="16" class="fill-raised stroke-line" />
	<circle cx="66" cy="291" r="17" class="fill-attn-fill" />
	<circle cx="66" cy="291" r="10" stroke-width="2.5" class="stroke-raised" />
	<rect x="94" y="278" width="68" height="8" rx="4" class="fill-line" />
	<rect x="94" y="294" width="44" height="8" rx="4" class="fill-primary" />
	<g stroke-width="3" class="stroke-primary">
		<rect x="94" y="310" width="20" height="12" rx="6" />
		<rect x="105" y="310" width="20" height="12" rx="6" />
	</g>

	<rect x="372" y="306" width="128" height="72" rx="16" class="fill-line" />
	<rect x="366" y="300" width="128" height="72" rx="16" class="fill-raised stroke-line" />
	<path d="M381 324l15-9 15 9z" class="fill-primary" />
	<g class="fill-primary">
		<rect x="384" y="327" width="4" height="14" rx="1" />
		<rect x="394" y="327" width="4" height="14" rx="1" />
		<rect x="404" y="327" width="4" height="14" rx="1" />
		<rect x="380" y="343" width="32" height="4" rx="1" />
	</g>
	<rect x="422" y="320" width="58" height="7" rx="3.5" class="fill-line" />
	<rect x="422" y="334" width="40" height="7" rx="3.5" class="fill-line" />
	<rect x="380" y="354" width="56" height="8" rx="4" class="fill-success-bg" />

	<path
		d="M226 70l4 10 10 4-10 4-4 10-4-10-10-4 10-4z M480 214l3 7 7 3-7 3-3 7-3-7-7-3 7-3z M70 196l3 6 6 3-6 3-3 6-3-6-6-3 6-3z"
		class="fill-attn-fill"
	/>
</svg>
