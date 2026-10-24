<script lang="ts">
	import type { Snippet } from 'svelte';
	import { CircleAlert, Unlink2 } from '@lucide/svelte';
	import { t } from '$lib/i18n';
	import LogoMark from './LogoMark.svelte';

	interface Props {
		role: 'user' | 'assistant';
		state?: 'streaming' | 'done' | 'no-data' | 'error';
		children: Snippet;
	}

	let { role, state = 'done', children }: Props = $props();
</script>

{#if role === 'user'}
	<div class="flex justify-end">
		<p class="sr-only">{t('ask.you')}</p>
		<div class="ml-12 rounded-lg rounded-br-xs bg-secondary px-4 py-3 text-secondary-on">
			{@render children()}
		</div>
	</div>
{:else}
	<div class="flex items-start gap-3" aria-busy={state === 'streaming'}>
		<span class="mt-1 shrink-0"><LogoMark class="h-6 w-6" /></span>
		<div
			class="min-w-0 flex-1 rounded-lg border px-4 py-3 {state === 'no-data'
				? 'border-dashed border-line-strong bg-sunken'
				: 'border-line bg-raised'}"
		>
			{#if state === 'no-data'}
				<p class="mb-2 flex items-center gap-2 font-bold">
					<Unlink2 size={16} strokeWidth={1.75} aria-hidden="true" />{t('ask.noData')}
				</p>
			{:else if state === 'error'}
				<CircleAlert size={16} strokeWidth={1.75} aria-hidden="true" class="mb-2 text-danger" />
			{/if}
			{@render children()}
			{#if state === 'streaming'}
				<span
					class="ml-0.5 inline-block h-4 w-0.5 animate-caret bg-fg align-middle"
					aria-hidden="true"
				></span>
			{/if}
		</div>
	</div>
{/if}
