<script lang="ts">
	import { ChartColumn, Table } from '@lucide/svelte';
	import { t } from '$lib/i18n';

	interface Props {
		pressed: boolean;
		controls: string;
		labelKembali?: string;
		ikonKembali?: typeof ChartColumn;
	}

	let {
		pressed = $bindable(),
		controls,
		labelKembali = t('action.viewAsChart'),
		ikonKembali: IkonKembali = ChartColumn
	}: Props = $props();
</script>

<button
	type="button"
	aria-pressed={pressed}
	aria-controls={controls}
	class="btn btn-outline self-start text-small shadow-1 aria-pressed:border-primary aria-pressed:bg-secondary aria-pressed:text-secondary-on"
	onclick={() => (pressed = !pressed)}
>
	{#if pressed}
		<IkonKembali size={18} strokeWidth={1.75} aria-hidden="true" />
		{labelKembali}
	{:else}
		<Table size={18} strokeWidth={1.75} aria-hidden="true" />
		{t('action.viewAsTable')}
	{/if}
</button>
