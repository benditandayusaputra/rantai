<script lang="ts">
	import '@fontsource-variable/plus-jakarta-sans';
	import '@fontsource-variable/newsreader';
	import '@fontsource/jetbrains-mono/400.css';
	import '@fontsource/jetbrains-mono/500.css';
	import '../app.css';
	import { ikutiTemaSistem } from '$lib/features/tema/tema';

	let { children } = $props();

	$effect(ikutiTemaSistem);
</script>

{@render children()}
