<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		label: string;
		class?: string;
		children: Snippet;
	}

	let { label, class: kelas = '', children }: Props = $props();
</script>

<div role="region" aria-label={label} tabindex="0" class="card overflow-x-auto {kelas}">
	{@render children()}
</div>
