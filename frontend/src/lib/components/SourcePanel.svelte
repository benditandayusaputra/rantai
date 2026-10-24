<script lang="ts">
	import { ExternalLink, X } from '@lucide/svelte';
	import type { Angka, Provenance } from '$lib/api/tipe';
	import { namaSumber, pemilikSumber } from '$lib/features/provenance/provenance';
	import { formatJam, formatTanggal } from '$lib/format';
	import { bahasaBawaan, t } from '$lib/i18n';
	import CertaintyBadge from './CertaintyBadge.svelte';

	interface Props {
		angka?: Angka;
		nilai?: string;
		provenance: Provenance | null;
		isGagal: boolean;
		open: boolean;
	}

	let { angka, nilai, provenance, isGagal, open = $bindable() }: Props = $props();

	const idJudul = $props.id();
	const level = $derived(angka?.label ?? provenance?.label);
	let dialog: HTMLDialogElement;
	let judul: HTMLHeadingElement;

	$effect(() => {
		if (open && !dialog.open) {
			dialog.showModal();
			judul.focus();
		} else if (!open && dialog.open) {
			dialog.close();
		}
	});

	function waktuAmbil(p: Provenance): string {
		const d = new Date(p.diambil_pada);
		return t('source.retrievedAt', { date: formatTanggal(d, bahasaBawaan), time: formatJam(d) });
	}
</script>

<dialog
	bind:this={dialog}
	onclose={() => (open = false)}
	aria-labelledby={idJudul}
	class="w-full max-w-none rounded-t-xl bg-raised p-5 text-fg shadow-3 backdrop:bg-scrim max-lg:mb-0 max-lg:mt-auto lg:max-w-sm lg:rounded-lg"
>
	<div class="flex items-start justify-between gap-3">
		<h2 id={idJudul} bind:this={judul} tabindex="-1" class="t-h3 pt-2">{t('source.title')}</h2>
		<button
			type="button"
			class="inline-flex min-h-touch min-w-touch items-center justify-center rounded-md hover:bg-sunken"
			aria-label={t('action.close')}
			onclick={() => (open = false)}
		>
			<X size={20} strokeWidth={1.75} aria-hidden="true" />
		</button>
	</div>
	{#if level}
		<div class="mt-3 flex flex-wrap items-center gap-2">
			<CertaintyBadge {level} size="md" />
			<span class="text-small text-muted">{t(`cert.${level}.desc`)}</span>
		</div>
	{/if}
	{#if angka}
		<p class="n-md mt-4">
			{nilai} <span class="text-small font-semibold text-muted">{angka.satuan}</span>
		</p>
	{/if}
	{#if provenance}
		<dl class="mt-4 space-y-3 text-small">
			<div>
				<dt class="t-over">{t('source.title')}</dt>
				<dd>
					{namaSumber(provenance.sumber)}
					<span class="text-muted">{pemilikSumber(provenance.sumber)}</span>
				</dd>
			</div>
			<div>
				<dt class="t-over">{t('source.retrieved')}</dt>
				<dd class="num">{waktuAmbil(provenance)}</dd>
			</div>
			{#if provenance.catatan}
				<div>
					<dt class="t-over">{t('source.formula')}</dt>
					<dd class="mt-1 rounded-md bg-sunken p-3 font-mono text-mono">{provenance.catatan}</dd>
				</div>
			{/if}
		</dl>
		{#if provenance.url}
			<a
				href={provenance.url}
				target="_blank"
				rel="noopener"
				class="mt-4 inline-flex min-h-touch items-center gap-2 font-bold text-secondary-on underline"
			>
				{t('action.openOriginal')}
				<ExternalLink size={16} strokeWidth={1.75} aria-hidden="true" />
				<span class="sr-only">{t('a11y.opensNewTab')}</span>
			</a>
		{/if}
	{:else if isGagal}
		<p role="alert" class="mt-4 text-small text-danger">{t('source.failed')}</p>
	{:else}
		<p aria-busy="true" class="mt-4 text-small text-muted">{t('source.loading')}</p>
	{/if}
</dialog>
