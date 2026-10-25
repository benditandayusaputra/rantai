<script lang="ts">
	import type { Angka } from '$lib/api/tipe';
	import { tampilkanAngka } from '$lib/features/kecamatan/kecamatan';
	import ProvenanceBadge from './ProvenanceBadge.svelte';

	interface Props {
		label: string;
		angka: Angka;
		size?: 'md' | 'lg';
		catatan?: string;
	}

	let { label, angka, size = 'md', catatan }: Props = $props();

	const tampil = $derived(tampilkanAngka(angka));
</script>

<div
	class="card flex flex-col gap-2 p-3 md:p-5 {angka.label === 'simulasi'
		? 'border-cert-sim bg-cert-sim-bg'
		: ''}"
>
	<p class="text-caption text-muted">{label}</p>
	<p class={size === 'lg' ? 'n-lg' : 'n-md'}>
		{tampil.nilai}
		{#if tampil.satuan}<span class="text-small font-semibold text-muted">{tampil.satuan}</span>{/if}
	</p>
	{#if catatan}<p class="text-caption text-muted">{catatan}</p>{/if}
	<div><ProvenanceBadge {angka} nilai={tampil.nilai} compact /></div>
</div>
