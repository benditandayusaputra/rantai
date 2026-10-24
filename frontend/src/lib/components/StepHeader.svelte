<script lang="ts">
	import { ArrowLeft, X } from '@lucide/svelte';
	import { t } from '$lib/i18n';

	interface Props {
		step: number;
		total?: number;
		stepLabel: string;
		onback: () => void;
		onclose: () => void;
	}

	let { step, total = 6, stepLabel, onback, onclose }: Props = $props();

	const teksLangkah = $derived(t('report.step', { n: step, total }));
</script>

<header class="sticky top-0 z-sticky border-b border-line bg-surface">
	<div class="mx-auto flex max-w-form items-center gap-3 px-4 pt-3">
		<button
			type="button"
			class="inline-flex h-11 w-11 shrink-0 items-center justify-center rounded-md hover:bg-sunken"
			aria-label={step === 1 ? t('action.close') : t('action.back')}
			onclick={step === 1 ? onclose : onback}
		>
			{#if step === 1}
				<X size={20} strokeWidth={1.75} aria-hidden="true" />
			{:else}
				<ArrowLeft size={20} strokeWidth={1.75} aria-hidden="true" />
			{/if}
		</button>
		<div class="flex flex-col">
			<span class="font-bold">{t('report.flowTitle')}</span>
			<span class="text-small text-muted">{teksLangkah}</span>
		</div>
	</div>
	<div
		role="progressbar"
		aria-valuemin={1}
		aria-valuemax={total}
		aria-valuenow={step}
		aria-valuetext="{teksLangkah}, {stepLabel}"
		class="mx-auto flex max-w-form gap-1.5 px-4 pb-3 pt-2"
	>
		{#each Array.from({ length: total }, (_, i) => i + 1) as urutan (urutan)}
			<span
				class="h-1.5 flex-1 rounded-pill {urutan < step
					? 'bg-primary'
					: urutan === step
						? 'border-ctrl border-primary bg-secondary'
						: 'border-ctrl border-dotted border-line-strong'}"
			></span>
		{/each}
	</div>
</header>
