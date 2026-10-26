<script lang="ts">
	let { class: kelas = 'w-full' }: { class?: string } = $props();

	const id = $props.id();
	const gedung = [
		{ x: 70, y: 250, l: 34, t: 144 },
		{ x: 108, y: 286, l: 26, t: 108 },
		{ x: 330, y: 270, l: 30, t: 124 },
		{ x: 364, y: 236, l: 36, t: 158 },
		{ x: 404, y: 300, l: 24, t: 94 }
	];
	const garisPembatas = [126, 150, 174, 198];
	const sudutBidik = ['M292 200v-10h10', 'M334 190h10v10', 'M344 218v10h-10', 'M302 228h-10v-10'];
	const saksi = ['fill-illu-skin', 'fill-gap-3', 'fill-primary'];
</script>

<svg viewBox="0 0 480 440" class={kelas} fill="none" aria-hidden="true">
	<defs>
		<clipPath id="{id}-foto">
			<rect x="256" y="96" width="114" height="140" rx="10" />
		</clipPath>
	</defs>

	<circle cx="240" cy="214" r="186" class="fill-secondary" />
	<ellipse cx="240" cy="394" rx="210" ry="10" class="fill-secondary-hover" />

	<g class="fill-secondary-hover">
		{#each gedung as g (g.x)}
			<rect x={g.x} y={g.y} width={g.l} height={g.t} rx="3" />
		{/each}
	</g>

	<g class="stroke-line-strong" stroke-linecap="round">
		<path d="M96 394V122" stroke-width="6" />
		<path d="M56 136h180" stroke-width="6" />
		<path d="M96 112L60 136M96 112l136 24" stroke-width="2" />
		<path d="M206 136v58" stroke-width="2" />
	</g>
	<rect x="198" y="194" width="16" height="14" rx="3" class="fill-attn-fill" />

	<polygon points="40,394 440,394 350,322 130,322" class="fill-line" />
	<ellipse cx="306" cy="372" rx="34" ry="9" class="fill-line-strong" />
	<path
		d="M278 368l-14-6M336 370l16-5M306 380l4 10"
		stroke-width="2.5"
		stroke-linecap="round"
		class="stroke-line-strong"
	/>

	<rect x="128" y="352" width="6" height="34" rx="2" class="fill-line-strong" />
	<rect x="196" y="352" width="6" height="34" rx="2" class="fill-line-strong" />
	<rect x="118" y="336" width="94" height="18" rx="4" class="fill-raised" />
	<g class="fill-attn-fill">
		{#each garisPembatas as x (x)}
			<polygon points="{x},354 {x + 12},354 {x + 22},336 {x + 10},336" />
		{/each}
	</g>
	<polygon points="394,388 422,388 412,338 404,338" class="fill-attn-fill" />
	<polygon points="400,366 416,366 414,356 402,356" class="fill-raised" />
	<rect x="386" y="386" width="44" height="7" rx="3.5" class="fill-attn-fill" />

	<g transform="rotate(-6 312 198)">
		<rect x="238" y="64" width="150" height="268" rx="24" class="fill-fg" />
		<rect x="246" y="74" width="134" height="248" rx="18" class="fill-raised" />
		<rect x="293" y="80" width="40" height="6" rx="3" class="fill-fg" />
		<g clip-path="url(#{id}-foto)">
			<rect x="256" y="96" width="114" height="140" class="fill-secondary" />
			<rect x="264" y="150" width="14" height="34" class="fill-secondary-hover" />
			<rect x="282" y="138" width="12" height="46" class="fill-secondary-hover" />
			<rect x="344" y="146" width="16" height="38" class="fill-secondary-hover" />
			<polygon points="250,236 376,236 340,180 286,180" class="fill-line" />
			<ellipse cx="318" cy="212" rx="18" ry="5" class="fill-line-strong" />
			<polygon points="346,230 358,230 354,204 350,204" class="fill-attn-fill" />
		</g>
		<g class="stroke-primary" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
			{#each sudutBidik as d (d)}<path {d} />{/each}
		</g>
		<rect x="256" y="246" width="114" height="22" rx="11" class="fill-sunken" />
		<circle cx="268" cy="257" r="5" class="fill-primary" />
		<rect x="279" y="254" width="70" height="6" rx="3" class="fill-line-strong" />
		<circle cx="313" cy="294" r="18" class="fill-primary" />
		<circle cx="313" cy="294" r="11" stroke-width="3" class="stroke-raised" />
	</g>

	<ellipse cx="190" cy="230" rx="16" ry="4" class="fill-line" />
	<path d="M190 216c-6-14-36-42-36-68a36 36 0 0 1 72 0c0 26-30 54-36 68z" class="fill-primary" />
	<circle cx="190" cy="148" r="14" class="fill-raised" />

	<path d="M412 70l4 11 11 4-11 4-4 11-4-11-11-4 11-4z" class="fill-attn-fill" />
	<path d="M438 118l2 6 6 2-6 2-2 6-2-6-6-2 6-2z" class="fill-attn-fill" />
	<path d="M218 66l2 6 6 2-6 2-2 6-2-6-6-2 6-2z" class="fill-primary" opacity="0.6" />

	<rect x="40" y="262" width="176" height="60" rx="16" class="fill-raised stroke-line" />
	{#each saksi as warna, i (warna)}
		<circle cx={70 + i * 20} cy="292" r="13" stroke-width="3" class="{warna} stroke-raised" />
	{/each}
	<rect x="126" y="282" width="68" height="8" rx="4" class="fill-line-strong" />
	<rect x="126" y="296" width="46" height="6" rx="3" class="fill-line" />
	<circle cx="210" cy="266" r="14" class="fill-success" />
	<path
		d="M203.5 266l4.5 4.5 8.5-8.5"
		stroke-width="3"
		stroke-linecap="round"
		stroke-linejoin="round"
		class="stroke-raised"
	/>
</svg>
