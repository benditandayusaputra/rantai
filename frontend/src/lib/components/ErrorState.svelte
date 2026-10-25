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

<div role="alert" class="card flex flex-col items-center gap-3 px-5 py-8 text-center">
	<span
		class="inline-flex h-12 w-12 items-center justify-center rounded-pill bg-danger-bg text-danger"
	>
		<TriangleAlert size={24} strokeWidth={1.75} aria-hidden="true" />
	</span>
	<p class="max-w-prose font-bold">{title}</p>
	{#if description}<p class="max-w-prose text-small text-muted">{description}</p>{/if}
	{#if onretry}
		<button type="button" class="btn btn-outline mt-1" onclick={onretry}>
			<RefreshCw size={18} strokeWidth={1.75} aria-hidden="true" />
			{t('action.retry')}
		</button>
	{/if}
</div>
