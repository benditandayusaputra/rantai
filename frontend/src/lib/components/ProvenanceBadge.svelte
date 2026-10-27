<script lang="ts">
	import { ChevronDown } from '@lucide/svelte';
	import type { Angka, Provenance } from '$lib/api/tipe';
	import { muatProvenance, namaSumber } from '$lib/features/provenance/provenance';
	import { formatTanggal } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import CertaintyBadge from './CertaintyBadge.svelte';
	import SourcePanel from './SourcePanel.svelte';

	interface Props {
		angka: Angka;
		nilai: string;
		compact?: boolean;
	}

	let { angka, nilai, compact = false }: Props = $props();

	let open = $state(false);
	let provenance = $state<Provenance | null>(null);
	let isGagal = $state(false);

	$effect(() => {
		muatProvenance(angka.provenance_id)
			.then((p) => (provenance = p))
			.catch(() => (isGagal = true));
	});

	const tingkat = $derived(t(`cert.${angka.label}`));
	const tanggal = $derived(provenance ? new Date(provenance.diambil_pada) : null);
	const labelAksesibel = $derived(
		provenance && tanggal
			? t('source.label', {
					name: namaSumber(provenance.sumber),
					level: tingkat,
					date: formatTanggal(tanggal, bahasaAktif, 'long')
				})
			: t('source.labelShort', { level: tingkat })
	);
</script>

<button
	type="button"
	class="inline-flex min-h-touch max-w-full items-center gap-1 rounded-sm border-ctrl border-line-strong bg-surface px-1.5 text-caption text-fg transition-colors duration-fast ease-standard hover:bg-sunken"
	aria-haspopup="dialog"
	aria-expanded={open}
	aria-label={labelAksesibel}
	onclick={() => (open = true)}
>
	<CertaintyBadge level={angka.label} class="shrink-0 border-0" />
	{#if compact}
		<span class="whitespace-nowrap">{t('source.openShort')}</span>
	{:else if provenance && tanggal}
		<span class="truncate px-1.5">
			{namaSumber(provenance.sumber)}{#if !compact}<span class="num text-muted">
					· {formatTanggal(tanggal, bahasaAktif)}</span
				>{/if}
		</span>
	{/if}
	<ChevronDown size={14} strokeWidth={1.75} aria-hidden="true" class="mr-0.5 shrink-0" />
</button>
<SourcePanel bind:open {angka} {nilai} {provenance} {isGagal} />
