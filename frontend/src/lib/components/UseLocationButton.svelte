<script lang="ts">
	import { LoaderCircle, LocateFixed } from '@lucide/svelte';
	import type { StatusLokasi } from '$lib/features/wilayah/lokasi';
	import { formatJarak } from '$lib/format';
	import { bahasaBawaan, t } from '$lib/i18n';

	interface Props {
		state: StatusLokasi;
		akurasiM?: number;
		onclick: () => void;
	}

	let { state, akurasiM = 0, onclick }: Props = $props();

	const id = $props.id();
	const pesan = $derived(
		{
			idle: '',
			locating: t('location.locating'),
			found: t('location.found', { m: Math.round(akurasiM) }),
			coarse: t('location.coarse', { dist: formatJarak(akurasiM, bahasaBawaan) }),
			denied: `${t('location.denied')}. ${t('location.deniedHint')}`,
			failed: t('location.failed')
		}[state]
	);
	const warnaPesan = $derived(
		state === 'coarse' ? 'text-attn' : state === 'found' ? 'text-success' : 'text-muted'
	);
</script>

<div class="flex flex-col gap-2">
	<button
		type="button"
		class="inline-flex min-h-touch items-center justify-center gap-2 self-start rounded-md border-ctrl px-4 font-semibold transition-colors duration-fast ease-standard disabled:border-dis-bg disabled:bg-dis-bg disabled:text-dis-fg {state ===
		'found'
			? 'border-success bg-success-bg text-success'
			: 'border-line-strong hover:border-fg hover:bg-sunken'}"
		disabled={state === 'denied' || state === 'locating'}
		aria-busy={state === 'locating'}
		aria-describedby={pesan ? `${id}-status` : undefined}
		{onclick}
	>
		{#if state === 'locating'}
			<LoaderCircle size={20} strokeWidth={1.75} aria-hidden="true" class="animate-spin" />
			{t('location.locating')}
		{:else}
			<LocateFixed size={20} strokeWidth={1.75} aria-hidden="true" />
			{t('location.use')}
		{/if}
	</button>
	<p id="{id}-status" class="text-small {warnaPesan}" aria-live="polite">{pesan}</p>
</div>
