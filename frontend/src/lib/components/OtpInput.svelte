<script lang="ts">
	import { bersihkanKode, panjangKode } from '$lib/features/akun/masuk';

	interface Props {
		id: string;
		value: string;
		error?: 'invalid' | 'expired' | null;
		locked?: boolean;
		describedby?: string;
		oncomplete?: (kode: string) => void;
	}

	let {
		id,
		value = $bindable(''),
		error = null,
		locked = false,
		describedby,
		oncomplete
	}: Props = $props();

	let isFokus = $state(false);
	const indeksAktif = $derived(Math.min(value.length, panjangKode - 1));

	function saatMengetik(event: Event & { currentTarget: HTMLInputElement }) {
		value = bersihkanKode(event.currentTarget.value);
		event.currentTarget.value = value;
		if (value.length === panjangKode) oncomplete?.(value);
	}
</script>

<div class="relative inline-flex gap-2">
	{#each { length: panjangKode }, i (i)}
		<span
			aria-hidden="true"
			data-aktif={(isFokus && !locked && i === indeksAktif) || undefined}
			data-error={error ? true : undefined}
			class="flex h-14 w-11 items-center justify-center rounded-md border-ctrl border-line-strong bg-raised font-mono text-h2 sm:w-14 data-[aktif]:border-focus data-[error]:border-danger data-[aktif]:shadow-ring"
			class:bg-dis-bg={locked}
			class:text-dis-fg={locked}>{value[i] ?? ''}</span
		>
	{/each}
	<input
		{id}
		{value}
		type="text"
		inputmode="numeric"
		autocomplete="one-time-code"
		pattern="[0-9]*"
		disabled={locked}
		aria-invalid={error ? 'true' : undefined}
		aria-describedby={describedby}
		class="absolute inset-0 h-full w-full cursor-text opacity-0 disabled:cursor-not-allowed"
		oninput={saatMengetik}
		onfocus={() => (isFokus = true)}
		onblur={() => (isFokus = false)}
	/>
</div>
