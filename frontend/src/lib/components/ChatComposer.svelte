<script lang="ts">
	import { MapPin, SendHorizontal, X } from '@lucide/svelte';
	import { t } from '$lib/i18n';

	interface Props {
		value: string;
		context?: { label: string; onremove: () => void };
		suggestions?: string[];
		disabled?: boolean;
		maxLength: number;
		onsubmit: (q: string) => void;
	}

	let {
		value = $bindable(),
		context,
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
		<ul aria-label={t('ask.suggestions')} class="-mx-1 flex gap-2 overflow-x-auto px-1 pb-1">
			{#each suggestions as saran (saran)}
				<li class="shrink-0">
					<button
						type="button"
						{disabled}
						class="min-h-touch max-w-66 rounded-lg border-ctrl border-line-strong px-4 py-2 text-left text-small hover:bg-sunken disabled:opacity-50"
						onclick={() => onsubmit(saran)}>{saran}</button
					>
				</li>
			{/each}
		</ul>
	{/if}
	<form
		class="rounded-lg border-ctrl border-line-strong bg-surface focus-within:border-primary"
		onsubmit={kirim}
	>
		{#if context}
			<div class="flex px-3 pt-3">
				<span
					class="inline-flex h-8 items-center gap-1 rounded-pill bg-secondary pl-3 text-small text-secondary-on"
				>
					<MapPin size={16} strokeWidth={1.75} aria-hidden="true" />
					{t('ask.context', { area: context.label })}
					<button
						type="button"
						class="inline-flex min-h-touch min-w-touch items-center justify-center"
						aria-label={t('ask.removeContext')}
						onclick={context.onremove}
					>
						<X size={16} strokeWidth={1.75} aria-hidden="true" />
					</button>
				</span>
			</div>
		{/if}
		<div class="flex items-end gap-2 p-2">
			<label class="flex-1">
				<span class="sr-only">{t('ask.placeholder')}</span>
				<textarea
					bind:value
					rows="2"
					maxlength={maxLength}
					placeholder={t('ask.placeholder')}
					aria-describedby={idPetunjuk}
					class="block w-full resize-none bg-transparent px-2 py-2 outline-none placeholder:text-muted"
					onkeydown={tanganiTombol}></textarea>
			</label>
			<button
				type="submit"
				disabled={disabled || isKosong}
				class="inline-flex min-h-touch min-w-touch items-center justify-center rounded-md bg-primary text-primary-on disabled:opacity-50"
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
