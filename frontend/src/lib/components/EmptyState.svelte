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
	class="flex flex-col items-center gap-3 rounded-lg border-ctrl border-dashed border-line-strong bg-surface px-5 py-8 text-center"
>
	<span class="inline-flex h-12 w-18 items-center justify-center rounded-md bg-sunken">
		<svg viewBox="0 0 72 24" class="h-5 w-15 text-line-strong" fill="none" aria-hidden="true">
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
	</span>
	<p class="max-w-prose font-bold">{title}</p>
	{#if description}<p class="max-w-prose text-small text-muted">{description}</p>{/if}
	{#if actions}
		<div class="mt-1 flex flex-wrap justify-center gap-2">{@render actions()}</div>
	{/if}
</div>
