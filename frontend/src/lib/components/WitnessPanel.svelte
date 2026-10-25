<script lang="ts">
	import { LoaderCircle, UserCheck, Users } from '@lucide/svelte';
	import { t } from '$lib/i18n';

	interface Props {
		count: number;
		needed?: number;
		canWitness: boolean;
		reason?: string;
		busy?: boolean;
		onwitness: () => void;
	}

	let { count, needed = 3, canWitness, reason, busy = false, onwitness }: Props = $props();

	const id = $props.id();
</script>

<section aria-labelledby="{id}-judul" class="card flex flex-col gap-4 p-5">
	<h2 id="{id}-judul" class="t-h3 flex items-center gap-3">
		<span class="icon-tile h-9 w-9"><Users size={18} strokeWidth={1.75} aria-hidden="true" /></span
		>{t('report.witness.title')}
	</h2>
	<div class="flex gap-1.5" aria-hidden="true">
		{#each Array.from({ length: needed }, (_, i) => i) as urutan (urutan)}
			<span
				class="h-2 flex-1 rounded-pill {urutan < count
					? 'bg-primary'
					: 'border-ctrl border-dotted border-line-strong'}"
			></span>
		{/each}
	</div>
	<p class="font-semibold num">{t('report.witness.count', { n: count, m: needed })}</p>
	<button
		type="button"
		class="btn btn-primary btn-lg btn-block"
		disabled={!canWitness || busy}
		aria-busy={busy}
		aria-describedby="{id}-syarat"
		onclick={onwitness}
	>
		{#if busy}
			<LoaderCircle size={20} strokeWidth={1.75} class="animate-spin" aria-hidden="true" />
		{:else}
			<UserCheck size={20} strokeWidth={1.75} aria-hidden="true" />
		{/if}
		{t('report.witness.cta')}
	</button>
	<p id="{id}-syarat" class="text-small text-muted">{reason ?? t('report.witness.rule')}</p>
</section>
