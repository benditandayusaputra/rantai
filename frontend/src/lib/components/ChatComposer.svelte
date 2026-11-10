<script lang="ts">
	import { MapPin, MessageCircle, SendHorizontal, X } from '@lucide/svelte';
	import { t } from '$lib/i18n';

	interface Props {
		value: string;
		context?: { label: string; onremove?: () => void };
		placeholder?: string;
		suggestions?: string[];
		disabled?: boolean;
		maxLength: number;
		onsubmit: (q: string) => void;
	}

	let {
		value = $bindable(),
		context,
		placeholder = t('ask.placeholder'),
		suggestions = [],
		disabled = false,
		maxLength,
		onsubmit
	}: Props = $props();

	const idPetunjuk = $props.id();
	const isKosong = $derived(value.trim() === '');

	function kirim(event: SubmitEvent) {
		event.preventDefault();
		if (isKosong || disabled) return;
		onsubmit(value.trim());
	}

	function tanganiTombol(event: KeyboardEvent) {
		if (event.key !== 'Enter' || event.shiftKey) return;
		event.preventDefault();
		(event.currentTarget as HTMLTextAreaElement).form?.requestSubmit();
	}
</script>

<div class="flex flex-col gap-3">
	{#if suggestions.length}
		<ul
			aria-label={t('ask.suggestions')}
			class="grid grid-cols-[repeat(auto-fill,minmax(min(100%,theme(spacing.80)),1fr))] gap-2"
		>
			{#each suggestions as saran (saran)}
				<li>
					<button
						type="button"
						{disabled}
						class="card-link flex h-full min-h-touch w-full items-start gap-3 px-4 py-3 text-left text-small font-semibold disabled:opacity-50"
						onclick={() => onsubmit(saran)}
					>
						<MessageCircle
							size={18}
							strokeWidth={1.75}
							aria-hidden="true"
							class="mt-0.5 shrink-0 text-muted"
						/>{saran}</button
					>
				</li>
			{/each}
		</ul>
	{/if}
	<form
		class="rounded-lg border-ctrl border-line-strong bg-raised shadow-1 transition-colors duration-fast ease-standard focus-within:border-focus hover:border-fg"
		onsubmit={kirim}
	>
		{#if context}
			<div class="flex px-3 pt-3">
				<span
					class="inline-flex h-8 items-center gap-1 rounded-md bg-secondary pl-3 text-small text-secondary-on {context.onremove
						? ''
						: 'pr-3'}"
				>
					<MapPin size={16} strokeWidth={1.75} aria-hidden="true" />
					{t('ask.context', { area: context.label })}
					{#if context.onremove}
						<button
							type="button"
							class="inline-flex min-h-touch min-w-touch items-center justify-center"
							aria-label={t('ask.changeArea')}
							title={t('ask.changeArea')}
							onclick={context.onremove}
						>
							<X size={16} strokeWidth={1.75} aria-hidden="true" />
						</button>
					{/if}
				</span>
			</div>
		{/if}
		<div class="flex items-end gap-2 p-2">
			<label class="flex-1">
				<span class="sr-only">{placeholder}</span>
				<textarea
					bind:value
					rows="2"
					maxlength={maxLength}
					{placeholder}
					{disabled}
					aria-describedby={idPetunjuk}
					class="block w-full resize-none bg-transparent px-2 py-2 outline-none placeholder:text-muted"
					onkeydown={tanganiTombol}></textarea>
			</label>
			<button
				type="submit"
				disabled={disabled || isKosong}
				class="btn btn-primary btn-icon"
				aria-label={t('ask.send')}
			>
				<SendHorizontal size={20} strokeWidth={1.75} aria-hidden="true" />
			</button>
		</div>
	</form>
	<p id={idPetunjuk} class="text-caption text-muted">
		{t('ask.caution')}
		{t('ask.tooLong', { n: maxLength })}
	</p>
</div>
