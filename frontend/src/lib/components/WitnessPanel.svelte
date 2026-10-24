<script lang="ts">
	import { LoaderCircle, Users } from '@lucide/svelte';
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

<section
	aria-labelledby="{id}-judul"
	class="flex flex-col gap-3 rounded-lg border-ctrl border-line bg-surface p-4"
>
	<h2 id="{id}-judul" class="flex items-center gap-2 t-h3">
		<Users size={20} strokeWidth={1.75} aria-hidden="true" />{t('report.witness.title')}
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
		class="inline-flex min-h-13 w-full items-center justify-center gap-2 rounded-md bg-primary px-4 font-semibold text-primary-on hover:bg-primary-hover disabled:bg-dis-bg disabled:text-dis-fg"
		disabled={!canWitness || busy}
		aria-busy={busy}
		aria-describedby="{id}-syarat"
		onclick={onwitness}
	>
		{#if busy}<LoaderCircle
				size={20}
				strokeWidth={1.75}
				class="animate-spin"
				aria-hidden="true"
			/>{/if}
		{t('report.witness.cta')}
	</button>
	<p id="{id}-syarat" class="text-small text-muted">{reason ?? t('report.witness.rule')}</p>
</section>
