<script lang="ts">
	import { inisial } from '$lib/features/kecamatan/kecamatan';

	interface Props {
		nama: string;
		foto?: string;
		class?: string;
	}

	let { nama, foto, class: kelas = 'h-11 w-11 text-small' }: Props = $props();

	let isGagal = $state(false);
</script>

<span
	class="inline-flex shrink-0 items-center justify-center overflow-hidden rounded-pill border border-line-strong bg-sunken font-semibold text-fg {kelas}"
	aria-hidden="true"
>
	{#if foto && !isGagal}
		<img
			src={foto}
			alt=""
			loading="lazy"
			decoding="async"
			referrerpolicy="no-referrer"
			class="h-full w-full object-cover object-top"
			onerror={() => (isGagal = true)}
		/>
	{:else}
		{inisial(nama)}
	{/if}
</span>
