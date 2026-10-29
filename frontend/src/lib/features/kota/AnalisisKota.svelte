<script lang="ts">
	import { ArrowRight, Info } from '@lucide/svelte';
	import type { PendidikanKota } from '$lib/api/tipe';
	import { teksJenjang, type TemaPendidikan } from '$lib/features/pendidikan/jenjang';
	import { formatAngka } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';
	import { rentangSebaran, susunSebaran, type CakupanSekolah } from './pendidikan-kota';

	let {
		data,
		tema,
		cakupan
	}: { data: PendidikanKota; tema: TemaPendidikan; cakupan: CakupanSekolah } = $props();
	const rentang = $derived(rentangSebaran(susunSebaran(data, tema, cakupan)));
	const r = $derived(data.ringkasan);
	const bilangan = (n: number) => formatAngka(n, bahasaAktif);
	const daftarNama = (nama: string[]) =>
		new Intl.ListFormat(bahasaAktif, { style: 'long', type: 'conjunction' }).format(nama);
	const level = $derived(t(teksJenjang[tema].jenjang));
	const sumberSPMB =
		'https://pusatinformasi.ult.kemendikdasmen.go.id/hc/id/articles/48105439825945-Pertanyaan-sering-ditanya-FAQ-seputar-Sistem-Penerimaan-Murid-Baru-SPMB';
</script>

<section aria-labelledby="judul-penjelasan-kota" class="card flex flex-col gap-5 p-4 md:p-6">
	<h2 id="judul-penjelasan-kota" class="t-h2">{t('city.explanation.title')}</h2>
	{#if rentang && data.sekolah.length}
		<p>
			{t(
				rentang.minimum === rentang.maksimum ? 'city.explanation.equal' : 'city.explanation.range',
				{
					min: bilangan(rentang.minimum),
					max: bilangan(rentang.maksimum),
					level,
					status: t(`city.status.${cakupan}`)
				}
			)}
		</p>
		{#if rentang.minimum !== rentang.maksimum}
			<dl class="grid gap-4 rounded-md bg-sunken p-4 md:grid-cols-2">
				<div>
					<dt class="text-caption font-semibold text-muted">{t('city.explanation.lowest')}</dt>
					<dd class="mt-1 font-semibold">{daftarNama(rentang.terendah)}</dd>
					<dd class="mt-1 text-small text-muted">
						{t('city.explanation.count', { n: bilangan(rentang.minimum), level })}
					</dd>
				</div>
				<div>
					<dt class="text-caption font-semibold text-muted">{t('city.explanation.highest')}</dt>
					<dd class="mt-1 font-semibold">{daftarNama(rentang.tertinggi)}</dd>
					<dd class="mt-1 text-small text-muted">
						{t('city.explanation.count', { n: bilangan(rentang.maksimum), level })}
					</dd>
				</div>
			</dl>
		{/if}
	{/if}
	<p class="text-small text-muted">{t('city.explanation.fairness')}</p>
	<p class="text-small text-muted">{t('city.explanation.capacity')}</p>
</section>

<section aria-labelledby="judul-dampak-kota" class="rounded-lg bg-info-bg p-4 md:p-6">
	<h2 id="judul-dampak-kota" class="t-h2">{t('city.impact.title')}</h2>
	<p class="mt-4 max-w-prose">{t('city.impact.body')}</p>
	<p class="mt-3 max-w-prose text-small">{t('city.impact.rules')}</p>
	<p class="mt-3 max-w-prose text-small text-muted">{t('city.impact.unknown')}</p>
	<a
		class="link mt-4 inline-flex min-h-touch items-center gap-2 text-small"
		href={sumberSPMB}
		target="_blank"
		rel="noopener"
		>{t('city.impact.source')}<ArrowRight size={16} aria-hidden="true" /><span class="sr-only"
			>{t('a11y.opensNewTab')}</span
		></a
	>
</section>

<section aria-labelledby="judul-kesimpulan-kota" class="card p-4 md:p-6">
	<h2 id="judul-kesimpulan-kota" class="t-h2">{t('city.conclusion.title')}</h2>
	{#if r.kebutuhan && r.kapasitas && r.kesenjangan}
		<p class="mt-4 font-semibold">
			{t(r.kesenjangan.nilai > 0 ? 'city.conclusion.short' : 'city.conclusion.enough', {
				n: bilangan(Math.abs(r.kesenjangan.nilai)),
				level
			})}
		</p>
		<p class="mt-3">
			{t('city.conclusion.local', { n: r.kecamatan_kurang, total: r.total_kecamatan, level })}
		</p>
	{:else}
		<p class="mt-4">
			{t('city.conclusion.incomplete', { n: r.kecamatan_dengan_data, total: r.total_kecamatan })}
		</p>
	{/if}
	<p class="mt-3 max-w-prose text-small text-muted">{t('city.conclusion.action')}</p>
</section>

<section aria-labelledby="judul-batasan-kota" class="flex flex-col gap-3 border-t border-line pt-5">
	<h2 id="judul-batasan-kota" class="t-h3 flex items-center gap-2">
		<Info size={20} aria-hidden="true" />{t('city.limits.title')}
	</h2>
	<ul class="list-disc space-y-2 pl-5 text-small text-muted">
		<li>{t('city.limits.residence')}</li>
		<li>{t('city.limits.private')}</li>
		<li>{t('city.limits.cohort')}</li>
		<li>{t('city.limits.inventory')}</li>
	</ul>
</section>
