<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		label: string;
		class?: string;
		children: Snippet;
	}

	let { label, class: kelas = '', children }: Props = $props();
</script>

<div
	role="region"
	aria-label={label}
	tabindex="0"
	class="overflow-x-auto rounded-lg border border-line bg-surface {kelas}"
>
	{@render children()}
</div>
