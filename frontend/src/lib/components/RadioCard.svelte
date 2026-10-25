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
	class="group flex min-h-14 cursor-pointer items-center gap-3 rounded-lg border-ctrl border-line-strong bg-raised p-4 transition-colors duration-fast ease-standard hover:border-fg has-[:checked]:border-primary has-[:checked]:bg-secondary has-[:focus-visible]:shadow-ring"
>
	<input type="radio" {name} {value} bind:group class="sr-only" />
	<span class="flex min-w-0 flex-1 flex-col items-start gap-1">
		{@render leading?.()}
		<span class="break-words font-semibold">{title}</span>
		{#if description}<span class="text-small text-muted">{description}</span>{/if}
	</span>
	<span
		aria-hidden="true"
		class="inline-flex h-6 w-6 shrink-0 items-center justify-center rounded-pill border-ctrl border-line-strong bg-raised group-has-[:checked]:border-primary group-has-[:checked]:bg-primary"
	>
		<span class="h-2.5 w-2.5 rounded-pill bg-primary-on opacity-0 group-has-[:checked]:opacity-100"
		></span>
	</span>
</label>
