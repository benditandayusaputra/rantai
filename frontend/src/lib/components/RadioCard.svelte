<script lang="ts" generics="T extends string">
	import type { Snippet } from 'svelte';

	interface Props {
		name: string;
		value: T;
		group: T | null;
		title: string;
		description?: string;
		leading?: Snippet;
	}

	let { name, value, group = $bindable(), title, description, leading }: Props = $props();
</script>

<label
	class="flex min-h-14 cursor-pointer items-center gap-3 rounded-lg border-ctrl border-line-strong bg-surface p-4 hover:border-fg has-[:checked]:border-secondary-on has-[:checked]:bg-secondary has-[:focus-visible]:shadow-ring"
>
	<input type="radio" {name} {value} bind:group class="sr-only" />
	{@render leading?.()}
	<span class="flex flex-col gap-1">
		<span class="font-semibold">{title}</span>
		{#if description}<span class="text-small text-muted">{description}</span>{/if}
	</span>
</label>
