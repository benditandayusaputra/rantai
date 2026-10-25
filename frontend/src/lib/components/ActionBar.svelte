<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<div class="glass sticky bottom-0 z-sticky border-t border-line pb-[env(safe-area-inset-bottom)]">
	<div class="mx-auto flex max-w-form flex-col gap-1 px-4 py-3">
		{@render children()}
	</div>
</div>
