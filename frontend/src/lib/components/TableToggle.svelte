<script lang="ts">
	import { ChartColumn, Table } from '@lucide/svelte';
	import { t } from '$lib/i18n';

	interface Props {
		pressed: boolean;
		controls: string;
		labelKembali?: string;
	}

	let { pressed = $bindable(), controls, labelKembali = t('action.viewAsChart') }: Props = $props();
</script>

<button
	type="button"
	aria-pressed={pressed}
	aria-controls={controls}
	class="inline-flex min-h-touch self-start items-center gap-2 rounded-md border-ctrl border-line-strong px-4 text-small font-semibold hover:border-fg hover:bg-sunken"
	onclick={() => (pressed = !pressed)}
>
	{#if pressed}
		<ChartColumn size={16} strokeWidth={1.75} aria-hidden="true" />
		{labelKembali}
	{:else}
		<Table size={16} strokeWidth={1.75} aria-hidden="true" />
		{t('action.viewAsTable')}
	{/if}
</button>
