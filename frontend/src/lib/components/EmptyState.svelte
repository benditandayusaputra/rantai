<script lang="ts">
	import type { Snippet } from 'svelte';
	import LogoMark from './LogoMark.svelte';

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
		<LogoMark class="h-5 w-auto fill-line-strong" />
	</span>
	<p class="max-w-prose font-bold">{title}</p>
	{#if description}<p class="max-w-prose text-small text-muted">{description}</p>{/if}
	{#if actions}
		<div class="mt-1 flex flex-wrap justify-center gap-2">{@render actions()}</div>
	{/if}
</div>
