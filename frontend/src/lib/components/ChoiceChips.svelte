<script lang="ts">
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
					class="inline-flex min-h-touch items-center rounded-pill border-ctrl border-line-strong px-4 text-small hover:bg-sunken aria-[current=true]:border-primary aria-[current=true]:bg-secondary aria-[current=true]:font-bold aria-[current=true]:text-secondary-on"
					>{item.label}</a
				>
			</li>
		{/each}
	</ul>
</nav>
