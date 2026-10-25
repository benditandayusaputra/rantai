<script lang="ts">
	import { t } from '$lib/i18n';

	interface Props {
		id: string;
		label: string;
		value: string;
		help?: string;
		rows?: number;
		maxlength: number;
	}

	let { id, label, value = $bindable(), help, rows = 5, maxlength }: Props = $props();

	const ambangPeringatan = 20;
	const sisa = $derived(maxlength - value.length);
</script>

<div class="flex flex-col gap-2">
	<label for={id} class="font-semibold">{label}</label>
	{#if help}<p id="{id}-bantuan" class="text-small text-muted">{help}</p>{/if}
	<textarea
		{id}
		{rows}
		{maxlength}
		bind:value
		aria-describedby={help ? `${id}-bantuan` : undefined}
		class="field h-auto min-h-32 resize-y py-3"></textarea>
	<p
		class="self-end text-small text-muted num"
		aria-live={sisa <= ambangPeringatan ? 'polite' : 'off'}
	>
		{t('report.counter', { n: value.length, max: maxlength })}
	</p>
</div>
