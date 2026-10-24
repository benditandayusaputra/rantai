<script lang="ts">
	import type { Label } from '$lib/api/tipe';
	import { t } from '$lib/i18n';
	import { kelasLabel } from './label';

	interface Props {
		n: number;
		level: Label;
		label: string;
		onopen: () => void;
	}

	let { n, level, label, onopen }: Props = $props();
</script>

<button
	type="button"
	class="relative mx-0.5 inline-flex h-5 min-w-5.5 items-center justify-center rounded-xs border-ctrl px-1 align-text-top text-caption font-bold leading-none after:absolute after:-inset-3 {kelasLabel[
		level
	]}"
	aria-haspopup="dialog"
	aria-label={t('ask.citation', { n, source: label })}
	onclick={onopen}>{n}</button
>
