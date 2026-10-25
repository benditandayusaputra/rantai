<script lang="ts">
	let { class: kelas = 'w-full' }: { class?: string } = $props();

	const batang = [
		{ x: 322, tinggi: 74, warna: 'fill-gap-2' },
		{ x: 346, tinggi: 46, warna: 'fill-gap-3' },
		{ x: 370, tinggi: 96, warna: 'fill-gap-1' },
		{ x: 394, tinggi: 60, warna: 'fill-gap-2' }
	];
	const gedung = [
		{ x: 128, y: 286, l: 30, t: 98 },
		{ x: 162, y: 318, l: 24, t: 66 },
		{ x: 300, y: 300, l: 34, t: 84 },
		{ x: 338, y: 262, l: 30, t: 122 },
		{ x: 372, y: 312, l: 28, t: 72 }
	];
</script>

<svg viewBox="0 0 480 440" class={kelas} fill="none" aria-hidden="true">
	<circle cx="250" cy="226" r="186" class="fill-secondary" />
	<ellipse cx="250" cy="388" rx="200" ry="10" class="fill-secondary-hover" />

	<g class="fill-secondary-hover">
		{#each gedung as g (g.x)}
			<rect x={g.x} y={g.y} width={g.l} height={g.t} rx="3" />
		{/each}
		<polygon points="62,386 132,386 120,366 74,366" />
		<polygon points="91,366 103,366 101,214 93,214" />
		<polygon points="86,214 108,214 104,205 90,205" />
	</g>
	<path d="M97 205c-8-8-4-18 0-25 4 7 8 17 0 25z" class="fill-attn-fill" />

	<rect x="306" y="102" width="150" height="196" rx="18" class="fill-line" />
	<rect x="298" y="94" width="150" height="196" rx="18" class="fill-raised stroke-line" />
	<rect x="314" y="112" width="20" height="14" rx="7" stroke-width="3.5" class="stroke-primary" />
	<rect x="326" y="112" width="20" height="14" rx="7" stroke-width="3.5" class="stroke-fg" />
	<rect x="356" y="114" width="72" height="9" rx="4.5" class="fill-line" />
	<rect x="314" y="138" width="56" height="7" rx="3.5" class="fill-sunken" />
	{#each batang as b (b.x)}
		<rect x={b.x} y={250 - b.tinggi} width="18" height={b.tinggi} rx="4" class={b.warna} />
	{/each}
	<rect x="314" y="250" width="118" height="2" rx="1" class="fill-line-strong" />
	<rect x="314" y="262" width="118" height="18" rx="9" class="fill-success-bg" />
	<circle cx="325" cy="271" r="5.5" class="fill-success" />
	<rect x="336" y="268" width="80" height="6" rx="3" class="fill-success" />

	<circle cx="126" cy="122" r="28" class="fill-raised stroke-line" />
	<path
		d="M114 122l8 8 16-16"
		stroke-width="4"
		stroke-linecap="round"
		stroke-linejoin="round"
		class="stroke-primary"
	/>

	<g stroke-width="5" class="stroke-attn-fill">
		<rect x="404" y="330" width="34" height="22" rx="11" />
		<rect x="424" y="330" width="34" height="22" rx="11" />
	</g>

	<path
		d="M262 70l4 10 10 4-10 4-4 10-4-10-10-4 10-4z M444 58l3 7 7 3-7 3-3 7-3-7-7-3 7-3z M60 248l3 6 6 3-6 3-3 6-3-6-6-3 6-3z"
		class="fill-attn-fill"
	/>

	<ellipse cx="212" cy="388" rx="62" ry="8" class="fill-line" />
	<path d="M168 294h88l-8 88h-24l-12-62-12 62h-24z" class="fill-illu-pants" />
	<rect x="170" y="378" width="36" height="13" rx="6.5" class="fill-fg" />
	<rect x="220" y="378" width="36" height="13" rx="6.5" class="fill-fg" />

	<path d="M172 204c0-12 8-18 18-18h44c10 0 18 6 18 18l6 96h-92z" class="fill-primary" />
	<path d="M199 186l13 17 13-17z" class="fill-illu-shade" />
	<rect x="202" y="168" width="20" height="22" rx="6" class="fill-illu-shade" />

	<path
		d="M176 204c-14 8-18 30-22 82"
		stroke-width="22"
		stroke-linecap="round"
		class="stroke-primary"
	/>
	<circle cx="154" cy="292" r="11" class="fill-illu-skin" />

	<path
		d="M246 206l30 42-10-46"
		stroke-width="22"
		stroke-linecap="round"
		stroke-linejoin="round"
		class="stroke-primary"
	/>
	<rect x="256" y="160" width="26" height="42" rx="6" class="fill-fg" />
	<rect x="259.5" y="165" width="19" height="30" rx="3" class="fill-secondary" />
	<circle cx="266" cy="200" r="11" class="fill-illu-skin" />

	<circle cx="183" cy="152" r="7" class="fill-illu-skin" />
	<circle cx="212" cy="146" r="31" class="fill-illu-skin" />
	<path
		d="M180 146c-2-26 16-40 36-38 22 2 32 20 28 40-8-12-20-16-32-14-12 1-22 6-32 12z"
		class="fill-illu-hair"
	/>
	<circle cx="236" cy="112" r="11" class="fill-illu-hair" />
	<circle cx="202" cy="150" r="3" class="fill-illu-hair" />
	<circle cx="225" cy="150" r="3" class="fill-illu-hair" />
	<path d="M204 162q9 7 18 0" stroke-width="3" stroke-linecap="round" class="stroke-illu-hair" />
	<circle cx="196" cy="160" r="4" class="fill-illu-shade" opacity="0.5" />
	<circle cx="231" cy="160" r="4" class="fill-illu-shade" opacity="0.5" />
</svg>
