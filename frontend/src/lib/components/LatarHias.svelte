<script lang="ts">
	import LogoMark from './LogoMark.svelte';

	let { class: kelas = '' }: { class?: string } = $props();
</script>

<div class="pointer-events-none absolute inset-0 -z-10 overflow-hidden {kelas}" aria-hidden="true">
	<div class="latar-titik absolute inset-0"></div>
	<div class="cahaya-primer absolute -left-48 -top-48 h-96 w-96 scale-150"></div>
	<div class="cahaya-sekunder absolute -bottom-48 -right-32 h-96 w-96 scale-150"></div>
	<LogoMark
		class="absolute -bottom-8 -right-24 h-auto w-full max-w-3xl -rotate-12 fill-primary opacity-5"
	/>
</div>
