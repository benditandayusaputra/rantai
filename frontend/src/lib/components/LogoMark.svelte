<script lang="ts">
	let { class: kelas = 'h-7 w-7' }: { class?: string } = $props();
</script>

<svg viewBox="0 0 32 32" class={kelas} fill="none" stroke-width="3.5" aria-hidden="true">
	<rect x="3" y="10" width="16" height="12" rx="6" class="stroke-primary" />
	<rect x="13" y="10" width="16" height="12" rx="6" class="stroke-fg" />
</svg>
