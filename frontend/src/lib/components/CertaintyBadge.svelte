<script lang="ts">
	import { BadgeCheck, EqualApproximately, FlaskConical, Sigma } from '@lucide/svelte';
	import type { Label } from '$lib/api/tipe';
	import { t } from '$lib/i18n';
	import { kelasLabel } from './label';

	interface Props {
		level: Label;
		size?: 'sm' | 'md';
		class?: string;
	}

	let { level, size = 'sm', class: kelasTambahan = '' }: Props = $props();

	const ikonLabel = {
		resmi: BadgeCheck,
		proksi: Sigma,
		estimasi: EqualApproximately,
		simulasi: FlaskConical
	} as const;

	const tinggi = { sm: 'h-5.5 px-1.5', md: 'h-7 px-2' } as const;
	const Ikon = $derived(ikonLabel[level]);
</script>

<span
	class="inline-flex shrink-0 items-center gap-1 whitespace-nowrap rounded-sm border-ctrl text-caption font-bold {tinggi[
		size
	]} {kelasLabel[level]} {kelasTambahan}"
	data-label={level}
>
	<Ikon size={16} strokeWidth={1.75} aria-hidden="true" />
	{t(`cert.${level}`)}
</span>
