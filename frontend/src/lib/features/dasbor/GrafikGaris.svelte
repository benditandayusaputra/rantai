<script lang="ts">
	import type { Porsi } from '$lib/api/tipe';
	import { susunGaris, ukuranGaris } from './dasbor';

	interface Props {
		isi: Porsi[];
		nilaiDari: (p: Porsi) => number;
		format: (n: number) => string;
		labelDari: (p: Porsi) => string;
	}

	let { isi, nilaiDari, format, labelDari }: Props = $props();

	const { lebar, tinggi, kiri, kanan, atas, bawah } = ukuranGaris;
	const garisDasar = tinggi - bawah;
	const garisBantu = [atas, (atas + garisDasar) / 2];

	const titik = $derived(
		susunGaris(isi.map(nilaiDari)).map((p, i) => {
			const porsi = isi[i] as Porsi;
			return { ...p, kunci: porsi.kunci, label: labelDari(porsi), nilai: format(nilaiDari(porsi)) };
		})
	);
	const jalurGaris = $derived(titik.map((p, i) => `${i === 0 ? 'M' : 'L'}${p.x},${p.y}`).join(''));
	const jalurArea = $derived.by(() => {
		const awal = titik[0];
		const akhir = titik.at(-1);
		if (!awal || !akhir) return '';
		return `${jalurGaris}L${akhir.x},${garisDasar}L${awal.x},${garisDasar}Z`;
	});

	const jarakLabel = 14;
	const jarakLabelBawah = 24;

	function isLabelDiBawah(i: number): boolean {
		const sebelum = titik[i - 1];
		const sesudah = titik[i + 1];
		const p = titik[i];
		const isLembah = Boolean(p && sebelum && sesudah && p.y > sebelum.y && p.y > sesudah.y);
		return isLembah && (p?.y ?? garisDasar) + jarakLabelBawah < garisDasar;
	}

	function jangkarNilai(i: number): 'start' | 'middle' | 'end' {
		if (titik.length < 2) return 'middle';
		if (i === 0) return 'start';
		return i === titik.length - 1 ? 'end' : 'middle';
	}
</script>

<figure>
	<svg viewBox="0 0 {lebar} {tinggi}" class="block h-auto w-full text-caption" aria-hidden="true">
		{#each garisBantu as y (y)}
			<line
				x1={kiri / 2}
				x2={lebar - kanan / 2}
				y1={y}
				y2={y}
				class="stroke-line"
				stroke-dasharray="4 4"
			/>
		{/each}
		<line
			x1={kiri / 2}
			x2={lebar - kanan / 2}
			y1={garisDasar}
			y2={garisDasar}
			stroke-width="1.5"
			class="stroke-line-strong"
		/>
		<path d={jalurArea} class="fill-chart" fill-opacity="0.12" />
		<path
			d={jalurGaris}
			fill="none"
			stroke-width="3"
			stroke-linecap="round"
			stroke-linejoin="round"
			class="stroke-chart"
		/>
		{#each titik as p, i (p.kunci)}
			<circle cx={p.x} cy={p.y} r="5.5" stroke-width="3" class="fill-raised stroke-chart" />
			<text
				x={p.x}
				y={isLabelDiBawah(i) ? p.y + jarakLabelBawah : p.y - jarakLabel}
				text-anchor={jangkarNilai(i)}
				class="num fill-fg font-bold">{p.nilai}</text
			>
			<text x={p.x} y={tinggi - 12} text-anchor="middle" class="fill-muted">{p.label}</text>
		{/each}
	</svg>
	<ul class="sr-only">
		{#each titik as p (p.kunci)}
			<li>{p.label}: {p.nilai}</li>
		{/each}
	</ul>
</figure>
