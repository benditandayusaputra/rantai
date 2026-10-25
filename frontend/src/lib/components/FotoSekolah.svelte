<script lang="ts">
	import { School } from '@lucide/svelte';
	import type { FotoSekolah } from '$lib/api/tipe';

	interface Props {
		foto?: FotoSekolah | null;
		class?: string;
	}

	let { foto, class: kelas = 'h-12 w-16' }: Props = $props();

	let isGagal = $state(false);
</script>

<span
	class="inline-flex shrink-0 items-center justify-center overflow-hidden rounded-md bg-sunken text-muted {kelas}"
	aria-hidden="true"
>
	{#if foto && !isGagal}
		<img
			src={foto.url}
			alt=""
			loading="lazy"
			decoding="async"
			referrerpolicy="no-referrer"
			class="h-full w-full object-cover"
			onerror={() => (isGagal = true)}
		/>
	{:else}
		<School size={18} strokeWidth={1.75} />
	{/if}
</span>
