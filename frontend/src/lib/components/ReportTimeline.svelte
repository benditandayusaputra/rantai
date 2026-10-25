<script lang="ts">
	import { Check, Clock } from '@lucide/svelte';
	import type { SaksiLaporan } from '$lib/api/tipe';
	import type { PeristiwaLinimasa } from '$lib/features/lapor/lapor';
	import { tampilanStatus } from '$lib/features/lapor/status';
	import { formatJam, formatTanggal } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';

	interface Props {
		events: PeristiwaLinimasa[];
		witnesses: SaksiLaporan;
	}

	let { events, witnesses }: Props = $props();
</script>

<ol class="flex flex-col">
	{#each events as peristiwa, i (`${peristiwa.status}-${i}`)}
		<li class="relative flex gap-3 pb-5 last:pb-0">
			{#if i < events.length - 1}
				<span
					aria-hidden="true"
					class="absolute left-4 top-9 h-full border-l-link {events[i + 1]?.isTerjadi
						? 'border-primary'
						: 'border-dotted border-line-strong'}"
				></span>
			{/if}
			<span
				class="relative z-raised inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-pill border-link bg-surface {peristiwa.isTerjadi
					? 'border-primary text-primary'
					: 'border-dotted border-line-strong text-muted'}"
			>
				{#if peristiwa.isTerjadi}
					<Check size={16} strokeWidth={2} aria-hidden="true" />
				{:else}
					<Clock size={16} strokeWidth={1.75} aria-hidden="true" />
				{/if}
			</span>
			<span class="flex flex-col gap-0.5 pt-1">
				<span class="font-bold {peristiwa.isTerjadi ? '' : 'text-muted'}"
					>{t(tampilanStatus[peristiwa.status].label, {
						n: witnesses.jumlah,
						m: witnesses.dibutuhkan
					})}</span
				>
				{#if peristiwa.pada}
					<span class="text-small text-muted num"
						>{formatTanggal(peristiwa.pada, bahasaAktif, 'long')} · {formatJam(
							peristiwa.pada
						)}</span
					>
				{:else}
					<span class="text-small text-muted">{t('chain.status.menunggu')}</span>
				{/if}
				{#if peristiwa.catatan}<span class="text-small text-muted">{peristiwa.catatan}</span>{/if}
			</span>
		</li>
	{/each}
</ol>
