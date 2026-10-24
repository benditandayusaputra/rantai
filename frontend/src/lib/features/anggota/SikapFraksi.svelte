<script lang="ts">
	import { ExternalLink } from '@lucide/svelte';
	import { ambilSikapFraksi, tangkap } from '$lib/api/klien';
	import type { SikapFraksi } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import { formatTanggal } from '$lib/format';
	import { bahasaBawaan, t } from '$lib/i18n';

	let { fraksi }: { fraksi: string } = $props();

	let daftar = $state<SikapFraksi[] | null>(null);

	$effect(() => {
		void tangkap(ambilSikapFraksi(fetch, fraksi)).then((hasil) => {
			daftar = hasil.ok ? hasil.data : [];
		});
	});
</script>

<section aria-labelledby="judul-sikap" class="flex flex-col gap-4">
	<h2 id="judul-sikap" class="t-h2">{t('stance.title', { fraksi })}</h2>
	<p class="text-small text-muted">{t('stance.lead')}</p>
	{#if daftar === null}
		<p class="text-muted" aria-live="polite">{t('state.loading')}</p>
	{:else if daftar.length === 0}
		<EmptyState title={t('stance.empty')} />
	{:else}
		<ul class="flex flex-col gap-3">
			{#each daftar as sikap (sikap.url)}
				<li class="flex flex-col gap-2 rounded-lg border border-line bg-surface p-4">
					<p class="flex flex-wrap items-center gap-2 text-small">
						<span class="rounded-xs bg-sunken px-2 py-0.5 text-caption font-bold uppercase"
							>{t(`stance.agenda.${sikap.agenda}`)}</span
						>
						<span class="font-semibold">{sikap.judul_agenda}</span>
						<span class="text-muted num"
							>{formatTanggal(new Date(sikap.tanggal), bahasaBawaan, 'long')}</span
						>
					</p>
					<p>{sikap.ringkasan}</p>
					<a
						href={sikap.url}
						target="_blank"
						rel="noopener"
						class="inline-flex items-center gap-1 self-start text-small font-semibold underline"
					>
						{t('stance.source')}<ExternalLink
							size={14}
							strokeWidth={1.75}
							aria-hidden="true"
						/><span class="sr-only">{t('a11y.opensNewTab')}</span>
					</a>
				</li>
			{/each}
		</ul>
	{/if}
</section>
