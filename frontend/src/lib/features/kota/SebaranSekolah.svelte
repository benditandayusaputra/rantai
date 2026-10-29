<script lang="ts">
	import type { PendidikanKota } from '$lib/api/tipe';
	import EmptyState from '$lib/components/EmptyState.svelte';
	import ProvenanceBadge from '$lib/components/ProvenanceBadge.svelte';
	import TableScroll from '$lib/components/TableScroll.svelte';
	import GrafikBatang from '$lib/features/dasbor/GrafikBatang.svelte';
	import { teksJenjang, type TemaPendidikan } from '$lib/features/pendidikan/jenjang';
	import { bahasaAktif, t } from '$lib/i18n';
	import { formatAngka, formatTanggal } from '$lib/format';
	import {
		bentukSekolah,
		jumlahSekolah,
		sumberSekolah,
		susunSebaran,
		type CakupanSekolah
	} from './pendidikan-kota';

	let {
		data,
		tema,
		cakupan = $bindable('negeri')
	}: {
		data: PendidikanKota;
		tema: TemaPendidikan;
		cakupan: CakupanSekolah;
	} = $props();
	const pilihan = ['negeri', 'swasta', 'semua'] as const;
	const sebaran = $derived(susunSebaran(data, tema, cakupan));
	const sumber = $derived(sumberSekolah(data));
	const bilangan = (n: number) => formatAngka(n, bahasaAktif);
</script>

<section aria-labelledby="judul-sebaran-kota" class="flex min-w-0 flex-col gap-5">
	<div class="flex flex-wrap items-end justify-between gap-4">
		<div>
			<h2 id="judul-sebaran-kota" class="t-h2">{t('city.distribution.title')}</h2>
			<p class="mt-2 text-small text-muted">{t('city.distribution.scope')}</p>
		</div>
		<div role="group" aria-label={t('city.schoolStatus')} class="flex flex-wrap gap-2">
			{#each pilihan as status (status)}
				<button
					type="button"
					class="chip"
					aria-pressed={cakupan === status}
					onclick={() => (cakupan = status)}>{t(`city.status.${status}`)}</button
				>
			{/each}
		</div>
	</div>
	{#if data.sekolah.length}
		<div class="card p-4 md:p-6">
			<h3 class="t-h3 mb-5">
				{t('city.distribution.chart', {
					level: t(teksJenjang[tema].jenjang),
					status: t(`city.status.${cakupan}`)
				})}
			</h3>
			<GrafikBatang
				isi={sebaran}
				nilaiDari={(p) => p.jumlah}
				namaDari={(p) => p.nama ?? p.kunci}
				format={bilangan}
			/>
		</div>
		<TableScroll label={t('city.distribution.table')}>
			<table class="w-full text-small">
				<caption class="p-4 text-left font-semibold"
					>{t('city.distribution.table')} · {t(`city.status.${cakupan}`)}</caption
				>
				<thead class="bg-sunken text-caption">
					<tr>
						<th scope="col" class="p-3 text-left">{t('spend.col.area')}</th>
						{#each bentukSekolah as bentuk (bentuk)}<th scope="col" class="p-3 text-right"
								>{bentuk}</th
							>{/each}
					</tr>
				</thead>
				<tbody>
					{#each data.kecamatan as kecamatan (kecamatan.kode)}
						<tr class="border-t border-line">
							<th scope="row" class="p-3 text-left font-semibold"
								><a
									class="link inline-flex min-h-touch items-center"
									href="/{bahasaAktif}/kecamatan/{kecamatan.kode}?bagian=tema&tema={tema}"
									>{kecamatan.nama}</a
								></th
							>
							{#each bentukSekolah as bentuk (bentuk)}<td class="num p-3 text-right"
									>{bilangan(jumlahSekolah(data, kecamatan.kode, [bentuk], cakupan))}</td
								>{/each}
						</tr>
					{/each}
				</tbody>
			</table>
		</TableScroll>
		<details class="rounded-lg border border-line p-4">
			<summary class="min-h-touch cursor-pointer py-2 font-semibold"
				>{t('city.distribution.sources')}</summary
			>
			<p class="my-3 max-w-prose text-small text-muted">{t('city.distribution.method')}</p>
			<ul class="flex flex-col gap-3">
				{#each sumber as s (s.negeri.provenance_id)}
					<li class="flex flex-wrap items-center gap-3 text-small">
						<span>{formatTanggal(new Date(s.diambil_pada), bahasaAktif)}</span><ProvenanceBadge
							angka={s.negeri}
							nilai={bilangan(s.negeri.nilai)}
						/>
					</li>
				{/each}
			</ul>
		</details>
	{:else}
		<EmptyState
			title={t('city.distribution.empty')}
			description={t('city.distribution.emptyBody')}
		/>
	{/if}
</section>
