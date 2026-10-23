<script lang="ts">
	import { RefreshCw, TriangleAlert } from '@lucide/svelte';
	import { t } from '$lib/i18n';

	interface Props {
		title: string;
		description?: string;
		onretry?: () => void;
	}

	let { title, description, onretry }: Props = $props();
</script>

<div
	role="alert"
	class="flex flex-col items-start gap-3 rounded-lg border border-line bg-surface p-5"
>
	<span
		class="inline-flex h-11 w-11 items-center justify-center rounded-pill bg-danger-bg text-danger"
	>
		<TriangleAlert size={20} strokeWidth={1.75} aria-hidden="true" />
	</span>
	<p class="font-bold">{title}</p>
	{#if description}<p class="text-small text-muted">{description}</p>{/if}
	{#if onretry}
		<button
			type="button"
			class="inline-flex min-h-touch items-center gap-2 rounded-md border-ctrl border-line-strong px-4 font-semibold hover:border-fg hover:bg-sunken"
			onclick={onretry}
		>
			<RefreshCw size={16} strokeWidth={1.75} aria-hidden="true" />
			{t('action.retry')}
		</button>
	{/if}
</div>
