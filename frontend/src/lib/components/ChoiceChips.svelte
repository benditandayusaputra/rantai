<script lang="ts">
	import { Check } from '@lucide/svelte';

	interface Props {
		labelledby: string;
		items: { key: string; href: string; label: string; current: boolean }[];
	}

	let { labelledby, items }: Props = $props();
</script>

<nav aria-labelledby={labelledby}>
	<ul class="flex flex-wrap gap-2">
		{#each items as item (item.key)}
			<li>
				<a
					href={item.href}
					aria-current={item.current}
					data-sveltekit-noscroll
					class="chip aria-[current=true]:border-primary aria-[current=true]:bg-secondary aria-[current=true]:text-secondary-on"
				>
					{#if item.current}<Check size={16} strokeWidth={2} aria-hidden="true" />{/if}
					{item.label}
				</a>
			</li>
		{/each}
	</ul>
</nav>
