<script lang="ts" module>
	export interface IsiPopup {
		keterangan: string;
		judul: string;
		tautan?: { href: string; label: string };
	}
</script>

<script lang="ts">
	import { ArrowRight } from '@lucide/svelte';

	let { keterangan, judul, tautan }: IsiPopup = $props();
</script>

<div class="flex max-w-60 flex-col gap-1 pr-4">
	<p class="text-caption text-muted">{keterangan}</p>
	<p class="text-small font-semibold">{judul}</p>
	{#if tautan}
		<a
			href={tautan.href}
			class="mt-1 inline-flex items-center gap-1 text-small font-semibold text-primary underline-offset-4 hover:underline"
		>
			{tautan.label}
			<ArrowRight size={14} strokeWidth={1.75} aria-hidden="true" />
		</a>
	{/if}
</div>
