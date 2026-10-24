<script lang="ts">
	import { t } from '$lib/i18n';
	import { draf } from './draf.svelte';

	let { onubah }: { onubah: (langkah: number) => void } = $props();

	const namaPaket = $derived(
		draf.paketTerdekat?.find((p) => p.kode === draf.paketKode)?.nama ?? t('report.s2.notListed')
	);
	const baris = $derived([
		{ langkah: 2, label: t('report.s6.project'), isi: namaPaket },
		{
			langkah: 3,
			label: t('report.s6.photos'),
			isi: t('report.s6.photoCount', { n: draf.foto.length })
		},
		{
			langkah: 4,
			label: t('report.s6.condition'),
			isi: draf.kondisi ? t(`report.cond.${draf.kondisi}`) : ''
		},
		{ langkah: 5, label: t('report.s6.note'), isi: draf.catatan.trim() || t('report.s6.noNote') }
	]);
</script>

<div class="flex flex-col gap-4">
	<h1 class="t-h1">{t('report.s6.title')}</h1>
	<dl class="flex flex-col divide-y divide-line rounded-lg border-ctrl border-line">
		{#each baris as b (b.langkah)}
			<div class="flex items-start gap-3 p-4">
				<div class="flex flex-1 flex-col gap-1">
					<dt class="text-small text-muted">{b.label}</dt>
					<dd class="break-words">{b.isi}</dd>
				</div>
				<button
					type="button"
					class="inline-flex min-h-touch items-center rounded-md px-2 font-semibold text-secondary-on underline-offset-4 hover:underline"
					aria-label="{t('action.edit')} {b.label}"
					onclick={() => onubah(b.langkah)}>{t('action.edit')}</button
				>
			</div>
		{/each}
		<div class="flex flex-col gap-1 p-4">
			<dt class="text-small text-muted">{t('report.s6.location')}</dt>
			<dd>{t('report.s6.locationValue')}</dd>
		</div>
	</dl>
	<label class="flex min-h-touch cursor-pointer items-start gap-3">
		<input type="checkbox" bind:checked={draf.setuju} class="mt-1 h-5.5 w-5.5 accent-primary" />
		<span>{t('report.s6.consent')}</span>
	</label>
</div>
