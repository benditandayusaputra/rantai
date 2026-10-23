<script lang="ts">
	import type { Angka } from '$lib/api/tipe';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import ProvenanceBadge from './ProvenanceBadge.svelte';

	interface Props {
		label: string;
		angka: Angka;
		size?: 'md' | 'lg';
	}

	let { label, angka, size = 'md' }: Props = $props();

	const tampil = $derived(tampilkanAngka(angka));
</script>

<div
	class="flex flex-col gap-2 rounded-lg border p-4 md:p-5 {angka.label === 'simulasi'
		? 'border-cert-sim bg-cert-sim-bg'
		: 'border-line bg-surface'}"
>
	<p class="text-caption text-muted">{label}</p>
	<p class={size === 'lg' ? 'n-lg' : 'n-md'}>
		{tampil.nilai}
		{#if tampil.satuan}<span class="text-small font-semibold text-muted">{tampil.satuan}</span>{/if}
	</p>
	<div><ProvenanceBadge {angka} nilai={tampil.nilai} compact /></div>
</div>
