<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		description?: string;
		actions?: Snippet;
	}

	let { title, description, actions }: Props = $props();
</script>

<div
	class="flex flex-col items-start gap-3 rounded-lg border-ctrl border-dashed border-line-strong p-5"
>
	<svg viewBox="0 0 72 24" class="h-6 w-18 text-line-strong" fill="none" aria-hidden="true">
		{#each [0, 1, 2] as urutan (urutan)}
			<rect
				x={4 + urutan * 22}
				y="4"
				width="20"
				height="16"
				rx="8"
				stroke="currentColor"
				stroke-width="3"
				stroke-dasharray="2 4"
			/>
		{/each}
	</svg>
	<p class="font-bold">{title}</p>
	{#if description}<p class="text-small text-muted">{description}</p>{/if}
	{@render actions?.()}
</div>
