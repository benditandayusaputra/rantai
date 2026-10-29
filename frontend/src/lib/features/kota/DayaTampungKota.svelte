<script lang="ts">
	import type { PendidikanKota } from '$lib/api/tipe';
	import CertaintyBadge from '$lib/components/CertaintyBadge.svelte';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import TableScroll from '$lib/components/TableScroll.svelte';
	import { teksJenjang, type TemaPendidikan } from '$lib/features/pendidikan/jenjang';
	import { formatAngka } from '$lib/format';
	import { bahasaAktif, t } from '$lib/i18n';

	let { data, tema }: { data: PendidikanKota; tema: TemaPendidikan } = $props();
	const r = $derived(data.ringkasan);
	const teks = $derived(teksJenjang[tema]);
	const bilangan = (n: number) => formatAngka(n, bahasaAktif);
	const persen = (n: number) => `${formatAngka(n * 100, bahasaAktif, 1)}%`;
	const cartes = $derived([
		{
			label: t(teks.dasborLulusan),
			angka: r.kebutuhan,
			nilai: r.kebutuhan ? bilangan(r.kebutuhan.nilai) : ''
		},
		{
			label: t(teks.dasborKursi),
			angka: r.kapasitas,
			nilai: r.kapasitas ? bilangan(r.kapasitas.nilai) : ''
		},
		{
			label: t('city.capacity.coverage'),
			angka: r.rasio,
			nilai: r.rasio ? persen(r.rasio.nilai) : ''
		}
	]);
</script>

<section aria-labelledby="judul-kapasitas-kota" class="flex min-w-0 flex-col gap-5">
	<div>
		<h2 id="judul-kapasitas-kota" class="t-h2">
			{t('city.capacity.title', { level: t(teks.jenjang) })}
		</h2>
		<p class="mt-2 text-small text-muted">
			{t('city.capacity.period', {
				period: data.periode || t('city.noPeriod'),
				n: r.kecamatan_dengan_data,
				total: r.total_kecamatan
			})}
		</p>
	</div>
	{#if r.kebutuhan && r.kapasitas && r.rasio}
		<dl class="grid gap-3 md:grid-cols-3">
			{#each cartes as kartu (kartu.label)}
				<div class="card flex flex-col gap-3 p-4 md:p-5">
					<dt class="text-small font-semibold text-muted">{kartu.label}</dt>
					<dd class="num font-heading text-h1 font-bold">{kartu.nilai}</dd>
					{#if kartu.angka}<dd><CertaintyBadge level={kartu.angka.label} /></dd>{/if}
				</div>
			{/each}
		</dl>
		<p class="text-caption text-muted">{t('city.capacity.method')}</p>
	{:else}
		<EmptyState
			title={t('city.capacity.incomplete')}
			description={t('city.capacity.incompleteBody')}
		/>
	{/if}
	<TableScroll label={t('city.capacity.table')}>
		<table class="w-full text-small">
			<caption class="p-4 text-left font-semibold">{t('city.capacity.table')}</caption>
			<thead class="bg-sunken text-caption">
				<tr
					><th scope="col" class="p-3 text-left">{t('spend.col.area')}</th><th
						scope="col"
						class="p-3 text-right">{t(teks.dasborLulusan)}</th
					><th scope="col" class="p-3 text-right">{t(teks.dasborKursi)}</th><th
						scope="col"
						class="p-3 text-right">{t('city.capacity.coverage')}</th
					><th scope="col" class="p-3 text-left">{t('city.capacity.source')}</th></tr
				>
			</thead>
			<tbody>
				{#each data.kecamatan as k (k.kode)}
					{@const i = k.indikator}
					<tr class="border-t border-line">
						<th scope="row" class="p-3 text-left font-semibold"
							><a
								href="/{bahasaAktif}/kecamatan/{k.kode}?bagian=tema&tema={tema}"
								class="link inline-flex min-h-touch items-center">{k.nama}</a
							></th
						>
						<td class="num p-3 text-right"
							>{i?.rasio ? bilangan(i.kebutuhan.nilai) : t('city.unavailable')}</td
						>
						<td class="num p-3 text-right"
							>{i ? bilangan(i.kapasitas.nilai) : t('city.unavailable')}</td
						>
						<td class="num p-3 text-right" class:text-attn={i?.rasio && i.rasio.nilai < 1}
							>{i?.rasio ? persen(i.rasio.nilai) : t('city.unavailable')}</td
						>
						<td class="p-3"
							>{#if i}<ProvenanceBadge
									angka={i.kapasitas}
									nilai={bilangan(i.kapasitas.nilai)}
									compact
								/>{/if}</td
						>
					</tr>
				{/each}
			</tbody>
		</table>
	</TableScroll>
</section>
