<script lang="ts">
	import { BadgeCheck, EqualApproximately, FlaskConical, Sigma } from '@lucide/svelte';
	import type { Label } from '$lib/api/tipe';
	import { t } from '$lib/i18n';

	interface Props {
		level: Label;
		size?: 'sm' | 'md';
		class?: string;
	}

	let { level, size = 'sm', class: kelasTambahan = '' }: Props = $props();

	const tampilan = {
		resmi: { ikon: BadgeCheck, kelas: 'border-cert-resmi bg-cert-resmi text-cert-resmi-on' },
		proksi: { ikon: Sigma, kelas: 'border-cert-proksi bg-cert-proksi-bg text-cert-proksi' },
		estimasi: {
			ikon: EqualApproximately,
			kelas: 'border-dashed border-cert-estimasi text-cert-estimasi'
		},
		simulasi: { ikon: FlaskConical, kelas: 'hatch-sim border-cert-sim text-cert-sim' }
	} as const;

	const tinggi = { sm: 'h-5.5 px-1.5', md: 'h-7 px-2' } as const;
	const Ikon = $derived(tampilan[level].ikon);
</script>

<span
	class="inline-flex shrink-0 items-center gap-1 whitespace-nowrap rounded-sm border-ctrl text-caption font-bold {tinggi[
		size
	]} {tampilan[level].kelas} {kelasTambahan}"
	data-label={level}
>
	<Ikon size={16} strokeWidth={1.75} aria-hidden="true" />
	{t(`cert.${level}`)}
</span>
